package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type UserDocumentSubTypeRepositoryInterface interface {
	GetAll(params *orgDTOs.UserDocumentSubTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentSubTypeDTO, error)
	GetByID(id uint) (*models.UserDocumentSubType, error)
	Create(item *models.UserDocumentSubType) (uint, error)
	Update(item *models.UserDocumentSubType) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type UserDocumentSubTypeRepository struct {
	db *gorm.DB
}

func NewUserDocumentSubTypeRepository(db *gorm.DB) *UserDocumentSubTypeRepository {
	return &UserDocumentSubTypeRepository{
		db: db,
	}
}

func (repo *UserDocumentSubTypeRepository) GetAll(params *orgDTOs.UserDocumentSubTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentSubTypeDTO, error) {
	var result []*orgDTOs.UserDocumentSubTypeDTO

	query := repo.db.Table("user_document_sub_types").
		Select(`
			user_document_sub_types.*, 
			user_document_types.title_en AS user_document_type_title_en, 
			user_document_types.title_ar AS user_document_type_title_ar`).
		Joins(`LEFT JOIN user_document_types ON user_document_types.id = user_document_sub_types.user_document_type_id`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(user_document_sub_types.title_en LIKE ? OR
			 user_document_sub_types.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&result).Error
	return result, err
}

func (repo *UserDocumentSubTypeRepository) GetByID(id uint) (*models.UserDocumentSubType, error) {
	var userDocumentSubType *models.UserDocumentSubType
	err := repo.db.Preload("UserDocumentType").First(&userDocumentSubType, id).Error
	return userDocumentSubType, err
}

func (repo *UserDocumentSubTypeRepository) Create(userDocumentSubType *models.UserDocumentSubType) (uint, error) {
	err := repo.db.Create(userDocumentSubType).Error
	return userDocumentSubType.ID, err
}

func (repo *UserDocumentSubTypeRepository) Update(userDocumentSubType *models.UserDocumentSubType) error {
	return repo.db.Model(&models.UserDocumentSubType{}).
		Where("id = ?", userDocumentSubType.ID).
		Select("UserDocumentTypeID", "TitleEn", "TitleAr", "Status").
		Updates(userDocumentSubType).Error
}

func (repo *UserDocumentSubTypeRepository) Delete(id uint) error {
	return repo.db.Delete(&models.UserDocumentSubType{}, id).Error
}

func (repo *UserDocumentSubTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.UserDocumentSubType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.UserDocumentSubType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *UserDocumentSubTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.UserDocumentSubType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.UserDocumentSubType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *UserDocumentSubTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.UserDocumentSubType{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *UserDocumentSubTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.UserDocumentSubType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
