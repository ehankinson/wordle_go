#!/usr/bin/env python3
"""
Wordly Web Player - Uses Selenium to play Wordly.org
Communicates with Go solver via stdin/stdout
"""

import json
import subprocess
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options

# Number of games to play consecutively
GAMES_TO_PLAY = 10

class WordlyWebPlayer:
    def __init__(self, wordly_url="https://wordly.org/"):
        self.wordly_url = wordly_url
        self.driver = None
        self.go_process = None

    def setup_driver(self, headless=False):
        """Setup Chrome WebDriver"""
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
        """Start the Go solver in automated mode"""
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
            print("Make sure you've compiled the Go program: go build -o wordle_solver main.go nyt_solver.go")
            return False

    def open_wordly(self):
        """Open Wordly website"""
        try:
            self.driver.get(self.wordly_url)
            print("Waiting for Wordly to load...")
            time.sleep(2)  # Let the page fully load
            
            # Wordly.org usually doesn't have a play button or modal
            # The game is ready immediately
            print("✅ Wordly loaded and ready to play!")
            return True
            
        except Exception as e:
            print(f"Error opening Wordly: {e}")
            return False

    def click_new_game(self):
        """Click the New game button to start a fresh game"""
        try:
            print("Looking for New game button...")
            new_game_button = self.driver.find_element(By.XPATH, "//button[text()='New game']")
            if new_game_button.is_displayed():
                new_game_button.click()
                print("✅ Clicked New game button")
                time.sleep(2)  # Wait for new game to load
                
                # Reset the Go solver for the new game
                self.reset_go_solver()
                return True
            else:
                print("New game button not visible")
                return False
        except Exception as e:
            print(f"Error clicking New game button: {e}")
            return False

    def reset_go_solver(self):
        """Reset the Go solver for a new game"""
        try:
            print("Resetting Go solver...")
            # Terminate the old process
            if self.go_process:
                self.go_process.terminate()
                self.go_process.wait()
            
            # Start a new Go solver process
            self.go_process = subprocess.Popen(
                ['./bin/wordle_solver', '--auto'],
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True,
                bufsize=1
            )
            print("✅ Go solver reset for new game")
            return True
        except Exception as e:
            print(f"Error resetting Go solver: {e}")
            return False

    def type_word(self, word):
        """Type a word into the Wordly game"""
        try:
            # Send keys directly to the body/document
            body = self.driver.find_element(By.TAG_NAME, "body")
            
            # Type the word
            for letter in word.lower():
                body.send_keys(letter)
                time.sleep(0.1)  # Small delay between letters
            
            # Press Enter to submit
            body.send_keys(Keys.RETURN)
            print(f"Typed word: {word}")
            
            # Wait for animation to complete
            time.sleep(1.5)
            return True
            
        except Exception as e:
            print(f"Error typing word '{word}': {e}")
            return False
    
    def check_word_accepted(self, row_number):
        """Check if the word was accepted by seeing if tiles have game colors"""
        try:
            # Find the current row's tiles
            rows = self.driver.find_elements(By.CSS_SELECTOR, "div.Row")
            if row_number >= len(rows):
                return False
                
            target_row = rows[row_number]
            tiles = target_row.find_elements(By.CSS_SELECTOR, "div.Row-letter")
            
            if len(tiles) != 5:
                return False
            
            # Check if any tile has a game color (not default background)
            for tile in tiles:
                bg_color = tile.value_of_css_property("background-color") or ""
                # Check for game colors: green, yellow, or gray
                # rgba(121, 184, 81) = green
                # rgba(243, 194, 55) = yellow  
                # rgba(164, 174, 196) = gray/absent
                if any(color in bg_color for color in ["121, 184, 81", "243, 194, 55", "164, 174, 196"]):
                    return True
            
            # If no tiles have game colors, word was not accepted
            return False
        except:
            return False

    def get_result_colors(self, row_number):
        """Extract the color results from the specified row"""
        try:
            print(f"Checking row {row_number} for results...")
            time.sleep(0.5)  # Reduced wait time
            
            # Find all rows using the exact Wordly.org selector
            rows = self.driver.find_elements(By.CSS_SELECTOR, "div.Row")
            print(f"Found {len(rows)} rows")
            
            if row_number >= len(rows):
                print(f"Row {row_number} doesn't exist (only {len(rows)} rows found)")
                return None
            
            # Get the specific row
            target_row = rows[row_number]
            
            # Find all tiles (Row-letter) in this row
            tiles = target_row.find_elements(By.CSS_SELECTOR, "div.Row-letter")
            print(f"Found {len(tiles)} tiles in row {row_number}")
            
            if len(tiles) != 5:
                print(f"Expected 5 tiles, found {len(tiles)}")
                return None
            
            # Extract colors from tiles
            result = ""
            for i, tile in enumerate(tiles):
                # Check data-animation attribute and background color
                data_animation = tile.get_attribute("data-animation") or ""
                classes = tile.get_attribute("class") or ""
                bg_color = tile.value_of_css_property("background-color") or ""
                text = tile.text or ""
                
                print(f"  Tile {i} ('{text}'): data-animation='{data_animation}', class='{classes}', bg='{bg_color}'")
                
                # Wordly.org uses class names for letter states
                if "letter-correct" in classes:
                    result += "g"
                elif "letter-elsewhere" in classes:
                    result += "y"
                elif "letter-absent" in classes:
                    result += "b"
                else:
                    # Fallback to color detection
                    if "rgba(121, 184, 81" in bg_color:  # Green
                        result += "g"
                    elif "rgba(243, 194, 55" in bg_color:  # Yellow
                        result += "y"
                    else:
                        result += "b"
            
            print(f"Row {row_number} result: {result}")
            return result
            
        except Exception as e:
            print(f"Error getting result colors: {e}")
            return None

    def communicate_with_go(self, validation_result):
        """Send validation result to Go solver and get next word"""
        try:
            # Send the validation result to Go
            self.go_process.stdin.write(validation_result + "\n")
            self.go_process.stdin.flush()
            
            # Read the response from Go
            response = self.go_process.stdout.readline().strip()
            
            if response.startswith("WORD:"):
                return response[5:]  # Return the word part
            elif response.startswith("SOLVED:"):
                parts = response.split(":")
                word = parts[1]
                attempts = parts[2]
                print(f"🎉 Puzzle solved with '{word}' in {attempts} attempts!")
                return None
            elif response.startswith("UPDATED:"):
                # Go solver updated its word list, wait for next word
                remaining_words = response[8:]
                print(f"Go solver updated word list: {remaining_words} words remaining")
                
                # Read the next response
                next_response = self.go_process.stdout.readline().strip()
                if next_response.startswith("WORD:"):
                    return next_response[5:]
                elif next_response.startswith("SOLVED:"):
                    parts = next_response.split(":")
                    word = parts[1]
                    attempts = parts[2]
                    print(f"🎉 Puzzle solved with '{word}' in {attempts} attempts!")
                    return None
            elif response.startswith("ERROR:"):
                print(f"Go solver error: {response[6:]}")
                return None
            elif response.startswith("FAILED:"):
                print("Go solver failed to find solution")
                return None
            else:
                print(f"Unexpected response from Go: {response}")
                return None
                
        except Exception as e:
            print(f"Error communicating with Go solver: {e}")
            return None

    def play_single_game(self, results):
        """Play a single game of Wordle"""
        try:
            # Get the first word from Go solver
            print("Waiting for first word from Go solver...")
            first_response = self.go_process.stdout.readline().strip()
            print(f"Go solver response: '{first_response}'")
            
            if not first_response.startswith("WORD:"):
                print(f"Expected first word, got: {first_response}")
                return False
            
            current_word = first_response[5:]
            print(f"Starting game with word: '{current_word}'")
            row = 0
            
            while current_word and row < 6:
                print(f"\nAttempt {row + 1}: Trying word '{current_word}'")
                
                # Type the word
                if not self.type_word(current_word):
                    break
                
                # Check if word was accepted by looking for colored tiles
                if not self.check_word_accepted(row):
                    print(f"❌ Word '{current_word}' not accepted by Wordly")
                    
                    # Delete the word
                    body = self.driver.find_element(By.TAG_NAME, "body")
                    for _ in range(5):  # Backspace 5 times
                        body.send_keys(Keys.BACKSPACE)
                        time.sleep(0.05)
                    
                    # Request new word from Go solver
                    print("Requesting new word from Go solver...")
                    self.go_process.stdin.write("SKIP\n")
                    self.go_process.stdin.flush()
                    
                    # Get new word
                    skip_response = self.go_process.stdout.readline().strip()
                    if skip_response.startswith("WORD:"):
                        current_word = skip_response[5:]
                        print(f"Got new word: '{current_word}'")
                        continue  # Try again with new word
                    else:
                        print(f"Unexpected response after SKIP: {skip_response}")
                        break
                
                # Get the color result
                result = self.get_result_colors(row)
                if not result:
                    print("Could not get result colors, stopping")
                    break
                
                print(f"Result: {result}")
                
                # Check if solved
                if result == "ggggg":
                    print(f"🎉 Congratulations! Solved in {row + 1} attempts!")
                    results["games_won"] += 1
                    results[f"{row + 1}_guess"] += 1
                    self.go_process.stdin.write("ggggg\n")
                    self.go_process.stdin.flush()
                    break
                
                # Send result to Go and get next word
                next_word = self.communicate_with_go(result)
                if not next_word:
                    break
                
                current_word = next_word
                row += 1
            
            if row >= 6:
                print("Game over - maximum attempts reached")
            
            # Keep browser open for a moment
            time.sleep(2)
            return True
            
        except Exception as e:
            print(f"Error during game play: {e}")
            return False

    def play_game(self, results):
        """Main game loop to play multiple games"""
        if not self.setup_driver():
            return False
        
        if not self.start_go_solver():
            self.driver.quit()
            return False
        
        if not self.open_wordly():
            self.cleanup()
            return False
        
        try:
            games_played = 0
            
            for game_num in range(GAMES_TO_PLAY):
                results["games_played"] += 1
                print(f"\n{'='*50}")
                print(f"GAME {game_num + 1} of {GAMES_TO_PLAY}")
                print('='*50)
                
                # Play a single game
                success = self.play_single_game(results)
                
                if success:
                    games_played += 1
                    print(f"✅ Game {game_num + 1} completed successfully!")
                    
                    # If not the last game, click "New game"
                    if game_num < GAMES_TO_PLAY - 1:
                        if self.click_new_game():
                            print("Starting new game...")
                        else:
                            print("Failed to start new game, stopping")
                            break
                else:
                    print(f"❌ Game {game_num + 1} failed")
                    break
            
            print(f"\n{'='*50}")
            print(f"FINAL RESULTS: Completed {games_played} out of {GAMES_TO_PLAY} games")
            print('='*50)
            
            return True
            
        except Exception as e:
            print(f"Error during multi-game play: {e}")
            return False
        finally:
            self.cleanup()

    def cleanup(self):
        """Clean up resources"""
        if self.go_process:
            self.go_process.terminate()
            self.go_process.wait()
        
        if self.driver:
            self.driver.quit()

def main():
    with open('data/wordly_results.json', 'r') as f:
        results = json.load(f)

    player = WordlyWebPlayer()
    success = player.play_game(results)
    
    if success:
        print("\nGame completed successfully!")
    else:
        print("\nGame failed to complete.")
    
    with open('data/wordly_results.json', 'w') as f:
        json.dump(results, f, indent=4)

if __name__ == "__main__":
    main()