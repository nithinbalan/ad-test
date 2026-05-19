package repositories

import (
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRISubCategoryRepositoryInterface interface {
	GetAll(params *dtos.KRISubCategoryListQueryParamsDTO) ([]*dtos.KRISubCategoryDTO, error)
	GetByID(id uint) (*models.KRISubCategory, error)
	Create(item *models.KRISubCategory) (*models.KRISubCategory, error)
	Update(id uint, item *models.KRISubCategory) (*models.KRISubCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(kriCategoryID uint, titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, kriCategoryID uint, titleEn string, titleAr string) (bool, error)
}

type KRISubCategoryRepository struct {
	db *gorm.DB
}

func NewKRISubCategoryRepository(db *gorm.DB) *KRISubCategoryRepository {
	return &KRISubCategoryRepository{db: db}
}

func (repo *KRISubCategoryRepository) GetAll(params *dtos.KRISubCategoryListQueryParamsDTO) ([]*dtos.KRISubCategoryDTO, error) {
	var items []*dtos.KRISubCategoryDTO

	query := repo.db.Table("kri_sub_categories").
		Select(`
			kri_sub_categories.*,
			kri_categories.title_en as kri_category_title_en,
			kri_categories.title_ar as kri_category_title_ar
		`).
		Joins(`
			LEFT JOIN kri_categories ON kri_categories.id = kri_sub_categories.kri_category_id
		`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(kri_sub_categories.title_en LIKE ? OR
			kri_sub_categories.title_ar LIKE ?)
		`, "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("kri_sub_categories.status = ?", params.Status)
	}

	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	result := query.Order(params.SortBy + " " + params.Sort).Find(&items)
	return items, result.Error
}

func (repo *KRISubCategoryRepository) GetByID(id uint) (*models.KRISubCategory, error) {
	var item *models.KRISubCategory
	result := repo.db.Preload("KRICategory").First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

func (repo *KRISubCategoryRepository) Create(item *models.KRISubCategory) (*models.KRISubCategory, error) {
	result := repo.db.Create(&item)
	return item, result.Error
}

func (repo *KRISubCategoryRepository) Update(id uint, updatedItem *models.KRISubCategory) (*models.KRISubCategory, error) {
	var item *models.KRISubCategory
	if err := repo.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Model(&item).Updates(updatedItem).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *KRISubCategoryRepository) Delete(id uint) error {
	var item *models.KRISubCategory
	if err := repo.db.First(&item, id).Error; err != nil {
		return err
	}
	result := repo.db.Delete(&models.KRISubCategory{}, id)
	return result.Error
}

func (repo *KRISubCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.KRISubCategory{}).Where("id IN ? AND status = ?", ids, "inactive").Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.KRISubCategory{}).Where("id IN ?", idsToUpdate).Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *KRISubCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.KRISubCategory{}).Where("id IN ? AND status = ?", ids, "active").Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.KRISubCategory{}).Where("id IN ?", idsToUpdate).Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *KRISubCategoryRepository) ExistsByTitle(kriCategoryID uint, titleEn string, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.KRISubCategory{}).
		Where("kri_category_id = ? AND (title_en = ? OR title_ar = ?)", kriCategoryID, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *KRISubCategoryRepository) ExistsByTitleExcludingID(id uint, kriCategoryID uint, titleEn string, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.KRISubCategory{}).
		Where("id <> ? AND kri_category_id = ? AND (title_en = ? OR title_ar = ?)", id, kriCategoryID, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
