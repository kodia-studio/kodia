# JWT Security Configuration

## Overview

This document explains how to properly configure JWT secrets for the Kodia Framework.

## Requirements

### Development Environment
- JWT secrets can be short or empty (for testing purposes)
- The framework will log a warning if secrets are weak
- This is acceptable for local development only

### Production Environment
- **MANDATORY**: Both `JWT_ACCESS_SECRET` and `JWT_REFRESH_SECRET` must be set
- **MANDATORY**: Both secrets must be at least 32 characters long
- Failure to meet these requirements will cause the application to exit with an error

## Generating Secure Secrets

### Using the provided script:
```bash
cd backend
bash scripts/generate-jwt-secrets.sh
```

### Manual generation using OpenSSL:
```bash
# Generate a 32-character random secret
openssl rand -base64 32

# Generate two secrets (one for access, one for refresh)
openssl rand -base64 32
openssl rand -base64 32
```

### Using alternative methods:
```bash
# Using Python
python3 -c "import secrets; print(secrets.token_urlsafe(32))"

# Using Node.js
node -e "console.log(require('crypto').randomBytes(32).toString('base64'))"
```

## Configuration

### 1. Development Setup
For local development, create a `.env` file in the backend directory:

```bash
# You can use simple secrets for development
APP_JWT_ACCESS_SECRET=dev-access-secret-key
APP_JWT_REFRESH_SECRET=dev-refresh-secret-key
```

### 2. Production Deployment
Before deploying to production:

1. Generate strong secrets using one of the methods above
2. Set environment variables (recommended):
   ```bash
   export APP_JWT_ACCESS_SECRET="your-generated-32-char-secret-1"
   export APP_JWT_REFRESH_SECRET="your-generated-32-char-secret-2"
   ```

3. Or add to your production `.env` file:
   ```bash
   APP_JWT_ACCESS_SECRET=<your-32-char-secret-1>
   APP_JWT_REFRESH_SECRET=<your-32-char-secret-2>
   ```

## Best Practices

✅ **DO:**
- Generate unique secrets for access and refresh tokens
- Use secrets at least 32 characters long in production
- Store secrets in environment variables (never hardcode)
- Rotate secrets periodically
- Use different secrets for different environments (dev, staging, prod)
- Add `.env` to `.gitignore` to prevent accidental commits

❌ **DON'T:**
- Use weak or predictable secrets
- Commit `.env` files to version control
- Use the same secret for access and refresh tokens
- Share secrets across environments
- Use default/placeholder secrets in production

## Secret Rotation

If you suspect a secret has been compromised:

1. Generate a new secret
2. Deploy with both old and new secrets (allows grace period)
3. Update configuration to use only the new secret
4. Monitor for authentication issues

## Troubleshooting

### Error: "JWT_ACCESS_SECRET must be set in production"
**Solution**: Set the `APP_JWT_ACCESS_SECRET` environment variable with a 32+ character secret

### Error: "JWT_ACCESS_SECRET must be at least 32 characters long"
**Solution**: Generate a new secret using the `generate-jwt-secrets.sh` script

### Development warning: "JWT secrets are weak"
**This is normal** in development mode. The application will still start and function correctly. This warning helps ensure you don't accidentally deploy development secrets to production.

## References

- [RFC 7518 - JSON Web Algorithms (JWA)](https://tools.ietf.org/html/rfc7518)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [OWASP JWT Security](https://cheatsheetseries.owasp.org/cheatsheets/JSON_Web_Token_for_Java_Cheat_Sheet.html)
