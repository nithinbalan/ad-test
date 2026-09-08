package services

import (
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"
)

type ProjectDecisionLogReviewServiceInterface interface {
	Submit(projectID uint, decisionLogID uint, currentUserID uint) error
	Resubmit(projectID uint, decisionLogID uint, currentUserID uint) error
	Approve(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) error
	Revert(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) error
	Reject(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) error
}

type ProjectDecisionLogReviewService struct {
	repo              repositories.ProjectDecisionLogReviewRepositoryInterface
	logService        logServices.ActivityLogServiceInterface
	projectDLSvc      ProjectDecisionLogServiceInterface
	workflowSvc       ProjectDecisionLogWorkflowServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewProjectDecisionLogReviewService(
	repo repositories.ProjectDecisionLogReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	projectDLSvc ProjectDecisionLogServiceInterface,
	workflowSvc ProjectDecisionLogWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProjectDecisionLogReviewServiceInterface {
	return &ProjectDecisionLogReviewService{
		repo:              repo,
		logService:        logService,
		projectDLSvc:      projectDLSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
	}
}

func (pdlrs *ProjectDecisionLogReviewService) createReviewWithHistory(workflowID uint, decisionLogID uint, level int, workflowStatusID uint, comment *string, currentUserID uint, tx utility.DatabaseTransaction) error {
	review := &models.ProjectDecisionLogWorkflowReview{
		ProjectDecisionLogWorkflowID: workflowID,
		WorkflowStatusID:             workflowStatusID,
		Comment:                      comment,
		CreatedBYID:                  currentUserID,
	}

	if err := pdlrs.repo.CreateReview(review, tx); err != nil {
		return err
	}

	history := &models.ProjectDecisionLogWorkflowHistory{
		ProjectDecisionLogID:         decisionLogID,
		ProjectDecisionLogWorkflowID: workflowID,
		WorkflowStatusID:             workflowStatusID,
		Level:                        level,
		Comment:                      comment,
		CreatedBYID:                  currentUserID,
	}

	return pdlrs.repo.CreateHistory(history, tx)
}

func (pdlrs *ProjectDecisionLogReviewService) Submit(projectID uint, decisionLogID uint, currentUserID uint) (err error) {
	tx := pdlrs.repo.BeginTransaction()

	var notify func()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			return
		}

		if notify != nil {
			go notify()
		}
	}()

	if err := pdlrs.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	approvedStatus, err := pdlrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	rejectedStatus, err := pdlrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel != nil {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already submitted"),
			"workflow_already_submitted",
			"project decision log",
		)
	}

	if pdl.WorkflowStatusID != nil && (*pdl.WorkflowStatusID == approvedStatus.ID || *pdl.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already completed"),
			"workflow_already_completed",
			"project decision log",
		)
	}

	if err := pdlrs.workflowSvc.DeleteAllWorkflowsForProjectDecisionLog(decisionLogID, tx); err != nil {
		return err
	}

	if err := pdlrs.workflowSvc.CreateProjectDecisionLogWorkflow(decisionLogID, currentUserID, tx); err != nil {
		return err
	}

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	currentDate := time.Now()
	level := 1
	pdl.WorkflowStatusID = &workflowStatus.ID
	pdl.SubmittedBYID = &currentUserID
	pdl.SubmittedAt = &currentDate
	pdl.NextReviewLevel = &level
	pdl.ApprovedBYID = nil
	pdl.ApprovedAt = nil
	pdl.RejectedBYID = nil
	pdl.RejectedAt = nil
	pdl.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
		return err
	}

	workflowLevel, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, level, tx)
	if err != nil {
		return err
	}

	if err := pdlrs.createReviewWithHistory(workflowLevel.ID, decisionLogID, level, workflowStatus.ID, nil, currentUserID, tx); err != nil {
		return err
	}

	notify = func() {
		_ = pdlrs.logService.LogActivity(
			models.CreateAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)

		pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
		if err == nil {
			pdlrs.notifyNextLevelUsers(projectID, decisionLogID, workflowLevel, pdlDetail, currentUserID)
		}
	}

	return err
}

func (pdlrs *ProjectDecisionLogReviewService) Resubmit(projectID uint, decisionLogID uint, currentUserID uint) (err error) {
	tx := pdlrs.repo.BeginTransaction()

	var notify func()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			return
		}

		if notify != nil {
			go notify()
		}
	}()

	if err := pdlrs.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	rejectedStatus, err := pdlrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	revertedStatus, err := pdlrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	if pdl.WorkflowStatusID == nil || (*pdl.WorkflowStatusID != rejectedStatus.ID && *pdl.WorkflowStatusID != revertedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("project decision log has not been rejected or reverted"),
			"workflow_not_rejected",
			"project decision log",
		)
	}

	if err := pdlrs.workflowSvc.DeleteAllWorkflowsForProjectDecisionLog(decisionLogID, tx); err != nil {
		return err
	}

	if err := pdlrs.workflowSvc.CreateProjectDecisionLogWorkflow(decisionLogID, currentUserID, tx); err != nil {
		return err
	}

	submitStatus, err := pdlrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	level := 1
	now := time.Now()

	pdl.WorkflowStatusID = &submitStatus.ID
	pdl.SubmittedBYID = &currentUserID
	pdl.SubmittedAt = &now
	pdl.NextReviewLevel = &level
	pdl.ApprovedBYID = nil
	pdl.ApprovedAt = nil
	pdl.RejectedBYID = nil
	pdl.RejectedAt = nil
	pdl.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
		return err
	}

	workflowLevel, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, level, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByType("resubmit")
	if err != nil {
		return err
	}

	if err := pdlrs.createReviewWithHistory(
		workflowLevel.ID,
		decisionLogID,
		level,
		workflowStatus.ID,
		nil,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	notify = func() {
		_ = pdlrs.logService.LogActivity(
			models.UpdateAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)

		pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
		if err == nil {
			pdlrs.sendResubmitNotification(projectID, decisionLogID, workflowLevel, pdlDetail, currentUserID)
		}
	}

	return err
}

func (pdlrs *ProjectDecisionLogReviewService) Approve(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) (err error) {
	tx := pdlrs.repo.BeginTransaction()

	var notify func()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			return
		}

		if notify != nil {
			go notify()
		}
	}()

	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	rejectedStatus, err := pdlrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel == nil || (pdl.WorkflowStatusID != nil && *pdl.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"project decision log",
		)
	}

	currentApprovalLevel := *pdl.NextReviewLevel
	currentWorkflow, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	pendingStatus, err := pdlrs.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.ProjectDecisionLogWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectDecisionLogWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pdlrs.createReviewWithHistory(currentWorkflow.ID, decisionLogID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			return err
		}

		projectDecisionLogWorkflow, err := pdlrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			return err
		}

		projectDecisionLogWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pdlrs.workflowSvc.UpdateProjectDecisionLogWorkflow(projectDecisionLogWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := pdlrs.workflowSvc.GetNextLevel(decisionLogID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			pdl.WorkflowStatusID = &workflowStatus.ID
			pdl.ApprovedBYID = &currentUserID
			pdl.ApprovedAt = &approvedAt
			pdl.RejectedBYID = nil
			pdl.RejectedAt = nil
			pdl.NextReviewLevel = nil
			pdl.UpdatedByID = &currentUserID

			if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
				return err
			}

			notify = func() {
				_ = pdlrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectDecisionLog{}.TableName(),
					decisionLogID,
					pdl,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
				if err == nil {
					pdlrs.sendApprovalNotification(projectID, decisionLogID, currentUserID, pdlDetail)
				}
			}
		} else {
			pdl.WorkflowStatusID = &pendingStatus.ID
			pdl.NextReviewLevel = &nextWorkflowLevel.Level
			pdl.UpdatedByID = &currentUserID

			if err := pdlrs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
				return err
			}

			if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
				return err
			}

			nextWf, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}

			notify = func() {
				_ = pdlrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectDecisionLog{}.TableName(),
					decisionLogID,
					pdl,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
				if err == nil {
					pdlrs.notifyNextLevelUsers(projectID, decisionLogID, nextWf, pdlDetail, currentUserID)
				}
			}
		}

	case "multiple-users", "role":
		var approvingParticipant *models.ProjectDecisionLogWorkflowParticipant
		for _, participant := range currentWorkflow.ProjectDecisionLogWorkflowParticipants {
			if participant.UserID == currentUserID {
				approvingParticipant = participant
				break
			}
		}
		if approvingParticipant == nil {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pdlrs.createReviewWithHistory(currentWorkflow.ID, decisionLogID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			var allApproved bool
			allApproved, err = pdlrs.workflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ProjectDecisionLogWorkflowParticipants, tx)
			if err != nil || !allApproved {
				if err == nil {
					notify = func() {
						_ = pdlrs.logService.LogActivity(
							models.ApproveAction,
							models.ProjectDecisionLog{}.TableName(),
							decisionLogID,
							pdl,
							currentUserID,
							"project", "project_register", "", nil,
						)
					}
				}
				return err
			}
		}

		projectDecisionLogWorkflow, err := pdlrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			return err
		}

		projectDecisionLogWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pdlrs.workflowSvc.UpdateProjectDecisionLogWorkflow(projectDecisionLogWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := pdlrs.workflowSvc.GetNextLevel(decisionLogID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			pdl.WorkflowStatusID = &workflowStatus.ID
			pdl.ApprovedBYID = &currentUserID
			pdl.ApprovedAt = &approvedAt
			pdl.RejectedBYID = nil
			pdl.RejectedAt = nil
			pdl.NextReviewLevel = nil
			pdl.UpdatedByID = &currentUserID

			if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
				return err
			}

			notify = func() {
				_ = pdlrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectDecisionLog{}.TableName(),
					decisionLogID,
					pdl,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
				if err == nil {
					pdlrs.sendApprovalNotification(projectID, decisionLogID, currentUserID, pdlDetail)
				}
			}
		} else {
			pdl.WorkflowStatusID = &pendingStatus.ID
			pdl.NextReviewLevel = &nextWorkflowLevel.Level
			pdl.UpdatedByID = &currentUserID

			if err := pdlrs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
				return err
			}

			if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
				return err
			}

			nextWf, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}

			notify = func() {
				_ = pdlrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectDecisionLog{}.TableName(),
					decisionLogID,
					pdl,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
				if err == nil {
					pdlrs.notifyNextLevelUsers(projectID, decisionLogID, nextWf, pdlDetail, currentUserID)
				}
			}
		}
	}

	return err
}

func (pdlrs *ProjectDecisionLogReviewService) Revert(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) (err error) {
	tx := pdlrs.repo.BeginTransaction()

	var notify func()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			return
		}

		if notify != nil {
			go notify()
		}
	}()

	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	rejectedStatus, err := pdlrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel == nil || (pdl.WorkflowStatusID != nil && *pdl.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "project decision log")
	}

	currentLevel := *pdl.NextReviewLevel

	currentWorkflow, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, currentLevel, tx)
	if err != nil {
		return err
	}

	allowed := false
	for _, participant := range currentWorkflow.ProjectDecisionLogWorkflowParticipants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}
	if !allowed {
		return utility.NewForbiddenError("you have no permission to revert")
	}

	revertStatus, err := pdlrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	if err := pdlrs.createReviewWithHistory(
		currentWorkflow.ID,
		decisionLogID,
		currentLevel,
		revertStatus.ID,
		comment,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.ProjectDecisionLogWorkflow
	if previousLevel > 0 {
		pendingStatus, err := pdlrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		previousLevelWorkflow, err = pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, previousLevel, tx)
		if err != nil {
			return err
		}

		if err := pdlrs.workflowSvc.DeleteReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := pdlrs.workflowSvc.UpdateWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			return err
		}
	}

	var submittedByID *uint
	if previousLevel == 0 {
		submittedByID = pdl.SubmittedBYID

		pdl.WorkflowStatusID = &revertStatus.ID
		pdl.NextReviewLevel = nil
		pdl.SubmittedBYID = nil
		pdl.SubmittedAt = nil
		pdl.ApprovedBYID = nil
		pdl.ApprovedAt = nil
		pdl.RejectedBYID = nil
		pdl.RejectedAt = nil
	} else {
		pendingStatus, err := pdlrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}
		pdl.WorkflowStatusID = &pendingStatus.ID
		pdl.NextReviewLevel = &previousLevel
	}
	pdl.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
		return err
	}

	if previousLevel > 0 {
		notify = func() {
			_ = pdlrs.logService.LogActivity(
				models.UpdateAction,
				models.ProjectDecisionLog{}.TableName(),
				decisionLogID,
				pdl,
				currentUserID,
				"project", "project_register", "", nil,
			)

			pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
			if err == nil {
				pdlrs.notifyNextLevelUsers(projectID, decisionLogID, previousLevelWorkflow, pdlDetail, currentUserID)
			}
		}
	} else {
		notify = func() {
			_ = pdlrs.logService.LogActivity(
				models.UpdateAction,
				models.ProjectDecisionLog{}.TableName(),
				decisionLogID,
				pdl,
				currentUserID,
				"project", "project_register", "", nil,
			)

			if submittedByID != nil && *submittedByID != 0 {
				pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
				if err == nil {
					pdlrs.sendRevertNotification(projectID, decisionLogID, currentUserID, pdlDetail)
				}
			}
		}
	}

	return err
}

func (pdlrs *ProjectDecisionLogReviewService) Reject(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) (err error) {
	tx := pdlrs.repo.BeginTransaction()

	var notify func()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		if err = tx.Commit(); err != nil {
			return
		}

		if notify != nil {
			go notify()
		}
	}()

	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel == nil ||
		(pdl.WorkflowStatusID != nil && *pdl.WorkflowStatusID == workflowStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"project decision log",
		)
	}

	currentReviewLevel := *pdl.NextReviewLevel

	currentWorkflow, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, currentReviewLevel, tx)
	if err != nil {
		return err
	}

	var rejectingParticipant *models.ProjectDecisionLogWorkflowParticipant

	for _, participant := range currentWorkflow.ProjectDecisionLogWorkflowParticipants {
		if participant.UserID == currentUserID {
			rejectingParticipant = participant
			break
		}
	}

	if rejectingParticipant == nil {
		return utility.NewForbiddenError("you have no permission to reject")
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	if err := pdlrs.createReviewWithHistory(
		currentWorkflow.ID,
		decisionLogID,
		currentReviewLevel,
		workflowStatus.ID,
		comment,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	currentWorkflow.WorkflowStatusID = workflowStatus.ID

	if err := pdlrs.workflowSvc.UpdateProjectDecisionLogWorkflow(currentWorkflow, tx); err != nil {
		return err
	}

	var submittedByID uint
	if pdl.SubmittedBYID != nil {
		submittedByID = *pdl.SubmittedBYID
	} else if pdl.RequestedByUserID != 0 {
		submittedByID = pdl.RequestedByUserID
	}

	now := time.Now()
	pdl.WorkflowStatusID = &workflowStatus.ID
	pdl.RejectedBYID = &currentUserID
	pdl.RejectedAt = &now
	pdl.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
		return err
	}

	notify = func() {
		_ = pdlrs.logService.LogActivity(
			models.RejectAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)

		if submittedByID != 0 {
			pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
			if err == nil {
				pdlrs.sendRejectNotification(projectID, decisionLogID, currentUserID, pdlDetail)
			}
		}
	}

	return err
}

func (pdlrs *ProjectDecisionLogReviewService) sendApprovalNotification(projectID uint, decisionLogID uint, currentUserID uint, pdlDetail *projectDTOs.ProjectDecisionLogDetailDTO) {
	go func() {
		if pdlDetail == nil || pdlDetail.RequestedByUser == nil || pdlDetail.RequestedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"decisionLogId":%d}`, projectID, decisionLogID)
		entity := "project-decision-log-review"
		action := "approval"

		approvedByUser, err := pdlrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{"error": err.Error()})
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(pdlDetail.RequestedByUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Project Decision Log Approval",
			HeadingAr:            "الموافقة على سجل قرارات المشروع",
			ReferenceNumber:      pdlDetail.ReferenceNumber,
			TitleEn:              pdlDetail.Decision,
			TitleAr:              pdlDetail.Decision,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Decision":              pdlDetail.Decision,
			"FirstNameEn":           pdlDetail.RequestedByUser.FirstNameEN,
			"LastNameEn":            pdlDetail.RequestedByUser.LastNameEN,
			"ReferenceNumber":       pdlDetail.ReferenceNumber,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/%d/decision-logs/%d", appURL, projectID, decisionLogID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{pdlDetail.RequestedByUser.Email},
			Subject:              "Project Decision Log Approval Notification",
			TemplateRelativePath: "project/templates/email/project_decision_log_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = pdlrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_decision_log")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{"error": err.Error()})
		}
	}()
}

func (pdlrs *ProjectDecisionLogReviewService) notifyNextLevelUsers(projectID uint, decisionLogID uint, nextWorkflow *models.ProjectDecisionLogWorkflow, pdlDetail *projectDTOs.ProjectDecisionLogDetailDTO, currentUserID uint) {
	go func() {
		if nextWorkflow == nil || pdlDetail == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"decisionLogId":%d}`, projectID, decisionLogID)
		entity := "project-decision-log-review"
		action := "review-user-assign"

		assignedBy, err := pdlrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
			return
		}

		for _, participant := range nextWorkflow.ProjectDecisionLogWorkflowParticipants {
			reviewUser, err := pdlrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Project Decision Log Review Assigned",
				HeadingAr:            "تم تعيين مراجعة سجل قرارات المشروع",
				ReferenceNumber:      pdlDetail.ReferenceNumber,
				TitleEn:              pdlDetail.Decision,
				TitleAr:              pdlDetail.Decision,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           reviewUser.FirstNameEn,
				"LastNameEn":            reviewUser.LastNameEn,
				"Decision":              pdlDetail.Decision,
				"ReferenceNumber":       pdlDetail.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/projects/%d/decision-logs/%d", appURL, projectID, decisionLogID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{reviewUser.Email},
				Subject:              "Project Decision Log Review Assigned",
				TemplateRelativePath: "project/templates/email/project_decision_log_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			_ = pdlrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_decision_log")
		}
	}()
}

func (pdlrs *ProjectDecisionLogReviewService) sendResubmitNotification(projectID uint, decisionLogID uint, workflowLevel *models.ProjectDecisionLogWorkflow, pdlDetail *projectDTOs.ProjectDecisionLogDetailDTO, currentUserID uint) {
	go func() {
		if workflowLevel == nil || pdlDetail == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"decisionLogId":%d}`, projectID, decisionLogID)
		entity := "project-decision-log-review"
		action := "resubmitted"

		resubmittedBy, err := pdlrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch resubmitting user", map[string]interface{}{"error": err.Error()})
			return
		}

		for _, participant := range workflowLevel.ProjectDecisionLogWorkflowParticipants {
			reviewUser, err := pdlrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Project Decision Log Review Resubmitted",
				HeadingAr:            "تمت إعادة إرسال مراجعة سجل قرارات المشروع",
				ReferenceNumber:      pdlDetail.ReferenceNumber,
				TitleEn:              pdlDetail.Decision,
				TitleAr:              pdlDetail.Decision,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":              reviewUser.FirstNameEn,
				"LastNameEn":               reviewUser.LastNameEn,
				"Decision":                 pdlDetail.Decision,
				"ReferenceNumber":          pdlDetail.ReferenceNumber,
				"ResubmittedByFirstNameEn": resubmittedBy.FirstNameEn,
				"ResubmittedByLastNameEn":  resubmittedBy.LastNameEn,
				"ResubmittedAt":            time.Now().Format("2006-01-02"),
				"RedirectURL":              fmt.Sprintf("%s/projects/%d/decision-logs/%d", appURL, projectID, decisionLogID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{reviewUser.Email},
				Subject:              "Project Decision Log Review Resubmitted",
				TemplateRelativePath: "project/templates/email/project_decision_log_workflow_resubmit_user_mail.html",
				EmailData:            emailData,
			}

			_ = pdlrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_decision_log")
		}
	}()
}

func (pdlrs *ProjectDecisionLogReviewService) sendRevertNotification(projectID uint, decisionLogID uint, currentUserID uint, pdlDetail *projectDTOs.ProjectDecisionLogDetailDTO) {
	go func() {
		if pdlDetail == nil || pdlDetail.RequestedByUser == nil || pdlDetail.RequestedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"decisionLogId":%d}`, projectID, decisionLogID)
		entity := "project-decision-log-review"
		action := "reverted"

		revertedByUser, err := pdlrs.userSvc.GetByID(currentUserID)
		if err != nil {
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(pdlDetail.RequestedByUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Project Decision Log Reverted",
			HeadingAr:            "تم إرجاع سجل قرارات المشروع",
			ReferenceNumber:      pdlDetail.ReferenceNumber,
			TitleEn:              pdlDetail.Decision,
			TitleAr:              pdlDetail.Decision,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Decision":              pdlDetail.Decision,
			"FirstNameEn":           pdlDetail.RequestedByUser.FirstNameEN,
			"LastNameEn":            pdlDetail.RequestedByUser.LastNameEN,
			"ReferenceNumber":       pdlDetail.ReferenceNumber,
			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
			"RevertedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/%d/decision-logs/%d", appURL, projectID, decisionLogID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{pdlDetail.RequestedByUser.Email},
			Subject:              "Project Decision Log Reverted Notification",
			TemplateRelativePath: "project/templates/email/project_decision_log_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		_ = pdlrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_decision_log")
	}()
}

func (pdlrs *ProjectDecisionLogReviewService) sendRejectNotification(projectID uint, decisionLogID uint, currentUserID uint, pdlDetail *projectDTOs.ProjectDecisionLogDetailDTO) {
	go func() {
		if pdlDetail == nil || pdlDetail.RequestedByUser == nil || pdlDetail.RequestedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"decisionLogId":%d}`, projectID, decisionLogID)
		entity := "project-decision-log-review"
		action := "rejected"

		rejectedByUser, err := pdlrs.userSvc.GetByID(currentUserID)
		if err != nil {
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(pdlDetail.RequestedByUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Project Decision Log Rejected",
			HeadingAr:            "تم رفض سجل قرارات المشروع",
			ReferenceNumber:      pdlDetail.ReferenceNumber,
			TitleEn:              pdlDetail.Decision,
			TitleAr:              pdlDetail.Decision,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Decision":              pdlDetail.Decision,
			"FirstNameEn":           pdlDetail.RequestedByUser.FirstNameEN,
			"LastNameEn":            pdlDetail.RequestedByUser.LastNameEN,
			"ReferenceNumber":       pdlDetail.ReferenceNumber,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/%d/decision-logs/%d", appURL, projectID, decisionLogID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{pdlDetail.RequestedByUser.Email},
			Subject:              "Project Decision Log Rejected Notification",
			TemplateRelativePath: "project/templates/email/project_decision_log_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		_ = pdlrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_decision_log")
	}()
}
