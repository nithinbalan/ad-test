package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	auditRepositories "jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditableItemServiceInterface interface {
	GetAll(params *auditDTOs.AuditableItemListQueryParamsDTO) (*auditDTOs.AuditableItemListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditableItem, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditableItemDetailDTO, error)
	Create(currentUserID uint, auditableItemFormDTO *auditDTOs.AuditableItemFormDTO) (uint, error)
	Update(id uint, currentUserID uint, audit *auditDTOs.AuditableItemFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, auditableItemStatusUpdateFormDTO *auditDTOs.AuditableItemStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type AuditableItemService struct {
	repo                 auditRepositories.AuditableItemRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	auditStatusSvc       AuditStatusServiceInterface
	auditorRiskRatingSvc AuditorRiskRatingServiceInterface
	organizationUnitSvc  settingServices.OrganizationUnitServiceInterface
	userSvc              organizationServices.UserServiceInterface
}

func NewAuditableItemService(
	repo auditRepositories.AuditableItemRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditStatusSvc AuditStatusServiceInterface,
	auditorRiskRatingSvc AuditorRiskRatingServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
	userSvc organizationServices.UserServiceInterface) AuditableItemServiceInterface {
	return &AuditableItemService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		auditStatusSvc:       auditStatusSvc,
		auditorRiskRatingSvc: auditorRiskRatingSvc,
		organizationUnitSvc:  organizationUnitSvc,
		userSvc:              userSvc,
	}
}

func (ais *AuditableItemService) GetAll(params *auditDTOs.AuditableItemListQueryParamsDTO) (*auditDTOs.AuditableItemListResponseDTO, error) {
	return ais.repo.GetAll(params)
}

func (ais AuditableItemService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditableItem, error) {
	auditableItem, err := ais.repo.GetByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditableItem, nil
}

func (ais AuditableItemService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditableItemDetailDTO, error) {
	var auditableItemDTO auditDTOs.AuditableItemDetailDTO

	// Get full details of the audit action plan
	auditableItem, err := ais.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the audit action plan into the DTO
	err = mapstructure.Decode(auditableItem, &auditableItemDTO)
	if err != nil {
		return nil, err
	}

	auditableItemDTO.CreatedAt = auditableItem.CreatedAt
	auditableItemDTO.UpdatedAt = auditableItem.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &auditableItemDTO, nil
}

func (ais *AuditableItemService) Create(currentUserID uint, auditableItemFormDTO *auditDTOs.AuditableItemFormDTO) (uint, error) {
	tx := ais.repo.BeginTransaction()
	if err := ais.validateAuditableItemFormDTO(auditableItemFormDTO); err != nil {
		return 0, err
	}
	auditStatus, err := ais.auditStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}
	var responsibleUnitID uint
	unitID, err := ais.userSvc.GetUserPrimaryUnitID(auditableItemFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible unit")
	}
	responsibleUnitID = unitID

	auditableItem := &models.AuditableItem{
		TitleEn:             &auditableItemFormDTO.TitleEn,
		TitleAr:             &auditableItemFormDTO.TitleAr,
		Comment:             auditableItemFormDTO.Comment,
		AuditorRiskRatingID: auditableItemFormDTO.AuditorRiskRatingID,
		ResponsibleUserID:   auditableItemFormDTO.ResponsibleUserID,
		ResponsibleUnitID:   responsibleUnitID,
		AuditStatusID:       uint(auditStatus.ID),
		CreatedBYID:         currentUserID,
	}

	id, err := ais.repo.Create(auditableItem, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := ais.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditableItem{}.TableName(),
			id,
			auditableItemFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditableItem{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ais *AuditableItemService) Update(id uint, currentUserID uint, auditableItemFormDTO *auditDTOs.AuditableItemFormDTO) error {
	tx := ais.repo.BeginTransaction()

	existingAuditableItem, err := ais.GetByID(id)

	if err != nil {
		return err
	}

	if err := ais.validateAuditableItemFormDTO(auditableItemFormDTO); err != nil {
		return err
	}
	if existingAuditableItem.ResponsibleUserID != auditableItemFormDTO.ResponsibleUserID {
		responsibleUnitID, err := ais.userSvc.GetUserPrimaryUnitID(auditableItemFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible unit")
		}
		existingAuditableItem.ResponsibleUnitID = responsibleUnitID
	}
	// Update Compliance ActionPlan
	existingAuditableItem.TitleEn = &auditableItemFormDTO.TitleEn
	existingAuditableItem.TitleAr = &auditableItemFormDTO.TitleAr
	existingAuditableItem.Comment = auditableItemFormDTO.Comment
	existingAuditableItem.AuditorRiskRatingID = auditableItemFormDTO.AuditorRiskRatingID
	existingAuditableItem.ResponsibleUserID = auditableItemFormDTO.ResponsibleUserID
	existingAuditableItem.UpdatedBYID = &currentUserID

	err = ais.repo.Update(existingAuditableItem, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (ais *AuditableItemService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ais.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {
		err := ais.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

	}

	if len(successfulIds) > 0 {
		tx.Commit()
		go func() {
			err := ais.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditableItem{}.TableName(),
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
	} else {
		tx.Rollback()
	}

	return failedIds, nil
}

func (ais *AuditableItemService) StatusUpdate(id uint, currentUserID uint, auditableItemStatusUpdateFormDTO *auditDTOs.AuditableItemStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	existingAuditableItem, err := ais.GetByID(id)

	if err != nil {
		return err
	}

	existingAuditableItem.AuditStatusID = auditableItemStatusUpdateFormDTO.AuditStatusID
	existingAuditableItem.UpdatedBYID = &currentUserID

	err = ais.repo.Update(existingAuditableItem, tx)

	if err != nil {
		return err
	}

	return nil
}

func (ais *AuditableItemService) validateAuditableItemFormDTO(auditableItemFormDTO *auditDTOs.AuditableItemFormDTO) error {
	if _, err := ais.auditorRiskRatingSvc.GetByID(auditableItemFormDTO.AuditorRiskRatingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditor risk rating does not exist"), "not_exist", "auditor risk rating")
	}
	if _, err := ais.userSvc.GetByID(auditableItemFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	return nil
}
