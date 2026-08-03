package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IdentityReviewCampaignServiceInterface interface {
	GetAll(params *dtos.IdentityReviewCampaignListQueryParamsDTO) (*dtos.IdentityReviewCampaignListResponseDTO, error)
	GetByID(id uint) (*models.IdentityReviewCampaign, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewCampaignDetailDTO, error)
	Create(currentUserID uint, formDTO *dtos.IdentityReviewCampaignFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewCampaignFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(id uint, token string) (string, error)
}

type IdentityReviewCampaignService struct {
	repo                    repositories.IdentityReviewCampaignRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	reviewTypeSvc           ReviewTypeServiceInterface
	reviewCampaignStatusSvc ReviewCampaignStatusServiceInterface
	userSvc                 organizationServices.UserServiceInterface
}

func NewIdentityReviewCampaignService(
	repo repositories.IdentityReviewCampaignRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	reviewTypeSvc ReviewTypeServiceInterface,
	reviewCampaignStatusSvc ReviewCampaignStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) IdentityReviewCampaignServiceInterface {
	return &IdentityReviewCampaignService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		reviewTypeSvc:           reviewTypeSvc,
		reviewCampaignStatusSvc: reviewCampaignStatusSvc,
		userSvc:                 userSvc,
	}
}

func (ircs *IdentityReviewCampaignService) GetAll(params *dtos.IdentityReviewCampaignListQueryParamsDTO) (*dtos.IdentityReviewCampaignListResponseDTO, error) {
	return ircs.repo.GetAll(params)
}

func (ircs *IdentityReviewCampaignService) GetByID(id uint) (*models.IdentityReviewCampaign, error) {
	campaign, err := ircs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return campaign, nil
}

func (ircs *IdentityReviewCampaignService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewCampaignDetailDTO, error) {
	campaign, err := ircs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var campaignDTO dtos.IdentityReviewCampaignDetailDTO
	if err := mapstructure.Decode(campaign, &campaignDTO); err != nil {
		return nil, err
	}

	campaignDTO.StartDate = campaign.StartDate.Format("2006-01-02")
	campaignDTO.DueDate = campaign.DueDate.Format("2006-01-02")
	campaignDTO.CreatedAt = campaign.CreatedAt
	campaignDTO.UpdatedAt = campaign.UpdatedAt

	if campaign.CreatedByUser != nil {
		if err := mapstructure.Decode(campaign.CreatedByUser, &campaignDTO.CreatedBy); err != nil {
			return nil, err
		}
	}

	if campaign.UpdatedByUser != nil {
		campaignDTO.UpdatedBy = &dtos.IdentityReviewCampaignUserDTO{}
		if err := mapstructure.Decode(campaign.UpdatedByUser, campaignDTO.UpdatedBy); err != nil {
			return nil, err
		}
	}

	return &campaignDTO, nil
}

func (ircs *IdentityReviewCampaignService) Create(currentUserID uint, formDTO *dtos.IdentityReviewCampaignFormDTO) (uint, error) {
	if len(formDTO.SystemRegisterIDs) == 0 {
		return 0, utility.NewValidationError(
			fmt.Errorf("at least one system register is required"),
			"required",
			"system_register_ids",
		)
	}

	tx := ircs.repo.BeginTransaction()

	if err := ircs.validateIdentityReviewCampaignFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	plannedStatus, err := ircs.reviewCampaignStatusSvc.GetByType("planned")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("failed to get planned review campaign status: %w", err)
	}

	dates, err := ircs.parseIdentityReviewCampaignDates(formDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	referenceNumber, err := ircs.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	campaign := &models.IdentityReviewCampaign{
		ReferenceNumber:    referenceNumber,
		CampaignName:       formDTO.CampaignName,
		ReviewTypeID:       formDTO.ReviewTypeID,
		Description:        formDTO.Description,
		StartDate:          *dates.StartDate,
		DueDate:            *dates.DueDate,
		ReviewerID:         formDTO.ReviewerID,
		ApproverID:         formDTO.ApproverID,
		CampaignStatusID:   plannedStatus.ID,
		ReminderDaysOffset: formDTO.ReminderDaysOffset,
		CreatedByID:        currentUserID,
	}

	id, err := ircs.repo.Create(campaign, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var systems []*models.IdentityReviewCampaignSystem
	for _, systemRegisterID := range formDTO.SystemRegisterIDs {
		systems = append(systems, &models.IdentityReviewCampaignSystem{
			IdentityReviewCampaignID: id,
			SystemRegisterID:         systemRegisterID,
			CreatedByID:              currentUserID,
		})
	}

	if err := ircs.repo.CreateSystems(systems, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IdentityReviewCampaignDocument
	for _, documentDTO := range formDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(documentDTO.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		documents = append(documents, &models.IdentityReviewCampaignDocument{
			IdentityReviewCampaignID: id,
			URL: utility.GenerateDocumentURL(
				"identity-access-review",
				"review-campaign",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			),
			Name:         documentDTO.FileName,
			Size:         documentDTO.Size,
			Token:        encryptedFilename,
			OriginalName: documentDTO.OriginalName,
			Type:         documentDTO.Type,
			CreatedByID:  currentUserID,
		})
	}

	if err := ircs.repo.CreateDocuments(documents, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, document := range documents {
		storagePath := filepath.Join(
			"identity-access-review",
			"review-campaign",
			strconv.FormatUint(uint64(id), 10),
			"documents",
		)
		if err := utility.StoreFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingCampaign, err := ircs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch identity review campaign", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		if err := ircs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewCampaign{}.TableName(),
			id,
			existingCampaign,
			currentUserID,
			"identity_access_review",
			"review_campaign",
			"campaign",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (ircs *IdentityReviewCampaignService) Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewCampaignFormDTO) error {
	tx := ircs.repo.BeginTransaction()

	campaign, err := ircs.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	oldCampaign, err := ircs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := ircs.validateIdentityReviewCampaignFormDTO(formDTO); err != nil {
		tx.Rollback()
		return err
	}

	dates, err := ircs.parseIdentityReviewCampaignDates(formDTO)
	if err != nil {
		tx.Rollback()
		return err
	}

	campaign.CampaignName = formDTO.CampaignName
	campaign.ReviewTypeID = formDTO.ReviewTypeID
	campaign.Description = formDTO.Description
	campaign.StartDate = *dates.StartDate
	campaign.DueDate = *dates.DueDate
	campaign.ReviewerID = formDTO.ReviewerID
	campaign.ApproverID = formDTO.ApproverID
	campaign.ReminderDaysOffset = formDTO.ReminderDaysOffset
	campaign.UpdatedByID = &currentUserID

	if err := ircs.repo.Update(campaign, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ircs.repo.DeleteSystemsByCampaignID(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	var systems []*models.IdentityReviewCampaignSystem
	for _, systemRegisterID := range formDTO.SystemRegisterIDs {
		systems = append(systems, &models.IdentityReviewCampaignSystem{
			IdentityReviewCampaignID: id,
			SystemRegisterID:         systemRegisterID,
			CreatedByID:              currentUserID,
		})
	}

	if err := ircs.repo.CreateSystems(systems, tx); err != nil {
		tx.Rollback()
		return err
	}

	var documentsToCreate []*models.IdentityReviewCampaignDocument
	var documentIDsToDelete []uint
	var documentsToDelete []*models.IdentityReviewCampaignDocument

	for _, documentDTO := range formDTO.Documents {
		if documentDTO.IsNewFile != nil && *documentDTO.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(documentDTO.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			documentsToCreate = append(documentsToCreate, &models.IdentityReviewCampaignDocument{
				IdentityReviewCampaignID: id,
				URL: utility.GenerateDocumentURL(
					"identity-access-review",
					"review-campaign",
					strconv.FormatUint(uint64(id), 10),
					"documents",
				),
				Name:         documentDTO.FileName,
				Size:         documentDTO.Size,
				Token:        encryptedFilename,
				OriginalName: documentDTO.OriginalName,
				Type:         documentDTO.Type,
				CreatedByID:  currentUserID,
			})
		}

		if documentDTO.IsDeletedFile != nil && *documentDTO.IsDeletedFile && documentDTO.ID != nil {
			existingDocument, err := ircs.repo.GetDocumentByID(*documentDTO.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if existingDocument.IdentityReviewCampaignID != id {
				tx.Rollback()
				return utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, *documentDTO.ID)
			}

			documentIDsToDelete = append(documentIDsToDelete, existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if err := ircs.repo.CreateDocuments(documentsToCreate, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := ircs.repo.DeleteDocuments(documentIDsToDelete, tx); err != nil {
		tx.Rollback()
		return err
	}

	storagePath := filepath.Join(
		"identity-access-review",
		"review-campaign",
		strconv.FormatUint(uint64(id), 10),
		"documents",
	)

	for _, document := range documentsToCreate {
		if err := utility.StoreFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		updatedCampaign, err := ircs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch identity review campaign", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		changedData := map[string]interface{}{
			"old": oldCampaign,
			"new": updatedCampaign,
		}
		if err := ircs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IdentityReviewCampaign{}.TableName(),
			id,
			changedData,
			currentUserID,
			"identity_access_review",
			"review_campaign",
			"campaign",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (ircs *IdentityReviewCampaignService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ircs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	var deletedCampaigns []interface{}
	documentsToDelete := make(map[uint][]*models.IdentityReviewCampaignDocument)

	for _, id := range ids {
		campaign, err := ircs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedCampaigns = append(deletedCampaigns, campaign)

		var documentIDs []uint
		for _, document := range campaign.Documents {
			documentIDs = append(documentIDs, document.ID)
			documentsToDelete[id] = append(documentsToDelete[id], document)
		}

		if err := ircs.repo.DeleteDocuments(documentIDs, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := ircs.repo.DeleteSystemsByCampaignID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := ircs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for campaignID, documents := range documentsToDelete {
		storagePath := filepath.Join(
			"identity-access-review",
			"review-campaign",
			strconv.FormatUint(uint64(campaignID), 10),
			"documents",
		)
		for _, document := range documents {
			if err := utility.RemoveFile(document.Name, storagePath); err != nil && !os.IsNotExist(err) {
				utility.LogError("Failed to delete file", map[string]interface{}{
					"error": err.Error(),
					"name":  document.Name,
				})
			}
		}
	}

	if len(successfulIDs) > 0 {
		go func() {
			if err := ircs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IdentityReviewCampaign{}.TableName(),
				successfulIDs,
				deletedCampaigns,
				currentUserID,
				nil,
			); err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	return failedIDs, nil
}

func (ircs *IdentityReviewCampaignService) validateIdentityReviewCampaignFormDTO(formDTO *dtos.IdentityReviewCampaignFormDTO) error {
	if _, err := ircs.reviewTypeSvc.GetByID(formDTO.ReviewTypeID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("review type does not exist"),
			"not_exist",
			"review_type_id",
		)
	}

	if _, err := ircs.userSvc.GetByID(formDTO.ReviewerID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("reviewer does not exist"),
			"not_exist",
			"reviewer_id",
		)
	}

	if _, err := ircs.userSvc.GetByID(formDTO.ApproverID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("approver does not exist"),
			"not_exist",
			"approver_id",
		)
	}

	if formDTO.ReminderDaysOffset != nil && *formDTO.ReminderDaysOffset < 0 {
		return utility.NewValidationError(
			fmt.Errorf("reminder days offset cannot be negative"),
			"min",
			"reminder_days_offset",
		)
	}

	return nil
}

func (ircs *IdentityReviewCampaignService) parseIdentityReviewCampaignDates(formDTO *dtos.IdentityReviewCampaignFormDTO) (*dtos.IdentityReviewCampaignParsedDates, error) {
	parsedDates := &dtos.IdentityReviewCampaignParsedDates{}

	startDate, err := utility.ParseDate(formDTO.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(
			fmt.Errorf("invalid start date"),
			"invalid_date",
			"start_date",
		)
	}
	parsedDates.StartDate = startDate

	dueDate, err := utility.ParseDate(formDTO.DueDate)
	if err != nil {
		return nil, utility.NewValidationError(
			fmt.Errorf("invalid due date"),
			"invalid_date",
			"due_date",
		)
	}
	parsedDates.DueDate = dueDate

	if parsedDates.DueDate.Before(*parsedDates.StartDate) {
		return nil, utility.NewValidationError(
			fmt.Errorf("due date cannot be before start date"),
			"invalid_date",
			"due_date",
		)
	}

	return parsedDates, nil
}

func (ircs *IdentityReviewCampaignService) GetDocumentPath(id uint, token string) (string, error) {
	campaign, err := ircs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utility.NewResourceNotFoundError(err, id)
		}
		return "", err
	}

	var matchedDocument *models.IdentityReviewCampaignDocument
	for _, document := range campaign.Documents {
		if document.Token == token {
			matchedDocument = document
			break
		}
	}
	if matchedDocument == nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" || decryptedFilename != matchedDocument.Name {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join(
		"identity-access-review",
		"review-campaign",
		strconv.FormatUint(uint64(id), 10),
		"documents",
	)
	fullPath := filepath.Join(config.GetEnv("STORAGE_PATH", ""), storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
