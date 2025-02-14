package services

import (
	"errors"
	"fmt"
	"jethur/config"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementServiceInterface interface {
	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementListDTO, error)
	GetByID(requirementId uint) (*models.ComplianceAssessmentRequirement, error)
	GetFullDetailsByID(assessmentId uint, requirementId uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error)
	Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirement *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error

	MarkAsCompleted(complianceAssessmentId uint, requirementId uint, currentUserID uint) error
	MarkAsOpen(complianceAssessmentId uint, requirementId uint, currentUserID uint) error

	CreateDocument(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, currentUserID uint, documents []complianceDTOs.ComplianceAssessmentRequirementDocumentFormDTO) error
	DeleteDocuments(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(complianceAssessmentID, requirementID, id uint, token string) (string, error)

	ApproveComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	RevertComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	ResubmitComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
}

// Service and implementations
type ComplianceAssessmentRequirementService struct {
	repo                                     complianceRepositories.ComplianceAssessmentRequirementRepositoryInterface
	activityLogSvc                           logServices.ActivityLogServiceInterface
	userSvc                                  organizationServices.UserServiceInterface
	complianceRegisterSvc                    ComplianceRegisterServiceInterface
	complianceAssessmentSvc                  ComplianceAssessmentServiceInterface
	complianceAssessmentRequirementStatusSvc ComplianceAssessmentRequirementStatusServiceInterface
}

func NewComplianceAssessmentRequirementService(
	repo complianceRepositories.ComplianceAssessmentRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
	complianceAssessmentRequirementStatusSvc ComplianceAssessmentRequirementStatusServiceInterface,
) ComplianceAssessmentRequirementServiceInterface {
	return &ComplianceAssessmentRequirementService{
		repo:                                     repo,
		activityLogSvc:                           activityLogSvc,
		userSvc:                                  userSvc,
		complianceRegisterSvc:                    complianceRegisterSvc,
		complianceAssessmentSvc:                  complianceAssessmentSvc,
		complianceAssessmentRequirementStatusSvc: complianceAssessmentRequirementStatusSvc,
	}
}

func (cras *ComplianceAssessmentRequirementService) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementListDTO, error) {

	if _, err := cras.complianceAssessmentSvc.GetByID(assessmentID); err != nil {
		return nil, err
	}

	return cras.repo.GetAll(assessmentID, params)
}

func (cras ComplianceAssessmentRequirementService) GetByID(id uint) (*models.ComplianceAssessmentRequirement, error) {
	complianceAssessmentRequirement, err := cras.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentRequirement, nil
}

func (cras ComplianceAssessmentRequirementService) GetFullDetailsByID(assessmentId uint, requirementId uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error) {
	var complianceAssessmentRequirementDTO complianceDTOs.ComplianceAssessmentRequirementDetailDTO

	complianceAssessment, err := cras.repo.GetFullDetailsByID(requirementId, tx...)
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

	complianceAssessmentRequirementDTO.CreatedAt = complianceAssessment.CreatedAt
	complianceAssessmentRequirementDTO.UpdatedAt = complianceAssessment.UpdatedAt
	// todo :- also add extra fields which required in details screen

	return &complianceAssessmentRequirementDTO, nil
}

func (cras ComplianceAssessmentRequirementService) Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirementUpdateFormDTO *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error {
	tx := cras.repo.BeginTransaction()

	existingComplianceAssessmentRequirement, err := cras.GetByID(requirementId)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceAssessmentRequirement.Score = &complianceRequirementUpdateFormDTO.Score
	existingComplianceAssessmentRequirement.Comments = complianceRequirementUpdateFormDTO.Comments
	existingComplianceAssessmentRequirement.ComplianceStatusID = &complianceRequirementUpdateFormDTO.ComplianceStatusID
	existingComplianceAssessmentRequirement.UpdatedBYID = &currentUserID

	err = cras.repo.Update(existingComplianceAssessmentRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"score":             oldComplianceAssessmentRequirement.Score,
			"comments":          oldComplianceAssessmentRequirement.Comments,
			"compliance_status": utility.GetFieldValue(oldComplianceAssessmentRequirement.ComplianceStatus),
		}

		newData := map[string]interface{}{
			"score":             existingAssessmentRequirement.Score,
			"comments":          existingAssessmentRequirement.Comments,
			"compliance_status": utility.GetFieldValue(existingAssessmentRequirement.ComplianceStatus),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cras.activityLogSvc.LogActivity(
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
func (cras *ComplianceAssessmentRequirementService) CreateDocument(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, currentUserID uint, requiredDocuments []complianceDTOs.ComplianceAssessmentRequirementDocumentFormDTO) error {
	tx := cras.repo.BeginTransaction()
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
		if err := cras.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
		//move file to storage
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
func (cras *ComplianceAssessmentRequirementService) DeleteDocuments(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := cras.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedDocument []map[string]interface{}

	documentsToDelete := make(map[uint][]*complianceDTOs.ComplianceAssessmentRequirementDocumentDTO)
	complianceAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentID, requirementID)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {

		for _, requiredDocumentTitle := range complianceAssessmentRequirement.RequiredDocumentTitles {
			for _, doc := range requiredDocumentTitle.Documents {
				if doc != nil && *doc.ID == id {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
				}
			}
		}

		err = cras.repo.DeleteDocuments(id)
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

	if len(successfulIds) > 0 {

		for _, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(complianceAssessmentID), 10), "requirements", strconv.FormatUint(uint64(requirementID), 10), "required-documents", strconv.FormatUint(uint64(requiredDocumentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}

		go func() {
			err := cras.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceAssessmentRequirementRequiredDocument{}.TableName(),
				successfulIds,
				deletedDocument,
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

func (cras *ComplianceAssessmentRequirementService) GetDocumentPath(complianceAssessmentID, requirementID, id uint, token string) (string, error) {

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

func (cras ComplianceAssessmentRequirementService) MarkAsCompleted(complianceAssessmentId uint, requirementId uint, currentUserID uint) error {
	tx := cras.repo.BeginTransaction()

	existingComplianceAssessmentRequirement, err := cras.GetByID(requirementId)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceAssessmentRequirement.IsCompleted = true
	existingComplianceAssessmentRequirement.UpdatedBYID = &currentUserID

	err = cras.repo.Update(existingComplianceAssessmentRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"is_completed": oldComplianceAssessmentRequirement.IsCompleted,
		}

		newData := map[string]interface{}{
			"is_completed": existingAssessmentRequirement.IsCompleted,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cras.activityLogSvc.LogActivity(
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

func (cras ComplianceAssessmentRequirementService) MarkAsOpen(complianceAssessmentId uint, requirementId uint, currentUserID uint) error {
	tx := cras.repo.BeginTransaction()

	existingComplianceAssessmentRequirement, err := cras.GetByID(requirementId)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceAssessmentRequirement.IsCompleted = false
	existingComplianceAssessmentRequirement.UpdatedBYID = &currentUserID

	err = cras.repo.Update(existingComplianceAssessmentRequirement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessmentRequirement, err := cras.GetFullDetailsByID(complianceAssessmentId, requirementId)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"is_completed": oldComplianceAssessmentRequirement.IsCompleted,
		}

		newData := map[string]interface{}{
			"is_completed": existingAssessmentRequirement.IsCompleted,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cras.activityLogSvc.LogActivity(
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

	requirement, err := cars.GetByID(requirementID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate if the current status allows approval
	if requirement.ComplianceAssessmentRequirementStatusID == &approvedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("requirement is already approved"), "invalid_state", "requirement")
	}

	assessment, err := cars.complianceAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if currentUserID != assessment.ReviewUserID {
		return utility.NewForbiddenError("User is not authorized to approve this assessment")
	}

	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementID:       requirementID,
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
	requirement.ComplianceAssessmentRequirementStatusID = &approvedStatus.ID
	requirement.UpdatedBYID = &currentUserID

	err = cars.repo.Update(requirement, tx)
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
	}()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) RevertComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error {
	tx := cars.repo.BeginTransaction()

	revertedStatus, err := cars.complianceAssessmentRequirementStatusSvc.GetByType("reverted")
	if err != nil {
		tx.Rollback()
		return err
	}

	requirement, err := cars.repo.GetByID(requirementID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate if the current status allows reverting
	if requirement.ComplianceAssessmentRequirementStatusID == &revertedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("requirement is already reverted"), "invalid_state", "requirement")
	}

	assessment, err := cars.complianceAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if currentUserID != assessment.ReviewUserID {
		return utility.NewForbiddenError("User is not authorized to revert this assessment")
	}

	// Add a comment to `compliance_assessment_requirement_comments`
	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementID:       requirementID,
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
	requirement.ComplianceAssessmentRequirementStatusID = &revertedStatus.ID
	requirement.UpdatedBYID = &currentUserID

	err = cars.repo.Update(requirement, tx)
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

	requirement, err := cars.repo.GetByID(requirementID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate if the current status allows resubmission
	if requirement.ComplianceAssessmentRequirementStatusID == &resubmittedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("requirement is already resubmitted"), "invalid_state", "requirement")
	}

	userPrimaryUnitID, err := cars.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return fmt.Errorf("failed to get user primary unit: %w", err)
	}

	assessment, err := cars.complianceAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if userPrimaryUnitID != assessment.PrimaryUnitID {
		return utility.NewForbiddenError("User is not authorized to resubmit this assessment")
	}

	comment := &models.ComplianceAssessmentRequirementComment{
		ComplianceAssessmentRequirementID:       requirementID,
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
	requirement.ComplianceAssessmentRequirementStatusID = &resubmittedStatus.ID
	requirement.UpdatedBYID = &currentUserID

	err = cars.repo.Update(requirement, tx)
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
	}()

	return nil
}
