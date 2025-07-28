package main

import (
    "os"
    "fmt"
    "bufio"
    "math/rand"
    
)

const VALID_WORDS_PATH = "words/all_valid_words.txt"
var LETTER_CONDITIONS = make(map[string]map[string]any)

func init() {
    for i := 'a'; i <= 'z'; i++ {
        letter := string(i)
        LETTER_CONDITIONS[letter] = map[string]any{
            "status":           nil,
            "correct_position": nil,
            "wrong_positions":  nil,
        }
    }
}



func get_random_word(word_list []string) string {
    return word_list[rand.Intn(len(word_list))]
}



func nyt_word_validator(final_word string, guessed_word string) string {
    nyt_string := []string{}
    
}



func get_valide_words() []string {
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



func get_letter_frequency(word_list []string) map[string]float64 {
    letter_count := make(map[string]int)
    for _ , word := range word_list {
        for _ , letter := range word {
            letter_count[string(letter)]++
        }
    }

    total_letter := len(letter_count) * 5
    letter_frequency := make(map[string]float64)
    for letter, count := range letter_count {
        letter_frequency[letter] = float64(count) / float64(total_letter)
    }
    
    return letter_frequency
}



func main() {
    word_list := get_valide_words()
    letter_frequency := get_letter_frequency(word_list)
    final_word := get_random_word(word_list)

    fmt.Println(len(word_list))
    fmt.Println(len(letter_frequency))
    fmt.Println(final_word)
}