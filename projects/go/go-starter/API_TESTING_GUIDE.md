# GET User by ID Endpoint - Testing Guide

## Endpoint Information

**URL**: `GET /api/v1/users/{id}`

**Description**: Retrieves user information by user ID

**Authentication**: Not required (public endpoint)

**Response Format**: JSON:API specification

---

## Quick Start

### 1. Start the Server

```bash
# Make sure your database is running
# Then start the server
go run cmd/server/main.go
```

### 2. Test the Endpoint

#### Using curl:

```bash
# Replace {user-id} with an actual UUID from your database
curl -X GET http://localhost:8080/api/v1/users/{user-id}
```

---

## Example Requests and Responses

### Success Response (200 OK)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

**Headers:**
```
Content-Type: application/vnd.api+json
Status: 200 OK
```

---

### User Not Found (404 Not Found)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/00000000-0000-0000-0000-000000000000
```

**Response:**
```json
{
  "errors": [
    {
      "status": "404",
      "code": "NOT_FOUND",
      "title": "Not Found",
      "detail": "User not found",
      "meta": {
        "request_id": "req_abc123def456"
      }
    }
  ]
}
```

**Headers:**
```
Content-Type: application/vnd.api+json
Status: 404 Not Found
```

---

### Invalid UUID Format (400 Bad Request)

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/invalid-uuid
```

**Response:**
```json
{
  "errors": [
    {
      "status": "400",
      "code": "INVALID_ID",
      "title": "Bad Request",
      "detail": "Invalid user ID format",
      "meta": {
        "request_id": "req_abc123def456"
      }
    }
  ]
}
```

**Headers:**
```
Content-Type: application/vnd.api+json
Status: 400 Bad Request
```

---

## Response Fields

The endpoint returns only minimal user information:

| Field | Type   | Description           |
|-------|--------|-----------------------|
| id    | string | User's UUID           |
| name  | string | User's full name      |
| email | string | User's email address  |

**Note**: Sensitive fields like `password_hash` are never returned.

---

## Testing with Real Data

### Step 1: Insert a Test User

First, you need to insert a test user into your database:

```bash
# Connect to your PostgreSQL database
psql postgres://postgres:postgres@localhost:5432/go_starter

# Insert a test user
INSERT INTO users (email, name, password_hash)
VALUES ('test@example.com', 'Test User', 'hashed_password_here')
RETURNING id;
```

Copy the returned UUID.

### Step 2: Test the Endpoint

```bash
# Use the UUID from the previous step
curl -X GET http://localhost:8080/api/v1/users/{UUID-from-step-1}
```

---

## Architecture Overview

The endpoint follows clean architecture principles:

```
HTTP Request
    ↓
Handler Layer (user_handler.go)
    ↓
Service Layer (user_service.go)
    ↓
Repository Layer (user_repository.go)
    ↓
Database (via sqlc)
```

### Files Created

1. **internal/repository/user_repository.go** - Database access layer
2. **internal/service/user_service.go** - Business logic layer
3. **internal/api/handlers/user_handler.go** - HTTP handler
4. **internal/api/handlers/response.go** - JSON:API helpers
5. **internal/api/handlers/user_dto.go** - Data transfer objects

### Files Modified

1. **internal/api/router.go** - Added user routes
2. **cmd/server/main.go** - Wired dependencies

---

## Error Codes

| Code            | HTTP Status | Description                    |
|-----------------|-------------|--------------------------------|
| INVALID_ID      | 400         | Invalid UUID format            |
| NOT_FOUND       | 404         | User does not exist            |
| INTERNAL_ERROR  | 500         | Unexpected server error        |

---

## JSON:API Specification

This endpoint follows the [JSON:API v1.1 specification](https://jsonapi.org/).

**Success responses** include:
- `data` object with `type`, `id`, and `attributes`

**Error responses** include:
- `errors` array with `status`, `code`, `title`, `detail`, and `meta`

---

## Next Steps

### Adding More Endpoints

You can easily add more user endpoints following the same pattern:

```go
// In internal/api/router.go
r.Route("/users", func(r chi.Router) {
    r.Get("/{id}", userHandler.GetUser)         // ✅ Implemented
    // r.Post("/", userHandler.CreateUser)      // TODO
    // r.Patch("/{id}", userHandler.UpdateUser) // TODO
    // r.Delete("/{id}", userHandler.DeleteUser)// TODO
    // r.Get("/", userHandler.ListUsers)        // TODO
})
```

### Adding Authentication

To add authentication:

1. Create authentication middleware
2. Apply to protected routes
3. Update handler to get user from context

```go
r.Route("/users", func(r chi.Router) {
    r.Use(middleware.RequireAuth) // Add authentication
    r.Get("/{id}", userHandler.GetUser)
})
```

---

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
pg_isready

# Test connection manually
psql postgres://postgres:postgres@localhost:5432/go_starter
```

### Server Won't Start

```bash
# Check if port 8080 is already in use
lsof -i :8080

# Kill the process if needed
kill -9 <PID>
```

### Build Errors

```bash
# Clean build cache
go clean -cache

# Download dependencies
go mod download

# Rebuild
go build -o bin/server ./cmd/server
```

---

## Health Check

Before testing the user endpoint, verify the server is running:

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"healthy"}
```
