package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ProcessDependencyRepositoryInterface interface {
	GetAll(params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyDTO, error)
	GetByID(id uint) (*models.ProcessDependency, error)
	Create(item *models.ProcessDependency) (uint, error)
	Update(item *models.ProcessDependency) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ProcessDependencyRepository struct {
	db *gorm.DB
}

func NewProcessDependencyRepository(db *gorm.DB) *ProcessDependencyRepository {
	return &ProcessDependencyRepository{
		db: db,
	}
}

func (repo *ProcessDependencyRepository) GetAll(params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyDTO, error) {
	var result []*orgDTOs.ProcessDependencyDTO

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

func (repo *ProcessDependencyRepository) GetByID(id uint) (*models.ProcessDependency, error) {
	var processDependency *models.ProcessDependency
	err := repo.db.Preload("ProcessDependencyType").First(&processDependency, id).Error
	return processDependency, err
}

func (repo *ProcessDependencyRepository) Create(processDependency *models.ProcessDependency) (uint, error) {
	err := repo.db.Create(processDependency).Error
	return processDependency.ID, err
}

func (repo *ProcessDependencyRepository) Update(processDependency *models.ProcessDependency) error {
	return repo.db.Model(&models.ProcessDependency{}).
		Where("id = ?", processDependency.ID).
		Select("ProcessDependencyTypeID", "TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status").
		Updates(processDependency).Error
}

func (repo *ProcessDependencyRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessDependency{}, id).Error
}

func (repo *ProcessDependencyRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProcessDependency{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProcessDependency{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProcessDependencyRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProcessDependency{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProcessDependency{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProcessDependencyRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessDependency{}).
		Where("(title_en = ? OR title_ar = ?)", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ProcessDependencyRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ProcessDependency{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
