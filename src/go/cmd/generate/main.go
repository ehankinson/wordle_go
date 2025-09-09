package main

import (
	"os"
	"fmt"
	"math"
	"time"
	"runtime"
	"strings"
	"math/rand"
	"encoding/json"
	"path/filepath"
	"wordle_go/solver"
)

func getCurrentDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

var currentDir = getCurrentDir()
var NUM_GAMES = 50000
var WORDLE_ATTEMPTS = 6
var STARTING_WORDS = solver.GetValidWords()
var DATA_FILE = filepath.Join(currentDir, "..", "..", "..", "..", "data", "trainning_set.json")

func main() {
	err := generateData(NUM_GAMES)
	if err != nil {
		fmt.Printf("Error generating data: %v\n", err)
		os.Exit(1)
	}
}



func getExplorationParams() solver.ExplorationData {
	r := rand.Float64()

	if r < 0.4 {
		return solver.ExplorationData{
			Temp:     0.3 + rand.Float64()*0.3,
			TopK:     5 + rand.Intn(5),
			Strategy: solver.EXPLOIT,
		}
	} else if r < 0.8 {
		return solver.ExplorationData{
			Temp:     0.7 + rand.Float64()*0.5,
			TopK:     10 + rand.Intn(10),
			Strategy: solver.EXPLORE,
		}
	} else {
		return solver.ExplorationData{
			Temp:     1.2 + rand.Float64()*0.8,
			TopK:     30 + rand.Intn(70),
			Strategy: solver.RANDOM,
		}
	}
}



func getWord(rankedWords []solver.WordScore, params solver.ExplorationData) string {
	// Handle empty word list
	if len(rankedWords) == 0 {
		return solver.GetRandomWord(STARTING_WORDS)
	}

	// Take only the Kth amount of words to select from
	gussableWords := rankedWords
	if len(rankedWords) > params.TopK {
		gussableWords = rankedWords[:params.TopK]
	}

	// Apply temperature-based selection
	weights := make([]float64, len(gussableWords))

	// Apply the 'softmax' function to allow better scaling on word probabilities
	// Smaller the 'temp' is the bigger the difference is from the best to worst word, and larger is the opposite
	sum := 0.0 // Sum of all weight so we can get a % scale from 0 to the total exp smoothing, instead of using a 0-1 scale
	for i, ws := range gussableWords {
		// Convert score to weight using temperature
		weight := math.Exp(ws.Prob / params.Temp)
		weights[i] = weight
		sum += weight
	}

	r := rand.Float64() * sum // Picking a random number between 0-sum [0-1]
	cumulative := 0.0
	for i, w := range weights {
		cumulative += w
		if r <= cumulative {
			return gussableWords[i].Word
		}
	}

	return gussableWords[0].Word
}



func calculateReward(guessNum int, wordsBefore int, wordsAfter int, feedback string, isFinal bool) float64 {
	reward := -1.0

	// Information reduction bonus (0 to 3 points)
	if wordsBefore > 0 {
		reductionRate := float64(wordsBefore-wordsAfter) / float64(wordsBefore)
		reward += reductionRate * 3.0
	}

	// Count new information
	greens := strings.Count(feedback, "g")
	yellows := strings.Count(feedback, "y")

	// Small bonuses for constraints discovered
	reward += float64(greens) * 0.3
	reward += float64(yellows) * 0.1

	// Terminal rewards
	if isFinal {
		if feedback == "ggggg" {
			reward += 10.0 + float64(7-guessNum)*2.0
		} else {
			reward -= 5.0
		}
	}

	return reward
}



func captureState(attemptNum int, wordList []string, letterConditions map[rune]*solver.LetterCondition, previousGuesses []string, previousFeedback []string) solver.State {
	absentLetters := []rune{}
	for letter, condition := range letterConditions {
		// We have no information about this letter
		if condition.Status == nil {
			continue
		}

		if !*condition.Status {
			absentLetters = append(absentLetters, letter)
		}
	}

	letterConstraints := solver.LetterConstraints{
		CorrectPositions: solver.GetFinalWord(),
		PresentLetters:   solver.GetKnownLetters(),
		AbsentLetters:    absentLetters,
	}
	return solver.State{
		AttemptNum:         attemptNum,
		LetterContraints:   letterConstraints,
		ValidWordsRemaning: len(wordList),
		PreviousGuesses:    previousGuesses,
		PreviousFeedback:   previousFeedback,
	}
}



func generateData(numGames int) error {
	startTime := time.Now()

	file, err := os.Create(DATA_FILE)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	file.WriteString("[\n")

	first := true
	totalExperiences := 0
	wins := 0

	for game := 0; game < numGames; game++ {
		solver.ResetGameState()

		wordList := STARTING_WORDS
		wordsBefore := len(wordList) // Number of words for reward calculation
		explorationParams := getExplorationParams()
		// finalWord := solver.GetRandomWord(wordList)
		finalWord := "quick"
		previousGuesses := []string{}
		previousFeedback := []string{}

		for att := 0; att < WORDLE_ATTEMPTS; att++ {
			currentState := captureState(att, wordList, solver.GetLetterConditions(), previousGuesses, previousFeedback)

			letterFrequency := solver.GetLetterFrequency(wordList)
			rankedWords := solver.RankedWords(wordList, letterFrequency)
			guessingWord := getWord(rankedWords, explorationParams)
			previousGuesses = append(previousGuesses, guessingWord)
			validationString := solver.NYTWordValidator(finalWord, guessingWord)
			previousFeedback = append(previousFeedback, validationString)

			solver.UpdateLetterConditions(validationString, guessingWord)
			wordList = solver.FilterWordList(wordList)
			wordsAfter := len(wordList)

			isWinner := validationString == "ggggg"
			isFinal := isWinner || att == WORDLE_ATTEMPTS-1

			reward := calculateReward(
				att+1,
				wordsBefore,
				wordsAfter,
				validationString,
				isFinal,
			)

			exp := solver.Experience{
				State:    currentState,
				Action:   guessingWord,
				Feedback: validationString,
				Reward:   reward,
				Done:     isFinal,
			}

			if !first {
				file.WriteString(",\n")
			}
			first = false

			jsonBytes, err := json.Marshal(exp)
			if err != nil {
				return fmt.Errorf("failed to marshal experience: %v", err)
			}

			file.WriteString("  ")
			file.Write(jsonBytes)

			totalExperiences++
			wordsBefore = wordsAfter

			if isWinner {
				wins++
				break
			}
		}

		if (game+1)%100 == 0 {
			elapsed := time.Since(startTime)
			fmt.Printf("Generated %d experiences (%.2f seconds elapsed)\n", totalExperiences, elapsed.Seconds())
			file.Sync()
		}
	}

	file.WriteString("\n]")

	totalElapsed := time.Since(startTime)
	fmt.Printf("\nCompleted! Generated %d experiences in %.2f seconds (%.2f min)\n",
		totalExperiences, totalElapsed.Seconds(), totalElapsed.Minutes())
	fmt.Printf("Games WON: %d, out of %d\n", wins, numGames)

	return nil
}