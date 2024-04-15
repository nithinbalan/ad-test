package repositories

import (
	"gorm.io/gorm"
)

type UserSettingRepositoryInterface interface {
}

type UserSettingRepository struct {
	db *gorm.DB
}

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{
		db: db,
	}
}
