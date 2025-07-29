package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

const VALID_WORDS_PATH = "words/all_valid_words.txt"
var LETTER_CONDITIONS = make(map[string]map[string]any)
var FINAL_WORD = make([]string, 5)

func init() {
	// Single letter checks
	for i := 'a'; i <= 'z'; i++ {
		letter := string(i)
		LETTER_CONDITIONS[letter] = map[string]any{
			"status":           nil,
			"correct_position": nil,
			"wrong_positions":  nil,
            "double": nil,
		}
	}
}



func reset_game_state() {
	// Reset letter conditions
	for i := 'a'; i <= 'z'; i++ {
		letter := string(i)
		LETTER_CONDITIONS[letter] = map[string]any{
			"status":           nil,
			"correct_position": nil,
			"wrong_positions":  nil,
            "double": nil,
		}
	}
	// Reset final word
	for i := range FINAL_WORD {
		FINAL_WORD[i] = ""
	}
}



func get_random_word(word_list []string) string {
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



func Update_letter_conditions(validation_string string, guessed_word string) {
    seen := []string{}
	for i := 0; i < len(guessed_word); i++ {
		letter_str := string(guessed_word[i])

		switch validation_string[i] {
            case 'g':
                // If the letter is in the correct position
                LETTER_CONDITIONS[letter_str]["status"] = true
                LETTER_CONDITIONS[letter_str]["correct_position"] = i
                
                // We will add this to the 'mock' final word to help us remove unwated words
                FINAL_WORD[i] = letter_str

                // If we have already seen the letter and it shouldn't be in the word
                // That means that duplicates are no longer allowed
                if contains_letter(letter_str, strings.Join(seen, "")) && !LETTER_CONDITIONS[letter_str]["status"].(bool) {
                    LETTER_CONDITIONS[letter_str]["double"] = false
                }

            case 'y':
                // If the letter is in the word, but incorrect position
                LETTER_CONDITIONS[letter_str]["status"] = true
                if LETTER_CONDITIONS[letter_str]["wrong_positions"] == nil {
                    LETTER_CONDITIONS[letter_str]["wrong_positions"] = make([]int, 0)
                }
                positions := LETTER_CONDITIONS[letter_str]["wrong_positions"].([]int)
                LETTER_CONDITIONS[letter_str]["wrong_positions"] = append(positions, i)

                // If we have already seen the letter and it shouldn't be in the word
                // That means that duplicates are no longer allowed
                if contains_letter(letter_str, strings.Join(seen, "")) && !LETTER_CONDITIONS[letter_str]["status"].(bool) {
                    LETTER_CONDITIONS[letter_str]["double"] = false
                }

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



func Get_letter_frequency(word_list []string) map[string]map[int]float64 {
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



func Get_best_word(word_list []string, letter_frequency map[string]map[int]float64) string {
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
        seen := []string{}
		add := true

		for pos, letter := range word {
			letter_str := string(letter)

            if FINAL_WORD[pos] != "" && FINAL_WORD[pos] != letter_str {
                add = false
                break
            }
			// If the letter is not in the word skip and don't add
			if LETTER_CONDITIONS[letter_str]["status"] != nil && !LETTER_CONDITIONS[letter_str]["status"].(bool) {
				add = false
				break
			}
			// IF the letter is in the word, and we have the correct position, but its in the incorrect position skip and don't add
			if LETTER_CONDITIONS[letter_str]["correct_position"] != nil && LETTER_CONDITIONS[letter_str]["correct_position"] != pos {
				add = false
				break
			}
			// If the letter is in the word, but we have the wrong position, and its the same we also skip and don't add
			if LETTER_CONDITIONS[letter_str]["wrong_positions"] == nil {
                seen = append(seen, letter_str)
				continue
			}

			wrong_positions := LETTER_CONDITIONS[letter_str]["wrong_positions"].([]int)
			if contains_number(pos, wrong_positions) {
				add = false
				break
            }

            // If we have already seen the letter, and doubles are false we need to skip the word
            if contains_letter(letter_str, strings.Join(seen, "")) {
                add = false
                break
            }

            seen = append(seen, letter_str)
		}

		if add {
			filtered_words = append(filtered_words, word)
		}
	}

	return filtered_words
}
