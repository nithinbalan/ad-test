package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ReminderContextRepositoryInterface interface {
	GetAllReminderContexts() ([]*models.ReminderContext, error)
	GetByID(id uint) (*models.ReminderContext, error)
}

// implementation
type ReminderContextRepository struct {
	db *gorm.DB
}

func NewReminderContextRepository(db *gorm.DB) *ReminderContextRepository {
	return &ReminderContextRepository{
		db: db,
	}
}

func (repo *ReminderContextRepository) GetAllReminderContexts() ([]*models.ReminderContext, error) {
	var contexts []*models.ReminderContext

	result := repo.db.Find(&contexts)

	if result.Error != nil {
		return nil, result.Error
	}

	return contexts, result.Error
}

func (repo *ReminderContextRepository) GetByID(id uint) (*models.ReminderContext, error) {
	var contexts *models.ReminderContext

	result := repo.db.First(&contexts, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return contexts, result.Error
}
