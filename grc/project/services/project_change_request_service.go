package services

import (
	"errors"
	"path/filepath"
	"strconv"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectChangeRequestServiceInterface interface {
	SetWorkflowService(workflowSvc ProjectChangeRequestWorkflowServiceInterface)
	GetAll(projectID uint, params *dtos.ProjectChangeRequestListQueryParamsDTO) (*dtos.ProjectChangeRequestListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectChangeRequest, error)
	GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectChangeRequestDetailDTO, error)
	Create(projectID uint, currentUserID uint, form *dtos.ProjectChangeRequestFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectChangeRequestFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetSummary(projectID uint) (*dtos.ProjectChangeRequestSummaryDTO, error)
	GetImpactAssessments(changeRequestID uint, tx ...utility.DatabaseTransaction) ([]*models.ProjectImpactAssessment, error)
}

type ProjectChangeRequestService struct {
	repo        repositories.ProjectChangeRequestRepositoryInterface
	activityLog logServices.ActivityLogServiceInterface
	userSvc     organizationServices.UserServiceInterface
	projectSvc  ProjectServiceInterface
	workflowSvc ProjectChangeRequestWorkflowServiceInterface
}

func (pcrs *ProjectChangeRequestService) SetWorkflowService(workflowSvc ProjectChangeRequestWorkflowServiceInterface) {
	pcrs.workflowSvc = workflowSvc
}

func NewProjectChangeRequestService(
	repo repositories.ProjectChangeRequestRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectSvc ProjectServiceInterface,
) ProjectChangeRequestServiceInterface {
	return &ProjectChangeRequestService{
		repo:        repo,
		activityLog: activityLog,
		userSvc:     userSvc,
		projectSvc:  projectSvc,
	}
}

func (pcrs *ProjectChangeRequestService) GetAll(projectID uint, params *dtos.ProjectChangeRequestListQueryParamsDTO) (*dtos.ProjectChangeRequestListResponseDTO, error) {
	if _, err := pcrs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pcrs.repo.GetAll(projectID, params)
}

func (pcrs *ProjectChangeRequestService) GetByID(projectID uint, id uint) (*models.ProjectChangeRequest, error) {
	if _, err := pcrs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	changeRequest, err := pcrs.repo.GetByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return changeRequest, nil
}

func (pcrs *ProjectChangeRequestService) GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectChangeRequestDetailDTO, error) {
	var detailDTO dtos.ProjectChangeRequestDetailDTO

	if _, err := pcrs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	changeRequest, err := pcrs.repo.GetFullDetailsByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	err = mapstructure.Decode(changeRequest, &detailDTO)
	if err != nil {
		return nil, err
	}

	if changeRequest.RequestedOn != nil {
		str := changeRequest.RequestedOn.Format("2006-01-02")
		detailDTO.RequestedOnString = &str
	}

	detailDTO.ImpactAssessments = make([]*dtos.ProjectImpactAssessmentDTO, 0, len(changeRequest.ImpactAssessments))
	for _, ia := range changeRequest.ImpactAssessments {
		iaDTO := &dtos.ProjectImpactAssessmentDTO{
			ID:                            ia.ID,
			ProjectChangeRequestID:        ia.ProjectChangeRequestID,
			ProjectImpactAssessmentTypeID: ia.ProjectImpactAssessmentTypeID,
			CurrentBudget:                 ia.CurrentBudget,
			RevisedBudget:                 ia.RevisedBudget,
			CurrentResourceCount:          ia.CurrentResourceCount,
			RevisedResourceCount:          ia.RevisedResourceCount,
			TechnicalComplexity:           ia.TechnicalComplexity,
		}
		if ia.CurrentDate != nil {
			cdStr := ia.CurrentDate.Format("2006-01-02")
			iaDTO.CurrentDateString = &cdStr
		}
		if ia.RevisedDate != nil {
			rdStr := ia.RevisedDate.Format("2006-01-02")
			iaDTO.RevisedDateString = &rdStr
		}
		if ia.ProjectImpactAssessmentType != nil {
			iaDTO.ProjectImpactAssessmentType = &dtos.ProjectImpactAssessmentTypeDTO{
				ID:      ia.ProjectImpactAssessmentType.ID,
				TitleEn: ia.ProjectImpactAssessmentType.TitleEn,
				TitleAr: ia.ProjectImpactAssessmentType.TitleAr,
				Type:    ia.ProjectImpactAssessmentType.Type,
				Status:  ia.ProjectImpactAssessmentType.Status,
			}
		}
		detailDTO.ImpactAssessments = append(detailDTO.ImpactAssessments, iaDTO)
	}

	detailDTO.CreatedAt = changeRequest.CreatedAt
	detailDTO.UpdatedAt = changeRequest.UpdatedAt

	return &detailDTO, nil
}

func (pcrs *ProjectChangeRequestService) Create(projectID uint, currentUserID uint, form *dtos.ProjectChangeRequestFormDTO) (uint, error) {
	project, err := pcrs.projectSvc.GetByID(projectID)
	if err != nil {
		return 0, err
	}

	requestedOn, err := utility.ParseDate(form.RequestedOn)
	if err != nil {
		return 0, utility.NewValidationError(err, "invalid_date", "requested_on")
	}

	refNumber, err := pcrs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	impactAssessments, err := parseImpactAssessments(form.ImpactAssessments)
	if err != nil {
		return 0, err
	}

	changeRequest := &models.ProjectChangeRequest{
		ReferenceNumber:       refNumber,
		ProjectID:             project.ID,
		Title:                 form.Title,
		Description:           form.Description,
		BusinessJustification: form.BusinessJustification,
		AnticipatedBenefits:   form.AnticipatedBenefits,
		RisksOfDeferral:       form.RisksOfDeferral,
		RequestedByUserID:     form.RequestedByUserID,
		RequestedOn:           requestedOn,
		CreatedBYID:           currentUserID,
		CreatedAt:             time.Now(),
	}

	tx := pcrs.repo.BeginTransaction()

	id, err := pcrs.repo.Create(changeRequest, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(impactAssessments) > 0 {
		for _, ia := range impactAssessments {
			ia.ProjectChangeRequestID = id
		}
		if err := pcrs.repo.CreateImpactAssessments(impactAssessments, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var documents []*models.ProjectChangeRequestDocument
	for _, doc := range form.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.ProjectChangeRequestDocument{
			ProjectChangeRequestID: id,
			URL:                    url,
			Name:                   doc.FileName,
			Size:                   doc.Size,
			Token:                  encryptedFileName,
			OriginalName:           &doc.OriginalName,
			Type:                   doc.Type,
			CreatedByID:            currentUserID,
			CreatedAt:              time.Now(),
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := pcrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	for _, doc := range documents {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	go func() {
		activityData := map[string]interface{}{
			"project_id":             projectID,
			"reference_number":       refNumber,
			"title":                  form.Title,
			"description":            form.Description,
			"business_justification": form.BusinessJustification,
			"anticipated_benefits":   form.AnticipatedBenefits,
			"risks_of_deferral":      form.RisksOfDeferral,
			"requested_by_user_id":   form.RequestedByUserID,
			"requested_on":           form.RequestedOn,
		}
		err := pcrs.activityLog.LogActivity(models.CreateAction, models.ProjectChangeRequest{}.TableName(), id, activityData, currentUserID, "project", "project_register", "", nil)
		if err != nil {
			utility.LogError("Failed to log create project change request activity", map[string]interface{}{
				"recordID": id,
				"error":    err.Error(),
			})
		}
	}()

	return id, nil
}

func (pcrs *ProjectChangeRequestService) Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectChangeRequestFormDTO) error {
	if _, err := pcrs.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	changeRequest, err := pcrs.repo.GetFullDetailsByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	oldData := map[string]interface{}{
		"title":                  changeRequest.Title,
		"description":            changeRequest.Description,
		"business_justification": changeRequest.BusinessJustification,
		"anticipated_benefits":   changeRequest.AnticipatedBenefits,
		"risks_of_deferral":      changeRequest.RisksOfDeferral,
		"requested_by_user_id":   changeRequest.RequestedByUserID,
		"requested_on":           changeRequest.RequestedOn,
	}

	requestedOn, err := utility.ParseDate(form.RequestedOn)
	if err != nil {
		return utility.NewValidationError(err, "invalid_date", "requested_on")
	}

	impactAssessments, err := parseImpactAssessments(form.ImpactAssessments)
	if err != nil {
		return err
	}

	changeRequest.Title = form.Title
	changeRequest.Description = form.Description
	changeRequest.BusinessJustification = form.BusinessJustification
	changeRequest.AnticipatedBenefits = form.AnticipatedBenefits
	changeRequest.RisksOfDeferral = form.RisksOfDeferral
	changeRequest.RequestedByUserID = form.RequestedByUserID
	changeRequest.RequestedOn = requestedOn
	changeRequest.UpdatedBYID = &currentUserID
	now := time.Now()
	changeRequest.UpdatedAt = &now

	tx := pcrs.repo.BeginTransaction()

	err = pcrs.repo.Update(changeRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrs.repo.DeleteImpactAssessments(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(impactAssessments) > 0 {
		for _, ia := range impactAssessments {
			ia.ProjectChangeRequestID = id
		}
		if err := pcrs.repo.CreateImpactAssessments(impactAssessments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Document handling
	var newDocsToCreate []*models.ProjectChangeRequestDocument
	var docsToDelete []*uint
	var documentsToDelete []*models.ProjectChangeRequestDocument

	for _, doc := range form.Documents {
		if (doc.IsNewFile != nil && *doc.IsNewFile) || (doc.ID == nil && doc.FileName != "") {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			document := &models.ProjectChangeRequestDocument{
				ProjectChangeRequestID: id,
				URL:                    url,
				Name:                   doc.FileName,
				Size:                   doc.Size,
				Token:                  encryptedFileName,
				OriginalName:           &doc.OriginalName,
				Type:                   doc.Type,
				CreatedByID:            currentUserID,
				CreatedAt:              time.Now(),
			}
			newDocsToCreate = append(newDocsToCreate, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := pcrs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return utility.NewResourceNotFoundError(err, id)
					}
					tx.Rollback()
					return err
				}
				documentsToDelete = append(documentsToDelete, existingDocument)
				docsToDelete = append(docsToDelete, &existingDocument.ID)
			}
		}
	}

	if len(docsToDelete) > 0 {
		if err := pcrs.repo.DeleteDocuments(docsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(newDocsToCreate) > 0 {
		if err := pcrs.repo.CreateDocuments(newDocsToCreate, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, doc := range newDocsToCreate {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to remove document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		newData := map[string]interface{}{
			"title":                  form.Title,
			"description":            form.Description,
			"business_justification": form.BusinessJustification,
			"anticipated_benefits":   form.AnticipatedBenefits,
			"risks_of_deferral":      form.RisksOfDeferral,
			"requested_by_user_id":   form.RequestedByUserID,
			"requested_on":           form.RequestedOn,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			err := pcrs.activityLog.LogActivity(models.UpdateAction, models.ProjectChangeRequest{}.TableName(), id, changes, currentUserID, "project", "project_register", "", nil)
			if err != nil {
				utility.LogError("Failed to log update project change request activity", map[string]interface{}{
					"recordID": id,
					"error":    err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pcrs *ProjectChangeRequestService) Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pcrs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		changeRequest, err := pcrs.repo.GetByID(projectID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		tx := pcrs.repo.BeginTransaction()

		// 1. Delete review comments and review comment documents
		if err := pcrs.repo.DeleteReviewComments(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 2. Delete workflow histories
		if err := pcrs.repo.DeleteWorkflowHistories(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 3. Delete workflows, reviews, and participants
		if pcrs.workflowSvc != nil {
			if err := pcrs.workflowSvc.DeleteAllWorkflowsForProjectChangeRequest(id, tx); err != nil {
				tx.Rollback()
				failedIds = append(failedIds, id)
				continue
			}
		}

		// 4. Delete impact assessments
		if err := pcrs.repo.DeleteImpactAssessments(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 5. Delete documents
		if err := pcrs.repo.DeleteDocumentsByChangeRequestID(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 6. Delete parent change request
		err = pcrs.repo.Delete(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := tx.Commit(); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedItems = append(deletedItems, map[string]interface{}{
			"id":               id,
			"reference_number": changeRequest.ReferenceNumber,
			"title":            changeRequest.Title,
		})

		folderPath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "change-requests", strconv.FormatUint(uint64(id), 10))
		if err := utility.RemoveFolder(folderPath); err != nil {
			utility.LogError("Failed to remove change request storage folder", map[string]interface{}{
				"changeRequestID": id,
				"path":            folderPath,
				"error":           err.Error(),
			})
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := pcrs.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectChangeRequest{}.TableName(),
				successfulIds,
				deletedItems,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete project change requests activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (pcrs *ProjectChangeRequestService) GetSummary(projectID uint) (*dtos.ProjectChangeRequestSummaryDTO, error) {
	if projectID != 0 {
		if _, err := pcrs.projectSvc.GetByID(projectID); err != nil {
			return nil, err
		}
	}

	return pcrs.repo.GetSummary(projectID)
}

func parseImpactAssessments(forms []dtos.ProjectImpactAssessmentFormDTO) ([]*models.ProjectImpactAssessment, error) {
	var impactAssessments []*models.ProjectImpactAssessment
	for _, f := range forms {
		var currentDate, revisedDate *time.Time
		if f.CurrentDate != nil && *f.CurrentDate != "" {
			cd, err := utility.ParseDate(*f.CurrentDate)
			if err != nil {
				return nil, utility.NewValidationError(err, "invalid_date", "current_date")
			}
			currentDate = cd
		}
		if f.RevisedDate != nil && *f.RevisedDate != "" {
			rd, err := utility.ParseDate(*f.RevisedDate)
			if err != nil {
				return nil, utility.NewValidationError(err, "invalid_date", "revised_date")
			}
			revisedDate = rd
		}
		ia := &models.ProjectImpactAssessment{
			ProjectImpactAssessmentTypeID: f.ProjectImpactAssessmentTypeID,
			CurrentDate:                   currentDate,
			RevisedDate:                   revisedDate,
			CurrentBudget:                 f.CurrentBudget,
			RevisedBudget:                 f.RevisedBudget,
			CurrentResourceCount:          f.CurrentResourceCount,
			RevisedResourceCount:          f.RevisedResourceCount,
			TechnicalComplexity:           f.TechnicalComplexity,
		}
		impactAssessments = append(impactAssessments, ia)
	}
	return impactAssessments, nil
}

func (pcrs *ProjectChangeRequestService) GetImpactAssessments(changeRequestID uint, tx ...utility.DatabaseTransaction) ([]*models.ProjectImpactAssessment, error) {
	return pcrs.repo.GetImpactAssessments(changeRequestID, tx...)
}
