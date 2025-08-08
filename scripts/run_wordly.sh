#!/bin/bash

# Run the setup script to activate conda and compile Go
source scripts/setup_env.sh

# If setup failed, exit
if [ $? -ne 0 ]; then
    echo "Setup failed, exiting..."
    exit 1
fi

# Run the Wordly web player
echo "Starting Wordly web player..."
python3 src/python/wordly_web_player.py

echo "Done!"