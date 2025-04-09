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
	GetFullDetailsByID(id uint) (*models.ProcessDependency, error)
	Create(item *models.ProcessDependency) (uint, error)
	Update(item *models.ProcessDependency) error
	Delete(id uint) error

	GetAssetByID(id uint) (*models.AssetRegister, error)
	UpdateProcessDependencyRecoveryObjective(processDependency *models.ProcessDependency, tx ...utility.DatabaseTransaction) error
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
	err := repo.db.Preload("ProcessDependencySubType").First(&processDependency, id).Error
	return processDependency, err
}

func (repo *ProcessDependencyRepository) GetFullDetailsByID(id uint) (*models.ProcessDependency, error) {

	var processDependency *models.ProcessDependency

	result := repo.db.
		Preload("Process").
		Preload("ProcessDependencySubType").
		Preload("DependentProcess").
		Preload("OrganizationUnit").
		Preload("Asset").
		First(&processDependency, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return processDependency, nil
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

func (repo *ProcessDependencyRepository) GetAssetByID(id uint) (*models.AssetRegister, error) {
	var assetRegister *models.AssetRegister
	result := repo.db.First(&assetRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetRegister, nil
}

func (repo *ProcessDependencyRepository) UpdateProcessDependencyRecoveryObjective(processDependency *models.ProcessDependency, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(processDependency).Error; err != nil {
		return err
	}
	return nil
}
