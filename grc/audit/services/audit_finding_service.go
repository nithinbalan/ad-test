package services

import (
	"errors"
	"fmt"
	"time"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditFindingServiceInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error)
	GetFullDetailsByID(auditID uint, id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error)
	GetByID(auditID uint, id uint) (*models.AuditFinding, error)
	Create(auditID uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) (uint, error)
	Update(auditID uint, id uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) error
	Delete(auditID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditFindingService struct {
	repo                  repositories.AuditFindingRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	auditSvc              AuditServiceInterface
	auditFindingStatusSvc AuditFindingStatusServiceInterface
	auditTimelineSvc      AuditTimelineServiceInterface
	auditableItemSvc      AuditableItemServiceInterface
	auditFindingCategory  AuditFindingCategoryServiceInterface
	severityLevelSvc      masterServices.SeverityLevelServiceInterface
}

func NewAuditFindingService(
	repo repositories.AuditFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditFindingStatusSvc AuditFindingStatusServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	auditFindingCategory AuditFindingCategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
) AuditFindingServiceInterface {
	return &AuditFindingService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		auditSvc:              auditSvc,
		auditFindingStatusSvc: auditFindingStatusSvc,
		auditTimelineSvc:      auditTimelineSvc,
		auditableItemSvc:      auditableItemSvc,
		auditFindingCategory:  auditFindingCategory,
		severityLevelSvc:      severityLevelSvc,
	}
}

func (afs *AuditFindingService) GetAll(auditID uint, params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error) {
	if _, err := afs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}

	return afs.repo.GetAll(params)
}

func (afs *AuditFindingService) GetByID(auditID uint, id uint) (*models.AuditFinding, error) {
	if _, err := afs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	auditFinding, err := afs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFinding, nil
}

func (afs *AuditFindingService) GetFullDetailsByID(auditID uint, id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error) {
	if _, err := afs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	var auditFindingDTO auditDTOs.AuditFindingDetailDTO

	auditFinding, err := afs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFinding, &auditFindingDTO)
	if err != nil {
		return nil, err
	}

	auditFindingDTO.CreatedAt = auditFinding.CreatedAt
	auditFindingDTO.UpdatedAt = auditFinding.UpdatedAt

	return &auditFindingDTO, nil
}

func (afs *AuditFindingService) Create(auditID uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) (uint, error) {
	if _, err := afs.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}

	tx := afs.repo.BeginTransaction()

	if err := afs.validateAuditFindingFormDTO(auditFindingFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditFindingStatus, err := afs.auditFindingStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	auditFinding := &models.AuditFinding{
		AuditID:                auditID,
		AuditableItemID:        auditFindingFormDTO.AuditableItemID,
		AuditFindingCategoryID: auditFindingFormDTO.AuditFindingCategoryID,
		SeverityLevelID:        auditFindingFormDTO.SeverityLevelID,
		Recommendations:        auditFindingFormDTO.Recommendations,
		AuditFindingStatusID:   auditFindingStatus.ID,
		CreatedAt:              time.Now(),
		CreatedByID:            currentUserID,
	}

	id, err := afs.repo.Create(auditFinding, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := afs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFinding{}.TableName(),
			id,
			auditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afs *AuditFindingService) Update(auditID uint, id uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) error {
	if _, err := afs.auditSvc.GetByID(auditID); err != nil {
		return err
	}
	tx := afs.repo.BeginTransaction()

	existingAuditFinding, err := afs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := afs.validateAuditFindingFormDTO(auditFindingFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditFinding.AuditableItemID = auditFindingFormDTO.AuditableItemID
	existingAuditFinding.AuditFindingCategoryID = auditFindingFormDTO.AuditFindingCategoryID
	existingAuditFinding.SeverityLevelID = auditFindingFormDTO.SeverityLevelID
	existingAuditFinding.Recommendations = auditFindingFormDTO.Recommendations
	existingAuditFinding.UpdatedByID = &currentUserID

	err = afs.repo.Update(existingAuditFinding, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := afs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditFinding{}.TableName(),
			id,
			auditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (afs *AuditFindingService) Delete(auditID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := afs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	tx := afs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		// Get existing audit finding details
		_, err := afs.GetFullDetailsByID(auditID, id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		// Delete the audit finding itself
		if err := afs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := afs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditFinding{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditFinding{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (afs *AuditFindingService) validateAuditFindingFormDTO(dto *auditDTOs.AuditFindingFormDTO) error {

	if _, err := afs.auditableItemSvc.GetByID(dto.AuditableItemID); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
	}

	if _, err := afs.auditFindingCategory.GetByID(dto.AuditFindingCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit finding category does not exist"), "not_exist", "audit finding category")
	}

	if _, err := afs.severityLevelSvc.GetByID(dto.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	return nil
}
