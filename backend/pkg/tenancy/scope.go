package tenancy

import (
	"gorm.io/gorm"
)

/**
 * TenantScope is a GORM global scope that automatically filters queries by tenant_id.
 */
type TenantScope struct {
	TenantID string
	Bypass   bool
}

func (s TenantScope) Name() string {
	return "kodia:tenancy"
}

func (s TenantScope) Initialize(db *gorm.DB) error {
	// Register callbacks for automatic tenant filtering
	db.Callback().Query().Before("gorm:query").Register("kodia:tenancy:query", s.applyScope)
	db.Callback().Update().Before("gorm:update").Register("kodia:tenancy:update", s.applyScope)
	db.Callback().Delete().Before("gorm:delete").Register("kodia:tenancy:delete", s.applyScope)
	db.Callback().Create().Before("gorm:create").Register("kodia:tenancy:create", s.applyCreate)
	
	return nil
}

func (s TenantScope) applyScope(db *gorm.DB) {
	if s.Bypass || s.TenantID == "" {
		return
	}

	// Double check if the model is Tenantable (has a tenant_id field)
	if _, ok := db.Statement.Schema.FieldsByDBName["tenant_id"]; ok {
		db.Where("tenant_id = ?", s.TenantID)
	}
}

func (s TenantScope) applyCreate(db *gorm.DB) {
	if s.Bypass || s.TenantID == "" {
		return
	}

	// Automatically set tenant_id on create if field exists
	if field, ok := db.Statement.Schema.FieldsByDBName["tenant_id"]; ok {
		field.Set(db.Statement.Context, db.Statement.ReflectValue, s.TenantID)
	}
}

/**
 * Filter is a manual scope that can be applied to a specific DB instance.
 */
func Filter(db *gorm.DB, tenantID string, isSuperAdmin bool) *gorm.DB {
	if isSuperAdmin || tenantID == "" {
		return db
	}

	return db.Where("tenant_id = ?", tenantID)
}
