#!/bin/bash

# Run the Wordly web player

echo "Starting Wordly web player..."

# Initialize conda for this shell session
eval "$(conda shell.bash hook)"

# Activate the wordle conda environment
echo "Activating 'wordle' conda environment..."
conda activate wordle

# Check if conda activation was successful
if [ $? -ne 0 ]; then
    echo "❌ Failed to activate conda environment 'wordle'"
    echo "🚀 Running environment setup script..."
    
    # Run the create_conda_env.sh script
    if [ -f "scripts/create_conda_env.sh" ]; then
        echo "Running scripts/create_conda_env.sh..."
        bash scripts/create_conda_env.sh
        
        # Check if the environment creation was successful
        if [ $? -ne 0 ]; then
            echo "❌ Failed to create conda environment"
            exit 1
        fi
        
        # Try to activate the environment again
        echo "🔄 Attempting to activate 'wordle' environment again..."
        conda activate wordle
        
        if [ $? -ne 0 ]; then
            echo "❌ Still failed to activate conda environment 'wordle' after creation"
            exit 1
        fi
    else
        echo "❌ create_conda_env.sh script not found"
        echo "Please create the environment with: conda create -n wordle python=3.9"
        exit 1
    fi
fi

echo "✅ Conda environment activated"

# Run the Wordly web player (no Go compilation needed)
python3 src/python/wordly.py

echo "Done!"