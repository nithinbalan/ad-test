package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type FeedbackOptionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.FeedbackOption, error)
	GetByID(id uint) (*models.FeedbackOption, error)
	GetByType(status string) (*models.FeedbackOption, error)
}

// implementation
type FeedbackOptionRepository struct {
	db *gorm.DB
}

func NewFeedbackOptionRepository(db *gorm.DB) *FeedbackOptionRepository {
	return &FeedbackOptionRepository{
		db: db,
	}
}

func (repo *FeedbackOptionRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.FeedbackOption, error) {
	var feedbackOption []*models.FeedbackOption

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
	result := query.Find(&feedbackOption)

	return feedbackOption, result.Error
}

func (repo *FeedbackOptionRepository) GetByID(id uint) (*models.FeedbackOption, error) {
	var status *models.FeedbackOption
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *FeedbackOptionRepository) GetByType(status string) (*models.FeedbackOption, error) {
	var feedbackOption *models.FeedbackOption

	result := repo.db.Where("type = ?", status).First(&feedbackOption)

	if result.Error != nil {
		return nil, result.Error
	}

	return feedbackOption, nil
}
