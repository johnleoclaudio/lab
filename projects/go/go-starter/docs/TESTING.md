# Testing Guidelines

This document covers testing strategies, patterns, and best practices for Go backend services.

## Table of Contents
1. [Testing Philosophy](#testing-philosophy)
2. [Test Types](#test-types)
3. [Unit Testing](#unit-testing)
4. [Integration Testing](#integration-testing)
5. [Test Coverage](#test-coverage)
6. [Mocking](#mocking)
7. [Common Patterns](#common-patterns)

---

## Testing Philosophy

### Test-Driven Development (TDD)

**Required workflow**: Write tests BEFORE implementation

```
1. RED   - Write failing test
2. GREEN - Write minimum code to pass
3. REFACTOR - Improve code while keeping tests green
```

### Testing Pyramid

```
        /\
       /  \  E2E (Few)
      /----\
     /      \  Integration (Some)
    /--------\
   /          \  Unit (Many)
  /____________\
```

- **Unit Tests**: 70% - Fast, isolated, test single functions
- **Integration Tests**: 20% - Test components working together
- **E2E Tests**: 10% - Test complete user flows

### Coverage Goals

```
Overall Project:     >75%
Business Logic:      >80%
Handlers:            >70%
Repositories:        >70%
```

---

## Test Types

### Unit Tests
- Test single function/method
- Mock all dependencies
- Fast (<10ms per test)
- No database, no network

### Integration Tests
- Test multiple components together
- Real database (testcontainers)
- Slower (100ms-1s per test)
- Test actual workflows

### End-to-End Tests
- Test complete API flows
- Real server, real database
- Slowest (1s+ per test)
- User perspective

---

## Unit Testing

### File Naming
```
user_service.go       -> user_service_test.go
user_repository.go    -> user_repository_test.go
user_handler.go       -> user_handler_test.go
```

### Table-Driven Tests Pattern

```go
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   *CreateUserRequest
        want    *User
        wantErr bool
        setup   func(*mock.UserRepository)
    }{
        {
            name: "success - valid user",
            input: &CreateUserRequest{
                Name:     "John Doe",
                Email:    "john@example.com",
                Password: "password123",
            },
            want: &User{
                ID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: false,
            setup: func(m *mock.UserRepository) {
                m.EXPECT().
                    GetByEmail(gomock.Any(), "john@example.com").
                    Return(nil, models.ErrNotFound)
                
                m.EXPECT().
                    Create(gomock.Any(), gomock.Any()).
                    Return(&User{
                        ID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
                        Name:  "John Doe",
                        Email: "john@example.com",
                    }, nil)
            },
        },
        {
            name: "error - email already exists",
            input: &CreateUserRequest{
                Name:     "John Doe",
                Email:    "existing@example.com",
                Password: "password123",
            },
            want:    nil,
            wantErr: true,
            setup: func(m *mock.UserRepository) {
                m.EXPECT().
                    GetByEmail(gomock.Any(), "existing@example.com").
                    Return(&User{ID: uuid.New()}, nil)
            },
        },
        {
            name: "error - repository failure",
            input: &CreateUserRequest{
                Name:     "John Doe",
                Email:    "john@example.com",
                Password: "password123",
            },
            want:    nil,
            wantErr: true,
            setup: func(m *mock.UserRepository) {
                m.EXPECT().
                    GetByEmail(gomock.Any(), gomock.Any()).
                    Return(nil, errors.New("database error"))
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := mock.NewUserRepository(ctrl)
            if tt.setup != nil {
                tt.setup(mockRepo)
            }
            
            service := NewUserService(mockRepo, slog.Default())
            
            // Execute
            got, err := service.CreateUser(context.Background(), tt.input)
            
            // Assert
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if !tt.wantErr {
                assert.Equal(t, tt.want.ID, got.ID)
                assert.Equal(t, tt.want.Name, got.Name)
                assert.Equal(t, tt.want.Email, got.Email)
            }
        })
    }
}
```

### Testing Handler Layer

```go
func TestUserHandler_CreateUser(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
        expectedBody   string
        setup          func(*mock.UserService)
    }{
        {
            name: "success",
            requestBody: `{
                "data": {
                    "type": "users",
                    "attributes": {
                        "name": "John Doe",
                        "email": "john@example.com",
                        "password": "password123"
                    }
                }
            }`,
            expectedStatus: http.StatusCreated,
            setup: func(m *mock.UserService) {
                m.EXPECT().
                    CreateUser(gomock.Any(), gomock.Any()).
                    Return(&User{
                        ID:    uuid.New(),
                        Name:  "John Doe",
                        Email: "john@example.com",
                    }, nil)
            },
        },
        {
            name:           "invalid JSON",
            requestBody:    `{invalid json}`,
            expectedStatus: http.StatusBadRequest,
            setup:          func(m *mock.UserService) {},
        },
        {
            name: "validation error",
            requestBody: `{
                "data": {
                    "type": "users",
                    "attributes": {
                        "name": "",
                        "email": "invalid-email"
                    }
                }
            }`,
            expectedStatus: http.StatusUnprocessableEntity,
            setup:          func(m *mock.UserService) {},
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockService := mock.NewUserService(ctrl)
            tt.setup(mockService)
            
            handler := NewUserHandler(mockService, slog.Default())
            
            // Create request
            req := httptest.NewRequest(http.MethodPost, "/api/v1/users", 
                strings.NewReader(tt.requestBody))
            req.Header.Set("Content-Type", "application/vnd.api+json")
            
            // Record response
            rr := httptest.NewRecorder()
            
            // Execute
            handler.CreateUser(rr, req)
            
            // Assert
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            if tt.expectedBody != "" {
                assert.Contains(t, rr.Body.String(), tt.expectedBody)
            }
        })
    }
}
```

### Testing Service Layer

```go
func TestUserService_GetUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  uuid.UUID
        want    *User
        wantErr error
        setup   func(*mock.UserRepository)
    }{
        {
            name:   "success",
            userID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
            want: &User{
                ID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: nil,
            setup: func(m *mock.UserRepository) {
                m.EXPECT().
                    GetByID(gomock.Any(), uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")).
                    Return(&User{
                        ID:    uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
                        Name:  "John Doe",
                        Email: "john@example.com",
                    }, nil)
            },
        },
        {
            name:    "not found",
            userID:  uuid.New(),
            want:    nil,
            wantErr: models.ErrNotFound,
            setup: func(m *mock.UserRepository) {
                m.EXPECT().
                    GetByID(gomock.Any(), gomock.Any()).
                    Return(nil, models.ErrNotFound)
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            
            mockRepo := mock.NewUserRepository(ctrl)
            tt.setup(mockRepo)
            
            service := NewUserService(mockRepo, slog.Default())
            
            got, err := service.GetUser(context.Background(), tt.userID)
            
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}
```

### Testing Repository Layer

```go
func TestUserRepository_Create(t *testing.T) {
    // Use pgxmock for database mocking
    mock, err := pgxmock.NewPool()
    require.NoError(t, err)
    defer mock.Close()
    
    queries := db.New(mock)
    repo := NewUserRepository(queries, slog.Default())
    
    user := &User{
        ID:           uuid.New(),
        Email:        "john@example.com",
        Name:         "John Doe",
        PasswordHash: "hashedpassword",
    }
    
    // Expect CreateUser query
    rows := pgxmock.NewRows([]string{"id", "email", "name", "created_at", "updated_at"}).
        AddRow(user.ID, user.Email, user.Name, time.Now(), time.Now())
    
    mock.ExpectQuery("INSERT INTO users").
        WithArgs(user.ID, user.Email, user.Name, user.PasswordHash).
        WillReturnRows(rows)
    
    // Execute
    result, err := repo.Create(context.Background(), user)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, user.ID, result.ID)
    assert.Equal(t, user.Email, result.Email)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

## Integration Testing

### Using Testcontainers

```go
func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }
    
    ctx := context.Background()
    
    // Start PostgreSQL container
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("postgres"),
        postgres.WithPassword("postgres"),
        testcontainers.WithWaitStrategy(
            wait.ForLog("database system is ready to accept connections").
                WithOccurrence(2).
                WithStartupTimeout(5*time.Second)),
    )
    require.NoError(t, err)
    defer container.Terminate(ctx)
    
    // Get connection string
    connStr, err := container.ConnectionString(ctx, "sslmode=disable")
    require.NoError(t, err)
    
    // Connect to database
    pool, err := pgxpool.New(ctx, connStr)
    require.NoError(t, err)
    defer pool.Close()
    
    // Run migrations
    err = runMigrations(connStr)
    require.NoError(t, err)
    
    // Create repository
    queries := db.New(pool)
    repo := NewUserRepository(queries, slog.Default())
    
    t.Run("Create and Get User", func(t *testing.T) {
        user := &User{
            ID:           uuid.New(),
            Email:        "test@example.com",
            Name:         "Test User",
            PasswordHash: "hashedpassword",
        }
        
        // Create
        created, err := repo.Create(ctx, user)
        require.NoError(t, err)
        assert.Equal(t, user.Email, created.Email)
        
        // Get
        retrieved, err := repo.GetByID(ctx, created.ID)
        require.NoError(t, err)
        assert.Equal(t, created.ID, retrieved.ID)
        assert.Equal(t, created.Email, retrieved.Email)
    })
    
    t.Run("Duplicate Email", func(t *testing.T) {
        user1 := &User{
            ID:           uuid.New(),
            Email:        "duplicate@example.com",
            Name:         "User 1",
            PasswordHash: "hash1",
        }
        
        _, err := repo.Create(ctx, user1)
        require.NoError(t, err)
        
        user2 := &User{
            ID:           uuid.New(),
            Email:        "duplicate@example.com",
            Name:         "User 2",
            PasswordHash: "hash2",
        }
        
        _, err = repo.Create(ctx, user2)
        require.Error(t, err) // Should fail on unique constraint
    })
}
```

### End-to-End API Tests

```go
func TestAPI_UserFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping E2E test")
    }
    
    // Setup test server
    server := setupTestServer(t)
    defer server.Close()
    
    client := server.Client()
    baseURL := server.URL + "/api/v1"
    
    t.Run("Complete User Lifecycle", func(t *testing.T) {
        // 1. Register user
        registerBody := `{
            "data": {
                "type": "users",
                "attributes": {
                    "name": "Test User",
                    "email": "test@example.com",
                    "password": "password123"
                }
            }
        }`
        
        resp := makeRequest(t, client, "POST", baseURL+"/users", registerBody)
        assert.Equal(t, http.StatusCreated, resp.StatusCode)
        
        var createResp struct {
            Data struct {
                ID         string `json:"id"`
                Attributes struct {
                    Email string `json:"email"`
                } `json:"attributes"`
            } `json:"data"`
        }
        json.NewDecoder(resp.Body).Decode(&createResp)
        userID := createResp.Data.ID
        
        // 2. Get user
        resp = makeRequest(t, client, "GET", baseURL+"/users/"+userID, "")
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        // 3. Update user
        updateBody := `{
            "data": {
                "type": "users",
                "id": "` + userID + `",
                "attributes": {
                    "name": "Updated Name"
                }
            }
        }`
        
        resp = makeRequest(t, client, "PATCH", baseURL+"/users/"+userID, updateBody)
        assert.Equal(t, http.StatusOK, resp.StatusCode)
        
        // 4. Delete user
        resp = makeRequest(t, client, "DELETE", baseURL+"/users/"+userID, "")
        assert.Equal(t, http.StatusNoContent, resp.StatusCode)
        
        // 5. Verify deleted
        resp = makeRequest(t, client, "GET", baseURL+"/users/"+userID, "")
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
    })
}

func makeRequest(t *testing.T, client *http.Client, method, url, body string) *http.Response {
    var reqBody io.Reader
    if body != "" {
        reqBody = strings.NewReader(body)
    }
    
    req, err := http.NewRequest(method, url, reqBody)
    require.NoError(t, err)
    
    if body != "" {
        req.Header.Set("Content-Type", "application/vnd.api+json")
    }
    req.Header.Set("Accept", "application/vnd.api+json")
    
    resp, err := client.Do(req)
    require.NoError(t, err)
    
    return resp
}
```

---

## Test Coverage

### Running Tests with Coverage

```bash
# Run all tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Check coverage by package
go test -coverprofile=coverage.out ./internal/service
go tool cover -func=coverage.out
```

### Coverage Requirements

```bash
# Fail if coverage below threshold
go test -cover -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' | \
awk '{if ($1 < 75) exit 1}'
```

### What to Cover

✅ **Must Cover**:
- All business logic paths
- Error handling paths
- Edge cases (nil, empty, boundary values)
- Validation logic

❌ **Can Skip**:
- Generated code (sqlc)
- Simple getters/setters
- Third-party library wrappers

---

## Mocking

### Using gomock

#### Generate Mocks

```bash
# Install mockgen
go install github.com/golang/mock/mockgen@latest

# Generate mocks
mockgen -source=internal/repository/user_repository.go \
    -destination=internal/repository/mock/user_repository_mock.go \
    -package=mock
```

#### Use in Tests

```go
func TestWithMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    
    mockRepo := mock.NewUserRepository(ctrl)
    
    // Set expectations
    mockRepo.EXPECT().
        GetByID(gomock.Any(), gomock.Any()).
        Return(&User{ID: uuid.New()}, nil).
        Times(1)
    
    // Use mock
    service := NewUserService(mockRepo, slog.Default())
    user, err := service.GetUser(context.Background(), uuid.New())
    
    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Using testify/mock

```go
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *User) (*User, error) {
    args := m.Called(ctx, user)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestWithTestify(t *testing.T) {
    mockRepo := new(MockUserRepository)
    
    mockRepo.On("Create", mock.Anything, mock.Anything).
        Return(&User{ID: uuid.New()}, nil)
    
    service := NewUserService(mockRepo, slog.Default())
    user, err := service.CreateUser(context.Background(), &CreateUserRequest{})
    
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
}
```

---

## Common Patterns

### Test Helpers

```go
// Helper to create test user
func createTestUser(t *testing.T) *User {
    t.Helper()
    return &User{
        ID:    uuid.New(),
        Name:  "Test User",
        Email: "test@example.com",
    }
}

// Helper to setup test context
func testContext() context.Context {
    ctx := context.Background()
    ctx = context.WithValue(ctx, requestIDKey, uuid.New().String())
    return ctx
}
```

### Table Test Subtests

```go
func TestValidation(t *testing.T) {
    tests := map[string]struct {
        input   string
        wantErr bool
    }{
        "valid email":   {"user@example.com", false},
        "missing @":     {"userexample.com", true},
        "missing domain": {"user@", true},
        "empty":         {"", true},
    }
    
    for name, tt := range tests {
        t.Run(name, func(t *testing.T) {
            err := ValidateEmail(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("got error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Parallel Tests

```go
func TestParallel(t *testing.T) {
    tests := []struct {
        name string
        // ...
    }{ /* ... */ }
    
    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Run in parallel
            // test implementation
        })
    }
}
```

---

## Checklist

Before committing code:

- [ ] All new code has tests
- [ ] Tests follow TDD (written before implementation)
- [ ] Table-driven tests used where appropriate
- [ ] Test coverage >80% for business logic
- [ ] Tests run fast (unit tests <10ms each)
- [ ] Integration tests use testcontainers
- [ ] Mocks used for dependencies
- [ ] Edge cases covered (nil, empty, boundary)
- [ ] Error paths tested
- [ ] Tests pass with `go test -race`
- [ ] Tests pass with `go test -short` (skip integration)
