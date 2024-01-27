#!/bin/bash
set -e

echo "Are you sure the Llama2 Model is running on 127.0.0.1:7878? (y/n)"
read answer

if [[ $answer = "n" ]]; then
  exit 1
fi

# Function to clean up the background Rust process
cleanup() {
    echo "Cleaning up..."
    kill $RUST_PID
    exit 0
}

# Trap to call cleanup when the script exits
trap cleanup EXIT

# Start Rust Backend
echo "Starting Rust Backend"
cd ../target/debug/
./milkyteadrop-local &
RUST_PID=$!

# Start Go Discord Bot
echo "Starting Discord Bot"
cd ../../src
go run main.go
cd ../

# Start Python Rest API
# echo "Now starting Python Rest API"
# cd ../models/stable_diffusion
# flask --app main run &
# The cleanup function will be called automatically when the script exits