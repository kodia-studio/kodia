# 📑 Documentation Index

Complete documentation for Kodia Framework. Start here and navigate to the guide you need.

---

## 🚀 Getting Started

**New to Kodia? Start here!**

- **[Getting Started](GETTING_STARTED.md)** - Step-by-step setup and first app
  - Installation
  - Project structure
  - Running development servers
  - Creating your first feature
  - Troubleshooting

---

## 🏗️ Architecture & Design

Understanding how Kodia is organized:

- **[Architecture Guide](ARCHITECTURE.md)** - System design and patterns
  - Layered architecture
  - Design patterns (Repository, DI, DTO, etc.)
  - Domain-Driven Design
  - Data flow diagrams
  - Component structure

---

## 🔧 Development Guides

Learn how to build features:

### Backend Development

- **[Backend Guide](BACKEND_GUIDE.md)** - REST API development
  - Creating handlers
  - Writing services
  - Working with databases
  - Authentication & authorization
  - Error handling
  - Best practices

- **[Rate Limiting Guide](../backend/docs/RATE_LIMITING.md)** - Prevent abuse
  - Rate limiting implementation
  - Attack prevention
  - Configuration

### Frontend Development

- **[Frontend Guide](FRONTEND_GUIDE.md)** - Enterprise SvelteKit development
  - Svelte 5 Runes & Reactivity
  - Layout System (Admin/App/Auth)
  - Premium Component Library (Forms, Tables, Charts)
  - Reactive Data Catching (Query Rune)
  - State management & API integration
  - Best practices

---

## 🛠️ Tools & CLI

Work with the Kodia CLI tool:

- **[CLI Reference](CLI_GUIDE.md)** - Command reference
  - `kodia generate crud` - Generate CRUD features
  - `kodia generate model` - Create models
  - `kodia generate migration` - Database migrations
  - `kodia generate middleware` - Custom middleware
  - All commands explained with examples

---

## 🧪 Testing

Write reliable, tested code:

- **[Testing Guide](TESTING.md)** - Testing strategies
  - Unit testing
  - Integration testing
  - End-to-end testing
  - Test structure
  - Mocking & fixtures
  - Coverage reporting
  - Testing best practices

---

## 🛡️ Security

Build secure applications:

- **[Security Best Practices](SECURITY.md)** - General security guidelines
- **[Advanced Security Features](ADVANCED_SECURITY.md)** - 2FA, ABAC, Audit Logs, **Social Auth, Passkeys, and API Keys**.
  - Authentication (JWT & Sessions)
  - Social Login (OAuth2)
  - Passkeys (WebAuthn)
  - Authorization (RBAC/ABAC)
  - Input validation & sanitization
  - SQL injection prevention
  - CORS & CSRF protection
  - Rate limiting
  - Scoped API Keys
  - Secrets management
  - Security headers
  - Audit logging
  - Production checklist

---

## 🏗️ Advanced Features

Scale your application with world-class background processing:

- **[Background Workers](WORKERS.md)** - Task Queues & Scheduling
  - **Scheduled Jobs** (Fluent API)
  - **Job Chaining**
  - **Job Batching**
  - **Queue Monitoring Dashboard**
  - Audit Logging

---

## 🚀 Deployment

Deploy to production:

- **[Deployment Guide](DEPLOYMENT.md)** - Production deployment
  - Docker deployment
  - Kubernetes setup
  - Cloud platform guides (AWS, GCP, Azure)
  - Environment configuration
  - Health checks

- **[Observability & Monitoring](OBSERVABILITY.md)** - Performance & Health
  - Distributed Tracing (OpenTelemetry)
  - Metrics Collection (Prometheus)
  - Real-time Health Checks
  - Error Tracking (Sentry)
  - Performance Profiling (pprof)
  - Scaling strategies

---

## 📖 Complete Component Docs

Detailed guides for each component:

### Backend
- [backend/README.md](../backend/README.md) - Backend layer overview
  - Project structure
  - Running the server
  - Creating handlers, services, repositories
  - Common tasks
  - Troubleshooting

### Frontend
- [frontend/README.md](../frontend/README.md) - Frontend layer overview
  - Project structure
  - Development setup
  - Component patterns
  - State management
  - Styling

### CLI
- [cli/README.md](../cli/README.md) - CLI tool
  - Installation
  - All commands
  - Scaffolding templates
  - Contributing to CLI

---

## ❓ FAQ & Troubleshooting

- **[FAQ](FAQ.md)** - Frequently asked questions
  - Common setup issues
  - Development questions
  - Deployment questions
  - Best practices

---

## 🤝 Contributing

Help improve Kodia:

- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute
  - Code of conduct
  - Development setup
  - Commit guidelines
  - Pull request process
  - Coding standards
  - Testing requirements

---

## 📄 Project Information

- **[README.md](../README.md)** - Framework overview
- **[CHANGELOG.md](../CHANGELOG.md)** - Version history and updates
- **[LICENSE](../LICENSE)** - MIT License

---

## 🔗 External Resources

- **[Kodia GitHub](https://github.com/kodia-studio/kodia)** - Source code
- **[Kodia Website](https://kodia.dev)** - Marketing site
- **[GitHub Discussions](https://github.com/kodia-studio/kodia/discussions)** - Community questions
- **[GitHub Issues](https://github.com/kodia-studio/kodia/issues)** - Bug reports & feature requests
- **[Discord](https://discord.gg/kodia)** - Community chat (coming soon)

---

## 🎯 Learning Path

### Beginner: Build Your First App
1. [Getting Started](GETTING_STARTED.md) - Installation & setup
2. [Architecture Guide](ARCHITECTURE.md) - Understand the design
3. [Backend Guide](BACKEND_GUIDE.md) - Create API endpoints
4. [Frontend Guide](FRONTEND_GUIDE.md) - Build UI
5. [CLI Reference](CLI_GUIDE.md) - Use scaffolding

### Intermediate: Production Ready
1. [Testing Guide](TESTING.md) - Write tests
2. [Security Guide](SECURITY.md) - Secure your app
3. [Deployment Guide](DEPLOYMENT.md) - Go live
4. [Performance Guide](DEPLOYMENT.md#performance) - Optimize

### Advanced: Scale & Contribute
1. [Advanced Architecture Patterns](ARCHITECTURE.md#design-patterns)
2. [Contributing Guide](../CONTRIBUTING.md) - Contribute to Kodia
3. [CLI Development](CLI_GUIDE.md#extending) - Extend the CLI
4. [Performance Tuning](DEPLOYMENT.md#scaling)

---

## 💡 Quick Links

**Common Questions:**
- How do I create a new feature? → [Getting Started: Your First Feature](GETTING_STARTED.md#your-first-feature)
- How do I protect an endpoint? → [Backend Guide: Authorization](BACKEND_GUIDE.md#authorization)
- How do I handle errors? → [Backend Guide: Error Handling](BACKEND_GUIDE.md#error-handling)
- How do I test my code? → [Testing Guide](TESTING.md)
- How do I deploy? → [Deployment Guide](DEPLOYMENT.md)
- How do I contribute? → [Contributing Guide](../CONTRIBUTING.md)

---

## 📞 Getting Help

Can't find what you need? Ask the community:

1. **[GitHub Discussions](https://github.com/kodia-studio/kodia/discussions)** - Ask questions
2. **[GitHub Issues](https://github.com/kodia-studio/kodia/issues)** - Report bugs
3. **[Discord](https://discord.gg/kodia)** - Chat with community (coming soon)
4. **[Email](mailto:support@kodia.dev)** - Contact support

For security issues, email: **security@kodia.dev**

---

## 📊 Documentation Statistics

- **15+** comprehensive guides
- **100+** code examples
- **50+** diagrams and illustrations
- **1000+** lines of documentation

---

**Happy coding! 🚀**

Last updated: 2024-04-19
