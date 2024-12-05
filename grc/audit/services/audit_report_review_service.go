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

	return &auditFinalReportWorkflowDTO, nil
}

func (arrs *AuditReportReviewService) Submit(auditID uint, currentUserID uint) error {
	tx := arrs.repo.BeginTransaction()
	existingAuditFinalReport, err := arrs.auditFinalReportSvc.GetByAuditID(auditID)
	if err != nil {
		tx.Rollback()
		return err
	}

	auditFinalReportStatus, err := arrs.auditFinalReportStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAuditFinalReport.AuditFinalReportStatusID = auditFinalReportStatus.ID

	currentDateStr := time.Now().Format("2006-01-02")
	submittedAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	var level int = 1

	existingAuditFinalReport.SubmittedBYID = &currentUserID
	existingAuditFinalReport.SubmittedAt = submittedAt
	existingAuditFinalReport.NextReviewLevel = level

	err = arrs.auditFinalReportSvc.Update(existingAuditFinalReport)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Fetch the workflow level details, including participants (review users)
	auditReportWorkflowLevel, err := arrs.auditReportWorkflowSvc.GetAuditReportWorkflowByLevel(existingAuditFinalReport.ID, level, tx)
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

	err = arrs.repo.CreateAuditReportReview(auditReportReview)
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
	tx.Commit()

	//to do
	//activity log

	// Prepare notification
	// todo

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
		//todo

		// arrs.sendApprovalNotification(id, currentUserID, document)

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
			//todo

			// arrs.notifyNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport)
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
		//todo
		// arrs.sendApprovalNotification(id, currentUserID, document)

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
					//todo

					// fmt.Println("124")
					// arrs.notifyNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport)
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
				//todo
				// arrs.notifyNextLevelUsers(auditFinalReport.ID, nextWorkflow, auditFinalReport)
			}
		}
	}

	return nil
}

//not using reject api
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
	//todo

	// arrs.sendRevertedNotification(id, currentUserID, auditFinalReportDetailDTO)

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
