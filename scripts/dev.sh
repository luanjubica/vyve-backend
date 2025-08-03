#!/bin/bash

set -e

echo "🔥 Starting Vyve Backend in development mode with hot reload..."

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "📦 Installing air for hot reload..."
    go install github.com/air-verse/air@latest
fi

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Running setup first..."
    ./scripts/setup.sh
fi

echo "🚀 Starting server with hot reload..."
echo "   Server will restart automatically when you save files"
echo "   Press Ctrl+C to stop"
echo ""

# Start with air for hot reload
air -c .air.toml
