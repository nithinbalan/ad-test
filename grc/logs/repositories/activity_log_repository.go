package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ActivityLogRepositoryInteface interface {
	SaveActivityLogs(logs []*models.ActivityLog, tx ...*gorm.DB) error
}

type ActivityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) *ActivityLogRepository {
	return &ActivityLogRepository{
		db: db,
	}
}

func (repo *ActivityLogRepository) SaveActivityLogs(logs []*models.ActivityLog, tx ...*gorm.DB) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Create(&logs).Error
}
