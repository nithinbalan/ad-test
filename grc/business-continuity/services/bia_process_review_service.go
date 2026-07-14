package services

import (
	"fmt"
	"time"

	businessContinuityRepositories "jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
)

type BiaProcessReviewServiceInterface interface {
	Submit(biaProcessID uint, currentUserID uint) error
	Approve(biaProcessID uint, currentUserID uint) error
	Reject(biaProcessID uint, currentUserID uint) error
	Reopen(biaProcessID uint, currentUserID uint) error
}

type BiaProcessReviewService struct {
	repo                  businessContinuityRepositories.BiaProcessReviewRepositoryInterface
	logService            logServices.ActivityLogServiceInterface
	biaProcessSvc         BiaProcessServiceInterface
	biaProcessWorkflowSvc BiaProcessWorkflowServiceInterface
	workflowStatusSvc     masterServices.WorkflowStatusServiceInterface
	userSvc               organizationSvc.UserServiceInterface
	notificationSvc       commonServices.NotificationServiceInterface
	processSvc            organizationSvc.ProcessServiceInterface
}

func NewBiaProcessReviewService(
	repo businessContinuityRepositories.BiaProcessReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	biaProcessSvc BiaProcessServiceInterface,
	biaProcessWorkflowSvc BiaProcessWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	processSvc organizationSvc.ProcessServiceInterface,
) BiaProcessReviewServiceInterface {
	return &BiaProcessReviewService{
		repo:                  repo,
		logService:            logService,
		biaProcessSvc:         biaProcessSvc,
		biaProcessWorkflowSvc: biaProcessWorkflowSvc,
		workflowStatusSvc:     workflowStatusSvc,
		userSvc:               userSvc,
		notificationSvc:       notificationSvc,
		processSvc:            processSvc,
	}
}

func (svc *BiaProcessReviewService) Submit(biaProcessID uint, currentUserID uint) error {
	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil {
		return err
	}

	if biaProcess.SubmittedByID != nil && biaProcess.WorkflowStatusID != nil {
		workflowStatus, err := svc.workflowStatusSvc.GetByID(*biaProcess.WorkflowStatusID)
		if err != nil {
			return err
		}
		if workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("bia_already_submitted")
		}
	}

	if err := svc.biaProcessWorkflowSvc.ValidateWorkflowBeforeCreate("bia"); err != nil {
		return err
	}

	tx := svc.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := svc.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := svc.biaProcessWorkflowSvc.DeleteAllWorkflowsForBiaProcess(biaProcessID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := svc.biaProcessWorkflowSvc.CreateBiaProcessWorkflow(biaProcessID, "bia", currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	biaProcess.WorkflowStatusID = &workflowStatus.ID
	biaProcess.SubmittedByID = &currentUserID
	biaProcess.SubmittedAt = &currentDate
	biaProcess.NextReviewLevel = &level
	biaProcess.UpdatedBYID = &currentUserID

	if err := svc.repo.UpdateBiaProcess(biaProcess, tx); err != nil {
		tx.Rollback()
		return err
	}

	biaProcessWorkflowLevel, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByLevel(biaProcessID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	review := &models.BiaProcessWorkflowReview{
		BiaProcessWorkflowID: biaProcessWorkflowLevel.ID,
		WorkflowStatusID:     workflowStatus.ID,
		CreatedBYID:          currentUserID,
	}

	if err := svc.repo.CreateBiaProcessReview(review, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := svc.biaProcessWorkflowSvc.CreateWorkflowHistory(biaProcessID, level, workflowStatus.ID, "", currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go svc.notifyReviewUsers(biaProcessID, biaProcessWorkflowLevel, currentUserID)

	return nil
}

func (svc *BiaProcessReviewService) Approve(biaProcessID uint, currentUserID uint) error {
	tx := svc.repo.BeginTransaction()

	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if biaProcess.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := biaProcess.NextReviewLevel
	currentWorkflow, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByLevel(biaProcessID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := svc.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "custom-roles", "hod":
		if len(currentWorkflow.BiaProcessWorkflowParticipants) == 0 {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("workflow has no participants"), "not_exist", "workflow")
		}
		if currentUserID != currentWorkflow.BiaProcessWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		review := &models.BiaProcessWorkflowReview{
			BiaProcessWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:     workflowStatus.ID,
			CreatedBYID:          currentUserID,
		}
		if err := svc.repo.CreateBiaProcessReview(review, tx); err != nil {
			tx.Rollback()
			return err
		}

		currentWorkflow2, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		currentWorkflow2.WorkflowStatusID = workflowStatus.ID
		if err := svc.biaProcessWorkflowSvc.UpdateBiaProcessWorkflow(currentWorkflow2, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := svc.biaProcessWorkflowSvc.GetNextLevel(biaProcessID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingBiaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingBiaProcess.WorkflowStatusID = &workflowStatus.ID
			existingBiaProcess.ApprovedByID = &currentUserID
			existingBiaProcess.ApprovedAt = &approvedAt
			existingBiaProcess.NextReviewLevel = nil
		} else {
			existingBiaProcess.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := svc.repo.UpdateBiaProcess(existingBiaProcess, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := svc.biaProcessWorkflowSvc.CreateWorkflowHistory(biaProcessID, *currentApprovalLevel, workflowStatus.ID, "", currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByLevel(biaProcessID, nextWorkflowLevel.Level)
			if err == nil {
				go svc.notifyReviewUsers(biaProcessID, nextWorkflow, currentUserID)
			}
		}
		go svc.sendApprovalNotification(biaProcessID, currentUserID)

		return nil

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.BiaProcessWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		review := &models.BiaProcessWorkflowReview{
			BiaProcessWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:     workflowStatus.ID,
			CreatedBYID:          currentUserID,
		}
		if err := svc.repo.CreateBiaProcessReview(review, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := svc.biaProcessWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.BiaProcessWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		currentWorkflow2, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		currentWorkflow2.WorkflowStatusID = workflowStatus.ID
		if err := svc.biaProcessWorkflowSvc.UpdateBiaProcessWorkflow(currentWorkflow2, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := svc.biaProcessWorkflowSvc.GetNextLevel(biaProcessID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingBiaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingBiaProcess.WorkflowStatusID = &workflowStatus.ID
			existingBiaProcess.ApprovedByID = &currentUserID
			existingBiaProcess.ApprovedAt = &approvedAt
			existingBiaProcess.NextReviewLevel = nil
		} else {
			existingBiaProcess.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := svc.repo.UpdateBiaProcess(existingBiaProcess, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := svc.biaProcessWorkflowSvc.CreateWorkflowHistory(biaProcessID, *currentApprovalLevel, workflowStatus.ID, "", currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByLevel(biaProcessID, nextWorkflowLevel.Level)
			if err == nil {
				go svc.notifyReviewUsers(biaProcessID, nextWorkflow, currentUserID)
			}
		}
		go svc.sendApprovalNotification(biaProcessID, currentUserID)

		return nil
	}

	return nil
}

func (svc *BiaProcessReviewService) Reject(biaProcessID uint, currentUserID uint) error {
	workflowStatus, err := svc.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	tx := svc.repo.BeginTransaction()

	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if biaProcess.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := biaProcess.NextReviewLevel
	currentWorkflow, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByLevel(biaProcessID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	review := &models.BiaProcessWorkflowReview{
		BiaProcessWorkflowID: currentWorkflow.ID,
		WorkflowStatusID:     workflowStatus.ID,
		CreatedBYID:          currentUserID,
	}
	if err := svc.repo.CreateBiaProcessReview(review, tx); err != nil {
		tx.Rollback()
		return err
	}

	currentWorkflow2, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	currentWorkflow2.WorkflowStatusID = workflowStatus.ID
	if err := svc.biaProcessWorkflowSvc.UpdateBiaProcessWorkflow(currentWorkflow2, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingBiaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingBiaProcess.WorkflowStatusID = &workflowStatus.ID
	existingBiaProcess.NextReviewLevel = nil
	existingBiaProcess.SubmittedAt = nil
	existingBiaProcess.SubmittedByID = nil
	existingBiaProcess.UpdatedBYID = &currentUserID
	if err := svc.repo.UpdateBiaProcess(existingBiaProcess, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := svc.biaProcessWorkflowSvc.CreateWorkflowHistory(biaProcessID, *currentApprovalLevel, workflowStatus.ID, "", currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	go svc.sendRevertNotification(biaProcessID, currentUserID)

	return tx.Commit()
}

func (svc *BiaProcessReviewService) Reopen(biaProcessID uint, currentUserID uint) error {
	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil {
		return err
	}

	if biaProcess.SubmittedByID != nil && biaProcess.WorkflowStatusID != nil {
		workflowStatus, err := svc.workflowStatusSvc.GetByID(*biaProcess.WorkflowStatusID)
		if err != nil {
			return err
		}
		if workflowStatus.Type != "approve" {
			return utility.NewForbiddenError("bia process already submitted")
		}
	}

	if err := svc.biaProcessWorkflowSvc.ValidateWorkflowBeforeCreate("bia"); err != nil {
		return err
	}

	tx := svc.repo.BeginTransaction()

	workflowStatus, err := svc.workflowStatusSvc.GetByType("reopen")
	if err != nil {
		tx.Rollback()
		return err
	}

	if biaProcess.ApprovedByID != nil {
		levelOneWorkflow, err := svc.biaProcessWorkflowSvc.GetBiaProcessWorkflowByLevel(biaProcessID, 1, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		review := &models.BiaProcessWorkflowReview{
			BiaProcessWorkflowID: levelOneWorkflow.ID,
			WorkflowStatusID:     workflowStatus.ID,
			CreatedBYID:          currentUserID,
		}
		if err := svc.repo.CreateBiaProcessReview(review, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := svc.biaProcessWorkflowSvc.CreateWorkflowHistory(biaProcessID, 1, workflowStatus.ID, "", currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	biaProcess.WorkflowStatusID = &workflowStatus.ID
	biaProcess.SubmittedByID = nil
	biaProcess.SubmittedAt = nil
	biaProcess.NextReviewLevel = nil
	biaProcess.ApprovedAt = nil
	biaProcess.ApprovedByID = nil
	biaProcess.UpdatedBYID = &currentUserID

	if err := svc.repo.UpdateBiaProcess(biaProcess, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go svc.sendReopenNotification(biaProcessID, currentUserID)

	return nil
}

func (svc *BiaProcessReviewService) getProcess(biaProcessID uint) (*models.Process, error) {
	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil {
		return nil, err
	}
	return svc.processSvc.GetByID(biaProcess.ProcessID)
}

func (svc *BiaProcessReviewService) notifyReviewUsers(biaProcessID uint, workflow *models.BiaProcessWorkflow, currentUserID uint) {
	process, err := svc.getProcess(biaProcessID)
	if err != nil {
		utility.LogError("Failed to fetch process for notification", map[string]interface{}{"error": err.Error()})
		return
	}
	title := process.TitleEn
	entity := "bia-process-review"
	action := "review-user-assign"
	urlParam := fmt.Sprintf(`{"id":%d}`, process.ID)
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%sbusiness-continuity/business-impact-analysis/%d/overview", appURL, process.ID)

	assignedBy, err := svc.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch assigning user", map[string]interface{}{"error": err.Error()})
		return
	}

	for _, participant := range workflow.BiaProcessWorkflowParticipants {
		responsibleUser, err := svc.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{"error": err.Error()})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "business_continuity",
			NotificationType:     string(commonDTOs.ReviewType),
			HeadingEn:            "BIA Process Review",
			HeadingAr:            "مراجعة عملية تقييم الأثر",
			TitleEn:              title,
			TitleAr:              title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"Title":                 title,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              "BIA Process: Review Assigned",
			TemplateRelativePath: "business-continuity/templates/email/bia_process_workflow_review_user_mail.html",
			EmailData:            emailData,
		}

		if err := svc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "bia"); err != nil {
			utility.LogError("Failed to send review notification", map[string]interface{}{"error": err.Error()})
		}
	}
}

func (svc *BiaProcessReviewService) sendApprovalNotification(biaProcessID uint, currentUserID uint) {
	process, err := svc.getProcess(biaProcessID)
	if err != nil {
		utility.LogError("Failed to fetch process for notification", map[string]interface{}{"error": err.Error()})
		return
	}
	title := process.TitleEn
	entity := "bia-process-review"
	action := "approval"
	urlParam := fmt.Sprintf(`{"id":%d}`, process.ID)
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%sbusiness-continuity/business-impact-analysis/%d/overview", appURL, process.ID)

	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil || biaProcess.CreatedBYID == nil {
		return
	}

	createdByUser, err := svc.userSvc.GetByID(*biaProcess.CreatedBYID)
	if err != nil {
		utility.LogError("Failed to fetch BIA process creator", map[string]interface{}{"error": err.Error()})
		return
	}

	approvedByUser, err := svc.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch approving user", map[string]interface{}{"error": err.Error()})
		return
	}

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*biaProcess.CreatedBYID},
		ModuleIdentifier:     "business_continuity",
		NotificationType:     string(commonDTOs.ApprovalType),
		HeadingEn:            "BIA Process Approved",
		HeadingAr:            "تمت الموافقة على عملية تقييم الأثر",
		TitleEn:              title,
		TitleAr:              title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           createdByUser.FirstNameEn,
		"LastNameEn":            createdByUser.LastNameEn,
		"Title":                 title,
		"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
		"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
		"ApprovedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{createdByUser.Email},
		Subject:              "BIA Process: Approved",
		TemplateRelativePath: "business-continuity/templates/email/bia_process_workflow_approve_user_mail.html",
		EmailData:            emailData,
	}

	if err := svc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "bia"); err != nil {
		utility.LogError("Failed to send approval notification", map[string]interface{}{"error": err.Error()})
	}
}

func (svc *BiaProcessReviewService) sendRevertNotification(biaProcessID uint, currentUserID uint) {
	process, err := svc.getProcess(biaProcessID)
	if err != nil {
		utility.LogError("Failed to fetch process for notification", map[string]interface{}{"error": err.Error()})
		return
	}
	title := process.TitleEn
	entity := "bia-process-review"
	action := "reverted"
	urlParam := fmt.Sprintf(`{"id":%d}`, process.ID)
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%sbusiness-continuity/business-impact-analysis/%d/overview", appURL, process.ID)

	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil || biaProcess.CreatedBYID == nil {
		return
	}

	createdByUser, err := svc.userSvc.GetByID(*biaProcess.CreatedBYID)
	if err != nil {
		utility.LogError("Failed to fetch BIA process creator", map[string]interface{}{"error": err.Error()})
		return
	}

	revertedByUser, err := svc.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{"error": err.Error()})
		return
	}

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*biaProcess.CreatedBYID},
		ModuleIdentifier:     "business_continuity",
		NotificationType:     string(commonDTOs.RejectType),
		HeadingEn:            "BIA Process Reverted",
		HeadingAr:            "تم إعادة عملية تقييم الأثر",
		TitleEn:              title,
		TitleAr:              title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           createdByUser.FirstNameEn,
		"LastNameEn":            createdByUser.LastNameEn,
		"Title":                 title,
		"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
		"RevertedByLastNameEn":  revertedByUser.LastNameEn,
		"RevertedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{createdByUser.Email},
		Subject:              "BIA Process: Reverted",
		TemplateRelativePath: "business-continuity/templates/email/bia_process_workflow_revert_user_mail.html",
		EmailData:            emailData,
	}

	if err := svc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "bia"); err != nil {
		utility.LogError("Failed to send revert notification", map[string]interface{}{"error": err.Error()})
	}
}

func (svc *BiaProcessReviewService) sendReopenNotification(biaProcessID uint, currentUserID uint) {
	process, err := svc.getProcess(biaProcessID)
	if err != nil {
		utility.LogError("Failed to fetch process for notification", map[string]interface{}{"error": err.Error()})
		return
	}
	title := process.TitleEn
	entity := "bia-process-review"
	action := "reopened"
	urlParam := fmt.Sprintf(`{"id":%d}`, process.ID)
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%sbusiness-continuity/business-impact-analysis/%d/overview", appURL, process.ID)

	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil || biaProcess.CreatedBYID == nil {
		return
	}

	createdByUser, err := svc.userSvc.GetByID(*biaProcess.CreatedBYID)
	if err != nil {
		utility.LogError("Failed to fetch BIA process creator", map[string]interface{}{"error": err.Error()})
		return
	}

	reopenedByUser, err := svc.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reopening user", map[string]interface{}{"error": err.Error()})
		return
	}

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*biaProcess.CreatedBYID},
		ModuleIdentifier:     "business_continuity",
		NotificationType:     string(commonDTOs.RevertType),
		HeadingEn:            "BIA Process Reopened",
		HeadingAr:            "تم إعادة فتح عملية تقييم الأثر",
		TitleEn:              title,
		TitleAr:              title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           createdByUser.FirstNameEn,
		"LastNameEn":            createdByUser.LastNameEn,
		"Title":                 title,
		"ReopenedByFirstNameEn": reopenedByUser.FirstNameEn,
		"ReopenedByLastNameEn":  reopenedByUser.LastNameEn,
		"ReopenedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{createdByUser.Email},
		Subject:              "BIA Process: Reopened",
		TemplateRelativePath: "business-continuity/templates/email/bia_process_workflow_reopen_user_mail.html",
		EmailData:            emailData,
	}

	if err := svc.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "bia"); err != nil {
		utility.LogError("Failed to send reopen notification", map[string]interface{}{"error": err.Error()})
	}
}
