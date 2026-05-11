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

type RiskUpdatePlanIntroductoryMeetingServiceInterface interface {
	// Introductory Meeting Document Operations
	GetIntroductoryMeetingDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error)
	CreateIntroductoryMeetingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetIntroductoryMeetingDocumentPath(planID uint, token string) (string, error)
	// Introductory Meeting - Operations
	MarkIntroductoryMeetingAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanIntroductoryMeetingService struct {
	repo              riskRepositories.RiskUpdatePlanIntroductoryMeetingRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
}

func NewRiskUpdatePlanIntroductoryMeetingService(
	repo riskRepositories.RiskUpdatePlanIntroductoryMeetingRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) RiskUpdatePlanIntroductoryMeetingServiceInterface {
	return &RiskUpdatePlanIntroductoryMeetingService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
	}
}

// Introductory Meeting Document Operations

func (rupims *RiskUpdatePlanIntroductoryMeetingService) GetIntroductoryMeetingDocumentPath(planID uint, token string) (string, error) {
	// Get timeline type for Introductory Meeting
	timelineType, err := rupims.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return "", err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupims.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	// Get all documents for this timeline
	documents, err := rupims.repo.GetIntroductoryMeetingDocumentsByTimelineID(timeline.ID)
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
			storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "introductory-meeting-documents")
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

// Introductory Meeting Document Operations

func (rupims *RiskUpdatePlanIntroductoryMeetingService) GetIntroductoryMeetingDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error) {
	// Get timeline type for Introductory Meeting
	timelineType, err := rupims.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupims.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	// Get documents by timeline ID
	documents, err := rupims.repo.GetIntroductoryMeetingDocumentsByTimelineID(timeline.ID)
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

func (rupims *RiskUpdatePlanIntroductoryMeetingService) CreateIntroductoryMeetingDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := rupims.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Introductory Meeting
	timelineType, err := rupims.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupims.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupims.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanIntroductoryMeetingDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanIntroductoryMeetingDocument

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
				"introductory-meeting-documents",
			)

			document := &models.RiskUpdatePlanIntroductoryMeetingDocument{
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
			existingDocument, err := rupims.repo.GetIntroductoryMeetingDocumentByID(*doc.ID)
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
			_, err := rupims.repo.CreateIntroductoryMeetingDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := rupims.repo.DeleteIntroductoryMeetingDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// File system operations - store new files
	for _, doc := range newDocuments {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "introductory-meeting-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// File system operations - delete files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register-update-plan", strconv.FormatUint(uint64(planID), 10), "introductory-meeting-documents")
		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Mark timeline as completed
	if err := rupims.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := rupims.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity
	go func() {
		err := rupims.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanIntroductoryMeetingDocument{}.TableName(),
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

func (rupims *RiskUpdatePlanIntroductoryMeetingService) MarkIntroductoryMeetingAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	// Get timeline type for Introductory Meeting
	timelineType, err := rupims.timelineTypeSvc.GetByType("introductory_meeting")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rupims.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rupims.repo.BeginTransaction()

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
	if err := rupims.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
