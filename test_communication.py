#!/usr/bin/env python3
"""
Test script to verify Python-Go communication works
"""

import subprocess
import time

def test_communication():
    print("Testing Python-Go communication...")
    
    # Start the Go solver
    try:
        go_process = subprocess.Popen(
            ['./wordle_solver', '--auto'],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1
        )
        
        print("✅ Go solver started")
        
        # Read the first word suggestion
        response = go_process.stdout.readline().strip()
        print(f"📤 Go says: {response}")
        
        if response.startswith("WORD:"):
            word = response[5:]
            print(f"🎯 Suggested word: {word}")
            
            # Simulate some validation results
            test_validations = ["bygby", "ggggg"]  # Example: mostly wrong, then solved
            
            for i, validation in enumerate(test_validations):
                print(f"\n🔄 Sending validation #{i+1}: {validation}")
                go_process.stdin.write(validation + "\n")
                go_process.stdin.flush()
                
                if validation == "ggggg":
                    # For solved case, Go will output SOLVED and exit
                    response = go_process.stdout.readline().strip()
                    print(f"📥 Go responds: {response}")
                    if response.startswith("SOLVED:"):
                        print("🎉 Communication test successful!")
                        break
                else:
                    # For other cases, expect UPDATED message then next WORD
                    response = go_process.stdout.readline().strip()
                    print(f"📥 Go responds: {response}")
                    
                    if response.startswith("UPDATED:"):
                        # Read the next word suggestion
                        next_response = go_process.stdout.readline().strip()
                        print(f"📥 Go continues: {next_response}")
                        if next_response.startswith("WORD:"):
                            next_word = next_response[5:]
                            print(f"🎯 Next suggested word: {next_word}")
                    elif response.startswith("ERROR:"):
                        print(f"❌ Error: {response}")
                        break
        
        go_process.terminate()
        go_process.wait()
        print("✅ Test completed")
        
    except Exception as e:
        print(f"❌ Test failed: {e}")

if __name__ == "__main__":
    test_communication()