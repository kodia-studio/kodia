# Social Login Plugin — Documentation

Complete guide to using the Kodia Social Login Plugin for Google and GitHub OAuth2 authentication.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [API Reference](#api-reference)
6. [Authentication Flow](#authentication-flow)
7. [Frontend Integration](#frontend-integration)
8. [Database Schema](#database-schema)
9. [Error Handling](#error-handling)
10. [Security Considerations](#security-considerations)
11. [Troubleshooting](#troubleshooting)

---

## Overview

The **Social Login Plugin** (`github.com/kodia-studio/authsocial`) provides seamless OAuth2 authentication with Google and GitHub for Kodia Framework applications.

### Key Features

- **Multi-provider support** — Google and GitHub out of the box
- **Automatic user creation** — Creates users from social profiles on first login
- **Account linking** — Links social accounts to existing email addresses
- **CSRF protection** — State parameter validation with Redis cache
- **JWT tokens** — Generates access and refresh tokens
- **Email fallback** — Fetches email from GitHub's user/emails endpoint if needed
- **Extensible design** — Add more OAuth providers by implementing `Provider` interface

### When to Use

Use this plugin when you need:
- Quick OAuth2 setup without building from scratch
- Support for multiple OAuth providers
- Automatic user account provisioning from social profiles
- Existing user linking to social accounts

---

## Architecture

### Components

```
┌─────────────────┐
│  OAuth Provider │
│  (Google/GitHub)│
└────────┬────────┘
         │
         │ Authorization Code
         │
┌────────▼────────────────────────┐
│  SocialHandler (HTTP)            │
│  - Redirect endpoint             │
│  - Callback endpoint             │
└────────┬──────────────────────────┘
         │
┌────────▼──────────────────────────┐
│  SocialAuthService                 │
│  - CSRF state management           │
│  - Provider Exchange & GetUser      │
│  - User find-or-create logic       │
│  - JWT token generation            │
└────────┬──────────────────────────┘
         │
┌────────▼──────────────────────────┐
│  Repositories                      │
│  - SocialAccountRepository         │
│  - UserRepository (existing)       │
└────────┬──────────────────────────┘
         │
┌────────▼──────────────────────────┐
│  Database                          │
│  - users (existing)                │
│  - social_accounts (new)           │
└────────────────────────────────────┘
```

### Data Flow

```
1. User clicks "Login with Google"
   ↓
2. GET /api/auth/social/google/redirect
   → Generate state, store in Redis (5min)
   → Redirect to Google OAuth URL
   ↓
3. User logs in on Google, authorizes app
   ↓
4. Google redirects to callback with code + state
   ↓
5. GET /api/auth/social/google/callback?code=...&state=...
   → Verify state (CSRF check)
   → Exchange code for token
   → Fetch user info from Google
   → Check if social account linked
     - Yes: Use existing user
     - No: Check if email exists
       - Yes: Link to existing user
       - No: Create new user
   → Create social_accounts record
   → Generate JWT tokens
   → Redirect to frontend with tokens
   ↓
6. Frontend stores tokens, user logged in
```

---

## Installation

### Prerequisites

- Kodia Framework project (v0.1.0+)
- Go 1.25+
- PostgreSQL database with `users` table
- Redis for state caching

### Step 1: Install Plugin

```bash
# Via CLI
kodia plugin install authsocial

# Or manually add to go.mod
require github.com/kodia-studio/authsocial v0.1.0
```

### Step 2: Register in `cmd/server/main.go`

```go
package main

import (
    // ... other imports
    authsocial "github.com/kodia-studio/authsocial"
    "github.com/kodia-studio/kodia/internal/providers"
)

func main() {
    // ... initialization code ...

    err = app.RegisterProviders(
        providers.NewDatabaseProvider(),
        providers.NewInfraProvider(),
        providers.NewHttpProvider(),
        providers.NewAuthProvider(),
        authsocial.NewServiceProvider(),  // ← Add this
        // ... other providers
    )
    if err != nil {
        log.Fatal(err)
    }

    app.Boot()
    app.Run()
}
```

**Important:** Register `authsocial.NewServiceProvider()` AFTER `providers.NewAuthProvider()` because the social login service depends on JWT manager.

### Step 3: Run Migrations

```bash
# Using kodia CLI
kodia db:migrate --path=vendor/github.com/kodia-studio/authsocial/migrations

# Or using golang-migrate directly
migrate -path vendor/github.com/kodia-studio/authsocial/migrations \
    -database "postgres://user:pass@localhost/dbname" up
```

This creates the `social_accounts` table.

---

## Configuration

### Environment Variables

```env
# Google OAuth (one or both providers must be configured)
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=https://yourapp.com/api/auth/social/google/callback

# GitHub OAuth
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=https://yourapp.com/api/auth/social/github/callback

# Application URLs
APP_BASE_URL=https://yourapp.com
FRONTEND_URL=https://yourfrontend.com
```

### Auto-Configuration

If `GOOGLE_REDIRECT_URL` or `GITHUB_REDIRECT_URL` are not set, they default to:

```
{APP_BASE_URL}/api/auth/social/{provider}/callback
```

Example: If `APP_BASE_URL=https://api.example.com`:
- Google callback: `https://api.example.com/api/auth/social/google/callback`
- GitHub callback: `https://api.example.com/api/auth/social/github/callback`

---

## API Reference

### Redirect Endpoint

Initiates the OAuth2 flow.

```http
GET /api/auth/social/:provider/redirect
```

**Parameters:**
- `:provider` (string) — `google` or `github`

**Response:**
- HTTP 302 redirect to OAuth provider

**Example:**
```bash
curl -L https://api.example.com/api/auth/social/google/redirect
# Redirects to: https://accounts.google.com/o/oauth2/v2/auth?...
```

**Usage in Frontend:**
```html
<a href="https://api.example.com/api/auth/social/google/redirect">
  Login with Google
</a>
```

### Callback Endpoint

Handles OAuth provider callback.

```http
GET /api/auth/social/:provider/callback
```

**Query Parameters:**
- `code` (string) — Authorization code from provider
- `state` (string) — CSRF state token
- `error` (string, optional) — Error from provider

**Response:**
- **Success:** HTTP 302 redirect to `{FRONTEND_URL}/auth/social/success?token=...&refresh=...`
- **Error:** HTTP 302 redirect to `{FRONTEND_URL}/auth/error?message=...`

**Success Redirect Example:**
```
https://yourfrontend.com/auth/social/success?
  token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...&
  refresh=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Error Redirect Example:**
```
https://yourfrontend.com/auth/error?message=invalid+state+parameter
```

---

## Authentication Flow

### Step-by-Step Flow

#### 1. Initiate Login

User clicks "Login with Google" button:

```
GET /api/auth/social/google/redirect
```

Backend:
- Generates random UUID as state token
- Stores in Redis with key `social_state:<uuid>` and 5-minute TTL
- Redirects user to Google's OAuth URL with state parameter

#### 2. User Authorization

User logs into Google and approves app permissions.

#### 3. OAuth Callback

Google redirects user back:

```
GET /api/auth/social/google/callback?code=...&state=...
```

Backend:
1. **CSRF Verification:** Check state in Redis
   - If missing/expired: reject with error
   - If valid: delete from Redis
2. **Exchange:** Call `provider.Exchange(code)` → get access token
3. **User Info:** Call `provider.GetUser(token)` → get user data
4. **Account Resolution:**
   ```
   IF social_account(provider=google, provider_id=X) exists:
       user = existing_user
   ELSE:
       IF user(email=Y) exists:
           user = existing_user
           CREATE social_account(user_id=user.id)
       ELSE:
           user = CREATE new_user(email=Y, name=Z, is_verified=true)
           CREATE social_account(user_id=user.id)
   ```
5. **Token Generation:**
   - Generate JWT access token (1 hour TTL)
   - Generate JWT refresh token (7 day TTL)
   - Store refresh token in database
6. **Redirect:** Redirect to frontend with tokens

#### 4. Frontend Token Storage

Frontend receives callback redirect:

```
https://frontend.com/auth/social/success?token=...&refresh=...
```

Frontend stores tokens:
- Access token → localStorage (or HttpOnly cookie)
- Refresh token → localStorage (or HttpOnly cookie)

---

## Frontend Integration

### Step 1: Login Links

Create links to initiate OAuth:

```svelte
<!-- src/routes/auth/login/+page.svelte -->
<script>
  const apiUrl = import.meta.env.VITE_API_URL || 'https://api.example.com';
</script>

<div class="login-container">
  <h1>Login to Your App</h1>
  
  <!-- Google Login -->
  <a href="{apiUrl}/api/auth/social/google/redirect" class="btn btn-google">
    <img src="/google-icon.svg" alt="Google" />
    Login with Google
  </a>

  <!-- GitHub Login -->
  <a href="{apiUrl}/api/auth/social/github/redirect" class="btn btn-github">
    <img src="/github-icon.svg" alt="GitHub" />
    Login with GitHub
  </a>
</div>

<style>
  .login-container { ... }
  .btn { ... }
</style>
```

### Step 2: Callback Handler Page

Create a page to handle OAuth callback:

```svelte
<!-- src/routes/auth/social/success/+page.svelte -->
<script lang="ts">
  import { authStore } from '$lib/stores/auth.store';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';

  const token = $page.url.searchParams.get('token');
  const refreshToken = $page.url.searchParams.get('refresh');

  if (token && refreshToken) {
    // Store tokens in auth store
    authStore.login(null, token, refreshToken);
    
    // Redirect to dashboard
    goto('/dashboard', { replaceHistory: true });
  } else {
    // No tokens, something went wrong
    goto('/auth/error?message=missing_tokens', { replaceHistory: true });
  }
</script>

<div class="loading">
  <p>Logging you in...</p>
</div>
```

### Step 3: Error Handler Page

Handle login errors:

```svelte
<!-- src/routes/auth/error/+page.svelte -->
<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

  const message = $page.url.searchParams.get('message') || 'Unknown error';
  const errorMap: Record<string, string> = {
    'invalid_state': 'Security validation failed. Please try again.',
    'invalid_code': 'Authorization failed. Please try again.',
    'email_not_found': 'Could not retrieve email from social profile.',
    'network_error': 'Network error. Please check your connection.',
  };

  const displayMessage = errorMap[message] || message;
</script>

<div class="error-container">
  <h1>Login Failed</h1>
  <p class="error-message">{displayMessage}</p>
  <a href="/auth/login" class="btn btn-primary">Back to Login</a>
</div>
```

### Step 4: Auth Store Integration

Update your auth store to handle social login tokens:

```typescript
// src/lib/stores/auth.store.ts
export const authStore = createAuthStore();

function createAuthStore() {
  // ... existing code ...
  
  return {
    // ... existing methods ...
    
    // Called from social callback handler
    login: (user: User | null, token: string, refreshToken: string) => {
      localStorage.setItem('access_token', token);
      localStorage.setItem('refresh_token', refreshToken);
      document.cookie = `access_token=${token}; path=/; max-age=${7 * 24 * 60 * 60}`;
      
      set({
        user, // Can be null for social login, fetch via /api/me later
        accessToken: token,
        refreshToken,
        isAuthenticated: true,
        isLoading: false
      });
    }
  };
}
```

---

## Database Schema

### social_accounts Table

```sql
CREATE TABLE social_accounts (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider    VARCHAR(50) NOT NULL,      -- 'google' or 'github'
    provider_id VARCHAR(255) NOT NULL,     -- Provider's unique user ID
    email       VARCHAR(255),              -- User's email from provider
    name        VARCHAR(255),              -- User's name from provider
    avatar_url  TEXT,                      -- User's avatar URL
    created_at  TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(provider, provider_id)
);

CREATE INDEX idx_social_accounts_user_id ON social_accounts(user_id);
```

### Relationships

```
users (1) ──────────────────── (many) social_accounts
  ↓
  id                            user_id
  email                         provider + provider_id (unique)
  name
  avatar_url
  is_verified
  ...
```

### Example Data

```sql
-- User who signed up via Google
INSERT INTO users (id, email, name, is_verified) 
VALUES ('550e8400-e29b-41d4-a716-446655440001', 'john@example.com', 'John Doe', true);

INSERT INTO social_accounts (user_id, provider, provider_id, email, name, avatar_url)
VALUES (
  '550e8400-e29b-41d4-a716-446655440001',
  'google',
  '118203558960537612345',
  'john@example.com',
  'John Doe',
  'https://lh3.googleusercontent.com/...'
);

-- Same user later linking GitHub account
INSERT INTO social_accounts (user_id, provider, provider_id, email, name, avatar_url)
VALUES (
  '550e8400-e29b-41d4-a716-446655440001',
  'github',
  '12345678',
  'john@example.com',
  'John Doe',
  'https://avatars.githubusercontent.com/u/12345678'
);
```

---

## Error Handling

### Common Errors

#### Invalid State Parameter
```
GET /api/auth/social/google/callback?code=...&state=old_state
```

**Response:**
```
302 → {FRONTEND_URL}/auth/error?message=invalid+or+expired+state+parameter
```

**Cause:** State was not stored or expired (>5 minutes)
**Solution:** User should restart login flow

#### Invalid Authorization Code
```
GET /api/auth/social/google/callback?code=invalid_code&state=...
```

**Response:**
```
302 → {FRONTEND_URL}/auth/error?message=failed+to+exchange+authorization+code
```

**Cause:** Code is invalid or expired
**Solution:** User should restart login flow

#### Missing Email
**For GitHub:** User has private email and hasn't made any email public

**Cause:** GitHub API returns no public email in `user/emails` endpoint
**Solution:** User must update privacy settings in GitHub settings

#### User Creation Failed
**Response:**
```
302 → {FRONTEND_URL}/auth/error?message=failed+to+create+user
```

**Cause:** Database error (constraint violation, etc.)
**Solution:** Check application logs

### Error Response Format

All errors redirect to:
```
{FRONTEND_URL}/auth/error?message={error_message}
```

Message is URL-encoded. Common error messages:

| Message | Meaning |
|---------|---------|
| `invalid_or_expired_state_parameter` | CSRF validation failed |
| `failed_to_exchange_authorization_code` | OAuth code exchange failed |
| `failed_to_get_user_information` | Couldn't fetch user data from provider |
| `failed_to_create_user` | Database error creating user |
| `failed_to_link_social_account` | Database error linking account |
| `unsupported_provider` | Provider other than google/github |

---

## Security Considerations

### CSRF Protection

**Implementation:**
- Generate random UUID for each redirect
- Store in Redis with 5-minute TTL
- Verify on callback before processing
- Delete state after validation

**Why needed:**
Prevents cross-site request forgery in OAuth flow

### State Parameter Validation

```
✓ Valid: State exists in Redis AND not expired
✗ Invalid: State not in Redis OR expired
✗ Invalid: State mismatch
```

### Token Security

- **Access Token:** Short-lived (1 hour), stored in localStorage or memory
- **Refresh Token:** Long-lived (7 days), should be stored securely
- **Best Practice:** Use HttpOnly cookies for refresh token

### Email Verification

- Social login users have `is_verified=true` automatically
- Rationale: Trusting OAuth providers for email verification
- Alternative: Set `is_verified=false` if your app requires additional verification

### Provider Credentials

**Never commit:**
```env
GOOGLE_CLIENT_SECRET=xxx      # ✗ DO NOT COMMIT
GITHUB_CLIENT_SECRET=xxx      # ✗ DO NOT COMMIT
```

**Always use:**
- Environment variables
- Secrets manager
- `.env.local` (add to `.gitignore`)

### Scopes

**Google Scopes:**
```
https://www.googleapis.com/auth/userinfo.email
https://www.googleapis.com/auth/userinfo.profile
```

Minimal required for basic profile info.

**GitHub Scopes:**
```
user:email
```

Allows reading user email.

---

## Troubleshooting

### "Invalid or expired state parameter"

**Symptoms:**
- Redirect shows `?message=invalid+or+expired+state+parameter`

**Causes:**
1. Redis cache not working
2. State stored but expired (>5 minutes elapsed)
3. Browser cache issue with state parameter

**Solutions:**
```bash
# 1. Check Redis connection
redis-cli ping  # Should return PONG

# 2. Check Redis has data
redis-cli KEYS "social_state:*"

# 3. Clear browser cache and retry
# Ctrl+Shift+Delete → Clear Cookies & Cache

# 4. Check server logs for Redis errors
docker logs backend-container | grep redis
```

### "Failed to get user information"

**Symptoms:**
- Callback fails with network error
- Provider API unreachable

**Causes:**
1. OAuth provider API down
2. Network connectivity issue
3. Invalid access token
4. Provider rate limiting

**Solutions:**
```bash
# 1. Check provider API status
# Google: https://www.google.com/appsstatus
# GitHub: https://www.githubstatus.com/

# 2. Verify network connectivity
curl -v https://api.github.com/user

# 3. Check logs for specific error
docker logs backend-container | grep "failed to get user"

# 4. Wait if rate limited (GitHub: 60 req/hour unauthenticated)
```

### "Email not found" (GitHub)

**Symptoms:**
- User created but without email
- Possible in user records: `email = ""`

**Causes:**
- GitHub user has no public email
- All emails are private

**Solutions:**
1. User updates GitHub settings:
   - Go to github.com/settings/emails
   - Make at least one email public
   - Try login again
2. Or provide custom email in app

### Port/URL Mismatch

**Symptoms:**
- Redirect URL invalid in OAuth provider console
- Callback fails with "Redirect URI mismatch"

**Causes:**
- Local dev: `http://localhost:3000` vs `http://127.0.0.1:3000`
- Prod: `https://app.com` but configured `http://app.com`

**Solution:**
Ensure redirect URLs match exactly in:
1. `.env` file (GOOGLE_REDIRECT_URL, GITHUB_REDIRECT_URL)
2. OAuth app settings (Google Cloud Console, GitHub Developer Settings)

Example:
```env
# .env
APP_BASE_URL=https://api.example.com
GOOGLE_REDIRECT_URL=https://api.example.com/api/auth/social/google/callback
GITHUB_REDIRECT_URL=https://api.example.com/api/auth/social/github/callback
```

Must match exactly in provider consoles.

---

## Next Steps

- [Plugin README](../plugins/authsocial/README.md) — Installation and usage guide
- [Framework Documentation](./ARCHITECTURE.md) — Overall framework design
- [JWT Security](./JWT_SECURITY.md) — Token security details
- [Database Guide](./ORM_GUIDE.md) — Database schema and queries
