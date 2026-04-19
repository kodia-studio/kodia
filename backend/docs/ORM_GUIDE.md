# 🗂️ Advanced ORM Guide - Kodia Framework

Complete guide to using GORM and advanced database patterns in Kodia Framework.

**Table of Contents:**
- [Base Models & Hooks](#base-models--hooks)
- [Query Scopes](#query-scopes)
- [Transactions](#transactions)
- [Pagination with Sorting & Filtering](#pagination-with-sorting--filtering)
- [Anti-Corruption Pattern](#anti-corruption-pattern)
- [Soft Deletes](#soft-deletes)
- [Seeders](#seeders)
- [Best Practices](#best-practices)
- [Common Patterns](#common-patterns)
- [Troubleshooting](#troubleshooting)

---

## Base Models & Hooks

### Using BaseModel

The `BaseModel` struct provides automatic UUID generation and timestamp management:

```go
import "github.com/kodia-studio/kodia/pkg/database"

// Your domain entity
type Article struct {
    database.BaseModel
    Title   string
    Content string
    Author  string
}

// Usage: GORM automatically generates UUID on creation
article := &Article{
    Title:   "My Article",
    Content: "Article content...",
    Author:  "John Doe",
}
db.Create(article)  // ID is auto-generated, CreatedAt/UpdatedAt set automatically
```

### Using SoftDeleteModel

For entities that need soft-delete support:

```go
type Post struct {
    database.SoftDeleteModel
    Title string
    Body  string
}

// Soft delete marks deleted_at timestamp
db.Delete(&post)  // Sets deleted_at to current time

// Query excludes soft-deleted records by default
var posts []Post
db.Find(&posts)  // Returns only posts where deleted_at IS NULL

// Include soft-deleted records
db.Unscoped().Find(&posts)  // Returns all posts including soft-deleted ones

// Find only soft-deleted records
db.Where("deleted_at IS NOT NULL").Find(&posts)
```

### GORM Hooks

The `BaseModel` uses a `BeforeSave` hook to generate UUIDs. You can add additional hooks to your models:

```go
// BeforeCreate hook
func (p *Post) BeforeCreate(tx *gorm.DB) error {
    if p.ID == "" {
        p.ID = uuid.New().String()
    }
    return nil
}

// AfterUpdate hook — update cache after updates
func (u *User) AfterUpdate(tx *gorm.DB) error {
    cache.Invalidate("user:" + u.ID)
    return nil
}

// BeforeDelete hook — backup before deletion
func (o *Order) BeforeDelete(tx *gorm.DB) error {
    return backupService.BackupOrder(o)
}
```

**Available hooks:**
- `BeforeCreate`, `AfterCreate`
- `BeforeUpdate`, `AfterUpdate`
- `BeforeSave`, `AfterSave`
- `BeforeDelete`, `AfterDelete`
- `BeforeFind`, `AfterFind`

---

## Query Scopes

Scopes are reusable query fragments that chain together:

```go
import "github.com/kodia-studio/kodia/pkg/database"
import "github.com/kodia-studio/kodia/pkg/pagination"
```

### Pagination Scope

```go
params := pagination.FromContext(c)  // Parse ?page=1&per_page=15&sort=name&sort_dir=asc&search=query

var users []User
var total int64

db.Model(&User{}).Count(&total)
db.Scopes(database.Paginate(params)).Find(&users)

// Respons dengan metadata
c.JSON(200, gin.H{
    "data": users,
    "meta": gin.H{
        "page": params.Page,
        "per_page": params.PerPage,
        "total": total,
        "total_pages": params.TotalPages(total),
    },
})
```

### SortBy Scope

Sort by whitelisted fields only (prevents SQL injection):

```go
// Whitelist allowed sort fields
allowedFields := []string{"name", "email", "created_at"}

var users []User
db.Scopes(
    database.SortBy(params.Sort, params.SortDir, allowedFields),
    database.Paginate(params),
).Find(&users)

// Example: ?sort=name&sort_dir=desc
// Generates: ORDER BY name DESC LIMIT 15 OFFSET 0
```

### Active Scope

Filter to active records only:

```go
var activeUsers []User
db.Scopes(database.Active()).Find(&activeUsers)
// Generates: SELECT * FROM users WHERE is_active = true
```

### Search Scope

Full-text-like search on a single column:

```go
var results []Product
db.Scopes(
    database.Search("name", "laptop"),
    database.Paginate(params),
).Find(&results)
// Generates: SELECT * FROM products WHERE name ILIKE '%laptop%' LIMIT 15
```

### Soft Delete Scopes

```go
var deletedPosts []Post
db.Scopes(database.OnlyTrashed()).Find(&deletedPosts)
// Returns only soft-deleted records

var allPosts []Post
db.Scopes(database.WithTrashed()).Find(&allPosts)
// Returns all records including soft-deleted ones
```

### Chaining Scopes

```go
var results []User
db.Scopes(
    database.Active(),
    database.SortBy("created_at", "desc", []string{"created_at", "name"}),
    database.Paginate(params),
).Find(&results)
```

---

## Transactions

Wrap multiple operations in a transaction for ACID guarantees:

### Basic Transaction

```go
import "github.com/kodia-studio/kodia/pkg/database"

err := database.WithTransaction(db, func(tx *gorm.DB) error {
    // Create user
    user := &User{Name: "John", Email: "john@example.com"}
    if err := tx.Create(user).Error; err != nil {
        return err  // Rolls back automatically
    }

    // Create welcome email record
    email := &Email{
        UserID: user.ID,
        Type:   "welcome",
        Status: "pending",
    }
    if err := tx.Create(email).Error; err != nil {
        return err  // Both Create() calls are rolled back
    }

    return nil  // Commits both operations
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

### Transaction with Isolation Level

```go
err := database.WithTransactionIsolation(db, gorm.LevelSerializable, func(tx *gorm.DB) error {
    // Critical business logic here
    var balance float64
    if err := tx.Model(&Account{}).Where("id = ?", accountID).
        Select("balance").Row().Scan(&balance); err != nil {
        return err
    }

    if balance < amount {
        return errors.New("insufficient balance")
    }

    return tx.Model(&Account{}).Where("id = ?", accountID).
        Update("balance", gorm.Expr("balance - ?", amount)).Error
})
```

### Safe Operations Helpers

```go
// SafeSave wraps Save in a transaction
err := database.SafeSave(db, &user)

// SafeDelete wraps Delete in a transaction
err := database.SafeDelete(db, &user)
```

---

## Pagination with Sorting & Filtering

### The Params Struct

The enhanced `Params` struct now includes sorting and searching:

```go
type Params struct {
    Page    int    // Current page (1-indexed)
    PerPage int    // Items per page (1-100)
    Sort    string // Field to sort by (e.g., "name")
    SortDir string // "asc" or "desc"
    Search  string // Free-text search query
}
```

### Parsing from Request

```go
// In your handler
params := pagination.FromContext(c)
// Parses: ?page=2&per_page=20&sort=created_at&sort_dir=desc&search=active

var products []Product
var total int64

db.Model(&Product{}).Count(&total)

db.Scopes(
    database.Search("name", params.Search),
    database.SortBy(params.Sort, params.SortDir, []string{"name", "price", "created_at"}),
    database.Paginate(params),
).Find(&products)

c.JSON(200, gin.H{
    "data": products,
    "meta": gin.H{
        "page": params.Page,
        "per_page": params.PerPage,
        "total": total,
        "total_pages": params.TotalPages(total),
    },
})
```

### Backward Compatibility

The original pagination still works:

```go
page := c.DefaultQuery("page", "1")
pageNum, _ := strconv.Atoi(page)
params := pagination.FromContext(c)
// If ?sort is not provided, params.Sort is ""
// If ?search is not provided, params.Search is ""
```

---

## Anti-Corruption Pattern

Separate GORM models from domain entities to keep domain layer clean:

### Example: User Repository

Domain entity (in `core/domain`):
```go
type User struct {
    ID        string
    Email     string
    Password  string
    Role      UserRole
    IsActive  bool
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}
```

GORM model (in `adapters/repository/postgres`):
```go
type gormUser struct {
    ID        string     `gorm:"column:id;primaryKey"`
    Email     string     `gorm:"column:email;uniqueIndex;not null"`
    Password  string     `gorm:"column:password;not null"`
    Role      string     `gorm:"column:role;not null;default:'user'"`
    IsActive  bool       `gorm:"column:is_active;not null;default:true"`
    CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt *time.Time `gorm:"column:deleted_at;index"`
}

func (gormUser) TableName() string { return "users" }

func (g *gormUser) toDomain() *domain.User {
    return &domain.User{
        ID:        g.ID,
        Email:     g.Email,
        Password:  g.Password,
        Role:      domain.UserRole(g.Role),
        IsActive:  g.IsActive,
        CreatedAt: g.CreatedAt,
        UpdatedAt: g.UpdatedAt,
        DeletedAt: g.DeletedAt,
    }
}

func fromDomainUser(u *domain.User) *gormUser {
    return &gormUser{
        ID:       u.ID,
        Email:    u.Email,
        Password: u.Password,
        Role:     string(u.Role),
        // ... etc
    }
}
```

Repository methods use the conversion:
```go
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
    m := fromDomainUser(user)
    return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
    var m gormUser
    if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, domain.ErrNotFound
        }
        return nil, err
    }
    return m.toDomain(), nil
}
```

**Benefits:**
- Domain logic stays free of ORM concerns
- Easier to test (mock domain objects, not GORM models)
- Easier to switch ORMs or databases later
- GORM tags don't pollute domain structs

---

## Soft Deletes

### Using SoftDeleteModel

The `SoftDeleteModel` embeds GORM's native `gorm.DeletedAt`:

```go
type Post struct {
    database.SoftDeleteModel
    Title   string
    Content string
}

// Soft delete
db.Delete(&post)  // Sets deleted_at

// Query excludes soft-deleted by default
var posts []Post
db.Find(&posts)  // Only non-deleted

// Include soft-deleted
db.Unscoped().Find(&posts)  // All records

// Only soft-deleted
db.Where("deleted_at IS NOT NULL").Find(&posts)

// Permanently delete (bypass soft-delete)
db.Unscoped().Delete(&post)  // Actually removes from database

// Restore soft-deleted record
db.Model(&post).Update("deleted_at", nil)
```

---

## Seeders

### Running Seeders

```bash
cd backend
go run cmd/seeder/main.go
```

Executes all registered seeders in `seeders/registry.go`:
- UserSeeder — creates 20 fake users
- ProductSeeder — creates 50 fake products

### Creating a Custom Seeder

```go
package seeders

import (
    "github.com/brianvoe/gofakeit/v7"
    "github.com/google/uuid"
    "github.com/kodia-studio/kodia/internal/core/domain"
    "gorm.io/gorm"
)

type BlogPostSeeder struct{}

func (s *BlogPostSeeder) Run(db *gorm.DB) error {
    posts := make([]interface{}, 30)

    for i := 0; i < 30; i++ {
        post := &domain.BlogPost{
            ID:    uuid.New().String(),
            Title: gofakeit.Sentence(5),
            Slug:  slugify(gofakeit.Sentence(5)),
            Body:  gofakeit.Paragraph(3),
        }
        posts[i] = post
    }

    return db.CreateInBatches(posts, 10).Error
}
```

### Register in Registry

```go
// seeders/registry.go
var Registry = []Seeder{
    &UserSeeder{},
    &ProductSeeder{},
    &BlogPostSeeder{},  // Add here
}
```

---

## Best Practices

### ✅ DO:

```go
// 1. Use BaseModel for automatic ID + timestamps
type Article struct {
    database.BaseModel
    Title string
}

// 2. Use scopes for reusable query fragments
db.Scopes(
    database.Active(),
    database.Paginate(params),
).Find(&items)

// 3. Use anti-corruption pattern for GORM models
type gormUser struct { /* GORM tags */ }
func (g *gormUser) toDomain() *domain.User { /* conversion */ }

// 4. Wrap critical operations in transactions
database.WithTransaction(db, func(tx *gorm.DB) error {
    // Multiple operations here
})

// 5. Use context throughout for timeouts/cancellation
db.WithContext(ctx).Create(&item)

// 6. Whitelist sortable fields
database.SortBy(params.Sort, params.SortDir, []string{"name", "email"})

// 7. Validate pagination parameters
if params.PerPage < 1 || params.PerPage > 100 {
    params.PerPage = 15
}
```

### ❌ DON'T:

```go
// ❌ Hardcode UUIDs — use BaseModel hook instead
user.ID = uuid.New().String()

// ❌ Skip transactions for multi-step operations
db.Create(&user)
db.Create(&welcome_email)  // Fails, but user already created

// ❌ Allow unsanitized sorting
db.Order(params.Sort).Find(&items)  // SQL injection risk

// ❌ Expose GORM models to domain layer
func (s *UserService) GetUser(id string) *gormUser {  // Wrong!
    // Should return domain.User, not gormUser
}

// ❌ Forget soft-delete when querying
db.Find(&items)  // May miss soft-deleted items in some contexts

// ❌ Use raw SQL without parameter binding
db.Where("email = " + email).First(&user)  // SQL injection!
// Use:
db.Where("email = ?", email).First(&user)
```

---

## Common Patterns

### Search with Multiple Fields

```go
var users []User
searchTerm := "john"
db.Where("name ILIKE ? OR email ILIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%").
    Find(&users)
```

### Batch Operations

```go
var items []Item

// Batch create with safety
db.CreateInBatches(items, 100)  // Insert in batches of 100

// Batch update
db.Model(&Item{}).Where("status = ?", "pending").
    Update("status", "processing")

// Batch delete
db.Where("created_at < ?", time.Now().AddDate(0, -1, 0)).
    Delete(&Item{})
```

### Count and Pagination

```go
var total int64
var items []Item

// Get total count
db.Model(&Item{}).Where("active = ?", true).Count(&total)

// Fetch paginated results
db.Where("active = ?", true).
    Offset((page - 1) * pageSize).
    Limit(pageSize).
    Find(&items)

return items, total, nil
```

### Exists Check

```go
var exists bool
db.Model(&User{}).
    Select("count(*) > 0").
    Where("email = ?", email).
    Find(&exists)

if exists {
    return errors.New("email already exists")
}
```

---

## Troubleshooting

### Issue: UUID not generated on create

**Cause:** BeforeSave hook not called
**Solution:** Use BaseModel or ensure BeforeSave hook is defined

```go
type MyEntity struct {
    database.BaseModel  // ✅ Now UUID auto-generates
}
```

### Issue: Soft-deleted records showing in queries

**Cause:** GORM soft-delete auto-filtering requires `gorm.DeletedAt`
**Solution:** Use SoftDeleteModel instead of manual time.Time field

```go
type Entity struct {
    database.SoftDeleteModel  // ✅ Auto-filters soft-deleted
}
```

### Issue: "record not found" but record exists

**Cause:** Likely soft-deleted
**Solution:** Use `Unscoped()` to include soft-deleted records

```go
db.Unscoped().Where("id = ?", id).First(&item)
```

### Issue: Sorting not working / SQL injection concerns

**Cause:** Unsanitized sort parameter
**Solution:** Use SortBy scope with whitelist

```go
database.SortBy(params.Sort, params.SortDir, []string{"name", "created_at"})
```

### Issue: Changes not persisted

**Cause:** Forgot to handle transaction or check error
**Solution:** Always check returned errors and use transactions

```go
if err := db.Save(&item).Error; err != nil {
    return err
}
```

---

## Cheat Sheet

| Task | Code |
|------|------|
| Create with auto ID | `db.Create(&item)` (if using BaseModel) |
| Find by ID | `db.Where("id = ?", id).First(&item)` |
| Find all paginated | `db.Scopes(database.Paginate(params)).Find(&items)` |
| Update | `db.Save(&item)` or `db.Model(&item).Update("field", value)` |
| Soft delete | `db.Delete(&item)` |
| Restore | `db.Model(&item).Update("deleted_at", nil)` |
| Hard delete | `db.Unscoped().Delete(&item)` |
| Count | `db.Model(&Item{}).Count(&total)` |
| Exists | Use Select("count(*) > 0") pattern |
| Transaction | `database.WithTransaction(db, func(tx *gorm.DB) error { ... })` |
| Batch create | `db.CreateInBatches(items, 100)` |
| Scope chain | `db.Scopes(Active(), Paginate(p), SortBy(...)).Find(&items)` |

---

**Next:** Configure your models using BaseModel or SoftDeleteModel, then use scopes for queries!
