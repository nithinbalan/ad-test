package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type DistributionStatusRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.DistributionStatus, error)
	GetByID(id uint) (*models.DistributionStatus, error)
}

type DistributionStatusRepository struct {
	db *gorm.DB
}

func NewDistributionStatusRepository(db *gorm.DB) DistributionStatusRepositoryInterface {
	return &DistributionStatusRepository{db: db}
}

func (repo *DistributionStatusRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.DistributionStatus, error) {
	var data []*models.DistributionStatus
	query := repo.db.Model(&models.DistributionStatus{})

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	err := query.Order(sortBy + " " + sort).Find(&data).Error
	return data, err
}

func (repo *DistributionStatusRepository) GetByID(id uint) (*models.DistributionStatus, error) {
	var data models.DistributionStatus
	err := repo.db.First(&data, id).Error
	if err != nil {
		return nil, err
	}
	return &data, nil
}