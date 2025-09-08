package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Build the path relative to this file's location
// Since main.go is in src/go, we go up 2 directories to reach project root
var VALID_WORDS_PATH = filepath.Join(get_current_dir(), "..", "..", "words", "all_valid_words.txt")

var LETTER_CONDITIONS = make(map[string]map[string]any)
var FINAL_WORD = make([]string, 5)
var KNOWN_LETTERS = []string{}

type word_score struct {
	word string
	prob float64
}

func init() {
	// Single letter checks
	for i := 'a'; i <= 'z'; i++ {
		letter := string(i)
		LETTER_CONDITIONS[letter] = map[string]any{
			"status":           nil,
			"correct_position": []int{},
			"wrong_positions":  []int{},
			"double":           true,
		}
	}
}



func reset_game_state() {
	// Reset letter conditions
	for i := 'a'; i <= 'z'; i++ {
		letter := string(i)
		LETTER_CONDITIONS[letter] = map[string]any{
			"status":           nil,
			"correct_position": []int{},
			"wrong_positions":  []int{},
			"double":           true,
		}
	}
	// Reset final word
	for i := range FINAL_WORD {
		FINAL_WORD[i] = ""
	}

	KNOWN_LETTERS = []string{}
}



func get_known_letters() []string {
	return KNOWN_LETTERS
}



func get_final_word() []string {
	return FINAL_WORD
}



func get_random_word(word_list []string) string {
	return word_list[rand.Intn(len(word_list))]
}



func nyt_word_validator(final_word string, guessed_word string) string {
	nyt_string := []string{}
	for i := 0; i < len(guessed_word); i++ {
		// If the letter is in the correct position
		if guessed_word[i] == final_word[i] {
			nyt_string = append(nyt_string, "g")
			// If the letter is in the word, but incorrect position
		} else if contains_letter(string(guessed_word[i]), final_word) {
			nyt_string = append(nyt_string, "y")
			// Letter is not in the word
		} else {
			nyt_string = append(nyt_string, "b")
		}
	}
	return strings.Join(nyt_string, "")
}



func update_letter_conditions(validation_string string, guessed_word string) {
	seen := []string{}
	for i := 0; i < len(guessed_word); i++ {
		letter_str := string(guessed_word[i])

		switch validation_string[i] {
		case 'g':
			// If the letter is in the correct position
			LETTER_CONDITIONS[letter_str]["status"] = true

			positions := LETTER_CONDITIONS[letter_str]["correct_position"].([]int)
			LETTER_CONDITIONS[letter_str]["correct_position"] = append(positions, i)

			// We will add this to the 'mock' final word to help us remove unwated words
			FINAL_WORD[i] = letter_str
			// Since we know that this letter is in the word, we should also remove all other words that don't include it
			KNOWN_LETTERS = append(KNOWN_LETTERS, letter_str)

			// If we have already seen the letter and it shouldn't be in the word
			// That means that duplicates are no longer allowed
			if contains_letter(letter_str, strings.Join(seen, "")) && !LETTER_CONDITIONS[letter_str]["status"].(bool) {
				LETTER_CONDITIONS[letter_str]["double"] = false
			}

		case 'y':
			// If the letter is in the word, but incorrect position
			LETTER_CONDITIONS[letter_str]["status"] = true

			positions := LETTER_CONDITIONS[letter_str]["wrong_positions"].([]int)
			LETTER_CONDITIONS[letter_str]["wrong_positions"] = append(positions, i)

			// If we have already seen the letter and it shouldn't be in the word
			// That means that duplicates are no longer allowed
			if contains_letter(letter_str, strings.Join(seen, "")) && !LETTER_CONDITIONS[letter_str]["status"].(bool) {
				LETTER_CONDITIONS[letter_str]["double"] = false
			}

			// Since we know that this letter is in the word, we should also remove all other words that don't include it
			KNOWN_LETTERS = append(KNOWN_LETTERS, letter_str)

		default:
			// If we have already looked at the letter and the new is a 'b' that means there are no doubles,
			// but could still be possible for that letter to be in the word
			if contains_letter(letter_str, strings.Join(seen, "")) {
				LETTER_CONDITIONS[letter_str]["double"] = false

				// Letter is not in the word
			} else {
				LETTER_CONDITIONS[letter_str]["status"] = false
			}
		}

		seen = append(seen, letter_str)
	}
}



func get_valide_words() []string {
	file, err := os.Open(VALID_WORDS_PATH)
	if err != nil {
		fmt.Printf("Error: Cannot find words file at '%s'\n", VALID_WORDS_PATH)
		fmt.Println("Please ensure you're running from the project root or src/go directory")
		fmt.Printf("Attempted paths: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	word_list := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word_list = append(word_list, scanner.Text())
	}

	return word_list
}



func get_letter_frequency(word_list []string) map[string]map[int]float64 {
	letter_count := make(map[string]map[int]float64)
	for _, word := range word_list {
		for i := 0; i < 5; i++ {
			single := string(word[i])
			if letter_count[single] == nil {
				letter_count[single] = make(map[int]float64)
			}
			letter_count[single][i]++
		}
	}

	total_letter := len(letter_count)
	letter_frequency := make(map[string]map[int]float64)
	for abr, positions := range letter_count {
		letter_frequency[abr] = make(map[int]float64)
		for position, count := range positions {
			letter_frequency[abr][position] = float64(count) / float64(total_letter)
		}
	}

	return letter_frequency
}



func get_best_word(word_list []string, letter_frequency map[string]map[int]float64) string {
	ranked_words := ranked_words(word_list, letter_frequency)
	return ranked_words[0].word
}



func ranked_words(word_list []string, letter_frequency map[string]map[int]float64) []word_score {
	ranked_words := make([]word_score, 0, len(letter_frequency))
	for _, word := range word_list {

		var prob float64 = 0.0

		for i, l := range word {
			letter_str := string(l)
			count := count_letters(letter_str, word)
			prob += letter_frequency[letter_str][i] / float64(count)
		}

		ranked_words = append(ranked_words, word_score{word: word, prob: prob})
	}

	sort.Slice(ranked_words, func(i, j int) bool {
		return ranked_words[i].prob > ranked_words[j].prob
	})

	return ranked_words
}



func filter_word_list(word_list []string) []string {
	filtered_words := []string{}

	for _, word := range word_list {
		// when the FINAL_WORD is being built we should skip any word which does not follow the skeleton
		skip := false

		for i := 0; i < 5; i++ {
			if FINAL_WORD[i] == "" {
				continue
			}

			if FINAL_WORD[i] != string(word[i]) {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		// The word does not contain one of the know letters in the word so we should skip
		if len(KNOWN_LETTERS) > 0 && !containes_letters(word, KNOWN_LETTERS) {
			continue
		}

		add := true

		for pos, letter := range word {
			letter_str := string(letter)

			// We have no information on this letter so we should skip
			if LETTER_CONDITIONS[letter_str]["status"] == nil {
				continue
			}

			if !LETTER_CONDITIONS[letter_str]["status"].(bool) {
				add = false
				break
			}

			wrong_positions := LETTER_CONDITIONS[letter_str]["wrong_positions"].([]int)
			// When the current position is in the wrong_positions we should skip
			if len(wrong_positions) > 0 && contains_number(pos, wrong_positions) {
				add = false
				break
			}
		}

		if add {
			filtered_words = append(filtered_words, word)
		}
	}

	return filtered_words
}
