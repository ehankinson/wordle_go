import os
import json

from pydantic import BaseModel
from driver import WordleDriver
from typing import Optional, List

CURR_DIR = os.path.dirname(os.path.abspath(__file__))


class Result(BaseModel):
    solved: bool
    word: str
    attempts: int
    guesses: List[str]
    remark: Optional[str]


def nyt_word_validator(final_word, guessed_word):
    nyt_string: List[str] = []
    for i in range(len(guessed_word)):
        if guessed_word[i] == final_word[i]:
            nyt_string.append("g")
        elif guessed_word[i] in final_word:
            nyt_string.append("y")
        else:
            nyt_string.append("b")
    return "".join(nyt_string)



def how_many_words_can_is_solve() -> None:
    attempts: List[Result] = []
    driver = WordleDriver()
    if not driver.start_go_solver():
        print("Failed to start Go solver")
        return

    try:
        with open(os.path.join(CURR_DIR, "..", "..", "words", "all_valid_words.txt"), "r") as f:
            words = f.read().splitlines()

        wins = 0
        for i, word in enumerate(words):  
            guess = driver.go_process.stdout.readline().strip()[5:]  
            guesses: List[str] = []
            solved = False
            remark = None
            
            for att in range(6):
                guesses.append(guess)
                feedback = nyt_word_validator(word, guess)
                if feedback == "ggggg":
                    solved = True
                    wins += 1
                    break

                guess = driver.get_next_word_from_solver(feedback)
                if guess is None and att < 5:  # Not the last attempt
                    remark = "Solver failed to find solution"
                    break

            res = Result(solved=solved, word=word, attempts=att + 1, guesses=guesses, remark=remark)
            attempts.append(res)
            driver.reset_solver()

        with open(os.path.join(CURR_DIR, "..", "..", "words", "total_results.json"), "w") as f:
            json.dump([result.model_dump() for result in attempts], f, indent=4)

        print(f"The current win rate is {(wins / len(words)) * 100:.2f}%")
        
    finally:
        # Clean up resources
        driver.cleanup()


if __name__ == "__main__":
    how_many_words_can_is_solve()

