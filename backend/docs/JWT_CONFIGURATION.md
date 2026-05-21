# JWT Token Configuration Guide

## Overview

Kodia Framework uses JWT (JSON Web Tokens) for authentication with configurable token expiry times. This guide explains how to configure token duration and best practices.

## Current Configuration

The default configuration provides a balance between security and user experience:

```env
APP_JWT_ACCESS_EXPIRY_HOURS=168      # 7 days
APP_JWT_REFRESH_EXPIRY_DAYS=90       # 90 days
```

| Token Type | Duration | Use Case |
|------------|----------|----------|
| Access Token | 7 days | Used for API requests, long-lived for better UX |
| Refresh Token | 90 days | Used to obtain new access tokens |

## How Token Expiry Works

### Access Token (Bearer Token)
- Used in every API request: `Authorization: Bearer <access_token>`
- Expires after 7 days (168 hours)
- User doesn't need to login again within 7 days
- Provides good user experience with reasonable security

### Refresh Token
- Stored securely (httpOnly cookie recommended in production)
- Expires after 90 days
- Used to obtain a new access token without re-entering credentials
- If expired, user must login again

## Configuration Options

### Development Environment

Default settings are suitable for development:
```env
APP_JWT_ACCESS_EXPIRY_HOURS=168      # 7 days - good for testing
APP_JWT_REFRESH_EXPIRY_DAYS=90       # 90 days
```

**Pros:**
- Users stay logged in for a week
- Good for testing without frequent re-logins
- Sufficient security for development

### Production Environment

For production, consider these security vs. UX trade-offs:

#### Option 1: Balanced (Recommended)
```env
APP_JWT_ACCESS_EXPIRY_HOURS=24       # 1 day
APP_JWT_REFRESH_EXPIRY_DAYS=30       # 30 days
```
**Best for:** Most applications
- More secure than 7 days
- Users still don't need to login often
- Good balance for financial/health apps

#### Option 2: High Security
```env
APP_JWT_ACCESS_EXPIRY_HOURS=1        # 1 hour
APP_JWT_REFRESH_EXPIRY_DAYS=7        # 7 days
```
**Best for:** High-security applications (banking, healthcare)
- Most secure option
- Users need to re-authenticate more often
- Risk of frequent re-logins

#### Option 3: Long-Lived Sessions
```env
APP_JWT_ACCESS_EXPIRY_HOURS=720      # 30 days
APP_JWT_REFRESH_EXPIRY_DAYS=365      # 1 year
```
**Best for:** Low-risk internal applications
- Excellent user experience
- Lower security - only use for trusted networks
- Good for internal dashboards

## Implementation Details

### Token Generation

When a user logs in, two tokens are generated:

```go
// Access token
accessToken, err := jwtManager.GenerateAccessToken(user.ID, int64(user.ExpiryHours))

// Refresh token
refreshToken, err := jwtManager.GenerateRefreshToken(user.ID, int64(user.ExpiryDays))
```

### Token Refresh Flow

```
1. Client tries API request with expired access token
2. Backend returns 401 Unauthorized
3. Client sends refresh token to `/api/auth/refresh`
4. Backend validates refresh token
5. Backend issues new access token
6. Client retries original request with new token
```

### Refresh Implementation (Frontend)

Modern frontends should implement automatic token refresh:

```typescript
// Intercept 401 responses
if (response.status === 401) {
    const newToken = await refreshAccessToken();
    return retryRequest(newToken);
}
```

## Security Best Practices

### 1. Secret Key Management
- Must be at least 32 characters in production
- Generate with: `openssl rand -base64 32`
- Store in environment variables, never in code
- Rotate regularly (requires re-authentication of all users)

```bash
# Generate secure secrets
openssl rand -base64 32
```

### 2. Token Storage

**Frontend:**
- **DO:** Store refresh token in httpOnly cookie (protected from XSS)
- **DON'T:** Store tokens in localStorage (vulnerable to XSS)
- **DON'T:** Store tokens in sessionStorage (lost on browser close)

**Backend:**
- Validate tokens before processing requests
- Check token signature and expiry
- Maintain blacklist for revoked tokens if needed

### 3. Token Rotation
- Implement token rotation on every refresh
- Issue new refresh token when access token expires
- Invalidate old refresh tokens

### 4. HTTPS Requirement
- Always use HTTPS in production
- JWT tokens can be decoded (not encrypted)
- HTTPS prevents man-in-the-middle attacks

### 5. CORS Configuration
- Restrict `CORS_ALLOWED_ORIGINS` to your frontend domain
- Prevent token exposure to unauthorized origins

```env
APP_CORS_ALLOWED_ORIGINS=https://yourdomain.com
```

## Monitoring Token Expiry

### Debug Token Contents

To inspect token claims (development only):

```go
claims, err := jwtManager.VerifyToken(tokenString)
if err == nil {
    log.Printf("Token expires at: %v", time.Unix(claims.ExpiresAt, 0))
}
```

### Log Token Metrics

Monitor token usage:
- Number of refresh requests
- Token expiry rate
- Failed token validations
- Token generation errors

## Troubleshooting

### Users Getting 401 Frequently
**Problem:** Access token duration too short
**Solution:** Increase `APP_JWT_ACCESS_EXPIRY_HOURS`

### Session Lasting Too Long
**Problem:** Refresh token duration too long
**Solution:** Decrease `APP_JWT_REFRESH_EXPIRY_DAYS`

### Token Not Recognized
**Check:**
1. Secret keys match between frontend and backend
2. Token hasn't expired
3. Token signature is valid
4. Clock synchronization (NTP) on server

### Can't Change Token Duration
**Remember:**
- Existing tokens retain their original expiry
- New setting applies only to newly issued tokens
- User must login again for changes to take effect

## Production Checklist

- [ ] Change default secret keys (`openssl rand -base64 32`)
- [ ] Choose appropriate expiry times
- [ ] Implement token refresh in frontend
- [ ] Use HTTPS only
- [ ] Configure CORS properly
- [ ] Implement token blacklist if needed
- [ ] Monitor token-related errors
- [ ] Test token expiry handling
- [ ] Document your JWT policy

## References

- [JWT.io](https://jwt.io) - JWT specification
- [RFC 7519](https://tools.ietf.org/html/rfc7519) - JSON Web Token (JWT)
- [OWASP Session Management](https://owasp.org/www-project-web-security-testing-guide/latest/4-Web_Application_Security_Testing/06-Session_Management_Testing/README)

---

© 2026 Kodia Studio. "Build like a user, code like a pro."
