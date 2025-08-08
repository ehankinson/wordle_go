# Wordle Go Solver

A Go-based Wordle solver with Python web automation that plays Wordle games automatically by combining algorithmic solving with browser interaction.


## Quick Start

```bash
# For NYT Wordle
./scripts/run_wordle.sh

# For Wordly.org
./scripts/run_wordly.sh
```

## How It Works

The Go solver calculates optimal guesses using information theory, while Python handles web automation via Selenium. They communicate through stdin/stdout to play games in real-time
