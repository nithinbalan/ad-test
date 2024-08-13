package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"

	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanChecklistServiceInterface interface {
	GetAll(auditPlanId uint, params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) ([]*auditDTOs.AuditPlanChecklistListDTO, error)
	GetByID(auditPlanId uint, id uint) (*models.AuditPlanChecklist, error)
	Create(auditPlanId uint, currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) (uint, error)
	Update(auditPlanId uint, id uint, currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) error
	Delete(auditPlanId uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditPlanChecklistService struct {
	repo             repositories.AuditPlanChecklistRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditPlanSvc     AuditPlanServiceInterface
	auditableItemSvc AuditableItemServiceInterface
}

func NewAuditPlanChecklistService(
	repo repositories.AuditPlanChecklistRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditPlanChecklistServiceInterface {
	return &AuditPlanChecklistService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditPlanSvc:     auditPlanSvc,
		auditableItemSvc: auditableItemSvc,
	}
}

func (aps *AuditPlanChecklistService) GetAll(auditPlanId uint, params *auditDTOs.AuditPlanChecklistListQueryParamsDTO) ([]*auditDTOs.AuditPlanChecklistListDTO, error) {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanId); err != nil {
		return nil, err
	}
	return aps.repo.GetAll(auditPlanId, params)
}

func (aps *AuditPlanChecklistService) GetByID(auditPlanId uint, id uint) (*models.AuditPlanChecklist, error) {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanId); err != nil {
		return nil, err
	}

	auditPlanChecklist, err := aps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlanChecklist, nil
}

func (aps *AuditPlanChecklistService) Create(auditPlanId uint, currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) (uint, error) {
	tx := aps.repo.BeginTransaction()

	if _, err := aps.auditPlanSvc.GetByID(auditPlanId); err != nil {
		return 0, err
	}

	if err := aps.validateAuditPlanChecklistFormDTO(auditPlanChecklistFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditPlanChecklist := &models.AuditPlanChecklist{
		AuditPlanID:     auditPlanId,
		AuditableItemID: auditPlanChecklistFormDTO.AuditableItemId,
		Question:        auditPlanChecklistFormDTO.Question,
		CreatedBYID:     currentUserID,
	}

	id, err := aps.repo.Create(auditPlanChecklist, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlanChecklist{}.TableName(),
			id,
			auditPlanChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (aps *AuditPlanChecklistService) Update(auditPlanId uint, id uint, currentUserID uint, auditPlanChecklistFormDTO *auditDTOs.AuditPlanChecklistFormDTO) error {
	tx := aps.repo.BeginTransaction()

	if _, err := aps.auditPlanSvc.GetByID(auditPlanId); err != nil {
		return err
	}

	existingAuditPlanChecklist, err := aps.GetByID(auditPlanId, id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := aps.validateAuditPlanChecklistFormDTO(auditPlanChecklistFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanChecklist.AuditableItemID = auditPlanChecklistFormDTO.AuditableItemId
	existingAuditPlanChecklist.Question = auditPlanChecklistFormDTO.Question
	existingAuditPlanChecklist.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlanChecklist, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlanChecklist{}.TableName(),
			id,
			auditPlanChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlanChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aps *AuditPlanChecklistService) Delete(auditPlanId uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := aps.repo.BeginTransaction()

	if _, err := aps.auditPlanSvc.GetByID(auditPlanId); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {

		if _, err := aps.GetByID(auditPlanId, id); err != nil {
			failedIDs = append(failedIDs, id)
		}

		if err := aps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := aps.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanChecklist{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanChecklist{}.TableName(),
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

func (aps *AuditPlanChecklistService) validateAuditPlanChecklistFormDTO(dto *auditDTOs.AuditPlanChecklistFormDTO) error {
	if _, err := aps.auditableItemSvc.GetByID(dto.AuditableItemId); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
	}

	return nil
}
