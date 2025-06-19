package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StakeholderCategory, error)
	GetByID(id uint) (*models.StakeholderCategory, error)
	GetByType(category string) (*models.StakeholderCategory, error)
}

type StakeholderCategoryRepository struct {
	db *gorm.DB
}

func NewStakeholderCategoryRepository(db *gorm.DB) *StakeholderCategoryRepository {
	return &StakeholderCategoryRepository{
		db: db,
	}
}

func (repo *StakeholderCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.StakeholderCategory, error) {
	var stakeholderCategories []*models.StakeholderCategory
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
	result := query.Find(&stakeholderCategories)

	return stakeholderCategories, result.Error
}

func (repo *StakeholderCategoryRepository) GetByID(id uint) (*models.StakeholderCategory, error) {
	var stakeholderCategory *models.StakeholderCategory
	result := repo.db.First(&stakeholderCategory, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return stakeholderCategory, nil
}

func (repo *StakeholderCategoryRepository) GetByType(category string) (*models.StakeholderCategory, error) {
	var stakeholderCategory *models.StakeholderCategory
	result := repo.db.Where("type = ?", category).First(&stakeholderCategory)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Stakeholder Category with type "+category)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return stakeholderCategory, nil
}
