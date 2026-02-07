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



func ContainsBytes(word [5]byte, knownLetters []byte) bool {
	for _, char := range word {
		for _, knowLetter := range knownLetters {
			if char == knowLetter {
				return true
			}
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



func GetValidWords() [][5]byte {
	file, err := os.Open(validWordsPath)
	if err != nil {
		fmt.Printf("Error: Cannot find words file at '%s'\n", validWordsPath)
		fmt.Println("Please ensure you're running from the project root or src/go directory")
		fmt.Printf("Attempted paths: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	wordList := [][5]byte{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		var word [5]byte
		copy(word[:], text)
		wordList = append(wordList, word)
	}

	return wordList
}
