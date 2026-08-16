package services

import (
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/identity-access-review/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	"jethur/utility"
)

type IdentityReviewCampaignReviewServiceInterface interface {
	Submit(campaignID uint, currentUserID uint) error
	Approve(campaignID uint, currentUserID uint) error
	Reject(campaignID uint, currentUserID uint) error
	Resubmit(campaignID uint, currentUserID uint) error
	Revert(campaignID uint, currentUserID uint) error
}

type IdentityReviewCampaignReviewService struct {
	repo              repositories.IdentityReviewCampaignReviewRepositoryInterface
	campaignSvc       IdentityReviewCampaignServiceInterface
	workflowSvc       IdentityReviewCampaignWorkflowServiceInterface
	campaignStatusSvc ReviewCampaignStatusServiceInterface
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewIdentityReviewCampaignReviewService(
	repo repositories.IdentityReviewCampaignReviewRepositoryInterface,
	campaignSvc IdentityReviewCampaignServiceInterface,
	workflowSvc IdentityReviewCampaignWorkflowServiceInterface,
	campaignStatusSvc ReviewCampaignStatusServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) IdentityReviewCampaignReviewServiceInterface {
	return &IdentityReviewCampaignReviewService{
		repo:              repo,
		campaignSvc:       campaignSvc,
		workflowSvc:       workflowSvc,
		campaignStatusSvc: campaignStatusSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		notificationSvc:   notificationSvc,
	}
}

func (ircrs *IdentityReviewCampaignReviewService) Submit(campaignID uint, currentUserID uint) (err error) {
	tx := ircrs.repo.BeginTransaction()

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

	if err := ircrs.workflowSvc.ValidateWorkflowBeforeCreate(campaignID, tx); err != nil {
		return err
	}

	campaign, err := ircrs.campaignSvc.GetByID(campaignID, tx)
	if err != nil {
		return err
	}

	if campaign.NextReviewLevel != nil {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already submitted"),
			"workflow_already_submitted",
			"identity review campaign",
		)
	}

	approved, err := ircrs.campaignStatusSvc.GetByType("approved")
	if err != nil {
		return err
	}

	if campaign.CampaignStatusID == approved.ID {
		return utility.NewValidationError(
			fmt.Errorf("workflow has already completed"),
			"workflow_already_completed",
			"identity review campaign",
		)
	}

	if err := ircrs.workflowSvc.DeleteAllWorkflowsForCampaign(campaignID, tx); err != nil {
		return err
	}

	if err := ircrs.workflowSvc.CreateCampaignWorkflow(campaignID, currentUserID, tx); err != nil {
		return err
	}

	submitted, err := ircrs.campaignStatusSvc.GetByType("submitted")
	if err != nil {
		return err
	}

	level := 1
	now := time.Now()

	campaign.CampaignStatusID = submitted.ID
	campaign.SubmittedBy = &currentUserID
	campaign.SubmittedAt = &now
	campaign.NextReviewLevel = &level
	campaign.UpdatedBYID = &currentUserID

	if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
		return err
	}

	workflowLevel, err := ircrs.workflowSvc.GetWorkflowByLevel(campaignID, level, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := ircrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	if err := ircrs.createReviewWithHistory(
		workflowLevel.ID,
		campaignID,
		level,
		workflowStatus.ID,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	notify = func() {
		ircrs.notifyWorkflowParticipants(campaign, workflowLevel, currentUserID)
	}

	return nil
}

func (ircrs *IdentityReviewCampaignReviewService) Approve(campaignID uint, currentUserID uint) (err error) {
	tx := ircrs.repo.BeginTransaction()

	var notifications []func()

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

		if len(notifications) > 0 {
			go func() {
				for _, notify := range notifications {
					notify()
				}
			}()
		}
	}()

	campaign, err := ircrs.campaignSvc.GetByID(campaignID, tx)
	if err != nil {
		return err
	}

	if campaign.NextReviewLevel == nil {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"identity review campaign",
		)
	}

	currentApprovalLevel := *campaign.NextReviewLevel
	currentWorkflow, err := ircrs.workflowSvc.GetWorkflowByLevel(campaignID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := ircrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		return err
	}

	pendingStatus, err := ircrs.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual", "hod":
		if len(currentWorkflow.IdentityReviewCampaignWorkflowParticipants) == 0 ||
			currentWorkflow.IdentityReviewCampaignWorkflowParticipants[0].UserID != currentUserID {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		approvingParticipant := currentWorkflow.IdentityReviewCampaignWorkflowParticipants[0]

		if approvingParticipant.WorkflowStatusID != nil {
			return utility.NewValidationError(
				fmt.Errorf("workflow participant has already reviewed"),
				"workflow_already_reviewed",
				"identity review campaign",
			)
		}

		if err := ircrs.createReviewWithHistory(
			currentWorkflow.ID,
			campaignID,
			currentApprovalLevel,
			workflowStatus.ID,
			currentUserID,
			tx,
		); err != nil {
			return err
		}

		approvingParticipant.WorkflowStatusID = &workflowStatus.ID

		if err := ircrs.workflowSvc.UpdateParticipant(approvingParticipant, tx); err != nil {
			return err
		}

		currentWorkflow.WorkflowStatusID = workflowStatus.ID

		if err := ircrs.workflowSvc.UpdateWorkflow(currentWorkflow, tx); err != nil {
			return err
		}

		notifications = append(notifications, func() {
			ircrs.sendApprovalNotification(currentUserID, campaign.CreatedBYID, campaign)
		})

	case "multiple-users", "role":
		var approvingParticipant *models.IdentityReviewCampaignWorkflowParticipant

		for _, participant := range currentWorkflow.IdentityReviewCampaignWorkflowParticipants {
			if participant.UserID == currentUserID {
				approvingParticipant = participant
				break
			}
		}

		if approvingParticipant == nil {
			return utility.NewForbiddenError("you have no permission to approve")
		}

		if approvingParticipant.WorkflowStatusID != nil {
			return utility.NewValidationError(
				fmt.Errorf("workflow participant has already reviewed"),
				"workflow_already_reviewed",
				"identity review campaign",
			)
		}

		if err := ircrs.createReviewWithHistory(
			currentWorkflow.ID,
			campaignID,
			currentApprovalLevel,
			workflowStatus.ID,
			currentUserID,
			tx,
		); err != nil {
			return err
		}

		approvingParticipant.WorkflowStatusID = &workflowStatus.ID

		if err := ircrs.workflowSvc.UpdateParticipant(approvingParticipant, tx); err != nil {
			return err
		}

		notifications = append(notifications, func() {
			ircrs.sendApprovalNotification(currentUserID, campaign.CreatedBYID, campaign)
		})

		if currentWorkflow.ApprovalTypes == "all" {
			allApproved, err := ircrs.workflowSvc.CheckAllParticipantsApproved(
				currentWorkflow.ID,
				currentWorkflow.IdentityReviewCampaignWorkflowParticipants,
				tx,
			)
			if err != nil {
				return err
			}

			if !allApproved {
				return nil
			}
		}

		currentWorkflow.WorkflowStatusID = workflowStatus.ID

		if err := ircrs.workflowSvc.UpdateWorkflow(currentWorkflow, tx); err != nil {
			return err
		}
	}

	nextWorkflowLevel, err := ircrs.workflowSvc.GetNextLevel(campaignID, currentApprovalLevel, tx)
	if err != nil {
		return err
	}

	if nextWorkflowLevel == nil {
		approvedStatus, err := ircrs.campaignStatusSvc.GetByType("approved")
		if err != nil {
			return err
		}

		campaign.CampaignStatusID = approvedStatus.ID
		campaign.NextReviewLevel = nil
		campaign.ReviewedBy = &currentUserID
		campaign.UpdatedBYID = &currentUserID

		if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
			return err
		}

		return nil
	}

	inProgressStatus, err := ircrs.campaignStatusSvc.GetByType("in_progress")
	if err != nil {
		return err
	}

	if err := ircrs.workflowSvc.UpdateWorkflowStatus(nextWorkflowLevel.ID, pendingStatus.ID, tx); err != nil {
		return err
	}

	campaign.CampaignStatusID = inProgressStatus.ID
	campaign.NextReviewLevel = &nextWorkflowLevel.Level
	campaign.UpdatedBYID = &currentUserID

	if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
		return err
	}

	nextWorkflow, err := ircrs.workflowSvc.GetWorkflowByLevel(campaignID, nextWorkflowLevel.Level, tx)
	if err != nil {
		return err
	}

	notifications = append(notifications, func() {
		ircrs.notifyWorkflowParticipants(campaign, nextWorkflow, currentUserID)
	})

	return nil
}

func (ircrs *IdentityReviewCampaignReviewService) Reject(campaignID uint, currentUserID uint) (err error) {
	tx := ircrs.repo.BeginTransaction()

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

	campaign, err := ircrs.campaignSvc.GetByID(campaignID, tx)
	if err != nil {
		return err
	}

	rejectedStatus, err := ircrs.campaignStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}

	if campaign.NextReviewLevel == nil ||
		campaign.CampaignStatusID == rejectedStatus.ID {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"identity review campaign",
		)
	}

	currentReviewLevel := *campaign.NextReviewLevel

	currentWorkflow, err := ircrs.workflowSvc.GetWorkflowByLevel(campaignID, currentReviewLevel, tx)
	if err != nil {
		return err
	}

	var rejectingParticipant *models.IdentityReviewCampaignWorkflowParticipant

	for _, participant := range currentWorkflow.IdentityReviewCampaignWorkflowParticipants {
		if participant.UserID == currentUserID {
			rejectingParticipant = participant
			break
		}
	}

	if rejectingParticipant == nil {
		return utility.NewForbiddenError("you have no permission to reject")
	}

	if rejectingParticipant.WorkflowStatusID != nil {
		return utility.NewValidationError(
			fmt.Errorf("workflow participant has already reviewed"),
			"workflow_already_reviewed",
			"identity review campaign",
		)
	}

	workflowStatus, err := ircrs.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	if err := ircrs.createReviewWithHistory(
		currentWorkflow.ID,
		campaignID,
		currentReviewLevel,
		workflowStatus.ID,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	rejectingParticipant.WorkflowStatusID = &workflowStatus.ID

	if err := ircrs.workflowSvc.UpdateParticipant(rejectingParticipant, tx); err != nil {
		return err
	}

	currentWorkflow.WorkflowStatusID = workflowStatus.ID

	if err := ircrs.workflowSvc.UpdateWorkflow(currentWorkflow, tx); err != nil {
		return err
	}

	var submittedByID uint

	if campaign.SubmittedBy != nil {
		submittedByID = *campaign.SubmittedBy
	}

	campaign.CampaignStatusID = rejectedStatus.ID
	campaign.UpdatedBYID = &currentUserID

	if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
		return err
	}

	if submittedByID != 0 {
		notify = func() {
			ircrs.sendRejectedNotification(currentUserID, submittedByID, campaign)
		}
	}

	return nil
}

func (ircrs *IdentityReviewCampaignReviewService) Resubmit(campaignID uint, currentUserID uint) (err error) {
	tx := ircrs.repo.BeginTransaction()

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

	if err := ircrs.workflowSvc.ValidateWorkflowBeforeCreate(campaignID, tx); err != nil {
		return err
	}

	campaign, err := ircrs.campaignSvc.GetByID(campaignID, tx)
	if err != nil {
		return err
	}

	rejectedStatus, err := ircrs.campaignStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}

	if campaign.CampaignStatusID != rejectedStatus.ID {
		return utility.NewValidationError(
			fmt.Errorf("campaign has not been rejected"),
			"workflow_not_rejected",
			"identity review campaign",
		)
	}

	if err := ircrs.workflowSvc.DeleteAllWorkflowsForCampaign(campaignID, tx); err != nil {
		return err
	}

	if err := ircrs.workflowSvc.CreateCampaignWorkflow(campaignID, currentUserID, tx); err != nil {
		return err
	}

	submittedStatus, err := ircrs.campaignStatusSvc.GetByType("submitted")
	if err != nil {
		return err
	}

	level := 1
	now := time.Now()

	campaign.CampaignStatusID = submittedStatus.ID
	campaign.SubmittedBy = &currentUserID
	campaign.SubmittedAt = &now
	campaign.NextReviewLevel = &level
	campaign.UpdatedBYID = &currentUserID

	if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
		return err
	}

	workflowLevel, err := ircrs.workflowSvc.GetWorkflowByLevel(campaignID, level, tx)
	if err != nil {
		return err
	}

	workflowStatus, err := ircrs.workflowStatusSvc.GetByType("resubmit")
	if err != nil {
		return err
	}

	if err := ircrs.createReviewWithHistory(
		workflowLevel.ID,
		campaignID,
		level,
		workflowStatus.ID,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	notify = func() {
		ircrs.sendResubmittedNotification(currentUserID, campaign, workflowLevel)
	}

	return nil
}

func (ircrs *IdentityReviewCampaignReviewService) Revert(campaignID uint, currentUserID uint) (err error) {
	tx := ircrs.repo.BeginTransaction()

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

	campaign, err := ircrs.campaignSvc.GetByID(campaignID, tx)
	if err != nil {
		return err
	}

	rejectedStatus, err := ircrs.campaignStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}

	if campaign.NextReviewLevel == nil ||
		campaign.CampaignStatusID == rejectedStatus.ID {
		return utility.NewValidationError(
			fmt.Errorf("workflow not active"),
			"workflow_not_active",
			"identity review campaign",
		)
	}

	currentLevel := *campaign.NextReviewLevel

	currentWorkflow, err := ircrs.workflowSvc.GetWorkflowByLevel(campaignID, currentLevel, tx)
	if err != nil {
		return err
	}

	allowed := false

	for _, participant := range currentWorkflow.IdentityReviewCampaignWorkflowParticipants {
		if participant.UserID == currentUserID {
			allowed = true
			break
		}
	}

	if !allowed {
		return utility.NewForbiddenError("you have no permission to revert")
	}

	revertStatus, err := ircrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	inProgressStatus, err := ircrs.campaignStatusSvc.GetByType("in_progress")
	if err != nil {
		return err
	}

	campaign.CampaignStatusID = inProgressStatus.ID

	if err := ircrs.createReviewWithHistory(
		currentWorkflow.ID,
		campaignID,
		currentLevel,
		revertStatus.ID,
		currentUserID,
		tx,
	); err != nil {
		return err
	}

	if err := ircrs.workflowSvc.UpdateWorkflowStatus(currentWorkflow.ID, revertStatus.ID, tx); err != nil {
		return err
	}

	previousLevel := currentLevel - 1

	var previousLevelWorkflow *models.IdentityReviewCampaignWorkflow

	if previousLevel > 0 {
		pendingStatus, err := ircrs.workflowStatusSvc.GetByType("pending")
		if err != nil {
			return err
		}

		previousLevelWorkflow, err = ircrs.workflowSvc.GetWorkflowByLevel(campaignID, previousLevel, tx)
		if err != nil {
			return err
		}

		if err := ircrs.workflowSvc.DeleteReviewsByWorkflowID(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := ircrs.workflowSvc.ResetParticipantStatuses(previousLevelWorkflow.ID, tx); err != nil {
			return err
		}

		if err := ircrs.workflowSvc.UpdateWorkflowStatus(previousLevelWorkflow.ID, pendingStatus.ID, tx); err != nil {
			return err
		}

		campaign.NextReviewLevel = &previousLevel
	} else {
		submittedByID := campaign.SubmittedBy

		campaign.SubmittedBy = nil
		campaign.SubmittedAt = nil
		campaign.NextReviewLevel = nil
		campaign.ReviewedBy = nil

		campaign.UpdatedBYID = &currentUserID

		if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
			return err
		}

		if submittedByID != nil {
			notify = func() {
				ircrs.sendRevertNotification(currentUserID, *submittedByID, campaign)
			}
		}

		return nil
	}

	campaign.UpdatedBYID = &currentUserID

	if err := ircrs.repo.UpdateCampaign(campaign, tx); err != nil {
		return err
	}

	notify = func() {
		ircrs.notifyRevertedToPreviousLevel(campaign, previousLevelWorkflow, currentUserID)
	}

	return nil
}

func (ircrs *IdentityReviewCampaignReviewService) notifyWorkflowParticipants(campaign *models.IdentityReviewCampaign, workflowLevel *models.IdentityReviewCampaignWorkflow, currentUserID uint) {
	assignedBy, err := ircrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch assigning user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, campaign.ID)
	entity := "identity-access-review-campaign"
	action := "review-user-assign"
	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/identity-access-review/campaigns/%d/overview", appURL, campaign.ID)
	notifiedUserIDs := make(map[uint]struct{})

	for _, participant := range workflowLevel.IdentityReviewCampaignWorkflowParticipants {
		if _, alreadyNotified := notifiedUserIDs[participant.UserID]; alreadyNotified {
			continue
		}
		notifiedUserIDs[participant.UserID] = struct{}{}

		responsibleUser, err := ircrs.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch workflow participant", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "identity_access_review",
			NotificationType:     string(commonDTOs.ReviewType),
			ReferenceNumber:      campaign.ReferenceNumber,
			HeadingEn:            "Identity Review Campaign Workflow",
			HeadingAr:            "سير عمل حملة مراجعة الهوية",
			TitleEn:              campaign.CampaignName,
			TitleAr:              campaign.CampaignName,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"CampaignName":          campaign.CampaignName,
			"ReferenceNumber":       campaign.ReferenceNumber,
			"AssignedByFirstNameEn": assignedBy.FirstNameEn,
			"AssignedByLastNameEn":  assignedBy.LastNameEn,
			"AssignedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Identity Review Campaign: %s - Review Assigned", campaign.CampaignName),
			TemplateRelativePath: "identity-access-review/templates/email/identity_review_campaign_workflow_review_user_mail.html",
			EmailData:            emailData,
		}

		if err := ircrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "identity_access_review_campaign"); err != nil {
			utility.LogError("Failed to send identity review workflow notification", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func (ircrs *IdentityReviewCampaignReviewService) notifyRevertedToPreviousLevel(campaign *models.IdentityReviewCampaign, workflowLevel *models.IdentityReviewCampaignWorkflow, currentUserID uint) {
	revertedByUser, err := ircrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, campaign.ID)

	entity := "identity-access-review-campaign"
	action := "reverted"

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/identity-access-review/campaigns/%d/overview", appURL, campaign.ID)
	notifiedUserIDs := make(map[uint]struct{})

	for _, participant := range workflowLevel.IdentityReviewCampaignWorkflowParticipants {
		if _, alreadyNotified := notifiedUserIDs[participant.UserID]; alreadyNotified {
			continue
		}
		notifiedUserIDs[participant.UserID] = struct{}{}

		responsibleUser, err := ircrs.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch identity review workflow participant", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "identity_access_review",
			NotificationType:     string(commonDTOs.RevertType),
			ReferenceNumber:      campaign.ReferenceNumber,
			HeadingEn:            "Identity Review Campaign Workflow",
			HeadingAr:            "سير عمل حملة مراجعة الهوية",
			TitleEn:              campaign.CampaignName,
			TitleAr:              campaign.CampaignName,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"CampaignName":          campaign.CampaignName,
			"ReferenceNumber":       campaign.ReferenceNumber,
			"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
			"RevertedByLastNameEn":  revertedByUser.LastNameEn,
			"RevertedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Identity Review Campaign: %s - Reverted", campaign.CampaignName),
			TemplateRelativePath: "identity-access-review/templates/email/identity_review_campaign_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		if err := ircrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "identity_access_review_campaign"); err != nil {
			utility.LogError("Failed to send identity review previous-level notification", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func (ircrs *IdentityReviewCampaignReviewService) sendApprovalNotification(currentUserID uint, createdByID uint, campaign *models.IdentityReviewCampaign) {
	approvedByUser, err := ircrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch approving user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	createdByUser, err := ircrs.userSvc.GetByID(createdByID)
	if err != nil {
		utility.LogError("Failed to fetch campaign created user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, campaign.ID)

	entity := "identity-access-review-campaign"
	action := "approval"

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{createdByID},
		ModuleIdentifier:     "identity_access_review",
		NotificationType:     string(commonDTOs.ApprovalType),
		ReferenceNumber:      campaign.ReferenceNumber,
		HeadingEn:            "Identity Review Campaign Workflow",
		HeadingAr:            "سير عمل حملة مراجعة الهوية",
		TitleEn:              campaign.CampaignName,
		TitleAr:              campaign.CampaignName,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/identity-access-review/campaigns/%d/overview", appURL, campaign.ID)

	emailData := map[string]interface{}{
		"FirstNameEn":           createdByUser.FirstNameEn,
		"LastNameEn":            createdByUser.LastNameEn,
		"CampaignName":          campaign.CampaignName,
		"ReferenceNumber":       campaign.ReferenceNumber,
		"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
		"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
		"ApprovedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{createdByUser.Email},
		Subject:              fmt.Sprintf("Identity Review Campaign: %s - Approved", campaign.CampaignName),
		TemplateRelativePath: "identity-access-review/templates/email/identity_review_campaign_workflow_approve_user_mail.html",
		EmailData:            emailData,
	}

	if err := ircrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "identity_access_review_campaign"); err != nil {
		utility.LogError("Failed to send identity review campaign approval notification", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (ircrs *IdentityReviewCampaignReviewService) sendRejectedNotification(currentUserID uint, submittedByID uint, campaign *models.IdentityReviewCampaign) {
	rejectedByUser, err := ircrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch rejecting user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	submittedByUser, err := ircrs.userSvc.GetByID(submittedByID)
	if err != nil {
		utility.LogError("Failed to fetch campaign submitted user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, campaign.ID)

	entity := "identity-access-review-campaign"
	action := "rejected"

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{submittedByID},
		ModuleIdentifier:     "identity_access_review",
		NotificationType:     string(commonDTOs.RejectType),
		ReferenceNumber:      campaign.ReferenceNumber,
		HeadingEn:            "Identity Review Campaign Workflow",
		HeadingAr:            "سير عمل حملة مراجعة الهوية",
		TitleEn:              campaign.CampaignName,
		TitleAr:              campaign.CampaignName,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/identity-access-review/campaigns/%d/overview", appURL, campaign.ID)

	emailData := map[string]interface{}{
		"FirstNameEn":           submittedByUser.FirstNameEn,
		"LastNameEn":            submittedByUser.LastNameEn,
		"CampaignName":          campaign.CampaignName,
		"ReferenceNumber":       campaign.ReferenceNumber,
		"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
		"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
		"RejectedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{submittedByUser.Email},
		Subject:              fmt.Sprintf("Identity Review Campaign: %s - Rejected", campaign.CampaignName),
		TemplateRelativePath: "identity-access-review/templates/email/identity_review_campaign_workflow_reject_user_mail.html",
		EmailData:            emailData,
	}

	if err := ircrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "identity_access_review_campaign"); err != nil {
		utility.LogError("Failed to send identity review campaign rejection notification", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (ircrs *IdentityReviewCampaignReviewService) sendResubmittedNotification(currentUserID uint, campaign *models.IdentityReviewCampaign, workflowLevel *models.IdentityReviewCampaignWorkflow) {
	resubmittedByUser, err := ircrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch resubmitting user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, campaign.ID)

	entity := "identity-access-review-campaign"
	action := "resubmitted"

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/identity-access-review/campaigns/%d/overview", appURL, campaign.ID)
	notifiedUserIDs := make(map[uint]struct{})

	for _, participant := range workflowLevel.IdentityReviewCampaignWorkflowParticipants {
		if _, alreadyNotified := notifiedUserIDs[participant.UserID]; alreadyNotified {
			continue
		}
		notifiedUserIDs[participant.UserID] = struct{}{}

		responsibleUser, err := ircrs.userSvc.GetByID(participant.UserID)
		if err != nil {
			utility.LogError("Failed to fetch identity review workflow participant", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{participant.UserID},
			ModuleIdentifier:     "identity_access_review",
			NotificationType:     string(commonDTOs.ReviewType),
			ReferenceNumber:      campaign.ReferenceNumber,
			HeadingEn:            "Identity Review Campaign Workflow",
			HeadingAr:            "سير عمل حملة مراجعة الهوية",
			TitleEn:              campaign.CampaignName,
			TitleAr:              campaign.CampaignName,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailData := map[string]interface{}{
			"FirstNameEn":              responsibleUser.FirstNameEn,
			"LastNameEn":               responsibleUser.LastNameEn,
			"CampaignName":             campaign.CampaignName,
			"ReferenceNumber":          campaign.ReferenceNumber,
			"ResubmittedByFirstNameEn": resubmittedByUser.FirstNameEn,
			"ResubmittedByLastNameEn":  resubmittedByUser.LastNameEn,
			"ResubmittedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":              redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Identity Review Campaign: %s - Resubmitted", campaign.CampaignName),
			TemplateRelativePath: "identity-access-review/templates/email/identity_review_campaign_workflow_resubmit_user_mail.html",
			EmailData:            emailData,
		}

		if err := ircrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "identity_access_review_campaign"); err != nil {
			utility.LogError("Failed to send identity review resubmission notification", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}
}

func (ircrs *IdentityReviewCampaignReviewService) sendRevertNotification(currentUserID uint, submittedByID uint, campaign *models.IdentityReviewCampaign) {
	revertedByUser, err := ircrs.userSvc.GetByID(currentUserID)
	if err != nil {
		utility.LogError("Failed to fetch reverting user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	submittedByUser, err := ircrs.userSvc.GetByID(submittedByID)
	if err != nil {
		utility.LogError("Failed to fetch campaign submitted user", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	urlParam := fmt.Sprintf(`{"identity_review_campaign_id":%d}`, campaign.ID)

	entity := "identity-access-review-campaign"
	action := "reverted"

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserIDs:              []uint{submittedByID},
		ModuleIdentifier:     "identity_access_review",
		NotificationType:     string(commonDTOs.RevertType),
		ReferenceNumber:      campaign.ReferenceNumber,
		HeadingEn:            "Identity Review Campaign Workflow",
		HeadingAr:            "سير عمل حملة مراجعة الهوية",
		TitleEn:              campaign.CampaignName,
		TitleAr:              campaign.CampaignName,
		Entity:               &entity,
		Action:               &action,
		URLParam:             urlParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	appURL := utility.GetAppURL()
	redirectURL := fmt.Sprintf("%s/identity-access-review/campaigns/%d/overview", appURL, campaign.ID)

	emailData := map[string]interface{}{
		"FirstNameEn":           submittedByUser.FirstNameEn,
		"LastNameEn":            submittedByUser.LastNameEn,
		"CampaignName":          campaign.CampaignName,
		"ReferenceNumber":       campaign.ReferenceNumber,
		"RevertedByFirstNameEn": revertedByUser.FirstNameEn,
		"RevertedByLastNameEn":  revertedByUser.LastNameEn,
		"RevertedAt":            time.Now().Format("2006-01-02"),
		"RedirectURL":           redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              []string{submittedByUser.Email},
		Subject:              fmt.Sprintf("Identity Review Campaign: %s - Reverted", campaign.CampaignName),
		TemplateRelativePath: "identity-access-review/templates/email/identity_review_campaign_workflow_revert_user_mail.html",
		EmailData:            emailData,
	}

	if err := ircrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "identity_access_review_campaign"); err != nil {
		utility.LogError("Failed to send identity review campaign revert notification", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (ircrs *IdentityReviewCampaignReviewService) createReviewWithHistory(workflowID uint, campaignID uint, level int, statusID uint, userID uint, tx utility.DatabaseTransaction) error {
	review := &models.IdentityReviewCampaignWorkflowReview{
		IdentityReviewCampaignWorkflowID: workflowID,
		WorkflowStatusID:                 statusID,
		CreatedByID:                      userID,
	}
	if err := ircrs.repo.CreateReview(review, tx); err != nil {
		return err
	}
	history := &models.IdentityReviewCampaignWorkflowHistory{
		IdentityReviewCampaignID: campaignID,
		Level:                    level,
		WorkflowStatusID:         statusID,
		CreatedByID:              userID,
		CreatedAt:                time.Now(),
	}
	return ircrs.repo.CreateHistory(history, tx)
}
