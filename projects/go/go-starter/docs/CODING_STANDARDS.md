# Coding Standards

This document defines Go coding conventions and best practices for this project.

## Table of Contents
1. [General Guidelines](#general-guidelines)
2. [Naming Conventions](#naming-conventions)
3. [Code Organization](#code-organization)
4. [Error Handling](#error-handling)
5. [Logging](#logging)
6. [Context Usage](#context-usage)
7. [Common Patterns](#common-patterns)
8. [Anti-Patterns to Avoid](#anti-patterns-to-avoid)

---

## General Guidelines

### Follow Go Idioms
- Read [Effective Go](https://golang.org/doc/effective_go)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (no exceptions)
- Run `go vet` before committing
- Use `golangci-lint` for comprehensive linting

### Code Quality Rules
```bash
# Format code
gofmt -s -w .

# Vet code
go vet ./...

# Lint code
golangci-lint run

# All must pass before commit
```

### Function Size
- Keep functions small and focused (< 50 lines ideal)
- Single responsibility per function
- If function does multiple things, split it
- Extract complex logic into helper functions

### Comments
```go
// ✅ Good: Explains WHY
// Hash password using bcrypt to prevent rainbow table attacks
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ❌ Bad: Explains WHAT (code already shows that)
// Generate password hash
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// ✅ Good: Package comment
// Package service implements business logic for user management.
// It orchestrates repository calls and enforces business rules.
package service
```

---

## Naming Conventions

### Packages
```go
// ✅ Good: lowercase, single word
package repository
package handler
package middleware

// ❌ Bad: mixed case, underscores, plural
package Repository
package repo_layer
package handlers
```

### Files
```go
// ✅ Good: snake_case
user_handler.go
user_service.go
user_repository.go
user_handler_test.go

// ❌ Bad: camelCase, PascalCase
userHandler.go
UserService.go
```

### Interfaces
```go
// ✅ Good: noun or adjective
type UserRepository interface { ... }
type Validator interface { ... }
type Readable interface { ... }

// ❌ Bad: "I" prefix (not Go style)
type IUserRepository interface { ... }

// ✅ Good: single-method interfaces end in "er"
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}
```

### Structs
```go
// ✅ Good: PascalCase for exported
type UserService struct {
    repo   repository.UserRepository
    logger *slog.Logger
}

// ✅ Good: camelCase for unexported
type userService struct {
    repo   repository.UserRepository
    logger *slog.Logger
}
```

### Functions and Methods
```go
// ✅ Good: PascalCase for exported
func NewUserService() UserService { ... }
func (s *UserService) CreateUser() { ... }

// ✅ Good: camelCase for unexported
func hashPassword(password string) string { ... }
func (s *userService) validateEmail(email string) error { ... }
```

### Variables
```go
// ✅ Good: camelCase, descriptive
var userRepository repository.UserRepository
var maxRetries = 3
var errNotFound = errors.New("not found")

// ❌ Bad: single letter for non-trivial scope
var u repository.UserRepository // What is 'u'?

// ✅ OK: single letter for short scope
for i := 0; i < 10; i++ {
    // 'i' is fine here
}

for _, user := range users {
    // OK in short loops
}
```

### Constants
```go
// ✅ Good: PascalCase or SCREAMING_SNAKE_CASE for groups
const MaxConnections = 100
const DefaultTimeout = 30 * time.Second

// ✅ Good: Grouped constants
const (
    StatusPending  = "pending"
    StatusActive   = "active"
    StatusInactive = "inactive"
)

// ✅ Good: iota for enums
type Status int

const (
    StatusPending Status = iota
    StatusActive
    StatusInactive
)
```

### Acronyms
```go
// ✅ Good: Consistent casing
type UserID uuid.UUID    // ID, not Id
var httpClient *http.Client  // HTTP, not Http
var apiURL string        // URL or API, not Url or Api

// Exception: At start of unexported name, all lowercase
var urlParser Parser     // url, not URL
var httpServer Server    // http, not HTTP
```

---

## Code Organization

### Import Grouping
```go
import (
    // 1. Standard library
    "context"
    "fmt"
    "time"
    
    // 2. External dependencies
    "github.com/google/uuid"
    "github.com/lib/pq"
    
    // 3. Internal packages
    "github.com/yourorg/project/internal/models"
    "github.com/yourorg/project/internal/repository"
)
```

### Package Structure
```go
// ✅ Good: Group related functionality
internal/
├── api/
│   ├── handlers/      // All HTTP handlers
│   └── middleware/    // All HTTP middleware
├── service/          // All business logic
├── repository/       // All data access
└── models/           // All domain models

// ❌ Bad: Feature-based without clear layers
internal/
├── user/
│   ├── handler.go
│   ├── service.go
│   └── repository.go
└── post/
    ├── handler.go
    ├── service.go
    └── repository.go
```

### File Organization
```go
// In each file:
// 1. Package declaration
package handler

// 2. Imports
import (...)

// 3. Constants
const (...)

// 4. Types
type UserHandler struct {...}

// 5. Constructor
func NewUserHandler() *UserHandler {...}

// 6. Methods
func (h *UserHandler) CreateUser() {...}

// 7. Helper functions (unexported)
func validateEmail() {...}
```

### Struct Field Ordering
```go
type UserHandler struct {
    // 1. Dependencies (injected)
    userService service.UserService
    logger      *slog.Logger
    validator   *validator.Validate
    
    // 2. Configuration
    config *Config
    
    // 3. State (if any, prefer stateless)
    cache map[string]interface{}
}
```

---

## Error Handling

### Rule: Always Handle Errors
```go
// ✅ Good: Handle error
result, err := doSomething()
if err != nil {
    return fmt.Errorf("do something: %w", err)
}

// ❌ Bad: Ignore error
result, _ := doSomething()

// ❌ Bad: Panic for normal errors
result, err := doSomething()
if err != nil {
    panic(err) // Only use panic for programming errors
}
```

### Wrap Errors with Context
```go
// ✅ Good: Wrap with context using %w
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    return user, nil
}

// Error chain is preserved:
// "failed to register: create user: insert user: connection refused"

// ❌ Bad: No wrapping (loses error chain)
if err != nil {
    return nil, err
}

// ❌ Bad: Using %v instead of %w (breaks errors.Is/As)
if err != nil {
    return nil, fmt.Errorf("create user: %v", err)
}
```

### Custom Error Types
```go
// ✅ Good: Custom errors for domain logic
var (
    ErrNotFound           = errors.New("resource not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrUnauthorized       = errors.New("unauthorized")
    ErrInvalidInput       = errors.New("invalid input")
)

// ✅ Good: Structured error with context
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Message)
}

// Usage: Check errors
if errors.Is(err, ErrNotFound) {
    // Handle not found
}

if errors.As(err, &validationErr) {
    // Handle validation error
}
```

### Error Logging
```go
// ✅ Good: Log before returning
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    reqID := middleware.GetRequestID(ctx)
    
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        s.logger.ErrorContext(ctx, "failed to create user",
            slog.String("request_id", reqID),
            slog.String("email", req.Email),
            slog.String("error", err.Error()))
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    return user, nil
}

// ❌ Bad: No logging
if err != nil {
    return nil, err
}
```

### Early Returns
```go
// ✅ Good: Early returns reduce nesting
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, http.StatusBadRequest, err)
        return
    }
    
    if err := h.validator.Struct(req); err != nil {
        respondError(w, http.StatusBadRequest, err)
        return
    }
    
    user, err := h.service.CreateUser(r.Context(), &req)
    if err != nil {
        respondError(w, http.StatusInternalServerError, err)
        return
    }
    
    respondJSON(w, http.StatusCreated, user)
}

// ❌ Bad: Nested if statements
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
        if err := h.validator.Struct(req); err == nil {
            user, err := h.service.CreateUser(r.Context(), &req)
            if err == nil {
                respondJSON(w, http.StatusCreated, user)
            } else {
                respondError(w, http.StatusInternalServerError, err)
            }
        } else {
            respondError(w, http.StatusBadRequest, err)
        }
    } else {
        respondError(w, http.StatusBadRequest, err)
    }
}
```

---

## Logging

### Use Structured Logging (slog)
```go
import "log/slog"

// ✅ Good: Structured logging with context
logger.InfoContext(ctx, "user created",
    slog.String("request_id", reqID),
    slog.String("user_id", user.ID.String()),
    slog.String("email", user.Email))

// ❌ Bad: String concatenation
logger.Info("user created: " + user.ID.String())

// ❌ Bad: Printf-style (not searchable)
logger.Infof("user created: %s", user.ID)
```

### Log Levels
```go
// Debug: Development information
logger.DebugContext(ctx, "processing request",
    slog.String("request_id", reqID),
    slog.Any("payload", req))

// Info: Normal operations
logger.InfoContext(ctx, "user created",
    slog.String("user_id", user.ID.String()))

// Warn: Recoverable issues
logger.WarnContext(ctx, "rate limit approaching",
    slog.String("user_id", user.ID.String()),
    slog.Int("requests", count))

// Error: Failures requiring attention
logger.ErrorContext(ctx, "failed to create user",
    slog.String("request_id", reqID),
    slog.String("error", err.Error()))
```

### Always Include Request ID
```go
// ✅ Good: Request ID in all logs
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    reqID := middleware.GetRequestID(ctx)
    
    s.logger.InfoContext(ctx, "creating user",
        slog.String("request_id", reqID),
        slog.String("email", req.Email))
    
    // ... implementation
    
    s.logger.InfoContext(ctx, "user created",
        slog.String("request_id", reqID),
        slog.String("user_id", user.ID.String()))
    
    return user, nil
}

// ❌ Bad: No request ID (can't trace logs)
s.logger.Info("creating user", slog.String("email", req.Email))
```

### Never Log Sensitive Data
```go
// ✅ Good: Omit sensitive data
logger.InfoContext(ctx, "user login",
    slog.String("request_id", reqID),
    slog.String("email", email))

// ❌ Bad: Logging password
logger.InfoContext(ctx, "user login",
    slog.String("email", email),
    slog.String("password", password)) // NEVER!

// ❌ Bad: Logging token
logger.InfoContext(ctx, "request authenticated",
    slog.String("token", token)) // NEVER!
```

### Log at Appropriate Places
```go
// ✅ Log at:
// - Service boundaries (entry/exit)
// - Before/after external calls
// - Error paths
// - Significant state changes

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    reqID := middleware.GetRequestID(ctx)
    
    // Entry point
    s.logger.InfoContext(ctx, "starting user creation",
        slog.String("request_id", reqID))
    
    // Before external call
    s.logger.DebugContext(ctx, "calling repository",
        slog.String("request_id", reqID))
    
    user, err := s.repo.Create(ctx, req)
    if err != nil {
        // Error path
        s.logger.ErrorContext(ctx, "repository error",
            slog.String("request_id", reqID),
            slog.String("error", err.Error()))
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    // Success path
    s.logger.InfoContext(ctx, "user created successfully",
        slog.String("request_id", reqID),
        slog.String("user_id", user.ID.String()))
    
    return user, nil
}
```

---

## Context Usage

### Always Pass Context
```go
// ✅ Good: Context as first parameter
func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    // Pass context through all layers
    return s.repo.Create(ctx, req)
}

// ❌ Bad: No context
func (s *userService) CreateUser(req *CreateUserRequest) (*User, error) {
    return s.repo.Create(req)
}
```

### Use Typed Context Keys
```go
// ✅ Good: Typed key prevents collisions
type ctxKey string

const (
    requestIDKey ctxKey = "request_id"
    userIDKey    ctxKey = "user_id"
)

func SetRequestID(ctx context.Context, reqID string) context.Context {
    return context.WithValue(ctx, requestIDKey, reqID)
}

func GetRequestID(ctx context.Context) string {
    if reqID, ok := ctx.Value(requestIDKey).(string); ok {
        return reqID
    }
    return ""
}

// ❌ Bad: String keys can collide
ctx = context.WithValue(ctx, "request_id", reqID)
```

### Set Timeouts
```go
// ✅ Good: Set timeout for long operations
func (s *userService) ProcessBatch(ctx context.Context, users []User) error {
    // Set 30 second timeout for batch processing
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel()
    
    for _, user := range users {
        if err := s.processUser(ctx, user); err != nil {
            return err
        }
        
        // Check if context canceled
        if ctx.Err() != nil {
            return fmt.Errorf("batch processing canceled: %w", ctx.Err())
        }
    }
    
    return nil
}
```

### Handle Context Cancellation
```go
// ✅ Good: Check context cancellation
func (s *userService) LongRunningTask(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return fmt.Errorf("task canceled: %w", ctx.Err())
        default:
            // Continue work
            if err := s.doWork(); err != nil {
                return err
            }
        }
    }
}
```

---

## Common Patterns

### Constructor Pattern
```go
// ✅ Good: Return interface, accept dependencies
func NewUserService(
    repo repository.UserRepository,
    logger *slog.Logger,
) service.UserService {
    return &userService{
        repo:   repo,
        logger: logger,
    }
}

// Interface for testing
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
}

// Concrete implementation
type userService struct {
    repo   repository.UserRepository
    logger *slog.Logger
}
```

### Options Pattern (for complex configuration)
```go
// ✅ Good: Options pattern for many optional parameters
type ServerOptions struct {
    port        int
    timeout     time.Duration
    maxRequests int
}

type ServerOption func(*ServerOptions)

func WithPort(port int) ServerOption {
    return func(o *ServerOptions) {
        o.port = port
    }
}

func WithTimeout(timeout time.Duration) ServerOption {
    return func(o *ServerOptions) {
        o.timeout = timeout
    }
}

func NewServer(opts ...ServerOption) *Server {
    options := &ServerOptions{
        port:    8080,           // defaults
        timeout: 30 * time.Second,
    }
    
    for _, opt := range opts {
        opt(options)
    }
    
    return &Server{options: options}
}

// Usage
server := NewServer(
    WithPort(9000),
    WithTimeout(60 * time.Second),
)
```

### Table-Driven Tests
```go
// ✅ Good: Table-driven tests
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"missing domain", "user@", true},
        {"empty", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## Anti-Patterns to Avoid

### Global Variables
```go
// ❌ Bad: Global state
var db *sql.DB

func GetUser(id string) (*User, error) {
    return db.Query(...)
}

// ✅ Good: Dependency injection
type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) GetUser(ctx context.Context, id string) (*User, error) {
    return r.db.QueryContext(ctx, ...)
}
```

### Panic for Normal Errors
```go
// ❌ Bad: Panic for recoverable errors
func GetUser(id string) *User {
    user, err := repo.GetUser(id)
    if err != nil {
        panic(err) // Don't panic!
    }
    return user
}

// ✅ Good: Return error
func GetUser(ctx context.Context, id string) (*User, error) {
    user, err := repo.GetUser(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get user: %w", err)
    }
    return user, nil
}

// ✅ OK: Panic for programming errors
func MustCompileRegex(pattern string) *regexp.Regexp {
    re, err := regexp.Compile(pattern)
    if err != nil {
        panic("invalid regex pattern") // OK for init-time errors
    }
    return re
}
```

### Ignoring Defer Errors
```go
// ❌ Bad: Ignore defer error
defer file.Close()

// ✅ Good: Handle defer error
defer func() {
    if err := file.Close(); err != nil {
        logger.Error("failed to close file", slog.String("error", err.Error()))
    }
}()

// ✅ Good: In function that returns error
func processFile(path string) (err error) {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer func() {
        if closeErr := file.Close(); closeErr != nil && err == nil {
            err = closeErr
        }
    }()
    
    // Process file
    return nil
}
```

### Init Functions with Side Effects
```go
// ❌ Bad: Heavy work in init
func init() {
    db = connectToDatabase()  // Bad: makes testing hard
    loadAllData()             // Bad: slows down imports
}

// ✅ Good: Explicit initialization
func main() {
    db, err := connectToDatabase()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // ... rest of main
}
```

---

## Checklist

Before committing code, verify:

- [ ] `gofmt -s -w .` applied
- [ ] `go vet ./...` passes
- [ ] `golangci-lint run` passes
- [ ] All errors handled (no `_ = err`)
- [ ] Errors wrapped with context (`%w`)
- [ ] Structured logging with request ID
- [ ] No sensitive data in logs
- [ ] Context passed through all layers
- [ ] Tests written (>80% coverage)
- [ ] No global variables
- [ ] No panics for normal errors
- [ ] Exported items have comments
