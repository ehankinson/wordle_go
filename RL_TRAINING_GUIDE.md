# Wordle RL Training Guide

## Overview

This guide documents the approach for training a Reinforcement Learning model to play Wordle, transitioning from the current heuristic-based Go solver to an AI agent that learns optimal strategies through experience.

## Current System vs RL Approach

### Current Heuristic Solver
- Uses letter frequency analysis at each position
- Deterministic: always picks the "best" word based on probability scores
- Filters words based on feedback (green/yellow/black)
- No learning or adaptation over time

### RL Approach
- Learns patterns from thousands of game experiences
- Balances exploration (trying new strategies) and exploitation (using known good strategies)
- Can adapt strategies based on game state
- Learns when to take risks vs play safe

## Training Data Structure

### Simplified Experience Format
```json
{
  "game_id": "001",
  "target_word": "CRANE",
  "success": true,
  "attempts": 3,
  "exploration_params": {
    "temperature": 0.8,
    "top_k": 15,
    "strategy": "balanced"
  },
  "experiences": [
    {
      "attempt": 0,
      "state": {
        "valid_words": 12972,
        "final_word": ["", "", "", "", ""],
        "known_letters": [],
        "excluded_letters": []
      },
      "guess": "SLATE",
      "feedback": "bbybb",
      "reward": -0.7,
      "done": false
    },
    {
      "attempt": 1,
      "state": {
        "valid_words": 183,
        "final_word": ["", "", "", "", ""],
        "known_letters": ["A"],
        "excluded_letters": ["S", "L", "T", "E"]
      },
      "guess": "BRAND",
      "feedback": "bgggy",
      "reward": 0.5,
      "done": false
    },
    {
      "attempt": 2,
      "state": {
        "valid_words": 4,
        "final_word": ["", "R", "A", "N", ""],
        "known_letters": ["R", "A", "N"],
        "excluded_letters": ["S", "L", "T", "E", "B", "D"]
      },
      "guess": "CRANE",
      "feedback": "ggggg",
      "reward": 16.0,
      "done": true
    }
  ]
}
```

## Key Concepts

### Temperature and Top-K Parameters

#### Temperature (0.1 to 2.0+)
Controls randomness in word selection probability distribution:
- **Low (0.1-0.5)**: Nearly deterministic, strongly favors best words
- **Medium (0.7-1.0)**: Balanced exploration/exploitation
- **High (1.5-2.0)**: More uniform distribution, explores weaker options

#### Top-K (5 to 100+)
Limits the candidate pool to K best options:
- **Small K (5-10)**: Conservative, only considers top choices
- **Medium K (10-30)**: Balanced approach
- **Large K (30-100)**: More exploratory, considers mediocre options

### The Softmax Selection Function

The word selection uses temperature-scaled softmax:

```go
func get_word(ranked_words []word_score, params exploration_data) string {
    // Limit to top K words
    candidates := ranked_words
    if len(ranked_words) > params.top_k {
        candidates = ranked_words[:params.top_k]
    }
    
    // Convert scores to weights using temperature
    weights := make([]float64, len(candidates))
    sum := 0.0
    
    for i, ws := range candidates {
        // Exponential scaling with temperature
        weight := math.Exp(ws.prob / params.temp)
        weights[i] = weight
        sum += weight
    }
    
    // Weighted random selection
    r := rand.Float64() * sum
    cumulative := 0.0
    for i, w := range weights {
        cumulative += w
        if r <= cumulative {
            return candidates[i].word
        }
    }
    
    return candidates[0].word
}
```

#### How It Works:
1. **Exponential scaling**: `math.Exp(score/temp)` converts scores to weights
2. **Temperature effect**: 
   - Low temp → large differences in weights → strong preference for best
   - High temp → similar weights → more random selection
3. **Probability selection**: Random value scaled to sum of weights

Example with words SLATE(0.89), CRANE(0.86), ADIEU(0.83):
- **Temp=0.5**: Weights [5.93, 5.58, 5.25] → 35%, 33%, 32% probability
- **Temp=2.0**: Weights [1.56, 1.54, 1.51] → 34%, 33%, 33% probability

## Reward Function Design

### Recommended Scoring Function

```go
func calculate_reward(guess_num int, words_before int, words_after int,
                      feedback string, is_final bool, won bool) float64 {
    reward := -1.0  // Base penalty per guess
    
    // Information reduction bonus (0 to 3 points)
    if words_before > 0 {
        reduction_rate := float64(words_before - words_after) / float64(words_before)
        reward += reduction_rate * 3.0
    }
    
    // Constraint discovery bonuses
    greens := strings.Count(feedback, "g")
    yellows := strings.Count(feedback, "y")
    reward += float64(greens) * 0.3  // Green letters worth more
    reward += float64(yellows) * 0.1  // Yellow letters still valuable
    
    // Terminal rewards
    if is_final {
        if won {
            // Win bonus + efficiency bonus
            reward += 10.0 + float64(7-guess_num) * 2.0
        } else {
            // Failure penalty
            reward -= 5.0
        }
    }
    
    return reward
}
```

### Reward Components:
- **Step penalty (-1)**: Encourages fewer guesses
- **Information gain (0-3)**: Rewards eliminating possibilities
- **Green letters (+0.3 each)**: Most valuable information
- **Yellow letters (+0.1 each)**: Still helpful
- **Win bonus (+10)**: Strong positive signal
- **Efficiency bonus (+2 per saved guess)**: Rewards quick solutions
- **Loss penalty (-5)**: Teaches that failing is worse than 6 guesses

## Exploration Strategy

### Mixed Strategy for Training Data

```go
type exploration_data struct {
    temp     float64
    top_k    int
    strategy string
}

func get_exploration_params() exploration_data {
    r := rand.Float64()
    
    if r < 0.4 {  // 40% - Exploitation (near-optimal)
        return exploration_data{
            temp:     0.3 + rand.Float64()*0.3,  // 0.3-0.6
            top_k:    5 + rand.Intn(5),          // 5-10
            strategy: "exploit",
        }
    } else if r < 0.8 {  // 40% - Balanced
        return exploration_data{
            temp:     0.7 + rand.Float64()*0.5,  // 0.7-1.2
            top_k:    10 + rand.Intn(20),        // 10-30
            strategy: "explore",
        }
    } else {  // 20% - High exploration
        return exploration_data{
            temp:     1.2 + rand.Float64()*0.8,  // 1.2-2.0
            top_k:    30 + rand.Intn(70),        // 30-100
            strategy: "random",
        }
    }
}
```

### Why This Distribution:
- **40% near-optimal**: Shows model what good play looks like
- **40% balanced**: Explores alternative strategies
- **20% high exploration**: Learns from suboptimal paths

## Complete Training Data Generator

```go
func generate_training_data(num_games int) []GameData {
    all_games := []GameData{}
    
    for game := 0; game < num_games; game++ {
        reset_game_state()  // Reset globals
        
        word_list := get_valide_words()
        words_before := len(word_list)
        exploration_params := get_exploration_params()
        final_word := get_random_word(word_list)
        
        game_experiences := []Experience{}
        previous_guesses := []string{}
        previous_feedback := []string{}
        success := false
        
        for att := 0; att < 6; att++ {
            // Capture current state
            state := State{
                Attempt:      att,
                ValidWords:   len(word_list),
                FinalWord:    FINAL_WORD,      // Global partial solution
                KnownLetters: KNOWN_LETTERS,   // Global known letters
            }
            
            // Make guess
            letter_frequency := get_letter_frequency(word_list)
            ranked_words := ranked_words(word_list, letter_frequency)
            guessing_word := get_word(ranked_words, exploration_params)
            validation := nyt_word_validator(final_word, guessing_word)
            
            // Update game state
            update_letter_conditions(validation, guessing_word)
            word_list = filter_word_list(word_list)
            words_after := len(word_list)
            
            // Calculate reward
            is_winner := validation == "ggggg"
            is_final := is_winner || att == 5
            reward := calculate_reward(
                att+1, words_before, words_after,
                validation, is_final, is_winner,
            )
            
            // Store experience
            exp := Experience{
                State:    state,
                Guess:    guessing_word,
                Feedback: validation,
                Reward:   reward,
                Done:     is_final,
            }
            game_experiences = append(game_experiences, exp)
            
            // Update tracking
            previous_guesses = append(previous_guesses, guessing_word)
            previous_feedback = append(previous_feedback, validation)
            words_before = words_after
            
            if is_winner {
                success = true
                break
            }
        }
        
        // Store complete game
        game_data := GameData{
            GameID:     fmt.Sprintf("game_%05d", game),
            TargetWord: final_word,
            Success:    success,
            Attempts:   len(game_experiences),
            Params:     exploration_params,
            Experiences: game_experiences,
        }
        all_games = append(all_games, game_data)
    }
    
    return all_games
}
```

## Using the Training Data

The generated JSON data can be used to train various RL models:

1. **Deep Q-Network (DQN)**: Learn Q-values for state-action pairs
2. **Policy Gradient**: Directly learn probability distribution over words
3. **Actor-Critic**: Combine value estimation with policy learning

The model learns to:
- Recognize game states (constraints, remaining words)
- Value different actions (which words lead to wins)
- Balance information gathering vs. exploiting known constraints
- Adapt strategy based on attempts remaining

## Key Insights

1. **Diversity is crucial**: Mix exploitation and exploration in training data
2. **Reward shaping matters**: Balance immediate vs. long-term rewards
3. **State representation**: Use same features your solver uses (FINAL_WORD, KNOWN_LETTERS)
4. **Temperature/K control**: Fine-tune exploration during data generation
5. **Learn from failures**: Include games that fail or take 5-6 attempts

## Next Steps

1. Generate 10,000+ games with varied exploration parameters
2. Implement PyTorch model to consume this data
3. Train with different architectures (DQN, A2C, PPO)
4. Evaluate against deterministic solver
5. Fine-tune reward function based on results