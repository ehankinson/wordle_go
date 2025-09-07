package main

import (
	"math"
	"math/rand"
	"strings"
)

var WORDLE_ATTEMPTS = 6
var STARTING_WORDS = Get_valide_words()

type Strategy int

const (
	EXPLOIT Strategy = iota
	EXPLORE
	RANDOM
)

type exploration_data struct {
	temp     float64
	top_k    int
	strategy Strategy
}

type experience struct {
	state state
	action string
	feedback string
	reward float64
	done bool
}

type state struct {
	attempt_number int
	letter_constraints letter_constraints
	valid_words_remaining int
	previous_guesses []string
	previous_feedback []string
}

type letter_constraints struct {
	correct_positions map[int]string
	present_letters []string
	absent_letters []string
}



func get_explration_params() exploration_data {
	r := rand.Float64()

	if r < 0.4 {
		return exploration_data {
			temp: 0.3 + rand.Float64() * 0.3,
			top_k: 5 + rand.Intn(5),
			strategy: EXPLOIT,
		}
	} else if r < 0.8 {
		return exploration_data {
			temp: 0.7 + rand.Float64() * 0.5,
			top_k: 10 + rand.Intn(10),
			strategy: EXPLORE,
		}
	} else {
		return exploration_data {
			temp: 1.2 + rand.Float64() * 0.8,
			top_k: 30 + rand.Intn(70),
			strategy: RANDOM,
		}
	}
}



func get_word(ranked_words []word_score, params exploration_data) string {
	// Take only the Kth amount of words to select from
	gussable_words := ranked_words
	if len(ranked_words) > params.top_k {
		gussable_words = ranked_words[:params.top_k]
	}

	// Apply temperature-based selection
	weights := make([]float64, len(gussable_words))
	
	// Apply the 'softmax' function to allow better scaling on word probabilities
	// Smaller the 'temp' is the bigger the difference is from the best to worst word, and larger is the opposite
	sum := 0.0 // Sum of all weight so we can get a % scale from 0 to the total exp smoothing, instead of using a 0-1 scale
	for i, ws := range gussable_words {
		// Convert score to weight using temperature
		weight := math.Exp(ws.prob / params.temp) 
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



func capture_state(attempt_num int, word_list []string, letter_conditions map[string]map[string]any, previous_guesses []string, previous_feedback []string) state {
	
}



func generate_data(num_games int) {
	for game := 0; game < num_games; game++ {
		word_list := STARTING_WORDS
		words_before := len(word_list) // Number of words for reward calculation
		exploration_params := get_explration_params()
		final_word := Get_random_word(word_list)

		game_experiences := []experience{}

		for att := 0; att < WORDLE_ATTEMPTS; att++ {
			current_state := capture_state(att, word_list, LETTER_CONDITIONS)
			
			letter_frequency := Get_letter_frequency(word_list)
			ranked_words := Ranked_words(word_list, letter_frequency)
			guessing_word := get_word(ranked_words, exploration_params)
			validation_string := Nyt_word_validator(final_word, guessing_word)

			Update_letter_conditions(validation_string, guessing_word)
			word_list = Filter_word_list(word_list)
			words_after := len(word_list)

			is_winner := validation_string == "ggggg"
			is_final := is_winner || att == WORDLE_ATTEMPTS - 1

			reward := calculate_reward(
				att + 1,
				words_before,
				words_after,
				validation_string,
				is_final,
			)

			exp := experience {
				state: current_state,
				action: guessing_word,
				feedback: validation_string,
				reward: reward,
				done: is_final,
			}

			game_experiences = append(game_experiences, exp)

			words_before = words_after

			if is_winner {
				break
			}
		}
	}
}
