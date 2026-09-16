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

type ProjectChangeRequestReviewServiceInterface interface {
	Submit(projectID uint, changeRequestID uint, currentUserID uint) error
	Resubmit(projectID uint, changeRequestID uint, currentUserID uint) error
	Approve(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) error
	Revert(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) error
	Reject(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) error
}

type ProjectChangeRequestReviewService struct {
	repo              repositories.ProjectChangeRequestReviewRepositoryInterface
	logService        logServices.ActivityLogServiceInterface
	projectCRSvc      ProjectChangeRequestServiceInterface
	workflowSvc       ProjectChangeRequestWorkflowServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
	projectSvc        ProjectServiceInterface
}

func NewProjectChangeRequestReviewService(
	repo repositories.ProjectChangeRequestReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	projectCRSvc ProjectChangeRequestServiceInterface,
	workflowSvc ProjectChangeRequestWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	projectSvc ProjectServiceInterface,
) ProjectChangeRequestReviewServiceInterface {
	return &ProjectChangeRequestReviewService{
		repo:              repo,
		logService:        logService,
		projectCRSvc:      projectCRSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
		projectSvc:        projectSvc,
	}
}

func (pcrrs *ProjectChangeRequestReviewService) createReviewWithHistory(
	workflowID uint,
	changeRequestID uint,
	level int,
	workflowStatusID uint,
	comment *string,
	currentUserID uint,
	tx utility.DatabaseTransaction,
) error {
	review := &models.ProjectChangeRequestWorkflowReview{
		ProjectChangeRequestWorkflowID: workflowID,
		WorkflowStatusID:               workflowStatusID,
		Comment:                        comment,
		CreatedBYID:                    currentUserID,
	}

	if err := pcrrs.repo.CreateReview(review, tx); err != nil {
		return err
	}

	history := &models.ProjectChangeRequestWorkflowHistory{
		ProjectChangeRequestID:         changeRequestID,
		ProjectChangeRequestWorkflowID: workflowID,
		WorkflowStatusID:               workflowStatusID,
		Level:                          level,
		Comment:                        comment,
		CreatedBYID:                    currentUserID,
	}

	return pcrrs.repo.CreateHistory(history, tx)
}

func (pcrrs *ProjectChangeRequestReviewService) Submit(projectID uint, changeRequestID uint, currentUserID uint) (err error) {
	tx := pcrrs.repo.BeginTransaction()

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

	if err := pcrrs.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	approvedStatus, err := pcrrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	rejectedStatus, err := pcrrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel != nil {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already submitted"),
			"workflow_already_submitted",
			"project change request",
		)
	}

	if pcr.WorkflowStatusID != nil && (*pcr.WorkflowStatusID == approvedStatus.ID || *pcr.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already completed"),
			"workflow_already_completed",
			"project change request",
		)
	}

	if err := pcrrs.workflowSvc.DeleteAllWorkflowsForProjectChangeRequest(changeRequestID, tx); err != nil {
		return err
	}

	if err := pcrrs.workflowSvc.CreateProjectChangeRequestWorkflow(changeRequestID, currentUserID, tx); err != nil {
		return err
	}

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	currentDate := time.Now()
	level := 1
	pcr.WorkflowStatusID = &workflowStatus.ID
	pcr.SubmittedBYID = &currentUserID
	pcr.SubmittedAt = &currentDate
	pcr.NextReviewLevel = &level
	pcr.ApprovedBYID = nil
	pcr.ApprovedAt = nil
	pcr.RejectedBYID = nil
	pcr.RejectedAt = nil
	pcr.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
		return err
	}

	workflowLevel, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, level, tx)
	if err != nil {
		return err
	}

	if err := pcrrs.createReviewWithHistory(workflowLevel.ID, changeRequestID, level, workflowStatus.ID, nil, currentUserID, tx); err != nil {
		return err
	}

	notify = func() {
		_ = pcrrs.logService.LogActivity(
			models.CreateAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)

		pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
		if err == nil {
			pcrrs.notifyNextLevelUsers(projectID, changeRequestID, workflowLevel, pcrDetail, currentUserID)
		}
	}

	return err
}

func (pcrrs *ProjectChangeRequestReviewService) Resubmit(projectID uint, changeRequestID uint, currentUserID uint) (err error) {
	tx := pcrrs.repo.BeginTransaction()

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

	if err := pcrrs.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	rejectedStatus, err := pcrrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	revertedStatus, err := pcrrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	if pcr.WorkflowStatusID == nil || (*pcr.WorkflowStatusID != rejectedStatus.ID && *pcr.WorkflowStatusID != revertedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("project change request has not been rejected or reverted"),
			"workflow_not_rejected",
			"project change request",
		)
	}

	if err := pcrrs.workflowSvc.DeleteAllWorkflowsForProjectChangeRequest(changeRequestID, tx); err != nil {
		return err
	}

	if err := pcrrs.workflowSvc.CreateProjectChangeRequestWorkflow(changeRequestID, currentUserID, tx); err != nil {
		return err
	}

	submitStatus, err := pcrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	level := 1
	now := time.Now()

	pcr.WorkflowStatusID = &submitStatus.ID
	pcr.SubmittedBYID = &currentUserID
	pcr.SubmittedAt = &now
	pcr.NextReviewLevel = &level
	pcr.ApprovedBYID = nil
	pcr.ApprovedAt = nil
	pcr.RejectedBYID = nil
	pcr.RejectedAt = nil
	pcr.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
		return err
	}

	workflowLevel, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, level, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByType("resubmit")
	if err != nil {
		return err
	}

	if err := pcrrs.createReviewWithHistory(
		workflowLevel.ID,
		changeRequestID,
		level,
		workflowStatus.ID,
		nil,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	notify = func() {
		_ = pcrrs.logService.LogActivity(
			models.UpdateAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)

		pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
		if err == nil {
			pcrrs.sendResubmitNotification(projectID, changeRequestID, workflowLevel, pcrDetail, currentUserID)
		}
	}

	return err
}

func (pcrrs *ProjectChangeRequestReviewService) Approve(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) (err error) {
	tx := pcrrs.repo.BeginTransaction()

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

	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	rejectedStatus, err := pcrrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel == nil || (pcr.WorkflowStatusID != nil && *pcr.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"project change request",
		)
	}

	currentApprovalLevel := *pcr.NextReviewLevel
	currentWorkflow, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	pendingStatus, err := pcrrs.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.ProjectChangeRequestWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectChangeRequestWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pcrrs.createReviewWithHistory(currentWorkflow.ID, changeRequestID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			return err
		}

		projectChangeRequestWorkflow, err := pcrrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			return err
		}

		projectChangeRequestWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(projectChangeRequestWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := pcrrs.workflowSvc.GetNextLevel(changeRequestID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			pcr.WorkflowStatusID = &workflowStatus.ID
			pcr.ApprovedBYID = &currentUserID
			pcr.ApprovedAt = &approvedAt
			pcr.RejectedBYID = nil
			pcr.RejectedAt = nil
			pcr.NextReviewLevel = nil
			pcr.UpdatedBYID = &currentUserID

			if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
				return err
			}

			if err := pcrrs.applyImpactAssessmentToProject(projectID, changeRequestID, currentUserID, tx); err != nil {
				return err
			}

			notify = func() {
				_ = pcrrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectChangeRequest{}.TableName(),
					changeRequestID,
					pcr,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
				if err == nil {
					pcrrs.sendApprovalNotification(projectID, changeRequestID, currentUserID, pcrDetail)
				}
			}
		} else {
			pcr.WorkflowStatusID = &pendingStatus.ID
			pcr.NextReviewLevel = &nextWorkflowLevel.Level
			pcr.UpdatedBYID = &currentUserID

			if err := pcrrs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
				return err
			}

			if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
				return err
			}

			nextWf, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}

			notify = func() {
				_ = pcrrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectChangeRequest{}.TableName(),
					changeRequestID,
					pcr,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
				if err == nil {
					pcrrs.notifyNextLevelUsers(projectID, changeRequestID, nextWf, pcrDetail, currentUserID)
				}
			}
		}

	case "multiple-users", "role":
		var approvingParticipant *models.ProjectChangeRequestWorkflowParticipant
		for _, participant := range currentWorkflow.ProjectChangeRequestWorkflowParticipants {
			if participant.UserID == currentUserID {
				approvingParticipant = participant
				break
			}
		}
		if approvingParticipant == nil {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pcrrs.createReviewWithHistory(currentWorkflow.ID, changeRequestID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			var allApproved bool
			allApproved, err = pcrrs.workflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ProjectChangeRequestWorkflowParticipants, tx)
			if err != nil || !allApproved {
				if err == nil {
					notify = func() {
						_ = pcrrs.logService.LogActivity(
							models.ApproveAction,
							models.ProjectChangeRequest{}.TableName(),
							changeRequestID,
							pcr,
							currentUserID,
							"project", "project_register", "", nil,
						)
					}
				}
				return err
			}
		}

		projectChangeRequestWorkflow, err := pcrrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			return err
		}

		projectChangeRequestWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(projectChangeRequestWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := pcrrs.workflowSvc.GetNextLevel(changeRequestID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			pcr.WorkflowStatusID = &workflowStatus.ID
			pcr.ApprovedBYID = &currentUserID
			pcr.ApprovedAt = &approvedAt
			pcr.RejectedBYID = nil
			pcr.RejectedAt = nil
			pcr.NextReviewLevel = nil
			pcr.UpdatedBYID = &currentUserID

			if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
				return err
			}

			if err := pcrrs.applyImpactAssessmentToProject(projectID, changeRequestID, currentUserID, tx); err != nil {
				return err
			}

			notify = func() {
				_ = pcrrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectChangeRequest{}.TableName(),
					changeRequestID,
					pcr,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
				if err == nil {
					pcrrs.sendApprovalNotification(projectID, changeRequestID, currentUserID, pcrDetail)
				}
			}
		} else {
			pcr.WorkflowStatusID = &pendingStatus.ID
			pcr.NextReviewLevel = &nextWorkflowLevel.Level
			pcr.UpdatedBYID = &currentUserID

			if err := pcrrs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
				return err
			}

			if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
				return err
			}

			nextWf, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}

			notify = func() {
				_ = pcrrs.logService.LogActivity(
					models.ApproveAction,
					models.ProjectChangeRequest{}.TableName(),
					changeRequestID,
					pcr,
					currentUserID,
					"project", "project_register", "", nil,
				)

				pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
				if err == nil {
					pcrrs.notifyNextLevelUsers(projectID, changeRequestID, nextWf, pcrDetail, currentUserID)
				}
			}
		}
	}

	return err
}

func (pcrrs *ProjectChangeRequestReviewService) Revert(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) (err error) {
	tx := pcrrs.repo.BeginTransaction()

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

	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	rejectedStatus, err := pcrrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel == nil || (pcr.WorkflowStatusID != nil && *pcr.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "project change request")
	}

	currentLevel := *pcr.NextReviewLevel

	currentWorkflow, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, currentLevel, tx)
	if err != nil {
		return err
	}

	allowed := false
	for _, participant := range currentWorkflow.ProjectChangeRequestWorkflowParticipants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}
	if !allowed {
		return utility.NewForbiddenError("you have no permission to revert")
	}

	revertStatus, err := pcrrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	if err := pcrrs.createReviewWithHistory(
		currentWorkflow.ID,
		changeRequestID,
		currentLevel,
		revertStatus.ID,
		comment,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	currentWorkflow.WorkflowStatusID = revertStatus.ID

	if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(currentWorkflow, tx); err != nil {
		return err
	}

	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.ProjectChangeRequestWorkflow
	if previousLevel > 0 {
		pendingStatus, err := pcrrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		previousLevelWorkflow, err = pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, previousLevel, tx)
		if err != nil {
			return err
		}

		if err := pcrrs.workflowSvc.DeleteReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := pcrrs.workflowSvc.UpdateWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			return err
		}
	}

	var submittedByID *uint
	if previousLevel == 0 {
		submittedByID = pcr.SubmittedBYID

		pcr.WorkflowStatusID = &revertStatus.ID
		pcr.NextReviewLevel = nil
		pcr.SubmittedBYID = nil
		pcr.SubmittedAt = nil
		pcr.ApprovedBYID = nil
		pcr.ApprovedAt = nil
		pcr.RejectedBYID = nil
		pcr.RejectedAt = nil
	} else {
		pendingStatus, err := pcrrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}
		pcr.WorkflowStatusID = &pendingStatus.ID
		pcr.NextReviewLevel = &previousLevel
	}
	pcr.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
		return err
	}

	if previousLevel > 0 {
		notify = func() {
			_ = pcrrs.logService.LogActivity(
				models.UpdateAction,
				models.ProjectChangeRequest{}.TableName(),
				changeRequestID,
				pcr,
				currentUserID,
				"project", "project_register", "", nil,
			)

			pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
			if err == nil {
				pcrrs.notifyNextLevelUsers(projectID, changeRequestID, previousLevelWorkflow, pcrDetail, currentUserID)
			}
		}
	} else {
		notify = func() {
			_ = pcrrs.logService.LogActivity(
				models.UpdateAction,
				models.ProjectChangeRequest{}.TableName(),
				changeRequestID,
				pcr,
				currentUserID,
				"project", "project_register", "", nil,
			)

			if submittedByID != nil && *submittedByID != 0 {
				pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
				if err == nil {
					pcrrs.sendRevertNotification(projectID, changeRequestID, currentUserID, pcrDetail)
				}
			}
		}
	}

	return err
}

func (pcrrs *ProjectChangeRequestReviewService) Reject(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) (err error) {
	tx := pcrrs.repo.BeginTransaction()

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

	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel == nil ||
		(pcr.WorkflowStatusID != nil && *pcr.WorkflowStatusID == workflowStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"project change request",
		)
	}

	currentReviewLevel := *pcr.NextReviewLevel

	currentWorkflow, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, currentReviewLevel, tx)
	if err != nil {
		return err
	}

	var rejectingParticipant *models.ProjectChangeRequestWorkflowParticipant

	for _, participant := range currentWorkflow.ProjectChangeRequestWorkflowParticipants {
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

	if err := pcrrs.createReviewWithHistory(
		currentWorkflow.ID,
		changeRequestID,
		currentReviewLevel,
		workflowStatus.ID,
		comment,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	currentWorkflow.WorkflowStatusID = workflowStatus.ID

	if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(currentWorkflow, tx); err != nil {
		return err
	}
	if err := pcrrs.workflowSvc.UpdateWorkflowStatus(currentWorkflow.ID, workflowStatus.ID, tx); err != nil {
		return err
	}

	var submittedByID uint
	if pcr.SubmittedBYID != nil {
		submittedByID = *pcr.SubmittedBYID
	} else if pcr.RequestedByUserID != 0 {
		submittedByID = pcr.RequestedByUserID
	}

	now := time.Now()
	pcr.WorkflowStatusID = &workflowStatus.ID
	pcr.RejectedBYID = &currentUserID
	pcr.RejectedAt = &now
	pcr.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
		return err
	}

	notify = func() {
		_ = pcrrs.logService.LogActivity(
			models.RejectAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)

		if submittedByID != 0 {
			pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
			if err == nil {
				pcrrs.sendRejectNotification(projectID, changeRequestID, currentUserID, pcrDetail)
			}
		}
	}

	return err
}

func (pcrrs *ProjectChangeRequestReviewService) sendApprovalNotification(projectID uint, changeRequestID uint, currentUserID uint, pcrDetail *projectDTOs.ProjectChangeRequestDetailDTO) {
	go func() {
		if pcrDetail == nil || pcrDetail.RequestedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"changeRequestId":%d}`, projectID, changeRequestID)
		entity := "project-change-request-review"
		action := "approval"

		approvedByUser, err := pcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{"error": err.Error()})
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(pcrDetail.RequestedByUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Project Change Request Approval",
			HeadingAr:            "الموافقة على طلب تغيير المشروع",
			ReferenceNumber:      pcrDetail.ReferenceNumber,
			TitleEn:              pcrDetail.Title,
			TitleAr:              pcrDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 pcrDetail.Title,
			"FirstNameEn":           pcrDetail.RequestedByUser.FirstNameEN,
			"LastNameEn":            pcrDetail.RequestedByUser.LastNameEN,
			"ReferenceNumber":       pcrDetail.ReferenceNumber,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/%d/change-requests/%d", appURL, projectID, changeRequestID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{pcrDetail.RequestedByUser.Email},
			Subject:              "Project Change Request Approval Notification",
			TemplateRelativePath: "project/templates/email/project_change_request_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{"error": err.Error()})
		}
	}()
}

func (pcrrs *ProjectChangeRequestReviewService) notifyNextLevelUsers(projectID uint, changeRequestID uint, nextWorkflow *models.ProjectChangeRequestWorkflow, pcrDetail *projectDTOs.ProjectChangeRequestDetailDTO, currentUserID uint) {
	go func() {
		if nextWorkflow == nil || pcrDetail == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"changeRequestId":%d}`, projectID, changeRequestID)
		entity := "project-change-request-review"
		action := "review-user-assign"

		assignedBy, err := pcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
			return
		}

		for _, participant := range nextWorkflow.ProjectChangeRequestWorkflowParticipants {
			reviewUser, err := pcrrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Project Change Request Review Assigned",
				HeadingAr:            "تم تعيين مراجعة طلب تغيير المشروع",
				ReferenceNumber:      pcrDetail.ReferenceNumber,
				TitleEn:              pcrDetail.Title,
				TitleAr:              pcrDetail.Title,
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
				"Title":                 pcrDetail.Title,
				"ReferenceNumber":       pcrDetail.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/projects/%d/change-requests/%d", appURL, projectID, changeRequestID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{reviewUser.Email},
				Subject:              "Project Change Request Review Assigned",
				TemplateRelativePath: "project/templates/email/project_change_request_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
		}
	}()
}

func (pcrrs *ProjectChangeRequestReviewService) sendResubmitNotification(projectID uint, changeRequestID uint, workflowLevel *models.ProjectChangeRequestWorkflow, pcrDetail *projectDTOs.ProjectChangeRequestDetailDTO, currentUserID uint) {
	go func() {
		if workflowLevel == nil || pcrDetail == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"changeRequestId":%d}`, projectID, changeRequestID)
		entity := "project-change-request-review"
		action := "resubmitted"

		resubmittedBy, err := pcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch resubmitting user", map[string]interface{}{"error": err.Error()})
			return
		}

		for _, participant := range workflowLevel.ProjectChangeRequestWorkflowParticipants {
			reviewUser, err := pcrrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Project Change Request Review Resubmitted",
				HeadingAr:            "تمت إعادة إرسال مراجعة طلب تغيير المشروع",
				ReferenceNumber:      pcrDetail.ReferenceNumber,
				TitleEn:              pcrDetail.Title,
				TitleAr:              pcrDetail.Title,
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
				"Title":                    pcrDetail.Title,
				"ReferenceNumber":          pcrDetail.ReferenceNumber,
				"ResubmittedByFirstNameEn": resubmittedBy.FirstNameEn,
				"ResubmittedByLastNameEn":  resubmittedBy.LastNameEn,
				"ResubmittedAt":            time.Now().Format("2006-01-02"),
				"RedirectURL":              fmt.Sprintf("%s/projects/%d/change-requests/%d", appURL, projectID, changeRequestID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{reviewUser.Email},
				Subject:              "Project Change Request Review Resubmitted",
				TemplateRelativePath: "project/templates/email/project_change_request_workflow_resubmit_user_mail.html",
				EmailData:            emailData,
			}

			_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
		}
	}()
}

func (pcrrs *ProjectChangeRequestReviewService) sendRevertNotification(projectID uint, changeRequestID uint, currentUserID uint, pcrDetail *projectDTOs.ProjectChangeRequestDetailDTO) {
	go func() {
		if pcrDetail == nil || pcrDetail.RequestedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"changeRequestId":%d}`, projectID, changeRequestID)
		entity := "project-change-request-review"
		action := "reverted"

		revertedByUser, err := pcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(pcrDetail.RequestedByUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Project Change Request Reverted",
			HeadingAr:            "تم إرجاع طلب تغيير المشروع",
			ReferenceNumber:      pcrDetail.ReferenceNumber,
			TitleEn:              pcrDetail.Title,
			TitleAr:              pcrDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 pcrDetail.Title,
			"FirstNameEn":           pcrDetail.RequestedByUser.FirstNameEN,
			"LastNameEn":            pcrDetail.RequestedByUser.LastNameEN,
			"ReferenceNumber":       pcrDetail.ReferenceNumber,
			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
			"RevertedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/%d/change-requests/%d", appURL, projectID, changeRequestID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{pcrDetail.RequestedByUser.Email},
			Subject:              "Project Change Request Reverted Notification",
			TemplateRelativePath: "project/templates/email/project_change_request_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
	}()
}

func (pcrrs *ProjectChangeRequestReviewService) sendRejectNotification(projectID uint, changeRequestID uint, currentUserID uint, pcrDetail *projectDTOs.ProjectChangeRequestDetailDTO) {
	go func() {
		if pcrDetail == nil || pcrDetail.RequestedByUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d,"changeRequestId":%d}`, projectID, changeRequestID)
		entity := "project-change-request-review"
		action := "rejected"

		rejectedByUser, err := pcrrs.userSvc.GetByID(currentUserID)
		if err != nil {
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(pcrDetail.RequestedByUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Project Change Request Rejected",
			HeadingAr:            "تم رفض طلب تغيير المشروع",
			ReferenceNumber:      pcrDetail.ReferenceNumber,
			TitleEn:              pcrDetail.Title,
			TitleAr:              pcrDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 pcrDetail.Title,
			"FirstNameEn":           pcrDetail.RequestedByUser.FirstNameEN,
			"LastNameEn":            pcrDetail.RequestedByUser.LastNameEN,
			"ReferenceNumber":       pcrDetail.ReferenceNumber,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/%d/change-requests/%d", appURL, projectID, changeRequestID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{pcrDetail.RequestedByUser.Email},
			Subject:              "Project Change Request Rejected Notification",
			TemplateRelativePath: "project/templates/email/project_change_request_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
	}()
}

func (pcrrs *ProjectChangeRequestReviewService) applyImpactAssessmentToProject(projectID uint, changeRequestID uint, currentUserID uint, tx utility.DatabaseTransaction) error {
	project, err := pcrrs.projectSvc.GetProjectByID(projectID, tx)
	if err != nil {
		return err
	}

	impactAssessments, err := pcrrs.projectCRSvc.GetImpactAssessments(changeRequestID, tx)
	if err != nil {
		return err
	}

	projectUpdated := false
	for _, ia := range impactAssessments {
		if ia.ProjectImpactAssessmentType != nil {
			switch ia.ProjectImpactAssessmentType.Type {
			case "timeline":
				if ia.RevisedDate != nil {
					project.EndDate = *ia.RevisedDate
					projectUpdated = true
				}
			case "budget":
				if ia.RevisedBudget != nil {
					project.Budget = ia.RevisedBudget
					projectUpdated = true
				}
			case "resources":
				if ia.RevisedResourceCount != nil {
					project.NumberOfResources = ia.RevisedResourceCount
					projectUpdated = true
				}
			}
		}
	}

	if projectUpdated {
		now := time.Now()
		project.UpdatedBYID = &currentUserID
		project.UpdatedAt = &now
		if err := pcrrs.projectSvc.UpdateProject(project, tx); err != nil {
			return err
		}
	}

	return nil
}
