#!/bin/bash

# Run the NYT Wordle solver

echo "Starting NYT Wordle Solver..."

# Always compile fresh to ensure latest changes are included
echo "Compiling NYT solver package..."
cd src/go

# Create bin directory if it doesn't exist
mkdir -p ../../bin

# Compile the NYT solver package
echo "  Building NYT solver..."
go build -o ../../bin/nyt_solver ./cmd/nyt

# Check if compilation was successful
if [ $? -ne 0 ]; then
    echo "❌ Failed to compile NYT solver package"
    echo "Make sure you have Go installed and the source files are present"
    exit 1
fi

echo "✅ NYT solver package compiled successfully"

# Return to project root and run the binary
cd ../..
echo "Running NYT Wordle Solver..."
./bin/nyt_solver

echo "Done!"