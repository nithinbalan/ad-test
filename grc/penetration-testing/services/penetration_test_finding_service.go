package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestFindingServiceInterface interface {
	GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error)
	GetByID(id uint) (*models.PenetrationTestFinding, error)
	GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestFindingDetailDTO, error)
	Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, statusID uint, percentage *uint, tx utility.DatabaseTransaction) error
}

type PenetrationTestFindingService struct {
	repo                repositories.PenetrationTestFindingRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	penetrationTestSvc  PenetrationTestServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
}

func NewPenetrationTestFindingService(
	repo repositories.PenetrationTestFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) PenetrationTestFindingServiceInterface {
	return &PenetrationTestFindingService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		penetrationTestSvc:  penetrationTestSvc,
		severityLevelSvc:    severityLevelSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
	}
}

func (ptfs *PenetrationTestFindingService) GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error) {
	return ptfs.repo.GetAll(params)
}

func (ptfs *PenetrationTestFindingService) GetByID(id uint) (*models.PenetrationTestFinding, error) {
	value, err := ptfs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (ptfs *PenetrationTestFindingService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestFindingDetailDTO, error) {
	value, err := ptfs.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &penetrationTestDTOs.PenetrationTestFindingDetailDTO{}
	if err := mapstructure.Decode(value, detail); err != nil {
		return nil, err
	}
	detail.CreatedAt = value.CreatedAt
	detail.UpdatedAt = value.UpdatedAt
	return detail, nil
}

func (ptfs *PenetrationTestFindingService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) (uint, error) {
	tx := ptfs.repo.BeginTransaction()
	if err := ptfs.validatePenetrationTestFindingFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	var targetRemediationDate *time.Time
	var err error
	if formDTO.TargetRemediationDate != nil && *formDTO.TargetRemediationDate != "" {
		targetRemediationDate, err = utility.ParseDate(*formDTO.TargetRemediationDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(err, "invalid_date", "target remediation date")
		}
	}
	value := &models.PenetrationTestFinding{
		PenetrationTestID:     formDTO.PenetrationTestID,
		Title:                 formDTO.Title,
		ActionPlanStatusID:    formDTO.ActionPlanStatusID,
		Percentage:            formDTO.Percentage,
		SeverityLevelID:       formDTO.SeverityLevelID,
		Description:           formDTO.Description,
		Recommendation:        formDTO.Recommendation,
		OwnerID:               formDTO.OwnerID,
		TargetRemediationDate: targetRemediationDate,
	}
	referenceNumber, err := ptfs.repo.GenerateReferenceNumber(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.ReferenceNumber = referenceNumber
	value.CreatedBYID = currentUserID
	value.CreatedAt = time.Now()

	id, err := ptfs.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if _, err := ptfs.repo.CreateInitialPenetrationTestFindingStatusUpdate(&models.PenetrationTestFindingStatusUpdate{
		PenetrationTestFindingID: id,
		ActionPlanStatusID:       formDTO.ActionPlanStatusID,
		Percentage:               formDTO.Percentage,
		Comment:                  "Initial action plan status",
		CreatedBYID:              currentUserID,
		CreatedAt:                time.Now(),
	}, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingEntry, err := ptfs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch penetration test finding", map[string]interface{}{"error": err.Error()})
			return
		}
		if err := ptfs.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTestFinding{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "create penetration test finding", nil); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return id, nil
}

func (ptfs *PenetrationTestFindingService) Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingFormDTO) error {
	tx := ptfs.repo.BeginTransaction()
	existing, err := ptfs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if form.ActionPlanStatusID == 0 {
		form.ActionPlanStatusID = existing.ActionPlanStatusID
	}
	if err := ptfs.validatePenetrationTestFindingFormDTO(form); err != nil {
		tx.Rollback()
		return err
	}

	var targetRemediationDate *time.Time
	if form.TargetRemediationDate != nil && *form.TargetRemediationDate != "" {
		targetRemediationDate, err = utility.ParseDate(*form.TargetRemediationDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid_date", "target remediation date")
		}
	}
	updated := &models.PenetrationTestFinding{
		PenetrationTestID:     form.PenetrationTestID,
		Title:                 form.Title,
		ActionPlanStatusID:    form.ActionPlanStatusID,
		Percentage:            form.Percentage,
		SeverityLevelID:       form.SeverityLevelID,
		Description:           form.Description,
		Recommendation:        form.Recommendation,
		OwnerID:               form.OwnerID,
		TargetRemediationDate: targetRemediationDate,
	}
	updated.ID = existing.ID
	updated.ReferenceNumber = existing.ReferenceNumber
	updated.CreatedBYID = existing.CreatedBYID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedBYID = &currentUserID
	now := time.Now()
	updated.UpdatedAt = &now

	if err := ptfs.repo.Update(updated, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		existingEntry, err := ptfs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch penetration test finding", map[string]interface{}{"error": err.Error()})
			return
		}
		if err := ptfs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestFinding{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "update penetration test finding", nil); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return nil
}

func (ptfs *PenetrationTestFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deleted []map[string]interface{}

	for _, id := range ids {
		value, err := ptfs.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		tx := ptfs.repo.BeginTransaction()
		if err := ptfs.repo.DeletePenetrationTestFindingStatusUpdates(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := ptfs.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deleted = append(deleted, map[string]interface{}{
			"id":    id,
			"title": value.Title,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := ptfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestFinding{}.TableName(),
				successfulIDs,
				deleted,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return failedIDs, nil
}

func (ptfs *PenetrationTestFindingService) StatusUpdate(id uint, currentUserID uint, statusID uint, percentage *uint, tx utility.DatabaseTransaction) error {
	finding, err := ptfs.repo.GetByID(id, tx)
	if err != nil {
		return err
	}
	finding.ActionPlanStatusID = statusID
	finding.Percentage = percentage
	finding.UpdatedBYID = &currentUserID
	now := time.Now()
	finding.UpdatedAt = &now
	return ptfs.repo.Update(finding, tx)
}

func (ptfs *PenetrationTestFindingService) validatePenetrationTestFindingFormDTO(formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) error {
	if _, err := ptfs.penetrationTestSvc.GetByID(formDTO.PenetrationTestID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test does not exist"), "not_exist", "penetration test")
	}
	if _, err := ptfs.severityLevelSvc.GetByID(formDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}
	if formDTO.ActionPlanStatusID == 0 {
		return utility.NewValidationError(fmt.Errorf("action plan status is required"), "required", "action plan status")
	}
	if _, err := ptfs.actionPlanStatusSvc.GetByID(formDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}
	if _, err := ptfs.userSvc.GetByID(formDTO.OwnerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("owner does not exist"), "not_exist", "owner")
	}
	return nil
}
