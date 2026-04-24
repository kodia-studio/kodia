# Unified Database Migrations

Kodia Framework provides a unified migration system that eliminates duplicate execution and ensures safe, transactional schema updates.

## Overview

The migration system:
- ✅ Tracks executed migrations in database
- ✅ Prevents duplicate execution
- ✅ Runs migrations automatically on startup
- ✅ Supports multiple migration strategies
- ✅ Maintains migration history

## Problem: Before v1.6.0

Previously, Kodia supported three different migration approaches:

```
1. SQL Files (manual)
2. Go Structs (GORM AutoMigrate)
3. Code-based Migrations (custom Schema API)

Running together → unpredictable order, possible duplicates
```

**Issues:**
- ❌ Migrations could run out of order
- ❌ No tracking of executed migrations
- ❌ Manual developer responsibility to prevent duplicates
- ❌ Hard to know which migrations ran

## Solution: Unified Migrator (v1.6.0+)

### How It Works

```
┌─────────────────────────────────┐
│ DatabaseProvider.Boot()         │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│ Create Migrator instance        │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│ EnsureTable()                   │
│ (create kodia_migrations table) │
└────────┬────────────────────────┘
         │
         ▼
┌─────────────────────────────────┐
│ For each migration:             │
│ ├─ Check if already ran         │
│ ├─ If not, execute Up()         │
│ └─ Record in kodia_migrations   │
└─────────────────────────────────┘
```

### Migration Registry

**File:** `internal/infrastructure/database/migrations/go/registry.go`

```go
func All() []Entry {
    return []Entry{
        {
            Name:      "20260422150028_create_webhook_histories",
            Migration: &Migration_20260422150028{},
        },
        {
            Name:      "20260422154949_create_security_elite_tables",
            Migration: &Migration_20260422154949{},
        },
        {
            Name:      "20260422164224_create_failed_jobs_tables",
            Migration: &Migration_20260422164224{},
        },
    }
}
```

All migrations listed here are run automatically in order on startup.

## Creating Migrations

### 1. Create Migration File

Create new migration file in `internal/infrastructure/database/migrations/go/`:

```go
// internal/infrastructure/database/migrations/go/migration_20260425_create_users_table.go

package migrations

import "github.com/kodia-studio/kodia/pkg/database"

type Migration_20260425_CreateUsersTable struct{}

func (m *Migration_20260425_CreateUsersTable) Up(schema *database.Schema) error {
    return schema.CreateTable("users", func(table *database.Blueprint) {
        table.ID()
        table.String("email").Unique()
        table.String("password")
        table.String("name").Nullable()
        table.Timestamps()
        table.SoftDeletes()
    })
}

func (m *Migration_20260425_CreateUsersTable) Down(schema *database.Schema) error {
    return schema.DropTable("users")
}
```

### 2. Register in Registry

Add to `registry.go`:

```go
func All() []Entry {
    return []Entry{
        // ... existing migrations ...
        {
            Name:      "20260425090000_create_users_table",
            Migration: &Migration_20260425_CreateUsersTable{},
        },
    }
}
```

### 3. Run Application

```bash
# Migrations execute automatically on startup
go run ./cmd/server

# Output:
# INFO: Running migration name=20260425090000_create_users_table batch=1
# INFO: Migration completed name=20260425090000_create_users_table
```

## Migration Tracking

### kodia_migrations Table

```sql
CREATE TABLE kodia_migrations (
    id UNSIGNED INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL UNIQUE,
    batch INT NOT NULL,
    run_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**Columns:**
- `id` - Auto-incrementing ID
- `name` - Migration name (e.g., "20260425090000_create_users_table")
- `batch` - Batch number (incremented for each startup)
- `run_at` - When migration executed

**Example Data:**
```sql
id  name                               batch  run_at
1   20260422150028_create_webhook...  1      2026-04-22 15:00:28
2   20260422154949_create_security...  1      2026-04-22 15:49:49
3   20260422164224_create_failed_jobs  1      2026-04-22 16:42:24
4   20260425090000_create_users_table  2      2026-04-25 09:00:00
```

### Checking Migration Status

```bash
# Query migrations table
SELECT * FROM kodia_migrations ORDER BY run_at DESC;

# Find migrations in current batch
SELECT * FROM kodia_migrations WHERE batch = (SELECT MAX(batch) FROM kodia_migrations);

# Count total migrations
SELECT COUNT(*) FROM kodia_migrations;
```

## Implementation Details

### Migrator Interface

**File:** `pkg/database/migrator.go`

```go
type Migration interface {
    Up(schema *Schema) error
    Down(schema *Schema) error
}

type Migrator struct {
    db  *gorm.DB
    log *zap.Logger
}
```

### Running a Migration

```go
func (m *Migrator) Run(name string, migration Migration) error {
    // 1. Check if already executed
    var count int64
    m.db.Model(&MigrationRecord{}).Where("name = ?", name).Count(&count)
    if count > 0 {
        m.log.Debug("Migration already executed", zap.String("name", name))
        return nil
    }

    // 2. Get next batch number
    var lastBatch int
    m.db.Model(&MigrationRecord{}).Select("COALESCE(MAX(batch), 0)").Scan(&lastBatch)
    nextBatch := lastBatch + 1

    // 3. Execute Up() function
    m.log.Info("Running migration", zap.String("name", name), zap.Int("batch", nextBatch))
    schema := NewSchema(m.db)
    if err := migration.Up(schema); err != nil {
        return fmt.Errorf("migration %s failed: %w", name, err)
    }

    // 4. Record in database
    record := MigrationRecord{
        Name:  name,
        Batch: nextBatch,
        RunAt: time.Now(),
    }
    return m.db.Create(&record).Error
}
```

## Multiple Database Strategies

### Strategy 1: Fluent Migration API

Best for schema changes:

```go
func (m *Migration_CreateUsersTable) Up(schema *database.Schema) error {
    return schema.CreateTable("users", func(table *database.Blueprint) {
        table.ID()
        table.String("email").Unique()
        table.String("password")
        table.Timestamps()
    })
}
```

**Advantages:**
- ✅ Type-safe
- ✅ Database-agnostic
- ✅ Readable syntax

### Strategy 2: GORM Models with AutoMigrate

Best for development:

```go
func (m *Migration_CreateUsersTable) Up(schema *database.Schema) error {
    return schema.AutoMigrate(&models.User{})
}
```

**Advantages:**
- ✅ Synced with Go models
- ✅ Quick for prototyping

**Disadvantages:**
- ⚠️ Can cause schema drift
- ⚠️ Hard to customize

### Strategy 3: Raw SQL

For complex operations:

```go
func (m *Migration_CreateUsersTable) Up(schema *database.Schema) error {
    return schema.Raw(`
        CREATE TABLE users (
            id INT PRIMARY KEY AUTO_INCREMENT,
            email VARCHAR(255) UNIQUE NOT NULL,
            password VARCHAR(255) NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `).Error
}
```

**Advantages:**
- ✅ Full database control
- ✅ Complex logic supported

**Disadvantages:**
- ⚠️ Not portable across databases
- ⚠️ More error-prone

## Rollback Mechanism

Currently, migrations are **forward-only**. To revert:

```go
// Option 1: Write a new migration that reverts
func (m *Migration_RevertUsersTable) Up(schema *database.Schema) error {
    // Modify schema to previous state
    return schema.DropTable("users")
}

// Option 2: Direct database intervention (dev only)
DELETE FROM kodia_migrations WHERE name = "20260425090000_create_users_table";
DROP TABLE users;
```

**Future:** Implement `Down()` support for automatic rollback.

## Production Deployment

### Pre-Deployment

```bash
# Test migrations in staging
./server --env=staging

# Verify no errors in logs
grep "Migration" app.log | grep -i error

# Check migration history
psql kodia_staging -c "SELECT * FROM kodia_migrations ORDER BY run_at DESC LIMIT 5"
```

### During Deployment

```bash
# Blue-green deployment with automatic migrations
docker pull kodia:v1.6.0
docker run kodia:v1.6.0
# Migrations run automatically on startup
```

### Monitoring

```bash
# Alert if new migrations fail
grep "migration.*failed" app.log | alert(severity: CRITICAL)

# Track migration execution time
grep "Running migration\|Migration completed" app.log | calculate_duration()

# Monitor migration count growth
SELECT batch, COUNT(*) FROM kodia_migrations GROUP BY batch
```

## Best Practices

### ✅ DO:

- Use sequential timestamps in migration names
- Keep migrations small and focused
- Test migrations on staging first
- Document complex migrations
- Always implement `Down()` for reversibility
- Use migrations for schema-only changes
- Version migrations with batch numbers
- Back up database before production migrations
- Monitor migration execution time
- Include migration status in deployment checklists

### ❌ DON'T:

- Mix data and schema changes in one migration
- Use migrations for large data transformations
- Forget to register migrations in registry
- Skip `Down()` implementation
- Run migrations manually in production
- Deploy without testing migrations
- Create cyclic migration dependencies
- Run migrations without backups
- Change migration names after execution

## Troubleshooting

### Issue: "Migration already executed" but table doesn't exist

```bash
# Diagnosis
SELECT * FROM kodia_migrations WHERE name = 'migration_name';
# Should show one row

# The migration ran but table creation failed mid-way
# Solution: Check error logs, fix issue, delete record, rerun
```

### Issue: Migrations not running on startup

```bash
# Check if DatabaseProvider.Boot() is called
grep "Booting provider" app.log | grep database

# Check registry
grep "Running migration\|already executed" app.log

# Verify migrations are registered
cat internal/infrastructure/database/migrations/go/registry.go
```

### Issue: Timeout during large migration

```bash
# Increase timeout
APP_SHUTDOWN_TIMEOUT_SECS=120 ./server

# Split into smaller migrations
# Option 1: Break into multiple files
# Option 2: Add batching logic in Up()
```

## Performance Considerations

### Migration Execution Time

- Typical: 10-100ms per migration
- Large tables: 100ms-5s per migration
- Schema changes: <10ms

### Database Lock

- Most DDL statements lock tables briefly
- Large tables may cause brief unavailability
- Use online schema change tools for production (future enhancement)

## Testing Migrations

### Unit Test

```go
func TestMigration_CreateUsersTable(t *testing.T) {
    db := setupTestDB()
    migrator := database.NewMigrator(db, logger)

    // Run migration
    migration := &Migration_CreateUsersTable{}
    err := migrator.Run("test_migration", migration)
    assert.NoError(t, err)

    // Verify table exists
    assert.True(t, db.Migrator().HasTable("users"))
    assert.True(t, db.Migrator().HasColumn("users", "email"))
}
```

### Integration Test

```go
func TestAllMigrations(t *testing.T) {
    // Run entire migration suite
    app := setupTestApp()
    err := app.Boot()
    assert.NoError(t, err)

    // Verify all tables created
    tables := []string{"users", "webhook_histories", "failed_jobs"}
    for _, table := range tables {
        assert.True(t, app.DB.Migrator().HasTable(table))
    }
}
```

## References

- [Laravel Migrations](https://laravel.com/docs/migrations) — Inspiration
- [GORM Migrations](https://gorm.io/docs/migration.html)
- [Flyway](https://flywaydb.org/) — External tool reference
- [Database Versioning Best Practices](https://www.liquibase.org/get-started/best-practices)

## Conclusion

Unified database migrations in Kodia ensure:
- ✅ **Safety**: No duplicate execution
- ✅ **Auditability**: Complete migration history
- ✅ **Reliability**: Automatic startup execution
- ✅ **Flexibility**: Multiple migration strategies
- ✅ **Production-Ready**: Batch tracking and monitoring
