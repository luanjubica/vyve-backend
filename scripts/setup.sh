#!/bin/bash

set -e

echo "🚀 Setting up Vyve Backend development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    echo "   Visit: https://golang.org/doc/install"
    exit 1
fi

echo "✅ Go found: $(go version)"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  IMPORTANT: Update the .env file with your actual Supabase values!"
    echo "   - Get your Supabase URL and keys from: https://supabase.com/dashboard"
fi

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod download
go mod tidy

# Install development tools
echo "🛠️  Installing development tools..."
if ! command -v air &> /dev/null; then
    echo "   Installing air for hot reload..."
    go install github.com/air-verse/air@latest
fi

# Create necessary directories
echo "📁 Creating runtime directories..."
mkdir -p logs tmp uploads

# Set script permissions
chmod +x scripts/*.sh

echo ""
echo "✅ Setup complete!"
echo ""
echo "🎯 Next steps:"
echo "=============="
echo "1. Update .env with your Supabase credentials:"
echo "   - SUPABASE_URL=https://your-project.supabase.co"
echo "   - SUPABASE_ANON_KEY=your-anon-key"
echo "   - SUPABASE_JWT_SECRET=your-jwt-secret"
echo "   - DATABASE_URL=your-supabase-db-url"
echo ""
echo "2. Run the application:"
echo "   make run          # Run once"
echo "   make dev          # Run with hot reload"
echo "   make docker-run   # Run with Docker"
echo ""
echo "3. Test the health endpoint:"
echo "   curl http://localhost:8080/health"
echo ""
echo "📚 More commands:"
echo "   make build        # Build binary"
echo "   make test         # Run tests"
echo "   make format       # Format code"
echo ""
