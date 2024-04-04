package utility

import "gorm.io/gorm"

// DatabaseTransaction abstracts transaction operations

type DatabaseTransaction interface {
	DB() *gorm.DB
	Commit() error
	Rollback() error
	// You can add other transaction-related methods here.
}

// GormDatabaseTransaction wraps a GORM DB transaction.
type GormDatabaseTransaction struct {
	tx *gorm.DB
}

func (gdt *GormDatabaseTransaction) DB() *gorm.DB {
	return gdt.tx
}

func (gdt *GormDatabaseTransaction) Commit() error {
	return gdt.tx.Commit().Error
}

func (gdt *GormDatabaseTransaction) Rollback() error {
	return gdt.tx.Rollback().Error
}

// NewGormDatabaseTransaction starts a new GORM transaction and returns our wrapper.
func NewGormDatabaseTransaction(db *gorm.DB) DatabaseTransaction {
	return &GormDatabaseTransaction{tx: db.Begin()}
}
