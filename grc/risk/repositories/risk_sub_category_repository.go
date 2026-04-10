package repositories

import (
	"jethur/grc/risk/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskSubCategoryRepositoryInterface interface {
	GetAll(params *dtos.RiskSubCategoryListQueryParamsDTO) ([]*dtos.RiskSubCategoryDTO, error)
	GetByID(id uint) (*models.RiskSubCategory, error)
	Create(item *models.RiskSubCategory) (*models.RiskSubCategory, error)
	Update(id uint, item *models.RiskSubCategory) (*models.RiskSubCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(riskCategoryID uint, titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, riskCategoryID uint, titleEn string, titleAr string) (bool, error)
}

type RiskSubCategoryRepository struct {
	db *gorm.DB
}

func NewRiskSubCategoryRepository(db *gorm.DB) *RiskSubCategoryRepository {
	return &RiskSubCategoryRepository{db: db}
}

func (repo *RiskSubCategoryRepository) GetAll(params *dtos.RiskSubCategoryListQueryParamsDTO) ([]*dtos.RiskSubCategoryDTO, error) {
	var items []*dtos.RiskSubCategoryDTO

	query := repo.db.Table("risk_sub_categories").
		Select(`
			risk_sub_categories.*,
			risk_categories.title_en as risk_category_title_en,
			risk_categories.title_ar as risk_category_title_ar
		`).
		Joins(`
			LEFT JOIN risk_categories ON risk_categories.id = risk_sub_categories.risk_category_id
		`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(risk_sub_categories.title_en LIKE ? OR
			risk_sub_categories.title_ar LIKE ?)
		`, "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("risk_sub_categories.status = ?", params.Status)
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

func (repo *RiskSubCategoryRepository) GetByID(id uint) (*models.RiskSubCategory, error) {
	var item *models.RiskSubCategory
	result := repo.db.Preload("RiskCategory").First(&item, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return item, nil
}

func (repo *RiskSubCategoryRepository) Create(item *models.RiskSubCategory) (*models.RiskSubCategory, error) {
	result := repo.db.Create(&item)
	return item, result.Error
}

func (repo *RiskSubCategoryRepository) Update(id uint, updatedItem *models.RiskSubCategory) (*models.RiskSubCategory, error) {
	var item *models.RiskSubCategory
	if err := repo.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Model(&item).Updates(updatedItem).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (repo *RiskSubCategoryRepository) Delete(id uint) error {
	var item *models.RiskSubCategory
	if err := repo.db.First(&item, id).Error; err != nil {
		return err
	}
	result := repo.db.Delete(&models.RiskSubCategory{}, id)
	return result.Error
}

func (repo *RiskSubCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.RiskSubCategory{}).Where("id IN ? AND status = ?", ids, "inactive").Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskSubCategory{}).Where("id IN ?", idsToUpdate).Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *RiskSubCategoryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.RiskSubCategory{}).Where("id IN ? AND status = ?", ids, "active").Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskSubCategory{}).Where("id IN ?", idsToUpdate).Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *RiskSubCategoryRepository) ExistsByTitle(riskCategoryID uint, titleEn string, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.RiskSubCategory{}).
		Where("risk_category_id = ? AND (title_en = ? OR title_ar = ?)", riskCategoryID, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *RiskSubCategoryRepository) ExistsByTitleExcludingID(id uint, riskCategoryID uint, titleEn string, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.RiskSubCategory{}).
		Where("id <> ? AND risk_category_id = ? AND (title_en = ? OR title_ar = ?)", id, riskCategoryID, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
