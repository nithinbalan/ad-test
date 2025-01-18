package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BcpStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpStatus, error)
	GetByID(id uint) (*models.BcpStatus, error)
	GetByType(status string) (*models.BcpStatus, error)
}

// implementation
type BcpStatusRepository struct {
	db *gorm.DB
}

func NewBcpStatusRepository(db *gorm.DB) *BcpStatusRepository {
	return &BcpStatusRepository{
		db: db,
	}
}

func (repo *BcpStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpStatus, error) {
	var bcpStatuses []*models.BcpStatus

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
	result := query.Find(&bcpStatuses)

	return bcpStatuses, result.Error
}

func (repo *BcpStatusRepository) GetByID(id uint) (*models.BcpStatus, error) {
	var status *models.BcpStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *BcpStatusRepository) GetByType(status string) (*models.BcpStatus, error) {
	var bcpStatus *models.BcpStatus

	result := repo.db.Where("type = ?", status).First(&bcpStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return bcpStatus, nil
}
