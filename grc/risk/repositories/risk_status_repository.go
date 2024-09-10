package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskStatus, error)
	GetByID(id uint) (*models.RiskStatus, error)
	GetByType(status string) (*models.RiskStatus, error)
}

// implementation
type RiskStatusRepository struct {
	db *gorm.DB
}

func NewRiskStatusRepository(db *gorm.DB) *RiskStatusRepository {
	return &RiskStatusRepository{
		db: db,
	}
}

func (repo *RiskStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskStatus, error) {

	var riskStatuses []*models.RiskStatus

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
	result := query.Find(&riskStatuses)

	return riskStatuses, result.Error
}

func (repo *RiskStatusRepository) GetByID(id uint) (*models.RiskStatus, error) {
	var status *models.RiskStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *RiskStatusRepository) GetByType(status string) (*models.RiskStatus, error) {
	var riskStatus *models.RiskStatus

	result := repo.db.Where("type = ?", status).First(&riskStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskStatus, nil
}
