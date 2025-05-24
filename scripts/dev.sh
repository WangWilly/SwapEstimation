#/bin/bash


export DB_SEED=true
export DB_IS_DEV=false
export LOG_FORMAT=console
# This script is used to run the dev server for the project.
# It sets up the environment and starts the server.
# Usage: ./scripts/dev.sh
# Make sure to run this script from the root of the project.
# Check if the script is being run from the root of the project
go run cmd/main.go
