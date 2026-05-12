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

type RiskUpdatePlanRiskWorkshopAssessmentServiceInterface interface {
	GetRiskWorkshopAssessmentDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error)
	CreateRiskWorkshopAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetRiskWorkshopAssessmentDocumentPath(planID uint, token string) (string, error)
	// Risk Workshop Assessment - Operations
	MarkRiskWorkshopAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanRiskWorkshopAssessmentService struct {
	repo              riskRepositories.RiskUpdatePlanRiskWorkshopAssessmentRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
}

func NewRiskUpdatePlanRiskWorkshopAssessmentService(
	repo riskRepositories.RiskUpdatePlanRiskWorkshopAssessmentRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) RiskUpdatePlanRiskWorkshopAssessmentServiceInterface {
	return &RiskUpdatePlanRiskWorkshopAssessmentService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
	}
}

func (ruprwas *RiskUpdatePlanRiskWorkshopAssessmentService) GetRiskWorkshopAssessmentDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error) {
	// Get timeline type for Risk Workshop Assessment
	timelineType, err := ruprwas.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruprwas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := ruprwas.repo.GetRiskWorkshopAssessmentDocumentsByTimelineID(timeline.ID)
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

func (ruprwas *RiskUpdatePlanRiskWorkshopAssessmentService) CreateRiskWorkshopAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := ruprwas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Workshop Assessment
	timelineType, err := ruprwas.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruprwas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := ruprwas.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanRiskWorkshopAssessmentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanRiskWorkshopAssessmentDocument

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
				"risk-workshop-assessment-documents",
			)

			document := &models.RiskUpdatePlanRiskWorkshopAssessmentDocument{
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
			existingDocument, err := ruprwas.repo.GetRiskWorkshopAssessmentDocumentByID(*doc.ID)
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
			_, err := ruprwas.repo.CreateRiskWorkshopAssessmentDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := ruprwas.repo.DeleteRiskWorkshopAssessmentDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-workshop-assessment-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-workshop-assessment-documents")
		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Mark timeline as completed
	if err := ruprwas.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := ruprwas.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := ruprwas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanRiskWorkshopAssessmentDocument{}.TableName(),
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
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (ruprwas *RiskUpdatePlanRiskWorkshopAssessmentService) GetRiskWorkshopAssessmentDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := ruprwas.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return "", err
	}

	timeline, err := ruprwas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := ruprwas.repo.GetRiskWorkshopAssessmentDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return "", err
	}

	for _, doc := range documents {
		if doc.Token != nil && *doc.Token == token {
			decryptedFilename, err := utility.DecryptData(*doc.Token)
			if err != nil {
				return "", err
			}

			storagePath := filepath.Join(
				"risk",
				"risk-register-update-plan",
				strconv.FormatUint(uint64(planID), 10),
				"risk-workshop-assessment-documents",
			)

			baseStoragePath := config.GetEnv("STORAGE_PATH", "")
			fullPath := filepath.Join(baseStoragePath, storagePath, string(decryptedFilename))

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return "", utility.NewResourceNotFoundError(err, 0)
			}

			return fullPath, nil
		}
	}

	return "", utility.NewResourceNotFoundError(fmt.Errorf("document not found"), 0)
}

func (ruprwas *RiskUpdatePlanRiskWorkshopAssessmentService) MarkRiskWorkshopAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := ruprwas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := ruprwas.timelineTypeSvc.GetByType("risk_workshop_assessment")
	if err != nil {
		return err
	}

	timeline, err := ruprwas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := ruprwas.repo.BeginTransaction()

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

	if err := ruprwas.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
