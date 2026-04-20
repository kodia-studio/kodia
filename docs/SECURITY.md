# 🛡️ Security Best Practices

This guide covers security features built into Kodia and best practices for production deployments.

**Table of Contents:**
- [Security Features](#security-features)
- [Authentication](#authentication)
- [Authorization](#authorization)
- [Input Validation](#input-validation)
- [SQL Injection Prevention](#sql-injection-prevention)
- [CORS & CSRF Protection](#cors--csrf-protection)
- [Rate Limiting](#rate-limiting)
- [Secrets Management](#secrets-management)
- [HTTPS & TLS](#https--tls)
- [Security Headers](#security-headers)
- [Audit Logging](#audit-logging)
- [Deployment Checklist](#deployment-checklist)
- [Vulnerability Reporting](#vulnerability-reporting)

---

## Security Features

Kodia includes security features by default:

| Feature | Status | Details |
|---------|--------|---------|
| JWT Authentication | ✅ Built-in | Access & refresh tokens with rotation |
| Two-Factor Auth (2FA) | ✅ Built-in | TOTP / Authenticator App support |
| RBAC & ABAC | ✅ Built-in | Granular permissions & policy engine |
| Rate Limiting | ✅ Built-in | Redis-based token bucket |
| Input Validation | ✅ Built-in | Server-side validation |
| CORS Protection | ✅ Built-in | Configurable origins with trace headers |
| SQL Injection Prevention | ✅ Built-in | Parameterized queries via GORM |
| Password Hashing | ✅ Built-in | Argon2/Bcrypt support |
| CSRF Protection | ✅ Built-in | Token validation |
| Security Headers | ✅ Built-in | CSP, HSTS, X-Frame-Options (DENY) |
| Audit Logging | ✅ Built-in | Immutable GORM-based action logs |
| Error Handling | ✅ Built-in | Sentry integration & recovery |
| IP Whitelisting | ✅ Built-in | CIDR/IP based route restrictions |
| Distributed Tracing | ✅ Built-in | OpenTelemetry (OTEL) integration |

---

## Authentication

### JWT Token-Based Authentication

Kodia uses **JWT (JSON Web Tokens)** with both access and refresh tokens:

```go
// Authentication flow
1. User logs in with credentials
2. Server validates credentials
3. Server issues:
   - Access token (short-lived: 15 minutes)
   - Refresh token (long-lived: 7 days)
4. Client stores tokens in secure storage
5. Client sends access token with each request
6. Server validates token signature and expiration
7. When access token expires, client uses refresh token to get new token
```

### Setup Authentication

```bash
# In backend/.env
JWT_ACCESS_SECRET=your-secret-key-here-minimum-32-characters
JWT_REFRESH_SECRET=your-refresh-secret-key-minimum-32-characters
JWT_ACCESS_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_DAYS=7
```

**Requirements:**
- ✅ Secrets must be 32+ characters
- ✅ Use strong random secrets (NOT words)
- ✅ Different secrets for access & refresh tokens
- ✅ Rotate secrets periodically in production

### Login Implementation

```go
// handlers/auth_handler.go
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: "Invalid request"})
        return
    }
    
    // 1. Find user by email
    user, err := h.userRepo.FindByEmail(c.Request.Context(), req.Email)
    if err != nil || user == nil {
        c.JSON(401, ErrorResponse{Error: "Invalid credentials"})
        return
    }
    
    // 2. Verify password
    if !h.authService.VerifyPassword(user.PasswordHash, req.Password) {
        c.JSON(401, ErrorResponse{Error: "Invalid credentials"})
        return
    }
    
    // 3. Generate tokens
    accessToken, refreshToken, err := h.jwtManager.GenerateTokens(user.ID)
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "Could not generate tokens"})
        return
    }
    
    // 4. Store refresh token in database
    err = h.refreshTokenRepo.Save(c.Request.Context(), &RefreshToken{
        UserID:    user.ID,
        Token:     refreshToken,
        ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
    })
    if err != nil {
        c.JSON(500, ErrorResponse{Error: "Could not save token"})
        return
    }
    
    c.JSON(200, LoginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    3600, // seconds
    })
}
```

### Protecting Endpoints

```go
// router.go
protected := api.Group("")
protected.Use(middleware.Auth(jwtManager))
{
    protected.GET("/me", handlers.GetCurrentUser)
    protected.POST("/posts", handlers.CreatePost)
    protected.DELETE("/posts/:id", handlers.DeletePost)
}

// Unauthenticated endpoints
api.POST("/login", handlers.Login)
api.POST("/register", handlers.Register)
api.GET("/posts", handlers.ListPosts)  // Public
```

### Token Refresh

```go
// Client-side: When access token expires
const refreshAccessToken = async () => {
    const response = await fetch('/api/auth/refresh', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
            refreshToken: localStorage.getItem('refreshToken')
        })
    })
    
    const data = await response.json()
    localStorage.setItem('accessToken', data.accessToken)
}

// Server-side: Refresh endpoint
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    var req RefreshTokenRequest
    c.BindJSON(&req)
    
    // 1. Validate refresh token exists
    storedToken, err := h.refreshTokenRepo.FindByToken(c.Request.Context(), req.RefreshToken)
    if err != nil {
        c.JSON(401, ErrorResponse{Error: "Invalid refresh token"})
        return
    }
    
    // 2. Check expiration
    if storedToken.ExpiresAt.Before(time.Now()) {
        c.JSON(401, ErrorResponse{Error: "Refresh token expired"})
        return
    }
    
    // 3. Generate new tokens
    newAccessToken, newRefreshToken, _ := h.jwtManager.GenerateTokens(storedToken.UserID)
    
    // 4. Rotate refresh token (invalidate old one)
    h.refreshTokenRepo.Delete(c.Request.Context(), storedToken.ID)
    h.refreshTokenRepo.Save(c.Request.Context(), &RefreshToken{
        UserID: storedToken.UserID,
        Token: newRefreshToken,
    })
    
    c.JSON(200, LoginResponse{
        AccessToken:  newAccessToken,
        RefreshToken: newRefreshToken,
    })
}
```

---

## Authorization (RBAC, ABAC & Permissions)

### Role-Based Access Control (RBAC)

```go
// User model with roles
type User struct {
    ID    string
    Email string
    Role  string  // "user", "admin", "moderator"
}

// Permission checking middleware
func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(*User)
        
        if user.Role != requiredRole {
            c.JSON(403, ErrorResponse{Error: "Forbidden"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Usage in routes
admin := api.Group("/admin")
admin.Use(middleware.Auth(jwtManager))
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/users", handlers.ListUsers)
    admin.DELETE("/users/:id", handlers.DeleteUser)
}
```

### Attribute-Based Access Control (ABAC)

```go
// Policy: Can user perform action on resource?
type Policy interface {
    Can(ctx context.Context, user *User, action string, resource interface{}) bool
}

// PostPolicy: Determine who can edit posts
type PostPolicy struct{}

func (p *PostPolicy) Can(ctx context.Context, user *User, action string, post *Post) bool {
    switch action {
    case "view":
        return true  // Anyone can view published posts
    case "create":
        return user != nil  // Authenticated users can create
    case "edit":
        return user.ID == post.AuthorID  // Only author can edit
    case "delete":
        return user.ID == post.AuthorID || user.Role == "admin"
    default:
        return false
    }
}

// Usage in handler
func (h *PostHandler) UpdatePost(c *gin.Context) {
    user := c.MustGet("user").(*User)
    post, _ := h.repo.FindByID(ctx, c.Param("id"))
    
    policy := &PostPolicy{}
    if !policy.Can(ctx, user, "edit", post) {
        c.JSON(403, ErrorResponse{Error: "Not authorized"})
        return
    }
    
    // Proceed with update
    h.service.UpdatePost(ctx, post)
}
```

---

## Input Validation

### Server-Side Validation

**Always validate on the server**, never trust client-side validation:

```go
// DTO with validation tags
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=5,max=200"`
    Content string `json:"content" validate:"required,min=20,max=10000"`
    Tags    []string `json:"tags" validate:"max=10,dive,required,min=2,max=50"`
}

// Middleware for automatic validation
func ValidateRequest(validate *validator.Validator) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req interface{}
        
        if err := c.BindJSON(&req); err != nil {
            c.JSON(400, ErrorResponse{Error: "Invalid JSON"})
            c.Abort()
            return
        }
        
        // Validate struct tags
        if err := validate.Struct(req); err != nil {
            c.JSON(400, ErrorResponse{
                Error: "Validation failed",
                Details: err.Error(),
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Usage
router.POST("/posts", ValidateRequest(validator), handlers.CreatePost)
```

### Sanitization

```go
import "github.com/microcosm-cc/bluemonday"

// Sanitize HTML content
func SanitizeHTML(dirty string) string {
    p := bluemonday.StrictPolicy()
    clean := p.Sanitize(dirty)
    return clean
}

// Usage in service
func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) error {
    // Remove dangerous HTML
    req.Content = SanitizeHTML(req.Content)
    // ... rest of logic
}
```

---

## SQL Injection Prevention

### Using GORM (Parameterized Queries)

```go
// ✅ GOOD: Parameterized query (safe)
var user User
db.Where("email = ?", userInput).First(&user)
// SQL: SELECT * FROM users WHERE email = $1

// ❌ BAD: String concatenation (unsafe)
var user User
db.Where(fmt.Sprintf("email = '%s'", userInput)).First(&user)
// Vulnerable to: ' OR '1'='1
```

### Never Use String Concatenation

```go
// ❌ NEVER DO THIS
query := "SELECT * FROM users WHERE email = '" + userInput + "'"
db.Raw(query).Scan(&users)

// ✅ ALWAYS DO THIS
db.Where("email = ?", userInput).Find(&users)
```

---

## CORS & CSRF Protection

### CORS Configuration

```go
// router.go
corsConfig := cors.Config{
    AllowOrigins:     cfg.CORS.AllowedOrigins,  // ["https://example.com"]
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           3600,  // 1 hour
}
engine.Use(cors.New(corsConfig))
```

### CSRF Token Protection

```go
// Middleware to generate CSRF token
func CSRFMiddleware(c *gin.Context) {
    token := generateSecureToken()
    c.SetCookie("csrf-token", token, 3600, "/", "example.com", true, true)
    c.Header("X-CSRF-Token", token)
    c.Next()
}

// Validate CSRF token on mutations
func ValidateCSRF(c *gin.Context) {
    if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
        c.Next()
        return
    }
    
    clientToken := c.GetHeader("X-CSRF-Token")
    cookieToken, _ := c.Cookie("csrf-token")
    
    if clientToken != cookieToken {
        c.JSON(403, ErrorResponse{Error: "CSRF validation failed"})
        c.Abort()
        return
    }
    
    c.Next()
}
```

---

## Rate Limiting

See [Rate Limiting Guide](../backend/docs/RATE_LIMITING.md) for detailed implementation.

```go
// Built-in rate limiters
auth := api.Group("/auth")
{
    authLimiter := middleware.AuthEndpointRateLimiter(redisClient, logger)
    auth.POST("/login", authLimiter.Middleware(), handlers.Login)      // 5 req/15min
    auth.POST("/register", authLimiter.Middleware(), handlers.Register)
}

api.GET("/posts", middleware.LooseRateLimiter(redisClient, logger).Middleware(), handlers.ListPosts)
// 100 req/minute
```

---

## Secrets Management

### Environment Variables

```bash
# ✅ Store secrets in .env (NOT in git)
JWT_ACCESS_SECRET=super-secret-key-here-32-characters-minimum
DATABASE_PASSWORD=database-password-here
API_KEY=third-party-api-key-here
SMTP_PASSWORD=email-password-here

# ❌ Never commit .env to git
echo ".env" >> .gitignore
echo ".env.local" >> .gitignore
```

### Loading Secrets

```go
import "github.com/joho/godotenv"

func Load(isDev bool) (*Config, error) {
    if isDev {
        godotenv.Load()  // Load from .env
    }
    
    cfg := &Config{
        JWT: JWTConfig{
            AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
            RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
        },
        Database: DatabaseConfig{
            Host:     os.Getenv("DATABASE_HOST"),
            Password: os.Getenv("DATABASE_PASSWORD"),
        },
    }
    
    // Validate secrets are set
    if len(cfg.JWT.AccessSecret) < 32 {
        return nil, fmt.Errorf("JWT_ACCESS_SECRET must be 32+ characters")
    }
    
    return cfg, nil
}
```

### Production Deployment Secrets

**Never put secrets in:**
- ❌ Code or git repository
- ❌ Docker images
- ❌ Config files
- ❌ Logs

**Use:**
- ✅ Environment variables
- ✅ Cloud key management (AWS Secrets Manager, GCP Secret Manager)
- ✅ Container orchestration secrets (Kubernetes Secrets)
- ✅ Secret management tools (HashiCorp Vault)

---

## HTTPS & TLS

### Forcing HTTPS

```go
// Middleware to redirect HTTP to HTTPS
func HTTPSRedirect() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Header.Get("X-Forwarded-Proto") != "https" {
            url := fmt.Sprintf("https://%s%s", c.Request.Host, c.Request.RequestURI)
            c.Redirect(301, url)
            return
        }
        c.Next()
    }
}

// Usage
if cfg.IsProduction() {
    engine.Use(HTTPSRedirect())
}
```

### Certificate Configuration

```bash
# Using Let's Encrypt with Docker
# In Dockerfile
RUN apt-get install certbot python3-certbot-nginx

# In docker-compose.yml
services:
  nginx:
    image: nginx:latest
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - /etc/letsencrypt:/etc/letsencrypt
    ports:
      - "443:443"
      - "80:80"
```

---

## Security Headers

### Built-In Headers

```go
// Middleware for security headers
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Prevent clickjacking
        c.Header("X-Frame-Options", "DENY")
        
        // Prevent MIME type sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        
        // Enable XSS protection
        c.Header("X-XSS-Protection", "1; mode=block")
        
        // Content Security Policy
        c.Header("Content-Security-Policy", "default-src 'self'")
        
        // HSTS (HTTP Strict Transport Security)
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        
        // Referrer Policy
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        
        // Permissions Policy
        c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        c.Next()
    }
}

// Usage
engine.Use(SecurityHeaders())
```

---

## Audit Logging

### Log Sensitive Operations

```go
// Log failed login attempts
func (h *AuthHandler) Login(c *gin.Context) {
    user, err := h.authService.Login(email, password)
    if err != nil {
        h.logger.Warn("Failed login attempt",
            zap.String("email", email),
            zap.String("ip", c.ClientIP()),
            zap.String("error", err.Error()),
        )
        return
    }
    
    h.logger.Info("User logged in",
        zap.String("user_id", user.ID),
        zap.String("email", user.Email),
        zap.String("ip", c.ClientIP()),
    )
}

// Log permission-related operations
func (h *UserHandler) DeleteUser(c *gin.Context) {
    actor := c.MustGet("user").(*User)
    target := c.Param("id")
    
    h.logger.Warn("User deletion requested",
        zap.String("actor_id", actor.ID),
        zap.String("actor_role", actor.Role),
        zap.String("target_user_id", target),
    )
    
    if err := h.service.DeleteUser(ctx, target); err != nil {
        return
    }
    
    h.logger.Info("User deleted successfully",
        zap.String("user_id", target),
    )
}
```

### Monitor Logs for Suspicious Activity

```bash
# Alert on multiple failed login attempts
grep "Failed login attempt" logs/app.log | \
  awk '{print $NF}' | \
  sort | uniq -c | \
  awk '$1 > 5 {print "ALERT: " $2 " has " $1 " failed attempts"}'
```

---

## Deployment Checklist

Before deploying to production:

- [ ] **Secrets**
  - [ ] JWT secrets are 32+ characters
  - [ ] Database password is strong
  - [ ] Secrets stored in environment, NOT in code
  - [ ] Secrets not in Docker images

- [ ] **HTTPS**
  - [ ] TLS certificate installed
  - [ ] HTTPS enforced on all endpoints
  - [ ] Certificate auto-renewal configured

- [ ] **Authentication**
  - [ ] Login rate limiting enabled
  - [ ] Password hashing uses bcrypt (cost 12+)
  - [ ] Refresh tokens stored securely
  - [ ] Session timeout configured

- [ ] **Authorization**
  - [ ] Role-based access control tested
  - [ ] Sensitive endpoints protected
  - [ ] Admin endpoints restricted

- [ ] **Input Validation**
  - [ ] All endpoints validate input
  - [ ] File uploads restricted
  - [ ] SQL injection testing done

- [ ] **Security Headers**
  - [ ] All headers configured
  - [ ] CSP policy strict
  - [ ] HSTS enabled

- [ ] **Rate Limiting**
  - [ ] Auth endpoints protected
  - [ ] API endpoints protected
  - [ ] Redis configured for rate limiting

- [ ] **Logging & Monitoring**
  - [ ] Audit logging enabled
  - [ ] Failed attempts logged
  - [ ] Alerts configured
  - [ ] Log retention policy set

- [ ] **Infrastructure**
  - [ ] Database backups enabled
  - [ ] Database connection pooling configured
  - [ ] Firewall rules restrictive
  - [ ] VPN access to database
  - [ ] Regular security updates

- [ ] **Testing**
  - [ ] Penetration testing done
  - [ ] Security tests in CI/CD
  - [ ] OWASP Top 10 reviewed

---

## Vulnerability Reporting

Found a security issue? **Please don't open a public issue.**

Instead, email: **security@kodia.dev** with:

- Description of vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We'll respond within 48 hours and work to resolve quickly.

---

**For more information:**
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://golang.org/doc/security)
- [Rate Limiting Guide](../backend/docs/RATE_LIMITING.md)
