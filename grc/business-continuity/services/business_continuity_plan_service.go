package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BusinessContinuityPlanServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO, currentUserID uint) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessContinuityPlanDetailDTO, error)
	GetByID(id uint) (*models.BusinessContinuityPlan, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BusinessContinuityPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, businessContinuityPlan *businessContinuityDTOs.BusinessContinuityPlanFormDTO) error
	MarkAsCompleted(id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	UpdateRecoveryObjectives(businessContinuityPlan *models.BusinessContinuityPlan, tx utility.DatabaseTransaction) error
	ExportBusinessContinuityPlan(data []*businessContinuityDTOs.BusinessContinuityPlanListDTO, currentUserID uint) (*os.File, error)

	GetDocumentPath(id uint, token string) (string, error)
	GetBcpHeatmap(bcpId uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error)
	GetTestSummary(bcpID uint) (*businessContinuityDTOs.BcpTestSummaryDTO, error)
}

// Service and implementations
type BusinessContinuityPlanService struct {
	repo                businessContinuityPlanRepositories.BusinessContinuityPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	businessFunctionSvc BusinessFunctionServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	bcpStatusSvc        BcpStatusServiceInterface
	riskHeatmapSvc      riskServices.RiskHeatmapServiceInterface
}

func NewBusinessContinuityPlanService(
	repo businessContinuityPlanRepositories.BusinessContinuityPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	businessFunctionSvc BusinessFunctionServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	bcpStatusSvc BcpStatusServiceInterface,
	riskHeatmapSvc riskServices.RiskHeatmapServiceInterface,
) BusinessContinuityPlanServiceInterface {
	return &BusinessContinuityPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		businessFunctionSvc: businessFunctionSvc,
		userSettingSvc:      userSettingSvc,
		bcpStatusSvc:        bcpStatusSvc,
		riskHeatmapSvc:      riskHeatmapSvc,
	}
}

func (bcps *BusinessContinuityPlanService) GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO, currentUserID uint) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error) {
	if currentUserID != 0 {
		isCyberUser, err := bcps.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			return nil, err
		}
		params.IsCyberUser = isCyberUser
	}
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

	// Load all equipment tools
	equipmentTools, err := bcps.repo.GetAllEquipmentTools(tx...)
	if err == nil {
		businessContinuityPlan.EquipmentTools = equipmentTools
	}

	// Load all outage buckets
	outageBuckets, err := bcps.repo.GetAllOutageBuckets(tx...)
	if err == nil {
		businessContinuityPlan.OutageBuckets = outageBuckets
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

	isCyberUser := false
	if currentUserID != 0 {
		isCyberUser, err = bcps.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			return 0, err
		}
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
		IsCyber:            isCyberUser,
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
			ImpactAnalysisLevelID:    existingProcess.Process.ImpactAnalysisLevelId,
			RtoFromValue:             existingProcess.Process.RtoFromValue,
			RtoToValue:               existingProcess.Process.RtoToValue,
			RtoUnitId:                existingProcess.Process.RtoUnitId,
			RtoCondition:             existingProcess.Process.RtoCondition,
			MtpdFromValue:            existingProcess.Process.MtpdFromValue,
			MtpdToValue:              existingProcess.Process.MtpdToValue,
			MtpdUnitId:               existingProcess.Process.MtpdUnitId,
			MtpdCondition:            existingProcess.Process.MtpdCondition,
			RpoFromValue:             existingProcess.Process.RpoFromValue,
			RpoToValue:               existingProcess.Process.RpoToValue,
			RpoUnitId:                existingProcess.Process.RpoUnitId,
			RpoCondition:             existingProcess.Process.RpoCondition,
			RspFromValue:             existingProcess.Process.RspFromValue,
			RspToValue:               existingProcess.Process.RspToValue,
			RspUnitID:                existingProcess.Process.RspUnitId,
			RspCondition:             existingProcess.Process.RspCondition,
			CreatedBYID:              currentUserID,
		}
		bcpProcessID, err := bcps.repo.CreateBcpProcesses(bcpProcess, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if err := bcps.snapshotProcessChildren(existingProcess.ProcessID, bcpProcessID, currentUserID, tx); err != nil {
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

		// Create a new StatusUpdateDocument instance
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
		// Save the StatusUpdateDocument
		if err := bcps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
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
			"business_continuity",
			"bcp",
			"business_continuity_plan",
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
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "business continuity plan responsible unit")
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
			var deleteProcessIDs []uint
			var deleteProcessPtrIDs []*uint
			for _, process := range oldBusinessContinuityPlan.Processes {
				deleteProcessIDs = append(deleteProcessIDs, *process.ID)
				deleteProcessPtrIDs = append(deleteProcessPtrIDs, process.ID)
			}
			if err := bcps.repo.DeleteBcpProcessChildrenByBcpProcessIDs(deleteProcessIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
			if err := bcps.repo.DeleteBcpProcesses(deleteProcessPtrIDs, tx); err != nil {
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
				ImpactAnalysisLevelID:    existingProcess.Process.ImpactAnalysisLevelId,
				RtoFromValue:             existingProcess.Process.RtoFromValue,
				RtoToValue:               existingProcess.Process.RtoToValue,
				RtoUnitId:                existingProcess.Process.RtoUnitId,
				RtoCondition:             existingProcess.Process.RtoCondition,
				MtpdFromValue:            existingProcess.Process.MtpdFromValue,
				MtpdToValue:              existingProcess.Process.MtpdToValue,
				MtpdUnitId:               existingProcess.Process.MtpdUnitId,
				MtpdCondition:            existingProcess.Process.MtpdCondition,
				RpoFromValue:             existingProcess.Process.RpoFromValue,
				RpoToValue:               existingProcess.Process.RpoToValue,
				RpoUnitId:                existingProcess.Process.RpoUnitId,
				RpoCondition:             existingProcess.Process.RpoCondition,
				RspFromValue:             existingProcess.Process.RspFromValue,
				RspToValue:               existingProcess.Process.RspToValue,
				RspUnitID:                existingProcess.Process.RspUnitId,
				RspCondition:             existingProcess.Process.RspCondition,
				CreatedBYID:              currentUserID,
			}
			bcpProcessID, err := bcps.repo.CreateBcpProcesses(bcpProcess, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if err := bcps.snapshotProcessChildren(existingProcess.ProcessID, bcpProcessID, currentUserID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		// update bcp status of old business function

		oldBusinessFunction, err := bcps.businessFunctionSvc.GetByID(oldBusinessContinuityPlan.BusinessFunctionID)
		if err != nil {
			tx.Rollback()
			return err
		}

		oldBusinessFunction.BcpStatus = false

		if err := bcps.businessFunctionSvc.UpdateBusinessFunction(oldBusinessFunction, tx); err != nil {
			tx.Rollback()
			return err
		}

		// update bcp status of new business function
		newBusinessFunction, err := bcps.businessFunctionSvc.GetByID(businessContinuityPlanFormDTO.BusinessFunctionID)
		if err != nil {
			tx.Rollback()
			return err
		}

		newBusinessFunction.BcpStatus = true

		if err := bcps.businessFunctionSvc.UpdateBusinessFunction(newBusinessFunction, tx); err != nil {
			tx.Rollback()
			return err
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

			// genereate URL
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

		// log changes
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
				"business_continuity",
				"bcp",
				"business_continuity_plan",
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

// snapshotProcessChildren copies all child-table data from the Process (by processID)
// into the corresponding BcpProcess child tables, so BCP data is isolated from future Process changes.
func (bcps *BusinessContinuityPlanService) snapshotProcessChildren(processID uint, bcpProcessID uint, currentUserID uint, tx utility.DatabaseTransaction) error {
	process, err := bcps.repo.GetProcessWithChildrenByID(processID)
	if err != nil {
		return err
	}

	// Execution records
	for _, er := range process.ProcessExecutionRecords {
		if err := bcps.repo.CreateBcpProcessExecutionRecord(&models.BcpProcessExecutionRecord{
			BcpProcessID:         bcpProcessID,
			BCMExecutionRecordID: er.BCMExecutionRecordID,
			CreatedBYID:          &currentUserID,
		}, tx); err != nil {
			return err
		}
	}

	// Manual procedure forms
	for _, mf := range process.ProcessManualProcedureForms {
		if err := bcps.repo.CreateBcpProcessManualProcedureForm(&models.BcpProcessManualProcedureForm{
			BcpProcessID:             bcpProcessID,
			BCMManualProcedureFormID: mf.BCMManualProcedureFormID,
			CreatedBYID:              &currentUserID,
		}, tx); err != nil {
			return err
		}
	}

	// System register (single)
	if process.ProcessSystemRegister != nil {
		if err := bcps.repo.CreateBcpProcessSystemRegister(&models.BcpProcessSystemRegister{
			BcpProcessID:     bcpProcessID,
			SystemRegisterID: process.ProcessSystemRegister.SystemRegisterID,
			CreatedBYID:      &currentUserID,
		}, tx); err != nil {
			return err
		}
	}

	// Internal dependencies
	for _, dep := range process.InternalDependencies {
		if err := bcps.repo.CreateBcpProcessInternalDependency(&models.BcpProcessInternalDependency{
			BcpProcessID:             bcpProcessID,
			SupportingUnitID:         dep.SupportingUnitID,
			DependencyTierID:         dep.DependencyTierID,
			DependencyPercentageFrom: dep.DependencyPercentageFrom,
			DependencyPercentageTo:   dep.DependencyPercentageTo,
			RequiredSupport:          dep.RequiredSupport,
			Notes:                    dep.Notes,
			CreatedBYID:              &currentUserID,
		}, tx); err != nil {
			return err
		}
	}

	// External dependencies
	for _, dep := range process.ExternalDependencies {
		if err := bcps.repo.CreateBcpProcessExternalDependency(&models.BcpProcessExternalDependency{
			BcpProcessID:             bcpProcessID,
			EntityName:               dep.EntityName,
			EntityTypeID:             dep.EntityTypeID,
			NeededSupportDescription: dep.NeededSupportDescription,
			DependencyTierID:         dep.DependencyTierID,
			SLARegisterID:            dep.SLARegisterID,
			VendorRTO:                dep.VendorRTO,
			VendorRTOUnitID:          dep.VendorRTOUnitID,
			InternalStakeholder:      dep.InternalStakeholder,
			ContactNumber:            dep.ContactNumber,
			ContactEmail:             dep.ContactEmail,
			CreatedBYID:              &currentUserID,
		}, tx); err != nil {
			return err
		}
	}

	// Employee analysis (with nested employees and tools)
	// Build a map: original ProcessEmployeeAnalysis.ID -> new BcpProcessEmployeeAnalysis.ID
	eaIDMap := make(map[uint]uint)
	for _, ea := range process.ProcessEmployeeAnalysis {
		bcpEaID, err := bcps.repo.CreateBcpProcessEmployeeAnalysis(&models.BcpProcessEmployeeAnalysis{
			BcpProcessID:  bcpProcessID,
			DesignationID: ea.DesignationID,
			NormalCount:   ea.NormalCount,
			CriticalCount: ea.CriticalCount,
			CreatedBYID:   &currentUserID,
		}, tx)
		if err != nil {
			return err
		}
		eaIDMap[ea.ID] = bcpEaID
		for _, emp := range ea.Employees {
			bcpEmpID, err := bcps.repo.CreateBcpProcessEmployeeAnalysisEmployee(&models.BcpProcessEmployeeAnalysisEmployee{
				BcpProcessEmployeeAnalysisID: bcpEaID,
				UserID:                       emp.UserID,
				CreatedBYID:                  &currentUserID,
			}, tx)
			if err != nil {
				return err
			}
			for _, tool := range emp.Tools {
				if err := bcps.repo.CreateBcpProcessEmployeeAnalysisEmployeeTool(&models.BcpProcessEmployeeAnalysisEmployeeTool{
					BcpProcessEmployeeAnalysisEmployeeID: bcpEmpID,
					EquipmentToolID:                      tool.EquipmentToolID,
					IsHas:                                tool.IsHas,
					IsNeeded:                             tool.IsNeeded,
					CreatedBYID:                          &currentUserID,
				}, tx); err != nil {
					return err
				}
			}
		}
	}

	// BCM Staff allocations (with outage buckets and their users)
	for _, sa := range process.BCMStaffAllocations {
		var bcpEaIDPtr *uint
		if sa.ProcessEmployeeAnalysisID != nil {
			if mappedID, ok := eaIDMap[*sa.ProcessEmployeeAnalysisID]; ok {
				bcpEaIDPtr = &mappedID
			}
		}
		bcpSaID, err := bcps.repo.CreateBcpProcessStaffAllocation(&models.BcpProcessStaffAllocation{
			BcpProcessID:                 bcpProcessID,
			BcpProcessEmployeeAnalysisID: bcpEaIDPtr,
			CreatedBYID:                  &currentUserID,
		}, tx)
		if err != nil {
			return err
		}
		for _, ob := range sa.OutageBuckets {
			bucketID := ob.OutageBucketID
			bcpObID, err := bcps.repo.CreateBcpProcessAllocationOutageBucket(&models.BcpProcessAllocationOutageBucket{
				BcpProcessAllocationID: bcpSaID,
				OutageBucketID:         bucketID,
				IsSkipped:              ob.IsSkipped,
				CreatedBYID:            &currentUserID,
			}, tx)
			if err != nil {
				return err
			}
			for _, u := range ob.Users {
				if err := bcps.repo.CreateBcpProcessAllocationOutageBucketUser(&models.BcpProcessAllocationOutageBucketUser{
					BcpProcessAllocationOutageBucketID: bcpObID,
					UserID:                             &u.UserID,
					IsRemote:                           u.IsRemote,
					CreatedBYID:                        &currentUserID,
				}, tx); err != nil {
					return err
				}
			}
		}
	}

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
			var deleteProcessIDs []uint
			var deleteProcessPtrIDs []*uint
			for _, process := range businessContinuityPlan.Processes {
				deleteProcessIDs = append(deleteProcessIDs, *process.ID)
				deleteProcessPtrIDs = append(deleteProcessPtrIDs, process.ID)
			}
			if err := bcps.repo.DeleteBcpProcessChildrenByBcpProcessIDs(deleteProcessIDs, tx); err != nil {
				tx.Rollback()
				failedIds = append(failedIds, id)
				continue
			}
			if err := bcps.repo.DeleteBcpProcesses(deleteProcessPtrIDs, tx); err != nil {
				tx.Rollback()
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Delete each document if documents exist
		if len(businessContinuityPlan.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range businessContinuityPlan.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
			}
			if err := bcps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Delete the business
		if err := bcps.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		// update bcp status of business function
		businessFunction, err := bcps.businessFunctionSvc.GetByID(businessContinuityPlan.BusinessFunctionID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		businessFunction.BcpStatus = false

		if err := bcps.businessFunctionSvc.UpdateBusinessFunction(businessFunction, tx); err != nil {
			tx.Rollback()
			return nil, err
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
			"Business Function",
			"Version",
			"Responsible User",
			"Status",
			"Purpose",
			"Scope",
			"Objectives",
			"Assumptions",
			"Created By",
			"Created At",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Business Function",
			"Version",
			"Responsible User",
			"Status",
			"Purpose",
			"Scope",
			"Objectives",
			"Assumptions",
			"Created By",
			"Created At",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, bcp := range data {
		createdDate := bcp.CreatedAt
		if t, err := time.Parse(time.RFC3339, bcp.CreatedAt); err == nil {
			createdDate = t.Format("2006-01-02")
		} else if t, err := time.Parse("2006-01-02T15:04:05.999Z", bcp.CreatedAt); err == nil {
			createdDate = t.Format("2006-01-02")
		}
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				bcp.ReferenceNumber,
				bcp.BusinessFunctionTitle,
				utility.GetIntValueAsString(bcp.Version),
				fmt.Sprintf("%s %s", bcp.ResponsibleUserFirstNameAr, bcp.ResponsibleUserLastNameAr),
				bcp.BcpStatusTitleAr,
				utility.GetStringValue(bcp.Purpose),
				utility.GetStringValue(bcp.Scope),
				utility.GetStringValue(bcp.Objectives),
				utility.GetStringValue(bcp.Assumptions),
				fmt.Sprintf("%s %s", bcp.CreatedUserFirstNameAr, bcp.CreatedUserLastNameAr),
				createdDate,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				bcp.ReferenceNumber,
				bcp.BusinessFunctionTitle,
				utility.GetIntValueAsString(bcp.Version),
				fmt.Sprintf("%s %s", bcp.ResponsibleUserFirstNameEn, bcp.ResponsibleUserLastNameEn),
				bcp.BcpStatusTitleEn,
				utility.GetStringValue(bcp.Purpose),
				utility.GetStringValue(bcp.Scope),
				utility.GetStringValue(bcp.Objectives),
				utility.GetStringValue(bcp.Assumptions),
				fmt.Sprintf("%s %s", bcp.CreatedUserFirstNameEn, bcp.CreatedUserLastNameEn),
				createdDate,
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

func (bcpc *BusinessContinuityPlanService) GetBcpHeatmap(bcpId uint) ([][]*riskDTOs.RiskHeatmapResponseDTO, error) {
	bcpRisks, err := bcpc.repo.GetBcpRisksByBcpId(bcpId)
	if err != nil {
		return nil, err
	}

	riskRegisterIDs := make([]uint, len(bcpRisks))
	for i, risk := range bcpRisks {
		riskRegisterIDs[i] = risk.RiskRegisterID
	}

	heatmap, err := bcpc.riskHeatmapSvc.GetRiskHeatmap(riskRegisterIDs)
	if err != nil {
		return nil, err
	}

	return heatmap, nil
}

func (bcs *BusinessContinuityPlanService) GetTestSummary(bcpID uint) (*businessContinuityDTOs.BcpTestSummaryDTO, error) {
	// Get executed tests from test_and_exercises
	executedTests, err := bcs.repo.GetTestAndExercisesByBcpID(bcpID)
	if err != nil {
		return nil, err
	}

	var successful, partiallySuccess, failure, notTested int

	for _, testExercise := range executedTests {
		if testExercise == nil || testExercise.TestAndExerciseTestResult == nil {
			continue
		}

		for _, sr := range testExercise.TestAndExerciseTestResult.TestResultSystemRegisters {
			if sr.TestResultOption == nil {
				continue
			}
			switch sr.TestResultOption.Type {
			case "successful":
				successful++
			case "partially-successful":
				partiallySuccess++
			case "failure":
				failure++
			case "not-tested":
				notTested++
			}
		}
	}

	executedCount := successful + partiallySuccess + failure

	// Calculate test score (percentage of successful tests out of tested ones)
	testScore := 0
	if executedCount > 0 {
		testScore = (successful * 100) / executedCount
	}

	return &businessContinuityDTOs.BcpTestSummaryDTO{
		TestScore:        testScore,
		ExecutedTests:    executedCount,
		SuccessfulTests:  successful,
		FailedTests:      failure,
		PartiallySuccess: partiallySuccess,
		NotTestedTests:   notTested,
	}, nil
}

func (bcps *BusinessContinuityPlanService) MarkAsCompleted(id uint, currentUserID uint) error {
	tx := bcps.repo.BeginTransaction()

	businessContinuityPlan, err := bcps.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	publishedStatus, err := bcps.bcpStatusSvc.GetByType("published")
	if err != nil {
		tx.Rollback()
		return err
	}

	oldBcpStatusID := businessContinuityPlan.BcpStatusID
	businessContinuityPlan.BcpStatusID = publishedStatus.ID
	businessContinuityPlan.UpdatedBYID = &currentUserID
	if err := bcps.repo.Update(businessContinuityPlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"bcp_status_id": oldBcpStatusID},
			map[string]interface{}{"bcp_status_id": publishedStatus.ID},
		)
		if len(changedData) == 0 {
			return
		}

		if err := bcps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.BusinessContinuityPlan{}.TableName(),
			id,
			changedData,
			currentUserID,
			"business_continuity",
			"bcp",
			"business_continuity_plan",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.BusinessContinuityPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}
