package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	complianceServices "jethur/grc/compliance/services"
	documentServices "jethur/grc/documents/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type AssetPolicyChecklistServiceInterface interface {
	GetByID(id uint) (*models.AssetPolicyChecklist, error)
	Create(formDTO *dtos.AssetPolicyChecklistFormDTO, currentUserID uint, assetPolicyID uint) (uint, error)
	Update(id uint, formDTO *dtos.AssetPolicyChecklistUpdateDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	AnswerAssetPolicyChecklist(id uint, formDTO *dtos.AssetPolicyChecklistAnswerFormDTO, currentUserID uint, assetId uint, assetPolicyId uint) error
	DeleteAnswerAssetPolicyChecklist(ids []uint, currentUserID uint, assetId uint, assetPolicyId uint) ([]uint, error)
	GetAssetPolicyChecklistSummary(assetID uint, policyID uint) (*dtos.AssetPolicyChecklistSummaryDTO, error)
}

type AssetPolicyChecklistService struct {
	repo                        repositories.AssetPolicyChecklistRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	assetRegisterSvc            AssetRegisterServiceInterface
	documentRequirementGroupSvc documentServices.DocumentRequirementGroupServiceInterface
	assetPolicySvc              AssetPolicyServiceInterface
	complianceStatusSvc         complianceServices.ComplianceStatusServiceInterface
}

func NewAssetPolicyChecklistService(
	repo repositories.AssetPolicyChecklistRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	documentRequirementGroupSvc documentServices.DocumentRequirementGroupServiceInterface,
	assetPolicySvc AssetPolicyServiceInterface,
	complianceStatusSvc complianceServices.ComplianceStatusServiceInterface,
) AssetPolicyChecklistServiceInterface {
	return &AssetPolicyChecklistService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		assetRegisterSvc:            assetRegisterSvc,
		documentRequirementGroupSvc: documentRequirementGroupSvc,
		assetPolicySvc:              assetPolicySvc,
		complianceStatusSvc:         complianceStatusSvc,
	}
}

func (apcs *AssetPolicyChecklistService) GetByID(id uint) (*models.AssetPolicyChecklist, error) {
	policy, err := apcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return policy, nil
}

func (apcs *AssetPolicyChecklistService) Create(formDTO *dtos.AssetPolicyChecklistFormDTO, currentUserID uint, assetPolicyID uint) (uint, error) {
	tx := apcs.repo.BeginTransaction()
	existingAssetPolicy, err := apcs.assetPolicySvc.GetByID(assetPolicyID)
	if err != nil {
		return 0, err
	}
	// Fetch QuestionGroupTitle from DocumentRequirementGroup
	questionGroup, err := apcs.documentRequirementGroupSvc.GetByID(existingAssetPolicy.DocumentID, formDTO.DocumentRequirementGroupID)
	if err != nil {
		tx.Rollback()
		return 0, errors.New("document requirement group not found")
	}

	// Create new checklist
	checklist := &models.AssetPolicyChecklist{
		AssetPolicyID:              assetPolicyID,
		DocumentRequirementGroupID: formDTO.DocumentRequirementGroupID,
		QuestionGroupTitle:         questionGroup.TitleEn,
		QuestionTitle:              formDTO.QuestionTitle,
		CreatedBYID:                currentUserID,
	}

	// Save the checklist
	id, err := apcs.repo.Create(checklist, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		err = apcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetPolicyChecklist{}.TableName(),
			id,
			checklist,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetPolicyChecklist{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()
	return id, nil
}
func (apcs *AssetPolicyChecklistService) Update(id uint, formDTO *dtos.AssetPolicyChecklistUpdateDTO, currentUserID uint) error {
	tx := apcs.repo.BeginTransaction()

	// Fetch existing checklist
	existingChecklist, err := apcs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return errors.New("checklist not found")
	}

	// Update fields
	existingChecklist.QuestionTitle = formDTO.QuestionTitle
	now := time.Now()
	existingChecklist.UpdatedAt = &now
	existingChecklist.UpdatedBYID = &currentUserID

	// Save updates
	if err := apcs.repo.Update(existingChecklist, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (apcs *AssetPolicyChecklistService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedChecklists []map[string]interface{}

	if len(ids) == 0 {
		return nil, nil
	}

	// Fetch all existing checklists one by one
	var existingChecklistIDs []uint
	for _, id := range ids {
		checklist, err := apcs.repo.GetByID(id)
		if err != nil {
			fmt.Printf("Failed to fetch checklist with ID %d: %v\n", id, err)
			failedIds = append(failedIds, id)
			continue
		}

		existingChecklistIDs = append(existingChecklistIDs, checklist.ID)
		deletedChecklists = append(deletedChecklists, map[string]interface{}{
			"id": checklist.ID,
		})
	}

	// If no checklists exist for deletion, return early
	if len(existingChecklistIDs) == 0 {
		return failedIds, nil
	}

	tx := apcs.repo.BeginTransaction()

	// Bulk Delete AssetPolicyChecklists
	if err := apcs.repo.DeleteAssetPolicyChecklists(existingChecklistIDs, tx); err != nil {
		tx.Rollback()
		fmt.Printf("Failed to delete checklists: %v\n", err)
		return failedIds, err
	}

	tx.Commit()
	successfulIds = existingChecklistIDs

	// Log Activity in Background
	if len(successfulIds) > 0 {
		go func() {
			err := apcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetPolicyChecklist{}.TableName(),
				successfulIds,
				deletedChecklists,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (apcs *AssetPolicyChecklistService) AnswerAssetPolicyChecklist(id uint, formDTO *dtos.AssetPolicyChecklistAnswerFormDTO, currentUserID uint, assetID uint, assetPolicyID uint) error {
	tx := apcs.repo.BeginTransaction()

	// Fetch existing checklist
	existingChecklist, err := apcs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update checklist fields
	existingChecklist.Justification = formDTO.Justification
	existingChecklist.ComplianceStatusID = &formDTO.ComplianceStatusID
	existingChecklist.AnsweredBYID = &currentUserID
	now := time.Now()
	existingChecklist.AnsweredAt = &now
	existingChecklist.UpdatedAt = &now
	existingChecklist.UpdatedBYID = &currentUserID

	// Save updates
	if err := apcs.repo.Update(existingChecklist, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Fetch all checklists under this AssetPolicyID
	checklists, err := apcs.assetPolicySvc.GetByAssetID(assetPolicyID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Fetch compliance statuses
	compliantStatus, err := apcs.complianceStatusSvc.GetByType("compliant")
	if err != nil {
		tx.Rollback()
		return err
	}
	notApplicableStatus, err := apcs.complianceStatusSvc.GetByType("not_applicable")
	if err != nil {
		tx.Rollback()
		return err
	}
	partiallyCompliantStatus, err := apcs.complianceStatusSvc.GetByType("partially_compliant")
	if err != nil {
		tx.Rollback()
		return err
	}
	notCompliantStatus, err := apcs.complianceStatusSvc.GetByType("not_implemented")
	if err != nil {
		tx.Rollback()
		return err
	}

	// ✅ Simplified logic using status set
	statusSet := map[uint]bool{}

	for _, checklist := range checklists {
		if checklist.ComplianceStatusID != nil {
			statusSet[*checklist.ComplianceStatusID] = true
		}
	}

	// Determine final status
	var finalStatusID *uint
	switch {
	case statusSet[notCompliantStatus.ID]:
		finalStatusID = &notCompliantStatus.ID
	case len(statusSet) == 1 && statusSet[partiallyCompliantStatus.ID]:
		finalStatusID = &partiallyCompliantStatus.ID
	case len(statusSet) == 1 && statusSet[compliantStatus.ID]:
		finalStatusID = &compliantStatus.ID
	case len(statusSet) == 1 && statusSet[notApplicableStatus.ID]:
		finalStatusID = &notApplicableStatus.ID
	}

	// ✅ Update Asset Policy status if applicable
	if finalStatusID != nil {
		updateDTO := &dtos.AssetPolicyStatusUpdateFormDTO{
			ComplianceStatusID: *finalStatusID,
		}

		if err := apcs.assetPolicySvc.StatusUpdate(assetPolicyID, currentUserID, updateDTO, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (apcs *AssetPolicyChecklistService) DeleteAnswerAssetPolicyChecklist(ids []uint, currentUserID uint, assetId uint, assetPolicyId uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var updatedChecklists []*models.AssetPolicyChecklist
	var updatedAssetPolicyChecklists []map[string]interface{}

	if len(ids) == 0 {
		return nil, nil
	}

	// Fetch all existing checklists and prepare updates
	for _, id := range ids {
		checklist, err := apcs.repo.GetByID(id)
		if err != nil {
			fmt.Printf("Checklist not found for ID %d: %v\n", id, err)
			failedIds = append(failedIds, id)
			continue
		}

		// Prepare checklist for update
		now := time.Now()
		checklist.Justification = nil
		checklist.ComplianceStatusID = nil
		checklist.AnsweredBYID = nil
		checklist.AnsweredAt = nil
		checklist.UpdatedAt = &now
		checklist.UpdatedBYID = &currentUserID

		updatedChecklists = append(updatedChecklists, checklist)
		updatedAssetPolicyChecklists = append(updatedAssetPolicyChecklists, map[string]interface{}{
			"id": checklist.ID,
		})
	}

	// If no checklists exist for updating, return early
	if len(updatedChecklists) == 0 {
		return failedIds, nil
	}

	tx := apcs.repo.BeginTransaction()

	// Bulk update checklists using the prepared model objects
	if err := apcs.repo.BulkUpdateAssetPolicyChecklists(updatedChecklists, tx); err != nil {
		tx.Rollback()
		fmt.Printf("Failed to unset asset policy checklists: %v\n", err)
		return failedIds, err
	}
	checklists, err := apcs.assetPolicySvc.GetByAssetID(assetPolicyId, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Fetch compliance statuses
	compliantStatus, err := apcs.complianceStatusSvc.GetByType("compliant")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	notApplicableStatus, err := apcs.complianceStatusSvc.GetByType("not_applicable")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	partiallyCompliantStatus, err := apcs.complianceStatusSvc.GetByType("partially_compliant")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// notCompliantStatus, err := apcs.complianceStatusSvc.GetByType("not_compliant")
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, err
	// }

	// Count occurrences of each status
	var compliantCount, notApplicableCount, partiallyCompliantCount int
	for _, checklist := range checklists {
		if checklist.ComplianceStatusID != nil {
			switch *checklist.ComplianceStatusID {
			case compliantStatus.ID:
				compliantCount++
			case notApplicableStatus.ID:
				notApplicableCount++
			case partiallyCompliantStatus.ID:
				partiallyCompliantCount++
				// case notCompliantStatus.ID:
				// 	notCompliantCount++
			}
		}
	}

	// **Determine Final Policy Status**
	var finalStatusID *uint

	// if notCompliantCount > 0 {
	// 	// If any checklist is **not_compliant**, then the overall policy is **not_compliant**
	// 	finalStatusID = &notCompliantStatus.ID
	// }
	if partiallyCompliantCount > 0 {
		// If all are **partially_compliant** or mixed with **compliant/not_applicable**
		finalStatusID = &partiallyCompliantStatus.ID
	} else if compliantCount > 0 && notApplicableCount == 0 {
		// If **all are compliant** (excluding not_applicable)
		finalStatusID = &compliantStatus.ID
	} else if notApplicableCount > 0 && compliantCount == 0 && partiallyCompliantCount == 0 {
		// If **all are not_applicable**
		finalStatusID = &notApplicableStatus.ID
	} else {
		// No status can be determined (fallback)
		finalStatusID = nil
	}

	// ✅ Update Asset Policy status if applicable
	if finalStatusID != nil {
		updateDTO := &dtos.AssetPolicyStatusUpdateFormDTO{
			ComplianceStatusID: *finalStatusID,
		}

		if err := apcs.assetPolicySvc.StatusUpdate(assetPolicyId, currentUserID, updateDTO, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	for _, checklist := range updatedChecklists {
		successfulIds = append(successfulIds, checklist.ID)
	}

	// Log Activity in Background
	if len(successfulIds) > 0 {
		go func() {
			err := apcs.activityLogSvc.LogBulkActivities(
				models.UpdateAction,
				models.AssetPolicyChecklist{}.TableName(),
				successfulIds,
				updatedAssetPolicyChecklists,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log unset activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (apcs *AssetPolicyChecklistService) GetAssetPolicyChecklistSummary(assetID uint, policyID uint) (*dtos.AssetPolicyChecklistSummaryDTO, error) {
	return apcs.repo.GetAssetPolicyChecklistSummary(policyID)
}
