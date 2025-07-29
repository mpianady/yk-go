#!/bin/bash

# Stop and remove all Docker containers defined in compose.yml
if ! docker-compose down; then
    echo "Error: Failed to stop Docker containers"
    exit 1
fi

echo "Docker containers stopped successfully"
