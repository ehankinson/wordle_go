package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"math/rand"
)

const VALID_WORDS_PATH = "words/all_valid_words.txt"

var LETTER_CONDITIONS = make(map[string]map[string]any)
var FINAL_WORD = make([]string, 5)
var KNOWN_LETTERS = []string{}

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



func Reset_game_state() {
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
	// Reset known letters
	KNOWN_LETTERS = []string{}
}



func Get_random_word(word_list []string) string {
	return word_list[rand.Intn(len(word_list))]
}



func contains_letter(letter string, word string) bool {
	for _, l := range word {
		if string(l) == letter {
			return true
		}
	}
	return false
}



func count_letters(letter string, letters string) int {
	count := 0
	for _, l := range letters {
		if string(l) == letter {
			count++
		}
	}

	return count
}



func contains_number(number int, numbers []int) bool {
	for _, item := range numbers {
		if number == item {
			return true
		}
	}

	return false
}



func containes_letters(word string, known_letters []string) bool {
	for _, know_letter := range known_letters {

		in_word := false
		for _, word_letter := range word {
			if know_letter == string(word_letter) {
				in_word = true
			}
		}

		if !in_word {
			return false
		}
	}

	return true
}



func Nyt_word_validator(final_word string, guessed_word string) string {
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



func Update_letter_conditions(validation_string string, guessed_word string) {
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



func Get_valide_words() []string {
	file, err := os.Open(VALID_WORDS_PATH)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return []string{}
	}

	word_list := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word_list = append(word_list, scanner.Text())
	}

	return word_list
}


func make_letter_dict(word_list []string) map[string]map[int]float64 {
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

	return letter_count
}



func get_letter_frequency(word_list []string) map[string]map[int]float64 {
	letter_count := make_letter_dict(word_list)

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



func get_remaning_letters(word_list []string) map[string]bool {
	letters := make(map[string]bool)
	for _, word := range word_list {
		for _, letter := range word {
			letter_str := string(letter)
			if LETTER_CONDITIONS[letter_str]["status"] == nil && !letters[letter_str]{
				letters[letter_str] = true
			}
		}
	}
	
	return letters
}



func best_common_letters_word(letters map[string]bool) string {
	all_words := Get_valide_words()
	max_count := 0
	best_word := ""

	keys := make([]string, 0, len(letters))
	for letter := range letters {
		keys = append(keys, letter)
	}
	keys_str := strings.Join(keys, "")

	for _, word := range all_words {
		count := 0
		for _, letter := range word {
			letter_str := string(letter)

			if contains_letter(letter_str, keys_str) {
				count++
			}
		}

		if count > max_count {
			max_count = count
			best_word = word
		}
	}

	return best_word
}



func letter_filter(word_list []string) map[string]map[int]float64 {
	letter_count := make_letter_dict(word_list)
	letter_frequency := make(map[string]map[int]float64)

	for letter, positions := range letter_count {
		letter_frequency[letter] = make(map[int]float64)

		// If the letter is in the known letters, or the letter is not in the word, we need to skip them,
		// since we want the best word with the most common letters that have not been used to filter out more letters
		known_letter := contains_letter(letter, strings.Join(KNOWN_LETTERS, ""))
		var good_letter bool = false

		if LETTER_CONDITIONS[letter]["status"] == nil {
			good_letter = true
		}

		for position, count := range positions {
			if known_letter || !good_letter {
				letter_frequency[letter][position] = 0.0
				continue
			}
			letter_frequency[letter][position] = float64(count) / float64(len(word_list))
		}
	}

	return letter_frequency
}



func word_probabilities(letter_frequency map[string]map[int]float64, word_list []string) map[string]float64 {
	word_probabilities := map[string]float64{}

	for _, word := range word_list {

		var prob float64 = 0.0
		for i, l := range word {
			letter_str := string(l)
			count := count_letters(letter_str, word)
			prob += letter_frequency[letter_str][i] / float64(count)
		}

		word_probabilities[word] = prob
	}

	return word_probabilities
}



func best_filter_word_probabilities() map[string]float64 {
	all_words := Get_valide_words()
	letter_frequency := letter_filter(all_words)
	word_probabilities := word_probabilities(letter_frequency, all_words)
	return word_probabilities
}



func best_word_probabilities(word_list []string) map[string]float64 {
	letter_frequency := get_letter_frequency(word_list)
	word_probabilities := word_probabilities(letter_frequency, word_list)
	return word_probabilities
}



func Get_best_word(word_list []string) string {
	// If the word list is less than 1000 words and known letters is less then 4,
	// we should use letters that are unknow to filter out letter
	var word_probabilities map[string]float64

	if len(word_list) < 100 {
		letters := get_remaning_letters(word_list)
		best_word := best_common_letters_word(letters)
		return best_word
	} else if len(word_list) < 1000 && len(KNOWN_LETTERS) < 4 {
		word_probabilities = best_filter_word_probabilities()
	} else {
		word_probabilities = best_word_probabilities(word_list)
	}

	best_word := ""
	max_prob := 0.0
	for word, probability := range word_probabilities {
		if probability > max_prob {
			max_prob = probability
			best_word = word
		}
	}

	return best_word
}



func Filter_word_list(word_list []string) []string {
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
