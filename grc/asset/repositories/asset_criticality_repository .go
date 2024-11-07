package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetCriticalityRepositoryInterface interface {
	Create(contribution *models.AssetCriticality) (uint, error)
}

type AssetCriticalityRepository struct {
	db *gorm.DB
}

func NewAssetCriticalityRepository(db *gorm.DB) *AssetCriticalityRepository {
	return &AssetCriticalityRepository{
		db: db,
	}
}

func (repo *AssetCriticalityRepository) Create(contribution *models.AssetCriticality) (uint, error) {
	if err := repo.db.Create(contribution).Error; err != nil {
		return 0, err
	}
	return contribution.ID, nil
}
