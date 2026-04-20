# 🔧 Backend - REST API Layer

The backend is a high-performance REST API built with **Go** and **Gin** framework, following **Clean Architecture** principles.

**Table of Contents:**
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Running the Server](#running-the-server)
- [API Endpoints](#api-endpoints)
- [Development Guide](#development-guide)
- [Common Tasks](#common-tasks)

---

## Quick Start

```bash
# 1. Install dependencies
go mod download

# 2. Setup environment
cp .env.example .env

# 3. Start PostgreSQL & Redis with Docker
docker-compose up -d

# 4. Run migrations
go run cmd/server/main.go migrate

# 5. Start development server
go run cmd/server/main.go

# Should see: Server is running on http://localhost:8080
```

---

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go             # Application entrypoint
│
├── internal/
│   ├── core/                   # Pure business logic (no framework dependencies)
│   │   ├── domain/
│   │   │   ├── user.go         # User entity with business logic
│   │   │   ├── post.go         # Post entity with business logic
│   │   │   ├── errors.go       # Domain-specific errors
│   │   │   └── events.go       # Domain events
│   │   ├── ports/
│   │   │   ├── repositories.go # Repository interfaces (contracts)
│   │   │   └── services.go     # Service interfaces
│   │   ├── services/
│   │   │   ├── auth_service.go # Authentication business logic
│   │   │   ├── user_service.go # User management logic
│   │   │   ├── post_service.go # Post management logic
│   │   │   └── post_repository_test.go
│   │   └── events/
│   │       └── dispatcher.go   # Event dispatching
│   │
│   ├── adapters/               # Framework-specific adapters
│   │   ├── http/
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go   # HTTP request handlers
│   │   │   │   ├── user_handler.go
│   │   │   │   ├── post_handler.go
│   │   │   │   └── *_test.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go           # JWT authentication
│   │   │   │   ├── cors.go           # CORS handling
│   │   │   │   ├── logger.go         # Request logging
│   │   │   │   ├── ratelimit.go      # Rate limiting
│   │   │   │   └── recovery.go       # Error recovery
│   │   │   ├── dto/
│   │   │   │   ├── auth_dto.go       # Request/Response DTOs
│   │   │   │   ├── user_dto.go
│   │   │   │   ├── post_dto.go
│   │   │   │   └── response.go       # Standard response wrapper
│   │   │   ├── validators/
│   │   │   │   ├── auth_validator.go # Input validation rules
│   │   │   │   ├── user_validator.go
│   │   │   │   └── post_validator.go
│   │   │   ├── cors_validator.go     # CORS configuration validation
│   │   │   ├── cors_validator_test.go
│   │   │   └── router.go             # Route definitions
│   │   │
│   │   └── repository/
│   │       ├── postgres/
│   │       │   ├── user_repository.go    # PostgreSQL implementation
│   │       │   ├── post_repository.go
│   │       │   └── refresh_token_repository.go
│   │       ├── cache/
│   │       │   └── redis_cache.go       # Redis caching
│   │       └── factory.go               # Repository factory
│   │
│   └── infrastructure/         # External integrations
│       ├── database/
│       │   ├── postgres.go          # Database connection
│       │   ├── migrations/          # SQL migration files
│       │   │   └── sql/
│       │   │       ├── 000001_create_users_table.up.sql
│       │   │       └── 000001_create_users_table.down.sql
│       │   ├── seeders/             # Test data seeders
│       │   └── database_test.go
│       ├── cache/
│       │   ├── redis.go             # Redis client
│       │   └── provider.go           # Cache interface
│       ├── logger/
│       │   └── logger.go            # Zap structured logging
│       ├── storage/
│       │   ├── s3_provider.go       # AWS S3 storage
│       │   ├── local_provider.go    # Local file storage
│       │   ├── path_validator.go
│       │   └── path_validator_test.go
│       ├── mailer/
│       │   ├── smtp_mailer.go       # SMTP email service
│       │   └── smtp_mailer_test.go
│       ├── worker/
│       │   └── asynq_provider.go    # Background job queue
│       └── events/
│           └── dispatcher.go         # Event dispatching
│
├── pkg/                        # Reusable utilities
│   ├── config/
│   │   └── config.go           # Configuration loading
│   ├── jwt/
│   │   └── jwt.go              # JWT token management
│   └── pathutil/
│       └── validator.go         # Path validation utilities
│
├── tests/                      # Test suites
│   ├── unit/
│   │   ├── services/           # Service unit tests
│   │   └── handlers/           # Handler unit tests
│   ├── integration/            # Database integration tests
│   │   ├── post_repository_test.go
│   │   ├── user_repository_test.go
│   │   └── fixtures.go         # Test data
│   ├── e2e/                    # End-to-end API tests
│   │   ├── auth_e2e_test.go
│   │   └── posts_e2e_test.go
│   ├── fixtures/               # Shared test fixtures
│   ├── helpers.go              # Test utilities
│   └── mocks/                  # Mock objects
│       ├── mock_repository.go
│       └── mock_service.go
│
├── docs/                       # Generated API documentation
│   └── swagger.json            # OpenAPI specification
│
├── .env.example                # Environment template
├── .dockerignore               # Docker build exclusions
├── Dockerfile                  # Container image
├── go.mod                      # Go module definition
├── go.sum                      # Dependency hashes
├── Makefile                    # Build commands
└── README.md                   # This file
```

---

## Running the Server

### Development Mode

```bash
# With auto-reload using air
go install github.com/cosmtrek/air@latest
air

# Or manual build & run
go run cmd/server/main.go

# Server listens on http://localhost:8080
```

### Production Build

```bash
# Build optimized binary
go build -o server cmd/server/main.go

# Run with optimizations
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server cmd/server/main.go

# Run binary
./server
```

### With Docker

```bash
# Build image
docker build -t kodia-backend .

# Run container
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/db" \
  -e REDIS_URL="redis://localhost:6379" \
  kodia-backend

# Or use docker-compose
docker-compose up backend
```

---

## API Endpoints

### Authentication Endpoints

```
POST   /api/auth/register        # User registration
POST   /api/auth/login           # User login (returns access + refresh tokens)
POST   /api/auth/refresh         # Refresh access token
POST   /api/auth/logout          # Logout current session
POST   /auth/logout-all          # Logout all sessions (requires auth)
GET    /api/auth/me              # Get current user (requires auth)
```

### User Endpoints

```
GET    /api/users/me             # Get current user profile (requires auth)
POST   /api/users/me/change-password  # Change password (requires auth)
GET    /api/users                # List users (admin only)
GET    /api/users/:id            # Get user by ID (admin only)
PATCH  /api/users/:id            # Update user (admin only)
DELETE /api/users/:id            # Delete user (admin only)
```

### Post Endpoints (Example)

```
GET    /api/posts                # List all posts
GET    /api/posts/:id            # Get post by ID
POST   /api/posts                # Create post (requires auth)
PUT    /api/posts/:id            # Update post (requires auth)
DELETE /api/posts/:id            # Delete post (requires auth)
```

### Health Check

```
GET    /api/health               # Health check endpoint
```

### API Documentation

```
GET    /api/docs                 # Interactive Swagger UI
GET    /api/docs/swagger.json    # OpenAPI specification
```

---

## Development Guide

### Creating a New Handler

```go
// handlers/post_handler.go
type PostHandler struct {
    service services.PostService
    validate *validator.Validator
    logger  *zap.Logger
}

func NewPostHandler(
    service services.PostService,
    validate *validator.Validator,
    logger *zap.Logger,
) *PostHandler {
    return &PostHandler{service, validate, logger}
}

// GET /api/posts
func (h *PostHandler) GetAll(c *gin.Context) {
    posts, err := h.service.GetAllPosts(c.Request.Context())
    if err != nil {
        h.logger.Error("failed to fetch posts", zap.Error(err))
        c.JSON(500, ErrorResponse{Error: "Internal server error"})
        return
    }
    
    c.JSON(200, posts)
}
```

### Creating a New Service

```go
// services/post_service.go
type PostService struct {
    repo   ports.PostRepository
    cache  ports.CacheProvider
    logger *zap.Logger
}

func NewPostService(
    repo ports.PostRepository,
    cache ports.CacheProvider,
    logger *zap.Logger,
) *PostService {
    return &PostService{repo, cache, logger}
}

func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
    // Validate business rules
    if len(req.Title) == 0 {
        return nil, ErrInvalidTitle
    }
    
    // Create domain object
    post := domain.NewPost(req.Title, req.Content)
    
    // Persist to database
    if err := s.repo.Save(ctx, post); err != nil {
        return nil, fmt.Errorf("failed to save post: %w", err)
    }
    
    return post, nil
}
```

### Creating a New Repository

```go
// repository/postgres/post_repository.go
type PostgresPostRepository struct {
    db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostgresPostRepository {
    return &PostgresPostRepository{db}
}

func (r *PostgresPostRepository) FindAll(ctx context.Context) ([]*Post, error) {
    var posts []*Post
    if err := r.db.WithContext(ctx).Find(&posts).Error; err != nil {
        return nil, err
    }
    return posts, nil
}
```

---

## Common Tasks

### Database Migrations

```bash
# Run migrations
go run cmd/server/main.go migrate

# Create migration
go run cmd/server/main.go migrate create_posts_table

# Rollback migration
go run cmd/server/main.go migrate rollback
```

### Running Tests

```bash
# Run all tests
make test

# Run specific test file
go test ./tests/unit/services/post_service_test.go

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Check security
gosec ./...

# Run all checks
make lint
```

### Debugging

```bash
# Run with debug logs
DEBUG=true go run cmd/server/main.go

# Use delve debugger
dlv debug cmd/server/main.go
(dlv) break main.main
(dlv) continue
```

---

## Environment Variables

See `.env.example` for all available options:

```bash
# Application
APP_NAME=Kodia
APP_ENV=development
APP_PORT=8080
APP_LOG_LEVEL=debug

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_NAME=kodia

# Cache
REDIS_URL=redis://localhost:6379

# JWT
JWT_ACCESS_SECRET=your-secret-key-here-32-chars-min
JWT_REFRESH_SECRET=your-refresh-secret-key-here-32-chars-min
JWT_ACCESS_EXPIRY_HOURS=24
JWT_REFRESH_EXPIRY_DAYS=7

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000

# Storage
STORAGE_PROVIDER=local  # or 's3'
STORAGE_LOCAL_DIR=./uploads
STORAGE_S3_BUCKET=my-bucket
STORAGE_S3_REGION=us-east-1

# Email
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password
MAIL_FROM_ADDRESS=noreply@kodia.dev
```

---

## Troubleshooting

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
APP_PORT=8081
```

### Database Connection Error

```bash
# Check PostgreSQL is running
docker ps | grep postgres

# Check connection string
psql postgresql://user:password@localhost:5432/kodia

# Reset database
docker-compose exec postgres dropdb kodia
docker-compose exec postgres createdb kodia
```

### Module Not Found

```bash
# Download all dependencies
go mod download

# Tidy module file
go mod tidy

# Verify module
go mod verify
```

---

## Performance Tips

1. **Use indexes** on frequently queried columns
2. **Cache results** using Redis for read-heavy operations
3. **Use pagination** to avoid loading all records
4. **Profile code** using pprof: `go tool pprof http://localhost:8080/debug/pprof`
5. **Batch operations** for bulk inserts/updates

---

## Further Reading

- [Backend Development Guide](../docs/BACKEND_GUIDE.md)
- [Architecture Guide](../docs/ARCHITECTURE.md)
- [Security Guide](../docs/SECURITY.md)
- [Testing Guide](../docs/TESTING.md)
- [Go Best Practices](https://golang.org/doc/effective_go)
- [Gin Documentation](https://gin-gonic.com/)
