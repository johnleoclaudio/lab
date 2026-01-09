# Database Operations with SQLAlchemy 2.0

Time: 1.5 hours | Focus: Modern async SQLAlchemy patterns for production

## Why SQLAlchemy 2.0?

- **Async support** built-in (asyncpg for PostgreSQL, aiomysql for MySQL)
- **Type safety** with modern Python type hints
- **Declarative models** with better syntax
- **Query API 2.0** - cleaner, more explicit
- **Migrations** with Alembic

## Setup

```bash
# PostgreSQL (recommended)
uv add sqlalchemy[asyncio] asyncpg alembic

# MySQL
uv add sqlalchemy[asyncio] aiomysql alembic

# SQLite (development only)
uv add sqlalchemy[asyncio] aiosqlite
```

## Database Connection Setup

```python
# database.py
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import DeclarativeBase
from typing import AsyncGenerator

# Database URL
DATABASE_URL = "postgresql+asyncpg://user:password@localhost:5432/dbname"

# Create async engine
engine = create_async_engine(
    DATABASE_URL,
    echo=True,  # Log SQL queries (disable in production)
    pool_size=20,  # Connection pool size
    max_overflow=10,  # Max connections above pool_size
    pool_pre_ping=True,  # Test connections before using
    pool_recycle=3600,  # Recycle connections after 1 hour
)

# Session factory
AsyncSessionLocal = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,  # Don't expire objects after commit
)

# Base class for models
class Base(DeclarativeBase):
    pass

# Dependency for FastAPI
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

## Defining Models

```python
# models.py
from sqlalchemy import String, Integer, Boolean, DateTime, ForeignKey, Text, Enum
from sqlalchemy.orm import Mapped, mapped_column, relationship
from datetime import datetime
from typing import Optional
import enum

class UserRole(str, enum.Enum):
    USER = "user"
    ADMIN = "admin"
    MODERATOR = "moderator"

class User(Base):
    __tablename__ = "users"
    
    # Primary key
    id: Mapped[int] = mapped_column(primary_key=True)
    
    # Required fields (not nullable)
    email: Mapped[str] = mapped_column(String(255), unique=True, index=True)
    username: Mapped[str] = mapped_column(String(50), unique=True, index=True)
    password_hash: Mapped[str] = mapped_column(String(255))
    
    # Optional fields (nullable)
    full_name: Mapped[Optional[str]] = mapped_column(String(100))
    
    # With defaults
    role: Mapped[UserRole] = mapped_column(default=UserRole.USER)
    is_active: Mapped[bool] = mapped_column(default=True)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    updated_at: Mapped[datetime] = mapped_column(
        default=datetime.utcnow,
        onupdate=datetime.utcnow
    )
    
    # Relationships
    posts: Mapped[list["Post"]] = relationship(back_populates="author", cascade="all, delete-orphan")
    profile: Mapped["Profile"] = relationship(back_populates="user", uselist=False)

class Profile(Base):
    __tablename__ = "profiles"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    user_id: Mapped[int] = mapped_column(ForeignKey("users.id"), unique=True)
    bio: Mapped[Optional[str]] = mapped_column(Text)
    avatar_url: Mapped[Optional[str]] = mapped_column(String(500))
    
    # Relationship
    user: Mapped["User"] = relationship(back_populates="profile")

class Post(Base):
    __tablename__ = "posts"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    title: Mapped[str] = mapped_column(String(200))
    content: Mapped[str] = mapped_column(Text)
    published: Mapped[bool] = mapped_column(default=False)
    created_at: Mapped[datetime] = mapped_column(default=datetime.utcnow)
    
    # Foreign key
    author_id: Mapped[int] = mapped_column(ForeignKey("users.id"))
    
    # Relationship
    author: Mapped["User"] = relationship(back_populates="posts")
    tags: Mapped[list["Tag"]] = relationship(secondary="post_tags", back_populates="posts")

class Tag(Base):
    __tablename__ = "tags"
    
    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(50), unique=True)
    
    # Many-to-many relationship
    posts: Mapped[list["Post"]] = relationship(secondary="post_tags", back_populates="tags")

# Association table for many-to-many
class PostTag(Base):
    __tablename__ = "post_tags"
    
    post_id: Mapped[int] = mapped_column(ForeignKey("posts.id"), primary_key=True)
    tag_id: Mapped[int] = mapped_column(ForeignKey("tags.id"), primary_key=True)
```

## CRUD Operations

### Create

```python
from sqlalchemy import select

async def create_user(db: AsyncSession, email: str, username: str, password: str) -> User:
    """Create a new user"""
    user = User(
        email=email,
        username=username,
        password_hash=hash_password(password)
    )
    db.add(user)
    await db.commit()
    await db.refresh(user)  # Refresh to get generated id
    return user

# Create with relationship
async def create_user_with_profile(
    db: AsyncSession,
    email: str,
    username: str,
    password: str,
    bio: str
) -> User:
    user = User(
        email=email,
        username=username,
        password_hash=hash_password(password),
        profile=Profile(bio=bio)  # Create related object
    )
    db.add(user)
    await db.commit()
    await db.refresh(user)
    return user

# Bulk insert
async def create_multiple_tags(db: AsyncSession, tag_names: list[str]) -> list[Tag]:
    tags = [Tag(name=name) for name in tag_names]
    db.add_all(tags)
    await db.commit()
    return tags
```

### Read (Query)

```python
from sqlalchemy import select, func
from sqlalchemy.orm import selectinload, joinedload

# Get by primary key
async def get_user_by_id(db: AsyncSession, user_id: int) -> User | None:
    return await db.get(User, user_id)

# Get one with filter
async def get_user_by_email(db: AsyncSession, email: str) -> User | None:
    result = await db.execute(
        select(User).where(User.email == email)
    )
    return result.scalar_one_or_none()

# Get multiple
async def get_all_users(db: AsyncSession, skip: int = 0, limit: int = 100) -> list[User]:
    result = await db.execute(
        select(User).offset(skip).limit(limit)
    )
    return result.scalars().all()

# With filtering
async def get_active_users(db: AsyncSession) -> list[User]:
    result = await db.execute(
        select(User).where(User.is_active == True)
    )
    return result.scalars().all()

# Multiple conditions
async def get_active_admins(db: AsyncSession) -> list[User]:
    result = await db.execute(
        select(User)
        .where(User.is_active == True)
        .where(User.role == UserRole.ADMIN)
    )
    return result.scalars().all()

# OR conditions
from sqlalchemy import or_

async def search_users(db: AsyncSession, query: str) -> list[User]:
    result = await db.execute(
        select(User).where(
            or_(
                User.username.ilike(f"%{query}%"),
                User.email.ilike(f"%{query}%")
            )
        )
    )
    return result.scalars().all()

# Order by
async def get_recent_users(db: AsyncSession, limit: int = 10) -> list[User]:
    result = await db.execute(
        select(User).order_by(User.created_at.desc()).limit(limit)
    )
    return result.scalars().all()

# Count
async def count_users(db: AsyncSession) -> int:
    result = await db.execute(select(func.count(User.id)))
    return result.scalar_one()
```

### Eager Loading (Avoid N+1 Queries)

```python
# ❌ BAD: N+1 query problem
async def get_users_with_posts_bad(db: AsyncSession) -> list[User]:
    result = await db.execute(select(User))
    users = result.scalars().all()
    
    # This triggers a separate query for EACH user!
    for user in users:
        _ = user.posts  # N queries
    
    return users

# ✅ GOOD: Eager load with selectinload (separate query, but only one)
async def get_users_with_posts_good(db: AsyncSession) -> list[User]:
    result = await db.execute(
        select(User).options(selectinload(User.posts))
    )
    return result.scalars().all()

# ✅ GOOD: Eager load with joinedload (single JOIN query)
async def get_users_with_profile(db: AsyncSession) -> list[User]:
    result = await db.execute(
        select(User).options(joinedload(User.profile))
    )
    return result.unique().scalars().all()  # unique() needed with joins

# Load multiple relationships
async def get_users_full(db: AsyncSession) -> list[User]:
    result = await db.execute(
        select(User)
        .options(selectinload(User.posts))
        .options(joinedload(User.profile))
    )
    return result.unique().scalars().all()

# Load nested relationships
async def get_posts_with_author_and_tags(db: AsyncSession) -> list[Post]:
    result = await db.execute(
        select(Post)
        .options(joinedload(Post.author))
        .options(selectinload(Post.tags))
    )
    return result.unique().scalars().all()
```

### Update

```python
# Update single object
async def update_user(db: AsyncSession, user_id: int, **kwargs) -> User | None:
    user = await db.get(User, user_id)
    if not user:
        return None
    
    # Update attributes
    for key, value in kwargs.items():
        setattr(user, key, value)
    
    await db.commit()
    await db.refresh(user)
    return user

# Bulk update
from sqlalchemy import update

async def deactivate_old_users(db: AsyncSession, cutoff_date: datetime) -> int:
    result = await db.execute(
        update(User)
        .where(User.created_at < cutoff_date)
        .values(is_active=False)
    )
    await db.commit()
    return result.rowcount  # Number of rows updated
```

### Delete

```python
# Delete single object
async def delete_user(db: AsyncSession, user_id: int) -> bool:
    user = await db.get(User, user_id)
    if not user:
        return False
    
    await db.delete(user)
    await db.commit()
    return True

# Bulk delete
from sqlalchemy import delete

async def delete_inactive_users(db: AsyncSession) -> int:
    result = await db.execute(
        delete(User).where(User.is_active == False)
    )
    await db.commit()
    return result.rowcount
```

## Advanced Queries

### Joins

```python
# Explicit join
async def get_posts_with_authors(db: AsyncSession) -> list[tuple[Post, User]]:
    result = await db.execute(
        select(Post, User)
        .join(User, Post.author_id == User.id)
        .where(Post.published == True)
    )
    return result.all()

# Left outer join
async def get_all_users_with_post_count(db: AsyncSession):
    result = await db.execute(
        select(User.username, func.count(Post.id).label("post_count"))
        .outerjoin(Post, User.id == Post.author_id)
        .group_by(User.id)
    )
    return result.all()
```

### Aggregations

```python
async def get_user_statistics(db: AsyncSession, user_id: int) -> dict:
    # Count posts
    post_count = await db.scalar(
        select(func.count(Post.id)).where(Post.author_id == user_id)
    )
    
    # Get date of first post
    first_post_date = await db.scalar(
        select(func.min(Post.created_at)).where(Post.author_id == user_id)
    )
    
    return {
        "post_count": post_count,
        "first_post_date": first_post_date
    }
```

### Subqueries

```python
# Users who have published posts
async def get_authors_with_published_posts(db: AsyncSession) -> list[User]:
    subq = select(Post.author_id).where(Post.published == True).distinct().subquery()
    
    result = await db.execute(
        select(User).where(User.id.in_(subq))
    )
    return result.scalars().all()
```

## Transactions

```python
# Manual transaction control
async def transfer_credits(
    db: AsyncSession,
    from_user_id: int,
    to_user_id: int,
    amount: int
) -> bool:
    try:
        # Start transaction (automatic with session)
        from_user = await db.get(User, from_user_id)
        to_user = await db.get(User, to_user_id)
        
        if from_user.credits < amount:
            raise ValueError("Insufficient credits")
        
        from_user.credits -= amount
        to_user.credits += amount
        
        await db.commit()  # Commit transaction
        return True
    except Exception as e:
        await db.rollback()  # Rollback on error
        raise

# Nested transactions (savepoints)
async def create_user_with_rollback(db: AsyncSession, email: str):
    async with db.begin_nested():  # Savepoint
        user = User(email=email)
        db.add(user)
        
        # If this fails, only this savepoint rolls back
        if await email_exists(db, email):
            raise ValueError("Email exists")
```

## Migrations with Alembic

```bash
# Initialize Alembic
alembic init alembic

# Edit alembic.ini - set database URL
# sqlalchemy.url = postgresql+asyncpg://user:pass@localhost/db

# Or use env.py for dynamic config
```

```python
# alembic/env.py
from sqlalchemy import pool
from sqlalchemy.ext.asyncio import async_engine_from_config
from alembic import context
from app.database import Base
from app.models import User, Post, Tag  # Import all models

config = context.config

# Set target metadata
target_metadata = Base.metadata

def run_migrations_offline():
    url = config.get_main_option("sqlalchemy.url")
    context.configure(
        url=url,
        target_metadata=target_metadata,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
    )
    with context.begin_transaction():
        context.run_migrations()

async def run_migrations_online():
    connectable = async_engine_from_config(
        config.get_section(config.config_ini_section),
        prefix="sqlalchemy.",
        poolclass=pool.NullPool,
    )
    
    async with connectable.connect() as connection:
        await connection.run_sync(do_run_migrations)
    
    await connectable.dispose()

def do_run_migrations(connection):
    context.configure(connection=connection, target_metadata=target_metadata)
    with context.begin_transaction():
        context.run_migrations()

if context.is_offline_mode():
    run_migrations_offline()
else:
    import asyncio
    asyncio.run(run_migrations_online())
```

```bash
# Generate migration
alembic revision --autogenerate -m "Create users table"

# Apply migration
alembic upgrade head

# Rollback one migration
alembic downgrade -1

# View current version
alembic current

# View migration history
alembic history
```

## Best Practices

```python
# 1. Use indexes for frequently queried columns
class User(Base):
    __tablename__ = "users"
    email: Mapped[str] = mapped_column(String(255), unique=True, index=True)
    username: Mapped[str] = mapped_column(String(50), index=True)

# 2. Use connection pooling
engine = create_async_engine(
    DATABASE_URL,
    pool_size=20,
    max_overflow=10,
    pool_pre_ping=True,
)

# 3. Always use async operations
# ❌ Don't use synchronous operations
user.posts  # Lazy loads synchronously

# ✅ Use eager loading
result = await db.execute(
    select(User).options(selectinload(User.posts))
)

# 4. Handle exceptions properly
async def get_user_safe(db: AsyncSession, user_id: int) -> User | None:
    try:
        user = await db.get(User, user_id)
        return user
    except Exception as e:
        logger.error(f"Error fetching user {user_id}: {e}")
        await db.rollback()
        raise

# 5. Use explicit column selection for large tables
# Instead of: select(User)
result = await db.execute(
    select(User.id, User.email, User.username)
    .where(User.is_active == True)
)

# 6. Close sessions properly (handled by get_db dependency)
async with AsyncSessionLocal() as session:
    try:
        # Do work
        await session.commit()
    except Exception:
        await session.rollback()
        raise
    finally:
        await session.close()
```

## Key Takeaways

1. **Always use async** - asyncpg for PostgreSQL, never sync drivers
2. **Eager load relationships** - avoid N+1 queries with selectinload/joinedload
3. **Use connection pooling** - configure pool_size and max_overflow
4. **Index frequently queried columns** - email, username, foreign keys
5. **Type hints with Mapped** - modern SQLAlchemy 2.0 syntax
6. **Migrations with Alembic** - version control your schema
7. **Handle transactions properly** - commit/rollback, use savepoints

Next: [05-testing-tooling.md](./05-testing-tooling.md) - Testing, linting, and modern Python tooling
