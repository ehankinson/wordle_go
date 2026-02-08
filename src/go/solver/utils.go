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



func HasCommonByte(word1 [5]byte, word2 [5]byte) bool {
	for _, x := range word1 {
		for _, y := range word2 {
			if x == y {
				return true
			}
		}
	}

	return false
}



func InKnownLetters(word [5]byte, knownLetters []byte) bool {
	for _, char := range word {
		for _, knowLetter := range knownLetters {
			if char == knowLetter {
				return true
			}
		}
	}

	return false
}



func ContainsByte(word [5]byte, b byte) bool {
	for _, x := range word {
		if x == b {
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



func Pattern(answer [5]byte, guess [5]byte) [5]byte {
	result := [5]byte{'b', 'b', 'b', 'b', 'b'}

	var count [26]int
	for i := 0; i < 5; i++ {
		count[answer[i]-'a']++
	}

	for i := 0; i < 5; i++ {
		if answer[i] == guess[i] {
			result[i] = 'g'
			count[answer[i]-'a']--
		}
	}

	for i := 0; i < 5; i++ {
		if result[i] == 'g' {
			continue
		}

		letter := guess[i] - 'a'
		if count[letter] > 0 {
			result[i] = 'y'
			count[letter]--
		}
	}

	return result
}



func FilterWords(feedback [5]byte, guess [5]byte, wordList [][5]byte) [][5]byte {
	results := [][5]byte{}

	for _, word := range wordList {

		pattern := Pattern(word, guess)
		if pattern == feedback {
			results = append(results, word)
		}
	}

	return results
}
