package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"os"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskActionPlanServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskActionPlan, error)
	GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskActionPlanDetailDTO, error)
	Create(riskID uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, risk *riskDTOs.RiskActionPlanFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(riskID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error

	ExportRiskActionPlan(data []*riskDTOs.RiskActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type RiskActionPlanService struct {
	repo                riskRepositories.RiskActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	riskRegisterSvc     RiskRegisterServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewRiskActionPlanService(
	repo riskRepositories.RiskActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface) RiskActionPlanServiceInterface {
	return &RiskActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		riskRegisterSvc:     riskRegisterSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (raps *RiskActionPlanService) GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error) {
	if riskID != 0 {
		if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
			return nil, err
		}
	}
	return raps.repo.GetAll(riskID, params)
}

func (raps RiskActionPlanService) GetByID(riskID uint, id uint) (*models.RiskActionPlan, error) {
	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskActionPlan, err := raps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlan, nil
}

func (raps RiskActionPlanService) GetFullDetailsByID(riskID uint, id uint) (*riskDTOs.RiskActionPlanDetailDTO, error) {
	var riskActionPlanDTO riskDTOs.RiskActionPlanDetailDTO

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskActionPlan, err := raps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskActionPlan, &riskActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &riskActionPlanDTO, nil
}

func (raps *RiskActionPlanService) Create(riskID uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) (uint, error) {
	tx := raps.repo.BeginTransaction()

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	if err := raps.validateRiskActionPlanFormDTO(RiskActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := raps.parseRiskActionPlanDates(RiskActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	riskActionPlan := &models.RiskActionPlan{
		RiskRegisterID:     riskID,
		Title:              RiskActionPlanFormDTO.Title,
		Description:        RiskActionPlanFormDTO.Description,
		SeverityLevelID:    RiskActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  RiskActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         RiskActionPlanFormDTO.Percentage,
		Comments:           RiskActionPlanFormDTO.Comments,
		ActionPlanStatusID: RiskActionPlanFormDTO.ActionPlanStatusID,
		CreatedBYID:        currentUserID,
	}

	id, err := raps.repo.Create(riskActionPlan)
	if err != nil {
		return 0, err
	}
	actionPlanStatusID := RiskActionPlanFormDTO.ActionPlanStatusID
	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   id,
		ActionPlanStatusID: actionPlanStatusID,
		Comment:            "Initial action plan status",
		CreatedByID:        currentUserID,
	}
	_, err = raps.repo.CreateInitialRiskActionPlanStatusUpdate(riskActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := raps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlan{}.TableName(),
			id,
			RiskActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (raps *RiskActionPlanService) Update(riskID uint, id uint, currentUserID uint, RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) error {
	tx := raps.repo.BeginTransaction()

	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingRiskActionPlan, err := raps.GetByID(riskID, id)

	if err != nil {
		return err
	}

	if err := raps.validateRiskActionPlanFormDTO(RiskActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := raps.parseRiskActionPlanDates(RiskActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Risk ActionPlan
	existingRiskActionPlan.Title = RiskActionPlanFormDTO.Title
	existingRiskActionPlan.Description = RiskActionPlanFormDTO.Description
	existingRiskActionPlan.SeverityLevelID = RiskActionPlanFormDTO.SeverityLevelID
	existingRiskActionPlan.ResponsibleUserID = RiskActionPlanFormDTO.ResponsibleUserID
	existingRiskActionPlan.TargetDate = dates.TargetDate
	existingRiskActionPlan.Percentage = RiskActionPlanFormDTO.Percentage
	existingRiskActionPlan.Comments = RiskActionPlanFormDTO.Comments
	existingRiskActionPlan.ActionPlanStatusID = RiskActionPlanFormDTO.ActionPlanStatusID
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = raps.repo.Update(existingRiskActionPlan, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (raps *RiskActionPlanService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := raps.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		riskActionPlan, err := raps.repo.GetFullDetailsByID(id)

		if err != nil {
			continue
		}
		if len(riskActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range riskActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := raps.repo.DeleteRiskActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = raps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := raps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskRegister{}.TableName(),
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

func (raps *RiskActionPlanService) StatusUpdate(riskID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error {

	existingRiskActionPlan, err := raps.GetByID(riskID, id)

	if err != nil {
		return err
	}

	existingRiskActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingRiskActionPlan.UpdatedBYID = &currentUserID

	err = raps.repo.Update(existingRiskActionPlan, tx)

	if err != nil {
		return err
	}

	return nil
}

func (raps *RiskActionPlanService) validateRiskActionPlanFormDTO(RiskActionPlanFormDTO *riskDTOs.RiskActionPlanFormDTO) error {
	if RiskActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := raps.severityLevelSvc.GetByID(*RiskActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := raps.actionPlanStatusSvc.GetByID(RiskActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (raps *RiskActionPlanService) parseRiskActionPlanDates(dto *riskDTOs.RiskActionPlanFormDTO) (*riskDTOs.RiskActionPlanParsedDates, error) {
	parsed := &riskDTOs.RiskActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (raps *RiskActionPlanService) ExportRiskActionPlan(data []*riskDTOs.RiskActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := raps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Risk",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
			"Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Action Plan",
			"Incident",
			"Description",
			"Severity",
			"Responsible User",
			"Target Date",
			"Status",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, actionPlan := range data {
		var responsibleUser, severity, status string

		if userSetting.Language == "ar" {
			responsibleUser = actionPlan.ResponsibleUserFirstNameAr
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleAr)
			status = actionPlan.ActionPlanStatusTitleAr
		} else {
			responsibleUser = actionPlan.ResponsibleUserFirstNameEn
			severity = utility.GetStringValue(actionPlan.SeverityLevelTitleEn)
			status = actionPlan.ActionPlanStatusTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			actionPlan.Title,
			actionPlan.RiskTitle,
			utility.GetStringValue(actionPlan.Description),
			severity,
			responsibleUser,
			utility.GetStringValue(actionPlan.TargetDate),
			status,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
