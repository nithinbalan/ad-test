package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentCategory, error)
	GetByID(id uint) (*models.DocumentCategory, error)
	Create(category *models.DocumentCategory) (*models.DocumentCategory, error)
	Update(id uint, category *models.DocumentCategory) (*models.DocumentCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error)
}

// Implementation
type DocumentCategoryRepository struct {
	db *gorm.DB
}

func NewDocumentCategoryRepository(db *gorm.DB) *DocumentCategoryRepository {
	return &DocumentCategoryRepository{
		db: db,
	}
}

func (repo *DocumentCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentCategory, error) {
	var documentCategories []*models.DocumentCategory
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
	result := query.Find(&documentCategories)

	return documentCategories, result.Error
}

func (repo *DocumentCategoryRepository) GetByID(id uint) (*models.DocumentCategory, error) {
	var category *models.DocumentCategory
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}

func (repo *DocumentCategoryRepository) Create(category *models.DocumentCategory) (*models.DocumentCategory, error) {
	result := repo.db.Create(&category)

	return category, result.Error
}

func (repo *DocumentCategoryRepository) Update(id uint, updatedCategory *models.DocumentCategory) (*models.DocumentCategory, error) {
	var category *models.DocumentCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&category).Updates(updatedCategory).Error; err != nil {
		return nil, err
	}

	return category, nil
}

func (repo *DocumentCategoryRepository) Delete(id uint) error {
	var category *models.DocumentCategory
	if err := repo.db.First(&category, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.DocumentCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DocumentCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.DocumentCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.DocumentCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *DocumentCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *DocumentCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.DocumentCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
