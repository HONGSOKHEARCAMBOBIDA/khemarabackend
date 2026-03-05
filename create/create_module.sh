#!/bin/bash

# Usage: ./create_module.sh Auth

MODULE_NAME=$1

if [ -z "$MODULE_NAME" ]; then
    echo "Please provide module name, e.g. ./create_module.sh Auth"
    exit 1
fi

# Lowercase for filenames
MODULE_LC=$(echo "$MODULE_NAME" | tr '[:upper:]' '[:lower:]')

# Folders to create
FOLDERS=("service" "controller" "model" "request" "response")

for folder in "${FOLDERS[@]}"; do
    mkdir -p "$folder"
done

# Create files
touch "service/${MODULE_LC}_service.go"
touch "controller/${MODULE_LC}_controller.go"
touch "model/${MODULE_LC}.go"
touch "request/${MODULE_LC}_request.go"
touch "response/${MODULE_LC}_response.go"

echo "✅ Module $MODULE_NAME Create Successfully!"