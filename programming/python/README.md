# Python Backend Development - One Day Crash Course

**Goal**: Learn the 20% of Python that covers 80% of backend development needs in 2025.

## Why Python for Backend?

- **FastAPI** dominates modern API development (replacing Flask/Django REST in many cases)
- **Async/await** is now standard (not optional)
- Type hints with **Pydantic V2** for data validation
- Best-in-class tooling: **uv** for dependency management, **Ruff** for linting
- Strong ecosystem for AI/ML integration

## Learning Path (6-8 hours)

1. **[01-python-essentials.md](./01-python-essentials.md)** (1.5 hours)
   - Core syntax, data structures, functions
   - List/dict comprehensions, unpacking, context managers
   - Type hints that actually matter

2. **[02-async-concurrency.md](./02-async-concurrency.md)** (1.5 hours)
   - async/await fundamentals
   - When to use async vs threads vs multiprocessing
   - Common patterns in backend services

3. **[03-fastapi-modern-backend.md](./03-fastapi-modern-backend.md)** (2 hours)
   - FastAPI with Pydantic V2
   - Dependency injection, middleware
   - Authentication, error handling
   - OpenAPI documentation

4. **[04-database-sqlalchemy.md](./04-database-sqlalchemy.md)** (1.5 hours)
   - SQLAlchemy 2.0 with async
   - Migrations with Alembic
   - Connection pooling, query optimization

5. **[05-testing-tooling.md](./05-testing-tooling.md)** (1 hour)
   - pytest patterns for APIs
   - uv for dependency management
   - Ruff for linting/formatting
   - Docker for deployment

6. **[06-production-ready.md](./06-production-ready.md)** (1 hour)
   - Project structure
   - Configuration management
   - Logging, monitoring
   - Deployment patterns

## Quick Start

```bash
# Install uv (modern pip replacement, 10-100x faster)
curl -LsSf https://astral.sh/uv/install.sh | sh

# Create a new project
uv init my-backend-api
cd my-backend-api

# Add dependencies
uv add fastapi uvicorn[standard] sqlalchemy[asyncio] pydantic pydantic-settings

# Add dev dependencies
uv add --dev pytest pytest-asyncio ruff
```

## 2025 Industry Standards

- **Package Manager**: uv (not pip/poetry)
- **Web Framework**: FastAPI (not Flask/Django for APIs)
- **Linter/Formatter**: Ruff (not flake8/black/isort)
- **ORM**: SQLAlchemy 2.0 async (or raw SQL with asyncpg)
- **Validation**: Pydantic V2 (10x faster than V1)
- **Testing**: pytest with pytest-asyncio
- **Type Checking**: pyright (faster) or mypy
- **Python Version**: 3.12+ (performance improvements)

## Next Steps After This Guide

- Build a full CRUD API with auth
- Learn Redis for caching/queues
- Study message queues (RabbitMQ, Kafka)
- Explore observability (OpenTelemetry)
- Try production deployment (Docker + Kubernetes)

Let's start with the fundamentals!
