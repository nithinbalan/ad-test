package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

// Interfaces
type ReminderContextRepositoryInterface interface {
	GetAllReminderContexts(params dtos.ReminderContextListQueryParamsDTO) ([]*models.ReminderContext, error)
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

func (repo *ReminderContextRepository) GetAllReminderContexts(params dtos.ReminderContextListQueryParamsDTO) ([]*models.ReminderContext, error) {
	var contexts []*models.ReminderContext

	query := repo.db.Model(&models.ReminderContext{})

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	err := query.Find(&contexts).Error
	return contexts, err
}

func (repo *ReminderContextRepository) GetByID(id uint) (*models.ReminderContext, error) {
	var contexts *models.ReminderContext

	result := repo.db.First(&contexts, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return contexts, result.Error
}
