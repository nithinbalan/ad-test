package services

import (
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type RiskUpdatePlanLeadershipAssessmentServiceInterface interface {
	GetLeadershipAssessmentUsers(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.LeaderShipAssessmentUserListResponseDTO, error)
	CreateLeadershipAssessment(planID uint, currentUserID uint) error
	MarkLeadershipAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
	GetLeadershipAssessmentDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error)
	CreateLeadershipAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetLeadershipAssessmentDocumentPath(planID uint, token string) (string, error)
}

type RiskUpdatePlanLeadershipAssessmentService struct {
	repo              riskRepositories.RiskUpdatePlanLeadershipAssessmentRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	userSvc           organizationServices.UserServiceInterface
}

func NewRiskUpdatePlanLeadershipAssessmentService(
	repo riskRepositories.RiskUpdatePlanLeadershipAssessmentRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) RiskUpdatePlanLeadershipAssessmentServiceInterface {
	return &RiskUpdatePlanLeadershipAssessmentService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
		userSvc:           userSvc,
	}
}

func (ruplas *RiskUpdatePlanLeadershipAssessmentService) GetLeadershipAssessmentUsers(planID uint, params *riskDTOs.CommonListQueryParamsDTO) (*riskDTOs.LeaderShipAssessmentUserListResponseDTO, error) {
	// Get timeline type for Leadership Assessment
	timelineType, err := ruplas.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruplas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	return ruplas.repo.GetLeadershipAssessmentUsers(timeline.ID, params)
}

func (ruplas *RiskUpdatePlanLeadershipAssessmentService) CreateLeadershipAssessment(planID uint, currentUserID uint) error {
	plan, err := ruplas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		return err
	}

	// Get timeline type for Leadership Assessment
	timelineType, err := ruplas.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return err
	}

	tx := ruplas.repo.BeginTransaction()

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruplas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Fetch users for the business unit
	users, err := ruplas.userSvc.GetUsersByOrganizationUnit(plan.BusinessUnitID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate that users exist for the business unit
	if len(users) == 0 {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("no users found for this business unit"), "not_exist", "leadership assessment users")
	}

	// Create leadership assessment mappings
	assessments := make([]*models.RiskUpdatePlanLeadershipAssessment, 0, len(users))
	for _, user := range users {
		assessments = append(assessments, &models.RiskUpdatePlanLeadershipAssessment{
			RiskUpdatePlanTimelineID: timeline.ID,
			UserID:                   user.ID,
			CreatedBYID:              &currentUserID,
			CreatedAt:                time.Now(),
		})
	}

	if err := ruplas.repo.CreateLeadershipAssessments(assessments, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Mark timeline as completed
	if err := ruplas.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Recalculate percentage
	if err := ruplas.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ruplas *RiskUpdatePlanLeadershipAssessmentService) MarkLeadershipAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := ruplas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := ruplas.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return err
	}

	timeline, err := ruplas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := ruplas.repo.BeginTransaction()

	if err := ruplas.repo.DeleteLeadershipAssessmentsByTimelineID(timeline.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

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

	if err := ruplas.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ruplas *RiskUpdatePlanLeadershipAssessmentService) GetLeadershipAssessmentDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error) {
	timelineType, err := ruplas.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return nil, err
	}

	timeline, err := ruplas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	documents, err := ruplas.repo.GetLeadershipDocumentsByTimelineID(timeline.ID)
	if err != nil {
		return nil, err
	}

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

func (ruplas *RiskUpdatePlanLeadershipAssessmentService) CreateLeadershipAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	// Verify plan exists
	_, err := ruplas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	// Get timeline type for Leadership Assessment
	timelineType, err := ruplas.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := ruplas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := ruplas.repo.BeginTransaction()

	// Handle documents
	var newDocuments []*models.RiskUpdatePlanLeadershipDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanLeadershipDocument

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
				"leadership-assessment-documents",
			)

			document := &models.RiskUpdatePlanLeadershipDocument{
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
			existingDocument, err := ruplas.repo.GetLeadershipDocumentByID(*doc.ID)
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
			_, err := ruplas.repo.CreateLeadershipDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// Delete documents
	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := ruplas.repo.DeleteLeadershipDocument(*docID, tx); err != nil {
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
			"leadership-assessment-documents",
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
			"leadership-assessment-documents",
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
		err := ruplas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskUpdatePlanLeadershipDocument{}.TableName(),
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

func (ruplas *RiskUpdatePlanLeadershipAssessmentService) GetLeadershipAssessmentDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := ruplas.timelineTypeSvc.GetByType("leadership_assessment")
	if err != nil {
		return "", err
	}

	timeline, err := ruplas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := ruplas.repo.GetLeadershipDocumentsByTimelineID(timeline.ID)
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
				"leadership-assessment-documents",
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
