package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestRetestServiceInterface interface {
	GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error)
	GetByID(id uint) (*models.PenetrationTestRetest, error)
	GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestRetestDetailDTO, error)
	Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestRetestFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFormDTO) error
	VerifyFinding(id uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestRetestFindingFormDTO) error
	Complete(id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type PenetrationTestRetestService struct {
	repo           repositories.PenetrationTestRetestRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestRetestService(repo repositories.PenetrationTestRetestRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestRetestServiceInterface {
	return &PenetrationTestRetestService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (ptrs *PenetrationTestRetestService) GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error) {
	return ptrs.repo.GetAll(params)
}

func (ptrs *PenetrationTestRetestService) GetByID(id uint) (*models.PenetrationTestRetest, error) {
	value, err := ptrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (ptrs *PenetrationTestRetestService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestRetestDetailDTO, error) {
	value, err := ptrs.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &penetrationTestDTOs.PenetrationTestRetestDetailDTO{}
	if err := mapstructure.Decode(value, detail); err != nil {
		return nil, err
	}
	detail.CreatedAt = value.CreatedAt
	detail.UpdatedAt = value.UpdatedAt
	detail.StartDate = value.StartDate
	detail.EndDate = value.EndDate
	detail.Findings = make([]*penetrationTestDTOs.PenetrationTestRetestFindingDTO, 0, len(value.RetestFindings))
	for _, retestFinding := range value.RetestFindings {
		finding := retestFinding.PenetrationTestFinding
		if finding == nil {
			continue
		}
		detail.Findings = append(detail.Findings, &penetrationTestDTOs.PenetrationTestRetestFindingDTO{
			ID:                       retestFinding.ID,
			PenetrationTestFindingID: retestFinding.PenetrationTestFindingID,
			FindingReferenceNumber:   finding.ReferenceNumber,
			FindingTitle:             finding.Title,
			IsVerified:               retestFinding.IsVerified,
			VerifiedAt:               retestFinding.VerifiedAt,
		})
		if !retestFinding.IsVerified && finding.SeverityLevel != nil {
			switch finding.SeverityLevel.Type {
			case "critical":
				detail.RemainingCritical++
			case "high":
				detail.RemainingHigh++
			case "medium":
				detail.RemainingMedium++
			case "low":
				detail.RemainingLow++
			}
			detail.RemainingTotal++
		}
	}
	return detail, nil
}

func (ptrs *PenetrationTestRetestService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestRetestFormDTO) (uint, error) {
	seen := make(map[uint]bool)
	for _, findingID := range formDTO.FindingIDs {
		if findingID == 0 || seen[findingID] {
			return 0, utility.NewValidationError(fmt.Errorf("finding IDs must be unique and valid"), "invalid", "finding_ids")
		}
		seen[findingID] = true
	}
	startDate, err := utility.ParseMandatoryDate(formDTO.StartDate)
	if err != nil {
		return 0, utility.NewValidationError(err, "invalid", "start_date")
	}
	var endDate *time.Time
	if formDTO.EndDate != nil && *formDTO.EndDate != "" {
		endDate, err = utility.ParseDate(*formDTO.EndDate)
		if err != nil {
			return 0, utility.NewValidationError(err, "invalid", "end_date")
		}
		if endDate.Before(startDate) {
			return 0, utility.NewValidationError(fmt.Errorf("end date cannot be before start date"), "invalid", "end_date")
		}
	}
	if formDTO.Status == "completed" && endDate == nil {
		return 0, utility.NewValidationError(fmt.Errorf("end date is required when completing a re-test"), "invalid", "end_date")
	}
	value := &models.PenetrationTestRetest{
		PenetrationTestID:         formDTO.PenetrationTestID,
		StartDate:                 startDate,
		EndDate:                   endDate,
		PenetrationTestExecutorID: formDTO.PenetrationTestExecutorID,
		Status:                    formDTO.Status,
		Notes:                     formDTO.Notes,
	}
	active, err := ptrs.repo.GetActiveByPenetrationTestID(value.PenetrationTestID)
	if err != nil {
		return 0, err
	}
	if active != nil {
		return 0, utility.NewValidationError(fmt.Errorf("an active re-test already exists"), "invalid", "penetration_test_id")
	}
	if err := ptrs.repo.ValidateReferences(value, formDTO.FindingIDs); err != nil {
		return 0, err
	}

	tx := ptrs.repo.BeginTransaction()
	value.ReferenceNumber, err = ptrs.repo.GenerateReferenceNumber(value.PenetrationTestID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.CreatedBYID = currentUserID
	value.CreatedAt = time.Now()
	id, err := ptrs.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, findingID := range formDTO.FindingIDs {
		_, err := ptrs.repo.CreateRetestFinding(&models.PenetrationTestRetestFinding{
			PenetrationTestRetestID:  value.ID,
			PenetrationTestFindingID: findingID,
			CreatedBYID:              currentUserID,
			CreatedAt:                time.Now(),
		}, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "create penetration test re-test", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return id, nil
}

func (ptrs *PenetrationTestRetestService) Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFormDTO) error {
	seen := make(map[uint]bool)
	for _, findingID := range form.FindingIDs {
		if findingID == 0 || seen[findingID] {
			return utility.NewValidationError(fmt.Errorf("finding IDs must be unique and valid"), "invalid", "finding_ids")
		}
		seen[findingID] = true
	}
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	startDate, err := utility.ParseMandatoryDate(form.StartDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(err, "invalid", "start_date")
	}
	var endDate *time.Time
	if form.EndDate != nil && *form.EndDate != "" {
		endDate, err = utility.ParseDate(*form.EndDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid", "end_date")
		}
		if endDate.Before(startDate) {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("end date cannot be before start date"), "invalid", "end_date")
		}
	}
	if form.Status == "completed" && endDate == nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("end date is required when completing a re-test"), "invalid", "end_date")
	}
	updated := &models.PenetrationTestRetest{
		PenetrationTestID:         form.PenetrationTestID,
		StartDate:                 startDate,
		EndDate:                   endDate,
		PenetrationTestExecutorID: form.PenetrationTestExecutorID,
		Status:                    form.Status,
		Notes:                     form.Notes,
	}
	if existing.Status == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be updated"), "invalid", "status")
	}
	if err := ptrs.repo.ValidateReferences(updated, form.FindingIDs); err != nil {
		tx.Rollback()
		return err
	}
	updated.ID = existing.ID
	updated.ReferenceNumber = existing.ReferenceNumber
	updated.CreatedBYID = existing.CreatedBYID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedBYID = &currentUserID
	now := time.Now()
	updated.UpdatedAt = &now
	if err := ptrs.repo.Update(updated, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := ptrs.repo.DeleteRetestFindings(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	for _, findingID := range form.FindingIDs {
		if _, err := ptrs.repo.CreateRetestFinding(&models.PenetrationTestRetestFinding{
			PenetrationTestRetestID:  id,
			PenetrationTestFindingID: findingID,
			CreatedBYID:              currentUserID,
			CreatedAt:                time.Now(),
		}, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "update penetration test re-test", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) VerifyFinding(id uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestRetestFindingFormDTO) error {
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if existing.Status == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be changed"), "invalid", "status")
	}
	for _, retestFinding := range existing.RetestFindings {
		if retestFinding.PenetrationTestFindingID == formDTO.FindingID {
			retestFinding.IsVerified = formDTO.IsVerified
			if formDTO.IsVerified {
				now := time.Now()
				retestFinding.VerifiedAt = &now
			} else {
				retestFinding.VerifiedAt = nil
			}
			if err := ptrs.repo.UpdateRetestFinding(retestFinding, tx); err != nil {
				tx.Rollback()
				return err
			}
			if err := tx.Commit(); err != nil {
				return utility.NewInternalServerError(err)
			}
			if ptrs.activityLogSvc != nil {
				go func() {
					existingEntry, err := ptrs.GetFullDetailsByID(id)
					if err != nil {
						utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
						return
					}
					if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "verify penetration test re-test finding", nil); err != nil {
						utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
					}
				}()
			}
			return nil
		}
	}
	tx.Rollback()
	return utility.NewValidationError(fmt.Errorf("finding is not linked to this re-test"), "invalid", "finding_id")
}

func (ptrs *PenetrationTestRetestService) Complete(id uint, currentUserID uint) error {
	tx := ptrs.repo.BeginTransaction()
	value, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if value.Status == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("re-test is already completed"), "invalid", "status")
	}
	now := time.Now()
	value.Status = "completed"
	value.EndDate = &now
	value.UpdatedBYID = &currentUserID
	value.UpdatedAt = &now
	if err := ptrs.repo.Update(value, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "complete penetration test re-test", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deleted []map[string]interface{}
	for _, id := range ids {
		value, err := ptrs.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		tx := ptrs.repo.BeginTransaction()
		if err := ptrs.repo.DeleteRetestFindings(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := ptrs.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		successfulIDs = append(successfulIDs, id)
		deleted = append(deleted, map[string]interface{}{"id": id, "reference_number": value.ReferenceNumber})
	}
	if len(successfulIDs) > 0 && ptrs.activityLogSvc != nil {
		go ptrs.activityLogSvc.LogBulkActivities(models.DeleteAction, models.PenetrationTestRetest{}.TableName(), successfulIDs, deleted, currentUserID, nil)
	}
	return failedIDs, nil
}
