# Skills Catalog

Skills are reusable capabilities that agents can execute. Each skill has defined inputs, processes, and outputs.

---

## Code Generation Skills

### skill: generate_endpoint
**Agent**: Go Backend Developer  
**Description**: Generate complete endpoint following TDD with all layers

**Input Parameters**:
```yaml
method: GET|POST|PUT|PATCH|DELETE
path: /api/v1/resource
request_schema:
  field_name:
    type: string|int|uuid|bool|time
    required: boolean
    validation: string (e.g., "email", "min=3,max=100")
response_schema: # same structure as request
auth_required: boolean
middleware: array of middleware names
```

**TDD Process**:
1. Generate test file first (red phase)
2. Generate handler skeleton
3. Generate service interface + implementation
4. Generate repository interface
5. Create SQL queries in `queries/*.sql`
6. Run `sqlc generate`
7. Implement repository using sqlc code
8. Verify tests pass (green phase)

**Output Files**:
```
internal/api/handlers/{resource}_handler.go
internal/api/handlers/{resource}_handler_test.go
internal/service/{resource}_service.go
internal/service/{resource}_service_test.go
internal/repository/{resource}_repository.go
internal/models/{resource}.go
queries/{resource}.sql
```

**Post-Generation Commands**:
```bash
sqlc generate
go test ./internal/api/handlers -v
go test ./internal/service -v
```

**Example Usage**:
```
Generate POST /api/v1/users endpoint:
- Request: {name: string (required, max=100), email: string (required, email), password: string (required, min=8)}
- Response: {id: uuid, name: string, email: string, created_at: time}
- Requires authentication: no (it's registration)
- Middleware: [RateLimitMiddleware]
```

---

### skill: generate_handler
**Agent**: Go Backend Developer  
**Description**: Generate HTTP handler with proper structure

**Input Parameters**:
```yaml
handler_name: string
http_method: string
service_dependency: string
request_dto: string
response_dto: string
```

**Generated Code Template**:
```go
type {Handler}Handler struct {
    service {Service}Service
    logger  *slog.Logger
    validator *validator.Validate
}

func New{Handler}Handler(service {Service}Service, logger *slog.Logger) *{Handler}Handler {
    return &{Handler}Handler{
        service: service,
        logger: logger,
        validator: validator.New(),
    }
}

func (h *{Handler}Handler) {MethodName}(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    reqID := middleware.GetRequestID(ctx)
    
    // Parse request
    var req {RequestDTO}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.ErrorContext(ctx, "failed to decode request",
            slog.String("request_id", reqID),
            slog.String("error", err.Error()))
        respondError(w, reqID, http.StatusBadRequest, "INVALID_JSON", "Invalid JSON format")
        return
    }
    
    // Validate
    if err := h.validator.Struct(req); err != nil {
        h.logger.WarnContext(ctx, "validation failed",
            slog.String("request_id", reqID),
            slog.String("error", err.Error()))
        respondValidationError(w, reqID, err)
        return
    }
    
    // Call service
    result, err := h.service.{Method}(ctx, &req)
    if err != nil {
        h.logger.ErrorContext(ctx, "service error",
            slog.String("request_id", reqID),
            slog.String("error", err.Error()))
        handleServiceError(w, reqID, err)
        return
    }
    
    // Success response (JSON:API format)
    respondJSON(w, http.StatusOK, map[string]interface{}{
        "data": map[string]interface{}{
            "type": "{resource_type}",
            "id": result.ID,
            "attributes": result,
        },
    })
}
```

**Includes**:
- Request ID propagation
- Structured logging with context
- Input validation
- Error handling with proper status codes
- JSON:API response format

---

### skill: generate_service
**Agent**: Go Backend Developer  
**Description**: Generate service layer with business logic

**Input Parameters**:
```yaml
service_name: string
methods:
  - name: string
    input: string
    output: string
    description: string
repository_dependencies: array of strings
```

**Generated Code Template**:
```go
// Interface
type {Service}Service interface {
    {Method}(ctx context.Context, input *{Input}) (*{Output}, error)
}

// Implementation
type {service}Service struct {
    repo {Repository}Repository
    logger *slog.Logger
}

func New{Service}Service(repo {Repository}Repository, logger *slog.Logger) {Service}Service {
    return &{service}Service{
        repo: repo,
        logger: logger,
    }
}

func (s *{service}Service) {Method}(ctx context.Context, input *{Input}) (*{Output}, error) {
    reqID := middleware.GetRequestID(ctx)
    
    s.logger.InfoContext(ctx, "starting {method}",
        slog.String("request_id", reqID))
    
    // Business logic here
    result, err := s.repo.{Method}(ctx, input)
    if err != nil {
        return nil, fmt.Errorf("failed to {action}: %w", err)
    }
    
    s.logger.InfoContext(ctx, "{method} completed",
        slog.String("request_id", reqID),
        slog.String("result_id", result.ID))
    
    return result, nil
}
```

**Includes**:
- Interface definition
- Dependency injection
- Context propagation
- Error wrapping
- Structured logging

---

### skill: generate_repository
**Agent**: Go Backend Developer  
**Description**: Generate repository interface (implementation uses sqlc)

**Input Parameters**:
```yaml
repository_name: string
model: string
operations: array of CRUD operations
```

**Generated Code Template**:
```go
// Repository interface
type {Model}Repository interface {
    Create(ctx context.Context, model *{Model}) (*{Model}, error)
    GetByID(ctx context.Context, id uuid.UUID) (*{Model}, error)
    List(ctx context.Context, offset, limit int) ([]*{Model}, error)
    Update(ctx context.Context, model *{Model}) (*{Model}, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

// Implementation using sqlc
type {model}Repository struct {
    queries *db.Queries
    logger  *slog.Logger
}

func New{Model}Repository(queries *db.Queries, logger *slog.Logger) {Model}Repository {
    return &{model}Repository{
        queries: queries,
        logger: logger,
    }
}

func (r *{model}Repository) Create(ctx context.Context, model *{Model}) (*{Model}, error) {
    reqID := middleware.GetRequestID(ctx)
    
    // Use sqlc-generated code
    result, err := r.queries.Create{Model}(ctx, db.Create{Model}Params{
        ID: uuid.New(),
        // map fields
    })
    if err != nil {
        r.logger.ErrorContext(ctx, "failed to create {model}",
            slog.String("request_id", reqID),
            slog.String("error", err.Error()))
        return nil, fmt.Errorf("create {model}: %w", err)
    }
    
    return &{Model}{
        ID: result.ID,
        // map fields
    }, nil
}
```

**Note**: Actual database operations MUST use sqlc-generated code

---

## Database Skills

### skill: generate_migration
**Agent**: Database Engineer  
**Description**: Generate migration pair (up and down)

**Input Parameters**:
```yaml
migration_name: string
tables:
  - name: string
    columns: array
    constraints: array
    indexes: array
```

**Output Files**:
```
migrations/{NNN}_{name}.up.sql
migrations/{NNN}_{name}.down.sql
```

**Up Migration Template**:
```sql
BEGIN;

CREATE TABLE IF NOT EXISTS {table_name} (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    {column_name} {type} {constraints},
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_{table}_{column} ON {table}({column});

COMMIT;
```

**Down Migration Template**:
```sql
BEGIN;

DROP INDEX IF EXISTS idx_{table}_{column};
DROP TABLE IF EXISTS {table_name};

COMMIT;
```

**Post-Generation**:
```bash
# Test up migration
migrate -path migrations -database "$DATABASE_URL" up 1

# Test down migration
migrate -path migrations -database "$DATABASE_URL" down 1
```

---

### skill: generate_sqlc_query
**Agent**: Database Engineer  
**Description**: Generate sqlc query in queries/*.sql

**Input Parameters**:
```yaml
query_name: string
operation: SELECT|INSERT|UPDATE|DELETE
table: string
return_type: :one|:many|:exec
parameters: array
```

**Query Templates**:

**SELECT (one)**:
```sql
-- name: Get{Model}ByID :one
SELECT id, name, email, created_at, updated_at
FROM {table}
WHERE id = $1;
```

**SELECT (many)**:
```sql
-- name: List{Models} :many
SELECT id, name, email, created_at, updated_at
FROM {table}
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
```

**INSERT**:
```sql
-- name: Create{Model} :one
INSERT INTO {table} (
    id, name, email, password_hash
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, name, email, created_at, updated_at;
```

**UPDATE**:
```sql
-- name: Update{Model} :one
UPDATE {table}
SET 
    name = $2,
    email = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, created_at, updated_at;
```

**DELETE**:
```sql
-- name: Delete{Model} :exec
DELETE FROM {table}
WHERE id = $1;
```

**Post-Generation**:
```bash
sqlc generate
```

---

## Testing Skills

### skill: generate_table_driven_test
**Agent**: Testing Agent  
**Description**: Generate comprehensive table-driven test

**Input Parameters**:
```yaml
function_name: string
package: string
test_cases:
  - name: string
    input: object
    expected: object
    want_error: boolean
mocks_needed: array of dependencies
```

**Generated Test Template**:
```go
func Test{Function}(t *testing.T) {
    tests := []struct {
        name    string
        input   {InputType}
        want    {OutputType}
        wantErr bool
        setup   func(*mock{Dependency})
    }{
        {
            name: "success - {scenario}",
            input: {InputType}{
                Field: "value",
            },
            want: {OutputType}{
                Field: "expected",
            },
            wantErr: false,
            setup: func(m *mock{Dependency}) {
                m.EXPECT().
                    Method(gomock.Any(), gomock.Any()).
                    Return(expectedValue, nil)
            },
        },
        {
            name: "error - {scenario}",
            input: {InputType}{
                Field: "invalid",
            },
            want: {OutputType}{},
            wantErr: true,
            setup: func(m *mock{Dependency}) {
                m.EXPECT().
                    Method(gomock.Any(), gomock.Any()).
                    Return(nil, errors.New("error"))
            },
        },
        {
            name: "edge case - {scenario}",
            input: {InputType}{
                Field: "",
            },
            want: {OutputType}{},
            wantErr: true,
            setup: nil, // No mock setup needed
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mock := mock.New{Dependency}(ctrl)
            if tt.setup != nil {
                tt.setup(mock)
            }
            
            // Execute
            svc := New{Service}(mock, slog.Default())
            got, err := svc.{Method}(context.Background(), tt.input)
            
            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Test Coverage Goals**:
- Happy path scenarios
- Error scenarios
- Edge cases (empty, nil, boundary values)
- Validation failures

---

### skill: generate_integration_test
**Agent**: Testing Agent  
**Description**: Generate integration test with real database

**Input Parameters**:
```yaml
endpoint: string
test_scenarios: array
```

**Generated Test Template**:
```go
func TestIntegration{Feature}(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    // Setup test database
    ctx := context.Background()
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer container.Terminate(ctx)
    
    // Get connection string
    connStr, err := container.ConnectionString(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    // Run migrations
    // ... migration code ...
    
    // Create test server
    server := setupTestServer(connStr)
    defer server.Close()
    
    tests := []struct {
        name           string
        method         string
        path           string
        body           interface{}
        expectedStatus int
        expectedBody   map[string]interface{}
    }{
        {
            name:   "create resource",
            method: "POST",
            path:   "/api/v1/resource",
            body: map[string]interface{}{
                "field": "value",
            },
            expectedStatus: http.StatusCreated,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Make request
            resp := makeRequest(t, server, tt.method, tt.path, tt.body)
            defer resp.Body.Close()
            
            // Assert status
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)
            
            // Assert body if provided
            if tt.expectedBody != nil {
                var got map[string]interface{}
                json.NewDecoder(resp.Body).Decode(&got)
                assert.Equal(t, tt.expectedBody, got)
            }
        })
    }
}
```

---

## Review Skills

### skill: review_sqlc_compliance
**Agent**: Code Reviewer  
**Description**: Scan codebase for manual SQL queries

**Process**:
1. Scan all `.go` files in `internal/`
2. Look for patterns:
   - `db.Query(`
   - `db.Exec(`
   - `tx.Query(`
   - `tx.Exec(`
   - Raw SQL string literals
3. Check if they're in repository layer and using sqlc

**Output Format**:
```
## sqlc Compliance Report

### ❌ Violations Found

File: internal/service/user_service.go:45
Issue: Manual SQL query detected
Code: db.Query("SELECT * FROM users WHERE id = $1", id)
Fix: Move to repository layer and define in queries/users.sql

File: internal/api/handlers/post_handler.go:78
Issue: Raw SQL in handler
Code: "INSERT INTO posts ..."
Fix: Create sqlc query and use repository

### ✅ Compliant Files
- internal/repository/user_repository.go (uses sqlc)
- internal/repository/post_repository.go (uses sqlc)
```

---

### skill: review_error_handling
**Agent**: Code Reviewer  
**Description**: Check error handling patterns

**Checks**:
- [ ] All errors are handled (no `_ = err`)
- [ ] Errors wrapped with context using `fmt.Errorf("context: %w", err)`
- [ ] Errors logged before returning
- [ ] Proper error types used
- [ ] HTTP handlers return appropriate status codes

**Output Format**:
```
## Error Handling Review

### Issues Found

Line 42: Error ignored
  _ = file.Close()
Fix: Handle or defer with error check

Line 58: Error not wrapped
  return err
Fix: return fmt.Errorf("failed to create user: %w", err)

Line 73: Error not logged
  return nil, err
Fix: Add logging before return

### Recommendations
- Consider creating custom error types for domain errors
- Add error categorization (NotFound, Validation, Internal)
```

---

### skill: review_context_propagation
**Agent**: Code Reviewer  
**Description**: Verify context usage through layers

**Checks**:
- [ ] Context passed to all layers
- [ ] Request ID stored in context (typed key)
- [ ] Request ID retrieved in logs
- [ ] Context timeouts set appropriately
- [ ] Context cancellation handled

**Output Format**:
```
## Context Propagation Review

### ✅ Correct Usage
- Handler passes context to service
- Service passes context to repository
- Request ID in all log entries

### ⚠️  Issues

Line 34: Context not passed to database call
  result, err := r.queries.GetUser(uuid)
Fix: r.queries.GetUser(ctx, uuid)

Line 89: String used as context key
  ctx = context.WithValue(ctx, "request_id", id)
Fix: Use typed key - type ctxKey string; const requestIDKey ctxKey = "request_id"

Line 102: Request ID not in log
  logger.Error("failed to process")
Fix: logger.ErrorContext(ctx, "failed to process", slog.String("request_id", reqID))
```

---

### skill: security_review
**Agent**: Code Reviewer  
**Description**: Security vulnerability scan

**Checks**:
- [ ] SQL injection prevention (sqlc usage)
- [ ] Input validation on all handlers
- [ ] Password hashing (bcrypt/argon2)
- [ ] Authentication middleware on protected routes
- [ ] No sensitive data in logs
- [ ] CORS configuration
- [ ] Rate limiting on public endpoints

**Output Format**:
```
## Security Review

### 🔴 Critical Issues

Line 67: Password stored in plaintext
  user.Password = req.Password
Fix: Hash with bcrypt before storing

Line 123: No input validation
  email := req.Email // directly used
Fix: Add email format validation

### 🟡 Warnings

Line 45: Sensitive data in logs
  logger.Info("user created", "password", user.Password)
Fix: Remove password from logs

Line 89: Missing rate limiting
  POST /api/v1/login endpoint
Fix: Add RateLimitMiddleware

### ✅ Good Practices
- Using sqlc (prevents SQL injection)
- JWT tokens properly validated
- HTTPS enforced
```

---

## Infrastructure Skills

### skill: generate_dockerfile
**Agent**: DevOps  
**Description**: Generate optimized multi-stage Dockerfile

**Input Parameters**:
```yaml
go_version: string
app_name: string
port: int
```

**Generated Dockerfile**:
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Copy migrations (if needed)
COPY --from=builder /app/migrations ./migrations

EXPOSE {port}

CMD ["./server"]
```

---

### skill: generate_docker_compose
**Agent**: DevOps  
**Description**: Generate docker-compose.yml for local development

**Input Parameters**:
```yaml
services: array (api, postgres, redis)
environment: dev|staging|prod
```

**Generated docker-compose.yml**:
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: ${DB_NAME}
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}
      REDIS_URL: redis://redis:6379
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - .:/app
    command: go run cmd/server/main.go

volumes:
  postgres_data:
  redis_data:
```

---

### skill: generate_makefile
**Agent**: DevOps  
**Description**: Generate Makefile with common commands

**Generated Makefile**:
```makefile
.PHONY: help build test run docker-up docker-down migrate-up migrate-down sqlc-gen lint fmt

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/server cmd/server/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out

run: ## Run the application
	go run cmd/server/main.go

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f api

migrate-up: ## Run database migrations up
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "$(DATABASE_URL)" down

migrate-create: ## Create new migration (usage: make migrate-create name=add_users_table)
	migrate create -ext sql -dir migrations -seq $(name)

sqlc-gen: ## Generate sqlc code
	sqlc generate

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	gofmt -s -w .
	go mod tidy

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out
```

---

## Composite Skills (Workflows)

### skill: scaffold_crud_resource
**Description**: Complete CRUD implementation
**Combines**: generate_migration + generate_sqlc_query + generate_endpoint + generate_tests

**Input**:
```yaml
resource_name: string (e.g., "post")
fields:
  - name: string
    type: string
    required: boolean
    validation: string
```

**Steps**:
1. Generate migration for table
2. Generate sqlc queries (Create, Get, List, Update, Delete)
3. Generate repository interface and implementation
4. Generate service layer
5. Generate handlers for all CRUD operations
6. Generate comprehensive tests
7. Update router

**Output**: Fully functional CRUD API

---

### skill: add_authentication
**Description**: Add JWT authentication to project

**Steps**:
1. Generate users table migration
2. Generate auth queries (create user, get by email)
3. Generate auth service (register, login, verify token)
4. Generate JWT middleware
5. Generate auth handlers (register, login, logout)
6. Add middleware to protected routes
7. Generate tests

**Output**: Complete authentication system

---

## Skill Usage Examples

### Example 1: New Endpoint
```
User: "Add POST /api/v1/posts endpoint for creating blog posts"

Agent: Go Backend Developer
Skill: generate_endpoint

Parameters:
- method: POST
- path: /api/v1/posts
- request_schema:
    title: {type: string, required: true, validation: "min=1,max=200"}
    content: {type: string, required: true}
    author_id: {type: uuid, required: true}
- response_schema:
    id: {type: uuid}
    title: {type: string}
    content: {type: string}
    author_id: {type: uuid}
    created_at: {type: time}
- auth_required: true
```

### Example 2: Code Review
```
User: "Review my user service for errors and security"

Agent: Code Reviewer
Skills: review_error_handling, security_review, review_context_propagation

Output: Comprehensive review report
```

### Example 3: Full Feature
```
User: "Implement complete comments feature for posts"

Agent: Go Backend Developer
Skill: scaffold_crud_resource

Parameters:
- resource_name: comment
- fields:
    post_id: {type: uuid, required: true}
    user_id: {type: uuid, required: true}
    content: {type: text, required: true, validation: "min=1,max=1000"}

Output: Complete CRUD with migrations, queries, tests
```
