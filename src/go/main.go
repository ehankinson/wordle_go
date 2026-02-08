package main

import (
	"fmt"
	"time"
	"wordle_go/solver"
)

func main() {
	fmt.Println("Welcome to GO Wordle Solver!")
	fmt.Println("\nThis is the main entry point for the Wordle solver project.")
	fmt.Println("You can run the following commands:")
	fmt.Println("  go run ./cmd/nyt        - Run the interactive NYT Wordle solver")
	fmt.Println("  go run ./cmd/generate   - Generate training data for ML")
	
	// Example of using the solver package
	fmt.Println("\nDemonstrating solver package:")
	words := solver.GetValidWords()
	
	start := time.Now()
	entropyMap := solver.GetEntropyMap(words)
	end := time.Now()
	bestWord := solver.GetBestEntropyWord(entropyMap)
	fmt.Printf("Time taken: %s\n", end.Sub(start))
	fmt.Printf("Best word: %s\n", string(bestWord[:]))

	feedback1 := [5]byte{'y', 'y', 'b', 'y', 'b'}
	words = solver.FilterWords(feedback1, bestWord, words)
	fmt.Printf("Words remaining: %d\n", len(words))
	// words = solver.FilterWords(feedback1, bestWord, words)
	// entropyMap = solver.GetEntropyMap(words)
	// bestWord = solver.GetBestEntropyWord(entropyMap)
	// fmt.Printf("Best word: %s\n", string(bestWord[:]))

	// feedback2 := [5]byte{'b', 'b', 'y', 'b', 'b'}
	// words = solver.FilterWords(feedback2, bestWord, words)
	// entropyMap = solver.GetEntropyMap(words)
	// bestWord = solver.GetBestEntropyWord(entropyMap)
	// fmt.Printf("Best word: %s\n", string(bestWord[:]))

	// feedback3 := [5]byte{'b', 'g', 'b', 'b', 'y'}
	// words = solver.FilterWords(feedback3, bestWord, words)
	// entropyMap = solver.GetEntropyMap(words)
	// bestWord = solver.GetBestEntropyWord(entropyMap)
	// fmt.Printf("Best word: %s\n", string(bestWord[:]))

	// feedback4 := [5]byte{'b', 'y', 'y', 'y', 'b'}
	// words = solver.FilterWords(feedback4, bestWord, words)
	// entropyMap = solver.GetEntropyMap(words)
	// bestWord = solver.GetBestEntropyWord(entropyMap)
	// fmt.Printf("Best word: %s\n", string(bestWord[:]))
}