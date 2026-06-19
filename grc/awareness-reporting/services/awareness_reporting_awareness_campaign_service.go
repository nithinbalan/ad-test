package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	awarenessRepositories "jethur/grc/awareness-reporting/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AwarenessCampaignServiceInterface interface {
	GetAll(params *awarenessDTOs.AwarenessCampaignListQueryParamsDTO) (*awarenessDTOs.AwarenessCampaignListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*awarenessDTOs.AwarenessCampaignDetailResponseDTO, error)
	GetByID(id uint) (*models.AwarenessCampaign, error)
	Create(currentUserID uint, formDTO *awarenessDTOs.AwarenessCampaignFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *awarenessDTOs.AwarenessCampaignFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(id uint, token string) (string, error)
}

// Service and implementations
type AwarenessCampaignService struct {
	repo                         awarenessRepositories.AwarenessCampaignRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	awarenessCampaignTypeSvc     AwarenessCampaignTypeServiceInterface
	awarenessPrioritySvc         AwarenessPriorityServiceInterface
	awarenessCampaignStatusSvc   AwarenessCampaignStatusServiceInterface
	awarenessEvaluationRatingSvc AwarenessEvaluationRatingServiceInterface
	organizationUnitSvc          settingsServices.OrganizationUnitServiceInterface
	userSvc                      organizationServices.UserServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
}

func NewAwarenessCampaignService(
	repo awarenessRepositories.AwarenessCampaignRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	awarenessCampaignTypeSvc AwarenessCampaignTypeServiceInterface,
	awarenessPrioritySvc AwarenessPriorityServiceInterface,
	awarenessCampaignStatusSvc AwarenessCampaignStatusServiceInterface,
	awarenessEvaluationRatingSvc AwarenessEvaluationRatingServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AwarenessCampaignServiceInterface {
	return &AwarenessCampaignService{
		repo:                         repo,
		activityLogSvc:               activityLogSvc,
		awarenessCampaignTypeSvc:     awarenessCampaignTypeSvc,
		awarenessPrioritySvc:         awarenessPrioritySvc,
		awarenessCampaignStatusSvc:   awarenessCampaignStatusSvc,
		awarenessEvaluationRatingSvc: awarenessEvaluationRatingSvc,
		organizationUnitSvc:          organizationUnitSvc,
		userSvc:                      userSvc,
		notificationSvc:              notificationSvc,
	}
}

func (acs *AwarenessCampaignService) GetAll(params *awarenessDTOs.AwarenessCampaignListQueryParamsDTO) (*awarenessDTOs.AwarenessCampaignListResponseDTO, error) {
	return acs.repo.GetAll(params)
}

func (acs *AwarenessCampaignService) GetByID(id uint) (*models.AwarenessCampaign, error) {
	campaign, err := acs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return campaign, nil
}

func (acs *AwarenessCampaignService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*awarenessDTOs.AwarenessCampaignDetailResponseDTO, error) {
	var campaignDTO awarenessDTOs.AwarenessCampaignDetailResponseDTO

	campaign, err := acs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(campaign, &campaignDTO); err != nil {
		return nil, err
	}

	campaignDTO.CreatedAt = campaign.CreatedAt
	campaignDTO.UpdatedAt = campaign.UpdatedAt

	if campaign.AwarenessCampaignDetails != nil {
		var detailDTO awarenessDTOs.AwarenessCampaignDetailDTO
		if err := mapstructure.Decode(campaign.AwarenessCampaignDetails, &detailDTO); err != nil {
			return nil, err
		}
		detailDTO.AwarenessCampaignID = campaign.AwarenessCampaignDetails.AwarenessCampaignID
		detailDTO.CreatedAt = campaign.AwarenessCampaignDetails.CreatedAt
		detailDTO.UpdatedAt = campaign.AwarenessCampaignDetails.UpdatedAt
		campaignDTO.Detail = &detailDTO
	}

	return &campaignDTO, nil
}

func (acs *AwarenessCampaignService) Create(currentUserID uint, formDTO *awarenessDTOs.AwarenessCampaignFormDTO) (uint, error) {
	tx := acs.repo.BeginTransaction()

	if err := acs.validateAwarenessCampaignFormDTO(formDTO); err != nil {
		return 0, err
	}

	dates, err := acs.parseAwarenessCampaignDates(formDTO, currentUserID)
	if err != nil {
		return 0, err
	}
	if dates.StartDate == nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("start date is required"), "required", "start_date")
	}

	referenceNumber, err := acs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	awarenessCampaign := &models.AwarenessCampaign{
		ReferenceNumber:           referenceNumber,
		AwarenessCampaignTypeID:   formDTO.AwarenessCampaignTypeID,
		CampaignTitle:             formDTO.CampaignTitle,
		Objective:                 formDTO.Objective,
		AwarenessPriorityID:       formDTO.AwarenessPriorityID,
		TargetAudience:            formDTO.TargetAudience,
		OwnerID:                   formDTO.OwnerID,
		OrganizationUnitID:        formDTO.OrganizationUnitID,
		StartDate:                 *dates.StartDate,
		EndDate:                   dates.EndDate,
		AwarenessCampaignStatusID: formDTO.AwarenessCampaignStatusID,
		LinkedSecurityControl:     formDTO.LinkedSecurityControl,
		CreatedByID:               currentUserID,
	}

	id, err := acs.repo.Create(awarenessCampaign, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if formDTO.Detail != nil {
		detail := &models.AwarenessCampaignDetail{
			AwarenessCampaignID:         id,
			AwarenessPackName:           formDTO.Detail.AwarenessPackName,
			RecipientCount:              formDTO.Detail.RecipientCount,
			SentCount:                   formDTO.Detail.SentCount,
			EmailOpenRate:               formDTO.Detail.EmailOpenRate,
			ClickRate:                   formDTO.Detail.ClickRate,
			DataSubmissionRate:          formDTO.Detail.DataSubmissionRate,
			ReportRate:                  formDTO.Detail.ReportRate,
			RiskClassificationID:        formDTO.Detail.RiskClassificationID,
			CampaignResult:              formDTO.Detail.CampaignResult,
			ImprovementRecommendations:  formDTO.Detail.ImprovementRecommendations,
			SubtypeRole:                 formDTO.Detail.SubtypeRole,
			Presenter:                   formDTO.Detail.Presenter,
			Agenda:                      formDTO.Detail.Agenda,
			Venue:                       formDTO.Detail.Venue,
			AttendeeCount:               formDTO.Detail.AttendeeCount,
			AwarenessEvaluationRatingID: formDTO.Detail.AwarenessEvaluationRatingID,
			CreatedByID:                 currentUserID,
		}

		if formDTO.Detail.SessionStartTime != nil && *formDTO.Detail.SessionStartTime != "" {
			sessionStart, err := utility.ParseDateTime(*formDTO.Detail.SessionStartTime)
			if err != nil {
				tx.Rollback()
				return 0, utility.NewValidationError(fmt.Errorf("invalid session start time: %w", err), "invalid_date", "session_start_time")
			}
			detail.SessionStartTime = sessionStart
		}
		if formDTO.Detail.SessionEndTime != nil && *formDTO.Detail.SessionEndTime != "" {
			sessionEnd, err := utility.ParseDateTime(*formDTO.Detail.SessionEndTime)
			if err != nil {
				tx.Rollback()
				return 0, utility.NewValidationError(fmt.Errorf("invalid session end time: %w", err), "invalid_date", "session_end_time")
			}
			detail.SessionEndTime = sessionEnd
		}

		detailID, err := acs.repo.CreateDetail(detail, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if len(formDTO.Detail.AwarenessChannelIDs) > 0 {
			var channels []*models.AwarenessCampaignChannel
			for _, channelID := range formDTO.Detail.AwarenessChannelIDs {
				if channelID == nil {
					continue
				}
				channels = append(channels, &models.AwarenessCampaignChannel{
					AwarenessCampaignDetailsID: detailID,
					AwarenessChannelID:         *channelID,
					CreatedByID:                currentUserID,
				})
			}
			if len(channels) > 0 {
				if err := acs.repo.CreateCampaignChannels(channels, tx); err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}

		var documents []*models.AwarenessCampaignDocument

		for _, doc := range formDTO.Detail.Documents {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			url := utility.GenerateDocumentURL("awareness", "awareness-campaign", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.AwarenessCampaignDocument{
				AwarenessCampaignTypeID:   formDTO.AwarenessCampaignTypeID,
				AwarenessCampaignDetailID: &detailID,
				URL:                       url,
				Name:                      doc.FileName,
				Size:                      doc.Size,
				Token:                     encryptedFilename,
				OriginalName:              doc.OriginalName,
				Type:                      doc.Type,
				CreatedByID:               currentUserID,
			}
			documents = append(documents, document)
		}

		if len(documents) > 0 {
			if err := acs.repo.CreateDocuments(documents, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
			for _, doc := range documents {
				storagePath := filepath.Join("awareness", "awareness-campaign", strconv.FormatUint(uint64(id), 10), "documents")
				if err := utility.StoreFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}

	tx.Commit()

	go func() {
		existingCampaign, err := acs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch awareness campaign", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = acs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessCampaign{}.TableName(),
			id,
			existingCampaign,
			currentUserID,
			"awareness", "awareness_campaign", "campaign", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AwarenessCampaign{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "awareness-campaign"
		action := "owner-assign"
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{formDTO.OwnerID},
			ModuleIdentifier:     "awareness",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Awareness Campaign Ownership Assigned",
			HeadingAr:            "تم تعيين ملكية حملة التوعية",
			TitleEn:              formDTO.CampaignTitle,
			TitleAr:              formDTO.CampaignTitle,
			ReferenceNumber:      awarenessCampaign.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}
		ownerUser, err := acs.userSvc.GetByID(formDTO.OwnerID)
		if err != nil {
			utility.LogError("Failed to fetch owner", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		createdBy, err := acs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created by user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		emailTo := []string{ownerUser.Email}
		subject := fmt.Sprintf("Awareness Campaign/%s: Ownership Assigned", formDTO.CampaignTitle)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                      ownerUser.FirstNameEn,
			"LastNameEn":                       ownerUser.LastNameEn,
			"AwarenessCampaignReferenceNumber": awarenessCampaign.ReferenceNumber,
			"AwarenessCampaignTitle":           formDTO.CampaignTitle,
			"AwarenessCampaignObjective":       formDTO.Objective,
			"CreatedByFirstNameEn":             createdBy.FirstNameEn,
			"CreatedByLastNameEn":              createdBy.LastNameEn,
			"AssignedOn":                       time.Now().Format("2006-01-02"),
			"RedirectURL":                      fmt.Sprintf("%s/awareness/campaign/%d/overview", appURL, id),
		}
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "awareness/templates/email/awareness_campaign_create_owner_mail.html",
			EmailData:            emailData,
		}
		err = acs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "awareness_campaign")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (acs *AwarenessCampaignService) Update(id uint, currentUserID uint, formDTO *awarenessDTOs.AwarenessCampaignFormDTO) error {
	tx := acs.repo.BeginTransaction()

	existingCampaign, err := acs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldOwnerID := existingCampaign.OwnerID

	oldCampaign, err := acs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := acs.validateAwarenessCampaignFormDTO(formDTO); err != nil {
		return err
	}

	dates, err := acs.parseAwarenessCampaignDates(formDTO, currentUserID)
	if err != nil {
		return err
	}

	existingCampaign.AwarenessCampaignTypeID = formDTO.AwarenessCampaignTypeID
	existingCampaign.CampaignTitle = formDTO.CampaignTitle
	existingCampaign.Objective = formDTO.Objective
	existingCampaign.AwarenessPriorityID = formDTO.AwarenessPriorityID
	existingCampaign.TargetAudience = formDTO.TargetAudience
	existingCampaign.OwnerID = formDTO.OwnerID
	existingCampaign.OrganizationUnitID = formDTO.OrganizationUnitID
	if dates.StartDate != nil {
		existingCampaign.StartDate = *dates.StartDate
	}
	existingCampaign.EndDate = dates.EndDate
	existingCampaign.AwarenessCampaignStatusID = formDTO.AwarenessCampaignStatusID
	existingCampaign.LinkedSecurityControl = formDTO.LinkedSecurityControl
	existingCampaign.UpdatedByID = &currentUserID

	if err := acs.repo.Update(existingCampaign, tx); err != nil {
		tx.Rollback()
		return err
	}

	if formDTO.Detail != nil {
		existingDetail, err := acs.repo.GetDetailByAwarenessCampaignID(id, tx)
		isNewDetail := false
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return err
			}
			existingDetail = &models.AwarenessCampaignDetail{
				AwarenessCampaignID: id,
				CreatedByID:         currentUserID,
			}
			isNewDetail = true
		}

		existingDetail.AwarenessPackName = formDTO.Detail.AwarenessPackName
		existingDetail.RecipientCount = formDTO.Detail.RecipientCount
		existingDetail.SentCount = formDTO.Detail.SentCount
		existingDetail.EmailOpenRate = formDTO.Detail.EmailOpenRate
		existingDetail.ClickRate = formDTO.Detail.ClickRate
		existingDetail.DataSubmissionRate = formDTO.Detail.DataSubmissionRate
		existingDetail.ReportRate = formDTO.Detail.ReportRate
		existingDetail.RiskClassificationID = formDTO.Detail.RiskClassificationID
		existingDetail.CampaignResult = formDTO.Detail.CampaignResult
		existingDetail.ImprovementRecommendations = formDTO.Detail.ImprovementRecommendations
		existingDetail.SubtypeRole = formDTO.Detail.SubtypeRole
		existingDetail.Presenter = formDTO.Detail.Presenter
		existingDetail.Agenda = formDTO.Detail.Agenda
		existingDetail.Venue = formDTO.Detail.Venue
		existingDetail.AttendeeCount = formDTO.Detail.AttendeeCount
		existingDetail.AwarenessEvaluationRatingID = formDTO.Detail.AwarenessEvaluationRatingID

		if formDTO.Detail.SessionStartTime != nil && *formDTO.Detail.SessionStartTime != "" {
			sessionStart, err := utility.ParseDateTime(*formDTO.Detail.SessionStartTime)
			if err != nil {
				tx.Rollback()
				return utility.NewValidationError(fmt.Errorf("invalid session start time: %w", err), "invalid_date", "session_start_time")
			}
			existingDetail.SessionStartTime = sessionStart
		}
		if formDTO.Detail.SessionEndTime != nil && *formDTO.Detail.SessionEndTime != "" {
			sessionEnd, err := utility.ParseDateTime(*formDTO.Detail.SessionEndTime)
			if err != nil {
				tx.Rollback()
				return utility.NewValidationError(fmt.Errorf("invalid session end time: %w", err), "invalid_date", "session_end_time")
			}
			existingDetail.SessionEndTime = sessionEnd
		}

		var detailID uint
		if isNewDetail {
			detailID, err = acs.repo.CreateDetail(existingDetail, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			existingDetail.UpdatedByID = &currentUserID
			if err := acs.repo.UpdateDetail(existingDetail, tx); err != nil {
				tx.Rollback()
				return err
			}
			detailID = existingDetail.ID

			if err := acs.repo.DeleteCampaignChannels(detailID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		if len(formDTO.Detail.AwarenessChannelIDs) > 0 {
			var channels []*models.AwarenessCampaignChannel
			for _, channelID := range formDTO.Detail.AwarenessChannelIDs {
				if channelID == nil {
					continue
				}
				channels = append(channels, &models.AwarenessCampaignChannel{
					AwarenessCampaignDetailsID: detailID,
					AwarenessChannelID:         *channelID,
					CreatedByID:                currentUserID,
				})
			}
			if len(channels) > 0 {
				if err := acs.repo.CreateCampaignChannels(channels, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		// documents update
		var documents []*models.AwarenessCampaignDocument
		var deleteDocumentIds []*uint
		var documentsToDelete []*models.AwarenessCampaignDocument

		for _, doc := range formDTO.Detail.Documents {
			if doc.IsNewFile != nil && *doc.IsNewFile {
				dataToEncrypt := []byte(doc.FileName)
				encryptedFilename, err := utility.EncryptData(dataToEncrypt)
				if err != nil {
					tx.Rollback()
					return err
				}
				url := utility.GenerateDocumentURL("awareness", "awareness-campaign", strconv.FormatUint(uint64(id), 10), "documents")

				document := &models.AwarenessCampaignDocument{
					AwarenessCampaignTypeID:   formDTO.AwarenessCampaignTypeID,
					AwarenessCampaignDetailID: &detailID,
					URL:                       url,
					Name:                      doc.FileName,
					Size:                      doc.Size,
					Token:                     encryptedFilename,
					OriginalName:              doc.OriginalName,
					Type:                      doc.Type,
					CreatedByID:               currentUserID,
				}
				documents = append(documents, document)

			} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
				if doc.ID != nil {
					existingDocument, err := acs.repo.GetDocumentByID(*doc.ID)
					if err != nil {
						tx.Rollback()
						return err
					}
					if existingDocument == nil {
						tx.Rollback()
						return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
					}
					deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
					documentsToDelete = append(documentsToDelete, existingDocument)
				}
			}
		}

		if len(documents) > 0 {
			if err := acs.repo.CreateDocuments(documents, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
		if len(deleteDocumentIds) > 0 {
			if err := acs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
		if len(documents) > 0 {
			for _, doc := range documents {
				storagePath := filepath.Join("awareness", "awareness-campaign", strconv.FormatUint(uint64(id), 10), "documents")
				if err := utility.StoreFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
		if len(documentsToDelete) > 0 {
			for _, document := range documentsToDelete {
				storagePath := filepath.Join("awareness", "awareness-campaign", strconv.FormatUint(uint64(id), 10), "documents")
				if err := utility.RemoveFile(document.Name, storagePath); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	tx.Commit()

	go func() {
		existingFullCampaign, err := acs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		oldData := map[string]interface{}{
			"awareness_campaign_type":   oldCampaign.AwarenessCampaignType.TitleEn,
			"awareness_campaign_status": oldCampaign.AwarenessCampaignStatus.TitleEn,
			"campaign_title":            oldCampaign.CampaignTitle,
			"objective":                 oldCampaign.Objective,
			"target_audience":           oldCampaign.TargetAudience,
			"start_date":                oldCampaign.StartDate,
			"end_date":                  oldCampaign.EndDate,
			"linked_security_control":   utility.GetStringValue(oldCampaign.LinkedSecurityControl),
		}

		newData := map[string]interface{}{
			"awareness_campaign_type":   existingFullCampaign.AwarenessCampaignType.TitleEn,
			"awareness_campaign_status": existingFullCampaign.AwarenessCampaignStatus.TitleEn,
			"campaign_title":            existingFullCampaign.CampaignTitle,
			"objective":                 existingFullCampaign.Objective,
			"target_audience":           existingFullCampaign.TargetAudience,
			"start_date":                existingFullCampaign.StartDate,
			"end_date":                  existingFullCampaign.EndDate,
			"linked_security_control":   utility.GetStringValue(existingFullCampaign.LinkedSecurityControl),
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := acs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AwarenessCampaign{}.TableName(),
				id,
				changedData,
				currentUserID,
				"awareness", "awareness_campaign", "campaign", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AwarenessCampaign{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				utility.LogError("Failed to log activity", fields)
			}
		}

		// Send notifications to old and new owner if ownership changed
		if oldOwnerID != existingCampaign.OwnerID {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "awareness-campaign"

			oldOwnerAction := "owner-remove"
			oldOwnerNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldOwnerID},
				ModuleIdentifier:     "awareness",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Awareness Campaign Ownership Removed",
				HeadingAr:            "تمت إزالة ملكية حملة التوعية",
				TitleEn:              existingCampaign.CampaignTitle,
				TitleAr:              existingCampaign.CampaignTitle,
				ReferenceNumber:      existingCampaign.ReferenceNumber,
				Entity:               &entity,
				Action:               &oldOwnerAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newOwnerAction := "owner-assign"
			newOwnerNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{existingCampaign.OwnerID},
				ModuleIdentifier:     "awareness",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Awareness Campaign Ownership Assigned",
				HeadingAr:            "تم تعيين ملكية حملة التوعية",
				TitleEn:              existingCampaign.CampaignTitle,
				TitleAr:              existingCampaign.CampaignTitle,
				ReferenceNumber:      existingCampaign.ReferenceNumber,
				Entity:               &entity,
				Action:               &newOwnerAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldOwnerUser, err := acs.userSvc.GetByID(oldOwnerID)
			if err != nil {
				utility.LogError("Failed to fetch old owner", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newOwnerUser, err := acs.userSvc.GetByID(existingCampaign.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch new owner", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			updatedByUser, err := acs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch updating user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			oldOwnerEmailTo := []string{oldOwnerUser.Email}
			oldOwnerSubject := fmt.Sprintf("Awareness Campaign/%s: Ownership Removed", existingCampaign.CampaignTitle)
			oldOwnerEmailData := map[string]interface{}{
				"FirstNameEn":                      oldOwnerUser.FirstNameEn,
				"LastNameEn":                       oldOwnerUser.LastNameEn,
				"AwarenessCampaignReferenceNumber": existingCampaign.ReferenceNumber,
				"AwarenessCampaignTitle":           existingCampaign.CampaignTitle,
				"AwarenessCampaignObjective":       existingCampaign.Objective,
				"UpdatedByFirstNameEn":             updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":              updatedByUser.LastNameEn,
				"RedirectURL":                      fmt.Sprintf("%s/awareness/campaign/%d/overview", appURL, id),
			}

			oldOwnerEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldOwnerEmailTo,
				Subject:              oldOwnerSubject,
				TemplateRelativePath: "awareness/templates/email/awareness_campaign_remove_owner_mail.html",
				EmailData:            oldOwnerEmailData,
			}

			newOwnerEmailTo := []string{newOwnerUser.Email}
			newOwnerSubject := fmt.Sprintf("Awareness Campaign/%s: Ownership Assigned", existingCampaign.CampaignTitle)
			newOwnerEmailData := map[string]interface{}{
				"FirstNameEn":                      newOwnerUser.FirstNameEn,
				"LastNameEn":                       newOwnerUser.LastNameEn,
				"AwarenessCampaignReferenceNumber": existingCampaign.ReferenceNumber,
				"AwarenessCampaignTitle":           existingCampaign.CampaignTitle,
				"AwarenessCampaignObjective":       existingCampaign.Objective,
				"AssignedByFirstNameEn":            updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":             updatedByUser.LastNameEn,
				"AssignedOn":                       time.Now().Format("2006-01-02"),
				"RedirectURL":                      fmt.Sprintf("%s/awareness/campaign/%d/overview", appURL, id),
			}

			newOwnerEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newOwnerEmailTo,
				Subject:              newOwnerSubject,
				TemplateRelativePath: "awareness/templates/email/awareness_campaign_change_owner_mail.html",
				EmailData:            newOwnerEmailData,
			}

			err = acs.notificationSvc.SendNotification(newOwnerNotificationDTO, newOwnerEmailNotificationDTO, "awareness_campaign")
			if err != nil {
				utility.LogError("Failed to send notifications to new owner", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = acs.notificationSvc.SendNotification(oldOwnerNotificationDTO, oldOwnerEmailNotificationDTO, "awareness_campaign")
			if err != nil {
				utility.LogError("Failed to send notifications to old owner", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (acs *AwarenessCampaignService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := acs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedCampaigns []map[string]interface{}

	var notifications []struct {
		OwnerID         uint
		TitleEn         string
		Objective       string
		ReferenceNumber string
	}

	documentsToDelete := make(map[uint][]*models.AwarenessCampaignDocument)

	for _, id := range ids {
		campaign, err := acs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedCampaigns = append(deletedCampaigns, map[string]interface{}{
			"id":    id,
			"title": campaign.CampaignTitle,
		})

		notifications = append(notifications, struct {
			OwnerID         uint
			TitleEn         string
			Objective       string
			ReferenceNumber string
		}{
			OwnerID:         campaign.OwnerID,
			TitleEn:         campaign.CampaignTitle,
			Objective:       campaign.Objective,
			ReferenceNumber: campaign.ReferenceNumber,
		})

		if campaign.AwarenessCampaignDetails != nil && len(campaign.AwarenessCampaignDetails.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range campaign.AwarenessCampaignDetails.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
				}
			}
			if err := acs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if campaign.AwarenessCampaignDetails != nil {
			if err := acs.repo.DeleteCampaignChannels(campaign.AwarenessCampaignDetails.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := acs.repo.DeleteDetail(campaign.AwarenessCampaignDetails.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := acs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		for campaignID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("awareness", "awareness-campaign", strconv.FormatUint(uint64(campaignID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					utility.LogError("Failed to delete file", map[string]interface{}{
						"error": err.Error(),
						"name":  doc.Name,
					})
				}
			}
		}
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := acs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessCampaign{}.TableName(),
				successfulIds,
				deletedCampaigns,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}

			// Send DB and email notifications for each successfully deleted campaign
			for _, notification := range notifications {
				urlParam := ""
				entity := "awareness-campaign"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{notification.OwnerID},
					ModuleIdentifier:     "awareness",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Awareness Campaign Deleted",
					HeadingAr:            "تم حذف حملة التوعية",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleEn,
					ReferenceNumber:      notification.ReferenceNumber,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := acs.userSvc.GetByID(notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Awareness Campaign/%s: Deleted", notification.TitleEn)
				templateRelativePath := "awareness/templates/email/awareness_campaign_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":                      ownerUser.FirstNameEn,
					"LastNameEn":                       ownerUser.LastNameEn,
					"AwarenessCampaignReferenceNumber": notification.ReferenceNumber,
					"AwarenessCampaignTitle":           notification.TitleEn,
					"AwarenessCampaignObjective":       notification.Objective,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = acs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "awareness_campaign")
				if err != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}()
	}

	return failedIds, nil
}

func (acs *AwarenessCampaignService) validateAwarenessCampaignFormDTO(formDTO *awarenessDTOs.AwarenessCampaignFormDTO) error {
	if _, err := acs.awarenessCampaignTypeSvc.GetByID(formDTO.AwarenessCampaignTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("awareness campaign type does not exist"), "not_exist", "awareness campaign type")
	}

	if formDTO.AwarenessPriorityID != nil {
		if _, err := acs.awarenessPrioritySvc.GetByID(*formDTO.AwarenessPriorityID); err != nil {
			return utility.NewValidationError(fmt.Errorf("awareness priority does not exist"), "not_exist", "awareness priority")
		}
	}

	if _, err := acs.awarenessCampaignStatusSvc.GetByID(formDTO.AwarenessCampaignStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("awareness campaign status does not exist"), "not_exist", "awareness campaign status")
	}

	if _, err := acs.userSvc.GetByID(formDTO.OwnerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("owner does not exist"), "not_exist", "owner")
	}

	if _, err := acs.organizationUnitSvc.GetByID(formDTO.OrganizationUnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization unit")
	}

	if formDTO.Detail != nil {
		if formDTO.Detail.RiskClassificationID != nil {
			if _, err := acs.awarenessPrioritySvc.GetByID(*formDTO.Detail.RiskClassificationID); err != nil {
				return utility.NewValidationError(fmt.Errorf("risk classification does not exist"), "not_exist", "risk classification")
			}
		}

		if formDTO.Detail.AwarenessEvaluationRatingID != nil {
			if _, err := acs.awarenessEvaluationRatingSvc.GetByID(*formDTO.Detail.AwarenessEvaluationRatingID); err != nil {
				return utility.NewValidationError(fmt.Errorf("awareness evaluation rating does not exist"), "not_exist", "awareness evaluation rating")
			}
		}
	}

	return nil
}

func (acs *AwarenessCampaignService) parseAwarenessCampaignDates(dto *awarenessDTOs.AwarenessCampaignFormDTO, currentUserID uint) (*awarenessDTOs.AwarenessCampaignParsedDates, error) {
	parsed := &awarenessDTOs.AwarenessCampaignParsedDates{}

	userPreference, err := acs.userSvc.GetCurrentUserDetails(currentUserID)
	if err != nil {
		return nil, err
	}

	timezoneName := userPreference.UserSetting.Timezone.Name

	if dto.StartDate != "" {
		parsed.StartDate, err = utility.ConvertToUTC(dto.StartDate, timezoneName)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date: %w", err), "invalid_date", "start_date")
		}
	}

	if dto.EndDate != nil && *dto.EndDate != "" {
		parsed.EndDate, err = utility.ConvertToUTC(*dto.EndDate, timezoneName)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date: %w", err), "invalid_date", "end_date")
		}
	}

	return parsed, nil
}

func (acs *AwarenessCampaignService) GetDocumentPath(id uint, token string) (string, error) {
	existingCampaign, err := acs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", utility.NewResourceNotFoundError(err, id)
		}
		return "", err
	}

	if existingCampaign.AwarenessCampaignDetails == nil || len(existingCampaign.AwarenessCampaignDetails.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the awareness campaign")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("awareness", "awareness-campaign", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
