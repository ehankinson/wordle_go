#!/bin/bash

echo "Setting up Wordle Web Player..."

# Compile the Go program
echo "Compiling Go solver..."
go build -o wordle_solver main.go nyt_solver.go

if [ $? -eq 0 ]; then
    echo "✅ Go solver compiled successfully"
else
    echo "❌ Failed to compile Go solver"
    exit 1
fi

# Check if Python has selenium
echo "Checking Python dependencies..."
python3 -c "import selenium" 2>/dev/null

if [ $? -eq 0 ]; then
    echo "✅ Selenium is installed"
else
    echo "❌ Selenium not found. Installing..."
    pip3 install selenium
fi

# Make Python script executable
chmod +x wordle_web_player.py

echo ""
echo "Setup complete! Usage:"
echo "1. Manual mode: ./wordle_solver"
echo "2. Web automation: python3 wordle_web_player.py"
echo ""
echo "Note: Make sure you have Chrome and ChromeDriver installed for web automation"