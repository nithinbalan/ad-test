package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditDraftFindingConversationServiceInterface interface {
	GetAll(auditDraftFindingID uint, queryParams *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error)
	GetByID(auditDraftFindingID uint, id uint) (*models.AuditDraftFindingConversation, error)
	GetFullDetailsByID(auditDraftFindingID uint, id uint) (*dtos.AuditDraftFindingConversationDetailDTO, error)
	Create(auditDraftFindingID uint, currentUserID uint, formDTO *dtos.AuditDraftFindingConversationFormDTO) (uint, error)
	Update(auditDraftFindingID uint, id uint, currentUserID uint, auditFinding *dtos.AuditDraftFindingConversationFormDTO) error
	Delete(auditDraftFindingID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(auditDraftFindingID uint, id uint, token string) (string, error)
}

type AuditDraftFindingConversationService struct {
	repo                       repositories.AuditDraftFindingConversationRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	auditDraftFindingSvc       AuditDraftFindingServiceInterface
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface
	auditSvc                   AuditServiceInterface
	userSvc                    organizationServices.UserServiceInterface
	notificationSvc            commonServices.NotificationServiceInterface
	auditDraftReportSvc        AuditDraftReportServiceInterface
}

func NewAuditDraftFindingConversationService(
	repo repositories.AuditDraftFindingConversationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditDraftFindingSvc AuditDraftFindingServiceInterface,
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface,
	auditSvc AuditServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	auditDraftReportSvc AuditDraftReportServiceInterface,
) AuditDraftFindingConversationServiceInterface {
	return &AuditDraftFindingConversationService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		auditDraftFindingSvc:       auditDraftFindingSvc,
		auditDraftFindingStatusSvc: auditDraftFindingStatusSvc,
		auditSvc:                   auditSvc,
		userSvc:                    userSvc,
		notificationSvc:            notificationSvc,
		auditDraftReportSvc:        auditDraftReportSvc,
	}
}

func (adfcs *AuditDraftFindingConversationService) GetAll(auditDraftFindingID uint, queryParams *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error) {
	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	return adfcs.repo.GetAll(auditDraftFindingID, queryParams)
}

func (adfcs AuditDraftFindingConversationService) GetByID(auditDraftFindingID uint, id uint) (*models.AuditDraftFindingConversation, error) {
	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	auditDraftFindingConversation, err := adfcs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftFindingConversation, nil
}

func (adfcs AuditDraftFindingConversationService) GetFullDetailsByID(auditDraftFindingID uint, id uint) (*dtos.AuditDraftFindingConversationDetailDTO, error) {
	var auditDraftFindingConversationDTO dtos.AuditDraftFindingConversationDetailDTO

	if _, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID); err != nil {
		return nil, err
	}

	auditDraftFindingConversation, err := adfcs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditDraftFindingConversation, &auditDraftFindingConversationDTO)
	if err != nil {
		return nil, err
	}

	auditDraftFindingConversationDTO.CreatedAt = auditDraftFindingConversation.CreatedAt
	auditDraftFindingConversationDTO.UpdatedAt = auditDraftFindingConversation.UpdatedAt

	return &auditDraftFindingConversationDTO, nil
}
func (adfcs *AuditDraftFindingConversationService) Create(auditDraftFindingID uint, currentUserID uint, formDTO *dtos.AuditDraftFindingConversationFormDTO) (uint, error) {
	tx := adfcs.repo.BeginTransaction()

	existingAuditDraftFinding, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID)
	if err != nil {
		return 0, err
	}
	// Check if the type is 'AR', restrict to auditees or audit members
	if formDTO.Type == "AR" {
		// Fetch the related audit by AuditDraftFindingID
		auditID := existingAuditDraftFinding.AuditID
		audit, err := adfcs.auditSvc.GetFullDetailsByID(auditID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Check if the current user is the TeamLead
		isAuthorizedUser := false
		if uint(audit.AuditPlan.TeamLead.ID) == currentUserID {
			isAuthorizedUser = true
		}

		// Check if the current user is one of the AuditPlan Members
		if !isAuthorizedUser {
			for _, member := range audit.AuditPlan.Members {
				if member.UserID == currentUserID {
					isAuthorizedUser = true
					break
				}
			}
		}

		// If the current user is not authorized, rollback and return forbidden error
		if !isAuthorizedUser {
			tx.Rollback()
			return 0, utility.NewForbiddenError("The current user is not authorized to perform this action.")
		}
	}
	conversation := &models.AuditDraftFindingConversation{
		AuditDraftFindingID: auditDraftFindingID,
		Comment:             formDTO.Comment,
		Type:                formDTO.Type,
		IsAccepted:          *formDTO.IsAccepted,
		IsDismissed:         *formDTO.IsDismissed,
		IsReopened:          *formDTO.IsReopened,
		CreatedBYID:         currentUserID,
	}
	id, err := adfcs.repo.Create(conversation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Check if the form indicates that the finding should be accepted
	if formDTO.IsAccepted != nil && *formDTO.IsAccepted {
		// Check if the finding is already accepted
		auditDraftFindingStatus, err := adfcs.auditDraftFindingStatusSvc.GetByType("accepted")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if existingAuditDraftFinding.AuditDraftFindingStatusID == auditDraftFindingStatus.ID {
			tx.Rollback()
			return 0, utility.NewForbiddenError("The finding has already been accepted.")
		}
		// Update Audit Draft Finding Status
		if err := adfcs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatus.ID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Check if the form indicates that the finding should be IsDismissed
	if formDTO.IsDismissed != nil && *formDTO.IsDismissed {
		// Check if the finding is already IsDismissed
		auditDraftFindingStatusDismissed, err := adfcs.auditDraftFindingStatusSvc.GetByType("dismissed")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if existingAuditDraftFinding.AuditDraftFindingStatus != nil && existingAuditDraftFinding.AuditDraftFindingStatus.IsFinalized {
			return 0, utility.NewForbiddenError("Action not allowed. The audit finding must be in not finalized' status before any further actions can be taken.")
		}
		// Update Audit Draft Finding Status
		if err := adfcs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusDismissed.ID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Check if the form indicates that the finding should be IsReopened
	if formDTO.IsReopened != nil && *formDTO.IsReopened {
		// Check if the finding is already IsReopened
		auditDraftFindingStatusOpen, err := adfcs.auditDraftFindingStatusSvc.GetByType("open")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		auditDraftFindingStatusDismissed, err := adfcs.auditDraftFindingStatusSvc.GetByType("dismissed")
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if existingAuditDraftFinding.AuditDraftFindingStatus != nil && existingAuditDraftFinding.AuditDraftFindingStatus.ID != auditDraftFindingStatusDismissed.ID {
			return 0, utility.NewForbiddenError("Action not allowed. The audit finding status must be 'Dismissed' before any further actions can be taken.")
		}
		// Update Audit Draft Finding Status
		if err := adfcs.auditDraftFindingSvc.StatusUpdate(auditDraftFindingID, currentUserID, auditDraftFindingStatusOpen.ID); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	var documents []*models.AuditDraftFindingConversationDocument
	for _, doc := range formDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(existingAuditDraftFinding.AuditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.AuditDraftFindingConversationDocument{
			AuditDraftFindingConversationID: id,
			URL:                             url,
			Name:                            doc.FileName,
			Size:                            doc.Size,
			Token:                           encryptedFilename,
			OriginalName:                    doc.OriginalName,
			Type:                            doc.Type,
			CreatedBYID:                     currentUserID,
		}

		documents = append(documents, document)
	}
	if len(documents) > 0 {
		if err := adfcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(existingAuditDraftFinding.AuditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	tx.Commit()

	go func() {
		auditID := existingAuditDraftFinding.AuditID
		auditDraftReport, err := adfcs.auditDraftReportSvc.GetByAuditID(auditID)
		if err != nil {
			utility.LogError("Failed to fetch draft report", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		if formDTO.IsReopened != nil && *formDTO.IsReopened {
			urlParam := fmt.Sprintf(`{"id":%d}`, auditDraftReport.ID)
			entity := "audit-draft-report"
			action := "reopened"

			// Create a map to track unique recipients
			recipientsMap := make(map[uint]bool)

			// Add participants excluding the current user
			for _, participant := range auditDraftReport.Audit.AuditPlan.Members {
				if participant.UserID != currentUserID {
					recipientsMap[participant.UserID] = true
				}
			}

			// Add InternalReviewSubmittedBYID
			recipientsMap[*auditDraftReport.Audit.AuditPlan.TeamLeadID] = true

			// Extract unique recipients into a slice
			var uniqueRecipients []uint
			for userID := range recipientsMap {
				uniqueRecipients = append(uniqueRecipients, userID)
			}

			// Now `uniqueRecipients` contains all unique UserIDs

			// Iterate through recipients
			for _, userID := range uniqueRecipients {
				responsibleUser, err := adfcs.userSvc.GetByID(userID)
				if err != nil {
					utility.LogError("Failed to fetch review user", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				assignedBy, err := adfcs.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch assigning user", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Audit Draft Report /%s: Findings Reopened Notification", auditDraftReport.Audit.AuditPlan.Title)
				appURL := utility.GetAppURL()

				emailData := map[string]interface{}{
					"FirstNameEn":           responsibleUser.FirstNameEn,
					"LastNameEn":            responsibleUser.LastNameEn,
					"AuditTitle":            auditDraftReport.Audit.AuditPlan.Title,
					"FindingTitle":          existingAuditDraftFinding.Title,
					"FindingDescription":    existingAuditDraftFinding.Description,
					"FindingCategoryTitle":  existingAuditDraftFinding.AuditFindingCategory.TitleEn,
					"ReopenedByFirstNameEn": assignedBy.FirstNameEn,
					"ReopenedByLastNameEn":  assignedBy.LastNameEn,
					"ReopenedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":           fmt.Sprintf("%s/audit/audits/%d/draft-report", appURL, auditID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "audit/templates/email/audit_draft_report_finding_reopened_user_mail.html",
					EmailData:            emailData,
				}

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               userID,
					ModuleIdentifier:     "audit",
					NotificationType:     string(commonDTOs.ReviewType),
					HeadingEn:            "Audit Draft Report Finding Reopened",
					HeadingAr:            "إعادة فتح نتيجة تقرير مسودة التدقيق",
					TitleEn:              existingAuditDraftFinding.Title,
					TitleAr:              existingAuditDraftFinding.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Send Notification and Email
				err = adfcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
		if formDTO.IsDismissed != nil && *formDTO.IsDismissed {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "audit-draft-report"
			action := "dismissed"

			existingAuditDraftReportFullDetails, err := adfcs.auditDraftReportSvc.GetByAuditID(auditID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			if existingAuditDraftReportFullDetails.Audit.AuditPlan.AuditableItems != nil {
				for _, item := range existingAuditDraftReportFullDetails.Audit.AuditPlan.AuditableItems {
					if item.AuditableItem.ResponsibleUser == nil {
						utility.LogError("Failed to fetch assigning user", map[string]interface{}{
							"error": err.Error(),
						})
						return
					}
					responsibleUser := item.AuditableItem.ResponsibleUser
					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserID:               responsibleUser.ID,
						ModuleIdentifier:     "audit_insight_draft_reports",
						NotificationType:     string(commonDTOs.ReviewType),
						HeadingEn:            "Audit Draft Report Finding Dismissed",
						HeadingAr:            "رفض نتيجة تقرير مسودة التدقيق",
						TitleEn:              existingAuditDraftReportFullDetails.Audit.AuditPlan.Title,
						TitleAr:              existingAuditDraftReportFullDetails.Audit.AuditPlan.Title,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					assignedBy, err := adfcs.userSvc.GetByID(currentUserID)
					if err != nil {
						utility.LogError("Failed to fetch assigning user", map[string]interface{}{
							"error": err.Error(),
						})
						return
					}

					emailTo := []string{responsibleUser.Email}
					subject := fmt.Sprintf("Audit Draft Report /%s: Finding /%s: Dismissed", existingAuditDraftReportFullDetails.Audit.AuditPlan.Title, existingAuditDraftFinding.Title)
					appURL := utility.GetAppURL()
					emailData := map[string]interface{}{
						"FirstNameEn":            responsibleUser.FirstNameEn,
						"LastNameEn":             responsibleUser.LastNameEn,
						"AuditTitle":             auditDraftReport.Audit.AuditPlan.Title,
						"FindingTitle":           existingAuditDraftFinding.Title,
						"FindingDescription":     existingAuditDraftFinding.Description,
						"FindingCategoryTitle":   existingAuditDraftFinding.AuditFindingCategory.TitleEn,
						"DismissedByFirstNameEn": assignedBy.FirstNameEn,
						"DismissedByLastNameEn":  assignedBy.LastNameEn,
						"DismissedOn":            time.Now().Format("2006-01-02"),
						"RedirectURL":            fmt.Sprintf("%s/audit-insight/draft-report/%d/details", appURL, id),
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: "audit/templates/email/audit_draft_report_finding_dismissed_mail.html",
						EmailData:            emailData,
					}

					// Send Notification and Email
					err = adfcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
					if err != nil {
						utility.LogError("Failed to send notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}
		existingAuditDraftFindingConversation, err := adfcs.GetFullDetailsByID(auditDraftFindingID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit draft finding", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = adfcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingConversation{}.TableName(),
			id,
			existingAuditDraftFindingConversation,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingConversation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (adfcs *AuditDraftFindingConversationService) Update(auditDraftFindingID uint, id uint, currentUserID uint, auditDraftFindingConversationFormDTO *dtos.AuditDraftFindingConversationFormDTO) error {
	tx := adfcs.repo.BeginTransaction()

	existingAuditDraftFinding, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID)
	if err != nil {
		return err
	}

	oldAuditDraftFindingConversation, err := adfcs.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditDraftFindingConversation, err := adfcs.GetByID(auditDraftFindingID, id)

	if err != nil {
		return err
	}

	// Update Audit Finding ActionPlan
	existingAuditDraftFindingConversation.Comment = auditDraftFindingConversationFormDTO.Comment
	existingAuditDraftFindingConversation.Type = auditDraftFindingConversationFormDTO.Type
	existingAuditDraftFindingConversation.UpdatedBYID = &currentUserID

	err = adfcs.repo.Update(existingAuditDraftFindingConversation, tx)

	if err != nil {
		return err
	}
	auditID := existingAuditDraftFinding.AuditID
	//update document
	var documents []*models.AuditDraftFindingConversationDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditDraftFindingConversationDocument

	for _, doc := range auditDraftFindingConversationFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new Control test Document instance
			document := &models.AuditDraftFindingConversationDocument{
				AuditDraftFindingConversationID: id,
				URL:                             url,
				Name:                            doc.FileName,
				Size:                            doc.Size,
				Token:                           encryptedFileName,
				OriginalName:                    doc.OriginalName,
				Type:                            doc.Type,
				CreatedBYID:                     currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := adfcs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := adfcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := adfcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	go func() {
		existingAuditDraftFindingConversation, err := adfcs.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"comment": oldAuditDraftFindingConversation.Comment,
			"type":    oldAuditDraftFindingConversation.Type,
		}

		newData := map[string]interface{}{
			"comment": existingAuditDraftFindingConversation.Comment,
			"type":    existingAuditDraftFindingConversation.Type,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := adfcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditDraftFindingConversation{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditDraftFindingConversation{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}

func (adfcs *AuditDraftFindingConversationService) Delete(auditDraftFindingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := adfcs.repo.BeginTransaction()
	existingAuditDraftFinding, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID)
	if err != nil {
		return nil, err
	}
	auditID := existingAuditDraftFinding.AuditID

	var failedIds []uint
	var successfulIds []uint
	var deletedAuditDraftFindingConversation []map[string]interface{}
	documentsToDelete := make(map[uint][]*dtos.AuditDraftFindingConversationDocumentDTO)

	for _, id := range ids {
		auditDraftFindingConversation, err := adfcs.GetFullDetailsByID(auditDraftFindingID, id)
		if err != nil {
			continue
		}
		if len(auditDraftFindingConversation.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range auditDraftFindingConversation.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
				}
			}
			if err := adfcs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}
		err = adfcs.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedAuditDraftFindingConversation = append(deletedAuditDraftFindingConversation, map[string]interface{}{
			"id":                  id,
			"comment":             auditDraftFindingConversation.Comment,
			"audit_draft_finding": auditDraftFindingConversation.AuditDraftFinding.Title,
		})
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}
	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for auditDraftFindingConversationID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(auditDraftFindingConversationID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := adfcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
				successfulIds,
				deletedAuditDraftFindingConversation,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (adfcs *AuditDraftFindingConversationService) GetDocumentPath(auditDraftFindingID uint, id uint, token string) (string, error) {
	existingAuditDraftFinding, err := adfcs.auditDraftFindingSvc.GetByID(auditDraftFindingID)
	if err != nil {
		return "", err
	}
	auditID := existingAuditDraftFinding.AuditID

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "draft-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
