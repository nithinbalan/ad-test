package services

import (
	"errors"
	"fmt"
	"math"
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
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	settingsDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementServiceInterface interface {
	GetAllClauses(assessmentID uint, currentUserID uint, complianceID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, error)

	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementListResponseDTO, error)
	GetByID(requirementId uint) (*models.ComplianceAssessmentRequirement, error)
	GetFullDetailsByID(assessmentId uint, requirementId uint, responsibleUnitID uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementDetailDTO, error)
	Update(complianceAssessmentId uint, requirementId uint, currentUserID uint, complianceRequirement *complianceDTOs.ComplianceAssessmentRequirementUpdateFormDTO) error

	MarkAsCompleted(complianceAssessmentId uint, requirementId uint, currentUserID uint, responsibleUnitID uint) error
	MarkAsOpen(complianceAssessmentId uint, requirementId uint, currentUserID uint, responsibleUnitID uint) error

	CreateDocument(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, currentUserID uint, documents []complianceDTOs.ComplianceAssessmentRequirementDocumentFormDTO) error
	DeleteDocuments(complianceAssessmentID uint, requirementID uint, requiredDocumentID uint, ids []uint, currentUserID uint) ([]uint, error)
	UpdateDocument(complianceAssessmentID uint, requirementID uint, documentID uint, currentUserID uint, formDTO *complianceDTOs.ComplianceAssessmentRequirementDocumentUpdateFormDTO) error

	GetDocumentPath(complianceAssessmentID, requirementID, id uint, token string) (string, error)

	ApproveComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	RevertComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	ResubmitComplianceRequirement(assessmentID uint, requirementID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentRequirementCommentFormDTO) error
	GetComplianceAssessmentRequirementResponsibleUnitByID(id uint) (*models.ComplianceAssessmentRequirementResponsibleUnit, error)

	MarkAsNotApplicable(assessmentID uint, requirementID uint, currentUserID uint, form *complianceDTOs.ComplianceAssessmentRequirementNotApplicableFormDTO) error
	AcknowledgeRequirement(assessmentID uint, requirementID uint, currentUserID uint, form *complianceDTOs.AcknowledgeRequirementFormDTO) error
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
	complianceStatusSvc                      ComplianceStatusServiceInterface
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
	complianceStatusSvc ComplianceStatusServiceInterface,
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
		complianceStatusSvc:                      complianceStatusSvc,
	}
}

func (cars *ComplianceAssessmentRequirementService) GetAllClauses(assessmentID uint, currentUserID uint, complianceID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, error) {
	clauses, err := cars.repo.GetAllClauses(assessmentID, currentUserID, complianceID)
	if err != nil {
		return nil, err
	}
	// If no clauses are found, return an empty slice
	if len(clauses) == 0 {
		return []*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO{}, nil
	}
	countsMap, err := cars.repo.GetClauseRequirementCounts(assessmentID)
	if err != nil {
		return nil, err
	}
	return cars.buildHierarchicalClauseStructure(clauses, countsMap)
}

func (cars *ComplianceAssessmentRequirementService) buildHierarchicalClauseStructure(clauses []*models.ComplianceClause, countMap map[uint]complianceDTOs.ClauseRequirementCountDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, error) {

	// 1) Create DTOs and index by ID
	clauseMap := make(map[uint]*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO, len(clauses))
	for _, clause := range clauses {
		dto := &complianceDTOs.ComplianceAssessmentRequirementClausesListDTO{
			ID:                    clause.ID,
			ClauseNo:              clause.ClauseNo,
			TitleEn:               clause.TitleEn,
			TitleAr:               clause.TitleAr,
			ChildClauses:          []*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO{},
			TotalRequirements:     0,
			CompletedRequirements: 0,
		}
		// seed leaf/self counts from DB (direct counts)
		if c, ok := countMap[clause.ID]; ok {
			dto.TotalRequirements = c.TotalRequirements
			dto.CompletedRequirements = c.CompletedRequirements
		}
		clauseMap[clause.ID] = dto
	}

	// 2) Wire up the hierarchy
	var roots []*complianceDTOs.ComplianceAssessmentRequirementClausesListDTO
	for _, clause := range clauses {
		if clause.ParentClauseID != nil {
			parent := clauseMap[*clause.ParentClauseID]
			parent.ChildClauses = append(parent.ChildClauses, clauseMap[clause.ID])
		} else {
			roots = append(roots, clauseMap[clause.ID])
		}
	}

	// 3) Roll-up: post-order sum of children into parent
	var rollup func(n *complianceDTOs.ComplianceAssessmentRequirementClausesListDTO) (tot, done int)
	rollup = func(n *complianceDTOs.ComplianceAssessmentRequirementClausesListDTO) (int, int) {
		// start from self (direct counts)
		total := n.TotalRequirements
		completed := n.CompletedRequirements

		// add children recursively
		for _, ch := range n.ChildClauses {
			ct, cd := rollup(ch)
			total += ct
			completed += cd
		}

		// write rolled-up numbers back to the node
		n.TotalRequirements = total
		n.CompletedRequirements = completed
		return total, completed
	}

	for _, r := range roots {
		rollup(r)
	}

	return roots, nil
}

func (cars *ComplianceAssessmentRequirementService) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementListResponseDTO, error) {
	if _, err := cars.complianceAssessmentSvc.GetByID(assessmentID); err != nil {
		return nil, err
	}

	response, err := cars.repo.GetAll(assessmentID, params)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return response, nil
	}

	assessmentRequirementIDs := make([]uint, 0, len(response.Items))
	for _, item := range response.Items {
		assessmentRequirementIDs = append(assessmentRequirementIDs, item.ID)
	}

	var unitRows []complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitListRowDTO
	if err := cars.repo.GetResponsibleUnitRowsByAssessmentRequirementIDs(assessmentRequirementIDs, params, &unitRows); err != nil {
		return nil, err
	}

	unitMap := make(map[uint][]*complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitListDTO)
	for _, row := range unitRows {
		unit := &settingsDTOs.OrganizationUnitDTO{
			ID:        int(row.ResponsibleUnitID),
			TitleEn:   row.TitleEn,
			TitleAr:   row.TitleAr,
			Status:    row.Status,
			ColorCode: row.ColorCode,
		}
		if row.ParentID != nil {
			unit.ParentUnit = &settingsDTOs.OrganizationUnitDTO{ID: int(*row.ParentID)}
		}

		unitDTO := &complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitListDTO{
			ID:                                 row.ID,
			AssessmentRequirementID:            row.AssessmentRequirementID,
			ResponsibleUnitID:                  row.ResponsibleUnitID,
			IsRegulationNotificationAssessment: row.IsRegulationNotificationAssessment,
			IsRiskBasedAssessment:              row.IsRiskBasedAssessment,
			ResponsibleUnit:                    unit,
			ResponsibleUserID:                  row.ResponsibleUserID,
			IsAcknowledged:                     row.IsAcknowledged,
			AcknowledgeUserID:                  row.AcknowledgeUserID,
			AcknowledgeComment:                 row.AcknowledgeComment,
		}

		if row.ResponsibleUserID != nil && row.ResponsibleUserFirstNameEn != nil {
			unitDTO.ResponsibleUser = &organizationDTOs.UserMinimalDTO{
				ID:          int(*row.ResponsibleUserID),
				FirstNameEN: strVal(row.ResponsibleUserFirstNameEn),
				LastNameEN:  strVal(row.ResponsibleUserLastNameEn),
				FirstNameAR: strVal(row.ResponsibleUserFirstNameAr),
				LastNameAR:  strVal(row.ResponsibleUserLastNameAr),
				Email:       strVal(row.ResponsibleUserEmail),
				Initial:     strVal(row.ResponsibleUserInitial),
				Color:       strVal(row.ResponsibleUserColor),
			}
		}

		if row.AcknowledgeUserID != nil && row.AcknowledgeUserFirstNameEn != nil {
			unitDTO.AcknowledgeUser = &organizationDTOs.UserMinimalDTO{
				ID:          int(*row.AcknowledgeUserID),
				FirstNameEN: strVal(row.AcknowledgeUserFirstNameEn),
				LastNameEN:  strVal(row.AcknowledgeUserLastNameEn),
				FirstNameAR: strVal(row.AcknowledgeUserFirstNameAr),
				LastNameAR:  strVal(row.AcknowledgeUserLastNameAr),
				Email:       strVal(row.AcknowledgeUserEmail),
				Initial:     strVal(row.AcknowledgeUserInitial),
				Color:       strVal(row.AcknowledgeUserColor),
			}
		}
		unitMap[row.AssessmentRequirementID] = append(unitMap[row.AssessmentRequirementID], unitDTO)
	}

	for _, item := range response.Items {
		if units, ok := unitMap[item.ID]; ok {
			item.ResponsibleUnits = units
		} else {
			item.ResponsibleUnits = []*complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitListDTO{}
		}
	}

	return response, nil
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

	// Get framework type from compliance register
	frameworkType := ""
	if len(assessmentDetails.Compliances) > 0 && assessmentDetails.Compliances[0].ComplianceRegister != nil &&
		assessmentDetails.Compliances[0].ComplianceRegister.ComplianceAssessmentFramework != nil &&
		assessmentDetails.Compliances[0].ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType != nil {
		frameworkType = assessmentDetails.Compliances[0].ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType.Type
	}

	if complianceAssessmentFrameworkValue.IsApplicable && frameworkType == "percentage" {
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
				"compliance", "compliance_assessment", "requirement", nil,
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

		// Delete from database (no transaction needed as repo method doesn't use it)
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

	// Remove files from storage (file cleanup should not affect DB operation)
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
					// Log error but don't fail the operation
					utility.LogError("Failed to delete file from storage", map[string]interface{}{
						"url":   doc.URL,
						"error": err.Error(),
					})
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
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
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
	// allComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, 0)
	// if err != nil {
	// 	return err
	// }
	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(complianceAssessmentId, requirementId, responsibleUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		return err
	}
	// if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID == nil {
	// 	return utility.NewForbiddenError("score_not_updated")
	// }

	// if complianceAssessmentRequirementResponsibleUnit.ResidualRiskScore == nil {
	// 	return utility.NewForbiddenError("score_not_updated")
	// }
	now := time.Now()
	complianceAssessmentRequirementResponsibleUnit.CompletedByUserID = &currentUserID
	complianceAssessmentRequirementResponsibleUnit.CompletedAt = &now
	complianceAssessmentRequirementResponsibleUnit.UpdatedBYID = &currentUserID

	err = cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// allCompleted := true
	// for _, ru := range allComplianceAssessmentRequirement.ResponsibleUnits {
	// 	if ru == nil {
	// 		continue
	// 	}
	// 	if ru.ID == complianceAssessmentRequirementResponsibleUnit.ID {
	// 		// this one is now completed
	// 		continue
	// 	}
	// 	if !ru.IsCompleted {
	// 		allCompleted = false
	// 		break
	// 	}
	// }

	// // Update parent CAR’s IsAllCompleted
	// car, err := cars.repo.GetByID(requirementId)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// car.IsAllCompleted = allCompleted
	// car.UpdatedBYID = &currentUserID

	// if err := cars.repo.Update(car, tx); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

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
				"compliance", "compliance_assessment", "requirement", nil,
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
	complianceAssessmentRequirementResponsibleUnit.CompletedByUserID = nil
	complianceAssessmentRequirementResponsibleUnit.CompletedAt = nil
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
				"compliance", "compliance_assessment", "requirement", nil,
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
		tx.Rollback()
		return err
	}
	oldComplianceAssessmentRequirement, err := cars.GetFullDetailsByID(assessmentID, requirementID, commentForm.ResponsibleUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}
	complianceAssessmentRequirementResponsibleUnit, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(oldComplianceAssessmentRequirement.ResponsibleUnit.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Fix: Dereference both sides for proper comparison
	if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID != nil &&
		*complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID == approvedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("requirement is already approved"), "invalid_state", "requirement")
	}
	// Check authorization based on assessment type
	if complianceAssessment.Type == "framework" {
		if len(complianceAssessment.Compliances) == 0 || complianceAssessment.Compliances[0].ComplianceRegister == nil {
			tx.Rollback()
			return utility.NewForbiddenError("compliance register not found")
		}
		if currentUserID != complianceAssessment.Compliances[0].ComplianceRegister.ResponsibleUserID {
			tx.Rollback()
			return utility.NewForbiddenError("User is not authorized to approve this assessment")
		}
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

	// Apply score update during approval if score is provided
	if commentForm.ComplianceStatusID != 0 {
		residualScore, err := cars.calculateResidualScore(oldComplianceAssessmentRequirement.ComplianceRequirement.InherentRiskScore, commentForm.Score)
		if err != nil {
			tx.Rollback()
			return err
		}
		ratingID, err := cars.getResidualRatingByScore(residualScore)
		if err != nil {
			tx.Rollback()
			return err
		}
		score := commentForm.Score
		statusID := commentForm.ComplianceStatusID
		complianceAssessmentRequirementResponsibleUnit.Score = &score
		complianceAssessmentRequirementResponsibleUnit.ComplianceStatusID = &statusID
		complianceAssessmentRequirementResponsibleUnit.ResidualRiskScore = &residualScore
		complianceAssessmentRequirementResponsibleUnit.ResidualRiskRatingID = ratingID
		complianceAssessmentRequirementResponsibleUnit.IsNotApplicable = commentForm.IsNotApplicable
	}

	// Update requirement status to "approved"

	complianceAssessmentRequirementResponsibleUnit.IsCompleted = true
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

	allCompletedNow := true
	for _, ru := range allComplianceAssessmentRequirement.ResponsibleUnits {
		if ru == nil {
			continue
		}
		if ru.ID == complianceAssessmentRequirementResponsibleUnit.ID {
			// this one is now approved; treat as completed
			continue
		}
		if !ru.IsCompleted {
			allCompletedNow = false
			break
		}
	}

	if allApproved {
		car.ComplianceAssessmentRequirementStatusID = &approvedStatus.ID
	} else {
		car.ComplianceAssessmentRequirementStatusID = nil
	}
	car.IsAllCompleted = allCompletedNow
	car.UpdatedBYID = &currentUserID

	if err := cars.repo.Update(car, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Move assessment from draft → in_progress only once
	assessmentModel, err := cars.complianceAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	assessmentInProgressStatus, err := cars.complianceAssessmentStatusSvc.GetByType("in_progress")
	if err != nil {
		tx.Rollback()
		return err
	}
	assessmentDraftStatus, err := cars.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return err
	}
	if assessmentModel.ComplianceAssessmentStatusID == assessmentDraftStatus.ID {
		assessmentModel.ComplianceAssessmentStatusID = assessmentInProgressStatus.ID
		if err := cars.repo.ComplianceAssessmentUpdateStatus(assessmentModel, tx); err != nil {
			tx.Rollback()
			return err
		}
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
			"compliance", "compliance_assessment", "requirement", nil,
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
			utility.LogError("Failed to fetch assessment details for notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get titles based on assessment type
		var titleEn, titleAr, referenceNumber string
		if assessmentDetails.Type == "framework" {
			if len(assessmentDetails.Compliances) > 0 && assessmentDetails.Compliances[0].ComplianceRegister != nil {
				titleEn = assessmentDetails.Compliances[0].ComplianceRegister.TitleEn
				titleAr = assessmentDetails.Compliances[0].ComplianceRegister.TitleAr
				referenceNumber = assessmentDetails.Compliances[0].ComplianceRegister.ReferenceNumber
			}
		} else {
			if assessmentDetails.ResponsibleOrganizationUnit != nil {
				titleEn = assessmentDetails.ResponsibleOrganizationUnit.TitleEn + "'s Assessment"
				titleAr = "تقييم " + assessmentDetails.ResponsibleOrganizationUnit.TitleAr
				referenceNumber = ""
			}
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
				TitleEn:              titleEn,
				TitleAr:              titleAr,
				ReferenceNumber:      referenceNumber,
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
					"TitleEn":             titleEn,
					"TitleAr":             titleAr,
					"ReferenceNumber":     referenceNumber,
					"Comment":             commentForm.Comment,
					"ApprovedByFirstName": createdBy.FirstNameEn,
					"ApprovedByLastName":  createdBy.LastNameEn,
					"ApprovalDate":        time.Now().Format("2006-01-02"),
					"RedirectURL":         fmt.Sprintf("%s/compliance/assessments/%d/overview", utility.GetAppURL(), assessmentID),
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
	// Validate if the current status allows reverting (Fix: dereference both sides)
	if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID != nil &&
		*complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID == revertedStatus.ID {
		return utility.NewValidationError(fmt.Errorf("requirement is already reverted"), "invalid_state", "requirement")
	}

	// Check authorization based on assessment type
	if complianceAssessment.Type == "framework" {
		if len(complianceAssessment.Compliances) == 0 || complianceAssessment.Compliances[0].ComplianceRegister == nil {
			return utility.NewForbiddenError("compliance register not found")
		}
		if currentUserID != complianceAssessment.Compliances[0].ComplianceRegister.ResponsibleUserID {
			return utility.NewForbiddenError("User is not authorized to revert this assessment")
		}
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

	car, err := cars.repo.GetByID(requirementID)
	if err != nil {
		tx.Rollback()
		return err
	}
	car.ComplianceAssessmentRequirementStatusID = nil
	car.UpdatedBYID = &currentUserID
	if err := cars.repo.Update(car, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Move assessment back from in_progress → draft on revert
	assessmentModel, err := cars.complianceAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	assessmentInProgressStatus, err := cars.complianceAssessmentStatusSvc.GetByType("in_progress")
	if err != nil {
		tx.Rollback()
		return err
	}
	assessmentDraftStatus, err := cars.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return err
	}
	if assessmentModel.ComplianceAssessmentStatusID == assessmentInProgressStatus.ID {
		assessmentModel.ComplianceAssessmentStatusID = assessmentDraftStatus.ID
		if err := cars.repo.ComplianceAssessmentUpdateStatus(assessmentModel, tx); err != nil {
			tx.Rollback()
			return err
		}
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
			"compliance", "compliance_assessment", "requirement", nil,
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
			utility.LogError("Failed to fetch assessment details for notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get titles based on assessment type
		var titleEn, titleAr, referenceNumber string
		if assessmentDetails.Type == "framework" {
			if len(assessmentDetails.Compliances) > 0 && assessmentDetails.Compliances[0].ComplianceRegister != nil {
				titleEn = assessmentDetails.Compliances[0].ComplianceRegister.TitleEn
				titleAr = assessmentDetails.Compliances[0].ComplianceRegister.TitleAr
				referenceNumber = assessmentDetails.Compliances[0].ComplianceRegister.ReferenceNumber
			}
		} else {
			if assessmentDetails.ResponsibleOrganizationUnit != nil {
				titleEn = assessmentDetails.ResponsibleOrganizationUnit.TitleEn + "'s Assessment"
				titleAr = "تقييم " + assessmentDetails.ResponsibleOrganizationUnit.TitleAr
				referenceNumber = ""
			}
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
				TitleEn:              titleEn,
				TitleAr:              titleAr,
				ReferenceNumber:      referenceNumber,
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
					"TitleEn":             titleEn,
					"TitleAr":             titleAr,
					"ReferenceNumber":     referenceNumber,
					"Comment":             commentForm.Comment,
					"RevertedByFirstName": createdBy.FirstNameEn,
					"RevertedByLastName":  createdBy.LastNameEn,
					"ReversionDate":       time.Now().Format("2006-01-02"),
					"RedirectURL":         fmt.Sprintf("%s/compliance/assessments/%d/overview", utility.GetAppURL(), assessmentID),
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
	// Validate if the current status allows resubmission (Fix: dereference both sides)
	if complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID != nil &&
		*complianceAssessmentRequirementResponsibleUnit.ComplianceAssessmentRequirementStatusID == resubmittedStatus.ID {
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
			"compliance", "compliance_assessment", "requirement", nil,
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
			utility.LogError("Failed to fetch assessment details for notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get titles based on assessment type
		var titleEn, titleAr, referenceNumber string
		if assessmentDetails.Type == "framework" {
			if len(assessmentDetails.Compliances) > 0 && assessmentDetails.Compliances[0].ComplianceRegister != nil {
				titleEn = assessmentDetails.Compliances[0].ComplianceRegister.TitleEn
				titleAr = assessmentDetails.Compliances[0].ComplianceRegister.TitleAr
				referenceNumber = assessmentDetails.Compliances[0].ComplianceRegister.ReferenceNumber
			}
		} else {
			if assessmentDetails.ResponsibleOrganizationUnit != nil {
				titleEn = assessmentDetails.ResponsibleOrganizationUnit.TitleEn + "'s Assessment"
				titleAr = "تقييم " + assessmentDetails.ResponsibleOrganizationUnit.TitleAr
				referenceNumber = ""
			}
		}

		createdBy, err := cars.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get responsible user ID based on assessment type
		var responsibleUserID uint
		if assessmentDetails.Type == "framework" {
			if len(assessmentDetails.Compliances) > 0 && assessmentDetails.Compliances[0].ComplianceRegister != nil {
				responsibleUserID = assessmentDetails.Compliances[0].ComplianceRegister.ResponsibleUserID
			}
		} else {
			// For department type, use review user ID
			responsibleUserID = assessmentDetails.ReviewUserID
		}

		responsibleUser, err := cars.userSvc.GetByID(responsibleUserID)
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
			TitleEn:              titleEn,
			TitleAr:              titleAr,
			ReferenceNumber:      referenceNumber,
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
				"TitleEn":                titleEn,
				"TitleAr":                titleAr,
				"ReferenceNumber":        referenceNumber,
				"Comment":                commentForm.Comment,
				"ResubmittedByFirstName": createdBy.FirstNameEn,
				"ResubmittedByLastName":  createdBy.LastNameEn,
				"ResubmissionDate":       time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/compliance/assessments/%d/overview", utility.GetAppURL(), assessmentID),
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

func (cars *ComplianceAssessmentRequirementService) calculateResidualScore(inherentScore *float64, scorePercent int) (float64, error) {
	inherent := 1.0
	if inherentScore != nil {
		inherent = *inherentScore
	}
	effectiveness := float64(scorePercent) / 100.0
	residual := 1.0 + (inherent-1.0)*(1.0-effectiveness)
	return math.Round(residual*100) / 100, nil
}

func (cars *ComplianceAssessmentRequirementService) getResidualRatingByScore(score float64) (*uint, error) {
	ratings, err := cars.repo.GetComplianceAssessmentRatings()
	if err != nil {
		return nil, err
	}
	for _, rating := range ratings {
		if rating.ScoreFrom != nil && rating.ScoreTo != nil {
			if score >= *rating.ScoreFrom && score <= *rating.ScoreTo {
				return &rating.ID, nil
			}
		}
	}
	return nil, fmt.Errorf("no residual risk rating found for score %.2f: check rating table ranges", score)
}

func (cars *ComplianceAssessmentRequirementService) MarkAsNotApplicable(assessmentID uint, requirementID uint, currentUserID uint, form *complianceDTOs.ComplianceAssessmentRequirementNotApplicableFormDTO) error {
	details, err := cars.GetFullDetailsByID(assessmentID, requirementID, form.ResponsibleUnitID)
	if err != nil {
		return err
	}

	if details.ResponsibleUnit == nil {
		return utility.NewResourceNotFoundError(fmt.Errorf("responsible unit not found for this requirement"), requirementID)
	}
	ru, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(details.ResponsibleUnit.ID)
	if err != nil {
		return err
	}

	tx := cars.repo.BeginTransaction()

	score := form.Score
	ru.IsNotApplicable = true
	ru.Comments = &form.Comments
	ru.Score = &score
	ru.UpdatedBYID = &currentUserID

	if err := cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(ru, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := cars.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentRequirementResponsibleUnit{}.TableName(),
			ru.ID,
			map[string]interface{}{
				"is_not_applicable": true,
				"comments":          form.Comments,
				"score":             form.Score,
			},
			currentUserID,
			"compliance", "compliance_assessment", "requirement", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ComplianceAssessmentRequirementResponsibleUnit{}.TableName(),
				"recordID":  ru.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) CalculateComplianceRequirementScores(assessmentID uint, changedRequirementID uint, tx ...utility.DatabaseTransaction) error {
	assessment, err := cars.complianceAssessmentSvc.GetFullDetailsByIDOptimized(assessmentID, tx...)
	if err != nil {
		return err
	}
	if assessment == nil || len(assessment.Compliances) == 0 || assessment.Compliances[0].ComplianceRegister == nil {
		return fmt.Errorf("assessment/register not found")
	}
	// Use first register for framework/scoring config (all linked registers share the same framework)
	register := assessment.Compliances[0].ComplianceRegister

	// ---------- Framework setup ----------
	frameworkType := ""
	fvMap := make(map[uint]int) // frameworkValueID -> numeric value
	maxFV := 0

	if f := register.ComplianceAssessmentFramework; f != nil {
		if f.AssessmentFrameworkType != nil {
			frameworkType = f.AssessmentFrameworkType.Type
		}
		for _, v := range f.Values {
			if v == nil || v.Value == nil {
				continue
			}
			fvMap[uint(v.ID)] = *v.Value
			if *v.Value > maxFV {
				maxFV = *v.Value
			}
		}
	}

	perReqMax := func() float64 {
		switch frameworkType {
		case "percentage":
			return 100.0
		case "binary":
			return 1.0
		case "point", "maturity":
			if maxFV > 0 {
				return float64(maxFV)
			}
			if frameworkType == "point" {
				return 2.0
			}
			return 5.0
		default:
			return 100.0
		}
	}()

	// ---------- Find the changed requirement ----------
	var changedReq *complianceDTOs.ComplianceAssessmentRequirementDTO
	for _, r := range assessment.AssessmentRequirements {
		if r != nil && r.ID != nil && *r.ID == changedRequirementID {
			changedReq = r
			break
		}
	}
	if changedReq == nil || changedReq.ComplianceRequirementID == 0 || changedReq.ComplianceRequirement.ID == 0 {
		return fmt.Errorf("requirement not found in assessment detail")
	}

	// Helper: does this applicable RU contribute to MaxPossible?
	// Your new rule: count MaxPossible **only if** the RU has a value/score.
	ruContributesMax := func(ru *complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitDTO) bool {
		// must be applicable
		isApplicable := true
		if ru.ComplianceAssessmentFrameworkValue != nil && ru.ComplianceAssessmentFrameworkValue.IsApplicable != nil {
			isApplicable = *ru.ComplianceAssessmentFrameworkValue.IsApplicable
		}
		if !isApplicable {
			return false
		}
		switch frameworkType {
		case "point", "maturity", "binary":
			return ru.ComplianceAssessmentFrameworkValue != nil // has a band/value id
		case "percentage":
			return ru.Score != nil // has an explicit number
		default:
			return ru.Score != nil
		}
	}

	// ---------- (A) Recompute the CHANGED requirement status/score ----------
	var total, maxPossible float64
	applicableSeen := false
	naCount := 0

	for _, ru := range changedReq.ResponsibleUnits {
		if ru == nil {
			continue
		}
		// applicability
		isApplicable := true
		if ru.ComplianceAssessmentFrameworkValue != nil && ru.ComplianceAssessmentFrameworkValue.IsApplicable != nil {
			isApplicable = *ru.ComplianceAssessmentFrameworkValue.IsApplicable
		}
		if !isApplicable {
			naCount++
			continue
		}
		applicableSeen = true

		// Max only if value/score present
		if ruContributesMax(ru) {
			maxPossible += perReqMax
		}

		// Score contribution
		switch frameworkType {
		case "point", "maturity":
			if ru.ComplianceAssessmentFrameworkValue != nil {
				if val, ok := fvMap[uint(ru.ComplianceAssessmentFrameworkValue.ID)]; ok {
					total += float64(val)
				}
			}
		case "percentage":
			if ru.Score != nil {
				total += float64(*ru.Score)
			}
		case "binary":
			if ru.ComplianceAssessmentFrameworkValue != nil {
				if raw, ok := fvMap[uint(ru.ComplianceAssessmentFrameworkValue.ID)]; ok && raw > 0 {
					total += 1.0
				}
			}
		default:
			if ru.Score != nil {
				total += float64(*ru.Score)
			}
		}
	}

	// All NA → mark requirement NA
	if !applicableSeen && naCount > 0 {
		notApplicable, err := cars.complianceStatusSvc.GetByType("not_applicable")
		if err != nil {
			return err
		}
		if err := cars.repo.UpdateComplianceAssessmentRequirementStatus(*changedReq.ID, notApplicable.ID, 0, tx...); err != nil {
			return err
		}
	} else {
		reqPct := 0.0
		if maxPossible > 0 {
			reqPct = (total / maxPossible) * 100.0
			// if reqPct == 100 {
			// 	reqPct = 99.99
			// }
		}
		reqScore := float64(int(reqPct*100)) / 100

		var reqStatusID uint
		if f := register.ComplianceAssessmentFramework; f != nil {
			for _, fv := range f.Values {
				if fv == nil || fv.ScoreMin == nil || fv.ScoreMax == nil {
					continue
				}
				min := float64(*fv.ScoreMin)
				max := float64(*fv.ScoreMax)
				if reqPct >= min && (reqPct < max || (reqPct == max && max == 100)) {
					reqStatusID = fv.ComplianceStatusID
					break
				}
			}
		}
		if err := cars.repo.UpdateComplianceAssessmentRequirementStatus(*changedReq.ID, reqStatusID, reqScore, tx...); err != nil {
			return err
		}
	}

	// ---------- (B) Single-pass compute ALL requirement scores (cache) ----------
	type reqScore struct {
		score   float64 // % rounded (0..100]
		counted bool    // counted in clause average
		clause  uint
	}
	reqScores := make(map[uint]reqScore) // reqID -> reqScore

	for _, r := range assessment.AssessmentRequirements {
		if r == nil || r.ComplianceRequirement.ID == 0 || r.ID == nil {
			continue
		}
		var s, m float64
		var anyApplicable bool
		na := 0

		for _, ru := range r.ResponsibleUnits {
			if ru == nil {
				continue
			}
			// applicability
			isApplicable := true
			if ru.ComplianceAssessmentFrameworkValue != nil && ru.ComplianceAssessmentFrameworkValue.IsApplicable != nil {
				isApplicable = *ru.ComplianceAssessmentFrameworkValue.IsApplicable
			}
			if !isApplicable {
				na++
				continue
			}
			anyApplicable = true

			// Max only if value/score present
			if ruContributesMax(ru) {
				m += perReqMax
			}

			// score
			switch frameworkType {
			case "point", "maturity":
				if ru.ComplianceAssessmentFrameworkValue != nil {
					if val, ok := fvMap[uint(ru.ComplianceAssessmentFrameworkValue.ID)]; ok {
						s += float64(val)
					}
				}
			case "percentage":
				if ru.Score != nil {
					s += float64(*ru.Score)
				}
			case "binary":
				if ru.ComplianceAssessmentFrameworkValue != nil {
					if raw, ok := fvMap[uint(ru.ComplianceAssessmentFrameworkValue.ID)]; ok && raw > 0 {
						s += 1.0
					}
				}
			default:
				if ru.Score != nil {
					s += float64(*ru.Score)
				}
			}
		}

		out := reqScore{clause: r.ComplianceRequirement.ComplianceClause.ID}
		// If all NA → exclude from averages
		if !anyApplicable && na > 0 {
			out.counted = false
			reqScores[*r.ID] = out
			continue
		}
		if m == 0 {
			// Applicable but no values yet → don’t count in averages
			out.counted = false
			reqScores[*r.ID] = out
			continue
		}
		pct := (s / m) * 100.0
		// if pct == 100 {
		// 	pct = 99.99
		// }
		out.score = float64(int(pct*100)) / 100
		out.counted = true
		reqScores[*r.ID] = out
	}
	// ---------- (C) Build clause graph & seed with SELF (direct) requirement scores ----------
	type cAgg struct {
		selfSum  float64 // sum of direct req scores under this clause
		selfCnt  int     // number of counted direct reqs
		final    float64 // computed final for this clause (after roll-up)
		hasFinal bool    // whether final is valid
	}

	agg := map[uint]*cAgg{}           // clauseID -> aggregates
	parentOf := map[uint]*uint{}      // childID -> parentID
	children := map[uint][]uint{}     // parentID -> []childIDs
	seenClause := map[uint]struct{}{} // all clauses present

	// Collect nodes + parent relations
	for _, r := range assessment.AssessmentRequirements {
		if r == nil || r.ComplianceRequirement.ID == 0 {
			continue
		}
		cl := r.ComplianceRequirement.ComplianceClause
		cid := cl.ID
		seenClause[cid] = struct{}{}
		if agg[cid] == nil {
			agg[cid] = &cAgg{}
		}
		if cl.ParentClauseID != nil {
			parentOf[cid] = cl.ParentClauseID
			children[*cl.ParentClauseID] = append(children[*cl.ParentClauseID], cid)
			seenClause[*cl.ParentClauseID] = struct{}{}
			if agg[*cl.ParentClauseID] == nil {
				agg[*cl.ParentClauseID] = &cAgg{}
			}
		} else {
			parentOf[cid] = nil
		}
	}

	// Seed SELF sums/counts from cached requirement scores (reqScores)
	for _, r := range assessment.AssessmentRequirements {
		if r == nil || r.ID == nil {
			continue
		}
		rs, ok := reqScores[*r.ID]
		if !ok || !rs.counted {
			continue
		}
		a := agg[rs.clause]
		a.selfSum += rs.score
		a.selfCnt++
	}

	// ---------- (D) Bottom-up roll-up with SIMPLE (non-weighted) child avg ----------
	var round2 = func(x float64) float64 {
		return math.Round(x*100) / 100
	}

	// find root clauses (no parent)
	roots := make([]uint, 0)
	for cid := range seenClause {
		if parentOf[cid] == nil {
			roots = append(roots, cid)
		}
	}

	visited := map[uint]bool{}
	var postOrder func(uint)
	postOrder = func(cid uint) {
		if visited[cid] {
			return
		}
		visited[cid] = true

		// process children first
		childSum := 0.0
		childCnt := 0
		for _, ch := range children[cid] {
			postOrder(ch)
			if agg[ch] != nil && agg[ch].hasFinal {
				childSum += agg[ch].final
				childCnt++
			}
		}

		// self average (if any)
		var (
			selfAvg  float64
			hasSelf  bool
			childAvg float64
			hasChild bool
		)

		if a := agg[cid]; a != nil {
			if a.selfCnt > 0 {
				selfAvg = round2(a.selfSum / float64(a.selfCnt))
				hasSelf = true
			}
		}
		if childCnt > 0 {
			childAvg = round2(childSum / float64(childCnt)) // SIMPLE average of direct children
			hasChild = true
		}

		// combine per your rule
		switch {
		case hasSelf && hasChild:
			agg[cid].final = round2((selfAvg + childAvg) / 2.0)
			agg[cid].hasFinal = true
		case hasSelf:
			agg[cid].final = selfAvg
			agg[cid].hasFinal = true
		case hasChild:
			agg[cid].final = childAvg
			agg[cid].hasFinal = true
		default:
			agg[cid].hasFinal = false
		}
	}

	// run post-order from each root
	for _, root := range roots {
		postOrder(root)
	}

	// Persist every clause that has a final
	for cid, a := range agg {
		var avg float64
		if a != nil && a.hasFinal {
			avg = a.final
		} else {
			avg = 0 // or skip entirely if you prefer not to touch empty clauses
		}

		var clauseStatusID uint
		if f := register.ComplianceAssessmentFramework; f != nil {
			for _, fv := range f.Values {
				if fv == nil || fv.ScoreMin == nil || fv.ScoreMax == nil {
					continue
				}
				min := float64(*fv.ScoreMin)
				max := float64(*fv.ScoreMax)
				if avg >= min && (avg < max || (avg == max && max == 100)) {
					clauseStatusID = fv.ComplianceStatusID
					break
				}
			}
		}
		if a != nil && a.hasFinal {
			if err := cars.repo.UpdateClauseComplianceStatus(cid, clauseStatusID, avg, tx...); err != nil {
				return err
			}
		}
	}

	// ---------- (E) Overall = SIMPLE average of root clause finals ----------
	var (
		overallSum float64
		overallCnt int
	)
	for _, root := range roots {
		if a := agg[root]; a != nil && a.hasFinal {
			overallSum += a.final
			overallCnt++
		}
	}
	overallPct := 0.0
	if overallCnt > 0 {
		overallPct = overallSum / float64(overallCnt)
	}
	overallScore := round2(overallPct)

	// Persist register
	var overallStatusID uint
	if f := register.ComplianceAssessmentFramework; f != nil {
		for _, fv := range f.Values {
			if fv == nil || fv.ScoreMin == nil || fv.ScoreMax == nil {
				continue
			}
			min := float64(*fv.ScoreMin)
			max := float64(*fv.ScoreMax)
			if overallPct >= min && (overallPct < max || (overallPct == max && max == 100)) {
				overallStatusID = fv.ComplianceStatusID
				break
			}
		}
	}
	// Update compliance status for all linked registers (framework: one, department: many)
	for _, c := range assessment.Compliances {
		if c == nil || c.ComplianceRegister == nil {
			continue
		}
		if err := cars.repo.UpdateComplianceStatus(uint(c.ComplianceRegister.ID), overallStatusID, overallScore, tx...); err != nil {
			return err
		}
	}

	return nil
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func (cars *ComplianceAssessmentRequirementService) AcknowledgeRequirement(assessmentID uint, requirementID uint, currentUserID uint, form *complianceDTOs.AcknowledgeRequirementFormDTO) error {
	if _, err := cars.complianceAssessmentSvc.GetByID(assessmentID); err != nil {
		return err
	}

	// Load the responsible unit record to check acknowledgment status and get its ID
	details, err := cars.GetFullDetailsByID(assessmentID, requirementID, form.ResponsibleUnitID)
	if err != nil {
		return err
	}

	if details.ResponsibleUnit == nil {
		return utility.NewForbiddenError("responsible_unit_not_assigned_to_requirement")
	}

	carru, err := cars.GetComplianceAssessmentRequirementResponsibleUnitByID(details.ResponsibleUnit.ID)
	if err != nil {
		return err
	}

	// Check if this responsible unit is eligible for acknowledgment (must be regulation notification based)
	if carru.IsRegulationNotificationAssessment == nil || !*carru.IsRegulationNotificationAssessment {
		return utility.NewForbiddenError("requirement_not_notification_based")
	}

	if carru.IsAcknowledged {
		return utility.NewValidationError(fmt.Errorf("requirement has already been acknowledged"), "already_acknowledged", "requirement")
	}

	now := time.Now()
	carru.IsAcknowledged = true
	carru.AcknowledgeUserID = &currentUserID
	carru.AcknowledgeComment = form.Comment
	carru.IsCompleted = true
	carru.CompletedByUserID = &currentUserID
	carru.CompletedAt = &now
	carru.UpdatedBYID = &currentUserID

	tx := cars.repo.BeginTransaction()
	if err := cars.repo.UpdateComplianceAssessmentRequirementResponsibleUnit(carru, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (cars *ComplianceAssessmentRequirementService) UpdateDocument(complianceAssessmentID uint, requirementID uint, documentID uint, currentUserID uint, formDTO *complianceDTOs.ComplianceAssessmentRequirementDocumentUpdateFormDTO) error {
	document, err := cars.repo.GetDocumentByID(formDTO.ID)
	if err != nil {
		return err
	}

	now := time.Now()

	if formDTO.Status == "approved" {
		document.ApprovedByID = &currentUserID
		document.UpdatedByID = currentUserID
		document.RejectedByID = nil
		document.ApprovedAt = &now
		document.UpdatedAt = now
		document.RejectedAt = nil
	}

	if formDTO.Status == "rejected" {
		document.RejectedByID = &currentUserID
		document.UpdatedByID = currentUserID
		document.ApprovedByID = nil
		document.RejectedAt = &now
		document.UpdatedAt = now
		document.ApprovedAt = nil
	}

	return cars.repo.UpdateDocument(document)
}
