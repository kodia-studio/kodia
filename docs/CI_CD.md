# CI/CD Pipeline

Kodia Framework uses GitHub Actions for automated testing, security scanning, building, and deployment.

---

## Overview

The CI/CD pipeline automates:
- **Static Analysis** — Go vet, linting, type checking
- **Security Scanning** — Vulnerability detection (Go modules, npm dependencies)
- **Testing** — Unit tests with race condition detection
- **Building** — Binary compilation and multi-platform Docker images
- **Releases** — Automated versioning, changelog, and artifact publishing

### Workflow Diagram

```
Push/PR to main/develop
    ↓
[CI Workflow] ← Runs on every commit
├─ Backend (lint, test, build)
├─ Frontend (lint, type-check, build)
├─ Documentation (verify existence)
└─ Docker (build without push)
    ↓
[All checks pass?]
├─ Yes → Ready to merge
└─ No → Review logs and fix
    ↓
[Push to main] ← Triggers release
    ↓
[Security Workflow] ← Parallel with CI
├─ Go vulnerability scan (govulncheck)
├─ Frontend audit (npm audit)
└─ Dependency review (PRs only)
    ↓
[Release Workflow] ← After CI passes
├─ Determine version (semantic-release)
├─ Build backend binary
├─ Build Docker images (amd64, arm64)
├─ Create GitHub Release
└─ Update CHANGELOG.md
```

---

## Workflows

### 1. CI Workflow (`.github/workflows/ci.yml`)

Continuous integration for every push and pull request.

**Triggers:**
- Push to `main`, `develop` branches (path-filtered)
- Pull requests to `main`, `develop` branches
- Path filters: `backend/**`, `frontend/**`, `cli/**`, `go.mod`, `go.sum`, `package.json`, `pnpm-lock.yaml`

**Environment:**
```bash
GO_VERSION=1.25
NODE_VERSION=22
```

#### Backend CI Job

```bash
✓ Go 1.25 setup
✓ go vet ./...                           # Code quality checks
✓ golangci-lint (with 5m timeout)        # Comprehensive linting
✓ Race detection tests                   # Concurrent access bugs
✓ Code coverage analysis                 # Coverage metrics
✓ Binary build verification              # Compile test
```

**Test Command:**
```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./pkg/... ./internal/...
```

Note: Tests scoped to `./pkg/...` and `./internal/...` to run unit tests only (excludes integration tests requiring external services).

#### Frontend CI Job

```bash
✓ Node 22 setup with npm cache
✓ npm ci (clean install)                 # Lockfile-based installation
✓ ESLint (npm run lint)                  # Code style
✓ TypeScript (npm run check)             # Type safety
✓ SvelteKit build (npm run build)        # Build verification
```

Frontend uses **npm** (package-lock.json) with caching for faster installs.

#### Documentation Check Job

Verifies required documentation files exist:
- `docs/DEPLOYMENT.md` — Deployment instructions
- `docs/TESTING.md` — Testing guide
- `docs/FRONTEND_GUIDE.md` — Frontend development
- `docs/ROADMAP.md` — Project roadmap
- `docs/INDEX.md` — Documentation index
- `docs/CI_CD.md` — CI/CD pipeline guide
- `docs/OPENAPI_GUIDE.md` — API documentation guide

#### Docker Build Job

```bash
✓ Set up Docker Buildx
✓ Build (no push) for testing
✓ Layer caching via GitHub Actions cache
```

Builds only test image to verify Dockerfile validity without publishing.

#### CI Results Job

Aggregates all job results:
- ✅ All pass → "Ready to merge"
- ❌ Any fail → Shows which job failed

---

### 2. Security Workflow (`.github/workflows/security.yml`)

Vulnerability scanning for Go modules and npm dependencies.

**Triggers:**
- Weekly: Monday 2:00 AM UTC
- Push to `main`, `develop` branches (path-filtered)
- Pull requests to `main`, `develop` branches
- Manual trigger via Actions tab

#### Go Vulnerability Scan Job

```bash
✓ Install govulncheck
✓ Scan Go module dependencies: govulncheck ./...
✓ Reports: CVEs, severity, fix recommendations
```

**Coverage:** All Go modules in `go.mod` and their dependencies.

#### Frontend Dependency Audit Job

```bash
✓ Node 22 setup
✓ npm ci (install from lockfile)
✓ npm audit --audit-level=moderate
✓ Report critical vulnerabilities and fail if found
```

**Audit Levels:**
- `low` — Development warnings
- `moderate` — Should fix soon
- `high` — Fix before deployment
- `critical` — Fix immediately

#### Dependency Review Job (PRs only)

```bash
✓ GitHub native dependency review action
✓ Blocks PRs introducing vulnerable dependencies
✓ Requires: fail-on-severity: moderate
```

Only runs on pull requests, not on push events.

---

### 3. Release Workflow (`.github/workflows/release.yml`)

Automated versioning, building, and publishing to GitHub Releases.

**Triggers:**
- Push to `main` branch when specific paths change
  - `backend/**`, `frontend/**`, `cli/**`, `docs/**`, `go.mod`, `CHANGELOG.md`
- Manual trigger via Actions tab with optional version override

**Path Filters:** Prevents unnecessary releases for unrelated changes.

#### Version Determination Job

```bash
✓ Install semantic-release + plugins
✓ Run dry-run to detect next version
✓ Parse conventional commits
✓ Extract version from output
```

**Semantic Versioning Rules:**

| Commit Type | Version Change | Example |
|---|---|---|
| `feat:` | Minor bump | v1.0.0 → v1.1.0 |
| `fix:` | Patch bump | v1.0.0 → v1.0.1 |
| `refactor:`, `docs:` | No version change | v1.0.0 → v1.0.0 |
| `BREAKING CHANGE:` in footer | Major bump | v1.0.0 → v2.0.0 |

**Example Commits:**
```bash
git commit -m "feat: add database migrations CLI"     # → v1.1.0
git commit -m "fix: config loading issue"            # → v1.0.1
git commit -m "feat!: rewrite auth system"           # → v2.0.0 (BREAKING)
git commit -m "docs: update README"                  # → no release
```

#### Build Backend Binary Job

```bash
✓ Go 1.25 setup
✓ Compile Linux binary (amd64)
✓ Inject version: -ldflags "-X main.Version=v1.8.0"
✓ Output: kodia-server-linux-amd64
✓ Upload as artifact for release
```

#### Build Docker Images Job

```bash
✓ Docker Buildx setup
✓ Multi-platform build: linux/amd64, linux/arm64
✓ Push to GitHub Container Registry (ghcr.io)
```

**Image Tags:**
```
ghcr.io/kodia-studio/kodia/backend:v1.8.0              # Semantic version
ghcr.io/kodia-studio/kodia/backend:1.8                 # Major.minor
ghcr.io/kodia-studio/kodia/backend:main                # Branch
ghcr.io/kodia-studio/kodia/backend:abc123def456        # Commit SHA
```

#### Create GitHub Release Job

```bash
✓ Download artifacts
✓ Generate changelog from commits since last tag
✓ Create GitHub Release with:
  - Release notes
  - Binary attachment
  - Docker image links
  - Documentation links
✓ Publish (not draft)
```

**Release Body Includes:**
- Version number
- Change summary (commits since last release)
- Binary download link
- Docker image pull command
- Documentation links

#### Update Changelog Job

```bash
✓ Prepend new version entry to CHANGELOG.md
✓ Include commits since last tag
✓ Commit with message: "chore: update changelog for v1.8.0"
✓ Push to main branch
```

#### Notify Job

```bash
✓ Summary message in workflow logs
✓ Links to artifacts and release page
✓ Confirms: docker image, binary, documentation
```

---

## Configuration & Versions

### Version Requirements

| Component | Version | Required |
|---|---|---|
| Go | 1.25.0+ | ✅ Backend builds |
| Node.js | 22+ | ✅ Frontend builds |
| npm | v10+ | ✅ Frontend package manager |
| Docker | 20.10+ | For local Docker builds |

### Backend Linting (`.golangci.yml`)

Configured linters detect:
- Missing error handling
- Unreachable code
- Style violations
- Import organization
- Variable shadowing
- Race conditions
- Security issues

Run locally:
```bash
golangci-lint run ./backend
```

### Frontend Linting (`.eslintrc` / `eslint.config.js`)

TypeScript ESLint with SvelteKit rules:
- Import organization
- Consistent naming
- Unused variables
- Async/await patterns
- Type safety

Run locally:
```bash
cd frontend && npm run lint
```

### Code Coverage

Coverage uploaded to **Codecov** on each CI run:
- Backend: Go test coverage (atomic mode)
- Files: `./pkg/...` and `./internal/...`
- Threshold: No minimum (for flexibility)

---

## Running Tests Locally

### Backend

```bash
cd backend

# Lint
go vet ./...
golangci-lint run

# Run tests (unit only, same as CI)
go test -v -race ./pkg/... ./internal/...

# Run with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./pkg/... ./internal/...
go tool cover -html=coverage.out

# Run integration tests (requires postgres + redis)
go test -v ./tests/...

# Run specific test
go test -v -run TestName ./...

# Run benchmarks
go test -bench=. ./pkg/...
```

**Note:** CI runs `./pkg/...` and `./internal/...` to exclude integration tests. Integration tests require external services (PostgreSQL, Redis).

### Frontend

```bash
cd frontend

# Install dependencies
npm ci

# Lint
npm run lint

# Type check
npm run check

# Build
npm run build

# Preview production build
npm run preview

# Development with hot reload
npm run dev
```

### Database Migrations

```bash
cd backend

# Check migration status
make migrate-status

# Run pending migrations
go run cmd/migrate/main.go up

# Rollback last batch
make migrate-rollback

# Reset (rollback all + re-run)
make migrate-reset

# Fresh (drop all + re-run)
make migrate-fresh

# Create new migration
make migration-create NAME=create_posts_table
```

---

## Docker Images

### Multi-Platform Support

Release workflow builds for:
- **linux/amd64** — Intel/AMD 64-bit servers, desktops
- **linux/arm64** — ARM 64-bit (M1/M2 Mac, Raspberry Pi 4+, AWS Graviton)

Both platforms published under the same tags.

### Image Registry

Published to **GitHub Container Registry (ghcr.io)**:
```
ghcr.io/kodia-studio/kodia/backend
```

**Authentication:**
```bash
# Login with GitHub personal access token
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin

# Or with read:packages scope in PAT
cat ~/path/to/token.txt | docker login ghcr.io -u USERNAME --password-stdin
```

### Image Tags

| Tag | When | Example |
|---|---|---|
| `v1.8.0` | Version release | Semantic version |
| `1.8` | Version release | Major.minor version |
| `main` | Push to main branch | Latest main build |
| `develop` | Push to develop branch | Latest develop build |
| `abc123de` | Every push | Commit SHA (12 chars) |

### Pulling Images

```bash
# Latest version
docker pull ghcr.io/kodia-studio/kodia/backend:v1.8.0

# From main branch
docker pull ghcr.io/kodia-studio/kodia/backend:main

# Specific commit
docker pull ghcr.io/kodia-studio/kodia/backend:abc123de

# Latest from develop
docker pull ghcr.io/kodia-studio/kodia/backend:develop
```

### Running Locally

```bash
# Using docker compose (includes frontend, postgres, redis)
docker compose up -d

# Access services
# Backend:  http://localhost:8080
# Frontend: http://localhost:3000
# API docs: http://localhost:8080/api/docs/
```

### Building Locally

```bash
# Build single platform
docker build -f backend/Dockerfile -t kodia:local ./backend

# Build multi-platform (requires buildx)
docker buildx build \
  -f backend/Dockerfile \
  --platform linux/amd64,linux/arm64 \
  -t kodia:local \
  ./backend
```

---

## GitHub Releases

### Release Content

Each release includes:

1. **Release Notes**
   - Version number (e.g., v1.8.0)
   - Change summary (commits since last release)
   - Links to full changelog

2. **Artifacts**
   - `kodia-server-linux-amd64` — Go binary for deployment
   - `kodia-server-linux-arm64` — ARM64 binary (optional)

3. **Docker Images**
   - Multi-platform builds in GitHub Container Registry
   - Accessible via tags (v1.8.0, 1.8, main, develop)

4. **Documentation**
   - [Deployment Guide](../DEPLOYMENT.md)
   - [Testing Guide](../TESTING.md)
   - [Frontend Guide](../FRONTEND_GUIDE.md)
   - [Roadmap](../ROADMAP.md)
   - [API Documentation](../OPENAPI_GUIDE.md)

### Accessing Releases

- **GitHub:** https://github.com/kodia-studio/kodia/releases
- **Latest:** https://github.com/kodia-studio/kodia/releases/latest
- **Assets:** Each release page shows downloadable binaries

### Downloading Binaries

```bash
# Download latest release
curl -LO https://github.com/kodia-studio/kodia/releases/download/v1.8.0/kodia-server-linux-amd64

# Make executable
chmod +x kodia-server-linux-amd64

# Run
./kodia-server-linux-amd64
```

---

## Secrets & Permissions

### Required Secrets

**`GITHUB_TOKEN`** (automatic)
- Provided by GitHub Actions automatically
- Scope: `contents:write` (for releases), `packages:write` (for Docker)
- Used by:
  - Docker login (registry: ghcr.io)
  - Creating GitHub releases
  - Committing changelog updates

No manual secrets needed for standard CI/CD.

### Optional Secrets

For external integrations (optional):
- `CODECOV_TOKEN` — Code coverage (if not using auto-detection)
- `SENTRY_AUTH_TOKEN` — Error tracking integration

### Workflow Permissions

Set at repository level:
```
Settings → Actions → General → Workflow permissions
```

Required:
- ✅ Read and write permissions
- ✅ Allow GitHub Actions to create pull requests (for changelog commits)

---

## Workflow Status & Monitoring

### GitHub Actions Dashboard

Monitor all workflows:
```
https://github.com/kodia-studio/kodia/actions
```

**Status Indicators:**
- 🟢 **All jobs passed** — Ready to merge/release
- 🟡 **In progress** — Workflows running
- 🔴 **Failed** — Check logs for errors
- ⚪ **Skipped** — Conditions not met

### Workflow Summary

Click on workflow run to see:
- **Job list** with status (✅/❌/⏳)
- **Step details** — What each step did
- **Logs** — Full command output and errors
- **Artifacts** — Uploaded files (coverage, binaries)
- **Timing** — Duration of each job
- **Annotations** — Linting warnings/errors

### Debugging Failed Workflows

1. Click the failed workflow run
2. Click the failed job
3. Click the failed step
4. Review the error output
5. Check the annotations for specific issues

### Common Failures

| Symptom | Likely Cause | Fix |
|---|---|---|
| Test fails locally, passes in CI | Dependency version | Run `go mod tidy`, pin versions |
| Linting error in CI only | Version difference | Update linter: `golangci-lint upgrade` |
| Docker build fails in CI | Missing file/layer | Check Dockerfile paths, verify .dockerignore |
| Release not created | No version bump | Check conventional commits, ensure `git log` readable |
| Coverage not uploading | Codecov token | Check token scope, verify path in upload step |

---

## Troubleshooting

### ❌ Tests Pass Locally but Fail in CI

**Problem:** Tests work on your machine but fail in GitHub Actions.

**Causes & Solutions:**
| Cause | Solution |
|---|---|
| Go version mismatch | Run `go version` — must be 1.25+ |
| Node version mismatch | Run `node -v` — must be 22+ |
| Stale test cache | Run `go clean -testcache` |
| Missing dependencies | Run `go mod tidy` then `go mod download` |
| npm cache stale | Delete `node_modules/`, run `npm ci` |
| Database not initialized | Ensure PostgreSQL/Redis running locally |

### ❌ Linting Errors in CI Only

**Problem:** Code passes local linting but fails in CI.

**Causes & Solutions:**
| Cause | Solution |
|---|---|
| Different linter version | Update local: `golangci-lint upgrade` |
| Missing .golangci.yml | Check file exists in `./backend/` |
| Import ordering | Run `goimports -w ./backend` |
| Formatting | Run `go fmt ./...` |

### ❌ Docker Build Fails in CI

**Problem:** Docker builds locally but fails in GitHub Actions.

**Diagnose:**
```bash
# Build locally with same Dockerfile
docker build -f backend/Dockerfile ./backend

# Check base image availability
docker pull golang:1.25-alpine
docker pull alpine:3.19

# Verify .dockerignore isn't excluding needed files
cat backend/.dockerignore
```

**Common Issues:**
- Base image tag changed (golang:1.25 vs 1.22)
- Missing files in context (check .dockerignore)
- Build args not passed
- Layer caching issues

### ❌ Release Not Created

**Problem:** Pushed to main but no release was created.

**Check:**
```bash
# Verify conventional commits
git log --oneline | head -10
# Should show commit messages like "feat:", "fix:", etc.

# Check if any change paths match release triggers
git show --name-only
# Should include: backend/**, frontend/**, go.mod, etc.

# Manual trigger if needed
# GitHub → Actions → Release → Run workflow
```

**Common Issues:**
- Commit message doesn't follow conventional format
- Changed files don't match path filters
- GitHub Actions logs show errors

**Fix:**
1. Make new commit with proper format: `git commit -m "feat: description"`
2. Push to main
3. Check Actions tab for new run
4. If still failing, check Actions → Release → logs

### ❌ Changelog Not Updating

**Problem:** Release created but CHANGELOG.md wasn't updated.

**Check:**
```bash
# Verify file exists
test -f CHANGELOG.md && echo "✓ File exists" || echo "✗ Missing"

# Check permissions
ls -l CHANGELOG.md

# Verify git config
git config user.name
git config user.email
```

**Common Issues:**
- File doesn't exist (create empty CHANGELOG.md)
- Wrong file path (should be at repo root)
- Git credentials misconfigured

### ❌ Security Scan False Positives

**Problem:** Vulnerability reported that's actually patched/ignored.

**For Go:**
```bash
# Run govulncheck locally
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Check if dependency needs update
go get -u github.com/package@latest
go mod tidy
```

**For npm:**
```bash
# Run audit locally
cd frontend
npm audit

# Update dependencies
npm update

# Or ignore known false positives
npm audit --json > audit.json
```

---

## Best Practices

### ✅ Commit Conventions

**Use conventional commit format:**
```bash
feat: add user authentication          # Minor version bump
fix: resolve race condition           # Patch version bump
refactor: simplify config loading     # No version bump
docs: update deployment guide         # No version bump
chore: update dependencies            # No version bump
feat!: rewrite API structure          # Major version bump (BREAKING)
```

**Structure:**
```
<type>[scope]: <description>

[optional body with more details]

[optional footer: BREAKING CHANGE: description]
```

**Examples:**
```bash
# Simple fix
git commit -m "fix: handle null pointer in user handler"

# Feature with scope
git commit -m "feat(auth): add two-factor authentication"

# Breaking change
git commit -m "feat!: change API response format

BREAKING CHANGE: /api/users response now returns 'userId' instead of 'id'
"
```

### ✅ Pre-Push Checklist

Before pushing to main:
1. ✓ Run tests locally: `go test -v -race ./pkg/... ./internal/...`
2. ✓ Run linter: `golangci-lint run ./backend`
3. ✓ Check frontend: `npm run lint && npm run check`
4. ✓ Build locally: `docker build -f backend/Dockerfile ./backend`
5. ✓ Verify commit message format
6. ✓ Update CHANGELOG.md manually if needed

### ❌ Don'ts

- 🚫 **Force push to main** — Breaks release history and GitHub Actions state
- 🚫 **Bypass CI checks** — Merge without passing tests
- 🚫 **Commit secrets** — API keys, tokens, credentials
- 🚫 **Large binary files** — Use Git LFS or remove before commit
- 🚫 **Mixed commits** — Don't mix features, fixes, and refactors in one commit
- 🚫 **Edit CHANGELOG.md in release PR** — Release workflow handles it
- 🚫 **Ignore linting warnings** — Fix the underlying issues

---

## Examples

### Creating a Feature Release

```bash
# 1. Create and test feature branch
git checkout -b feature/new-api-endpoint

# 2. Make changes, test locally
npm run lint
go test -v ./...

# 3. Commit with conventional format
git commit -m "feat: add new user list endpoint with pagination"

# 4. Push and create PR
git push origin feature/new-api-endpoint
# → Create PR on GitHub

# 5. PR checks run automatically (CI workflow)
# → If all pass, can merge

# 6. Push to main triggers Release workflow
git checkout main
git merge feature/new-api-endpoint
git push origin main
# → GitHub Actions:
#    1. Detects "feat:" commit → version bump
#    2. Builds binaries and Docker images
#    3. Creates v1.8.0 release (bumped from v1.7.0)
#    4. Updates CHANGELOG.md
#    5. Publishes Docker images to ghcr.io
```

### Hotfix Process

```bash
# 1. Create hotfix branch from main (not develop)
git checkout -b hotfix/critical-bug main

# 2. Fix the issue
git commit -m "fix: prevent memory leak in connection pool"

# 3. Test thoroughly
go test -v ./...

# 4. Merge to main and develop
git push origin hotfix/critical-bug
# → Create PR, ensure tests pass

git checkout main
git merge hotfix/critical-bug
git push origin main
# → Triggers Release (patch version: v1.7.1)

git checkout develop
git merge hotfix/critical-bug
git push origin develop
```

---

## Workflow Files

### File Locations

| Workflow | File | Triggers |
|---|---|---|
| CI | `.github/workflows/ci.yml` | Push/PR to main, develop |
| Security | `.github/workflows/security.yml` | Push/PR, weekly, manual |
| Release | `.github/workflows/release.yml` | Push to main, manual |

### Modifying Workflows

**Warning:** Workflows are critical. Before changing:

1. Test changes locally (if possible)
2. Review changes carefully
3. Create PR for approval before merging
4. Document any new environment variables or secrets
5. Update this documentation

---

## Performance Tips

### Faster CI Runs

1. **Use path filters** — Prevents unnecessary runs
2. **Cache dependencies** — Configured automatically
3. **Parallel jobs** — CI runs backend + frontend simultaneously
4. **Skip steps** — Use `if: conditions` to skip unnecessary steps

### Caching Strategy

| Cache | Hit Rate | Size |
|---|---|---|
| Go modules | 95%+ | ~500MB |
| npm cache | 90%+ | ~200MB |
| Docker layers | 80%+ | ~2GB |

Cache is automatically managed by GitHub Actions.

---

**Last Updated**: April 2026  
**Framework Version**: v1.8.0+  
**Go Version**: 1.25  
**Node Version**: 22
