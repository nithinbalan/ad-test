package repositories

import (
	"errors"
	"jethur/grc/organization/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ControlImplementationStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlImplementationStatus, error)
	GetByID(id uint) (*models.ControlImplementationStatus, error)
	GetByType(status string) (*models.ControlImplementationStatus, error)
}

// implementation
type ControlImplementationStatusRepository struct {
	db *gorm.DB
}

func NewControlImplementationStatusRepository(db *gorm.DB) *ControlImplementationStatusRepository {
	return &ControlImplementationStatusRepository{
		db: db,
	}
}

func (repo *ControlImplementationStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlImplementationStatus, error) {

	var controlImplementationStatuses []*models.ControlImplementationStatus

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
	result := query.Find(&controlImplementationStatuses)

	return controlImplementationStatuses, result.Error
}

func (repo *ControlImplementationStatusRepository) GetByID(id uint) (*models.ControlImplementationStatus, error) {
	var status *models.ControlImplementationStatus
	result := repo.db.First(&status, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ControlImplementationStatusRepository) GetByType(status string) (*models.ControlImplementationStatus, error) {
	var controlImplementationStatus *models.ControlImplementationStatus

	result := repo.db.Where("status = ?", status).First(&controlImplementationStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Control implementation status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return controlImplementationStatus, nil
}
