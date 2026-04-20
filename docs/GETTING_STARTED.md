# 🚀 Getting Started with Kodia

Welcome to Kodia! This guide will walk you through setting up your development environment and creating your first application.

**Table of Contents:**
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Project Structure](#project-structure)
- [Your First Feature](#your-first-feature)
- [Running Tests](#running-tests)
- [Deployment](#deployment)
- [Next Steps](#next-steps)

---

## Prerequisites

Before you begin, ensure you have installed:

| Tool | Version | Purpose |
|------|---------|---------|
| **Go** | 1.26+ | Backend runtime |
| **Node.js** | 25+ | Frontend runtime |
| **Docker** | Latest | Container runtime |
| **Docker Compose** | 2.0+ | Multi-container orchestration |
| **PostgreSQL** | 15+ | Primary database (or MySQL 8+) |
| **Redis** | 7+ | Caching & session storage |

**Installation links:**
- [Install Go](https://go.dev/doc/install)
- [Install Node.js](https://nodejs.org/)
- [Install Docker](https://docs.docker.com/get-docker/)
- [Install PostgreSQL](https://www.postgresql.org/download/)

---

## Installation

### Step 1: Clone Repository

```bash
# Clone the Kodia repository
git clone https://github.com/kodia-studio/kodia.git my-awesome-app
cd my-awesome-app

# Verify directory structure
ls -la
# Output should show: backend/ frontend/ cli/ docker-compose.yml Makefile
```

### Step 2: Setup Environment Variables

```bash
# Copy environment templates
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Edit backend configuration
nano backend/.env
# Key settings to verify:
# - DATABASE_URL=postgres://user:password@localhost:5432/kodia
# - REDIS_URL=redis://localhost:6379
# - JWT_ACCESS_SECRET=your-secret-key-here (min 32 chars)
# - APP_PORT=8080

# Edit frontend configuration
nano frontend/.env
# Key settings:
# - VITE_API_BASE_URL=http://localhost:8080/api
```

### Step 3: Install Dependencies

```bash
# Backend dependencies
cd backend
go mod download
cd ..

# Frontend dependencies
cd frontend
npm install
cd ..

# CLI dependencies
cd cli
go mod download
cd ..
```

### Step 4: Start Services with Docker

```bash
# Start all services (PostgreSQL, Redis, Backend, Frontend)
docker-compose up -d

# Verify services are running
docker-compose ps
# Output:
# NAME                COMMAND                  SERVICE      STATUS
# kodia-postgres      "docker-entrypoint.s…"   postgres     Up 2 seconds
# kodia-redis         "redis-server"           redis        Up 2 seconds

# Wait 5 seconds for databases to be ready
sleep 5

# Check PostgreSQL is accessible
psql postgresql://postgres:password@localhost:5432/kodia -c "SELECT 1"
```

### Step 5: Run Database Migrations

```bash
# Navigate to backend
cd backend

# Run migrations
go run cmd/server/main.go migrate

# Or using make
make migrate

# Verify migration success - check if users table exists
psql postgresql://postgres:password@localhost:5432/kodia -c "\dt"
```

### Step 6: Start Development Servers

```bash
# Terminal 1: Start backend server
cd backend
go run cmd/server/main.go
# Output: Server is running on http://localhost:8080

# Terminal 2: Start frontend dev server
cd frontend
npm run dev
# Output: 
# VITE v... dev server running at:
# > Local:   http://localhost:5173/

# Terminal 3 (Optional): Watch CLI changes
cd cli
go run kodia/main.go --version
```

### Step 7: Verify Installation

```bash
# Test backend API
curl http://localhost:8080/api/health
# Expected response: {"success": true, "message": "Kodia Backend is healthy"}

# Test frontend
open http://localhost:5173
# Should see login page

# Test API documentation
open http://localhost:8080/api/docs
# Should see Swagger UI
```

---

## Project Structure

```
kodia/
├── backend/                     # Go/Gin REST API
│   ├── cmd/
│   │   └── server/
│   │       └── main.go         # Application entrypoint
│   ├── internal/
│   │   ├── core/               # Business logic (pure Go)
│   │   │   ├── domain/         # Entities & value objects
│   │   │   ├── ports/          # Interface definitions
│   │   │   ├── services/       # Business logic
│   │   │   └── events/         # Domain events
│   │   ├── adapters/           # Framework adapters
│   │   │   ├── http/           # HTTP handlers (Gin)
│   │   │   └── repository/     # Data access
│   │   └── infrastructure/     # External services
│   │       ├── database/       # PostgreSQL/MySQL
│   │       ├── cache/          # Redis
│   │       ├── storage/        # File storage
│   │       ├── mailer/         # Email
│   │       └── logger/         # Logging
│   ├── pkg/                     # Shared packages
│   ├── tests/                   # Test suites
│   ├── docs/                    # Generated API docs
│   ├── Dockerfile              # Container image
│   ├── .env.example            # Environment template
│   ├── go.mod                  # Go module definition
│   └── Makefile                # Build commands
│
├── frontend/                    # SvelteKit + Tailwind application
│   ├── src/
│   │   ├── lib/
│   │   │   ├── api/            # API client & requests
│   │   │   ├── components/     # Reusable components
│   │   │   ├── stores/         # Svelte stores (state)
│   │   │   └── utils/          # Helper functions
│   │   ├── routes/             # Page components (file-based routing)
│   │   │   ├── +page.svelte    # Home page
│   │   │   ├── login/          # Login flow
│   │   │   ├── register/       # Registration
│   │   │   └── dashboard/      # Protected dashboard
│   │   └── app.css             # Global styles
│   ├── tests/                  # Component & E2E tests
│   ├── Dockerfile
│   ├── .env.example
│   ├── package.json
│   └── vite.config.ts
│
├── cli/                        # Kodia CLI tool
│   ├── internal/
│   │   ├── commands/           # CLI command definitions
│   │   ├── scaffolding/        # Code generation templates
│   │   └── validation/         # Input validation
│   ├── kodia/
│   │   └── main.go
│   └── Makefile
│
├── docker-compose.yml          # Local dev environment
├── Makefile                    # Main development commands
├── README.md                   # Framework overview
├── CONTRIBUTING.md             # Contribution guidelines
├── CHANGELOG.md                # Version history
└── docs/                       # Documentation
    ├── GETTING_STARTED.md     # This file
    ├── ARCHITECTURE.md        # System design
    ├── BACKEND_GUIDE.md       # Backend development
    ├── FRONTEND_GUIDE.md      # Frontend development
    ├── CLI_GUIDE.md           # CLI commands
    ├── DEPLOYMENT.md          # Production deployment
    ├── SECURITY.md            # Security best practices
    ├── TESTING.md             # Testing strategies
    └── FAQ.md                 # Common questions
```

---

## Your First Feature

Let's create a **Blog Post** feature using Kodia's scaffolding:

### Step 1: Generate CRUD Scaffolding

```bash
# Generate complete post management feature
cd cli
go run kodia/main.go generate crud posts \
  --fields=title:string,slug:string,content:text,published:boolean \
  --with-tests \
  --with-validation

# This generates:
# ✅ Database migration: migrations/202404191200_create_posts_table.sql
# ✅ Domain model: backend/internal/core/domain/post.go
# ✅ Repository: backend/internal/adapters/repository/postgres/post_repository.go
# ✅ Service: backend/internal/core/services/post_service.go
# ✅ HTTP Handler: backend/internal/adapters/http/handlers/post_handler.go
# ✅ DTOs: backend/internal/adapters/http/dto/post_dto.go
# ✅ Validator: backend/internal/adapters/http/validators/post_validator.go
# ✅ Tests: backend/tests/integration/post_repository_test.go
# ✅ Frontend components: frontend/src/lib/components/Post*.svelte
```

### Step 2: Review Generated Code

```bash
# Check backend files
ls -la backend/internal/core/domain/post.go
ls -la backend/internal/core/services/post_service.go
cat backend/internal/adapters/http/handlers/post_handler.go

# Check frontend files
ls -la frontend/src/lib/components/PostList.svelte
ls -la frontend/src/lib/components/PostForm.svelte

# Check tests
cat backend/tests/integration/post_repository_test.go
```

### Step 3: Run Migrations

```bash
cd backend
go run cmd/server/main.go migrate
# Output: Running migration: 202404191200_create_posts_table.sql... OK
```

### Step 4: Customize the Feature

Edit the generated service to add custom business logic:

```go
// backend/internal/core/services/post_service.go

func (s *PostService) GetPublishedPosts(ctx context.Context) ([]domain.Post, error) {
    // Add custom filtering
    return s.repo.FindByStatus(ctx, "published")
}

func (s *PostService) GenerateSlug(title string) string {
    // Auto-generate URL-friendly slug
    slug := strings.ToLower(title)
    slug = strings.ReplaceAll(slug, " ", "-")
    return slug
}
```

Edit frontend components:

```svelte
<!-- frontend/src/lib/components/PostList.svelte -->
<script>
  import { onMount } from 'svelte'
  import { getAllPosts } from '$lib/api/posts'

  let posts = []
  let loading = true

  onMount(async () => {
    posts = await getAllPosts()
    loading = false
  })
</script>

<div class="posts-grid">
  {#each posts as post (post.id)}
    <article class="post-card">
      <h2>{post.title}</h2>
      <p>{post.content}</p>
    </article>
  {/each}
</div>
```

### Step 5: Test Your Feature

```bash
# Run tests
cd backend
make test-integration

# Test API manually
curl -X GET http://localhost:8080/api/posts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Or use the Swagger UI
open http://localhost:8080/api/docs
```

### Step 6: Register Routes

The scaffolding automatically adds routes to `backend/internal/adapters/http/router.go`:

```go
// Auth routes with rate limiting
api := engine.Group("/api")
{
    // Posts CRUD
    api.GET("/posts", h.postHandler.GetAll)           // Public
    api.GET("/posts/:id", h.postHandler.GetByID)      // Public
    
    protected := api.Group("")
    protected.Use(middleware.Auth(jwtManager))
    {
        protected.POST("/posts", h.postHandler.Create)   // Authenticated
        protected.PUT("/posts/:id", h.postHandler.Update)
        protected.DELETE("/posts/:id", h.postHandler.Delete)
    }
}
```

---

## Running Tests

### Backend Tests

```bash
# Run all backend tests
cd backend
make test

# Run specific test type
make test-unit              # Unit tests
make test-integration       # Database integration tests
make test-coverage          # With coverage report

# Watch mode - rerun on file changes
make test-watch

# Generate coverage report
make test-coverage
open coverage.html
```

### Frontend Tests

```bash
# Run all component tests
cd frontend
npm run test

# Run E2E tests
npm run test:e2e

# Watch mode
npm run test:watch
```

---

## Common Development Tasks

### Create Database

```bash
# Create new database
createdb kodia_dev

# Or using Docker
docker exec kodia-postgres createdb -U postgres kodia_dev
```

### Reset Database

```bash
# Drop and recreate (be careful in production!)
cd backend
make db-reset

# Or manually
dropdb kodia_dev
createdb kodia_dev
go run cmd/server/main.go migrate
```

### Access Database Shell

```bash
# PostgreSQL interactive shell
psql postgresql://postgres:password@localhost:5432/kodia

# Common commands:
# \dt          - List all tables
# \d posts     - Describe posts table
# SELECT * FROM posts;  - Query data
# \q           - Quit
```

### View Logs

```bash
# Backend logs
docker logs -f kodia-backend

# Frontend build output
cd frontend
npm run build

# See all logs
docker-compose logs -f
```

---

## Troubleshooting

### "Connection refused" on Database

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Restart services
docker-compose restart

# Check environment variables
cat backend/.env | grep DATABASE_URL
```

### "Port already in use"

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
DATABASE_PORT=5433  # Change to different port
```

### "JWT token invalid"

```bash
# Regenerate JWT secret (ensure it's 32+ characters)
# In backend/.env:
JWT_ACCESS_SECRET=your-new-secret-key-here-make-it-very-long-32-chars

# Restart backend
```

### Frontend API calls fail

```bash
# Check CORS configuration
cat backend/.env | grep CORS

# Check API base URL in frontend
cat frontend/.env | grep VITE_API_BASE_URL

# Verify backend is running
curl http://localhost:8080/api/health
```

---

## Next Steps

1. **Read the [Architecture Guide](ARCHITECTURE.md)** to understand system design
2. **Explore [Backend Development](BACKEND_GUIDE.md)** for API creation
3. **Learn [Frontend Development](FRONTEND_GUIDE.md)** for UI building
4. **Review [Security Best Practices](SECURITY.md)** before production
5. **Follow [Deployment Guide](DEPLOYMENT.md)** to go live
6. **Join our community** for support and discussions

---

## Resources

- 📖 [Full Documentation](../)
- 🎓 [Tutorial Videos](https://youtube.com/@kodia) (coming soon)
- 💬 [GitHub Discussions](https://github.com/kodia-studio/kodia/discussions)
- 🐛 [Issue Tracker](https://github.com/kodia-studio/kodia/issues)
- 📧 [Email Support](mailto:support@kodia.dev)

---

**Stuck? Check [FAQ.md](FAQ.md) or open an issue on GitHub!**
