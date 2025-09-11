package services

import (
	"errors"

	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ReminderServiceInterface interface {
	GetAll(params dtos.ReminderListQueryParamsDTO) ([]*models.Reminder, error)
	GetByID(id uint) (*models.Reminder, error)
	GetByType(status string) (*models.Reminder, error) // Renamed method

	GetAllEnabledModules() ([]*dtos.ModuleListDTO, error)
	GetAllRemindersByModuleID(moduleID uint) ([]*dtos.GroupedReminderDTO, error)
}

// Service and implementations
type ReminderService struct {
	repo repositories.ReminderRepositoryInterface
}

func NewReminderService(repo repositories.ReminderRepositoryInterface) ReminderServiceInterface {
	return &ReminderService{
		repo: repo,
	}
}

func (rs *ReminderService) GetAll(params dtos.ReminderListQueryParamsDTO) ([]*models.Reminder, error) {
	return rs.repo.GetAll(params)
}

func (rs ReminderService) GetByID(id uint) (*models.Reminder, error) {
	reminder, err := rs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return reminder, nil
}

func (rs *ReminderService) GetByType(status string) (*models.Reminder, error) {
	statusResult, err := rs.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}

func (rs *ReminderService) GetAllEnabledModules() ([]*dtos.ModuleListDTO, error) {
	return rs.repo.GetAllEnabledModules()
}

func (rs *ReminderService) GetAllRemindersByModuleID(moduleID uint) ([]*dtos.GroupedReminderDTO, error) {
	reminders, err := rs.repo.GetReminderByModuleID(moduleID)
	if err != nil {
		return nil, err
	}

	reminderDTOs := []*dtos.GroupedReminderDTO{}
	for _, reminder := range reminders {

		reminderDTO := &dtos.GroupedReminderDTO{
			ID:            reminder.ID,
			TitleEn:       reminder.TitleEn,
			TitleAr:       reminder.TitleAr,
			Type:          reminder.Type,
			ReminderRules: []dtos.ReminderRuleDetailDTO{},
		}

		for _, rule := range reminder.ReminderRules {
			reminderDTO.ReminderRules = append(reminderDTO.ReminderRules, dtos.ReminderRuleDetailDTO{
				ID:                 rule.ID,
				TitleEn:            rule.TitleEn,
				TitleAr:            rule.TitleAr,
				NotificationMethod: rule.NotificationMethod,
				Type:               rule.Type,
				DayOffset:          rule.DayOffset,
				TriggerConditionID: rule.TriggerConditionID,
				TriggerCondition: dtos.TriggerConditionDTO{
					ID:      rule.TriggerCondition.ID,
					TitleEn: rule.TriggerCondition.TitleEn,
					TitleAr: rule.TriggerCondition.TitleAr,
					Type:    rule.TriggerCondition.Type,
				},
				ReminderContextID: rule.ReminderContextID,
				ReminderContext: dtos.ReminderContextDTO{
					ID:       rule.ReminderContext.ID,
					ModuleID: rule.ReminderContext.ModuleID,
					TitleEn:  rule.ReminderContext.TitleEn,
					TitleAr:  rule.ReminderContext.TitleAr,
					Type:     rule.ReminderContext.Type,
					Status:   rule.ReminderContext.Status,
				},
				CreatedAt: rule.CreatedAt,
			})
		}

		reminderDTOs = append(reminderDTOs, reminderDTO)
	}
	return reminderDTOs, err
}
