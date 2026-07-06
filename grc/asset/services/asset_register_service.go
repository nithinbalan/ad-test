package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetRegisterServiceInterface interface {
	GetAll(params *dtos.AssetRegisterListQueryParamsDTO) (*dtos.AssetRegisterListResponseDTO, error)
	GetByID(id uint) (*models.AssetRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetRegisterDetailDTO, error)
	Create(assetForm *dtos.AssetRegisterFormDTO, currentUserID uint) (uint, error)
	Update(id uint, assetForm *dtos.AssetRegisterFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error) 
	Unarchive(ids []uint, currentUserID uint) ([]uint, error)

	UpdateTotalScore(assetID uint, totalScore float64, tx ...utility.DatabaseTransaction) error
	GetControlSummary(assetID uint) (map[string]interface{}, error)
	UpdateAssetRecoveryObjective(asset *models.AssetRegister, tx utility.DatabaseTransaction) error
	GetDocumentPath(id uint, token string) (string, error)

	ExportAsset(data []*dtos.AssetRegisterListDTO, currentUserID uint) (*os.File, error)
}

type AssetRegisterService struct {
	repo                                  repositories.AssetRegisterRepositoryInterface
	activityLogSvc                        logServices.ActivityLogServiceInterface
	assetTypeSvc                          AssetTypeServiceInterface
	assetCategorySvc                      AssetCategoryServiceInterface
	assetStatusSvc                        AssetStatusServiceInterface
	userSvc                               organizationServices.UserServiceInterface
	assetCriticalityMatrixRatingScaleSvc  AssetCriticalityMatrixRatingScaleServiceInterface
	assetCriticalityMatrixSvc             AssetCriticalityMatrixServiceInterface
	userSettingSvc                        organizationServices.UserSettingServiceInterface
	controlImplementationStatusSvc        organizationServices.ControlImplementationStatusServiceInterface
	controlEfficiencyStatusSvc            organizationServices.ControlEfficiencyStatusServiceInterface
	assetDeviceCategorySvc                AssetDeviceCategoryServiceInterface
	assetCriticalityMatrixFactorRatingSvc AssetCriticalityMatrixFactorRatingServiceInterface
	confidentialityClassificationSvc      ConfidentialityClassificationServiceInterface
	assetEnvironmentSvc                   AssetEnvironmentServiceInterface
}

func NewAssetRegisterService(
	repo repositories.AssetRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetTypeSvc AssetTypeServiceInterface,
	assetCategorySvc AssetCategoryServiceInterface,
	assetStatusSvc AssetStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	assetCriticalityMatrixRatingScaleSvc AssetCriticalityMatrixRatingScaleServiceInterface,
	assetCriticalityMatrixSvc AssetCriticalityMatrixServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	controlImplementationStatusSvc organizationServices.ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusSvc organizationServices.ControlEfficiencyStatusServiceInterface,
	assetDeviceCategorySvc AssetDeviceCategoryServiceInterface,
	assetCriticalityMatrixFactorRatingSvc AssetCriticalityMatrixFactorRatingServiceInterface,
	confidentialityClassificationSvc ConfidentialityClassificationServiceInterface,
	assetEnvironmentSvc AssetEnvironmentServiceInterface,
) AssetRegisterServiceInterface {
	return &AssetRegisterService{
		repo:                                  repo,
		activityLogSvc:                        activityLogSvc,
		assetTypeSvc:                          assetTypeSvc,
		assetCategorySvc:                      assetCategorySvc,
		assetStatusSvc:                        assetStatusSvc,
		userSvc:                               userSvc,
		assetCriticalityMatrixRatingScaleSvc:  assetCriticalityMatrixRatingScaleSvc,
		assetCriticalityMatrixSvc:             assetCriticalityMatrixSvc,
		userSettingSvc:                        userSettingSvc,
		controlImplementationStatusSvc:        controlImplementationStatusSvc,
		controlEfficiencyStatusSvc:            controlEfficiencyStatusSvc,
		assetDeviceCategorySvc:                assetDeviceCategorySvc,
		assetCriticalityMatrixFactorRatingSvc: assetCriticalityMatrixFactorRatingSvc,
		confidentialityClassificationSvc:      confidentialityClassificationSvc,
		assetEnvironmentSvc:                   assetEnvironmentSvc,
	}
}

func (ars *AssetRegisterService) GetAll(params *dtos.AssetRegisterListQueryParamsDTO) (*dtos.AssetRegisterListResponseDTO, error) {
	return ars.repo.GetAll(params)
}

func (ars *AssetRegisterService) GetByID(id uint) (*models.AssetRegister, error) {
	asset, err := ars.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return asset, nil
}

func (as *AssetRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetRegisterDetailDTO, error) {
	var assetDetailDTO dtos.AssetRegisterDetailDTO

	asset, err := as.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(asset, &assetDetailDTO)
	if err != nil {
		return nil, err
	}

	if asset.AcquisitionDate != nil {
		acquisitionDateString := asset.AcquisitionDate.Format("2006-01-02")
		assetDetailDTO.AcquisitionDateString = &acquisitionDateString
	}
	if asset.AssignmentDate != nil {
		assessmentDateString := asset.AssignmentDate.Format("2006-01-02")
		assetDetailDTO.AssignmentDateString = &assessmentDateString
	}
	if asset.WarrantyExpiryDate != nil {
		warrantyExpiryDateString := asset.WarrantyExpiryDate.Format("2006-01-02")
		assetDetailDTO.WarrantyExpiryDateString = &warrantyExpiryDateString
	}
		if asset.InventoryTime != nil {
		inventoryTimeString := asset.InventoryTime.Format("2006-01-02T15:04:05")
		assetDetailDTO.InventoryTimeString = &inventoryTimeString
	}
	if asset.DiscoveryTime != nil {
		discoveryTimeString := asset.DiscoveryTime.Format("2006-01-02T15:04:05")
		assetDetailDTO.DiscoveryTimeString = &discoveryTimeString
	}
	if asset.LastResponseTime != nil {
		lastResponseTimeString := asset.LastResponseTime.Format("2006-01-02T15:04:05")
		assetDetailDTO.LastResponseTimeString = &lastResponseTimeString
	}
	// Add additional fields as necessary for the details screen
	// For example, setting relationship data or additional fields as required by the DTO
	assetDetailDTO.CreatedAt = asset.CreatedAt
	return &assetDetailDTO, nil
}

func (ars *AssetRegisterService) Create(assetRegisterFormDTO *dtos.AssetRegisterFormDTO, currentUserID uint) (uint, error) {
	tx := ars.repo.BeginTransaction()

	if err := ars.validateAssetRegisterFormDTO(assetRegisterFormDTO); err != nil {
		return 0, err
	}
	dates, err := ars.parseAssetRegisterDates(assetRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := ars.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	matrixID, err := ars.assetCriticalityMatrixSvc.GetCriticalityMatrixByCategoryID(assetRegisterFormDTO.AssetCategoryID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if matrixID == 0 {
		matrixID, err = ars.assetCriticalityMatrixSvc.GetDefaultCriticalityMatrixID()
		if err != nil || matrixID == 0 {
			tx.Rollback()
			return 0, utility.NewValidationError(fmt.Errorf("no criticality matrix found for asset category or default"), "not_exist", "criticality matrix")
		}
	}

	organizationUnitID, err := ars.userSvc.GetUserPrimaryUnitID(assetRegisterFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	asset := &models.AssetRegister{
		ReferenceNumber:          referenceNumber,
		TitleEn:                  assetRegisterFormDTO.TitleEn,
		TitleAr:                  assetRegisterFormDTO.TitleAr,
		AssetCategoryID:          assetRegisterFormDTO.AssetCategoryID,
		AssetSubCategoryID:       assetRegisterFormDTO.AssetSubCategoryID,
		AssetStatusID:            assetRegisterFormDTO.AssetStatusID,
		ResponsibleUserID:        assetRegisterFormDTO.ResponsibleUserID,
		ResponsibleUnitID:        organizationUnitID,
		Location:                 assetRegisterFormDTO.Location,
		Description:              assetRegisterFormDTO.Description,
		AssetCriticalityMatrixID: matrixID,
		AcquisitionDate:          dates.AcquisitionDate,
		AssignmentDate:           dates.AssignmentDate,
		WarrantyExpiryDate:       dates.WarrantyExpiryDate,
		Manufacturer:             assetRegisterFormDTO.Manufacturer,
		Model:                    assetRegisterFormDTO.Model,
		SerialNumber:             assetRegisterFormDTO.SerialNumber,

		TechnicalOwnerID:                assetRegisterFormDTO.TechnicalOwnerID,
		AdministrativeOwnerID:           assetRegisterFormDTO.AdministrativeOwnerID,
		Hostname:                        assetRegisterFormDTO.Hostname,
		IPAddress:                       assetRegisterFormDTO.IPAddress,
		MACAddress:                      assetRegisterFormDTO.MACAddress,
		FQDN:                            assetRegisterFormDTO.FQDN,
		DeviceCategoryID:                assetRegisterFormDTO.DeviceCategoryID,
		CriticalityClassificationID:     assetRegisterFormDTO.CriticalityClassificationID,
		ConfidentialityClassificationID: assetRegisterFormDTO.ConfidentialityClassificationID,
		Region:                          assetRegisterFormDTO.Region,
		AssignedUserID:                  assetRegisterFormDTO.AssignedUserID,
		OperatingSystem:                 assetRegisterFormDTO.OperatingSystem,
		AssetEnvironmentID:              assetRegisterFormDTO.AssetEnvironmentID,
		InventoryTime:                   dates.InventoryTime,
		DiscoveryTime:                   dates.DiscoveryTime,
		LastResponseTime:                dates.LastResponseTime,
		IsArchived:                      false,
		CreatedBYID:                     currentUserID,
	}

	assetID, err := ars.repo.Create(asset, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Handle Asset Documents Upload
	var assetDocuments []*models.AssetDocument

	for _, doc := range assetRegisterFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the storage URL
		url := utility.GenerateDocumentURL("assets", "asset-register", strconv.FormatUint(uint64(assetID), 10), "documents")

		// Create a new AssetDocument instance
		document := &models.AssetDocument{
			AssetRegisterID: assetID,
			URL:             url,
			Name:            doc.FileName,
			Size:            doc.Size,
			Token:           encryptedFilename,
			OriginalName:    doc.OriginalName,
			Type:            doc.Type,
			CreatedBYID:     currentUserID,
		}

		// Append to the slice
		assetDocuments = append(assetDocuments, document)
	}

	// Save documents in the database
	if len(assetDocuments) > 0 {
		if err := ars.repo.CreateAssetDocuments(assetDocuments, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		// Store files in the system
		for _, doc := range assetDocuments {
			storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	tx.Commit()

	go func() {
		existingAsset, err := ars.GetFullDetailsByID(assetID)
		if err != nil {
			utility.LogError("Failed to fetch asset register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = ars.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetRegister{}.TableName(),
			assetID,
			existingAsset,
			currentUserID,
			"asset",
			"asset_register",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetRegister{}.TableName(),
				"recordID":  assetID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()

	return assetID, nil
}

func (ars *AssetRegisterService) Update(id uint, assetRegisterFormDTO *dtos.AssetRegisterFormDTO, currentUserID uint) error {
	tx := ars.repo.BeginTransaction()
	existingAsset, err := ars.repo.GetByID(id)
	if err != nil {
		return err
	}

	oldAssetRegister, err := ars.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// oldResponsibleUserID := existingAsset.ResponsibleUserID

	if err := ars.validateAssetRegisterFormDTO(assetRegisterFormDTO); err != nil {
		return err
	}
	dates, err := ars.parseAssetRegisterDates(assetRegisterFormDTO)
	if err != nil {
		return err
	}

	if existingAsset.ResponsibleUserID != assetRegisterFormDTO.ResponsibleUserID {
		organizationUnitID, err := ars.userSvc.GetUserPrimaryUnitID(assetRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "control owner unit")
		}
		existingAsset.ResponsibleUnitID = organizationUnitID
	}

	// Update fields
	existingAsset.TitleEn = assetRegisterFormDTO.TitleEn
	existingAsset.TitleAr = assetRegisterFormDTO.TitleAr
	existingAsset.AssetCategoryID = assetRegisterFormDTO.AssetCategoryID
	existingAsset.AssetSubCategoryID = assetRegisterFormDTO.AssetSubCategoryID
	existingAsset.AssetStatusID = assetRegisterFormDTO.AssetStatusID
	existingAsset.ResponsibleUserID = assetRegisterFormDTO.ResponsibleUserID
	existingAsset.Location = assetRegisterFormDTO.Location
	existingAsset.Description = assetRegisterFormDTO.Description
	existingAsset.AcquisitionDate = dates.AcquisitionDate
	existingAsset.AssignmentDate = dates.AssignmentDate
	existingAsset.WarrantyExpiryDate = dates.WarrantyExpiryDate
	existingAsset.Manufacturer = assetRegisterFormDTO.Manufacturer
	existingAsset.Model = assetRegisterFormDTO.Model
	existingAsset.SerialNumber = assetRegisterFormDTO.SerialNumber

	existingAsset.TechnicalOwnerID = assetRegisterFormDTO.TechnicalOwnerID
	existingAsset.AdministrativeOwnerID = assetRegisterFormDTO.AdministrativeOwnerID
	existingAsset.Hostname = assetRegisterFormDTO.Hostname
	existingAsset.IPAddress = assetRegisterFormDTO.IPAddress
	existingAsset.MACAddress = assetRegisterFormDTO.MACAddress
	existingAsset.FQDN = assetRegisterFormDTO.FQDN
	existingAsset.DeviceCategoryID = assetRegisterFormDTO.DeviceCategoryID
	existingAsset.CriticalityClassificationID = assetRegisterFormDTO.CriticalityClassificationID
	existingAsset.ConfidentialityClassificationID = assetRegisterFormDTO.ConfidentialityClassificationID
	existingAsset.Region = assetRegisterFormDTO.Region
	existingAsset.AssignedUserID = assetRegisterFormDTO.AssignedUserID
	existingAsset.OperatingSystem = assetRegisterFormDTO.OperatingSystem
	existingAsset.AssetEnvironmentID = assetRegisterFormDTO.AssetEnvironmentID
	existingAsset.InventoryTime = dates.InventoryTime
	existingAsset.DiscoveryTime = dates.DiscoveryTime
	existingAsset.LastResponseTime = dates.LastResponseTime

	existingAsset.UpdatedBYID = &currentUserID

	err = ars.repo.Update(existingAsset, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle Asset Document Updates
	var assetDocuments []*models.AssetDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AssetDocument

	for _, doc := range assetRegisterFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate document URL
			url := utility.GenerateDocumentURL("assets", "asset-register", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new AssetDocument instance
			document := &models.AssetDocument{
				AssetRegisterID: id,
				URL:             url,
				Name:            doc.FileName,
				Size:            doc.Size,
				Token:           encryptedFilename,
				OriginalName:    doc.OriginalName,
				Type:            doc.Type,
				CreatedBYID:     currentUserID,
			}

			assetDocuments = append(assetDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := ars.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	// Save new documents
	if len(assetDocuments) > 0 {
		if err := ars.repo.CreateAssetDocuments(assetDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete documents
	if len(deleteDocumentIds) > 0 {
		if err := ars.repo.DeleteAssetDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store new files in the system
	if len(assetDocuments) > 0 {
		for _, doc := range assetDocuments {
			storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(existingAsset.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Remove deleted files from the system
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(existingAsset.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()
	go func() {
		existingAsset, err := ars.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"asset_category":     oldAssetRegister.AssetCategory.TitleEn,
			"asset_sub_category": oldAssetRegister.AssetSubCategory.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", oldAssetRegister.ResponsibleUser.FirstNameEN, oldAssetRegister.ResponsibleUser.LastNameEN),
			"title_en":           oldAssetRegister.TitleEn,
			"title_ar":           oldAssetRegister.TitleAr,
			"description":        oldAssetRegister.Description,
			"asset_status":       oldAssetRegister.AssetStatus.TitleEn,
			"location":           oldAssetRegister.Location,
			"acquisition_date":   oldAssetRegister.AcquisitionDateString,
			"hostname":           oldAssetRegister.Hostname,
			"ip_address":         oldAssetRegister.IPAddress,
			"mac_address":        oldAssetRegister.MACAddress,
			"fqdn":               oldAssetRegister.FQDN,
			"region":             oldAssetRegister.Region,
			"operating_system":   oldAssetRegister.OperatingSystem,
		}

		newData := map[string]interface{}{
			"asset_category":     existingAsset.AssetCategory.TitleEn,
			"asset_sub_category": existingAsset.AssetSubCategory.TitleEn,
			"responsible_user":   fmt.Sprintf("%s %s", existingAsset.ResponsibleUser.FirstNameEN, existingAsset.ResponsibleUser.LastNameEN),
			"title_en":           existingAsset.TitleEn,
			"title_ar":           existingAsset.TitleAr,
			"description":        existingAsset.Description,
			"asset_status":       existingAsset.AssetStatus.TitleEn,
			"location":           existingAsset.Location,
			"acquisition_date":   existingAsset.AcquisitionDateString,
			"hostname":           existingAsset.Hostname,
			"ip_address":         existingAsset.IPAddress,
			"mac_address":        existingAsset.MACAddress,
			"fqdn":               existingAsset.FQDN,
			"region":             existingAsset.Region,
			"operating_system":   existingAsset.OperatingSystem,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ars.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AssetRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset",
				"asset_register",
				"",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AssetRegister{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (as *AssetRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAssetRegister []map[string]interface{}
	documentsToDelete := make(map[uint][]*models.AssetDocument)

	tx := as.repo.BeginTransaction()

	for _, id := range ids {
		asset, err := as.GetFullDetailsByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Store documents for deletion (Type Conversion from DTO to Model)
		if len(asset.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range asset.Documents {
				if doc != nil && doc.ID != nil {
					convertedDoc := &models.AssetDocument{
						ID:           *doc.ID,
						Name:         doc.Name,
						OriginalName: doc.OriginalName,
						Token:        doc.Token,
						Type:         doc.Type,
					}
					documentsToDelete[id] = append(documentsToDelete[id], convertedDoc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}

			if err := as.repo.DeleteAssetDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		// Delete asset record
		err = as.repo.Delete(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		// Store deleted asset details for logging
		deletedAssetRegister = append(deletedAssetRegister, map[string]interface{}{
			"id":       id,
			"title_en": asset.TitleEn,
			"title_ar": asset.TitleAr,
		})
	}

	tx.Commit()

	// Remove documents and folders after successful deletion
	if len(successfulIds) > 0 {
		for assetID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					fmt.Printf("Failed to remove file %s: %v\n", doc.Name, err)
				}
			}
			destinationPath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				fmt.Printf("Failed to remove folder %s: %v\n", destinationPath, err)
			}
		}
	}

	// Log activity in the background
	if len(successfulIds) > 0 {
		go func() {
			err := as.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetRegister{}.TableName(),
				successfulIds,
				deletedAssetRegister,
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

func (ars *AssetRegisterService) UpdateTotalScore(assetID uint, totalScore float64, tx ...utility.DatabaseTransaction) error {
	existingAsset, err := ars.repo.GetByID(assetID)
	if err != nil {
		return err
	}

	if totalScore == 0 {
		existingAsset.AssetCriticalityMatrixRatingScaleID = nil
	} else {
		ratingScale, err := ars.assetCriticalityMatrixRatingScaleSvc.GetRatingScaleByScore(existingAsset.AssetCriticalityMatrixID, totalScore)
		if err != nil {
			return err
		}
		existingAsset.AssetCriticalityMatrixRatingScaleID = &ratingScale.ID

	}
	existingAsset.TotalScore = &totalScore

	err = ars.repo.Update(existingAsset, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (ars *AssetRegisterService) validateAssetRegisterFormDTO(assetRegisterFormDTO *dtos.AssetRegisterFormDTO) error {
	if _, err := ars.assetCategorySvc.GetByID(uint(assetRegisterFormDTO.AssetCategoryID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset category does not exist"), "not_exist", "asset category")
	}

	if _, err := ars.userSvc.GetByID(uint(assetRegisterFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := ars.assetStatusSvc.GetByID(uint(assetRegisterFormDTO.AssetStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset status does not exist"), "not_exist", "asset status")
	}
	if assetRegisterFormDTO.TechnicalOwnerID != nil {
		if _, err := ars.userSvc.GetByID(*assetRegisterFormDTO.TechnicalOwnerID); err != nil {
			return utility.NewValidationError(fmt.Errorf("technical owner does not exist"), "not_exist", "technical_owner")
		}
	}

	if assetRegisterFormDTO.AdministrativeOwnerID != nil {
		if _, err := ars.userSvc.GetByID(*assetRegisterFormDTO.AdministrativeOwnerID); err != nil {
			return utility.NewValidationError(fmt.Errorf("administrative owner does not exist"), "not_exist", "administrative_owner")
		}
	}

	if assetRegisterFormDTO.AssignedUserID != nil {
		if _, err := ars.userSvc.GetByID(*assetRegisterFormDTO.AssignedUserID); err != nil {
			return utility.NewValidationError(fmt.Errorf("assigned user does not exist"), "not_exist", "assigned_user")
		}
	}

	if assetRegisterFormDTO.DeviceCategoryID != nil {
		if _, err := ars.assetDeviceCategorySvc.GetByID(*assetRegisterFormDTO.DeviceCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("device category does not exist"), "not_exist", "device_category")
		}
	}

	if assetRegisterFormDTO.CriticalityClassificationID != nil {
		if _, err := ars.assetCriticalityMatrixFactorRatingSvc.GetByID(0, *assetRegisterFormDTO.CriticalityClassificationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("criticality classification does not exist"), "not_exist", "criticality_classification")
		}
	}

	if assetRegisterFormDTO.ConfidentialityClassificationID != nil {
		if _, err := ars.confidentialityClassificationSvc.GetByID(*assetRegisterFormDTO.ConfidentialityClassificationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("confidentiality classification does not exist"), "not_exist", "confidentiality_classification")
		}
	}

	if assetRegisterFormDTO.AssetEnvironmentID != nil {
		if _, err := ars.assetEnvironmentSvc.GetByID(*assetRegisterFormDTO.AssetEnvironmentID); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset environment does not exist"), "not_exist", "asset_environment")
		}
	}

	return nil
}

func (ars *AssetRegisterService) parseAssetRegisterDates(assetRegisterFormDTO *dtos.AssetRegisterFormDTO) (*dtos.AssetRegisterParsedDates, error) {
	parsed := &dtos.AssetRegisterParsedDates{}

	var err error
	if assetRegisterFormDTO.AcquisitionDate != nil {
		parsed.AcquisitionDate, err = utility.ParseDate(*assetRegisterFormDTO.AcquisitionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid acquisition date"), "invalid_date", "acquisition_date")
		}
	}

	if assetRegisterFormDTO.AssignmentDate != nil {
		parsed.AssignmentDate, err = utility.ParseDate(*assetRegisterFormDTO.AssignmentDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid assignment date"), "invalid_date", "assignment_date")
		}
	}

	if assetRegisterFormDTO.WarrantyExpiryDate != nil {
		parsed.WarrantyExpiryDate, err = utility.ParseDate(*assetRegisterFormDTO.WarrantyExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid warranty expiry date"), "invalid_date", "warranty_expiry_date")
		}
	}

	if assetRegisterFormDTO.InventoryTime != nil {
		parsed.InventoryTime, err = utility.ParseDateTime(*assetRegisterFormDTO.InventoryTime)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid inventory time"), "invalid_date", "inventory_time")
		}
	}

	if assetRegisterFormDTO.DiscoveryTime != nil {
		parsed.DiscoveryTime, err = utility.ParseDateTime(*assetRegisterFormDTO.DiscoveryTime)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid discovery time"), "invalid_date", "discovery_time")
		}
	}

	if assetRegisterFormDTO.LastResponseTime != nil {
		parsed.LastResponseTime, err = utility.ParseDateTime(*assetRegisterFormDTO.LastResponseTime)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last response time"), "invalid_date", "last_response_time")
		}
	}

	return parsed, nil
}

func (ars *AssetRegisterService) GetControlSummary(assetID uint) (map[string]interface{}, error) {
	assetControls, err := ars.repo.GetControlSummary(assetID)
	if err != nil {
		return nil, err
	}

	if len(assetControls) == 0 {
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

	for _, ac := range assetControls {
		// Safe nil check for Control
		if ac.Control == nil {
			continue
		}

		// Implementation status
		if ac.Control.ControlImplementationStatus != nil {
			implStatus := ac.Control.ControlImplementationStatus.Type
			implementationStatusCount[implStatus]++
			implementationDetails[implStatus] = *ac.Control.ControlImplementationStatus
		}

		// Efficiency status
		if ac.Control.ControlEfficiencyStatus != nil {
			effStatus := ac.Control.ControlEfficiencyStatus.Type
			efficiencyStatusCount[effStatus]++
			efficiencyDetails[effStatus] = *ac.Control.ControlEfficiencyStatus
		}
	}

	totalControls := len(assetControls)

	// Avoid division by zero
	totalControlsFloat := float64(totalControls)
	if totalControls == 0 {
		totalControlsFloat = 1
	}

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
		"total_control_score":       (float64(implementationStatusCount["implemented"]) / totalControlsFloat) * 100,
		"implementation_status":     implementationStatus,
		"efficiency_percentage":     (float64(efficiencyStatusCount["efficient"]) / totalControlsFloat) * 100,
		"implementation_percentage": (float64(implementationStatusCount["implemented"]) / totalControlsFloat) * 100,
		"total_efficiency":          efficiencyStatus,
		"total_control_count":       totalControls,
	}

	return summary, nil
}

func (ars *AssetRegisterService) UpdateAssetRecoveryObjective(assetRegister *models.AssetRegister, tx utility.DatabaseTransaction) error {
	err := ars.repo.Update(assetRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ars *AssetRegisterService) GetDocumentPath(id uint, token string) (string, error) {
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (ars *AssetRegisterService) ExportAsset(data []*dtos.AssetRegisterListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := ars.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"الرقم التسلسلي",
			"الرقم المرجعي",
			"العنوان",
			"الفئة",
			"حالة الأصل",
			"الموقع",
			"المستخدم المسؤول",
			"تاريخ الاقتناء",
			"المالك الفني",
			"المالك الإداري",
			"المستخدم المعين",
			"اسم المضيف",
			"عنوان IP",
			"عنوان MAC",
			"FQDN",
			"المنطقة",
			"نظام التشغيل",
			"فئة الجهاز",
			"تصنيف السرية",
			"بيئة الأصل",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Category",
			"Asset Status",
			"Location",
			"Responsible User",
			"Acquisition Date",
			"Technical Owner",
			"Administrative Owner",
			"Assigned User",
			"Hostname",
			"IP Address",
			"MAC Address",
			"FQDN",
			"Region",
			"Operating System",
			"Device Category",
			"Confidentiality Classification",
			"Asset Environment",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, asset := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				asset.ReferenceNumber,
				asset.TitleAr,
				utility.GetStringValue(asset.AssetCategoryTitleAr),
				asset.AssetStatusTitleAr,
				utility.GetStringValue(asset.Location),
				fmt.Sprintf("%s %s", asset.ResponsibleUserFirstNameAr, asset.ResponsibleUserLastNameAr),
				asset.AcquisitionDate,
				fmt.Sprintf("%s %s",
					utility.GetStringValue(asset.TechnicalOwnerFirstNameAr),
					utility.GetStringValue(asset.TechnicalOwnerLastNameAr)),
				fmt.Sprintf("%s %s",
					utility.GetStringValue(asset.AdministrativeOwnerFirstNameAr),
					utility.GetStringValue(asset.AdministrativeOwnerLastNameAr)),
				fmt.Sprintf("%s %s",
					utility.GetStringValue(asset.AssignedUserFirstNameAr),
					utility.GetStringValue(asset.AssignedUserLastNameAr)),
				utility.GetStringValue(asset.Hostname),
				utility.GetStringValue(asset.IPAddress),
				utility.GetStringValue(asset.MACAddress),
				utility.GetStringValue(asset.FQDN),
				utility.GetStringValue(asset.Region),
				utility.GetStringValue(asset.OperatingSystem),
				utility.GetStringValue(asset.DeviceCategoryTitleAr),
				utility.GetStringValue(asset.ConfidentialityClassificationTitleAr),
				utility.GetStringValue(asset.AssetEnvironmentTitleAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				asset.ReferenceNumber,
				asset.TitleEn,
				utility.GetStringValue(asset.AssetCategoryTitleEn),
				asset.AssetStatusTitleEn,
				utility.GetStringValue(asset.Location),
				fmt.Sprintf("%s %s", asset.ResponsibleUserFirstNameEn, asset.ResponsibleUserLastNameEn),
				asset.AcquisitionDate,
				fmt.Sprintf("%s %s",
					utility.GetStringValue(asset.TechnicalOwnerFirstNameEn),
					utility.GetStringValue(asset.TechnicalOwnerLastNameEn)),
				fmt.Sprintf("%s %s",
					utility.GetStringValue(asset.AdministrativeOwnerFirstNameEn),
					utility.GetStringValue(asset.AdministrativeOwnerLastNameEn)),
				fmt.Sprintf("%s %s",
					utility.GetStringValue(asset.AssignedUserFirstNameEn),
					utility.GetStringValue(asset.AssignedUserLastNameEn)),
				utility.GetStringValue(asset.Hostname),
				utility.GetStringValue(asset.IPAddress),
				utility.GetStringValue(asset.MACAddress),
				utility.GetStringValue(asset.FQDN),
				utility.GetStringValue(asset.Region),
				utility.GetStringValue(asset.OperatingSystem),
				utility.GetStringValue(asset.DeviceCategoryTitleEn),
				utility.GetStringValue(asset.ConfidentialityClassificationTitleEn),
				utility.GetStringValue(asset.AssetEnvironmentTitleEn),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (ars *AssetRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ars.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := ars.repo.ForceDelete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		return ids, err
	}

	go func() {
		if err := ars.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.AssetRegister{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		); err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}
func (ars *AssetRegisterService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ars.repo.BeginTransaction() 

	var failedIds []uint
	var successfulIds []uint
	var archivedAssetRegister []map[string]interface{}

	for _, id := range ids {
		asset, err := ars.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue 
		}

		if err := ars.repo.Archive(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		archivedAssetRegister = append(archivedAssetRegister, map[string]interface{}{
			"id":       id,
			"title_en": asset.TitleEn,
			"title_ar": asset.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback() 
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		utility.LogError("Failed to commit asset archive transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ars.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.AssetRegister{}.TableName(),
				successfulIds,
				archivedAssetRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}
func (ars *AssetRegisterService) Unarchive(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ars.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var unarchivedAssetRegister []map[string]interface{}

	for _, id := range ids {
		asset, err := ars.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := ars.repo.Unarchive(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		unarchivedAssetRegister = append(unarchivedAssetRegister, map[string]interface{}{
			"id":       id,
			"title_en": asset.TitleEn,
			"title_ar": asset.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		utility.LogError("Failed to commit asset unarchive transaction", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ars.activityLogSvc.LogBulkActivities(
				models.UnarchiveAction,
				models.AssetRegister{}.TableName(),
				successfulIds,
				unarchivedAssetRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk unarchive activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}