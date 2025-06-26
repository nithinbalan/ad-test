package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"sort"

	"gorm.io/gorm"
)

// Interfaces
type EscalationRepositoryInterface interface {
	GetAll(params dtos.EscalationListQueryParamsDTO) ([]*models.Escalation, error)
	GetByID(id uint) (*models.Escalation, error)
	GetFullDetailsByID(id uint) (*models.Escalation, error)

	GetAllEnabledModules() ([]*dtos.ModuleListDTO, error)
	GetEscalationByModuleID(moduleID uint) ([]*models.Escalation, error)
}

// implementation
type EscalationRepository struct {
	db *gorm.DB
}

func NewEscalationRepository(db *gorm.DB) *EscalationRepository {
	return &EscalationRepository{
		db: db,
	}
}

func (repo *EscalationRepository) GetAll(params dtos.EscalationListQueryParamsDTO) ([]*models.Escalation, error) {
	var escalations []*models.Escalation

	query := repo.db

	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&escalations)

	return escalations, result.Error
}

func (repo *EscalationRepository) GetByID(id uint) (*models.Escalation, error) {
	var status *models.Escalation

	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *EscalationRepository) GetFullDetailsByID(id uint) (*models.Escalation, error) {

	var escalation *models.Escalation

	result := repo.db.Preload("EscalationLevels", func(db *gorm.DB) *gorm.DB {
		return db.Order("day_offset ASC")
	}).Preload("EscalationLevels.Designation").
		First(&escalation, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return escalation, nil
}

func (repo *EscalationRepository) GetAllEnabledModules() ([]*dtos.ModuleListDTO, error) {
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
		`).Joins(`JOIN escalations ON escalations.module_id = modules.id`).
		Where("modules.module_id IS NULL").Where("modules.is_enable=?", true).Group("modules.id, modules.identifier, modules.title_en, modules.title_ar, modules.icon, modules.route, modules.is_enable, modules.[order]")

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

func (repo *EscalationRepository) GetEscalationByModuleID(moduleID uint) ([]*models.Escalation, error) {
	var summary []*models.Escalation
	err := repo.db.
		Preload("EscalationLevels").
		Where("module_id = ?", moduleID).
		Find(&summary).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}
