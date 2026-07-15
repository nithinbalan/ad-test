package services

import (
	"fmt"
	"time"

	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
)

type AssetRiskReviewServiceInterface interface {
	Submit(riskID uint, currentUserID uint) error
	Approve(riskID uint, currentUserID uint) error
	Reject(riskID uint, currentUserID uint) error
}

type AssetRiskReviewService struct {
	repo              repositories.AssetRiskReviewRepositoryInterface
	riskSvc           AssetRiskRegisterServiceInterface
	riskWorkflowSvc   AssetRiskWorkflowServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewAssetRiskReviewService(
	repo repositories.AssetRiskReviewRepositoryInterface,
	riskSvc AssetRiskRegisterServiceInterface,
	riskWorkflowSvc AssetRiskWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AssetRiskReviewServiceInterface {
	return &AssetRiskReviewService{
		repo:              repo,
		riskSvc:           riskSvc,
		riskWorkflowSvc:   riskWorkflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
	}
}

func (arrs *AssetRiskReviewService) Submit(riskID uint, currentUserID uint) (err error) {
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

	if err := arrs.riskWorkflowSvc.ValidateWorkflowBeforeCreate("asset_risk", tx); err != nil {
		return err
	}

	riskRegister, err := arrs.riskSvc.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	submittedRiskStatus, err := arrs.riskSvc.GetRiskStatusByType("submitted", tx)
	if err != nil {
		return err
	}

	identifiedRiskStatus, err := arrs.riskSvc.GetRiskStatusByType("identified", tx)
	if err != nil {
		return err
	}

	if riskRegister.NextReviewLevel != nil {
		return utility.NewValidationError(fmt.Errorf("risk has already been submitted"), "already_submitted", "risk")
	}

	if riskRegister.RiskStatusID != identifiedRiskStatus.ID {
		return utility.NewValidationError(fmt.Errorf("workflow has already completed"), "already_completed", "risk")
	}

	if err := arrs.riskWorkflowSvc.CreateRiskWorkflow(riskRegister.ID, "asset_risk", currentUserID, tx); err != nil {
		return err
	}

	level := 1
	now := time.Now()

	riskRegister.RiskStatusID = submittedRiskStatus.ID
	riskRegister.SubmittedBy = &currentUserID
	riskRegister.SubmittedAt = &now
	riskRegister.NextReviewLevel = &level
	riskRegister.UpdatedBYID = &currentUserID

	if err := arrs.repo.UpdateRisk(riskRegister, tx); err != nil {
		return err
	}

	riskWorkflowLevel, err := arrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskRegister.ID, level, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := arrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	riskReview := &models.RiskWorkflowReview{
		RiskWorkflowID:   riskWorkflowLevel.ID,
		WorkflowStatusID: workflowStatus.ID,
		CreatedBYID:      currentUserID,
	}

	if err := arrs.repo.CreateRiskReview(riskReview, tx); err != nil {
		return err
	}

	notify = func() {
		arrs.notifyWorkflowParticipants(riskRegister, riskWorkflowLevel, currentUserID)
	}

	return err
}

func (arrs *AssetRiskReviewService) Approve(riskID uint, currentUserID uint) (err error) {
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

	riskRegister, err := arrs.riskSvc.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	if riskRegister.NextReviewLevel == nil {
		return utility.NewValidationError(fmt.Errorf("risk has not been submitted for approval"), "not_submitted", "risk")
	}
	currentApprovalLevel := *riskRegister.NextReviewLevel

	currentRiskWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := arrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	switch currentRiskWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentRiskWorkflow.RiskWorkflowParticipants[0].UserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		riskReview := &models.RiskWorkflowReview{
			RiskWorkflowID:   currentRiskWorkflow.ID,
			WorkflowStatusID: workflowStatus.ID,
			CreatedBYID:      currentUserID,
		}
		if err = arrs.repo.CreateRiskReview(riskReview, tx); err != nil {
			return err
		}

		riskWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByID(currentRiskWorkflow.ID, tx)
		if err != nil {
			return err
		}
		riskWorkflow.WorkflowStatusID = workflowStatus.ID
		if err = arrs.riskWorkflowSvc.UpdateRiskWorkflow(riskWorkflow, tx); err != nil {
			return err
		}

		nextRiskWorkflowLevel, err := arrs.riskWorkflowSvc.GetNextLevel(riskRegister.ID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextRiskWorkflowLevel == nil {
			openStatus, err := arrs.riskSvc.GetRiskStatusByType("open", tx)
			if err != nil {
				return err
			}
			riskRegister.RiskStatusID = openStatus.ID
			riskRegister.NextReviewLevel = nil
			if err := arrs.repo.UpdateRisk(riskRegister, tx); err != nil {
				return err
			}

			riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID, tx)
			if err != nil {
				return err
			}
			notify = func() {
				arrs.sendApprovalNotification(currentUserID, riskDetail)
				arrs.sendResponsibleUserAssignNotification(currentUserID, riskDetail)
			}
		} else {
			riskRegister.NextReviewLevel = &nextRiskWorkflowLevel.Level
			if err := arrs.repo.UpdateRisk(riskRegister, tx); err != nil {
				return err
			}
			nextWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, nextRiskWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}
			notify = func() {
				arrs.notifyWorkflowParticipants(riskRegister, nextWorkflow, currentUserID)
			}
		}

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentRiskWorkflow.RiskWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		riskReview := &models.RiskWorkflowReview{
			RiskWorkflowID:   currentRiskWorkflow.ID,
			WorkflowStatusID: workflowStatus.ID,
			CreatedBYID:      currentUserID,
		}
		if err = arrs.repo.CreateRiskReview(riskReview, tx); err != nil {
			return err
		}

		if currentRiskWorkflow.ApprovalTypes == "all" {
			var allApproved bool
			allApproved, err = arrs.riskWorkflowSvc.CheckAllParticipantsApproved(currentRiskWorkflow.ID, currentRiskWorkflow.RiskWorkflowParticipants, tx)
			if err != nil || !allApproved {
				return err
			}
		}

		riskWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByID(currentRiskWorkflow.ID, tx)
		if err != nil {
			return err
		}
		riskWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := arrs.riskWorkflowSvc.UpdateRiskWorkflow(riskWorkflow, tx); err != nil {
			return err
		}

		nextRiskWorkflowLevel, err := arrs.riskWorkflowSvc.GetNextLevel(riskRegister.ID, currentApprovalLevel, tx)
		if err != nil {
			return err
		}

		if nextRiskWorkflowLevel == nil {
			openStatus, err := arrs.riskSvc.GetRiskStatusByType("open", tx)
			if err != nil {
				return err
			}
			riskRegister.RiskStatusID = openStatus.ID
			riskRegister.NextReviewLevel = nil
			if err := arrs.repo.UpdateRisk(riskRegister, tx); err != nil {
				return err
			}

			riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID, tx)
			if err != nil {
				return err
			}
			notify = func() {
				arrs.sendApprovalNotification(currentUserID, riskDetail)
				arrs.sendResponsibleUserAssignNotification(currentUserID, riskDetail)
			}
		} else {
			riskRegister.NextReviewLevel = &nextRiskWorkflowLevel.Level
			if err := arrs.repo.UpdateRisk(riskRegister, tx); err != nil {
				return err
			}

			nextWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, nextRiskWorkflowLevel.Level, tx)
			if err != nil {
				return err
			}
			notify = func() {
				arrs.notifyWorkflowParticipants(riskRegister, nextWorkflow, currentUserID)
			}
		}
	}

	return err
}

func (arrs *AssetRiskReviewService) Reject(riskID uint, currentUserID uint) (err error) {
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

	riskRegister, err := arrs.riskSvc.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	if riskRegister.NextReviewLevel == nil {
		return utility.NewValidationError(fmt.Errorf("risk has not been submitted for approval"), "not_submitted", "risk")
	}

	currentRiskWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByLevel(riskID, *riskRegister.NextReviewLevel, tx)
	if err != nil {
		return err
	}

	allowed := false
	for _, participant := range currentRiskWorkflow.RiskWorkflowParticipants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}
	if !allowed {
		return utility.NewForbiddenError("you have no permission to reject")
	}

	workflowStatus, err := arrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	riskWorkflowReview := &models.RiskWorkflowReview{
		RiskWorkflowID:   currentRiskWorkflow.ID,
		WorkflowStatusID: workflowStatus.ID,
		CreatedBYID:      currentUserID,
	}

	if err := arrs.repo.CreateRiskReview(riskWorkflowReview, tx); err != nil {
		return err
	}

	riskWorkflow, err := arrs.riskWorkflowSvc.GetRiskWorkflowByID(currentRiskWorkflow.ID, tx)
	if err != nil {
		return err
	}

	riskWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := arrs.riskWorkflowSvc.UpdateRiskWorkflow(riskWorkflow, tx); err != nil {
		return err
	}

	riskStatus, err := arrs.riskSvc.GetRiskStatusByType("rejected", tx)
	if err != nil {
		return err
	}

	riskRegister.RiskStatusID = riskStatus.ID
	riskRegister.NextReviewLevel = nil
	if err := arrs.repo.UpdateRisk(riskRegister, tx); err != nil {
		return err
	}

	riskDetail, err := arrs.riskSvc.GetFullDetailsByID(riskID, tx)
	if err != nil {
		return err
	}
	notify = func() {
		arrs.sendRejectedNotification(currentUserID, riskDetail)
	}

	return err
}

func (arrs *AssetRiskReviewService) notifyWorkflowParticipants(riskRegister *models.RiskRegister, workflowLevel *models.RiskWorkflow, currentUserID uint) {
	assignedBy, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch assigning user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskRegister.ID)
	entity := "asset-risk-review"
	action := "review-user-assign"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/asset-risk-register/%d/workflow", appURL, riskRegister.ID)

	for _, participant := range workflowLevel.RiskWorkflowParticipants {
		responsibleUser, err := arrs.userSvc.GetByID(participant.User.ID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.User.ID},
			ModuleIdentifier:     "vulnerability",
			NotificationType:     string(commonDTOs.ReviewType),
			ReferenceNumber:      *riskRegister.ReferenceNumber,
			HeadingEn:            "Asset Risk Acceptance Workflow",
			HeadingAr:            "سير عمل قبول مخاطر الأصول",
			TitleEn:              riskRegister.Title,
			TitleAr:              riskRegister.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"Title":                 riskRegister.Title,
			"ReferenceNumber":       riskRegister.ReferenceNumber,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Asset Risk Acceptance Workflow: %s - Review Assigned", riskRegister.Title),
			TemplateRelativePath: "asset/templates/email/asset_risk_workflow_review_user_mail.html",
			EmailData:            emailData,
		}

		if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "vulnerability"); err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func (arrs *AssetRiskReviewService) sendResponsibleUserAssignNotification(currentUserID uint, riskDetail *assetRiskDTOs.AssetRiskRegisterDetailDTO) {
	responsibleUserID := uint(riskDetail.ResponsibleUser.ID)

	responsibleUser, err := arrs.userSvc.GetByID(responsibleUserID)
	if err != nil {
		utility.LogError("Failed to fetch responsible user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.ID)
	entity := "asset-risk-register"
	action := "responsible-user-assign"

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{responsibleUserID},
		ModuleIdentifier:     "asset_risk",
		NotificationType:     string(commonDTOs.AssignmentType),
		HeadingEn:            "Risk Register Responsibility Assigned",
		HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
		ReferenceNumber:      riskDetail.ReferenceNumber,
		TitleEn:              riskDetail.Title,
		TitleAr:              riskDetail.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	emailData := map[string]interface{}{
		"FirstNameEn":          responsibleUser.FirstNameEn,
		"LastNameEn":           responsibleUser.LastNameEn,
		"RiskReferenceNumber":  riskDetail.ReferenceNumber,
		"RiskTitle":            riskDetail.Title,
		"RiskJustification":    utility.GetStringValue(riskDetail.Justification),
		"CreatedByFirstNameEn": riskDetail.CreatedByUser.FirstNameEN,
		"CreatedByLastNameEn":  riskDetail.CreatedByUser.LastNameEN,
		"AssignedOn":           time.Now().Format("2006-01-02"),
		"RedirectURL":          fmt.Sprintf("%s/asset/asset-risk-register/%d/overview", appURL, riskDetail.ID),
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{responsibleUser.Email},
		Subject:              fmt.Sprintf("Asset Risk Register/%s: Responsibility Assigned", riskDetail.Title),
		TemplateRelativePath: "asset/templates/email/asset_risk_register_create_responsible_user_mail.html",
		EmailData:            emailData,
	}

	if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk"); err != nil {
		utility.LogError("Failed to send responsible user assign notification", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (arrs *AssetRiskReviewService) sendApprovalNotification(currentUserID uint, riskDetail *assetRiskDTOs.AssetRiskRegisterDetailDTO) {
	approvedByUser, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch approving user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.ID)
	entity := "asset-risk-review"
	action := "approval"
	responsibleUserID := uint(riskDetail.ResponsibleUser.ID)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{responsibleUserID},
		ModuleIdentifier:     "vulnerability",
		NotificationType:     string(commonDTOs.ApprovalType),
		HeadingEn:            "Asset Risk Acceptance Workflow",
		HeadingAr:            "سير عمل قبول مخاطر الأصول",
		ReferenceNumber:      riskDetail.ReferenceNumber,
		TitleEn:              riskDetail.Title,
		TitleAr:              riskDetail.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/asset-risk-register/%d/overview", appURL, riskDetail.ID)

	emailData := map[string]interface{}{
		"FirstNameEn":           riskDetail.ResponsibleUser.FirstNameEN,
		"LastNameEn":            riskDetail.ResponsibleUser.LastNameEN,
		"Title":                 riskDetail.Title,
		"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
		"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
		"ApprovedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{riskDetail.ResponsibleUser.Email},
		Subject:              fmt.Sprintf("Asset Risk Acceptance Workflow: %s - Approved", riskDetail.Title),
		TemplateRelativePath: "asset/templates/email/asset_risk_workflow_approve_user_mail.html",
		EmailData:            emailData,
	}

	if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "vulnerability"); err != nil {
		utility.LogError("Failed to send approval notification email", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (arrs *AssetRiskReviewService) sendRejectedNotification(currentUserID uint, riskDetail *assetRiskDTOs.AssetRiskRegisterDetailDTO) {
	rejectedByUser, err := arrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch review performed user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, riskDetail.ID)
	entity := "asset-risk-review"
	action := "rejected"
	responsibleUserID := uint(riskDetail.ResponsibleUser.ID)

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{responsibleUserID},
		ModuleIdentifier:     "vulnerability",
		NotificationType:     string(commonDTOs.RejectType),
		HeadingEn:            "Asset Risk Acceptance Workflow",
		HeadingAr:            "سير عمل قبول مخاطر الأصول",
		ReferenceNumber:      riskDetail.ReferenceNumber,
		TitleEn:              riskDetail.Title,
		TitleAr:              riskDetail.Title,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/asset/asset-risk-register/%d/workflow", appURL, riskDetail.ID)

	emailData := map[string]interface{}{
		"FirstNameEn":           riskDetail.ResponsibleUser.FirstNameEN,
		"LastNameEn":            riskDetail.ResponsibleUser.LastNameEN,
		"Title":                 riskDetail.Title,
		"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
		"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
		"RejectedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{riskDetail.ResponsibleUser.Email},
		Subject:              fmt.Sprintf("Asset Risk Acceptance Workflow: %s - Rejected", riskDetail.Title),
		TemplateRelativePath: "asset/templates/email/asset_risk_workflow_reject_user_mail.html",
		EmailData:            emailData,
	}

	if err := arrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "vulnerability"); err != nil {
		utility.LogError("Failed to send rejected notification email", map[string]interface{}{
			"error": err.Error(),
		})
	}
}
