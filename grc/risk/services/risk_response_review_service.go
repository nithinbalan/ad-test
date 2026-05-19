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

type RiskResponseReviewServiceInterface interface {
	Submit(riskResponseID uint, currentUserID uint) error
	Approve(riskResponseID uint, currentUserID uint) error
	Reject(riskResponseID uint, currentUserID uint) error
	Reopen(riskResponseID uint, currentUserID uint) error
}

type RiskResponseReviewService struct {
	repo                    repositories.RiskResponseReviewRepositoryInterface
	logService              logServices.ActivityLogServiceInterface
	riskResponseSvc         RiskResponseServiceInterface
	riskResponseWorkflowSvc RiskResponseWorkflowServiceInterface
	workflowStatusSvc       masterServices.WorkflowStatusServiceInterface
	userSvc                 organizationSvc.UserServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
	taskSvc                 taskServices.TaskServiceInterface
}

func NewRiskResponseReviewService(
	repo repositories.RiskResponseReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	riskResponseSvc RiskResponseServiceInterface,
	riskResponseWorkflowSvc RiskResponseWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) RiskResponseReviewServiceInterface {
	return &RiskResponseReviewService{
		repo:                    repo,
		logService:              logService,
		riskResponseSvc:         riskResponseSvc,
		riskResponseWorkflowSvc: riskResponseWorkflowSvc,
		workflowStatusSvc:       workflowStatusSvc,
		userSvc:                 userSvc,
		notificationSvc:         notificationSvc,
		taskSvc:                 taskSvc,
	}
}

func (rrs *RiskResponseReviewService) Submit(riskResponseID uint, currentUserID uint) error {
	// Get full details for validation and notifications
	riskResponseDetail, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
	if err != nil {
		return err
	}

	//Validate residual risk rating - must be high or critical
	if riskResponseDetail.RiskRegister != nil && riskResponseDetail.RiskRegister.ResidualRiskRating != nil {
		ratingType := riskResponseDetail.RiskRegister.ResidualRiskRating.Type
		if ratingType != "high" && ratingType != "critical" {
			return utility.NewValidationError(fmt.Errorf("residual risk rating must be 'high' or 'critical' to submit risk response for review"), "residual_risk_rating", "residual_risk_rating")
		}

		// If residual risk is critical or high, check if at least one action plan exists
		if len(riskResponseDetail.RiskRegister.ActionPlans) == 0 {
			return utility.NewValidationError(fmt.Errorf("at least one action plan is required when residual risk rating is 'high' or 'critical'"), "residual_risk_rating", "residual_risk_rating")
		}
	}

	// Get the model for database operations
	riskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID)
	if err != nil {
		return err
	}

	if riskResponse.SubmittedByID != nil && riskResponse.WorkflowStatusID != nil {
		workflowStatus, err := rrs.workflowStatusSvc.GetByID(*riskResponse.WorkflowStatusID)
		if err != nil {
			return err
		}

		if workflowStatus.Type != "revert" {
			return utility.NewForbiddenError("risk_response_already_submitted")
		}
	}

	err = rrs.riskResponseWorkflowSvc.ValidateWorkflowBeforeCreate("risk_response")
	if err != nil {
		return err
	}

	tx := rrs.repo.BeginTransaction()
	currentDate := time.Now()
	level := 1

	workflowStatus, err := rrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = rrs.riskResponseWorkflowSvc.DeleteAllWorkflowsForRiskResponse(riskResponseID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = rrs.riskResponseWorkflowSvc.CreateRiskResponseWorkflow(riskResponseID, "risk_response", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskResponse.WorkflowStatusID = &workflowStatus.ID
	riskResponse.SubmittedByID = &currentUserID
	riskResponse.SubmittedAt = &currentDate
	riskResponse.NextReviewLevel = &level
	riskResponse.UpdatedBYID = &currentUserID

	err = rrs.repo.UpdateRiskResponse(riskResponse, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskResponseWorkflowLevel, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponseID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskResponseReview := &models.RiskResponseWorkflowReview{
		RiskResponseWorkflowID: riskResponseWorkflowLevel.ID,
		WorkflowStatusID:       workflowStatus.ID,
		CreatedBYID:            currentUserID,
	}

	err = rrs.repo.CreateRiskResponseReview(riskResponseReview, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Send notifications asynchronously
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, riskResponseDetail.RiskRegister.ID)
		entity := "risk-response-review"
		action := "review-user-assign"

		// Get risk details from DTO for notifications
		referenceNumber := riskResponseDetail.RiskRegister.ReferenceNumber
		riskTitle := riskResponseDetail.RiskRegister.Title

		for _, participant := range riskResponseWorkflowLevel.RiskResponseWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      referenceNumber,
				HeadingEn:            "Risk Response Review",
				HeadingAr:            "مراجعة استجابة المخاطر",
				TitleEn:              riskTitle,
				TitleAr:              riskTitle,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := rrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := rrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := "Risk Response: Review Assigned"
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 riskTitle,
				"ReferenceNumber":       referenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskResponseDetail.RiskRegister.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_response_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_response")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (rrs *RiskResponseReviewService) Approve(riskResponseID uint, currentUserID uint) error {
	tx := rrs.repo.BeginTransaction()

	riskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if riskResponse.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := riskResponse.NextReviewLevel
	currentWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponseID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := rrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentWorkflow.RiskResponseWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		riskResponseReview := &models.RiskResponseWorkflowReview{
			RiskResponseWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:       workflowStatus.ID,
			CreatedBYID:            currentUserID,
		}
		if err := rrs.repo.CreateRiskResponseReview(riskResponseReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		riskResponseWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		riskResponseWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := rrs.riskResponseWorkflowSvc.UpdateRiskResponseWorkflow(riskResponseWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := rrs.riskResponseWorkflowSvc.GetNextLevel(riskResponseID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingRiskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingRiskResponse.WorkflowStatusID = &workflowStatus.ID
			existingRiskResponse.ApprovedByID = &currentUserID
			existingRiskResponse.ApprovedAt = &approvedAt
			existingRiskResponse.NextReviewLevel = nil
		} else {
			existingRiskResponse.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := rrs.repo.UpdateRiskResponse(existingRiskResponse, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction before sending notifications
		if err := tx.Commit(); err != nil {
			return err
		}

		// Send notifications after transaction commit
		if nextWorkflowLevel != nil {
			nextWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponseID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}

			riskDetails, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
			if err != nil {
				return err
			}
			rrs.notifyNextLevelUsers(riskResponseID, nextWorkflow, riskDetails, currentUserID)
		}

		riskDetails, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
		if err != nil {
			return err
		}
		rrs.sendApprovalNotification(riskResponseID, currentUserID, riskDetails)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentWorkflow.RiskResponseWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		riskResponseReview := &models.RiskResponseWorkflowReview{
			RiskResponseWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:       workflowStatus.ID,
			CreatedBYID:            currentUserID,
		}
		if err := rrs.repo.CreateRiskResponseReview(riskResponseReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := rrs.riskResponseWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.RiskResponseWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			if !allApproved {
				return tx.Commit()
			}
		}

		riskResponseWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		riskResponseWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := rrs.riskResponseWorkflowSvc.UpdateRiskResponseWorkflow(riskResponseWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextWorkflowLevel, err := rrs.riskResponseWorkflowSvc.GetNextLevel(riskResponseID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingRiskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			existingRiskResponse.WorkflowStatusID = &workflowStatus.ID
			existingRiskResponse.ApprovedByID = &currentUserID
			existingRiskResponse.ApprovedAt = &approvedAt
			existingRiskResponse.NextReviewLevel = nil
		} else {
			existingRiskResponse.NextReviewLevel = &nextWorkflowLevel.Level
		}

		if err := rrs.repo.UpdateRiskResponse(existingRiskResponse, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction before sending notifications
		if err := tx.Commit(); err != nil {
			return err
		}

		// Send notifications after transaction commit
		if nextWorkflowLevel != nil {
			nextWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponseID, nextWorkflowLevel.Level)
			if err != nil {
				return err
			}
			riskDetails, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
			if err != nil {
				return err
			}
			rrs.notifyNextLevelUsers(riskResponseID, nextWorkflow, riskDetails, currentUserID)
		}

		riskDetails, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
		if err != nil {
			return err
		}
		rrs.sendApprovalNotification(riskResponseID, currentUserID, riskDetails)
	}

	return nil
}

func (rrs *RiskResponseReviewService) Reject(riskResponseID uint, currentUserID uint) error {
	// Get full details for notifications
	riskResponseDetail, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
	if err != nil {
		return err
	}

	workflowStatus, err := rrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	tx := rrs.repo.BeginTransaction()

	riskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if riskResponse.NextReviewLevel == nil {
		tx.Rollback()
		return utility.NewForbiddenError("review level not found")
	}

	currentApprovalLevel := riskResponse.NextReviewLevel
	currentWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponseID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskResponseWorkflowReview := &models.RiskResponseWorkflowReview{
		RiskResponseWorkflowID: currentWorkflow.ID,
		WorkflowStatusID:       workflowStatus.ID,
		CreatedBYID:            currentUserID,
	}

	if err := rrs.repo.CreateRiskResponseReview(riskResponseWorkflowReview, tx); err != nil {
		tx.Rollback()
		return err
	}

	riskResponseWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	riskResponseWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := rrs.riskResponseWorkflowSvc.UpdateRiskResponseWorkflow(riskResponseWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRiskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRiskResponse.WorkflowStatusID = &workflowStatus.ID
	existingRiskResponse.NextReviewLevel = nil
	existingRiskResponse.SubmittedAt = nil
	existingRiskResponse.SubmittedByID = nil
	existingRiskResponse.UpdatedBYID = &currentUserID
	if err := rrs.repo.UpdateRiskResponse(existingRiskResponse, tx); err != nil {
		tx.Rollback()
		return err
	}

	rrs.sendRejectedNotification(riskResponseID, currentUserID, riskResponseDetail)

	return tx.Commit()
}

func (rrs *RiskResponseReviewService) Reopen(riskResponseID uint, currentUserID uint) error {
	// Get full details for notifications
	riskResponseDetail, err := rrs.riskResponseSvc.GetFullDetailsByID(riskResponseID)
	if err != nil {
		return err
	}

	// Get the model for database operations
	riskResponse, err := rrs.riskResponseSvc.GetByID(riskResponseID)
	if err != nil {
		return err
	}

	if riskResponse.SubmittedByID != nil && riskResponse.WorkflowStatusID != nil {
		workflowStatus, err := rrs.workflowStatusSvc.GetByID(*riskResponse.WorkflowStatusID)
		if err != nil {
			return err
		}

		if workflowStatus.Type != "approve" {
			return utility.NewForbiddenError("risk response already submitted")
		}
	}

	if err := rrs.riskResponseWorkflowSvc.ValidateWorkflowBeforeCreate("risk_response"); err != nil {
		return err
	}

	tx := rrs.repo.BeginTransaction()

	workflowStatus, err := rrs.workflowStatusSvc.GetByType("reopen")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Check if risk response was submitted (has workflow history)
	if riskResponse.ApprovedByID != nil {
		// Get level 1 workflow to attach reopen history entry
		levelOneWorkflow, err := rrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponseID, 1, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create workflow review record for reopen history at level 1
		riskResponseReview := &models.RiskResponseWorkflowReview{
			RiskResponseWorkflowID: levelOneWorkflow.ID,
			WorkflowStatusID:       workflowStatus.ID,
			CreatedBYID:            currentUserID,
		}
		if err := rrs.repo.CreateRiskResponseReview(riskResponseReview, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	riskResponse.WorkflowStatusID = &workflowStatus.ID
	riskResponse.SubmittedByID = nil
	riskResponse.SubmittedAt = nil
	riskResponse.NextReviewLevel = nil
	riskResponse.ApprovedAt = nil
	riskResponse.ApprovedByID = nil
	riskResponse.UpdatedBYID = &currentUserID

	if err := rrs.repo.UpdateRiskResponse(riskResponse, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {

		urlParam := fmt.Sprintf(`{"id":%d}`, riskResponseDetail.RiskRegister.ID)
		entity := "risk-response-review"
		action := "reopened"

		if riskResponse.CreatedBYID == nil {
			return
		}

		createdByUser, err := rrs.userSvc.GetByID(*riskResponse.CreatedBYID)
		if err != nil {
			utility.LogError("Failed to fetch risk response creator", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get risk details from DTO for notifications
		var referenceNumber string
		var riskTitle string
		if riskResponseDetail.RiskRegister != nil {
			referenceNumber = riskResponseDetail.RiskRegister.ReferenceNumber
			riskTitle = riskResponseDetail.RiskRegister.Title
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*riskResponse.CreatedBYID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.RevertType),
			HeadingEn:            "Risk Response Reopened",
			HeadingAr:            "تم إعادة فتح استجابة المخاطر",
			ReferenceNumber:      referenceNumber,
			TitleEn:              riskTitle,
			TitleAr:              riskTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		reopenedByUser, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch risk response reopen user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailTo := []string{createdByUser.Email}
		subject := "Risk Response: Review Reopened"
		redirectURL := fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskResponseDetail.RiskRegister.ID)

		emailData := map[string]interface{}{
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"Title":                 riskTitle,
			"ReferenceNumber":       referenceNumber,
			"ReopenedByFirstNameEn": reopenedByUser.FirstNameEn,
			"ReopenedByLastNameEn":  reopenedByUser.LastNameEn,
			"ReopenedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_response_workflow_reopen_user_mail.html",
			EmailData:            emailData,
		}

		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_response")
		if err != nil {
			utility.LogError("Failed to send risk response reopened notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (rrs *RiskResponseReviewService) sendApprovalNotification(id uint, currentUserID uint, riskDetail *riskDTOs.RiskResponseDetailDTO) {
	go func() {
		if riskDetail.CreatedBy == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.RiskRegister.ID)
		entity := "risk-response-review"
		action := "approval"

		// Get risk details for notifications
		var referenceNumber string
		var riskTitle string
		if riskDetail.RiskRegister != nil {
			referenceNumber = riskDetail.RiskRegister.ReferenceNumber
			riskTitle = riskDetail.RiskRegister.Title
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(riskDetail.CreatedBy.ID)},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Risk Response Approval",
			HeadingAr:            "Risk Response Approval",
			ReferenceNumber:      referenceNumber,
			TitleEn:              riskTitle,
			TitleAr:              riskTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		approvedByUser, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"Title":                 riskTitle,
			"Description":           "",
			"FirstNameEn":           riskDetail.CreatedBy.FirstNameEN,
			"LastNameEn":            riskDetail.CreatedBy.LastNameEN,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskDetail.RiskRegister.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{riskDetail.CreatedBy.Email},
			Subject:              "Risk Response Approval Notification",
			TemplateRelativePath: "risk/templates/email/risk_response_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_response")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (rrs *RiskResponseReviewService) notifyNextLevelUsers(id uint, nextWorkflowLevel *models.RiskResponseWorkflow, riskDetail *riskDTOs.RiskResponseDetailDTO, currentUserID uint) {
	go func() {
		if riskDetail.CreatedBy == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.RiskRegister.ID)
		entity := "risk-response-review"
		action := "review-user-assign"

		// Get risk details for notifications
		var referenceNumber string
		var riskTitle string
		if riskDetail.RiskRegister != nil {
			referenceNumber = riskDetail.RiskRegister.ReferenceNumber
			riskTitle = riskDetail.RiskRegister.Title
		}

		for _, participant := range nextWorkflowLevel.RiskResponseWorkflowParticipants {
			responsibleUser, err := rrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(riskDetail.CreatedBy.ID)
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Risk Response Review",
				HeadingAr:            "مراجعة استجابة المخاطر",
				ReferenceNumber:      referenceNumber,
				TitleEn:              riskTitle,
				TitleAr:              riskTitle,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			assignedBy, err := rrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 riskTitle,
				"Description":           "",
				"ResponsibleUserName":   riskDetail.CreatedBy.FirstNameEN + " " + riskDetail.CreatedBy.LastNameEN,
				"ResponsibleUserEmail":  riskDetail.CreatedBy.Email,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskDetail.RiskRegister.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              "Risk Response Review Assigned",
				TemplateRelativePath: "risk/templates/email/risk_response_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_response")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (rrs *RiskResponseReviewService) sendRejectedNotification(id uint, currentUserID uint, riskResponseDetail *riskDTOs.RiskResponseDetailDTO) {
	go func() {
		// Get the risk response model to get CreatedBYID
		riskResponse, err := rrs.riskResponseSvc.GetByID(id)
		if err != nil {
			return
		}

		if riskResponse.CreatedBYID == nil {
			return
		}

		createdByUser, err := rrs.userSvc.GetByID(*riskResponse.CreatedBYID)
		if err != nil {
			utility.LogError("Failed to fetch risk response creator", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, riskResponseDetail.RiskRegister.ID)
		entity := "risk-response-review"
		action := "rejected"

		// Get risk details from DTO for notifications
		var referenceNumber string
		var riskTitle string
		if riskResponseDetail.RiskRegister != nil {
			referenceNumber = riskResponseDetail.RiskRegister.ReferenceNumber
			riskTitle = riskResponseDetail.RiskRegister.Title
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*riskResponse.CreatedBYID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Risk Response Rejected",
			HeadingAr:            "تم رفض استجابة المخاطر",
			ReferenceNumber:      referenceNumber,
			TitleEn:              riskTitle,
			TitleAr:              riskTitle,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := rrs.userSvc.GetByID(currentUserID)
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
			"Title":                 riskTitle,
			"ReferenceNumber":       referenceNumber,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%srisk/risk-register/%d/risk-response", appURL, riskResponseDetail.RiskRegister.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{createdByUser.Email},
			Subject:              "Risk Response: Rejected Notification",
			TemplateRelativePath: "risk/templates/email/risk_response_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_response")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
