# Wordle Go Solver

A high-performance Wordle solver written in Go that uses information theory and letter frequency analysis to find optimal word guesses. The solver integrates with web browsers through Python/Selenium to automatically play Wordle games.

## Current Features

- Algorithmic solver using entropy-based word selection
- Automated web playing for NYT Wordle and Wordly.org
- Interactive command-line solver for manual play
- ~95% win rate with average solve in 3-4 guesses

## Future Plans

- Reinforcement Learning (RL) model to improve solving strategy
- AI agent that learns optimal guess patterns through self-play
- Performance comparison between algorithmic and ML approaches

## Quick Start

```bash
# Play NYT Wordle automatically in browser
./scripts/run_wordle.sh

# Play Wordly.org automatically in browser
./scripts/run_wordly.sh

# Interactive command-line solver
./scripts/run_nyt_solver.sh
```

## Requirements

- Go 1.19+
- Python 3.9+ with Selenium
- Chrome/Firefox browser with appropriate driver
