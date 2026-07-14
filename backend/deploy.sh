#!/bin/bash

set -e

echo "Starting Kshanik backend deployment..."

sudo apt-get update
sudo apt-get install -y ca-certificates curl

if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker "$USER"
fi

if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
    echo "Edit .env with production values, then run ./deploy.sh again."
    exit 0
fi

set -a
. ./.env
set +a

if [ ! -f Caddyfile ]; then
    echo "Error: Caddyfile not found. Run this script from the backend directory."
    exit 1
fi

echo "Registering systemd service..."
sudo cp kshanik-stack.service /etc/systemd/system/kshanik-stack.service
sudo systemctl daemon-reload
sudo systemctl enable kshanik-stack
sudo systemctl restart kshanik-stack

echo "Deployment complete."
echo "-------------------------------------------------------"
echo "1. Run 'sudo journalctl -u kshanik-stack -f' to see logs."
echo "2. Run 'sudo docker compose ps' to check container status."
echo "3. Run 'curl https://$API_DOMAIN/health' to verify HTTPS."
echo "-------------------------------------------------------"
