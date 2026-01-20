# Security Best Practices

This document covers security guidelines and best practices for Go backend services.

## Table of Contents
1. [Authentication & Authorization](#authentication--authorization)
2. [Password Security](#password-security)
3. [Input Validation](#input-validation)
4. [SQL Injection Prevention](#sql-injection-prevention)
5. [Sensitive Data Handling](#sensitive-data-handling)
6. [HTTPS & TLS](#https--tls)
7. [CORS Configuration](#cors-configuration)
8. [Rate Limiting](#rate-limiting)
9. [Security Headers](#security-headers)
10. [Common Vulnerabilities](#common-vulnerabilities)

---

## Authentication & Authorization

### JWT Authentication

#### Token Generation
```go
import (
    "github.com/golang-jwt/jwt/v5"
    "time"
)

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID, email, role string, secret []byte) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "your-app-name",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secret)
}
```

#### Token Validation
```go
func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return secret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

#### Authentication Middleware
```go
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Extract token from header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                respondError(w, "", http.StatusUnauthorized, "UNAUTHORIZED", "Missing authorization header")
                return
            }
            
            // Extract token (Bearer <token>)
            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                respondError(w, "", http.StatusUnauthorized, "UNAUTHORIZED", "Invalid authorization format")
                return
            }
            
            tokenString := parts[1]
            
            // Validate token
            claims, err := ValidateToken(tokenString, secret)
            if err != nil {
                respondError(w, "", http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
                return
            }
            
            // Add claims to context
            ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
            ctx = context.WithValue(ctx, userRoleKey, claims.Role)
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### Authorization Middleware
```go
func RequireRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole := GetUserRole(r.Context())
            
            allowed := false
            for _, role := range roles {
                if userRole == role {
                    allowed = true
                    break
                }
            }
            
            if !allowed {
                respondError(w, "", http.StatusForbidden, "FORBIDDEN", "Insufficient permissions")
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Usage
router.Handle("/api/v1/admin/users", 
    AuthMiddleware(secret)(
        RequireRole("admin")(
            http.HandlerFunc(adminHandler.ListUsers))))
```

### Resource Ownership Check
```go
func (s *postService) UpdatePost(ctx context.Context, postID uuid.UUID, req *UpdatePostRequest) (*Post, error) {
    // Get requesting user
    userID := GetUserID(ctx)
    userRole := GetUserRole(ctx)
    
    // Get existing post
    post, err := s.repo.GetByID(ctx, postID)
    if err != nil {
        return nil, err
    }
    
    // Check ownership or admin
    if post.AuthorID != userID && userRole != "admin" {
        return nil, ErrForbidden
    }
    
    // Proceed with update
    return s.repo.Update(ctx, postID, req)
}
```

---

## Password Security

### Never Store Plain Text Passwords
```go
// ❌ NEVER DO THIS
user.Password = req.Password

// ✅ ALWAYS hash passwords
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    // Cost 12 is a good balance (can adjust based on security needs)
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func VerifyPassword(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

### Password Requirements
```go
func ValidatePassword(password string) error {
    if len(password) < 8 {
        return errors.New("password must be at least 8 characters")
    }
    
    hasUpper := false
    hasLower := false
    hasNumber := false
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        }
    }
    
    if !hasUpper || !hasLower || !hasNumber {
        return errors.New("password must contain uppercase, lowercase, and number")
    }
    
    return nil
}
```

### Using Argon2 (More Secure Alternative)
```go
import "golang.org/x/crypto/argon2"

type PasswordConfig struct {
    time    uint32
    memory  uint32
    threads uint8
    keyLen  uint32
    saltLen uint32
}

func DefaultPasswordConfig() *PasswordConfig {
    return &PasswordConfig{
        time:    1,
        memory:  64 * 1024,
        threads: 4,
        keyLen:  32,
        saltLen: 16,
    }
}

func HashPasswordArgon2(password string, config *PasswordConfig) (string, error) {
    salt := make([]byte, config.saltLen)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }
    
    hash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)
    
    // Encode as: $argon2id$salt$hash
    encoded := fmt.Sprintf("$argon2id$%s$%s",
        base64.RawStdEncoding.EncodeToString(salt),
        base64.RawStdEncoding.EncodeToString(hash))
    
    return encoded, nil
}
```

---

## Input Validation

### Validate All Inputs
```go
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"omitempty,gte=0,lte=150"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondError(w, "", http.StatusBadRequest, "INVALID_JSON", err.Error())
        return
    }
    
    // Validate
    validate := validator.New()
    if err := validate.Struct(req); err != nil {
        respondValidationError(w, "", err)
        return
    }
    
    // Additional custom validation
    if err := ValidatePassword(req.Password); err != nil {
        respondError(w, "", http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
        return
    }
    
    // Process request
    // ...
}
```

### Sanitize User Input
```go
import "html"

func SanitizeString(input string) string {
    // Remove leading/trailing whitespace
    sanitized := strings.TrimSpace(input)
    
    // Escape HTML to prevent XSS
    sanitized = html.EscapeString(sanitized)
    
    return sanitized
}

// Usage
user.Name = SanitizeString(req.Name)
```

### Validate UUIDs
```go
func ValidateUUID(id string) (uuid.UUID, error) {
    parsed, err := uuid.Parse(id)
    if err != nil {
        return uuid.Nil, errors.New("invalid UUID format")
    }
    return parsed, nil
}
```

### Validate File Uploads
```go
func ValidateFileUpload(file multipart.File, header *multipart.FileHeader) error {
    // Check file size (10MB limit)
    if header.Size > 10*1024*1024 {
        return errors.New("file size exceeds 10MB limit")
    }
    
    // Check file type
    buffer := make([]byte, 512)
    if _, err := file.Read(buffer); err != nil {
        return err
    }
    file.Seek(0, 0) // Reset file pointer
    
    contentType := http.DetectContentType(buffer)
    allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
    
    allowed := false
    for _, t := range allowedTypes {
        if contentType == t {
            allowed = true
            break
        }
    }
    
    if !allowed {
        return fmt.Errorf("file type %s not allowed", contentType)
    }
    
    return nil
}
```

---

## SQL Injection Prevention

### ALWAYS Use sqlc (Never Manual SQL)
```go
// ✅ CORRECT: Using sqlc-generated code
result, err := r.queries.GetUserByEmail(ctx, email)

// ❌ WRONG: Manual SQL (vulnerable to injection)
query := "SELECT * FROM users WHERE email = '" + email + "'"
db.Query(query)

// ❌ WRONG: Even with placeholders, use sqlc
query := "SELECT * FROM users WHERE email = $1"
db.Query(query, email)
```

### Why sqlc Prevents SQL Injection
- Type-safe queries at compile time
- Automatic parameterization
- No string concatenation
- Review queries in .sql files
- Catches errors before runtime

---

## Sensitive Data Handling

### Never Log Sensitive Data
```go
// ❌ NEVER log passwords, tokens, or PII
logger.Info("user login", 
    slog.String("email", email),
    slog.String("password", password)) // NEVER!

logger.Info("request authenticated",
    slog.String("token", token)) // NEVER!

// ✅ Log only safe information
logger.InfoContext(ctx, "user login attempt",
    slog.String("request_id", reqID),
    slog.String("email", email))

logger.InfoContext(ctx, "authentication successful",
    slog.String("request_id", reqID),
    slog.String("user_id", userID))
```

### Never Return Sensitive Data
```go
// ✅ Exclude sensitive fields from JSON
type User struct {
    ID           uuid.UUID `json:"id"`
    Email        string    `json:"email"`
    Name         string    `json:"name"`
    PasswordHash string    `json:"-"` // Never expose
    CreatedAt    time.Time `json:"created_at"`
}
```

### Use Environment Variables for Secrets
```go
// ❌ NEVER hardcode secrets
const jwtSecret = "my-secret-key"
const apiKey = "abc123"

// ✅ Load from environment
type Config struct {
    JWTSecret    string `mapstructure:"JWT_SECRET"`
    DatabaseURL  string `mapstructure:"DATABASE_URL"`
    APIKey       string `mapstructure:"API_KEY"`
}

func LoadConfig() (*Config, error) {
    viper.AutomaticEnv()
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    
    // Validate required secrets
    if cfg.JWTSecret == "" {
        return nil, errors.New("JWT_SECRET is required")
    }
    
    return &cfg, nil
}
```

### Redact Sensitive Data in Logs
```go
func RedactEmail(email string) string {
    parts := strings.Split(email, "@")
    if len(parts) != 2 {
        return "***"
    }
    
    name := parts[0]
    if len(name) > 2 {
        name = name[:2] + "***"
    } else {
        name = "***"
    }
    
    return name + "@" + parts[1]
}

// Usage
logger.Info("password reset requested",
    slog.String("email", RedactEmail(email)))
// Output: "jo***@example.com"
```

---

## HTTPS & TLS

### Enforce HTTPS
```go
func RedirectToHTTPS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("X-Forwarded-Proto") != "https" {
            url := "https://" + r.Host + r.URL.String()
            http.Redirect(w, r, url, http.StatusMovedPermanently)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### TLS Configuration
```go
func setupTLSServer(addr string, handler http.Handler) *http.Server {
    tlsConfig := &tls.Config{
        MinVersion:               tls.VersionTLS12,
        CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
        PreferServerCipherSuites: true,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
        },
    }
    
    return &http.Server{
        Addr:         addr,
        Handler:      handler,
        TLSConfig:    tlsConfig,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
}
```

---

## CORS Configuration

### Secure CORS Setup
```go
import "github.com/rs/cors"

func setupCORS() *cors.Cors {
    return cors.New(cors.Options{
        AllowedOrigins: []string{
            "https://yourdomain.com",
            "https://app.yourdomain.com",
        },
        AllowedMethods: []string{
            http.MethodGet,
            http.MethodPost,
            http.MethodPut,
            http.MethodPatch,
            http.MethodDelete,
            http.MethodOptions,
        },
        AllowedHeaders: []string{
            "Accept",
            "Authorization",
            "Content-Type",
            "X-Request-ID",
        },
        ExposedHeaders: []string{
            "X-Request-ID",
        },
        AllowCredentials: true,
        MaxAge:           300, // 5 minutes
    })
}

// ❌ NEVER use wildcard in production
AllowedOrigins: []string{"*"} // Dangerous!
```

---

## Rate Limiting

### Per-IP Rate Limiting
```go
import "golang.org/x/time/rate"

type IPRateLimiter struct {
    ips map[string]*rate.Limiter
    mu  *sync.RWMutex
    r   rate.Limit
    b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    return &IPRateLimiter{
        ips: make(map[string]*rate.Limiter),
        mu:  &sync.RWMutex{},
        r:   r,
        b:   b,
    }
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()
    
    limiter, exists := i.ips[ip]
    if !exists {
        limiter = rate.NewLimiter(i.r, i.b)
        i.ips[ip] = limiter
    }
    
    return limiter
}

func RateLimitMiddleware(limiter *IPRateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := getIP(r)
            
            if !limiter.GetLimiter(ip).Allow() {
                respondError(w, "", http.StatusTooManyRequests, 
                    "RATE_LIMIT_EXCEEDED", "Too many requests")
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func getIP(r *http.Request) string {
    // Check X-Forwarded-For header
    forwarded := r.Header.Get("X-Forwarded-For")
    if forwarded != "" {
        return strings.Split(forwarded, ",")[0]
    }
    
    // Check X-Real-IP header
    realIP := r.Header.Get("X-Real-IP")
    if realIP != "" {
        return realIP
    }
    
    // Fallback to RemoteAddr
    ip, _, _ := net.SplitHostPort(r.RemoteAddr)
    return ip
}
```

### Per-User Rate Limiting
```go
func RateLimitByUser(next http.Handler) http.Handler {
    limiters := make(map[string]*rate.Limiter)
    mu := &sync.RWMutex{}
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := GetUserID(r.Context())
        
        mu.Lock()
        limiter, exists := limiters[userID]
        if !exists {
            limiter = rate.NewLimiter(rate.Every(time.Minute), 100) // 100 requests per minute
            limiters[userID] = limiter
        }
        mu.Unlock()
        
        if !limiter.Allow() {
            respondError(w, "", http.StatusTooManyRequests, 
                "RATE_LIMIT_EXCEEDED", "Too many requests")
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## Security Headers

### Add Security Headers
```go
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Prevent clickjacking
        w.Header().Set("X-Frame-Options", "DENY")
        
        // Prevent MIME sniffing
        w.Header().Set("X-Content-Type-Options", "nosniff")
        
        // Enable XSS protection
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        
        // Referrer policy
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Content Security Policy
        w.Header().Set("Content-Security-Policy", 
            "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
        
        // HSTS (only if using HTTPS)
        if r.TLS != nil {
            w.Header().Set("Strict-Transport-Security", 
                "max-age=31536000; includeSubDomains; preload")
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## Common Vulnerabilities

### Prevent Timing Attacks
```go
import "crypto/subtle"

// ✅ Use constant-time comparison for sensitive data
func CompareTokens(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// ❌ Don't use regular comparison
if token == expectedToken { } // Vulnerable to timing attack
```

### Prevent Path Traversal
```go
func SecureFilePath(basePath, userPath string) (string, error) {
    // Join paths
    fullPath := filepath.Join(basePath, userPath)
    
    // Clean path
    cleanPath := filepath.Clean(fullPath)
    
    // Verify path is within base directory
    if !strings.HasPrefix(cleanPath, basePath) {
        return "", errors.New("invalid file path")
    }
    
    return cleanPath, nil
}
```

### Prevent XXE (XML External Entity)
```go
import "encoding/xml"

func ParseXML(data []byte, v interface{}) error {
    decoder := xml.NewDecoder(bytes.NewReader(data))
    
    // Disable external entities
    decoder.Strict = false
    decoder.Entity = xml.HTMLEntity
    
    return decoder.Decode(v)
}
```

---

## Checklist

Security review before deployment:

- [ ] All passwords hashed with bcrypt/argon2
- [ ] JWT tokens validated properly
- [ ] Input validation on all endpoints
- [ ] sqlc used for all database queries (no manual SQL)
- [ ] Sensitive data not logged
- [ ] Sensitive fields excluded from JSON responses
- [ ] Secrets loaded from environment variables
- [ ] HTTPS enforced
- [ ] CORS configured with specific origins
- [ ] Rate limiting implemented
- [ ] Security headers added
- [ ] Authentication required on protected routes
- [ ] Authorization checks for resource access
- [ ] No timing attack vulnerabilities
- [ ] File upload validation implemented
- [ ] Path traversal prevention in file operations
