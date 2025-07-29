package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBlack  = "\033[30m"
	
	// Background colors for better visibility
	BgGreen  = "\033[42m"
	BgYellow = "\033[43m"
	BgRed    = "\033[41m"
	BgBlack  = "\033[40m"
)

func colorize_letter(letter byte, validation byte) string {
	letterStr := string(letter)
	switch validation {
	case 'g':
		return BgGreen + ColorWhite + " " + letterStr + " " + ColorReset
	case 'y':
		return BgYellow + ColorBlack + " " + letterStr + " " + ColorReset
	case 'b':
		return BgBlack + ColorWhite + " " + letterStr + " " + ColorReset
	default:
		return " " + letterStr + " "
	}
}

func display_colored_word(word string, validation string) string {
	result := ""
	for i := 0; i < len(word) && i < len(validation); i++ {
		result += colorize_letter(word[i], validation[i])
	}
	return result
}

func remove_word(word_list []string, word string) []string {
	result := make([]string, 0, len(word_list))
	for _, w := range word_list {
		if w != word {
			result = append(result, w)
		}
	}
	return result
}



func play_single_game(reader *bufio.Reader) {
	reset_game_state()
	word_list := Get_valide_words()
	letter_frequency := Get_letter_frequency(word_list)

	fmt.Println(ColorCyan + "\n=== Starting New Wordle Game ===" + ColorReset)
	fmt.Println("Enter validation string: " + BgGreen + ColorWhite + " g " + ColorReset + " = green (correct), " + BgYellow + ColorBlack + " y " + ColorReset + " = yellow (wrong position), " + BgBlack + ColorWhite + " b " + ColorReset + " = black (not in word)")
	fmt.Println(ColorYellow + "Type 'skip' if the suggested word doesn't exist" + ColorReset)
	fmt.Println(ColorGreen + "Type 'ggggg' when you solve the puzzle" + ColorReset)

	for attempt := 1; attempt <= 6; attempt++ {
		fmt.Printf(ColorPurple + "\n--- Attempt %d/6 ---\n" + ColorReset, attempt)
		fmt.Printf(ColorBlue + "Words remaining: %d\n" + ColorReset, len(word_list))

		// Handle case where word doesn't exist
		for {
			best_word := Get_best_word(word_list, letter_frequency)
			fmt.Printf(ColorCyan + "Suggested word: " + ColorWhite + "%s" + ColorReset + "\n", best_word)

			fmt.Print("Enter validation (or 'skip'): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "skip" {
				fmt.Println("Removing word from dictionary...")
				word_list = remove_word(word_list, best_word)
				letter_frequency = Get_letter_frequency(word_list)
				continue
			}

			// Check if solved
			if input == "ggggg" {
				fmt.Printf(ColorGreen + "\n🎉 Congratulations! You solved it with '%s' in %d attempts!\n" + ColorReset, best_word, attempt)
				return
			}

			// Validate input
			if len(input) != 5 {
				fmt.Println(ColorRed + "Please enter exactly 5 characters (g/y/b)" + ColorReset)
				continue
			}

			valid := true
			for _, char := range input {
				if char != 'g' && char != 'y' && char != 'b' {
					fmt.Println(ColorRed + "Please use only 'g', 'y', or 'b' characters" + ColorReset)
					valid = false
					break
				}
			}

			if valid {
				fmt.Printf("Result: %s\n", display_colored_word(best_word, input))
				Update_letter_conditions(input, best_word)
				word_list = Filter_word_list(word_list)
				letter_frequency = Get_letter_frequency(word_list)
				break
			}
		}

		if len(word_list) == 0 {
			fmt.Println(ColorRed + "No more words available. Something might be wrong with the input." + ColorReset)
			return
		}
	}

	fmt.Println(ColorRed + "Game over! Maximum attempts reached." + ColorReset)
}



func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(ColorPurple + "Welcome to the NYT Wordle Solver!" + ColorReset)

	for {
		play_single_game(reader)

		fmt.Print(ColorYellow + "\nWould you like to play another game? (y/n): " + ColorReset)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println(ColorCyan + "Thanks for playing!" + ColorReset)
			break
		}
	}
}