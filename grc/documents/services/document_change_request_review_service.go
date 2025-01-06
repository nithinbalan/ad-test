package services

import (
	"encoding/json"
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestReviewServiceInterface interface {
	GetComments(id uint) (*documentDTOs.DocumentChangeRequestWorkflowDTO, error)
	AddComment(id uint, currentUserID uint, documentCRReviewCommentFormDTO *documentDTOs.DocumentCRReviewCommentFormDTO) error
}

type DocumentChangeRequestReviewService struct {
	repo                           repositories.DocumentChangeRequestReviewRepositoryInterface
	logService                     logServices.ActivityLogServiceInterface
	documentSvc                    DocumentServiceInterface
	documentChangeRequestSvc       DocumentChangeRequestServiceInterface
	documentCRWorkflowSvc          DocumentChangeRequestWorkflowServiceInterface
	documentCRWorkflowStatusSvc    DocumentChangeRequestWorkflowStatusServiceInterface
	documentChangeRequestStatusSvc DocumentChangeRequestStatusServiceInterface
	userSvc                        organizationSvc.UserServiceInterface
	notificationSvc                commonServices.NotificationServiceInterface
	taskSvc                        taskServices.TaskServiceInterface
}

func NewDocumentChangeRequestReviewService(
	repo repositories.DocumentChangeRequestReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	documentChangeRequestSvc DocumentChangeRequestServiceInterface,
	documentCRWorkflowSvc DocumentChangeRequestWorkflowServiceInterface,
	documentCRWorkflowStatusSvc DocumentChangeRequestWorkflowStatusServiceInterface,
	documentChangeRequestStatusSvc DocumentChangeRequestStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface) DocumentChangeRequestReviewServiceInterface {
	return &DocumentChangeRequestReviewService{
		repo:                           repo,
		logService:                     logService,
		documentSvc:                    documentSvc,
		documentChangeRequestSvc:       documentChangeRequestSvc,
		documentCRWorkflowSvc:          documentCRWorkflowSvc,
		documentCRWorkflowStatusSvc:    documentCRWorkflowStatusSvc,
		documentChangeRequestStatusSvc: documentChangeRequestStatusSvc,
		userSvc:                        userSvc,
		notificationSvc:                notificationSvc,
		taskSvc:                        taskSvc,
	}
}

func (dcrrs *DocumentChangeRequestReviewService) GetComments(id uint) (*documentDTOs.DocumentChangeRequestWorkflowDTO, error) {
	var documentChangeRequestWorkflowDTO documentDTOs.DocumentChangeRequestWorkflowDTO

	documentCRFullDetails, err := dcrrs.documentChangeRequestSvc.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	if documentCRFullDetails.DocumentChangeRequestStatus.Type == "draft" {
		return nil, utility.NewForbiddenError("Document must be submitted for review to access comments.")
	}

	currentApprovalLevel := documentCRFullDetails.NextReviewLevel

	documentCRWorkflows, err := dcrrs.repo.GetComments(id, *currentApprovalLevel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(documentCRWorkflows, &documentChangeRequestWorkflowDTO)
	if err != nil {
		return nil, err
	}

	return &documentChangeRequestWorkflowDTO, nil
}

func (dcrrs *DocumentChangeRequestReviewService) AddComment(id uint, currentUserID uint, documentCRReviewCommentFormDTO *documentDTOs.DocumentCRReviewCommentFormDTO) error {
	documentChangeRequest, err := dcrrs.documentChangeRequestSvc.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	currentApprovalLevel := documentChangeRequest.NextReviewLevel
	currentDocumentCRWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByLevel(id, *currentApprovalLevel)
	if err != nil {
		return err
	}

	if documentCRReviewCommentFormDTO.IsApproved {
		return dcrrs.handleApprovedCase(id, currentUserID, documentCRReviewCommentFormDTO, documentChangeRequest, currentDocumentCRWorkflow, *currentApprovalLevel)
	} else if documentCRReviewCommentFormDTO.IsRejected {
		return dcrrs.handleRejectedCase(id, currentUserID, documentCRReviewCommentFormDTO, documentChangeRequest, currentDocumentCRWorkflow)
	} else if documentCRReviewCommentFormDTO.IsReverted {
		return dcrrs.handleRevertedCase(id, currentUserID, documentChangeRequest, documentCRReviewCommentFormDTO, currentDocumentCRWorkflow)
	} else {
		return dcrrs.handleCommentOnlyCase(id, currentUserID, documentChangeRequest, documentCRReviewCommentFormDTO, currentDocumentCRWorkflow)
	}
}

func (dcrrs *DocumentChangeRequestReviewService) handleApprovedCase(id uint, currentUserID uint, documentCRReviewCommentFormDTO *documentDTOs.DocumentCRReviewCommentFormDTO, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO, currentDocumentCRWorkflow *models.DocumentChangeRequestWorkflow, currentApprovalLevel int) error {
	documentCRWorkflowStatus, err := dcrrs.documentCRWorkflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}
	switch currentDocumentCRWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentDocumentCRWorkflow.DocumentChangeRequestWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		documentCRWorkflowReview := &models.DocumentChangeRequestWorkflowReview{
			DocumentChangeRequestWorkflowID:       currentDocumentCRWorkflow.ID,
			DocumentChangeRequestWorkflowStatusID: documentCRWorkflowStatus.ID,
			Comment:                               documentCRReviewCommentFormDTO.Comment,
			CreatedBYID:                           currentUserID,
		}

		err = dcrrs.repo.CreateDocumentWorkflowReview(documentCRWorkflowReview)
		if err != nil {
			return err
		}

		documentCRWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByID(currentDocumentCRWorkflow.ID)
		if err != nil {
			return err
		}

		documentCRWorkflow.DocumentChangeRequestWorkflowStatusID = documentCRWorkflowStatus.ID
		err = dcrrs.documentCRWorkflowSvc.UpdateDocumentCRWorkflow(documentCRWorkflow)
		if err != nil {
			return err
		}

		_, err = dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, currentUserID, currentDocumentCRWorkflow.ID)
		if err != nil {
			return err
		}

		dcrrs.sendApprovalNotification(id, currentUserID, documentChangeRequest)

		nextDocCRWorkflowLevel, err := dcrrs.documentCRWorkflowSvc.GetNextLevel(id, currentApprovalLevel)
		if err != nil {
			return err
		}

		if nextDocCRWorkflowLevel == nil {
			documentCRStatus, err := dcrrs.documentChangeRequestStatusSvc.GetByType("approved")
			if err != nil {
				return err
			}
			existingDocumentCR, err := dcrrs.documentChangeRequestSvc.GetByID(id)
			if err != nil {
				return err
			}

			existingDocumentCR.DocumentChangeRequestStatusID = documentCRStatus.ID
			err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingDocumentCR)
			if err != nil {
				return err
			}
		} else {
			existingDocumentCR, err := dcrrs.documentChangeRequestSvc.GetByID(id)
			if err != nil {
				return err
			}
			existingDocumentCR.NextReviewLevel = &nextDocCRWorkflowLevel.Level
			err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingDocumentCR)
			if err != nil {
				return err
			}
			nextWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByLevel(id, nextDocCRWorkflowLevel.Level)
			if err != nil {
				return err
			}
			err = dcrrs.CreateTasksForNextLevelUsers(id, nextWorkflow, documentChangeRequest)
			if err != nil {
				return err
			}

			dcrrs.notifyNextLevelUsers(id, nextWorkflow, documentChangeRequest)
		}
	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentDocumentCRWorkflow.DocumentChangeRequestWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}

		if !allowed {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		documentCRWorkflowReview := &models.DocumentChangeRequestWorkflowReview{
			DocumentChangeRequestWorkflowID:       currentDocumentCRWorkflow.ID,
			DocumentChangeRequestWorkflowStatusID: documentCRWorkflowStatus.ID,
			Comment:                               documentCRReviewCommentFormDTO.Comment,
			CreatedBYID:                           currentUserID,
		}

		err = dcrrs.repo.CreateDocumentWorkflowReview(documentCRWorkflowReview)
		if err != nil {
			return err
		}

		dcrrs.sendApprovalNotification(id, currentUserID, documentChangeRequest)

		if currentDocumentCRWorkflow.ApprovalTypes == "all" {
			_, err := dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, currentUserID, currentDocumentCRWorkflow.ID)
			if err != nil {
				return err
			}
			allApproved, err := dcrrs.documentCRWorkflowSvc.CheckAllParticipantsApproved(currentDocumentCRWorkflow.ID, currentDocumentCRWorkflow.DocumentChangeRequestWorkflowParticipants)
			if err != nil {
				return err
			}

			if !allApproved {
				return nil
			} else {
				documentCRWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByID(currentDocumentCRWorkflow.ID)
				if err != nil {
					return err
				}

				documentCRWorkflow.DocumentChangeRequestWorkflowStatusID = documentCRWorkflowStatus.ID
				err = dcrrs.documentCRWorkflowSvc.UpdateDocumentCRWorkflow(documentCRWorkflow)
				if err != nil {
					return err
				}

				nextDocCRWorkflowLevel, err := dcrrs.documentCRWorkflowSvc.GetNextLevel(id, currentApprovalLevel)
				if err != nil {
					return err
				}

				if nextDocCRWorkflowLevel == nil {
					documentChangeRequestStatus, err := dcrrs.documentChangeRequestStatusSvc.GetByType("approved")
					if err != nil {
						return err
					}
					existingDocumentCR, err := dcrrs.documentChangeRequestSvc.GetByID(id)
					if err != nil {
						return err
					}

					existingDocumentCR.DocumentChangeRequestStatusID = documentChangeRequestStatus.ID
					err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingDocumentCR)
					if err != nil {
						return err
					}
				} else {
					existingDocumentCR, err := dcrrs.documentChangeRequestSvc.GetByID(id)
					if err != nil {
						return err
					}

					existingDocumentCR.NextReviewLevel = &nextDocCRWorkflowLevel.Level
					err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingDocumentCR)
					if err != nil {
						return err
					}
					nextWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByLevel(id, nextDocCRWorkflowLevel.Level)
					if err != nil {
						return err
					}
					err = dcrrs.CreateTasksForNextLevelUsers(id, nextWorkflow, documentChangeRequest)
					if err != nil {
						return err
					}
					dcrrs.notifyNextLevelUsers(id, nextWorkflow, documentChangeRequest)
				}
			}
		} else {
			documentCRWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByID(currentDocumentCRWorkflow.ID)
			if err != nil {
				return err
			}

			documentCRWorkflow.DocumentChangeRequestWorkflowStatusID = documentCRWorkflowStatus.ID
			err = dcrrs.documentCRWorkflowSvc.UpdateDocumentCRWorkflow(documentCRWorkflow)
			if err != nil {
				return err
			}

			task, err := dcrrs.taskSvc.UpdateDocumentTaskStatus(id, currentUserID, currentDocumentCRWorkflow.ID)
			if err != nil {
				return err
			}

			// Parse the related_task_ids from the task
			var relatedTaskIDs []uint
			if task.RelatedTaskIDs != nil {
				if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
					return err
				}

				// Close each related task
				for _, relatedTaskID := range relatedTaskIDs {
					relatedTask, err := dcrrs.taskSvc.GetByID(relatedTaskID)
					if err != nil {
						return err
					}
					_, err = dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, relatedTask.ResponsibleUserID, currentDocumentCRWorkflow.ID)
					if err != nil {
						return err
					}
				}
			}

			nextDocCRWorkflowLevel, err := dcrrs.documentCRWorkflowSvc.GetNextLevel(id, currentApprovalLevel)
			if err != nil {
				return err
			}

			if nextDocCRWorkflowLevel == nil {
				documentChangeRequestStatus, err := dcrrs.documentChangeRequestStatusSvc.GetByType("approved")
				if err != nil {
					return err
				}
				existingDocumentCR, err := dcrrs.documentChangeRequestSvc.GetByID(id)
				if err != nil {
					return err
				}

				existingDocumentCR.DocumentChangeRequestStatusID = documentChangeRequestStatus.ID
				err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingDocumentCR)
				if err != nil {
					return err
				}
			} else {
				existingDocumentCR, err := dcrrs.documentChangeRequestSvc.GetByID(id)
				if err != nil {
					return err
				}

				existingDocumentCR.NextReviewLevel = &nextDocCRWorkflowLevel.Level
				err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingDocumentCR)
				if err != nil {
					return err
				}

				nextWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByLevel(id, nextDocCRWorkflowLevel.Level)
				if err != nil {
					return err
				}

				err = dcrrs.CreateTasksForNextLevelUsers(id, nextWorkflow, documentChangeRequest)
				if err != nil {
					return err
				}

				dcrrs.notifyNextLevelUsers(id, nextWorkflow, documentChangeRequest)
			}
		}
	}

	return nil
}

func (dcrrs *DocumentChangeRequestReviewService) handleRejectedCase(id uint, currentUserID uint, documentCRReviewCommentFormDTO *documentDTOs.DocumentCRReviewCommentFormDTO, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO, currentDocumentCRWorkflow *models.DocumentChangeRequestWorkflow) error {
	documentCRWorkflowStatus, err := dcrrs.documentCRWorkflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	documentCRWorkflowReview := &models.DocumentChangeRequestWorkflowReview{
		DocumentChangeRequestWorkflowID:       currentDocumentCRWorkflow.ID,
		DocumentChangeRequestWorkflowStatusID: documentCRWorkflowStatus.ID,
		Comment:                               documentCRReviewCommentFormDTO.Comment,
		CreatedBYID:                           currentUserID,
	}

	err = dcrrs.repo.CreateDocumentWorkflowReview(documentCRWorkflowReview)
	if err != nil {
		return err
	}

	documentCRWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByID(currentDocumentCRWorkflow.ID)
	if err != nil {
		return err
	}

	documentCRWorkflow.DocumentChangeRequestWorkflowStatusID = documentCRWorkflowStatus.ID
	err = dcrrs.documentCRWorkflowSvc.UpdateDocumentCRWorkflow(documentCRWorkflow)
	if err != nil {
		return err
	}

	documentChangeRequestStatus, err := dcrrs.documentChangeRequestStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}
	existingCRDocument, err := dcrrs.documentChangeRequestSvc.GetByID(id)
	if err != nil {
		return err
	}

	existingCRDocument.DocumentChangeRequestStatusID = documentChangeRequestStatus.ID
	err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingCRDocument)
	if err != nil {
		return err
	}

	task, err := dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, currentUserID, currentDocumentCRWorkflow.ID)
	if err != nil {
		return err
	}

	// Parse the related_task_ids from the task
	var relatedTaskIDs []uint
	if task.RelatedTaskIDs != nil {
		if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
			return err
		}

		// Close each related task
		for _, relatedTaskID := range relatedTaskIDs {
			relatedTask, err := dcrrs.taskSvc.GetByID(relatedTaskID)
			if err != nil {
				return err
			}
			_, err = dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, relatedTask.ResponsibleUserID, currentDocumentCRWorkflow.ID)
			if err != nil {
				return err
			}
		}
	}

	dcrrs.sendRejectedNotification(id, currentUserID, documentChangeRequest)

	return nil
}

func (dcrrs *DocumentChangeRequestReviewService) handleRevertedCase(id, currentUserID uint, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO, documentCRReviewCommentFormDTO *documentDTOs.DocumentCRReviewCommentFormDTO, currentDocumentCRWorkflow *models.DocumentChangeRequestWorkflow) error {
	documentCRWorkflowStatus, err := dcrrs.documentCRWorkflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}
	documentCRWorkflowReview := &models.DocumentChangeRequestWorkflowReview{
		DocumentChangeRequestWorkflowID:       currentDocumentCRWorkflow.ID,
		DocumentChangeRequestWorkflowStatusID: documentCRWorkflowStatus.ID,
		Comment:                               documentCRReviewCommentFormDTO.Comment,
		CreatedBYID:                           currentUserID,
	}

	err = dcrrs.repo.CreateDocumentWorkflowReview(documentCRWorkflowReview)
	if err != nil {
		return err
	}

	documentCRWorkflow, err := dcrrs.documentCRWorkflowSvc.GetDocumentCRWorkflowByID(currentDocumentCRWorkflow.ID)
	if err != nil {
		return err
	}

	documentCRWorkflow.DocumentChangeRequestWorkflowStatusID = documentCRWorkflowStatus.ID
	err = dcrrs.documentCRWorkflowSvc.UpdateDocumentCRWorkflow(documentCRWorkflow)
	if err != nil {
		return err
	}

	documentChangeRequestStatus, err := dcrrs.documentChangeRequestStatusSvc.GetByType("reverted")
	if err != nil {
		return err
	}
	existingCRDocument, err := dcrrs.documentChangeRequestSvc.GetByID(id)
	if err != nil {
		return err
	}

	existingCRDocument.DocumentChangeRequestStatusID = documentChangeRequestStatus.ID
	err = dcrrs.documentChangeRequestSvc.UpdateDocumentChangeRequest(existingCRDocument)
	if err != nil {
		return err
	}

	task, err := dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, currentUserID, currentDocumentCRWorkflow.ID)
	if err != nil {
		return err
	}

	// Parse the related_task_ids from the task
	var relatedTaskIDs []uint
	if task.RelatedTaskIDs != nil {
		if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
			return err
		}

		// Close each related task
		for _, relatedTaskID := range relatedTaskIDs {
			relatedTask, err := dcrrs.taskSvc.GetByID(relatedTaskID)
			if err != nil {
				return err
			}
			_, err = dcrrs.taskSvc.UpdateDocumentChangeRequestTaskStatus(id, relatedTask.ResponsibleUserID, currentDocumentCRWorkflow.ID)
			if err != nil {
				return err
			}
		}
	}

	dcrrs.sendRevertedNotification(id, currentUserID, documentChangeRequest)

	return nil
}

func (dcrrs *DocumentChangeRequestReviewService) handleCommentOnlyCase(id, currentUserID uint, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO, documentCRReviewCommentFormDTO *documentDTOs.DocumentCRReviewCommentFormDTO, currentDocumentCRWorkflow *models.DocumentChangeRequestWorkflow) error {
	documentCRWorkflowStatus, err := dcrrs.documentCRWorkflowStatusSvc.GetByType("comment")
	if err != nil {
		return err
	}
	documentCRWorkflowReview := &models.DocumentChangeRequestWorkflowReview{
		DocumentChangeRequestWorkflowID:       currentDocumentCRWorkflow.ID,
		DocumentChangeRequestWorkflowStatusID: documentCRWorkflowStatus.ID,
		Comment:                               documentCRReviewCommentFormDTO.Comment,
		CreatedBYID:                           currentUserID,
	}

	err = dcrrs.repo.CreateDocumentWorkflowReview(documentCRWorkflowReview)
	if err != nil {
		return err
	}

	dcrrs.sendCommentNotification(id, currentUserID, documentCRReviewCommentFormDTO.Comment, documentChangeRequest)

	return nil
}

func (dcrrs *DocumentChangeRequestReviewService) sendApprovalNotification(id uint, currentUserID uint, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "document-change-request"
		action := "approval"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(documentChangeRequest.Document.ResponsibleUser.ID),
			ModuleIdentifier:     "document",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Document change Request Approval",
			HeadingAr:            "الموافقة على طلب تغيير المستند",
			ReferenceNumber:      documentChangeRequest.ReferenceNumber,
			TitleEn:              documentChangeRequest.Document.TitleEn,
			TitleAr:              documentChangeRequest.Document.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		approvedByUser, err := dcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{documentChangeRequest.Document.ResponsibleUser.Email}
		subject := fmt.Sprintf("Document Change Request/%s: Approval Notification", documentChangeRequest.Document.TitleEn)
		urlParam = fmt.Sprintf("%s/document/change-request/%d", appURL, id)

		emailData := map[string]interface{}{
			"FirstNameEn":           documentChangeRequest.Document.ResponsibleUser.FirstNameEN,
			"LastNameEn":            documentChangeRequest.Document.ResponsibleUser.LastNameEN,
			"DocumentTitle":         documentChangeRequest.Document.TitleEn,
			"DocumentType":          documentChangeRequest.Document.DocumentType.TitleEn,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "document/templates/email/document_change_request_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = dcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (dcrrs *DocumentChangeRequestReviewService) notifyNextLevelUsers(id uint, nextDocWorkflowLevel *models.DocumentChangeRequestWorkflow, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO) {
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "document-change-request"
		action := "review-user-assign"

		for _, participant := range nextDocWorkflowLevel.DocumentChangeRequestWorkflowParticipants {
			responsibleUser, err := dcrrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(documentChangeRequest.Document.ResponsibleUser.ID)
			// Prepare database notification DTO
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               participant.UserID,
				ModuleIdentifier:     "document",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Document Change Request Review",
				HeadingAr:            "مراجعة طلب تغيير المستند",
				ReferenceNumber:      documentChangeRequest.ReferenceNumber,
				TitleEn:              documentChangeRequest.Document.TitleEn,
				TitleAr:              documentChangeRequest.Document.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			// Prepare email notification DTO
			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Document Change Request/%s: Review Assigned", documentChangeRequest.Document.TitleEn)
			appURL := utility.GetAppURL()

			emailData := map[string]interface{}{
				"FirstNameEn":               responsibleUser.FirstNameEn,
				"LastNameEn":                responsibleUser.LastNameEn,
				"DocumentTitle":             documentChangeRequest.Document.TitleEn,
				"DocumentType":              documentChangeRequest.Document.DocumentType.TitleEn,
				"DocumentChangeRequestType": documentChangeRequest.DocumentChangeRequestType.TitleEn,
				"ReasonForChange":           documentChangeRequest.ReasonForChange,
				"AssignedByFirstNameEn":     documentChangeRequest.Document.ResponsibleUser.FirstNameEN,
				"AssignedByLastNameEn":      documentChangeRequest.Document.ResponsibleUser.LastNameEN,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%s/document/change-request/%d", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "document/templates/email/document_change_request_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send both database and email notifications
			err = dcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (dcrrs *DocumentChangeRequestReviewService) CreateTasksForNextLevelUsers(changeRequestID uint, workflowLevel *models.DocumentChangeRequestWorkflow, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO) error {
	urlParam := fmt.Sprintf(`{"id":%d}`, changeRequestID)
	var taskIDs []uint

	createdByID := uint(documentChangeRequest.Document.ResponsibleUser.ID)
	// Loop through each participant and create a task
	for _, participant := range workflowLevel.DocumentChangeRequestWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:                "document",
			DocumentChangeRequestID:         &changeRequestID,
			DocumentChangeRequestWorkflowID: &workflowLevel.ID,
			// Title:              documentChangeRequest.Document.TitleEn,
			// Description:        documentChangeRequest.Document.DescriptionEn,
			URLParam:          &urlParam,
			QueryParam:        nil,
			ResponsibleUserID: participant.User.ID,
		}

		taskID, err := dcrrs.taskSvc.Create(taskDTO, createdByID)
		if err != nil {
			return err
		}
		taskIDs = append(taskIDs, taskID)
	}

	// Set up related task IDs for each task
	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			return err
		}

		err = dcrrs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	return nil
}

func (dcrrs *DocumentChangeRequestReviewService) sendRejectedNotification(id uint, currentUserID uint, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "document-change-request"
		action := "rejected"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(documentChangeRequest.Document.ResponsibleUser.ID),
			ModuleIdentifier:     "document",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Document Change Request Rejected",
			HeadingAr:            "تم رفض الوثيقة",
			ReferenceNumber:      documentChangeRequest.Document.ReferenceNumber,
			TitleEn:              documentChangeRequest.Document.TitleEn,
			TitleAr:              documentChangeRequest.Document.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		rejectedByUser, err := dcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{documentChangeRequest.Document.ResponsibleUser.Email}
		subject := fmt.Sprintf("Document Change Request/%s: Rejected Notification", documentChangeRequest.Document.TitleEn)
		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

		emailData := map[string]interface{}{
			"FirstNameEn":               documentChangeRequest.Document.ResponsibleUser.FirstNameEN,
			"LastNameEn":                documentChangeRequest.Document.ResponsibleUser.LastNameEN,
			"DocumentTitle":             documentChangeRequest.Document.TitleEn,
			"DocumentType":              documentChangeRequest.Document.DocumentType.TitleEn,
			"DocumentChangeRequestType": documentChangeRequest.DocumentChangeRequestType.TitleEn,
			"RejectedByFirstNameEn":     rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":      rejectedByUser.LastNameEn,
			"ApprovedAt":                time.Now().Format("2006-01-02"),
			"RedirectURL":               urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "document/templates/email/document_change_request_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = dcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (dcrrs *DocumentChangeRequestReviewService) sendRevertedNotification(id uint, currentUserID uint, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "document-change-request"
		action := "reverted"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(documentChangeRequest.Document.ResponsibleUser.ID),
			ModuleIdentifier:     "document",
			NotificationType:     string(commonDTOs.RevertType),
			HeadingEn:            "Document Change Request Reverted",
			HeadingAr:            "تم التراجع عن طلب تغيير المستند",
			ReferenceNumber:      documentChangeRequest.ReferenceNumber,
			TitleEn:              documentChangeRequest.Document.TitleEn,
			TitleAr:              documentChangeRequest.Document.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the user who performed the reversion
		revertedByUser, err := dcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch reversion performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{documentChangeRequest.Document.ResponsibleUser.Email}
		subject := fmt.Sprintf("Document Change Request/%s: Reverted Notification", documentChangeRequest.Document.TitleEn)
		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

		emailData := map[string]interface{}{
			"FirstNameEn":               documentChangeRequest.Document.ResponsibleUser.FirstNameEN,
			"LastNameEn":                documentChangeRequest.Document.ResponsibleUser.LastNameEN,
			"DocumentTitle":             documentChangeRequest.Document.TitleEn,
			"DocumentType":              documentChangeRequest.Document.DocumentType.TitleEn,
			"DocumentChangeRequestType": documentChangeRequest.DocumentChangeRequestType.TitleEn,
			"RevertedByFirstNameEn":     revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":      revertedByUser.LastNameEn,
			"RevertedAt":                time.Now().Format("2006-01-02"),
			"RedirectURL":               urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "document/templates/email/document_change_request_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = dcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
		if err != nil {
			utility.LogError("Failed to send reverted notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (dcrrs *DocumentChangeRequestReviewService) sendCommentNotification(id uint, currentUserID uint, commentText string, documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "document-change-request"
		action := "comment"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(documentChangeRequest.Document.ResponsibleUser.ID),
			ModuleIdentifier:     "document",
			NotificationType:     string(commonDTOs.CommentType),
			HeadingEn:            "New Document Change Request Comment",
			HeadingAr:            "تعليق جديد على طلب تغيير المستند",
			ReferenceNumber:      documentChangeRequest.ReferenceNumber,
			TitleEn:              documentChangeRequest.Document.TitleEn,
			TitleAr:              documentChangeRequest.Document.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the user who made the comment
		commentByUser, err := dcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch user who commented", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{documentChangeRequest.Document.ResponsibleUser.Email}
		subject := fmt.Sprintf("Document/%s: New Comment Notification", documentChangeRequest.Document.TitleEn)
		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

		emailData := map[string]interface{}{
			"FirstNameEn":               documentChangeRequest.Document.ResponsibleUser.FirstNameEN,
			"LastNameEn":                documentChangeRequest.Document.ResponsibleUser.LastNameEN,
			"DocumentTitle":             documentChangeRequest.Document.TitleEn,
			"DocumentType":              documentChangeRequest.Document.DocumentType.TitleEn,
			"DocumentChangeRequestType": documentChangeRequest.DocumentChangeRequestType.TitleEn,
			"CommentByFirstNameEn":      commentByUser.FirstNameEn,
			"CommentByLastNameEn":       commentByUser.LastNameEn,
			"CommentText":               commentText,
			"CommentedAt":               time.Now().Format("2006-01-02"),
			"RedirectURL":               urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "document/templates/email/document_change_request_workflow_comment_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = dcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
		if err != nil {
			utility.LogError("Failed to send comment notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
