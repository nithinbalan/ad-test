package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type InternalCommitteeStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.InternalCommitteeStatus, error)
	GetByID(id uint) (*models.InternalCommitteeStatus, error)
	GetByType(typeFilter string) (*models.InternalCommitteeStatus, error)
}

// implementation
type InternalCommitteeStatusRepository struct {
	db *gorm.DB
}

func NewInternalCommitteeStatusRepository(db *gorm.DB) *InternalCommitteeStatusRepository {
	return &InternalCommitteeStatusRepository{
		db: db,
	}
}

func (repo *InternalCommitteeStatusRepository) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.InternalCommitteeStatus, error) {
	var internalCommitteeStatuses []*models.InternalCommitteeStatus

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

	if typeFilter == "draft" || typeFilter == "submitted" || typeFilter == "in_review" || typeFilter == "in_progress" || typeFilter == "approved" || typeFilter == "reformation" || typeFilter == "disbanded" {
		query = query.Where("type = ?", typeFilter)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&internalCommitteeStatuses)

	return internalCommitteeStatuses, result.Error
}

func (repo *InternalCommitteeStatusRepository) GetByID(id uint) (*models.InternalCommitteeStatus, error) {
	var internalCommitteeStatus *models.InternalCommitteeStatus

	result := repo.db.First(&internalCommitteeStatus, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return internalCommitteeStatus, result.Error
}

func (repo *InternalCommitteeStatusRepository) GetByType(typeFilter string) (*models.InternalCommitteeStatus, error) {
	var internalCommitteeStatus *models.InternalCommitteeStatus

	result := repo.db.Where("type = ?", typeFilter).First(&internalCommitteeStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return internalCommitteeStatus, nil
}
