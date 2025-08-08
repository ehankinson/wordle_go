#!/bin/bash

# Initialize conda for this shell session
eval "$(conda shell.bash hook)"

# Activate the wordle conda environment
echo "Activating 'wordle' conda environment..."
conda activate wordle

# Check if conda activation was successful
if [ $? -ne 0 ]; then
    echo "Failed to activate conda environment 'wordle'"
    echo "Make sure you have created the environment with: conda create -n wordle python=3.9"
    exit 1
fi

echo "✅ Conda environment activated"

# Compile the Go programs
echo "Compiling Go solver..."
go build -o bin/wordle_solver src/go/main.go src/go/nyt_solver.go

# Check if compilation was successful
if [ $? -ne 0 ]; then
    echo "Failed to compile Go programs"
    echo "Make sure you have Go installed and the source files are present"
    exit 1
fi

echo "✅ Go solver compiled successfully"