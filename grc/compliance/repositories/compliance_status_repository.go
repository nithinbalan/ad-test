package repositories

import (
	"jethur/grc/compliance/models"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error)
}

// implementation
type ComplianceStatusRepository struct {
	db *gorm.DB
}

func NewComplianceStatusRepository(db *gorm.DB) *ComplianceStatusRepository {
	return &ComplianceStatusRepository{
		db: db,
	}
}

func (repo *ComplianceStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error) {

	var complianceStatuses []*models.ComplianceStatus

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
	result := query.Find(&complianceStatuses)

	return complianceStatuses, result.Error
}
