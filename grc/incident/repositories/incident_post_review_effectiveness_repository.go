package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type IncidentPostReviewEffectivenessRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentPostReviewEffectiveness, error)
	GetByID(id uint) (*models.IncidentPostReviewEffectiveness, error)
}

// implementation
type IncidentPostReviewEffectivenessRepository struct {
	db *gorm.DB
}

func NewIncidentPostReviewEffectivenessRepository(db *gorm.DB) *IncidentPostReviewEffectivenessRepository {
	return &IncidentPostReviewEffectivenessRepository{
		db: db,
	}
}

func (repo *IncidentPostReviewEffectivenessRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentPostReviewEffectiveness, error) {

	var incidentPostReviewEffectivenesses []*models.IncidentPostReviewEffectiveness

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
	result := query.Find(&incidentPostReviewEffectivenesses)

	return incidentPostReviewEffectivenesses, result.Error
}

func (repo *IncidentPostReviewEffectivenessRepository) GetByID(id uint) (*models.IncidentPostReviewEffectiveness, error) {
	var status *models.IncidentPostReviewEffectiveness
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}
