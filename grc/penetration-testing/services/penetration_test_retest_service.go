package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestRetestServiceInterface interface {
	GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error)
	GetByID(id uint) (*models.PenetrationTestRetest, error)
	GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestRetestDetailDTO, error)
	Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestRetestFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFormDTO) error
	GetEligibleFindings(id uint) ([]*penetrationTestDTOs.PenetrationTestFindingListDTO, error)
	UpdateFindings(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFindingUpdateFormDTO) error
	DeleteFinding(id uint, retestFindingID uint, currentUserID uint) error
	SaveFinding(id uint, retestFindingID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFindingSaveFormDTO) error
	DecideFinding(id uint, retestFindingID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFindingResultFormDTO) error
	GetDocumentPath(id uint, retestFindingID uint, token string) (string, error)
	Complete(id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type PenetrationTestRetestService struct {
	repo                      repositories.PenetrationTestRetestRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	penetrationTestSvc        PenetrationTestServiceInterface
	penetrationTestFindingSvc PenetrationTestFindingServiceInterface
	executorSvc               PenetrationTestExecutorServiceInterface
	actionPlanStatusSvc       masterServices.ActionPlanStatusServiceInterface
	statusUpdateRepo          repositories.PenetrationTestFindingStatusUpdateRepositoryInterface
	retestResultSvc           PenetrationTestRetestResultServiceInterface
	retestStatusSvc           PenetrationTestRetestStatusServiceInterface
}

func NewPenetrationTestRetestService(
	repo repositories.PenetrationTestRetestRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
	penetrationTestFindingSvc PenetrationTestFindingServiceInterface,
	executorSvc PenetrationTestExecutorServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	statusUpdateRepo repositories.PenetrationTestFindingStatusUpdateRepositoryInterface,
	retestResultSvc PenetrationTestRetestResultServiceInterface,
	retestStatusSvc PenetrationTestRetestStatusServiceInterface,
) PenetrationTestRetestServiceInterface {
	return &PenetrationTestRetestService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		penetrationTestSvc:        penetrationTestSvc,
		penetrationTestFindingSvc: penetrationTestFindingSvc,
		executorSvc:               executorSvc,
		actionPlanStatusSvc:       actionPlanStatusSvc,
		statusUpdateRepo:          statusUpdateRepo,
		retestResultSvc:           retestResultSvc,
		retestStatusSvc:           retestStatusSvc,
	}
}

func (ptrs *PenetrationTestRetestService) GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error) {
	return ptrs.repo.GetAll(params)
}

func (ptrs *PenetrationTestRetestService) GetByID(id uint) (*models.PenetrationTestRetest, error) {
	value, err := ptrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (ptrs *PenetrationTestRetestService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestRetestDetailDTO, error) {
	value, err := ptrs.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &penetrationTestDTOs.PenetrationTestRetestDetailDTO{}
	if err := mapstructure.Decode(value, detail); err != nil {
		return nil, err
	}
	detail.CreatedAt = value.CreatedAt
	detail.UpdatedAt = value.UpdatedAt
	detail.StartDate = value.StartDate
	detail.EndDate = value.EndDate
	if value.PenetrationTestRetestStatus != nil {
		detail.PenetrationTestRetestStatus = &penetrationTestDTOs.PenetrationTestRetestStatusDTO{
			ID:        value.PenetrationTestRetestStatus.ID,
			TitleEn:   value.PenetrationTestRetestStatus.TitleEn,
			TitleAr:   value.PenetrationTestRetestStatus.TitleAr,
			ColorCode: value.PenetrationTestRetestStatus.ColorCode,
			Type:      value.PenetrationTestRetestStatus.Type,
		}
	}
	detail.Findings = make([]*penetrationTestDTOs.PenetrationTestRetestFindingDTO, 0, len(value.RetestFindings))
	for _, retestFinding := range value.RetestFindings {
		finding := retestFinding.PenetrationTestFinding
		if finding == nil {
			continue
		}
		findingDTO := &penetrationTestDTOs.PenetrationTestRetestFindingDTO{
			ID:                       retestFinding.ID,
			PenetrationTestFindingID: retestFinding.PenetrationTestFindingID,
			FindingReferenceNumber:   finding.ReferenceNumber,
			FindingTitle:             finding.Title,
			Comment:                  retestFinding.Comment,
			CommentAt:                retestFinding.CommentAt,
			ResultAt:                 retestFinding.ResultAt,
			Documents:                retestFinding.Documents,
		}
		if retestFinding.RetestResult != nil {
			findingDTO.Result = &penetrationTestDTOs.PenetrationTestRetestResultDTO{
				ID:      retestFinding.RetestResult.ID,
				TitleEn: retestFinding.RetestResult.TitleEn,
				TitleAr: retestFinding.RetestResult.TitleAr,
				Type:    retestFinding.RetestResult.Type,
			}
		}
		if retestFinding.CommentByUser != nil {
			findingDTO.CommentByUser = &organizationDTOs.UserMinimalDTO{}
			if err := mapstructure.Decode(retestFinding.CommentByUser, findingDTO.CommentByUser); err != nil {
				return nil, err
			}
		}
		if retestFinding.ResultByUser != nil {
			findingDTO.ResultByUser = &organizationDTOs.UserMinimalDTO{}
			if err := mapstructure.Decode(retestFinding.ResultByUser, findingDTO.ResultByUser); err != nil {
				return nil, err
			}
		}
		detail.Findings = append(detail.Findings, findingDTO)
		if (retestFinding.RetestResult == nil || retestFinding.RetestResult.Type != "pass") && finding.SeverityLevel != nil {
			switch finding.SeverityLevel.Type {
			case "critical":
				detail.RemainingCritical++
			case "high":
				detail.RemainingHigh++
			case "medium":
				detail.RemainingMedium++
			case "low":
				detail.RemainingLow++
			}
			detail.RemainingTotal++
		}
	}
	return detail, nil
}

func (ptrs *PenetrationTestRetestService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestRetestFormDTO) (uint, error) {
	if err := ptrs.validatePenetrationTestRetestFormDTO(formDTO); err != nil {
		return 0, err
	}
	startDate, err := utility.ParseMandatoryDate(formDTO.StartDate)
	if err != nil {
		return 0, utility.NewValidationError(err, "invalid", "start_date")
	}
	var endDate *time.Time
	if formDTO.EndDate != nil && *formDTO.EndDate != "" {
		endDate, err = utility.ParseDate(*formDTO.EndDate)
		if err != nil {
			return 0, utility.NewValidationError(err, "invalid", "end_date")
		}
		if endDate.Before(startDate) {
			return 0, utility.NewValidationError(fmt.Errorf("end date cannot be before start date"), "invalid", "end_date")
		}
	}
	retestStatus, err := ptrs.retestStatusSvc.GetByType("draft")
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("invalid status"), "invalid", "status")
	}
	value := &models.PenetrationTestRetest{
		PenetrationTestID:             formDTO.PenetrationTestID,
		StartDate:                     startDate,
		EndDate:                       endDate,
		PenetrationTestExecutorID:     formDTO.PenetrationTestExecutorID,
		PenetrationTestRetestStatusID: retestStatus.ID,
		Notes:                         formDTO.Notes,
	}
	active, err := ptrs.repo.GetActiveByPenetrationTestID(value.PenetrationTestID)
	if err != nil {
		return 0, err
	}
	if active != nil {
		return 0, utility.NewValidationError(fmt.Errorf("an active re-test already exists"), "retest_already_active", "re-test")
	}
	tx := ptrs.repo.BeginTransaction()
	value.ReferenceNumber, err = ptrs.repo.GenerateReferenceNumber(value.PenetrationTestID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.CreatedBYID = currentUserID
	value.CreatedAt = time.Now()
	id, err := ptrs.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	for _, findingID := range formDTO.FindingIDs {
		_, err := ptrs.repo.CreateRetestFinding(&models.PenetrationTestRetestFinding{
			PenetrationTestRetestID:  value.ID,
			PenetrationTestFindingID: findingID,
		}, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err := ptrs.penetrationTestSvc.UpdateStatus(value.PenetrationTestID, "under retest", tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "create penetration test re-test", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return id, nil
}

func (ptrs *PenetrationTestRetestService) Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFormDTO) error {
	if err := ptrs.validatePenetrationTestRetestFormDTO(form); err != nil {
		return err
	}
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	startDate, err := utility.ParseMandatoryDate(form.StartDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(err, "invalid", "start_date")
	}
	var endDate *time.Time
	if form.EndDate != nil && *form.EndDate != "" {
		endDate, err = utility.ParseDate(*form.EndDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid", "end_date")
		}
		if endDate.Before(startDate) {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("end date cannot be before start date"), "invalid", "end_date")
		}
	}
	updated := &models.PenetrationTestRetest{
		PenetrationTestID:             form.PenetrationTestID,
		StartDate:                     startDate,
		EndDate:                       endDate,
		PenetrationTestExecutorID:     form.PenetrationTestExecutorID,
		PenetrationTestRetestStatusID: existing.PenetrationTestRetestStatusID,
		Notes:                         form.Notes,
	}
	if existing.PenetrationTestRetestStatus != nil && existing.PenetrationTestRetestStatus.Type == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be updated"), "retest_already_completed", "re-test")
	}
	updated.ID = existing.ID
	updated.ReferenceNumber = existing.ReferenceNumber
	updated.CreatedBYID = existing.CreatedBYID
	updated.CreatedAt = existing.CreatedAt
	updated.UpdatedBYID = &currentUserID
	now := time.Now()
	updated.UpdatedAt = &now
	if err := ptrs.repo.Update(updated, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := ptrs.repo.DeleteRetestFindings(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	for _, findingID := range form.FindingIDs {
		if _, err := ptrs.repo.CreateRetestFinding(&models.PenetrationTestRetestFinding{
			PenetrationTestRetestID:  id,
			PenetrationTestFindingID: findingID,
		}, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "update penetration test re-test", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) GetEligibleFindings(id uint) ([]*penetrationTestDTOs.PenetrationTestFindingListDTO, error) {
	existing, err := ptrs.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	closedStatus, err := ptrs.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return nil, err
	}
	return ptrs.repo.GetEligibleOpenFindings(id, existing.PenetrationTestID, closedStatus.ID)
}

func (ptrs *PenetrationTestRetestService) UpdateFindings(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFindingUpdateFormDTO) error {
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if existing.PenetrationTestRetestStatus != nil && existing.PenetrationTestRetestStatus.Type == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be changed"), "retest_already_completed", "re-test")
	}
	assigned := make(map[uint]bool)
	for _, retestFinding := range existing.RetestFindings {
		assigned[retestFinding.PenetrationTestFindingID] = true
	}
	seen := make(map[uint]bool)
	for _, findingID := range form.FindingIDs {
		if findingID == 0 || seen[findingID] {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("finding IDs must be unique and valid"), "invalid", "finding IDs")
		}
		seen[findingID] = true
		if assigned[findingID] {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("finding is already assigned to this re-test"), "invalid", "finding IDs")
		}
		finding, err := ptrs.penetrationTestFindingSvc.GetByID(findingID)
		if err != nil || finding.PenetrationTestID != existing.PenetrationTestID {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("finding does not belong to penetration test"), "not_exist", "finding IDs")
		}
		if _, err := ptrs.repo.CreateRetestFinding(&models.PenetrationTestRetestFinding{
			PenetrationTestRetestID:  id,
			PenetrationTestFindingID: findingID,
		}, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "assign penetration test re-test findings", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) DeleteFinding(id uint, retestFindingID uint, currentUserID uint) error {
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if existing.PenetrationTestRetestStatus != nil && existing.PenetrationTestRetestStatus.Type == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be changed"), "retest_already_completed", "re-test")
	}
	retestFinding, err := ptrs.repo.GetRetestFindingByID(id, retestFindingID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, retestFindingID)
		}
		return err
	}
	storagePath := filepath.Join("penetration-testing", "penetration-test-retests", strconv.FormatUint(uint64(id), 10), "findings", strconv.FormatUint(uint64(retestFindingID), 10), "documents")
	if len(retestFinding.Documents) > 0 {
		documentIDs := make([]*uint, 0, len(retestFinding.Documents))
		for _, doc := range retestFinding.Documents {
			documentIDs = append(documentIDs, &doc.ID)
		}
		if err := ptrs.repo.DeleteDocuments(documentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
		for _, doc := range retestFinding.Documents {
			if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if err := ptrs.repo.DeleteRetestFindingByID(id, retestFindingID, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "unassign penetration test re-test finding", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) UpdateStatus(id uint, status string, tx ...utility.DatabaseTransaction) error {
	statusMaster, err := ptrs.retestStatusSvc.GetByType(status)
	if err != nil {
		return err
	}
	return ptrs.repo.UpdateStatus(id, statusMaster.ID, tx...)
}

func (ptrs *PenetrationTestRetestService) SaveFinding(id uint, retestFindingID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFindingSaveFormDTO) error {
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if existing.PenetrationTestRetestStatus != nil && existing.PenetrationTestRetestStatus.Type == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be changed"), "retest_already_completed", "re-test")
	}
	retestFinding, err := ptrs.repo.GetRetestFindingByID(id, retestFindingID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, retestFindingID)
		}
		return err
	}
	now := time.Now()
	retestFinding.Comment = &form.Comment
	retestFinding.CommentBYID = &currentUserID
	retestFinding.CommentAt = &now
	if err := ptrs.repo.UpdateRetestFinding(retestFinding, tx); err != nil {
		tx.Rollback()
		return err
	}
	var documents []*models.PenetrationTestRetestFindingDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.PenetrationTestRetestFindingDocument
	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("penetration-testing", "penetration-test-retests", uint64(id), "findings", uint64(retestFindingID), "documents")
			documents = append(documents, &models.PenetrationTestRetestFindingDocument{
				PenetrationTestRetestFindingID: retestFindingID,
				URL:                            url,
				Name:                           doc.FileName,
				Size:                           doc.Size,
				Token:                          encryptedFilename,
				OriginalName:                   doc.OriginalName,
				Type:                           doc.Type,
				CreatedBYID:                    currentUserID,
			})
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := ptrs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utility.NewResourceNotFoundError(err, *doc.ID)
				}
				return err
			}
			documentsToDelete = append(documentsToDelete, existingDocument)
			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
		}
	}
	if len(documents) > 0 {
		if err := ptrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIDs) > 0 {
		if err := ptrs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	storagePath := filepath.Join("penetration-testing", "penetration-test-retests", strconv.FormatUint(uint64(id), 10), "findings", strconv.FormatUint(uint64(retestFindingID), 10), "documents")
	for _, doc := range documents {
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, doc := range documentsToDelete {
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := ptrs.UpdateStatus(id, "in-progress", tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "save penetration test re-test finding", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) DecideFinding(id uint, retestFindingID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestRetestFindingResultFormDTO) error {
	tx := ptrs.repo.BeginTransaction()
	existing, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if existing.PenetrationTestRetestStatus != nil && existing.PenetrationTestRetestStatus.Type == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("completed re-tests cannot be changed"), "retest_already_completed", "re-test")
	}
	retestFinding, err := ptrs.repo.GetRetestFindingByID(id, retestFindingID, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, retestFindingID)
		}
		return err
	}
	if retestFinding.Comment == nil || *retestFinding.Comment == "" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("comment is required before deciding a finding"), "required", "comment")
	}
	result, err := ptrs.retestResultSvc.GetByID(form.ResultID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewValidationError(fmt.Errorf("invalid result"), "invalid", "result")
		}
		return err
	}
	now := time.Now()
	retestFinding.RetestResultID = &result.ID
	retestFinding.ResultAt = &now
	retestFinding.ResultBYID = &currentUserID
	if err := ptrs.repo.UpdateRetestFinding(retestFinding, tx); err != nil {
		tx.Rollback()
		return err
	}
	if result.Type == "pass" {
		closedStatus, err := ptrs.actionPlanStatusSvc.GetByType("close")
		if err != nil {
			tx.Rollback()
			return err
		}
		percentage := uint(100)
		if err := ptrs.penetrationTestFindingSvc.StatusUpdate(retestFinding.PenetrationTestFindingID, currentUserID, closedStatus.ID, &percentage, tx); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := ptrs.statusUpdateRepo.Create(&models.PenetrationTestFindingStatusUpdate{
			PenetrationTestFindingID: retestFinding.PenetrationTestFindingID,
			ActionPlanStatusID:       closedStatus.ID,
			Percentage:               &percentage,
			Comment:                  *retestFinding.Comment,
			CreatedBYID:              currentUserID,
			CreatedAt:                time.Now(),
		}, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	allDecided := len(existing.RetestFindings) > 0
	for _, finding := range existing.RetestFindings {
		if finding.ID == retestFinding.ID {
			continue
		}
		if finding.RetestResultID == nil {
			allDecided = false
			break
		}
	}
	if allDecided {
		completedStatus, err := ptrs.retestStatusSvc.GetByType("completed")
		if err != nil {
			tx.Rollback()
			return err
		}
		now := time.Now()
		existing.PenetrationTestRetestStatusID, existing.EndDate = completedStatus.ID, &now
		existing.UpdatedBYID, existing.UpdatedAt = &currentUserID, &now
		if err := ptrs.repo.Update(existing, tx); err != nil {
			tx.Rollback()
			return err
		}
	} else if err := ptrs.UpdateStatus(id, "in-progress", tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "decide penetration test re-test finding", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) GetDocumentPath(id uint, retestFindingID uint, token string) (string, error) {
	value, err := ptrs.repo.GetRetestFindingByID(id, retestFindingID)
	if err != nil {
		return "", err
	}
	if len(value.Documents) == 0 {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("no documents found"), retestFindingID)
	}
	filenameBytes, err := utility.DecryptData(token)
	if err != nil || string(filenameBytes) == "" {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("resource not found"), retestFindingID)
	}
	storagePath := filepath.Join("penetration-testing", "penetration-test-retests", strconv.FormatUint(uint64(id), 10), "findings", strconv.FormatUint(uint64(retestFindingID), 10), "documents")
	fullPath := filepath.Join(config.GetEnv("STORAGE_PATH", ""), storagePath, string(filenameBytes))
	if _, err := os.Stat(fullPath); err != nil {
		return "", utility.NewResourceNotFoundError(err, retestFindingID)
	}
	return fullPath, nil
}

func (ptrs *PenetrationTestRetestService) Complete(id uint, currentUserID uint) error {
	tx := ptrs.repo.BeginTransaction()
	value, err := ptrs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if value.PenetrationTestRetestStatus != nil && value.PenetrationTestRetestStatus.Type == "completed" {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("re-test is already completed"), "retest_already_completed", "re-test")
	}
	completedStatus, err := ptrs.retestStatusSvc.GetByType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}
	now := time.Now()
	value.PenetrationTestRetestStatusID = completedStatus.ID
	value.EndDate = &now
	value.UpdatedBYID = &currentUserID
	value.UpdatedAt = &now
	if err := ptrs.repo.Update(value, tx); err != nil {
		tx.Rollback()
		return err
	}
	closedStatus, err := ptrs.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		tx.Rollback()
		return err
	}
	percentage := uint(100)
	for _, retestFinding := range value.RetestFindings {
		if retestFinding.RetestResult == nil || retestFinding.RetestResult.Type != "pass" {
			continue
		}
		if retestFinding.PenetrationTestFinding != nil && retestFinding.PenetrationTestFinding.ActionPlanStatusID == closedStatus.ID {
			continue
		}
		if err := ptrs.penetrationTestFindingSvc.StatusUpdate(retestFinding.PenetrationTestFindingID, currentUserID, closedStatus.ID, &percentage, tx); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := ptrs.statusUpdateRepo.Create(&models.PenetrationTestFindingStatusUpdate{
			PenetrationTestFindingID: retestFinding.PenetrationTestFindingID,
			ActionPlanStatusID:       closedStatus.ID,
			Percentage:               &percentage,
			Comment:                  "Finding verified and closed through re-test completion",
			CreatedBYID:              currentUserID,
			CreatedAt:                time.Now(),
		}, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptrs.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptrs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test re-test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptrs.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestRetest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "complete penetration test re-test", nil); err != nil {
				utility.LogError("Failed to log penetration test re-test activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptrs *PenetrationTestRetestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deleted []map[string]interface{}
	for _, id := range ids {
		value, err := ptrs.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		tx := ptrs.repo.BeginTransaction()
		if err := ptrs.repo.DeleteRetestFindings(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := ptrs.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		successfulIDs = append(successfulIDs, id)
		deleted = append(deleted, map[string]interface{}{"id": id, "reference_number": value.ReferenceNumber})
	}
	if len(successfulIDs) > 0 && ptrs.activityLogSvc != nil {
		go ptrs.activityLogSvc.LogBulkActivities(models.DeleteAction, models.PenetrationTestRetest{}.TableName(), successfulIDs, deleted, currentUserID, nil)
	}
	return failedIDs, nil
}

func (ptrs *PenetrationTestRetestService) validatePenetrationTestRetestFormDTO(formDTO *penetrationTestDTOs.PenetrationTestRetestFormDTO) error {
	if _, err := ptrs.penetrationTestSvc.GetByID(formDTO.PenetrationTestID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test does not exist"), "not_exist", "penetration test")
	}
	if _, err := ptrs.executorSvc.GetByID(formDTO.PenetrationTestExecutorID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test executor does not exist"), "not_exist", "penetration test executor")
	}
	seen := make(map[uint]bool)
	for _, findingID := range formDTO.FindingIDs {
		if findingID == 0 || seen[findingID] {
			return utility.NewValidationError(fmt.Errorf("finding IDs must be unique and valid"), "invalid", "finding IDs")
		}
		finding, err := ptrs.penetrationTestFindingSvc.GetByID(findingID)
		if err != nil || finding.PenetrationTestID != formDTO.PenetrationTestID {
			return utility.NewValidationError(fmt.Errorf("finding does not belong to penetration test"), "not_exist", "finding IDs")
		}
		seen[findingID] = true
	}
	return nil
}
