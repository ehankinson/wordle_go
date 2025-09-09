#!/bin/bash

# Initialize conda for this shell session
eval "$(conda shell.bash hook)"

# Activate the wordle conda environment
echo "Activating 'wordle' conda environment..."
conda activate wordle

# Check if conda activation was successful
if [ $? -ne 0 ]; then
    echo "❌ Failed to activate conda environment 'wordle'"
    echo "🚀 Running environment setup script..."
    
    # Get the directory where this script is located
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    
    # Run the create_conda_env.sh script
    if [ -f "$SCRIPT_DIR/create_conda_env.sh" ]; then
        echo "Running $SCRIPT_DIR/create_conda_env.sh..."
        bash "$SCRIPT_DIR/create_conda_env.sh"
        
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
        echo "❌ create_conda_env.sh script not found at $SCRIPT_DIR/create_conda_env.sh"
        echo "Please make sure you have created the environment with: conda create -n wordle python=3.9"
        exit 1
    fi
fi

echo "✅ Conda environment activated"

# Compile the Go programs
echo "Compiling Go programs..."

# Change to the Go source directory
cd src/go

# Build the main program
echo "  Building main program..."
go build -o ../../bin/wordle_main .

# Build the NYT solver
echo "  Building NYT solver..."
go build -o ../../bin/nyt_solver ./cmd/nyt

# Build the data generator
echo "  Building data generator..."
go build -o ../../bin/generate_data ./cmd/generate

# Return to original directory
cd ../..

# Check if compilation was successful
if [ $? -ne 0 ]; then
    echo "Failed to compile Go programs"
    echo "Make sure you have Go installed and the source files are present"
    exit 1
fi

echo "✅ Go programs compiled successfully"