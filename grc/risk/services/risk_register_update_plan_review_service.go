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
	"jethur/grc/risk/repositories"
	"jethur/utility"
)

type RiskRegisterUpdatePlanReviewServiceInterface interface {
	Submit(riskRegisterUpdatePlanID uint, currentUserID uint) error
	Approve(riskRegisterUpdatePlanID uint, currentUserID uint) error
	Reject(riskRegisterUpdatePlanID uint, currentUserID uint) error
}

type RiskRegisterUpdatePlanReviewService struct {
	repo                              repositories.RiskRegisterUpdatePlanReviewRepositoryInterface
	logService                        logServices.ActivityLogServiceInterface
	riskRegisterUpdatePlanSvc         RiskRegisterUpdatePlanServiceInterface
	riskRegisterUpdatePlanWorkflowSvc RiskRegisterUpdatePlanWorkflowServiceInterface
	workflowStatusSvc                 masterServices.WorkflowStatusServiceInterface
	userSvc                           organizationSvc.UserServiceInterface
	notificationSvc                   commonServices.NotificationServiceInterface
	planStatusSvc                     RiskRegisterUpdatePlanStatusServiceInterface
	riskRegisterSvc                   RiskRegisterServiceInterface
}

func NewRiskRegisterUpdatePlanReviewService(
	repo repositories.RiskRegisterUpdatePlanReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	riskRegisterUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	riskRegisterUpdatePlanWorkflowSvc RiskRegisterUpdatePlanWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	planStatusSvc RiskRegisterUpdatePlanStatusServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
) RiskRegisterUpdatePlanReviewServiceInterface {
	return &RiskRegisterUpdatePlanReviewService{
		repo:                              repo,
		logService:                        logService,
		riskRegisterUpdatePlanSvc:         riskRegisterUpdatePlanSvc,
		riskRegisterUpdatePlanWorkflowSvc: riskRegisterUpdatePlanWorkflowSvc,
		workflowStatusSvc:                 workflowStatusSvc,
		userSvc:                           userSvc,
		notificationSvc:                   notificationSvc,
		planStatusSvc:                     planStatusSvc,
		riskRegisterSvc:                   riskRegisterSvc,
	}
}

func (rruprs *RiskRegisterUpdatePlanReviewService) Submit(riskRegisterUpdatePlanID uint, currentUserID uint) error {
	riskRegisterUpdatePlan, err := rruprs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
	if err != nil {
		return err
	}

	if riskRegisterUpdatePlan.SubmittedBYID != nil && riskRegisterUpdatePlan.WorkflowStatusID != nil {
		workflowStatus, err := rruprs.workflowStatusSvc.GetByID(*riskRegisterUpdatePlan.WorkflowStatusID)
		if err != nil {
			return err
		}

		if workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("risk_register_update_plan_already_submitted")
		}
	}

	err = rruprs.riskRegisterUpdatePlanWorkflowSvc.ValidateWorkflowBeforeCreate("risk_register_update_plan")
	if err != nil {
		return err
	}

	tx := rruprs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := rruprs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = rruprs.riskRegisterUpdatePlanWorkflowSvc.DeleteAllWorkflowsForRiskRegisterUpdatePlan(riskRegisterUpdatePlanID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = rruprs.riskRegisterUpdatePlanWorkflowSvc.CreateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanID, "risk_register_update_plan", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskRegisterUpdatePlan.WorkflowStatusID = &workflowStatus.ID
	riskRegisterUpdatePlan.SubmittedBYID = &currentUserID
	riskRegisterUpdatePlan.SubmittedAt = &currentDate
	riskRegisterUpdatePlan.NextReviewLevel = &level
	riskRegisterUpdatePlan.UpdatedBYID = &currentUserID

	err = rruprs.repo.UpdateRiskRegisterUpdatePlan(riskRegisterUpdatePlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskRegisterUpdatePlanWorkflowLevel, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create submit review in both current and history tables
	if err := rruprs.createReviewWithHistory(riskRegisterUpdatePlanWorkflowLevel.ID, riskRegisterUpdatePlanID, level, workflowStatus.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, riskRegisterUpdatePlan.ID)
		entity := "risk-register-update-plan-review"
		action := "review-user-assign"

		planTitle := fmt.Sprintf("Update Plan: %s to %s", riskRegisterUpdatePlan.StartDate.Format("2006-01-02"), riskRegisterUpdatePlan.EndDate.Format("2006-01-02"))

		// Notify Risk Champion
		riskChampionNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{riskRegisterUpdatePlan.RiskChampionID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.ReviewType),
			HeadingEn:            "Risk Register Update Plan Submitted",
			HeadingAr:            "تم تقديم خطة تحديث سجل المخاطر",
			TitleEn:              planTitle,
			TitleAr:              planTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		riskChampion, err := rruprs.userSvc.GetByID(riskRegisterUpdatePlan.RiskChampionID)
		if err != nil {
			utility.LogError("Failed to fetch risk champion", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			submittedBy, err := rruprs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch submitting user", map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				emailTo := []string{riskChampion.Email}
				subject := "Risk Register Update Plan: Submitted"
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":            riskChampion.FirstNameEn,
					"LastNameEn":             riskChampion.LastNameEn,
					"Title":                  planTitle,
					"SubmittedByFirstNameEn": submittedBy.FirstNameEn,
					"SubmittedByLastNameEn":  submittedBy.LastNameEn,
					"SubmittedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%srisk/risk-register-update-plan/%d/overview", appURL, riskRegisterUpdatePlan.ID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "risk/templates/email/risk_register_update_plan_submit_champion_mail.html",
					EmailData:            emailData,
				}

				rruprs.notificationSvc.SendNotification(riskChampionNotificationDTO, emailNotificationDTO, "risk_register_update_plan")
			}
		}

		// Notify workflow participants
		for _, participant := range riskRegisterUpdatePlanWorkflowLevel.RiskRegisterUpdatePlanWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Risk Register Update Plan Review",
				HeadingAr:            "مراجعة خطة تحديث سجل المخاطر",
				TitleEn:              planTitle,
				TitleAr:              planTitle,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := rruprs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := rruprs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := "Risk Register Update Plan: Review Assigned"
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 planTitle,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/risk-register-update-plan/%d/overview", appURL, riskRegisterUpdatePlan.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_register_update_plan_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = rruprs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_register_update_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (rruprs *RiskRegisterUpdatePlanReviewService) Approve(riskRegisterUpdatePlanID uint, currentUserID uint) error {
	riskRegisterUpdatePlan, err := rruprs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
	if err != nil {
		return err
	}

	tx := rruprs.repo.BeginTransaction()

	if riskRegisterUpdatePlan.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := riskRegisterUpdatePlan.NextReviewLevel
	currentWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := rruprs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentWorkflow.RiskRegisterUpdatePlanWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		// Create approval review in both current and history tables
		if err := rruprs.createReviewWithHistory(currentWorkflow.ID, riskRegisterUpdatePlanID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Update participant status
		for _, participant := range currentWorkflow.RiskRegisterUpdatePlanWorkflowParticipants {
			if participant.UserID == currentUserID {
				participant.WorkflowStatusID = workflowStatus.ID
				if err := rruprs.repo.UpdateRiskRegisterUpdatePlanWorkflowParticipant(participant, tx); err != nil {
					tx.Rollback()
					return err
				}
				break
			}
		}

		riskRegisterUpdatePlanWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		riskRegisterUpdatePlanWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := rruprs.riskRegisterUpdatePlanWorkflowSvc.UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetNextLevel(riskRegisterUpdatePlanID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingRiskRegisterUpdatePlan, err := rruprs.riskRegisterUpdatePlanSvc.GetByPlanID(riskRegisterUpdatePlanID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingRiskRegisterUpdatePlan.WorkflowStatusID = &workflowStatus.ID
			existingRiskRegisterUpdatePlan.ApprovedBYID = &currentUserID
			existingRiskRegisterUpdatePlan.ApprovedAt = &approvedAt
			existingRiskRegisterUpdatePlan.NextReviewLevel = nil
		} else {
			existingRiskRegisterUpdatePlan.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := rruprs.repo.UpdateRiskRegisterUpdatePlan(existingRiskRegisterUpdatePlan, tx); err != nil {
			tx.Rollback()
			return err
		}

		// If fully approved, update percentage to include final approval step
		if nextWorkflowLevel == nil {
			if err := rruprs.riskRegisterUpdatePlanSvc.CalculateAndUpdatePercentage(riskRegisterUpdatePlanID, tx); err != nil {
				tx.Rollback()
				return err
			}

			// Update plan status to completed
			completedStatus, err := rruprs.planStatusSvc.GetByType("completed")
			if err != nil {
				tx.Rollback()
				return err
			}

			if err := rruprs.riskRegisterUpdatePlanSvc.UpdatePlanStatus(riskRegisterUpdatePlanID, uint(completedStatus.ID), tx); err != nil {
				tx.Rollback()
				return err
			}

			// Get all risks under the plan and generate reference numbers
			risks, err := rruprs.riskRegisterUpdatePlanSvc.GetRisksByPlanID(riskRegisterUpdatePlanID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, risk := range risks {
				referenceNumber, err := rruprs.riskRegisterSvc.GenerateReferenceNumber()
				if err != nil {
					tx.Rollback()
					return err
				}

				risk.ReferenceNumber = &referenceNumber
				if err := rruprs.riskRegisterSvc.UpdateRisk(risk, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}

			planDetails, err := rruprs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
			if err != nil {
				return err
			}
			rruprs.notifyNextLevelUsers(riskRegisterUpdatePlanID, nextWorkflow, planDetails, currentUserID)
		}

		planDetails, err := rruprs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
		if err != nil {
			return err
		}
		rruprs.sendApprovalNotification(riskRegisterUpdatePlanID, currentUserID, planDetails)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.RiskRegisterUpdatePlanWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		// Create approval review in both current and history tables
		if err := rruprs.createReviewWithHistory(currentWorkflow.ID, riskRegisterUpdatePlanID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Update participant status
		for _, participant := range currentWorkflow.RiskRegisterUpdatePlanWorkflowParticipants {
			if participant.UserID == currentUserID {
				participant.WorkflowStatusID = workflowStatus.ID
				if err := rruprs.repo.UpdateRiskRegisterUpdatePlanWorkflowParticipant(participant, tx); err != nil {
					tx.Rollback()
					return err
				}
				break
			}
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.RiskRegisterUpdatePlanWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		riskRegisterUpdatePlanWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		riskRegisterUpdatePlanWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := rruprs.riskRegisterUpdatePlanWorkflowSvc.UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetNextLevel(riskRegisterUpdatePlanID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingRiskRegisterUpdatePlan, err := rruprs.riskRegisterUpdatePlanSvc.GetByPlanID(riskRegisterUpdatePlanID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingRiskRegisterUpdatePlan.WorkflowStatusID = &workflowStatus.ID
			existingRiskRegisterUpdatePlan.ApprovedBYID = &currentUserID
			existingRiskRegisterUpdatePlan.ApprovedAt = &approvedAt
			existingRiskRegisterUpdatePlan.NextReviewLevel = nil
		} else {
			existingRiskRegisterUpdatePlan.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := rruprs.repo.UpdateRiskRegisterUpdatePlan(existingRiskRegisterUpdatePlan, tx); err != nil {
			tx.Rollback()
			return err
		}

		// If fully approved, update percentage to include final approval step
		if nextWorkflowLevel == nil {
			if err := rruprs.riskRegisterUpdatePlanSvc.CalculateAndUpdatePercentage(riskRegisterUpdatePlanID, tx); err != nil {
				tx.Rollback()
				return err
			}

			// Update plan status to completed
			completedStatus, err := rruprs.planStatusSvc.GetByType("completed")
			if err != nil {
				tx.Rollback()
				return err
			}

			if err := rruprs.riskRegisterUpdatePlanSvc.UpdatePlanStatus(riskRegisterUpdatePlanID, uint(completedStatus.ID), tx); err != nil {
				tx.Rollback()
				return err
			}

			// Get all risks under the plan and generate reference numbers
			risks, err := rruprs.riskRegisterUpdatePlanSvc.GetRisksByPlanID(riskRegisterUpdatePlanID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, risk := range risks {
				referenceNumber, err := rruprs.riskRegisterSvc.GenerateReferenceNumber()
				if err != nil {
					tx.Rollback()
					return err
				}

				risk.ReferenceNumber = &referenceNumber
				if err := rruprs.riskRegisterSvc.UpdateRisk(risk, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}
			planDetails, err := rruprs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
			if err != nil {
				return err
			}
			rruprs.notifyNextLevelUsers(riskRegisterUpdatePlanID, nextWorkflow, planDetails, currentUserID)
		}

		planDetails, err := rruprs.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
		if err != nil {
			return err
		}
		rruprs.sendApprovalNotification(riskRegisterUpdatePlanID, currentUserID, planDetails)
	}

	return nil
}

func (rruprs *RiskRegisterUpdatePlanReviewService) Reject(riskRegisterUpdatePlanID uint, currentUserID uint) error {
	workflowStatus, err := rruprs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	tx := rruprs.repo.BeginTransaction()

	riskRegisterUpdatePlan, err := rruprs.riskRegisterUpdatePlanSvc.GetByPlanID(riskRegisterUpdatePlanID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if riskRegisterUpdatePlan.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := riskRegisterUpdatePlan.NextReviewLevel
	currentWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create revert review in both current and history tables
	if err := rruprs.createReviewWithHistory(currentWorkflow.ID, riskRegisterUpdatePlanID, currentWorkflow.Level, workflowStatus.ID, currentUserID, tx); err != nil {
		tx.Rollback()
		return err
	}

	riskRegisterUpdatePlanWorkflow, err := rruprs.riskRegisterUpdatePlanWorkflowSvc.GetRiskRegisterUpdatePlanWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskRegisterUpdatePlanWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := rruprs.riskRegisterUpdatePlanWorkflowSvc.UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRiskRegisterUpdatePlan, err := rruprs.riskRegisterUpdatePlanSvc.GetByPlanID(riskRegisterUpdatePlanID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRiskRegisterUpdatePlan.WorkflowStatusID = &workflowStatus.ID
	existingRiskRegisterUpdatePlan.NextReviewLevel = nil
	existingRiskRegisterUpdatePlan.SubmittedAt = nil
	existingRiskRegisterUpdatePlan.SubmittedBYID = nil
	existingRiskRegisterUpdatePlan.UpdatedBYID = &currentUserID
	if err := rruprs.repo.UpdateRiskRegisterUpdatePlan(existingRiskRegisterUpdatePlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	rruprs.sendRejectedNotification(riskRegisterUpdatePlanID, currentUserID, riskRegisterUpdatePlan)

	return tx.Commit()
}

func (rruprs *RiskRegisterUpdatePlanReviewService) sendApprovalNotification(id uint, currentUserID uint, planDetail *models.RiskRegisterUpdatePlan) {
	go func() {

		urlParam := fmt.Sprintf(`{"id":%d}`, planDetail.ID)
		entity := "risk-register-update-plan-review"
		action := "approval"

		planTitle := fmt.Sprintf("Update Plan: %s to %s", planDetail.StartDate.Format("2006-01-02"), planDetail.EndDate.Format("2006-01-02"))

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{planDetail.CreatedBYID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Risk Register Update Plan Approval",
			HeadingAr:            "Risk Register Update Plan Approval",
			TitleEn:              planTitle,
			TitleAr:              planTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		createdByUser, err := rruprs.userSvc.GetByID(planDetail.CreatedBYID)
		if err != nil {
			utility.LogError("Failed to fetch plan creator", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		approvedByUser, err := rruprs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 planTitle,
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%srisk/risk-register-update-plan/%d/overview", appURL, planDetail.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{createdByUser.Email},
			Subject:              "Risk Register Update Plan Approval Notification",
			TemplateRelativePath: "risk/templates/email/risk_register_update_plan_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = rruprs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_register_update_plan")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (rruprs *RiskRegisterUpdatePlanReviewService) notifyNextLevelUsers(id uint, nextWorkflowLevel *models.RiskRegisterUpdatePlanWorkflow, planDetail *models.RiskRegisterUpdatePlan, currentUserID uint) {
	go func() {

		urlParam := fmt.Sprintf(`{"id":%d}`, planDetail.ID)
		entity := "risk-register-update-plan-review"
		action := "review-user-assign"

		planTitle := fmt.Sprintf("Update Plan: %s to %s", planDetail.StartDate.Format("2006-01-02"), planDetail.EndDate.Format("2006-01-02"))

		for _, participant := range nextWorkflowLevel.RiskRegisterUpdatePlanWorkflowParticipants {
			responsibleUser, err := rruprs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := planDetail.CreatedBYID
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Risk Register Update Plan Review",
				HeadingAr:            "مراجعة خطة تحديث سجل المخاطر",
				TitleEn:              planTitle,
				TitleAr:              planTitle,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			assignedBy, err := rruprs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByUser, err := rruprs.userSvc.GetByID(planDetail.CreatedBYID)
			if err != nil {
				utility.LogError("Failed to fetch plan creator", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 planTitle,
				"ResponsibleUserName":   createdByUser.FirstNameEn + " " + createdByUser.LastNameEn,
				"ResponsibleUserEmail":  createdByUser.Email,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/risk-register-update-plan/%d/overview", appURL, planDetail.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              "Risk Register Update Plan Review Assigned",
				TemplateRelativePath: "risk/templates/email/risk_register_update_plan_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = rruprs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_register_update_plan")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (rruprs *RiskRegisterUpdatePlanReviewService) sendRejectedNotification(id uint, currentUserID uint, riskRegisterUpdatePlan *models.RiskRegisterUpdatePlan) {
	go func() {
		createdByUser, err := rruprs.userSvc.GetByID(riskRegisterUpdatePlan.CreatedBYID)
		if err != nil {
			utility.LogError("Failed to fetch plan creator", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, riskRegisterUpdatePlan.ID)
		entity := "risk-register-update-plan-review"
		action := "rejected"

		planTitle := fmt.Sprintf("Update Plan: %s to %s", riskRegisterUpdatePlan.StartDate.Format("2006-01-02"), riskRegisterUpdatePlan.EndDate.Format("2006-01-02"))

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{riskRegisterUpdatePlan.CreatedBYID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Risk Register Update Plan Reverted",
			HeadingAr:            "تم رفض خطة تحديث سجل المخاطر",
			TitleEn:              planTitle,
			TitleAr:              planTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := rruprs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"Title":                 planTitle,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%srisk/risk-register-update-plan/%d/overview", appURL, riskRegisterUpdatePlan.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{createdByUser.Email},
			Subject:              "Risk Register Update Plan: Reverted Notification",
			TemplateRelativePath: "risk/templates/email/risk_register_update_plan_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		err = rruprs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_register_update_plan")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

// Helper function to create review in both current and history tables
func (rruprs *RiskRegisterUpdatePlanReviewService) createReviewWithHistory(workflowID uint, planID uint, level int, statusID uint, userID uint, tx utility.DatabaseTransaction) error {
	// Create in current reviews table
	planReview := &models.RiskRegisterUpdatePlanWorkflowReview{
		RiskRegisterUpdatePlanWorkflowID: workflowID,
		WorkflowStatusID:                 statusID,
		CreatedBYID:                      userID,
	}
	if err := rruprs.repo.CreateRiskRegisterUpdatePlanReview(planReview, tx); err != nil {
		return err
	}

	// Create in history table for audit trail (no FK to workflow to avoid deletion conflicts)
	historyReview := &models.RiskRegisterUpdatePlanWorkflowHistory{
		RiskRegisterUpdatePlanID: planID,
		Level:                    level,
		WorkflowStatusID:         statusID,
		CreatedBYID:              userID,
		CreatedAt:                time.Now(),
	}
	if err := rruprs.repo.CreateRiskRegisterUpdatePlanWorkflowHistory(historyReview, tx); err != nil {
		return err
	}

	return nil
}
