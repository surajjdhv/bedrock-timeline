#!/bin/bash

set -e

PROJECT_DIR="/home/suraj/bedrock-timeline"
INSTALL_DIR="/opt/bedrock-timeline"
SERVICE_NAME="bedrock-timeline"

echo "Building..."
cd "$PROJECT_DIR"
make build

echo "Deploying..."
sudo systemctl stop "$SERVICE_NAME"
sudo cp "$PROJECT_DIR/build/bedrock-timeline" "$INSTALL_DIR/"
sudo chown bedrock:bedrock "$INSTALL_DIR/bedrock-timeline"

# Copy .env file if it exists in project but not in install dir
if [ -f "$PROJECT_DIR/.env" ] && [ ! -f "$INSTALL_DIR/.env" ]; then
	echo "Copying .env file..."
	sudo cp "$PROJECT_DIR/.env" "$INSTALL_DIR/.env"
	sudo chown bedrock:bedrock "$INSTALL_DIR/.env"
fi

echo "Starting service..."
sudo systemctl start "$SERVICE_NAME"

echo "Checking status..."
sudo systemctl status "$SERVICE_NAME" --no-pager

echo ""
echo "Deploy complete! Service is running on port 3000."