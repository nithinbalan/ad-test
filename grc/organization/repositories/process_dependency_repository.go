package repositories

import (
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ProcessDependencyRepositoryInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyListDTO, error)
	GetByID(id uint) (*models.ProcessDependency, error)
	GetFullDetailsByID(id uint) (*models.ProcessDependency, error)
	Create(item *models.ProcessDependency) (uint, error)
	Update(updatedProcessDependency *models.ProcessDependency) error
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

func (repo *ProcessDependencyRepository) GetAll(processID uint, params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyListDTO, error) {
	var result []*orgDTOs.ProcessDependencyListDTO

	query := repo.db.Table("process_dependencies").
		Select(`
			process_dependencies.id,
			process_dependencies.process_id,
			parent_process.title_en AS process_title_en,
			parent_process.title_ar AS process_title_ar,

			process_dependencies.process_dependency_sub_type_id,
			process_dependency_sub_types.title_en AS process_dependency_sub_type_title_en,
			process_dependency_sub_types.title_ar AS process_dependency_sub_type_title_ar,

			process_dependency_types.title_en AS process_dependency_type_title_en,
			process_dependency_types.title_ar AS process_dependency_type_title_ar,

			process_dependencies.dependent_process_id,
			dependent_process.title_en AS dependent_process_title_en,
			dependent_process.title_ar AS dependent_process_title_ar,

			process_dependencies.organization_unit_id,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,

			process_dependencies.asset_id,
			asset_register.title_en AS asset_title_en,
			asset_register.title_ar AS asset_title_ar,

			process_dependencies.title,
			process_dependencies.description
		`).
		Joins(`LEFT JOIN processes AS parent_process ON parent_process.id = process_dependencies.process_id`).
		Joins(`LEFT JOIN process_dependency_sub_types ON process_dependency_sub_types.id = process_dependencies.process_dependency_sub_type_id`).
		Joins(`LEFT JOIN process_dependency_types ON process_dependency_types.id = process_dependency_sub_types.process_dependency_type_id`).
		Joins(`LEFT JOIN processes AS dependent_process ON dependent_process.id = process_dependencies.dependent_process_id`).
		Joins(`LEFT JOIN organization_units ON organization_units.id = process_dependencies.organization_unit_id`).
		Joins(`LEFT JOIN asset_register ON asset_register.id = process_dependencies.asset_id`).
		Where("process_dependencies.process_id = ?", processID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			process_dependencies.title LIKE ? OR
			process_dependencies.description LIKE ? OR
			process_dependency_sub_types.title_en LIKE ? OR
			process_dependency_types.title_en LIKE ? OR
			dependent_process.title_en LIKE ? OR
			organization_units.title_en LIKE ? OR
			asset_register.title_en LIKE ?
		`, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "process_dependencies.title"
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
	err := repo.db.First(&processDependency, id).Error
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

func (repo *ProcessDependencyRepository) Update(updatedProcessDependency *models.ProcessDependency) error {

	if err := repo.db.Save(updatedProcessDependency).Error; err != nil {
		return err
	}

	return nil
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
