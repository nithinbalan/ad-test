package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditRiskRatingRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditRiskRating, error)
	GetByID(id uint) (*models.AuditRiskRating, error)
	GetByType(riskRatingType string) (*models.AuditRiskRating, error)
}

// implementation
type AuditRiskRatingRepository struct {
	db *gorm.DB
}

func NewAuditRiskRatingRepository(db *gorm.DB) *AuditRiskRatingRepository {
	return &AuditRiskRatingRepository{
		db: db,
	}
}

func (repo *AuditRiskRatingRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditRiskRating, error) {

	var auditRiskRating []*models.AuditRiskRating

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
	result := query.Find(&auditRiskRating)

	return auditRiskRating, result.Error
}

func (repo *AuditRiskRatingRepository) GetByID(id uint) (*models.AuditRiskRating, error) {
	var auditRiskRating *models.AuditRiskRating
	result := repo.db.First(&auditRiskRating, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditRiskRating, nil
}

func (repo *AuditRiskRatingRepository) GetByType(riskRatingType string) (*models.AuditRiskRating, error) {
	var auditRiskRating *models.AuditRiskRating

	result := repo.db.Where("type = ?", riskRatingType).First(&auditRiskRating)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Risk Rating with type "+riskRatingType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditRiskRating, nil
}
