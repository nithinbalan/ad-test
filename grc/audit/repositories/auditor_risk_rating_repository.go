package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditorRiskRatingRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditorRiskRating, error)
	GetByID(id uint) (*models.AuditorRiskRating, error)
	GetByType(riskRatingType string) (*models.AuditorRiskRating, error)
}

// implementation
type AuditorRiskRatingRepository struct {
	db *gorm.DB
}

func NewAuditorRiskRatingRepository(db *gorm.DB) *AuditorRiskRatingRepository {
	return &AuditorRiskRatingRepository{
		db: db,
	}
}

func (repo *AuditorRiskRatingRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditorRiskRating, error) {

	var auditorRiskRating []*models.AuditorRiskRating

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
	result := query.Find(&auditorRiskRating)

	return auditorRiskRating, result.Error
}

func (repo *AuditorRiskRatingRepository) GetByID(id uint) (*models.AuditorRiskRating, error) {
	var auditorRiskRating *models.AuditorRiskRating
	result := repo.db.First(&auditorRiskRating, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditorRiskRating, nil
}

func (repo *AuditorRiskRatingRepository) GetByType(riskRatingType string) (*models.AuditorRiskRating, error) {
	var auditorRiskRating *models.AuditorRiskRating

	result := repo.db.Where("type = ?", riskRatingType).First(&auditorRiskRating)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Risk Rating with type "+riskRatingType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditorRiskRating, nil
}
