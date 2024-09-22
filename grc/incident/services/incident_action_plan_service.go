package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentActionPlanServiceInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentActionPlanListQueryParamsDTO) (*incidentDTOs.IncidentActionPlanListResponseDTO, error)
	GetByID(incidentID uint, id uint) (*models.IncidentActionPlan, error)
	GetFullDetailsByID(incidentID uint, id uint) (*incidentDTOs.IncidentActionPlanDetailDTO, error)
	Create(incidentID uint, currentUserID uint, IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) (uint, error)
	Update(incidentID uint, id uint, currentUserID uint, incident *incidentDTOs.IncidentActionPlanFormDTO) error
	Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(incidentID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error

	ExportIncidentActionPlan(data []*incidentDTOs.IncidentActionPlanListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type IncidentActionPlanService struct {
	repo                incidentRepositories.IncidentActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	incidentTimelineSvc IncidentTimelineServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewIncidentActionPlanService(
	repo incidentRepositories.IncidentActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	incidentTimelineSvc IncidentTimelineServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) IncidentActionPlanServiceInterface {
	return &IncidentActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentRegisterSvc: incidentRegisterSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		incidentTimelineSvc: incidentTimelineSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (iaps *IncidentActionPlanService) GetAll(incidentID uint, params *incidentDTOs.IncidentActionPlanListQueryParamsDTO) (*incidentDTOs.IncidentActionPlanListResponseDTO, error) {
	if incidentID != 0 {
		if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
			return nil, err
		}
	}
	return iaps.repo.GetAll(incidentID, params)
}

func (iaps IncidentActionPlanService) GetByID(incidentID uint, id uint) (*models.IncidentActionPlan, error) {
	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentActionPlan, err := iaps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentActionPlan, nil
}

func (iaps IncidentActionPlanService) GetFullDetailsByID(incidentID uint, id uint) (*incidentDTOs.IncidentActionPlanDetailDTO, error) {
	var incidentActionPlanDTO incidentDTOs.IncidentActionPlanDetailDTO

	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentActionPlan, err := iaps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(incidentActionPlan, &incidentActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen
	if incidentActionPlan.TargetDate != nil {
		targetDate := incidentActionPlan.TargetDate.Format("2006-01-02")
		incidentActionPlanDTO.TargetDateString = &targetDate
	}

	incidentActionPlanDTO.CreatedAt = incidentActionPlan.CreatedAt
	incidentActionPlanDTO.UpdatedAt = incidentActionPlan.UpdatedAt

	return &incidentActionPlanDTO, nil
}

func (iaps *IncidentActionPlanService) Create(incidentID uint, currentUserID uint, IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) (uint, error) {
	tx := iaps.repo.BeginTransaction()
	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}

	if err := iaps.validateIncidentActionPlanFormDTO(IncidentActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := iaps.parseIncidentActionPlanDates(IncidentActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	incidentActionPlan := &models.IncidentActionPlan{
		IncidentRegisterID: incidentID,
		Title:              IncidentActionPlanFormDTO.Title,
		Description:        IncidentActionPlanFormDTO.Description,
		SeverityLevelID:    IncidentActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  IncidentActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         IncidentActionPlanFormDTO.Percentage,
		Comments:           IncidentActionPlanFormDTO.Comments,
		ActionPlanStatusID: IncidentActionPlanFormDTO.ActionPlanStatusID,
		CreatedBYID:        currentUserID,
	}

	id, err := iaps.repo.Create(incidentActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	actionPlanStatusID := IncidentActionPlanFormDTO.ActionPlanStatusID
	incidentActionPlanUpdate := &models.IncidentActionPlanStatusUpdate{
		IncidentActionPlanID: id,
		ActionPlanStatusID:   actionPlanStatusID,
		Comment:              "Initial action plan status",
		CreatedByID:          currentUserID,
	}
	_, err = iaps.repo.CreateInitialIncidentActionPlanStatusUpdate(incidentActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//Update timeline
	//todo
	actinPlanCount, err := iaps.repo.GetCountByRiskID(incidentID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if actinPlanCount == 1 {
		timelineType, err := iaps.incidentTimelineSvc.GetByType("action-plans")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		timelineStatus, err := iaps.incidentTimelineSvc.GetByStatusType("completed")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		currentTime := time.Now()
		timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
			PerformedByID:    &currentUserID,
			TimelineStatusID: timelineStatus.ID,
			PerformedAt:      &currentTime,
		}
		if err := iaps.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := iaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentActionPlan{}.TableName(),
			id,
			IncidentActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (iaps *IncidentActionPlanService) Update(incidentID uint, id uint, currentUserID uint, IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) error {
	tx := iaps.repo.BeginTransaction()

	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	existingIncidentActionPlan, err := iaps.GetByID(incidentID, id)

	if err != nil {
		return err
	}

	if err := iaps.validateIncidentActionPlanFormDTO(IncidentActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := iaps.parseIncidentActionPlanDates(IncidentActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Incident ActionPlan
	existingIncidentActionPlan.Title = IncidentActionPlanFormDTO.Title
	existingIncidentActionPlan.Description = IncidentActionPlanFormDTO.Description
	existingIncidentActionPlan.SeverityLevelID = IncidentActionPlanFormDTO.SeverityLevelID
	existingIncidentActionPlan.ResponsibleUserID = IncidentActionPlanFormDTO.ResponsibleUserID
	existingIncidentActionPlan.TargetDate = dates.TargetDate
	existingIncidentActionPlan.Percentage = IncidentActionPlanFormDTO.Percentage
	existingIncidentActionPlan.Comments = IncidentActionPlanFormDTO.Comments
	existingIncidentActionPlan.ActionPlanStatusID = IncidentActionPlanFormDTO.ActionPlanStatusID
	existingIncidentActionPlan.UpdatedBYID = &currentUserID

	err = iaps.repo.Update(existingIncidentActionPlan, tx)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (iaps *IncidentActionPlanService) Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := iaps.repo.BeginTransaction()

	if _, err := iaps.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		incidentActionPlan, err := iaps.repo.GetFullDetailsByID(id)

		if err != nil {
			continue
		}
		if len(incidentActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range incidentActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := iaps.repo.DeleteIncidentActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}
		//Delete-> update timeline performedAt as null
		//todo
		actinPlanCount, err := iaps.repo.GetCountByRiskID(incidentID, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if actinPlanCount == 1 {
			timelineType, err := iaps.incidentTimelineSvc.GetByType("action-plans")
			if err != nil {
				tx.Rollback()
				return failedIds, err
			}
			timelineStatus, err := iaps.incidentTimelineSvc.GetByStatusType("pending")
			if err != nil {
				tx.Rollback()
				return failedIds, err
			}
			timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
				PerformedByID:    nil,
				TimelineStatusID: timelineStatus.ID,
				PerformedAt:      nil,
			}
			if err := iaps.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
				tx.Rollback()
				return failedIds, err
			}
		}
		err = iaps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := iaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentRegister{}.TableName(),
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

func (iaps *IncidentActionPlanService) StatusUpdate(incidentID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error {

	existingIncidentActionPlan, err := iaps.GetByID(incidentID, id)

	if err != nil {
		return err
	}

	existingIncidentActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingIncidentActionPlan.UpdatedBYID = &currentUserID

	err = iaps.repo.Update(existingIncidentActionPlan, tx)

	if err != nil {
		return err
	}

	return nil
}

func (iaps *IncidentActionPlanService) validateIncidentActionPlanFormDTO(IncidentActionPlanFormDTO *incidentDTOs.IncidentActionPlanFormDTO) error {
	if IncidentActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := iaps.severityLevelSvc.GetByID(*IncidentActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := iaps.actionPlanStatusSvc.GetByID(IncidentActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (iaps *IncidentActionPlanService) parseIncidentActionPlanDates(dto *incidentDTOs.IncidentActionPlanFormDTO) (*incidentDTOs.IncidentActionPlanParsedDates, error) {
	parsed := &incidentDTOs.IncidentActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (iaps *IncidentActionPlanService) ExportIncidentActionPlan(data []*incidentDTOs.IncidentActionPlanListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := iaps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
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
			actionPlan.IncidentTitle,
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
