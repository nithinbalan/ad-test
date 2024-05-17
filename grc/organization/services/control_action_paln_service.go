package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlActionPlanServiceInterface interface {
	GetAll(controlID uint, params *organizationDTOs.ControlActionPlanListQueryParamsDTO) (*organizationDTOs.ControlActionPlanListResponseDTO, error)
	GetByID(controlID uint, id uint) (*models.ControlActionPlan, error)
	GetFullDetailsByID(controlID uint, id uint) (*organizationDTOs.ControlActionPlanDetailDTO, error)
	Create(controlID uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) (uint, error)
	Update(controlID uint, id uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) error
	Delete(controlID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ControlActionPlanService struct {
	repo                organizationRepositories.ControlActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	controlSvc          ControlServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             UserServiceInterface
}

func NewControlActionPlanService(
	repo organizationRepositories.ControlActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlSvc ControlServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc UserServiceInterface) ControlActionPlanServiceInterface {
	return &ControlActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		controlSvc:          controlSvc,
		severityLevelSvc:    severityLevelSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
	}
}

func (caps *ControlActionPlanService) GetAll(controlID uint, params *organizationDTOs.ControlActionPlanListQueryParamsDTO) (*organizationDTOs.ControlActionPlanListResponseDTO, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return caps.repo.GetAll(controlID, params)
}

func (caps ControlActionPlanService) GetByID(controlID uint, id uint) (*models.ControlActionPlan, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	controlActionPlan, err := caps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return controlActionPlan, nil
}

func (caps ControlActionPlanService) GetFullDetailsByID(controlID uint, id uint) (*organizationDTOs.ControlActionPlanDetailDTO, error) {
	var controlActionPlanDTO organizationDTOs.ControlActionPlanDetailDTO

	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	controlActionPlan, err := caps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(controlActionPlan, &controlActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &controlActionPlanDTO, nil
}

func (caps *ControlActionPlanService) Create(controlID uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) (uint, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return 0, err
	}

	if err := caps.validateControlActionPlanFormDTO(controlID, ControlActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caps.parseControlActionPlanDates(ControlActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	controlActionPlan := &models.ControlActionPlan{
		ControlID:          controlID,
		Title:              ControlActionPlanFormDTO.Title,
		Description:        ControlActionPlanFormDTO.Description,
		SeverityLevelID:    ControlActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  ControlActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         *ControlActionPlanFormDTO.Percentage,
		ActionPlanStatusID: ControlActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:        currentUserID,
	}

	id, err := caps.repo.Create(controlActionPlan)
	if err != nil {
		return 0, err
	}

	go func() {
		err := caps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlActionPlan{}.TableName(),
			id,
			ControlActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (caps *ControlActionPlanService) Update(controlID uint, id uint, currentUserID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) error {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return err
	}

	existingControlActionPlan, err := caps.GetByID(controlID, id)

	if err != nil {
		return err
	}

	if err := caps.validateControlActionPlanFormDTO(controlID, ControlActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caps.parseControlActionPlanDates(ControlActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingControlActionPlan.Title = ControlActionPlanFormDTO.Title
	existingControlActionPlan.Description = ControlActionPlanFormDTO.Description
	existingControlActionPlan.SeverityLevelID = ControlActionPlanFormDTO.SeverityLevelID
	existingControlActionPlan.ResponsibleUserID = ControlActionPlanFormDTO.ResponsibleUserID
	existingControlActionPlan.TargetDate = dates.TargetDate
	existingControlActionPlan.Percentage = *ControlActionPlanFormDTO.Percentage
	existingControlActionPlan.ActionPlanStatusID = ControlActionPlanFormDTO.ActionPlanStatusID
	existingControlActionPlan.UpdatedByID = &currentUserID

	err = caps.repo.Update(existingControlActionPlan)

	if err != nil {
		return err
	}

	return nil
}

func (caps *ControlActionPlanService) Delete(controlID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caps.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := caps.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = caps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ControlActionPlan{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (caps *ControlActionPlanService) validateControlActionPlanFormDTO(controlID uint, ControlActionPlanFormDTO *organizationDTOs.ControlActionPlanFormDTO) error {
	if ControlActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := caps.severityLevelSvc.GetByID(*ControlActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := caps.userSvc.GetByID(uint(ControlActionPlanFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := caps.actionPlanStatusSvc.GetByID(ControlActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (caps *ControlActionPlanService) parseControlActionPlanDates(dto *organizationDTOs.ControlActionPlanFormDTO) (*organizationDTOs.ControlActionPlanParsedDates, error) {
	parsed := &organizationDTOs.ControlActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
