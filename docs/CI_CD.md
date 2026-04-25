# CI/CD Pipeline

Kodia Framework uses GitHub Actions for continuous integration and continuous deployment.

---

## Overview

The CI/CD pipeline automates:
- **Testing** — Run tests on every pull request
- **Linting** — Check code style and quality
- **Building** — Build binaries and Docker images
- **Type Checking** — Verify TypeScript types
- **Releases** — Automated versioning and releases

---

## Workflows

### CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and pull request to `main` and `develop` branches.

#### Jobs

1. **Backend CI**
   - Go 1.22 setup
   - golangci-lint checks
   - Unit tests with race condition detection
   - Code coverage upload to Codecov
   - Binary build verification

2. **Frontend CI**
   - Node.js 20 setup
   - ESLint linting
   - TypeScript type checking
   - SvelteKit build verification

3. **Documentation Check**
   - Verify all required documentation files exist
   - Files checked: DEPLOYMENT.md, TESTING.md, FRONTEND_GUIDE.md, ROADMAP.md, INDEX.md

4. **Docker Build**
   - Build Docker image without pushing
   - Test Dockerfile compilation
   - Cache build layers

5. **CI Results**
   - Aggregate all check results
   - Fail if any job fails
   - Pass only if all checks succeed

#### Triggers

```yaml
on:
  push:
    branches: [main, develop]
    paths:
      - 'backend/**'
      - 'frontend/**'
      - '.github/workflows/ci.yml'
      - 'go.mod'
      - 'package.json'
      - 'pnpm-lock.yaml'
  pull_request:
    branches: [main, develop]
    # Same paths as push
```

#### Environment Variables

```bash
GO_VERSION=1.22
NODE_VERSION=20
```

---

### Release Workflow (`.github/workflows/release.yml`)

Runs on push to `main` branch when changes are detected. Can also be triggered manually.

#### Jobs

1. **Version Determination**
   - Uses semantic-release to determine next version
   - Based on conventional commits
   - Outputs: `version`, `published` flag

2. **Build Backend**
   - Compile Go server binary
   - Linux x64 architecture
   - Inject version into binary

3. **Build Docker Image**
   - Multi-stage Docker build
   - Push to GitHub Container Registry
   - Tag with version, branch, and commit SHA

4. **Create Release**
   - Generate changelog from commits
   - Create GitHub Release
   - Attach server binary
   - Link to documentation

5. **Update Changelog**
   - Prepend new version to CHANGELOG.md
   - Commit and push to main
   - Git credentials from GitHub Actions token

6. **Notify**
   - Summary of release
   - Links to artifacts

#### Version Numbering

Uses **Semantic Versioning** with **Conventional Commits**:

- `feat:` → Minor version bump (v1.0.0 → v1.1.0)
- `fix:` → Patch version bump (v1.0.0 → v1.0.1)
- `BREAKING CHANGE:` → Major version bump (v1.0.0 → v2.0.0)

Example:
```bash
git commit -m "feat: add deployment guide"    # v1.7.0 → v1.8.0
git commit -m "fix: environment variables"    # v1.7.0 → v1.7.1
git commit -m "feat!: rewrite API structure"  # v1.7.0 → v2.0.0
```

---

## Configuration Files

### Backend Linting (`.golangci.yml`)

```yaml
linters:
  enable:
    - bodyclose
    - cyclonumber
    - deadcode
    - errcheck
    - goconst
    - gocritic
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gosimple
    - govet
    - ineffassign
    - lll
    - misspell
    - staticcheck
    - unused
    - varcheck
```

Excludes test files and certain checks for specific paths.

### Frontend Linting (`.eslintrc` / `eslint.config.js`)

Uses TypeScript ESLint with SvelteKit rules:
- Check for import sorting (isort-style)
- Enforce consistent code style
- Detect unused variables
- Check async/await patterns

### Code Coverage

Coverage reports uploaded to **Codecov**:
- Backend: Go test coverage
- Frontend: Optional, can be added
- Threshold: No minimum enforced (for flexibility)

---

## Running Tests Locally

### Backend

```bash
cd backend

# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run linter
golangci-lint run

# Run specific test
go test -v -run TestName ./...
```

### Frontend

```bash
cd frontend

# Run linter
pnpm run lint

# Run type check
pnpm run check

# Build
pnpm run build

# Preview build
pnpm run preview
```

---

## Docker Image Publishing

### Docker Image Naming

```
ghcr.io/kodia-studio/kodia/backend:v1.7.0
ghcr.io/kodia-studio/kodia/backend:1.7
ghcr.io/kodia-studio/kodia/backend:main
ghcr.io/kodia-studio/kodia/backend:abc123def456 (commit SHA)
```

### Tagging Strategy

- **Latest version**: `v1.7.0`
- **Major.minor**: `1.7`
- **Branch**: `main`, `develop`
- **Commit**: Short SHA (first 12 chars)

### Pulling Images

```bash
# Latest version
docker pull ghcr.io/kodia-studio/kodia/backend:v1.7.0

# From main branch
docker pull ghcr.io/kodia-studio/kodia/backend:main

# Specific commit
docker pull ghcr.io/kodia-studio/kodia/backend:abc123def456
```

---

## GitHub Release Assets

Each release includes:

1. **Server Binary** (`kodia-server-linux-amd64`)
   - Compiled Go binary
   - Ready to deploy on Linux
   - Includes version information

2. **Docker Image**
   - Published to GitHub Container Registry
   - Can be pulled with `docker pull`

3. **Documentation Links**
   - Deployment Guide
   - Testing Guide
   - Frontend Guide
   - Roadmap

---

## Secrets & Permissions

### Required Secrets

- **`GITHUB_TOKEN`** (automatic)
  - Used for Docker login
  - Used for creating releases
  - Used for committing changelog

### Permissions

The workflows use minimal required permissions:

```yaml
permissions:
  contents: read      # For checkout
  packages: write     # For Docker push (optional)
  checks: write       # For check runs (optional)
```

---

## Troubleshooting

### Tests Failing Locally but Passing in CI

- Check Go version: `go version` should be 1.22+
- Check Node version: `node -v` should be 20+
- Clear cache: `go clean -testcache` / `pnpm store prune`
- Rebuild: `go mod tidy` / `pnpm install`

### Docker Build Failing in CI

- Check Dockerfile: `docker build -f backend/Dockerfile ./backend`
- Check base image availability
- Verify layer caching

### Release Not Being Created

- Check conventional commits: `git log --oneline | head -10`
- Verify `main` branch is the target
- Check GitHub Actions logs for errors
- Try manual trigger: Actions tab → Release → Run workflow

### Changelog Not Updating

- Verify git credentials are configured
- Check file permissions
- Verify commit message format

---

## Best Practices

✅ **Do**:
- Keep commits small and focused
- Use conventional commit format
- Write clear commit messages
- Include breaking changes in commit footer
- Run tests before pushing
- Review CI logs on failure

❌ **Don't**:
- Force push to main (breaks release history)
- Bypass CI checks
- Ignore test failures
- Commit sensitive data
- Mix feature and fix commits

---

## Monitoring

### GitHub Actions Dashboard

Monitor workflow status at:
```
https://github.com/kodia-studio/kodia/actions
```

### Workflow Runs

Each workflow run shows:
- ✅ Passed jobs (green)
- ❌ Failed jobs (red)
- ⏳ In-progress jobs (yellow)
- ⊘ Skipped jobs (gray)

### Logs

Click on job to view:
- Step-by-step execution
- Error messages
- Timing information
- Artifacts

---

## Next Steps

1. **Ensure conventions**: Write commits using conventional format
2. **Run tests locally**: Before pushing to ensure CI passes
3. **Monitor releases**: Check Actions tab for release progress
4. **Use images**: Pull Docker images from GitHub Container Registry

---

**Last Updated**: April 2026  
**Framework Version**: v1.7.0+
