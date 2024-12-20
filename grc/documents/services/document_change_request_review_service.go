package services

import (
	"errors"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	organizationSvc "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestReviewServiceInterface interface {
	GetComments(id uint) (*documentDTOs.DocumentChangeRequestWorkflowDTO, error)
	//AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error
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

func (drs *DocumentChangeRequestReviewService) GetComments(id uint) (*documentDTOs.DocumentChangeRequestWorkflowDTO, error) {
	var documentChangeRequestWorkflowDTO documentDTOs.DocumentChangeRequestWorkflowDTO

	documentCRFullDetails, err := drs.documentChangeRequestSvc.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	if documentCRFullDetails.DocumentChangeRequestStatus.Type == "draft" {
		return nil, utility.NewForbiddenError("Document must be submitted for review to access comments.")
	}

	currentApprovalLevel := documentCRFullDetails.NextReviewLevel

	documentCRWorkflows, err := drs.repo.GetComments(id, *currentApprovalLevel)
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

// func (drs *DocumentChangeRequestReviewService) AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error {
// 	document, err := drs.documentSvc.GetFullDetailsByID(id, nil)
// 	if err != nil {
// 		return err
// 	}

// 	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel
// 	currentDocumentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
// 	if err != nil {
// 		return err
// 	}

// 	if documentReviewCommentFormDTO.IsApproved {
// 		return drs.handleApprovedCase(id, currentUserID, documentReviewCommentFormDTO, document, currentDocumentWorkflow, *currentApprovalLevel)
// 	} else if documentReviewCommentFormDTO.IsRejected {
// 		return drs.handleRejectedCase(id, currentUserID, documentReviewCommentFormDTO, document, currentDocumentWorkflow)
// 	} else if documentReviewCommentFormDTO.IsReverted {
// 		return drs.handleRevertedCase(id, currentUserID, document, documentReviewCommentFormDTO, currentDocumentWorkflow)
// 	} else {
// 		return drs.handleCommentOnlyCase(id, currentUserID, document, documentReviewCommentFormDTO, currentDocumentWorkflow)
// 	}
// }

// func (drs *DocumentChangeRequestReviewService) handleApprovedCase(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, document *documentDTOs.DocumentDetailsDTO, currentDocumentWorkflow *models.DocumentWorkflow, currentApprovalLevel int) error {
// 	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("approve")
// 	if err != nil {
// 		return err
// 	}
// 	switch currentDocumentWorkflow.WorkflowLevelType.Type {
// 	case "individual":
// 		if currentUserID != currentDocumentWorkflow.DocumentWorkflowParticipants[0].UserID {
// 			return utility.NewForbiddenError("you have no permission to approve")
// 		}

// 		documentWorkflowReview := &models.DocumentWorkflowReview{
// 			DocumentWorkflowID:       currentDocumentWorkflow.ID,
// 			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
// 			Comment:                  documentReviewCommentFormDTO.Comment,
// 			CreatedBYID:              currentUserID,
// 		}

// 		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
// 		if err != nil {
// 			return err
// 		}

// 		documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
// 		if err != nil {
// 			return err
// 		}

// 		documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
// 		err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = drs.taskSvc.UpdateDocumentTaskStatus(id, currentUserID, currentDocumentWorkflow.ID)
// 		if err != nil {
// 			return err
// 		}

// 		drs.sendApprovalNotification(id, currentUserID, document)

// 		nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, currentApprovalLevel)
// 		if err != nil {
// 			return err
// 		}

// 		if nextDocWorkflowLevel == nil {
// 			documentStatus, err := drs.documentStatusSvc.GetByType("approved")
// 			if err != nil {
// 				return err
// 			}
// 			existingDocument, err := drs.documentSvc.GetByID(id)
// 			if err != nil {
// 				return err
// 			}

// 			existingDocument.DocumentStatusID = documentStatus.ID
// 			err = drs.documentSvc.UpdateDocument(existingDocument)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
// 			if err != nil {
// 				return err
// 			}
// 			existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
// 			err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
// 			if err != nil {
// 				return err
// 			}
// 			nextWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, nextDocWorkflowLevel.Level)
// 			if err != nil {
// 				return err
// 			}
// 			err = drs.CreateTasksForNextLevelUsers(id, nextWorkflow, document)
// 			if err != nil {
// 				return err
// 			}

// 			drs.notifyNextLevelUsers(id, nextWorkflow, document)
// 		}
// 	case "multiple-users", "role":
// 		allowed := false
// 		for _, participant := range currentDocumentWorkflow.DocumentWorkflowParticipants {
// 			if participant.UserID == currentUserID {
// 				allowed = true
// 				break
// 			}
// 		}

// 		if !allowed {
// 			return utility.NewForbiddenError("you have no permission to approve")
// 		}

// 		documentWorkflowReview := &models.DocumentWorkflowReview{
// 			DocumentWorkflowID:       currentDocumentWorkflow.ID,
// 			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
// 			Comment:                  documentReviewCommentFormDTO.Comment,
// 			CreatedBYID:              currentUserID,
// 		}

// 		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
// 		if err != nil {
// 			return err
// 		}

// 		drs.sendApprovalNotification(id, currentUserID, document)

// 		if currentDocumentWorkflow.ApprovalTypes == "all" {
// 			_, err := drs.taskSvc.UpdateDocumentTaskStatus(id, currentUserID, currentDocumentWorkflow.ID)
// 			if err != nil {
// 				return err
// 			}
// 			allApproved, err := drs.documentWorkflowSvc.CheckAllParticipantsApproved(currentDocumentWorkflow.ID, currentDocumentWorkflow.DocumentWorkflowParticipants)
// 			if err != nil {
// 				return err
// 			}

// 			if !allApproved {
// 				return nil
// 			} else {
// 				documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
// 				if err != nil {
// 					return err
// 				}

// 				documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
// 				err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
// 				if err != nil {
// 					return err
// 				}

// 				nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, currentApprovalLevel)
// 				if err != nil {
// 					return err
// 				}

// 				if nextDocWorkflowLevel == nil {
// 					documentStatus, err := drs.documentStatusSvc.GetByType("approved")
// 					if err != nil {
// 						return err
// 					}
// 					existingDocument, err := drs.documentSvc.GetByID(id)
// 					if err != nil {
// 						return err
// 					}

// 					existingDocument.DocumentStatusID = documentStatus.ID
// 					err = drs.documentSvc.UpdateDocument(existingDocument)
// 					if err != nil {
// 						return err
// 					}
// 				} else {
// 					existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
// 					if err != nil {
// 						return err
// 					}

// 					existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
// 					err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
// 					if err != nil {
// 						return err
// 					}
// 					nextWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, nextDocWorkflowLevel.Level)
// 					if err != nil {
// 						return err
// 					}
// 					err = drs.CreateTasksForNextLevelUsers(id, nextWorkflow, document)
// 					if err != nil {
// 						return err
// 					}
// 					fmt.Println("124")
// 					drs.notifyNextLevelUsers(id, nextWorkflow, document)
// 				}
// 			}
// 		} else {
// 			documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
// 			if err != nil {
// 				return err
// 			}

// 			documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
// 			err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
// 			if err != nil {
// 				return err
// 			}

// 			task, err := drs.taskSvc.UpdateDocumentTaskStatus(id, currentUserID, currentDocumentWorkflow.ID)
// 			if err != nil {
// 				return err
// 			}

// 			// Parse the related_task_ids from the task
// 			var relatedTaskIDs []uint
// 			if task.RelatedTaskIDs != nil {
// 				if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
// 					return err
// 				}

// 				// Close each related task
// 				for _, relatedTaskID := range relatedTaskIDs {
// 					relatedTask, err := drs.taskSvc.GetByID(relatedTaskID)
// 					if err != nil {
// 						return err
// 					}
// 					_, err = drs.taskSvc.UpdateDocumentTaskStatus(id, relatedTask.ResponsibleUserID, currentDocumentWorkflow.ID)
// 					if err != nil {
// 						return err
// 					}
// 				}
// 			}

// 			nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, currentApprovalLevel)
// 			if err != nil {
// 				return err
// 			}

// 			if nextDocWorkflowLevel == nil {
// 				documentStatus, err := drs.documentStatusSvc.GetByType("approved")
// 				if err != nil {
// 					return err
// 				}
// 				existingDocument, err := drs.documentSvc.GetByID(id)
// 				if err != nil {
// 					return err
// 				}

// 				existingDocument.DocumentStatusID = documentStatus.ID
// 				err = drs.documentSvc.UpdateDocument(existingDocument)
// 				if err != nil {
// 					return err
// 				}
// 			} else {
// 				existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
// 				if err != nil {
// 					return err
// 				}

// 				existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
// 				err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
// 				if err != nil {
// 					return err
// 				}

// 				nextWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, nextDocWorkflowLevel.Level)
// 				if err != nil {
// 					return err
// 				}

// 				err = drs.CreateTasksForNextLevelUsers(id, nextWorkflow, document)
// 				if err != nil {
// 					return err
// 				}

// 				drs.notifyNextLevelUsers(id, nextWorkflow, document)
// 			}
// 		}
// 	}

// 	return nil
// }

// func (drs *DocumentChangeRequestReviewService) handleRejectedCase(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, document *documentDTOs.DocumentDetailsDTO, currentDocumentWorkflow *models.DocumentWorkflow) error {
// 	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("reject")
// 	if err != nil {
// 		return err
// 	}

// 	documentWorkflowReview := &models.DocumentWorkflowReview{
// 		DocumentWorkflowID:       currentDocumentWorkflow.ID,
// 		DocumentWorkflowStatusID: documentWorkflowStatus.ID,
// 		Comment:                  documentReviewCommentFormDTO.Comment,
// 		CreatedBYID:              currentUserID,
// 	}

// 	err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
// 	if err != nil {
// 		return err
// 	}

// 	documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
// 	if err != nil {
// 		return err
// 	}

// 	documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
// 	err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
// 	if err != nil {
// 		return err
// 	}

// 	documentStatus, err := drs.documentStatusSvc.GetByType("rejected")
// 	if err != nil {
// 		return err
// 	}
// 	existingDocument, err := drs.documentSvc.GetByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	existingDocument.DocumentStatusID = documentStatus.ID
// 	err = drs.documentSvc.UpdateDocument(existingDocument)
// 	if err != nil {
// 		return err
// 	}

// 	task, err := drs.taskSvc.UpdateDocumentTaskStatus(id, currentUserID, currentDocumentWorkflow.ID)
// 	if err != nil {
// 		return err
// 	}

// 	// Parse the related_task_ids from the task
// 	var relatedTaskIDs []uint
// 	if task.RelatedTaskIDs != nil {
// 		if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
// 			return err
// 		}

// 		// Close each related task
// 		for _, relatedTaskID := range relatedTaskIDs {
// 			relatedTask, err := drs.taskSvc.GetByID(relatedTaskID)
// 			if err != nil {
// 				return err
// 			}
// 			_, err = drs.taskSvc.UpdateDocumentTaskStatus(id, relatedTask.ResponsibleUserID, currentDocumentWorkflow.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	drs.sendRejectedNotification(id, currentUserID, document)

// 	return nil
// }

// func (drs *DocumentChangeRequestReviewService) handleRevertedCase(id, currentUserID uint, document *documentDTOs.DocumentDetailsDTO, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, currentDocumentWorkflow *models.DocumentWorkflow) error {
// 	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("revert")
// 	if err != nil {
// 		return err
// 	}
// 	documentWorkflowReview := &models.DocumentWorkflowReview{
// 		DocumentWorkflowID:       currentDocumentWorkflow.ID,
// 		DocumentWorkflowStatusID: documentWorkflowStatus.ID,
// 		Comment:                  documentReviewCommentFormDTO.Comment,
// 		CreatedBYID:              currentUserID,
// 	}

// 	err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
// 	if err != nil {
// 		return err
// 	}

// 	documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
// 	if err != nil {
// 		return err
// 	}

// 	documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
// 	err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
// 	if err != nil {
// 		return err
// 	}

// 	documentStatus, err := drs.documentStatusSvc.GetByType("reverted")
// 	if err != nil {
// 		return err
// 	}
// 	existingDocument, err := drs.documentSvc.GetByID(id)
// 	if err != nil {
// 		return err
// 	}

// 	existingDocument.DocumentStatusID = documentStatus.ID
// 	err = drs.documentSvc.UpdateDocument(existingDocument)
// 	if err != nil {
// 		return err
// 	}

// 	task, err := drs.taskSvc.UpdateDocumentTaskStatus(id, currentUserID, currentDocumentWorkflow.ID)
// 	if err != nil {
// 		return err
// 	}

// 	// Parse the related_task_ids from the task
// 	var relatedTaskIDs []uint
// 	if task.RelatedTaskIDs != nil {
// 		if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
// 			return err
// 		}

// 		// Close each related task
// 		for _, relatedTaskID := range relatedTaskIDs {
// 			relatedTask, err := drs.taskSvc.GetByID(relatedTaskID)
// 			if err != nil {
// 				return err
// 			}
// 			_, err = drs.taskSvc.UpdateDocumentTaskStatus(id, relatedTask.ResponsibleUserID, currentDocumentWorkflow.ID)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}

// 	drs.sendRevertedNotification(id, currentUserID, document)

// 	return nil
// }

// func (drs *DocumentChangeRequestReviewService) handleCommentOnlyCase(id, currentUserID uint, document *documentDTOs.DocumentDetailsDTO, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO, currentDocumentWorkflow *models.DocumentWorkflow) error {
// 	documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("comment")
// 	if err != nil {
// 		return err
// 	}
// 	documentWorkflowReview := &models.DocumentWorkflowReview{
// 		DocumentWorkflowID:       currentDocumentWorkflow.ID,
// 		DocumentWorkflowStatusID: documentWorkflowStatus.ID,
// 		Comment:                  documentReviewCommentFormDTO.Comment,
// 		CreatedBYID:              currentUserID,
// 	}

// 	err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
// 	if err != nil {
// 		return err
// 	}

// 	drs.sendCommentNotification(id, currentUserID, documentReviewCommentFormDTO.Comment, document)

// 	return nil
// }

// func (drs *DocumentChangeRequestReviewService) sendApprovalNotification(id uint, currentUserID uint, document *documentDTOs.DocumentDetailsDTO) {
// 	go func() {
// 		// Prepare URL and parameters for the notification
// 		urlParam := fmt.Sprintf(`{"id":%d}`, id)
// 		entity := "document"
// 		action := "approval"

// 		// Create the database notification DTO
// 		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
// 			UserID:               uint(document.ResponsibleUser.ID),
// 			ModuleIdentifier:     "document",
// 			NotificationType:     string(commonDTOs.ApprovalType),
// 			HeadingEn:            "Document Approval",
// 			HeadingAr:            "مراجعة الوثيقة",
// 			ReferenceNumber:      document.ReferenceNumber,
// 			TitleEn:              document.TitleEn,
// 			TitleAr:              document.TitleAr,
// 			Entity:               &entity,
// 			Action:               &action,
// 			URLParam:             urlParam,
// 			IsSystemNotification: false,
// 			CreatedByID:          &currentUserID,
// 		}

// 		// Get the approving user's details
// 		approvedByUser, err := drs.userSvc.GetByID(currentUserID)
// 		if err != nil {
// 			utility.LogError("Failed to fetch approving user", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		// Prepare email details
// 		appURL := utility.GetAppURL()
// 		emailTo := []string{document.ResponsibleUser.Email}
// 		subject := fmt.Sprintf("Document/%s: Approval Notification", document.TitleEn)
// 		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

// 		emailData := map[string]interface{}{
// 			"FirstNameEn":           document.ResponsibleUser.FirstNameEN,
// 			"LastNameEn":            document.ResponsibleUser.LastNameEN,
// 			"DocumentTitle":         document.TitleEn,
// 			"DocumentType":          document.DocumentType.TitleEn,
// 			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
// 			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
// 			"ApprovedAt":            time.Now().Format("2006-01-02"),
// 			"RedirectURL":           urlParam,
// 		}

// 		// Create the email notification DTO
// 		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
// 			EmailTo:              emailTo,
// 			Subject:              subject,
// 			TemplateRelativePath: "document/templates/email/document_workflow_approve_user_mail.html",
// 			EmailData:            emailData,
// 		}

// 		// Send the notification
// 		err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
// 		if err != nil {
// 			utility.LogError("Failed to send approval notification email", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 		}
// 	}()
// }

// func (drs *DocumentChangeRequestReviewService) notifyNextLevelUsers(id uint, nextDocWorkflowLevel *models.DocumentWorkflow, document *documentDTOs.DocumentDetailsDTO) {
// 	go func() {
// 		urlParam := fmt.Sprintf(`{"id":%d}`, id)
// 		entity := "document"
// 		action := "review-user-assign"

// 		for _, participant := range nextDocWorkflowLevel.DocumentWorkflowParticipants {
// 			responsibleUser, err := drs.userSvc.GetByID(participant.UserID)
// 			if err != nil {
// 				utility.LogError("Failed to fetch review user", map[string]interface{}{
// 					"error": err.Error(),
// 				})
// 				continue
// 			}

// 			createdByID := uint(document.ResponsibleUser.ID)
// 			// Prepare database notification DTO
// 			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
// 				UserID:               participant.UserID,
// 				ModuleIdentifier:     "document",
// 				NotificationType:     string(commonDTOs.ReviewType),
// 				HeadingEn:            "Document Review",
// 				HeadingAr:            "مراجعة الوثيقة",
// 				ReferenceNumber:      document.ReferenceNumber,
// 				TitleEn:              document.TitleEn,
// 				TitleAr:              document.TitleAr,
// 				Entity:               &entity,
// 				Action:               &action,
// 				URLParam:             urlParam,
// 				IsSystemNotification: false,
// 				CreatedByID:          &createdByID,
// 			}

// 			// Prepare email notification DTO
// 			emailTo := []string{responsibleUser.Email}
// 			subject := fmt.Sprintf("Document/%s: Review Assigned", document.TitleEn)
// 			appURL := utility.GetAppURL()

// 			emailData := map[string]interface{}{
// 				"FirstNameEn":           responsibleUser.FirstNameEn,
// 				"LastNameEn":            responsibleUser.LastNameEn,
// 				"DocumentTitle":         document.TitleEn,
// 				"DocumentType":          document.DocumentType.TitleEn,
// 				"AssignedByFirstNameEn": document.ResponsibleUser.FirstNameEN,
// 				"AssignedByLastNameEn":  document.ResponsibleUser.LastNameEN,
// 				"AssignedOn":            time.Now().Format("2006-01-02"),
// 				"RedirectURL":           fmt.Sprintf("%s/document/register/%d", appURL, id),
// 			}

// 			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
// 				EmailTo:              emailTo,
// 				Subject:              subject,
// 				TemplateRelativePath: "document/templates/email/document_workflow_review_user_mail.html",
// 				EmailData:            emailData,
// 			}

// 			// Send both database and email notifications
// 			err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
// 			if err != nil {
// 				utility.LogError("Failed to send next level review notification", map[string]interface{}{
// 					"error": err.Error(),
// 				})
// 			}
// 		}
// 	}()
// }

// func (drs *DocumentChangeRequestReviewService) CreateTasksForNextLevelUsers(documentID uint, workflowLevel *models.DocumentWorkflow, document *documentDTOs.DocumentDetailsDTO) error {
// 	urlParam := fmt.Sprintf(`{"id":%d}`, documentID)
// 	var taskIDs []uint

// 	createdByID := uint(document.ResponsibleUser.ID)
// 	// Loop through each participant and create a task
// 	for _, participant := range workflowLevel.DocumentWorkflowParticipants {
// 		taskDTO := commonDTOs.TaskFormDTO{
// 			ModuleIdentifier:   "document",
// 			DocumentID:         &documentID,
// 			DocumentWorkflowID: &workflowLevel.ID,
// 			Title:              document.TitleEn,
// 			Description:        document.DescriptionEn,
// 			URLParam:           &urlParam,
// 			QueryParam:         nil,
// 			ResponsibleUserID:  participant.User.ID,
// 		}

// 		taskID, err := drs.taskSvc.Create(taskDTO, createdByID)
// 		if err != nil {
// 			return err
// 		}
// 		taskIDs = append(taskIDs, taskID)
// 	}

// 	// Set up related task IDs for each task
// 	for _, taskID := range taskIDs {
// 		relatedTaskIDs := []uint{}
// 		for _, otherTaskID := range taskIDs {
// 			if otherTaskID != taskID {
// 				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
// 			}
// 		}

// 		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
// 		if err != nil {
// 			return err
// 		}

// 		err = drs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func (drs *DocumentChangeRequestReviewService) sendRejectedNotification(id uint, currentUserID uint, document *documentDTOs.DocumentDetailsDTO) {
// 	go func() {
// 		// Prepare URL and parameters for the notification
// 		urlParam := fmt.Sprintf(`{"id":%d}`, id)
// 		entity := "document"
// 		action := "rejected"

// 		// Create the database notification DTO
// 		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
// 			UserID:               uint(document.ResponsibleUser.ID),
// 			ModuleIdentifier:     "document",
// 			NotificationType:     string(commonDTOs.RejectType),
// 			HeadingEn:            "Document Rejected",
// 			HeadingAr:            "تم رفض الوثيقة",
// 			ReferenceNumber:      document.ReferenceNumber,
// 			TitleEn:              document.TitleEn,
// 			TitleAr:              document.TitleAr,
// 			Entity:               &entity,
// 			Action:               &action,
// 			URLParam:             urlParam,
// 			IsSystemNotification: false,
// 			CreatedByID:          &currentUserID,
// 		}

// 		// Get the approving user's details
// 		rejectedByUser, err := drs.userSvc.GetByID(currentUserID)
// 		if err != nil {
// 			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		// Prepare email details
// 		appURL := utility.GetAppURL()
// 		emailTo := []string{document.ResponsibleUser.Email}
// 		subject := fmt.Sprintf("Document/%s: Rejected Notification", document.TitleEn)
// 		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

// 		emailData := map[string]interface{}{
// 			"FirstNameEn":           document.ResponsibleUser.FirstNameEN,
// 			"LastNameEn":            document.ResponsibleUser.LastNameEN,
// 			"DocumentTitle":         document.TitleEn,
// 			"DocumentType":          document.DocumentType.TitleEn,
// 			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
// 			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
// 			"ApprovedAt":            time.Now().Format("2006-01-02"),
// 			"RedirectURL":           urlParam,
// 		}

// 		// Create the email notification DTO
// 		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
// 			EmailTo:              emailTo,
// 			Subject:              subject,
// 			TemplateRelativePath: "document/templates/email/document_workflow_reject_user_mail.html",
// 			EmailData:            emailData,
// 		}

// 		// Send the notification
// 		err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
// 		if err != nil {
// 			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 		}
// 	}()
// }

// func (drs *DocumentChangeRequestReviewService) sendRevertedNotification(id uint, currentUserID uint, document *documentDTOs.DocumentDetailsDTO) {
// 	go func() {
// 		// Prepare URL and parameters for the notification
// 		urlParam := fmt.Sprintf(`{"id":%d}`, id)
// 		entity := "document"
// 		action := "reverted"

// 		// Create the database notification DTO
// 		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
// 			UserID:               uint(document.ResponsibleUser.ID),
// 			ModuleIdentifier:     "document",
// 			NotificationType:     string(commonDTOs.RevertType),
// 			HeadingEn:            "Document Reverted",
// 			HeadingAr:            "تم إعادة الوثيقة",
// 			ReferenceNumber:      document.ReferenceNumber,
// 			TitleEn:              document.TitleEn,
// 			TitleAr:              document.TitleAr,
// 			Entity:               &entity,
// 			Action:               &action,
// 			URLParam:             urlParam,
// 			IsSystemNotification: false,
// 			CreatedByID:          &currentUserID,
// 		}

// 		// Get the user who performed the reversion
// 		revertedByUser, err := drs.userSvc.GetByID(currentUserID)
// 		if err != nil {
// 			utility.LogError("Failed to fetch reversion performed user", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		// Prepare email details
// 		appURL := utility.GetAppURL()
// 		emailTo := []string{document.ResponsibleUser.Email}
// 		subject := fmt.Sprintf("Document/%s: Reverted Notification", document.TitleEn)
// 		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

// 		emailData := map[string]interface{}{
// 			"FirstNameEn":           document.ResponsibleUser.FirstNameEN,
// 			"LastNameEn":            document.ResponsibleUser.LastNameEN,
// 			"DocumentTitle":         document.TitleEn,
// 			"DocumentType":          document.DocumentType.TitleEn,
// 			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
// 			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
// 			"RevertedAt":            time.Now().Format("2006-01-02"),
// 			"RedirectURL":           urlParam,
// 		}

// 		// Create the email notification DTO
// 		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
// 			EmailTo:              emailTo,
// 			Subject:              subject,
// 			TemplateRelativePath: "document/templates/email/document_workflow_revert_user_mail.html",
// 			EmailData:            emailData,
// 		}

// 		// Send the notification
// 		err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
// 		if err != nil {
// 			utility.LogError("Failed to send reverted notification email", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 		}
// 	}()
// }

// func (drs *DocumentChangeRequestReviewService) sendCommentNotification(id uint, currentUserID uint, commentText string, document *documentDTOs.DocumentDetailsDTO) {
// 	go func() {
// 		// Prepare URL and parameters for the notification
// 		urlParam := fmt.Sprintf(`{"id":%d}`, id)
// 		entity := "document"
// 		action := "comment"

// 		// Create the database notification DTO
// 		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
// 			UserID:               uint(document.ResponsibleUser.ID),
// 			ModuleIdentifier:     "document",
// 			NotificationType:     string(commonDTOs.CommentType),
// 			HeadingEn:            "New Document Comment",
// 			HeadingAr:            "تعليق جديد على الوثيقة",
// 			ReferenceNumber:      document.ReferenceNumber,
// 			TitleEn:              document.TitleEn,
// 			TitleAr:              document.TitleAr,
// 			Entity:               &entity,
// 			Action:               &action,
// 			URLParam:             urlParam,
// 			IsSystemNotification: false,
// 			CreatedByID:          &currentUserID,
// 		}

// 		// Get the user who made the comment
// 		commentByUser, err := drs.userSvc.GetByID(currentUserID)
// 		if err != nil {
// 			utility.LogError("Failed to fetch user who commented", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 			return
// 		}

// 		// Prepare email details
// 		appURL := utility.GetAppURL()
// 		emailTo := []string{document.ResponsibleUser.Email}
// 		subject := fmt.Sprintf("Document/%s: New Comment Notification", document.TitleEn)
// 		urlParam = fmt.Sprintf("%s/document/register/%d", appURL, id)

// 		emailData := map[string]interface{}{
// 			"FirstNameEn":          document.ResponsibleUser.FirstNameEN,
// 			"LastNameEn":           document.ResponsibleUser.LastNameEN,
// 			"DocumentTitle":        document.TitleEn,
// 			"DocumentType":         document.DocumentType.TitleEn,
// 			"CommentByFirstNameEn": commentByUser.FirstNameEn,
// 			"CommentByLastNameEn":  commentByUser.LastNameEn,
// 			"CommentText":          commentText,
// 			"CommentedAt":          time.Now().Format("2006-01-02"),
// 			"RedirectURL":          urlParam,
// 		}

// 		// Create the email notification DTO
// 		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
// 			EmailTo:              emailTo,
// 			Subject:              subject,
// 			TemplateRelativePath: "document/templates/email/document_workflow_comment_user_mail.html",
// 			EmailData:            emailData,
// 		}

// 		// Send the notification
// 		err = drs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
// 		if err != nil {
// 			utility.LogError("Failed to send comment notification email", map[string]interface{}{
// 				"error": err.Error(),
// 			})
// 		}
// 	}()
// }
