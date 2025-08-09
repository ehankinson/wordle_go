#!/usr/bin/env python3
"""
NYT Wordle Web Player - Refactored to use base driver
"""

import time

from driver import WordleDriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC


class NYTWordleDriver(WordleDriver):
    
    def __init__(self):
        super().__init__("https://www.nytimes.com/games/wordle/index.html", "nty")



    def handle_initial_setup(self):
        """Handle NYT Wordle Play button and any popups"""
        try:
            # Just click the Play button directly
            print("Clicking Play button...")
            play_button = WebDriverWait(self.driver, 5).until(
                EC.element_to_be_clickable((By.CSS_SELECTOR, "button[data-testid='Play']"))
            )
            play_button.click()
            print("✅ Clicked Play button!")
            time.sleep(0.5)
            
            # Close any popups that appear after clicking Play
            try:
                close_icon = WebDriverWait(self.driver, 0.5).until(
                    EC.element_to_be_clickable((By.CSS_SELECTOR, "[data-testid='icon-close']"))
                )
                close_icon.click()
                print("Closed popup")
            except:
                pass  # No popup to close
            
            print("✅ NYT Wordle ready to play!")
            time.sleep(0.5)
            return True
            
        except Exception as e:
            print(f"No Play button found or already in game: {e}")
            # Game might already be ready if no Play button
            return True



    def click_new_game(self):
        """NYT Wordle doesn't have a new game button (one game per day)"""
        print("NYT Wordle only allows one game per day")
        return False



    def nyt(self):
        """Play a single NYT Wordle game"""
        if not self.game_setup():
            return False
        
        if not self.handle_initial_setup():
            return False
        
        if not self.play_game():
            return False
        
        time.sleep(3)
        self.cleanup()
        return True



def main():
    player = NYTWordleDriver()
    if player.nyt():
        print("\n✅ Game completed successfully!")
    else:
        print("\n❌ Failed to complete game")



if __name__ == "__main__":
    main()