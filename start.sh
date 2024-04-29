#!/bin/bash

# Number of nodes
number_of_clients=10
num_iterations=$((number_of_clients+1))


# Loop through the range of iterations
for ((i=0; i<num_iterations; i++)); do
    # Format the command with the current index
    base_command="cd Desktop/fl-app-channel && ./app-channel demo --config config/peer_%d.yaml --log-level trace --log-file logs/peer_%d.log"
    current_command=$(printf "$base_command" "$i" "$i")

    # Open a new tab in the terminal and execute the command
    osascript -e "tell application \"Terminal\" to do script \"$current_command\""
done

