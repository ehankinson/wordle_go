package main

import (
	"runtime"
	"path/filepath"
)

// Get the directory of this source file at runtime (like Python's __file__)
func get_current_dir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
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
