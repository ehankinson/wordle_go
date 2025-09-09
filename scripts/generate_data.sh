#!/bin/bash

# Run the data generation script for ML training

echo "Starting data generation for ML training..."

# Always compile fresh to ensure latest changes are included
echo "Compiling generate package..."
cd src/go

# Create bin directory if it doesn't exist
mkdir -p ../../bin

# Compile the generate package
echo "  Building data generator..."
go build -o ../../bin/generate_data ./cmd/generate

# Check if compilation was successful
if [ $? -ne 0 ]; then
    echo "❌ Failed to compile generate package"
    echo "Make sure you have Go installed and the source files are present"
    exit 1
fi

echo "✅ Generate package compiled successfully"

# Return to project root and run the binary
cd ../..
echo "Running data generation..."
./bin/generate_data

echo "Done!"