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

type RiskUpdatePlanRiskAssessmentWorkshopServiceInterface interface {
	GetRiskAssessmentWorkshopDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error)
	CreateRiskAssessmentWorkshopDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetRiskAssessmentWorkshopDocumentPath(planID uint, token string) (string, error)
	// Risk Assessment Workshop - Operations
	MarkRiskAssessmentWorkshopAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanRiskAssessmentWorkshopService struct {
	repo              riskRepositories.RiskUpdatePlanRiskAssessmentWorkshopRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
}

func NewRiskUpdatePlanRiskAssessmentWorkshopService(
	repo riskRepositories.RiskUpdatePlanRiskAssessmentWorkshopRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) RiskUpdatePlanRiskAssessmentWorkshopServiceInterface {
	return &RiskUpdatePlanRiskAssessmentWorkshopService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
	}
}

func (rupraws *RiskUpdatePlanRiskAssessmentWorkshopService) GetRiskAssessmentWorkshopDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error) {
	// Get timeline type for Risk Assessment Workshop
	timelineType, err := rupraws.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupraws.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := rupraws.repo.GetRiskAssessmentWorkshopDocumentsByTimelineID(timeline.ID)
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

func (rupraws *RiskUpdatePlanRiskAssessmentWorkshopService) CreateRiskAssessmentWorkshopDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := rupraws.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Risk Assessment Workshop
	timelineType, err := rupraws.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupraws.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupraws.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanRiskAssessmentWorkshopDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanRiskAssessmentWorkshopDocument

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
				"risk-assessment-workshop-documents",
			)

			document := &models.RiskUpdatePlanRiskAssessmentWorkshopDocument{
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
			existingDocument, err := rupraws.repo.GetRiskAssessmentWorkshopDocumentByID(*doc.ID)
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
			_, err := rupraws.repo.CreateRiskAssessmentWorkshopDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := rupraws.repo.DeleteRiskAssessmentWorkshopDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-assessment-workshop-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "risk-assessment-workshop-documents")
		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Mark timeline as completed
	if err := rupraws.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rupraws.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := rupraws.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanRiskAssessmentWorkshopDocument{}.TableName(),
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

func (rupraws *RiskUpdatePlanRiskAssessmentWorkshopService) GetRiskAssessmentWorkshopDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := rupraws.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return "", err
	}

	timeline, err := rupraws.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := rupraws.repo.GetRiskAssessmentWorkshopDocumentsByTimelineID(timeline.ID)
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
				"risk-assessment-workshop-documents",
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

func (rupraws *RiskUpdatePlanRiskAssessmentWorkshopService) MarkRiskAssessmentWorkshopAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rupraws.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rupraws.timelineTypeSvc.GetByType("risk_assessment_workshop")
	if err != nil {
		return err
	}

	timeline, err := rupraws.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupraws.repo.BeginTransaction()

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

	if err := rupraws.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
