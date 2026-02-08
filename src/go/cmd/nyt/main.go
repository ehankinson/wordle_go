package main

import (
	"os"
	"fmt"
	"bufio"
	"strings"
	"wordle_go/solver"
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



func colorizeLetter(letter byte, validation byte) string {
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



func displayColoredWord(word string, validation string) string {
	result := ""
	for i := 0; i < len(word) && i < len(validation); i++ {
		result += colorizeLetter(word[i], validation[i])
	}
	return result
}



func removeWord(wordList [][5]byte, word [5]byte) [][5]byte {
	result := make([][5]byte, 0, len(wordList))

	for _, w := range wordList {
		if w != word {
			result = append(result, w)
		}
	}
	return result
}



func PlayEntropyGame(wordList [][5]byte) [5]byte{
	if len(wordList) < 2 {
		return wordList[0]
	}
	entropyMap := solver.GetEntropyMap(wordList)
	bestWord := solver.GetBestEntropyWord(entropyMap)
	return bestWord
}



func PlayLetterProbabilityGame(wordList [][5]byte) [5]byte {
	letterFrequency := solver.GetLetterFrequency(wordList)
	bestWord := solver.GetBestProbabilityWord(wordList, letterFrequency)
	return bestWord
}



func playSingleGame(reader *bufio.Reader, mode string) {
	wordList := solver.GetValidWords()

	fmt.Println(ColorCyan + "\n=== Starting New Wordle Game ===" + ColorReset)
	fmt.Println("Enter validation string: " + BgGreen + ColorWhite + " g " + ColorReset + " = green (correct), " + BgYellow + ColorBlack + " y " + ColorReset + " = yellow (wrong position), " + BgBlack + ColorWhite + " b " + ColorReset + " = black (not in word)")
	fmt.Println(ColorYellow + "Type 'skip' if the suggested word doesn't exist" + ColorReset)
	fmt.Println(ColorGreen + "Type 'ggggg' when you solve the puzzle" + ColorReset)

	for attempt := 1; attempt <= 6; attempt++ {
		fmt.Printf(ColorPurple+"\n--- Attempt %d/6 ---\n"+ColorReset, attempt)
		fmt.Printf(ColorBlue+"Words remaining: %d\n"+ColorReset, len(wordList))

		if attempt == 5 {
			fmt.Printf("Failed: Maximum attempts reached\n")
		}

		// Handle case where word doesn't exist
		for {
			bestWord := [5]byte{}
			if mode == "entropy" {
				bestWord = PlayEntropyGame(wordList)
			} else {
				bestWord = PlayLetterProbabilityGame(wordList)
			}

			fmt.Printf(ColorCyan+"Suggested word: "+ColorWhite+"%s"+ColorReset+"\n", string(bestWord[:]))

			fmt.Print("Enter validation (or 'skip'): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "skip" {
				fmt.Println("Removing word from dictionary...")
				wordList = removeWord(wordList, bestWord)
				continue
			}

			// Check if solved
			if input == "ggggg" {
				fmt.Printf(ColorGreen+"\n🎉 Congratulations! You solved it with '%s' in %d attempts!\n"+ColorReset, bestWord, attempt)
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
				originalLength := len(wordList)

				feedback := [5]byte{}
				for i := 0; i < 5; i++ {
					feedback[i] = input[i]
				}

				wordList = solver.FilterWords(feedback, bestWord, wordList)
				fmt.Printf("UPDATED:%d\nHas been filtered by %d words", len(wordList), originalLength - len(wordList))
				break
			}
		}

		if len(wordList) == 0 {
			fmt.Println(ColorRed + "No more words available. Something might be wrong with the input." + ColorReset)
			return
		}
	}

	fmt.Println(ColorRed + "Game over! Maximum attempts reached." + ColorReset)
}



func playAutomatedGame(mode string) {
	wordList := solver.GetValidWords()

	reader := bufio.NewReader(os.Stdin)

	for attempt := 1; attempt <= 6; attempt++ {
		// Handle case where word doesn't exist
		for {
			bestWord := [5]byte{}
			if mode == "entropy" {
				bestWord = PlayEntropyGame(wordList)
			} else {
				bestWord = PlayLetterProbabilityGame(wordList)
			}

			// Output suggested word for Python to use
			fmt.Printf("WORD:%s\n", string(bestWord[:]))

			// Wait for validation input from Python
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("ERROR:Failed to read input\n")
				return
			}
			input = strings.TrimSpace(input)

			if input == "SKIP" {
				wordList = removeWord(wordList, bestWord)
				continue
			}

			// Check if solved
			if input == "ggggg" {
				fmt.Printf("SOLVED:%s:%d\n", bestWord, attempt)
				return
			}

			// Validate input
			if len(input) != 5 {
				fmt.Printf("ERROR:Invalid input length\n")
				continue
			}

			valid := true
			for _, char := range input {
				if char != 'g' && char != 'y' && char != 'b' {
					fmt.Printf("ERROR:Invalid characters in input\n")
					valid = false
					break
				}
			}

			if valid {
				feedback := [5]byte{}
				for i := 0; i < 5; i++ {
					feedback[i] = input[i]
				}

				wordList = solver.FilterWords(feedback, bestWord, wordList)
				fmt.Printf("UPDATED:%d\n", len(wordList))
				break
			}
		}

		if len(wordList) == 0 {
			fmt.Printf("ERROR:No more words available\n")
			return
		}
	}

	fmt.Printf("FAILED:Maximum attempts reached\n")
}



func main() {
	// Check for automated mode
	if len(os.Args) > 1 && os.Args[1] == "--auto" {
		mode := os.Args[2]
		playAutomatedGame(mode)
		return
	}

	// Original interactive mode
	reader := bufio.NewReader(os.Stdin)
	mode := os.Args[1]

	fmt.Println(ColorPurple + "Welcome to the NYT Wordle Solver!" + ColorReset)

	for {
		playSingleGame(reader, mode)

		fmt.Print(ColorYellow + "\nWould you like to play another game? (y/n): " + ColorReset)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "y" && response != "yes" {
			fmt.Println(ColorCyan + "Thanks for playing!" + ColorReset)
			break
		}
	}
}
