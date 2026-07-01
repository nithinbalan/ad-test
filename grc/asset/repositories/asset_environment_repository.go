package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetEnvironmentRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetEnvironment, error)
	GetByID(id uint) (*models.AssetEnvironment, error)
	GetByType(status string) (*models.AssetEnvironment, error)
}

type AssetEnvironmentRepository struct {
	db *gorm.DB
}

func NewAssetEnvironmentRepository(db *gorm.DB) *AssetEnvironmentRepository {
	return &AssetEnvironmentRepository{
		db: db,
	}
}

func (repo *AssetEnvironmentRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AssetEnvironment, error) {
	var assetEnvironments []*models.AssetEnvironment

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Find(&assetEnvironments)

	return assetEnvironments, result.Error
}

func (repo *AssetEnvironmentRepository) GetByID(id uint) (*models.AssetEnvironment, error) {
	var assetEnvironment *models.AssetEnvironment

	result := repo.db.First(&assetEnvironment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return assetEnvironment, nil
}

func (repo *AssetEnvironmentRepository) GetByType(status string) (*models.AssetEnvironment, error) {
	var assetEnvironment *models.AssetEnvironment

	result := repo.db.Where("type = ?", status).First(&assetEnvironment)

	if result.Error != nil {
		return nil, result.Error
	}

	return assetEnvironment, nil
}