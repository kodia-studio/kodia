# Role-Based Access Control (RBAC)

Kodia Framework provides a complete Role-Based Access Control (RBAC) system as an alternative to Attribute-Based Access Control (ABAC). RBAC is simpler and more intuitive for most applications where you need to organize users into roles with specific permissions.

---

## Overview

The RBAC system consists of three core concepts:

- **Roles**: Named groups of permissions (e.g., `admin`, `editor`, `viewer`)
- **Permissions**: Granular action names (e.g., `posts.create`, `users.delete`)
- **User Roles**: The assignment of roles to users

### Key Features

- **In-memory engine** for fast permission checks at request time
- **Database-backed** for persistent role and permission definitions
- **Wildcard permissions** — a `*` permission grants all actions
- **Middleware integration** — built into HTTP request handlers
- **Async sync** — loads from database at startup
- **Complementary to ABAC** — use both for different use cases

---

## Quick Start

### 1. Define Roles in Your Application

Typically done during application setup or via a seeder:

```go
roleService := app.Resolve[ports.RoleService]("role_service")

// Create roles with permissions
adminRole, _ := roleService.CreateRole(ctx,
    "admin",
    "Administrator with full access",
    []string{"*"}, // Wildcard: all permissions
)

editorRole, _ := roleService.CreateRole(ctx,
    "editor",
    "Can create and edit content",
    []string{
        "posts.create",
        "posts.edit",
        "posts.publish",
        "comments.moderate",
    },
)

readerRole, _ := roleService.CreateRole(ctx,
    "reader",
    "Read-only access",
    []string{
        "posts.view",
        "comments.view",
    },
)
```

### 2. Assign Roles to Users

```go
// Assign role to a user
_ = roleService.AssignRole(ctx, userID, "editor")

// Assign multiple roles to a user
_ = roleService.AssignRole(ctx, userID, "admin")
_ = roleService.AssignRole(ctx, userID, "editor")

// Revoke a role
_ = roleService.RevokeRole(ctx, userID, "editor")
```

### 3. Enforce RBAC in Handlers

```go
func (h *MyHandler) CreatePost(c *gin.Context) {
    // Require the user to have "editor" role
    if err := middleware.RequireRole("editor")(c); err != nil {
        return // Automatically returns 403
    }

    // Handler logic...
}

// Or with multiple roles (user must have at least one)
func (h *MyHandler) AdminDashboard(c *gin.Context) {
    if err := middleware.RequireRole("admin", "super_admin")(c); err != nil {
        return
    }

    // Handler logic...
}
```

### 4. Check Permissions at Request Time

```go
// In a handler
func (h *MyHandler) DeleteUser(c *gin.Context) {
    // Get RBAC engine
    rbac := app.Resolve[*policy.RBACEngine]("rbac")

    // Get user's roles from JWT or context
    userRoles := c.GetStringSlice("user_roles") // ["editor", "admin"]

    // Check if user can perform action
    if !rbac.Can(userRoles, policy.Permission("users.delete")) {
        response.Forbidden(c, "You do not have permission to delete users")
        return
    }

    // Perform deletion...
}
```

---

## Architecture

### RBAC Engine (`pkg/policy/rbac.go`)

In-memory registry of roles and their permissions. Loaded from database at startup.

```go
// Create engine
engine := policy.NewRBACEngine()

// Define roles
engine.DefineRole("admin",
    policy.Permission("*"),  // All permissions
)

engine.DefineRole("user",
    policy.Permission("profile.read"),
    policy.Permission("profile.update"),
    policy.Permission("posts.create"),
)

// Check permissions
can := engine.Can([]string{"user"}, policy.Permission("posts.create"))
// Result: true
```

### Role Service (`internal/core/services/role_service.go`)

Business logic for managing roles and syncing with the RBAC engine.

```go
type RoleService interface {
    CreateRole(ctx, name, description string, permissions []string) (*domain.RoleEntity, error)
    AssignRole(ctx, userID, roleName string) error
    RevokeRole(ctx, userID, roleName string) error
    GetUserRoles(ctx, userID string) ([]string, error)
    GetAllRoles(ctx) ([]*domain.RoleEntity, error)
    DeleteRole(ctx, id string) error
    SyncEngineFromDB(ctx) error  // Called at startup
}
```

### Database Schema

Four tables manage RBAC:

- **roles** — `id`, `name`, `description`, `created_at`, `updated_at`, `deleted_at`
- **permissions** — `id`, `name`, `description`, `group`, `created_at`, `deleted_at`
- **role_permissions** — join table linking roles to permissions
- **user_roles** — join table linking users to roles

---

## HTTP API

All role management endpoints are in `/api/admin/roles` and protected with `admin` role requirement.

### Create Role

```http
POST /api/admin/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "moderator",
  "description": "Can moderate user content",
  "permissions": ["posts.moderate", "comments.moderate", "users.ban"]
}
```

### List All Roles

```http
GET /api/admin/roles
Authorization: Bearer <token>
```

### Delete Role

```http
DELETE /api/admin/roles/:id
Authorization: Bearer <token>
```

### Assign Role to User

```http
POST /api/admin/users/:user_id/roles
Authorization: Bearer <token>
Content-Type: application/json

{
  "role_name": "editor"
}
```

### Revoke Role from User

```http
DELETE /api/admin/users/:user_id/roles/:role
Authorization: Bearer <token>
```

### Get User Roles

```http
GET /api/admin/users/:user_id/roles
Authorization: Bearer <token>
```

---

## Permission Naming Convention

Use a dot-separated hierarchical naming scheme:

```
<resource>.<action>

Examples:
  users.read
  users.create
  users.update
  users.delete
  users.ban

  posts.create
  posts.edit
  posts.delete
  posts.publish

  admin.access
  admin.audit
```

Organize permissions by group in the database:

```go
permRepo.Create(ctx, &domain.PermissionEntity{
    ID:    kodia.NewID(),
    Name:  "posts.publish",
    Group: "posts",
    Description: "Ability to publish posts",
})
```

---

## RBAC vs ABAC

| Aspect | RBAC | ABAC |
|---|---|---|
| **Complexity** | Simple | Complex |
| **Setup time** | Fast | Slower |
| **Use case** | Most applications | Fine-grained policies |
| **Performance** | Excellent | Good |
| **Example** | "admin", "editor" | "User who created the post AND IP in whitelist AND time < 5pm" |

**Use RBAC when:**
- You have clear role boundaries (admin, user, editor, viewer)
- Users have consistent permissions within roles
- You want simplicity and maintainability

**Use ABAC when:**
- Permissions depend on complex context (time, IP, user attributes)
- You need fine-grained, dynamic policy evaluation
- Rules can't be simplified into static roles

You can combine both: use RBAC for role assignment, and ABAC for nuanced permission checks.

---

## Middleware Usage

### RequireRole

Ensure user has at least one of the given roles:

```go
func MyHandler(c *gin.Context) {
    if err := middleware.RequireRole("admin", "super_admin")(c); err != nil {
        return // 403 Forbidden
    }

    // Handler logic...
}
```

Can also be used as middleware on routes:

```go
router.POST("/api/admin/users", middleware.RequireRole("admin"), adminHandler.CreateUser)
```

### RequirePermission

Check user has at least one permission:

```go
func MyHandler(c *gin.Context) {
    if err := middleware.RequirePermission("posts.delete", "admin.access")(c); err != nil {
        return // 403 Forbidden
    }

    // Handler logic...
}
```

---

## Best Practices

✅ **Do:**
- Use descriptive role and permission names
- Organize permissions hierarchically by resource (`posts.create`, `posts.edit`)
- Use wildcard `*` sparingly (only for super admin)
- Sync the RBAC engine after role/permission changes
- Create a seeder to initialize default roles and permissions
- Document your permission structure

❌ **Don't:**
- Create overly granular roles for every user combination
- Mix RBAC and ABAC without clear boundaries
- Store permissions on the User entity directly (use roles)
- Forget to sync the RBAC engine when manually updating the database
- Hard-code role names in your handlers (use constants)

---

## Examples

### Creating a Permission Structure

```go
// In a seeder
permissions := []struct {
    name  string
    group string
}{
    // User permissions
    {"users.view", "users"},
    {"users.create", "users"},
    {"users.edit", "users"},
    {"users.delete", "users"},
    {"users.ban", "users"},

    // Post permissions
    {"posts.create", "posts"},
    {"posts.edit", "posts"},
    {"posts.delete", "posts"},
    {"posts.publish", "posts"},

    // Admin permissions
    {"admin.access", "admin"},
    {"audit.view", "admin"},
}

for _, p := range permissions {
    permRepo.Create(ctx, &domain.PermissionEntity{
        ID:    kodia.NewID(),
        Name:  p.name,
        Group: p.group,
    })
}
```

### Checking Permissions in Handlers

```go
func (h *PostHandler) DeletePost(c *gin.Context) {
    postID := c.Param("id")
    userID := c.GetString("user_id")

    // Get RBAC engine
    rbac := h.app.Resolve[*policy.RBACEngine]("rbac")
    userRoles := c.GetStringSlice("user_roles")

    // User can delete if they have the permission OR created the post themselves
    canDelete := rbac.Can(userRoles, policy.Permission("posts.delete"))
    isAuthor, _ := h.postService.IsAuthor(ctx, postID, userID)

    if !canDelete && !isAuthor {
        response.Forbidden(c, "You cannot delete this post")
        return
    }

    h.postService.DeletePost(ctx, postID)
    response.OK(c, "Post deleted", nil)
}
```

---

**Last Updated**: April 2026  
**Framework Version**: v1.7.0+
