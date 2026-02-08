package solver

import (
	"math"
)

func patterIndex(pattern [5]byte) int {
	index := 0
	for i := 0; i < 5; i++ {
		index *= 3

		switch pattern[i] {
		case 'b':
			index += 0
		case 'y':
			index += 1
		case 'g':
			index += 2
		}
	}
	return index
}



func EntropyWordBuckets(currentWord [5]byte, wordList [][5]byte) [243]int {
	var count [243]int

	for _, word := range wordList {
		
		result := [5]byte{}
		for i := 0; i < 5; i++ {
			if word[i] == currentWord[i] {
				result[i] = 'g'
			} else if ContainsByte(word, currentWord[i]) {
				result[i] = 'y'
			} else {
				result[i] = 'b'
			}
		}

		index := patterIndex(result)
		count[index]++
	}

	return count
}



func GetEntropyMap(wordList [][5]byte) map[[5]byte]float64 {
	totalWords := float64(len(wordList))
	entropyMap := make(map[[5]byte]float64)
	
	for _, word := range wordList {
		entropy := 0.0
		bucket := EntropyWordBuckets(word, wordList)
		for i := 0; i < 243; i++ {
			count := float64(bucket[i])
			if count == 0 {
				continue
			}
			entropy += (count / totalWords) * math.Log2(totalWords / count)
		}
		entropyMap[word] = entropy
	}

	return entropyMap
}



func GetBestEntropyWord(entropyMap map[[5]byte]float64) [5]byte {
	bestEntropy := 0.0
	bestWord := [5]byte{}
	for word, entropy := range entropyMap {
		if entropy > bestEntropy {
			bestEntropy = entropy
			bestWord = word
		}
	}

	return bestWord
}
