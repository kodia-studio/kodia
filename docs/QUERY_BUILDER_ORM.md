# 🔨 Database Query Builder & Advanced ORM Guide

**Enterprise-grade fluent query building** with type-safe generics and immutable patterns. Build complex SQL queries without writing raw SQL while maintaining full type safety.

**Table of Contents:**
- [What is QueryBuilder?](#what-is-querybuilder)
- [Why Use QueryBuilder?](#why-use-querybuilder)
- [Quick Start](#quick-start)
- [Filtering](#filtering)
- [Selection & Ordering](#selection--ordering)
- [Eager Loading](#eager-loading)
- [Pagination](#pagination)
- [Soft Deletes](#soft-deletes)
- [Advanced Operations](#advanced-operations)
- [Generic BaseRepository](#generic-baserepository)
- [Best Practices](#best-practices)
- [Examples & Patterns](#examples--patterns)

---

## What is QueryBuilder?

QueryBuilder is a **fluent, chainable wrapper** around GORM's `*gorm.DB` that provides a clean, expressive API for building database queries in Go.

### Features

| Feature | Description |
|---------|-------------|
| **Fluent Chaining** | `Where(...).OrderDesc(...).Limit(...).Find(...)` |
| **Type-Safe** | Compile-time type checking with Go generics |
| **Immutable** | Each method returns a new QueryBuilder instance |
| **GORM Native** | Works seamlessly with GORM v1.31.1+ |
| **SQL Injection Safe** | Parameterized queries—never concatenate user input |
| **No String Building** | Methods generate SQL, not string concatenation |

### Example

```go
// Before: verbose raw GORM
var users []gormUser
baseQuery := r.db.WithContext(ctx).Model(&gormUser{})
if err := baseQuery.Where("role = ?", "admin").
    OrderBy("created_at DESC").
    Limit(20).
    Find(&users).Error; err != nil {
    return nil, err
}

// After: fluent QueryBuilder
var users []gormUser
if err := r.Query().
    Where("role = ?", "admin").
    OrderDesc("created_at").
    Limit(20).
    Find(ctx, &users); err != nil {
    return nil, err
}
```

---

## Why Use QueryBuilder?

### ✅ Developer Experience
- **Readable syntax** — Queries read like English
- **IDE autocomplete** — All methods discoverable
- **No string formatting** — Safe from SQL injection
- **Chainable API** — Build queries progressively

### ✅ Code Quality
- **Less boilerplate** — No repetitive `db.WithContext(ctx).Where(...)` patterns
- **Consistent patterns** — All repos use the same query style
- **Easier refactoring** — Queries are composable, not scattered
- **Type-safe** — Compile errors for typos, not runtime errors

### ✅ Performance
- **Lazy evaluation** — Queries execute only at terminal operations
- **Single query execution** — No redundant queries
- **Eager loading support** — Prevent N+1 query problems
- **Index-friendly** — Compatible with database indexes

---

## Quick Start

### Basic Query

```go
import "github.com/kodia-studio/kodia/pkg/database"

// Create a QueryBuilder
qb := database.NewQuery(db).
    Where("status = ?", "active").
    OrderDesc("created_at")

// Execute: Find all matching records
var users []gormUser
if err := qb.Find(ctx, &users); err != nil {
    log.Fatal(err)
}
```

### Using BaseRepository

```go
type UserRepository struct {
    database.BaseRepository[gormUser]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{
        BaseRepository: database.NewBaseRepository[gormUser](db),
    }
}

// Now use Query() method
func (r *UserRepository) FindActive(ctx context.Context) ([]gormUser, error) {
    var users []gormUser
    err := r.Query().
        Where("is_active = ?", true).
        OrderAsc("name").
        Find(ctx, &users)
    return users, err
}
```

---

## Filtering

### Basic WHERE Conditions

```go
// Single condition
r.Query().Where("role = ?", "admin").Find(ctx, &users)

// Multiple AND conditions
r.Query().
    Where("role = ?", "admin").
    Where("is_active = ?", true).
    Find(ctx, &users)

// OR conditions
r.Query().
    Where("role = ?", "admin").
    OrWhere("role = ?", "moderator").
    Find(ctx, &users)

// NOT conditions
r.Query().
    Not("role = ?", "guest").
    Find(ctx, &users)
```

### IN / NOT IN Clauses

```go
// IN
statuses := []string{"active", "pending", "approved"}
r.Query().WhereIn("status", statuses).Find(ctx, &users)

// NOT IN
roles := []string{"guest", "banned"}
r.Query().WhereNotIn("role", roles).Find(ctx, &users)
```

### Range Queries

```go
// BETWEEN
r.Query().
    WhereBetween("created_at", startDate, endDate).
    Find(ctx, &users)

// NULL checks
r.Query().WhereNull("deleted_at").Find(ctx, &users)
r.Query().WhereNotNull("avatar_url").Find(ctx, &users)
```

---

## Selection & Ordering

### Column Selection

```go
// Select specific columns
var users []gormUser
r.Query().
    Select("id", "name", "email").
    Find(ctx, &users)
```

### Ordering

```go
// Ascending
r.Query().OrderAsc("name").Find(ctx, &users)

// Descending
r.Query().OrderDesc("created_at").Find(ctx, &users)

// Raw ORDER BY
r.Query().OrderRaw("RAND()").Find(ctx, &users)

// Multiple order clauses
r.Query().
    OrderDesc("created_at").
    OrderAsc("name").
    Find(ctx, &users)
```

### Limiting & Offsetting

```go
// Limit results
r.Query().Limit(10).Find(ctx, &users)

// Offset (for manual pagination)
r.Query().Offset(20).Limit(10).Find(ctx, &users)

// Using pagination helper
params := &pagination.Params{Page: 2, PerPage: 20}
r.Query().Paginate(params).Find(ctx, &users)
```

---

## Eager Loading

### Preload (Prevent N+1 Queries)

```go
// Load relations
type gormUser struct {
    ID    string
    Name  string
    Roles []gormRole `gorm:"many2many:user_roles;"`
}

// Without preload: 1 + N queries
var users []gormUser
r.Query().Find(ctx, &users)  // 1 query, then N queries for each user's roles

// With preload: 2 queries total
var users []gormUser
r.Query().Preload("Roles").Find(ctx, &users)

// Nested preload
r.Query().
    Preload("Roles").
    Preload("Roles.Permissions").
    Find(ctx, &users)

// Conditional preload
r.Query().
    Preload("Roles", "is_active = ?", true).
    Find(ctx, &users)
```

### Joins (for filtering on relations)

```go
// INNER JOIN
var users []gormUser
r.Query().
    Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
    Joins("INNER JOIN roles ON roles.id = user_roles.role_id").
    Where("roles.name = ?", "admin").
    Find(ctx, &users)

// LEFT JOIN
r.Query().
    Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
    Find(ctx, &users)

// Combine Joins + Preload for efficiency
r.Query().
    Joins("INNER JOIN user_roles ON user_roles.user_id = users.id").
    Where("user_roles.role_id = ?", adminRoleID).
    Preload("Roles").
    Find(ctx, &users)
```

---

## Pagination

### Paginated Queries

```go
type User struct {
    ID    string
    Name  string
    Email string
}

// Basic pagination
params := &pagination.Params{Page: 1, PerPage: 20}
var users []gormUser
total, err := r.Query().
    Where("is_active = ?", true).
    OrderDesc("created_at").
    FindPaginated(ctx, &users, params)

// Result: users = 20 records, total = total count across all pages
```

### Handling Pagination in HTTP Handlers

```go
func (h *UserHandler) ListUsers(c *gin.Context) {
    page := c.DefaultQuery("page", "1")
    perPage := c.DefaultQuery("per_page", "20")
    
    params := &pagination.Params{}
    params.Parse(page, perPage)  // Parse from query strings
    
    users, total, err := h.userRepo.FindAll(c.Request.Context(), params)
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to fetch users"})
        return
    }
    
    c.JSON(200, gin.H{
        "data": users,
        "total": total,
        "page": params.Page,
        "per_page": params.PerPage,
    })
}
```

---

## Soft Deletes

### Default Behavior (Exclude Soft-Deleted)

```go
// These automatically exclude soft-deleted records (deleted_at IS NULL)
r.Query().Find(ctx, &users)
r.Query().Where("role = ?", "admin").Find(ctx, &users)
r.Query().Count(ctx)
```

### Include Soft-Deleted Records

```go
// Include deleted records
var users []gormUser
r.Query().WithTrashed().Find(ctx, &users)
```

### Only Soft-Deleted Records

```go
// Get only deleted records
var deletedUsers []gormUser
r.Query().OnlyTrashed().Find(ctx, &deletedUsers)

// With pagination
total, err := r.Query().
    OnlyTrashed().
    OrderDesc("deleted_at").
    FindPaginated(ctx, &deletedUsers, params)
```

---

## Advanced Operations

### Aggregation

```go
// Count records
count, err := r.Query().Where("role = ?", "admin").Count(ctx)

// Check existence
exists, err := r.Query().Where("email = ?", "john@example.com").Exists(ctx)

// Pluck (single column)
var emails []string
r.Query().Where("is_active = ?", true).Pluck(ctx, "email", &emails)
```

### Distinct & Grouping

```go
// DISTINCT
var roles []gormRole
r.Query().Distinct("role").Find(ctx, &roles)

// GROUP BY with HAVING
type RoleCount struct {
    Role  string
    Count int
}

r.Query().
    GroupBy("role").
    Having("COUNT(*) > ?", 5).
    Find(ctx, &groups)
```

### Custom Scopes

```go
// Define a scope function
func Active() func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("is_active = ?", true)
    }
}

// Use it with QueryBuilder
r.Query().Scopes(Active()).Find(ctx, &users)
```

---

## Generic BaseRepository

### What is BaseRepository?

`BaseRepository[T]` is a generic struct that provides common CRUD operations for any GORM model type `T`.

### Setup

```go
import "github.com/kodia-studio/kodia/pkg/database"

type UserRepository struct {
    database.BaseRepository[gormUser]  // Embed generic repository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{
        BaseRepository: database.NewBaseRepository[gormUser](db),
    }
}
```

### Available Methods

| Method | Purpose | Returns |
|--------|---------|---------|
| `Query()` | Create a QueryBuilder | `*QueryBuilder` |
| `DB()` | Access underlying GORM instance | `*gorm.DB` |
| `RawCreate()` | Insert a record | `error` |
| `RawFindByID()` | Find by primary key | `(*T, error)` |
| `RawSave()` | Create or update | `error` |
| `RawDelete()` | Soft delete | `error` |
| `RawRestore()` | Un-delete | `error` |
| `RawForceDelete()` | Permanent delete | `error` |
| `RawCount()` | Count active records | `(int64, error)` |
| `RawExists()` | Check if exists | `(bool, error)` |

### Example Usage

```go
func (r *UserRepository) FindAdmins(ctx context.Context) ([]*domain.User, error) {
    var models []gormUser
    if err := r.Query().
        Where("role = ?", "admin").
        Preload("Roles").
        Find(ctx, &models); err != nil {
        return nil, err
    }
    
    // Convert to domain
    users := make([]*domain.User, len(models))
    for i, m := range models {
        users[i] = m.toDomain()
    }
    return users, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
    return r.RawDelete(ctx, id)  // Soft delete
}

func (r *UserRepository) RestoreUser(ctx context.Context, id string) error {
    return r.RawRestore(ctx, id)
}

func (r *UserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
    return r.Query().Where("role = ?", role).Count(ctx)
}
```

---

## Best Practices

### ✅ DO

- **Use QueryBuilder for filtering** — More readable than raw GORM
- **Preload related data** — Prevent N+1 query problems
- **Leverage pagination** — For large result sets
- **Use soft deletes** — For user-facing data
- **Name queries for clarity** — Extract complex queries into methods

```go
// ✅ Good
func (r *UserRepository) FindActiveAdminsByRole(ctx context.Context, role string) ([]gormUser, error) {
    return r.Query().
        Where("is_active = ?", true).
        Where("role = ?", role).
        OrderAsc("name").
        Find(ctx, &users)
}

// ✅ Also good
adminQuery := r.Query().
    Where("is_active = ?", true).
    Where("role = ?", "admin")
```

### ❌ DON'T

- **Don't mix raw GORM with QueryBuilder** — Choose one pattern per repository
- **Don't forget context** — Always pass `ctx` to terminal operations
- **Don't assume no records** — Check error for `gorm.ErrRecordNotFound`
- **Don't build WHERE dynamically with strings** — Use parameterized queries

```go
// ❌ Wrong: string concatenation (SQL injection risk)
role := userInput  // "admin OR 1=1"
query := fmt.Sprintf("WHERE role = '%s'", role)  // VULNERABLE!

// ✅ Right: parameterized
r.Query().Where("role = ?", userInput).Find(ctx, &users)
```

---

## Examples & Patterns

### Pattern 1: Simple Find

```go
user, err := r.FindByID(ctx, "user-123")
if err != nil {
    if errors.Is(err, domain.ErrNotFound) {
        // User doesn't exist
        return nil, nil
    }
    return nil, err
}
```

### Pattern 2: Paginated List with Filters

```go
func (r *UserRepository) FindByStatus(ctx context.Context, status string, params *pagination.Params) ([]*domain.User, int64, error) {
    var models []gormUser
    total, err := r.Query().
        Where("status = ?", status).
        OrderDesc("created_at").
        FindPaginated(ctx, &models, params)
    if err != nil {
        return nil, 0, err
    }
    
    users := make([]*domain.User, len(models))
    for i, m := range models {
        users[i] = m.toDomain()
    }
    return users, total, nil
}
```

### Pattern 3: Complex Query with Joins

```go
func (r *RoleRepository) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
    var roles []gormRole
    err := r.Query().
        Joins("INNER JOIN user_roles ON user_roles.role_id = roles.id").
        Where("user_roles.user_id = ?", userID).
        Find(ctx, &roles)
    if err != nil {
        return nil, err
    }
    
    names := make([]string, len(roles))
    for i, role := range roles {
        names[i] = role.Name
    }
    return names, nil
}
```

### Pattern 4: Soft Delete Workflow

```go
// Delete
err := r.Delete(ctx, userID)

// Check deleted
deletedUsers, total, _ := r.FindTrashed(ctx, params)

// Restore
err := r.Restore(ctx, userID)

// Force delete (permanent)
err := r.ForceDelete(ctx, userID)
```

### Pattern 5: Eager Loading

```go
type gormUser struct {
    ID    string
    Name  string
    Roles []gormRole `gorm:"many2many:user_roles;"`
}

// Load user with all roles
var user gormUser
err := r.Query().
    Preload("Roles").
    Preload("Roles.Permissions").
    Where("id = ?", userID).
    First(ctx, &user)
```

---

## See Also

- [Database Migrations](./DATABASE_MIGRATIONS.md) — Schema management
- [Soft Deletes Guide](./DATABASE_SOFT_DELETES.md) — Soft delete patterns
- [Validation Layer](./VALIDATION_LAYER.md) — Input validation
- [GORM Documentation](https://gorm.io/docs) — Official GORM guide

---

**Last Updated:** April 29, 2026  
**Version:** 1.0  
**Status:** Stable
