package repositories

import (
	"jethur/grc/logs/models"

	"gorm.io/gorm"
)

type ErrorLogRepositoryInterface interface {
	LogError(errorLog models.ErrorLog) error
}

type ErrorLogRepository struct {
	db *gorm.DB
}

func NewErrorLogRepository(db *gorm.DB) ErrorLogRepositoryInterface {
	return &ErrorLogRepository{
		db: db,
	}
}

func (repo *ErrorLogRepository) LogError(errorLog models.ErrorLog) error {
	return repo.db.Create(&errorLog).Error
}
