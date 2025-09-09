package services

import (
	"errors"
	"fmt"

	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	settingDTOs "jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ReminderRuleServiceInterface interface {
	GetAll(params *settingDTOs.ReminderRuleListQueryParamsDTO) ([]*settingDTOs.ReminderRuleListDTO, error)
	GetFullDetailsByID(id uint) (*settingDTOs.ReminderRuleDetailDTO, error)
	GetByID(id uint) (*models.ReminderRule, error)
	Create(currentUserID uint, reminderRuleFormDTO *settingDTOs.ReminderRuleFormDTO) (uint, error)
	Update(id uint, currentUserID uint, reminderRuleFormDTO *settingDTOs.ReminderRuleFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type ReminderRuleService struct {
	repo                repositories.ReminderRuleRepositoryInterface
	reminderSvc         ReminderServiceInterface
	triggerConditionSvc TriggerConditionServiceInterface
	activityLogSvc      activityLogServices.ActivityLogServiceInterface
	reminderContextSvc  ReminderContextServiceInterface
}

func NewReminderRuleService(
	repo repositories.ReminderRuleRepositoryInterface,
	reminderSvc ReminderServiceInterface,
	triggerConditionSvc TriggerConditionServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	reminderContextSvc ReminderContextServiceInterface,
) ReminderRuleServiceInterface {
	return &ReminderRuleService{
		repo:                repo,
		reminderSvc:         reminderSvc,
		triggerConditionSvc: triggerConditionSvc,
		activityLogSvc:      activityLogSvc,
		reminderContextSvc:  reminderContextSvc,
	}
}

func (rrs *ReminderRuleService) GetAll(params *settingDTOs.ReminderRuleListQueryParamsDTO) ([]*settingDTOs.ReminderRuleListDTO, error) {
	return rrs.repo.GetAll(params)
}

func (rrs *ReminderRuleService) GetByID(id uint) (*models.ReminderRule, error) {
	reminderRule, err := rrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return reminderRule, nil
}

func (rrs *ReminderRuleService) GetFullDetailsByID(id uint) (*settingDTOs.ReminderRuleDetailDTO, error) {
	var reminderRuleDetailDTO settingDTOs.ReminderRuleDetailDTO

	reminderRule, err := rrs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(reminderRule, &reminderRuleDetailDTO)
	if err != nil {
		return nil, err
	}

	return &reminderRuleDetailDTO, nil
}

func (rrs *ReminderRuleService) Create(currentUserID uint, dto *settingDTOs.ReminderRuleFormDTO) (uint, error) {
	if err := rrs.validateReminderRuleFormDTO(dto); err != nil {
		return 0, err
	}

	rule := &models.ReminderRule{
		TitleEn:            dto.TitleEn,
		TitleAr:            dto.TitleAr,
		ReminderID:         dto.ReminderID,
		TriggerConditionID: dto.TriggerConditionID,
		NotificationMethod: dto.NotificationMethod,
		DayOffset:          dto.DayOffset,
		Type:               dto.Type,
		ReminderContextID:  dto.ReminderContextID,
		CreatedBYID:        currentUserID,
	}

	id, err := rrs.repo.Create(rule)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rrs *ReminderRuleService) Update(id uint, currentUserID uint, dto *settingDTOs.ReminderRuleFormDTO) error {
	oldData, err := rrs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	existing, err := rrs.GetByID(id)
	if err != nil {
		return err
	}

	if err := rrs.validateReminderRuleFormDTO(dto); err != nil {
		return err
	}

	existing.ReminderID = dto.ReminderID
	existing.TriggerConditionID = dto.TriggerConditionID
	existing.NotificationMethod = dto.NotificationMethod
	existing.TitleEn = dto.TitleEn
	existing.TitleAr = dto.TitleAr
	existing.DayOffset = dto.DayOffset
	existing.UpdatedBYID = &currentUserID
	existing.ReminderContextID= dto.ReminderContextID

	err = rrs.repo.Update(existing)
	if err != nil {
		return err
	}

	go func() {
		// Fetch updated full details after update
		newData, err := rrs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		oldFields := map[string]interface{}{
			"reminder":            oldData.Reminder.TitleEn,
			"trigger_condition":   oldData.TriggerCondition.TitleEn,
			"title_en":            oldData.TitleEn,
			"title_ar":            oldData.TitleAr,
			"notification_method": oldData.NotificationMethod,
			"day_offset":          oldData.DayOffset,
		}

		newFields := map[string]interface{}{
			"reminder":            newData.Reminder.TitleEn,
			"trigger_condition":   newData.TriggerCondition.TitleEn,
			"title_en":            newData.TitleEn,
			"title_ar":            newData.TitleAr,
			"notification_method": newData.NotificationMethod,
			"day_offset":          newData.DayOffset,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)

		if len(changedData) > 0 {
			err := rrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ReminderRule{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log reminder rule update", map[string]interface{}{
					"action":    models.UpdateAction,
					"table":     models.ReminderRule{}.TableName(),
					"record_id": id,
					"user_id":   currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (rrs *ReminderRuleService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := rrs.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := rrs.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedItems = append(deletedItems, map[string]interface{}{
				"id":       id,
				"title_en": item.TitleEn,
			})
		}
	}

	if len(successfulIds) > 0 {
		go rrs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ReminderRule{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (rrs *ReminderRuleService) validateReminderRuleFormDTO(reminderRuleFormDTO *settingDTOs.ReminderRuleFormDTO) error {
	// Validate reminder
	if _, err := rrs.reminderSvc.GetByID(reminderRuleFormDTO.ReminderID); err != nil {
		return utility.NewValidationError(fmt.Errorf("reminder does not exist"), "not_exist", "process dependency sub type")
	}

	if _, err := rrs.triggerConditionSvc.GetByID(reminderRuleFormDTO.TriggerConditionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("trigger condition does not exist"), "not_exist", "process dependency sub type")
	}

	if _, err := rrs.reminderContextSvc.GetByID(reminderRuleFormDTO.ReminderContextID); err != nil {
		return utility.NewValidationError(fmt.Errorf("reminder context does not exist"), "not_exist", "process dependency sub type")
	}

	return nil
}
