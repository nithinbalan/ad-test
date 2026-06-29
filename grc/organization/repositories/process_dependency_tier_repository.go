package repositories

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessDependencyTierRepositoryInterface interface {
	GetAll(params *dtos.ProcessDependencyTierListQueryParamsDTO) ([]*models.ProcessDependencyTier, error)
	GetByID(id uint) (*models.ProcessDependencyTier, error)
	Create(tier *models.ProcessDependencyTier) (uint, error)
	Update(tier *models.ProcessDependencyTier) error
	Delete(id uint) error
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type ProcessDependencyTierRepository struct {
	db *gorm.DB
}

func NewProcessDependencyTierRepository(db *gorm.DB) *ProcessDependencyTierRepository {
	return &ProcessDependencyTierRepository{
		db: db,
	}
}

func (repo *ProcessDependencyTierRepository) GetAll(params *dtos.ProcessDependencyTierListQueryParamsDTO) ([]*models.ProcessDependencyTier, error) {
	var tiers []*models.ProcessDependencyTier

	query := repo.db

	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&tiers)

	return tiers, result.Error
}

func (repo *ProcessDependencyTierRepository) GetByID(id uint) (*models.ProcessDependencyTier, error) {
	var tier *models.ProcessDependencyTier

	result := repo.db.First(&tier, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return tier, nil
}

func (repo *ProcessDependencyTierRepository) Create(tier *models.ProcessDependencyTier) (uint, error) {
	if err := repo.db.Create(tier).Error; err != nil {
		return 0, err
	}

	return tier.ID, nil
}

func (repo *ProcessDependencyTierRepository) Update(tier *models.ProcessDependencyTier) error {
	return repo.db.Save(tier).Error
}

func (repo *ProcessDependencyTierRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessDependencyTier{}, id).Error
}

func (repo *ProcessDependencyTierRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessDependencyTier{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	return count > 0, err
}

func (repo *ProcessDependencyTierRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessDependencyTier{}).
		Where("id != ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	return count > 0, err
}
