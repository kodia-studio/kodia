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
// This hook is called before every save operation (Create, Update).
func (m *BaseModel) BeforeSave(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

// Timestamps embeds commonly-used timestamp fields.
// Use this for simple cases where you just need created/updated timestamps.
type Timestamps struct {
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
