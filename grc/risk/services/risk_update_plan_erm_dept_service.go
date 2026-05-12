package services

import (
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type RiskUpdatePlanERMDeptServiceInterface interface {
	CreateERMDepartmentEvaluation(planID uint, currentUserID uint) error
	MapRisksToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapRisksToUpdatePlanDTO) error
	GetERMDeptEvaluationRisks(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error)
	MapProcessesToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapProcessesToUpdatePlanDTO) error
	GetERMDeptEvaluationProcesses(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error)
	MapControlsToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapControlsToUpdatePlanDTO) error
	GetERMDeptEvaluationControls(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error)
	MarkERMDeptEvaluationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error

	// ERM Department Evaluation - Documents
	GetERMDeptEvaluationDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error)
	CreateERMDeptEvaluationDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetERMDeptEvaluationDocumentPath(planID uint, token string) (string, error)
}

type RiskUpdatePlanERMDeptService struct {
	repo              riskRepositories.RiskUpdatePlanERMDeptRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
}

func NewRiskUpdatePlanERMDeptService(
	repo riskRepositories.RiskUpdatePlanERMDeptRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) RiskUpdatePlanERMDeptServiceInterface {
	return &RiskUpdatePlanERMDeptService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
	}
}

func (rupeds *RiskUpdatePlanERMDeptService) CreateERMDepartmentEvaluation(planID uint, currentUserID uint) error {
	plan, err := rupeds.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		return err
	}

	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID (pre-validation)
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// All pre-validation passed, start transaction
	tx := rupeds.repo.BeginTransaction()

	// Fetch risks, processes, and controls for the business unit
	riskIDs, err := rupeds.repo.GetRisksByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	processIDs, err := rupeds.repo.GetProcessesByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	controlIDs, err := rupeds.repo.GetControlsByBusinessUnitID(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate that ERM department evaluation data exists
	if len(riskIDs) == 0 && len(processIDs) == 0 && len(controlIDs) == 0 {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("no ERM department evaluation data found for this business unit"), "not_exist", "ERM department evaluation data")
	}

	// Create risk mappings
	riskMappings := make([]*models.RiskUpdatePlanRisk, 0, len(riskIDs))
	for _, riskID := range riskIDs {
		riskMappings = append(riskMappings, &models.RiskUpdatePlanRisk{
			RiskUpdatePlanTimelineID: timeline.ID,
			RiskRegisterID:           riskID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := rupeds.repo.CreateRiskMappings(riskMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create process mappings
	processMappings := make([]*models.RiskUpdatePlanProcess, 0, len(processIDs))
	for _, processID := range processIDs {
		processMappings = append(processMappings, &models.RiskUpdatePlanProcess{
			RiskUpdatePlanTimelineID: timeline.ID,
			ProcessID:                processID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := rupeds.repo.CreateProcessMappings(processMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create control mappings
	controlMappings := make([]*models.RiskUpdatePlanControl, 0, len(controlIDs))
	for _, controlID := range controlIDs {
		controlMappings = append(controlMappings, &models.RiskUpdatePlanControl{
			RiskUpdatePlanTimelineID: timeline.ID,
			ControlID:                controlID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := rupeds.repo.CreateControlMappings(controlMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Mark timeline as completed
	if err := rupeds.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rupeds.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ============================================
// ERM Department Evaluation - Risks
// ============================================

func (rupeds *RiskUpdatePlanERMDeptService) MapRisksToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapRisksToUpdatePlanDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID (pre-validation)
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// All pre-validation passed, start transaction
	tx := rupeds.repo.BeginTransaction()

	// Delete all existing risk mappings for this timeline
	if err := rupeds.repo.DeleteAllRiskMappingsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new risk mappings
	riskMappings := make([]*models.RiskUpdatePlanRisk, 0, len(formDTO.RiskRegisterIDs))
	for _, riskID := range formDTO.RiskRegisterIDs {
		riskMappings = append(riskMappings, &models.RiskUpdatePlanRisk{
			RiskUpdatePlanTimelineID: timeline.ID,
			RiskRegisterID:           riskID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := rupeds.repo.CreateRiskMappings(riskMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rupeds *RiskUpdatePlanERMDeptService) GetERMDeptEvaluationRisks(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationListResponseDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return rupeds.repo.GetERMDeptEvaluationRisks(timeline.ID, params)
}

// ============================================
// ERM Department Evaluation - Processes
// ============================================

func (rupeds *RiskUpdatePlanERMDeptService) MapProcessesToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapProcessesToUpdatePlanDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID (pre-validation)
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// All pre-validation passed, start transaction
	tx := rupeds.repo.BeginTransaction()

	// Delete all existing process mappings for this timeline
	if err := rupeds.repo.DeleteAllProcessMappingsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new process mappings
	processMappings := make([]*models.RiskUpdatePlanProcess, 0, len(formDTO.ProcessIDs))
	for _, processID := range formDTO.ProcessIDs {
		processMappings = append(processMappings, &models.RiskUpdatePlanProcess{
			RiskUpdatePlanTimelineID: timeline.ID,
			ProcessID:                processID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := rupeds.repo.CreateProcessMappings(processMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rupeds *RiskUpdatePlanERMDeptService) GetERMDeptEvaluationProcesses(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationProcessListResponseDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return rupeds.repo.GetERMDeptEvaluationProcesses(timeline.ID, params)
}

// ============================================
// ERM Department Evaluation - Controls
// ============================================

func (rupeds *RiskUpdatePlanERMDeptService) MapControlsToTimeline(planID uint, currentUserID uint, formDTO *riskDTOs.MapControlsToUpdatePlanDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID (pre-validation)
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	// All pre-validation passed, start transaction
	tx := rupeds.repo.BeginTransaction()

	// Delete all existing control mappings for this timeline
	if err := rupeds.repo.DeleteAllControlMappingsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new control mappings
	controlMappings := make([]*models.RiskUpdatePlanControl, 0, len(formDTO.ControlIDs))
	for _, controlID := range formDTO.ControlIDs {
		controlMappings = append(controlMappings, &models.RiskUpdatePlanControl{
			RiskUpdatePlanTimelineID: timeline.ID,
			ControlID:                controlID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := rupeds.repo.CreateControlMappings(controlMappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rupeds *RiskUpdatePlanERMDeptService) GetERMDeptEvaluationControls(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.ERMDeptEvaluationControlListResponseDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return rupeds.repo.GetERMDeptEvaluationControls(timeline.ID, params)
}

// ============================================
// ERM Department Evaluation Operations
// ============================================

func (rupeds *RiskUpdatePlanERMDeptService) MarkERMDeptEvaluationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupeds.repo.BeginTransaction()

	// Update timeline - mark as not applicable and completed
	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rupeds.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// ERM Department Evaluation Document Operations

func (rupeds *RiskUpdatePlanERMDeptService) GetERMDeptEvaluationDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents
	documents, err := rupeds.repo.GetERMDeptDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

	// Convert to DTOs
	var documentDTOs []*riskDTOs.CommonDocumentDTO
	for _, doc := range documents {
		dto := &riskDTOs.CommonDocumentDTO{
			ID:           &doc.ID,
			URL:          doc.URL,
			Name:         doc.Name,
			Token:        doc.Token,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			Size:         doc.Size,
		}
		documentDTOs = append(documentDTOs, dto)
	}

	return documentDTOs, nil
}

func (rupeds *RiskUpdatePlanERMDeptService) CreateERMDeptEvaluationDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {

	// Verify plan exists
	_, err := rupeds.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupeds.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanErmDeptDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanErmDeptDocument

	now := time.Now()

	for _, doc := range documents {

		// NEW FILE
		if doc.IsNewFile != nil && *doc.IsNewFile {

			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"erm-dept-evaluation-documents",
			)

			document := &models.RiskUpdatePlanErmDeptDocument{
				RiskUpdatePlanTimelineID: timeline.ID,
				URL:                      &url,
				Name:                     &doc.FileName,
				Size:                     &doc.Size,
				Token:                    &encryptedFilename,
				OriginalName:             &doc.OriginalName,
				Type:                     &doc.Type,
				CreatedBYID:              &currentUserID,
				CreatedAt:                now,
			}

			newDocuments = append(newDocuments, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {

			// DELETE FILE
			existingDocument, err := rupeds.repo.GetERMDeptDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Verify document belongs to this timeline
			if existingDocument.RiskUpdatePlanTimelineID != timeline.ID {
				tx.Rollback()
				return utility.NewResourceNotFoundError(fmt.Errorf("document not found"), *doc.ID)
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	// Create new documents
	if len(newDocuments) > 0 {
		for _, doc := range newDocuments {
			_, err := rupeds.repo.CreateERMDeptDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := rupeds.repo.DeleteERMDeptDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"erm-dept-evaluation-documents",
		)

		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"erm-dept-evaluation-documents",
		)

		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	// Log activity
	go func() {
		err := rupeds.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanErmDeptDocument{}.TableName(),
			planID,
			map[string]interface{}{
				"new_documents_count":     len(newDocuments),
				"deleted_documents_count": len(documentsToDelete),
			},
			currentUserID,
			"risk", "risk_register_update_plan", "documents_update", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action": models.UpdateAction,
				"planID": planID,
				"userID": currentUserID,
				"error":  err.Error(),
			})
		}
	}()

	return nil
}

func (rupeds *RiskUpdatePlanERMDeptService) GetERMDeptEvaluationDocumentPath(planID uint, token string) (string, error) {
	// Get timeline type for ERM Department Evaluation
	timelineType, err := rupeds.timelineTypeSvc.GetByType("erm_department_evaluation")
	if err != nil {
		return "", err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupeds.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	// Get all documents for this timeline
	documents, err := rupeds.repo.GetERMDeptDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	// Find document with matching token
	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			// Decrypt filename
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			// Construct file path
			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"erm-dept-evaluation-documents",
			)
			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			// Check if the file exists
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(errors.New("document not found"), 0)
}
