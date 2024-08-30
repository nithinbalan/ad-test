package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type StatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Status, error)
	GetByID(id uint) (*models.Status, error)
	GetByType(strStatus string) (*models.Status, error)
}

// Implementation
type StatusRepository struct {
	db *gorm.DB
}

func NewStatusRepository(db *gorm.DB) *StatusRepository {
	return &StatusRepository{
		db: db,
	}
}

func (repo *StatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Status, error) {
	var statuses []*models.Status

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
	result := query.Find(&statuses)

	return statuses, result.Error
}

func (repo *StatusRepository) GetByID(id uint) (*models.Status, error) {
	var status *models.Status
	result := repo.db.First(&status, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *StatusRepository) GetByType(strStatus string) (*models.Status, error) {
	var status *models.Status

	result := repo.db.Where("type = ?", strStatus).First(&status)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Status with status "+strStatus)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}
