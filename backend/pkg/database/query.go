package database

import (
	"context"
	"strings"

	"github.com/kodia-studio/kodia/pkg/pagination"
	"gorm.io/gorm"
)

// QueryBuilder is a fluent, chainable wrapper around *gorm.DB.
// Each method returns a new QueryBuilder — the original is never mutated.
//
// Example:
//   users, err := db.Query().
//       Where("role = ?", "admin").
//       WhereIn("status", []string{"active", "pending"}).
//       Preload("Roles").
//       OrderDesc("created_at").
//       Find(ctx, &results)
type QueryBuilder struct {
	db *gorm.DB
}

// NewQuery creates a QueryBuilder from a *gorm.DB instance.
func NewQuery(db *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: db}
}

// --- Filtering ---

// Where adds a WHERE condition.
//   q.Where("role = ?", "admin")
func (q *QueryBuilder) Where(query interface{}, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Where(query, args...)}
}

// WhereIn adds a WHERE field IN (...) clause.
//   q.WhereIn("status", []string{"active", "pending"})
func (q *QueryBuilder) WhereIn(field string, values interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Where(field+" IN ?", values)}
}

// WhereNotIn adds a WHERE field NOT IN (...) clause.
func (q *QueryBuilder) WhereNotIn(field string, values interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Where(field+" NOT IN ?", values)}
}

// WhereBetween adds a WHERE field BETWEEN from AND to clause.
func (q *QueryBuilder) WhereBetween(field string, from, to interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Where(field+" BETWEEN ? AND ?", from, to)}
}

// WhereNull adds a WHERE field IS NULL clause.
func (q *QueryBuilder) WhereNull(field string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Where(field + " IS NULL")}
}

// WhereNotNull adds a WHERE field IS NOT NULL clause.
func (q *QueryBuilder) WhereNotNull(field string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Where(field + " IS NOT NULL")}
}

// OrWhere adds an OR WHERE condition.
func (q *QueryBuilder) OrWhere(query interface{}, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Or(query, args...)}
}

// Not adds a NOT condition.
func (q *QueryBuilder) Not(query interface{}, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Not(query, args...)}
}

// --- Selection & Ordering ---

// Select specifies which columns to retrieve.
//   q.Select("id", "name", "email")
func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Select(fields)}
}

// OrderAsc orders results by field ascending.
func (q *QueryBuilder) OrderAsc(field string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Order(field + " ASC")}
}

// OrderDesc orders results by field descending.
func (q *QueryBuilder) OrderDesc(field string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Order(field + " DESC")}
}

// OrderRaw sets a raw ORDER BY expression.
func (q *QueryBuilder) OrderRaw(expr string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Order(expr)}
}

// Limit sets the maximum number of rows to return.
func (q *QueryBuilder) Limit(n int) *QueryBuilder {
	return &QueryBuilder{db: q.db.Limit(n)}
}

// Offset sets the number of rows to skip.
func (q *QueryBuilder) Offset(n int) *QueryBuilder {
	return &QueryBuilder{db: q.db.Offset(n)}
}

// Paginate applies offset+limit from a *pagination.Params.
func (q *QueryBuilder) Paginate(p *pagination.Params) *QueryBuilder {
	if p == nil {
		return q
	}
	return &QueryBuilder{db: q.db.Offset(p.Offset()).Limit(p.Limit())}
}

// Distinct adds DISTINCT to the query.
func (q *QueryBuilder) Distinct(args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Distinct(args...)}
}

// GroupBy adds a GROUP BY clause.
func (q *QueryBuilder) GroupBy(fields ...string) *QueryBuilder {
	return &QueryBuilder{db: q.db.Group(strings.Join(fields, ", "))}
}

// Having adds a HAVING clause (use after GroupBy).
func (q *QueryBuilder) Having(query string, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Having(query, args...)}
}

// --- Eager Loading & Joins ---

// Preload preloads a relation (eager loading).
//   q.Preload("Roles")
//   q.Preload("Roles", "name = ?", "admin")   // conditional preload
func (q *QueryBuilder) Preload(relation string, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Preload(relation, args...)}
}

// Joins adds a JOIN clause.
//   q.Joins("INNER JOIN roles ON roles.id = user_roles.role_id")
func (q *QueryBuilder) Joins(query string, args ...interface{}) *QueryBuilder {
	return &QueryBuilder{db: q.db.Joins(query, args...)}
}

// --- Soft Delete ---

// WithTrashed includes soft-deleted records in query results.
func (q *QueryBuilder) WithTrashed() *QueryBuilder {
	return &QueryBuilder{db: q.db.Unscoped()}
}

// OnlyTrashed returns only soft-deleted records.
func (q *QueryBuilder) OnlyTrashed() *QueryBuilder {
	return &QueryBuilder{db: q.db.Unscoped().Where("deleted_at IS NOT NULL")}
}

// --- Scope Composition ---

// Scopes applies one or more GORM scope functions.
//   q.Scopes(database.Active(), database.Search("name", term))
func (q *QueryBuilder) Scopes(fns ...func(*gorm.DB) *gorm.DB) *QueryBuilder {
	return &QueryBuilder{db: q.db.Scopes(fns...)}
}

// --- Terminal Operations ---

// First retrieves the first matching record (ORDER BY primary key).
// Returns gorm.ErrRecordNotFound if no result.
func (q *QueryBuilder) First(ctx context.Context, dest interface{}) error {
	return q.db.WithContext(ctx).First(dest).Error
}

// Find retrieves all matching records into dest (slice pointer).
func (q *QueryBuilder) Find(ctx context.Context, dest interface{}) error {
	return q.db.WithContext(ctx).Find(dest).Error
}

// FindPaginated retrieves paginated records and returns total count.
// Returns (total, error).
//
// Example:
//   var users []gormUser
//   total, err := db.Query().
//       Where("role = ?", "admin").
//       FindPaginated(ctx, &users, params)
func (q *QueryBuilder) FindPaginated(ctx context.Context, dest interface{}, p *pagination.Params) (int64, error) {
	var total int64
	base := q.db.WithContext(ctx)

	if err := base.Count(&total).Error; err != nil {
		return 0, err
	}

	if p == nil {
		p = &pagination.Params{Page: 1, PerPage: pagination.MaxPerPage}
	}

	return total, base.Offset(p.Offset()).Limit(p.Limit()).Find(dest).Error
}

// Count returns the number of matching rows.
func (q *QueryBuilder) Count(ctx context.Context) (int64, error) {
	var count int64
	return count, q.db.WithContext(ctx).Count(&count).Error
}

// Exists returns true if at least one matching row exists.
func (q *QueryBuilder) Exists(ctx context.Context) (bool, error) {
	count, err := q.Count(ctx)
	return count > 0, err
}

// Pluck retrieves a single column from matching rows into dest (slice pointer).
//   var emails []string
//   q.Where("role = ?", "admin").Pluck(ctx, "email", &emails)
func (q *QueryBuilder) Pluck(ctx context.Context, field string, dest interface{}) error {
	return q.db.WithContext(ctx).Pluck(field, dest).Error
}

// Create inserts a new record.
func (q *QueryBuilder) Create(ctx context.Context, value interface{}) error {
	return q.db.WithContext(ctx).Create(value).Error
}

// Save creates or updates a record (upsert by primary key).
func (q *QueryBuilder) Save(ctx context.Context, value interface{}) error {
	return q.db.WithContext(ctx).Save(value).Error
}

// Update sets a single column on matching rows.
//   q.Where("id = ?", id).Update(ctx, "is_read", true)
func (q *QueryBuilder) Update(ctx context.Context, column string, value interface{}) error {
	return q.db.WithContext(ctx).Update(column, value).Error
}

// Updates sets multiple columns on matching rows using a map or struct.
func (q *QueryBuilder) Updates(ctx context.Context, values interface{}) error {
	return q.db.WithContext(ctx).Updates(values).Error
}

// Delete removes matching rows (soft-delete if model has DeletedAt).
func (q *QueryBuilder) Delete(ctx context.Context, value interface{}) error {
	return q.db.WithContext(ctx).Delete(value).Error
}

// DB returns the underlying *gorm.DB for advanced or unsupported operations.
// Use sparingly.
func (q *QueryBuilder) DB() *gorm.DB {
	return q.db
}
