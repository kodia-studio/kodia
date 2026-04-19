package database

import "gorm.io/gorm"

// WithTransaction wraps a function in a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
//
// Example:
//
//	err := database.WithTransaction(db, func(tx *gorm.DB) error {
//		if err := tx.Create(&user1).Error; err != nil {
//			return err
//		}
//		if err := tx.Create(&user2).Error; err != nil {
//			return err
//		}
//		return nil  // Commits both operations
//	})
func WithTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.Transaction(fn)
}

// SafeDelete performs a soft delete via transaction.
// Useful for critical deletes that need rollback protection.
func SafeDelete(db *gorm.DB, value interface{}) error {
	return WithTransaction(db, func(tx *gorm.DB) error {
		return tx.Delete(value).Error
	})
}

// SafeSave performs a save operation via transaction.
// Useful for critical updates that need rollback protection.
func SafeSave(db *gorm.DB, value interface{}) error {
	return WithTransaction(db, func(tx *gorm.DB) error {
		return tx.Save(value).Error
	})
}
