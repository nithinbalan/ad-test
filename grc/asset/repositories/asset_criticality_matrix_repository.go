package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetCriticalityMatrixRepositoryInterface interface {
	GetFirstMatrixID() (uint, error)
}

type AssetCriticalityMatrixRepository struct {
	db *gorm.DB
}

func NewAssetCriticalityMatrixRepository(db *gorm.DB) AssetCriticalityMatrixRepositoryInterface {
	return &AssetCriticalityMatrixRepository{
		db: db,
	}
}

func (repo *AssetCriticalityMatrixRepository) GetFirstMatrixID() (uint, error) {
	var matrix models.AssetCriticalityMatrix
	err := repo.db.First(&matrix).Error
	if err != nil {
		return 0, err
	}
	return matrix.ID, nil
}
