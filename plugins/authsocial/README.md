# Kodia Social Login Plugin

Complete OAuth2 social authentication plugin for Google and GitHub integration with Kodia Framework.

## Features

- ✅ Google OAuth2 integration
- ✅ GitHub OAuth2 integration
- ✅ Automatic user creation from social profiles
- ✅ Account linking for existing users
- ✅ CSRF protection with state parameter
- ✅ JWT token generation (access + refresh)
- ✅ Email verification from GitHub emails endpoint

## Installation

### 1. Install the Plugin

```bash
kodia plugin install authsocial
```

Or manually add to `go.mod`:

```
require github.com/kodia-studio/authsocial v0.1.0
```

### 2. Register in `cmd/server/main.go`

Import the plugin:

```go
import authsocial "github.com/kodia-studio/authsocial"
```

Add to `app.RegisterProviders()` (after `AuthProvider`):

```go
app.RegisterProviders(
    providers.NewDatabaseProvider(),
    providers.NewInfraProvider(),
    providers.NewHttpProvider(),
    providers.NewAuthProvider(),
    authsocial.NewServiceProvider(),  // Add this
    // ... other providers
)
```

### 3. Set Environment Variables

```env
# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=https://yourapp.com/api/auth/social/google/callback

# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_REDIRECT_URL=https://yourapp.com/api/auth/social/github/callback
```

**Optional:** If `GOOGLE_REDIRECT_URL` and `GITHUB_REDIRECT_URL` are not set, they default to:
- `{APP_BASE_URL}/api/auth/social/google/callback`
- `{APP_BASE_URL}/api/auth/social/github/callback`

### 4. Run Database Migrations

```bash
kodia db:migrate --path=vendor/github.com/kodia-studio/authsocial/migrations
```

Or if using `golang-migrate`:

```bash
migrate -path vendor/github.com/kodia-studio/authsocial/migrations -database "postgres://user:pass@localhost/dbname" up
```

## Setup OAuth Applications

### Google

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Create a new OAuth 2.0 Client ID (Web application)
3. Set Authorized redirect URI to: `https://yourapp.com/api/auth/social/google/callback`
4. Copy Client ID and Client Secret to your `.env`

### GitHub

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Create a new OAuth App
3. Set Authorization callback URL to: `https://yourapp.com/api/auth/social/github/callback`
4. Copy Client ID and Client Secret to your `.env`

## API Endpoints

### Initiate Login

```
GET /api/auth/social/:provider/redirect
```

Initiates the OAuth flow and redirects to the provider's login page.

**Parameters:**
- `:provider` - `google` or `github`

**Example:**
```html
<a href="/api/auth/social/google/redirect">Login with Google</a>
<a href="/api/auth/social/github/redirect">Login with GitHub</a>
```

### OAuth Callback

```
GET /api/auth/social/:provider/callback
```

Handles the OAuth callback from the provider. Automatically:
- Exchanges authorization code for token
- Fetches user information
- Creates user account (if new)
- Links social account
- Redirects to frontend with authentication tokens

**Query Parameters:**
- `code` - Authorization code from provider
- `state` - CSRF protection token
- `error` - Error message (if login failed)

**Redirect on Success:**
```
{FRONTEND_URL}/auth/social/success?token=<JWT_ACCESS_TOKEN>&refresh=<REFRESH_TOKEN>
```

**Redirect on Error:**
```
{FRONTEND_URL}/auth/error?message=<ERROR_MESSAGE>
```

## Frontend Integration

### Step 1: Handle OAuth Redirect

```svelte
<a href="https://yourapi.com/api/auth/social/google/redirect" class="btn">
  Login with Google
</a>
```

### Step 2: Handle Callback

After user logs in with the provider, they're redirected to:

```
https://yourfrontend.com/auth/social/success?token=eyJhbGc...&refresh=eyJhbGc...
```

Create a page to handle this:

```svelte
<!-- src/routes/auth/social/success/+page.svelte -->
<script>
  import { authStore } from '$lib/stores/auth.store';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';

  const token = $page.url.searchParams.get('token');
  const refreshToken = $page.url.searchParams.get('refresh');

  if (token && refreshToken) {
    // Save tokens to store
    authStore.login(null, token, refreshToken);
    // Redirect to dashboard
    goto('/dashboard');
  }
</script>

<div>Logging you in...</div>
```

### Step 3: Error Handling

Handle errors at:

```
https://yourfrontend.com/auth/error?message=<ERROR_MESSAGE>
```

```svelte
<!-- src/routes/auth/error/+page.svelte -->
<script>
  import { page } from '$app/stores';

  const message = $page.url.searchParams.get('message');
</script>

<div class="error">
  <h1>Login Failed</h1>
  <p>{message}</p>
  <a href="/login">Back to Login</a>
</div>
```

## How It Works

1. **Redirect Phase:**
   - Generate CSRF state token, store in Redis (5min TTL)
   - Redirect user to OAuth provider

2. **Provider Login:**
   - User logs in and authorizes app
   - Provider redirects back with `code` and `state`

3. **Callback Phase:**
   - Verify state token (CSRF protection)
   - Exchange code for access token
   - Fetch user info from provider
   - Check if social account already linked
   - If not linked, check if email exists (link) or create new user
   - Generate JWT access + refresh tokens
   - Redirect to frontend with tokens

4. **Frontend:**
   - Receive tokens and save to localStorage
   - User is now authenticated

## Database Schema

### social_accounts table

```sql
CREATE TABLE social_accounts (
    id          UUID PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES users(id),
    provider    VARCHAR(50),      -- 'google' or 'github'
    provider_id VARCHAR(255),     -- Provider's user ID
    email       VARCHAR(255),
    name        VARCHAR(255),
    avatar_url  TEXT,
    created_at  TIMESTAMPTZ,
    UNIQUE(provider, provider_id)
);
```

## Configuration

The plugin reads configuration from environment variables:

| Variable | Description | Default |
|---|---|---|
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | (required if using Google) |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | (required if using Google) |
| `GOOGLE_REDIRECT_URL` | Google OAuth redirect URL | `{APP_BASE_URL}/api/auth/social/google/callback` |
| `GITHUB_CLIENT_ID` | GitHub OAuth Client ID | (required if using GitHub) |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth Client Secret | (required if using GitHub) |
| `GITHUB_REDIRECT_URL` | GitHub OAuth redirect URL | `{APP_BASE_URL}/api/auth/social/github/callback` |

**Note:** At least one provider must be configured.

## Troubleshooting

### "Invalid or expired state parameter"

- The state was not stored or expired
- Check Redis cache is working
- Check state parameter is being passed correctly

### "Failed to get user information"

- OAuth credentials might be incorrect
- Network connectivity issue
- Provider API might be down

### "Email not found"

- GitHub user might not have public email
- Plugin automatically fetches from GitHub emails endpoint
- If still failing, user must update GitHub email visibility

## Architecture

The plugin follows Kodia's Hexagonal Architecture pattern:

- **Entity:** `SocialAccount` domain entity
- **Repository:** `SocialAccountRepository` interface with GORM implementation
- **Service:** `SocialAuthService` handles OAuth flow logic
- **Handler:** `SocialHandler` HTTP endpoints
- **Provider:** `ServiceProvider` framework integration

## Development

### Testing Locally

1. Set env vars with test OAuth credentials
2. Run migrations: `kodia db:migrate --path=plugins/authsocial/migrations`
3. Start server: `kodia dev`
4. Visit: `http://localhost:3000/api/auth/social/google/redirect`

## License

MIT License - Kodia Studio
