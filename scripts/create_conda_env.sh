#!/bin/bash

# Script to create conda environment 'wordle' if it doesn't exist
# and install all required packages

set -e  # Exit on any error

echo "🔍 Checking for conda installation..."

# Check if conda is available
if ! command -v conda &> /dev/null; then
    echo "❌ Conda is not installed or not in PATH"
    echo "Please install Anaconda or Miniconda first"
    exit 1
fi

echo "✅ Conda found"

# Initialize conda for this shell session
eval "$(conda shell.bash hook)"

# Check if wordle environment already exists
echo "🔍 Checking if 'wordle' environment exists..."

if conda env list | grep -q "^wordle\s"; then
    echo "✅ 'wordle' environment already exists"
    echo "To recreate it, run: conda env remove -n wordle"
    
    # Activate existing environment
    echo "🚀 Activating existing 'wordle' environment..."
    conda activate wordle
    
    # Check if we should update packages
    echo "📦 Would you like to update packages? (y/n)"
    read -p "Update packages from requirements.txt? [y/N]: " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "📦 Updating packages from requirements.txt..."
        if [ -f "requirements.txt" ]; then
            pip install -r requirements.txt
        elif [ -f "requirments.txt" ]; then
            echo "⚠️  Found 'requirments.txt' (misspelled) - installing essential packages manually..."
            pip install selenium pygame requests python-dotenv tqdm matplotlib PyQt5 torch torchvision torchaudio
        else
            echo "⚠️  requirements.txt not found, skipping package installation"
        fi
    fi
else
    echo "📦 Creating new 'wordle' conda environment..."
    
    # Create the environment with Python 3.9 (as mentioned in setup_env.sh)
    conda create -n wordle python=3.9 -y
    
    # Activate the new environment
    echo "🚀 Activating 'wordle' environment..."
    conda activate wordle
    
    # Install packages from requirements.txt if it exists
    echo "📦 Installing packages from requirements.txt..."
    if [ -f "requirements.txt" ]; then
        echo "Installing Python packages..."
        pip install -r requirements.txt
        echo "✅ All packages installed successfully"
    elif [ -f "requirments.txt" ]; then
        echo "⚠️  Found 'requirments.txt' (misspelled) - installing essential packages manually..."
        pip install selenium pygame requests python-dotenv tqdm matplotlib PyQt5 torch torchvision torchaudio
        echo "✅ Essential packages installed successfully"
    else
        echo "⚠️  requirements.txt not found in project root"
        echo "You may need to install packages manually"
    fi
fi

# Verify the environment is active
if [[ "$CONDA_DEFAULT_ENV" == "wordle" ]]; then
    echo "✅ Successfully activated 'wordle' environment"
    echo "🐍 Python version: $(python --version)"
    echo "📍 Environment location: $CONDA_PREFIX"
else
    echo "❌ Failed to activate 'wordle' environment"
    exit 1
fi

echo ""
echo "🎉 Setup complete! The 'wordle' conda environment is ready to use."
echo ""
echo "To activate this environment in the future, run:"
echo "  conda activate wordle"
echo ""
echo "To deactivate the environment, run:"
echo "  conda deactivate"
