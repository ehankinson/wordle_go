package solver

import (
	"os"
	"fmt"
	"bufio"
	"slices"
	"runtime"
	"path/filepath"
)

var validWordsPath = filepath.Join(getCurrentDir(), "..", "..", "..", "words", "all_valid_words.txt")


func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}



func ContainsRunes(word string, runes []rune) bool {
	for _, char := range word {
		if slices.Contains(runes, char) {
			return true
		}
	}

	return false
}

func ContainsRune(word string, r rune) bool {
	for _, char := range word {
		if r == char {
			return true
		}
	}

	return false
}



func ContainsNumber(n int, numbers []int) bool {
	return slices.Contains(numbers, n)
}



func CountRunes(word string, r rune) int {
	count := 0
	for _, char := range word {
		if r == char {
			count++
		}
	}

	return count
}



func GetValidWords() []string {
	file, err := os.Open(validWordsPath)
	if err != nil {
		fmt.Printf("Error: Cannot find words file at '%s'\n", validWordsPath)
		fmt.Println("Please ensure you're running from the project root or src/go directory")
		fmt.Printf("Attempted paths: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	wordList := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wordList = append(wordList, scanner.Text())
	}

	return wordList
}
