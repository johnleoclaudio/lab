# API Standards

This document defines REST API design principles and JSON:API response format standards.

## Table of Contents
1. [REST Principles](#rest-principles)
2. [JSON:API Format](#jsonapi-format)
3. [HTTP Methods](#http-methods)
4. [Status Codes](#status-codes)
5. [Versioning](#versioning)
6. [Request/Response Examples](#requestresponse-examples)
7. [Error Handling](#error-handling)

---

## REST Principles

### Resource-Based URLs

```
✅ Good: Use nouns, not verbs
GET    /api/v1/users
POST   /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id

❌ Bad: Verbs in URL
GET    /api/v1/getUsers
POST   /api/v1/createUser
GET    /api/v1/getUserById/:id
```

### Nested Resources

```
✅ Good: Show relationships
GET    /api/v1/users/:user_id/posts
POST   /api/v1/users/:user_id/posts
GET    /api/v1/posts/:post_id/comments

⚠️  Limit nesting to 2 levels:
GET    /api/v1/users/:user_id/posts/:post_id
❌  /api/v1/users/:user_id/posts/:post_id/comments/:comment_id/replies
```

### Collection Filtering

```
GET /api/v1/posts?status=published
GET /api/v1/posts?author_id=123
GET /api/v1/posts?tag=golang&status=published
```

### Pagination

```
GET /api/v1/posts?page=2&per_page=20
GET /api/v1/posts?offset=40&limit=20
```

### Sorting

```
GET /api/v1/posts?sort=created_at
GET /api/v1/posts?sort=-created_at  # Descending
GET /api/v1/posts?sort=author_name,created_at
```

### Field Selection

```
GET /api/v1/users?fields=id,name,email
GET /api/v1/posts?fields=id,title
```

---

## JSON:API Format

We follow the [JSON:API specification](https://jsonapi.org/) v1.1.

### Why JSON:API?

- ✅ Standardized structure
- ✅ Built-in error handling
- ✅ Relationship handling
- ✅ Pagination support
- ✅ Widely adopted

### Basic Structure

#### Single Resource Response

```json
{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

#### Collection Response

```json
{
  "data": [
    {
      "type": "users",
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "attributes": {
        "name": "John Doe",
        "email": "john@example.com"
      }
    },
    {
      "type": "users",
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "attributes": {
        "name": "Jane Smith",
        "email": "jane@example.com"
      }
    }
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20
  }
}
```

### With Relationships

```json
{
  "data": {
    "type": "posts",
    "id": "abc123",
    "attributes": {
      "title": "My Post",
      "content": "Post content here"
    },
    "relationships": {
      "author": {
        "data": {
          "type": "users",
          "id": "user123"
        }
      },
      "comments": {
        "data": [
          {"type": "comments", "id": "comment1"},
          {"type": "comments", "id": "comment2"}
        ]
      }
    }
  },
  "included": [
    {
      "type": "users",
      "id": "user123",
      "attributes": {
        "name": "John Doe"
      }
    },
    {
      "type": "comments",
      "id": "comment1",
      "attributes": {
        "content": "Great post!"
      }
    }
  ]
}
```

### Error Response

```json
{
  "errors": [
    {
      "status": "400",
      "code": "VALIDATION_ERROR",
      "title": "Invalid input",
      "detail": "The 'email' field must be a valid email address",
      "source": {
        "pointer": "/data/attributes/email"
      }
    }
  ]
}
```

---

## HTTP Methods

### GET - Retrieve Resources

```
GET /api/v1/users          # List users
GET /api/v1/users/:id      # Get single user

Response: 200 OK
```

**Rules**:
- Must be idempotent (same result on multiple calls)
- No request body
- Should not modify data

### POST - Create Resource

```
POST /api/v1/users

Request:
{
  "data": {
    "type": "users",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com",
      "password": "securepassword"
    }
  }
}

Response: 201 Created
Location: /api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**Rules**:
- Not idempotent (creates new resource each time)
- Returns 201 Created with Location header
- Returns created resource in response body

### PUT - Replace Resource

```
PUT /api/v1/users/:id

Request:
{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe Updated",
      "email": "john.new@example.com"
    }
  }
}

Response: 200 OK
```

**Rules**:
- Idempotent (same result on multiple calls)
- Replaces entire resource
- All fields should be provided

### PATCH - Partial Update

```
PATCH /api/v1/users/:id

Request:
{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe Updated"
    }
  }
}

Response: 200 OK
```

**Rules**:
- Idempotent
- Updates only provided fields
- Preferred over PUT for partial updates

### DELETE - Remove Resource

```
DELETE /api/v1/users/:id

Response: 204 No Content
```

**Rules**:
- Idempotent
- Returns 204 No Content (no body)
- May return 200 OK with deleted resource

---

## Status Codes

### 2xx Success

```
200 OK              - Successful GET, PUT, PATCH, DELETE
201 Created         - Successful POST (resource created)
204 No Content      - Successful DELETE (no response body)
```

### 3xx Redirection

```
301 Moved Permanently  - Resource permanently moved
302 Found              - Temporary redirect
304 Not Modified       - Cached version is still valid
```

### 4xx Client Errors

```
400 Bad Request        - Invalid request format/syntax
401 Unauthorized       - Authentication required
403 Forbidden          - Authenticated but not authorized
404 Not Found          - Resource doesn't exist
405 Method Not Allowed - HTTP method not supported
409 Conflict           - Resource conflict (e.g., duplicate email)
422 Unprocessable Entity - Validation errors
429 Too Many Requests  - Rate limit exceeded
```

### 5xx Server Errors

```
500 Internal Server Error - Unexpected server error
502 Bad Gateway           - Invalid response from upstream
503 Service Unavailable   - Server temporarily unavailable
504 Gateway Timeout       - Upstream server timeout
```

---

## Versioning

### URL Path Versioning (Preferred)

```
✅ Good: Version in path
/api/v1/users
/api/v2/users

✅ Clear and explicit
✅ Easy to route
✅ Easy to cache
```

### Version Strategy

```
v1 - Initial version
v2 - Breaking changes (e.g., field renamed, removed)
v1.1 - Non-breaking changes (new optional fields)
```

### Deprecation

```
# Header to warn about deprecation
Deprecation: true
Sunset: Sat, 31 Dec 2024 23:59:59 GMT
Link: </api/v2/users>; rel="successor-version"
```

---

## Request/Response Examples

### Create User

**Request**:
```http
POST /api/v1/users HTTP/1.1
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json

{
  "data": {
    "type": "users",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com",
      "password": "securepassword123"
    }
  }
}
```

**Response**:
```http
HTTP/1.1 201 Created
Content-Type: application/vnd.api+json
Location: /api/v1/users/550e8400-e29b-41d4-a716-446655440000

{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### Get User

**Request**:
```http
GET /api/v1/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Accept: application/vnd.api+json
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/vnd.api+json

{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    }
  }
}
```

### List Users with Pagination

**Request**:
```http
GET /api/v1/users?page=2&per_page=20 HTTP/1.1
Accept: application/vnd.api+json
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/vnd.api+json

{
  "data": [
    {
      "type": "users",
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "attributes": {
        "name": "John Doe",
        "email": "john@example.com"
      }
    }
  ],
  "meta": {
    "total": 100,
    "page": 2,
    "per_page": 20,
    "total_pages": 5
  },
  "links": {
    "self": "/api/v1/users?page=2&per_page=20",
    "first": "/api/v1/users?page=1&per_page=20",
    "prev": "/api/v1/users?page=1&per_page=20",
    "next": "/api/v1/users?page=3&per_page=20",
    "last": "/api/v1/users?page=5&per_page=20"
  }
}
```

### Update User

**Request**:
```http
PATCH /api/v1/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
Content-Type: application/vnd.api+json

{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe Updated"
    }
  }
}
```

**Response**:
```http
HTTP/1.1 200 OK
Content-Type: application/vnd.api+json

{
  "data": {
    "type": "users",
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "attributes": {
      "name": "John Doe Updated",
      "email": "john@example.com",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T11:45:00Z"
    }
  }
}
```

### Delete User

**Request**:
```http
DELETE /api/v1/users/550e8400-e29b-41d4-a716-446655440000 HTTP/1.1
```

**Response**:
```http
HTTP/1.1 204 No Content
```

---

## Error Handling

### Error Response Structure

```json
{
  "errors": [
    {
      "status": "400",
      "code": "ERROR_CODE",
      "title": "Short error title",
      "detail": "Detailed error message",
      "source": {
        "pointer": "/data/attributes/field_name"
      }
    }
  ]
}
```

### Common Error Codes

```go
const (
    ErrCodeValidation      = "VALIDATION_ERROR"
    ErrCodeNotFound        = "NOT_FOUND"
    ErrCodeUnauthorized    = "UNAUTHORIZED"
    ErrCodeForbidden       = "FORBIDDEN"
    ErrCodeConflict        = "CONFLICT"
    ErrCodeInternal        = "INTERNAL_ERROR"
    ErrCodeRateLimit       = "RATE_LIMIT_EXCEEDED"
    ErrCodeInvalidJSON     = "INVALID_JSON"
)
```

### Validation Errors

```json
{
  "errors": [
    {
      "status": "422",
      "code": "VALIDATION_ERROR",
      "title": "Validation failed",
      "detail": "The 'email' field must be a valid email address",
      "source": {
        "pointer": "/data/attributes/email"
      }
    },
    {
      "status": "422",
      "code": "VALIDATION_ERROR",
      "title": "Validation failed",
      "detail": "The 'password' field must be at least 8 characters",
      "source": {
        "pointer": "/data/attributes/password"
      }
    }
  ]
}
```

### Not Found Error

```json
{
  "errors": [
    {
      "status": "404",
      "code": "NOT_FOUND",
      "title": "Resource not found",
      "detail": "User with ID '550e8400-e29b-41d4-a716-446655440000' does not exist"
    }
  ]
}
```

### Unauthorized Error

```json
{
  "errors": [
    {
      "status": "401",
      "code": "UNAUTHORIZED",
      "title": "Authentication required",
      "detail": "You must provide a valid authentication token"
    }
  ]
}
```

### Conflict Error

```json
{
  "errors": [
    {
      "status": "409",
      "code": "CONFLICT",
      "title": "Resource conflict",
      "detail": "A user with email 'john@example.com' already exists"
    }
  ]
}
```

### Rate Limit Error

```json
{
  "errors": [
    {
      "status": "429",
      "code": "RATE_LIMIT_EXCEEDED",
      "title": "Too many requests",
      "detail": "You have exceeded the rate limit. Please try again in 60 seconds"
    }
  ]
}
```

### Internal Server Error

```json
{
  "errors": [
    {
      "status": "500",
      "code": "INTERNAL_ERROR",
      "title": "Internal server error",
      "detail": "An unexpected error occurred. Please try again later",
      "meta": {
        "request_id": "req_abc123def456"
      }
    }
  ]
}
```

### Implementation in Go

```go
type JSONAPIError struct {
    Status string                 `json:"status"`
    Code   string                 `json:"code"`
    Title  string                 `json:"title"`
    Detail string                 `json:"detail"`
    Source *JSONAPIErrorSource    `json:"source,omitempty"`
    Meta   map[string]interface{} `json:"meta,omitempty"`
}

type JSONAPIErrorSource struct {
    Pointer string `json:"pointer,omitempty"`
}

type JSONAPIErrorResponse struct {
    Errors []JSONAPIError `json:"errors"`
}

func respondError(w http.ResponseWriter, reqID string, status int, code, detail string) {
    w.Header().Set("Content-Type", "application/vnd.api+json")
    w.WriteHeader(status)
    
    response := JSONAPIErrorResponse{
        Errors: []JSONAPIError{
            {
                Status: fmt.Sprintf("%d", status),
                Code:   code,
                Title:  http.StatusText(status),
                Detail: detail,
                Meta: map[string]interface{}{
                    "request_id": reqID,
                },
            },
        },
    }
    
    json.NewEncoder(w).Encode(response)
}

func respondValidationError(w http.ResponseWriter, reqID string, err error) {
    w.Header().Set("Content-Type", "application/vnd.api+json")
    w.WriteHeader(http.StatusUnprocessableEntity)
    
    var errors []JSONAPIError
    
    // Parse validation errors
    if validationErrs, ok := err.(validator.ValidationErrors); ok {
        for _, e := range validationErrs {
            errors = append(errors, JSONAPIError{
                Status: "422",
                Code:   "VALIDATION_ERROR",
                Title:  "Validation failed",
                Detail: fmt.Sprintf("The '%s' field %s", e.Field(), e.Tag()),
                Source: &JSONAPIErrorSource{
                    Pointer: fmt.Sprintf("/data/attributes/%s", e.Field()),
                },
            })
        }
    }
    
    response := JSONAPIErrorResponse{Errors: errors}
    json.NewEncoder(w).Encode(response)
}
```

---

## Headers

### Request Headers

```http
Content-Type: application/vnd.api+json
Accept: application/vnd.api+json
Authorization: Bearer <token>
X-Request-ID: <uuid>  # Optional, generated if not provided
```

### Response Headers

```http
Content-Type: application/vnd.api+json
X-Request-ID: <uuid>
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1640000000
```

---

## Checklist

Before deploying API endpoints:

- [ ] URL follows REST principles (nouns, not verbs)
- [ ] Proper HTTP method used (GET/POST/PUT/PATCH/DELETE)
- [ ] Response follows JSON:API format
- [ ] Appropriate status code returned
- [ ] Errors follow JSON:API error format
- [ ] Request ID included in all responses
- [ ] Validation errors list all fields
- [ ] API versioned in URL path (/api/v1/)
- [ ] Authentication/authorization applied
- [ ] Rate limiting configured
- [ ] Documentation updated (OpenAPI/Swagger)
