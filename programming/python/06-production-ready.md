# Production-Ready Backend Structure

Time: 1 hour | Focus: Project structure, deployment, and best practices

## Recommended Project Structure

```
my-backend-api/
├── app/
│   ├── __init__.py
│   ├── main.py              # FastAPI app entry point
│   ├── config.py            # Configuration management
│   ├── database.py          # Database connection and session
│   │
│   ├── models/              # SQLAlchemy models
│   │   ├── __init__.py
│   │   ├── user.py
│   │   ├── post.py
│   │   └── tag.py
│   │
│   ├── schemas/             # Pydantic schemas (request/response)
│   │   ├── __init__.py
│   │   ├── user.py
│   │   ├── post.py
│   │   └── common.py
│   │
│   ├── api/                 # API routes
│   │   ├── __init__.py
│   │   ├── deps.py          # Shared dependencies
│   │   └── v1/
│   │       ├── __init__.py
│   │       ├── users.py
│   │       ├── posts.py
│   │       └── auth.py
│   │
│   ├── services/            # Business logic
│   │   ├── __init__.py
│   │   ├── user_service.py
│   │   ├── post_service.py
│   │   └── auth_service.py
│   │
│   ├── crud/                # Database operations (alternative to services)
│   │   ├── __init__.py
│   │   ├── user.py
│   │   └── post.py
│   │
│   └── utils/               # Utility functions
│       ├── __init__.py
│       ├── security.py      # Password hashing, JWT
│       ├── email.py
│       └── validators.py
│
├── tests/
│   ├── __init__.py
│   ├── conftest.py          # Pytest fixtures
│   ├── test_api/
│   │   ├── test_users.py
│   │   └── test_auth.py
│   └── test_services/
│       └── test_user_service.py
│
├── alembic/                 # Database migrations
│   ├── versions/
│   └── env.py
│
├── scripts/                 # Utility scripts
│   ├── init_db.py
│   └── seed_data.py
│
├── .env.example             # Environment variables template
├── .gitignore
├── alembic.ini
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── pyproject.toml           # Project config and dependencies
├── uv.lock                  # Locked dependencies
└── README.md
```

## Complete Example Files

### 1. app/main.py

```python
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager

from app.config import settings
from app.database import engine
from app.api.v1 import users, posts, auth

@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup
    print("Starting up...")
    yield
    # Shutdown
    print("Shutting down...")
    await engine.dispose()

app = FastAPI(
    title=settings.app_name,
    version="1.0.0",
    lifespan=lifespan,
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.allowed_origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(auth.router, prefix="/api/v1/auth", tags=["auth"])
app.include_router(users.router, prefix="/api/v1/users", tags=["users"])
app.include_router(posts.router, prefix="/api/v1/posts", tags=["posts"])

@app.get("/")
async def root():
    return {"message": "Welcome to the API"}

@app.get("/health")
async def health_check():
    return {"status": "healthy"}
```

### 2. app/config.py

```python
from pydantic_settings import BaseSettings, SettingsConfigDict
from typing import List

class Settings(BaseSettings):
    # App
    app_name: str = "My Backend API"
    debug: bool = False
    environment: str = "development"
    
    # Database
    database_url: str
    
    # Security
    secret_key: str
    algorithm: str = "HS256"
    access_token_expire_minutes: int = 30
    
    # CORS
    allowed_origins: List[str] = ["http://localhost:3000"]
    
    # Redis (optional)
    redis_url: str | None = None
    
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        case_sensitive=False,
    )

settings = Settings()
```

### 3. app/database.py

```python
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase
from typing import AsyncGenerator

from app.config import settings

engine = create_async_engine(
    settings.database_url,
    echo=settings.debug,
    pool_size=20,
    max_overflow=10,
    pool_pre_ping=True,
)

AsyncSessionLocal = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,
)

class Base(DeclarativeBase):
    pass

async def get_db() -> AsyncGenerator[AsyncSession, None]:
    async with AsyncSessionLocal() as session:
        try:
            yield session
            await session.commit()
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()
```

### 4. app/models/user.py

```python
from sqlalchemy import String, Boolean, DateTime, Enum as SQLEnum
from sqlalchemy.orm import Mapped, mapped_column, relationship
from datetime import datetime
import enum

from app.database import Base

class UserRole(str, enum.Enum):
    USER = "user"
    ADMIN = "admin"

class User(Base):
    __tablename__ = "users"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    email: Mapped[str] = mapped_column(String(255), unique=True, index=True)
    username: Mapped[str] = mapped_column(String(50), unique=True, index=True)
    password_hash: Mapped[str] = mapped_column(String(255))
    full_name: Mapped[str | None] = mapped_column(String(100))
    role: Mapped[UserRole] = mapped_column(SQLEnum(UserRole), default=UserRole.USER)
    is_active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    updated_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow,
        onupdate=datetime.utcnow
    )
    
    # Relationships
    posts: Mapped[list["Post"]] = relationship(back_populates="author")
```

### 5. app/schemas/user.py

```python
from pydantic import BaseModel, EmailStr, Field, ConfigDict
from datetime import datetime
from app.models.user import UserRole

class UserBase(BaseModel):
    email: EmailStr
    username: str = Field(min_length=3, max_length=50)

class UserCreate(UserBase):
    password: str = Field(min_length=8, max_length=100)
    full_name: str | None = None

class UserUpdate(BaseModel):
    email: EmailStr | None = None
    full_name: str | None = None
    password: str | None = None

class UserResponse(UserBase):
    id: int
    full_name: str | None
    role: UserRole
    is_active: bool
    created_at: datetime
    
    model_config = ConfigDict(from_attributes=True)

class UserInDB(UserResponse):
    password_hash: str
```

### 6. app/api/deps.py

```python
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from sqlalchemy.ext.asyncio import AsyncSession
from jose import JWTError, jwt

from app.database import get_db
from app.config import settings
from app.models.user import User
from sqlalchemy import select

security = HTTPBearer()

async def get_current_user(
    credentials: HTTPAuthorizationCredentials = Depends(security),
    db: AsyncSession = Depends(get_db)
) -> User:
    """Get current authenticated user"""
    try:
        token = credentials.credentials
        payload = jwt.decode(token, settings.secret_key, algorithms=[settings.algorithm])
        user_id: int = payload.get("sub")
        if user_id is None:
            raise HTTPException(status_code=401, detail="Invalid token")
    except JWTError:
        raise HTTPException(status_code=401, detail="Invalid token")
    
    result = await db.execute(select(User).where(User.id == user_id))
    user = result.scalar_one_or_none()
    
    if user is None:
        raise HTTPException(status_code=404, detail="User not found")
    
    return user

async def get_current_active_user(
    current_user: User = Depends(get_current_user)
) -> User:
    """Ensure user is active"""
    if not current_user.is_active:
        raise HTTPException(status_code=403, detail="Inactive user")
    return current_user

async def get_admin_user(
    current_user: User = Depends(get_current_active_user)
) -> User:
    """Ensure user is admin"""
    from app.models.user import UserRole
    if current_user.role != UserRole.ADMIN:
        raise HTTPException(status_code=403, detail="Admin access required")
    return current_user
```

### 7. app/api/v1/users.py

```python
from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select

from app.database import get_db
from app.models.user import User
from app.schemas.user import UserCreate, UserResponse, UserUpdate
from app.api.deps import get_current_active_user, get_admin_user
from app.utils.security import get_password_hash

router = APIRouter()

@router.post("/", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
async def create_user(
    user_in: UserCreate,
    db: AsyncSession = Depends(get_db)
):
    """Create new user"""
    # Check if email exists
    result = await db.execute(select(User).where(User.email == user_in.email))
    if result.scalar_one_or_none():
        raise HTTPException(status_code=400, detail="Email already registered")
    
    # Create user
    user = User(
        email=user_in.email,
        username=user_in.username,
        password_hash=get_password_hash(user_in.password),
        full_name=user_in.full_name
    )
    db.add(user)
    await db.commit()
    await db.refresh(user)
    return user

@router.get("/me", response_model=UserResponse)
async def get_current_user_info(
    current_user: User = Depends(get_current_active_user)
):
    """Get current user info"""
    return current_user

@router.get("/{user_id}", response_model=UserResponse)
async def get_user(
    user_id: int,
    db: AsyncSession = Depends(get_db),
    _: User = Depends(get_current_active_user)  # Require auth
):
    """Get user by ID"""
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return user

@router.get("/", response_model=list[UserResponse])
async def list_users(
    skip: int = 0,
    limit: int = 100,
    db: AsyncSession = Depends(get_db),
    _: User = Depends(get_admin_user)  # Admin only
):
    """List all users (admin only)"""
    result = await db.execute(select(User).offset(skip).limit(limit))
    return result.scalars().all()

@router.patch("/{user_id}", response_model=UserResponse)
async def update_user(
    user_id: int,
    user_update: UserUpdate,
    db: AsyncSession = Depends(get_db),
    current_user: User = Depends(get_current_active_user)
):
    """Update user (own profile or admin)"""
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    
    # Check permissions
    from app.models.user import UserRole
    if user.id != current_user.id and current_user.role != UserRole.ADMIN:
        raise HTTPException(status_code=403, detail="Not authorized")
    
    # Update fields
    update_data = user_update.model_dump(exclude_unset=True)
    if "password" in update_data:
        update_data["password_hash"] = get_password_hash(update_data.pop("password"))
    
    for field, value in update_data.items():
        setattr(user, field, value)
    
    await db.commit()
    await db.refresh(user)
    return user

@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(
    user_id: int,
    db: AsyncSession = Depends(get_db),
    _: User = Depends(get_admin_user)  # Admin only
):
    """Delete user (admin only)"""
    user = await db.get(User, user_id)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    
    await db.delete(user)
    await db.commit()
    return None
```

### 8. app/utils/security.py

```python
from datetime import datetime, timedelta
from typing import Any
from jose import jwt
from passlib.context import CryptContext

from app.config import settings

pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

def verify_password(plain_password: str, hashed_password: str) -> bool:
    return pwd_context.verify(plain_password, hashed_password)

def get_password_hash(password: str) -> str:
    return pwd_context.hash(password)

def create_access_token(subject: int | str, expires_delta: timedelta | None = None) -> str:
    if expires_delta:
        expire = datetime.utcnow() + expires_delta
    else:
        expire = datetime.utcnow() + timedelta(minutes=settings.access_token_expire_minutes)
    
    to_encode = {"exp": expire, "sub": str(subject)}
    encoded_jwt = jwt.encode(to_encode, settings.secret_key, algorithm=settings.algorithm)
    return encoded_jwt
```

## Logging Configuration

```python
# app/logging_config.py
import logging
from app.config import settings

def setup_logging():
    logging.basicConfig(
        level=logging.DEBUG if settings.debug else logging.INFO,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )

# Use in main.py
from app.logging_config import setup_logging
setup_logging()

# Use in code
import logging
logger = logging.getLogger(__name__)

logger.info("User created successfully")
logger.error(f"Failed to create user: {error}")
```

## Deployment Checklist

### Environment Variables (Production)

```bash
# .env.production
APP_NAME="My Backend API"
DEBUG=false
ENVIRONMENT=production
DATABASE_URL=postgresql+asyncpg://user:pass@db.example.com/prod
SECRET_KEY=super-secret-key-change-this
ALLOWED_ORIGINS=["https://myapp.com"]
```

### Production Dockerfile

```dockerfile
FROM python:3.12-slim as base

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Install uv
COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

# Copy dependency files
COPY pyproject.toml uv.lock ./

# Install dependencies (production only)
RUN uv sync --frozen --no-dev

# Copy application
COPY . .

# Create non-root user
RUN useradd -m -u 1000 appuser && chown -R appuser:appuser /app
USER appuser

EXPOSE 8000

CMD ["uv", "run", "uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--workers", "4"]
```

### Health Checks

```python
# Add to main.py
@app.get("/health")
async def health_check(db: AsyncSession = Depends(get_db)):
    """Health check endpoint"""
    try:
        # Check database
        await db.execute(text("SELECT 1"))
        
        # Check Redis (if using)
        # await redis.ping()
        
        return {"status": "healthy", "database": "connected"}
    except Exception as e:
        logger.error(f"Health check failed: {e}")
        raise HTTPException(status_code=503, detail="Service unavailable")
```

## Production Tips

1. **Use environment variables** - Never hardcode secrets
2. **Enable connection pooling** - Configure pool_size and max_overflow
3. **Add rate limiting** - Use slowapi or nginx
4. **Implement caching** - Redis for frequently accessed data
5. **Monitor performance** - Use APM tools (DataDog, New Relic)
6. **Set up logging** - Centralized logging with ELK or CloudWatch
7. **Database indexes** - Index all foreign keys and frequently queried columns
8. **Use migrations** - Always use Alembic for schema changes
9. **Write tests** - Aim for 80%+ coverage
10. **Document API** - FastAPI auto-generates docs, keep them updated

## Deployment Options

### 1. Docker + Cloud Run (GCP)

```bash
# Build and push
docker build -t gcr.io/project-id/my-api .
docker push gcr.io/project-id/my-api

# Deploy
gcloud run deploy my-api \
  --image gcr.io/project-id/my-api \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated
```

### 2. Kubernetes

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-api
  template:
    metadata:
      labels:
        app: my-api
    spec:
      containers:
      - name: api
        image: my-api:latest
        ports:
        - containerPort: 8000
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: db-secrets
              key: url
```

### 3. Traditional VPS (Digital Ocean, etc.)

```bash
# On server
git clone repo
cd repo
uv sync
uv run alembic upgrade head

# Use systemd service
sudo systemctl start my-api
sudo systemctl enable my-api

# Or use supervisor
supervisorctl start my-api
```

## Key Takeaways

1. **Organize by feature** - models, schemas, routes, services per domain
2. **Separate concerns** - config, database, business logic, API routes
3. **Use dependency injection** - for auth, database, shared logic
4. **Type everything** - models, schemas, function signatures
5. **Test thoroughly** - unit tests, integration tests, end-to-end
6. **Log appropriately** - info for important events, error with context
7. **Monitor in production** - health checks, metrics, alerts
8. **Document as you go** - FastAPI docs, code comments, README

Congratulations! You now have the 20% of Python knowledge that covers 80% of backend development needs in 2025.
