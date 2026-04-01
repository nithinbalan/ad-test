package services

import (
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
)

type InternalCommitteeWorkflowReviewServiceInterface interface {
	Submit(internalCommitteeID uint, currentUserID uint) error
	Approve(internalCommitteeID uint, currentUserID uint) error
	Revert(internalCommitteeID uint, currentUserID uint) error
}

type InternalCommitteeWorkflowReviewService struct {
	repo                 repositories.InternalCommitteeWorkflowReviewRepositoryInterface
	logService           logServices.ActivityLogServiceInterface
	committeeSvc         InternalCommitteeServiceInterface
	committeeWorkflowSvc InternalCommitteeWorkflowServiceInterface
	workflowStatusSvc    masterServices.WorkflowStatusServiceInterface
	committeeStatusSvc   InternalCommitteeStatusServiceInterface
	userSvc              UserServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
}

func NewInternalCommitteeWorkflowReviewService(
	repo repositories.InternalCommitteeWorkflowReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	committeeSvc InternalCommitteeServiceInterface,
	committeeWorkflowSvc InternalCommitteeWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	committeeStatusSvc InternalCommitteeStatusServiceInterface,
	userSvc UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) InternalCommitteeWorkflowReviewServiceInterface {
	return &InternalCommitteeWorkflowReviewService{
		repo:                 repo,
		logService:           logService,
		committeeSvc:         committeeSvc,
		committeeWorkflowSvc: committeeWorkflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		committeeStatusSvc:   committeeStatusSvc,
		userSvc:              userSvc,
		notificationSvc:      notificationSvc,
	}
}

func (icwrs *InternalCommitteeWorkflowReviewService) Submit(committeeID uint, currentUserID uint) error {
	committeeDTO, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID)
	if err != nil {
		return err
	}

	if committeeDTO.LatestVersion == nil {
		return fmt.Errorf("no latest version found for committee")
	}

	var level int

	// Check if already submitted
	if committeeDTO.LatestVersion.SubmittedByUser != nil {
		// Already submitted - check if it was reverted
		if committeeDTO.LatestVersion.WorkflowStatus == nil {
			return utility.NewForbiddenError("committee_already_submitted")
		}

		if committeeDTO.LatestVersion.WorkflowStatus.Type == "revert" {
			// Reverted - check if NextReviewLevel is set
			if committeeDTO.LatestVersion.NextReviewLevel == nil {
				return fmt.Errorf("cannot resubmit: next review level is not set for reverted committee")
			}
			level = *committeeDTO.LatestVersion.NextReviewLevel
		} else {
			// Already submitted and not reverted
			return utility.NewForbiddenError("committee_already_submitted")
		}
	} else {
		// Never submitted before - this is a new submission
		level = 1
	}

	tx := icwrs.repo.BeginTransaction()

	// Validate that all departments have participants before submission
	allDepartmentsHaveParticipants, err := icwrs.committeeSvc.ValidateAllDepartmentsHaveParticipants(committeeDTO.LatestVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !allDepartmentsHaveParticipants {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("all departments must have at least one participant"), "missing_participants", "committee")
	}

	currentDateStr := time.Now()

	err = icwrs.committeeWorkflowSvc.ValidateWorkflowBeforeCreate("committee_management", tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the model to update
	committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
	if err != nil {
		tx.Rollback()
		return err
	}

	var latestVersion *models.InternalCommitteeVersion
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersion = version
			break
		}
	}

	if latestVersion == nil {
		tx.Rollback()
		return fmt.Errorf("no latest version model found")
	}

	// Check if workflows exist for this version
	existingWorkflows, err := icwrs.committeeWorkflowSvc.GetWorkflowListData(committeeID, latestVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Only create workflows if they don't exist
	if len(existingWorkflows.InternalCommitteeWorkflows) == 0 {
		// Create new workflow only for first-time submission
		if err := icwrs.committeeWorkflowSvc.CreateInternalCommitteeWorkflow(committeeID, "committee_management", currentUserID, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	submitStatus, err := icwrs.workflowStatusSvc.GetByType("submit")
	if err != nil {
		tx.Rollback()
		return err
	}

	committeeInReviewStatus, err := icwrs.committeeStatusSvc.GetByType("in_review")
	if err != nil {
		tx.Rollback()
		return err
	}

	latestVersion.WorkflowStatusID = &submitStatus.ID
	latestVersion.SubmittedByID = &currentUserID
	latestVersion.SubmittedAt = &currentDateStr
	latestVersion.NextReviewLevel = &level
	latestVersion.UpdatedBYID = &currentUserID
	latestVersion.CommitteeStatusID = committeeInReviewStatus.ID

	if err := icwrs.repo.UpdateInternalCommitteeVersion(latestVersion, tx); err != nil {
		tx.Rollback()
		return err
	}

	committeeWorkflowLevel, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByLevel(latestVersion.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	committeeReview := &models.InternalCommitteeWorkflowReview{
		InternalCommitteeWorkflowID: committeeWorkflowLevel.ID,
		WorkflowStatusID:            submitStatus.ID,
		CreatedBYID:                 currentUserID,
	}

	if err := icwrs.repo.CreateInternalCommitteeWorkflowReview(committeeReview, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Send notifications asynchronously
	go func() {
		committeeDetail, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID)
		if err != nil {
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, committeeID)
		entity := "internal-committee-review"
		action := "review-user-assign"

		for _, participant := range committeeWorkflowLevel.InternalCommitteeWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      committeeDetail.ReferenceNumber,
				HeadingEn:            "Internal Committee Review",
				HeadingAr:            "مراجعة اللجنة الداخلية",
				TitleEn:              committeeDetail.LatestVersion.CommitteeNameEn,
				TitleAr:              committeeDetail.LatestVersion.CommitteeNameAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := icwrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			assignedBy, err := icwrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Internal Committee %s: Review Assigned", committeeDetail.LatestVersion.CommitteeNameEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 committeeDetail.LatestVersion.CommitteeNameEn,
				"ReferenceNumber":       committeeDetail.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/organization/internal-committee/%d/workflow", appURL, committeeID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/committee_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = icwrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (icwrs *InternalCommitteeWorkflowReviewService) Approve(committeeID uint, currentUserID uint) error {
	tx := icwrs.repo.BeginTransaction()

	committeeDetail, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentApprovalLevel := committeeDetail.LatestVersion.NextReviewLevel
	if currentApprovalLevel == nil {
		tx.Rollback()
		return fmt.Errorf("no review level found")
	}

	currentCommitteeWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByLevel(committeeDetail.LatestVersion.ID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := icwrs.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentCommitteeWorkflow.WorkflowLevelType.Type {
	case "individual":
		if len(currentCommitteeWorkflow.InternalCommitteeWorkflowParticipants) == 0 || currentUserID != currentCommitteeWorkflow.InternalCommitteeWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		committeeReview := &models.InternalCommitteeWorkflowReview{
			InternalCommitteeWorkflowID: currentCommitteeWorkflow.ID,
			WorkflowStatusID:            workflowStatus.ID,
			CreatedBYID:                 currentUserID,
		}
		if err := icwrs.repo.CreateInternalCommitteeWorkflowReview(committeeReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		committeeWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByID(currentCommitteeWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		committeeWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := icwrs.committeeWorkflowSvc.UpdateInternalCommitteeWorkflow(committeeWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextCommitteeWorkflowLevel, err := icwrs.committeeWorkflowSvc.GetNextLevel(committeeID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextCommitteeWorkflowLevel == nil {
			// Final approval - set committee status to approved
			approvedCommitteeStatus, err := icwrs.committeeStatusSvc.GetByType("approved")
			if err != nil {
				tx.Rollback()
				return err
			}

			approvedWorkflowStatus, err := icwrs.workflowStatusSvc.GetByType("approve")
			if err != nil {
				tx.Rollback()
				return err
			}

			committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Update latest version
			for _, version := range committee.InternalCommitteeVersions {
				if version.IsLatestVersion != nil && *version.IsLatestVersion {
					version.CommitteeStatusID = approvedCommitteeStatus.ID
					version.WorkflowStatusID = &approvedWorkflowStatus.ID
					approvedAt := time.Now()
					version.ApprovedAt = &approvedAt
					version.ApprovedByID = &currentUserID
					if err := icwrs.repo.UpdateInternalCommitteeVersion(version, tx); err != nil {
						tx.Rollback()
						return err
					}
					break
				}
			}
		} else {
			// Move to next level - set status to in_review
			inReviewStatus, err := icwrs.committeeStatusSvc.GetByType("in_review")
			if err != nil {
				tx.Rollback()
				return err
			}

			committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, version := range committee.InternalCommitteeVersions {
				if version.IsLatestVersion != nil && *version.IsLatestVersion {
					version.CommitteeStatusID = inReviewStatus.ID
					version.NextReviewLevel = &nextCommitteeWorkflowLevel.Level
					if err := icwrs.repo.UpdateInternalCommitteeVersion(version, tx); err != nil {
						tx.Rollback()
						return err
					}
					break
				}
			}

			nextWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByLevel(committeeDetail.LatestVersion.ID, nextCommitteeWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			icwrs.notifyNextLevelUsers(committeeID, nextWorkflow, committeeDetail)
		}

		committeeDetails, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		icwrs.sendApprovalNotification(committeeID, currentUserID, committeeDetails)

	case "multiple-users", "role":
		allowed := false
		for _, participant := range currentCommitteeWorkflow.InternalCommitteeWorkflowParticipants {
			if participant.UserID == currentUserID {
				allowed = true
				break
			}
		}
		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		committeeReview := &models.InternalCommitteeWorkflowReview{
			InternalCommitteeWorkflowID: currentCommitteeWorkflow.ID,
			WorkflowStatusID:            workflowStatus.ID,
			CreatedBYID:                 currentUserID,
		}
		if err := icwrs.repo.CreateInternalCommitteeWorkflowReview(committeeReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentCommitteeWorkflow.ApprovalTypes == "all" {
			allApproved, err := icwrs.committeeWorkflowSvc.CheckAllParticipantsApproved(currentCommitteeWorkflow.ID, currentCommitteeWorkflow.InternalCommitteeWorkflowParticipants, tx)
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
		committeeWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByID(currentCommitteeWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		committeeWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := icwrs.committeeWorkflowSvc.UpdateInternalCommitteeWorkflow(committeeWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		nextCommitteeWorkflowLevel, err := icwrs.committeeWorkflowSvc.GetNextLevel(committeeID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextCommitteeWorkflowLevel == nil {
			// Final approval
			approvedCommitteeStatus, err := icwrs.committeeStatusSvc.GetByType("approved")
			if err != nil {
				tx.Rollback()
				return err
			}

			approvedWorkflowStatus, err := icwrs.workflowStatusSvc.GetByType("approve")
			if err != nil {
				tx.Rollback()
				return err
			}

			committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, version := range committee.InternalCommitteeVersions {
				if version.IsLatestVersion != nil && *version.IsLatestVersion {
					version.CommitteeStatusID = approvedCommitteeStatus.ID
					version.WorkflowStatusID = &approvedWorkflowStatus.ID
					approvedAt := time.Now()
					version.ApprovedAt = &approvedAt
					version.ApprovedByID = &currentUserID
					if err := icwrs.repo.UpdateInternalCommitteeVersion(version, tx); err != nil {
						tx.Rollback()
						return err
					}
					break
				}
			}
		} else {
			// Move to next level - set status to in_review
			inReviewStatus, err := icwrs.committeeStatusSvc.GetByType("in_review")
			if err != nil {
				tx.Rollback()
				return err
			}

			committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
			if err != nil {
				tx.Rollback()
				return err
			}

			var latestVersionID uint
			for _, version := range committee.InternalCommitteeVersions {
				if version.IsLatestVersion != nil && *version.IsLatestVersion {
					latestVersionID = version.ID
					version.CommitteeStatusID = inReviewStatus.ID
					version.NextReviewLevel = &nextCommitteeWorkflowLevel.Level
					if err := icwrs.repo.UpdateInternalCommitteeVersion(version, tx); err != nil {
						tx.Rollback()
						return err
					}
					break
				}
			}

			nextWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByLevel(latestVersionID, nextCommitteeWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			icwrs.notifyNextLevelUsers(committeeID, nextWorkflow, committeeDetail)
		}

		committeeDetails, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		icwrs.sendApprovalNotification(committeeID, currentUserID, committeeDetails)
	}

	return tx.Commit()
}

func (icwrs *InternalCommitteeWorkflowReviewService) Revert(committeeID uint, currentUserID uint) error {
	revertStatus, err := icwrs.workflowStatusSvc.GetByType("revert")
	if err != nil {
		return err
	}

	pendingStatus, err := icwrs.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	tx := icwrs.repo.BeginTransaction()

	committeeDetail, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentApprovalLevel := committeeDetail.LatestVersion.NextReviewLevel
	if currentApprovalLevel == nil {
		tx.Rollback()
		return fmt.Errorf("no review level found")
	}

	currentLevel := *currentApprovalLevel

	// 1. Mark current level (N) as reverted
	currentCommitteeWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByLevel(committeeDetail.LatestVersion.ID, currentLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create revert review for current level
	committeeWorkflowReview := &models.InternalCommitteeWorkflowReview{
		InternalCommitteeWorkflowID: currentCommitteeWorkflow.ID,
		WorkflowStatusID:            revertStatus.ID,
		CreatedBYID:                 currentUserID,
	}

	if err := icwrs.repo.CreateInternalCommitteeWorkflowReview(committeeWorkflowReview, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Update current workflow status to reverted
	currentCommitteeWorkflow.WorkflowStatusID = revertStatus.ID
	if err := icwrs.committeeWorkflowSvc.UpdateInternalCommitteeWorkflow(currentCommitteeWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get all workflows for this version to reset levels
	committeeWorkflowList, err := icwrs.committeeWorkflowSvc.GetWorkflowListData(committeeID, committeeDetail.LatestVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Calculate previous level
	previousLevel := currentLevel - 1

	// 2. Reset only level N-1 (previous level) and levels > N
	for _, workflow := range committeeWorkflowList.InternalCommitteeWorkflows {
		// Skip levels below previous level - they should remain approved
		if workflow.Level < previousLevel {
			continue
		}

		// Skip current level - already marked as reverted above
		if workflow.Level == currentLevel {
			continue
		}

		// Reset level N-1 (previous level) and all levels > N
		if workflow.Level == previousLevel || workflow.Level > currentLevel {

			// Don't delete reviews - keep them for history
			// Just update workflow status to pending
			workflow.WorkflowStatusID = pendingStatus.ID
			if err := icwrs.committeeWorkflowSvc.UpdateInternalCommitteeWorkflow(workflow, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 4. Update committee version based on revert level
	committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Determine committee status and submission fields based on revert level
	var committeeStatus *models.InternalCommitteeStatus
	if previousLevel == 0 {
		// Reverting to level 0 means fully reverted - reset submission fields
		committeeStatus, err = icwrs.committeeStatusSvc.GetByType("revert")
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, version := range committee.InternalCommitteeVersions {
			if version.IsLatestVersion != nil && *version.IsLatestVersion {
				version.CommitteeStatusID = committeeStatus.ID
				version.WorkflowStatusID = &revertStatus.ID
				version.NextReviewLevel = &previousLevel
				// Reset submission fields for Level 0 revert
				version.SubmittedByID = nil
				version.SubmittedAt = nil
				if err := icwrs.repo.UpdateInternalCommitteeVersion(version, tx); err != nil {
					tx.Rollback()
					return err
				}
				break
			}
		}
	} else {
		// Reverting to any level > 0 means in review - keep submission fields
		committeeStatus, err = icwrs.committeeStatusSvc.GetByType("in_review")
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, version := range committee.InternalCommitteeVersions {
			if version.IsLatestVersion != nil && *version.IsLatestVersion {
				version.CommitteeStatusID = committeeStatus.ID
				version.WorkflowStatusID = &revertStatus.ID
				version.NextReviewLevel = &previousLevel
				// Do NOT reset submission fields for Level > 0 revert
				if err := icwrs.repo.UpdateInternalCommitteeVersion(version, tx); err != nil {
					tx.Rollback()
					return err
				}
				break
			}
		}
	}

	tx.Commit()

	// Notifications should be sent AFTER commit succeeds
	icwrs.sendRevertNotification(committeeID, currentUserID, committeeDetail)

	return nil
}

func (icwrs *InternalCommitteeWorkflowReviewService) sendApprovalNotification(id uint, currentUserID uint, committeeDetail *organizationDTOs.InternalCommitteeDetailDTO) {
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, committeeDetail.ID)
		entity := "internal-committee-review"
		action := "approval"

		// Get created by user ID
		createdByUserID := committeeDetail.CreatedByUser.ID

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(createdByUserID)},
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Internal Committee Approval",
			HeadingAr:            "الموافقة على اللجنة الداخلية",
			ReferenceNumber:      committeeDetail.ReferenceNumber,
			TitleEn:              committeeDetail.LatestVersion.CommitteeNameEn,
			TitleAr:              committeeDetail.LatestVersion.CommitteeNameAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		approvedByUser, err := icwrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdByUser, err := icwrs.userSvc.GetByID(uint(createdByUserID))
		if err != nil {
			utility.LogError("Failed to fetch created by user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailTo := []string{createdByUser.Email}
		subject := fmt.Sprintf("Internal Committee %s: Approval Notification", committeeDetail.LatestVersion.CommitteeNameEn)
		redirectURL := fmt.Sprintf("%s/organization/internal-committee/%d/workflow", appURL, committeeDetail.ID)

		if committeeDetail.LatestVersion.CommitteeStatus.Type == "approved" {
			redirectURL = fmt.Sprintf("%s/organization/internal-committee/%d/overview", appURL, id)
		}

		emailData := map[string]interface{}{
			"Title":                 committeeDetail.LatestVersion.CommitteeNameEn,
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/committee_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		err = icwrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (icwrs *InternalCommitteeWorkflowReviewService) notifyNextLevelUsers(id uint, nextCommitteeWorkflowLevel *models.InternalCommitteeWorkflow, committeeDetail *organizationDTOs.InternalCommitteeDetailDTO) {
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, committeeDetail.ID)
		entity := "internal-committee-review"
		action := "review-user-assign"

		for _, participant := range nextCommitteeWorkflowLevel.InternalCommitteeWorkflowParticipants {
			responsibleUser, err := icwrs.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(committeeDetail.CreatedByUser.ID)

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Internal Committee Review",
				HeadingAr:            "مراجعة اللجنة الداخلية",
				ReferenceNumber:      committeeDetail.ReferenceNumber,
				TitleEn:              committeeDetail.LatestVersion.CommitteeNameEn,
				TitleAr:              committeeDetail.LatestVersion.CommitteeNameAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Internal Committee %s: Review Assigned", committeeDetail.LatestVersion.CommitteeNameEn)
			appURL := utility.GetAppURL()
			redirectURL := fmt.Sprintf("%s/organization/internal-committee/%d/workflow", appURL, committeeDetail.ID)

			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 committeeDetail.LatestVersion.CommitteeNameEn,
				"AssignedByFirstNameEn": committeeDetail.CreatedByUser.FirstNameEN,
				"AssignedByLastNameEn":  committeeDetail.CreatedByUser.LastNameEN,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           redirectURL,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/committee_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = icwrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (icwrs *InternalCommitteeWorkflowReviewService) sendRevertNotification(id uint, currentUserID uint, committee *organizationDTOs.InternalCommitteeDetailDTO) {
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "internal-committee-review"
		action := "reverted"

		createdByUserID := uint(committee.CreatedByUser.ID)

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{createdByUserID},
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Internal Committee Rejected",
			HeadingAr:            "تم رفض اللجنة الداخلية",
			ReferenceNumber:      committee.ReferenceNumber,
			TitleEn:              committee.LatestVersion.CommitteeNameEn,
			TitleAr:              committee.LatestVersion.CommitteeNameAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := icwrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdByUser, err := icwrs.userSvc.GetByID(createdByUserID)
		if err != nil {
			utility.LogError("Failed to fetch created by user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		emailTo := []string{createdByUser.Email}
		subject := fmt.Sprintf("Internal Committee %s: Rejected Notification", committee.LatestVersion.CommitteeNameEn)
		redirectURL := fmt.Sprintf("%s/organization/internal-committee/%d/workflow", appURL, committee.ID)

		emailData := map[string]interface{}{
			"FirstNameEn":           createdByUser.FirstNameEn,
			"LastNameEn":            createdByUser.LastNameEn,
			"Title":                 committee.LatestVersion.CommitteeNameEn,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"RejectedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectURL,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/committee_workflow_revert_user_mail.html",
			EmailData:            emailData,
		}

		err = icwrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
