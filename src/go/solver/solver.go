package solver

import (
	"os"
	"fmt"
	"sort"
	"bufio"
	"runtime"
	"strings"
	"math/rand"
	"path/filepath"
)

// Build the path relative to this file's location
var validWordsPath = filepath.Join(getCurrentDir(), "..", "..", "..", "words", "all_valid_words.txt")

var letterConditions = make(map[rune]*LetterCondition)
var finalWord = make([]rune, 5)
var knownLetters = []rune{}
var trueVal = true
var falseVal = false

func init() {
	ResetGameState()
}

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}



func ResetGameState() {
	// Reset letter conditions
	for letter := 'a'; letter <= 'z'; letter++ {
		letterConditions[letter] = &LetterCondition{
			Status:           nil,
			CorrectPositions: make([]int, 5),
			WrongPositions:   make([]int, 5),
			Double:           &trueVal,
		}
	}
	// Reset final word
	finalWord = make([]rune, 5)
	knownLetters = []rune{}
}



func GetKnownLetters() []rune {
	return knownLetters
}



func GetFinalWord() []rune {
	return finalWord
}



func GetRandomWord(wordList []string) string {
	return wordList[rand.Intn(len(wordList))]
}



func NYTWordValidator(finalWord string, guessedWord string) string {
	nytString := []string{}
	for i := 0; i < len(guessedWord); i++ {
		// If the letter is in the correct position
		if guessedWord[i] == finalWord[i] {
			nytString = append(nytString, "g")
			// If the letter is in the word, but incorrect position
		} else if ContainsRune(finalWord, rune(guessedWord[i])) {
			nytString = append(nytString, "y")
			// Letter is not in the word
		} else {
			nytString = append(nytString, "b")
		}
	}
	return strings.Join(nytString, "")
}



func UpdateLetterConditions(validationString string, guessedWord string) {
	seen := make(map[rune]*bool)
	for index, letter := range guessedWord {
		switch rune(validationString[index]) {
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



func GetLetterFrequency(wordList []string) map[rune]map[int]float64 {
	letterCount := make(map[rune]map[int]float64)
	for _, word := range wordList {
		for index, letter := range word {
			if letterCount[letter] == nil {
				letterCount[letter] = make(map[int]float64)
			}
			letterCount[letter][index]++
		}
	}

	totalLetter := len(letterCount)
	letterFrequency := make(map[rune]map[int]float64)
	for abr, positions := range letterCount {
		letterFrequency[abr] = make(map[int]float64)
		for position, count := range positions {
			letterFrequency[abr][position] = float64(count) / float64(totalLetter)
		}
	}

	return letterFrequency
}



func GetBestWord(wordList []string, letterFrequency map[rune]map[int]float64) string {
	rankedWords := RankedWords(wordList, letterFrequency)
	return rankedWords[0].Word
}



func RankedWords(wordList []string, letterFrequency map[rune]map[int]float64) []WordScore {
	rankedWords := make([]WordScore, 0, len(letterFrequency))
	for _, word := range wordList {

		var prob float64 = 0.0

		for i, letter := range word {
			count := CountRunes(word, letter)
			prob += letterFrequency[letter][i] / float64(count)
		}

		rankedWords = append(rankedWords, WordScore{Word: word, Prob: prob})
	}

	sort.Slice(rankedWords, func(i, j int) bool {
		return rankedWords[i].Prob > rankedWords[j].Prob
	})

	return rankedWords
}



func FilterWordList(wordList []string) []string {
	filteredWords := []string{}

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
		if len(knownLetters) > 0 && !ContainsRunes(word, knownLetters) {
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



// Export access to letter conditions for state capture
func GetLetterConditions() map[rune]*LetterCondition {
	return letterConditions
}