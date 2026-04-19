package database

import (
	"fmt"

	"github.com/kodia-studio/kodia/pkg/pagination"
	"gorm.io/gorm"
)

// Paginate is a GORM scope for pagination.
// Usage: db.Scopes(Paginate(params)).Find(&items)
func Paginate(p *pagination.Params) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(p.Offset()).Limit(p.Limit())
	}
}

// SortBy is a GORM scope for dynamic sorting with injection prevention.
// Only fields in the allowed list can be sorted to prevent SQL injection.
// Usage: db.Scopes(SortBy("name", "asc", []string{"name", "email", "created_at"})).Find(&items)
func SortBy(field, dir string, allowed []string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Validate field against whitelist
		allowed_map := make(map[string]bool)
		for _, f := range allowed {
			allowed_map[f] = true
		}

		if !allowed_map[field] {
			return db // Silently ignore invalid sort field
		}

		// Validate direction
		if dir != "asc" && dir != "desc" {
			dir = "asc"
		}

		return db.Order(fmt.Sprintf("%s %s", field, dir))
	}
}

// Active is a GORM scope that filters to active records only.
// Usage: db.Scopes(Active()).Find(&users)
func Active() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_active = ?", true)
	}
}

// Search is a GORM scope for full-text-like search on a single column.
// Uses ILIKE for PostgreSQL case-insensitive search.
// Usage: db.Scopes(Search("name", "john")).Find(&users)
func Search(col, term string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if term == "" {
			return db
		}
		// ILIKE for PostgreSQL (case-insensitive), LIKE for MySQL
		// Parameter binding prevents SQL injection
		return db.Where(fmt.Sprintf("%s ILIKE ?", col), "%"+term+"%")
	}
}

// WithTrashed is a GORM scope that includes soft-deleted records in the query.
// Usage: db.Scopes(WithTrashed()).Find(&items)
func WithTrashed() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Unscoped()
	}
}

// OnlyTrashed is a GORM scope that returns only soft-deleted records.
// Usage: db.Scopes(OnlyTrashed()).Find(&items)
func OnlyTrashed() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NOT NULL")
	}
}
