# 🏗️ Architecture Guide

This document explains the design philosophy, architectural layers, and patterns used in Kodia Framework.

**Table of Contents:**
- [Design Philosophy](#design-philosophy)
- [Layered Architecture](#layered-architecture)
- [Component Diagram](#component-diagram)
- [Data Flow](#data-flow)
- [Domain-Driven Design](#domain-driven-design)
- [Design Patterns](#design-patterns)
- [Dependency Injection](#dependency-injection)
- [Testing Architecture](#testing-architecture)

---

## Design Philosophy

Kodia follows **Clean Architecture** principles with these core tenets:

### 1. Separation of Concerns

Each layer has a specific responsibility:
- **HTTP Layer** - Handle incoming requests
- **Service Layer** - Implement business logic
- **Repository Layer** - Data access and persistence
- **Infrastructure Layer** - External services (DB, cache, email)

### 2. Dependency Inversion

High-level modules don't depend on low-level modules. Both depend on abstractions (interfaces):

```go
// ❌ BAD: Service depends on concrete repository
type PostService struct {
    repo *PostgresPostRepository  // Tightly coupled
}

// ✅ GOOD: Service depends on interface
type PostService struct {
    repo ports.PostRepository  // Loosely coupled via interface
}
```

### 3. Convention Over Configuration

Sensible defaults reduce boilerplate:
- Models go in `domain/`
- Services go in `services/`
- Handlers go in `adapters/http/handlers/`
- Tests follow the same structure with `_test.go` suffix

### 4. Explicit Dependencies

All dependencies are passed to constructors, never hidden in globals:

```go
// ✅ GOOD: Explicit dependencies
func NewPostService(
    repo ports.PostRepository,
    cache ports.CacheProvider,
    logger *zap.Logger,
) *PostService {
    return &PostService{repo, cache, logger}
}
```

---

## Layered Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  HTTP Layer                             │
│  (Handlers, Middleware, CORS, Request/Response DTOs)   │
└────────────────────┬────────────────────────────────────┘
                     │ (HTTP Request/Response)
┌────────────────────▼────────────────────────────────────┐
│                Service Layer                            │
│  (Business Logic, Validation, Authorization)            │
└────────────────────┬────────────────────────────────────┘
                     │ (Domain Objects, Business Events)
┌────────────────────▼────────────────────────────────────┐
│              Repository Layer                           │
│  (Database Queries, Data Mapping, ORM)                 │
└────────────────────┬────────────────────────────────────┘
                     │ (SQL)
┌────────────────────▼────────────────────────────────────┐
│           Infrastructure Layer                          │
│  (PostgreSQL, Redis, S3, SMTP, External APIs)          │
└─────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

#### 1. HTTP Layer (`adapters/http/`)

**Purpose:** Handle HTTP requests and responses

**Components:**
- **Handlers** - Process incoming requests, call services, return responses
- **Middleware** - Cross-cutting concerns (auth, logging, CORS)
- **DTOs** - Data Transfer Objects for request/response
- **Validators** - Input validation rules

**Example:**
```go
// POST /api/posts
func (h *PostHandler) Create(c *gin.Context) {
    // 1. Extract & validate request
    var req dto.CreatePostRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 2. Call service
    post, err := h.service.Create(c.Request.Context(), req)
    
    // 3. Return response
    c.JSON(201, post)
}
```

**Key Rules:**
- ✅ Handle HTTP protocol details
- ✅ Call services for business logic
- ✅ Validate input before passing to service
- ❌ No direct database access
- ❌ No business logic in handlers

---

#### 2. Service Layer (`core/services/`)

**Purpose:** Implement business logic and use cases

**Components:**
- **Services** - Orchestrate operations across repositories
- **Domain Events** - Events published when important things happen
- **Authorization** - Check permissions before operations

**Example:**
```go
// PostService implements business rules for posts
type PostService struct {
    repo       ports.PostRepository
    authSvc    ports.AuthService
    eventBus   ports.EventBus
    logger     *zap.Logger
}

// CreatePost implements the use case: Create a new post
func (s *PostService) CreatePost(ctx context.Context, req *PostCreateRequest) (*Post, error) {
    // 1. Authorize: only authenticated users can create posts
    user := ctx.Value("user").(*User)
    if !s.authSvc.Can(user, "create_post") {
        return nil, ErrUnauthorized
    }
    
    // 2. Validate business rules
    if len(req.Title) == 0 {
        return nil, ErrInvalidTitle
    }
    
    // 3. Create domain object
    post := domain.NewPost(req.Title, req.Content, user.ID)
    
    // 4. Persist to database
    if err := s.repo.Save(ctx, post); err != nil {
        return nil, err
    }
    
    // 5. Publish domain event
    s.eventBus.Publish(domain.PostCreatedEvent{PostID: post.ID})
    
    s.logger.Info("Post created", zap.String("title", post.Title))
    return post, nil
}
```

**Key Rules:**
- ✅ Implement use cases/business rules
- ✅ Orchestrate across multiple repositories
- ✅ Publish domain events
- ✅ Validate business invariants
- ❌ Don't handle HTTP details
- ❌ Don't access external services directly (inject them)
- ❌ Don't return HTTP status codes

---

#### 3. Repository Layer (`adapters/repository/`)

**Purpose:** Abstract data access from services

**Components:**
- **Repository Interfaces** - Define data access contracts
- **Repository Implementations** - Concrete implementations (PostgreSQL, MySQL)
- **Query Builders** - Type-safe query construction
- **Mappers** - Map between domain objects and database records

**Example:**
```go
// ports/repository.go - Interface definition
type PostRepository interface {
    Save(ctx context.Context, post *Post) error
    FindByID(ctx context.Context, id string) (*Post, error)
    FindAll(ctx context.Context, limit int, offset int) ([]*Post, error)
    Update(ctx context.Context, post *Post) error
    Delete(ctx context.Context, id string) error
}

// repository/postgres/post_repository.go - Implementation
type PostgresPostRepository struct {
    db *gorm.DB
}

func (r *PostgresPostRepository) FindAll(ctx context.Context, limit, offset int) ([]*Post, error) {
    var posts []*Post
    if err := r.db.WithContext(ctx).
        Limit(limit).
        Offset(offset).
        Find(&posts).Error; err != nil {
        return nil, err
    }
    return posts, nil
}
```

**Key Rules:**
- ✅ Implement repository interfaces
- ✅ Use ORM (GORM) for queries
- ✅ Handle database errors gracefully
- ❌ No business logic
- ❌ Return only domain objects
- ❌ Don't know about HTTP or services

---

#### 4. Infrastructure Layer (`infrastructure/`)

**Purpose:** Integrate with external systems

**Components:**
- **Database** - PostgreSQL/MySQL connection & configuration
- **Cache** - Redis for caching & sessions
- **Storage** - File storage (local/S3)
- **Mailer** - SMTP email service
- **Logger** - Structured logging (Zap)
- **Config** - Environment configuration

**Example:**
```go
// infrastructure/database/postgres.go
func New(cfg *config.Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    return db, nil
}
```

---

## Component Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                    Client (Frontend)                         │
│                                                              │
└─────────────────────────┬──────────────────────────────────┘
                          │ HTTP/JSON
┌─────────────────────────▼──────────────────────────────────┐
│                  HTTP Handler Layer                         │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────┐   │
│  │PostHandler  │  │UserHandler   │  │AuthHandler     │   │
│  └──────┬──────┘  └──────┬───────┘  └────────┬───────┘   │
│         │                │                    │            │
├─────────┼────────────────┼────────────────────┼───────────┤
│         └────────────────┼────────────────────┘           │
│                         │ Call Services                    │
│         ┌───────────────▼────────────────┐                │
│         │    Service Layer              │                │
│         │                               │                │
│         │  ┌──────────┐  ┌───────────┐  │                │
│         │  │PostSvc   │  │UserSvc    │  │                │
│         │  └─────┬────┘  └─────┬─────┘  │                │
│         └────────┼─────────────┼────────┘                │
│                  │             │ Call Repositories       │
├──────────────────┼─────────────┼─────────────────────────┤
│                  │             │                         │
│         ┌────────▼─────┐  ┌────▼──────────┐             │
│         │PostRepository│  │UserRepository │             │
│         └────────┬─────┘  └────┬──────────┘             │
│                  │             │                         │
├──────────────────┼─────────────┼─────────────────────────┤
│                  │             │                         │
│         ┌────────▼─────────────▼──────┐                 │
│         │   Infrastructure Layer      │                 │
│         │  (PostgreSQL, Redis, etc)   │                 │
│         └────────────────────────────┘                 │
│                                                        │
└────────────────────────────────────────────────────────┘
```

---

## Data Flow

### Creating a Post (Request → Response)

```
1. REQUEST ARRIVES
   POST /api/posts
   {"title": "My First Post", "content": "Hello world"}

2. HTTP HANDLER
   PostHandler.Create(c *gin.Context)
   ├─ Extract request body
   ├─ Validate input format (DTO validation)
   └─ Call PostService.CreatePost()

3. SERVICE LAYER
   PostService.CreatePost(ctx, req)
   ├─ Check authorization (is user authenticated?)
   ├─ Validate business rules (title not empty, etc)
   ├─ Create domain object: Post
   └─ Call PostRepository.Save(post)

4. REPOSITORY LAYER
   PostRepository.Save(ctx, post)
   ├─ Map domain.Post to database model
   ├─ Execute SQL: INSERT INTO posts (...)
   └─ Return saved post with ID

5. SERVICE RETURNS
   ├─ Publish domain event: PostCreated
   ├─ Log operation
   └─ Return post to handler

6. HANDLER RETURNS
   c.JSON(201, post)
   └─ Return JSON response to client

7. RESPONSE SENT
   HTTP/1.1 201 Created
   {"id": "123", "title": "My First Post", ...}
```

---

## Domain-Driven Design

Kodia uses Domain-Driven Design (DDD) principles:

### Domain Objects (core/domain/)

```go
// Post is an aggregate root representing a blog post
type Post struct {
    ID        string
    Title     string
    Slug      string
    Content   string
    Author    *User
    Published bool
    Views     int
    CreatedAt time.Time
    UpdatedAt time.Time
}

// Business logic lives in the domain
func (p *Post) Publish() error {
    if p.Published {
        return fmt.Errorf("post already published")
    }
    p.Published = true
    return nil
}

// Domain events
type PostPublishedEvent struct {
    PostID    string
    PublishedAt time.Time
}
```

### Repository Pattern

```go
// Repositories abstract data access
type PostRepository interface {
    // Query operations
    FindByID(ctx context.Context, id string) (*Post, error)
    FindBySlug(ctx context.Context, slug string) (*Post, error)
    FindByAuthor(ctx context.Context, authorID string) ([]*Post, error)
    
    // Mutation operations
    Save(ctx context.Context, post *Post) error
    Update(ctx context.Context, post *Post) error
    Delete(ctx context.Context, id string) error
}
```

### Service Interfaces

```go
// Services define use cases
type PostService interface {
    CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error)
    GetPost(ctx context.Context, id string) (*Post, error)
    UpdatePost(ctx context.Context, id string, req *UpdatePostRequest) (*Post, error)
    DeletePost(ctx context.Context, id string) error
    PublishPost(ctx context.Context, id string) error
}
```

---

## Design Patterns

### 1. Repository Pattern

Abstraction layer for data access:

```go
type Repository interface {
    Save(ctx context.Context, entity Entity) error
    FindByID(ctx context.Context, id string) (Entity, error)
    Delete(ctx context.Context, id string) error
}
```

**Benefits:**
- Easy to test (mock repository)
- Easy to switch databases
- Encapsulates query logic

---

### 2. Service Locator / Dependency Injection

All dependencies injected via constructors:

```go
// Bad: Hidden dependencies
var db = gorm.Open(...)  // Global

// Good: Explicit dependencies
func NewPostService(repo PostRepository, cache Cache, logger Logger) *PostService {
    return &PostService{repo, cache, logger}
}
```

---

### 3. DTO Pattern

Transfer objects for HTTP layer:

```go
// DTO: Request from client
type CreatePostRequest struct {
    Title   string `json:"title" validate:"required,min=5"`
    Content string `json:"content" validate:"required,min=20"`
}

// DTO: Response to client
type PostResponse struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
}

// Convert domain object to response DTO
func ToPostResponse(p *domain.Post) *PostResponse {
    return &PostResponse{
        ID:        p.ID,
        Title:     p.Title,
        Content:   p.Content,
        Author:    p.Author.Name,
        CreatedAt: p.CreatedAt,
    }
}
```

---

### 4. Factory Pattern

Create complex objects:

```go
// PostFactory creates new posts with business logic
type PostFactory struct {
    slugGenerator SlugGenerator
}

func (f *PostFactory) CreatePost(title, content string, author *User) *Post {
    return &Post{
        ID:        uuid.New().String(),
        Title:     title,
        Content:   content,
        Slug:      f.slugGenerator.Generate(title),
        Author:    author,
        CreatedAt: time.Now(),
    }
}
```

---

### 5. Middleware Pattern

Cross-cutting concerns:

```go
// Logger middleware logs all requests
func LoggerMiddleware(logger Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        duration := time.Since(start)
        
        logger.Info("HTTP request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
        )
    }
}
```

---

## Dependency Injection

Kodia uses **constructor-based dependency injection** (not service locator):

```go
// main.go - Wire all dependencies at startup
func main() {
    // 1. Initialize infrastructure
    db := database.New(cfg)
    cache := cache.New(cfg)
    logger := logger.New(cfg)
    
    // 2. Create repositories
    postRepo := postgres.NewPostRepository(db)
    userRepo := postgres.NewUserRepository(db)
    
    // 3. Create services (inject dependencies)
    postSvc := services.NewPostService(postRepo, cache, logger)
    userSvc := services.NewUserService(userRepo, cache, logger)
    
    // 4. Create handlers (inject services)
    postHandler := handlers.NewPostHandler(postSvc, logger)
    userHandler := handlers.NewUserHandler(userSvc, logger)
    
    // 5. Setup router
    router := http.NewRouter(postHandler, userHandler)
    
    // 6. Start server
    router.Run(":8080")
}
```

**Benefits:**
- Explicit dependencies
- Easy to test (inject mocks)
- No hidden global state
- Clear object lifecycle

---

## Testing Architecture

### Unit Tests

Test individual functions in isolation:

```go
// service_test.go
func TestCreatePost(t *testing.T) {
    // 1. Setup (create mocks)
    mockRepo := &MockPostRepository{}
    mockCache := &MockCache{}
    mockLogger := &MockLogger{}
    
    // 2. Create service with mocks
    svc := NewPostService(mockRepo, mockCache, mockLogger)
    
    // 3. Test business logic
    post, err := svc.CreatePost(context.Background(), &CreatePostRequest{
        Title:   "Test Post",
        Content: "Test content",
    })
    
    // 4. Assert
    assert.NoError(t, err)
    assert.Equal(t, "Test Post", post.Title)
    assert.True(t, mockRepo.SaveCalled)
}
```

### Integration Tests

Test database interactions:

```go
// repository_test.go - Uses real database
func TestPostRepositorySave(t *testing.T) {
    // Use testcontainers to spin up real PostgreSQL
    db := testDB.NewPostgresContainer()
    defer db.Terminate()
    
    repo := NewPostRepository(db)
    
    post := &Post{ID: "1", Title: "Test"}
    err := repo.Save(context.Background(), post)
    
    assert.NoError(t, err)
    
    // Verify in database
    retrieved, _ := repo.FindByID(context.Background(), "1")
    assert.Equal(t, "Test", retrieved.Title)
}
```

### E2E Tests

Test complete user flows:

```go
// e2e_test.go
func TestCreatePostFlow(t *testing.T) {
    // 1. Start server
    server := startTestServer()
    defer server.Close()
    
    // 2. Register user
    user := registerUser(server)
    
    // 3. Login
    token := login(server, user)
    
    // 4. Create post
    resp, err := http.Post(
        fmt.Sprintf("%s/api/posts", server.URL),
        "application/json",
        createPostPayload(),
    )
    
    // 5. Assert
    assert.Equal(t, 201, resp.StatusCode)
}
```

---

## Best Practices

### ✅ DO

- Use interfaces for dependencies
- Keep services focused on business logic
- Validate at boundaries (HTTP layer)
- Log important operations
- Use context for cancellation/timeout
- Test at multiple levels (unit, integration, E2E)
- Keep domain logic in domain objects
- Use domain events for complex operations

### ❌ DON'T

- Pass HTTP details to services
- Access database directly from handlers
- Use global variables
- Ignore error handling
- Mix concerns across layers
- Return framework-specific errors from services
- Test only at E2E level
- Put business logic in handlers

---

**For more details, see related guides:**
- [Backend Development Guide](BACKEND_GUIDE.md)
- [Testing Guide](TESTING.md)
- [Security Guide](SECURITY.md)
