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

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskUpdatePlanRiskTeamAssessmentServiceInterface interface {
	GetRiskTeamAssessmentUser(planID uint) (*riskDTOs.RiskTeamAssessmentUserListDTO, error)
	CreateRiskTeamAssessment(planID uint, currentUserID uint) error
	MarkRiskTeamAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
	GetRiskTeamAssessmentDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error)
	CreateRiskTeamAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error
	GetRiskTeamAssessmentDocumentPath(planID uint, token string) (string, error)
}

type RiskUpdatePlanRiskTeamAssessmentService struct {
	repo              riskRepositories.RiskUpdatePlanRiskTeamAssessmentRepositoryInterface
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc   RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc   RiskRegisterServiceInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	userSvc           organizationServices.UserServiceInterface
}

func NewRiskUpdatePlanRiskTeamAssessmentService(
	repo riskRepositories.RiskUpdatePlanRiskTeamAssessmentRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) RiskUpdatePlanRiskTeamAssessmentServiceInterface {
	return &RiskUpdatePlanRiskTeamAssessmentService{
		repo:              repo,
		riskUpdatePlanSvc: riskUpdatePlanSvc,
		timelineTypeSvc:   timelineTypeSvc,
		riskRegisterSvc:   riskRegisterSvc,
		activityLogSvc:    activityLogSvc,
		userSvc:           userSvc,
	}
}

func (ruprtas *RiskUpdatePlanRiskTeamAssessmentService) GetRiskTeamAssessmentUser(planID uint) (*riskDTOs.RiskTeamAssessmentUserListDTO, error) {

	timelineType, err := ruprtas.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return nil, err
	}

	timeline, err := ruprtas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	assessments, err := ruprtas.repo.GetRiskTeamAssessmentUsers(timeline.ID)
	if err != nil {
		return nil, err
	}

	if len(assessments) == 0 {
		return nil, nil
	}

	var dto riskDTOs.RiskTeamAssessmentUserListDTO

	err = mapstructure.Decode(assessments[0], &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (ruprtas *RiskUpdatePlanRiskTeamAssessmentService) CreateRiskTeamAssessment(planID uint, currentUserID uint) error {
	plan, err := ruprtas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		return err
	}

	timelineType, err := ruprtas.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return err
	}

	tx := ruprtas.repo.BeginTransaction()

	timeline, err := ruprtas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate that risk champion exists
	if plan.RiskChampionID == 0 {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("no risk champion assigned for this plan"), "not_exist", "risk champion")
	}

	assessments := make([]*models.RiskUpdatePlanRiskTeamAssessment, 0, 1)
	assessments = append(assessments, &models.RiskUpdatePlanRiskTeamAssessment{
		RiskUpdatePlanTimelineID: timeline.ID,
		RiskChampionID:           plan.RiskChampionID,
		CreatedBYID:              &currentUserID,
		CreatedAt:                time.Now(),
	})

	if err := ruprtas.repo.CreateRiskTeamAssessments(assessments, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ruprtas.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ruprtas.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ruprtas *RiskUpdatePlanRiskTeamAssessmentService) MarkRiskTeamAssessmentAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := ruprtas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := ruprtas.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return err
	}

	timeline, err := ruprtas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := ruprtas.repo.BeginTransaction()

	if err := ruprtas.repo.DeleteRiskTeamAssessmentsByTimelineID(timeline.ID, tx); err != nil {
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

	if err := ruprtas.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (ruprtas *RiskUpdatePlanRiskTeamAssessmentService) GetRiskTeamAssessmentDocuments(planID uint) ([]*riskDTOs.CommonDocumentDTO, error) {
	timelineType, err := ruprtas.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return nil, err
	}

	timeline, err := ruprtas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	documents, err := ruprtas.repo.GetRiskTeamDocumentsByTimelineID(timeline.ID)
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

func (ruprtas *RiskUpdatePlanRiskTeamAssessmentService) CreateRiskTeamAssessmentDocuments(planID uint, currentUserID uint, documents []*commonDTOs.DocumentFormDTO) error {
	_, err := ruprtas.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := ruprtas.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return err
	}

	timeline, err := ruprtas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := ruprtas.repo.BeginTransaction()

	var newDocuments []*models.RiskUpdatePlanRiskTeamDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.RiskUpdatePlanRiskTeamDocument

	now := time.Now()

	for _, doc := range documents {

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
				"risk-team-assessment-documents",
			)

			document := &models.RiskUpdatePlanRiskTeamDocument{
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

			existingDocument, err := ruprtas.repo.GetRiskTeamDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			if existingDocument.RiskUpdatePlanTimelineID != timeline.ID {
				tx.Rollback()
				return utility.NewResourceNotFoundError(fmt.Errorf("document not found"), *doc.ID)
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(newDocuments) > 0 {
		for _, doc := range newDocuments {
			_, err := ruprtas.repo.CreateRiskTeamDocument(doc, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(deleteDocumentIDs) > 0 {
		for _, docID := range deleteDocumentIDs {
			if err := ruprtas.repo.DeleteRiskTeamDocument(*docID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"risk-team-assessment-documents",
		)

		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join(
			"risk",
			"risk-register-update-plan",
			strconv.FormatUint(uint64(planID), 10),
			"risk-team-assessment-documents",
		)

		if doc.Name != nil {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	return nil
}

func (ruprtas *RiskUpdatePlanRiskTeamAssessmentService) GetRiskTeamAssessmentDocumentPath(planID uint, token string) (string, error) {
	timelineType, err := ruprtas.timelineTypeSvc.GetByType("risk_team_assessment")
	if err != nil {
		return "", err
	}

	timeline, err := ruprtas.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return "", err
	}

	documents, err := ruprtas.repo.GetRiskTeamDocumentsByTimelineID(timeline.ID)
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
				"risk-team-assessment-documents",
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
