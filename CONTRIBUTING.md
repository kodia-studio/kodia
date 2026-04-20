# 🤝 Contributing to Kodia Framework

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to Kodia.

**Table of Contents:**
- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)

---

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inspiring community for all. Please read and adhere to our [Code of Conduct](CODE_OF_CONDUCT.md).

**In summary:**
- ✅ Be respectful and inclusive
- ✅ Welcome different perspectives
- ✅ Provide constructive feedback
- ❌ No harassment or discrimination
- ❌ No spam or advertising

---

## Getting Started

### 1. Fork the Repository

```bash
# Click "Fork" button on GitHub
# Then clone YOUR fork
git clone https://github.com/YOUR-USERNAME/kodia.git
cd kodia

# Add upstream remote
git remote add upstream https://github.com/kodia-studio/kodia.git
```

### 2. Create a Feature Branch

```bash
# Update from upstream
git fetch upstream
git checkout main
git merge upstream/main

# Create feature branch
git checkout -b feature/my-feature
# or
git checkout -b fix/bug-description
```

### 3. Make Your Changes

```bash
# Edit files
# Test your changes
# Commit with clear messages

# Keep your branch up to date
git fetch upstream
git rebase upstream/main
```

---

## Development Setup

### Prerequisites

- Go 1.26+
- Node.js 25+
- Docker & Docker Compose
- PostgreSQL 15+
- Redis 7+

### Install Dependencies

```bash
# Backend
cd backend && go mod download

# Frontend
cd ../frontend && npm install

# CLI
cd ../cli && go mod download
```

### Setup Environment

```bash
# Copy environment templates
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Start Docker services
docker-compose up -d
```

### Verify Setup

```bash
# Run tests
make test

# Start dev servers
make dev

# Should see:
# ✅ Backend running on :8080
# ✅ Frontend running on :5173
# ✅ All services healthy
```

---

## Making Changes

### Choose an Issue Type

Before starting, check if an issue exists:

1. **Bug Fix** - Fix existing issue
2. **Feature** - Add new functionality  
3. **Documentation** - Improve docs
4. **Performance** - Optimize code

### Work on the Right Component

**Backend Changes** (`backend/`):
```bash
cd backend

# Run backend tests
go test ./...

# Build
go build -o ./server cmd/server/main.go

# Check formatting
go fmt ./...

# Run linter
golangci-lint run
```

**Frontend Changes** (`frontend/`):
```bash
cd frontend

# Run component tests
npm run test

# Build
npm run build

# Check formatting
npm run format

# Lint
npm run lint
```

**CLI Changes** (`cli/`):
```bash
cd cli

# Run tests
go test ./...

# Build
go build -o kodia kodia/main.go

# Test command
./kodia --help
```

---

## Commit Guidelines

### Commit Message Format

Use clear, descriptive commit messages:

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation
- **style**: Formatting (no code change)
- **refactor**: Code restructure (no behavior change)
- **perf**: Performance improvement
- **test**: Add or update tests
- **chore**: Build, dependencies, etc.

### Scope

- **backend**: Backend-specific changes
- **frontend**: Frontend-specific changes
- **cli**: CLI-specific changes
- **docs**: Documentation changes
- **core**: Core framework changes

### Examples

```bash
git commit -m "feat(backend): Add post publish endpoint"
git commit -m "fix(frontend): Fix login form validation error"
git commit -m "docs(security): Add OWASP security guidelines"
git commit -m "test(cli): Add scaffolding command tests"
```

### Body

Explain **what** and **why**, not **how**:

```
feat(backend): Add email verification

Users can now verify their email address after registration.
This adds an extra security layer to prevent account takeover.

- Added verification_token column to users table
- Added verify_email endpoint
- Send verification email on registration
- Add migration: 202404190000_add_email_verification.sql

Fixes #1234
```

---

## Pull Request Process

### 1. Create Pull Request

```bash
# Push your branch
git push origin feature/my-feature

# Open PR on GitHub with:
# - Clear title
# - Description of changes
# - Link to related issues
# - Testing instructions
```

### 2. PR Description Template

```markdown
## Description
Brief description of changes.

## Related Issues
Fixes #123
Related to #456

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation
- [ ] Breaking change

## Testing
How to test these changes:
1. Step 1
2. Step 2
3. Step 3

## Checklist
- [ ] Tests pass
- [ ] Code follows style guide
- [ ] Documentation updated
- [ ] No breaking changes
```

### 3. Code Review

- ✅ Expect feedback and questions
- ✅ Respond respectfully to reviews
- ✅ Make requested changes
- ✅ Push new commits (don't force push)
- ❌ Don't take feedback personally

### 4. Merge

Once approved by maintainers:
- PRs will be squashed and merged
- Commit message will be cleaned up
- Branch will be deleted

---

## Coding Standards

### Go Code Style

Follow [Effective Go](https://golang.org/doc/effective_go):

```go
// ✅ Good: Clear, idiomatic Go
func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
    post := domain.NewPost(req.Title, req.Content)
    
    if err := s.repo.Save(ctx, post); err != nil {
        return nil, fmt.Errorf("failed to save post: %w", err)
    }
    
    return post, nil
}

// ❌ Bad: Unclear naming, inefficient
func (s *PostService) c(ctx context.Context, r *CreatePostRequest) (*Post, error) {
    p := &Post{}
    p.Title = r.Title
    // ...
}
```

### JavaScript/TypeScript Style

Follow [SvelteKit best practices](https://kit.svelte.dev/docs):

```typescript
// ✅ Good: Type-safe, readable
interface Post {
  id: string
  title: string
  content: string
}

export async function load({ fetch }) {
  const response = await fetch('/api/posts')
  const posts: Post[] = await response.json()
  return { posts }
}

// ❌ Bad: Any types, unclear
const load = async (ctx) => {
  const resp = await ctx.fetch('/api/posts')
  const data = await resp.json()
  return { data }
}
```

### Documentation

```go
// ✅ Good: Exported functions have comments
// PostService handles all post-related business logic
type PostService struct {
    repo   ports.PostRepository
    cache  ports.CacheProvider
}

// CreatePost creates a new post in the database
func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
    // ...
}

// ❌ Bad: Missing documentation
type PostService struct {
    repo ports.PostRepository
}

func (s *PostService) CreatePost(ctx context.Context, req *CreatePostRequest) (*Post, error) {
}
```

---

## Testing Requirements

### Backend Tests

All backend code must have tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

**Test file locations:**
- Service tests: `tests/services/`
- Repository tests: `tests/repositories/`
- Handler tests: `tests/handlers/`
- Integration tests: `tests/integration/`

**Minimum coverage: 80%**

### Frontend Tests

Test all components and pages:

```bash
# Run component tests
npm run test

# Run E2E tests
npm run test:e2e

# Check coverage
npm run test:coverage
```

### Example Test

```go
// backend/tests/services/post_service_test.go
func TestCreatePost(t *testing.T) {
    // Arrange
    mockRepo := &MockPostRepository{}
    mockCache := &MockCache{}
    service := NewPostService(mockRepo, mockCache)
    
    req := &CreatePostRequest{
        Title:   "Test Post",
        Content: "Test content",
    }
    
    // Act
    post, err := service.CreatePost(context.Background(), req)
    
    // Assert
    require.NoError(t, err)
    require.Equal(t, "Test Post", post.Title)
    require.True(t, mockRepo.SaveCalled)
}
```

---

## Documentation

### Update Relevant Docs

When making changes, update documentation:

| Change | Documentation |
|--------|---------------|
| New feature | Update FEATURES.md |
| API change | Update backend/docs/ |
| CLI command | Update CLI_GUIDE.md |
| Config option | Update config examples |
| Breaking change | Update CHANGELOG.md |

### Document New Features

```markdown
## New Feature: Post Publishing

Posts can now be published with a scheduled date.

### Usage

```bash
kodia generate crud posts --fields=published_at:datetime
```

### Example

```go
post.PublishAt(time.Now().Add(24 * time.Hour))
```
```

---

## Review Criteria

Your PR will be merged if it:

- ✅ Follows commit guidelines
- ✅ Includes tests (80%+ coverage)
- ✅ Passes all CI checks
- ✅ Updates documentation
- ✅ Has clear PR description
- ✅ Addresses reviewer feedback
- ✅ Includes no breaking changes (without discussion)

---

## Common Mistakes to Avoid

- ❌ Large PRs (>500 lines) - Split into smaller PRs
- ❌ Missing tests - Every change needs tests
- ❌ Weak commit messages - Be descriptive
- ❌ Ignored linter warnings - Fix all warnings
- ❌ Breaking changes without discussion - Talk to maintainers first
- ❌ Merge main into your branch - Use rebase instead
- ❌ Force push to PR branch - Avoid rewriting history

---

## Getting Help

- 📖 [Documentation](docs/)
- 💬 [GitHub Discussions](https://github.com/kodia-studio/kodia/discussions)
- 🐛 [Issue Tracker](https://github.com/kodia-studio/kodia/issues)
- 💬 [Discord Community](https://discord.gg/kodia) (coming soon)

---

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- GitHub repository

Thank you for contributing! 🎉
