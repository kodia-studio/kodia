package database

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

// ErrNotFound is returned when a record is not found.
var ErrNotFound = errors.New("record not found")

// BaseRepository provides generic CRUD operations for a GORM model T.
// Embed this in concrete repositories to get access to Query() and basic CRUD.
//
// Example:
//   type UserRepository struct {
//       database.BaseRepository[gormUser]
//   }
//
//   func NewUserRepository(db *gorm.DB) *UserRepository {
//       return &UserRepository{BaseRepository: database.NewBaseRepository[gormUser](db)}
//   }
//
//   // Then inside methods:
//   var users []gormUser
//   total, err := r.Query().Where("role = ?", "admin").FindPaginated(ctx, &users, params)
type BaseRepository[T any] struct {
	db *gorm.DB
}

// NewBaseRepository creates a new BaseRepository for GORM model type T.
func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return BaseRepository[T]{db: db}
}

// Query returns a fresh QueryBuilder scoped to model type T.
// This is the primary entry point for building queries.
func (r *BaseRepository[T]) Query() *QueryBuilder {
	var m T
	return NewQuery(r.db.Model(&m))
}

// DB returns the underlying *gorm.DB.
func (r *BaseRepository[T]) DB() *gorm.DB {
	return r.db
}

// RawCreate inserts a new record at the GORM model level.
func (r *BaseRepository[T]) RawCreate(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Create(model).Error
}

// RawFindByID retrieves a record by primary key at the GORM model level.
// Returns ErrNotFound if not found.
func (r *BaseRepository[T]) RawFindByID(ctx context.Context, id string) (*T, error) {
	var model T
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return &model, err
}

// RawSave creates or updates a record by primary key.
func (r *BaseRepository[T]) RawSave(ctx context.Context, model *T) error {
	return r.db.WithContext(ctx).Save(model).Error
}

// RawDelete soft-deletes by ID (or hard-deletes if model has no DeletedAt).
func (r *BaseRepository[T]) RawDelete(ctx context.Context, id string) error {
	var model T
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model).Error
}

// RawRestore recovers a soft-deleted record by ID.
func (r *BaseRepository[T]) RawRestore(ctx context.Context, id string) error {
	var model T
	return r.db.WithContext(ctx).Unscoped().
		Model(&model).Where("id = ?", id).Update("deleted_at", nil).Error
}

// RawForceDelete permanently removes a record by ID.
func (r *BaseRepository[T]) RawForceDelete(ctx context.Context, id string) error {
	var model T
	return r.db.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&model).Error
}

// RawCount returns the total number of active (non-soft-deleted) records.
func (r *BaseRepository[T]) RawCount(ctx context.Context) (int64, error) {
	var count int64
	var model T
	return count, r.db.WithContext(ctx).Model(&model).Count(&count).Error
}

// RawExists returns true if a record with the given ID exists (not soft-deleted).
func (r *BaseRepository[T]) RawExists(ctx context.Context, id string) (bool, error) {
	var count int64
	var model T
	err := r.db.WithContext(ctx).Model(&model).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
