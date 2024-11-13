package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetCriticalityMatrixRatingScaleRepositoryInterface interface {
	GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixRatingScale, error)
	GetByID(id uint) (*models.AssetCriticalityMatrixRatingScale, error)
	Create(scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error)
	Update(id uint, scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error)
	Delete(id uint) error

	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)

	GetRatingScaleByScore(matrixID uint, score float64) (*models.AssetCriticalityMatrixRatingScale, error)
}

type AssetCriticalityMatrixRatingScaleRepository struct {
	db *gorm.DB
}

func NewAssetCriticalityMatrixRatingScaleRepository(db *gorm.DB) *AssetCriticalityMatrixRatingScaleRepository {
	return &AssetCriticalityMatrixRatingScaleRepository{db: db}
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) GetAll(matrixID uint, sort, sortBy, keyword, status string) ([]*models.AssetCriticalityMatrixRatingScale, error) {
	var scales []*models.AssetCriticalityMatrixRatingScale
	query := repo.db.Model(&models.AssetCriticalityMatrixRatingScale{}).
		Where("asset_criticality_matrix_id = ?", matrixID)

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
	result := query.Find(&scales)

	return scales, result.Error
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) GetByID(id uint) (*models.AssetCriticalityMatrixRatingScale, error) {
	var scale models.AssetCriticalityMatrixRatingScale
	result := repo.db.First(&scale, id)
	return &scale, result.Error
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) Create(scale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error) {
	result := repo.db.Create(scale)
	return scale, result.Error
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) Update(id uint, updatedScale *models.AssetCriticalityMatrixRatingScale) (*models.AssetCriticalityMatrixRatingScale, error) {
	var scale models.AssetCriticalityMatrixRatingScale
	if err := repo.db.First(&scale, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&scale).Updates(updatedScale).Error; err != nil {
		return nil, err
	}

	return &scale, nil
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) Delete(id uint) error {
	result := repo.db.Delete(&models.AssetCriticalityMatrixRatingScale{}, id)
	return result.Error
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetCriticalityMatrixRatingScale{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AssetCriticalityMatrixRatingScale{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssetCriticalityMatrixRatingScaleRepository) GetRatingScaleByScore(matrixID uint, score float64) (*models.AssetCriticalityMatrixRatingScale, error) {
	var ratingScale models.AssetCriticalityMatrixRatingScale
	err := repo.db.
		Where("asset_criticality_matrix_id = ? AND score_from <= ? AND score_to >= ?", matrixID, score, score).
		First(&ratingScale).Error
	if err != nil {
		return nil, err
	}
	return &ratingScale, nil
}
