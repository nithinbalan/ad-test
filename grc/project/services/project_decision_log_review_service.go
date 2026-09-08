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

func (pdlrs *ProjectDecisionLogReviewService) Submit(projectID uint, decisionLogID uint, currentUserID uint) error {
	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	if pdl.WorkflowStatusID != nil {
		workflowStatus, err := pdlrs.workflowStatusSvc.GetByID(*pdl.WorkflowStatusID)
		if err == nil && workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("project decision log has already been submitted")
		}
	}

	if err := pdlrs.workflowSvc.ValidateWorkflowBeforeCreate(); err != nil {
		return err
	}

	tx := pdlrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.workflowSvc.DeleteAllWorkflowsForProjectDecisionLog(decisionLogID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.workflowSvc.CreateProjectDecisionLogWorkflow(decisionLogID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	pdl.WorkflowStatusID = &workflowStatus.ID
	pdl.SubmittedBYID = &currentUserID
	pdl.SubmittedAt = &currentDate
	pdl.NextReviewLevel = &level
	pdl.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
		tx.Rollback()
		return err
	}

	workflowLevel, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.createReviewWithHistory(workflowLevel.ID, decisionLogID, level, workflowStatus.ID, nil, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pdlrs.logService.LogActivity(
			models.CreateAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
	if err == nil {
		pdlrs.notifyNextLevelUsers(projectID, decisionLogID, workflowLevel, pdlDetail, currentUserID)
	}

	return nil
}

func (pdlrs *ProjectDecisionLogReviewService) Resubmit(projectID uint, decisionLogID uint, currentUserID uint) error {
	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	if pdl.WorkflowStatusID == nil {
		return utility.NewForbiddenError("project decision log must be reverted before resubmitting")
	}

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByID(*pdl.WorkflowStatusID)
	if err != nil || workflowStatus.Type != "revert" {
		return utility.NewForbiddenError("project decision log is not in reverted state")
	}

	if err := pdlrs.workflowSvc.ValidateWorkflowBeforeCreate(); err != nil {
		return err
	}

	tx := pdlrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	submitStatus, err := pdlrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.workflowSvc.DeleteAllWorkflowsForProjectDecisionLog(decisionLogID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.workflowSvc.CreateProjectDecisionLogWorkflow(decisionLogID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	pdl.WorkflowStatusID = &submitStatus.ID
	pdl.SubmittedBYID = &currentUserID
	pdl.SubmittedAt = &currentDate
	pdl.NextReviewLevel = &level
	pdl.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(pdl, tx); err != nil {
		tx.Rollback()
		return err
	}

	workflowLevel, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.createReviewWithHistory(workflowLevel.ID, decisionLogID, level, submitStatus.ID, nil, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pdlrs.logService.LogActivity(
			models.UpdateAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
	if err == nil {
		pdlrs.sendResubmitNotification(projectID, decisionLogID, workflowLevel, pdlDetail, currentUserID)
	}

	return nil
}

func (pdlrs *ProjectDecisionLogReviewService) Approve(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) error {
	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel == nil {
		return utility.NewForbiddenError("no pending review level found")
	}

	tx := pdlrs.repo.BeginTransaction()
	currentApprovalLevel := pdl.NextReviewLevel
	currentWorkflow, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := pdlrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	var nextWorkflowLevel *models.ProjectDecisionLogWorkflow

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "custom-roles", "hod":
		if len(currentWorkflow.ProjectDecisionLogWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectDecisionLogWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pdlrs.createReviewWithHistory(currentWorkflow.ID, decisionLogID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		projectDecisionLogWorkflow, err := pdlrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		projectDecisionLogWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pdlrs.workflowSvc.UpdateProjectDecisionLogWorkflow(projectDecisionLogWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err = pdlrs.workflowSvc.GetNextLevel(decisionLogID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingPDL, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingPDL.WorkflowStatusID = &workflowStatus.ID
			existingPDL.ApprovedBYID = &currentUserID
			existingPDL.ApprovedAt = &approvedAt
			existingPDL.NextReviewLevel = nil
		} else {
			existingPDL.NextReviewLevel = &nextWorkflowLevel.Level
		}
		existingPDL.UpdatedByID = &currentUserID

		if err := pdlrs.repo.UpdateDecisionLog(existingPDL, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ProjectDecisionLogWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pdlrs.createReviewWithHistory(currentWorkflow.ID, decisionLogID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := pdlrs.workflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ProjectDecisionLogWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		projectDecisionLogWorkflow, err := pdlrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		projectDecisionLogWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pdlrs.workflowSvc.UpdateProjectDecisionLogWorkflow(projectDecisionLogWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err = pdlrs.workflowSvc.GetNextLevel(decisionLogID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingPDL, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingPDL.WorkflowStatusID = &workflowStatus.ID
			existingPDL.ApprovedBYID = &currentUserID
			existingPDL.ApprovedAt = &approvedAt
			existingPDL.NextReviewLevel = nil
		} else {
			existingPDL.NextReviewLevel = &nextWorkflowLevel.Level
		}
		existingPDL.UpdatedByID = &currentUserID

		if err := pdlrs.repo.UpdateDecisionLog(existingPDL, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go func() {
		_ = pdlrs.logService.LogActivity(
			models.ApproveAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
	if err == nil {
		if nextWorkflowLevel != nil {
			nextWf, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, nextWorkflowLevel.Level)
			if err == nil {
				pdlrs.notifyNextLevelUsers(projectID, decisionLogID, nextWf, pdlDetail, currentUserID)
			}
		} else {
			pdlrs.sendApprovalNotification(projectID, decisionLogID, currentUserID, pdlDetail)
		}
	}

	return nil
}

func (pdlrs *ProjectDecisionLogReviewService) Revert(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) error {
	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel == nil {
		return utility.NewForbiddenError("no pending review level found")
	}

	tx := pdlrs.repo.BeginTransaction()
	currentApprovalLevel := pdl.NextReviewLevel
	currentWorkflow, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	revertStatus, err := pdlrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		tx.Rollback()
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "custom-roles", "hod":
		if len(currentWorkflow.ProjectDecisionLogWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectDecisionLogWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to revert")
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ProjectDecisionLogWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to revert")
		}
	}

	if err := pdlrs.workflowSvc.DeleteReviewsByWorkflowID(currentWorkflow.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.createReviewWithHistory(currentWorkflow.ID, decisionLogID, currentWorkflow.Level, revertStatus.ID, comment, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := pdlrs.workflowSvc.UpdateWorkflowStatus(currentWorkflow.ID, revertStatus.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	currentLevel := *currentApprovalLevel
	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.ProjectDecisionLogWorkflow
	if previousLevel > 0 {
		pendingStatus, err := pdlrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			tx.Rollback()
			return err
		}

		previousLevelWorkflow, err = pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, previousLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := pdlrs.workflowSvc.DeleteReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := pdlrs.workflowSvc.UpdateWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingPDL, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if previousLevel == 0 {
		existingPDL.WorkflowStatusID = &revertStatus.ID
		existingPDL.NextReviewLevel = nil
	} else {
		existingPDL.NextReviewLevel = &previousLevel
	}
	existingPDL.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(existingPDL, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pdlrs.logService.LogActivity(
			models.RejectAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
	if err == nil {
		if previousLevel > 0 {
			pdlrs.notifyNextLevelUsers(projectID, decisionLogID, previousLevelWorkflow, pdlDetail, currentUserID)
		} else {
			pdlrs.sendRevertNotification(projectID, decisionLogID, currentUserID, pdlDetail)
		}
	}

	return nil
}

func (pdlrs *ProjectDecisionLogReviewService) Reject(projectID uint, decisionLogID uint, currentUserID uint, form *projectDTOs.ProjectDecisionLogReviewFormDTO) error {
	pdl, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		return err
	}

	if pdl.NextReviewLevel == nil {
		return utility.NewForbiddenError("no pending review level found")
	}

	tx := pdlrs.repo.BeginTransaction()
	currentApprovalLevel := pdl.NextReviewLevel
	currentWorkflow, err := pdlrs.workflowSvc.GetWorkflowByLevel(decisionLogID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	rejectStatus, err := pdlrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		tx.Rollback()
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "custom-roles", "hod":
		if len(currentWorkflow.ProjectDecisionLogWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectDecisionLogWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to reject")
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ProjectDecisionLogWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to reject")
		}
	}

	if err := pdlrs.createReviewWithHistory(currentWorkflow.ID, decisionLogID, currentWorkflow.Level, rejectStatus.ID, comment, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	projectDecisionLogWorkflow, err := pdlrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	projectDecisionLogWorkflow.WorkflowStatusID = rejectStatus.ID
	if err := pdlrs.workflowSvc.UpdateProjectDecisionLogWorkflow(projectDecisionLogWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingPDL, err := pdlrs.projectDLSvc.GetByID(projectID, decisionLogID)
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	existingPDL.WorkflowStatusID = &rejectStatus.ID
	existingPDL.RejectedBYID = &currentUserID
	existingPDL.RejectedAt = &now
	existingPDL.NextReviewLevel = nil
	existingPDL.UpdatedByID = &currentUserID

	if err := pdlrs.repo.UpdateDecisionLog(existingPDL, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pdlrs.logService.LogActivity(
			models.RejectAction,
			models.ProjectDecisionLog{}.TableName(),
			decisionLogID,
			pdl,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pdlDetail, err := pdlrs.projectDLSvc.GetFullDetailsByID(projectID, decisionLogID)
	if err == nil {
		pdlrs.sendRejectNotification(projectID, decisionLogID, currentUserID, pdlDetail)
	}

	return nil
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
