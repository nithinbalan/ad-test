package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"time"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditChecklistServiceInterface interface {
	GetAll(auditId uint, params *auditDTOs.AuditChecklistListQueryParamsDTO) ([]*auditDTOs.AuditChecklistListDTO, error)
	GetByID(auditId uint, id uint) (*models.AuditChecklist, error)
	GetFullDetailsByID(auditId uint, id uint) (*auditDTOs.AuditChecklistDetailDTO, error)
	Create(auditId uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) (uint, error) 
	Update(auditId uint, id uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) error
	Delete(auditId uint, ids []uint, currentUserID uint) ([]uint, error)
	Answer(auditId uint, id uint, currentUserID uint, auditChecklistAnswerFormDTO *auditDTOs.AuditChecklistAnswerFormDTO) error
	DeleteAnswer(auditId uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AuditChecklistService struct {
	repo             repositories.AuditChecklistRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditSvc         AuditServiceInterface
	auditableItemSvc AuditableItemServiceInterface
	conformanceSvc   masterServices.ConformanceServiceInterface
	statusSvc        masterServices.StatusServiceInterface
}

func NewAuditChecklistService(
	repo repositories.AuditChecklistRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	conformanceSvc masterServices.ConformanceServiceInterface,
	statusSvc masterServices.StatusServiceInterface,
) AuditChecklistServiceInterface {
	return &AuditChecklistService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditSvc:         auditSvc,
		auditableItemSvc: auditableItemSvc,
		conformanceSvc:   conformanceSvc,
		statusSvc:        statusSvc,
	}
}

func (acs *AuditChecklistService) GetAll(auditId uint, params *auditDTOs.AuditChecklistListQueryParamsDTO) ([]*auditDTOs.AuditChecklistListDTO, error) {
	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}
	return acs.repo.GetAll(auditId, params)
}

func (acs *AuditChecklistService) GetByID(auditId uint, id uint) (*models.AuditChecklist, error) {
	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	auditChecklist, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditChecklist, nil
}

func (acs *AuditChecklistService) GetFullDetailsByID(auditId uint, id uint) (*auditDTOs.AuditChecklistDetailDTO, error) {
	var auditChecklitDetailDTO auditDTOs.AuditChecklistDetailDTO

	checklist, err := acs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(checklist, &auditChecklitDetailDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &auditChecklitDetailDTO, nil
}

func (acs *AuditChecklistService) Create(auditId uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) (uint, error) {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return 0, err
	}

	if err := acs.validateAuditChecklistFormDTO(auditChecklistFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	status, err := acs.statusSvc.GetByType("pending")
	if err != nil {
		return 0,err
	}
	auditChecklist := &models.AuditChecklist{
		AuditID:     auditId,
		AuditableItemID: auditChecklistFormDTO.AuditableItemID,
		Question:        auditChecklistFormDTO.Question,
		StatusID:        status.ID,
		CreatedBYID:     currentUserID,
	}

	id, err := acs.repo.Create(auditChecklist, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditChecklist{}.TableName(),
			id,
			auditChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}


func (acs *AuditChecklistService) Update(auditId uint, id uint, currentUserID uint, auditChecklistFormDTO *auditDTOs.AuditChecklistFormDTO) error {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return err
	}

	existingAuditChecklist, err := acs.GetByID(auditId, id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := acs.validateAuditChecklistFormDTO(auditChecklistFormDTO); err != nil {
		tx.Rollback()
		return err
	}
	existingAuditChecklist.AuditableItemID = auditChecklistFormDTO.AuditableItemID
	existingAuditChecklist.Question = auditChecklistFormDTO.Question
	existingAuditChecklist.UpdatedBYID = &currentUserID

	err = acs.repo.Update(existingAuditChecklist, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditChecklist{}.TableName(),
			id,
			auditChecklistFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (acs *AuditChecklistService) Delete(auditId uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {

		if _, err := acs.GetByID(auditId, id); err != nil {
			failedIDs = append(failedIDs, id)
		}

		if err := acs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := acs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditChecklist{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditChecklist{}.TableName(),
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

func (acs *AuditChecklistService) Answer(auditId uint, id uint, currentUserID uint, auditChecklistAnswerFormDTO *auditDTOs.AuditChecklistAnswerFormDTO) error {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return err
	}

	existingAuditChecklist, err := acs.GetByID(auditId, id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := acs.validateAuditChecklistAnswerFormDTO(auditChecklistAnswerFormDTO); err != nil {
		tx.Rollback()
		return err
	}
	now := time.Now()
	status, err := acs.statusSvc.GetByType("completed")
	if err != nil {
		return err
	}
	existingAuditChecklist.ConformanceID = &auditChecklistAnswerFormDTO.ConformanceID
	existingAuditChecklist.Comment = auditChecklistAnswerFormDTO.Comment
	existingAuditChecklist.AnsweredBYID = &currentUserID
	existingAuditChecklist.StatusID = status.ID
	existingAuditChecklist.AnsweredAt = &now

	err = acs.repo.Update(existingAuditChecklist, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := acs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditChecklist{}.TableName(),
			id,
			auditChecklistAnswerFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (acs *AuditChecklistService) DeleteAnswer(auditId uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := acs.repo.BeginTransaction()

	if _, err := acs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	for _, id := range ids {
		existingAuditChecklist, err := acs.GetByID(auditId, id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			return nil, utility.NewResourceNotFoundError(err, id)
		}

		existingAuditChecklist.ConformanceID = nil
		existingAuditChecklist.Comment = nil
		existingAuditChecklist.AnsweredBYID = nil
		existingAuditChecklist.AnsweredAt = nil

		err = acs.repo.Update(existingAuditChecklist, tx)
		if err != nil {
			return nil,err
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := acs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditChecklist{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditChecklist{}.TableName(),
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

func (acs *AuditChecklistService) validateAuditChecklistFormDTO(dto *auditDTOs.AuditChecklistFormDTO) error {
	if _, err := acs.auditableItemSvc.GetByID(*dto.AuditableItemID); err != nil {
		return utility.NewValidationError(fmt.Errorf("Auditable item does not exist"), "not_exist", "auditable item")
	}

	return nil
}


func (acs *AuditChecklistService) validateAuditChecklistAnswerFormDTO(dto *auditDTOs.AuditChecklistAnswerFormDTO) error {
	if _, err := acs.conformanceSvc.GetByID(dto.ConformanceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("Conformance does not exist"), "not_exist", "conformance")
	}

	return nil
}
