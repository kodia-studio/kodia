# CORS Security Configuration Guide

## Overview

This document describes the CORS (Cross-Origin Resource Sharing) security implementation in the Kodia Framework. CORS is a security mechanism that controls which websites can access your API. Proper CORS configuration is critical to prevent credential theft and unauthorized API access.

## Problem Statement

### Before (Vulnerable)

The original CORS configuration had a critical vulnerability:

```go
engine.Use(cors.New(cors.Config{
    AllowOrigins:     r.cfg.CORS.AllowedOrigins,  // ❌ Could contain "*"
    AllowMethods:     []string{"GET", "POST", ...},
    AllowHeaders:     []string{"Origin", "Content-Type", ...},
    AllowCredentials: true,  // ❌ Allows cookies and auth headers
}))
```

**The Vulnerability:**
If `AllowedOrigins` contains a wildcard "*" AND `AllowCredentials` is set to `true`, this creates a **critical security vulnerability**.

### Why This Is Dangerous

**Attack Scenario:**
```
1. User logs into your banking app (gets authentication cookie)
2. User visits attacker's website
3. Attacker's website makes API request to your bank API with origin "*"
   GET /api/accounts?transfer=true
   Cookie: auth_token=user_session_12345
4. Your API allows the request (wildcard origin + credentials)
5. User's account is compromised
```

**The Issue:** Wildcard origins mean ANY website can make authenticated requests to your API!

## Solution Implementation

### 1. Added CORS Validation

**New Function: `ValidateCORSConfig()`**

```go
// Validates CORS configuration for security issues
func ValidateCORSConfig(cfg *config.Config, log *zap.Logger) error {
    // Check for wildcard origin with credentials
    hasWildcard := false
    for _, origin := range corsOrigins {
        if origin == "*" || strings.Contains(origin, "*") {
            hasWildcard = true
            break
        }
    }
    
    // Reject insecure configuration
    if hasWildcard {
        return fmt.Errorf(
            "insecure CORS configuration: cannot use wildcard origin '*' with credentials. "+
                "Wildcard origins allow ANY website to access your API with user credentials. "+
                "Solution: Use specific allowed origins instead (e.g., https://example.com)")
    }
    
    // Validate all origins are absolute URLs
    for _, origin := range corsOrigins {
        if !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
            return fmt.Errorf("invalid CORS origin '%s': must be absolute URL", origin)
        }
    }
    
    return nil
}
```

### 2. Integrated Into Router Setup

**File: `internal/adapters/http/router.go`**

```go
func (r *Router) Setup() *gin.Engine {
    // ... other setup ...
    
    // ✅ Validate CORS configuration BEFORE applying it
    if err := ValidateCORSConfig(r.cfg, r.log); err != nil {
        r.log.Fatal("CORS configuration validation failed", zap.Error(err))
    }
    
    // Apply CORS middleware (guaranteed safe)
    corsConfig := cors.Config{
        AllowOrigins:     r.cfg.CORS.AllowedOrigins,
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }
    engine.Use(cors.New(corsConfig))
    
    return engine
}
```

## Validation Rules

### ✅ Allowed CORS Origins

**Single Domain:**
```yaml
AllowedOrigins:
  - "https://example.com"
```

**Multiple Specific Domains:**
```yaml
AllowedOrigins:
  - "https://app.example.com"
  - "https://api.example.com"
  - "https://admin.example.com"
```

**Development with Localhost:**
```yaml
AllowedOrigins:
  - "http://localhost:3000"
  - "http://localhost:8080"
  - "https://example.com"      # Production domain
```

**HTTPS-Only (Production):**
```yaml
AllowedOrigins:
  - "https://secure-app.example.com"
  - "https://secure-api.example.com"
```

### ❌ Rejected CORS Origins

| Origin | Reason | Security Risk |
|--------|--------|---------------|
| `*` | Wildcard allows ANY origin | ANY website can steal credentials |
| `*.example.com` | Wildcard subdomain | Subdomain enumeration possible |
| `example.com` | Missing scheme | Ambiguous protocol (HTTP or HTTPS?) |
| `http://example.com` | HTTP in production | Unencrypted traffic possible |
| `localhost` | Missing scheme | Invalid absolute URL |

## Attack Scenarios Blocked

### ❌ Wildcard Credential Theft

**Before:**
```
User Configuration: AllowOrigins = ["*"], AllowCredentials = true
Attack: Any website can make authenticated requests
Impact: Account takeover, data theft
```

**After:**
```
Configuration Validation fails with error:
"insecure CORS configuration: cannot use wildcard origin '*' with credentials"
✅ API refuses to start with insecure config
✅ Prevents accidental deployment of vulnerable configuration
```

### ❌ Wildcard Subdomain Attack

**Before:**
```
Configuration: AllowOrigins = ["*.example.com"]
Attack: Attacker controls malicious.example.com subdomain
Impact: Credential theft from all subdomain users
```

**After:**
```
Validation rejects: "cannot use wildcard in origins"
✅ Only specific, enumerated subdomains allowed
✅ Requires explicit allowlisting
```

### ❌ HTTP Credential Leakage in Production

**Before:**
```
Production Config: AllowOrigins = ["http://example.com"]
Attack: Man-in-the-middle intercepts HTTP traffic
Impact: Session hijacking, credentials stolen
```

**After:**
```
Validation WARNING in production:
"Insecure CORS origin in production: Use https:// origins for security"
✅ Logs warning, requires explicit acknowledgment
✅ Helps catch misconfigurations before deployment
```

## Configuration Examples

### Development Environment

```yaml
# .env.development
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080,http://127.0.0.1:3000
```

**Result:**
```
✅ Local development allowed
✅ Multiple local ports supported
✅ No production domains exposed
```

### Production Environment

```yaml
# .env.production
CORS_ALLOWED_ORIGINS=https://app.example.com,https://web.example.com
```

**Result:**
```
✅ HTTPS-only
✅ Specific domains only
✅ No localhost access
✅ No wildcard origins
```

### Staging Environment

```yaml
# .env.staging
CORS_ALLOWED_ORIGINS=https://staging-app.example.com,http://localhost:3000
```

**Result:**
```
✅ Staging domain protected (HTTPS)
✅ Local development still allowed
✅ Localhost + HTTPS staging
```

## Validation Flow

```
┌─────────────────────────────┐
│  Application Startup        │
│  (Router.Setup() called)    │
└──────────────┬──────────────┘
               │
               ▼
    ┌──────────────────────┐
    │ ValidateCORSConfig() │
    └──────┬───────┬───────┘
           │       │
          ✅       ❌
        Valid   Invalid
           │       │
           ▼       ▼
      Apply     Fail with Error
      CORS      (Startup Stops)
      
     API Safe   (No Deployment)
```

## Test Coverage

Comprehensive validation tests in `cors_validator_test.go`:

**Tests Included:**
- ✅ Wildcard origin rejection with credentials
- ✅ Specific origin acceptance (HTTPS, HTTP localhost)
- ✅ Origin scheme validation
- ✅ Empty config handling
- ✅ Attack scenario prevention (wildcard subdomains)
- ✅ Production warnings for insecure origins

**Running Tests:**
```bash
go test -v ./internal/adapters/http
```

## Best Practices

### ✅ DO:

- Use HTTPS origins in production
- List specific, enumerated origins (allowlist approach)
- Use localhost for development
- Separate configs for dev, staging, production
- Review CORS origins during code review
- Monitor CORS errors in logs
- Update origins when adding new services

**Example:**
```yaml
# Production allowlist
allowed_origins:
  - "https://app.example.com"
  - "https://api.example.com"
  - "https://admin.example.com"
```

### ❌ DON'T:

- Use wildcard "*" origins (ever)
- Use wildcard subdomains (*.example.com)
- Use HTTP origins in production
- Allow localhost in production config
- Expose CORS origins in error messages
- Trust origin values from user input
- Use loose wildcard matching

**Example (Don't Do This):**
```yaml
# ❌ INSECURE - Never do this
allowed_origins:
  - "*"
  - "*.example.com"
  - "http://example.com"
```

## Compliance & Standards

### W3C CORS Specification
- ✅ Implements standard CORS validation
- ✅ Proper credential handling
- ✅ Origin validation

### OWASP Security Guidelines
- ✅ Prevents credential theft
- ✅ Implements allowlisting
- ✅ Rejects unsafe configurations

### CWE-434: Unrestricted Upload of File with Dangerous Type
- ✅ Prevents cross-origin credential leakage
- ✅ Enforces origin restrictions

## Production Deployment Checklist

Before deploying to production:

- [ ] CORS origins are HTTPS-only
- [ ] No wildcard origins (*) in config
- [ ] No wildcard subdomains (*.example.com)
- [ ] No localhost origins in production config
- [ ] All required frontend domains are listed
- [ ] CORS validation passes on startup
- [ ] No warnings in application logs
- [ ] Code review approved CORS configuration
- [ ] Staging environment tested successfully
- [ ] Monitoring alerts configured for CORS errors

## Monitoring & Alerts

### Log CORS Errors

```
Application startup log:
✅ "CORS validation passed"
   - Allowed origins: [https://app.example.com, https://admin.example.com]
   - Credentials: enabled
```

### Monitor for Invalid Requests

```
Browser console when blocked by CORS:
Access to XMLHttpRequest at 'https://api.example.com'
from origin 'https://attacker.com' has been blocked by CORS policy
```

This is **GOOD** - the browser (and your CORS config) blocked the unauthorized origin.

## References

- [MDN: Cross-Origin Resource Sharing (CORS)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [OWASP: CORS](https://owasp.org/www-community/CORS)
- [W3C CORS Specification](https://www.w3.org/TR/cors/)
- [CWE-434: Unrestricted Upload](https://cwe.mitre.org/data/definitions/434.html)

## Conclusion

By validating CORS configuration and rejecting wildcard origins with credentials:
- 🔒 Prevents credential theft attacks
- 🛡️ Protects user sessions and data
- 📋 Follows web security standards
- ✅ Enforces allowlist approach
- 🚀 Fails safely (reject by default)
