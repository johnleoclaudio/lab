# Architecture Overview

## Project Structure

```
project-root/
├── cmd/                    # Application entry points
│   └── server/
│       └── main.go        # Main application
├── internal/              # Private application code
│   ├── api/              # HTTP layer
│   │   ├── handlers/     # HTTP request handlers
│   │   │   ├── user_handler.go
│   │   │   ├── user_handler_test.go
│   │   │   └── ...
│   │   ├── middleware/   # HTTP middleware
│   │   │   ├── auth.go
│   │   │   ├── request_id.go
│   │   │   ├── logging.go
│   │   │   └── cors.go
│   │   └── router.go     # Route definitions
│   ├── service/          # Business logic layer
│   │   ├── user_service.go
│   │   ├── user_service_test.go
│   │   └── ...
│   ├── repository/       # Data access layer
│   │   ├── user_repository.go
│   │   ├── user_repository_test.go
│   │   └── ...
│   ├── models/           # Domain models
│   │   ├── user.go
│   │   ├── errors.go
│   │   └── ...
│   └── config/           # Configuration
│       └── config.go
├── migrations/           # Database migrations
│   ├── 001_create_users_table.up.sql
│   ├── 001_create_users_table.down.sql
│   └── ...
├── queries/              # sqlc query definitions
│   ├── users.sql
│   └── ...
├── db/                   # sqlc generated code (gitignored)
│   └── *.go
├── docs/                 # Documentation
│   ├── AGENTS.md
│   ├── SKILLS.md
│   ├── ARCHITECTURE.md
│   ├── CODING_STANDARDS.md
│   ├── DATABASE.md
│   ├── API_STANDARDS.md
│   ├── TESTING.md
│   └── SECURITY.md
├── scripts/              # Build/deployment scripts
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── sqlc.yaml
├── go.mod
└── go.sum
```

## Layered Architecture

### Overview
The application follows a clean, layered architecture with clear separation of concerns:

```
┌─────────────────────────────────────┐
│         HTTP Client                  │
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│      Handler Layer                   │  - Parse HTTP requests
│   (internal/api/handlers)            │  - Validate input
│                                      │  - Call services
│                                      │  - Format HTTP responses
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│      Service Layer                   │  - Business logic
│   (internal/service)                 │  - Orchestrate repositories
│                                      │  - Enforce business rules
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│    Repository Layer                  │  - Data access via sqlc
│   (internal/repository)              │  - Database operations
│                                      │  - Transaction management
└─────────────┬───────────────────────┘
              │
┌─────────────▼───────────────────────┐
│         Database                     │  - PostgreSQL
└─────────────────────────────────────┘
```

### Dependency Flow

**Rule**: Dependencies always point inward

```
Handlers    →    Services    →    Repositories    →    Database
(external)       (business)        (data access)        (storage)
```

- Handlers depend on Services
- Services depend on Repositories
- Repositories depend on Database (via sqlc)
- **Never**: Services calling handlers, Repositories calling services

### Cross-Cutting Concerns

```
┌─────────────────────────────────────────────────┐
│              Middleware Layer                    │
│  - Request ID generation                        │
│  - Logging                                      │
│  - Authentication/Authorization                 │
│  - CORS                                         │
│  - Rate Limiting                                │
│  - Error Recovery                               │
└─────────────────────────────────────────────────┘
           │                    │
           ▼                    ▼
      Handlers  ──────────►  Services  ──────────►  Repositories
```

## Layer Responsibilities

### 1. Handler Layer (`internal/api/handlers/`)

**Purpose**: HTTP request/response handling

**Responsibilities**:
- Parse HTTP requests (JSON, form data, URL params)
- Validate input using struct tags
- Extract authentication/authorization info
- Call appropriate service methods
- Format responses (JSON:API format)
- Set proper HTTP status codes
- Handle HTTP-specific errors

**What Handlers Should NOT Do**:
- ❌ Contain business logic
- ❌ Access database directly
- ❌ Call other handlers
- ❌ Perform complex data transformations

**Example Structure**:
```go
type UserHandler struct {
    userService service.UserService
    logger      *slog.Logger
    validator   *validator.Validate
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Get context and request ID
    ctx := r.Context()
    reqID := middleware.GetRequestID(ctx)
    
    // 2. Parse request
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, reqID, http.StatusBadRequest, "INVALID_JSON", err)
        return
    }
    
    // 3. Validate
    if err := h.validator.Struct(req); err != nil {
        respondValidationError(w, reqID, err)
        return
    }
    
    // 4. Call service
    user, err := h.userService.CreateUser(ctx, &req)
    if err != nil {
        handleServiceError(w, reqID, err)
        return
    }
    
    // 5. Respond with JSON:API format
    respondJSON(w, http.StatusCreated, JSONAPIResponse{
        Data: JSONAPIData{
            Type: "users",
            ID: user.ID,
            Attributes: user,
        },
    })
}
```

### 2. Service Layer (`internal/service/`)

**Purpose**: Business logic and orchestration

**Responsibilities**:
- Implement business rules
- Orchestrate multiple repository calls
- Perform data transformations
- Validate business constraints
- Handle transactions (if needed)
- Log business events

**What Services Should NOT Do**:
- ❌ Parse HTTP requests
- ❌ Format HTTP responses
- ❌ Execute SQL queries directly
- ❌ Handle HTTP status codes

**Example Structure**:
```go
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id uuid.UUID) (*User, error)
}

type userService struct {
    userRepo repository.UserRepository
    logger   *slog.Logger
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error) {
    reqID := middleware.GetRequestID(ctx)
    
    // Business logic: Check if email already exists
    existing, err := s.userRepo.GetByEmail(ctx, req.Email)
    if err != nil && !errors.Is(err, ErrNotFound) {
        return nil, fmt.Errorf("check existing user: %w", err)
    }
    if existing != nil {
        return nil, ErrEmailAlreadyExists
    }
    
    // Business logic: Hash password
    hashedPassword, err := hashPassword(req.Password)
    if err != nil {
        return nil, fmt.Errorf("hash password: %w", err)
    }
    
    // Create user
    user := &User{
        ID:           uuid.New(),
        Email:        req.Email,
        Name:         req.Name,
        PasswordHash: hashedPassword,
    }
    
    created, err := s.userRepo.Create(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("create user: %w", err)
    }
    
    s.logger.InfoContext(ctx, "user created",
        slog.String("request_id", reqID),
        slog.String("user_id", created.ID.String()))
    
    return created, nil
}
```

### 3. Repository Layer (`internal/repository/`)

**Purpose**: Data access abstraction

**Responsibilities**:
- Define data access interfaces
- Implement using sqlc-generated code
- Handle database errors
- Manage database connections
- Execute queries and transactions

**What Repositories Should NOT Do**:
- ❌ Contain business logic
- ❌ Write manual SQL queries
- ❌ Validate business rules
- ❌ Call other repositories directly

**Example Structure**:
```go
// Interface
type UserRepository interface {
    Create(ctx context.Context, user *User) (*User, error)
    GetByID(ctx context.Context, id uuid.UUID) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    List(ctx context.Context, offset, limit int) ([]*User, error)
    Update(ctx context.Context, user *User) (*User, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

// Implementation using sqlc
type userRepository struct {
    queries *db.Queries
    logger  *slog.Logger
}

func (r *userRepository) Create(ctx context.Context, user *User) (*User, error) {
    reqID := middleware.GetRequestID(ctx)
    
    // Use sqlc-generated code (NEVER manual SQL)
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
    
    return &User{
        ID:        result.ID,
        Email:     result.Email,
        Name:      result.Name,
        CreatedAt: result.CreatedAt,
        UpdatedAt: result.UpdatedAt,
    }, nil
}
```

### 4. Models Layer (`internal/models/`)

**Purpose**: Domain entities and data transfer objects

**Responsibilities**:
- Define domain entities
- Define request/response DTOs
- Define custom error types
- Provide utility methods

**Example Structure**:
```go
// Domain entity
type User struct {
    ID           uuid.UUID `json:"id"`
    Email        string    `json:"email"`
    Name         string    `json:"name"`
    PasswordHash string    `json:"-"` // Never expose in JSON
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}

// Request DTO
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password" validate:"required,min=8"`
}

// Response DTO
type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
}

// Custom errors
var (
    ErrNotFound           = errors.New("not found")
    ErrEmailAlreadyExists = errors.New("email already exists")
    ErrUnauthorized       = errors.New("unauthorized")
)
```

## Dependency Injection

**Pattern**: Constructor injection using interfaces

### Benefits:
- Testability (easy to mock)
- Flexibility (swap implementations)
- Clear dependencies
- No global state

### Example:
```go
// main.go
func main() {
    // Initialize database
    conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    queries := db.New(conn)
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(queries, logger)
    postRepo := repository.NewPostRepository(queries, logger)
    
    // Initialize services (inject repositories)
    userService := service.NewUserService(userRepo, logger)
    postService := service.NewPostService(postRepo, userRepo, logger)
    
    // Initialize handlers (inject services)
    userHandler := handlers.NewUserHandler(userService, logger)
    postHandler := handlers.NewPostHandler(postService, logger)
    
    // Setup router
    router := setupRouter(userHandler, postHandler)
    
    // Start server
    http.ListenAndServe(":8080", router)
}
```

## Context Propagation

**Pattern**: Pass `context.Context` through all layers

### Context Contents:
- Request ID (unique per request)
- User authentication info
- Timeouts and deadlines
- Cancellation signals

### Implementation:
```go
// 1. Middleware adds request ID to context
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        reqID := r.Header.Get("X-Request-ID")
        if reqID == "" {
            reqID = uuid.New().String()
        }
        
        ctx := context.WithValue(r.Context(), requestIDKey, reqID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// 2. Extract in any layer
func GetRequestID(ctx context.Context) string {
    if reqID, ok := ctx.Value(requestIDKey).(string); ok {
        return reqID
    }
    return ""
}

// 3. Use in logs
logger.InfoContext(ctx, "processing request",
    slog.String("request_id", GetRequestID(ctx)))
```

## Error Handling Strategy

### Error Wrapping:
```go
// Repository layer
if err != nil {
    return nil, fmt.Errorf("create user: %w", err)
}

// Service layer
if err != nil {
    return nil, fmt.Errorf("failed to register user: %w", err)
}

// Handler layer receives full error chain
```

### Error Classification:
```go
// Custom error types for different scenarios
var (
    ErrNotFound      = errors.New("resource not found")
    ErrUnauthorized  = errors.New("unauthorized")
    ErrForbidden     = errors.New("forbidden")
    ErrValidation    = errors.New("validation failed")
    ErrConflict      = errors.New("resource conflict")
)

// In handler, map to HTTP status
func handleServiceError(w http.ResponseWriter, reqID string, err error) {
    switch {
    case errors.Is(err, ErrNotFound):
        respondError(w, reqID, http.StatusNotFound, "NOT_FOUND", err)
    case errors.Is(err, ErrUnauthorized):
        respondError(w, reqID, http.StatusUnauthorized, "UNAUTHORIZED", err)
    case errors.Is(err, ErrValidation):
        respondError(w, reqID, http.StatusBadRequest, "VALIDATION_ERROR", err)
    default:
        respondError(w, reqID, http.StatusInternalServerError, "INTERNAL_ERROR", err)
    }
}
```

## Transaction Management

**Pattern**: Handle transactions in service layer when needed

```go
func (s *userService) CreateUserWithProfile(ctx context.Context, req *CreateUserRequest) error {
    // Begin transaction
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback(ctx)
    
    queries := s.queries.WithTx(tx)
    
    // Create user
    user, err := queries.CreateUser(ctx, db.CreateUserParams{...})
    if err != nil {
        return fmt.Errorf("create user: %w", err)
    }
    
    // Create profile
    _, err = queries.CreateProfile(ctx, db.CreateProfileParams{
        UserID: user.ID,
        ...
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

## Testing Architecture

### Unit Tests (Per Layer):
```
Handler Tests → Mock Services
Service Tests → Mock Repositories  
Repository Tests → Mock Database (pgxmock) or Real DB (testcontainers)
```

### Integration Tests:
```
HTTP Request → Full Stack → Real Database
```

### Test Organization:
```
internal/
├── api/handlers/
│   ├── user_handler.go
│   └── user_handler_test.go        # Unit tests with mock service
├── service/
│   ├── user_service.go
│   └── user_service_test.go        # Unit tests with mock repository
└── repository/
    ├── user_repository.go
    └── user_repository_test.go      # Unit tests with pgxmock

tests/
└── integration/
    └── user_api_test.go              # Integration tests
```

## Configuration Management

**Pattern**: Environment-based configuration

```go
type Config struct {
    ServerAddress string `mapstructure:"SERVER_ADDRESS"`
    DatabaseURL   string `mapstructure:"DATABASE_URL"`
    RedisURL      string `mapstructure:"REDIS_URL"`
    JWTSecret     string `mapstructure:"JWT_SECRET"`
    LogLevel      string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (*Config, error) {
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    return &cfg, nil
}
```

## Application Startup Flow

```
1. Load Configuration
   └─> Read .env file
   └─> Validate required vars

2. Initialize Logger
   └─> Set log level
   └─> Configure structured logging

3. Connect to Database
   └─> Establish connection pool
   └─> Run migrations (optional)

4. Initialize Dependencies (bottom-up)
   └─> Create repositories
   └─> Create services
   └─> Create handlers

5. Setup Router
   └─> Register routes
   └─> Add middleware

6. Start HTTP Server
   └─> Listen on configured port
   └─> Handle graceful shutdown
```

## Graceful Shutdown

```go
func main() {
    server := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
    
    // Start server in goroutine
    go func() {
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exiting")
}
```
