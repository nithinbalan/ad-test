package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentReviewServiceInterface interface {
	GetComments(id uint) (*documentDTOs.DocumentVersionWorkflowDTO, error)
	AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error
}

type DocumentReviewService struct {
	repo                      repositories.DocumentReviewRepositoryInterface
	logService                logServices.ActivityLogServiceInterface
	documentSvc               DocumentServiceInterface
	documentWorkflowSvc       DocumentWorkflowServiceInterface
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface
	documentStatusSvc         DocumentStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
}

func NewDocumentReviewService(
	repo repositories.DocumentReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	documentWorkflowSvc DocumentWorkflowServiceInterface,
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface,
	documentStatusSvc DocumentStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface) DocumentReviewServiceInterface {
	return &DocumentReviewService{
		repo:                      repo,
		logService:                logService,
		documentSvc:               documentSvc,
		documentWorkflowSvc:       documentWorkflowSvc,
		documentWorkflowStatusSvc: documentWorkflowStatusSvc,
		documentStatusSvc:         documentStatusSvc,
		userSvc:                   userSvc,
		notificationSvc:           notificationSvc,
	}
}

func (drs *DocumentReviewService) GetComments(id uint) (*documentDTOs.DocumentVersionWorkflowDTO, error) {
	var documentVersionWorkflowDTO documentDTOs.DocumentVersionWorkflowDTO

	document, err := drs.documentSvc.GetFullDetailsByID(id, nil)
	if err != nil {
		return nil, err
	}

	if document.DocumentStatus.Type == "draft" {
		return nil, utility.NewForbiddenError("Document must be submitted for review to access comments.")
	}

	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel

	documentVersionWorkflows, err := drs.repo.GetComments(document.LatestDocumentVersion.ID, *currentApprovalLevel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(documentVersionWorkflows, &documentVersionWorkflowDTO)
	if err != nil {
		return nil, err
	}

	return &documentVersionWorkflowDTO, nil
}

func (drs *DocumentReviewService) AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error {
	document, err := drs.documentSvc.GetFullDetailsByID(id, nil)
	if err != nil {
		return err
	}

	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel
	currentDocumentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
	if err != nil {
		return err
	}

	if documentReviewCommentFormDTO.IsApproved {
		return drs.handleApprovedCase(id, currentUserID, documentReviewCommentFormDTO, document, currentDocumentWorkflow, *currentApprovalLevel)
	} else if documentReviewCommentFormDTO.IsRejected {
		return drs.handleRejectedCase(id, currentUserID, documentReviewCommentFormDTO, currentDocumentWorkflow)
	} else if documentReviewCommentFormDTO.IsReverted {
		return drs.handleRevertedCase(currentUserID, documentReviewCommentFormDTO, currentDocumentWorkflow)
	} else {
		return drs.handleCommentOnlyCase(currentUserID, documentReviewCommentFormDTO, currentDocumentWorkflow)
	}
}

func (drs *DocumentReviewService) handleApprovedCase(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, document *documentDTOs.DocumentDetailsDTO, currentDocumentWorkflow *models.DocumentWorkflow, currentApprovalLevel int) error {
	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}
	switch currentDocumentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentDocumentWorkflow.DocumentWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}

		documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
		if err != nil {
			return err
		}

		documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
		err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
		if err != nil {
			return err
		}

		drs.sendApprovalNotification(id, currentUserID, document)

		nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, currentApprovalLevel)
		if err != nil {
			return err
		}

		if nextDocWorkflowLevel == nil {
			documentStatus, err := drs.documentStatusSvc.GetByType("approved")
			if err != nil {
				return err
			}
			existingDocument, err := drs.documentSvc.GetByID(id)
			if err != nil {
				return err
			}

			existingDocument.DocumentStatusID = documentStatus.ID
			err = drs.documentSvc.UpdateDocument(existingDocument)
			if err != nil {
				return err
			}
		} else {
			existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
			if err != nil {
				return err
			}

			existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
			err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
			if err != nil {
				return err
			}

			drs.notifyNextLevelUsers(id, nextDocWorkflowLevel, currentUserID, document)
		}
	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentDocumentWorkflow.DocumentWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}

		if !allowed {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}

		drs.sendApprovalNotification(id, currentUserID, document)

		if currentDocumentWorkflow.ApprovalTypes == "all" {
			allApproved, err := drs.documentWorkflowSvc.CheckAllParticipantsApproved(currentDocumentWorkflow.ID, currentDocumentWorkflow.DocumentWorkflowParticipants)
			if err != nil {
				return err
			}

			if !allApproved {
				return nil
			} else {
				documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
				if err != nil {
					return err
				}

				documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
				err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
				if err != nil {
					return err
				}

				nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, currentApprovalLevel)
				if err != nil {
					return err
				}

				if nextDocWorkflowLevel == nil {
					documentStatus, err := drs.documentStatusSvc.GetByType("approved")
					if err != nil {
						return err
					}
					existingDocument, err := drs.documentSvc.GetByID(id)
					if err != nil {
						return err
					}

					existingDocument.DocumentStatusID = documentStatus.ID
					err = drs.documentSvc.UpdateDocument(existingDocument)
					if err != nil {
						return err
					}
				} else {
					existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
					if err != nil {
						return err
					}

					existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
					err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
					if err != nil {
						return err
					}

					drs.notifyNextLevelUsers(id, nextDocWorkflowLevel, currentUserID, document)
				}
			}
		} else {
			documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
			if err != nil {
				return err
			}

			documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
			err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
			if err != nil {
				return err
			}

			nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, currentApprovalLevel)
			if err != nil {
				return err
			}

			if nextDocWorkflowLevel == nil {
				documentStatus, err := drs.documentStatusSvc.GetByType("approved")
				if err != nil {
					return err
				}
				existingDocument, err := drs.documentSvc.GetByID(id)
				if err != nil {
					return err
				}

				existingDocument.DocumentStatusID = documentStatus.ID
				err = drs.documentSvc.UpdateDocument(existingDocument)
				if err != nil {
					return err
				}
			} else {
				existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
				if err != nil {
					return err
				}

				existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
				err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
				if err != nil {
					return err
				}

				drs.notifyNextLevelUsers(id, nextDocWorkflowLevel, currentUserID, document)
			}
		}
	}

	return nil
}

func (drs *DocumentReviewService) handleRejectedCase(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, currentDocumentWorkflow *models.DocumentWorkflow) error {
	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	documentWorkflowReview := &models.DocumentWorkflowReview{
		DocumentWorkflowID:       currentDocumentWorkflow.ID,
		DocumentWorkflowStatusID: documentWorkflowStatus.ID,
		Comment:                  documentReviewCommentFormDTO.Comment,
		CreatedBYID:              currentUserID,
	}

	err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
	if err != nil {
		return err
	}

	documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
	if err != nil {
		return err
	}

	documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
	err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
	if err != nil {
		return err
	}

	documentStatus, err := drs.documentStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}
	existingDocument, err := drs.documentSvc.GetByID(id)
	if err != nil {
		return err
	}

	existingDocument.DocumentStatusID = documentStatus.ID
	err = drs.documentSvc.UpdateDocument(existingDocument)
	if err != nil {
		return err
	}

	return nil
}

func (drs *DocumentReviewService) handleRevertedCase(currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, currentDocumentWorkflow *models.DocumentWorkflow) error {
	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}
	documentWorkflowReview := &models.DocumentWorkflowReview{
		DocumentWorkflowID:       currentDocumentWorkflow.ID,
		DocumentWorkflowStatusID: documentWorkflowStatus.ID,
		Comment:                  documentReviewCommentFormDTO.Comment,
		CreatedBYID:              currentUserID,
	}

	err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
	if err != nil {
		return err
	}

	return nil
}

func (drs *DocumentReviewService) handleCommentOnlyCase(currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, currentDocumentWorkflow *models.DocumentWorkflow) error {
	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("comment")
	if err != nil {
		return err
	}
	documentWorkflowReview := &models.DocumentWorkflowReview{
		DocumentWorkflowID:       currentDocumentWorkflow.ID,
		DocumentWorkflowStatusID: documentWorkflowStatus.ID,
		Comment:                  documentReviewCommentFormDTO.Comment,
		CreatedBYID:              currentUserID,
	}

	err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
	if err != nil {
		return err
	}

	return nil
}

func (drs *DocumentReviewService) sendApprovalNotification(id uint, currentUserID uint, document *documentDTOs.DocumentDetailsDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "document"
		action := "approval"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(document.ResponsibleUser.ID),
			ModuleIdentifier:     "document",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Document Approval",
			HeadingAr:            "مراجعة الوثيقة",
			TitleEn:              document.TitleEn,
			TitleAr:              document.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		approvedByUser, err := drs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{document.ResponsibleUser.Email}
		subject := fmt.Sprintf("Document/%s: Approval Notification", document.TitleEn)
		urlParam = fmt.Sprintf("%s/documents/%d", appURL, id)

		emailData := map[string]interface{}{
			"FirstNameEn":           document.ResponsibleUser.FirstNameEN,
			"LastNameEn":            document.ResponsibleUser.LastNameEN,
			"DocumentTitle":         document.TitleEn,
			"DocumentType":          document.DocumentType.TitleEn,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "document/templates/email/document_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (drs *DocumentReviewService) notifyNextLevelUsers(id uint, nextDocWorkflowLevel *models.DocumentWorkflow, currentUserID uint, document *documentDTOs.DocumentDetailsDTO) {
	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	entity := "document"
	action := "review-user-assign"

	for _, participant := range nextDocWorkflowLevel.DocumentWorkflowParticipants {
		responsibleUser, err := drs.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		// Prepare database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               participant.UserID,
			ModuleIdentifier:     "document",
			NotificationType:     string(commonDTOs.ReviewType),
			HeadingEn:            "Document Review",
			HeadingAr:            "مراجعة الوثيقة",
			TitleEn:              document.TitleEn,
			TitleAr:              document.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             fmt.Sprintf(`{"id":%d}`, id),
			IsSystemNotification: false,
			CreatedByID:          &currentUserID, // Assuming CreatedByID is not needed for next-level notifications
		}

		// Prepare email notification DTO
		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Document/%s: Review Assigned", document.TitleEn)

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"DocumentTitle":         document.TitleEn,
			"DocumentType":          document.DocumentType.TitleEn,
			"AssignedByFirstNameEn": document.ResponsibleUser.FirstNameEN,
			"AssignedByLastNameEn":  document.ResponsibleUser.LastNameEN,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "document/templates/email/document_workflow_review_user_mail.html",
			EmailData:            emailData,
		}

		// Send both database and email notifications
		err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
		if err != nil {
			utility.LogError("Failed to send next level review notification", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}
