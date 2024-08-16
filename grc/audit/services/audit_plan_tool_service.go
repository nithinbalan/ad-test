package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanToolServiceInterface interface {
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanToolListQueryParamsDTO) (*auditDTOs.AuditPlanToolListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanToolDetailDTO, error)
	GetByID(id uint) (*models.AuditPlanTool, error)
	Create(auditPlanID uint, currentUserID uint, auditPlanToolFormDTO *auditDTOs.AuditPlanToolFormDTO) (uint, error)
	Update(id uint, auditPlanID uint, currentUserID uint, auditPlanToolFormDTO *auditDTOs.AuditPlanToolFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditPlanToolService struct {
	repo           repositories.AuditPlanToolRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	auditPlanSvc   AuditPlanServiceInterface
	toolTypeSvc    masterServices.ToolTypeServiceInterface
}

func NewAuditPlanToolService(
	repo repositories.AuditPlanToolRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	toolTypeSvc masterServices.ToolTypeServiceInterface,
) AuditPlanToolServiceInterface {
	return &AuditPlanToolService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		auditPlanSvc:   auditPlanSvc,
		toolTypeSvc:    toolTypeSvc,
	}
}

func (apts *AuditPlanToolService) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanToolListQueryParamsDTO) (*auditDTOs.AuditPlanToolListResponseDTO, error) {
	if _, err := apts.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}
	return apts.repo.GetAll(auditPlanID, params)
}

func (apts *AuditPlanToolService) GetByID(id uint) (*models.AuditPlanTool, error) {
	auditPlanTool, err := apts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlanTool, nil
}

func (apts *AuditPlanToolService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanToolDetailDTO, error) {
	var auditPlanToolDTO auditDTOs.AuditPlanToolDetailDTO

	auditPlanTool, err := apts.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlanTool, &auditPlanToolDTO)
	if err != nil {
		return nil, err
	}

	auditPlanToolDTO.CreatedAt = auditPlanTool.CreatedAt
	auditPlanToolDTO.UpdatedAt = auditPlanTool.UpdatedAt

	return &auditPlanToolDTO, nil
}

func (apts *AuditPlanToolService) Create(auditPlanID uint, currentUserID uint, auditPlanToolFormDTO *auditDTOs.AuditPlanToolFormDTO) (uint, error) {
	tx := apts.repo.BeginTransaction()
	if _, err := apts.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}
	if err := apts.validateAuditPlanToolFormDTO(auditPlanToolFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditPlanTool := &models.AuditPlanTool{
		AuditPlanID: auditPlanID,
		ToolTypeID:  auditPlanToolFormDTO.ToolTypeID,
		Title:       auditPlanToolFormDTO.Title,
		Purpose:     auditPlanToolFormDTO.Purpose,
		CreatedBYID: currentUserID,
	}

	id, err := apts.repo.Create(auditPlanTool, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := apts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlanTool{}.TableName(),
			id,
			auditPlanToolFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanTool{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (apts *AuditPlanToolService) Update(id uint, auditPlanID uint, currentUserID uint, auditPlanToolFormDTO *auditDTOs.AuditPlanToolFormDTO) error {
	tx := apts.repo.BeginTransaction()
	if _, err := apts.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return err
	}
	if err := apts.validateAuditPlanToolFormDTO(auditPlanToolFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanTool, err := apts.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditPlanTool.ToolTypeID = auditPlanToolFormDTO.ToolTypeID
	existingAuditPlanTool.Title = auditPlanToolFormDTO.Title
	existingAuditPlanTool.Purpose = auditPlanToolFormDTO.Purpose
	existingAuditPlanTool.UpdatedBYID = &currentUserID

	err = apts.repo.Update(existingAuditPlanTool, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := apts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlanTool{}.TableName(),
			id,
			auditPlanToolFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlanTool{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (apts *AuditPlanToolService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := apts.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		_, err := apts.repo.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := apts.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := apts.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanTool{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanTool{}.TableName(),
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

func (apts *AuditPlanToolService) validateAuditPlanToolFormDTO(dto *auditDTOs.AuditPlanToolFormDTO) error {
	if _, err := apts.toolTypeSvc.GetByID(dto.ToolTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("tool type does not exist"), "not_exist", "tool type")
	}

	return nil
}
