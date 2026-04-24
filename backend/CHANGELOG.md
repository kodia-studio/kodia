# Changelog - Kodia Backend 🐨⚙️

All notable changes to the Kodia Backend kernel and official providers.

## [1.6.0] - 2026-04-24

### Added
- **Request Tracing**: Native X-Request-ID middleware for end-to-end request tracking with auto-injection into responses and logs
- **Global Rate Limiting**: LooseRateLimiter (100 req/min) now applied to all `/api/*` routes automatically
- **Graceful Shutdown Improvements**: Configurable timeout (APP_SHUTDOWN_TIMEOUT_SECS), SIGHUP signal support, enhanced logging
- **Unified Database Migrations**: New Migrator system with automatic execution tracking and duplicate prevention
- **Type-Safe DI Container**: Generic `Resolve[T]()` and `MustResolve[T]()` helpers for compile-time safe dependency injection
- **Standardized Error Responses**: All responses include `request_id` and optional `error_code` fields for consistency

### Improved
- **Error Response Format**: Recovery and RateLimit middleware now use response envelope instead of raw gin.H
- **Request Logging**: All structured logs include `request_id` field for request tracing
- **Code Quality**: Replaced all `MustGet()` calls with `MustResolve[T]()` in 8 provider files
- **Documentation**: New comprehensive guides for request tracing, graceful shutdown, database migrations

### Technical Details
- RequestID middleware registered first, ensures all subsequent middleware has tracing context
- Global rate limiter can be layered with endpoint-specific limiters
- Graceful shutdown uses context timeouts for safe connection draining
- Database migrations tracked in `kodia_migrations` table with batch support
- Type-safe container uses Go 1.18+ generics

### Files Modified
- `pkg/response/response.go` — Added RequestID, ErrorCode fields, send() helper
- `pkg/kodia/app.go` — Configurable timeout, SIGHUP support, enhanced logging
- `pkg/config/config.go` — Added ShutdownTimeoutSecs field
- `internal/adapters/http/router.go` — RequestID first, global rate limiter
- `internal/adapters/http/middleware/recovery.go` — Uses response.InternalServerError
- `internal/adapters/http/middleware/ratelimit.go` — Uses response.TooManyRequests
- `internal/adapters/http/middleware/logger.go` — Includes request_id in logs

### Files Created
- `internal/adapters/http/middleware/request_id.go` — X-Request-ID middleware
- `pkg/kodia/container.go` — Type-safe Resolve[T], MustResolve[T]
- `pkg/database/migrator.go` — Unified migration runner
- `internal/infrastructure/database/migrations/go/registry.go` — Migration registry
- `docs/REQUEST_TRACING.md` — Complete request tracing documentation
- `docs/GRACEFUL_SHUTDOWN.md` — Graceful shutdown guide
- `docs/DATABASE_MIGRATIONS.md` — Unified migrations documentation

### Updated Providers
- `auth_provider.go`, `user_provider.go`, `graphql_provider.go`
- `websocket_provider.go`, `notification_provider.go`, `pulse_provider.go`
- `realtime_provider.go`, `observability_provider.go`, `database_provider.go`

### Build Quality
- 0 compilation errors (`go build ./...`)
- 0 vet warnings (`go vet ./...`)
- Backward compatible with v1.5.0

## [1.5.0] - 2026-04-23

### Added
- **Institutional Plugin System**: New formal `Plugin` interface with metadata support for third-party extensions.
- **Global Hook System**: High-performance, thread-safe `HookManager` for cross-plugin communication.
- **Public Benchmarks**: Comprehensive performance data comparison against Laravel, Next.js, and FastAPI.
- **Documentation**: New `PLUGIN_GUIDE.md` and updated `plugins.md` documentation.

### Improved
- **Kernel Resiliency**: Enhanced graceful shutdown logic with priority-based cleanup tasks.
- **Documentation Engine**: Optimized markdown parsing and sidebar navigation for the documentation portal.

## [1.4.0] - 2026-04-15

### Added
- **Query Monitoring**: Automatic N+1 query detection in development mode.
- **Response Caching**: Redis-backed middleware for high-performance API caching.
- **Advanced Security**: New middleware for Rate Limiting and Fingerprinting.

### Fixed
- Fixed race condition in WebSocket connection pool.
- Corrected JWT expiration handling for long-lived refresh tokens.

## [1.3.0] - 2026-03-30

### Added
- **GORM v2 Integration**: Switched to the latest GORM for better performance and type safety.
- **Auto Migrations**: Support for safe, transactional database schema updates.
- **Mail Provider**: Built-in support for SMTP and Mailgun.

---
© 2026 Kodia Studio. "Build like a user, code like a pro."
