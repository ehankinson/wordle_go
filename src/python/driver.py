#!/usr/bin/env python3
"""
Base WebDriver class for Wordle-like games
Provides common functionality for web automation
"""

import time
import subprocess

from abc import ABC
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.chrome.options import Options

CHECKS = {
    "nty": {
        "correct": "correct",
        "present": "present",
        "absent": "absent",
        "does_not_exist": "tbd",
        "row_selector": "[role='group']:nth-child({row_number})",
        "tile_selector": "[data-testid='tile']",
        "game_grid": "[data-testid='game-board']",
        "attibute_state": "data-state"
    },
    "wordly": {
        "correct": "Row-letter letter-correct",
        "present": "Row-letter letter-elsewhere",
        "absent": "Row-letter letter-absent",
        "does_not_exist": "Row-letter selected",
        "row_selector": "div.Row:nth-child({row_number})",
        "tile_selector": "div.Row-letter",
        "game_grid": "div.game_rows",
        "attibute_state": "class"
    }
}

class WordleDriver(ABC):
    """Base class for Wordle web drivers"""
    
    def __init__(self, game_url, wordle_type):
        self.driver = None
        self.go_process = None
        self.game_url = game_url
        self.wordle_type = wordle_type
        self.config = CHECKS[wordle_type]
    


    def setup_chrome_driver(self, headless=False):
        """Setup Chrome WebDriver with common options"""
        chrome_options = Options()
        if headless:
            chrome_options.add_argument("--headless")
        chrome_options.add_argument("--no-sandbox")
        chrome_options.add_argument("--disable-dev-shm-usage")
        chrome_options.add_argument("--disable-gpu")
        
        try:
            self.driver = webdriver.Chrome(options=chrome_options)
            self.driver.implicitly_wait(10)
            return True
        except Exception as e:
            print(f"Error setting up Chrome driver: {e}")
            print("Make sure you have Chrome and ChromeDriver installed")
            return False



    def start_go_solver(self):
        """Start the Go solver process"""
        try:
            self.go_process = subprocess.Popen(
                ['./bin/wordle_solver', '--auto'],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                bufsize=1
            )
            return True
        except Exception as e:
            print(f"Error starting Go solver: {e}")
            print("Make sure you've compiled the Go program")
            return False



    def open_game(self):
        """Open the game website"""
        try:
            self.driver.get(self.game_url)
            print(f"Opening {self.game_url}...")
            time.sleep(0.5)  # Let page load
            return True
        except Exception as e:
            print(f"Error opening game: {e}")
            return False



    def type_word(self, word):
        """Type a word into the game"""
        try:
            body = self.driver.find_element(By.TAG_NAME, "body")
            
            # Type each letter
            for letter in word.lower():
                body.send_keys(letter)
                time.sleep(0.1)
            
            # Submit the word
            body.send_keys(Keys.RETURN)
            print(f"Typed word: {word}")
            
            # Wait for animation
            time.sleep(0.5)
            return True
            
        except Exception as e:
            print(f"Error typing word '{word}': {e}")
            return False



    def delete_word(self, length=5):
        """Delete a word by sending backspace keys"""
        try:
            body = self.driver.find_element(By.TAG_NAME, "body")
            for _ in range(length):
                body.send_keys(Keys.BACKSPACE)
                time.sleep(0.05)
            return True
        except Exception as e:
            print(f"Error deleting word: {e}")
            return False



    def communicate_with_solver(self, message):
        """Send message to Go solver and get response"""
        try:
            self.go_process.stdin.write(message + "\n")
            self.go_process.stdin.flush()
            response = self.go_process.stdout.readline().strip()
            return response
        except Exception as e:
            print(f"Error communicating with solver: {e}")
            return None



    def get_next_word_from_solver(self, validation_result):
        """Get next word from solver based on validation result"""
        response = self.communicate_with_solver(validation_result)
        
        if response and response.startswith("WORD:"):
            return response[5:]
        elif response and response.startswith("SOLVED:"):
            parts = response.split(":")
            word = parts[1] if len(parts) > 1 else ""
            attempts = parts[2] if len(parts) > 2 else ""
            print(f"🎉 Puzzle solved with '{word}' in {attempts} attempts!")
            return None
        elif response and response.startswith("UPDATED:"):
            # Solver updated, get next word
            remaining = response[8:]
            print(f"Solver updated: {remaining} words remaining")
            next_response = self.go_process.stdout.readline().strip()
            if next_response.startswith("WORD:"):
                return next_response[5:]
        elif response and response.startswith("ERROR:"):
            print(f"Solver error: {response[6:]}")
        elif response and response.startswith("FAILED:"):
            print("Solver failed to find solution")
        
        return None



    def game_won(self, result: str, row: int) -> bool:
        if result == "ggggg":
            print(f"🎉 Solved in {row + 1} attempts!")
            self.communicate_with_solver("ggggg")
            
            # Try to share results
            try:
                share_button = self.driver.find_element(
                    By.CSS_SELECTOR, 
                    "button[data-testid='share-button']"
                )
                if share_button.is_displayed():
                    share_button.click()
                    print("✅ Results copied to clipboard!")
            except:
                pass
            
            return True
        return False
    


    def skip_word(self, current_word):
        print(f"❌ Word '{current_word}' not accepted")
        self.delete_word()
        
        # Request new word
        skip_response = self.communicate_with_solver("SKIP")
        if skip_response.startswith("WORD:"):
            new_word = skip_response[5:]
            return new_word
        else:
            return None



    def cleanup(self):
        """Clean up resources"""
        if self.go_process:
            self.go_process.terminate()
            self.go_process.wait()
        
        if self.driver:
            self.driver.quit()



    def reset_solver(self):
        """Reset the Go solver for a new game"""
        try:
            if self.go_process:
                self.go_process.terminate()
                self.go_process.wait()
            
            return self.start_go_solver()
        except Exception as e:
            print(f"Error resetting solver: {e}")
            return False
    


    def get_tile_colors(self, row_number):
        """Get color results from specified row, returns None if word wasn't accepted"""
        try:
            print(f"Checking row {row_number} for results...")
            time.sleep(1.5)  # Wait for animation
            
            # Directly get the specific row
            row_selector = self.config["row_selector"].format(row_number=row_number + 1)
            target_row = self.driver.find_element(By.CSS_SELECTOR, row_selector)
            tiles = target_row.find_elements(By.CSS_SELECTOR, self.config["tile_selector"])
            
            if len(tiles) != 5:
                return None
            
            result = ""
            for i, tile in enumerate(tiles):
                state = tile.get_attribute(self.config["attibute_state"])

                if i == 0 and state == self.config["does_not_exist"]:
                    return None
                
                if state == self.config["correct"]:
                    result += "g"
                elif state == self.config["present"]:
                    result += "y"
                elif state == self.config["absent"]:
                    result += "b"
                else:
                    print(f"❌ Word not accepted - tile {i} has state '{state}'")
                    return None
            
            print(f"Row {row_number} result: {result}")
            return result
            
        except Exception as e:
            print(f"Error getting colors: {e}")
            return None
    


    def get_game_grid(self):
        """Get the game grid element"""
        return self.driver.find_element(By.CSS_SELECTOR, self.config["game_grid"])
    


    def game_setup(self) -> bool:
        """Set up the game"""
        if not self.setup_chrome_driver():
            return False
        
        if not self.start_go_solver():
            self.cleanup()
            return False
        
        if not self.open_game():
            self.cleanup()
            return False

        return True



    def play_game(self, results: dict = None) -> bool:
        # Get first word from solver
        print("Getting first word from solver...")
        first_response = self.go_process.stdout.readline().strip()
        
        if not first_response.startswith("WORD:"):
            print(f"Expected first word, got: {first_response}")
            return False
        
        current_word = first_response[5:]
        print(f"Starting with: '{current_word}'")
        row = 0
        
        while current_word and row < 6:
            print(f"\nAttempt {row + 1}: '{current_word}'")
            
            if not self.type_word(current_word):
                break
            
            # Get color result (this also checks if word was accepted)
            result = self.get_tile_colors(row)
            
            if result is None:
                current_word = self.skip_word(current_word)
                if current_word is None:
                    print("No more words to try")
                    break
                continue
            
            print(f"Result: {result}")
            
            # Check if solved
            if self.game_won(result, row):
                break
            
            # Get next word
            next_word = self.get_next_word_from_solver(result)
            if not next_word:
                break
            
            current_word = next_word
            row += 1
        
        if row >= 6 and result != "ggggg":
            if results is not None:
                results["games_played"] += 1
                results["wins"] += 1
                results[row + 1] += 1

            print("Game over - max attempts reached")
        
        return True