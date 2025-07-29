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



func contains(letter string, word string) bool {
    for _, l := range word {
        if string(l) == letter {
            return true
        }
    }
    return false
}



func nyt_word_validator(final_word string, guessed_word string) string {
    nyt_string := []string{}
    for i := 0; i < len(guessed_word); i++ {
        if guessed_word[i] == final_word[i] {
            nyt_string = append(nyt_string, "g")
        } else if contains(string(guessed_word[i]), final_word) {
            nyt_string = append(nyt_string, "y")
        } else {
            nyt_string = append(nyt_string, "b")
        }
    }
    return strings.Join(nyt_string, "")
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



func get_letter_frequency(word_list []string) map[string]map[int]float64 {
    letter_count := make(map[string]map[int]int)
    for _, word := range word_list {
        for i := 0; i < 5; i++ {
            single := string(word[i])
            if letter_count[single] == nil {
                letter_count[single] = make(map[int]int)
            }
            letter_count[single][i]++
        }

        for i := 0; i < 4; i++ {
            double := string(word[i:i+2])
            if letter_count[double] == nil {
                letter_count[double] = make(map[int]int)
            }
            letter_count[double][i]++
        }

        for i := 0; i < 3; i++ {
            triple := string(word[i:i+3])
            if letter_count[triple] == nil {
                letter_count[triple] = make(map[int]int)
            }
            letter_count[triple][i]++
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
    word_probabilities := map[string]float64{}
    for _, word := range word_list {
        var prob float64 = 0
        for i := 0; i < 5; i++ {
            letter := string(word[i])
            prob += letter_frequency[letter][i]
        }
        for i := 0; i < 4; i++ {
            letter := string(word[i:i+2])
            prob += letter_frequency[letter][i]
        }
        for i := 0; i < 3; i++ {
            letter := string(word[i:i+3])
            prob += letter_frequency[letter][i]
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



func main() {
    word_list := get_valide_words()
    letter_frequency := get_letter_frequency(word_list)
    final_word := get_random_word(word_list)
    guessed_word := get_best_word(word_list, letter_frequency)
    validation_string := nyt_word_validator(final_word, guessed_word)

    fmt.Println(len(word_list))
    fmt.Println(len(letter_frequency))
    fmt.Println(final_word)
    fmt.Printf("The best word is: %s\n", guessed_word)
    fmt.Println(validation_string)
}