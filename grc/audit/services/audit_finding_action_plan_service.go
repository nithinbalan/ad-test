package services

import (
	"errors"
	"fmt"
	auditFindingDTOs "jethur/grc/audit/dtos"
	auditFindingRepositories "jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditFindingActionPlanServiceInterface interface {
	GetAll(auditID uint, findingID uint, params *auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditFindingActionPlanListResponseDTO, error)
	GetByID(auditID uint, findingID uint, id uint) (*models.AuditFindingActionPlan, error)
	GetByFindingID(findingID uint, id uint) (*models.AuditFindingActionPlan, error)
	GetFullDetailsByID(auditID uint, findingID uint, id uint) (*auditFindingDTOs.AuditFindingActionPlanDetailDTO, error)
	Create(auditID uint, findingID uint, currentUserID uint, AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) (uint, error)
	Update(auditID uint, findingID uint, id uint, currentUserID uint, auditFinding *auditFindingDTOs.AuditFindingActionPlanFormDTO) error
	Delete(auditID uint, findingID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(findingID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error
}

// Service and implementations
type AuditFindingActionPlanService struct {
	repo                auditFindingRepositories.AuditFindingActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	auditFindingSvc     AuditFindingServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	auditTimelineSvc    AuditTimelineServiceInterface
}

func NewAuditFindingActionPlanService(
	repo auditFindingRepositories.AuditFindingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
) AuditFindingActionPlanServiceInterface {
	return &AuditFindingActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		auditFindingSvc:     auditFindingSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		severityLevelSvc:    severityLevelSvc,
		auditTimelineSvc:    auditTimelineSvc,
	}
}

func (afaps *AuditFindingActionPlanService) GetAll(auditID uint, findingID uint, params *auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditFindingActionPlanListResponseDTO, error) {
	if auditID != 0 && findingID != 0 {
		if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
			return nil, err
		}
	}
	return afaps.repo.GetAll(findingID, params)
}

func (afaps AuditFindingActionPlanService) GetByID(auditID uint, findingID uint, id uint) (*models.AuditFindingActionPlan, error) {
	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditFindingActionPlan, err := afaps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingActionPlan, nil
}

func (afaps AuditFindingActionPlanService) GetByFindingID(findingID uint, id uint) (*models.AuditFindingActionPlan, error) {
	auditFindingActionPlan, err := afaps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingActionPlan, nil
}

func (afaps AuditFindingActionPlanService) GetFullDetailsByID(auditID uint, findingID uint, id uint) (*auditFindingDTOs.AuditFindingActionPlanDetailDTO, error) {
	var auditFindingActionPlanDTO auditFindingDTOs.AuditFindingActionPlanDetailDTO

	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditFindingActionPlan, err := afaps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFindingActionPlan, &auditFindingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	if auditFindingActionPlan.TargetDate != nil {
		targetDate := auditFindingActionPlan.TargetDate.Format("2006-01-02")
		auditFindingActionPlanDTO.TargetDateString = &targetDate
	}

	auditFindingActionPlanDTO.CreatedAt = auditFindingActionPlan.CreatedAt
	auditFindingActionPlanDTO.UpdatedAt = auditFindingActionPlan.UpdatedAt

	return &auditFindingActionPlanDTO, nil
}

func (afaps *AuditFindingActionPlanService) Create(auditID uint, findingID uint, currentUserID uint, AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) (uint, error) {
	tx := afaps.repo.BeginTransaction()
	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}

	if err := afaps.validateAuditFindingActionPlanFormDTO(AuditFindingActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := afaps.parseAuditFindingActionPlanDates(AuditFindingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	auditFindingActionPlan := &models.AuditFindingActionPlan{
		AuditFindingID:     findingID,
		Title:              AuditFindingActionPlanFormDTO.Title,
		Description:        AuditFindingActionPlanFormDTO.Description,
		SeverityLevelID:    AuditFindingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  AuditFindingActionPlanFormDTO.ResponsibleUserID,
		TargetDate:         dates.TargetDate,
		Percentage:         AuditFindingActionPlanFormDTO.Percentage,
		Comments:           AuditFindingActionPlanFormDTO.Comments,
		ActionPlanStatusID: AuditFindingActionPlanFormDTO.ActionPlanStatusID,
		CreatedBYID:        currentUserID,
	}

	id, err := afaps.repo.Create(auditFindingActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := AuditFindingActionPlanFormDTO.ActionPlanStatusID
	auditFindingActionPlanUpdate := &models.AuditFindingActionPlanStatusUpdate{
		AuditFindingActionPlanID: id,
		ActionPlanStatusID:       actionPlanStatusID,
		Comment:                  "Initial action plan status",
		CreatedByID:              currentUserID,
	}

	_, err = afaps.repo.CreateInitialAuditFindingActionPlanStatusUpdate(auditFindingActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := afaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingActionPlan{}.TableName(),
			id,
			AuditFindingActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afaps *AuditFindingActionPlanService) Update(auditID uint, findingID uint, id uint, currentUserID uint, AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) error {
	tx := afaps.repo.BeginTransaction()

	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	existingAuditFindingActionPlan, err := afaps.GetByID(auditID, findingID, id)

	if err != nil {
		return err
	}

	if err := afaps.validateAuditFindingActionPlanFormDTO(AuditFindingActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := afaps.parseAuditFindingActionPlanDates(AuditFindingActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Audit Finding ActionPlan
	existingAuditFindingActionPlan.Title = AuditFindingActionPlanFormDTO.Title
	existingAuditFindingActionPlan.Description = AuditFindingActionPlanFormDTO.Description
	existingAuditFindingActionPlan.SeverityLevelID = AuditFindingActionPlanFormDTO.SeverityLevelID
	existingAuditFindingActionPlan.ResponsibleUserID = AuditFindingActionPlanFormDTO.ResponsibleUserID
	existingAuditFindingActionPlan.TargetDate = dates.TargetDate
	existingAuditFindingActionPlan.Percentage = AuditFindingActionPlanFormDTO.Percentage
	existingAuditFindingActionPlan.Comments = AuditFindingActionPlanFormDTO.Comments
	existingAuditFindingActionPlan.ActionPlanStatusID = AuditFindingActionPlanFormDTO.ActionPlanStatusID
	existingAuditFindingActionPlan.UpdatedBYID = &currentUserID

	err = afaps.repo.Update(existingAuditFindingActionPlan, tx)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (afaps *AuditFindingActionPlanService) Delete(auditID uint, findingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := afaps.repo.BeginTransaction()

	if _, err := afaps.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		auditFindingActionPlan, err := afaps.GetFullDetailsByID(auditID, findingID, id)

		if err != nil {
			continue
		}
		if len(auditFindingActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range auditFindingActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := afaps.repo.DeleteAuditFindingActionPlanStatusUpdates(deleteStatusUpdateIds); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}
		err = afaps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := afaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
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

func (afaps *AuditFindingActionPlanService) StatusUpdate(findingID uint, id uint, currentUserID uint, ActionPlanStatusID uint, tx utility.DatabaseTransaction) error {

	existingAuditFindingActionPlan, err := afaps.GetByFindingID(findingID, id)

	if err != nil {
		return err
	}

	existingAuditFindingActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingAuditFindingActionPlan.UpdatedBYID = &currentUserID

	err = afaps.repo.Update(existingAuditFindingActionPlan, tx)

	if err != nil {
		return err
	}

	return nil
}

func (afaps *AuditFindingActionPlanService) validateAuditFindingActionPlanFormDTO(AuditFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) error {
	if AuditFindingActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := afaps.severityLevelSvc.GetByID(*AuditFindingActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := afaps.actionPlanStatusSvc.GetByID(AuditFindingActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (afaps *AuditFindingActionPlanService) parseAuditFindingActionPlanDates(dto *auditFindingDTOs.AuditFindingActionPlanFormDTO) (*auditFindingDTOs.AuditFindingActionPlanParsedDates, error) {
	parsed := &auditFindingDTOs.AuditFindingActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
