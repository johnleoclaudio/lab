# Testing, Tooling, and Modern Python Workflow

Time: 1 hour | Focus: 2025 industry-standard tools and practices

## Modern Python Tooling Stack (2025)

| Tool | Purpose | Why It's Better |
|------|---------|-----------------|
| **uv** | Package management | 10-100x faster than pip, replaces pip/poetry/virtualenv |
| **Ruff** | Linting + Formatting | 100x faster than flake8+black+isort, all-in-one |
| **pytest** | Testing | Industry standard, plugin ecosystem |
| **pyright** | Type checking | Faster than mypy, better VS Code integration |
| **docker** | Containerization | Consistent deployments |

## uv - The Modern Package Manager

```bash
# Install uv (replaces pip, poetry, virtualenv)
curl -LsSf https://astral.sh/uv/install.sh | sh

# Create new project
uv init my-backend-api
cd my-backend-api

# Add production dependencies
uv add fastapi uvicorn[standard] sqlalchemy[asyncio] pydantic pydantic-settings

# Add development dependencies
uv add --dev pytest pytest-asyncio pytest-cov ruff

# Run your app (auto-creates venv)
uv run uvicorn app:app

# Run scripts
uv run python main.py

# Lock dependencies (like package-lock.json)
uv lock

# Sync dependencies (install from lock file)
uv sync

# Update dependencies
uv lock --upgrade

# Remove package
uv remove package-name
```

### pyproject.toml (Modern Python Config)

```toml
[project]
name = "my-backend-api"
version = "0.1.0"
description = "Modern FastAPI backend"
authors = [{name = "Your Name", email = "you@example.com"}]
requires-python = ">=3.12"
dependencies = [
    "fastapi>=0.109.0",
    "uvicorn[standard]>=0.27.0",
    "sqlalchemy[asyncio]>=2.0.25",
    "asyncpg>=0.29.0",
    "pydantic>=2.5.0",
    "pydantic-settings>=2.1.0",
]

[project.optional-dependencies]
dev = [
    "pytest>=7.4.0",
    "pytest-asyncio>=0.23.0",
    "pytest-cov>=4.1.0",
    "ruff>=0.1.0",
]

[tool.uv]
dev-dependencies = [
    "pytest>=7.4.0",
    "pytest-asyncio>=0.23.0",
]

[build-system]
requires = ["hatchling"]
build-backend = "hatchling.build"
```

## Ruff - Lightning Fast Linting and Formatting

```bash
# Install (already included with uv add --dev ruff)
uv add --dev ruff

# Check code
uv run ruff check .

# Auto-fix issues
uv run ruff check --fix .

# Format code
uv run ruff format .

# Run both
uv run ruff check --fix . && uv run ruff format .
```

### Ruff Configuration

```toml
# Add to pyproject.toml
[tool.ruff]
line-length = 100
target-version = "py312"

[tool.ruff.lint]
select = [
    "E",   # pycodestyle errors
    "W",   # pycodestyle warnings
    "F",   # pyflakes
    "I",   # isort
    "B",   # flake8-bugbear
    "C4",  # flake8-comprehensions
    "UP",  # pyupgrade
]
ignore = [
    "E501",  # line too long (handled by formatter)
    "B008",  # do not perform function calls in argument defaults
]

[tool.ruff.lint.per-file-ignores]
"__init__.py" = ["F401"]  # Ignore unused imports in __init__.py

[tool.ruff.format]
quote-style = "double"
indent-style = "space"
```

## pytest - Testing Framework

### Basic Test Structure

```python
# tests/test_users.py
import pytest
from httpx import AsyncClient
from app.main import app

@pytest.mark.asyncio
async def test_create_user():
    async with AsyncClient(app=app, base_url="http://test") as client:
        response = await client.post(
            "/users",
            json={
                "email": "test@example.com",
                "name": "Test User",
                "password": "SecurePass123"
            }
        )
    
    assert response.status_code == 201
    data = response.json()
    assert data["email"] == "test@example.com"
    assert "id" in data
    assert "password" not in data  # Should not return password

@pytest.mark.asyncio
async def test_get_user():
    async with AsyncClient(app=app, base_url="http://test") as client:
        # Create user first
        create_response = await client.post(
            "/users",
            json={"email": "test@example.com", "name": "Test", "password": "pass123"}
        )
        user_id = create_response.json()["id"]
        
        # Get user
        response = await client.get(f"/users/{user_id}")
    
    assert response.status_code == 200
    assert response.json()["id"] == user_id

def test_password_validation():
    """Test without async"""
    from app.utils import validate_password
    
    assert validate_password("SecurePass123") is True
    assert validate_password("weak") is False
```

### Fixtures for Test Setup

```python
# tests/conftest.py
import pytest
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from app.database import Base, get_db
from app.main import app

# Test database URL (use in-memory SQLite for speed)
TEST_DATABASE_URL = "sqlite+aiosqlite:///:memory:"

@pytest.fixture
async def test_db():
    """Create a fresh database for each test"""
    engine = create_async_engine(TEST_DATABASE_URL, echo=False)
    
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.create_all)
    
    AsyncTestSession = async_sessionmaker(
        engine, class_=AsyncSession, expire_on_commit=False
    )
    
    async with AsyncTestSession() as session:
        yield session
    
    async with engine.begin() as conn:
        await conn.run_sync(Base.metadata.drop_all)
    
    await engine.dispose()

@pytest.fixture
def override_get_db(test_db):
    """Override FastAPI dependency"""
    async def _override_get_db():
        yield test_db
    
    app.dependency_overrides[get_db] = _override_get_db
    yield
    app.dependency_overrides.clear()

@pytest.fixture
async def test_user(test_db):
    """Create a test user"""
    from app.models import User
    user = User(
        email="test@example.com",
        username="testuser",
        password_hash="hashed_password"
    )
    test_db.add(user)
    await test_db.commit()
    await test_db.refresh(user)
    return user

# Use fixtures in tests
@pytest.mark.asyncio
async def test_with_fixture(test_user, override_get_db):
    async with AsyncClient(app=app, base_url="http://test") as client:
        response = await client.get(f"/users/{test_user.id}")
    
    assert response.status_code == 200
    assert response.json()["email"] == test_user.email
```

### Mocking External Dependencies

```python
from unittest.mock import AsyncMock, patch

@pytest.mark.asyncio
async def test_external_api_call():
    """Mock external API calls"""
    with patch('app.services.external_api.fetch_data') as mock_fetch:
        # Setup mock
        mock_fetch.return_value = {"data": "mocked"}
        
        # Call function that uses external API
        result = await fetch_data()
        
        # Assertions
        assert result == {"data": "mocked"}
        mock_fetch.assert_called_once()

@pytest.mark.asyncio
async def test_database_error_handling():
    """Test error handling"""
    mock_db = AsyncMock()
    mock_db.execute.side_effect = Exception("Database error")
    
    with pytest.raises(Exception, match="Database error"):
        await get_users(mock_db)
```

### Parameterized Tests

```python
@pytest.mark.parametrize("email,expected", [
    ("valid@example.com", True),
    ("invalid.email", False),
    ("no-at-sign.com", False),
    ("", False),
])
def test_email_validation(email, expected):
    from app.utils import is_valid_email
    assert is_valid_email(email) == expected

@pytest.mark.asyncio
@pytest.mark.parametrize("user_role,expected_status", [
    ("admin", 200),
    ("user", 403),
    ("guest", 401),
])
async def test_role_based_access(user_role, expected_status):
    # Test authorization based on role
    pass
```

### Running Tests

```bash
# Run all tests
uv run pytest

# Run specific test file
uv run pytest tests/test_users.py

# Run specific test
uv run pytest tests/test_users.py::test_create_user

# Run with coverage
uv run pytest --cov=app --cov-report=html

# Run in parallel (faster)
uv add --dev pytest-xdist
uv run pytest -n auto

# Run with output
uv run pytest -v -s

# Run only failed tests
uv run pytest --lf

# Watch mode (rerun on file changes)
uv add --dev pytest-watch
uv run ptw
```

### pytest Configuration

```toml
# Add to pyproject.toml
[tool.pytest.ini_options]
testpaths = ["tests"]
python_files = ["test_*.py"]
python_classes = ["Test*"]
python_functions = ["test_*"]
asyncio_mode = "auto"
addopts = [
    "--strict-markers",
    "--strict-config",
    "-ra",
    "--cov=app",
    "--cov-report=term-missing",
]
```

## Environment Configuration

```python
# app/config.py
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):
    # Database
    database_url: str
    
    # JWT
    secret_key: str
    algorithm: str = "HS256"
    access_token_expire_minutes: int = 30
    
    # App
    app_name: str = "My Backend API"
    debug: bool = False
    
    # External services
    redis_url: str | None = None
    
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
    )

# Create global settings instance
settings = Settings()
```

```bash
# .env file (never commit this!)
DATABASE_URL=postgresql+asyncpg://user:pass@localhost/db
SECRET_KEY=your-secret-key-here
DEBUG=true
REDIS_URL=redis://localhost:6379
```

```bash
# .env.example (commit this as template)
DATABASE_URL=postgresql+asyncpg://user:pass@localhost/db
SECRET_KEY=change-me-in-production
DEBUG=false
REDIS_URL=redis://localhost:6379
```

## Docker for Development and Production

```dockerfile
# Dockerfile
FROM python:3.12-slim

WORKDIR /app

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

# Copy dependency files
COPY pyproject.toml uv.lock ./

# Install dependencies
RUN uv sync --frozen --no-dev

# Copy application
COPY . .

# Expose port
EXPOSE 8000

# Run application
CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8000:8000"
    environment:
      - DATABASE_URL=postgresql+asyncpg://postgres:postgres@db:5432/app
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    volumes:
      - ./app:/app/app  # Hot reload in development
  
  db:
    image: postgres:16-alpine
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
      - POSTGRES_DB=app
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
  
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

```bash
# Run with docker-compose
docker-compose up

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f app

# Run migrations
docker-compose exec app uv run alembic upgrade head

# Run tests in container
docker-compose exec app uv run pytest
```

## Makefile for Common Tasks

```makefile
# Makefile
.PHONY: help install dev test lint format clean run migrate

help:
	@echo "Available commands:"
	@echo "  make install  - Install dependencies"
	@echo "  make dev      - Install dev dependencies"
	@echo "  make test     - Run tests"
	@echo "  make lint     - Run linter"
	@echo "  make format   - Format code"
	@echo "  make run      - Run development server"
	@echo "  make migrate  - Run database migrations"

install:
	uv sync

dev:
	uv sync --all-extras

test:
	uv run pytest -v --cov=app

test-watch:
	uv run ptw

lint:
	uv run ruff check .

format:
	uv run ruff format .
	uv run ruff check --fix .

clean:
	find . -type d -name __pycache__ -exec rm -rf {} +
	find . -type f -name "*.pyc" -delete
	rm -rf .pytest_cache .coverage htmlcov

run:
	uv run uvicorn app.main:app --reload

migrate:
	uv run alembic upgrade head

migrate-create:
	@read -p "Enter migration message: " message; \
	uv run alembic revision --autogenerate -m "$$message"

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app
```

## Pre-commit Hooks (Optional)

```bash
# Install pre-commit
uv add --dev pre-commit

# Create .pre-commit-config.yaml
```

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/astral-sh/ruff-pre-commit
    rev: v0.1.9
    hooks:
      - id: ruff
        args: [--fix]
      - id: ruff-format

  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-added-large-files
```

```bash
# Install hooks
uv run pre-commit install

# Run manually
uv run pre-commit run --all-files
```

## CI/CD with GitHub Actions

```yaml
# .github/workflows/test.yml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:16
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Install uv
        uses: astral-sh/setup-uv@v1
      
      - name: Set up Python
        run: uv python install 3.12
      
      - name: Install dependencies
        run: uv sync --all-extras
      
      - name: Run linter
        run: uv run ruff check .
      
      - name: Run tests
        run: uv run pytest --cov=app --cov-report=xml
        env:
          DATABASE_URL: postgresql+asyncpg://postgres:postgres@localhost:5432/test
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.xml
```

## Key Takeaways

1. **uv replaces pip, poetry, virtualenv** - faster, simpler
2. **Ruff for linting and formatting** - 100x faster than black+flake8
3. **pytest with async support** - pytest-asyncio for testing APIs
4. **Use fixtures and mocks** - isolate tests, don't hit real services
5. **pydantic-settings for config** - type-safe environment variables
6. **Docker compose** for local development with services
7. **Makefile for common tasks** - consistent commands across team

Next: [06-production-ready.md](./06-production-ready.md) - Production-ready project structure and deployment
