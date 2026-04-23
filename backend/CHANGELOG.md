# Changelog - Kodia Backend 🐨⚙️

All notable changes to the Kodia Backend kernel and official providers.

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
