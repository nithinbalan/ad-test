package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type UserDocumentTypeRepositoryInterface interface {
	GetAll(params *orgDTOs.UserDocumentTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentTypeDTO, error)
	GetByID(id uint) (*models.UserDocumentType, error)
	Create(item *models.UserDocumentType) (uint, error)
	Update(item *models.UserDocumentType) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type UserDocumentTypeRepository struct {
	db *gorm.DB
}

func NewUserDocumentTypeRepository(db *gorm.DB) *UserDocumentTypeRepository {
	return &UserDocumentTypeRepository{
		db: db,
	}
}

func (repo *UserDocumentTypeRepository) GetAll(params *orgDTOs.UserDocumentTypeListQueryParamsDTO) ([]*orgDTOs.UserDocumentTypeDTO, error) {
	var result []*orgDTOs.UserDocumentTypeDTO

	query := repo.db.Table("user_document_types").
		Select(`
			user_document_types.*`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(user_document_types.title_en LIKE ? OR
			 user_document_types.title_ar LIKE ?)`,
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

func (repo *UserDocumentTypeRepository) GetByID(id uint) (*models.UserDocumentType, error) {
	var userDocumentType *models.UserDocumentType
	err := repo.db.First(&userDocumentType, id).Error
	return userDocumentType, err
}

func (repo *UserDocumentTypeRepository) Create(userDocumentType *models.UserDocumentType) (uint, error) {
	err := repo.db.Create(userDocumentType).Error
	return userDocumentType.ID, err
}

func (repo *UserDocumentTypeRepository) Update(userDocumentType *models.UserDocumentType) error {
	return repo.db.Model(&models.UserDocumentType{}).
		Where("id = ?", userDocumentType.ID).
		Select("TitleEn", "TitleAr", "Status").
		Updates(userDocumentType).Error
}

func (repo *UserDocumentTypeRepository) Delete(id uint) error {
	return repo.db.Delete(&models.UserDocumentType{}, id).Error
}

func (repo *UserDocumentTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.UserDocumentType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.UserDocumentType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *UserDocumentTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.UserDocumentType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.UserDocumentType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *UserDocumentTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.UserDocumentType{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *UserDocumentTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.UserDocumentType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
