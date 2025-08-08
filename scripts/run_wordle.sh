#!/bin/bash

# Run the setup script to activate conda and compile Go
source scripts/setup_env.sh

# If setup failed, exit
if [ $? -ne 0 ]; then
    echo "Setup failed, exiting..."
    exit 1
fi

# Run the Wordle web player
echo "Starting Wordle web player..."
python3 src/python/wordle_web_player.py

echo "Done!"