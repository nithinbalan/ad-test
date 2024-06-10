package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskResponseStrategyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskResponseStrategy, error)
	GetByID(id uint) (*models.RiskResponseStrategy, error)
	Create(strategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error)
	Update(id uint, strategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type RiskResponseStrategyRepository struct {
	db *gorm.DB
}

func NewRiskResponseStrategyRepository(db *gorm.DB) *RiskResponseStrategyRepository {
	return &RiskResponseStrategyRepository{
		db: db,
	}
}

func (repo *RiskResponseStrategyRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskResponseStrategy, error) {

	var riskResponseStrategies []*models.RiskResponseStrategy

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
	result := query.Find(&riskResponseStrategies)

	return riskResponseStrategies, result.Error
}

func (repo *RiskResponseStrategyRepository) GetByID(id uint) (*models.RiskResponseStrategy, error) {
	var strategy *models.RiskResponseStrategy

	result := repo.db.First(&strategy, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategy, nil
}

func (repo *RiskResponseStrategyRepository) Create(strategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error) {
	result := repo.db.Create(&strategy)

	return strategy, result.Error
}

func (repo *RiskResponseStrategyRepository) Update(id uint, updatedStrategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error) {
	var strategy *models.RiskResponseStrategy

	if err := repo.db.First(&strategy, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&strategy).Updates(updatedStrategy).Error; err != nil {
		return nil, err
	}

	return strategy, nil
}

func (repo *RiskResponseStrategyRepository) Delete(id uint) error {

	var strategy *models.RiskResponseStrategy
	if err := repo.db.First(&strategy, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskResponseStrategy{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskResponseStrategyRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.RiskResponseStrategy{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskResponseStrategy{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *RiskResponseStrategyRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.RiskResponseStrategy{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskResponseStrategy{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *RiskResponseStrategyRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskResponseStrategy{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *RiskResponseStrategyRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.RiskResponseStrategy{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
