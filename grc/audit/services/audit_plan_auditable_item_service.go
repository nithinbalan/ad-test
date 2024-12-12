package services

import (
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanAuditableItemServiceInterface interface {
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanAuditableItemListQueryParamsDTO) (*auditDTOs.AuditPlanAuditableItemListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanAuditableItem, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanAuditableItemDetailDTO, error)
	Create(auditPlanID uint, currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanAuditableItemFormDTO) (uint, error)
	Delete(auditPlanID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetByAuditableItemID(auditableItemID uint) (*models.AuditPlanAuditableItem, error)
}

type AuditPlanAuditableItemService struct {
	repo             repositories.AuditPlanAuditableItemRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditPlanSvc     AuditPlanServiceInterface
	auditableItemSvc AuditableItemServiceInterface
}

func NewAuditPlanAuditableItemService(
	repo repositories.AuditPlanAuditableItemRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditPlanAuditableItemServiceInterface {
	return &AuditPlanAuditableItemService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditPlanSvc:     auditPlanSvc,
		auditableItemSvc: auditableItemSvc,
	}
}

func (apais *AuditPlanAuditableItemService) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanAuditableItemListQueryParamsDTO) (*auditDTOs.AuditPlanAuditableItemListResponseDTO, error) {
	if _, err := apais.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}
	return apais.repo.GetAll(auditPlanID, params)
}

func (apais *AuditPlanAuditableItemService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanAuditableItem, error) {
	auditPlanAuditableItem, err := apais.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlanAuditableItem, nil
}

func (apais *AuditPlanAuditableItemService) GetByAuditableItemID(auditableItemID uint) (*models.AuditPlanAuditableItem, error) {
	auditPlanAuditableItem, err := apais.repo.GetByAuditableItemID(auditableItemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditableItemID)
		}
		return nil, err
	}
	return auditPlanAuditableItem, nil
}

func (apais AuditPlanAuditableItemService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanAuditableItemDetailDTO, error) {
	var auditableItemDTO auditDTOs.AuditPlanAuditableItemDetailDTO

	// Get full details of the audit action plan
	auditableItem, err := apais.repo.GetFullDetailsByID(id, tx...)
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

	// TODO: Add any additional fields required for the details screen

	return &auditableItemDTO, nil
}

func (apais *AuditPlanAuditableItemService) Create(auditPlanID uint, currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanAuditableItemFormDTO) (uint, error) {
	tx := apais.repo.BeginTransaction()

	var createdIds []uint
	var addedAuditPlanAuditableItem []map[string]interface{}

	if _, err := apais.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}
	if err := apais.validateAuditPlanAuditableItemFormDTO(auditPlanFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Check if any of the AuditableItemIDs are already added to the table with this AuditPlanID
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		var existingItem models.AuditPlanAuditableItem
		err := apais.repo.GetByAuditPlanAndAuditableItemID(auditPlanID, *itemID, &existingItem, tx)
		if err == nil {
			tx.Rollback()
			return 0, utility.NewForbiddenError(fmt.Sprintf("The auditable item with ID %d is already added to this plan", *itemID))
		}
	}

	//var auditPlanAuditableItems []*models.AuditPlanAuditableItem
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		auditableItem, err := apais.auditableItemSvc.GetByID(*itemID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		auditPlanAuditableItem := &models.AuditPlanAuditableItem{
			AuditPlanID:     auditPlanID,
			AuditableItemID: *itemID,
			CreatedBYID:     currentUserID,
		}

		// Insert the AuditPlanAuditableItem and get the generated ID
		auditPlanAuditableItemID, err := apais.repo.Create(auditPlanAuditableItem, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if auditableItem.Type == "process" {
			riskIDs, err := apais.GetProcessAssociatedRisks(*auditableItem.ProcessID, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			// Insert associated risks for the newly inserted AuditPlanAuditableItem

			for _, riskID := range riskIDs {
				risk := &models.AuditPlanAuditableItemRisk{
					AuditPlanAuditableItemID: auditPlanAuditableItemID,
					RiskRegisterID:           riskID,
					CreatedByID:              currentUserID,
				}
				auditPlanAuditableItemRiskID, err := apais.repo.CreateAuditableItemRisk(risk, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}

				controlIDs, err := apais.GetRiskControls(risk.RiskRegisterID, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}

				// Prepare risk control entries
				var riskControlsToInsert []*models.AuditPlanAuditableItemRiskControl
				for _, controlID := range controlIDs {
					riskControl := &models.AuditPlanAuditableItemRiskControl{
						AuditPlanAuditableItemRiskID: auditPlanAuditableItemRiskID,
						ControlID:                    controlID,
						CreatedByID:                  currentUserID,
					}
					riskControlsToInsert = append(riskControlsToInsert, riskControl)
				}

				if len(riskControlsToInsert) > 0 {
					if err := apais.repo.CreateAuditableItemRiskControls(riskControlsToInsert, tx); err != nil {
						tx.Rollback()
						return 0, err
					}
				}

			}

		}

		existingAuditableItem, err := apais.GetFullDetailsByID(auditPlanAuditableItemID)
		if err != nil {
			utility.LogError("Failed to fetch audit", map[string]interface{}{
				"error": err.Error(),
			})
			return 0, err
		}

		addedAuditPlanAuditableItem = append(addedAuditPlanAuditableItem, map[string]interface{}{
			"id":                        auditPlanAuditableItemID,
			"audit_plan_auditable_item": existingAuditableItem,
		})

		createdIds = append(createdIds, auditPlanAuditableItemID)
	}

	tx.Commit()

	go func() {
		err := apais.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.AuditPlanAuditableItem{}.TableName(),
			createdIds, //todo
			addedAuditPlanAuditableItem,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanAuditableItem{}.TableName(),
				"recordID":  auditPlanID, //todo
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return auditPlanID, nil
}

func (apais *AuditPlanAuditableItemService) Delete(auditPlanID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := apais.repo.BeginTransaction()
	if _, err := apais.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}
	var failedIDs []uint
	var successfulIDs []uint
	var deletedAuditPlanAuditableItem []map[string]interface{}
	for _, id := range ids {
		// Get existing audit plan details
		auditableItem, err := apais.GetFullDetailsByID(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		// Delete the audit plan itself
		if err := apais.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedAuditPlanAuditableItem = append(deletedAuditPlanAuditableItem, map[string]interface{}{
			"id":         id,
			"audit_plan": auditableItem.AuditPlan.Title,
		})
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := apais.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlanAuditableItem{}.TableName(),
					id,
					deletedAuditPlanAuditableItem,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlanAuditableItem{}.TableName(),
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

func (apais *AuditPlanAuditableItemService) validateAuditPlanAuditableItemFormDTO(dto *auditDTOs.AuditPlanAuditableItemFormDTO) error {
	// Check if AuditableItemIds is nil
	if dto.AuditableItemIds == nil {
		return utility.NewValidationError(fmt.Errorf("auditable item IDs are required"), "required", "auditable item IDs")
	}

	// Check if the list is empty
	if len(dto.AuditableItemIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("at least one auditable item ID is required"), "required", "auditable item ID")
	}

	// Validate each AuditableItemId
	for _, itemID := range dto.AuditableItemIds {
		if itemID != nil {
			if _, err := apais.auditableItemSvc.GetByID(*itemID); err != nil {
				return utility.NewValidationError(fmt.Errorf("auditable item with ID %d does not exist", *itemID), "not_exist", "auditable item")
			}
		}
	}

	return nil
}

// New function to get risks by process ID
func (apais *AuditPlanAuditableItemService) GetProcessAssociatedRisks(processID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	return apais.repo.GetProcessAssociatedRisks(processID, tx...)
}

func (apais *AuditPlanAuditableItemService) GetRiskControls(riskID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	return apais.repo.GetRiskControls(riskID, tx...)
}
