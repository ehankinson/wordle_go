import os
import json
import numpy as np

from typing import Set, Dict


CURR_DIR = os.path.dirname(os.path.abspath(__file__))
MAX_ATTEMPTS = 6.0
MAX_POSITIONS = 5.0
MAX_LETTERS = 26.0
MAX_VALID_WORDS = 14855.0  # Total valid Wordle words (from your data analysis)
STATE_VECTOR_SIZE = 1000


class LoadData:

    def __init__(self, data_path: str = "data/trainning_set.json"):
        path: str = os.path.join(CURR_DIR, '..', '..', '..', data_path) if data_path != 'data/trainning_set.json' else data_path
        self.data = self._load_json(path)
        self.num_games = len(self.data)
        self._build_results()
    


    def __repr__(self):
        return f"Experience count: {self.get_experience_count()}, Vocab size: {self.get_vocab_size()}"
            
        
    
    def _load_json(self, path: str) -> list[dict]:
        with open(path, 'r') as f:
            return json.load(f)
    


    def _build_results(self) -> None:
        self.total_wins = 0
        self.max_reward = 0
        self.average_reward = 0
        self.experience_count = 0
        self.unique_words: Set[str] = set()
        self.starting_words: Dict[str, Dict[str, int]] = {}
        self.games_won: Dict[int, int] = {1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0}
        
        for game in self.data:
            did_win = False
            starting_word = ""

            for att, exp in enumerate(game):
                word = exp["action"]
                reward = exp["reward"]
                if att == 0:
                    starting_word = word
                    if starting_word not in self.starting_words:
                        self.starting_words[starting_word] = {'games_started': 0, 'games_won': 0}

                    self.starting_words[starting_word]['games_started'] += 1
                
                if att == len(game) - 1:
                    did_win = exp["feedback"] == "ggggg"
                    if did_win:
                        self.starting_words[starting_word]['games_won'] += 1
                        self.games_won[att + 1] += 1
                        self.total_wins += 1
                
                if word not in self.unique_words:
                    self.unique_words.add(word)
                
                self.max_reward = max(self.max_reward, reward)
                self.average_reward += reward
                
            self.experience_count += len(game)
        self.average_reward /= self.experience_count
    


    def _encode_feedback(self, letter_constrains) -> np.ndarray:
        matrix = np.zeros((26, 5))
        correct_positions = letter_constrains["correct_positions"]
        for pos, ascii_code in enumerate(correct_positions):
            if ascii_code != 0:
                letter_idx = ascii_code - ord('a')
                matrix[letter_idx, pos] = 2
        
        for ascii_code in letter_constrains["absent_letters"]:
            letter_idx = ascii_code - ord('a')
            for pos in range(5):
                if matrix[letter_idx, pos] == 0:
                    matrix[letter_idx, pos] = -1

        for ascii_code in letter_constrains["present_letters"]:
            letter_idx = ascii_code - ord('a')
            for pos in range(5):
                if matrix[letter_idx, pos] == 0:
                    matrix[letter_idx, pos] = 1
        
        return matrix.flatten()
    


    def _encode_previous_guesses(self, guesses: list[str]) -> np.ndarray:
        encoded = np.zeros((5, 5, 26))

        for guess_ids, word in enumerate(guesses):
            for pos, letter in enumerate(word):
                letter_idx = ord(letter) - ord('a')
                encoded[guess_ids, pos, letter_idx] = 1
        
        return encoded.flatten()



    def get_experience_count(self) -> int:
        return self.experience_count
    


    def get_vocab_size(self) -> int:
        return len(self.unique_words)
    


    def most_common_starting_words(self) -> list[str]:
        return sorted(self.starting_words.items(), key=lambda x: x[1]['games_started'], reverse=True)[:10]
    


    def starting_word_win_rate(self) -> dict[str, float]:
        starting_word_win_rate = {}
        for starting_word, data in self.starting_words.items():
            starting_word_win_rate[starting_word] = data["games_won"] / data["games_started"]

        return starting_word_win_rate
    


    def get_average_reward(self) -> float:
        return self.average_reward
    


    def get_max_reward(self) -> float:
        return self.max_reward
    


    def get_win_distribution(self) -> dict[int, str]:
        win_distribution = {}
        for game in self.games_won:
            win_distribution[game] = f"{self.games_won[game] / self.total_wins * 100:.2f}%"
        
        return win_distribution
    


    def encode_state(self, state: dict) -> np.ndarray:
        feedback_encoding = self._encode_feedback(state["state"]["letter_constraints"])
        previous_guesses_encoding = self._encode_previous_guesses(state["state"]["previous_guesses"])
        numerical = np.array([
            state["state"]["attempt_number"] / MAX_ATTEMPTS,
            state["state"]["valid_words_remaining"] / MAX_VALID_WORDS,
            (6 - state["state"]["attempt_number"]) / MAX_ATTEMPTS,
            len(state["state"]["letter_constraints"]["absent_letters"]) / MAX_LETTERS,
            len(state["state"]["letter_constraints"]["present_letters"]) / MAX_LETTERS,
            sum(1 for pos in state["state"]["letter_constraints"]["correct_positions"] if pos != 0) / MAX_POSITIONS,
        ])

        full_encoding = np.concatenate([feedback_encoding, previous_guesses_encoding, numerical])

        if len(full_encoding) != STATE_VECTOR_SIZE:
            padding = np.zeros(STATE_VECTOR_SIZE - len(full_encoding))
            full_encoding = np.concatenate([full_encoding, padding])

        return full_encoding
                    



def print_results(what_it_is: str, data: any):
    print("=" * 100)
    print(what_it_is)
    if isinstance(data, dict):
        for key, value in data.items():
            print(f"{key}: {value}")
    elif isinstance(data, list):
        for item in data:
            print(item)
    else:
        print(data)


if __name__ == "__main__":
    load_data = LoadData()
    state = load_data.data[0][5]
    load_data.encode_state(state)
    print_results("Experience count", load_data.get_experience_count())
    print_results("Vocab size", load_data.get_vocab_size()) 
    print_results("Most common starting words", load_data.most_common_starting_words())
    print_results("Starting word win rate", load_data.starting_word_win_rate())
    print_results("Average reward", load_data.get_average_reward())
    print_results("Max reward", load_data.get_max_reward())
    print_results("Win distribution", load_data.get_win_distribution())