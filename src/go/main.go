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
	fmt.Printf("Loaded %d valid Wordle words\n", len(words))

	letterFrequency := solver.GetLetterFrequency(words)
	bestWord := solver.GetBestWord(words, letterFrequency)
	fmt.Printf("Best word: %s\n", bestWord)
	fmt.Println(len(letterFrequency))
	
	// Get a random word
	// randomWord := solver.GetRandomWord(words)
	// fmt.Printf("Random word: %s\n", randomWord)

	startTime := time.Now()
	patterns := solver.GetAllLetterPatterns()
	elapsed := time.Since(startTime)
	fmt.Println(len(patterns))
	fmt.Printf("Time taken: %s\n", elapsed)
}