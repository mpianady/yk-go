#!/bin/bash

# Make sure to set executable permissions: chmod +x bin/start.sh

# Stop any existing containers
docker-compose down

# Build and start containers in detached mode
if ! docker-compose up -d; then
    echo "Error: Failed to start Docker containers"
    exit 1
fi

echo "Docker containers started successfully"