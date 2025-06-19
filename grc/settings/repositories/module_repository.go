package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"sort"

	"gorm.io/gorm"
)

// Interfaces
type ModuleRepositoryInterface interface {
	GetModuleIDByIdentifier(identifier string) (uint, error)
	GetByID(id uint) (*models.Module, error)
	GetByIdentifier(identifier string) (*models.Module, error)
	GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error)
	GetFullDetailsByID(id uint) (*models.Module, error)
	Update(moduleId uint, module *models.Module) (uint, error)

	GetSystemPreference() (*models.SystemPreference, error)
	GetAllModules() ([]*dtos.ModuleListDTO, error)
}

// Implementation
type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{
		db: db,
	}
}

func (repo *ModuleRepository) GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error) {
	var modules []models.Module

	query := repo.db.Preload("SubModules", func(db *gorm.DB) *gorm.DB {
		if params.IsActiveSubModule {
			return db.Where("[modules].is_enable = ?", true).Order("[modules].[order] ASC")
		}
		return db.Order("[modules].[order] ASC")
	}).
		Where("module_id IS NULL").
		Select(`
			modules.id,
			modules.identifier,
			modules.title_en,
			modules.title_ar,
			modules.icon,
			modules.route,
			modules.is_enable,
			[modules].[order]  -- Escape the 'order' column
		`)

	// Apply keyword search
	if params.Keyword != "" {
		query = query.Where(`
			(modules.title_en LIKE ? OR
			modules.title_ar LIKE ? OR
			modules.identifier LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Apply sorting based on the provided fields
	if params.SortBy != "" && (params.Sort == "asc" || params.Sort == "desc") {
		// Escape the SortBy column if it is "order"
		if params.SortBy == "order" {
			params.SortBy = "[modules].[order]"
		}
		query = query.Order(params.SortBy + " " + params.Sort)
	} else {
		query = query.Order("[modules].[order] asc")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "active":
			query = query.Where("[modules].is_enable = ?", true)
		case "inactive":
			query = query.Where("[modules].is_enable = ?", false)
		}
	}

	if err := query.Find(&modules).Error; err != nil {
		return nil, err
	}

	var moduleList []dtos.ModuleListDTO
	for _, module := range modules {
		var submoduleDTOs []dtos.SubModuleListDTO
		if module.SubModules != nil {
			for _, submodule := range *module.SubModules {
				submoduleDTO := dtos.SubModuleListDTO{
					ID:         submodule.ID,
					Identifier: submodule.Identifier,
					TitleEn:    submodule.TitleEn,
					TitleAr:    submodule.TitleAr,
					Order:      submodule.Order,
					Icon:       submodule.Icon,
					Route:      submodule.Route,
					IsEnable:   submodule.IsEnable,
				}
				submoduleDTOs = append(submoduleDTOs, submoduleDTO)
			}

			sort.SliceStable(submoduleDTOs, func(i, j int) bool {
				return submoduleDTOs[i].Order < submoduleDTOs[j].Order
			})
		}

		moduleDTO := dtos.ModuleListDTO{
			ID:         module.ID,
			Identifier: module.Identifier,
			TitleEn:    module.TitleEn,
			TitleAr:    module.TitleAr,
			Order:      module.Order,
			Icon:       module.Icon,
			Route:      module.Route,
			IsEnable:   module.IsEnable,
			SubModules: submoduleDTOs,
		}

		moduleList = append(moduleList, moduleDTO)
	}

	sort.SliceStable(moduleList, func(i, j int) bool {
		return moduleList[i].Order < moduleList[j].Order
	})

	return &moduleList, nil
}

func (repo *ModuleRepository) GetModuleIDByIdentifier(identifier string) (uint, error) {
	var module models.Module
	if err := repo.db.Where("identifier = ?", identifier).First(&module).Error; err != nil {
		return 0, err
	}
	return module.ID, nil
}

func (repo *ModuleRepository) GetByID(id uint) (*models.Module, error) {
	var module *models.Module
	result := repo.db.First(&module, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return module, nil
}

func (repo *ModuleRepository) GetByIdentifier(identifier string) (*models.Module, error) {
	var module *models.Module
	if err := repo.db.Where("identifier = ?", identifier).First(&module).Error; err != nil {
		return nil, err
	}

	return module, nil
}

func (repo *ModuleRepository) GetFullDetailsByID(id uint) (*models.Module, error) {
	var module *models.Module // Declare the variable inside the function

	result := repo.db.
		Preload("SubModules").
		Preload("Parent").
		First(&module, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return module, nil
}

func (repo *ModuleRepository) Update(moduleId uint, module *models.Module) (uint, error) {
	db := repo.db
	if err := db.Save(module).Error; err != nil {
		return 0, err
	}

	return moduleId, nil
}

func (repo *ModuleRepository) GetSystemPreference() (*models.SystemPreference, error) {
	var systemPreference *models.SystemPreference
	result := repo.db.
		Preload("Timezone").
		Preload("Currency").
		First(&systemPreference)

	if result.Error != nil {
		return nil, result.Error
	}

	return systemPreference, nil
}

func (repo *ModuleRepository) GetAllModules() ([]*dtos.ModuleListDTO, error) {
	var modules []models.Module

	query := repo.db.Table("modules").
		Select(`
			modules.id,
			modules.identifier,
			modules.title_en,
			modules.title_ar,
			modules.icon,
			modules.route,
			modules.is_enable,
			[modules].[order]  -- Escape the 'order' column
		`).Where("module_id IS NULL")

	if err := query.Find(&modules).Error; err != nil {
		return nil, err
	}

	var moduleList []*dtos.ModuleListDTO
	for _, module := range modules {
		var submoduleDTOs []dtos.SubModuleListDTO
		if module.SubModules != nil {
			for _, submodule := range *module.SubModules {
				submoduleDTO := dtos.SubModuleListDTO{
					ID:         submodule.ID,
					Identifier: submodule.Identifier,
					TitleEn:    submodule.TitleEn,
					TitleAr:    submodule.TitleAr,
					Order:      submodule.Order,
					Icon:       submodule.Icon,
					Route:      submodule.Route,
					IsEnable:   submodule.IsEnable,
				}
				submoduleDTOs = append(submoduleDTOs, submoduleDTO)
			}

			sort.SliceStable(submoduleDTOs, func(i, j int) bool {
				return submoduleDTOs[i].Order < submoduleDTOs[j].Order
			})
		}

		moduleDTO := dtos.ModuleListDTO{
			ID:         module.ID,
			Identifier: module.Identifier,
			TitleEn:    module.TitleEn,
			TitleAr:    module.TitleAr,
			Order:      module.Order,
			Icon:       module.Icon,
			Route:      module.Route,
			IsEnable:   module.IsEnable,
			SubModules: submoduleDTOs,
		}

		moduleList = append(moduleList, &moduleDTO)
	}

	sort.SliceStable(moduleList, func(i, j int) bool {
		return moduleList[i].Order < moduleList[j].Order
	})

	return moduleList, nil
}
