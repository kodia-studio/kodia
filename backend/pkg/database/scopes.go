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

// Paginator holds pagination metadata and data.
type Paginator struct {
	Total       int64       `json:"total"`
	PerPage     int         `json:"per_page"`
	CurrentPage int         `json:"current_page"`
	LastPage    int         `json:"last_page"`
	From        int         `json:"from"`
	To          int         `json:"to"`
	Data        interface{} `json:"data"`
}

// PaginateFluent performs pagination and returns metadata.
// Example: paginator, err := database.PaginateFluent(db, params, &users)
func PaginateFluent(db *gorm.DB, p *pagination.Params, dest interface{}) (*Paginator, error) {
	var total int64
	db.Model(dest).Count(&total)

	err := db.Offset(p.Offset()).Limit(p.Limit()).Find(dest).Error
	if err != nil {
		return nil, err
	}

	lastPage := int((total + int64(p.Limit()) - 1) / int64(p.Limit()))
	from := p.Offset() + 1
	to := p.Offset() + p.Limit()
	if int64(to) > total {
		to = int(total)
	}

	return &Paginator{
		Total:       total,
		PerPage:     p.Limit(),
		CurrentPage: p.Page,
		LastPage:    lastPage,
		From:        from,
		To:          to,
		Data:        dest,
	}, nil
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
