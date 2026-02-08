package solver

import (
	"sort"
	"strings"
)



func NYTWordValidator(finalWord string, guessedWord string) string {
	nytString := []string{}
	for i := 0; i < len(guessedWord); i++ {
		// If the letter is in the correct position
		if guessedWord[i] == finalWord[i] {
			nytString = append(nytString, "g")
			// If the letter is in the word, but incorrect position
		// } else if ContainsByte(finalWord, rune(guessedWord[i])) {
		// 	nytString = append(nytString, "y")
			// Letter is not in the word
		} else {
			nytString = append(nytString, "b")
		}
	}
	return strings.Join(nytString, "")
}



func GetLetterFrequency(wordList [][5]byte) [26][5]float64 {
	var letterFrequency [26][5]float64
	for _, word := range wordList {
		for index, letter := range word {
			letterIndex := int(letter - 'a')
			letterFrequency[letterIndex][index] += 1.0
		}
	}

	totalWords := float64(len(wordList))
	for i := 0; i < 26; i++ {
		for j := 0; j < 5; j++ {
			letterFrequency[i][j] /= totalWords
		}
	}

	return letterFrequency
}



func GetBestProbabilityWord(wordList [][5]byte, letterFrequency [26][5]float64) [5]byte {
	rankedWords := RankedWords(wordList, letterFrequency)
	return rankedWords[0].Word
}



func RankedWords(wordList [][5]byte, letterFrequency [26][5]float64) []WordScore {
	totalWords := len(wordList)
	rankedWords := make([]WordScore, totalWords)

	for wordIndex, word := range wordList {
		var prob = 0.0
		var count [26]int
		for letterIndex, letter := range word {
			count[letter - 'a']++
			index := int(letter - 'a')
			prob += letterFrequency[index][letterIndex] / float64(count[letter - 'a'])
		}
		rankedWords[wordIndex].Word = word
		rankedWords[wordIndex].Prob = prob
	}

	sort.Slice(rankedWords, func(i, j int) bool {
		return rankedWords[i].Prob > rankedWords[j].Prob
	})

	return rankedWords
}
