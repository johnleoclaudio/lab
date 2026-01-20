# Go Starter Project

A production-ready Go backend starter template with clean architecture, following best practices for API development.

## Features

- Clean layered architecture (Handler → Service → Repository)
- PostgreSQL with sqlc for type-safe database queries
- JWT authentication ready
- Structured logging with slog
- Docker and docker-compose setup
- Database migrations
- Comprehensive middleware (logging, recovery, CORS, request ID)
- Environment-based configuration
- Makefile for common tasks
- Health check endpoint

## Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP request handlers
│   │   ├── middleware/  # HTTP middleware
│   │   └── router.go    # Route definitions
│   ├── service/         # Business logic layer
│   ├── repository/      # Data access layer
│   ├── models/          # Domain models
│   └── config/          # Configuration
├── migrations/          # Database migrations
├── queries/             # sqlc query definitions
├── docs/                # Documentation
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── sqlc.yaml
```

## Prerequisites

- Go 1.22 or higher
- Docker and Docker Compose
- PostgreSQL (or use Docker)
- Make

## Quick Start

### 1. Clone and Setup

```bash
# Update go.mod with your module path
# Edit go.mod and replace github.com/yourusername/go-starter with your path

# Install development tools and setup environment
make setup
```

This will:
- Install required tools (migrate, sqlc, golangci-lint)
- Copy .env.example to .env
- Start Docker containers
- Run database migrations
- Generate sqlc code

### 2. Manual Setup (Alternative)

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your settings
# Update DATABASE_URL, JWT_SECRET, etc.

# Install tools
make install-tools

# Start database
make docker-up

# Run migrations
make migrate-up

# Generate sqlc code
make sqlc-gen

# Run the application
make run
```

## Development

### Running the Application

```bash
# Run locally
make run

# Build binary
make build

# Run with Docker
make docker-up
```

### Database Operations

```bash
# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Generate sqlc code (after modifying queries/)
make sqlc-gen
```

### Testing

```bash
# Run all tests with coverage
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint
```

## API Endpoints

### Health Check
```
GET /health
```

Response:
```json
{
  "status": "healthy"
}
```

## Configuration

Configuration is managed through environment variables. Copy `.env.example` to `.env` and update values:

```bash
# Server
SERVER_ADDRESS=:8080
SERVER_ENV=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## Architecture

This project follows a clean layered architecture:

### Handler Layer
- Parses HTTP requests
- Validates input
- Calls service layer
- Formats HTTP responses
- **Never** contains business logic or accesses database directly

### Service Layer
- Implements business logic
- Orchestrates repository calls
- Validates business rules
- Handles transactions
- **Never** handles HTTP specifics

### Repository Layer
- Data access abstraction
- Uses sqlc-generated code (NEVER manual SQL)
- Handles database operations
- **Never** contains business logic

## Adding Features

### Add a New Endpoint

Refer to `docs/QUICK_START.md` for detailed examples. Basic steps:

1. Create migration in `migrations/`
2. Add queries in `queries/`
3. Run `make sqlc-gen`
4. Create repository in `internal/repository/`
5. Create service in `internal/service/`
6. Create handler in `internal/api/handlers/`
7. Register route in `internal/api/router.go`
8. Write tests

### Example: Add Posts Resource

See `docs/QUICK_START.md` section "Add Complete CRUD Resource" for AI prompts to generate complete features.

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[README.md](docs/README.md)** - Documentation overview
- **[QUICK_START.md](docs/QUICK_START.md)** - Ready-to-use prompts for AI agents
- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** - System design and structure
- **[AGENTS.md](docs/AGENTS.md)** - AI agent roles and workflows
- **[SKILLS.md](docs/SKILLS.md)** - Reusable capabilities catalog
- **[CODING_STANDARDS.md](docs/CODING_STANDARDS.md)** - Go conventions
- **[DATABASE.md](docs/DATABASE.md)** - sqlc and migration guidelines
- **[API_STANDARDS.md](docs/API_STANDARDS.md)** - REST and JSON:API spec
- **[TESTING.md](docs/TESTING.md)** - Testing strategies
- **[SECURITY.md](docs/SECURITY.md)** - Security best practices

## Docker Support

### Development with Docker

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

### Build Docker Image

```bash
docker build -t go-starter .
```

## Contributing

1. Follow the coding standards in `docs/CODING_STANDARDS.md`
2. Write tests for new features
3. Ensure `make lint` passes
4. Update documentation as needed

## License

MIT License

## Support

For questions or issues, please refer to the documentation in the `docs/` directory or open an issue.
