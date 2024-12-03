package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditDocumentRequestConversationServiceInterface interface {
	GetAll(auditID uint, auditDocumentRequestID uint, queryParams *dtos.AuditDocumentRequestConversationListQueryParamsDTO) ([]*dtos.AuditDocumentRequestConversationListDTO, error)
	Create(auditID uint, auditDocumentRequestID uint, currentUserID uint, formDTO *dtos.AuditDocumentRequestConversationFormDTO, source string) (uint, error)
	GetDocumentPath(auditId uint, auditDocumentRequestID uint, conversationID uint, token string) (string, error)
}

type AuditDocumentRequestConversationService struct {
	repo                          repositories.AuditDocumentRequestConversationRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	auditSvc                      AuditServiceInterface
	auditDocumentRequestSvc       AuditDocumentRequestServiceInterface
	auditDocumentRequestStatusSvc AuditDocumentRequestStatusServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	notificationSvc               commonServices.NotificationServiceInterface
}

func NewAuditDocumentRequestConversationService(
	repo repositories.AuditDocumentRequestConversationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDocumentRequestSvc AuditDocumentRequestServiceInterface,
	auditDocumentRequestStatusSvc AuditDocumentRequestStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditDocumentRequestConversationServiceInterface {
	return &AuditDocumentRequestConversationService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		auditSvc:                      auditSvc,
		auditDocumentRequestSvc:       auditDocumentRequestSvc,
		auditDocumentRequestStatusSvc: auditDocumentRequestStatusSvc,
		userSvc:                       userSvc,
		notificationSvc:               notificationSvc,
	}
}

func (adrcs *AuditDocumentRequestConversationService) GetAll(auditID uint, auditDocumentRequestID uint, queryParams *dtos.AuditDocumentRequestConversationListQueryParamsDTO) ([]*dtos.AuditDocumentRequestConversationListDTO, error) {
	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return nil, err
	}
	if _, err := adrcs.auditDocumentRequestSvc.GetByID(auditDocumentRequestID); err != nil {
		return nil, err
	}

	return adrcs.repo.GetAll(auditDocumentRequestID, queryParams)
}

func (adrcs *AuditDocumentRequestConversationService) Create(auditID uint, auditDocumentRequestID uint, currentUserID uint, formDTO *dtos.AuditDocumentRequestConversationFormDTO, source string) (uint, error) {
	tx := adrcs.repo.BeginTransaction()

	if _, err := adrcs.auditSvc.GetByID(auditID); err != nil {
		return 0, err
	}
	documentRequest, err := adrcs.auditDocumentRequestSvc.GetByID(auditDocumentRequestID)
	if err != nil {
		return 0, err
	}

	audit, err := adrcs.auditDocumentRequestSvc.GetFullDetailsByID(auditDocumentRequestID)
	if err != nil {
		return 0, err
	}

	if err := adrcs.validateAuditDocumentRequestConversationFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	now := time.Now()

	conversation := &models.AuditDocumentRequestConversation{
		AuditDocumentRequestID: auditDocumentRequestID,
		Comment:                formDTO.Comment,
		CommentedBYID:          &currentUserID,
		CommentedAt:            now,
	}

	id, err := adrcs.repo.Create(conversation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.AuditDocumentRequestConversationDocument
	for _, doc := range formDTO.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "document-requests", strconv.FormatUint(uint64(auditDocumentRequestID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.AuditDocumentRequestConversationDocument{
			AuditDocumentRequestConversationID: id,
			URL:                                url,
			Name:                               doc.FileName,
			Size:                               doc.Size,
			Token:                              encryptedFilename,
			OriginalName:                       doc.OriginalName,
			Type:                               doc.Type,
			CreatedByID:                        currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := adrcs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditID), 10), "document-requests", strconv.FormatUint(uint64(auditDocumentRequestID), 10), "conversations", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	auditDocumentRequestClosedStatus, err := adrcs.auditDocumentRequestStatusSvc.GetByType("closed")
	if err != nil {
		return 0, err
	}
	// Check if IsClose is not nil and its value is true
	if formDTO.IsClose != nil && *formDTO.IsClose {
		if documentRequest.AuditDocumentRequestStatusID != uint(auditDocumentRequestClosedStatus.ID) {
			documentRequest.AuditDocumentRequestStatusID = uint(auditDocumentRequestClosedStatus.ID)
			documentRequest.UpdatedAt = &now
			documentRequest.UpdatedBYID = &currentUserID

			err = adrcs.repo.AuditDocumentRequestStatusUpdate(documentRequest, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	} else {
		// Reopen the status when a new comment is added, if the status is closed
		if documentRequest.AuditDocumentRequestStatusID == uint(auditDocumentRequestClosedStatus.ID) {
			auditDocumentRequestOpenStatus, err := adrcs.auditDocumentRequestStatusSvc.GetByType("open")
			if err != nil {
				return 0, err
			}
			if documentRequest.AuditDocumentRequestStatusID != uint(auditDocumentRequestOpenStatus.ID) {
				documentRequest.AuditDocumentRequestStatusID = uint(auditDocumentRequestOpenStatus.ID)
				documentRequest.UpdatedAt = &now
				documentRequest.UpdatedBYID = &currentUserID

				err = adrcs.repo.AuditDocumentRequestStatusUpdate(documentRequest, tx)
				if err != nil {
					tx.Rollback()
					return 0, err
				}
			}

		}
	}

	// todo -send notifications to users
	tx.Commit()

	go func() {
		err := adrcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDocumentRequestConversation{}.TableName(),
			id,
			conversation,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDocumentRequestConversation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		if source == "audit-insight" {
			if *documentRequest.RequestedBYID == 0 {
				utility.LogError("Responsible user is missing", map[string]interface{}{
					"documentRequestID": id,
				})
				return
			}

			requestedUserID := documentRequest.RequestedBYID
			requestedBy, err := adrcs.userSvc.GetByID(*requestedUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"userID": requestedUserID,
					"error":  err.Error(),
				})
				return
			}

			var notificationHeadingEn, notificationHeadingAr, emailSubject, emailTemplate, action string

			// Content for new comment (not closed)
			action = "new-comment"
			notificationHeadingEn = "New Comment Added to Audit Insight Document Request"
			notificationHeadingAr = "تم إضافة تعليق جديد إلى طلب وثائق رؤية التدقيق"
			emailSubject = fmt.Sprintf("New Audit-Insight Document Request Comment: %s", documentRequest.Title)
			emailTemplate = "audit/templates/email/audit_insight_document_request_conversation.html"

			urlParam := fmt.Sprintf(`{"audit_id":%d,"request_id":%d}`, auditID, id)
			entity := "audit-insight-document-request-conversation"
			appURL := utility.GetAppURL()

			dbNotification := &commonDTOs.DBNotificationDTO{
				UserID:               requestedBy.ID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      audit.Audit.AuditPlan.ReferenceNumber,
				HeadingEn:            notificationHeadingEn,
				HeadingAr:            notificationHeadingAr,
				TitleEn:              documentRequest.Title,
				TitleAr:              documentRequest.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{requestedBy.Email},
				Subject:              emailSubject,
				TemplateRelativePath: emailTemplate,
				EmailData: map[string]interface{}{
					"FirstNameEn":     requestedBy.FirstNameEn,
					"LastNameEn":      requestedBy.LastNameEn,
					"AuditTitle":      audit.Audit.AuditPlan.Title,
					"Title":           documentRequest.Title,
					"Comment":         formDTO.Comment,
					"Description":     documentRequest.Description,
					"ResponsibleUnit": audit.ResponsibleUnit.TitleEn,
					"ResponsibleUser": fmt.Sprintf("%s %s", audit.ResponsibleUser.FirstNameEN, audit.ResponsibleUser.LastNameEN),
					"CommenttedAt":    time.Now().Format("2006-01-02"),
					"RedirectURL":     fmt.Sprintf("%saudit-insight/document-requests", appURL),
				},
			}

			// Send Notifications
			if err := adrcs.notificationSvc.SendNotification(dbNotification, emailNotification, "audit"); err != nil {
				utility.LogError("Failed to send document request notifications", map[string]interface{}{
					"documentRequestID": id,
					"userID":            documentRequest.RequestedByUser.ID,
					"error":             err.Error(),
				})
			}
		} else {
			if documentRequest.ResponsibleUserID == 0 {
				utility.LogError("Responsible user is missing", map[string]interface{}{
					"documentRequestID": id,
				})
				return
			}

			responsibleUserID := documentRequest.ResponsibleUserID
			responsibleUser, err := adrcs.userSvc.GetByID(responsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"userID": responsibleUserID,
					"error":  err.Error(),
				})
				return
			}

			var notificationHeadingEn, notificationHeadingAr, emailSubject, emailTemplate, action string

			if formDTO.IsClose != nil && *formDTO.IsClose {
				// Content for closed status
				action = "status-closed"
				notificationHeadingEn = "Audit Document Request Closed"
				notificationHeadingAr = "تم إغلاق طلب وثائق التدقيق"
				emailSubject = fmt.Sprintf("Audit Document Request Closed: %s", documentRequest.Title)
				emailTemplate = "audit/templates/email/audit_document_request_conversation_closed.html"

			} else {
				// Content for new comment (not closed)
				action = "new-comment"
				notificationHeadingEn = "New Comment Added to Audit Document Request"
				notificationHeadingAr = "تم إضافة تعليق جديد إلى طلب وثائق التدقيق"
				emailSubject = fmt.Sprintf("New Audit Document Request Comment: %s", documentRequest.Title)
				emailTemplate = "audit/templates/email/audit_document_request_conversation.html"
			}

			urlParam := fmt.Sprintf(`{"audit_id":%d,"request_id":%d}`, auditID, id)
			entity := "audit-document-request-conversation"
			appURL := utility.GetAppURL()

			dbNotification := &commonDTOs.DBNotificationDTO{
				UserID:               responsibleUser.ID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      audit.Audit.AuditPlan.ReferenceNumber,
				HeadingEn:            notificationHeadingEn,
				HeadingAr:            notificationHeadingAr,
				TitleEn:              documentRequest.Title,
				TitleAr:              documentRequest.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{responsibleUser.Email},
				Subject:              emailSubject,
				TemplateRelativePath: emailTemplate,
				EmailData: map[string]interface{}{
					"FirstNameEn":     responsibleUser.FirstNameEn,
					"LastNameEn":      responsibleUser.LastNameEn,
					"AuditTitle":      audit.Audit.AuditPlan.Title,
					"Title":           documentRequest.Title,
					"Comment":         formDTO.Comment,
					"Description":     documentRequest.Description,
					"ResponsibleUnit": audit.ResponsibleUnit.TitleEn,
					"RequestedUser":   fmt.Sprintf("%s %s", audit.RequestedByUser.FirstNameEN, audit.RequestedByUser.LastNameEN),
					"CommenttedAt":    time.Now().Format("2006-01-02"),
					"RedirectURL":     fmt.Sprintf("%saudit/audits/%d/document-request", appURL, auditID),
				},
			}

			// Send Notifications
			if err := adrcs.notificationSvc.SendNotification(dbNotification, emailNotification, "audit"); err != nil {
				utility.LogError("Failed to send document request notifications", map[string]interface{}{
					"documentRequestID": id,
					"userID":            documentRequest.ResponsibleUser.ID,
					"error":             err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (adrcs *AuditDocumentRequestConversationService) validateAuditDocumentRequestConversationFormDTO(formDTO *dtos.AuditDocumentRequestConversationFormDTO) error {
	if len(formDTO.Comment) == 0 {
		return errors.New("comment is required")
	}
	if len(formDTO.Comment) > 500 {
		return errors.New("comment exceeds maximum length of 500 characters")
	}
	return nil
}

func (adrcs *AuditDocumentRequestConversationService) GetDocumentPath(auditId uint, auditDocumentRequestID uint, conversationID uint, token string) (string, error) {
	_, err := adrcs.auditSvc.GetByID(auditId)
	if err != nil {
		return "", err
	}
	_, err = adrcs.auditDocumentRequestSvc.GetByID(auditDocumentRequestID)
	if err != nil {
		return "", err
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audits", strconv.FormatUint(uint64(auditId), 10), "document-requests", strconv.FormatUint(uint64(auditDocumentRequestID), 10), "conversations", strconv.FormatUint(uint64(conversationID), 10), "documents")

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
