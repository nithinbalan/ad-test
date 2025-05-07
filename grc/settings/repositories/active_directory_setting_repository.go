package repositories

import (
	"gorm.io/gorm"
)

// Interfaces
type ActiveDirectorySettingRepositoryInterface interface {
}

// Implementation
type ActiveDirectorySettingRepository struct {
	db *gorm.DB
}

func NewActiveDirectorySettingRepository(db *gorm.DB) *ActiveDirectorySettingRepository {
	return &ActiveDirectorySettingRepository{
		db: db,
	}
}
