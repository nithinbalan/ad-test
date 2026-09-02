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
}

func NewProjectChangeRequestReviewService(
	repo repositories.ProjectChangeRequestReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	projectCRSvc ProjectChangeRequestServiceInterface,
	workflowSvc ProjectChangeRequestWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProjectChangeRequestReviewServiceInterface {
	return &ProjectChangeRequestReviewService{
		repo:              repo,
		logService:        logService,
		projectCRSvc:      projectCRSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
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

func (pcrrs *ProjectChangeRequestReviewService) Submit(projectID uint, changeRequestID uint, currentUserID uint) error {
	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	if pcr.WorkflowStatusID != nil {
		workflowStatus, err := pcrrs.workflowStatusSvc.GetByID(*pcr.WorkflowStatusID)
		if err == nil && workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("project change request has already been submitted")
		}
	}

	if err := pcrrs.workflowSvc.ValidateWorkflowBeforeCreate(); err != nil {
		return err
	}

	tx := pcrrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrrs.workflowSvc.DeleteAllWorkflowsForProjectChangeRequest(changeRequestID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrrs.workflowSvc.CreateProjectChangeRequestWorkflow(changeRequestID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	pcr.WorkflowStatusID = &workflowStatus.ID
	pcr.SubmittedBYID = &currentUserID
	pcr.SubmittedAt = &currentDate
	pcr.NextReviewLevel = &level
	pcr.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
		tx.Rollback()
		return err
	}

	workflowLevel, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrrs.createReviewWithHistory(workflowLevel.ID, changeRequestID, level, workflowStatus.ID, nil, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcrrs.logService.LogActivity(
			models.CreateAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
	if err == nil {
		pcrrs.notifyNextLevelUsers(projectID, changeRequestID, workflowLevel, pcrDetail, currentUserID)
	}

	return nil
}

func (pcrrs *ProjectChangeRequestReviewService) Resubmit(projectID uint, changeRequestID uint, currentUserID uint) error {
	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	if pcr.WorkflowStatusID == nil {
		return utility.NewForbiddenError("project change request must be reverted before resubmitting")
	}

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByID(*pcr.WorkflowStatusID)
	if err != nil || workflowStatus.Type != "revert" {
		return utility.NewForbiddenError("project change request is not in reverted state")
	}

	if err := pcrrs.workflowSvc.ValidateWorkflowBeforeCreate(); err != nil {
		return err
	}

	tx := pcrrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	submitStatus, err := pcrrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrrs.workflowSvc.DeleteAllWorkflowsForProjectChangeRequest(changeRequestID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrrs.workflowSvc.CreateProjectChangeRequestWorkflow(changeRequestID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	pcr.WorkflowStatusID = &submitStatus.ID
	pcr.SubmittedBYID = &currentUserID
	pcr.SubmittedAt = &currentDate
	pcr.NextReviewLevel = &level
	pcr.UpdatedBYID = &currentUserID
	pcr.RevertedBYID = nil
	pcr.RevertedAt = nil

	if err := pcrrs.repo.UpdateChangeRequest(pcr, tx); err != nil {
		tx.Rollback()
		return err
	}

	workflowLevel, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pcrrs.createReviewWithHistory(workflowLevel.ID, changeRequestID, level, submitStatus.ID, nil, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcrrs.logService.LogActivity(
			models.UpdateAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
	if err == nil {
		pcrrs.notifyNextLevelUsers(projectID, changeRequestID, workflowLevel, pcrDetail, currentUserID)
	}

	return nil
}

func (pcrrs *ProjectChangeRequestReviewService) Approve(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) error {
	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel == nil {
		return utility.NewForbiddenError("no pending review level found")
	}

	tx := pcrrs.repo.BeginTransaction()
	currentApprovalLevel := pcr.NextReviewLevel
	currentWorkflow, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := pcrrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	var comment *string
	if form != nil && form.Comment != nil && *form.Comment != "" {
		comment = form.Comment
	}

	var nextWorkflowLevel *models.ProjectChangeRequestWorkflow

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.ProjectChangeRequestWorkflowParticipants) == 0 ||
			currentUserID != currentWorkflow.ProjectChangeRequestWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pcrrs.createReviewWithHistory(currentWorkflow.ID, changeRequestID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		projectChangeRequestWorkflow, err := pcrrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		projectChangeRequestWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(projectChangeRequestWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err = pcrrs.workflowSvc.GetNextLevel(changeRequestID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingPCR, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingPCR.WorkflowStatusID = &workflowStatus.ID
			existingPCR.ApprovedBYID = &currentUserID
			existingPCR.ApprovedAt = &approvedAt
			existingPCR.NextReviewLevel = nil
		} else {
			existingPCR.NextReviewLevel = &nextWorkflowLevel.Level
		}
		existingPCR.UpdatedBYID = &currentUserID

		if err := pcrrs.repo.UpdateChangeRequest(existingPCR, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ProjectChangeRequestWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := pcrrs.createReviewWithHistory(currentWorkflow.ID, changeRequestID, currentWorkflow.Level, workflowStatus.ID, comment, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := pcrrs.workflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ProjectChangeRequestWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		projectChangeRequestWorkflow, err := pcrrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		projectChangeRequestWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(projectChangeRequestWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err = pcrrs.workflowSvc.GetNextLevel(changeRequestID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingPCR, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingPCR.WorkflowStatusID = &workflowStatus.ID
			existingPCR.ApprovedBYID = &currentUserID
			existingPCR.ApprovedAt = &approvedAt
			existingPCR.NextReviewLevel = nil
		} else {
			existingPCR.NextReviewLevel = &nextWorkflowLevel.Level
		}
		existingPCR.UpdatedBYID = &currentUserID

		if err := pcrrs.repo.UpdateChangeRequest(existingPCR, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	go func() {
		_ = pcrrs.logService.LogActivity(
			models.ApproveAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
	if err == nil {
		if nextWorkflowLevel != nil {
			nextWf, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, nextWorkflowLevel.Level)
			if err == nil {
				pcrrs.notifyNextLevelUsers(projectID, changeRequestID, nextWf, pcrDetail, currentUserID)
			}
		} else {
			pcrrs.sendApprovalNotification(projectID, changeRequestID, currentUserID, pcrDetail)
		}
	}

	return nil
}

func (pcrrs *ProjectChangeRequestReviewService) Revert(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) error {
	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel == nil {
		return utility.NewForbiddenError("no pending review level found")
	}

	tx := pcrrs.repo.BeginTransaction()
	currentApprovalLevel := pcr.NextReviewLevel
	currentWorkflow, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	revertStatus, err := pcrrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to revert")
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ProjectChangeRequestWorkflowParticipants {
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

	if err := pcrrs.createReviewWithHistory(currentWorkflow.ID, changeRequestID, currentWorkflow.Level, revertStatus.ID, comment, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	projectChangeRequestWorkflow, err := pcrrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	projectChangeRequestWorkflow.WorkflowStatusID = revertStatus.ID
	if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(projectChangeRequestWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	currentLevel := *currentApprovalLevel
	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.ProjectChangeRequestWorkflow
	if previousLevel > 0 {
		pendingStatus, err := pcrrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			tx.Rollback()
			return err
		}

		previousLevelWorkflow, err = pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, previousLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		var reviewIDs []*uint
		for _, r := range previousLevelWorkflow.ProjectChangeRequestWorkflowReviews {
			reviewIDs = append(reviewIDs, &r.ID)
		}
		if len(reviewIDs) > 0 {
			if err := pcrrs.workflowSvc.DeleteWorkflowReviews(reviewIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		previousLevelWorkflow.WorkflowStatusID = pendingStatus.ID
		if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(previousLevelWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingPCR, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if previousLevel == 0 {
		now := time.Now()
		existingPCR.WorkflowStatusID = &revertStatus.ID
		existingPCR.RevertedBYID = &currentUserID
		existingPCR.RevertedAt = &now
		existingPCR.NextReviewLevel = nil
	} else {
		existingPCR.NextReviewLevel = &previousLevel
	}
	existingPCR.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(existingPCR, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcrrs.logService.LogActivity(
			models.RejectAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
	if err == nil {
		if previousLevel > 0 {
			pcrrs.notifyNextLevelUsers(projectID, changeRequestID, previousLevelWorkflow, pcrDetail, currentUserID)
		} else {
			pcrrs.sendRevertNotification(projectID, changeRequestID, currentUserID, pcrDetail)
		}
	}

	return nil
}

func (pcrrs *ProjectChangeRequestReviewService) Reject(projectID uint, changeRequestID uint, currentUserID uint, form *projectDTOs.ProjectChangeRequestReviewFormDTO) error {
	pcr, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		return err
	}

	if pcr.NextReviewLevel == nil {
		return utility.NewForbiddenError("no pending review level found")
	}

	tx := pcrrs.repo.BeginTransaction()
	currentApprovalLevel := pcr.NextReviewLevel
	currentWorkflow, err := pcrrs.workflowSvc.GetWorkflowByLevel(changeRequestID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	rejectStatus, err := pcrrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		tx.Rollback()
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
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to reject")
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.ProjectChangeRequestWorkflowParticipants {
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

	if err := pcrrs.createReviewWithHistory(currentWorkflow.ID, changeRequestID, currentWorkflow.Level, rejectStatus.ID, comment, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	projectChangeRequestWorkflow, err := pcrrs.workflowSvc.GetWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	projectChangeRequestWorkflow.WorkflowStatusID = rejectStatus.ID
	if err := pcrrs.workflowSvc.UpdateProjectChangeRequestWorkflow(projectChangeRequestWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingPCR, err := pcrrs.projectCRSvc.GetByID(projectID, changeRequestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	existingPCR.WorkflowStatusID = &rejectStatus.ID
	existingPCR.RejectedBYID = &currentUserID
	existingPCR.RejectedAt = &now
	existingPCR.NextReviewLevel = nil
	existingPCR.UpdatedBYID = &currentUserID

	if err := pcrrs.repo.UpdateChangeRequest(existingPCR, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcrrs.logService.LogActivity(
			models.RejectAction,
			models.ProjectChangeRequest{}.TableName(),
			changeRequestID,
			pcr,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	pcrDetail, err := pcrrs.projectCRSvc.GetFullDetailsByID(projectID, changeRequestID)
	if err == nil {
		pcrrs.sendRejectNotification(projectID, changeRequestID, currentUserID, pcrDetail)
	}

	return nil
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

		err = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_change_request")
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

			_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_change_request")
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

		_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_change_request")
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

		_ = pcrrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "project_change_request")
	}()
}
