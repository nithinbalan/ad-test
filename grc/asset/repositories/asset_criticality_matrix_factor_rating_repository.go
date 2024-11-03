package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetCriticalityMatrixFactorRatingRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error)
	GetByID(id uint) (*models.AssetCriticalityMatrixFactorRating, error)
	Create(rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error)
	Update(id uint, rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error)
	Delete(id uint) error
}

type AssetCriticalityMatrixFactorRatingRepository struct {
	db *gorm.DB
}

func NewAssetCriticalityMatrixFactorRatingRepository(db *gorm.DB) *AssetCriticalityMatrixFactorRatingRepository {
	return &AssetCriticalityMatrixFactorRatingRepository{db: db}
}

func (repo *AssetCriticalityMatrixFactorRatingRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixFactorRating, error) {
	var ratings []*models.AssetCriticalityMatrixFactorRating
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
	result := query.Find(&ratings)

	return ratings, result.Error
}

func (repo *AssetCriticalityMatrixFactorRatingRepository) GetByID(id uint) (*models.AssetCriticalityMatrixFactorRating, error) {
	var rating models.AssetCriticalityMatrixFactorRating
	result := repo.db.First(&rating, id)
	return &rating, result.Error
}

func (repo *AssetCriticalityMatrixFactorRatingRepository) Create(rating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	result := repo.db.Create(rating)
	return rating, result.Error
}

func (repo *AssetCriticalityMatrixFactorRatingRepository) Update(id uint, updatedRating *models.AssetCriticalityMatrixFactorRating) (*models.AssetCriticalityMatrixFactorRating, error) {
	var rating models.AssetCriticalityMatrixFactorRating
	if err := repo.db.First(&rating, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&rating).Updates(updatedRating).Error; err != nil {
		return nil, err
	}

	return &rating, nil
}

func (repo *AssetCriticalityMatrixFactorRatingRepository) Delete(id uint) error {
	result := repo.db.Delete(&models.AssetCriticalityMatrixFactorRating{}, id)
	return result.Error
}
