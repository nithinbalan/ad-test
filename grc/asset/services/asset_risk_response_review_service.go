package services

import (
	"fmt"
	"time"

	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
)

type AssetRiskResponseReviewServiceInterface interface {
	Submit(riskID uint, currentUserID uint) error
	Approve(riskID uint, currentUserID uint) error
	Revert(riskID uint, currentUserID uint) error
}

type AssetRiskResponseReviewService struct {
	repo                    repositories.AssetRiskResponseReviewRepositoryInterface
	riskSvc                 AssetRiskRegisterServiceInterface
	riskResponseSvc         AssetRiskResponseServiceInterface
	riskResponseWorkflowSvc AssetRiskResponseWorkflowServiceInterface
	workflowStatusSvc       masterServices.WorkflowStatusServiceInterface
	userSvc                 organizationSvc.UserServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
}

func NewAssetRiskResponseReviewService(
	repo repositories.AssetRiskResponseReviewRepositoryInterface,
	riskSvc AssetRiskRegisterServiceInterface,
	riskResponseSvc AssetRiskResponseServiceInterface,
	riskResponseWorkflowSvc AssetRiskResponseWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AssetRiskResponseReviewServiceInterface {
	return &AssetRiskResponseReviewService{
		repo:                    repo,
		riskSvc:                 riskSvc,
		riskResponseSvc:         riskResponseSvc,
		riskResponseWorkflowSvc: riskResponseWorkflowSvc,
		workflowStatusSvc:       workflowStatusSvc,
		userSvc:                 userSvc,
		notificationSvc:         notificationSvc,
	}
}

func (arrs *AssetRiskResponseReviewService) Submit(riskID uint, currentUserID uint) (err error) {
	tx := arrs.repo.BeginTransaction()

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

	if err := arrs.riskResponseWorkflowSvc.ValidateWorkflowBeforeCreate("asset_risk_response", tx); err != nil {
		return err
	}

	riskResponse, err := arrs.riskResponseSvc.GetByRiskRegisterID(riskID, tx)
	if err != nil {
		return err
	}

	if riskResponse.NextReviewLevel != nil {
		return utility.NewValidationError(fmt.Errorf("risk response has already been submitted"), "workflow_already_submitted", "risk_response")
	}

	if riskResponse.ApprovedByID != nil {
		return utility.NewValidationError(fmt.Errorf("risk response workflow has already been completed"), "workflow_already_completed", "risk_response")
	}

	if err := arrs.riskResponseWorkflowSvc.DeleteAllWorkflowsForRiskResponse(riskID, tx); err != nil {
		return err
	}

	if err := arrs.riskResponseWorkflowSvc.CreateRiskResponseWorkflow(riskResponse.ID, "asset_risk_response", currentUserID, tx); err != nil {
		return err
	}

	submitStatus, err := arrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	level := 1
	now := time.Now()

	riskResponse.WorkflowStatusID = &submitStatus.ID
	riskResponse.SubmittedByID = &currentUserID
	riskResponse.SubmittedAt = &now
	riskResponse.NextReviewLevel = &level
	riskResponse.UpdatedBYID = &currentUserID

	if err := arrs.repo.UpdateRiskResponse(riskResponse, tx); err != nil {
		return err
	}

	riskResponseWorkflowLevel, err := arrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponse.ID, level, tx)
	if err != nil {
		return err
	}

	if err := arrs.createReviewWithHistory(riskResponseWorkflowLevel.ID, riskResponse.ID, level, submitStatus.ID, currentUserID, tx); err != nil {
		return err
	}

	notify = func() {
		arrs.notifyWorkflowParticipants(riskID, riskResponseWorkflowLevel, currentUserID)
	}

	return err
}

func (arrs *AssetRiskResponseReviewService) Approve(riskID uint, currentUserID uint) (err error) {
	tx := arrs.repo.BeginTransaction()

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

	riskResponse, err := arrs.riskResponseSvc.GetByRiskRegisterID(riskID, tx)
	if err != nil {
		return err
	}

	if riskResponse.NextReviewLevel == nil {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "risk_response")
	}
	currentApprovalLevel := *riskResponse.NextReviewLevel

	currentWorkflow, err := arrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponse.ID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	approveStatus, err := arrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	pendingStatus, err := arrs.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if len(currentWorkflow.RiskResponseWorkflowParticipants) == 0 || currentUserID != currentWorkflow.RiskResponseWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := arrs.createReviewWithHistory(currentWorkflow.ID, riskResponse.ID, currentApprovalLevel, approveStatus.ID, currentUserID, tx); err != nil {
			return err
		}

		currentWorkflow.WorkflowStatusID = approveStatus.ID
		if err := arrs.riskResponseWorkflowSvc.UpdateRiskResponseWorkflow(currentWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := arrs.riskResponseWorkflowSvc.GetNextLevel(riskResponse.ID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			riskResponse.WorkflowStatusID = &approveStatus.ID
			riskResponse.ApprovedByID = &currentUserID
			riskResponse.ApprovedAt = &approvedAt
			riskResponse.NextReviewLevel = nil
		} else {
			riskResponse.WorkflowStatusID = &pendingStatus.ID
			riskResponse.NextReviewLevel = &nextWorkflowLevel.Level
		}
		riskResponse.UpdatedBYID = &currentUserID

		if err := arrs.repo.UpdateRiskResponse(riskResponse, tx); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := arrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponse.ID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}
			notify = func() {
				arrs.notifyWorkflowParticipants(riskID, nextWorkflow, currentUserID)
			}
		} else {
			notify = func() {
				arrs.sendApprovalNotification(riskID, currentUserID, riskResponse.CreatedBYID)
			}
		}

	case "multiple-users", "role":
		var approvingParticipant *models.RiskResponseWorkflowParticipant
		for _, participant := range currentWorkflow.RiskResponseWorkflowParticipants {
			if participant.UserID == currentUserID {
				approvingParticipant = participant
				break
			}
		}
		if approvingParticipant == nil {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if err := arrs.createReviewWithHistory(currentWorkflow.ID, riskResponse.ID, currentApprovalLevel, approveStatus.ID, currentUserID, tx); err != nil {
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := arrs.riskResponseWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.RiskResponseWorkflowParticipants, tx)
			if err != nil || !allApproved {
				return err
			}
		}

		currentWorkflow.WorkflowStatusID = approveStatus.ID
		if err := arrs.riskResponseWorkflowSvc.UpdateRiskResponseWorkflow(currentWorkflow, tx); err != nil {
			return err
		}

		nextWorkflowLevel, err := arrs.riskResponseWorkflowSvc.GetNextLevel(riskResponse.ID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextWorkflowLevel == nil {
			approvedAt := time.Now()
			riskResponse.WorkflowStatusID = &approveStatus.ID
			riskResponse.ApprovedByID = &currentUserID
			riskResponse.ApprovedAt = &approvedAt
			riskResponse.NextReviewLevel = nil
		} else {
			riskResponse.WorkflowStatusID = &pendingStatus.ID
			riskResponse.NextReviewLevel = &nextWorkflowLevel.Level
		}
		riskResponse.UpdatedBYID = &currentUserID

		if err := arrs.repo.UpdateRiskResponse(riskResponse, tx); err != nil {
			return err
		}

		if nextWorkflowLevel != nil {
			nextWorkflow, err := arrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponse.ID, nextWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}
			notify = func() {
				arrs.notifyWorkflowParticipants(riskID, nextWorkflow, currentUserID)
			}
		} else {
			notify = func() {
				arrs.sendApprovalNotification(riskID, currentUserID, riskResponse.CreatedBYID)
			}
		}
	}

	return err
}

func (arrs *AssetRiskResponseReviewService) Revert(riskID uint, currentUserID uint) (err error) {
	tx := arrs.repo.BeginTransaction()

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

	riskResponse, err := arrs.riskResponseSvc.GetByRiskRegisterID(riskID, tx)
	if err != nil {
		return err
	}

	if riskResponse.NextReviewLevel == nil {
		return utility.NewValidationError(fmt.Errorf("workflow not active"), "workflow_not_active", "risk_response")
	}
	currentLevel := *riskResponse.NextReviewLevel

	currentWorkflow, err := arrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponse.ID, currentLevel, tx)
	if err != nil {
		return err
	}

	allowed := false
	for _, participant := range currentWorkflow.RiskResponseWorkflowParticipants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}
	if !allowed {
		return utility.NewForbiddenError("you have no permission to revert")
	}

	revertStatus, err := arrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	if err := arrs.createReviewWithHistory(currentWorkflow.ID, riskResponse.ID, currentLevel, revertStatus.ID, currentUserID, tx); err != nil {
		return err
	}

	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.RiskResponseWorkflow
	if previousLevel > 0 {
		pendingStatus, err := arrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		previousLevelWorkflow, err = arrs.riskResponseWorkflowSvc.GetRiskResponseWorkflowByLevel(riskResponse.ID, previousLevel, tx)
		if err != nil {
			return err
		}

		if err := arrs.riskResponseWorkflowSvc.DeleteRiskResponseWorkflowReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := arrs.riskResponseWorkflowSvc.UpdateRiskResponseWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			return err
		}
	}

	var submittedByID *uint
	if previousLevel == 0 {
		submittedByID = riskResponse.SubmittedByID

		riskResponse.WorkflowStatusID = &revertStatus.ID
		riskResponse.NextReviewLevel = nil
		riskResponse.SubmittedByID = nil
		riskResponse.SubmittedAt = nil
	} else {
		riskResponse.NextReviewLevel = &previousLevel
	}
	riskResponse.UpdatedBYID = &currentUserID

	if err := arrs.repo.UpdateRiskResponse(riskResponse, tx); err != nil {
		return err
	}

	if previousLevel > 0 {
		notify = func() {
			arrs.notifyRevertedToPreviousLevel(riskID, previousLevelWorkflow, currentUserID)
		}
	} else {
		notify = func() {
			arrs.sendRevertNotification(riskID, currentUserID, submittedByID, riskResponse.CreatedBYID)
		}
	}

	return err
}

func (arrs *AssetRiskResponseReviewService) notifyWorkflowParticipants(riskID uint, workflowLevel *models.RiskResponseWorkflow, currentUserID uint) {
	assignedBy, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch assigning user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID)
	if err != nil {
		utility.LogError("Failed to fetch risk register details", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskID)
	entity := "asset-risk-response-review"
	action := "review-user-assign"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/risks/%d/risk-response", appURL, riskID)

	for _, participant := range workflowLevel.RiskResponseWorkflowParticipants {
		responsibleUser, err := arrs.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.ReviewType),
			ReferenceNumber:      riskDetail.ReferenceNumber,
			HeadingEn:            "ISMS Risk Response Workflow",
			HeadingAr:            "سير عمل معالجة مخاطر ISMS",
			TitleEn:              riskDetail.Title,
			TitleAr:              riskDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"Title":                 riskDetail.Title,
			"ReferenceNumber":       riskDetail.ReferenceNumber,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("ISMS Risk Response Workflow: %s - Review Assigned", riskDetail.Title),
			TemplateRelativePath: "asset/templates/email/asset_risk_response_workflow_review_user_mail.html",
			EmailData:            emailData,
		}

		if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func (arrs *AssetRiskResponseReviewService) notifyRevertedToPreviousLevel(riskID uint, workflowLevel *models.RiskResponseWorkflow, currentUserID uint) {
	revertedByUser, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID)
	if err != nil {
		utility.LogError("Failed to fetch risk register details", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskID)
	entity := "asset-risk-response-review"
	action := "reverted"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/risks/%d/risk-response", appURL, riskID)

	for _, participant := range workflowLevel.RiskResponseWorkflowParticipants {
		responsibleUser, err := arrs.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.RevertType),
			ReferenceNumber:      riskDetail.ReferenceNumber,
			HeadingEn:            "ISMS Risk Response Workflow",
			HeadingAr:            "سير عمل معالجة مخاطر ISMS",
			TitleEn:              riskDetail.Title,
			TitleAr:              riskDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"Title":                 riskDetail.Title,
			"ReferenceNumber":       riskDetail.ReferenceNumber,
			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
			"RevertedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("ISMS Risk Response Workflow: %s - Reverted", riskDetail.Title),
			TemplateRelativePath: "asset/templates/email/asset_risk_response_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func (arrs *AssetRiskResponseReviewService) sendApprovalNotification(riskID uint, currentUserID uint, createdByID *uint) {
	if createdByID == nil {
		return
	}

	approvedByUser, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch approving user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	createdByUser, err := arrs.userSvc.GetByID(*createdByID)
	if err != nil {
		utility.LogError("Failed to fetch risk response creator", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID)
	if err != nil {
		utility.LogError("Failed to fetch risk register details", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskID)
	entity := "asset-risk-response-review"
	action := "approval"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/risks/%d/risk-response", appURL, riskID)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*createdByID},
		ModuleIdentifier:     "asset_risk",
		NotificationType:     string(commonDTOs.ApprovalType),
		HeadingEn:            "ISMS Risk Response Workflow",
		HeadingAr:            "سير عمل معالجة مخاطر ISMS",
		ReferenceNumber:      riskDetail.ReferenceNumber,
		TitleEn:              riskDetail.Title,
		TitleAr:              riskDetail.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           createdByUser.FirstNameEn,
		"LastNameEn":            createdByUser.LastNameEn,
		"Title":                 riskDetail.Title,
		"ReferenceNumber":       riskDetail.ReferenceNumber,
		"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
		"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
		"ApprovedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{createdByUser.Email},
		Subject:              fmt.Sprintf("ISMS Risk Response Workflow: %s - Approved", riskDetail.Title),
		TemplateRelativePath: "asset/templates/email/asset_risk_response_workflow_approve_user_mail.html",
		EmailData:            emailData,
	}

	if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); err != nil {
		utility.LogError("Failed to send approval notification email", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (arrs *AssetRiskResponseReviewService) sendRevertNotification(riskID uint, currentUserID uint, submittedByID *uint, createdByID *uint) {
	notifyUserID := submittedByID
	if notifyUserID == nil {
		notifyUserID = createdByID
	}
	if notifyUserID == nil {
		return
	}

	revertedByUser, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	notifyUser, err := arrs.userSvc.GetByID(*notifyUserID)
	if err != nil {
		utility.LogError("Failed to fetch submitted user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID)
	if err != nil {
		utility.LogError("Failed to fetch risk register details", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskID)
	entity := "asset-risk-response-review"
	action := "reverted"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/risks/%d/risk-response", appURL, riskID)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{*notifyUserID},
		ModuleIdentifier:     "asset_risk",
		NotificationType:     string(commonDTOs.RevertType),
		HeadingEn:            "ISMS Risk Response Workflow",
		HeadingAr:            "سير عمل معالجة مخاطر ISMS",
		ReferenceNumber:      riskDetail.ReferenceNumber,
		TitleEn:              riskDetail.Title,
		TitleAr:              riskDetail.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	emailData := map[string]interface{}{
		"FirstNameEn":           notifyUser.FirstNameEn,
		"LastNameEn":            notifyUser.LastNameEn,
		"Title":                 riskDetail.Title,
		"ReferenceNumber":       riskDetail.ReferenceNumber,
		"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
		"RevertedByLastNameEn":  revertedByUser.LastNameEn,
		"RevertedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{notifyUser.Email},
		Subject:              fmt.Sprintf("ISMS Risk Response Workflow: %s - Reverted", riskDetail.Title),
		TemplateRelativePath: "asset/templates/email/asset_risk_response_workflow_revert_user_mail.html",
		EmailData:            emailData,
	}

	if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); err != nil {
		utility.LogError("Failed to send reverted notification email", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// Helper function to create review in both current and history tables
func (arrs *AssetRiskResponseReviewService) createReviewWithHistory(workflowID uint, riskResponseID uint, level int, statusID uint, userID uint, tx utility.DatabaseTransaction) error {
	riskResponseReview := &models.RiskResponseWorkflowReview{
		RiskResponseWorkflowID: workflowID,
		WorkflowStatusID:       statusID,
		CreatedBYID:            userID,
	}
	if err := arrs.repo.CreateRiskResponseReview(riskResponseReview, tx); err != nil {
		return err
	}

	historyReview := &models.RiskResponseWorkflowHistory{
		RiskResponseID:   riskResponseID,
		Level:            level,
		WorkflowStatusID: statusID,
		CreatedBYID:      userID,
		CreatedAt:        time.Now(),
	}
	if err := arrs.repo.CreateRiskResponseWorkflowHistory(historyReview, tx); err != nil {
		return err
	}

	return nil
}
