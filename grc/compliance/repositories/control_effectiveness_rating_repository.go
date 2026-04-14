package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ControlEffectivenessRatingRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlEffectivenessRating, error)
	GetByID(id uint) (*models.ControlEffectivenessRating, error)
	GetByType(typeFilter string) (*models.ControlEffectivenessRating, error)
}

// implementation
type ControlEffectivenessRatingRepository struct {
	db *gorm.DB
}

func NewControlEffectivenessRatingRepository(db *gorm.DB) *ControlEffectivenessRatingRepository {
	return &ControlEffectivenessRatingRepository{
		db: db,
	}
}

func (repo *ControlEffectivenessRatingRepository) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlEffectivenessRating, error) {
	var controlEffectivenessRatings []*models.ControlEffectivenessRating

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

	if typeFilter == "control_missing" || typeFilter == "Ineffective" || typeFilter == "improvement_needed" || typeFilter == "partially_effective" || typeFilter == "very_effective" {
		query = query.Where("type = ?", typeFilter)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&controlEffectivenessRatings)

	return controlEffectivenessRatings, result.Error
}

func (repo *ControlEffectivenessRatingRepository) GetByID(id uint) (*models.ControlEffectivenessRating, error) {
	var controlEffectivenessRating *models.ControlEffectivenessRating

	result := repo.db.First(&controlEffectivenessRating, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlEffectivenessRating, result.Error
}

func (repo *ControlEffectivenessRatingRepository) GetByType(typeFilter string) (*models.ControlEffectivenessRating, error) {
	var controlEffectivenessRating *models.ControlEffectivenessRating

	result := repo.db.Where("type = ?", typeFilter).First(&controlEffectivenessRating)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlEffectivenessRating, nil
}
