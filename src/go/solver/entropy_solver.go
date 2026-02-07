package solver


func GetAllLetterPatterns() [243][5]byte {
	total := 243 // 3^5
	states := [3]byte{'b', 'y', 'g'}
	patterns := [243][5]byte{}

	for i := 0; i < total; i++ {
		count := i
		var pattern [5]byte

		for j := 0; j < 5; j++ {
			pattern[4 - j] = states[count % 3]
			count /= 3
		}

		patterns[i] = pattern
	}
	
	return patterns
}



// func filterWords(pattern [5]byte, wordList []string) []string {

// }



// func GetWordEntropy(word string, patterns [243][5]byte) float64 {
// 	entropy := 0.0
// 	for _, pattern := range patterns {

// 	}
// }



// func GetEntropyMap(wordList []string) map[string]float64 {
// 	entropyMap := make(map[string]float64)
// 	for _, word := range wordList {
// 		entropyMap[word] = GetWordEntropy(word)
// 	}
// }