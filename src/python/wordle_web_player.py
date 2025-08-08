#!/usr/bin/env python3
"""
Wordle Web Player - Uses Selenium to play Wordle on the web
Communicates with Go solver via stdin/stdout
"""

import subprocess
import time
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service

class WordleWebPlayer:
    def __init__(self, wordle_url="https://www.nytimes.com/games/wordle/index.html"):
        self.wordle_url = wordle_url
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
                ['./wordle_solver', '--auto'],
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



    def open_wordle(self):
        """Open Wordle website and handle initial popup"""
        try:
            self.driver.get(self.wordle_url)
            time.sleep(1.5)
            
            # Try to close any welcome/instruction popups
            try:
                close_buttons = self.driver.find_elements(By.CSS_SELECTOR, "[data-testid='icon-close'], .close, .modal-close, button[aria-label='Close']")
                for button in close_buttons:
                    if button.is_displayed():
                        button.click()
                        time.sleep(1)
                        break
            except:
                pass
            
            # Look for and click the Play button
            try:
                print("Looking for Play button...")
                
                # Use the exact selector from the NYT Wordle structure
                play_button_clicked = False
                try:
                    play_button = self.driver.find_element(By.CSS_SELECTOR, "button[data-testid='Play']")
                    if play_button.is_displayed() and play_button.is_enabled():
                        print("✅ Found Play button instantly!")
                        play_button.click()
                        play_button_clicked = True
                    else:
                        print("Play button not clickable, continuing...")
                except:
                    print("Play button not found, continuing...")
                
                if play_button_clicked:
                    print("Waiting for game board to load...")
                    
                    # Handle "How To Play" modal that might appear after clicking Play
                    try:
                        print("Looking for 'How To Play' modal close button...")
                        time.sleep(1)  # Brief wait for modal to appear
                        
                        # Use the exact selector from the HTML structure
                        try:
                            close_button = self.driver.find_element(By.CSS_SELECTOR, "[data-testid='icon-close']")
                            if close_button.is_displayed():
                                print("✅ Found modal close button instantly!")
                                close_button.click()
                                time.sleep(0.5)  # Brief wait for modal to close
                            else:
                                print("Modal close button not visible, continuing...")
                        except:
                            print("Modal close button not found, continuing...")
                            
                    except Exception as modal_e:
                        print(f"Error handling modal: {modal_e}")
                        # Continue anyway
                else:
                    print("Play button not found, proceeding anyway...")
                    
            except Exception as e:
                print(f"Error clicking Play button: {e}")
                # Continue anyway as the game might already be accessible
            
            # Wait for the exact NYT Wordle board to appear
            print("Looking for game board...")
            try:
                # Use the exact selector from the HTML structure
                game_board = WebDriverWait(self.driver, 1.5).until(
                    EC.presence_of_element_located((By.CSS_SELECTOR, "div[class*='Board-module_board']"))
                )
                print("✅ Found game board instantly!")
            except:
                print("Could not find game board, page may not be loaded properly")
                raise Exception("NYT Wordle board not found")
            
            print("Successfully loaded Wordle page and ready to play!")
            return True
            
        except Exception as e:
            print(f"Error opening Wordle: {e}")
            return False



    def type_word(self, word):
        """Type a word into the Wordle game"""
        try:
            # Find the active input or just send keys to the body
            body = self.driver.find_element(By.TAG_NAME, "body")
            
            # Clear any existing input and type the word
            for letter in word.lower():
                body.send_keys(letter)
                time.sleep(0.1)
            
            # Press Enter to submit
            body.send_keys(Keys.RETURN)
            print("Word submitted, waiting for animation to complete...")
            time.sleep(3)  # Wait longer for Wordle animation to complete
            
            print(f"Typed word: {word}")
            return True
            
        except Exception as e:
            print(f"Error typing word '{word}': {e}")
            return False



    def get_result_colors(self, row_number):
        """Extract the color results from the specified row"""
        try:
            # Wait for Wordle animations to complete and tiles to show final colors
            print("Waiting for tile colors to finalize...")
            time.sleep(2)  # Give extra time for animations
            
            # Wait for tiles to have final states (not 'empty' or 'tbd')
            max_attempts = 10
            for attempt in range(max_attempts):
                try:
                    rows = self.driver.find_elements(By.CSS_SELECTOR, "div[class*='Row-module_row']")
                    if row_number < len(rows):
                        target_row = rows[row_number]
                        tiles = target_row.find_elements(By.CSS_SELECTOR, "div[class*='Tile-module_tile']")
                        
                        # Check if tiles have final states
                        final_states_count = 0
                        for tile in tiles:
                            state = tile.get_attribute("data-state") or ""
                            if state in ["correct", "present", "absent"]:
                                final_states_count += 1
                        
                        if final_states_count >= 3:  # At least 3 tiles should have final states
                            print(f"✅ Tiles have final states after {attempt + 1} attempts")
                            break
                        else:
                            print(f"Waiting for final states... attempt {attempt + 1}/10 (found {final_states_count} final states)")
                            time.sleep(0.5)
                except:
                    time.sleep(0.5)
                    continue
            
            # Perfect selector based on the actual NYT Wordle structure
            print(f"Looking for tiles in row {row_number}...")
            
            try:
                # Find all rows, then get the specific row, then get its tiles
                rows = self.driver.find_elements(By.CSS_SELECTOR, "div[class*='Row-module_row']")
                print(f"Found {len(rows)} total rows")
                
                if row_number >= len(rows):
                    print(f"Row {row_number} doesn't exist (only {len(rows)} rows found)")
                    return None
                
                # Get the specific row
                target_row = rows[row_number]
                
                # Find all tiles in this specific row
                tiles = target_row.find_elements(By.CSS_SELECTOR, "div[class*='Tile-module_tile']")
                print(f"Found {len(tiles)} tiles in row {row_number}")
                
                if len(tiles) != 5:
                    print(f"Expected 5 tiles, found {len(tiles)}")
                    return None
                
                print(f"✅ Successfully found 5 tiles in row {row_number}")
                
            except Exception as e:
                print(f"Error finding tiles: {e}")
                return None
            
            # Now read the colors from each tile
            result = ""
            for i, tile in enumerate(tiles):
                # Debug: Show all attributes for this tile
                classes = tile.get_attribute("class") or ""
                data_state = tile.get_attribute("data-state") or ""
                aria_label = tile.get_attribute("aria-label") or ""
                bg_color = tile.value_of_css_property("background-color") or ""
                letter = tile.text or ""
                
                print(f"  Tile {i} ('{letter}'): class='{classes}', data-state='{data_state}', aria-label='{aria_label}', bg-color='{bg_color}'")
                
                # Optimized: Check data-state first (most reliable for NYT Wordle)
                if data_state == "correct":
                    tile_color = "green"
                    result += "g"
                elif data_state == "present":
                    tile_color = "yellow"
                    result += "y"
                elif data_state == "absent":
                    tile_color = "gray"
                    result += "b"
                else:
                    # Fallback to original detection if data-state is unexpected
                    if "correct" in classes.lower() or "correct" in aria_label.lower():
                        tile_color = "green"
                        result += "g"
                    elif "present" in classes.lower() or "present" in aria_label.lower():
                        tile_color = "yellow"
                        result += "y"
                    else:
                        tile_color = "gray (fallback)"
                        result += "b"
                
                print(f"    → Determined color: {tile_color}")
            
            print(f"Row {row_number} final result: {result}")
            return result
            
        except Exception as e:
            print(f"Error getting result colors for row {row_number}: {e}")
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
                # Go solver updated its word list, now wait for the next word
                remaining_words = response[8:]
                print(f"Go solver updated word list: {remaining_words} words remaining")
                
                # Read the next response (should be the next WORD:)
                next_response = self.go_process.stdout.readline().strip()
                print(f"Next Go response: '{next_response}'")
                
                if next_response.startswith("WORD:"):
                    return next_response[5:]
                elif next_response.startswith("SOLVED:"):
                    parts = next_response.split(":")
                    word = parts[1]
                    attempts = parts[2]
                    print(f"🎉 Puzzle solved with '{word}' in {attempts} attempts!")
                    return None
                else:
                    print(f"Expected WORD after UPDATED, got: {next_response}")
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



    def play_game(self):
        """Main game loop"""
        if not self.setup_driver():
            return False
        
        if not self.start_go_solver():
            self.driver.quit()
            return False
        
        if not self.open_wordle():
            self.cleanup()
            return False
        
        try:
            # Get the first word from Go solver
            print("Waiting for first word from Go solver...")
            first_response = self.go_process.stdout.readline().strip()
            print(f"Go solver response: '{first_response}'")
            
            if not first_response.startswith("WORD:"):
                print(f"Expected first word, got: {first_response}")
                # Check if Go process has any errors
                if self.go_process.poll() is not None:
                    stderr_output = self.go_process.stderr.read()
                    print(f"Go process stderr: {stderr_output}")
                return False
            
            current_word = first_response[5:]
            print(f"Starting game with word: '{current_word}'")
            row = 0
            
            while current_word and row < 6:
                print(f"\nAttempt {row + 1}: Trying word '{current_word}'")
                
                # Type the word in the browser
                if not self.type_word(current_word):
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
            
            # Keep browser open for a moment to see the result
            time.sleep(1.5)
            return True
            
        except Exception as e:
            print(f"Error during game play: {e}")
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
    print("Starting Wordle Web Player...")
    print("Make sure you have:")
    print("1. Chrome browser installed")
    print("2. ChromeDriver installed and in PATH")
    print("3. Go solver compiled: go build -o wordle_solver main.go nyt_solver.go")
    print("4. selenium installed: pip install selenium")
    print()
    
    player = WordleWebPlayer()
    
    # You can change the URL here for different Wordle sites
    # Examples:
    # player.wordle_url = "https://wordlegame.org/"  # Alternative Wordle site
    # player.wordle_url = "https://www.powerlanguage.co.uk/wordle/"  # Original Wordle
    
    success = player.play_game()
    
    if success:
        print("\nGame completed successfully!")
    else:
        print("\nGame failed to complete.")

if __name__ == "__main__":
    main()