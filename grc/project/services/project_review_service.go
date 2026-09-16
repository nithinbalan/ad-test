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

type ProjectReviewServiceInterface interface {
	Submit(projectID uint, currentUserID uint) error
	Resubmit(projectID uint, currentUserID uint) error
	Approve(projectID uint, currentUserID uint, form *projectDTOs.ProjectWorkflowReviewFormDTO) error
	Revert(projectID uint, currentUserID uint, form *projectDTOs.ProjectWorkflowReviewFormDTO) error
	Reject(projectID uint, currentUserID uint, form *projectDTOs.ProjectWorkflowReviewFormDTO) error
}

type ProjectReviewService struct {
	repo              repositories.ProjectReviewRepositoryInterface
	logService        logServices.ActivityLogServiceInterface
	projectSvc        ProjectServiceInterface
	workflowSvc       ProjectWorkflowServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewProjectReviewService(
	repo repositories.ProjectReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	projectSvc ProjectServiceInterface,
	workflowSvc ProjectWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProjectReviewServiceInterface {
	return &ProjectReviewService{
		repo:              repo,
		logService:        logService,
		projectSvc:        projectSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
	}
}

func (prs *ProjectReviewService) createReviewWithHistory(
	workflowID uint,
	projectID uint,
	level int,
	workflowStatusID uint,
	comment *string,
	currentUserID uint,
	tx utility.DatabaseTransaction,
) error {
	review := &models.ProjectWorkflowReview{
		ProjectWorkflowID: workflowID,
		WorkflowStatusID:  workflowStatusID,
		Comment:           comment,
		CreatedBYID:       currentUserID,
	}

	if err := prs.repo.CreateReview(review, tx); err != nil {
		return err
	}

	history := &models.ProjectWorkflowHistory{
		ProjectID:         projectID,
		ProjectWorkflowID: workflowID,
		WorkflowStatusID:  workflowStatusID,
		Level:             level,
		Comment:           comment,
		CreatedBYID:       currentUserID,
	}

	return prs.repo.CreateHistory(history, tx)
}

func (prs *ProjectReviewService) Submit(projectID uint, currentUserID uint) (err error) {
	tx := prs.repo.BeginTransaction()

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

	if err := prs.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	project, err := prs.projectSvc.GetProjectByID(projectID, tx)
	if err != nil {
		return err
	}

	approvedStatus, err := prs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	rejectedStatus, err := prs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if project.NextReviewLevel != nil {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already submitted"),
			"workflow_already_submitted",
			"project",
		)
	}

	if project.WorkflowStatusID != nil && (*project.WorkflowStatusID == approvedStatus.ID || *project.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already completed"),
			"workflow_already_completed",
			"project",
		)
	}

	if err := prs.workflowSvc.DeleteAllWorkflowsForProject(projectID, tx); err != nil {
		return err
	}

	if err := prs.workflowSvc.CreateProjectWorkflow(projectID, currentUserID, tx); err != nil {
		return err
	}

	workflowStatus, err := prs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	currentDate := time.Now()
	level := 1
	project.WorkflowStatusID = &workflowStatus.ID
	project.SubmittedBYID = &currentUserID
	project.SubmittedAt = &currentDate
	project.NextReviewLevel = &level
	project.ApprovedBYID = nil
	project.ApprovedAt = nil
	project.RejectedBYID = nil
	project.RejectedAt = nil
	project.UpdatedBYID = &currentUserID

	if err := prs.repo.UpdateProject(project, tx); err != nil {
		return err
	}

	workflowLevel, err := prs.workflowSvc.GetWorkflowByLevel(projectID, level, tx)
	if err != nil {
		return err
	}

	if err := prs.createReviewWithHistory(workflowLevel.ID, projectID, level, workflowStatus.ID, nil, currentUserID, tx); err != nil {
		return err
	}

	notify = func() {
		_ = prs.logService.LogActivity(
			models.CreateAction,
			models.Project{}.TableName(),
			projectID,
			project,
			currentUserID,
			"project", "project_register", "", nil,
		)

		projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
		if err == nil {
			prs.notifyNextLevelUsers(projectID, workflowLevel, projectDetail, currentUserID)
		}
	}

	return err
}

func (prs *ProjectReviewService) Resubmit(projectID uint, currentUserID uint) (err error) {
	tx := prs.repo.BeginTransaction()

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

	if err := prs.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	project, err := prs.projectSvc.GetProjectByID(projectID, tx)
	if err != nil {
		return err
	}

	rejectedStatus, err := prs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	revertedStatus, err := prs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	if project.WorkflowStatusID == nil || (*project.WorkflowStatusID != rejectedStatus.ID && *project.WorkflowStatusID != revertedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("project has not been rejected or reverted"),
			"workflow_not_rejected",
			"project",
		)
	}

	if err := prs.workflowSvc.DeleteAllWorkflowsForProject(projectID, tx); err != nil {
		return err
	}

	if err := prs.workflowSvc.CreateProjectWorkflow(projectID, currentUserID, tx); err != nil {
		return err
	}

	submitStatus, err := prs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	level := 1
	now := time.Now()

	project.WorkflowStatusID = &submitStatus.ID
	project.SubmittedBYID = &currentUserID
	project.SubmittedAt = &now
	project.NextReviewLevel = &level
	project.ApprovedBYID = nil
	project.ApprovedAt = nil
	project.RejectedBYID = nil
	project.RejectedAt = nil
	project.UpdatedBYID = &currentUserID

	if err := prs.repo.UpdateProject(project, tx); err != nil {
		return err
	}

	workflowLevel, err := prs.workflowSvc.GetWorkflowByLevel(projectID, level, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := prs.workflowStatusSvc.GetByType("resubmit")
	if err != nil {
		return err
	}

	if err := prs.createReviewWithHistory(
		workflowLevel.ID,
		projectID,
		level,
		workflowStatus.ID,
		nil,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	notify = func() {
		_ = prs.logService.LogActivity(
			models.UpdateAction,
			models.Project{}.TableName(),
			projectID,
			project,
			currentUserID,
			"project", "project_register", "", nil,
		)

		projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
		if err == nil {
			prs.sendResubmitNotification(projectID, workflowLevel, projectDetail, currentUserID)
		}
	}

	return err
}

func (prs *ProjectReviewService) Approve(projectID uint, currentUserID uint, form *projectDTOs.ProjectWorkflowReviewFormDTO) (err error) {
	tx := prs.repo.BeginTransaction()

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

	project, err := prs.projectSvc.GetProjectByID(projectID, tx)
	if err != nil {
		return err
	}

	rejectedStatus, err := prs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if project.NextReviewLevel == nil || (project.WorkflowStatusID != nil && *project.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"project",
		)
	}

	currentApprovalLevel := *project.NextReviewLevel
	currentWorkflow, err := prs.workflowSvc.GetWorkflowByLevel(projectID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := prs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	pendingStatus, err := prs.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.ProjectWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := prs.createReviewWithHistory(currentWorkflow.ID, projectID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			return err
		}

		projectWorkflow, err := prs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			return err
		}

		projectWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := prs.workflowSvc.UpdateProjectWorkflow(projectWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := prs.workflowSvc.GetNextLevel(projectID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			project.WorkflowStatusID = &workflowStatus.ID
			project.ApprovedBYID = &currentUserID
			project.ApprovedAt = &approvedAt
			project.RejectedBYID = nil
			project.RejectedAt = nil
			project.NextReviewLevel = nil
			project.UpdatedBYID = &currentUserID

			if err := prs.repo.UpdateProject(project, tx); err != nil {
				return err
			}

			notify = func() {
				_ = prs.logService.LogActivity(
					models.ApproveAction,
					models.Project{}.TableName(),
					projectID,
					project,
					currentUserID,
					"project", "project_register", "", nil,
				)

				projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
				if err == nil {
					prs.sendApprovalNotification(projectID, currentUserID, projectDetail)
				}
			}
		} else {
			project.WorkflowStatusID = &pendingStatus.ID
			project.NextReviewLevel = &nextWorkflowLevel.Level
			project.UpdatedBYID = &currentUserID

			if err := prs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
				return err
			}

			if err := prs.repo.UpdateProject(project, tx); err != nil {
				return err
			}

			nextWf, err := prs.workflowSvc.GetWorkflowByLevel(projectID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}

			notify = func() {
				_ = prs.logService.LogActivity(
					models.ApproveAction,
					models.Project{}.TableName(),
					projectID,
					project,
					currentUserID,
					"project", "project_register", "", nil,
				)

				projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
				if err == nil {
					prs.notifyNextLevelUsers(projectID, nextWf, projectDetail, currentUserID)
				}
			}
		}

	case "multiple-users", "role":
		var approvingParticipant *models.ProjectWorkflowParticipant
		for _, participant := range currentWorkflow.ProjectWorkflowParticipants {
			if participant.UserID == currentUserID {
				approvingParticipant = participant
				break
			}
		}
		if approvingParticipant == nil {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := prs.createReviewWithHistory(currentWorkflow.ID, projectID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			var allApproved bool
			allApproved, err = prs.workflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ProjectWorkflowParticipants, tx)
			if err != nil || !allApproved {
				if err == nil {
					notify = func() {
						_ = prs.logService.LogActivity(
							models.ApproveAction,
							models.Project{}.TableName(),
							projectID,
							project,
							currentUserID,
							"project", "project_register", "", nil,
						)
					}
				}
				return err
			}
		}

		projectWorkflow, err := prs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			return err
		}

		projectWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := prs.workflowSvc.UpdateProjectWorkflow(projectWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := prs.workflowSvc.GetNextLevel(projectID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			project.WorkflowStatusID = &workflowStatus.ID
			project.ApprovedBYID = &currentUserID
			project.ApprovedAt = &approvedAt
			project.RejectedBYID = nil
			project.RejectedAt = nil
			project.NextReviewLevel = nil
			project.UpdatedBYID = &currentUserID

			if err := prs.repo.UpdateProject(project, tx); err != nil {
				return err
			}

			notify = func() {
				_ = prs.logService.LogActivity(
					models.ApproveAction,
					models.Project{}.TableName(),
					projectID,
					project,
					currentUserID,
					"project", "project_register", "", nil,
				)

				projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
				if err == nil {
					prs.sendApprovalNotification(projectID, currentUserID, projectDetail)
				}
			}
		} else {
			project.WorkflowStatusID = &pendingStatus.ID
			project.NextReviewLevel = &nextWorkflowLevel.Level
			project.UpdatedBYID = &currentUserID

			if err := prs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
				return err
			}

			if err := prs.repo.UpdateProject(project, tx); err != nil {
				return err
			}

			nextWf, err := prs.workflowSvc.GetWorkflowByLevel(projectID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}

			notify = func() {
				_ = prs.logService.LogActivity(
					models.ApproveAction,
					models.Project{}.TableName(),
					projectID,
					project,
					currentUserID,
					"project", "project_register", "", nil,
				)

				projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
				if err == nil {
					prs.notifyNextLevelUsers(projectID, nextWf, projectDetail, currentUserID)
				}
			}
		}
	}

	return err
}

func (prs *ProjectReviewService) Revert(projectID uint, currentUserID uint, form *projectDTOs.ProjectWorkflowReviewFormDTO) (err error) {
	tx := prs.repo.BeginTransaction()

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

	project, err := prs.projectSvc.GetProjectByID(projectID, tx)
	if err != nil {
		return err
	}

	rejectedStatus, err := prs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if project.NextReviewLevel == nil || (project.WorkflowStatusID != nil && *project.WorkflowStatusID == rejectedStatus.ID) {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "project")
	}

	currentLevel := *project.NextReviewLevel

	currentWorkflow, err := prs.workflowSvc.GetWorkflowByLevel(projectID, currentLevel, tx)
	if err != nil {
		return err
	}

	allowed := false
	for _, participant := range currentWorkflow.ProjectWorkflowParticipants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}
	if !allowed {
		return utility.NewForbiddenError("you have no permission to revert")
	}

	revertStatus, err := prs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	if err := prs.createReviewWithHistory(
		currentWorkflow.ID,
		projectID,
		currentLevel,
		revertStatus.ID,
		comment,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.ProjectWorkflow
	if previousLevel > 0 {
		pendingStatus, err := prs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		previousLevelWorkflow, err = prs.workflowSvc.GetWorkflowByLevel(projectID, previousLevel, tx)
		if err != nil {
			return err
		}

		if err := prs.workflowSvc.DeleteReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := prs.workflowSvc.UpdateWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			return err
		}
	}

	var submittedByID *uint
	if previousLevel == 0 {
		submittedByID = project.SubmittedBYID

		project.WorkflowStatusID = &revertStatus.ID
		project.NextReviewLevel = nil
		project.SubmittedBYID = nil
		project.SubmittedAt = nil
		project.ApprovedBYID = nil
		project.ApprovedAt = nil
		project.RejectedBYID = nil
		project.RejectedAt = nil
	} else {
		pendingStatus, err := prs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}
		project.WorkflowStatusID = &pendingStatus.ID
		project.NextReviewLevel = &previousLevel
	}
	project.UpdatedBYID = &currentUserID

	if err := prs.repo.UpdateProject(project, tx); err != nil {
		return err
	}

	if previousLevel > 0 {
		notify = func() {
			_ = prs.logService.LogActivity(
				models.UpdateAction,
				models.Project{}.TableName(),
				projectID,
				project,
				currentUserID,
				"project", "project_register", "", nil,
			)

			projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
			if err == nil {
				prs.notifyNextLevelUsers(projectID, previousLevelWorkflow, projectDetail, currentUserID)
			}
		}
	} else {
		notify = func() {
			_ = prs.logService.LogActivity(
				models.UpdateAction,
				models.Project{}.TableName(),
				projectID,
				project,
				currentUserID,
				"project", "project_register", "", nil,
			)

			if submittedByID != nil && *submittedByID != 0 {
				projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
				if err == nil {
					prs.sendRevertNotification(projectID, currentUserID, projectDetail)
				}
			}
		}
	}

	return err
}

func (prs *ProjectReviewService) Reject(projectID uint, currentUserID uint, form *projectDTOs.ProjectWorkflowReviewFormDTO) (err error) {
	tx := prs.repo.BeginTransaction()

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

	project, err := prs.projectSvc.GetProjectByID(projectID, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := prs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if project.NextReviewLevel == nil ||
		(project.WorkflowStatusID != nil && *project.WorkflowStatusID == workflowStatus.ID) {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"project",
		)
	}

	currentReviewLevel := *project.NextReviewLevel

	currentWorkflow, err := prs.workflowSvc.GetWorkflowByLevel(projectID, currentReviewLevel, tx)
	if err != nil {
		return err
	}

	var rejectingParticipant *models.ProjectWorkflowParticipant

	for _, participant := range currentWorkflow.ProjectWorkflowParticipants {
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

	if err := prs.createReviewWithHistory(
		currentWorkflow.ID,
		projectID,
		currentReviewLevel,
		workflowStatus.ID,
		comment,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	currentWorkflow.WorkflowStatusID = workflowStatus.ID

	if err := prs.workflowSvc.UpdateProjectWorkflow(currentWorkflow, tx); err != nil {
		return err
	}
	if err := prs.workflowSvc.UpdateWorkflowStatus(currentWorkflow.ID, workflowStatus.ID, tx); err != nil {
		return err
	}

	var submittedByID uint
	if project.SubmittedBYID != nil {
		submittedByID = *project.SubmittedBYID
	} else if project.ResponsibleUserID != 0 {
		submittedByID = project.ResponsibleUserID
	}

	now := time.Now()
	project.WorkflowStatusID = &workflowStatus.ID
	project.RejectedBYID = &currentUserID
	project.RejectedAt = &now
	project.UpdatedBYID = &currentUserID

	if err := prs.repo.UpdateProject(project, tx); err != nil {
		return err
	}

	notify = func() {
		_ = prs.logService.LogActivity(
			models.RejectAction,
			models.Project{}.TableName(),
			projectID,
			project,
			currentUserID,
			"project", "project_register", "", nil,
		)

		if submittedByID != 0 {
			projectDetail, err := prs.projectSvc.GetFullDetailsByID(projectID)
			if err == nil {
				prs.sendRejectNotification(projectID, currentUserID, projectDetail)
			}
		}
	}

	return err
}

func (prs *ProjectReviewService) sendApprovalNotification(projectID uint, currentUserID uint, projectDetail *projectDTOs.ProjectDetailDTO) {
	go func() {
		if projectDetail == nil || projectDetail.ResponsibleUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d}`, projectID)
		entity := "project-review"
		action := "approval"

		approvedByUser, err := prs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{"error": err.Error()})
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(projectDetail.ResponsibleUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Project Workflow Approval",
			HeadingAr:            "الموافقة على سير عمل المشروع",
			ReferenceNumber:      fmt.Sprintf("PRJ-%d", projectDetail.ID),
			TitleEn:              projectDetail.TitleEn,
			TitleAr:              projectDetail.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 projectDetail.TitleEn,
			"FirstNameEn":           projectDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            projectDetail.ResponsibleUser.LastNameEN,
			"ReferenceNumber":       fmt.Sprintf("PRJ-%d", projectDetail.ID),
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/project-register/%d/overview", appURL, projectID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{projectDetail.ResponsibleUser.Email},
			Subject:              "Project Workflow Approval Notification",
			TemplateRelativePath: "project/templates/email/project_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = prs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{"error": err.Error()})
		}
	}()
}

func (prs *ProjectReviewService) notifyNextLevelUsers(projectID uint, nextWorkflow *models.ProjectWorkflow, projectDetail *projectDTOs.ProjectDetailDTO, currentUserID uint) {
	go func() {
		if nextWorkflow == nil || projectDetail == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d}`, projectID)
		entity := "project-review"
		action := "review-user-assign"

		assignedBy, err := prs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
			return
		}

		for _, participant := range nextWorkflow.ProjectWorkflowParticipants {
			reviewUser, err := prs.userSvc.GetByID(participant.UserID)
			if err != nil {
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Project Review Assigned",
				HeadingAr:            "تم تعيين مراجعة المشروع",
				ReferenceNumber:      fmt.Sprintf("PRJ-%d", projectDetail.ID),
				TitleEn:              projectDetail.TitleEn,
				TitleAr:              projectDetail.TitleAr,
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
				"Title":                 projectDetail.TitleEn,
				"ReferenceNumber":       fmt.Sprintf("PRJ-%d", projectDetail.ID),
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/projects/project-register/%d/overview", appURL, projectID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{reviewUser.Email},
				Subject:              "Project Review Assigned",
				TemplateRelativePath: "project/templates/email/project_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			_ = prs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
		}
	}()
}

func (prs *ProjectReviewService) sendResubmitNotification(projectID uint, workflowLevel *models.ProjectWorkflow, projectDetail *projectDTOs.ProjectDetailDTO, currentUserID uint) {
	go func() {
		if workflowLevel == nil || projectDetail == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d}`, projectID)
		entity := "project-review"
		action := "resubmitted"

		resubmittedBy, err := prs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch resubmitting user", map[string]interface{}{"error": err.Error()})
			return
		}

		for _, participant := range workflowLevel.ProjectWorkflowParticipants {
			reviewUser, err := prs.userSvc.GetByID(participant.UserID)
			if err != nil {
				continue
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "project",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Project Review Resubmitted",
				HeadingAr:            "تمت إعادة إرسال مراجعة المشروع",
				ReferenceNumber:      fmt.Sprintf("PRJ-%d", projectDetail.ID),
				TitleEn:              projectDetail.TitleEn,
				TitleAr:              projectDetail.TitleAr,
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
				"Title":                    projectDetail.TitleEn,
				"ReferenceNumber":          fmt.Sprintf("PRJ-%d", projectDetail.ID),
				"ResubmittedByFirstNameEn": resubmittedBy.FirstNameEn,
				"ResubmittedByLastNameEn":  resubmittedBy.LastNameEn,
				"ResubmittedAt":            time.Now().Format("2006-01-02"),
				"RedirectURL":              fmt.Sprintf("%s/projects/project-register/%d/overview", appURL, projectID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{reviewUser.Email},
				Subject:              "Project Review Resubmitted",
				TemplateRelativePath: "project/templates/email/project_workflow_resubmit_user_mail.html",
				EmailData:            emailData,
			}

			_ = prs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
		}
	}()
}

func (prs *ProjectReviewService) sendRevertNotification(projectID uint, currentUserID uint, projectDetail *projectDTOs.ProjectDetailDTO) {
	go func() {
		if projectDetail == nil || projectDetail.ResponsibleUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d}`, projectID)
		entity := "project-review"
		action := "reverted"

		revertedByUser, err := prs.userSvc.GetByID(currentUserID)
		if err != nil {
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(projectDetail.ResponsibleUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Project Workflow Reverted",
			HeadingAr:            "تم إرجاع سير عمل المشروع",
			ReferenceNumber:      fmt.Sprintf("PRJ-%d", projectDetail.ID),
			TitleEn:              projectDetail.TitleEn,
			TitleAr:              projectDetail.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 projectDetail.TitleEn,
			"FirstNameEn":           projectDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            projectDetail.ResponsibleUser.LastNameEN,
			"ReferenceNumber":       fmt.Sprintf("PRJ-%d", projectDetail.ID),
			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
			"RevertedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/project-register/%d/overview", appURL, projectID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{projectDetail.ResponsibleUser.Email},
			Subject:              "Project Workflow Reverted Notification",
			TemplateRelativePath: "project/templates/email/project_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		_ = prs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
	}()
}

func (prs *ProjectReviewService) sendRejectNotification(projectID uint, currentUserID uint, projectDetail *projectDTOs.ProjectDetailDTO) {
	go func() {
		if projectDetail == nil || projectDetail.ResponsibleUser.ID == 0 {
			return
		}

		urlParam := fmt.Sprintf(`{"projectId":%d}`, projectID)
		entity := "project-review"
		action := "rejected"

		rejectedByUser, err := prs.userSvc.GetByID(currentUserID)
		if err != nil {
			return
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(projectDetail.ResponsibleUser.ID)},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Project Workflow Rejected",
			HeadingAr:            "تم رفض سير عمل المشروع",
			ReferenceNumber:      fmt.Sprintf("PRJ-%d", projectDetail.ID),
			TitleEn:              projectDetail.TitleEn,
			TitleAr:              projectDetail.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 projectDetail.TitleEn,
			"FirstNameEn":           projectDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            projectDetail.ResponsibleUser.LastNameEN,
			"ReferenceNumber":       fmt.Sprintf("PRJ-%d", projectDetail.ID),
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/projects/project-register/%d/overview", appURL, projectID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{projectDetail.ResponsibleUser.Email},
			Subject:              "Project Workflow Rejected Notification",
			TemplateRelativePath: "project/templates/email/project_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		_ = prs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_management")
	}()
}
