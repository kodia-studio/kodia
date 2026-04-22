# Advanced Security Features

Kodia provides a suite of enterprise-grade security features to protect your application from modern threats while providing fine-grained access control.

## Overview

The Kodia security architecture is designed to be **Defense in Depth**, combining multiple layers of protection:

1.  **Transport Security**: HSTS and automated Security Headers.
2.  **Authentication Hardening**: 2FA (TOTP), Refresh Token Rotation, and Session Tracking.
3.  **Access Control**: Granular RBAC and dynamic ABAC.
4.  **Monitoring**: Immutable Audit Logging.
5.  **Connectivity Control**: IP Whitelisting.

---

## 1. Role-Based Access Control (RBAC)

Kodia goes beyond simple "Role" strings by implementing a permission-based system. Each user can have a set of permissions stored in their JWT claims.

### Middleware Usage

```go
// Allow users with 'admin' role
api.GET("/settings", middleware.RequireRole("admin"))

// Require at least one of these permissions
api.POST("/posts", middleware.RequirePermission("posts:write", "posts:admin"))

// Require ALL of these permissions
api.DELETE("/posts", middleware.RequireAllPermissions("posts:write", "posts:delete"))
```

---

## 2. Two-Factor Authentication (2FA)

Protect user accounts with TOTP (Time-based One-Time Password) support.

### Setup Flow
Kodia provides utilities in `pkg/auth2fa` to handle the heavy lifting:

1.  **Generate Secret**: Create a secret for the user and show a QR Code.
2.  **Verify & Enable**: Confirm the user has the code before activating 2FA.

```go
secret, url, err := auth2fa.GenerateSecret(user.Email)
qrCodeBase64, err := auth2fa.GenerateQRCodeBase64(url)
```

---

## 3. Refresh Token Rotation & Reuse Detection

To prevent attackers from using stolen refresh tokens, Kodia implements **Automatic Rotation**.

-   Each time a refresh token is used, it is invalidated and a new one is issued.
-   **Reuse Detection**: If a previously used refresh token is presented again, Kodia detects a potential breach and **invalidates all active sessions** for that user immediately.

---

## 4. Session & Device Tracking

Monitor exactly where and when your users are logged in.

-   **Automatic Tracking**: Kodia records IP, UserAgent, and Device information for every session.
-   **Revocation**: Revoke specific devices or "Logout from all devices" with a single command.

---

## 5. Audit Logging

Every critical action should be traceable. The `AuditLogger` provides an immutable trail of events.

```go
auditLogger.LogAction(
    userID, 
    userEmail, 
    "POST-123", 
    audit.ActionDelete, 
    audit.StatusSuccess, 
    "User deleted a project", 
    ip, 
    ua,
)
```

---

## 6. Security Headers

Kodia automatically injects recommended security headers through the `SecurityHeaders` middleware:

-   `Content-Security-Policy` (CSP)
-   `Strict-Transport-Security` (HSTS)
-   `X-Frame-Options: DENY` (Prevent Clickjacking)
-   `X-Content-Type-Options: nosniff`
-   `Referrer-Policy`: Controls how much referrer information is included.

---

## 7. IP Whitelisting

Restrict access to sensitive routes (like `/admin` or `/internal-api`) only to specific IP addresses or CIDR ranges.

```go
allowedIPs := []string{"192.168.1.1", "10.0.0.0/24"}
admin.Use(middleware.IPWhitelist(allowedIPs))
```

---

## 8. Attribute-Based Access Control (ABAC)

For complex scenarios where roles aren't enough, Kodia's ABAC engine allows you to define policies based on attributes.

```go
evaluator := policy.NewEvaluator()
evaluator.AddPolicy(policy.Policy{
    Name: "Allow editing own posts",
    Effect: policy.EffectAllow,
    Condition: func(sub, obj, env policy.Attributes) bool {
        return sub["id"] == obj["author_id"]
    },
})
```

---

## 9. Social Authentication (OAuth2)

Kodia supports world-class social login providers out of the box. This allows users to register and login using their existing accounts from trusted platforms.

- **Providers**: Google, GitHub, and **Discord**.
- **Security**: State-based CSRF protection and automatic user account linking.
- **Implementation**: Easily add new providers by implementing the `AuthProvider` port.

---

## 10. Passkeys (WebAuthn)

Stay ahead of the curve with **Passkeys**. Kodia supports WebAuthn, allowing users to login using biometrics (FaceID, TouchID) or security keys without needing a password.

- **Phishing Resistant**: Cryptographically tied to your specific domain.
- **Seamless UX**: One-touch login for mobile and desktop users.
- **Enterprise Ready**: Full support for multi-device credential synchronization.

---

## 11. Scoped API Keys

For programmatic access, Kodia provides a robust API Key management system with granular scoping.

- **Granular Scopes**: Define exactly what an API Key can do (e.g., `posts:read`, `users:write`).
- **Rotation**: Easily generate and revoke keys from the administrative dashboard.
- **Audit Trail**: Every request made with an API Key is tied back to the owner for full accountability.

