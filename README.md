# 🚀 Kodia Framework

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Svelte Version](https://img.shields.io/badge/Svelte-5-FF3E00?style=flat&logo=svelte)](https://svelte.dev/)
[![Tailwind Version](https://img.shields.io/badge/Tailwind-4-38B2AC?style=flat&logo=tailwind-css)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GitHub Stars](https://img.shields.io/github/stars/kodia-studio/kodia?style=social)](https://github.com/kodia-studio/kodia)

**Kodia** is an **opinionated, batteries-included fullstack framework** designed for rapid development, security-first architecture, and exceptional developer experience (DX).

Built with the power of **Go/Gin** on the backend, **SvelteKit** on the frontend, and a **powerful CLI tool** for instant scaffolding—Kodia brings the productivity of Laravel to the Go ecosystem.

---

| Feature | Description |
|---------|-------------|
| 🏗️ **Clean Architecture** | Backend follows SOLID principles with clear separation of concerns |
| 🔐 **Enterprise Security** | 2FA (TOTP), RBAC & ABAC Policy Engine, Token Rotation, CSRF, XSS protection |
| 🎨 **Elite Design System** | High-fidelity glassmorphism, **bg-hive** patterns, and Tailwind v4 |
| 📊 **Kodia UI Suite v2** | Pre-wired Svelte 5 components with artisanal focus and premium animations |
| 🔍 **Kodia Search** | Official Meilisearch plugin with async background indexing |
| 🌍 **Localization & Tenancy** | SaaS-ready with automated database isolation and reactive i18n |
| 📊 **Full Observability** | OpenTelemetry, Prometheus, Sentry, pprof profiling |
| ⚡ **Developer Experience** | `kodia` CLI for instant feature scaffolding and ecosystem management |
| 📦 **Batteries Included** | Auth, CRUD, middleware, validation, file uploads, email—all out-of-the-box |
| 🛳️ **Kodia Sail** | One-command Docker infrastructure (Postgres, Redis, Meilisearch, Mailpit) |
| 📊 **Kodia Pulse** | Real-time administrative monitoring dashboard (CPU, RAM, Logs) |
| 🎨 **Modern Frontend** | SvelteKit 5, Tailwind CSS v4, premium component library included |
| 🐳 **Production Ready** | Docker, CI/CD templates, health checks, structured logging |

---

## 🎨 Elite Design System

Kodia v2 introduces the **Elite Design System**, a high-fidelity visual language meticulously synchronized with the `kodia.id` design authority.

- **Artisanal Glassmorphism**: Advanced `backdrop-blur` and artisanal border transitions for a premium, institutional feel.
- **Hive Architecture**: Signature `bg-hive` patterns and holographic radial glows across all layouts.
- **Premium Auth Suite**: Redesigned Login, Register, and 2FA stages that prioritize security without sacrificing "wow" factor.
- **Standardized Navigation**: All internal documentation and community links are synchronized to the central `kodia.id` portal.

---

## 🎯 Philosophy

> **"Opinionated by default, flexible by design."**

Like Laravel, Kodia provides a highly structured framework with best practices baked in—but remains flexible enough to adapt to your needs.

- **Convention over Configuration** - Sensible defaults reduce boilerplate
- **Batteries Included** - Everything you need from day one
- **Developer Friendly** - Rapid scaffolding and intuitive APIs
- **Production Grade** - Built with observability and security from the start
- **Type Safe** - Leverage Go's strong typing and SvelteKit's type safety

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.26 or higher
- **Node.js** 25 or higher  
- **Docker** & **Docker Compose**
- **PostgreSQL** 15+ (or MySQL 8+)

### 1️⃣ Installation

The easiest way to get started is by installing our CLI:

```bash
# Install Kodia CLI
go install github.com/kodia-studio/cli/kodia@latest

# Create a new project
kodia new my-app
cd my-app
```

### 2️⃣ Run Infrastructure (Sail)

Launch your entire development stack with zero local configuration:

```bash
# Start PostgreSQL, Redis, Meilisearch, and Mailpit
kodia sail up

# In separate terminals:
kodia dev
```

### 3️⃣ Check Health

```bash
curl http://localhost:5173              # Frontend (SvelteKit)
curl http://localhost:8080/api/health   # Backend (Gin)
curl http://localhost:8025              # Mailpit Dashboard
```

### 3️⃣ Development

```bash
# Start development servers
make dev

# In separate terminal, watch the CLI tool
cd cli && make watch
```

**Access the app:**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- API Docs: http://localhost:8080/api/docs (Swagger)
- PostgreSQL: localhost:5432

### 4️⃣ Generate Your First Feature

```bash
# Generate a CRUD feature
kodia generate crud posts \
  --fields=title:string,content:text,author:references:users \
  --with-tests \
  --with-migrations \
  --with-validation

# This generates:
# - Database migration
# - Model (domain entity)
# - Repository (data layer)
# - Service (business logic)
# - Handler (HTTP endpoints)
# - Request/Response DTOs
# - Validation rules
# - Unit tests
# - Frontend components
```

---

## 📂 Project Structure

```
kodia/
├── backend/                  # Go/Gin REST API
│   ├── cmd/server/          # Application entrypoint
│   ├── internal/
│   │   ├── core/            # Business logic (services, domain)
│   │   ├── adapters/        # HTTP handlers, repositories
│   │   └── infrastructure/  # DB, cache, mail, storage
│   ├── pkg/                 # Shared packages
│   ├── tests/               # Test suites
│   ├── docs/                # API documentation (generated)
│   └── Dockerfile
│
├── frontend/                # SvelteKit + Tailwind application
│   ├── src/
│   │   ├── lib/            # Reusable components & utilities
│   │   ├── routes/         # Page components (file-based routing)
│   │   └── styles/         # Global styles
│   ├── tests/              # Component & E2E tests
│   └── Dockerfile
│
├── cli/                    # Kodia CLI tool
│   ├── internal/
│   │   ├── commands/       # CLI commands
│   │   ├── scaffolding/    # Code generation templates
│   │   └── validation/     # Input validation
│   └── kodia/main.go
│
├── docker-compose.yml      # Local development environment
├── Makefile                # Development commands
├── CONTRIBUTING.md         # Contributing guidelines
├── CHANGELOG.md            # Version history
└── docs/                   # Comprehensive documentation
    ├── GETTING_STARTED.md
    ├── ARCHITECTURE.md
    ├── BACKEND_GUIDE.md
    ├── FRONTEND_GUIDE.md
    ├── CLI_GUIDE.md
    ├── DEPLOYMENT.md
    ├── SECURITY.md
    ├── TESTING.md
    └── FAQ.md
```

---

## 📚 Documentation

Complete documentation is available in the `docs/` directory:

| Document | Purpose |
|----------|---------|
| [Getting Started](docs/GETTING_STARTED.md) | Step-by-step setup and first app |
| [Architecture](docs/ARCHITECTURE.md) | System design and patterns |
| [Backend Guide](docs/BACKEND_GUIDE.md) | API development, routing, middleware |
| [Frontend Guide](docs/FRONTEND_GUIDE.md) | SvelteKit, components, state management |
| [CLI Reference](docs/CLI_GUIDE.md) | All scaffolding commands |
| [Deployment](docs/DEPLOYMENT.md) | Docker, Kubernetes, cloud platforms |
| [Security](docs/SECURITY.md) | Authentication, CORS, rate limiting, best practices |
| [Testing](docs/TESTING.md) | Unit, integration, and E2E testing |
| [FAQ](docs/FAQ.md) | Common questions and troubleshooting |

---

## 🔧 Common Commands

```bash
# Development
make dev              # Start all dev servers
make docker-up        # Start Docker services
make docker-down      # Stop Docker services
make migrate          # Run database migrations

# Backend
make backend-test     # Run backend tests
make backend-build    # Build backend binary
backend/.env setup    # Edit configuration

# Frontend
make frontend-dev     # Start frontend dev server
make frontend-build   # Build for production
make frontend-test    # Run component tests

# CLI
make cli-build        # Build CLI binary
make cli-test         # Test CLI tool
kodia --help          # See all commands

# Database
make db-reset         # Reset database to fresh state
make db-seed          # Run database seeders
make db-shell         # Connect to database shell

# Testing
make test             # Run all tests
make test-coverage    # Generate coverage report
make test-e2e         # Run end-to-end tests
```

---

## 🏗️ Architecture Overview

Kodia follows **Clean Architecture** principles:

```
┌─────────────────────────────────────────────┐
│         HTTP Layer (Gin)                    │
│  (Handlers, Middleware, CORS, Auth)        │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│        Service Layer (Business Logic)       │
│  (Auth, User, Post services, etc.)         │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│    Repository Layer (Data Access)           │
│  (PostgreSQL, MySQL, Redis)                │
└─────────────────────────────────────────────┘
```

**Benefits:**
- ✅ Easy testing (mock repositories)
- ✅ Decoupled from frameworks
- ✅ Clear separation of concerns
- ✅ Flexible database switching

Read [Architecture Guide](docs/ARCHITECTURE.md) for detailed explanation.

---

## 🔐 Security Features

Kodia is **secure by default**:

- ✅ **JWT Authentication** with access/refresh tokens
- ✅ **Rate Limiting** on sensitive endpoints (5 req/15min on auth)
- ✅ **CORS Protection** with configurable origins
- ✅ **Input Validation** using modern validators
- ✅ **SQL Injection Prevention** via parameterized queries
- ✅ **CSRF Protection** via token validation
- ✅ **Password Hashing** with bcrypt
- ✅ **Secure Headers** (CSP, HSTS, X-Frame-Options)
- ✅ **Audit Logging** for sensitive operations
- ✅ **Environment Variable Management** (no hardcoded secrets)

See [Security Guide](docs/SECURITY.md) for implementation details.

---

## 💡 Example: Building a Blog

Generate a complete blog feature in seconds:

```bash
# 1. Create post model with all CRUD operations
kodia generate crud posts \
  --fields=title:string,slug:string,content:text,published:boolean,author:references:users \
  --with-tests \
  --with-validation \
  --with-migrations

# This generates:
# ✅ Database migration (create posts table)
# ✅ Post model (domain entity)
# ✅ PostRepository (database queries)
# ✅ PostService (business logic)
# ✅ PostHandler (REST endpoints)
# ✅ Request/Response DTOs
# ✅ Validation rules
# ✅ Unit & integration tests
# ✅ PostList & PostDetail frontend components

# 2. Create comment feature
kodia generate crud comments \
  --fields=content:text,post:references:posts,author:references:users \
  --with-tests

# 3. Add middleware for auth protection
kodia generate middleware RequireAuth

# 4. Generate API docs automatically
kodia docs generate

# 5. Run tests
make test

# 6. Deploy
docker build -t my-blog .
docker push my-blog:latest
```

---

## 🧪 Testing

Kodia comes with comprehensive testing support:

```bash
# Run all tests
make test

# Run specific test type
make test-unit          # Unit tests only
make test-integration   # Database integration tests
make test-e2e           # End-to-end API tests
make test-coverage      # With coverage report

# Watch mode (rerun on file changes)
make test-watch
```

**Test Structure:**
- Backend: `backend/tests/` with unit, integration, fixtures
- Frontend: `frontend/tests/` with component and E2E tests
- Examples: See [Testing Guide](docs/TESTING.md)

---

## 🚀 Deployment

Deploy to your favorite platform:

**Supported Platforms:**
- ✅ Docker & Kubernetes
- ✅ Heroku
- ✅ Railway
- ✅ Fly.io
- ✅ AWS (ECS, Lambda)
- ✅ GCP (Cloud Run, App Engine)
- ✅ Azure (App Service)
- ✅ DigitalOcean App Platform

See [Deployment Guide](docs/DEPLOYMENT.md) for step-by-step instructions.

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

**How to contribute:**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📊 Comparison with Other Frameworks

| Feature | Kodia | Laravel | Next.js |
|---------|-------|---------|---------|
| Language | Go | PHP | JavaScript |
| Performance | ⚡⚡⚡ Very Fast | ⚡⚡ Good | ⚡⚡ Good |
| Learning Curve | 🟢 Easy | 🟢 Easy | 🟡 Moderate |
| Type Safety | 🟢 Yes (Go) | 🟡 Partial | 🟢 Yes (TS) |
| CLI Scaffolding | ✅ Yes | ✅ Yes | ❌ No |
| Built-in Auth | ✅ Yes | ✅ Yes | ❌ Optional |
| Database ORM | ✅ GORM | ✅ Eloquent | ❌ Optional |
| Rate Limiting | ✅ Yes | ✅ Yes | ❌ Optional |
| Background Jobs | ✅ Yes | ✅ Yes | ❌ Optional |
| Hot Reload | ✅ Yes | ✅ Yes | ✅ Yes |

---

## 📦 Kodia Ecosystem

### Official Plugins

Kodia is infinitely extensible. Install official plugins via CLI:

- 💳 `payment` - Stripe & Midtrans integration (Trial: `kodia plugin install payment`)
- 📧 `notifications` - Email, SMS, Push notifications
- 🗂️ `storage` - AWS S3, Google Cloud Storage
- 🔍 `search` - Elasticsearch, Meilisearch
- 📊 `@kodia/analytics` - Google Analytics, Mixpanel
- 🌐 `@kodia/i18n` - Internationalization support

### Community Packages

We're building an ecosystem. [Publish your package](docs/PUBLISHING_PACKAGES.md)!

---

## 🐛 Reporting Issues

Found a bug? Please [open an issue](https://github.com/kodia-studio/kodia/issues) with:

- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Environment details (Go version, OS, etc.)
- Screenshots if applicable

---

## 📄 License

Kodia Framework is licensed under the [MIT License](LICENSE).

---

## 🌟 Support

- 📖 [Full Documentation](docs/)
- 💬 [GitHub Discussions](https://github.com/kodia-studio/kodia/discussions)
- 🐛 [Issue Tracker](https://github.com/kodia-studio/kodia/issues)
- 💬 [Discord Community](https://discord.gg/kodia) (coming soon)
- 📧 [Email Support](mailto:support@kodia.dev)

---

## 🙏 Acknowledgments

Kodia is inspired by the best practices from:
- **Laravel** - PHP's most elegant framework
- **Rails** - Ruby's convention over configuration
- **Next.js** - Modern fullstack development
- **Spring Boot** - Enterprise Java patterns
- **Django** - Python's batteries-included philosophy

---

**Ready to build something amazing? [Get Started Now!](docs/GETTING_STARTED.md)**

---

Made with ❤️ by the Kodia Team
