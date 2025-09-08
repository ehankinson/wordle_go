//go:build generate
// +build generate

package main

import (
	"os"
	"fmt"
	"math"
	"strings"
	"math/rand"
	"encoding/json"
	"path/filepath"
)

var current_dir = get_current_dir()
var NUM_GAMES = 50000
var WORDLE_ATTEMPTS = 6
var STARTING_WORDS = get_valide_words()
var DATA_FILE = filepath.Join(current_dir, "..", "..", "data", "trainning_set.json")

type Strategy int

const (
	EXPLOIT Strategy = iota
	EXPLORE
	RANDOM
)

type ExplorationData struct {
	Temp     float64
	TopK     int
	Strategy Strategy
}

type Experience struct {
	State    State   `json:"state"`
	Action   string  `json:"action"`
	Feedback string  `json:"feedback"`
	Reward   float64 `json:"reward"`
	Done     bool    `json:"done"`
}

type State struct {
	AttemptNum         int               `json:"attempt_number"`
	LetterContraints   LetterConstraints `json:"letter_constraints"`
	ValidWordsRemaning int               `json:"valid_words_remaining"`
	PreviousGuesses    []string          `json:"previous_guesses"`
	PreviousFeedback   []string          `json:"previous_feedback"`
}

type LetterConstraints struct {
	CorrectPositions []string `json:"correct_positions"`
	PresentLetters   []string `json:"present_letters"`
	AbsentLetters    []string `json:"absent_letters"`
}

func main() {
	err := generate_data(NUM_GAMES)
	if err != nil {
		fmt.Printf("Error generating data: %v\n", err)
		os.Exit(1)
	}
}



func get_explration_params() ExplorationData {
	r := rand.Float64()

	if r < 0.4 {
		return ExplorationData{
			Temp:     0.3 + rand.Float64()*0.3,
			TopK:     5 + rand.Intn(5),
			Strategy: EXPLOIT,
		}
	} else if r < 0.8 {
		return ExplorationData{
			Temp:     0.7 + rand.Float64()*0.5,
			TopK:     10 + rand.Intn(10),
			Strategy: EXPLORE,
		}
	} else {
		return ExplorationData{
			Temp:     1.2 + rand.Float64()*0.8,
			TopK:     30 + rand.Intn(70),
			Strategy: RANDOM,
		}
	}
}



func get_word(ranked_words []word_score, params ExplorationData) string {
	// Take only the Kth amount of words to select from
	gussable_words := ranked_words
	if len(ranked_words) > params.TopK {
		gussable_words = ranked_words[:params.TopK]
	}

	// Apply temperature-based selection
	weights := make([]float64, len(gussable_words))

	// Apply the 'softmax' function to allow better scaling on word probabilities
	// Smaller the 'temp' is the bigger the difference is from the best to worst word, and larger is the opposite
	sum := 0.0 // Sum of all weight so we can get a % scale from 0 to the total exp smoothing, instead of using a 0-1 scale
	for i, ws := range gussable_words {
		// Convert score to weight using temperature
		weight := math.Exp(ws.prob / params.Temp)
		weights[i] = weight
		sum += weight
	}

	r := rand.Float64() * sum // Picking a random number between 0-sum [0-1]
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			return gussable_words[i].word
		}
	}

	return gussable_words[0].word
}



func calculate_reward(guess_num int, words_before int, words_after int, feedback string, is_final bool) float64 {
	reward := -1.0

	// Information reduction bonus (0 to 3 points)
	if words_before > 0 {
		reduction_rate := float64(words_before-words_after) / float64(words_before)
		reward += reduction_rate * 3.0
	}

	// Count new information
	greens := strings.Count(feedback, "g")
	yellows := strings.Count(feedback, "y")

	// Small bonuses for constraints discovered
	reward += float64(greens) * 0.3
	reward += float64(yellows) * 0.1

	// Terminal rewards
	if is_final {
		if feedback == "ggggg" {
			reward += 10.0 + float64(7-guess_num)*2.0
		} else {
			reward -= 5.0
		}
	}

	return reward
}



func capture_state(attempt_num int, word_list []string, letter_conditions map[string]map[string]any, previous_guesses []string, previous_feedback []string) State {
	absent_letters := []string{}
	for letter, condition := range letter_conditions {
		// We have no information about this letter
		if condition["status"] == nil {
			continue
		}

		if !condition["status"].(bool) {
			absent_letters = append(absent_letters, letter)
		}
	}

	letter_constraints := LetterConstraints{
		CorrectPositions: get_final_word(),
		PresentLetters:   get_known_letters(),
		AbsentLetters:    absent_letters,
	}
	return State{
		AttemptNum:         attempt_num,
		LetterContraints:   letter_constraints,
		ValidWordsRemaning: len(word_list),
		PreviousGuesses:    previous_guesses,
		PreviousFeedback:   previous_feedback,
	}
}



func generate_data(num_games int) error {
	file, err := os.Create(DATA_FILE)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	file.WriteString("[\n")

	first := true
	total_experiences := 0

	for game := 0; game < num_games; game++ {
		reset_game_state()

		word_list := STARTING_WORDS
		words_before := len(word_list) // Number of words for reward calculation
		exploration_params := get_explration_params()
		final_word := get_random_word(word_list)
		previous_guesses := []string{}
		previous_feedback := []string{}

		for att := 0; att < WORDLE_ATTEMPTS; att++ {
			current_state := capture_state(att, word_list, LETTER_CONDITIONS, previous_guesses, previous_feedback)

			letter_frequency := get_letter_frequency(word_list)
			ranked_words := ranked_words(word_list, letter_frequency)
			guessing_word := get_word(ranked_words, exploration_params)
			previous_guesses = append(previous_guesses, guessing_word)
			validation_string := nyt_word_validator(final_word, guessing_word)
			previous_feedback = append(previous_feedback, validation_string)

			update_letter_conditions(validation_string, guessing_word)
			word_list = filter_word_list(word_list)
			words_after := len(word_list)

			is_winner := validation_string == "ggggg"
			is_final := is_winner || att == WORDLE_ATTEMPTS-1

			reward := calculate_reward(
				att+1,
				words_before,
				words_after,
				validation_string,
				is_final,
			)

			exp := Experience{
				State:    current_state,
				Action:   guessing_word,
				Feedback: validation_string,
				Reward:   reward,
				Done:     is_final,
			}

			if !first {
				file.WriteString(",\n")
			}
			first = false

			json_bytes, err := json.Marshal(exp)
			if err != nil {
				return fmt.Errorf("failed to marshal experience: %v", err)
			}

			file.WriteString("  ")
			file.Write(json_bytes)

			total_experiences++
			words_before = words_after

			if is_winner {
				break
			}
		}

		if (game+1)%100 == 0 {
			fmt.Printf("Generated %d experiences\n", total_experiences)
			file.Sync()
		}
	}

	file.WriteString("\n]")
	fmt.Printf("Generated %d experiences\n", total_experiences)
	return nil
}
