package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

func remove_word(word_list []string, word string) []string {
	result := make([]string, 0, len(word_list))
	for _, w := range word_list {
		if w != word {
			result = append(result, w)
		}
	}

	return result
}

func play_game(reader *bufio.Reader, validation_string string) string {
	word_list := Get_valide_words()
	letter_frequency := Get_letter_frequency(word_list)
	best_word := ""
	fmt.Println("If the word given does not exists, please right skip and a new work will be provided")
	for i := 0; i < 6; i++ {
		fmt.Printf("There was a total of %d words to choose from\n", len(word_list))
		validation_string := "skip"
		for {
			if validation_string != "skip" { break }

			if validation_string == "solved" {
				return validation_string
			}

			best_word = Get_best_word(word_list, letter_frequency)
			// best_word = "cooly"
			fmt.Printf("The Best Word Is: %s\n", best_word)
			
			fmt.Print("Enter the nyt validation please: ")
			validation_string, _ = reader.ReadString('\n')
			validation_string = strings.TrimSpace(validation_string)

			if validation_string == "skip" {
				word_list = remove_word(word_list, best_word)
				letter_frequency = Get_letter_frequency(word_list)
			}
		}

		Update_letter_conditions(validation_string, best_word)
		word_list = Filter_word_list(word_list)
		letter_frequency = Get_letter_frequency(word_list)
		fmt.Println("")
	}

	return ""
}


func Solve_nyt() {
	validation_string := "skip"
	reader := bufio.NewReader(os.Stdin)
	for {
		if validation_string == "solved" {
			fmt.Print("Would you like to continue (Y/n)")
			cont, _ := reader.ReadString('\n')
			cont = strings.TrimSpace(cont)

			if cont == "n" { 
				break 
			} else {
				validation_string = play_game(reader, validation_string)
			} 
		} else {
			validation_string = play_game(reader, validation_string)
		}
	}
}



func main() {
    Solve_nyt()
}