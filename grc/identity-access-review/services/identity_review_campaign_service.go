package services

import (
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IdentityReviewCampaignServiceInterface interface {
	GetAll(params *dtos.IdentityReviewCampaignListQueryParamsDTO) (*dtos.IdentityReviewCampaignListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewCampaignDetailDTO, error)
	Create(currentUserID uint, formDTO *dtos.IdentityReviewCampaignFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *dtos.IdentityReviewCampaignFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(id uint, token string) (string, error)
	CalculateAndUpdateCompletionPercentage(id uint, tx ...utility.DatabaseTransaction) error
	ExportZip(id uint, currentUserID uint) (string, error)
	SetIdentityReviewCampaignWorkflowService(svc IdentityReviewCampaignWorkflowServiceInterface)
}

func (ircs *IdentityReviewCampaignService) CalculateAndUpdateCompletionPercentage(id uint, tx ...utility.DatabaseTransaction) error {
	return ircs.repo.CalculateAndUpdateCompletionPercentage(id, tx...)
}

type IdentityReviewCampaignService struct {
	repo                              repositories.IdentityReviewCampaignRepositoryInterface
	activityLogSvc                    logServices.ActivityLogServiceInterface
	reviewTypeSvc                     ReviewTypeServiceInterface
	reviewCampaignStatusSvc           ReviewCampaignStatusServiceInterface
	userSvc                           organizationServices.UserServiceInterface
	reportTemplateSettingSvc          settingsServices.ReportTemplateSettingServiceInterface
	identityReviewCampaignWorkflowSvc IdentityReviewCampaignWorkflowServiceInterface
	notificationSvc                   commonServices.NotificationServiceInterface
}

func (ircs *IdentityReviewCampaignService) SetIdentityReviewCampaignWorkflowService(svc IdentityReviewCampaignWorkflowServiceInterface) {
	ircs.identityReviewCampaignWorkflowSvc = svc
}

func NewIdentityReviewCampaignService(
	repo repositories.IdentityReviewCampaignRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	reviewTypeSvc ReviewTypeServiceInterface,
	reviewCampaignStatusSvc ReviewCampaignStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	reportTemplateSettingSvc settingsServices.ReportTemplateSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) IdentityReviewCampaignServiceInterface {
	return &IdentityReviewCampaignService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		reviewTypeSvc:            reviewTypeSvc,
		reviewCampaignStatusSvc:  reviewCampaignStatusSvc,
		userSvc:                  userSvc,
		reportTemplateSettingSvc: reportTemplateSettingSvc,
		notificationSvc:          notificationSvc,
	}
}

func (ircs *IdentityReviewCampaignService) GetAll(params *dtos.IdentityReviewCampaignListQueryParamsDTO) (*dtos.IdentityReviewCampaignListResponseDTO, error) {
	return ircs.repo.GetAll(params)
}

func (ircs *IdentityReviewCampaignService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaign, error) {
	campaign, err := ircs.repo.GetByID(id, tx...)
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
	if formDTO.AssetID == 0 {
		return 0, utility.NewValidationError(
			fmt.Errorf("asset is required"),
			"required",
			"asset_id",
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
		WorkflowID:         formDTO.WorkflowID,
		CampaignStatusID:   plannedStatus.ID,
		ReminderDaysOffset: formDTO.ReminderDaysOffset,
		CreatedByID:        currentUserID,
	}

	id, err := ircs.repo.Create(campaign, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	assetID := formDTO.AssetID
	systems := []*models.IdentityReviewCampaignSystem{{
		IdentityReviewCampaignID: id,
		AssetRegisterID:          &assetID,
		CreatedByID:              currentUserID,
	}}

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

		urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, id)
		entity := "identity-access-review-campaign"
		action := "reviewer-assign"
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{campaign.ReviewerID},
			ModuleIdentifier:     "identity_access_review",
			NotificationType:     string(commonDTOs.AssignmentType),
			ReferenceNumber:      campaign.ReferenceNumber,
			HeadingEn:            "Identity Access Review Campaign Assigned",
			HeadingAr:            "تم تعيين حملة مراجعة الهوية والوصول",
			TitleEn:              campaign.CampaignName,
			TitleAr:              campaign.CampaignName,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		if err := ircs.notificationSvc.SendNotification(dbNotificationDTO, nil, "identity_access_review_campaign"); err != nil {
			utility.LogError("Failed to send identity access review campaign creation notification", map[string]interface{}{
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

	if campaign.SubmittedAt != nil && campaign.WorkflowID != formDTO.WorkflowID {
		tx.Rollback()
		return utility.NewValidationError(
			fmt.Errorf("workflow cannot be changed after campaign submission"),
			"workflow_already_submitted",
			"workflow_id",
		)
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
	campaign.WorkflowID = formDTO.WorkflowID
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

	assetID := formDTO.AssetID
	systems := []*models.IdentityReviewCampaignSystem{{
		IdentityReviewCampaignID: id,
		AssetRegisterID:          &assetID,
		CreatedByID:              currentUserID,
	}}

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

func (ircs *IdentityReviewCampaignService) ExportZip(id uint, currentUserID uint) (string, error) {
	if _, err := ircs.reportTemplateSettingSvc.GetImagePath("logo"); err != nil {
		return "", utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := ircs.reportTemplateSettingSvc.GetImagePath("cover-page-bg"); err != nil {
		return "", utility.NewForbiddenError("cover_image_missing")
	}

	campaign, err := ircs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utility.NewResourceNotFoundError(err, id)
		}
		return "", err
	}

	users, err := ircs.repo.GetExportUsers(id)
	if err != nil {
		return "", err
	}

	workflow, err := ircs.identityReviewCampaignWorkflowSvc.GetWorkflowList(id)
	if err != nil {
		return "", err
	}

	completionPercentage := campaign.CompletionPercentage
	if completionPercentage < 0 {
		completionPercentage = 0
	} else if completionPercentage > 100 {
		completionPercentage = 100
	}

	reviewerName := "-"
	if campaign.Reviewer != nil {
		reviewerName = campaign.Reviewer.FirstNameEn + " " + campaign.Reviewer.LastNameEn
	}

	reviewTypeTitleEn := "-"
	if campaign.ReviewType != nil {
		reviewTypeTitleEn = campaign.ReviewType.TitleEn
	}

	reminderDaysOffset := "-"
	if campaign.ReminderDaysOffset != nil {
		reminderDaysOffset = fmt.Sprintf("%d days", *campaign.ReminderDaysOffset)
	}

	assetTitleEn := "-"
	if len(campaign.Systems) > 0 && campaign.Systems[0].AssetRegister != nil {
		assetTitleEn = campaign.Systems[0].AssetRegister.TitleEn
	}

	campaignStatusTitleEn := "-"
	campaignStatusColorCode := "#6b7280"
	if campaign.CampaignStatus != nil {
		campaignStatusTitleEn = campaign.CampaignStatus.TitleEn
		campaignStatusColorCode = campaign.CampaignStatus.ColorCode
	}

	viewModel := dtos.IdentityReviewCampaignExportPDFDTO{
		CampaignName:            campaign.CampaignName,
		ReferenceNumber:         campaign.ReferenceNumber,
		Description:             campaign.Description,
		StartDate:               campaign.StartDate.Format("02 Jan 2006"),
		DueDate:                 campaign.DueDate.Format("02 Jan 2006"),
		GeneratedAt:             time.Now().Format("02 Jan 2006"),
		ReviewTypeTitleEn:       reviewTypeTitleEn,
		ReminderDaysOffset:      reminderDaysOffset,
		AssetTitleEn:            assetTitleEn,
		CampaignStatusTitleEn:   campaignStatusTitleEn,
		CampaignStatusColorCode: campaignStatusColorCode,
		ReviewerName:            reviewerName,
		CompletionPercentage:    completionPercentage,
		CompletionProgressDash:  float64(completionPercentage) * 4.555,
		Users:                   users,
		Workflow:                workflow,
	}

	bodyContent, err := utility.GeneratePDFContent(viewModel, []string{"identity-access-review/templates/pdf/identity_review_campaign_body.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF body: %w", err)
	}
	coverContent, err := utility.GeneratePDFContent(viewModel, []string{"identity-access-review/templates/pdf/identity_review_campaign_cover.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF cover: %w", err)
	}
	headerContent, err := utility.GeneratePDFContent(viewModel, []string{"identity-access-review/templates/pdf/identity_review_campaign_header.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF header: %w", err)
	}
	footerContent, err := utility.GeneratePDFContent(viewModel, []string{"identity-access-review/templates/pdf/identity_review_campaign_footer.html"}, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate PDF footer: %w", err)
	}

	exportDir, err := os.MkdirTemp("", fmt.Sprintf("identity-review-campaign-%d-", id))
	if err != nil {
		return "", fmt.Errorf("failed to create temporary export directory: %w", err)
	}
	defer os.RemoveAll(exportDir)

	pdfFileName := fmt.Sprintf("IdentityAccessReviewCampaign_%d.pdf", id)
	pdfFilePath := filepath.Join(exportDir, pdfFileName)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, bodyContent, pdfFilePath, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	filesToZip := []string{pdfFilePath}
	usedFileNames := map[string]bool{pdfFileName: true}
	for _, document := range campaign.Documents {
		documentPath, err := ircs.GetDocumentPath(id, document.Token)
		if err != nil {
			return "", fmt.Errorf("failed to include document %d in export: %w", document.ID, err)
		}

		documentFileName := filepath.Base(document.OriginalName)
		if documentFileName == "." || documentFileName == string(filepath.Separator) || documentFileName == "" {
			documentFileName = filepath.Base(document.Name)
		}
		if usedFileNames[documentFileName] {
			documentFileName = fmt.Sprintf("%d_%s", document.ID, documentFileName)
		}
		usedFileNames[documentFileName] = true

		stagedDocumentPath := filepath.Join(exportDir, documentFileName)
		if err := utility.CopyFile(documentPath, stagedDocumentPath); err != nil {
			return "", fmt.Errorf("failed to stage document %d for export: %w", document.ID, err)
		}
		filesToZip = append(filesToZip, stagedDocumentPath)
	}

	zipDir := filepath.Join("storage", "zip")
	if err := os.MkdirAll(zipDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create ZIP directory: %w", err)
	}

	zipFilePath := filepath.Join(zipDir, fmt.Sprintf("identity_review_campaign_%d.zip", id))
	if err := utility.CreateZipFile(zipFilePath, filesToZip); err != nil {
		_ = os.Remove(zipFilePath)
		return "", fmt.Errorf("failed to create ZIP file: %w", err)
	}

	return zipFilePath, nil
}
