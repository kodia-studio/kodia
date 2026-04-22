# Kodia Framework Documentation

Welcome to the Kodia Framework documentation. This guide helps you navigate all available documentation.

## Quick Navigation

### For New Developers

**Start here:**
1. [Architecture Overview](ARCHITECTURE.md) — Understand the framework structure
2. [ORM Guide](ORM_GUIDE.md) — Database basics
3. [Validation Layer Quick Reference](VALIDATION_QUICK_REFERENCE.md) — Request validation

### For Feature Developers

**Building Features:**
- [Routing & Controllers](ARCHITECTURE.md#http-layer)
- [Validation Layer](VALIDATION_LAYER.md) — Comprehensive validation guide
- [ORM Features](ORM_GUIDE.md)
- [JWT Security](JWT_SECURITY.md)
- [WebSocket Implementation](WEBSOCKET_GUIDE.md)

### For Security-Focused

**Security Best Practices:**
- [JWT Security](JWT_SECURITY.md)
- [CORS Security](CORS_SECURITY.md)
- [File Upload Security](FILE_UPLOAD_SECURITY.md)
- [Email Template Security](EMAIL_TEMPLATE_SECURITY.md)
- [Logging Security](LOGGING_SECURITY.md)

### For DevOps/Infrastructure

**Infrastructure & Deployment:**
- [Rate Limiting](RATE_LIMITING.md)
- [CORS Configuration](CORS_SECURITY.md)
- [Logging Setup](LOGGING_SECURITY.md)

---

## Documentation by Topic

### Core Framework

| Document | Purpose | Audience |
|----------|---------|----------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | Framework design & patterns | All developers |
| [ORM_GUIDE.md](ORM_GUIDE.md) | Database operations with GORM | Backend developers |
| [VALIDATION_LAYER.md](VALIDATION_LAYER.md) | Request validation system | Backend developers |
| [VALIDATION_QUICK_REFERENCE.md](VALIDATION_QUICK_REFERENCE.md) | Validation cheat sheet | Backend developers |
| [VALIDATION_IMPLEMENTATION_SUMMARY.md](VALIDATION_IMPLEMENTATION_SUMMARY.md) | Implementation details | Framework maintainers |
| [NOTIFICATION_SYSTEM.md](NOTIFICATION_SYSTEM.md) | Notification storage & delivery | Backend developers |

### Features & Integration

| Document | Purpose | Audience |
|----------|---------|----------|
| [NOTIFICATION_SYSTEM.md](NOTIFICATION_SYSTEM.md) | Notifications & real-time delivery | Backend developers |
| [SOCIAL_LOGIN.md](SOCIAL_LOGIN.md) | OAuth2 social authentication | Backend developers |
| [WEBSOCKET_GUIDE.md](WEBSOCKET_GUIDE.md) | Real-time communication | Backend developers |
| [JWT_SECURITY.md](JWT_SECURITY.md) | Token-based authentication | Backend developers |

### Security

| Document | Purpose | Audience |
|----------|---------|----------|
| [JWT_SECURITY.md](JWT_SECURITY.md) | Secure token handling | Backend developers |
| [CORS_SECURITY.md](CORS_SECURITY.md) | Cross-origin requests | DevOps/Backend |
| [FILE_UPLOAD_SECURITY.md](FILE_UPLOAD_SECURITY.md) | Safe file uploads | Backend developers |
| [EMAIL_TEMPLATE_SECURITY.md](EMAIL_TEMPLATE_SECURITY.md) | Secure email rendering | Backend developers |
| [LOGGING_SECURITY.md](LOGGING_SECURITY.md) | Secure logging practices | DevOps/Backend |

### Infrastructure

| Document | Purpose | Audience |
|----------|---------|----------|
| [RATE_LIMITING.md](RATE_LIMITING.md) | Request throttling | DevOps/Backend |

---

## Learning Paths

### 🚀 Getting Started

1. [ARCHITECTURE.md](ARCHITECTURE.md) — Understand the big picture
2. [VALIDATION_QUICK_REFERENCE.md](VALIDATION_QUICK_REFERENCE.md) — Learn validation basics
3. [ORM_GUIDE.md](ORM_GUIDE.md) — Master database operations

**Time:** ~1 hour

### 🔒 Security Hardening

1. [JWT_SECURITY.md](JWT_SECURITY.md) — Secure authentication
2. [CORS_SECURITY.md](CORS_SECURITY.md) — Safe cross-origin access
3. [FILE_UPLOAD_SECURITY.md](FILE_UPLOAD_SECURITY.md) — Secure file handling
4. [EMAIL_TEMPLATE_SECURITY.md](EMAIL_TEMPLATE_SECURITY.md) — Safe email templates
5. [LOGGING_SECURITY.md](LOGGING_SECURITY.md) — Secure logging

**Time:** ~2 hours

### 🛠️ Advanced Features

1. [VALIDATION_LAYER.md](VALIDATION_LAYER.md) — Advanced validation
2. [NOTIFICATION_SYSTEM.md](NOTIFICATION_SYSTEM.md) — Notifications & real-time delivery
3. [WEBSOCKET_GUIDE.md](WEBSOCKET_GUIDE.md) — Real-time features
4. [SOCIAL_LOGIN.md](SOCIAL_LOGIN.md) — OAuth2 integration

**Time:** ~2.5 hours

### 🚢 Production Deployment

1. [RATE_LIMITING.md](RATE_LIMITING.md) — Protect your API
2. [LOGGING_SECURITY.md](LOGGING_SECURITY.md) — Monitor safely
3. [CORS_SECURITY.md](CORS_SECURITY.md) — Configure properly

**Time:** ~1 hour

---

## Recent Documentation

### Newly Added (Latest)

🆕 **[NOTIFICATION_SYSTEM.md](NOTIFICATION_SYSTEM.md)**
- Real-time notification delivery
- Persistent notification storage
- Async email notifications
- WebSocket integration

🆕 **[VALIDATION_LAYER.md](VALIDATION_LAYER.md)**
- Comprehensive request validation system
- 4 custom validation rules
- Eliminates boilerplate code
- Security-focused design

🆕 **[VALIDATION_QUICK_REFERENCE.md](VALIDATION_QUICK_REFERENCE.md)**
- Quick lookup guide
- Common patterns
- Troubleshooting tips

---

## Common Tasks

### "How do I send notifications?"

→ [NOTIFICATION_SYSTEM.md](NOTIFICATION_SYSTEM.md) — Full guide  
→ [WEBSOCKET_GUIDE.md](WEBSOCKET_GUIDE.md) — Real-time delivery

### "How do I validate a request?"

→ [VALIDATION_LAYER.md](VALIDATION_LAYER.md) — Full guide  
→ [VALIDATION_QUICK_REFERENCE.md](VALIDATION_QUICK_REFERENCE.md) — Quick lookup

### "How do I secure API endpoints?"

→ [JWT_SECURITY.md](JWT_SECURITY.md) — Authentication  
→ [CORS_SECURITY.md](CORS_SECURITY.md) — Cross-origin access  

### "How do I handle file uploads safely?"

→ [FILE_UPLOAD_SECURITY.md](FILE_UPLOAD_SECURITY.md)

### "How do I set up real-time features?"

→ [WEBSOCKET_GUIDE.md](WEBSOCKET_GUIDE.md)

### "How do I add OAuth2 login?"

→ [SOCIAL_LOGIN.md](SOCIAL_LOGIN.md)

### "How do I protect my API from abuse?"

→ [RATE_LIMITING.md](RATE_LIMITING.md)

### "How do I work with databases?"

→ [ORM_GUIDE.md](ORM_GUIDE.md)

### "What's the overall architecture?"

→ [ARCHITECTURE.md](ARCHITECTURE.md)

---

## Documentation Standards

All documentation follows these standards:

✅ **Structure**
- Quick summary at the top
- Table of contents
- Code examples
- Best practices section

✅ **Content**
- Real-world examples
- Clear explanations
- Links to related docs
- Troubleshooting section

✅ **Code Examples**
- Runnable code
- Common patterns shown
- ✓ Do's and ✗ Don'ts

---

## Contributing to Documentation

To add or update documentation:

1. Follow the structure of existing docs
2. Include code examples where applicable
3. Add a "See Also" section with related docs
4. Update this README with new topics
5. Keep language clear and concise

---

## Version Information

| Component | Version | Link |
|-----------|---------|------|
| Kodia Framework | Latest | [github.com/kodia-studio/kodia](https://github.com/kodia-studio/kodia) |
| Go | 1.25.0+ | [golang.org](https://golang.org) |
| Gin Framework | 1.12.0+ | [gin-gonic.com](https://gin-gonic.com) |
| GORM | 1.31.1+ | [gorm.io](https://gorm.io) |

---

## Quick Reference Links

### Popular Features
- [Notifications](NOTIFICATION_SYSTEM.md)
- [Validation](VALIDATION_LAYER.md)
- [Authentication](JWT_SECURITY.md)
- [Database](ORM_GUIDE.md)
- [WebSockets](WEBSOCKET_GUIDE.md)
- [File Upload](FILE_UPLOAD_SECURITY.md)
- [OAuth2](SOCIAL_LOGIN.md)

### Security Topics
- [JWT](JWT_SECURITY.md)
- [CORS](CORS_SECURITY.md)
- [File Uploads](FILE_UPLOAD_SECURITY.md)
- [Logging](LOGGING_SECURITY.md)
- [Email](EMAIL_TEMPLATE_SECURITY.md)

### Infrastructure
- [Rate Limiting](RATE_LIMITING.md)
- [Logging](LOGGING_SECURITY.md)
- [CORS](CORS_SECURITY.md)

---

## Need Help?

**Documentation Issues:**
- Report unclear explanations
- Suggest missing topics
- Request code examples
- [Open an issue](https://github.com/kodia-studio/kodia/issues)

**Framework Issues:**
- [GitHub Issues](https://github.com/kodia-studio/kodia/issues)
- [Discussions](https://github.com/kodia-studio/kodia/discussions)

---

**Last Updated:** April 2026  
**Documentation Status:** Complete and Up-to-Date ✅
