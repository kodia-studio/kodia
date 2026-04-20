# 📚 Backend Development Guide

Complete guide to building REST APIs with Kodia's Go/Gin backend.

**Table of Contents:**
- [HTTP Handlers](#http-handlers)
- [Routing](#routing)
- [Middleware](#middleware)
- [Services](#services)
- [Database & Repositories](#database--repositories)
- [Validation](#validation)
- [Error Handling](#error-handling)
- [Authentication & Authorization](#authentication--authorization)
- [Logging](#logging)
- [Best Practices](#best-practices)

---

## HTTP Handlers

Handlers process HTTP requests and return responses.

### Creating a Handler

```go
// handlers/user_handler.go
type UserHandler struct {
    service    services.UserService
    validate   *validator.Validator
    logger     *zap.Logger
}

// Constructor (dependency injection)
func NewUserHandler(
    service services.UserService,
    validate *validator.Validator,
    logger *zap.Logger,
) *UserHandler {
    return &UserHandler{service, validate, logger}
}

// GET /api/users
func (h *UserHandler) GetAll(c *gin.Context) {
    users, err := h.service.GetAllUsers(c.Request.Context())
    if err != nil {
        h.logger.Error("failed to get users", zap.Error(err))
        c.JSON(500, ErrorResponse{Error: "Internal server error"})
        return
    }
    
    c.JSON(200, users)
}

// GET /api/users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
    id := c.Param("id")
    
    user, err := h.service.GetUserByID(c.Request.Context(), id)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            c.JSON(404, ErrorResponse{Error: "User not found"})
            return
        }
        c.JSON(500, ErrorResponse{Error: "Internal server error"})
        return
    }
    
    c.JSON(200, user)
}

// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
    var req dto.CreateUserRequest
    
    // Parse request body
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: "Invalid request body"})
        return
    }
    
    // Validate input
    if err := h.validate.Struct(req); err != nil {
        c.JSON(400, ErrorResponse{Error: "Validation failed"})
        return
    }
    
    // Call service
    user, err := h.service.CreateUser(c.Request.Context(), &req)
    if err != nil {
        h.logger.Error("failed to create user", zap.Error(err))
        c.JSON(500, ErrorResponse{Error: "Failed to create user"})
        return
    }
    
    c.JSON(201, user)
}
```

### Handler Best Practices

✅ **DO:**
```go
// Handle errors gracefully
if err != nil {
    h.logger.Error("operation failed", zap.Error(err))
    c.JSON(500, ErrorResponse{Error: "Internal server error"})
    return
}

// Validate input early
if err := h.validate.Struct(req); err != nil {
    c.JSON(400, ErrorResponse{Error: "Validation failed"})
    return
}

// Use context for timeouts
user, err := h.service.GetUser(c.Request.Context(), id)

// Log important operations
h.logger.Info("user created", zap.String("user_id", user.ID))
```

❌ **DON'T:**
```go
// Expose internal errors
c.JSON(500, err)  // Shows database errors to client!

// Mix business logic with HTTP
if user.Email == "" {  // Should be in service
    return
}

// Forget to validate
h.service.UpdateUser(c.Request.Context(), req)  // What if req is invalid?

// Return wrong status codes
c.JSON(200, ErrorResponse{Error: "Not found"})  // Should be 404
```

---

## Routing

Define API routes in the router.

### Basic Routing

```go
// router.go
func (r *Router) Setup() *gin.Engine {
    engine := gin.New()
    
    api := engine.Group("/api")
    {
        // Public endpoints
        api.GET("/health", r.healthHandler.Check)
        api.POST("/auth/login", r.authHandler.Login)
        api.POST("/auth/register", r.authHandler.Register)
        
        // Protected endpoints
        protected := api.Group("")
        protected.Use(middleware.Auth(r.jwtManager))
        {
            protected.GET("/users/me", r.userHandler.GetMe)
            protected.POST("/posts", r.postHandler.Create)
        }
        
        // Admin endpoints
        admin := api.Group("")
        admin.Use(middleware.Auth(r.jwtManager))
        admin.Use(middleware.RequireRole("admin"))
        {
            admin.GET("/users", r.userHandler.GetAll)
            admin.DELETE("/users/:id", r.userHandler.Delete)
        }
    }
    
    return engine
}
```

### Route Groups

```go
// Organize routes logically
auth := api.Group("/auth")
{
    auth.POST("/login", h.Login)
    auth.POST("/register", h.Register)
    auth.POST("/refresh", h.Refresh)
}

posts := api.Group("/posts")
posts.Use(middleware.Auth(jwtManager))
{
    posts.POST("", h.Create)
    posts.PUT("/:id", h.Update)
    posts.DELETE("/:id", h.Delete)
}
```

---

## Middleware

Cross-cutting concerns like auth, logging, CORS.

### Using Middleware

```go
// Apply globally
engine.Use(middleware.Logger(logger))
engine.Use(middleware.Recovery())

// Apply to groups
protected := api.Group("")
protected.Use(middleware.Auth(jwtManager))
{
    protected.POST("/posts", h.CreatePost)
}
```

### Creating Custom Middleware

```go
// middleware/custom.go
func RequireRole(requiredRole string) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := c.MustGet("user").(*domain.User)
        
        if user.Role != requiredRole {
            c.JSON(403, ErrorResponse{Error: "Forbidden"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

// Usage
admin := api.Group("")
admin.Use(middleware.Auth(jwtManager))
admin.Use(middleware.RequireRole("admin"))
```

### Built-in Middleware

```go
// Authentication
protected.Use(middleware.Auth(jwtManager))

// Rate limiting
auth.POST("/login", middleware.AuthEndpointRateLimiter(redisClient, logger).Middleware())

// CORS
engine.Use(cors.New(corsConfig))

// Logging
engine.Use(middleware.Logger(logger))

// Error recovery
engine.Use(middleware.Recovery())
```

---

## Services

Business logic layer - where most of the work happens.

### Creating a Service

```go
// services/post_service.go
type PostService struct {
    repo       ports.PostRepository
    cache      ports.CacheProvider
    eventBus   ports.EventBus
    logger     *zap.Logger
}

func NewPostService(
    repo ports.PostRepository,
    cache ports.CacheProvider,
    eventBus ports.EventBus,
    logger *zap.Logger,
) *PostService {
    return &PostService{repo, cache, eventBus, logger}
}

// CreatePost implements the create post use case
func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
    // 1. Validate business rules
    if len(req.Title) == 0 {
        return nil, fmt.Errorf("title is required")
    }
    if len(req.Content) < 10 {
        return nil, fmt.Errorf("content must be at least 10 characters")
    }
    
    // 2. Create domain object
    post := domain.NewPost(req.Title, req.Content)
    
    // 3. Persist to database
    if err := s.repo.Save(ctx, post); err != nil {
        return nil, fmt.Errorf("failed to save post: %w", err)
    }
    
    // 4. Publish domain event
    s.eventBus.Publish(domain.PostCreatedEvent{
        PostID:    post.ID,
        CreatedAt: time.Now(),
    })
    
    // 5. Log operation
    s.logger.Info("post created", zap.String("post_id", post.ID))
    
    return post, nil
}

// GetPostByID gets a post from cache or database
func (s *PostService) GetPostByID(ctx context.Context, id string) (*Post, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("post:%s", id)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var post Post
        json.Unmarshal(cached, &post)
        return &post, nil
    }
    
    // Query database
    post, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }
    
    // Cache result for 1 hour
    if data, err := json.Marshal(post); err == nil {
        s.cache.Set(ctx, cacheKey, data, 3600)
    }
    
    return post, nil
}
```

### Service Best Practices

✅ **DO:**
```go
// Implement use cases
func (s *Service) PerformBusinessAction() error

// Use interfaces for dependencies
type Service struct {
    repo ports.Repository
    cache ports.CacheProvider
}

// Handle context for timeouts
func (s *Service) DoSomething(ctx context.Context) error {
    return s.repo.Query(ctx)
}

// Publish domain events
s.eventBus.Publish(ImportantEvent{})

// Log operations
s.logger.Info("action completed", zap.String("id", id))
```

❌ **DON'T:**
```go
// Mix HTTP with business logic
func (s *Service) CreatePost() *http.Response

// Use global variables
var globalRepo = NewRepository()  // Hard to test

// Ignore errors
s.repo.Save(ctx, obj)  // What if it fails?

// Retrieve HTTP headers
user := req.Header.Get("Authorization")  // Should be passed in
```

---

## Database & Repositories

Abstract data access with repositories.

### Repository Interface

```go
// ports/repository.go
type PostRepository interface {
    // Queries
    FindByID(ctx context.Context, id string) (*Post, error)
    FindAll(ctx context.Context, limit int, offset int) ([]*Post, error)
    FindByAuthor(ctx context.Context, authorID string) ([]*Post, error)
    
    // Mutations
    Save(ctx context.Context, post *Post) error
    Update(ctx context.Context, post *Post) error
    Delete(ctx context.Context, id string) error
}
```

### PostgreSQL Implementation

```go
// repository/postgres/post_repository.go
type PostgresPostRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostgresPostRepository {
    return &PostgresPostRepository{db}
}

// FindByID retrieves a single post
func (r *PostgresPostRepository) FindByID(ctx context.Context, id string) (*Post, error) {
    var post Post
    if err := r.db.WithContext(ctx).First(&post, "id = ?", id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrPostNotFound
        }
        return nil, fmt.Errorf("failed to find post: %w", err)
    }
    return &post, nil
}

// FindAll retrieves paginated posts
func (r *PostgresPostRepository) FindAll(ctx context.Context, limit int, offset int) ([]*Post, error) {
    var posts []*Post
    if err := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Order("created_at DESC").
        Find(&posts).Error; err != nil {
        return nil, fmt.Errorf("failed to find posts: %w", err)
    }
    return posts, nil
}

// Save creates a new post
func (r *PostgresPostRepository) Save(ctx context.Context, post *Post) error {
    if err := r.db.WithContext(ctx).Create(post).Error; err != nil {
        return fmt.Errorf("failed to save post: %w", err)
    }
    return nil
}

// Update modifies an existing post
func (r *PostgresPostRepository) Update(ctx context.Context, post *Post) error {
    if err := r.db.WithContext(ctx).Save(post).Error; err != nil {
        return fmt.Errorf("failed to update post: %w", err)
    }
    return nil
}

// Delete removes a post
func (r *PostgresPostRepository) Delete(ctx context.Context, id string) error {
    if err := r.db.WithContext(ctx).Delete(&Post{}, "id = ?", id).Error; err != nil {
        return fmt.Errorf("failed to delete post: %w", err)
    }
    return nil
}
```

### Using Repositories in Services

```go
func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
    post := domain.NewPost(req.Title, req.Content)
    
    // Repository handles the SQL
    if err := s.repo.Save(ctx, post); err != nil {
        return nil, err
    }
    
    return post, nil
}
```

---

## Validation

Validate user input server-side.

### Using Struct Tags

```go
// dto/user_dto.go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=3,max=100"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8,max=256"`
    Age      int    `json:"age" validate:"gte=18"`
}
```

### Validation in Handlers

```go
func (h *UserHandler) Create(c *gin.Context) {
    var req dto.CreateUserRequest
    
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{Error: "Invalid JSON"})
        return
    }
    
    // Validate request
    if err := h.validate.Struct(req); err != nil {
        validationErrors := err.(validator.ValidationErrors)
        c.JSON(400, ErrorResponse{
            Error: "Validation failed",
            Details: formatValidationErrors(validationErrors),
        })
        return
    }
    
    // Proceed if valid
    user, _ := h.service.CreateUser(c.Request.Context(), &req)
    c.JSON(201, user)
}

// Helper to format validation errors
func formatValidationErrors(errors validator.ValidationErrors) map[string]string {
    result := make(map[string]string)
    for _, err := range errors {
        result[err.Field()] = err.Tag()
    }
    return result
}
```

---

## Error Handling

Handle errors gracefully without exposing internals.

### Domain Errors

```go
// domain/errors.go
var (
    ErrUserNotFound      = errors.New("user not found")
    ErrInvalidEmail      = errors.New("invalid email format")
    ErrDuplicateEmail    = errors.New("email already exists")
    ErrInvalidPassword   = errors.New("invalid password")
    ErrUnauthorized      = errors.New("unauthorized")
    ErrForbidden         = errors.New("forbidden")
)
```

### Error Handling in Handlers

```go
func (h *UserHandler) Create(c *gin.Context) {
    user, err := h.service.CreateUser(c.Request.Context(), &req)
    
    if err != nil {
        // Match domain errors to HTTP status codes
        switch {
        case errors.Is(err, domain.ErrDuplicateEmail):
            c.JSON(409, ErrorResponse{Error: "Email already exists"})
        case errors.Is(err, domain.ErrInvalidEmail):
            c.JSON(400, ErrorResponse{Error: "Invalid email"})
        default:
            // Never expose internal errors to client
            h.logger.Error("failed to create user", zap.Error(err))
            c.JSON(500, ErrorResponse{Error: "Internal server error"})
        }
        return
    }
    
    c.JSON(201, user)
}
```

### HTTP Status Codes

```
200 OK                  - Request succeeded
201 Created             - Resource created
204 No Content          - Successful with no response body
400 Bad Request         - Invalid input
401 Unauthorized        - Authentication required
403 Forbidden           - Insufficient permissions
404 Not Found           - Resource not found
409 Conflict            - Resource already exists
422 Unprocessable       - Validation failed
429 Too Many Requests   - Rate limited
500 Internal Server     - Server error
503 Service Unavailable - Service down
```

---

## Authentication & Authorization

Secure your API.

### JWT Authentication

```go
// Use the auth middleware
protected := api.Group("")
protected.Use(middleware.Auth(jwtManager))
{
    protected.GET("/me", h.GetCurrentUser)
    protected.POST("/posts", h.CreatePost)
}

// Access authenticated user in handler
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
    user := c.MustGet("user").(*User)  // Set by Auth middleware
    c.JSON(200, user)
}
```

### Authorization with Roles

```go
// Require specific role
admin := api.Group("")
admin.Use(middleware.Auth(jwtManager))
admin.Use(middleware.RequireRole("admin"))
{
    admin.GET("/users", h.ListUsers)
    admin.DELETE("/users/:id", h.DeleteUser)
}
```

### Advanced Authorization with Policies

```go
// Check if user can perform action on resource
func (h *PostHandler) Update(c *gin.Context) {
    user := c.MustGet("user").(*User)
    postID := c.Param("id")
    post, _ := h.repo.FindByID(c.Request.Context(), postID)
    
    policy := &PostPolicy{}
    if !policy.Can(c.Request.Context(), user, "edit", post) {
        c.JSON(403, ErrorResponse{Error: "Forbidden"})
        return
    }
    
    // User is authorized to edit
    h.service.UpdatePost(c.Request.Context(), post)
}
```

---

## Logging

Track important events and errors.

### Structured Logging with Zap

```go
// Log different levels
h.logger.Info("user logged in",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)

h.logger.Warn("multiple failed login attempts",
    zap.String("ip", c.ClientIP()),
    zap.Int("attempts", 5),
)

h.logger.Error("failed to save post",
    zap.Error(err),
    zap.String("post_id", post.ID),
)

// Debug logs (only in development)
h.logger.Debug("database query",
    zap.String("query", sql),
    zap.Int64("duration_ms", duration),
)
```

### What to Log

✅ **DO log:**
- Failed authentication attempts
- Successful user actions (created, deleted, etc.)
- Errors with context
- Performance-critical operations

❌ **DON'T log:**
- User passwords
- API keys
- Credit card numbers
- Raw request bodies with sensitive data

---

## Best Practices

### Project Organization

```go
// ✅ Organize by feature
api/
├── auth/
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── dto.go
├── posts/
│   ├── handler.go
│   ├── service.go
│   ├── repository.go
│   └── dto.go
└── users/
    ├── handler.go
    ├── service.go
    ├── repository.go
    └── dto.go
```

### Error Wrapping

```go
// ✅ Preserve error chain with %w
if err != nil {
    return fmt.Errorf("failed to save user: %w", err)
}

// ❌ Don't lose error context
if err != nil {
    return err  // Lost context about what failed
}
```

### Context Usage

```go
// ✅ Pass context for timeouts/cancellation
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.FindByID(ctx, id)
}

// ❌ Ignore context
func (s *Service) GetUser(id string) (*User, error) {  // No timeout support
    return s.repo.FindByID(context.Background(), id)
}
```

---

**Next: [Testing Guide](TESTING.md) | [Security Guide](SECURITY.md)**
