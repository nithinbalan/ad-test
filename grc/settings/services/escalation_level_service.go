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

type EscalationLevelServiceInterface interface {
	GetAll(params *settingDTOs.EscalationLevelListQueryParamsDTO) ([]*settingDTOs.EscalationLevelListDTO, error)
	GetFullDetailsByID(id uint) (*settingDTOs.EscalationLevelDetailDTO, error)
	GetByID(id uint) (*models.EscalationLevel, error)
	Create(currentUserID uint, escalationRuleFormDTO *settingDTOs.EscalationLevelFormDTO) (uint, error)
	Update(id uint, currentUserID uint, escalationRuleFormDTO *settingDTOs.EscalationLevelFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type EscalationLevelService struct {
	repo           repositories.EscalationLevelRepositoryInterface
	escalationSvc  EscalationServiceInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewEscalationLevelService(
	repo repositories.EscalationLevelRepositoryInterface,
	escalationSvc EscalationServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
) EscalationLevelServiceInterface {
	return &EscalationLevelService{
		repo:           repo,
		escalationSvc:  escalationSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (rrs *EscalationLevelService) GetAll(params *settingDTOs.EscalationLevelListQueryParamsDTO) ([]*settingDTOs.EscalationLevelListDTO, error) {
	return rrs.repo.GetAll(params)
}

func (rrs *EscalationLevelService) GetByID(id uint) (*models.EscalationLevel, error) {
	escalationRule, err := rrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return escalationRule, nil
}

func (rrs *EscalationLevelService) GetFullDetailsByID(id uint) (*settingDTOs.EscalationLevelDetailDTO, error) {
	var escalationRuleDetailDTO settingDTOs.EscalationLevelDetailDTO

	escalationRule, err := rrs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(escalationRule, &escalationRuleDetailDTO)
	if err != nil {
		return nil, err
	}

	return &escalationRuleDetailDTO, nil
}

func (rrs *EscalationLevelService) Create(currentUserID uint, dto *settingDTOs.EscalationLevelFormDTO) (uint, error) {

	if err := rrs.validateEscalationLevelFormDTO(dto); err != nil {
		return 0, err
	}

	level := &models.EscalationLevel{
		EscalationID:       dto.EscalationID,
		EscalationTo:       dto.EscalationTo,
		NotificationMethod: dto.NotificationMethod,
		DayOffset:          dto.DayOffset,
		CreatedBYID:        currentUserID,
	}

	id, err := rrs.repo.Create(level)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rrs *EscalationLevelService) Update(id uint, currentUserID uint, dto *settingDTOs.EscalationLevelFormDTO) error {

	oldData, err := rrs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	existing, err := rrs.GetByID(id)
	if err != nil {
		return err
	}

	if err := rrs.validateEscalationLevelFormDTO(dto); err != nil {
		return err
	}

	existing.EscalationID = dto.EscalationID
	existing.EscalationTo = dto.EscalationTo
	existing.NotificationMethod = dto.NotificationMethod
	existing.DayOffset = dto.DayOffset
	existing.UpdatedBYID = &currentUserID

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
			"escalation":          oldData.Escalation.TitleEn,
			"escalation_to":       oldData.EscalationTo,
			"notification_method": oldData.NotificationMethod,
			"day_offset":          oldData.DayOffset,
		}

		newFields := map[string]interface{}{
			"escalation":          newData.Escalation.TitleEn,
			"escalation_tp":       newData.EscalationTo,
			"notification_method": newData.NotificationMethod,
			"day_offset":          newData.DayOffset,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)

		if len(changedData) > 0 {
			err := rrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.EscalationLevel{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log escalation level update", map[string]interface{}{
					"action":    models.UpdateAction,
					"table":     models.EscalationLevel{}.TableName(),
					"record_id": id,
					"user_id":   currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (rrs *EscalationLevelService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

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
				"title_en": item.Escalation.TitleEn,
			})
		}
	}

	if len(successfulIds) > 0 {
		go rrs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.EscalationLevel{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (rrs *EscalationLevelService) validateEscalationLevelFormDTO(escalationRuleFormDTO *settingDTOs.EscalationLevelFormDTO) error {
	// Validate escalation
	if _, err := rrs.escalationSvc.GetByID(escalationRuleFormDTO.EscalationID); err != nil {
		return utility.NewValidationError(fmt.Errorf("escalation does not exist"), "not_exist", "escalation")
	}

	return nil
}
