package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ProcessDependencySubTypeRepositoryInterface interface {
	GetAll(params *orgDTOs.ProcessDependencySubTypeListQueryParamsDTO) ([]*orgDTOs.ProcessDependencySubTypeDTO, error)
	GetByID(id uint) (*models.ProcessDependencySubType, error)
	Create(item *models.ProcessDependencySubType) (uint, error)
	Update(item *models.ProcessDependencySubType) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ProcessDependencySubTypeRepository struct {
	db *gorm.DB
}

func NewProcessDependencySubTypeRepository(db *gorm.DB) *ProcessDependencySubTypeRepository {
	return &ProcessDependencySubTypeRepository{
		db: db,
	}
}

func (repo *ProcessDependencySubTypeRepository) GetAll(params *orgDTOs.ProcessDependencySubTypeListQueryParamsDTO) ([]*orgDTOs.ProcessDependencySubTypeDTO, error) {
	var result []*orgDTOs.ProcessDependencySubTypeDTO

	query := repo.db.Table("process_dependency_sub_types").
		Select(`
			process_dependency_sub_types.*, 
			process_dependency_types.title_en AS process_dependency_type_title_en, 
			process_dependency_types.title_ar AS process_dependency_type_title_ar`).
		Joins(`LEFT JOIN process_dependency_types ON process_dependency_types.id = process_dependency_sub_types.process_dependency_type_id`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
			(process_dependency_sub_types.title_en LIKE ? OR
			 process_dependency_sub_types.title_ar LIKE ? OR
			 process_dependency_sub_types.description_en LIKE ? OR
			 process_dependency_sub_types.description_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
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

func (repo *ProcessDependencySubTypeRepository) GetByID(id uint) (*models.ProcessDependencySubType, error) {
	var processDependencySubType *models.ProcessDependencySubType
	err := repo.db.Preload("ProcessDependencyType").First(&processDependencySubType, id).Error
	return processDependencySubType, err
}

func (repo *ProcessDependencySubTypeRepository) Create(processDependencySubType *models.ProcessDependencySubType) (uint, error) {
	err := repo.db.Create(processDependencySubType).Error
	return processDependencySubType.ID, err
}

func (repo *ProcessDependencySubTypeRepository) Update(processDependencySubType *models.ProcessDependencySubType) error {
	return repo.db.Model(&models.ProcessDependencySubType{}).
		Where("id = ?", processDependencySubType.ID).
		Select("ProcessDependencyTypeID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status").
		Updates(processDependencySubType).Error
}

func (repo *ProcessDependencySubTypeRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessDependencySubType{}, id).Error
}

func (repo *ProcessDependencySubTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProcessDependencySubType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProcessDependencySubType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProcessDependencySubTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProcessDependencySubType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProcessDependencySubType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProcessDependencySubTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessDependencySubType{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ProcessDependencySubTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessDependencySubType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
