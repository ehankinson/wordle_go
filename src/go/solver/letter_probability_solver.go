package solver

import (
	"math/rand"
	"sort"
	"strings"
)



var letterConditions = make(map[byte]*LetterCondition)
var finalWord [5]byte
var knownLetters []byte
var trueVal = true
var falseVal = false



func init() {
	ResetGameState()
}



func ResetGameState() {
	// Reset letter conditions
	for letter := byte('a'); letter <= byte('z'); letter++ {
		letterConditions[letter] = &LetterCondition{
			Status:           nil,
			Double:           &trueVal,
			CorrectPositions: make([]int, 5),
			WrongPositions:   make([]int, 5),
		}
	}
	// Reset final word
	finalWord = [5]byte{}
	knownLetters = []byte{}
}



func GetKnownLetters() []byte {
	return knownLetters
}



func GetFinalWord() [5]byte {
	return finalWord
}



func GetRandomWord(wordList [][5]byte) [5]byte {
	return wordList[rand.Intn(len(wordList))]
}



func GetLetterConditions() map[byte]*LetterCondition {
	return letterConditions
}



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



func UpdateLetterConditions(validationString string, guessedWord [5]byte) {
	seen := make(map[byte]*bool)
	validationBytes := []byte(validationString)

	for index, letter := range guessedWord {
		switch validationBytes[index] {
		case 'g':
			// If the letter is in the correct position
			letterConditions[letter].Status = &trueVal

			positions := letterConditions[letter].CorrectPositions
			positions = append(positions, index)
			letterConditions[letter].CorrectPositions = positions

			// We will add this to the 'mock' final word to help us remove unwated words
			finalWord[index] = letter
			// Since we know that this letter is in the word, we should also remove all other words that don't include it
			knownLetters = append(knownLetters, letter)

			// If we have already seen the letter and it shouldn't be in the word
			// That means that duplicates are no longer allowed

			if _, exists := seen[letter]; exists && !*letterConditions[letter].Status {
				letterConditions[letter].Double = &falseVal
			}

		case 'y':
			// If the letter is in the word, but incorrect position
			letterConditions[letter].Status = &trueVal

			positions := letterConditions[letter].WrongPositions
			positions = append(positions, index)
			letterConditions[letter].WrongPositions = positions

			// If we have already seen the letter and it shouldn't be in the word
			// That means that duplicates are no longer allowed
			if _, exists := seen[letter]; exists && !*letterConditions[letter].Status {
				letterConditions[letter].Double = &falseVal
			}

			// Since we know that this letter is in the word, we should also remove all other words that don't include it
			knownLetters = append(knownLetters, letter)

		default:
			// If we have already looked at the letter and the new is a 'b' that means there are no doubles,
			// but could still be possible for that letter to be in the word
			if _, exists := seen[letter]; exists {
				letterConditions[letter].Double = &falseVal
			} else {
				letterConditions[letter].Status = &falseVal
			}
		}

		seen[letter] = &trueVal
	}
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



func GetBestWord(wordList [][5]byte, letterFrequency [26][5]float64) [5]byte {
	rankedWords := RankedWords(wordList, letterFrequency)
	return rankedWords[0].Word
}



func RankedWords(wordList [][5]byte, letterFrequency [26][5]float64) []WordScore {
	totalWords := len(wordList)
	rankedWords := make([]WordScore, totalWords)

	for wordIndex, word := range wordList {
		var prob float64 = 0.0
		for letterIndex, letter := range word {
			index := int(letter - 'a')
			prob += letterFrequency[index][letterIndex]
		}
		rankedWords[wordIndex].Word = word
		rankedWords[wordIndex].Prob = prob
	}

	sort.Slice(rankedWords, func(i, j int) bool {
		return rankedWords[i].Prob > rankedWords[j].Prob
	})

	return rankedWords
}



func FilterWordList(wordList [][5]byte) [][5]byte {
	filteredWords := make([][5]byte, 0, len(wordList)/4)

	for _, word := range wordList {
		// when the finalWord is being built we should skip any word which does not follow the skeleton
		skip := false

		for index, letter := range word {
			if finalWord[index] == 0 {
				continue
			}

			if finalWord[index] != letter {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		// The word does not contain one of the know letters in the word, so we should skip
		if len(knownLetters) > 0 && !InKnownLetters(word, knownLetters) {
			continue
		}

		add := true

		for pos, letter := range word {

			// We have no information on this letter so we should skip
			if letterConditions[letter].Status == nil {
				continue
			}

			if !*letterConditions[letter].Status {
				add = false
				break
			}

			wrongPositions := letterConditions[letter].WrongPositions
			// When the current position is in the wrongPositions we should skip
			if len(wrongPositions) > 0 && ContainsNumber(pos, wrongPositions) {
				add = false
				break
			}
		}

		if add {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}
