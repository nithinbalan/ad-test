package services

import (
	"fmt"
	"path/filepath"
	"strconv"
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

type ProjectSecurityGateChecklistExceptionReviewServiceInterface interface {
	Submit(securityGateID uint, checklistID uint, currentUserID uint) error
	Approve(securityGateID uint, checklistID uint, currentUserID uint, form *projectDTOs.ProjectSecurityGateChecklistExceptionReviewFormDTO) error
	Revert(securityGateID uint, checklistID uint, currentUserID uint, form *projectDTOs.ProjectSecurityGateChecklistExceptionReviewFormDTO) error
}

type ProjectSecurityGateChecklistExceptionReviewService struct {
	repo              repositories.ProjectSecurityGateChecklistExceptionReviewRepositoryInterface
	logService        logServices.ActivityLogServiceInterface
	checklistSvc      ProjectSecurityGateChecklistServiceInterface
	securityGateSvc   ProjectSecurityGateServiceInterface
	workflowSvc       ProjectSecurityGateChecklistExceptionWorkflowServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewProjectSecurityGateChecklistExceptionReviewService(
	repo repositories.ProjectSecurityGateChecklistExceptionReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	checklistSvc ProjectSecurityGateChecklistServiceInterface,
	securityGateSvc ProjectSecurityGateServiceInterface,
	workflowSvc ProjectSecurityGateChecklistExceptionWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProjectSecurityGateChecklistExceptionReviewServiceInterface {
	return &ProjectSecurityGateChecklistExceptionReviewService{
		repo:              repo,
		logService:        logService,
		checklistSvc:      checklistSvc,
		securityGateSvc:   securityGateSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
	}
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) Submit(securityGateID uint, checklistID uint, currentUserID uint) (err error) {
	tx := psgcers.repo.BeginTransaction()

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

	checklist, err := psgcers.checklistSvc.GetByID(securityGateID, checklistID, tx)
	if err != nil {
		return err
	}

	if !checklist.IsException {
		return utility.NewValidationError(
			fmt.Errorf("checklist item is not marked as an exception"),
			"checklist_not_marked_as_exception",
			"exception",
		)
	}

	if checklist.NextReviewLevel != nil {
		return utility.NewValidationError(
			fmt.Errorf("exception has already been submitted"),
			"workflow_already_submitted",
			"checklist exception",
		)
	}

	if checklist.ApprovedBYID != nil {
		return utility.NewValidationError(
			fmt.Errorf("exception workflow has already been completed"),
			"workflow_already_completed",
			"checklist exception",
		)
	}

	if err := psgcers.workflowSvc.ValidateWorkflowBeforeCreate(tx); err != nil {
		return err
	}

	if err := psgcers.workflowSvc.DeleteAllWorkflowsForChecklist(checklistID, tx); err != nil {
		return err
	}

	if err := psgcers.workflowSvc.CreateWorkflow(checklistID, currentUserID, tx); err != nil {
		return err
	}

	submitStatus, err := psgcers.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	oldData := map[string]interface{}{
		"workflow_status_id": nil,
		"next_review_level":  nil,
		"submitted_by_id":    nil,
	}
	if checklist.WorkflowStatusID != nil {
		oldData["workflow_status_id"] = *checklist.WorkflowStatusID
	}
	if checklist.NextReviewLevel != nil {
		oldData["next_review_level"] = *checklist.NextReviewLevel
	}
	if checklist.SubmittedBYID != nil {
		oldData["submitted_by_id"] = *checklist.SubmittedBYID
	}

	level := 1
	now := time.Now()

	checklist.WorkflowStatusID = &submitStatus.ID
	checklist.SubmittedBYID = &currentUserID
	checklist.SubmittedAt = &now
	checklist.NextReviewLevel = &level
	checklist.UpdatedBYID = &currentUserID
	checklist.UpdatedAt = &now

	newData := map[string]interface{}{
		"workflow_status_id": submitStatus.ID,
		"next_review_level":  level,
		"submitted_by_id":    currentUserID,
	}
	changes := utility.CompareFieldChanges(oldData, newData)

	if err := psgcers.repo.UpdateChecklist(checklist, tx); err != nil {
		return err
	}

	workflowLevel, err := psgcers.workflowSvc.GetWorkflowByLevel(checklistID, level, tx)
	if err != nil {
		return err
	}

	if err := psgcers.createReviewWithHistory(securityGateID, workflowLevel.ID, checklistID, level, submitStatus.ID, nil, nil, currentUserID, tx); err != nil {
		return err
	}

	notify = func() {
		if len(changes) > 0 {
			_ = psgcers.logService.LogActivity(
				models.UpdateAction,
				models.ProjectSecurityGateChecklist{}.TableName(),
				checklistID,
				changes,
				currentUserID,
				"project", "project_security_gate", "project_security_gate_checklist_exception_submit", nil,
			)
		}

		psgcers.notifyWorkflowParticipants(securityGateID, checklistID, workflowLevel, currentUserID)
	}

	return err
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) Approve(securityGateID uint, checklistID uint, currentUserID uint, form *projectDTOs.ProjectSecurityGateChecklistExceptionReviewFormDTO) (err error) {
	tx := psgcers.repo.BeginTransaction()

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

	checklist, err := psgcers.checklistSvc.GetByID(securityGateID, checklistID, tx)
	if err != nil {
		return err
	}

	if checklist.NextReviewLevel == nil {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "checklist exception")
	}
	currentApprovalLevel := *checklist.NextReviewLevel

	currentWorkflow, err := psgcers.workflowSvc.GetWorkflowByLevel(checklistID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	approveStatus, err := psgcers.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	pendingStatus, err := psgcers.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	var formDocuments []*commonDTOs.DocumentFormDTO
	if form != nil {
		formDocuments = form.Documents
	}

	alreadyApproved, err := psgcers.workflowSvc.HasUserApproved(currentWorkflow.ID, currentUserID, tx)
	if err != nil {
		return err
	}
	if alreadyApproved {
		return utility.NewValidationError(
			fmt.Errorf("user has already approved this level"),
			"workflow_already_approved",
			"checklist exception",
		)
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.Participants) == 0 || currentUserID != currentWorkflow.Participants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := psgcers.createReviewWithHistory(securityGateID, currentWorkflow.ID, checklistID, currentApprovalLevel, approveStatus.ID, comment, formDocuments, currentUserID, tx); err != nil {
			return err
		}

	case "multiple-users", "role":
		var approvingParticipant *models.ProjectSecurityGateChecklistExceptionWorkflowParticipant
		for _, participant := range currentWorkflow.Participants {
			if participant.UserID == currentUserID {
				approvingParticipant = participant
				break
			}
		}
		if approvingParticipant == nil {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := psgcers.createReviewWithHistory(securityGateID, currentWorkflow.ID, checklistID, currentApprovalLevel, approveStatus.ID, comment, formDocuments, currentUserID, tx); err != nil {
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, allApprovedErr := psgcers.workflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.Participants, tx)
			if allApprovedErr != nil {
				return allApprovedErr
			}

			// The remaining approvers still have to act, so the level stays open.
			if !allApproved {
				return err
			}
		}

	default:
		return utility.NewForbiddenError("you have no permission to approve")
	}

	currentWorkflow.WorkflowStatusID = approveStatus.ID
	currentWorkflow.UpdatedBYID = &currentUserID
	if err := psgcers.workflowSvc.UpdateWorkflow(currentWorkflow, tx); err != nil {
		return err
	}

	nextWorkflowLevel, err := psgcers.workflowSvc.GetNextLevel(checklistID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	oldData := map[string]interface{}{
		"workflow_status_id": nil,
		"next_review_level":  nil,
		"approved_by_id":     nil,
	}
	if checklist.WorkflowStatusID != nil {
		oldData["workflow_status_id"] = *checklist.WorkflowStatusID
	}
	if checklist.NextReviewLevel != nil {
		oldData["next_review_level"] = *checklist.NextReviewLevel
	}
	if checklist.ApprovedBYID != nil {
		oldData["approved_by_id"] = *checklist.ApprovedBYID
	}

	now := time.Now()
	if nextWorkflowLevel == nil {
		checklist.WorkflowStatusID = &approveStatus.ID
		checklist.ApprovedBYID = &currentUserID
		checklist.ApprovedAt = &now
		checklist.NextReviewLevel = nil
	} else {
		checklist.WorkflowStatusID = &pendingStatus.ID
		checklist.NextReviewLevel = &nextWorkflowLevel.Level

		if err := psgcers.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
			return err
		}
	}
	checklist.UpdatedBYID = &currentUserID
	checklist.UpdatedAt = &now

	newData := map[string]interface{}{
		"workflow_status_id": nil,
		"next_review_level":  nil,
		"approved_by_id":     nil,
	}
	if checklist.WorkflowStatusID != nil {
		newData["workflow_status_id"] = *checklist.WorkflowStatusID
	}
	if checklist.NextReviewLevel != nil {
		newData["next_review_level"] = *checklist.NextReviewLevel
	}
	if checklist.ApprovedBYID != nil {
		newData["approved_by_id"] = *checklist.ApprovedBYID
	}
	changes := utility.CompareFieldChanges(oldData, newData)

	if err := psgcers.repo.UpdateChecklist(checklist, tx); err != nil {
		return err
	}

	if nextWorkflowLevel != nil {
		nextWorkflow, err := psgcers.workflowSvc.GetWorkflowByLevel(checklistID, nextWorkflowLevel.Level, tx)
		if err != nil {
			return err
		}

		notify = func() {
			if len(changes) > 0 {
				_ = psgcers.logService.LogActivity(
					models.ApproveAction,
					models.ProjectSecurityGateChecklist{}.TableName(),
					checklistID,
					changes,
					currentUserID,
					"project", "project_security_gate", "project_security_gate_checklist_exception_approve", nil,
				)
			}

			psgcers.notifyWorkflowParticipants(securityGateID, checklistID, nextWorkflow, currentUserID)
		}
	} else {
		submittedByID := checklist.SubmittedBYID

		notify = func() {
			if len(changes) > 0 {
				_ = psgcers.logService.LogActivity(
					models.ApproveAction,
					models.ProjectSecurityGateChecklist{}.TableName(),
					checklistID,
					changes,
					currentUserID,
					"project", "project_security_gate", "project_security_gate_checklist_exception_approve", nil,
				)
			}

			psgcers.sendApprovalNotification(securityGateID, checklistID, currentUserID, submittedByID, checklist.CreatedBYID)
		}
	}

	return err
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) Revert(securityGateID uint, checklistID uint, currentUserID uint, form *projectDTOs.ProjectSecurityGateChecklistExceptionReviewFormDTO) (err error) {
	tx := psgcers.repo.BeginTransaction()

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

	checklist, err := psgcers.checklistSvc.GetByID(securityGateID, checklistID, tx)
	if err != nil {
		return err
	}

	if checklist.NextReviewLevel == nil {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "checklist exception")
	}
	currentLevel := *checklist.NextReviewLevel

	currentWorkflow, err := psgcers.workflowSvc.GetWorkflowByLevel(checklistID, currentLevel, tx)
	if err != nil {
		return err
	}

	allowed := false
	for _, participant := range currentWorkflow.Participants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}
	if !allowed {
		return utility.NewForbiddenError("you have no permission to revert")
	}

	revertStatus, err := psgcers.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	var formDocuments []*commonDTOs.DocumentFormDTO
	if form != nil {
		formDocuments = form.Documents
	}

	if err := psgcers.workflowSvc.DeleteReviewsByWorkflowID(currentWorkflow.ID, tx); err != nil {
		return err
	}

	if err := psgcers.createReviewWithHistory(securityGateID, currentWorkflow.ID, checklistID, currentLevel, revertStatus.ID, comment, formDocuments, currentUserID, tx); err != nil {
		return err
	}

	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.ProjectSecurityGateChecklistExceptionWorkflow
	if previousLevel > 0 {
		pendingStatus, err := psgcers.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		previousLevelWorkflow, err = psgcers.workflowSvc.GetWorkflowByLevel(checklistID, previousLevel, tx)
		if err != nil {
			return err
		}

		if err := psgcers.workflowSvc.DeleteReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := psgcers.workflowSvc.UpdateWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			return err
		}
	}

	oldData := map[string]interface{}{
		"workflow_status_id": nil,
		"next_review_level":  nil,
		"submitted_by_id":    nil,
		"approved_by_id":     nil,
	}
	if checklist.WorkflowStatusID != nil {
		oldData["workflow_status_id"] = *checklist.WorkflowStatusID
	}
	if checklist.NextReviewLevel != nil {
		oldData["next_review_level"] = *checklist.NextReviewLevel
	}
	if checklist.SubmittedBYID != nil {
		oldData["submitted_by_id"] = *checklist.SubmittedBYID
	}
	if checklist.ApprovedBYID != nil {
		oldData["approved_by_id"] = *checklist.ApprovedBYID
	}

	now := time.Now()
	var submittedByID *uint
	if previousLevel == 0 {
		submittedByID = checklist.SubmittedBYID

		checklist.WorkflowStatusID = &revertStatus.ID
		checklist.NextReviewLevel = nil
		checklist.SubmittedBYID = nil
		checklist.SubmittedAt = nil
		checklist.ApprovedBYID = nil
		checklist.ApprovedAt = nil
	} else {
		pendingStatus, err := psgcers.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		checklist.WorkflowStatusID = &pendingStatus.ID
		checklist.NextReviewLevel = &previousLevel
	}
	checklist.UpdatedBYID = &currentUserID
	checklist.UpdatedAt = &now

	newData := map[string]interface{}{
		"workflow_status_id": nil,
		"next_review_level":  nil,
		"submitted_by_id":    nil,
		"approved_by_id":     nil,
	}
	if checklist.WorkflowStatusID != nil {
		newData["workflow_status_id"] = *checklist.WorkflowStatusID
	}
	if checklist.NextReviewLevel != nil {
		newData["next_review_level"] = *checklist.NextReviewLevel
	}
	if checklist.SubmittedBYID != nil {
		newData["submitted_by_id"] = *checklist.SubmittedBYID
	}
	if checklist.ApprovedBYID != nil {
		newData["approved_by_id"] = *checklist.ApprovedBYID
	}
	changes := utility.CompareFieldChanges(oldData, newData)

	if err := psgcers.repo.UpdateChecklist(checklist, tx); err != nil {
		return err
	}

	createdByID := checklist.CreatedBYID

	if previousLevel > 0 {
		notify = func() {
			if len(changes) > 0 {
				_ = psgcers.logService.LogActivity(
					models.UpdateAction,
					models.ProjectSecurityGateChecklist{}.TableName(),
					checklistID,
					changes,
					currentUserID,
					"project", "project_security_gate", "project_security_gate_checklist_exception_revert", nil,
				)
			}

			psgcers.notifyRevertedToPreviousLevel(securityGateID, checklistID, previousLevelWorkflow, currentUserID)
		}
	} else {
		notify = func() {
			if len(changes) > 0 {
				_ = psgcers.logService.LogActivity(
					models.UpdateAction,
					models.ProjectSecurityGateChecklist{}.TableName(),
					checklistID,
					changes,
					currentUserID,
					"project", "project_security_gate", "project_security_gate_checklist_exception_revert", nil,
				)
			}

			psgcers.sendRevertNotification(securityGateID, checklistID, currentUserID, submittedByID, createdByID)
		}
	}

	return err
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) notifyWorkflowParticipants(securityGateID uint, checklistID uint, workflowLevel *models.ProjectSecurityGateChecklistExceptionWorkflow, currentUserID uint) {
	if workflowLevel == nil {
		return
	}

	assignedBy, err := psgcers.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
		return
	}

	securityGate, err := psgcers.securityGateSvc.GetByID(securityGateID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistDetail, err := psgcers.checklistSvc.GetFullDetailsByID(securityGateID, checklistID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate checklist details", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistTitle := ""
	if checklistDetail.Requirement != nil {
		checklistTitle = checklistDetail.Requirement.TitleEn
	}

	urlParam := fmt.Sprintf(`{"securityGateId":%d,"checklistId":%d}`, securityGateID, checklistID)
	entity := "project-security-gate-checklist-exception-review"
	action := "review-user-assign"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/projects/security-gate/%d/checklists/%d", appURL, securityGateID, checklistID)

	for _, participant := range workflowLevel.Participants {
		reviewUser, err := psgcers.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{"error": err.Error()})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.ReviewType),
			ReferenceNumber:      securityGate.ProjectTitleEn,
			HeadingEn:            "Project Security Gate Exception Review Assigned",
			HeadingAr:            "تم تعيين مراجعة استثناء بوابة أمان المشروع",
			TitleEn:              checklistTitle,
			TitleAr:              checklistTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           reviewUser.FirstNameEn,
			"LastNameEn":            reviewUser.LastNameEn,
			"ProjectTitle":          securityGate.ProjectTitleEn,
			"Checklist":             checklistTitle,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{reviewUser.Email},
			Subject:              fmt.Sprintf("Project Security Gate Exception: %s - Review Assigned", securityGate.ProjectTitleEn),
			TemplateRelativePath: "project/templates/email/project_security_gate_checklist_exception_workflow_review_user_mail.html",
			EmailData:            emailData,
		}

		if err := psgcers.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_security_gate"); err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{"error": err.Error()})
		}
	}
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) notifyRevertedToPreviousLevel(securityGateID uint, checklistID uint, workflowLevel *models.ProjectSecurityGateChecklistExceptionWorkflow, currentUserID uint) {
	if workflowLevel == nil {
		return
	}

	revertedByUser, err := psgcers.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{"error": err.Error()})
		return
	}

	securityGate, err := psgcers.securityGateSvc.GetByID(securityGateID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistDetail, err := psgcers.checklistSvc.GetFullDetailsByID(securityGateID, checklistID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate checklist details", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistTitle := ""
	if checklistDetail.Requirement != nil {
		checklistTitle = checklistDetail.Requirement.TitleEn
	}

	urlParam := fmt.Sprintf(`{"securityGateId":%d,"checklistId":%d}`, securityGateID, checklistID)
	entity := "project-security-gate-checklist-exception-review"
	action := "reverted"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/projects/security-gate/%d/checklists/%d", appURL, securityGateID, checklistID)

	for _, participant := range workflowLevel.Participants {
		reviewUser, err := psgcers.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{"error": err.Error()})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "project",
			NotificationType:     string(commonDTOs.RevertType),
			ReferenceNumber:      securityGate.ProjectTitleEn,
			HeadingEn:            "Project Security Gate Exception Reverted",
			HeadingAr:            "تم إرجاع استثناء بوابة أمان المشروع",
			TitleEn:              checklistTitle,
			TitleAr:              checklistTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           reviewUser.FirstNameEn,
			"LastNameEn":            reviewUser.LastNameEn,
			"ProjectTitle":          securityGate.ProjectTitleEn,
			"Checklist":             checklistTitle,
			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
			"RevertedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{reviewUser.Email},
			Subject:              fmt.Sprintf("Project Security Gate Exception: %s - Reverted", securityGate.ProjectTitleEn),
			TemplateRelativePath: "project/templates/email/project_security_gate_checklist_exception_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		if err := psgcers.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_security_gate"); err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{"error": err.Error()})
		}
	}
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) sendApprovalNotification(securityGateID uint, checklistID uint, currentUserID uint, submittedByID *uint, createdByID *uint) {
	notifyUserID := submittedByID
	if notifyUserID == nil {
		notifyUserID = createdByID
	}
	if notifyUserID == nil {
		return
	}

	approvedByUser, err := psgcers.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch approving user", map[string]interface{}{"error": err.Error()})
		return
	}

	notifyUser, err := psgcers.userSvc.GetByID(*notifyUserID)
	if err != nil {
		utility.LogError("Failed to fetch submitted user", map[string]interface{}{"error": err.Error()})
		return
	}

	securityGate, err := psgcers.securityGateSvc.GetByID(securityGateID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistDetail, err := psgcers.checklistSvc.GetFullDetailsByID(securityGateID, checklistID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate checklist details", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistTitle := ""
	if checklistDetail.Requirement != nil {
		checklistTitle = checklistDetail.Requirement.TitleEn
	}

	urlParam := fmt.Sprintf(`{"securityGateId":%d,"checklistId":%d}`, securityGateID, checklistID)
	entity := "project-security-gate-checklist-exception-review"
	action := "approval"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/projects/security-gate/%d/checklists/%d", appURL, securityGateID, checklistID)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*notifyUserID},
		ModuleIdentifier:     "project",
		NotificationType:     string(commonDTOs.ApprovalType),
		ReferenceNumber:      securityGate.ProjectTitleEn,
		HeadingEn:            "Project Security Gate Exception Approved",
		HeadingAr:            "تمت الموافقة على استثناء بوابة أمان المشروع",
		TitleEn:              checklistTitle,
		TitleAr:              checklistTitle,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           notifyUser.FirstNameEn,
		"LastNameEn":            notifyUser.LastNameEn,
		"ProjectTitle":          securityGate.ProjectTitleEn,
		"Checklist":             checklistTitle,
		"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
		"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
		"ApprovedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{notifyUser.Email},
		Subject:              fmt.Sprintf("Project Security Gate Exception: %s - Approved", securityGate.ProjectTitleEn),
		TemplateRelativePath: "project/templates/email/project_security_gate_checklist_exception_workflow_approve_user_mail.html",
		EmailData:            emailData,
	}

	if err := psgcers.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_security_gate"); err != nil {
		utility.LogError("Failed to send approval notification email", map[string]interface{}{"error": err.Error()})
	}
}

func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) sendRevertNotification(securityGateID uint, checklistID uint, currentUserID uint, submittedByID *uint, createdByID *uint) {
	notifyUserID := submittedByID
	if notifyUserID == nil {
		notifyUserID = createdByID
	}
	if notifyUserID == nil {
		return
	}

	revertedByUser, err := psgcers.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{"error": err.Error()})
		return
	}

	notifyUser, err := psgcers.userSvc.GetByID(*notifyUserID)
	if err != nil {
		utility.LogError("Failed to fetch submitted user", map[string]interface{}{"error": err.Error()})
		return
	}

	securityGate, err := psgcers.securityGateSvc.GetByID(securityGateID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistDetail, err := psgcers.checklistSvc.GetFullDetailsByID(securityGateID, checklistID)
	if err != nil {
		utility.LogError("Failed to fetch project security gate checklist details", map[string]interface{}{"error": err.Error()})
		return
	}

	checklistTitle := ""
	if checklistDetail.Requirement != nil {
		checklistTitle = checklistDetail.Requirement.TitleEn
	}

	urlParam := fmt.Sprintf(`{"securityGateId":%d,"checklistId":%d}`, securityGateID, checklistID)
	entity := "project-security-gate-checklist-exception-review"
	action := "reverted"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/projects/security-gate/%d/checklists/%d", appURL, securityGateID, checklistID)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*notifyUserID},
		ModuleIdentifier:     "project",
		NotificationType:     string(commonDTOs.RevertType),
		ReferenceNumber:      securityGate.ProjectTitleEn,
		HeadingEn:            "Project Security Gate Exception Reverted",
		HeadingAr:            "تم إرجاع استثناء بوابة أمان المشروع",
		TitleEn:              checklistTitle,
		TitleAr:              checklistTitle,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           notifyUser.FirstNameEn,
		"LastNameEn":            notifyUser.LastNameEn,
		"ProjectTitle":          securityGate.ProjectTitleEn,
		"Checklist":             checklistTitle,
		"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
		"RevertedByLastNameEn":  revertedByUser.LastNameEn,
		"RevertedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{notifyUser.Email},
		Subject:              fmt.Sprintf("Project Security Gate Exception: %s - Reverted", securityGate.ProjectTitleEn),
		TemplateRelativePath: "project/templates/email/project_security_gate_checklist_exception_workflow_revert_user_mail.html",
		EmailData:            emailData,
	}

	if err := psgcers.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_security_gate"); err != nil {
		utility.LogError("Failed to send reverted notification email", map[string]interface{}{"error": err.Error()})
	}
}

// Helper function to create review in both current and history tables
func (psgcers *ProjectSecurityGateChecklistExceptionReviewService) createReviewWithHistory(securityGateID uint, workflowID uint, checklistID uint, level int, workflowStatusID uint, comment *string, formDocuments []*commonDTOs.DocumentFormDTO, currentUserID uint, tx utility.DatabaseTransaction) error {
	review := &models.ProjectSecurityGateChecklistExceptionWorkflowReview{
		ExceptionWorkflowID: workflowID,
		WorkflowStatusID:    workflowStatusID,
		Comment:             comment,
		CreatedBYID:         currentUserID,
		CreatedAt:           time.Now(),
	}

	if err := psgcers.repo.CreateReview(review, tx); err != nil {
		return err
	}

	history := &models.ProjectSecurityGateChecklistExceptionWorkflowHistory{
		ChecklistID:      checklistID,
		WorkflowStatusID: workflowStatusID,
		Level:            level,
		Comment:          comment,
		CreatedBYID:      currentUserID,
		CreatedAt:        time.Now(),
	}

	if err := psgcers.repo.CreateHistory(history, tx); err != nil {
		return err
	}

	if len(formDocuments) == 0 {
		return nil
	}

	var documents []*models.ProjectSecurityGateChecklistExceptionWorkflowReviewDocument

	for _, doc := range formDocuments {
		encryptedFileName, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			return err
		}

		url := utility.GenerateDocumentURL(
			"project",
			"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
			"checklists", strconv.FormatUint(uint64(checklistID), 10),
			"reviews", strconv.FormatUint(uint64(history.ID), 10),
			"documents",
		)

		documents = append(documents, &models.ProjectSecurityGateChecklistExceptionWorkflowReviewDocument{
			HistoryID:    history.ID,
			URL:          url,
			Name:         doc.FileName,
			Size:         doc.Size,
			Token:        encryptedFileName,
			OriginalName: doc.OriginalName,
			Type:         doc.Type,
			CreatedBYID:  currentUserID,
			CreatedAt:    time.Now(),
		})
	}

	if err := psgcers.repo.CreateReviewDocuments(documents, tx); err != nil {
		return err
	}

	storagePath := filepath.Join(
		"project",
		"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
		"checklists", strconv.FormatUint(uint64(checklistID), 10),
		"reviews", strconv.FormatUint(uint64(history.ID), 10),
		"documents",
	)

	for _, document := range documents {
		if err := utility.StoreFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	return nil
}
