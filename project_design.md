# Rancangan Proyek: Framework Kodia

> Framework fullstack profesional berbasis **Golang Gin** + **SvelteKit + Tailwind + Bits UI**, siap pakai, mudah dikembangkan, dan terdokumentasi dengan standar industry.

---

## 🎯 Filosofi & Visi

**Framework Kodia** dirancang dengan filosofi utama:

> **"Opinionated by default, flexible by design."**

Seperti Laravel yang memberikan struktur yang opinionated namun tetap fleksibel, Kodia menyediakan:

| Prinsip | Penjelasan |
|---|---|
| **Convention over Configuration** | Struktur folder dan penamaan yang sudah disepakati sehingga developer tidak perlu konfigurasi manual setiap kali memulai proyek |
| **Batteries Included** | Auth, CRUD, Middleware, Pagination, File Upload, Email sudah tersedia out-of-the-box |
| **Separation of Concerns** | Backend dan frontend terpisah secara jelas, namun terintegrasi dengan mulus via REST API + Type-safe contracts |
| **Developer Experience (DX)** | Kodia CLI memungkinkan scaffold fitur baru dalam hitungan detik, seperti `artisan` di Laravel |
| **Production Ready** | Docker, CI/CD config, health checks, logging structured, dan observability sudah terenkapsulasi |

---

## 🧩 Komponen Utama

Framework Kodia terdiri dari **3 komponen inti**:

```
framework-kodia/
├── backend/          ← Golang Gin (REST API + WebSocket)
├── frontend/         ← SvelteKit + Tailwind v4 + Bits UI
├── kodia-cli/        ← CLI tool (Go) untuk scaffolding
├── docker-compose.yml
└── README.md
```

---

## 1️⃣ Backend — Golang Gin

### Arsitektur: Clean Architecture

Mengikuti prinsip **Clean Architecture** (Uncle Bob) yang diimplementasikan dalam konteks Golang & Gin.

```
backend/
├── cmd/
│   └── server/
│       └── main.go                  ← Entry point: wire semua dependency
│
├── internal/
│   ├── core/                        ← Business Logic Layer (pure Go, no framework)
│   │   ├── domain/                  ← Entities & Value Objects
│   │   │   ├── user.go
│   │   │   ├── errors.go            ← Domain errors custom
│   │   │   └── ...
│   │   ├── ports/                   ← Interface definitions (contracts)
│   │   │   ├── repositories.go      ← Repository interfaces
│   │   │   └── services.go          ← Service interfaces
│   │   └── services/                ← Use Cases / Business rules
│   │       ├── auth_service.go
│   │       ├── user_service.go
│   │       └── ...
│   │
│   ├── adapters/                    ← Adapter Layer
│   │   ├── http/                    ← HTTP Handlers (Gin)
│   │   │   ├── handlers/
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── user_handler.go
│   │   │   │   └── ...
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go          ← JWT Auth middleware
│   │   │   │   ├── cors.go
│   │   │   │   ├── logger.go        ← Structured logging
│   │   │   │   ├── rate_limiter.go
│   │   │   │   └── recovery.go
│   │   │   ├── dto/                 ← Request/Response DTOs
│   │   │   │   ├── auth_dto.go
│   │   │   │   ├── user_dto.go
│   │   │   │   └── response.go      ← Standard API response wrapper
│   │   │   ├── validators/          ← Input validation
│   │   │   └── router.go            ← Route definitions
│   │   │
│   │   └── repository/              ← Database implementations
│   │       ├── postgres/
│   │       │   ├── user_repo.go
│   │       │   └── ...
│   │       └── cache/
│   │           └── redis_cache.go
│   │
│   └── infrastructure/              ← Infrastructure Layer
│       ├── database/
│       │   ├── postgres.go          ← PostgreSQL connection
│       │   ├── migrations/          ← SQL migration files
│       │   │   ├── 000001_create_users_table.up.sql
│       │   │   └── 000001_create_users_table.down.sql
│       │   └── seeders/
│       ├── cache/
│       │   └── redis.go
│       ├── mailer/
│       │   └── smtp.go
│       ├── storage/
│       │   └── local.go             ← File storage (local/S3)
│       └── logger/
│           └── zap.go               ← Structured logging (Zap)
│
├── pkg/                             ← Reusable public packages
│   ├── jwt/                         ← JWT utilities
│   ├── hash/                        ← Password hashing (bcrypt/argon2)
│   ├── pagination/                  ← Standard pagination helper
│   ├── response/                    ← JSON response builder
│   ├── validator/                   ← Custom validation rules
│   └── config/                      ← Config loader (viper)
│
├── configs/
│   ├── app.go                       ← App configuration struct
│   └── config.yaml.example          ← Config template
│
├── tests/
│   ├── unit/
│   ├── integration/
│   └── fixtures/
│
├── Makefile                         ← Dev commands (make run, make test, make migrate)
├── Dockerfile
├── .env.example
├── go.mod
└── go.sum
```

### Teknologi Stack Backend

| Komponen | Library | Alasan |
|---|---|---|
| HTTP Router | `gin-gonic/gin` | Performa tinggi, ekosistem luas |
| Database ORM | `jackc/pgx` + `sqlc` | Type-safe SQL, tidak magic |
| Migrations | `golang-migrate` | Standard migration tool |
| Config | `spf13/viper` | Multi-source config (env, yaml, flags) |
| Validation | `go-playground/validator` | Battle-tested validation |
| JWT | `golang-jwt/jwt` | JWT v5, secure by default |
| Password | `golang.org/x/crypto` (bcrypt) | NIST-recommended |
| Logging | `uber-go/zap` | Structured logging, high performance |
| Testing | `testify` + `testcontainers` | Integration testing dengan DB nyata |
| Docs | `swaggo/swag` | Auto-generate Swagger/OpenAPI |
| Cache | Redis via `redis/go-redis` | Session, rate limiting, cache layer |

### Standard API Response Format

Semua API response menggunakan format standar yang konsisten:

```json
{
  "success": true,
  "message": "Data berhasil diambil",
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 15,
    "total": 100,
    "total_pages": 7
  }
}
```

Error response:
```json
{
  "success": false,
  "message": "Validasi gagal",
  "errors": {
    "email": ["Email tidak valid"],
    "password": ["Password minimal 8 karakter"]
  }
}
```

---

## 2️⃣ Frontend — SvelteKit + Tailwind + Bits UI

### Arsitektur: Feature-First

```
frontend/
├── src/
│   ├── lib/
│   │   ├── components/
│   │   │   ├── ui/                  ← Bits UI primitives (copy-owned)
│   │   │   │   ├── button/
│   │   │   │   ├── dialog/
│   │   │   │   ├── input/
│   │   │   │   ├── table/
│   │   │   │   ├── toast/
│   │   │   │   └── ...
│   │   │   └── shared/              ← Komponen reusable custom
│   │   │       ├── AppHeader.svelte
│   │   │       ├── AppSidebar.svelte
│   │   │       ├── DataTable.svelte  ← Generic table dengan sorting/filter
│   │   │       ├── PageHeader.svelte
│   │   │       ├── ConfirmDialog.svelte
│   │   │       └── ...
│   │   │
│   │   ├── api/                     ← API client layer
│   │   │   ├── client.ts            ← Base fetch wrapper (auth, error handling)
│   │   │   ├── auth.ts
│   │   │   ├── users.ts
│   │   │   └── ...
│   │   │
│   │   ├── stores/                  ← Svelte stores (global state)
│   │   │   ├── auth.store.ts        ← Auth state (user, token)
│   │   │   ├── toast.store.ts       ← Toast notification state
│   │   │   └── theme.store.ts       ← Dark/light mode
│   │   │
│   │   ├── types/                   ← TypeScript type definitions
│   │   │   ├── api.types.ts         ← API response types
│   │   │   ├── auth.types.ts
│   │   │   └── ...
│   │   │
│   │   ├── utils/                   ← Helper functions
│   │   │   ├── format.ts            ← Date, currency, number formatting
│   │   │   ├── validation.ts
│   │   │   └── ...
│   │   │
│   │   └── server/                  ← Server-only code (SvelteKit)
│   │       ├── api.ts               ← Server-side API calls
│   │       └── auth.ts              ← Session management
│   │
│   ├── routes/
│   │   ├── (auth)/                  ← Auth route group (no layout)
│   │   │   ├── login/
│   │   │   │   ├── +page.svelte
│   │   │   │   └── +page.server.ts
│   │   │   ├── register/
│   │   │   └── forgot-password/
│   │   │
│   │   ├── (app)/                   ← App route group (with main layout)
│   │   │   ├── +layout.svelte       ← Main app layout (sidebar + header)
│   │   │   ├── +layout.server.ts    ← Auth guard
│   │   │   ├── dashboard/
│   │   │   │   ├── +page.svelte
│   │   │   │   └── +page.server.ts
│   │   │   └── [feature]/           ← Feature-based routing
│   │   │
│   │   ├── api/                     ← SvelteKit API routes (BFF pattern)
│   │   │   └── [...path]/
│   │   │       └── +server.ts       ← Proxy ke backend (opsional)
│   │   │
│   │   ├── +layout.svelte           ← Root layout
│   │   └── +error.svelte            ← Error page
│   │
│   ├── app.html
│   ├── app.css                      ← Global CSS + Tailwind base
│   └── hooks.server.ts              ← Auth session hooks
│
├── static/
│   ├── fonts/
│   └── images/
│
├── svelte.config.js
├── tailwind.config.ts               ← Custom design tokens
├── vite.config.ts
├── tsconfig.json
├── package.json
└── .env.example
```

### Design System

Kodia Frontend menggunakan design system yang terstruktur berbasis Tailwind v4:

| Token | Contoh |
|---|---|
| **Brand Color** | `--color-primary-*` (customizable per project) |
| **Typography** | Inter / Geist (via Google Fonts atau self-hosted) |
| **Spacing** | Tailwind default scale |
| **Border Radius** | `--radius-*` consistent tokens |
| **Dark Mode** | Class-based (`dark:`) dengan Svelte store |

---

## 3️⃣ Kodia CLI

CLI tool yang ditulis dalam Go, berfungsi sebagai `artisan`-nya Kodia.

### Struktur CLI

```
kodia-cli/
├── cmd/
│   └── kodia/
│       ├── main.go                  ← Entry point
│       └── root.go                  ← Root command (cobra)
│
├── internal/
│   ├── commands/                    ← Semua command implementations
│   │   ├── new.go                   ← kodia new <project-name>
│   │   ├── make/
│   │   │   ├── handler.go           ← kodia make:handler <Name>
│   │   │   ├── service.go           ← kodia make:service <Name>
│   │   │   ├── repository.go        ← kodia make:repository <Name>
│   │   │   ├── migration.go         ← kodia make:migration <name>
│   │   │   ├── page.go              ← kodia make:page <route>
│   │   │   ├── component.go         ← kodia make:component <Name>
│   │   │   └── feature.go           ← kodia make:feature <Name> (full scaffold)
│   │   ├── db/
│   │   │   ├── migrate.go           ← kodia db:migrate
│   │   │   ├── rollback.go          ← kodia db:rollback
│   │   │   └── seed.go              ← kodia db:seed
│   │   └── dev.go                   ← kodia dev (start both servers)
│   │
│   ├── scaffolding/                 ← Template engine untuk code generation
│   │   ├── templates/               ← Go template files (.tmpl)
│   │   │   ├── handler.tmpl
│   │   │   ├── service.tmpl
│   │   │   ├── repository.tmpl
│   │   │   ├── migration.tmpl
│   │   │   ├── svelte-page.tmpl
│   │   │   └── svelte-component.tmpl
│   │   └── generator.go             ← Template rendering engine
│   │
│   └── config/
│       └── detect.go                ← Auto-detect project root & config
│
├── Makefile
├── go.mod
└── go.sum
```

### Daftar Commands CLI

```bash
# Membuat proyek baru
kodia new my-app                     # Full-stack project baru
kodia new my-app --backend-only      # Hanya backend
kodia new my-app --frontend-only     # Hanya frontend

# Scaffold Backend
kodia make:handler User              # Buat handler CRUD lengkap
kodia make:service User              # Buat service + interface
kodia make:repository User           # Buat repository + interface
kodia make:migration create_users    # Buat file migration SQL
kodia make:feature Product           # Full scaffold: handler + service + repo + migration + route + page

# Database
kodia db:migrate                     # Jalankan semua migration
kodia db:rollback                    # Rollback 1 step
kodia db:rollback --steps=3          # Rollback N steps
kodia db:seed                        # Jalankan seeders
kodia db:fresh                       # Drop all + re-migrate + seed

# Scaffold Frontend
kodia make:page products/index       # Buat route SvelteKit
kodia make:component ProductCard     # Buat Svelte component
kodia make:api-client products       # Buat API client TypeScript

# Development
kodia dev                            # Start backend + frontend bersamaan
kodia dev --backend                  # Hanya backend
kodia dev --frontend                 # Hanya frontend

# Utility
kodia version                        # Print versi framework
kodia update                         # Update ke versi terbaru
```

---

## 🔄 Alur Komunikasi Frontend–Backend

```
                    ┌─────────────────────────────────────┐
                    │           SvelteKit Frontend         │
                    │                                     │
                    │  +page.server.ts  →  API Client     │
                    │  (load function)     ($lib/api/*)   │
                    └──────────────────┬──────────────────┘
                                       │ HTTP/HTTPS
                                       │ (Bearer JWT)
                    ┌──────────────────▼──────────────────┐
                    │            Gin Backend               │
                    │                                     │
                    │  Router → Middleware → Handler       │
                    │               ↓                     │
                    │           Service Layer              │
                    │               ↓                     │
                    │        Repository Layer              │
                    │               ↓                     │
                    │     PostgreSQL / Redis               │
                    └─────────────────────────────────────┘
```

### Fitur Out-of-the-Box

| Fitur | Backend | Frontend |
|---|---|---|
| **Authentication** | JWT (access + refresh token) | Login/Register page + auth store |
| **Authorization** | RBAC middleware | Role-based UI gating |
| **CRUD Pattern** | Service + Repository + Handler | DataTable + Form + API client |
| **File Upload** | Multipart upload endpoint | File input component |
| **Pagination** | Standard pagination di response | Paginator component |
| **Dark Mode** | — | Svelte store + Tailwind class strategy |
| **Toast Notification** | — | Bits UI + custom store |
| **Form Validation** | go-playground/validator | Client-side + zod |
| **API Documentation** | Swagger/OpenAPI (swaggo) | — |
| **WebSocket** | Gin-based WebSocket | Svelte reactive store |
| **Rate Limiting** | Redis-based middleware | — |
| **Health Check** | `/api/health` endpoint | — |
| **Audit Log** | Middleware-based audit trail | — |

---

## 📦 Strategi Instalasi

### Cara 1: Via Kodia CLI (Direkomendasikan)
```bash
# Install CLI
curl -fsSL https://kodia.dev/install.sh | bash
# atau
go install github.com/kodia/cli@latest

# Buat proyek baru
kodia new my-awesome-app
cd my-awesome-app

# Mulai development
kodia dev
```

### Cara 2: Manual / Git Template
```bash
git clone https://github.com/kodia/framework my-app
cd my-app
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
docker-compose up -d   # PostgreSQL + Redis
kodia db:migrate
kodia dev
```

---

## 📚 Sistem Dokumentasi

Terinspirasi dari dokumentasi Laravel, Tailwind CSS, dan SvelteKit:

```
docs/
├── getting-started/
│   ├── installation.md
│   ├── directory-structure.md
│   ├── configuration.md
│   └── quick-start.md
│
├── backend/
│   ├── routing.md
│   ├── middleware.md
│   ├── handlers.md
│   ├── services.md
│   ├── repositories.md
│   ├── database.md
│   ├── migrations.md
│   ├── authentication.md
│   ├── authorization.md
│   ├── validation.md
│   ├── file-upload.md
│   ├── caching.md
│   ├── logging.md
│   ├── testing.md
│   └── deployment.md
│
├── frontend/
│   ├── structure.md
│   ├── routing.md
│   ├── api-client.md
│   ├── state-management.md
│   ├── components.md
│   ├── forms.md
│   ├── authentication.md
│   ├── dark-mode.md
│   └── testing.md
│
├── cli/
│   ├── installation.md
│   ├── commands.md
│   └── scaffolding.md
│
└── deployment/
    ├── docker.md
    ├── server.md
    └── ci-cd.md
```

Dokumentasi akan dibangun sebagai **website statis** (menggunakan VitePress atau Astro Starlight) — seperti:
- laravel.com/docs
- tailwindcss.com/docs
- kit.svelte.dev

---

## 🗺️ Roadmap Versi

### v0.1.0 — Foundation (MVP)
- [x] Struktur folder backend + frontend ditetapkan
- [ ] Backend boilerplate: Auth (JWT), User CRUD, Middleware stack
- [ ] Frontend boilerplate: Auth pages, Dashboard layout, API client
- [ ] Kodia CLI: `kodia new`, `kodia dev`, `kodia make:feature`
- [ ] Docker setup (PostgreSQL + Redis + Backend + Frontend)

### v0.2.0 — DX Enhancement
- [ ] `kodia make:*` commands lengkap (handler, service, repo, page, component)
- [ ] Database migration tooling
- [ ] Auto-generated Swagger docs
- [ ] WebSocket support
- [ ] File upload (local + S3)

### v0.3.0 — Production Ready
- [ ] Multi-tenancy support
- [ ] Email templating
- [ ] Job queue (Redis-based)
- [ ] Scheduled tasks (cron jobs)
- [ ] CI/CD templates (GitHub Actions)

### v1.0.0 — Stable Release
- [ ] Dokumentasi lengkap (versi website)
- [ ] Test coverage ≥ 80%
- [ ] Performance benchmarks
- [ ] Community plugins/packages support

---

## ⚙️ Konvensi & Standar Kode

### Backend (Go)

| Hal | Konvensi |
|---|---|
| **Penamaan Package** | lowercase, singular (`user`, `product`) |
| **Penamaan File** | snake_case (`user_handler.go`, `auth_service.go`) |
| **Penamaan Struct** | PascalCase (`UserHandler`, `AuthService`) |
| **Penamaan Interface** | PascalCase + akhiran `Repository`/`Service` (`UserRepository`) |
| **Error Handling** | Selalu explicit, gunakan `domain.ErrNotFound` bukan string |
| **Logging** | Selalu structured (key-value), tidak boleh `fmt.Println` |
| **Test Files** | Sesuai package (`user_service_test.go`) |

### Frontend (SvelteKit/TypeScript)

| Hal | Konvensi |
|---|---|
| **Komponen** | PascalCase (`ProductCard.svelte`) |
| **Route Files** | SvelteKit convention (`+page.svelte`, `+layout.svelte`) |
| **API Functions** | camelCase, verb-first (`getUsers`, `createProduct`) |
| **Types** | PascalCase + suffix `Type`/`Dto` (`UserType`, `CreateProductDto`) |
| **Stores** | camelCase + suffix `.store.ts` (`auth.store.ts`) |

---

## 🐳 Docker Setup

```yaml
# docker-compose.yml (rancangan lengkap)
services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: kodia_db
      POSTGRES_USER: kodia
      POSTGRES_PASSWORD: secret
    ports: ["5432:5432"]
    volumes: [postgres_data:/var/lib/postgresql/data]

  redis:
    image: redis:7-alpine
    ports: ["6379:6379"]

  backend:
    build: ./backend
    depends_on: [postgres, redis]
    env_file: ./backend/.env
    ports: ["8080:8080"]

  frontend:
    build: ./frontend
    depends_on: [backend]
    ports: ["3000:3000"]
```

---

## ❓ Pertanyaan Terbuka untuk Diskusi

> [!IMPORTANT]
> Poin-poin di bawah perlu dikonfirmasi sebelum implementasi dimulai.

1. **Nama Domain Framework** — Apakah nama resmi framework ini `Kodia Framework` atau ada nama lain? Ini akan mempengaruhi nama package Go, nama npm scope, dan URL dokumentasi.

2. **Database Target** — Apakah fokus hanya PostgreSQL, atau perlu support MySQL juga? (Laravel support keduanya)

3. **Multi-tenancy** — Apakah fitur multi-tenancy (satu instance, banyak tenant) perlu ada dari v0.1, atau bisa ditunda ke v1.0?

4. **Authentication Strategy** — Apakah session-based + JWT, atau murni JWT saja? Dan apakah perlu OAuth (Login with Google/GitHub) dari awal?

5. **Prioritas v0.1** — Dari roadmap di atas, fitur mana yang paling kritis untuk dibangun pertama? Apakah mulai dari CLI, backend boilerplate, atau frontend boilerplate?
