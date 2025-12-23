#!/usr/bin/env bash

# Exit if no argument is provided
if [ -z "$1" ]; then
  echo "Usage: ./setenv.sh [config-name]"
  exit 1
fi

ENV=$1
SOURCE_FILE="config/config.${ENV}.env"
TARGET_FILE=".env"

# Check if the source file exists
if [ ! -f "$SOURCE_FILE" ]; then
  echo "Error: $SOURCE_FILE does not exist."
  exit 1
fi

# Copy the file
cp "$SOURCE_FILE" "$TARGET_FILE"

echo "Environment set to $SOURCE_FILE"
