package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type AssetStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetStatus, error)
	GetByID(id uint) (*models.AssetStatus, error)
	GetByType(status string) (*models.AssetStatus, error)
}

// Implementation
type AssetStatusRepository struct {
	db *gorm.DB
}

func NewAssetStatusRepository(db *gorm.DB) *AssetStatusRepository {
	return &AssetStatusRepository{
		db: db,
	}
}

func (repo *AssetStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetStatus, error) {
	var assetStatuses []*models.AssetStatus

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
	result := query.Find(&assetStatuses)

	return assetStatuses, result.Error
}

func (repo *AssetStatusRepository) GetByID(id uint) (*models.AssetStatus, error) {
	var status *models.AssetStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *AssetStatusRepository) GetByType(status string) (*models.AssetStatus, error) {
	var assetStatus *models.AssetStatus

	result := repo.db.Where("type = ?", status).First(&assetStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return assetStatus, nil
}
