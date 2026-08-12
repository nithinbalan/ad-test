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
	SetInternalCommitteeService(svc InternalCommitteeServiceInterface)
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
	organizationUnitSvc  OrganizationUnitsServiceInterface
	committeeTermSvc     CommitteeTermServiceInterface
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
	organizationUnitSvc OrganizationUnitsServiceInterface,
	committeeTermSvc CommitteeTermServiceInterface,
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
		organizationUnitSvc:  organizationUnitSvc,
		committeeTermSvc:     committeeTermSvc,
	}
}

func (icwrs *InternalCommitteeWorkflowReviewService) SetInternalCommitteeService(svc InternalCommitteeServiceInterface) {
	icwrs.committeeSvc = svc
}

func (icwrs *InternalCommitteeWorkflowReviewService) Submit(committeeID uint, currentUserID uint) error {
	committeeDTO, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID)
	if err != nil {
		return err
	}

	// parentUnit, err := icwrs.organizationUnitSvc.GetParentOrganizationUnitAndCurrentUnitHead(committeeDTO.RequestingEntityID)
	// if err != nil {
	// 	return err
	// }

	// parentHead, err := icwrs.organizationUnitSvc.GetParentOrganizationUnitAndCurrentUnitHead(parentUnit.OrganizationUnit.ID)
	// if err != nil {
	// 	return err
	// }

	// if parentHead.HeadUser == nil {
	// 	return utility.NewForbiddenError(
	// 		fmt.Sprintf(
	// 			"No head user is assigned to the parent organization unit '%s'. Please contact the administrator.",
	// 			parentUnit.OrganizationUnit.TitleEn,
	// 		),
	// 	)
	// }

	// if parentHead.HeadUser.ID != currentUserID {
	// 	return utility.NewForbiddenError(
	// 		fmt.Sprintf(
	// 			"You are not authorized to submit this request. Only the head of the parent organization unit '%s' can perform this action.",
	// 			parentUnit.OrganizationUnit.TitleEn,
	// 		),
	// 	)
	// }

	if committeeDTO.LatestVersion == nil {
		return fmt.Errorf("no latest version found for committee")
	}

	if committeeDTO.LatestVersion.CreatedByUser.ID != int(currentUserID) {
		return utility.NewForbiddenError(
			"You are not authorized to submit this request. Only the Requested by user can perform this action.",
		)
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

	// Create submit review in both current and history tables
	if err := icwrs.createReviewWithHistory(committeeWorkflowLevel.ID, submitStatus.ID, currentUserID, tx); err != nil {
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
				"RedirectURL":           fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeID),
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
	commettetermtype, err := icwrs.committeeTermSvc.GetByType("permanent")
	if err != nil {
		return err
	}

	tx := icwrs.repo.BeginTransaction()

	committeeDetail, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	commetteterm := committeeDetail.LatestVersion.CommitteeTerm

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
	case "individual", "custom-roles", "hod":
		if len(currentCommitteeWorkflow.InternalCommitteeWorkflowParticipants) == 0 || currentUserID != currentCommitteeWorkflow.InternalCommitteeWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		// Create approval review in both current and history tables
		if err := icwrs.createReviewWithHistory(currentCommitteeWorkflow.ID, workflowStatus.ID, currentUserID, tx); err != nil {
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
			// Final approval - check if it's disband or regular approval
			committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
			if err != nil {
				tx.Rollback()
				return err
			}

			approvedWorkflowStatus, err := icwrs.workflowStatusSvc.GetByType("approve")
			if err != nil {
				tx.Rollback()
				return err
			}

			// Update latest version
			for _, version := range committee.InternalCommitteeVersions {
				if version.IsLatestVersion != nil && *version.IsLatestVersion {
					// Check if disband is requested
					var finalCommitteeStatus *models.InternalCommitteeStatus
					if version.IsDisbandRequested {
						// Set to disbanded
						finalCommitteeStatus, err = icwrs.committeeStatusSvc.GetByType("disbanded")
						if err != nil {
							tx.Rollback()
							return err
						}
					} else {
						// Set to approved
						finalCommitteeStatus, err = icwrs.committeeStatusSvc.GetByType("approved")
						if err != nil {
							tx.Rollback()
							return err
						}
					}

					if commetteterm.ID == commettetermtype.ID {
						endDate := time.Now().AddDate(3, 0, 0)
						version.TermEndDate = &endDate
					}
					approvedAt := time.Now()
					version.TermStartDate = &approvedAt
					version.CommitteeStatusID = finalCommitteeStatus.ID
					version.WorkflowStatusID = &approvedWorkflowStatus.ID
					version.NextReviewLevel = nil
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
			// Move to next level - set status to in_review and workflow status to submit
			inReviewStatus, err := icwrs.committeeStatusSvc.GetByType("in_review")
			if err != nil {
				tx.Rollback()
				return err
			}

			submitStatus, err := icwrs.workflowStatusSvc.GetByType("submit")
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
					version.WorkflowStatusID = &submitStatus.ID
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

	case "multiple-users", "role", "user-and-reporting-user":
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

		// Create approval review in both current and history tables
		if err := icwrs.createReviewWithHistory(currentCommitteeWorkflow.ID, workflowStatus.ID, currentUserID, tx); err != nil {
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
				// Update workflow status to pending if it was reverted
				if currentCommitteeWorkflow.WorkflowStatus != nil && currentCommitteeWorkflow.WorkflowStatus.Type == "revert" {
					pendingStatus, err := icwrs.workflowStatusSvc.GetByType("pending")
					if err != nil {
						tx.Rollback()
						return err
					}

					committeeWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByID(currentCommitteeWorkflow.ID, tx)
					if err != nil {
						tx.Rollback()
						return err
					}

					committeeWorkflow.WorkflowStatusID = pendingStatus.ID
					if err := icwrs.committeeWorkflowSvc.UpdateInternalCommitteeWorkflow(committeeWorkflow, tx); err != nil {
						tx.Rollback()
						return err
					}
				}

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
			// Final approval - check if it's disband or regular approval
			committee, err := icwrs.committeeSvc.GetLatestVersionByID(committeeID)
			if err != nil {
				tx.Rollback()
				return err
			}

			approvedWorkflowStatus, err := icwrs.workflowStatusSvc.GetByType("approve")
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, version := range committee.InternalCommitteeVersions {
				if version.IsLatestVersion != nil && *version.IsLatestVersion {
					// Check if disband is requested
					var finalCommitteeStatus *models.InternalCommitteeStatus
					if version.IsDisbandRequested {
						// Set to disbanded
						finalCommitteeStatus, err = icwrs.committeeStatusSvc.GetByType("disbanded")
						if err != nil {
							tx.Rollback()
							return err
						}
					} else {
						// Set to approved
						finalCommitteeStatus, err = icwrs.committeeStatusSvc.GetByType("approved")
						if err != nil {
							tx.Rollback()
							return err
						}
					}
					if commetteterm.ID == commettetermtype.ID {
						endDate := time.Now().AddDate(3, 0, 0)
						version.TermEndDate = &endDate
					}
					version.CommitteeStatusID = finalCommitteeStatus.ID
					approvedAt := time.Now()
					version.TermStartDate = &approvedAt
					version.WorkflowStatusID = &approvedWorkflowStatus.ID
					version.NextReviewLevel = nil
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
			// Move to next level - set status to in_review and workflow status to submit
			inReviewStatus, err := icwrs.committeeStatusSvc.GetByType("in_review")
			if err != nil {
				tx.Rollback()
				return err
			}

			submitStatus, err := icwrs.workflowStatusSvc.GetByType("submit")
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
					version.WorkflowStatusID = &submitStatus.ID
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

	// Create revert review in both current and history tables
	if err := icwrs.createReviewWithHistory(currentCommitteeWorkflow.ID, revertStatus.ID, currentUserID, tx); err != nil {
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

	// 2. Reset level N-1 (previous level), level N (reverted level), and levels > N
	for _, workflow := range committeeWorkflowList.InternalCommitteeWorkflows {
		// Skip levels below previous level - they should remain approved
		if workflow.Level < previousLevel {
			continue
		}

		// Reset level N-1 (previous level), level N (reverted level), and all levels > N
		if workflow.Level >= previousLevel {

			// Delete old approval reviews from current table (history is preserved)
			var reviewsToDelete []*uint
			for _, review := range workflow.InternalCommitteeWorkflowReviews {
				if review.WorkflowStatus != nil && review.WorkflowStatus.Type == "approve" {
					reviewsToDelete = append(reviewsToDelete, &review.ID)
				}
			}

			if len(reviewsToDelete) > 0 {
				if err := icwrs.committeeWorkflowSvc.DeleteInternalCommitteeReviews(reviewsToDelete, tx); err != nil {
					tx.Rollback()
					return err
				}
			}

			// Update workflow status to pending (except current level which is already reverted)
			if workflow.Level != currentLevel {
				workflow.WorkflowStatusID = pendingStatus.ID
				if err := icwrs.committeeWorkflowSvc.UpdateInternalCommitteeWorkflow(workflow, tx); err != nil {
					tx.Rollback()
					return err
				}
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
				version.NextReviewLevel = nil
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

	// Fetch fresh committee data after commit to get updated NextReviewLevel
	updatedCommitteeDetail, err := icwrs.committeeSvc.GetFullDetailsByID(committeeID)
	if err != nil {
		utility.LogError("Failed to fetch updated committee details for notification", map[string]interface{}{
			"error": err.Error(),
		})
		return nil
	}

	// Notifications should be sent AFTER commit succeeds
	icwrs.sendRevertNotification(committeeID, currentUserID, updatedCommitteeDetail)

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
		redirectURL := fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeDetail.ID)

		if committeeDetail.LatestVersion.CommitteeStatus.Type == "approved" {
			redirectURL = fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeDetail.ID)
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
			redirectURL := fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeDetail.ID)

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

		// Get previous level (NextReviewLevel) participants to notify
		var notificationUserIDs []uint
		var emailRecipients []string

		if committee.LatestVersion.NextReviewLevel != nil {
			// Get the workflow for the previous level that needs re-approval
			previousLevelWorkflow, err := icwrs.committeeWorkflowSvc.GetInternalCommitteeWorkflowByLevel(
				committee.LatestVersion.ID,
				*committee.LatestVersion.NextReviewLevel,
			)
			if err == nil && previousLevelWorkflow != nil {
				// Collect participant user IDs and emails
				for _, participant := range previousLevelWorkflow.InternalCommitteeWorkflowParticipants {
					notificationUserIDs = append(notificationUserIDs, participant.UserID)
					if participant.User != nil {
						emailRecipients = append(emailRecipients, participant.User.Email)
					}
				}
			}
		}

		// If no participants found, fallback to created by user
		if len(notificationUserIDs) == 0 {
			createdByUserID := uint(committee.CreatedByUser.ID)
			notificationUserIDs = []uint{createdByUserID}
			emailRecipients = []string{committee.CreatedByUser.Email}
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              notificationUserIDs,
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

		appURL := utility.GetAppURL()
		subject := fmt.Sprintf("Internal Committee %s: Rejected Notification", committee.LatestVersion.CommitteeNameEn)
		redirectURL := fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committee.ID)

		// Send individual emails to each recipient with their name
		for _, userID := range notificationUserIDs {
			recipientUser, err := icwrs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch recipient user", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
				continue
			}

			emailData := map[string]interface{}{
				"FirstNameEn":           recipientUser.FirstNameEn,
				"LastNameEn":            recipientUser.LastNameEn,
				"Title":                 committee.LatestVersion.CommitteeNameEn,
				"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
				"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
				"RejectedAt":            time.Now().Format("2006-01-02"),
				"RedirectURL":           redirectURL,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{recipientUser.Email},
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/committee_workflow_revert_user_mail.html",
				EmailData:            emailData,
			}

			err = icwrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
			if err != nil {
				utility.LogError("Failed to send rejected notification email", map[string]interface{}{
					"error": err.Error(),
					"email": recipientUser.Email,
				})
			}
		}
	}()
}

// Helper function to create review in both current and history tables
func (icwrs *InternalCommitteeWorkflowReviewService) createReviewWithHistory(workflowID uint, statusID uint, userID uint, tx utility.DatabaseTransaction) error {
	// Create in current reviews table
	committeeReview := &models.InternalCommitteeWorkflowReview{
		InternalCommitteeWorkflowID: workflowID,
		WorkflowStatusID:            statusID,
		CreatedBYID:                 userID,
	}
	if err := icwrs.repo.CreateInternalCommitteeWorkflowReview(committeeReview, tx); err != nil {
		return err
	}

	// Create in history table for audit trail
	historyReview := &models.InternalCommitteeWorkflowReviewHistory{
		InternalCommitteeWorkflowID: workflowID,
		WorkflowStatusID:            statusID,
		CreatedBYID:                 userID,
		CreatedAt:                   time.Now(),
	}
	if err := icwrs.repo.CreateInternalCommitteeWorkflowReviewHistory(historyReview, tx); err != nil {
		return err
	}

	return nil
}
