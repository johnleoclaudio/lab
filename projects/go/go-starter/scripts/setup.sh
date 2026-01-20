#!/bin/bash
set -e

echo "🚀 Go Starter Project Setup"
echo "============================"
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ .env file created. Please update it with your settings."
    echo ""
fi

# Check if required tools are installed
echo "🔍 Checking required tools..."

check_tool() {
    if command -v $1 &> /dev/null; then
        echo "  ✅ $1 is installed"
        return 0
    else
        echo "  ❌ $1 is NOT installed"
        return 1
    fi
}

MISSING_TOOLS=0

if ! check_tool go; then
    echo "     Install from: https://golang.org/dl/"
    MISSING_TOOLS=1
fi

if ! check_tool docker; then
    echo "     Install from: https://docs.docker.com/get-docker/"
    MISSING_TOOLS=1
fi

if ! check_tool migrate; then
    echo "     Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"
    MISSING_TOOLS=1
fi

if ! check_tool sqlc; then
    echo "     Install with: go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"
    MISSING_TOOLS=1
fi

echo ""

if [ $MISSING_TOOLS -eq 1 ]; then
    echo "⚠️  Some tools are missing. Install them before proceeding."
    echo ""
    read -p "Do you want to install Go tools (migrate, sqlc)? (y/n) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "📦 Installing Go tools..."
        go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
        go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
        echo "✅ Go tools installed"
    fi
fi

echo ""
echo "📦 Downloading Go dependencies..."
go mod download
go mod tidy
echo "✅ Dependencies ready"

echo ""
echo "🐳 Starting Docker containers..."
docker-compose up -d
echo "✅ Docker containers started"

echo ""
echo "⏳ Waiting for database to be ready..."
sleep 5

echo ""
echo "🗄️  Running database migrations..."
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/go_starter?sslmode=disable"
migrate -path migrations -database "$DATABASE_URL" up
echo "✅ Migrations complete"

echo ""
echo "⚙️  Generating sqlc code..."
sqlc generate
echo "✅ sqlc code generated"

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "  1. Update .env with your configuration"
echo "  2. Run 'make run' to start the server"
echo "  3. Visit http://localhost:8080/health to verify"
echo ""
echo "Available commands:"
echo "  make help       - Show all available commands"
echo "  make run        - Run the application"
echo "  make test       - Run tests"
echo "  make docker-up  - Start Docker services"
echo ""
