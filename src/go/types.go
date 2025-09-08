package main

type Strategy int

const (
	EXPLOIT Strategy = iota
	EXPLORE
	RANDOM
)

type ExplorationData struct {
	Temp     float64
	TopK     int
	Strategy Strategy
}

type Experience struct {
	State    State   `json:"state"`
	Action   string  `json:"action"`
	Feedback string  `json:"feedback"`
	Reward   float64 `json:"reward"`
	Done     bool    `json:"done"`
}

type State struct {
	AttemptNum         int               `json:"attempt_number"`
	LetterContraints   LetterConstraints `json:"letter_constraints"`
	ValidWordsRemaning int               `json:"valid_words_remaining"`
	PreviousGuesses    []string          `json:"previous_guesses"`
	PreviousFeedback   []string          `json:"previous_feedback"`
}

type LetterConstraints struct {
	CorrectPositions []string `json:"correct_positions"`
	PresentLetters   []string `json:"present_letters"`
	AbsentLetters    []string `json:"absent_letters"`
}