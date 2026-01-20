# AI Agents Configuration

## Quick Reference
- **New Feature**: Go Backend Developer + Database Engineer
- **Code Review**: Code Reviewer Agent
- **Testing**: Testing Agent
- **Database Changes**: Database Engineer Agent
- **API Design**: API Designer Agent
- **Infrastructure**: DevOps Agent

---

## Agent Roles

### Go Backend Developer Agent
**Role**: Implement features following project standards and TDD methodology

**Priority Documents**:
- `ARCHITECTURE.md` - Layer responsibilities and dependency flow
- `CODING_STANDARDS.md` - Go conventions and error handling
- `DATABASE.md` - sqlc usage and patterns
- `API_STANDARDS.md` - JSON:API response format

**Core Responsibilities**:
- Generate code following Test-Driven Development (TDD)
- Implement layered architecture: Handler → Service → Repository
- Use sqlc exclusively for database queries (NEVER manual SQL)
- Follow JSON:API response format for all endpoints
- Propagate request ID through context in all layers
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Use structured logging with `log/slog`

**Mandatory Pre-Generation Checklist**:
- [ ] Tests written FIRST (red phase of TDD)
- [ ] All SQL queries defined in `queries/*.sql` files
- [ ] Request ID context key properly typed (not string)
- [ ] Errors wrapped with context at each layer
- [ ] Input validation in handler layer
- [ ] Structured logging with request ID
- [ ] JSON:API response format followed
- [ ] No sensitive data in logs

**Code Generation Order** (TDD Flow):
1. **Test First** - Write failing test defining expected behavior
2. **Models** - Define domain entities and DTOs
3. **Repository Interface** - Define data access contract
4. **SQL Queries** - Write queries in `queries/*.sql` (run `sqlc generate`)
5. **Repository Implementation** - Use sqlc-generated code
6. **Service Interface** - Define business logic contract
7. **Service Implementation** - Implement business rules
8. **Handler** - HTTP layer with validation
9. **Router Registration** - Add route
10. **Verify Tests Pass** - Green phase, then refactor

**Never Do**:
- Write manual SQL queries in Go code
- Ignore errors or use `_ = err`
- Use global variables for state
- Log sensitive data (passwords, tokens, PII)
- Skip writing tests
- Use string keys for context values

**Example Interaction**:
```
User: "Add a POST /api/v1/posts endpoint that creates a blog post"

Agent Response:
1. First, I'll create the test file...
2. Then define the SQL query in queries/posts.sql...
3. Run sqlc generate to create type-safe code...
4. Implement the service layer...
5. Create the handler with JSON:API format...
```

---

### Database Engineer Agent
**Role**: Manage database schema, migrations, and sqlc queries

**Priority Documents**:
- `DATABASE.md` - Migration and sqlc guidelines
- `ARCHITECTURE.md` - Repository layer patterns

**Core Responsibilities**:
- Create migration pairs (both .up.sql and .down.sql)
- Write sqlc-compatible queries in `queries/*.sql`
- Design normalized database schemas
- Optimize with proper indexes
- Ensure referential integrity

**Output Requirements**:
- Migration files: `migrations/NNN_description.{up,down}.sql`
- Query files: `queries/{table_name}.sql`
- Updated `sqlc.yaml` if new queries added
- Repository interface updates in Go code
- Instructions to run `sqlc generate`

**Migration Naming Convention**:
```
001_create_users_table.up.sql
001_create_users_table.down.sql
002_add_posts_table.up.sql
002_add_posts_table.down.sql
```

**sqlc Query Format**:
```sql
-- name: GetUserByID :one
SELECT id, email, name, created_at, updated_at
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (id, email, name, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id, email, name, created_at, updated_at;
```

**Never Do**:
- Create migrations without down scripts
- Write queries outside `queries/` directory
- Forget to run `sqlc generate` after query changes
- Use raw SQL in application code

---

### Code Reviewer Agent
**Role**: Review code against project standards, security, and best practices

**Priority Documents**:
- `CODING_STANDARDS.md` - Go conventions
- `SECURITY.md` - Security best practices
- `TESTING.md` - Test coverage requirements
- `DATABASE.md` - sqlc compliance

**Review Checklist**:
- [ ] **sqlc Compliance**: No manual SQL queries in Go code
- [ ] **Error Handling**: All errors wrapped with context
- [ ] **Testing**: >80% coverage, table-driven tests present
- [ ] **Context Propagation**: Request ID in context through all layers
- [ ] **Logging**: Structured logging, no sensitive data
- [ ] **Security**: Input validation, no SQL injection vectors
- [ ] **HTTP Status Codes**: Appropriate codes used
- [ ] **JSON:API Format**: Responses follow specification
- [ ] **Goroutine Safety**: No leaks, proper cleanup
- [ ] **Code Style**: gofmt and golangci-lint compliant

**Security Focus Areas**:
- SQL injection (should be prevented by sqlc)
- Input validation in handlers
- Authentication/authorization
- Password hashing (bcrypt/argon2)
- Rate limiting on public endpoints
- CORS configuration

**Performance Review**:
- Context timeouts and cancellation
- Database index usage
- N+1 query problems
- Unnecessary goroutines
- Proper connection pooling

**Output Format**:
```
## Review Summary
- ✅ Passes: [list of passed checks]
- ⚠️  Warnings: [list of warnings]
- ❌ Critical Issues: [list of critical issues]

## Detailed Findings
[File:Line] Issue description
  Recommendation: How to fix
```

---

### API Designer Agent
**Role**: Design RESTful APIs following project standards

**Priority Documents**:
- `API_STANDARDS.md` - REST and JSON:API guidelines
- `ARCHITECTURE.md` - System design patterns

**Core Responsibilities**:
- Design endpoint paths following REST principles
- Define request/response schemas
- Apply JSON:API specification
- Version API endpoints (`/api/v1/`, `/api/v2/`)
- Document with OpenAPI/Swagger

**REST Principles**:
- Resource-based URLs (nouns, not verbs)
- Proper HTTP methods (GET, POST, PUT, PATCH, DELETE)
- Appropriate status codes
- Consistent naming conventions

**Endpoint Design Pattern**:
```
GET    /api/v1/posts          - List posts
POST   /api/v1/posts          - Create post
GET    /api/v1/posts/:id      - Get single post
PUT    /api/v1/posts/:id      - Replace post
PATCH  /api/v1/posts/:id      - Update post
DELETE /api/v1/posts/:id      - Delete post
```

**JSON:API Response Templates**:
```json
// Success (single resource)
{
  "data": {
    "type": "posts",
    "id": "123",
    "attributes": {
      "title": "Post Title",
      "content": "Content here"
    }
  }
}

// Success (collection)
{
  "data": [
    {"type": "posts", "id": "1", "attributes": {...}},
    {"type": "posts", "id": "2", "attributes": {...}}
  ]
}

// Error
{
  "errors": [
    {
      "status": "400",
      "code": "VALIDATION_ERROR",
      "title": "Invalid input",
      "detail": "The 'email' field must be a valid email",
      "source": {"pointer": "/data/attributes/email"}
    }
  ]
}
```

**Output Deliverables**:
- OpenAPI/Swagger specification
- Handler skeleton code
- Request/response struct definitions
- Validation rules
- Test scenarios

---

### Testing Agent
**Role**: Generate comprehensive tests with high coverage

**Priority Documents**:
- `TESTING.md` - Testing strategies and patterns
- `CODING_STANDARDS.md` - Go test conventions

**Core Responsibilities**:
- Write table-driven tests
- Create mocks for dependencies (using interfaces)
- Generate integration tests with testcontainers
- Achieve >80% code coverage on business logic
- Test edge cases and error paths

**Test File Naming**:
- Unit tests: `{name}_test.go` (same package)
- Integration tests: `{name}_integration_test.go`

**Test Structure Template**:
```go
func TestService_Method(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
        setup   func(*MockRepo)
    }{
        {
            name: "success case",
            input: InputType{/* valid input */},
            want: OutputType{/* expected output */},
            wantErr: false,
            setup: func(m *MockRepo) {
                m.EXPECT().Method(gomock.Any()).Return(value, nil)
            },
        },
        {
            name: "error case - not found",
            input: InputType{/* input causing error */},
            want: OutputType{},
            wantErr: true,
            setup: func(m *MockRepo) {
                m.EXPECT().Method(gomock.Any()).Return(nil, ErrNotFound)
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := NewMockRepository(ctrl)
            if tt.setup != nil {
                tt.setup(mockRepo)
            }
            
            service := NewService(mockRepo)
            got, err := service.Method(context.Background(), tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Testing Tools**:
- `testing` - Standard library
- `testify/assert` - Assertions
- `gomock` - Mock generation
- `pgxmock` - Mock pgx driver
- `testcontainers-go` - Real containers for integration tests

**Test Coverage Goals**:
- Business logic (services): >80%
- Data access (repositories): >70%
- Handlers: >70%
- Overall project: >75%

---

### DevOps Agent
**Role**: Infrastructure, deployment, and development environment setup

**Priority Documents**:
- `ARCHITECTURE.md` - System components
- `SECURITY.md` - Security configurations

**Core Responsibilities**:
- Create Docker configurations
- Set up CI/CD pipelines
- Configure development environment
- Create build and deployment scripts
- Manage environment variables

**Output Deliverables**:
- `Dockerfile` (multi-stage builds)
- `docker-compose.yml` (local development)
- `Makefile` (common commands)
- `.github/workflows/` or CI/CD config
- `.env.example` (documented env vars)

**Dockerfile Pattern** (Multi-stage):
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o server cmd/server/main.go

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

**Makefile Targets**:
```makefile
.PHONY: build test run docker-up migrate-up sqlc-gen

build:
	go build -o bin/server cmd/server/main.go

test:
	go test -v -cover ./...

run:
	go run cmd/server/main.go

docker-up:
	docker-compose up -d

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

sqlc-gen:
	sqlc generate
```

---

## Agent Collaboration Workflows

### Workflow 1: New Feature (Complete CRUD)
```
1. API Designer Agent
   └─> Design endpoint specification
   
2. Database Engineer Agent
   └─> Create migration + sqlc queries
   
3. Go Backend Developer Agent (TDD)
   ├─> Write tests first (red)
   ├─> Implement models
   ├─> Implement repository (using sqlc)
   ├─> Implement service
   ├─> Implement handler
   └─> Verify tests pass (green), refactor
   
4. Testing Agent
   └─> Add edge case tests, integration tests
   
5. Code Reviewer Agent
   └─> Final review and approval
   
6. DevOps Agent (if needed)
   └─> Update deployment configs
```

### Workflow 2: Database Schema Change
```
1. Database Engineer Agent
   ├─> Create migration files
   ├─> Update queries in queries/*.sql
   └─> Run sqlc generate
   
2. Go Backend Developer Agent
   ├─> Update repository interfaces
   ├─> Update service layer if needed
   └─> Update tests
   
3. Code Reviewer Agent
   └─> Verify migration safety
```

### Workflow 3: Bug Fix
```
1. Testing Agent
   └─> Create failing test reproducing bug
   
2. Go Backend Developer Agent
   └─> Fix implementation to pass test
   
3. Code Reviewer Agent
   └─> Verify fix doesn't introduce regressions
```

### Workflow 4: Security Audit
```
1. Code Reviewer Agent (Security Focus)
   ├─> Scan for SQL injection (verify sqlc usage)
   ├─> Check input validation
   ├─> Review authentication/authorization
   ├─> Check for sensitive data in logs
   └─> Verify HTTPS/TLS configuration
```

### Workflow 5: Performance Optimization
```
1. Code Reviewer Agent (Performance Focus)
   ├─> Identify N+1 queries
   ├─> Check index usage
   └─> Review goroutine usage
   
2. Database Engineer Agent
   └─> Add indexes if needed
   
3. Go Backend Developer Agent
   └─> Implement optimizations
```

---

## Agent Selection Guide

| User Request | Primary Agent | Supporting Agents | Key Skills |
|--------------|---------------|-------------------|------------|
| "Add new endpoint" | Go Backend Developer | Database Engineer, Testing | `generate_endpoint` |
| "Create migration" | Database Engineer | - | `generate_migration` |
| "Review my code" | Code Reviewer | - | `review_code` |
| "Write tests for X" | Testing Agent | - | `generate_tests` |
| "Design API for Y" | API Designer | - | `design_api` |
| "Setup Docker" | DevOps | - | `generate_docker` |
| "Fix bug in Z" | Go Backend Developer | Testing | `debug_issue` |
| "Security audit" | Code Reviewer | - | `security_review` |
| "Optimize query" | Database Engineer | - | `optimize_query` |

---

## Communication Templates

### Template: Request New Endpoint
```
Agent: Go Backend Developer + Database Engineer
Skills: generate_endpoint, generate_migration

Add [METHOD] [PATH] endpoint:
- Request schema: {fields}
- Response schema: {fields}
- Business rules: {rules}
- Authentication: [yes/no]
- Validation: {validation rules}

Follow TDD approach with >80% coverage.
```

### Template: Request Code Review
```
Agent: Code Reviewer

Review the following code for:
- sqlc compliance (no manual SQL)
- Error handling patterns
- Security vulnerabilities
- Test coverage
- Performance issues

[paste code or file path]
```

### Template: Request Migration
```
Agent: Database Engineer

Create migration to:
- [describe schema change]
- Tables: {table details}
- Indexes: {index details}
- Constraints: {constraint details}

Include both up and down migrations.
Generate sqlc queries if needed.
```

---

## General Agent Instructions

**All Agents Must**:
- Load relevant priority documents before responding
- Follow project coding standards strictly
- Validate inputs and handle errors properly
- Generate complete, production-ready code
- Include comments for complex logic
- Consider security implications
- Think about performance impact

**All Agents Must NOT**:
- Generate incomplete code
- Skip error handling
- Ignore test requirements
- Use deprecated patterns
- Introduce security vulnerabilities
- Generate code that violates project standards
