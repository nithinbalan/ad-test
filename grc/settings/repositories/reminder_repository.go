package repositories

import (
	"sort"

	"jethur/grc/models"
	"jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

// Interfaces
type ReminderRepositoryInterface interface {
	GetAll(params dtos.ReminderListQueryParamsDTO) ([]*models.Reminder, error)
	GetByID(id uint) (*models.Reminder, error)
	GetByType(status string) (*models.Reminder, error)

	GetAllEnabledModules() ([]*dtos.ModuleListDTO, error)
	GetReminderByModuleID(moduleID uint) ([]*models.Reminder, error)
}

// implementation
type ReminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) *ReminderRepository {
	return &ReminderRepository{
		db: db,
	}
}

func (repo *ReminderRepository) GetAll(params dtos.ReminderListQueryParamsDTO) ([]*models.Reminder, error) {
	var reminders []*models.Reminder

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
	result := query.Find(&reminders)

	return reminders, result.Error
}

func (repo *ReminderRepository) GetByID(id uint) (*models.Reminder, error) {
	var status *models.Reminder

	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ReminderRepository) GetByType(status string) (*models.Reminder, error) {
	var reminder *models.Reminder

	result := repo.db.Where("type = ?", status).First(&reminder)

	if result.Error != nil {
		return nil, result.Error
	}

	return reminder, nil
}

func (repo *ReminderRepository) GetAllEnabledModules() ([]*dtos.ModuleListDTO, error) {
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
		`).Joins(`JOIN reminders ON reminders.module_id = modules.id`).
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

func (repo *ReminderRepository) GetReminderByModuleID(moduleID uint) ([]*models.Reminder, error) {
	var summary []*models.Reminder
	err := repo.db.
		Preload("ReminderRules").
		Preload("ReminderRules", func(db *gorm.DB) *gorm.DB {
			return db.Order("day_offset ASC").Preload("TriggerCondition").Preload("ReminderContext")
		}).
		Where("module_id = ?", moduleID).
		Find(&summary).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}
