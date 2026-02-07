package solver

import "slices"

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
