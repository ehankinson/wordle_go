package solver

func ContainsRunes(word string, runes []rune) bool {
	for _, char := range word {
		for _, r := range runes {
			if char == r {
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
	for _, num := range numbers {
		if n == num {
			return true
		}
	}

	return false
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