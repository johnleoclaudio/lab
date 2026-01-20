# Quick Start Guide for AI Agents

This guide provides ready-to-use prompts for common development tasks. Copy and customize these prompts for your specific needs.

---

## Table of Contents
1. [Project Setup](#project-setup)
2. [Adding Features](#adding-features)
3. [Database Operations](#database-operations)
4. [Testing](#testing)
5. [Code Review](#code-review)
6. [Infrastructure](#infrastructure)

---

## Project Setup

### Initialize New Project
```
Agent: DevOps Agent + Go Backend Developer Agent
Skills: generate_project_structure, generate_dockerfile, generate_makefile

Create a new Go backend project with:
- Project name: [your-project-name]
- Database: PostgreSQL
- Authentication: JWT
- Caching: Redis

Include:
- Complete directory structure
- Docker and docker-compose setup
- Makefile with common commands
- CI/CD pipeline template
- .env.example with documented variables
```

### Setup Database
```
Agent: Database Engineer Agent
Skills: generate_migration, generate_sqlc_query

Create initial database setup:
- Users table (id, email, password_hash, name, created_at, updated_at)
- Include unique index on email
- Create sqlc queries for user CRUD operations
- Generate both up and down migrations
```

---

## Adding Features

### Add Simple CRUD Endpoint
```
Agent: Go Backend Developer Agent
Skills: generate_endpoint

Add [METHOD] /api/v1/[resource] endpoint:
- Request schema:
  * field1: string (required, max=100)
  * field2: int (required, min=1)
  * field3: uuid (required)
- Response schema:
  * id: uuid
  * field1: string
  * field2: int
  * created_at: time
- Authentication: [yes/no]
- Follow TDD approach
- Include comprehensive tests
```

**Example - Blog Posts**:
```
Agent: Go Backend Developer Agent
Skills: generate_endpoint

Add POST /api/v1/posts endpoint:
- Request schema:
  * title: string (required, min=1, max=200)
  * content: string (required, min=1)
  * author_id: uuid (required)
  * tags: []string (optional, max=10)
- Response schema:
  * id: uuid
  * title: string
  * content: string
  * author_id: uuid
  * tags: []string
  * created_at: time
  * updated_at: time
- Authentication: yes
- Validation: Check author exists before creating post
- Follow TDD approach with >80% coverage
```

### Add Complete CRUD Resource
```
Agent: Go Backend Developer Agent + Database Engineer Agent
Skills: scaffold_crud_resource

Implement complete CRUD for [resource_name]:

Fields:
- [field1]: [type] ([constraints])
- [field2]: [type] ([constraints])
- [field3]: [type] ([constraints])

Requirements:
- All CRUD endpoints (Create, Read, List, Update, Delete)
- Database migration with indexes
- sqlc queries for all operations
- Service layer with business logic
- Full test coverage (>80%)
- JSON:API response format
- Authentication required for Create/Update/Delete
```

**Example - Comments System**:
```
Agent: Go Backend Developer Agent + Database Engineer Agent
Skills: scaffold_crud_resource

Implement complete CRUD for comments:

Fields:
- post_id: uuid (required, foreign key to posts)
- user_id: uuid (required, foreign key to users)
- content: text (required, min=1, max=1000)
- parent_id: uuid (optional, for nested comments)

Requirements:
- GET /api/v1/posts/:post_id/comments (list comments for post)
- POST /api/v1/comments (create comment)
- PATCH /api/v1/comments/:id (update own comment)
- DELETE /api/v1/comments/:id (delete own comment)
- Business rules:
  * Users can only edit/delete their own comments
  * Cannot comment on non-existent posts
  * Parent comment must exist if parent_id provided
- Include index on post_id and user_id
- Full test coverage including auth checks
```

### Add Authentication
```
Agent: Go Backend Developer Agent
Skills: add_authentication

Implement JWT authentication system:
- User registration endpoint (POST /api/v1/auth/register)
- Login endpoint (POST /api/v1/auth/login)
- JWT token generation and validation
- Authentication middleware
- Password hashing with bcrypt
- Token refresh mechanism
- Include comprehensive tests

Password requirements:
- Minimum 8 characters
- At least one uppercase, lowercase, number

JWT configuration:
- Secret from environment variable
- Token expiry: 24 hours
- Refresh token expiry: 7 days
```

### Add Authorization/Permissions
```
Agent: Go Backend Developer Agent
Skills: generate_middleware, generate_service

Implement role-based access control (RBAC):
- Roles: admin, user, guest
- Permissions table in database
- Role assignment to users
- Authorization middleware
- Permission checking in service layer

Endpoints:
- POST /api/v1/users/:id/roles (assign role - admin only)
- GET /api/v1/users/:id/permissions (list permissions)

Protect existing endpoints:
- /api/v1/posts: Create (authenticated), Update/Delete (owner or admin)
- /api/v1/users: Create (public), Update/Delete (owner or admin)
```

---

## Database Operations

### Create Migration
```
Agent: Database Engineer Agent
Skills: generate_migration

Create migration to [action]:
- Migration name: [descriptive_name]
- Tables:
  * [table1]: [columns with types and constraints]
  * [table2]: [columns with types and constraints]
- Indexes:
  * [table.column] for [reason]
- Foreign keys:
  * [table.column] references [other_table.column]

Include both up and down migrations.
Generate corresponding sqlc queries.
```

**Example - Add Posts Table**:
```
Agent: Database Engineer Agent
Skills: generate_migration

Create migration to add posts table:
- Migration name: create_posts_table
- Tables:
  * posts:
    - id: UUID PRIMARY KEY
    - title: VARCHAR(200) NOT NULL
    - content: TEXT NOT NULL
    - author_id: UUID NOT NULL (FK to users.id)
    - published: BOOLEAN DEFAULT false
    - created_at: TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    - updated_at: TIMESTAMP WITH TIME ZONE DEFAULT NOW()
- Indexes:
  * posts.author_id (for filtering by author)
  * posts.published, posts.created_at (for listing published posts)
- Foreign keys:
  * posts.author_id references users.id ON DELETE CASCADE

Include both up and down migrations.
Generate sqlc queries for:
- CreatePost
- GetPostByID
- ListPostsByAuthor
- UpdatePost
- DeletePost
- PublishPost
```

### Add Column to Existing Table
```
Agent: Database Engineer Agent
Skills: generate_migration

Create migration to add [column_name] to [table_name]:
- Column: [name] [type] [constraints]
- Default value: [value]
- Nullable: [yes/no]
- Index: [yes/no]

Update sqlc queries if needed.
```

### Add Database Index
```
Agent: Database Engineer Agent
Skills: generate_migration

Create migration to add index on [table].[columns]:
- Index name: idx_[table]_[columns]
- Columns: [column1, column2, ...]
- Type: [BTREE/HASH/GIN/etc]
- Reason: [performance improvement for specific query]

Include analysis of query performance impact.
```

### Complex Query
```
Agent: Database Engineer Agent
Skills: generate_sqlc_query

Create sqlc query to [description]:
- Query name: [CamelCaseName]
- Operation: [SELECT/INSERT/UPDATE/DELETE]
- Tables: [tables involved]
- Joins: [join conditions]
- Filters: [WHERE conditions]
- Return type: :one/:many/:exec

Optimize for performance with proper indexes.
```

**Example - Get User with Post Count**:
```
Agent: Database Engineer Agent
Skills: generate_sqlc_query

Create sqlc query to get user with their post count:
- Query name: GetUserWithPostCount
- Operation: SELECT
- Tables: users, posts
- Joins: LEFT JOIN posts ON users.id = posts.author_id
- Filters: users.id = $1
- Return type: :one
- Return fields:
  * user.id, user.name, user.email
  * COUNT(posts.id) as post_count

Include query in queries/users.sql
```

---

## Testing

### Generate Unit Tests
```
Agent: Testing Agent
Skills: generate_table_driven_test

Generate unit tests for [function/method]:
- Package: [package_path]
- Function: [function_name]
- Test cases:
  * Success case: [description]
  * Error case: [description]
  * Edge case: [description]
- Mocks needed: [dependencies to mock]
- Coverage goal: >80%
```

**Example - Service Method Tests**:
```
Agent: Testing Agent
Skills: generate_table_driven_test

Generate unit tests for UserService.CreateUser:
- Package: internal/service
- Function: CreateUser
- Test cases:
  * Success: Valid user creation
  * Error: Email already exists
  * Error: Invalid email format (caught by validation)
  * Error: Repository failure
  * Edge: Empty name
  * Edge: Very long email
- Mocks needed: UserRepository
- Coverage goal: >80%

Use gomock for repository mocking.
Include setup and teardown logic.
```

### Generate Integration Tests
```
Agent: Testing Agent
Skills: generate_integration_test

Generate integration tests for [feature]:
- Endpoint: [HTTP method] [path]
- Test scenarios:
  * [scenario 1]: [expected behavior]
  * [scenario 2]: [expected behavior]
- Use testcontainers for real database
- Test complete request/response cycle
- Verify database state after operations
```

**Example - User Registration Flow**:
```
Agent: Testing Agent
Skills: generate_integration_test

Generate integration tests for user registration:
- Endpoint: POST /api/v1/auth/register
- Test scenarios:
  * Success: New user registration
  * Error: Duplicate email
  * Error: Invalid email format
  * Error: Password too short
  * Success: Login with created user
- Use testcontainers for PostgreSQL
- Verify:
  * User created in database
  * Password properly hashed
  * Response contains user ID and no password
  * Can login with credentials
```

### Generate Test Fixtures
```
Agent: Testing Agent

Create test fixtures for [resource]:
- Valid fixtures: [list scenarios]
- Invalid fixtures: [list edge cases]
- Helper functions for fixture setup
- Cleanup functions for test teardown
```

---

## Code Review

### General Code Review
```
Agent: Code Reviewer Agent
Skills: review_code, review_error_handling, review_context_propagation

Review the following code:

[paste code or provide file path]

Check for:
- sqlc compliance (no manual SQL)
- Error handling patterns
- Context propagation
- Security vulnerabilities
- Test coverage
- Code style and conventions
- Performance issues

Provide specific line-by-line feedback.
```

### Security Audit
```
Agent: Code Reviewer Agent
Skills: security_review

Perform security audit on [component/file]:

[paste code or provide file path]

Focus on:
- SQL injection vulnerabilities
- Input validation
- Authentication/authorization
- Password handling
- Sensitive data exposure
- CORS configuration
- Rate limiting
- XSS prevention
```

### Performance Review
```
Agent: Code Reviewer Agent
Skills: performance_review

Review [component] for performance issues:

[paste code or provide file path]

Check for:
- N+1 query problems
- Missing database indexes
- Inefficient algorithms
- Unnecessary goroutines
- Memory leaks
- Connection pool configuration
- Cache usage opportunities
```

### sqlc Compliance Check
```
Agent: Code Reviewer Agent
Skills: review_sqlc_compliance

Scan the codebase for sqlc compliance violations:
- Check all files in internal/repository
- Flag any manual SQL queries
- Verify all database operations use sqlc-generated code
- List violations with file:line references

Provide recommendations for fixing violations.
```

---

## Infrastructure

### Setup Docker
```
Agent: DevOps Agent
Skills: generate_dockerfile, generate_docker_compose

Create Docker setup:
- Multi-stage Dockerfile for Go app
- docker-compose.yml with:
  * Go API service
  * PostgreSQL
  * Redis (optional)
- Development and production configurations
- Health checks for all services
- Volume mounts for development
```

### Create Makefile
```
Agent: DevOps Agent
Skills: generate_makefile

Create Makefile with commands for:
- build: Build the application
- test: Run tests with coverage
- run: Run the application locally
- docker-up: Start Docker containers
- docker-down: Stop Docker containers
- migrate-up: Run migrations
- migrate-down: Rollback migrations
- sqlc-gen: Generate sqlc code
- lint: Run linter
- fmt: Format code
- clean: Clean build artifacts
```

### Setup CI/CD
```
Agent: DevOps Agent
Skills: generate_ci_cd

Create CI/CD pipeline for [platform]:
- Platform: [GitHub Actions/GitLab CI/CircleCI]
- Pipeline stages:
  * Lint and format check
  * Run tests
  * Build Docker image
  * Run migrations
  * Deploy to [staging/production]
- Include:
  * Environment-specific configurations
  * Secret management
  * Automated rollback on failure
```

---

## Agent Selection Matrix

| Task Type | Primary Agent | Supporting Agents | Typical Duration |
|-----------|---------------|-------------------|------------------|
| New endpoint (simple) | Go Backend Developer | - | 5-10 min |
| New endpoint (complex) | Go Backend Developer | Database Engineer, Testing | 15-30 min |
| Database migration | Database Engineer | - | 5 min |
| Full CRUD resource | Go Backend Developer | Database Engineer, Testing | 20-40 min |
| Authentication | Go Backend Developer | - | 30-60 min |
| Code review | Code Reviewer | - | 5-15 min |
| Test generation | Testing Agent | - | 10-20 min |
| Infrastructure setup | DevOps | - | 15-30 min |

---

## Best Practices for Prompts

### ✅ Good Prompts
- **Specific**: Include exact field names, types, and constraints
- **Complete**: Provide all necessary context and requirements
- **Clear**: Use precise language, avoid ambiguity
- **Structured**: Break down complex requests into steps

### ❌ Avoid
- Vague requests like "add a feature"
- Missing critical details (data types, validation rules)
- Requesting multiple unrelated tasks at once
- Skipping testing or security requirements

### Template Structure
```
Agent: [Specific Agent]
Skills: [Relevant Skills]

[Clear task description]

Requirements:
- [Requirement 1]
- [Requirement 2]
- [Requirement 3]

[Optional: Specific code examples or patterns to follow]
```

---

## Progressive Complexity

### Level 1: Simple Tasks (Start Here)
- Add single endpoint with basic CRUD
- Create simple migration
- Generate unit tests for one function
- Review small code snippet

### Level 2: Moderate Tasks
- Add complete CRUD resource
- Implement authentication
- Create complex database queries
- Generate integration tests

### Level 3: Complex Tasks
- Multi-table operations with transactions
- Role-based access control
- Performance optimization
- Complete feature with multiple endpoints

### Level 4: Advanced Tasks
- Microservices architecture
- Event-driven systems
- Complex authorization rules
- Full system security audit

---

## Troubleshooting

### If Generated Code Doesn't Compile
```
Agent: Code Reviewer Agent

The following generated code has compilation errors:

[paste error messages]
[paste relevant code]

Please fix the issues and explain what went wrong.
```

### If Tests Are Failing
```
Agent: Testing Agent

The following tests are failing:

[paste test failures]

Please:
1. Identify why tests are failing
2. Fix the test cases
3. Verify edge cases are covered
```

### If Migrations Fail
```
Agent: Database Engineer Agent

Migration failed with error:

[paste error message]

Migration file: [file path]

Please:
1. Identify the issue
2. Fix the migration
3. Ensure up and down migrations work correctly
```
