# 🗑️ Soft Deletes Guide

**Soft deletion** is a data safety pattern that marks records as deleted without removing them from the database. This guide covers implementing and using soft deletes in Kodia Framework.

**Table of Contents:**
- [What is Soft Delete?](#what-is-soft-delete)
- [Why Use Soft Deletes?](#why-use-soft-deletes)
- [Quick Start](#quick-start)
- [How It Works](#how-it-works)
- [Using Soft Deletes](#using-soft-deletes)
- [Advanced Operations](#advanced-operations)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

---

## What is Soft Delete?

Soft deletion marks a record as deleted by setting a `deleted_at` timestamp, rather than physically removing it from the database. The record remains stored but is automatically excluded from all normal queries.

### Comparison

| Operation | Hard Delete | Soft Delete |
|-----------|------------|------------|
| **Action** | Removes row from DB | Sets `deleted_at` timestamp |
| **Recovery** | ❌ Data is lost | ✅ Can restore anytime |
| **Query Impact** | Data appears gone | Requires explicit filtering |
| **Audit Trail** | No history | Full history preserved |
| **Use Case** | Permanent removal | Safety/reversibility |

### Example: Users Table

```sql
-- Before soft delete
SELECT * FROM users WHERE id = '123';
-- Output: { id: 123, name: "John", deleted_at: NULL }

-- After soft delete (delete called)
SELECT * FROM users WHERE id = '123';
-- Output: (empty - automatically filtered)

-- Access soft-deleted records
SELECT * FROM users WHERE deleted_at IS NOT NULL;
-- Output: { id: 123, name: "John", deleted_at: "2026-04-29 10:30:00" }

-- Restore
UPDATE users SET deleted_at = NULL WHERE id = '123';
-- Now visible in normal queries again
```

---

## Why Use Soft Deletes?

### ✅ Data Safety
- **Accidental Deletion Recovery** — Restore deleted data without database backup
- **Audit Compliance** — Maintain complete record history for compliance requirements
- **Undo/Restore UX** — Offer users a "restore deleted item" feature

### ✅ Business Logic
- **Analytics** — Track which customers churned without losing their data
- **Referential Integrity** — Keep relationships intact (user exists but is inactive)
- **Business Intelligence** — Historical analysis without data loss

### ✅ Operational
- **Zero Downtime** — No migration required; just flag as deleted
- **Reversible** — Undo mistakes immediately
- **Compliance** — Meet GDPR/data retention requirements (can restore if user requests)

---

## Quick Start

### Enable Soft Deletes on a Model

Soft deletes are built into Kodia. Any entity using GORM automatically supports them.

#### Step 1: Use GORM's DeletedAt in Your Model

```go
package postgres

import (
    "time"
    "gorm.io/gorm"
)

// gormUser is a GORM model with soft delete support
type gormUser struct {
    ID        string         `gorm:"column:id;primaryKey"`
    Name      string         `gorm:"column:name"`
    Email     string         `gorm:"column:email"`
    CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`  // ← Soft delete field
}

func (gormUser) TableName() string { return "users" }
```

#### Step 2: Database Schema

Ensure your migration includes the `deleted_at` column:

```go
// internal/infrastructure/database/migrations/go/...
func (m *MyMigration) Up(schema *database.Schema) error {
    return schema.Create("users", func(table *database.Blueprint) {
        table.ID()
        table.String("name")
        table.String("email").Unique()
        table.SoftDeletes()  // ← Adds deleted_at timestamp + index
        table.Timestamps()
    })
}
```

#### Step 3: Use in Repository

```go
func (r *UserRepository) Delete(ctx context.Context, id string) error {
    // GORM automatically sets deleted_at to current time
    return r.db.WithContext(ctx).Where("id = ?", id).Delete(&gormUser{}).Error
}
```

---

## How It Works

### Automatic Filtering

GORM automatically filters out soft-deleted records in all queries:

```go
// These queries exclude soft-deleted records automatically
db.Find(&users)                           // ✅ Excludes deleted
db.Where("role = ?", "admin").Find(&...)  // ✅ Excludes deleted
db.First(&user)                           // ✅ Excludes deleted
db.Count(&count)                          // ✅ Excludes deleted
```

### Type Conversion

Kodia provides helper functions to convert between domain and GORM layers:

```go
// Domain layer uses *time.Time
type User struct {
    ID        string
    DeletedAt *time.Time  // nil = not deleted, value = deleted at this time
}

// GORM layer uses gorm.DeletedAt
type gormUser struct {
    DeletedAt gorm.DeletedAt  // Built-in type with Valid flag
}

// Conversion helpers
database.ToDeletedAt(u.DeletedAt)        // *time.Time → gorm.DeletedAt
database.FromDeletedAt(g.DeletedAt)      // gorm.DeletedAt → *time.Time
```

---

## Using Soft Deletes

### Basic Operations

#### Delete (Soft)

```go
// Soft delete - sets deleted_at timestamp
err := userRepo.Delete(ctx, userID)
if err != nil {
    log.Error("failed to delete user", err)
}

// User still exists in DB but is hidden from normal queries
```

#### Restore

```go
// Un-delete - clears the deleted_at timestamp
err := userRepo.Restore(ctx, userID)
if err != nil {
    log.Error("failed to restore user", err)
}

// User is now visible in normal queries again
```

#### Force Delete (Permanent)

```go
// Permanently remove from database (cannot be undone!)
err := userRepo.ForceDelete(ctx, userID)
if err != nil {
    log.Error("failed to permanently delete user", err)
}
```

### Querying Soft-Deleted Records

#### Exclude Soft-Deleted (Default Behavior)

```go
// All these automatically exclude soft-deleted records
users, total, err := userRepo.FindAll(ctx, params)
user, err := userRepo.FindByID(ctx, id)
exists, err := userRepo.ExistsByEmail(ctx, email)
count, err := userRepo.CountByRole(ctx, role)
```

#### Include Soft-Deleted Records

```go
// Use WithTrashed scope to include deleted records
db.Scopes(database.WithTrashed()).Find(&users)

// Now queries will return ALL records (including deleted)
```

#### Fetch Only Soft-Deleted Records

```go
// Use OnlyTrashed scope to get ONLY deleted records
db.Scopes(database.OnlyTrashed()).Find(&deletedUsers)

// Or use repository method
deletedUsers, total, err := userRepo.FindTrashed(ctx, params)
```

### Example: Complete Workflow

```go
// 1. Create user
user := &domain.User{Name: "John", Email: "john@example.com"}
err := userRepo.Create(ctx, user)

// 2. User is visible
foundUser, _ := userRepo.FindByID(ctx, user.ID)
fmt.Println(foundUser.Name)  // Output: "John"

// 3. Soft delete user
userRepo.Delete(ctx, user.ID)

// 4. User is hidden from normal queries
foundUser, err := userRepo.FindByID(ctx, user.ID)
// err == domain.ErrNotFound  (because it's soft-deleted)

// 5. Admin restores user
userRepo.Restore(ctx, user.ID)

// 6. User is visible again
foundUser, _ := userRepo.FindByID(ctx, user.ID)
fmt.Println(foundUser.Name)  // Output: "John"
```

---

## Advanced Operations

### Check if Record is Deleted

```go
user, _ := userRepo.FindByID(ctx, id)  // With WithTrashed scope

// Using the IsDeleted() method
if user.IsDeleted() {
    fmt.Println("User has been deleted")
    fmt.Printf("Deleted at: %v\n", user.DeletedAt)
}
```

### Batch Restore Deleted Records

```go
// Restore all users deleted in the last 7 days
sevenDaysAgo := time.Now().AddDate(0, 0, -7)

db.Unscoped().
    Model(&gormUser{}).
    Where("deleted_at > ?", sevenDaysAgo).
    Update("deleted_at", nil)
```

### Permanent Cleanup (Hard Delete Old Records)

```go
// Permanently delete records that were soft-deleted over 1 year ago
oneYearAgo := time.Now().AddDate(-1, 0, 0)

db.Unscoped().
    Where("deleted_at < ?", oneYearAgo).
    Delete(&gormUser{})
```

### Count Deleted Records

```go
var count int64
db.Unscoped().Model(&gormUser{}).Where("deleted_at IS NOT NULL").Count(&count)
fmt.Printf("Total deleted users: %d\n", count)
```

### List Deleted Records with Details

```go
var deletedUsers []gormUser
db.Scopes(database.OnlyTrashed()).
    Order("deleted_at DESC").
    Find(&deletedUsers)

for _, user := range deletedUsers {
    fmt.Printf("%s deleted at %v\n", user.Name, user.DeletedAt.Time)
}
```

---

## Best Practices

### ✅ DO

- **Use soft deletes by default** for user-facing entities (users, posts, files)
- **Provide restore functionality** in your UI for better UX
- **Log deletions** for audit trails
- **Implement cleanup jobs** to permanently delete old soft-deleted records
- **Use meaningful timestamps** to understand when records were deleted
- **Test restore paths** thoroughly

### ❌ DON'T

- **Force delete without confirmation** — Always warn users that force delete is permanent
- **Expose DeletedAt to API responses** — Filter it out in DTOs (or only show to admins)
- **Mix hard and soft deletes** — Choose one pattern per entity
- **Forget to migrate old data** — Ensure `deleted_at` column is added to existing tables
- **Rely solely on soft deletes** — Combine with regular backups for true data safety
- **Query with AND deleted_at IS NULL** — GORM does this automatically; you'll double-filter

### Consistent Naming

| Operation | Method | Behavior |
|-----------|--------|----------|
| **Soft Delete** | `Delete(ctx, id)` | Sets `deleted_at` timestamp |
| **Restore** | `Restore(ctx, id)` | Clears `deleted_at` (sets to NULL) |
| **Force Delete** | `ForceDelete(ctx, id)` | Permanent removal (cannot undo) |

---

## Troubleshooting

### "Record Not Found" When It Should Exist

**Problem:** Record was soft-deleted.

```go
user, err := userRepo.FindByID(ctx, id)
// err == domain.ErrNotFound
```

**Solution:** Use `WithTrashed()` scope if you need to access soft-deleted records:

```go
var user gormUser
db.Scopes(database.WithTrashed()).First(&user, "id = ?", id)
```

### OnlyTrashed Returns Empty Results

**Problem:** Called without `Unscoped()`.

```go
// ❌ Wrong - GORM's automatic filter blocks results
db.Where("deleted_at IS NOT NULL").Find(&users)  // Empty!

// ✅ Correct
db.Scopes(database.OnlyTrashed()).Find(&users)  // Works!
```

### Soft Delete Not Working

**Problem:** Model uses `*time.Time` instead of `gorm.DeletedAt`.

```go
// ❌ Wrong - manual handling required
type User struct {
    DeletedAt *time.Time
}

// ✅ Correct
type User struct {
    DeletedAt gorm.DeletedAt
}
```

### Deleted Records Appear in Count

**Problem:** Counted with `Unscoped()`.

```go
// ❌ Includes deleted records
db.Unscoped().Model(&User{}).Count(&total)

// ✅ Excludes deleted records (default)
db.Model(&User{}).Count(&total)
```

---

## API Examples

### REST Endpoint Pattern

```go
// Soft delete
DELETE /api/users/:id
Response: 204 No Content

// Restore deleted user
POST /api/users/:id/restore
Response: 200 { user object }

// List deleted users (admin only)
GET /api/admin/users?deleted=only
Response: 200 [ { deleted users } ]

// Permanently delete (admin only, requires confirmation)
DELETE /api/admin/users/:id/force
Body: { "confirm": true }
Response: 204 No Content
```

### Example Handler

```go
// POST /api/users/:id/restore
func (h *UserHandler) Restore(c *gin.Context) {
    id := c.Param("id")
    
    if err := h.userRepo.Restore(c.Request.Context(), id); err != nil {
        c.JSON(500, ErrorResponse{Error: "Failed to restore user"})
        return
    }
    
    user, _ := h.userRepo.FindByID(c.Request.Context(), id)
    c.JSON(200, user)
}
```

---

## Laravel Equivalents

If you're migrating from Laravel, here's the mapping:

| Laravel | Kodia Go |
|---------|----------|
| `$user->delete()` | `repo.Delete(ctx, id)` |
| `$user->restore()` | `repo.Restore(ctx, id)` |
| `$user->forceDelete()` | `repo.ForceDelete(ctx, id)` |
| `$user->trashed()` | `user.IsDeleted()` |
| `User::withTrashed()->get()` | `db.Scopes(database.WithTrashed()).Find(...)` |
| `User::onlyTrashed()->get()` | `db.Scopes(database.OnlyTrashed()).Find(...)` |
| `User::restore()` | `db.Unscoped().Model(...).Update("deleted_at", nil)` |
| `$table->softDeletes()` | `table.SoftDeletes()` |

---

## See Also

- [Backend Guide](BACKEND_GUIDE.md) — General database patterns
- [Database Migrations](../backend/docs/DATABASE_MIGRATIONS.md) — Schema management
- [GORM Documentation](https://gorm.io/docs/delete.html) — Official GORM soft delete guide
- [Architecture Guide](ARCHITECTURE.md) — Repository pattern details

---

**Last updated:** April 29, 2026
