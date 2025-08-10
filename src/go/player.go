package main

import (
	"fmt"
)

func sum(arr []int) int {
	sum := 0
	for _, num := range arr {
		sum += num
	}
	return sum
}



func test_games(games_to_play int) {
	games_won := 0
	final_guess := []int{}

	for i := 0; i < games_to_play; i++ {
		// fmt.Printf("Playing game %d of %d\n", i+1, games_to_play)

		Reset_game_state()
		word_list := Get_valide_words()
		// fmt.Printf("Length of word list: %d\n", len(word_list))
		word := Get_random_word(word_list)
		word = "tases"

		
		for i := 0; i < 6; i++ {
			best_word := Get_best_word(word_list)
			fmt.Printf("Best word: %s\n", best_word)
			fmt.Printf("Amount of words: %d\n", len(word_list))
			// fmt.Printf("Best word: %s\n", best_word)
			validation_string := Nyt_word_validator(word, best_word)
			if validation_string == "ggggg" {
				games_won++
				final_guess = append(final_guess, i)
				// fmt.Printf("Game won in %d guesses\n", i)
				break
			}
			Update_letter_conditions(validation_string, best_word)
			word_list = Filter_word_list(word_list)
			// fmt.Printf("Length of word list: %d\n", len(word_list))

			// if i == 5 {
			// 	fmt.Printf("Word: %s\n", word)
			// 	fmt.Printf("Game lost\n")
			// }
		}
	}

	average_guess := float64(sum(final_guess)) / float64(games_won)
	fmt.Printf("\nGames won: %d out of %d\n", games_won, games_to_play)
	fmt.Printf("Average guess: %f\n", average_guess)

}



func main() {
	games_to_play := 1
	test_games(games_to_play)
}
