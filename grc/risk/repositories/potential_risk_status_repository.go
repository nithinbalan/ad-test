package repositories

import (
	"errors"
	"jethur/grc/risk/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type PotentialRiskStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PotentialRiskStatus, error)
	GetByID(id uint) (*models.PotentialRiskStatus, error)
	GetByStatus(status string) (*models.PotentialRiskStatus, error)
}

// Implementation
type PotentialRiskStatusRepository struct {
	db *gorm.DB
}

func NewPotentialRiskStatusRepository(db *gorm.DB) *PotentialRiskStatusRepository {
	return &PotentialRiskStatusRepository{
		db: db,
	}
}

func (repo *PotentialRiskStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PotentialRiskStatus, error) {
	var potentialRiskStatuses []*models.PotentialRiskStatus

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
	result := query.Find(&potentialRiskStatuses)

	return potentialRiskStatuses, result.Error
}

func (repo *PotentialRiskStatusRepository) GetByID(id uint) (*models.PotentialRiskStatus, error) {
	var potentialRiskStatus *models.PotentialRiskStatus
	result := repo.db.First(&potentialRiskStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return potentialRiskStatus, result.Error
}

func (repo *PotentialRiskStatusRepository) GetByStatus(status string) (*models.PotentialRiskStatus, error) {
	var potentialRiskStatus *models.PotentialRiskStatus

	result := repo.db.Where("type = ?", status).First(&potentialRiskStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Potential Risk Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return potentialRiskStatus, nil
}
