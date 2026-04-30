package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel is an embedded struct that provides common fields and GORM hooks
// for automatic UUID generation and timestamp management.
type BaseModel struct {
	ID        string    `gorm:"column:id;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// SoftDeleteModel embeds BaseModel and adds soft-delete support via gorm.DeletedAt.
// Use this when entities should be soft-deleted (marked as deleted, not actually removed).
type SoftDeleteModel struct {
	BaseModel
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at,omitempty"`
}

// BeforeSave generates a UUID for the ID field if it's empty.
func (m *BaseModel) BeforeSave(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// AfterUpdate is a placeholder for post-update logic.
func (m *BaseModel) AfterUpdate(tx *gorm.DB) error {
	// Custom post-update logic here
	return nil
}

// BeforeDelete is a placeholder for pre-deletion logic.
func (m *BaseModel) BeforeDelete(tx *gorm.DB) error {
	// Custom pre-deletion logic here (e.g. cascading deletes)
	return nil
}

// Timestamps embeds commonly-used timestamp fields.
// Use this for simple cases where you just need created/updated timestamps.
type Timestamps struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// ToDeletedAt converts *time.Time (domain layer) to gorm.DeletedAt (GORM model).
// Used in repository converters to transform domain entities to GORM models.
func ToDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

// FromDeletedAt converts gorm.DeletedAt (GORM model) to *time.Time (domain layer).
// Used in repository converters to transform GORM models to domain entities.
func FromDeletedAt(d gorm.DeletedAt) *time.Time {
	if !d.Valid {
		return nil
	}
	return &d.Time
}
