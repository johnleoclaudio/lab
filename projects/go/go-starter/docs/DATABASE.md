# Database Guidelines

This document covers database design, migrations, and sqlc usage. **Critical**: All database queries MUST use sqlc-generated code.

## Table of Contents
1. [Technology Stack](#technology-stack)
2. [sqlc Usage (CRITICAL)](#sqlc-usage-critical)
3. [Migrations](#migrations)
4. [Schema Design](#schema-design)
5. [Indexes](#indexes)
6. [Transactions](#transactions)
7. [Common Patterns](#common-patterns)

---

## Technology Stack

### Database
- **PostgreSQL 15+** - Primary database
- **pgx/v5** - Recommended driver (better performance than lib/pq)

### ORM/Query Builder
- **sqlc** (REQUIRED) - Type-safe SQL code generation
- **golang-migrate** - Database migrations

### Connection
```go
import (
    "github.com/jackc/pgx/v5/pgxpool"
)

// Connection pool configuration
config, err := pgxpool.ParseConfig(databaseURL)
if err != nil {
    return err
}

// Pool settings
config.MaxConns = 25
config.MinConns = 5
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute

pool, err := pgxpool.NewWithConfig(context.Background(), config)
```

---

## sqlc Usage (CRITICAL)

### Rule: ALL Database Queries MUST Use sqlc

**Never write manual SQL in Go code. Always define queries in `.sql` files and generate type-safe code with sqlc.**

### Why sqlc?
- ✅ Type safety at compile time
- ✅ SQL injection prevention
- ✅ Better performance (no reflection)
- ✅ Explicit query definitions
- ✅ Easy to review and test

### Setup

#### 1. Install sqlc
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

#### 2. Configure sqlc.yaml
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "queries"
    schema: "migrations"
    gen:
      go:
        package: "db"
        out: "db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
```

#### 3. Project Structure
```
project/
├── migrations/           # Schema definitions
│   ├── 001_init.up.sql
│   └── 001_init.down.sql
├── queries/              # Query definitions
│   ├── users.sql
│   └── posts.sql
├── db/                   # Generated code (gitignored)
│   ├── db.go
│   ├── models.go
│   ├── users.sql.go
│   └── posts.sql.go
└── sqlc.yaml
```

### Writing Queries

#### Query Annotations
```sql
-- name: GetUser :one
-- :one - Returns single row (or error if not found)
-- :many - Returns multiple rows (can be empty slice)
-- :exec - Executes query, returns error only
-- :execrows - Returns number of affected rows
```

#### Example: queries/users.sql
```sql
-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, email, name, password_hash, created_at, updated_at
FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT id, email, name, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    id, email, name, password_hash
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, email, name, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET 
    name = $2,
    email = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, email, name, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);
```

#### Generate Code
```bash
# Generate type-safe Go code
sqlc generate

# This creates files in db/ directory
# - models.go (struct definitions)
# - users.sql.go (query implementations)
```

### Using Generated Code

#### In Repository
```go
package repository

import (
    "context"
    "fmt"
    "log/slog"
    
    "github.com/google/uuid"
    "github.com/yourorg/project/db"
    "github.com/yourorg/project/internal/models"
)

type UserRepository interface {
    Create(ctx context.Context, user *models.User) (*models.User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
    GetByEmail(ctx context.Context, email string) (*models.User, error)
    List(ctx context.Context, offset, limit int) ([]*models.User, error)
    Update(ctx context.Context, user *models.User) (*models.User, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
    queries *db.Queries
    logger  *slog.Logger
}

func NewUserRepository(queries *db.Queries, logger *slog.Logger) UserRepository {
    return &userRepository{
        queries: queries,
        logger:  logger,
    }
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
    reqID := middleware.GetRequestID(ctx)
    
    // Use sqlc-generated CreateUser function
    result, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        ID:           user.ID,
        Email:        user.Email,
        Name:         user.Name,
        PasswordHash: user.PasswordHash,
    })
    if err != nil {
        r.logger.ErrorContext(ctx, "failed to create user",
            slog.String("request_id", reqID),
            slog.String("error", err.Error()))
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    return &models.User{
        ID:        result.ID,
        Email:     result.Email,
        Name:      result.Name,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
    }, nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
    result, err := r.queries.GetUserByID(ctx, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, models.ErrNotFound
        }
        return nil, fmt.Errorf("get user: %w", err)
    }
    
    return &models.User{
        ID:        result.ID,
        Email:     result.Email,
        Name:      result.Name,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
    }, nil
}
```

### Advanced Queries

#### Joins
```sql
-- name: GetUserWithPosts :many
SELECT 
    u.id, u.email, u.name,
    p.id AS post_id, p.title, p.content
FROM users u
LEFT JOIN posts p ON u.id = p.author_id
WHERE u.id = $1;
```

#### Aggregations
```sql
-- name: GetUserStats :one
SELECT 
    u.id,
    u.name,
    COUNT(p.id) AS post_count,
    COUNT(c.id) AS comment_count
FROM users u
LEFT JOIN posts p ON u.id = p.author_id
LEFT JOIN comments c ON u.id = c.user_id
WHERE u.id = $1
GROUP BY u.id, u.name;
```

#### Conditional Queries
```sql
-- name: SearchUsers :many
SELECT id, email, name, created_at
FROM users
WHERE 
    ($1::text IS NULL OR email ILIKE '%' || $1 || '%')
    AND ($2::text IS NULL OR name ILIKE '%' || $2 || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;
```

#### Bulk Insert
```sql
-- name: BatchInsertUsers :copyfrom
INSERT INTO users (
    id, email, name, password_hash
) VALUES (
    $1, $2, $3, $4
);
```

### What NOT to Do

```go
// ❌ NEVER: Manual SQL in Go code
func (r *userRepository) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
    var user User
    err := r.db.QueryRow(ctx, 
        "SELECT id, email FROM users WHERE id = $1", 
        id,
    ).Scan(&user.ID, &user.Email)
    return &user, err
}

// ❌ NEVER: String concatenation (SQL injection!)
query := "SELECT * FROM users WHERE email = '" + email + "'"

// ❌ NEVER: Building queries dynamically in code
query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", tableName)

// ✅ ALWAYS: Use sqlc-generated code
result, err := r.queries.GetUserByID(ctx, id)
```

---

## Migrations

### Migration Files

#### Naming Convention
```
NNN_description.up.sql
NNN_description.down.sql

Examples:
001_create_users_table.up.sql
001_create_users_table.down.sql
002_add_posts_table.up.sql
002_add_posts_table.down.sql
```

#### Create Migration
```bash
# Using golang-migrate
migrate create -ext sql -dir migrations -seq create_users_table

# This creates:
# migrations/000001_create_users_table.up.sql
# migrations/000001_create_users_table.down.sql
```

### Migration Structure

#### Up Migration (001_create_users_table.up.sql)
```sql
BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at);

-- Trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
```

#### Down Migration (001_create_users_table.down.sql)
```sql
BEGIN;

DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;

COMMIT;
```

### Migration Best Practices

#### Always Use Transactions
```sql
BEGIN;
-- All changes here
COMMIT;
```

#### Always Create Down Migration
```sql
-- Every up migration must have a corresponding down
-- Test that down migration actually reverses the up migration
```

#### Make Migrations Idempotent
```sql
-- ✅ Good: Idempotent
CREATE TABLE IF NOT EXISTS users (...);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
DROP TABLE IF EXISTS users;

-- ❌ Bad: Not idempotent
CREATE TABLE users (...);  -- Fails if table exists
DROP TABLE users;          -- Fails if table doesn't exist
```

#### One Logical Change Per Migration
```sql
-- ✅ Good: Single purpose
-- 001_create_users_table.sql
-- 002_add_posts_table.sql
-- 003_add_user_role_column.sql

-- ❌ Bad: Multiple unrelated changes
-- 001_create_all_tables.sql (users, posts, comments, tags...)
```

### Running Migrations

```bash
# Up (apply all pending migrations)
migrate -path migrations -database "$DATABASE_URL" up

# Up (apply specific number)
migrate -path migrations -database "$DATABASE_URL" up 2

# Down (rollback one migration)
migrate -path migrations -database "$DATABASE_URL" down 1

# Force version (if migrations are stuck)
migrate -path migrations -database "$DATABASE_URL" force 5

# Version (check current version)
migrate -path migrations -database "$DATABASE_URL" version
```

### In Code (Optional)
```go
import "github.com/golang-migrate/migrate/v4"

func runMigrations(databaseURL string) error {
    m, err := migrate.New(
        "file://migrations",
        databaseURL,
    )
    if err != nil {
        return err
    }
    
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    
    return nil
}
```

---

## Schema Design

### Data Types

#### Recommended Types
```sql
-- IDs
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- Strings
email VARCHAR(255) NOT NULL
name VARCHAR(100) NOT NULL
description TEXT

-- Numbers
age INTEGER
price DECIMAL(10, 2)
count BIGINT

-- Boolean
is_active BOOLEAN DEFAULT true

-- Timestamps
created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()

-- JSON
metadata JSONB
```

#### UUID vs Serial
```sql
-- ✅ Preferred: UUID (better for distributed systems)
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- ⚠️  Alternative: Serial (simpler, sequential)
id SERIAL PRIMARY KEY

-- ⚠️  Alternative: BigSerial (for very large tables)
id BIGSERIAL PRIMARY KEY
```

### Constraints

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL CHECK (length(name) > 0),
    age INTEGER CHECK (age >= 0 AND age <= 150),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Foreign Keys

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    -- Foreign key with cascading delete
    CONSTRAINT fk_posts_author
        FOREIGN KEY (author_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- Alternatives for ON DELETE:
-- CASCADE - Delete posts when user is deleted
-- SET NULL - Set author_id to NULL when user is deleted
-- RESTRICT - Prevent deletion if posts exist
-- NO ACTION - Same as RESTRICT
```

### Default Values

```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    published BOOLEAN DEFAULT false,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

---

## Indexes

### When to Add Indexes

✅ Add index when:
- Column used in WHERE clause frequently
- Column used in JOIN conditions
- Column used in ORDER BY
- Foreign key columns
- Unique constraints need enforcement

❌ Don't add index when:
- Table is very small (< 1000 rows)
- Column values are not selective (e.g., boolean)
- Column is rarely queried
- Write performance is critical

### Index Types

#### B-Tree (Default)
```sql
-- Most common, good for equality and range queries
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

#### Composite Index
```sql
-- For queries filtering by multiple columns
CREATE INDEX idx_posts_author_created ON posts(author_id, created_at);

-- Order matters! This is good for:
-- WHERE author_id = ? AND created_at > ?
-- WHERE author_id = ?
-- But NOT optimized for: WHERE created_at > ?
```

#### Partial Index
```sql
-- Index only subset of rows
CREATE INDEX idx_published_posts ON posts(created_at) 
WHERE published = true;

-- Good for: SELECT * FROM posts WHERE published = true ORDER BY created_at
```

#### GIN Index (for JSONB and arrays)
```sql
CREATE INDEX idx_posts_metadata ON posts USING GIN (metadata);

-- Good for: SELECT * FROM posts WHERE metadata @> '{"category": "tech"}';
```

#### Full Text Search
```sql
-- Add tsvector column
ALTER TABLE posts ADD COLUMN search_vector tsvector;

-- Update trigger
CREATE TRIGGER update_search_vector
BEFORE INSERT OR UPDATE ON posts
FOR EACH ROW EXECUTE FUNCTION
tsvector_update_trigger(search_vector, 'pg_catalog.english', title, content);

-- GIN index for fast search
CREATE INDEX idx_posts_search ON posts USING GIN (search_vector);

-- Query
SELECT * FROM posts 
WHERE search_vector @@ to_tsquery('postgresql & database');
```

### Index Naming Convention
```sql
-- Pattern: idx_{table}_{columns}
idx_users_email
idx_posts_author_created
idx_comments_post_user
```

---

## Transactions

### When to Use Transactions

Use transactions when:
- Multiple related database operations must succeed together
- Data consistency is critical
- Operations span multiple tables

### In sqlc

#### Define Transaction Query
```sql
-- queries/transactions.sql

-- name: CreateUserWithProfile :exec
-- This will be used within a transaction
INSERT INTO users (id, email, name, password_hash)
VALUES ($1, $2, $3, $4);

-- name: CreateProfile :exec
INSERT INTO profiles (user_id, bio, avatar_url)
VALUES ($1, $2, $3);
```

#### Use in Repository
```go
func (r *userRepository) CreateUserWithProfile(
    ctx context.Context,
    user *models.User,
    profile *models.Profile,
) error {
    // Begin transaction
    tx, err := r.pool.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx) // Rollback if not committed
    
    // Create queries with transaction
    queries := r.queries.WithTx(tx)
    
    // Create user
    err = queries.CreateUserWithProfile(ctx, db.CreateUserWithProfileParams{
        ID:           user.ID,
        Email:        user.Email,
        Name:         user.Name,
        PasswordHash: user.PasswordHash,
    })
    if err != nil {
        return fmt.Errorf("create user: %w", err)
    }
    
    // Create profile
    err = queries.CreateProfile(ctx, db.CreateProfileParams{
        UserID:    user.ID,
        Bio:       profile.Bio,
        AvatarUrl: profile.AvatarURL,
    })
    if err != nil {
        return fmt.Errorf("create profile: %w", err)
    }
    
    // Commit transaction
    if err := tx.Commit(ctx); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }
    
    return nil
}
```

### Transaction Isolation Levels
```go
// Default: Read Committed
tx, err := pool.BeginTx(ctx, pgx.TxOptions{
    IsoLevel: pgx.ReadCommitted,
})

// Serializable (strictest, slowest)
tx, err := pool.BeginTx(ctx, pgx.TxOptions{
    IsoLevel: pgx.Serializable,
})
```

---

## Common Patterns

### Soft Delete
```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_posts_deleted_at ON posts(deleted_at) WHERE deleted_at IS NULL;

-- Queries
-- name: SoftDeletePost :exec
UPDATE posts SET deleted_at = NOW() WHERE id = $1;

-- name: ListActivePosts :many
SELECT * FROM posts WHERE deleted_at IS NULL;
```

### Audit Trail
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by UUID,
    deleted_at TIMESTAMP WITH TIME ZONE,
    deleted_by UUID
);
```

### Enum Types
```sql
-- Option 1: CHECK constraint
CREATE TABLE posts (
    status VARCHAR(20) CHECK (status IN ('draft', 'published', 'archived'))
);

-- Option 2: PostgreSQL ENUM (preferred)
CREATE TYPE post_status AS ENUM ('draft', 'published', 'archived');

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    status post_status DEFAULT 'draft'
);
```

### Pagination
```sql
-- name: ListPostsPaginated :many
SELECT id, title, created_at
FROM posts
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts WHERE deleted_at IS NULL;
```

---

## Checklist

Before deploying database changes:

- [ ] Migration has both up and down files
- [ ] Migrations use transactions (BEGIN/COMMIT)
- [ ] Migrations are idempotent (IF EXISTS/IF NOT EXISTS)
- [ ] Indexes created for foreign keys
- [ ] Indexes created for frequently queried columns
- [ ] Queries defined in queries/*.sql (not in Go code)
- [ ] `sqlc generate` runs without errors
- [ ] Generated code compiles
- [ ] Repository uses sqlc-generated code only
- [ ] No manual SQL in Go files
- [ ] Migration tested (up then down)
