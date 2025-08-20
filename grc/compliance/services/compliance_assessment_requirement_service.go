package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementServiceInterface interface {
	GetAllClauses(assessmentID uint, currentUserID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, error)

	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementListResponseDTO, error)
	GetByID(requirementId uint) (*models.ComplianceAssessmentRequirement, error)
	GetFullDetailsByID(assessmentId uint, requirementId uint, responsibleUnitID uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error)
	Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirement *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error

	MarkAsCompleted(complianceAssessmentId uint, requirementId uint, currentUserID uint, responsibleUnitID uint) error
	MarkAsOpen(complianceAssessmentId uint, requirementId uint, currentUserID uint, responsibleUnitID uint) error

	CreateDocument(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, currentUserID uint, documents []complianceDTOs.ComplianceAssessmentRequirementDocumentFormDTO) error
	DeleteDocuments(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(complianceAssessmentID, requirementID, id uint, token string) (string, error)

	ApproveComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	RevertComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	ResubmitComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	GetComplianceAssessmentRequirementResponsibleUnitByID(id uint) (*models.ComplianceAssessmentRequirementResponsibleUnit, error)
}

// Service and implementations
type ComplianceAssessmentRequirementService struct {
	repo                                     complianceRepositories.ComplianceAssessmentRequirementRepositoryInterface
	activityLogSvc                           logServices.ActivityLogServiceInterface
	userSvc                                  organizationServices.UserServiceInterface
	complianceRegisterSvc                    ComplianceRegisterServiceInterface
	complianceAssessmentSvc                  ComplianceAssessmentServiceInterface
	complianceAssessmentRequirementStatusSvc ComplianceAssessmentRequirementStatusServiceInterface
	notificationSvc                          commonServices.NotificationServiceInterface
	complianceAssessmentFrameworkSvc         ComplianceAssessmentFrameworkServiceInterface
	complianceAssessmentStatusSvc            ComplianceAssessmentStatusServiceInterface
}

func NewComplianceAssessmentRequirementService(
	repo complianceRepositories.ComplianceAssessmentRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
	complianceAssessmentRequirementStatusSvc ComplianceAssessmentRequirementStatusServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface,
	complianceAssessmentStatusSvc ComplianceAssessmentStatusServiceInterface,
) ComplianceAssessmentRequirementServiceInterface {
	return &ComplianceAssessmentRequirementService{
		repo:                                     repo,
		activityLogSvc:                           activityLogSvc,
		userSvc:                                  userSvc,
		complianceRegisterSvc:                    complianceRegisterSvc,
		complianceAssessmentSvc:                  complianceAssessmentSvc,
		complianceAssessmentRequirementStatusSvc: complianceAssessmentRequirementStatusSvc,
		notificationSvc:                          notificationSvc,
		complianceAssessmentFrameworkSvc:         complianceAssessmentFrameworkSvc,
		complianceAssessmentStatusSvc:            complianceAssessmentStatusSvc,
	}
}

func (cars *ComplianceAssessmentRequirementService) GetAllClauses(assessmentID uint, currentUserID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, error) {
	clauses, err := cars.repo.GetAllClauses(assessmentID, currentUserID)
	if err != nil {
		return nil, err
	}

	// If no clauses are found, return an empty slice
	if len(clauses) == 0 {
		return []*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO{}, nil
	}

	return cars.buildHierarchicalClauseStructure(clauses)

}
func (cars *ComplianceAssessmentRequirementService) buildHierarchicalClauseStructure(clauses []*models.ComplianceClause) ([]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, error) {
	// Build a proper hierarchy by filtering out child clauses from the root level
	var rootClauses []*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO
	clauseMap := make(map[uint]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO)

	// Convert all clauses into DTOs and map them by ID
	for _, clause := range clauses {
		dto := &complianceDTOs.ComplianceAssessmentRequirementClausesListDTO{
			ID:           clause.ID,
			ClauseNo:     clause.ClauseNo,
			TitleEn:      clause.TitleEn,
			TitleAr:      clause.TitleAr,
			ChildClauses: []*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO{},
		}
		clauseMap[clause.ID] = dto
	}

	// Build hierarchy by linking children to their parent
	for _, clause := range clauses {
		if clause.ParentClauseID != nil {
			parent := clauseMap[*clause.ParentClauseID]
			parent.ChildClauses = append(parent.ChildClauses, clauseMap[clause.ID])
		} else {
			rootClauses = append(rootClauses, clauseMap[clause.ID])
		}
	}

	return rootClauses, nil
}
func (cars *ComplianceAssessmentRequirementService) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementListResponseDTO, error) {
	if _, err := cars.complianceAssessmentSvc.GetByID(assessmentID); err != nil {
		return nil, err
	}

	return cars.repo.GetAll(assessmentID, params)
}

func (cars ComplianceAssessmentRequirementService) GetByID(id uint) (*models.ComplianceAssessmentRequirement, error) {
	complianceAssessmentRequirement, err := cars.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentRequirement, nil
}

func (cars ComplianceAssessmentRequirementService) GetFullDetailsByID(assessmentId uint, requirementId uint, responsibleUnitID uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error) {
	var complianceAssessmentRequirementDTO complianceDTOs.ComplianceAssessmentRequirementDetailDTO

	complianceAssessment, err := cars.repo.GetFullDetailsByID(requirementId, responsibleUnitID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, requirementId)
		}
		return nil, err
	}

	err = mapstructure.Decode(complianceAssessment, &complianceAssessmentRequirementDTO)
	if err != nil {
		return nil, err
	}
	var filteredResponsibleUnits []*complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitDTO
	// If responsibleUnitID is provided, filter to return only that specific responsibleUnit
	for _, responsibleUnit := range complianceAssessmentRequirementDTO.ResponsibleUnits {
		if responsibleUnit.ResponsibleUnitID == responsibleUnitID {
			filteredResponsibleUnits = append(filteredResponsibleUnits, responsibleUnit)
			break // Only one responsibleUnit should match
		}
	}

	// ✅ iterate over each responsible unit
	for i, ru := range complianceAssessment.ResponsibleUnits {
		// Add TargetDate to each action plan
		for j, actionPlan := range ru.ActionPlans {
			if actionPlan.TargetDate != nil {
				targetDate := actionPlan.TargetDate.Format("2006-01-02")

				// Defensive check: ensure DTO has the same structure
				if len(complianceAssessmentRequirementDTO.ResponsibleUnits) > i &&
					len(complianceAssessmentRequirementDTO.ResponsibleUnits[i].ActionPlans) > j {
					complianceAssessmentRequirementDTO.ResponsibleUnits[i].ActionPlans[j].TargetDateString = &targetDate
				}
			}
		}

		// Add CreatedAt to each comment
		for j, comment := range ru.RequirementComments {
			if len(complianceAssessmentRequirementDTO.ResponsibleUnits) > i &&
				len(complianceAssessmentRequirementDTO.ResponsibleUnits[i].RequirementComments) > j {
				complianceAssessmentRequirementDTO.ResponsibleUnits[i].RequirementComments[j].CreatedAt = comment.CreatedAt
			}
		}
	}

	// Check if any responsibleUnit was found before accessing the first element
	if len(filteredResponsibleUnits) > 0 {
		complianceAssessmentRequirementDTO.ResponsibleUnit = filteredResponsibleUnits[0]
		complianceAssessmentRequirementDTO.CompletedByUser = filteredResponsibleUnits[0].CompletedByUser
	}

	complianceAssessmentRequirementDTO.CreatedAt = complianceAssessment.CreatedAt
	complianceAssessmentRequirementDTO.UpdatedAt = complianceAssessment.UpdatedAt
	// todo :- also add extra fields which required in details screen

	return &complianceAssessmentRequirementDTO, nil
}

func (cars ComplianceAssessmentRequirementService) Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirementUpdateFormDTO *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error {
	assessmentInProgressStatus, err := cars.complianceAssessmentStatusSvc.GetByType("in_progress")
	if err != nil {
		return err
	}

	assessmentInitiatedStatus, err := cars.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}
	assessmentDetails, err := cars.complianceAssessmentSvc.GetFullDetailsByID(complianceAssessmentId)
	if err != nil {
		return err
	}

	// existingComplianceAssessmentRequirement, err := cars.GetByID(requirementId)
	// if err != nil {
	// 	return err
	// }

	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, complianceRequirementUpdateFormDTO.ResponsibleUnitID)
	if err != nil {
		return err
	}
	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	complianceAssessmentFrameworkValue, err := cars.complianceAssessmentFrameworkSvc.GetByFrameworkValueID(complianceRequirementUpdateFormDTO.AssessmentFrameworkValueID)
	if err != nil {
		return err
	}

	complianceAssessment, err := cars.complianceAssessmentSvc.GetByID(complianceAssessmentId)
	if err != nil {
		return err
	}
	tx := cars.repo.BeginTransaction()
	// existingComplianceAssessmentRequirement.ComplianceStatusID = &complianceRequirementUpdateFormDTO.ComplianceStatusID
	// existingComplianceAssessmentRequirement.UpdatedBYID = &currentUserID

	// err = cars.repo.Update(existingComplianceAssessmentRequirement, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	if complianceAssessment.ComplianceAssessmentStatusID == assessmentInitiatedStatus.ID {
		complianceAssessment.ComplianceAssessmentStatusID = assessmentInProgressStatus.ID
		err = cars.repo.ComplianceAssessmentUpdateStatus(complianceAssessment, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if complianceAssessmentFrameworkValue.IsApplicable && assessmentDetails.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType.Type == "percentage" {
		complianceAssessmentRequirementResponsibleUnit.Score = &complianceRequirementUpdateFormDTO.Score
	} else {
		complianceAssessmentRequirementResponsibleUnit.Score = nil
	}
	complianceAssessmentRequirementResponsibleUnit.Comments = complianceRequirementUpdateFormDTO.Comments
	complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID = &complianceRequirementUpdateFormDTO.AssessmentFrameworkValueID
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID

	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"score":                                 complianceAssessmentRequirementResponsibleUnit.Score,
			"comments":                              complianceAssessmentRequirementResponsibleUnit.Comments,
			"compliance_assessment_framework_value": utility.GetFieldValue(complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValue),
		}

		newData := map[string]interface{}{
			"score":                                 complianceRequirementUpdateFormDTO.Score,
			"comments":                              complianceRequirementUpdateFormDTO.Comments,
			"compliance_assessment_framework_value": utility.GetFieldValue(complianceAssessmentFrameworkValue),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cars.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessmentRequirement{}.TableName(),
				requirementId,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
					"recordID":  requirementId,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) CreateDocument(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, currentUserID uint, requiredDocuments []complianceDTOs.ComplianceAssessmentRequirementDocumentFormDTO) error {
	tx := cars.repo.BeginTransaction()
	var documents []*models.ComplianceAssessmentRequirementRequiredDocument
	for _, doc := range requiredDocuments {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		url := utility.GenerateDocumentURL("compliance", "compliance-assessment", strconv.FormatUint(uint64(complianceAssessmentID), 10), "requirements", strconv.FormatUint(uint64(requirementID), 10), "required-documents", strconv.FormatUint(uint64(requiredDocumentID), 10), "documents")

		document := &models.ComplianceAssessmentRequirementRequiredDocument{
			ComplianceAssessmentRequirementDocumentTitleID: requiredDocumentID,
			Language:     doc.Language,
			URL:          url,
			Name:         doc.FileName,
			Size:         doc.Size,
			Token:        encryptedFilename,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			CreatedByID:  currentUserID,
		}

		documents = append(documents, document)

	}

	if len(documents) > 0 {
		if err := cars.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(complianceAssessmentID), 10), "requirements", strconv.FormatUint(uint64(requirementID), 10), "required-documents", strconv.FormatUint(uint64(requiredDocumentID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}

		}
	}
	tx.Commit()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) DeleteDocuments(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, ids []uint, currentUserID uint) ([]uint, error) {

	tx := cars.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedDocument []map[string]interface{}

	documentsToDelete := make(map[uint][]*complianceDTOs.ComplianceAssessmentRequirementDocumentDTO)

	// Fetch full detail
	complianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentID, requirementID, 0)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		// Loop through each responsible unit
		for _, ru := range complianceAssessmentRequirement.ResponsibleUnits {
			for _, requiredDocumentTitle := range ru.RequiredDocumentTitles {
				for _, doc := range requiredDocumentTitle.Documents {
					if doc != nil && *doc.ID == id {
						documentsToDelete[id] = append(documentsToDelete[id], doc)
					}
				}
			}
		}

		// Delete from database
		err = cars.repo.DeleteDocuments(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedDocument = append(deletedDocument, map[string]interface{}{
			"id":                                id,
			"compliance_assessment_requirement": complianceAssessmentRequirement.ComplianceRequirement.TitleEn,
		})
	}

	// Remove files from storage
	if len(successfulIds) > 0 {
		for _, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join(
					"compliance", "compliance-assessment",
					strconv.FormatUint(uint64(complianceAssessmentID), 10),
					"requirements",
					strconv.FormatUint(uint64(requirementID), 10),
					"required-documents",
					strconv.FormatUint(uint64(requiredDocumentID), 10),
					"documents",
				)
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}

		// Log deletion activity
		go func() {
			err := cars.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceAssessmentRequirementRequiredDocument{}.TableName(),
				successfulIds,
				deletedDocument,
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

func (cars *ComplianceAssessmentRequirementService) GetDocumentPath(complianceAssessmentID, requirementID, id uint, token string) (string, error) {
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(complianceAssessmentID), 10), "requirements", strconv.FormatUint(uint64(requirementID), 10), "required-documents", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (cars ComplianceAssessmentRequirementService) MarkAsCompleted(complianceAssessmentId uint, requirementId uint, currentUserID uint, responsibleUnitID uint) error {
	tx := cars.repo.BeginTransaction()

	_, err := cars.GetByID(requirementId)
	if err != nil {
		return err
	}
	allComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, 0)
	if err != nil {
		return err
	}
	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, responsibleUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	now := time.Now()
	complianceAssessmentRequirementResponsibleUnit.IsCompleted = true
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID
	complianceAssessmentRequirementResponsibleUnit.CompletedByUserID = &currentUserID
	complianceAssessmentRequirementResponsibleUnit.CompletedAt = &now

	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	allCompleted := true
	for _, ru := range allComplianceAssessmentRequirement.ResponsibleUnits {
		if ru == nil {
			continue
		}
		if ru.ID == complianceAssessmentRequirementResponsibleUnit.ID {
			// this one is now completed
			continue
		}
		if !ru.IsCompleted {
			allCompleted = false
			break
		}
	}

	// Update parent CAR’s IsAllCompleted
	car, err := cars.repo.GetByID(requirementId)
	if err != nil {
		tx.Rollback()
		return err
	}
	car.IsAllCompleted = allCompleted
	car.UpdatedBYID = &currentUserID

	if err := cars.repo.Update(car, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, responsibleUnitID)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"is_completed": oldComplianceAssessmentRequirement.ResponsibleUnit.IsCompleted,
		}

		newData := map[string]interface{}{
			"is_completed": existingAssessmentRequirement.ResponsibleUnit.IsCompleted,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cars.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessmentRequirement{}.TableName(),
				requirementId,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
					"recordID":  requirementId,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cars ComplianceAssessmentRequirementService) MarkAsOpen(complianceAssessmentId uint, requirementId uint, currentUserID uint, responsibleUnitID uint) error {
	tx := cars.repo.BeginTransaction()

	_, err := cars.GetByID(requirementId)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, responsibleUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	complianceAssessmentRequirementResponsibleUnit.IsCompleted = false
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID

	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	car, err := cars.repo.GetByID(requirementId)
	if err != nil {
		tx.Rollback()
		return err
	}
	car.IsAllCompleted = false
	car.UpdatedBYID = &currentUserID

	if err := cars.repo.Update(car, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		existingAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, responsibleUnitID)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"is_completed": oldComplianceAssessmentRequirement.ResponsibleUnit.IsCompleted,
		}

		newData := map[string]interface{}{
			"is_completed": existingAssessmentRequirement.ResponsibleUnit.IsCompleted,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cars.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessmentRequirement{}.TableName(),
				requirementId,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
					"recordID":  requirementId,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) ApproveComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error {
	tx := cars.repo.BeginTransaction()

	approvedStatus, err := cars.complianceAssessmentRequirementStatusSvc.GetByType("approved")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = cars.GetByID(requirementID)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessment, err := cars.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	allComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(assessmentID, requirementID, 0)
	if err != nil {
		return err
	}
	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(assessmentID, requirementID, commentForm.ResponsibleUnitID)
	if err != nil {
		return err
	}
	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID == &approvedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("requirement is already approved"), "invalid_state", "requirement")
	}
	if currentUserID != complianceAssessment.ComplianceRegister.ResponsibleUserID {
		return utility.NewForbiddenError("User is not authorized to approve this assessment")
	}

	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementResponsibleUnitID: complianceAssessmentRequirementResponsibleUnit.ID,
		Comment:                                 commentForm.Comment,
		CreatedBYID:                             currentUserID,
		ComplianceAssessmentRequirementStatusID: &approvedStatus.ID,
	}

	err = cars.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update requirement status to "approved"
	complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID = &approvedStatus.ID
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID

	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// assume you've just set the target RU to "approved"
	allApproved := true
	for _, ru := range allComplianceAssessmentRequirement.ResponsibleUnits {
		if ru == nil {
			continue
		}
		if ru.ID == complianceAssessmentRequirementResponsibleUnit.ID {
			// we just approved this one; treat as approved
			continue
		}
		st := ru.ComplianceAssessmentRequirementStatus
		if st == nil || st.Type != "approved" {
			allApproved = false
			break
		}
	}

	car, err := cars.repo.GetByID(requirementID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if allApproved {
		car.ComplianceAssessmentRequirementStatusID = &approvedStatus.ID
	} else {
		car.ComplianceAssessmentRequirementStatusID = nil
	}
	car.UpdatedBYID = &currentUserID

	if err := cars.repo.Update(car, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := cars.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentRequirement{}.TableName(),
			requirementID,
			map[string]interface{}{
				"status":  "approved",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
				"recordID":  requirementID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		assessmentDetails, err := cars.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
		if err != nil {
			tx.Rollback()
			return
		}

		createdBy, err := cars.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		if complianceAssessmentRequirementResponsibleUnit.CompletedByUserID != nil {
			completedUser, err := cars.userSvc.GetByID(*complianceAssessmentRequirementResponsibleUnit.CompletedByUserID)
			if err != nil {
				utility.LogError("Failed to fetch submitted user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Prepare notification and email
			urlParam := fmt.Sprintf(`{"assessment_id":%d, "requirement_id":%d}`, assessmentID, requirementID)
			entity := "compliance-assessment-requirement"
			action := "approved"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{completedUser.ID},
				ModuleIdentifier:     "compliance_register",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Compliance Requirement Approved",
				HeadingAr:            "تم الموافقة على متطلبات الامتثال",
				TitleEn:              assessmentDetails.ComplianceRegister.TitleEn,
				TitleAr:              assessmentDetails.ComplianceRegister.TitleAr,
				ReferenceNumber:      assessmentDetails.ComplianceRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{completedUser.Email},
				Subject:              fmt.Sprintf("Compliance Requirement #%d Approved", requirementID),
				TemplateRelativePath: "compliance/templates/email/compliance_requirement_approved.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":         completedUser.FirstNameEn,
					"LastNameEn":          completedUser.LastNameEn,
					"TitleEn":             assessmentDetails.ComplianceRegister.TitleEn,
					"TitleAr":             assessmentDetails.ComplianceRegister.TitleAr,
					"ReferenceNumber":     assessmentDetails.ComplianceRegister.ReferenceNumber,
					"Comment":             commentForm.Comment,
					"ApprovedByFirstName": createdBy.FirstNameEn,
					"ApprovedByLastName":  createdBy.LastNameEn,
					"ApprovalDate":        time.Now().Format("2006-01-02"),
					"RedirectURL":         fmt.Sprintf("%s/compliance/assessment/%d/overview", utility.GetAppURL(), assessmentID),
				},
			}

			// **Send Both Notifications at ONCE**
			err = cars.notificationSvc.SendNotification(dbNotificationDTO, emailNotification, "compliance_assessment")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) RevertComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error {

	revertedStatus, err := cars.complianceAssessmentRequirementStatusSvc.GetByType("reverted")
	if err != nil {
		return err
	}
	complianceAssessment, err := cars.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}
	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(assessmentID, requirementID, commentForm.ResponsibleUnitID)
	if err != nil {
		return err
	}
	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	// Validate if the current status allows reverting
	if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID == &revertedStatus.ID {
		return utility.NewValidationError(fmt.Errorf("requirement is already reverted"), "invalid_state", "requirement")
	}

	if currentUserID != complianceAssessment.ComplianceRegister.ResponsibleUserID {
		return utility.NewForbiddenError("User is not authorized to revert this assessment")
	}
	tx := cars.repo.BeginTransaction()

	// Add a comment to `compliance_assessment_requirement_comments`
	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementResponsibleUnitID: complianceAssessmentRequirementResponsibleUnit.ID,
		Comment:                                 commentForm.Comment,
		CreatedBYID:                             currentUserID,
		CreatedAt:                               time.Now(),
		ComplianceAssessmentRequirementStatusID: &revertedStatus.ID,
	}

	err = cars.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update requirement status to "reverted"
	complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID = &revertedStatus.ID
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID

	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// car, err := cars.repo.GetByID(requirementID)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// car.ComplianceAssessmentRequirementStatusID = nil
	// car.UpdatedBYID = &currentUserID

	// if err := cars.repo.Update(car, tx); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	tx.Commit()

	// Log the action
	go func() {
		err := cars.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentRequirement{}.TableName(),
			requirementID,
			map[string]interface{}{
				"status":  "reverted",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
				"recordID":  requirementID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		assessmentDetails, err := cars.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
		if err != nil {
			tx.Rollback()
			return
		}

		createdBy, err := cars.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		if complianceAssessmentRequirementResponsibleUnit.CompletedByUserID != nil {

			completedUser, err := cars.userSvc.GetByID(*complianceAssessmentRequirementResponsibleUnit.CompletedByUserID)
			if err != nil {
				utility.LogError("Failed to fetch submitted user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Prepare notification and email
			urlParam := fmt.Sprintf(`{"assessment_id":%d, "requirement_id":%d}`, assessmentID, requirementID)
			entity := "compliance-assessment-requirement"
			action := "reverted"

			dbNotification := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{completedUser.ID},
				ModuleIdentifier:     "compliance_register",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Requirement Reverted",
				HeadingAr:            "تم إرجاع متطلب الامتثال",
				TitleEn:              assessmentDetails.ComplianceRegister.TitleEn,
				TitleAr:              assessmentDetails.ComplianceRegister.TitleAr,
				ReferenceNumber:      assessmentDetails.ComplianceRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{completedUser.Email},
				Subject:              fmt.Sprintf("Compliance Requirement #%d Reverted", requirementID),
				TemplateRelativePath: "compliance/templates/email/compliance_requirement_reverted.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":         completedUser.FirstNameEn,
					"LastNameEn":          completedUser.LastNameEn,
					"TitleEn":             assessmentDetails.ComplianceRegister.TitleEn,
					"TitleAr":             assessmentDetails.ComplianceRegister.TitleAr,
					"ReferenceNumber":     assessmentDetails.ComplianceRegister.ReferenceNumber,
					"Comment":             commentForm.Comment,
					"RevertedByFirstName": createdBy.FirstNameEn,
					"RevertedByLastName":  createdBy.LastNameEn,
					"ReversionDate":       time.Now().Format("2006-01-02"),
					"RedirectURL":         fmt.Sprintf("%s/compliance/assessment/%d/overview", utility.GetAppURL(), assessmentID),
				},
			}

			// **Send Both Notifications at ONCE**
			err = cars.notificationSvc.SendNotification(dbNotification, emailNotification, "compliance_assessment")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) ResubmitComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error {
	tx := cars.repo.BeginTransaction()

	resubmittedStatus, err := cars.complianceAssessmentRequirementStatusSvc.GetByType("resubmitted")
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessmentRequirement, err := cars.GetFullDetailsByID(assessmentID, requirementID, commentForm.ResponsibleUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(complianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	// Validate if the current status allows resubmission
	if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID == &resubmittedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("requirement is already resubmitted"), "invalid_state", "requirement")
	}

	// userPrimaryUnitID, err := cars.userSvc.GetUserPrimaryUnitID(currentUserID)
	// if err != nil {
	// 	return fmt.Errorf("failed to get user primary unit: %w", err)
	// }

	// if userPrimaryUnitID != assessment.PrimaryUnitID {
	// 	return utility.NewForbiddenError("User is not authorized to resubmit this assessment")
	// }

	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementResponsibleUnitID: complianceAssessmentRequirementResponsibleUnit.ID,
		Comment:                                 commentForm.Comment,
		CreatedBYID:                             currentUserID,
		ComplianceAssessmentRequirementStatusID: &resubmittedStatus.ID,
	}

	err = cars.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update requirement status to "resubmitted"
	complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID = &resubmittedStatus.ID
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID
	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := cars.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentRequirement{}.TableName(),
			requirementID,
			map[string]interface{}{
				"status":  "resubmitted",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ComplianceAssessmentRequirement{}.TableName(),
				"recordID":  requirementID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		assessmentDetails, err := cars.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
		if err != nil {
			tx.Rollback()
			return
		}

		createdBy, err := cars.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		responsibleUser, err := cars.userSvc.GetByID(assessmentDetails.ComplianceRegister.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"assessment_id":%d, "requirement_id":%d}`, assessmentID, requirementID)
		entity := "compliance-assessment-requirement"
		action := "resubmitted"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{responsibleUser.ID},
			ModuleIdentifier:     "compliance_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Requirement Resubmitted",
			HeadingAr:            "إعادة تقديم متطلب الامتثال",
			TitleEn:              assessmentDetails.ComplianceRegister.TitleEn,
			TitleAr:              assessmentDetails.ComplianceRegister.TitleAr,
			ReferenceNumber:      assessmentDetails.ComplianceRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Compliance Requirement #%d Resubmitted", requirementID),
			TemplateRelativePath: "compliance/templates/email/compliance_requirement_resubmitted.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"TitleEn":                assessmentDetails.ComplianceRegister.TitleEn,
				"TitleAr":                assessmentDetails.ComplianceRegister.TitleAr,
				"ReferenceNumber":        assessmentDetails.ComplianceRegister.ReferenceNumber,
				"Comment":                commentForm.Comment,
				"ResubmittedByFirstName": createdBy.FirstNameEn,
				"ResubmittedByLastName":  createdBy.LastNameEn,
				"ResubmissionDate":       time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/compliance/assessment/%d/overview", utility.GetAppURL(), assessmentID),
			},
		}

		// **Send Both Notifications at ONCE**
		err = cars.notificationSvc.SendNotification(dbNotification, emailNotification, "compliance_assessment")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (cars ComplianceAssessmentRequirementService) GetComplianceAssessmentRequirementResponsibleUnitByID(id uint) (*models.ComplianceAssessmentRequirementResponsibleUnit, error) {
	complianceAssessmentRequirementResponsibleUnit, err := cars.repo.GetComplianceAssessmentRequirementResponsibleUnitByID(id)
	if err != nil {
		return nil, err
	}
	return complianceAssessmentRequirementResponsibleUnit, nil
}
