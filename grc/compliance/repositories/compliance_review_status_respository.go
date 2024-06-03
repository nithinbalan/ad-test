package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceReviewStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceReviewStatus, error)
	GetByID(id uint) (*models.ComplianceReviewStatus, error)
	GetByType(status string) (*models.ComplianceReviewStatus, error)
}

// Implementation
type ComplianceReviewStatusRepository struct {
	db *gorm.DB
}

func NewComplianceReviewStatusRepository(db *gorm.DB) *ComplianceReviewStatusRepository {
	return &ComplianceReviewStatusRepository{
		db: db,
	}
}

func (repo *ComplianceReviewStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceReviewStatus, error) {

	var complianceReviewStatuses []*models.ComplianceReviewStatus

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
	result := query.Find(&complianceReviewStatuses)

	return complianceReviewStatuses, result.Error
}

func (repo *ComplianceReviewStatusRepository) GetByID(id uint) (*models.ComplianceReviewStatus, error) {
	var complianceReviewStatus *models.ComplianceReviewStatus
	result := repo.db.First(&complianceReviewStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return complianceReviewStatus, result.Error
}

func (repo *ComplianceReviewStatusRepository) GetByType(status string) (*models.ComplianceReviewStatus, error) {
	var complianceReviewStatus *models.ComplianceReviewStatus

	result := repo.db.Where("type = ?", status).First(&complianceReviewStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Compliance Review Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return complianceReviewStatus, nil
}
