# Logging Security Guide

## Overview

This document describes security measures implemented in the Kodia Framework to prevent information disclosure through application logs. The focus is on preventing accidental exposure of sensitive data (credentials, infrastructure details, etc.) in debug and info logs.

## Problem Statement

### Before (Vulnerable)

The database connection logging included sensitive infrastructure information:

```go
log.Debug("Attempting database connection",
    zap.String("driver", cfg.Database.Driver),
    zap.String("host", cfg.Database.Host),           // ❌ Infrastructure detail
    zap.Int("port", cfg.Database.Port),              // ❌ Infrastructure detail
    zap.String("user", cfg.Database.User),           // ❌ Credential
    zap.String("database", cfg.Database.Name),
    zap.String("ssl_mode", cfg.Database.SSLMode),    // ❌ Infrastructure detail
)
```

**Risks:**
- ❌ Logs exposing database credentials if intercepted
- ❌ Infrastructure details (host, port) exposing network topology
- ❌ Database usernames in logs searchable by attackers
- ❌ Log aggregation systems exposing sensitive data to unauthorized users
- ❌ Compliance violations (PCI-DSS, HIPAA if logs contain credentials)

### After (Secure)

```go
log.Debug("Attempting database connection",
    zap.String("driver", cfg.Database.Driver),       // ✅ Non-sensitive
    zap.String("database", cfg.Database.Name),       // ✅ Non-sensitive
)
```

## Security Implementation

### Sensitive vs Non-Sensitive Data

**Sensitive Data (DO NOT LOG):**
- ❌ Database host/IP address
- ❌ Database port number
- ❌ Database username
- ❌ Database password
- ❌ SSL/TLS mode or certificate paths
- ❌ API keys or tokens
- ❌ Session tokens or cookies
- ❌ Encryption keys
- ❌ Email addresses (PII)
- ❌ Phone numbers (PII)
- ❌ Physical addresses (PII)
- ❌ Social security numbers
- ❌ Credit card numbers

**Non-Sensitive Data (Safe to Log):**
- ✅ Application version
- ✅ Environment name (development, staging, production)
- ✅ Database driver type (postgres, mysql)
- ✅ Database name (logical identifier)
- ✅ HTTP status codes
- ✅ Response times/latency
- ✅ Feature flags (enabled/disabled)
- ✅ Error types (not error details if they contain sensitive data)

## Changes Made

### Database Connection Logging

**File: `internal/infrastructure/database/database.go`**

**Debug Log (Line 36-39):**
```go
log.Debug("Attempting database connection",
    zap.String("driver", cfg.Database.Driver),
    zap.String("database", cfg.Database.Name),
)
```

**Info Log (Line 75-78):**
```go
log.Info("Database connected",
    zap.String("driver", cfg.Database.Driver),
    zap.String("database", cfg.Database.Name),
)
```

**Removed Fields:**
- ❌ `host` - Infrastructure detail
- ❌ `port` - Infrastructure detail
- ❌ `user` - Database credential
- ❌ `ssl_mode` - Infrastructure detail

## Attack Scenarios Blocked

### Log Interception

❌ **Blocked**: Attacker intercepts logs and finds credentials
```
BEFORE: [DEBUG] Database connection attempt host=secure-db.company.com port=5432 user=admin_prod password=X
        ↓
        Attacker uses credentials for unauthorized access

AFTER:  [DEBUG] Database connection attempt driver=postgres database=company_db
        ↓
        No useful attack information exposed
```

### Log Aggregation Exposure

❌ **Blocked**: Unauthorized user views logs in centralized logging system
```
BEFORE: Logs aggregated in ELK/Datadog with sensitive fields
        Junior developer or contractor can search for: "user=" or "password="
        Finds production credentials in logs

AFTER:  No credentials stored in logs
        Log search results only show non-sensitive information
```

### Supply Chain Compromise

❌ **Blocked**: Attacker gains temporary access to log files
```
BEFORE: Attacker finds: host=10.0.1.50, user=prod_user
        Can scan network at 10.0.1.x or target that user account

AFTER:  No network topology or credential information in logs
        Infrastructure details protected
```

## Logging Best Practices

### General Rules

✅ **DO:**
- Log errors with context (error type, operation attempted)
- Log important milestones (startup, shutdown, migrations)
- Log non-sensitive config values (environment, version)
- Log request IDs for traceability
- Use structured logging (not string concatenation)
- Include timestamps and log levels
- Log access attempts (successful and failed) with IP/user agent
- Sanitize user input if logged

❌ **DON'T:**
- Log passwords, keys, tokens, or secrets
- Log full request/response bodies containing PII
- Log database credentials or connection strings
- Log infrastructure details (IPs, ports, host names)
- Log personally identifiable information (names, emails, IDs)
- Log third-party API keys
- Use log level INFO for sensitive data (use DEBUG with caution)
- Log raw SQL queries that might contain sensitive values

### Structured Logging Pattern

**Good Example:**
```go
log.Info("User login successful",
    zap.String("user_id", userId),        // ✅ Non-PII identifier
    zap.String("method", "password"),     // ✅ Non-sensitive info
    zap.String("ip_address", ipAddr),     // ✅ Non-sensitive info
    // ❌ DON'T: zap.String("password", password)
    // ❌ DON'T: zap.String("email", email)
)
```

**Bad Example:**
```go
log.Info("Login: " + email + ":" + password)    // ❌ Contains credentials
log.Debug("Full request: " + requestBody)       // ❌ May contain sensitive data
```

## Test Coverage

Comprehensive tests verify sensitive data is not logged:

```bash
go test -v ./internal/infrastructure/database
```

Tests verify:
- ✅ Host/IP addresses not logged
- ✅ Port numbers not logged
- ✅ Database usernames not logged
- ✅ SSL mode not logged
- ✅ Driver type IS logged (non-sensitive)
- ✅ Database name IS logged (non-sensitive)

## Sensitive Data Patterns to Detect

### Email Address
```go
// ❌ DON'T
log.Info("User action", zap.String("email", user.Email))

// ✅ DO
log.Info("User action", zap.String("user_id", user.ID))
```

### Phone Number
```go
// ❌ DON'T
log.Info("Contact update", zap.String("phone", phone))

// ✅ DO
log.Info("Contact update", zap.String("contact_type", "phone"))
```

### Connection Strings
```go
// ❌ DON'T
log.Debug("Database URL", zap.String("dsn", cfg.Database.DSN()))

// ✅ DO
log.Debug("Database connecting", 
    zap.String("driver", cfg.Database.Driver),
    zap.String("database", cfg.Database.Name),
)
```

### API Keys
```go
// ❌ DON'T
log.Info("API initialized", zap.String("key", apiKey))

// ✅ DO
log.Info("API initialized", zap.String("provider", "stripe"))
```

## Third-Party Service Integration

### Log Aggregation Services (ELK, Datadog, etc.)

When sending logs to third-party services:

✅ **DO:**
- Configure redaction rules for sensitive patterns
- Use log aggregation service's PII detection features
- Implement log filtering before transmission
- Audit who has access to logs
- Use VPC/private endpoints for log transmission

❌ **DON'T:**
- Send raw logs without sanitization
- Trust that sensitive data won't be exposed
- Store logs in third-party systems without encryption
- Grant log access to all team members

### Example Datadog Scrubbing Rules
```yaml
logs:
  processing:
    rules:
      - type: mask_sequences
        name: "mask_passwords"
        pattern: "password[=:]\S+"
        replace_placeholder: "password=[REDACTED]"
      - type: mask_sequences
        name: "mask_tokens"
        pattern: "token[=:]\S+"
        replace_placeholder: "token=[REDACTED]"
```

## Compliance Implications

### PCI-DSS (Payment Card Industry)
- ❌ Requirement 3.4: Render PAN unreadable in logs
- ❌ Requirement 6.5.10: No hardcoded credentials
- ✅ Our fix: No credentials in logs

### HIPAA (Healthcare)
- ❌ Requirement: Protect PHI (Protected Health Information)
- ❌ Requirement: Audit controls for data access
- ✅ Our fix: No credentials or infrastructure details exposing patient data access

### GDPR (General Data Protection Regulation)
- ❌ Requirement: Protect personal data
- ❌ Requirement: Data minimization principle
- ✅ Our fix: Only non-PII data logged

## Incident Response

If sensitive data was logged:

1. **Immediate Actions:**
   - Identify affected logs and time period
   - Check if sensitive data is searchable in log system
   - Rotate affected credentials immediately
   - Check access logs for suspicious queries

2. **Investigation:**
   - Determine what was exposed
   - Verify credentials haven't been used maliciously
   - Check for unauthorized access attempts
   - Review who had log access during exposure period

3. **Remediation:**
   - Remove sensitive data from logs
   - Deploy code fix (as in this commit)
   - Rebuild/redeploy application
   - Update log retention policies
   - Implement log redaction rules

4. **Prevention:**
   - Code review checklist: "No credentials in logs"
   - Automated scanning for common patterns
   - Security training for developers
   - Regular log audits

## Performance Considerations

Removing sensitive fields from logs has **positive** performance impact:
- ✅ Fewer data fields to serialize
- ✅ Faster log transmission
- ✅ Reduced log storage requirements
- ✅ Faster log queries (fewer fields to index)

## Monitoring & Alerts

### Log Monitoring Rules

Create alerts for:
- ✅ Failed database connection attempts
- ✅ Multiple failed login attempts
- ✅ Connection pool exhaustion
- ✅ Migration failures

Avoid alerting on:
- ❌ Credential exposure patterns (since we don't log them)
- ❌ Infrastructure details (since we don't log them)

## Code Review Checklist

When reviewing code with logging:

- [ ] No credentials in logs (passwords, keys, tokens)
- [ ] No PII in logs (emails, phone numbers, addresses)
- [ ] No infrastructure details (IPs, ports, hostnames)
- [ ] No API keys or secrets
- [ ] No full request/response bodies with sensitive data
- [ ] Structured logging used (not string concatenation)
- [ ] Log levels appropriate (DEBUG for details, INFO for events)
- [ ] Error messages don't expose internals
- [ ] Third-party service integrations sanitize logs

## References

- [OWASP: Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)
- [CWE-532: Insertion of Sensitive Information into Log File](https://cwe.mitre.org/data/definitions/532.html)
- [NIST: Application and Database Security](https://csrc.nist.gov/publications/detail/sp/800-53/rev-5)
- [PCI-DSS: Logging Requirements](https://www.pcisecuritystandards.org/)

## Conclusion

By removing sensitive information from logs while maintaining useful debugging information, we achieve:
- 🔒 Better security posture
- 📋 Compliance with regulations
- 🔍 Still useful debugging information
- ⚡ Improved performance
- 🛡️ Reduced attack surface
