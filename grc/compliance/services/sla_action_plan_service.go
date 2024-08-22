package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SLAActionPlanServiceInterface interface {
	GetAll(slaID uint, params *complianceDTOs.SLAActionPlanListQueryParamsDTO) (*complianceDTOs.SLAActionPlanListResponseDTO, error)
	GetByID(slaID uint, id uint) (*models.SLAActionPlan, error)
	GetFullDetailsByID(slaID uint, id uint) (*complianceDTOs.SLAActionPlanDetailDTO, error)
	Create(slaID uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) (uint, error)
	Update(slaID uint, id uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) error
	Delete(slaID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(slaID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error
	GetSLAActionPlanCountByStatus(slaID uint) ([]complianceDTOs.SLAActionPlanStatusCountResponseDTO, error)
}

// Service and implementations
type SLAActionPlanService struct {
	repo                complianceRepositories.SLAActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	slaRegisterSvc      SLARegisterServiceInterface
	slaCheckQuestionSvc SLACheckQuestionServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
}

func NewSLAActionPlanService(
	repo complianceRepositories.SLAActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaRegisterSvc SLARegisterServiceInterface,
	slaCheckQuestionSvc SLACheckQuestionServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface) SLAActionPlanServiceInterface {
	return &SLAActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		slaRegisterSvc:      slaRegisterSvc,
		slaCheckQuestionSvc: slaCheckQuestionSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
	}
}

func (saps *SLAActionPlanService) GetAll(slaID uint, params *complianceDTOs.SLAActionPlanListQueryParamsDTO) (*complianceDTOs.SLAActionPlanListResponseDTO, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}
	return saps.repo.GetAll(slaID, params)
}

func (saps *SLAActionPlanService) GetByID(slaID uint, id uint) (*models.SLAActionPlan, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	slaActionPlan, err := saps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return slaActionPlan, nil
}

func (saps *SLAActionPlanService) GetFullDetailsByID(slaID uint, id uint) (*complianceDTOs.SLAActionPlanDetailDTO, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	var slaActionPlanDTO complianceDTOs.SLAActionPlanDetailDTO

	slaActionPlan, err := saps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(slaActionPlan, &slaActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo: add extra fields required in the details screen

	return &slaActionPlanDTO, nil
}

func (saps *SLAActionPlanService) Create(slaID uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) (uint, error) {
	tx := saps.repo.BeginTransaction()
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return 0, err
	}

	if err := saps.validateSLAActionPlanFormDTO(slaID, slaActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := saps.parseSLAActionPlanDates(slaActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	slaActionPlan := &models.SLAActionPlan{
		SLARegisterID:      slaID,
		SLACheckQuestionID: slaActionPlanFormDTO.SLACheckQuestionID,
		Title:              slaActionPlanFormDTO.Title,
		Description:        slaActionPlanFormDTO.Description,
		SeverityLevelID:    slaActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  slaActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         slaActionPlanFormDTO.Percentage,
		Comments:           slaActionPlanFormDTO.Comments,
		ActionPlanStatusID: slaActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:        currentUserID,
	}

	id, err := saps.repo.Create(slaActionPlan)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := slaActionPlanFormDTO.ActionPlanStatusID
	slaActionPlanUpdate := &models.SLAActionPlanStatusUpdate{
		SLAActionPlanID:    id,
		ActionPlanStatusID: actionPlanStatusID,
		Comment:            "Initial action plan status",
		CreatedByID:        currentUserID,
	}

	_, err = saps.repo.CreateInitialSLAActionPlanStatusUpdate(slaActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := saps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLAActionPlan{}.TableName(),
			id,
			slaActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLAActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (saps *SLAActionPlanService) Update(slaID uint, id uint, currentUserID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) error {

	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return err
	}

	existingSLAActionPlan, err := saps.GetByID(slaID, id)

	if err != nil {
		return err
	}

	if err := saps.validateSLAActionPlanFormDTO(slaID, slaActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := saps.parseSLAActionPlanDates(slaActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update SLA ActionPlan
	existingSLAActionPlan.Title = slaActionPlanFormDTO.Title
	existingSLAActionPlan.SLACheckQuestionID = slaActionPlanFormDTO.SLACheckQuestionID
	existingSLAActionPlan.Description = slaActionPlanFormDTO.Description
	existingSLAActionPlan.SeverityLevelID = slaActionPlanFormDTO.SeverityLevelID
	existingSLAActionPlan.ResponsibleUserID = slaActionPlanFormDTO.ResponsibleUserID
	existingSLAActionPlan.TargetDate = dates.TargetDate
	existingSLAActionPlan.Percentage = slaActionPlanFormDTO.Percentage
	existingSLAActionPlan.Comments = slaActionPlanFormDTO.Comments
	existingSLAActionPlan.ActionPlanStatusID = slaActionPlanFormDTO.ActionPlanStatusID
	existingSLAActionPlan.UpdatedByID = &currentUserID

	err = saps.repo.Update(existingSLAActionPlan)

	if err != nil {
		return err
	}

	return nil
}

func (saps *SLAActionPlanService) Delete(slaID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := saps.repo.BeginTransaction()
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		slaActionPlan, err := saps.GetFullDetailsByID(slaID, id)

		if err != nil {
			continue
		}

		if len(slaActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range slaActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := saps.repo.DeleteSLAActionPlanStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = saps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := saps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SLARegister{}.TableName(),
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

func (saps *SLAActionPlanService) StatusUpdate(slaID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error {
	existingSLAActionPlan, err := saps.GetByID(slaID, id)

	if err != nil {
		return err
	}

	existingSLAActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingSLAActionPlan.UpdatedByID = &currentUserID

	err = saps.repo.Update(existingSLAActionPlan, tx)

	if err != nil {
		return err
	}

	return nil
}

func (saps *SLAActionPlanService) validateSLAActionPlanFormDTO(slaID uint, slaActionPlanFormDTO *complianceDTOs.SLAActionPlanFormDTO) error {

	if slaActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := saps.severityLevelSvc.GetByID(*slaActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if slaActionPlanFormDTO.SLACheckQuestionID != nil {
		if _, err := saps.slaCheckQuestionSvc.GetQuestionByID(slaID, *slaActionPlanFormDTO.SLACheckQuestionID); err != nil {
			return utility.NewValidationError(fmt.Errorf("SLA requirement does not exist"), "not_exist", "SLA requirement")
		}
	}

	if _, err := saps.actionPlanStatusSvc.GetByID(slaActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (saps *SLAActionPlanService) parseSLAActionPlanDates(dto *complianceDTOs.SLAActionPlanFormDTO) (*complianceDTOs.SLAActionPlanParsedDates, error) {
	parsed := &complianceDTOs.SLAActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (saps *SLAActionPlanService) GetSLAActionPlanCountByStatus(slaID uint) ([]complianceDTOs.SLAActionPlanStatusCountResponseDTO, error) {
	if _, err := saps.slaRegisterSvc.GetByID(slaID); err != nil {
		return nil, err
	}

	response, err := saps.repo.GetSLAActionPlanCountByStatus(slaID)
	if err != nil {
		return nil, err
	}

	return response, nil
}
