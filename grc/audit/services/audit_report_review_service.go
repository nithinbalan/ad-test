package services

import (
	"encoding/json"
	"errors"
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditReportReviewServiceInterface interface {
	GetComments(auditID uint) (*auditDTOs.AuditFinalReportWorkflowDTO, error)
	Submit(auditID uint, currentUserID uint) error
	AddComment(id uint, currentUserID uint, auditReportReviewCommentFormDTO *auditDTOs.AuditReportReviewCommentFormDTO) error
}

type AuditReportReviewService struct {
	repo                      repositories.AuditReportReviewRepositoryInterface
	logService                logServices.ActivityLogServiceInterface
	auditFinalReportSvc       AuditFinalReportServiceInterface
	auditReportWorkflowSvc    AuditReportWorkflowServiceInterface
	workflowStatusSvc         masterServices.WorkflowStatusServiceInterface
	auditFinalReportStatusSvc AuditFinalReportStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	taskSvc                   taskServices.TaskServiceInterface
}

func NewAuditReportReviewService(
	repo repositories.AuditReportReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	auditFinalReportSvc AuditFinalReportServiceInterface,
	auditReportWorkflowSvc AuditReportWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	auditFinalReportStatusSvc AuditFinalReportStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) AuditReportReviewServiceInterface {
	return &AuditReportReviewService{
		repo:                      repo,
		logService:                logService,
		auditFinalReportSvc:       auditFinalReportSvc,
		auditReportWorkflowSvc:    auditReportWorkflowSvc,
		workflowStatusSvc:         workflowStatusSvc,
		auditFinalReportStatusSvc: auditFinalReportStatusSvc,
		userSvc:                   userSvc,
		notificationSvc:           notificationSvc,
		taskSvc:                   taskSvc,
	}
}

func (arrs *AuditReportReviewService) GetComments(auditID uint) (*auditDTOs.AuditFinalReportWorkflowDTO, error) {
	var auditFinalReportWorkflowDTO auditDTOs.AuditFinalReportWorkflowDTO

	auditFinalReport, err := arrs.auditFinalReportSvc.GetByAuditID(auditID)
	if err != nil {
		return nil, err
	}
	auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("draft")
	if err != nil {
		return nil, err
	}
	if auditFinalReport.AuditFinalReportStatusID == auditFinalReportStatus.ID {
		return nil, utility.NewForbiddenError("Audit final report must be submitted for review to access comments.")
	}

	currentApprovalLevel := auditFinalReport.NextReviewLevel

	auditFinalReportWorkflows, err := arrs.repo.GetComments(auditFinalReport.ID, currentApprovalLevel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditID)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFinalReportWorkflows, &auditFinalReportWorkflowDTO)
	if err != nil {
		return nil, err
	}
	// Explicitly map `id` to `audit_final_report_id`
	auditFinalReportWorkflowDTO.AuditFinalReportID = auditFinalReportWorkflowDTO.ID
	return &auditFinalReportWorkflowDTO, nil
}

func (arrs *AuditReportReviewService) Submit(auditID uint, currentUserID uint) error {
	existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByAuditID(auditID)
	if err != nil {
		return err
	}
	auditFinalReport, err := arrs.auditFinalReportSvc.GetByID(existingAuditFinalReport.ID)
	if err != nil {
		return err
	}
	tx := arrs.repo.BeginTransaction()
	currentDateStr := time.Now()

	// currentDateStr := time.Now().Format("2006-01-02")
	// submittedAt, err := utility.ParseDate(currentDateStr)
	// if err != nil {
	// 	return fmt.Errorf("error parsing current date: %w", err)
	// }
	//begin tx
	
	auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("in-review")
	if err != nil {
		return err
	}
	var level int
	if existingAuditFinalReport.AuditFinalReportStatus.Type == "revert" {
		level = existingAuditFinalReport.NextReviewLevel
	} else {
		// For non-reverted report, start from level 1
		level = 1
	}

	auditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID
	auditFinalReport.SubmittedBYID = &currentUserID
	auditFinalReport.SubmittedAt = &currentDateStr
	auditFinalReport.NextReviewLevel = level

	err = arrs.auditFinalReportSvc.Update(auditFinalReport,tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Fetch the workflow level details, including participants (review users)
	auditReportWorkflowLevel, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByLevel(auditFinalReport.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	workflowStatus, err := arrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}
	auditReportReview := &models.AuditReportReview{
		AuditReportWorkflowID: auditReportWorkflowLevel.ID,
		WorkflowStatusID:      workflowStatus.ID,
		CreatedBYID:           currentUserID,
	}

	err = arrs.repo.CreateAuditReportReview(auditReportReview, tx)
	if err != nil {
		return err
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, existingAuditFinalReport.ID)
	var taskIDs []uint
	auditFinalReportID := &existingAuditFinalReport.ID
	auditReportWorkflowID := &auditReportWorkflowLevel.ID
	for _, participant := range auditReportWorkflowLevel.AuditReportWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:      "audit",
			AuditFinalReportID:    auditFinalReportID,
			AuditReportWorkflowID: auditReportWorkflowID,
			Title:                 existingAuditFinalReport.Audit.AuditPlan.Title,
			Description:           existingAuditFinalReport.Audit.AuditPlan.Description,
			URLParam:              &urlParam,
			QueryParam:            nil,
			ResponsibleUserID:     participant.User.ID,
		}

		taskID, err := arrs.taskSvc.Create(taskDTO, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		taskIDs = append(taskIDs, taskID)

	}

	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			tx.Rollback()
			return err
		}

		err = arrs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON), tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	//add audit report edit history
	auditReportEditHistoryDTO := &auditDTOs.AuditReportEditHistoryDTO{
		AuditID:            auditID,
		AuditFinalReportID: existingAuditFinalReport.ID,
		ActionType:         "submit",
	}

	err = arrs.auditFinalReportSvc.CreateAuditReportEditHistory(auditReportEditHistoryDTO, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	//to do
	//activity log

	// Prepare notification
	// todo
	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, auditID)
		entity := "audit-final-report"
		action := "review-user-assign"

		for _, participant := range auditReportWorkflowLevel.AuditReportWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               participant.User.ID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      existingAuditFinalReport.ReferenceNumber,
				HeadingEn:            "Audit Final Report Review",
				HeadingAr:            "مراجعة التقرير النهائي للتدقيق",
				TitleEn:              existingAuditFinalReport.Audit.AuditPlan.Title,
				TitleAr:              existingAuditFinalReport.Audit.AuditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := arrs.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := arrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Audit Final Report /%s: Review Assigned", existingAuditFinalReport.Audit.AuditPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"AuditTitle":            existingAuditFinalReport.Audit.AuditPlan.Title,
				"ReferenceNumber":       existingAuditFinalReport.Audit.AuditPlan.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/audits/%d", appURL, auditID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_report_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (arrs *AuditReportReviewService) AddComment(id uint, currentUserID uint, auditReportReviewCommentFormDTO *auditDTOs.AuditReportReviewCommentFormDTO) error {
	auditFinalReport, err := arrs.auditFinalReportSvc.GetFullDetailsByID(id, currentUserID)
	if err != nil {
		return err
	}
	auditFinalReportID := auditFinalReport.ID
	currentApprovalLevel := auditFinalReport.NextReviewLevel
	currentAuditReportWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByLevel(auditFinalReport.ID, currentApprovalLevel)
	if err != nil {
		return err
	}

	if auditReportReviewCommentFormDTO.IsApproved {
		return arrs.handleApprovedCase(auditFinalReportID, currentUserID, auditReportReviewCommentFormDTO, auditFinalReport, currentAuditReportWorkflow, currentApprovalLevel)
	} else if auditReportReviewCommentFormDTO.IsRejected {
		return arrs.handleRejectedCase(auditFinalReportID, currentUserID, auditReportReviewCommentFormDTO, auditFinalReport, currentAuditReportWorkflow)
	} else if auditReportReviewCommentFormDTO.IsReverted {
		return arrs.handleRevertedCase(auditFinalReportID, currentUserID, auditFinalReport, auditReportReviewCommentFormDTO, currentAuditReportWorkflow)
	} else {
		return arrs.handleCommentOnlyCase(auditFinalReportID, currentUserID, auditFinalReport, auditReportReviewCommentFormDTO, currentAuditReportWorkflow)
	}
}

func (arrs *AuditReportReviewService) handleApprovedCase(id uint, currentUserID uint, auditReportReviewCommentFormDTO *auditDTOs.AuditReportReviewCommentFormDTO, auditFinalReport *auditDTOs.AuditFinalReportDetailDTO, currentAuditReportWorkflow *models.AuditReportWorkflow, currentApprovalLevel int) error {
	workflowStatus, err := arrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}
	switch currentAuditReportWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentAuditReportWorkflow.AuditReportWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		auditReportReview := &models.AuditReportReview{
			AuditReportWorkflowID: currentAuditReportWorkflow.ID,
			WorkflowStatusID:      workflowStatus.ID,
			Comment:               &auditReportReviewCommentFormDTO.Comment,
			CreatedBYID:           currentUserID,
		}

		err = arrs.repo.CreateAuditReportReview(auditReportReview)
		if err != nil {
			return err
		}

		auditReportWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByID(currentAuditReportWorkflow.ID)
		if err != nil {
			return err
		}

		auditReportWorkflow.WorkflowStatusID = workflowStatus.ID
		err = arrs.auditReportWorkflowSvc.UpdateAuditReportWorkflow(auditReportWorkflow)
		if err != nil {
			return err
		}
		_, err = arrs.taskSvc.UpdateAuditReportTaskStatus(id, currentUserID, currentAuditReportWorkflow.ID)
		if err != nil {
			return err
		}

		arrs.sendApprovalNotification(id, currentUserID, auditFinalReport)

		nextAuditReportWorkflowLevel, err := arrs.auditReportWorkflowSvc.GetNextLevel(auditFinalReport.ID, currentApprovalLevel)
		if err != nil {
			return err
		}

		if nextAuditReportWorkflowLevel == nil {
			auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("approve")
			if err != nil {
				return err
			}
			existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
			if err != nil {
				return err
			}

			existingAuditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID
			err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
			if err != nil {
				return err
			}
		} else {
			existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
			if err != nil {
				return err
			}
			existingAuditFinalReport.NextReviewLevel = nextAuditReportWorkflowLevel.Level
			err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
			if err != nil {
				return err
			}

			nextWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByLevel(auditFinalReport.ID, nextAuditReportWorkflowLevel.Level)
			if err != nil {
				return err
			}
			err = arrs.CreateTasksForNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport, currentUserID)
			if err != nil {
				return err
			}

			arrs.notifyNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport)
		}
	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentAuditReportWorkflow.AuditReportWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}

		if !allowed {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		auditReportReview := &models.AuditReportReview{
			AuditReportWorkflowID: currentAuditReportWorkflow.ID,
			WorkflowStatusID:      workflowStatus.ID,
			Comment:               &auditReportReviewCommentFormDTO.Comment,
			CreatedBYID:           currentUserID,
		}

		err = arrs.repo.CreateAuditReportReview(auditReportReview)
		if err != nil {
			return err
		}
		arrs.sendApprovalNotification(id, currentUserID, auditFinalReport)

		if currentAuditReportWorkflow.ApprovalType == "all" {
			_, err := arrs.taskSvc.UpdateAuditReportTaskStatus(id, currentUserID, currentAuditReportWorkflow.ID)
			if err != nil {
				return err
			}
			allApproved, err := arrs.auditReportWorkflowSvc.CheckAllParticipantsApproved(currentAuditReportWorkflow.ID, currentAuditReportWorkflow.AuditReportWorkflowParticipants)
			if err != nil {
				return err
			}

			if !allApproved {
				return nil
			} else {
				documentWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByID(currentAuditReportWorkflow.ID)
				if err != nil {
					return err
				}

				documentWorkflow.WorkflowStatusID = workflowStatus.ID
				err = arrs.auditReportWorkflowSvc.UpdateAuditReportWorkflow(documentWorkflow)
				if err != nil {
					return err
				}

				nextAuditReportWorkflowLevel, err := arrs.auditReportWorkflowSvc.GetNextLevel(auditFinalReport.ID, currentApprovalLevel)
				if err != nil {
					return err
				}

				if nextAuditReportWorkflowLevel == nil {
					auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("approve")
					if err != nil {
						return err
					}
					existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
					if err != nil {
						return err
					}

					existingAuditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID
					err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
					if err != nil {
						return err
					}
				} else {
					existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
					if err != nil {
						return err
					}

					existingAuditFinalReport.NextReviewLevel = nextAuditReportWorkflowLevel.Level
					err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
					if err != nil {
						return err
					}
					nextWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByLevel(auditFinalReport.ID, nextAuditReportWorkflowLevel.Level)
					if err != nil {
						return err
					}
					err = arrs.CreateTasksForNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport, currentUserID)
					if err != nil {
						return err
					}

					// fmt.Println("124")
					arrs.notifyNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport)
				}
			}
		} else {
			documentWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByID(currentAuditReportWorkflow.ID)
			if err != nil {
				return err
			}

			documentWorkflow.WorkflowStatusID = workflowStatus.ID
			err = arrs.auditReportWorkflowSvc.UpdateAuditReportWorkflow(documentWorkflow)
			if err != nil {
				return err
			}

			task, err := arrs.taskSvc.UpdateAuditReportTaskStatus(id, currentUserID, currentAuditReportWorkflow.ID)
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
					relatedTask, err := arrs.taskSvc.GetByID(relatedTaskID)
					if err != nil {
						return err
					}
					_, err = arrs.taskSvc.UpdateAuditReportTaskStatus(id, relatedTask.ResponsibleUserID, currentAuditReportWorkflow.ID)
					if err != nil {
						return err
					}
				}
			}

			nextAuditReportWorkflowLevel, err := arrs.auditReportWorkflowSvc.GetNextLevel(auditFinalReport.ID, currentApprovalLevel)
			if err != nil {
				return err
			}

			if nextAuditReportWorkflowLevel == nil {
				auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("approve")
				if err != nil {
					return err
				}
				existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
				if err != nil {
					return err
				}

				existingAuditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID
				err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
				if err != nil {
					return err
				}
			} else {
				existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
				if err != nil {
					return err
				}

				existingAuditFinalReport.NextReviewLevel = nextAuditReportWorkflowLevel.Level
				err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
				if err != nil {
					return err
				}
				nextWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByLevel(auditFinalReport.ID, nextAuditReportWorkflowLevel.Level)
				if err != nil {
					return err
				}

				err = arrs.CreateTasksForNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport, currentUserID)
				if err != nil {
					return err
				}
				arrs.notifyNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport)
			}
		}
	}
	//add audit report edit history
	auditReportEditHistoryDTO := &auditDTOs.AuditReportEditHistoryDTO{
		AuditID:            auditFinalReport.AuditID,
		AuditFinalReportID: auditFinalReport.ID,
		ActionType:         "approve",
		Comment: &auditReportReviewCommentFormDTO.Comment,
	}

	err = arrs.auditFinalReportSvc.CreateAuditReportEditHistory(auditReportEditHistoryDTO, currentUserID)
	if err != nil {
		return err
	}

	return nil
}

// not using reject api
func (arrs *AuditReportReviewService) handleRejectedCase(id uint, currentUserID uint, auditReportReviewCommentFormDTO *auditDTOs.AuditReportReviewCommentFormDTO, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO, currentAuditReportWorkflow *models.AuditReportWorkflow) error {
	workflowStatus, err := arrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}
	auditReportReview := &models.AuditReportReview{
		AuditReportWorkflowID: currentAuditReportWorkflow.ID,
		WorkflowStatusID:      workflowStatus.ID,
		Comment:               &auditReportReviewCommentFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}

	err = arrs.repo.CreateAuditReportReview(auditReportReview)
	if err != nil {
		return err
	}

	auditReportWorkflow, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByID(currentAuditReportWorkflow.ID)
	if err != nil {
		return err
	}

	auditReportWorkflow.WorkflowStatusID = workflowStatus.ID
	err = arrs.auditReportWorkflowSvc.UpdateAuditReportWorkflow(auditReportWorkflow)
	if err != nil {
		return err
	}

	auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}
	existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByID(id)
	if err != nil {
		return err
	}

	existingAuditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID
	err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
	if err != nil {
		return err
	}
	task, err := arrs.taskSvc.UpdateAuditReportTaskStatus(id, currentUserID, currentAuditReportWorkflow.ID)
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
			relatedTask, err := arrs.taskSvc.GetByID(relatedTaskID)
			if err != nil {
				return err
			}
			_, err = arrs.taskSvc.UpdateAuditReportTaskStatus(id, relatedTask.ResponsibleUserID, currentAuditReportWorkflow.ID)
			if err != nil {
				return err
			}
		}
	}
	//todo

	// arrs.sendRejectedNotification(id, currentUserID, auditFinalReportDetailDTO)

	return nil
}

func (arrs *AuditReportReviewService) handleRevertedCase(id, currentUserID uint, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO, auditReportReviewCommentFormDTO *auditDTOs.AuditReportReviewCommentFormDTO, currentAuditReportWorkflow *models.AuditReportWorkflow) error {
	workflowStatus, err := arrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}
	auditReportReview := &models.AuditReportReview{
		AuditReportWorkflowID: currentAuditReportWorkflow.ID,
		WorkflowStatusID:      workflowStatus.ID,
		Comment:               &auditReportReviewCommentFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}

	err = arrs.repo.CreateAuditReportReview(auditReportReview)
	if err != nil {
		return err
	}
	task, err := arrs.taskSvc.UpdateAuditReportTaskStatus(id, currentUserID, currentAuditReportWorkflow.ID)
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
			relatedTask, err := arrs.taskSvc.GetByID(relatedTaskID)
			if err != nil {
				return err
			}
			_, err = arrs.taskSvc.UpdateAuditReportTaskStatus(id, relatedTask.ResponsibleUserID, currentAuditReportWorkflow.ID)
			if err != nil {
				return err
			}
		}
	}
	//add audit report edit history
	auditReportEditHistoryDTO := &auditDTOs.AuditReportEditHistoryDTO{
		AuditID:            auditFinalReportDetailDTO.AuditID,
		AuditFinalReportID: auditFinalReportDetailDTO.ID,
		ActionType:         "revert",
		Comment:         &auditReportReviewCommentFormDTO.Comment,
	}

	err = arrs.auditFinalReportSvc.CreateAuditReportEditHistory(auditReportEditHistoryDTO, currentUserID)
	if err != nil {
		return err
	}
	//todo
	auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}
	auditFinalReport, err := arrs.auditFinalReportSvc.GetByID(auditFinalReportDetailDTO.ID)
	if err != nil {
		return err
	}
	auditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID
	auditFinalReport.UpdatedBYID = &currentUserID

	err = arrs.auditFinalReportSvc.Update(auditFinalReport)
	if err != nil {
		return err
	}

	arrs.sendRevertedNotification(id, currentUserID, auditFinalReportDetailDTO)

	return nil
}

func (arrs *AuditReportReviewService) handleCommentOnlyCase(id, currentUserID uint, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO, auditReportReviewCommentFormDTO *auditDTOs.AuditReportReviewCommentFormDTO, currentAuditReportWorkflow *models.AuditReportWorkflow) error {
	workflowStatus, err := arrs.workflowStatusSvc.GetByType("comment")
	if err != nil {
		return err
	}

	auditReportReview := &models.AuditReportReview{
		AuditReportWorkflowID: currentAuditReportWorkflow.ID,
		WorkflowStatusID:      workflowStatus.ID,
		Comment:               &auditReportReviewCommentFormDTO.Comment,
		CreatedBYID:           currentUserID,
	}

	err = arrs.repo.CreateAuditReportReview(auditReportReview)
	if err != nil {
		return err
	}
	//todo
	// arrs.sendCommentNotification(id, currentUserID, auditReportReviewCommentFormDTO.Comment, auditFinalReportDetailDTO)

	return nil
}

func (arrs *AuditReportReviewService) CreateTasksForNextLevelUsers(auditFinalReportID uint, workflowLevel *models.AuditReportWorkflow, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO, currentUserID uint) error {
	urlParam := fmt.Sprintf(`{"id":%d}`, auditFinalReportID)
	var taskIDs []uint
	createdByID := currentUserID
	// Loop through each participant and create a task
	for _, participant := range workflowLevel.AuditReportWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:      "audit",
			AuditFinalReportID:    &auditFinalReportID,
			AuditReportWorkflowID: &workflowLevel.ID,
			Title:                 auditFinalReportDetailDTO.Title,
			Description:           auditFinalReportDetailDTO.Description,
			URLParam:              &urlParam,
			QueryParam:            nil,
			ResponsibleUserID:     participant.User.ID,
		}

		taskID, err := arrs.taskSvc.Create(taskDTO, createdByID)
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

		err = arrs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	return nil
}

func (arrs *AuditReportReviewService) sendApprovalNotification(id uint, currentUserID uint, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, auditFinalReportDetailDTO.AuditID)
		entity := "audit-final-report"
		action := "approval"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(auditFinalReportDetailDTO.PreparedBy.ID),
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Audit Final Report Approval",
			HeadingAr:            "اعتماد التقرير النهائي للتدقيق",
			ReferenceNumber:      auditFinalReportDetailDTO.ReferenceNumber,
			TitleEn:              auditFinalReportDetailDTO.Title,
			TitleAr:              auditFinalReportDetailDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		approvedByUser, err := arrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{auditFinalReportDetailDTO.PreparedBy.Email}
		subject := fmt.Sprintf("Audit Final Report/%s: Approval Notification", auditFinalReportDetailDTO.Title)
		urlParam = fmt.Sprintf("%s/audits/%d/final-reports/%d", appURL, auditFinalReportDetailDTO.AuditID, id)

		emailData := map[string]interface{}{
			"AuditTitle":            auditFinalReportDetailDTO.Title,
			"Description":           auditFinalReportDetailDTO.Description,
			"FirstNameEn":           auditFinalReportDetailDTO.PreparedBy.FirstNameEN,
			"LastNameEn":            auditFinalReportDetailDTO.PreparedBy.LastNameEN,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "audit/templates/email/audit_report_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (arrs *AuditReportReviewService) notifyNextLevelUsers(id uint, nextDocWorkflowLevel *models.AuditReportWorkflow, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO) {
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, auditFinalReportDetailDTO.AuditID)
		entity := "audit-final-report"
		action := "review-user-assign"

		for _, participant := range nextDocWorkflowLevel.AuditReportWorkflowParticipants {
			responsibleUser, err := arrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(auditFinalReportDetailDTO.PreparedBy.ID)
			// Prepare database notification DTO
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               participant.UserID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Audit Report Review",
				HeadingAr:            "مراجعة الوثيقة",
				ReferenceNumber:      auditFinalReportDetailDTO.ReferenceNumber,
				TitleEn:              auditFinalReportDetailDTO.Title,
				TitleAr:              auditFinalReportDetailDTO.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			// Prepare email notification DTO
			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Document/%s: Review Assigned", auditFinalReportDetailDTO.Title)
			appURL := utility.GetAppURL()
			urlParam = fmt.Sprintf("%s/audits/%d/final-reports/%d", appURL, auditFinalReportDetailDTO.AuditID, id)

			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 auditFinalReportDetailDTO.Title,
				"Description":           auditFinalReportDetailDTO.Description,
				"AssignedByFirstNameEn": auditFinalReportDetailDTO.PreparedBy.FirstNameEN,
				"AssignedByLastNameEn":  auditFinalReportDetailDTO.PreparedBy.LastNameEN,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           urlParam,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_report_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send both database and email notifications
			err = arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (arrs *AuditReportReviewService) sendRevertedNotification(id uint, currentUserID uint, auditFinalReportDetailDTO *auditDTOs.AuditFinalReportDetailDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, auditFinalReportDetailDTO.AuditID)
		entity := "audit-final-report"
		action := "reverted"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               uint(auditFinalReportDetailDTO.PreparedBy.ID),
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Audit Final Report Reverted",
			HeadingAr:            "تم إرجاع التقرير النهائي للتدقيق",
			ReferenceNumber:      auditFinalReportDetailDTO.ReferenceNumber,
			TitleEn:              auditFinalReportDetailDTO.Title,
			TitleAr:              auditFinalReportDetailDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the user who performed the reversion
		approvedByUser, err := arrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch reversion performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{auditFinalReportDetailDTO.PreparedBy.Email}
		subject := fmt.Sprintf("Audit Final Report/%s: Approval Notification", auditFinalReportDetailDTO.Title)
		urlParam = fmt.Sprintf("%s/audits/%d/final-reports/%d", appURL, auditFinalReportDetailDTO.AuditID, id)

		emailData := map[string]interface{}{
			"AuditTitle":            auditFinalReportDetailDTO.Title,
			"Description":           auditFinalReportDetailDTO.Description,
			"FirstNameEn":           auditFinalReportDetailDTO.PreparedBy.FirstNameEN,
			"LastNameEn":            auditFinalReportDetailDTO.PreparedBy.LastNameEN,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "audit/templates/email/audit_report_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
		if err != nil {
			utility.LogError("Failed to send reverted notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
