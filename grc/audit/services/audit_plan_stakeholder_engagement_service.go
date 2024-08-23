package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanStakeholderEngagementServiceInterface interface {
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanStakeholderEngagementListQueryParamsDTO) (*auditDTOs.AuditPlanStakeholderEngagementListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlanStakeholderEngagement, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanStakeholderEngagementDetailDTO, error)
	Create(auditPlanID uint, currentUserID uint, formDTO *auditDTOs.AuditPlanStakeholderEngagementFormDTO) (uint, error)
	Update(id uint, auditPlanID uint, currentUserID uint, formDTO *auditDTOs.AuditPlanStakeholderEngagementFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditPlanStakeholderEngagementService struct {
	repo                   repositories.AuditPlanStakeholderEngagementRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	auditPlanSvc           AuditPlanServiceInterface
	auditEngagementTypeSvc AuditEngagementTypeServiceInterface
	userSvc                organizationServices.UserServiceInterface
	auditScheduleStatusSvc AuditScheduleStatusServiceInterface
}

func NewAuditPlanStakeholderEngagementService(
	repo repositories.AuditPlanStakeholderEngagementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	auditEngagementTypeSvc AuditEngagementTypeServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditScheduleStatusSvc AuditScheduleStatusServiceInterface,
) AuditPlanStakeholderEngagementServiceInterface {
	return &AuditPlanStakeholderEngagementService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		auditPlanSvc:           auditPlanSvc,
		auditEngagementTypeSvc: auditEngagementTypeSvc,
		userSvc:                userSvc,
		auditScheduleStatusSvc: auditScheduleStatusSvc,
	}
}

func (apses *AuditPlanStakeholderEngagementService) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanStakeholderEngagementListQueryParamsDTO) (*auditDTOs.AuditPlanStakeholderEngagementListResponseDTO, error) {
	if _, err := apses.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}
	return apses.repo.GetAll(auditPlanID, params)
}

func (apses *AuditPlanStakeholderEngagementService) GetByID(id uint) (*models.AuditPlanStakeholderEngagement, error) {
	engagement, err := apses.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return engagement, nil
}

func (apses *AuditPlanStakeholderEngagementService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanStakeholderEngagementDetailDTO, error) {
	var engagementDTO auditDTOs.AuditPlanStakeholderEngagementDetailDTO

	engagement, err := apses.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(engagement, &engagementDTO)
	if err != nil {
		return nil, err
	}
	engagementDTO.CreatedAt = engagement.CreatedAt
	engagementDTO.UpdatedAt = engagement.UpdatedAt

	return &engagementDTO, nil
}

func (apses *AuditPlanStakeholderEngagementService) Create(auditPlanID uint, currentUserID uint, formDTO *auditDTOs.AuditPlanStakeholderEngagementFormDTO) (uint, error) {
	tx := apses.repo.BeginTransaction()
	if _, err := apses.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}

	if err := apses.validateFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := apses.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}
	auditScheduleStatus, err := apses.auditScheduleStatusSvc.GetByType("scheduled")
	if err != nil {
		return 0, err
	}
	engagement := &models.AuditPlanStakeholderEngagement{
		ReferenceNumber:       referenceNumber,
		AuditPlanID:           auditPlanID,
		AuditEngagementTypeID: formDTO.AuditEngagementTypeID,
		Title:                 formDTO.Title,
		Objectives:            formDTO.Objectives,
		ResponsibleUserID:     formDTO.ResponsibleUserID,
		AuditScheduleStatusID: uint(auditScheduleStatus.ID),
		CreatedByID:           currentUserID,
	}

	id, err := apses.repo.Create(engagement, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := apses.activityLogSvc.LogActivity(
			models.CreateAction,
			engagement.TableName(),
			id,
			formDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": engagement.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (apses *AuditPlanStakeholderEngagementService) Update(id uint, auditPlanID uint, currentUserID uint, formDTO *auditDTOs.AuditPlanStakeholderEngagementFormDTO) error {
	tx := apses.repo.BeginTransaction()
	if _, err := apses.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return err
	}

	if err := apses.validateFormDTO(formDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingEngagement, err := apses.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingEngagement.AuditEngagementTypeID = formDTO.AuditEngagementTypeID
	existingEngagement.Title = formDTO.Title
	existingEngagement.Objectives = formDTO.Objectives
	existingEngagement.ResponsibleUserID = formDTO.ResponsibleUserID
	existingEngagement.UpdatedByID = &currentUserID

	err = apses.repo.Update(existingEngagement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := apses.activityLogSvc.LogActivity(
			models.UpdateAction,
			existingEngagement.TableName(),
			id,
			formDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": existingEngagement.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (apses *AuditPlanStakeholderEngagementService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := apses.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		_, err := apses.repo.GetByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := apses.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := apses.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanStakeholderEngagement{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanStakeholderEngagement{}.TableName(),
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

func (apses *AuditPlanStakeholderEngagementService) validateFormDTO(formDTO *auditDTOs.AuditPlanStakeholderEngagementFormDTO) error {
	if _, err := apses.auditEngagementTypeSvc.GetByID(formDTO.AuditEngagementTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit engagement type does not exist"), "not_exist", "audit engagement type")
	}
	if _, err := apses.userSvc.GetByID(*formDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}
	return nil
}
