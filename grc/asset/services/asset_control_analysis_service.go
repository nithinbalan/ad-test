package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// Service interface
type AssetControlAnalysisServiceInterface interface {
	GetAll(assetID uint, params *dtos.AssetControlAnalysisListQueryParamsDTO) ([]*dtos.AssetControlAnalysisListDTO, error)
	GetByID(assetID uint, id uint) (*models.AssetControl, error)
	GetFullDetailsByID(assetID uint, id uint) (*dtos.AssetControlAnalysisDetailDTO, error)
	Create(assetID uint, currentUserID uint, formDTO *dtos.AssetControlFormDTO) error
	Update(assetID uint, id uint, currentUserID uint, updateFormDTO *dtos.AssetControlUpdateFormDTO) error
	Delete(assetID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetControlSummary(assetID uint) (map[string]interface{}, error)
}

// Service implementation
type AssetControlAnalysisService struct {
	repo                           repositories.AssetControlAnalysisRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	controlSvc                     organizationServices.ControlServiceInterface
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface
	controlEfficiencyStatusSvc     organizationServices.ControlEfficiencyStatusServiceInterface
	assetRegisterSvc               AssetRegisterServiceInterface
}

func NewAssetControlAnalysisService(
	repo repositories.AssetControlAnalysisRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusSvc organizationServices.ControlEfficiencyStatusServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
) AssetControlAnalysisServiceInterface {
	return &AssetControlAnalysisService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusSvc:     controlEfficiencyStatusSvc,
		assetRegisterSvc:               assetRegisterSvc,
	}
}

func (acas *AssetControlAnalysisService) GetAll(assetID uint, params *dtos.AssetControlAnalysisListQueryParamsDTO) ([]*dtos.AssetControlAnalysisListDTO, error) {
	return acas.repo.GetAll(assetID, params)
}

func (acas *AssetControlAnalysisService) GetByID(assetID uint, id uint) (*models.AssetControl, error) {
	assetControlAnalysis, err := acas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assetControlAnalysis, nil
}

func (acas *AssetControlAnalysisService) GetFullDetailsByID(assetID uint, id uint) (*dtos.AssetControlAnalysisDetailDTO, error) {
	var assetControlAnalysisDTO dtos.AssetControlAnalysisDetailDTO

	assetControlAnalysis, err := acas.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(assetControlAnalysis, &assetControlAnalysisDTO)
	if err != nil {
		return nil, err
	}

	return &assetControlAnalysisDTO, nil
}

func (acas *AssetControlAnalysisService) Create(assetID uint, currentUserID uint, formDTO *dtos.AssetControlFormDTO) error {
	tx := acas.repo.BeginTransaction()

	// Ensure at least one control is provided
	if len(formDTO.ControlIds) == 0 {
		return errors.New("at least one control must be selected")
	}

	if err := acas.validateAssetControlFormDTO(formDTO); err != nil {
		return err
	}

	defaultImplementationStatus, err := acas.controlImplementationStatusSvc.GetByType("not-implemented")
	if err != nil {
		return err
	}

	defaultEfficiencyStatus, err := acas.controlEfficiencyStatusSvc.GetByType("inefficient")
	if err != nil {
		return err
	}

	var controlAnalysis []*models.AssetControl
	var ids []uint

	for _, controlID := range formDTO.ControlIds {

		assetControlAnalysis := &models.AssetControl{
			AssetRegisterID:               assetID,
			ControlID:                     controlID,
			ControlImplementationStatusID: &defaultImplementationStatus.ID,
			ControlEfficiencyStatusID:     &defaultEfficiencyStatus.ID,
			CreatedBYID:                   currentUserID,
		}

		createdID, err := acas.repo.Create(assetControlAnalysis)
		if err != nil {
			tx.Rollback()
			return err
		}

		controlAnalysis = append(controlAnalysis, assetControlAnalysis)
		ids = append(ids, createdID)
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		err := acas.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.AssetControl{}.TableName(),
			formDTO.ControlIds,
			controlAnalysis,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetControl{}.TableName(),
				"recordID":  formDTO.ControlIds,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (acas *AssetControlAnalysisService) Update(assetID uint, id uint, currentUserID uint, updateFormDTO *dtos.AssetControlUpdateFormDTO) error {
	tx := acas.repo.BeginTransaction()

	if err := acas.validateAssetControlUpdateFormDTO(updateFormDTO); err != nil {
		return err
	}

	existingAssetControlAnalysis, err := acas.GetByID(assetID, id)
	if err != nil {
		return err
	}

	oldControlAnalysis, err := acas.GetFullDetailsByID(assetID, id)
	if err != nil {
		return err
	}

	// Update fields based on the DTO
	existingAssetControlAnalysis.ControlEfficiencyStatusID = &updateFormDTO.ControlEfficiencyStatusID
	existingAssetControlAnalysis.ControlImplementationStatusID = &updateFormDTO.ControlImplementationStatusID
	existingAssetControlAnalysis.Comments = updateFormDTO.Comments
	existingAssetControlAnalysis.UpdatedByID = &currentUserID

	// Save changes
	err = acas.repo.Update(existingAssetControlAnalysis, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.AssetControlDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.AssetControlDocument

	fmt.Println(updateFormDTO.Documents)

	// Iterate over the provided documents
	for _, doc := range updateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				return err
			}

			// Generate document URL
			url := utility.GenerateDocumentURL("asset", "asset-control", strconv.FormatUint(uint64(id), 10), "documents")

			// Create document model
			document := &models.AssetControlDocument{
				AssetControlID: id,
				URL:            url,
				Name:           doc.FileName,
				Size:           doc.Size,
				Token:          encryptedFilename,
				OriginalName:   doc.OriginalName,
				Type:           doc.Type,
				CreatedByID:    currentUserID,
			}

			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle file deletion
			if doc.ID != nil {
				existingDocument, err := acas.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			}
		}
	}

	// Save new documents if any
	if len(newDocuments) > 0 {
		if err := acas.repo.CreateDocuments(newDocuments, tx); err != nil {
			return err
		}
	}

	// Delete specified documents
	if len(deleteDocumentIDs) > 0 {
		if err := acas.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			return err
		}
	}

	// Store uploaded files
	if len(newDocuments) > 0 {
		for _, doc := range newDocuments {
			storagePath := filepath.Join("asset", "asset-control", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				return err
			}
		}
	}

	// Remove deleted files from storage
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("asset", "asset-control", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				return err
			}
		}
	}

	tx.Commit()

	// Fetch updated data for logging
	existingControlAnalysis, err := acas.GetFullDetailsByID(assetID, id)
	if err != nil {
		return err
	}

	// Log changes asynchronously
	go func() {
		changedData := utility.CompareFieldChanges(map[string]interface{}{
			"control_efficiency_status":     oldControlAnalysis.ControlEfficiencyStatusID,
			"control_implementation_status": oldControlAnalysis.ControlImplementationStatusID,
			"comments":                      oldControlAnalysis.Comments,
		}, map[string]interface{}{
			"control_efficiency_status":     existingControlAnalysis.ControlEfficiencyStatusID,
			"control_implementation_status": existingControlAnalysis.ControlImplementationStatusID,
			"comments":                      existingControlAnalysis.Comments,
		})

		if len(changedData) > 0 {
			err := acas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AssetControl{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AssetControl{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (acas *AssetControlAnalysisService) Delete(assetID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := acas.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedControlAnalysis []map[string]interface{}

	for _, id := range ids {
		assetControlAnalysis, err := acas.GetByID(assetID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = acas.repo.Delete(assetControlAnalysis.ID, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedControlAnalysis = append(deletedControlAnalysis, map[string]interface{}{
			"id":       id,
			"asset_id": assetID,
		})
	}

	tx.Commit()

	go func() {
		err := acas.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.AssetControl{}.TableName(),
			successfulIds,
			deletedControlAnalysis,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log bulk delete activity", nil)
		}
	}()

	return failedIds, nil
}

func (acas *AssetControlAnalysisService) validateAssetControlFormDTO(assetControlFormDTO *dtos.AssetControlFormDTO) error {

	controlIds := assetControlFormDTO.ControlIds

	controlIDs := make(map[uint]bool)

	for _, controlId := range controlIds {
		if controlIDs[controlId] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", controlId), "duplicate", "control")
		}

		if _, err := acas.controlSvc.GetByID(controlId); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlId), "not_exist", "control")
		}

		controlIDs[controlId] = true
	}

	return nil
}

func (acas *AssetControlAnalysisService) validateAssetControlUpdateFormDTO(assetControlUpdateFormDTO *dtos.AssetControlUpdateFormDTO) error {

	if _, err := acas.controlImplementationStatusSvc.GetByID(uint(assetControlUpdateFormDTO.ControlImplementationStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control implementation status does not exist"), "not_exist", "control implementation status")
	}

	if _, err := acas.controlImplementationStatusSvc.GetByID(uint(assetControlUpdateFormDTO.ControlEfficiencyStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control efficiency status does not exist"), "not_exist", "control efficiency status")
	}

	return nil
}

func (acas *AssetControlAnalysisService) GetControlSummary(assetID uint) (map[string]interface{}, error) {
	assetControlAnalyses, err := acas.repo.GetControlSummary(assetID)
	if err != nil {
		return nil, err
	}

	if len(assetControlAnalyses) == 0 {
		return map[string]interface{}{
			"total_control_score":       0,
			"implementation_status":     []map[string]interface{}{},
			"efficiency_percentage":     0,
			"implementation_percentage": 0,
			"total_efficiency":          []map[string]interface{}{},
			"total_control_count":       0,
		}, nil
	}

	implementationStatusCount := make(map[string]int)
	efficiencyStatusCount := make(map[string]int)
	implementationDetails := make(map[string]models.ControlImplementationStatus)
	efficiencyDetails := make(map[string]models.ControlEfficiencyStatus)

	for _, aca := range assetControlAnalyses {
		implStatus := aca.ControlImplementationStatus.Type
		implementationStatusCount[implStatus]++
		implementationDetails[implStatus] = *aca.ControlImplementationStatus

		effStatus := aca.ControlEfficiencyStatus.Type
		efficiencyStatusCount[effStatus]++
		efficiencyDetails[effStatus] = *aca.ControlEfficiencyStatus
	}

	totalControls := len(assetControlAnalyses)

	implementationStatus := make([]map[string]interface{}, 0)
	for status, count := range implementationStatusCount {
		implDetail := implementationDetails[status]
		implementationStatus = append(implementationStatus, map[string]interface{}{
			"id":         implDetail.ID,
			"type":       implDetail.Type,
			"title_en":   implDetail.TitleEn,
			"title_ar":   implDetail.TitleAr,
			"count":      count,
			"color_code": implDetail.ColorCode,
		})
	}

	efficiencyStatus := make([]map[string]interface{}, 0)
	for status, count := range efficiencyStatusCount {
		effDetail := efficiencyDetails[status]
		efficiencyStatus = append(efficiencyStatus, map[string]interface{}{
			"id":         effDetail.ID,
			"type":       effDetail.Type,
			"title_en":   effDetail.TitleEn,
			"title_ar":   effDetail.TitleAr,
			"count":      count,
			"color_code": effDetail.ColorCode,
		})
	}

	summary := map[string]interface{}{
		"total_control_score":       (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"implementation_status":     implementationStatus,
		"efficiency_percentage":     (float64(efficiencyStatusCount["efficient"]) / float64(totalControls)) * 100,
		"implementation_percentage": (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"total_efficiency":          efficiencyStatus,
		"total_control_count":       totalControls,
	}

	return summary, nil
}
