package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestFindingServiceInterface interface {
	GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error)
	GetByID(id uint) (*models.PenetrationTestFinding, error)
	GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestFindingDetailDTO, error)
	Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type PenetrationTestFindingService struct {
	repo           repositories.PenetrationTestFindingRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestFindingService(repo repositories.PenetrationTestFindingRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestFindingServiceInterface {
	return &PenetrationTestFindingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (service *PenetrationTestFindingService) GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error) {
	return service.repo.GetAll(params)
}

func (service *PenetrationTestFindingService) GetByID(id uint) (*models.PenetrationTestFinding, error) {
	value, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (service *PenetrationTestFindingService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestFindingDetailDTO, error) {
	value, err := service.GetByID(id)
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

func (service *PenetrationTestFindingService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) (uint, error) {
	tx := service.repo.BeginTransaction()
	var targetRemediationDate *time.Time
	var err error
	if formDTO.TargetRemediationDate != nil && *formDTO.TargetRemediationDate != "" {
		targetRemediationDate, err = utility.ParseDate(*formDTO.TargetRemediationDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(err, "invalid", "target_remediation_date")
		}
	}
	value := &models.PenetrationTestFinding{
		PenetrationTestID:     formDTO.PenetrationTestID,
		Title:                 formDTO.Title,
		SeverityLevelID:       formDTO.SeverityLevelID,
		Description:           formDTO.Description,
		Recommendation:        formDTO.Recommendation,
		OwnerID:               formDTO.OwnerID,
		TargetRemediationDate: targetRemediationDate,
	}

	if err := service.repo.ValidateReferences(value); err != nil {
		tx.Rollback()
		return 0, err
	}

	referenceNumber, err := service.repo.GenerateReferenceNumber(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.ReferenceNumber = referenceNumber
	value.CreatedBYID = currentUserID
	value.CreatedAt = time.Now()

	id, err := service.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	if service.activityLogSvc != nil {
		go func() {
			existingEntry, err := service.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test finding", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := service.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTestFinding{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "create penetration test finding", nil); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return id, nil
}

func (service *PenetrationTestFindingService) Update(id uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFindingFormDTO) error {
	tx := service.repo.BeginTransaction()
	existing, err := service.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var targetRemediationDate *time.Time
	if formDTO.TargetRemediationDate != nil && *formDTO.TargetRemediationDate != "" {
		targetRemediationDate, err = utility.ParseDate(*formDTO.TargetRemediationDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid", "target_remediation_date")
		}
	}
	updated := &models.PenetrationTestFinding{
		PenetrationTestID:     formDTO.PenetrationTestID,
		Title:                 formDTO.Title,
		SeverityLevelID:       formDTO.SeverityLevelID,
		Description:           formDTO.Description,
		Recommendation:        formDTO.Recommendation,
		OwnerID:               formDTO.OwnerID,
		TargetRemediationDate: targetRemediationDate,
	}
	if err := service.repo.ValidateReferences(updated); err != nil {
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

	if err := service.repo.Update(updated, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	if service.activityLogSvc != nil {
		go func() {
			existingEntry, err := service.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test finding", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := service.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestFinding{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "update penetration test finding", nil); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (service *PenetrationTestFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deleted []map[string]interface{}

	for _, id := range ids {
		value, err := service.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		tx := service.repo.BeginTransaction()
		if err := service.repo.Delete(id, tx); err != nil {
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
			"id":               id,
			"reference_number": value.ReferenceNumber,
			"title":            value.Title,
		})
	}

	if len(successfulIDs) > 0 && service.activityLogSvc != nil {
		go service.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.PenetrationTestFinding{}.TableName(),
			successfulIDs,
			deleted,
			currentUserID,
			nil,
		)
	}
	return failedIDs, nil
}
