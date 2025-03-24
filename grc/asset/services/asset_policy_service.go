package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	documentServices "jethur/grc/documents/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetPolicyServiceInterface interface {
	GetAll(params *dtos.AssetPolicyListQueryParamsDTO, assetRegisterID uint) (*dtos.AssetPolicyListResponseDTO, error)
	GetByID(id uint) (*models.AssetPolicy, error)
	GetByAssetID(assetId uint, tx ...utility.DatabaseTransaction) ([]*models.AssetPolicy, error)
	GetFullDetailsByID(assetRegisterID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetPolicyDetailDTO, error)
	Update(assetPolicyForm *dtos.AssetPolicyFormDTO, currentUserID uint, assetId uint) ([]uint, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetAssetPolicyComplianceSummary(assetID uint) (*dtos.AssetPolicyByComplianceDTO, error)
	StatusUpdate(id uint, currentUserID uint, assetPolicyStatusUpdateFormDTO *dtos.AssetPolicyStatusUpdateFormDTO, tx utility.DatabaseTransaction) error
}

type AssetPolicyService struct {
	repo                        repositories.AssetPolicyRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	assetRegisterSvc            AssetRegisterServiceInterface
	documentSvc                 documentServices.DocumentServiceInterface
	documentRequirementGroupSvc documentServices.DocumentRequirementGroupServiceInterface
	documentRequirementSvc      documentServices.DocumentRequirementServiceInterface
}

func NewAssetPolicyService(
	repo repositories.AssetPolicyRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
	documentRequirementGroupSvc documentServices.DocumentRequirementGroupServiceInterface,
	documentRequirementSvc documentServices.DocumentRequirementServiceInterface,
) AssetPolicyServiceInterface {
	return &AssetPolicyService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		assetRegisterSvc:            assetRegisterSvc,
		documentSvc:                 documentSvc,
		documentRequirementGroupSvc: documentRequirementGroupSvc,
		documentRequirementSvc:      documentRequirementSvc,
	}
}

func (aps *AssetPolicyService) GetAll(params *dtos.AssetPolicyListQueryParamsDTO, assetRegisterID uint) (*dtos.AssetPolicyListResponseDTO, error) {
	if _, err := aps.assetRegisterSvc.GetByID(assetRegisterID); err != nil {
		return nil, err
	}
	return aps.repo.GetAll(params, assetRegisterID)
}
func (aps *AssetPolicyService) GetByID(id uint) (*models.AssetPolicy, error) {
	policy, err := aps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return policy, nil
}
func (aps *AssetPolicyService) GetByAssetID(assetId uint, tx ...utility.DatabaseTransaction) ([]*models.AssetPolicy, error) {
	policies, err := aps.repo.GetByAssetID(assetId, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, assetId)
		}
		return nil, err
	}
	return policies, nil
}

func (aps *AssetPolicyService) GetFullDetailsByID(assetRegisterID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetPolicyDetailDTO, error) {
	var assetPolicyDetailDTO dtos.AssetPolicyDetailDTO
	if _, err := aps.assetRegisterSvc.GetByID(assetRegisterID); err != nil {
		return nil, err
	}
	assetPolicy, err := aps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(assetPolicy, &assetPolicyDetailDTO)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}
	assetPolicyDetailDTO.CreatedAt = assetPolicy.CreatedAt
	assetPolicyDetailDTO.UpdatedAt = &assetPolicy.UpdatedAt

	return &assetPolicyDetailDTO, nil
}

func (aps *AssetPolicyService) Update(assetPolicyForm *dtos.AssetPolicyFormDTO, currentUserID uint, assetId uint) ([]uint, error) {
	tx := aps.repo.BeginTransaction()

	var policyIDs []uint
	var deletePolicyIDs []uint

	// Fetch existing asset policies for this asset
	existingPolicies, err := aps.repo.GetByAssetID(assetId, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, documentID := range assetPolicyForm.DocumentIDs {
		if _, err := aps.documentSvc.GetByID(documentID, tx); err != nil {
			tx.Rollback()
			return nil, utility.NewValidationError(fmt.Errorf("document id %d does not exist", documentID), "not_exist", "document id")
		}
	}

	existingPolicyMap := make(map[uint]uint) // Map[DocumentID]AssetPolicyID
	for _, policy := range existingPolicies {
		existingPolicyMap[policy.DocumentID] = policy.ID
	}

	newDocumentIDSet := make(map[uint]bool) // Set of incoming DocumentIDs
	for _, documentID := range assetPolicyForm.DocumentIDs {
		newDocumentIDSet[documentID] = true
	}

	// Identify policies and checklists to delete
	for documentID, policyID := range existingPolicyMap {
		if !newDocumentIDSet[documentID] {
			deletePolicyIDs = append(deletePolicyIDs, policyID)
		}
	}

	// If there are policies to delete, remove them in bulk
	if len(deletePolicyIDs) > 0 {
		// First delete associated checklists
		if err := aps.repo.DeleteAssetPolicyChecklists(deletePolicyIDs, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		// Then delete the asset policies
		if err := aps.repo.DeleteAssetPolicies(deletePolicyIDs, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// ADD new policies (and related checklists) that are in DTO but not in DB
	for _, documentID := range assetPolicyForm.DocumentIDs {
		if _, exists := existingPolicyMap[documentID]; !exists {
			assetPolicy := &models.AssetPolicy{
				AssetRegisterID: assetId,
				DocumentID:      documentID,
				CreatedBYID:     currentUserID,
			}

			policyID, err := aps.repo.Create(assetPolicy, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			// Fetch the document and its requirement groups
			document, err := aps.documentSvc.GetFullDetailsByID(documentID, &currentUserID, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			// Store checklists in batch
			var checklists []*models.AssetPolicyChecklist

			// Loop through DocumentRequirementGroups
			for _, requirementGroup := range document.DocumentRequirementGroups {
				// Loop through DocumentRequirements
				for _, requirement := range requirementGroup.DocumentRequirements {
					checklists = append(checklists, &models.AssetPolicyChecklist{
						AssetPolicyID:              policyID,
						DocumentRequirementGroupID: requirementGroup.ID,
						QuestionGroupTitle:         requirementGroup.TitleEn,
						QuestionTitle:              requirement.TitleEn,
						CreatedBYID:                currentUserID,
					})
				}
			}

			// Save all checklists in bulk
			if len(checklists) > 0 {
				if err := aps.repo.CreateAssetPolicyChecklists(checklists, tx); err != nil {
					tx.Rollback()
					return nil, err
				}
			}

			policyIDs = append(policyIDs, policyID)
		}
	}

	tx.Commit()

	// Log Activity in Background
	go func() {
		for _, policyID := range policyIDs {
			existingPolicy, err := aps.GetFullDetailsByID(assetId, policyID)
			if err != nil {
				utility.LogError("Failed to fetch asset policy", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}
			err = aps.activityLogSvc.LogActivity(
				models.CreateAction,
				models.AssetPolicy{}.TableName(),
				policyID,
				existingPolicy,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.AssetPolicy{}.TableName(),
					"recordID":  policyID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				utility.LogError("Failed to log asset policy activity", fields)
			}
		}
	}()

	return policyIDs, nil
}

func (aps *AssetPolicyService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAssetPolicy []map[string]interface{}

	if len(ids) == 0 {
		return nil, nil
	}

	// Fetch all existing policies one by one
	var existingPolicyIDs []uint
	for _, id := range ids {
		policy, err := aps.repo.GetByID(id)
		if err != nil {
			fmt.Printf("Failed to fetch asset policy with ID %d: %v\n", id, err)
			failedIds = append(failedIds, id)
			continue
		}

		existingPolicyIDs = append(existingPolicyIDs, policy.ID)
		deletedAssetPolicy = append(deletedAssetPolicy, map[string]interface{}{
			"id": policy.ID,
		})
	}

	// If no policies exist for deletion, return early
	if len(existingPolicyIDs) == 0 {
		return failedIds, nil
	}

	tx := aps.repo.BeginTransaction()

	// Bulk Delete AssetPolicyChecklists
	if err := aps.repo.DeleteAssetPolicyChecklists(existingPolicyIDs, tx); err != nil {
		tx.Rollback()
		fmt.Printf("Failed to delete checklists for policies: %v\n", err)
		return failedIds, err
	}
	if err := aps.repo.DeleteAssetPolicyComments(existingPolicyIDs, tx); err != nil {
		tx.Rollback()
		fmt.Printf("Failed to delete asset policies comments: %v\n", err)
		return failedIds, err
	}
	//  Delete AssetPolicies
	if err := aps.repo.DeleteAssetPolicies(existingPolicyIDs, tx); err != nil {
		tx.Rollback()
		fmt.Printf("Failed to delete asset policies: %v\n", err)
		return failedIds, err
	}

	tx.Commit()
	successfulIds = existingPolicyIDs

	// Log Activity in Background
	if len(successfulIds) > 0 {
		go func() {
			err := aps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetPolicy{}.TableName(),
				successfulIds,
				deletedAssetPolicy,
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

func (aps *AssetPolicyService) GetAssetPolicyComplianceSummary(assetID uint) (*dtos.AssetPolicyByComplianceDTO, error) {
	if _, err := aps.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return aps.repo.GetAssetPolicyComplianceSummary(assetID)
}

func (aps *AssetPolicyService) StatusUpdate(id uint, currentUserID uint, assetPolicyStatusUpdateFormDTO *dtos.AssetPolicyStatusUpdateFormDTO, tx utility.DatabaseTransaction) error {
	// Fetch the existing asset policy record
	existingAssetPolicy, err := aps.GetByID(id)
	if err != nil {
		return err
	}

	// Update asset policy status, comments, and percentage
	existingAssetPolicy.ComplianceStatusID = &assetPolicyStatusUpdateFormDTO.ComplianceStatusID
	existingAssetPolicy.UpdatedBYID = currentUserID

	// Update the record in the repository
	err = aps.repo.Update(existingAssetPolicy, tx)
	if err != nil {
		return err
	}

	return nil
}
