#!/usr/bin/env python3
"""
Wordly Web Player - Refactored to use base driver
"""

import json
import time
from selenium.webdriver.common.by import By
from driver import WordleDriver


GAMES_TO_PLAY = 10


class WordlyDriver(WordleDriver):
    """Driver for Wordly.org"""
    
    def __init__(self):
        super().__init__("https://wordly.org/", "wordly")



    def click_new_game(self):
        """Click the New game button"""
        try:
            print("Looking for New game button...")
            new_game_button = self.driver.find_element(By.XPATH, "//button[text()='New game']")
            if new_game_button.is_displayed():
                new_game_button.click()
                print("✅ Clicked New game button")
                time.sleep(0.25)
                self.reset_solver()
                return True
            return False
        except Exception as e:
            print(f"Error clicking New game: {e}")
            return False



    def play_games(self, num_games=GAMES_TO_PLAY):
        """Play multiple games"""
        # Load results
        with open('data/wordly_results.json', 'r') as f:
            results = json.load(f)
        
        if not self.game_setup():
            return False
        
        try:
            for game_num in range(num_games):
                results["games_played"] += 1
                print(f"\n{'='*50}")
                print(f"GAME {game_num + 1} of {num_games}")
                print('='*50)
                
                if self.play_game(results):
                    print(f"✅ Game {game_num + 1} completed!")
                    
                    if game_num < num_games - 1:
                        if not self.click_new_game():
                            break
                else:
                    print(f"❌ Game {game_num + 1} failed")
                    break
            
            # Save results
            with open('data/wordly_results.json', 'w') as f:
                json.dump(results, f, indent=4)
            
            return True
            
        except Exception as e:
            print(f"Error: {e}")
            return False
        finally:
            self.cleanup()


def main():
    player = WordlyDriver()
    if player.play_games():
        print("\n✅ All games completed!")
    else:
        print("\n❌ Failed to complete games")


if __name__ == "__main__":
    main()