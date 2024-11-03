package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetCriticalityMatrixFactorRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactor, error)
	GetByID(id uint) (*models.AssetCriticalityMatrixFactor, error)
	Create(factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error)
	Update(id uint, factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error)
	Delete(id uint) error
}

type AssetCriticalityMatrixFactorRepository struct {
	db *gorm.DB
}

func NewAssetCriticalityMatrixFactorRepository(db *gorm.DB) *AssetCriticalityMatrixFactorRepository {
	return &AssetCriticalityMatrixFactorRepository{db: db}
}

func (repo *AssetCriticalityMatrixFactorRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactor, error) {
	var factors []*models.AssetCriticalityMatrixFactor
	query := repo.db

	if keyword != "" {
		query = query.Where("title_en LIKE ? OR title_ar LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
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
	result := query.Find(&factors)

	return factors, result.Error
}

func (repo *AssetCriticalityMatrixFactorRepository) GetByID(id uint) (*models.AssetCriticalityMatrixFactor, error) {
	var factor models.AssetCriticalityMatrixFactor
	result := repo.db.First(&factor, id)
	return &factor, result.Error
}

func (repo *AssetCriticalityMatrixFactorRepository) Create(factor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error) {
	result := repo.db.Create(factor)
	return factor, result.Error
}

func (repo *AssetCriticalityMatrixFactorRepository) Update(id uint, updatedFactor *models.AssetCriticalityMatrixFactor) (*models.AssetCriticalityMatrixFactor, error) {
	var factor models.AssetCriticalityMatrixFactor
	if err := repo.db.First(&factor, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&factor).Updates(updatedFactor).Error; err != nil {
		return nil, err
	}

	return &factor, nil
}

func (repo *AssetCriticalityMatrixFactorRepository) Delete(id uint) error {
	result := repo.db.Delete(&models.AssetCriticalityMatrixFactor{}, id)
	return result.Error
}
