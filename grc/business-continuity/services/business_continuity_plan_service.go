package services

import (
	"errors"
	"fmt"
	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BusinessContinuityPlanServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessContinuityPlanDetailDTO, error)
	GetByID(id uint) (*models.BusinessContinuityPlan, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BusinessContinuityPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, businessContinuityPlan *businessContinuityDTOs.BusinessContinuityPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	UpdateRecoveryObjectives(businessContinuityPlan *models.BusinessContinuityPlan, tx utility.DatabaseTransaction) error
	ExportBusinessContinuityPlan(data []*businessContinuityDTOs.BusinessContinuityPlanListDTO, currentUserID uint) (*os.File, error)

	GetDocumentPath(id uint, token string) (string, error)
}

// Service and implementations
type BusinessContinuityPlanService struct {
	repo                businessContinuityPlanRepositories.BusinessContinuityPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	businessFunctionSvc BusinessFunctionServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	bcpStatusSvc        BcpStatusServiceInterface
}

func NewBusinessContinuityPlanService(
	repo businessContinuityPlanRepositories.BusinessContinuityPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	businessFunctionSvc BusinessFunctionServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	bcpStatusSvc BcpStatusServiceInterface,
) BusinessContinuityPlanServiceInterface {
	return &BusinessContinuityPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		businessFunctionSvc: businessFunctionSvc,
		userSettingSvc:      userSettingSvc,
		bcpStatusSvc:        bcpStatusSvc,
	}
}

func (bcps *BusinessContinuityPlanService) GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error) {
	return bcps.repo.GetAll(params)
}

func (bcps BusinessContinuityPlanService) GetByID(id uint) (*models.BusinessContinuityPlan, error) {
	businessContinuityPlan, err := bcps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return businessContinuityPlan, nil
}

func (bcps BusinessContinuityPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessContinuityPlanDetailDTO, error) {
	// Initialize your DTO
	var businessContinuityPlanDTO businessContinuityDTOs.BusinessContinuityPlanDetailDTO

	businessContinuityPlan, err := bcps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(businessContinuityPlan, &businessContinuityPlanDTO)
	if err != nil {
		return nil, err
	}

	businessContinuityPlanDTO.CreatedAt = businessContinuityPlan.CreatedAt
	businessContinuityPlanDTO.UpdatedAt = businessContinuityPlan.UpdatedAt

	return &businessContinuityPlanDTO, nil
}

func (bcps *BusinessContinuityPlanService) Create(currentUserID uint, businessContinuityPlanFormDTO *businessContinuityDTOs.BusinessContinuityPlanFormDTO) (uint, error) {

	tx := bcps.repo.BeginTransaction()

	exists, err := bcps.repo.GetBusinessContinuityPlanByBusinessFunctionID(businessContinuityPlanFormDTO.BusinessFunctionID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if exists != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("the business function already exists"), "unique", "business function")
	}

	referenceNumber, err := bcps.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	bcpStatus, err := bcps.bcpStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	responsibleUnitID, err := bcps.userSvc.GetUserPrimaryUnitID(*businessContinuityPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}

	businessContinuityPlan := &models.BusinessContinuityPlan{
		ReferenceNumber:    referenceNumber,
		BusinessFunctionID: businessContinuityPlanFormDTO.BusinessFunctionID,
		Version:            businessContinuityPlanFormDTO.Version,
		Purpose:            businessContinuityPlanFormDTO.Purpose,
		Scope:              businessContinuityPlanFormDTO.Scope,
		Objectives:         businessContinuityPlanFormDTO.Objectives,
		Assumptions:        businessContinuityPlanFormDTO.Assumptions,
		BcpStatusID:        bcpStatus.ID,
		ResponsibleUserID:  *businessContinuityPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:  &responsibleUnitID,
		CreatedBYID:        currentUserID,
	}
	id, err := bcps.repo.Create(businessContinuityPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	businessFunction, err := bcps.businessFunctionSvc.GetByID(businessContinuityPlanFormDTO.BusinessFunctionID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	businessFunction.BcpStatus = true

	if err := bcps.businessFunctionSvc.UpdateBusinessFunction(businessFunction, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	businessFunctionDetails, err := bcps.businessFunctionSvc.GetFullDetailsByID(businessContinuityPlanFormDTO.BusinessFunctionID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, existingProcess := range businessFunctionDetails.Processes {
		bcpProcess := &models.BcpProcess{
			ProcessID:                existingProcess.ProcessID,
			BusinessContinuityPlanID: id,
			RiskLevelID:              existingProcess.Process.RiskLevelId,
			Rto:                      existingProcess.Process.Rto,
			RtoUnitId:                existingProcess.Process.RtoUnitId,
			Mtpd:                     existingProcess.Process.Mtpd,
			MtpdUnitId:               existingProcess.Process.MtpdUnitId,
			Rpo:                      existingProcess.Process.Rpo,
			RpoUnitId:                existingProcess.Process.RpoUnitId,
			CreatedBYID:              currentUserID,
		}
		err := bcps.repo.CreateBcpProcesses(bcpProcess, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var documents []*models.BcpDocument

	// Iterate over each document in the array
	for _, doc := range businessContinuityPlanFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("business-continuity", "business-continuity-plan", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new ComplianceRequirementStatusUpdateDocument instance
		document := &models.BcpDocument{
			BusinessContinuityPlanID: id,
			URL:                      url,
			Name:                     doc.FileName,
			Size:                     doc.Size,
			Token:                    encryptedFilename,
			OriginalName:             doc.OriginalName,
			Type:                     doc.Type,
			CreatedByID:              currentUserID,
		}
		// Append the document to the slice
		documents = append(documents, document)

	}
	if len(documents) > 0 {
		// Save the ComplianceRequirementStatusUpdateDocument
		if err := bcps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("business-continuity", "business-continuity-plan", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

	}

	tx.Commit()
	go func() {
		existingBusinessContinuityPlan, err := bcps.GetFullDetailsByID(businessFunction.ID)
		if err != nil {
			utility.LogError("Failed to fetch business continuity plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = bcps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BusinessContinuityPlan{}.TableName(),
			id,
			existingBusinessContinuityPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BusinessContinuityPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (bcps *BusinessContinuityPlanService) Update(id uint, currentUserID uint, businessContinuityPlanFormDTO *businessContinuityDTOs.BusinessContinuityPlanFormDTO) error {
	tx := bcps.repo.BeginTransaction()

	// Check for uniqueness of the BusinessFunctionID
	duplicatePlan, err := bcps.repo.GetBusinessContinuityPlanByBusinessFunctionID(businessContinuityPlanFormDTO.BusinessFunctionID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if duplicatePlan != nil && duplicatePlan.ID != id {
		tx.Rollback()
		return utility.NewValidationError(
			fmt.Errorf("a business continuity plan with this business function already exists"),
			"unique",
			"business function",
		)
	}

	existingBusinessContinuityPlan, err := bcps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldBusinessContinuityPlan, err := bcps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingBusinessContinuityPlan.ResponsibleUserID != *businessContinuityPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := bcps.userSvc.GetUserPrimaryUnitID(existingBusinessContinuityPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingBusinessContinuityPlan.ResponsibleUnitID = &responsibleUnitID
	}

	existingBusinessContinuityPlan.BusinessFunctionID = businessContinuityPlanFormDTO.BusinessFunctionID
	existingBusinessContinuityPlan.Version = businessContinuityPlanFormDTO.Version
	existingBusinessContinuityPlan.Purpose = businessContinuityPlanFormDTO.Purpose
	existingBusinessContinuityPlan.Scope = businessContinuityPlanFormDTO.Scope
	existingBusinessContinuityPlan.Objectives = businessContinuityPlanFormDTO.Objectives
	existingBusinessContinuityPlan.Assumptions = businessContinuityPlanFormDTO.Assumptions
	existingBusinessContinuityPlan.ResponsibleUserID = *businessContinuityPlanFormDTO.ResponsibleUserID
	existingBusinessContinuityPlan.UpdatedBYID = &currentUserID

	err = bcps.repo.Update(existingBusinessContinuityPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if businessContinuityPlanFormDTO.BusinessFunctionID != oldBusinessContinuityPlan.BusinessFunctionID {
		if len(oldBusinessContinuityPlan.Processes) > 0 {
			var deleteprocessIDs []*uint
			for _, process := range oldBusinessContinuityPlan.Processes {
				deleteprocessIDs = append(deleteprocessIDs, process.ID)
			}
			if err := bcps.repo.DeleteBcpProcesses(deleteprocessIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
		businessFunctionDetails, err := bcps.businessFunctionSvc.GetFullDetailsByID(businessContinuityPlanFormDTO.BusinessFunctionID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, existingProcess := range businessFunctionDetails.Processes {
			bcpProcess := &models.BcpProcess{
				ProcessID:                existingProcess.ProcessID,
				BusinessContinuityPlanID: id,
				RiskLevelID:              existingProcess.Process.RiskLevelId,
				Rto:                      existingProcess.Process.Rto,
				RtoUnitId:                existingProcess.Process.RtoUnitId,
				Mtpd:                     existingProcess.Process.Mtpd,
				MtpdUnitId:               existingProcess.Process.MtpdUnitId,
				Rpo:                      existingProcess.Process.Rpo,
				RpoUnitId:                existingProcess.Process.RpoUnitId,
				CreatedBYID:              currentUserID,
			}
			err := bcps.repo.CreateBcpProcesses(bcpProcess, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	var documents []*models.BcpDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.BcpDocument

	for _, doc := range businessContinuityPlanFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("business-continuity", "business-continuity-plan", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.BcpDocument{
				BusinessContinuityPlanID: id,
				URL:                      url,
				Name:                     doc.FileName,
				Size:                     doc.Size,
				Token:                    encryptedFileName,
				OriginalName:             doc.OriginalName,
				Type:                     doc.Type,
				CreatedByID:              currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := bcps.repo.GetDocumentByID(*doc.ID)
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

	if len(documents) > 0 {
		if err := bcps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := bcps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("business-continuity", "business-continuity-plan", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "business-continuity-plan", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingBusinessContinuityPlan, err := bcps.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"business_function": oldBusinessContinuityPlan.BusinessFunction.Title,
			"version":           oldBusinessContinuityPlan.Version,
			"scope":             oldBusinessContinuityPlan.Scope,
			"purpose":           oldBusinessContinuityPlan.Purpose,
			"objectives":        oldBusinessContinuityPlan.Objectives,
			"assumptions":       oldBusinessContinuityPlan.Assumptions,
		}

		newData := map[string]interface{}{
			"business_function": existingBusinessContinuityPlan.BusinessFunction.Title,
			"version":           existingBusinessContinuityPlan.Version,
			"scope":             existingBusinessContinuityPlan.Scope,
			"purpose":           existingBusinessContinuityPlan.Purpose,
			"objectives":        existingBusinessContinuityPlan.Objectives,
			"assumptions":       existingBusinessContinuityPlan.Assumptions,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := bcps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BusinessContinuityPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BusinessContinuityPlan{}.TableName(),
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

func (bcps *BusinessContinuityPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := bcps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBusinessContinuityPlan []map[string]interface{}

	documentsToDelete := make(map[uint][]*businessContinuityDTOs.BcpDocumentDTO)

	for _, id := range ids {
		businessContinuityPlan, err := bcps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBusinessContinuityPlan = append(deletedBusinessContinuityPlan, map[string]interface{}{
			"id":                id,
			"business_function": businessContinuityPlan.BusinessFunction.Title,
			"version":           businessContinuityPlan.Version,
		})

		if len(businessContinuityPlan.Processes) > 0 {
			var deleteprocessIDs []*uint
			for _, process := range businessContinuityPlan.Processes {
				deleteprocessIDs = append(deleteprocessIDs, process.ID)
			}
			if err := bcps.repo.DeleteBcpProcesses(deleteprocessIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete each Compliance document if documents exist
		if len(businessContinuityPlan.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range businessContinuityPlan.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
			}
			if err := bcps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the business
		if err := bcps.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := bcps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BusinessContinuityPlan{}.TableName(),
				successfulIds,
				deletedBusinessContinuityPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

		}()
	}

	return failedIds, nil
}

func (bfs *BusinessContinuityPlanService) UpdateRecoveryObjectives(businessContinuityPlan *models.BusinessContinuityPlan, tx utility.DatabaseTransaction) error {

	err := bfs.repo.Update(businessContinuityPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (bfs *BusinessContinuityPlanService) ExportBusinessContinuityPlan(data []*businessContinuityDTOs.BusinessContinuityPlanListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := bfs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"RTO",
			"MTPD",
			"RPO",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"RTO",
			"MTPD",
			"RPO",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, bcp := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				bcp.ReferenceNumber,
				bcp.BusinessFunctionTitle,
				utility.GetIntValueAsString(bcp.Rto) + " " + utility.GetStringValue(bcp.RtoUnitTitleAr),
				utility.GetIntValueAsString(bcp.Mtpd) + " " + utility.GetStringValue(bcp.MtpdUnitTitleAr),
				utility.GetIntValueAsString(bcp.Rpo) + " " + utility.GetStringValue(bcp.RpoUnitTitleAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				bcp.ReferenceNumber,
				bcp.BusinessFunctionTitle,
				utility.GetIntValueAsString(bcp.Rto) + " " + utility.GetStringValue(bcp.RtoUnitTitleEn),
				utility.GetIntValueAsString(bcp.Mtpd) + " " + utility.GetStringValue(bcp.MtpdUnitTitleEn),
				utility.GetIntValueAsString(bcp.Rpo) + " " + utility.GetStringValue(bcp.RpoUnitTitleEn),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (crs *BusinessContinuityPlanService) GetDocumentPath(id uint, token string) (string, error) {
	existingBcp, err := crs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	if len(existingBcp.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the business continuity plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("business-continuity", "business-continuity-plan", strconv.FormatUint(uint64(id), 10), "documents")

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
