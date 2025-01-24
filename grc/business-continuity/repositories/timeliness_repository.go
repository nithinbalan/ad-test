package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type TimelinessRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Timeliness, error)
	GetByID(id uint) (*models.Timeliness, error)
	GetByType(status string) (*models.Timeliness, error)
}

// implementation
type TimelinessRepository struct {
	db *gorm.DB
}

func NewTimelinessRepository(db *gorm.DB) *TimelinessRepository {
	return &TimelinessRepository{
		db: db,
	}
}

func (repo *TimelinessRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Timeliness, error) {
	var timeliness []*models.Timeliness

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&timeliness)

	return timeliness, result.Error
}

func (repo *TimelinessRepository) GetByID(id uint) (*models.Timeliness, error) {
	var status *models.Timeliness
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *TimelinessRepository) GetByType(status string) (*models.Timeliness, error) {
	var timelinessType *models.Timeliness

	result := repo.db.Where("type = ?", status).First(&timelinessType)

	if result.Error != nil {
		return nil, result.Error
	}

	return timelinessType, nil
}
