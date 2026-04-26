#!/bin/bash

# Kshanik Search - Fault Tolerant Deployment Script
# To be run on the Oracle Cloud Instance

set -e

echo "🚀 Starting Fault-Tolerant Deployment..."

# 1. Update and Install Docker dependencies
sudo apt-get update
sudo apt-get install -y ca-certificates curl gnupg lsb-release

# 2. Install Docker (if not present)
if ! command -v docker &> /dev/null; then
    echo "📦 Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
fi

# 3. Create .env if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating default .env file..."
    cp .env.example .env
    echo "⚠️ Please edit the .env file with your production secrets after this script finishes."
fi

# 4. Set up Caddyfile
if [ ! -f Caddyfile ]; then
    echo "❌ Error: Caddyfile not found! Please ensure you are in the kshanik_search/backend directory."
    exit 1
fi

# 5. Register Systemd Service
echo "⚙️ Registering Systemd service..."
sudo cp kshanik-stack.service /etc/systemd/system/kshanik-stack.service
sudo systemctl daemon-reload
sudo systemctl enable kshanik-stack

# 6. Start the Stack
echo "🔄 Starting Docker containers..."
sudo docker compose up -d

echo "✅ Deployment Complete!"
echo "-------------------------------------------------------"
echo "1. Run 'sudo journalctl -u kshanik-stack -f' to see logs."
echo "2. Edit 'Caddyfile' to add your real domain."
echo "3. Run 'sudo docker compose restart caddy' after editing."
echo "-------------------------------------------------------"
