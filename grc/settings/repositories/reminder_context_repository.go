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

	// Keyword search (title_en / title_ar)
	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Status filter
	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	// Other filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Sorting
	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute query
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
