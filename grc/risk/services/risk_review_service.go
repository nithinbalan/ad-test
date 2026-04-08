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
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
)

type RiskReviewServiceInterface interface {
	Submit(riskID uint, currentUserID uint) error
	Approve(riskID uint, currentUserID uint) error
	Reject(riskID uint, currentUserID uint) error
}

type RiskReviewService struct {
	repo              repositories.RiskReviewRepositoryInterface
	logService        logServices.ActivityLogServiceInterface
	riskSvc           RiskRegisterServiceInterface
	riskWorkflowSvc   RiskWorkflowServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	riskStatusSvc     RiskStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
	taskSvc           taskServices.TaskServiceInterface
}

func NewRiskReviewService(
	repo repositories.RiskReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	riskSvc RiskRegisterServiceInterface,
	riskWorkflowSvc RiskWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) RiskReviewServiceInterface {
	return &RiskReviewService{
		repo:              repo,
		logService:        logService,
		riskSvc:           riskSvc,
		riskWorkflowSvc:   riskWorkflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		riskStatusSvc:     riskStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
		taskSvc:           taskSvc,
	}
}

func (rrs *RiskReviewService) Submit(riskID uint, currentUserID uint) error {
	riskRegister, err := rrs.riskSvc.GetByID(riskID)
	if err != nil {
		return err
	}

	tx := rrs.repo.BeginTransaction()
	currentDateStr := time.Now()

	riskStatus, err := rrs.riskStatusSvc.GetByType("submitted")
	if err != nil {
		return err
	}
	var level int

	level = 1

	riskRegister.RiskStatusID = riskStatus.ID
	riskRegister.SubmittedBy = &currentUserID
	riskRegister.SubmittedAt = &currentDateStr
	riskRegister.NextReviewLevel = &level
	riskRegister.UpdatedBYID = &currentUserID

	err = rrs.repo.UpdateRisk(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Fetch the workflow level details, including participants (review users)
	riskWorkflowLevel, err := rrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskRegister.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	workflowStatus, err := rrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}
	riskReview := &models.RiskWorkflowReview{
		RiskWorkflowID:   riskWorkflowLevel.ID,
		WorkflowStatusID: workflowStatus.ID,
		CreatedBYID:      currentUserID,
	}

	err = rrs.repo.CreateRiskReview(riskReview, tx)
	if err != nil {
		return err
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskRegister.ID)
	tx.Commit()

	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, riskID)
		entity := "risk-review"
		action := "review-user-assign"

		for _, participant := range riskWorkflowLevel.RiskWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.User.ID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      riskRegister.ReferenceNumber,
				HeadingEn:            "Risk Review",
				HeadingAr:            "مراجعة المخاطر",
				TitleEn:              riskRegister.Title,
				TitleAr:              riskRegister.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := rrs.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := rrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Risk %s: Review Assigned", riskRegister.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 riskRegister.Title,
				"ReferenceNumber":       riskRegister.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, riskID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (rrs *RiskReviewService) Approve(riskID uint, currentUserID uint) error {
	tx := rrs.repo.BeginTransaction()

	riskRegister, err := rrs.riskSvc.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentApprovalLevel := riskRegister.NextReviewLevel
	currentRiskWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := rrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentRiskWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentRiskWorkflow.RiskWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		riskReview := &models.RiskWorkflowReview{
			RiskWorkflowID:   currentRiskWorkflow.ID,
			WorkflowStatusID: workflowStatus.ID,
			CreatedBYID:      currentUserID,
		}
		if err := rrs.repo.CreateRiskReview(riskReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		riskWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByID(currentRiskWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		riskWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := rrs.riskWorkflowSvc.UpdateRiskWorkflow(riskWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextRiskWorkflowLevel, err := rrs.riskWorkflowSvc.GetNextLevel(uint(riskRegister.ID), *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextRiskWorkflowLevel == nil {
			riskStatus, err := rrs.riskStatusSvc.GetByType("open")
			if err != nil {
				tx.Rollback()
				return err
			}
			existingRisk, err := rrs.riskSvc.GetByID(riskID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			existingRisk.RiskStatusID = riskStatus.ID
			if err := rrs.repo.UpdateRisk(existingRisk, tx); err != nil {
				tx.Rollback()
				return err
			}
		} else {
			existingRisk, err := rrs.riskSvc.GetByID(riskID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			existingRisk.NextReviewLevel = &nextRiskWorkflowLevel.Level
			if err := rrs.repo.UpdateRisk(existingRisk, tx); err != nil {
				tx.Rollback()
				return err
			}

			nextWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, nextRiskWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			rrs.notifyNextLevelUsers(riskID, nextWorkflow, riskRegister)
		}

		riskDetails, err := rrs.riskSvc.GetFullDetailsByID(riskID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		rrs.sendApprovalNotification(riskID, currentUserID, riskDetails)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentRiskWorkflow.RiskWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		riskReview := &models.RiskWorkflowReview{
			RiskWorkflowID:   currentRiskWorkflow.ID,
			WorkflowStatusID: workflowStatus.ID,
			CreatedBYID:      currentUserID,
		}
		if err := rrs.repo.CreateRiskReview(riskReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentRiskWorkflow.ApprovalTypes == "all" {
			allApproved, err := rrs.riskWorkflowSvc.CheckAllParticipantsApproved(currentRiskWorkflow.ID, currentRiskWorkflow.RiskWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				tx.Commit()
				return nil
			}
		}

		// update workflow status
		riskWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByID(currentRiskWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		riskWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := rrs.riskWorkflowSvc.UpdateRiskWorkflow(riskWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextRiskWorkflowLevel, err := rrs.riskWorkflowSvc.GetNextLevel(uint(riskRegister.ID), *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextRiskWorkflowLevel == nil {
			riskStatus, err := rrs.riskStatusSvc.GetByType("open")
			if err != nil {
				tx.Rollback()
				return err
			}
			existingRisk, err := rrs.riskSvc.GetByID(riskID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			existingRisk.RiskStatusID = riskStatus.ID
			if err := rrs.repo.UpdateRisk(existingRisk, tx); err != nil {
				tx.Rollback()
				return err
			}
		} else {
			existingRisk, err := rrs.riskSvc.GetByID(riskID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			existingRisk.NextReviewLevel = &nextRiskWorkflowLevel.Level
			if err := rrs.repo.UpdateRisk(existingRisk, tx); err != nil {
				tx.Rollback()
				return err
			}

			nextWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByLevel(uint(riskRegister.ID), nextRiskWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			rrs.notifyNextLevelUsers(uint(riskRegister.ID), nextWorkflow, riskRegister)
		}

		riskDetails, err := rrs.riskSvc.GetFullDetailsByID(riskID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		rrs.sendApprovalNotification(riskID, currentUserID, riskDetails)
	}

	// ✅ Only commit once at the end
	return tx.Commit()
}

func (rrs *RiskReviewService) Reject(riskID uint, currentUserID uint) error {
	riskWorkflowStatus, err := rrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	riskStatus, err := rrs.riskStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}

	tx := rrs.repo.BeginTransaction()

	riskRegister, err := rrs.riskSvc.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentApprovalLevel := riskRegister.NextReviewLevel
	currentRiskWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskWorkflowReview := &models.RiskWorkflowReview{
		RiskWorkflowID:   currentRiskWorkflow.ID,
		WorkflowStatusID: riskWorkflowStatus.ID,
		CreatedBYID:      currentUserID,
	}

	if err := rrs.repo.CreateRiskReview(riskWorkflowReview, tx); err != nil {
		tx.Rollback()
		return err
	}

	riskWorkflow, err := rrs.riskWorkflowSvc.GetRiskWorkflowByID(currentRiskWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskWorkflow.WorkflowStatusID = riskWorkflowStatus.ID
	if err := rrs.riskWorkflowSvc.UpdateRiskWorkflow(riskWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRisk, err := rrs.riskSvc.GetByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRisk.RiskStatusID = riskStatus.ID
	if err := rrs.repo.UpdateRisk(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction after all DB operations succeed
	tx.Commit()

	// Notifications should be sent AFTER commit succeeds
	rrs.sendRejectedNotification(riskID, currentUserID, riskRegister)

	return nil
}

func (rrs *RiskReviewService) sendApprovalNotification(id uint, currentUserID uint, riskDetail *riskDTOs.RiskRegisterDetailDTO) {
	go func() {
		// Only send notification if responsible user exists
		if riskDetail.ResponsibleUser == nil {
			return
		}

		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.ID)
		entity := "risk-review"
		action := "approval"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(riskDetail.ResponsibleUser.ID)},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Risk Approval",
			HeadingAr:            "الموافقة على المخاطر",
			ReferenceNumber:      riskDetail.ReferenceNumber,
			TitleEn:              riskDetail.Title,
			TitleAr:              riskDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		approvedByUser, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{riskDetail.ResponsibleUser.Email}
		subject := fmt.Sprintf("Risk %s: Approval Notification", riskDetail.Title)
		urlParam = fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, riskDetail.ID)

		if riskDetail.RiskStatus.Type == "open" {
			urlParam = fmt.Sprintf("%s/risk/risk-register/%d/overview", appURL, id)
		}

		emailData := map[string]interface{}{
			"Title":                 riskDetail.Title,
			"Description":           riskDetail.Description,
			"FirstNameEn":           riskDetail.ResponsibleUser.FirstNameEN,
			"LastNameEn":            riskDetail.ResponsibleUser.LastNameEN,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (rrs *RiskReviewService) notifyNextLevelUsers(id uint, nextRiskWorkflowLevel *models.RiskWorkflow, riskDetail *riskDTOs.RiskRegisterDetailDTO) {
	go func() {
		// Only send notification if responsible user exists
		if riskDetail.ResponsibleUser == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.ID)
		entity := "risk-review"
		action := "review-user-assign"

		for _, participant := range nextRiskWorkflowLevel.RiskWorkflowParticipants {
			responsibleUser, err := rrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(riskDetail.ResponsibleUserID)
			// Prepare database notification DTO
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Risk Review",
				HeadingAr:            "مراجعة المخاطر",
				ReferenceNumber:      riskDetail.ReferenceNumber,
				TitleEn:              riskDetail.Title,
				TitleAr:              riskDetail.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			// Prepare email notification DTO
			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Risk %s: Review Assigned", riskDetail.Title)
			appURL := utility.GetAppURL()
			urlParam = fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, riskDetail.ID)

			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 riskDetail.Title,
				"Description":           riskDetail.Description,
				"AssignedByFirstNameEn": riskDetail.ResponsibleUser.FirstNameEN,
				"AssignedByLastNameEn":  riskDetail.ResponsibleUser.LastNameEN,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           urlParam,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send both database and email notifications
			err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (rrs *RiskReviewService) sendRejectedNotification(id uint, currentUserID uint, risk *riskDTOs.RiskRegisterDetailDTO) {
	go func() {
		// Only send notification if responsible user exists
		if risk.ResponsibleUser == nil {
			return
		}

		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "risk-review"
		action := "rejected"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(risk.ResponsibleUser.ID)},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Risk Rejected",
			HeadingAr:            "تم رفض المخاطر",
			ReferenceNumber:      risk.ReferenceNumber,
			TitleEn:              risk.Title,
			TitleAr:              risk.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		rejectedByUser, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{risk.ResponsibleUser.Email}
		subject := fmt.Sprintf("Risk %s: Rejected Notification", risk.Title)
		urlParam = fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, risk.ID)

		emailData := map[string]interface{}{
			"FirstNameEn":           risk.ResponsibleUser.FirstNameEN,
			"LastNameEn":            risk.ResponsibleUser.LastNameEN,
			"Title":                 risk.Title,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
