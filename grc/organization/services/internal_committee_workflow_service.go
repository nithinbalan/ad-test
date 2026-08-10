package services

import (
	"fmt"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"sort"
	"time"
)

type InternalCommitteeWorkflowServiceInterface interface {
	SetInternalCommitteeService(svc InternalCommitteeServiceInterface)
	GetInternalCommitteeWorkflowList(internalCommitteeID uint) (*organizationDTOs.InternalCommitteeWorkflowListDTO, error)
	GetPreviousVersionWorkflowHistory(internalCommitteeID uint) ([]*organizationDTOs.InternalCommitteeWorkflowHistoryDTO, error)
	GetWorkflowListData(internalCommitteeID uint, latestVersionID uint, tx ...utility.DatabaseTransaction) (*models.InternalCommittee, error)
	CreateInternalCommitteeWorkflow(internalCommitteeWorkflowID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateInternalCommitteeWorkflow(internalCommitteeWorkflow *models.InternalCommitteeWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, internalCommitteeWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetInternalCommitteeWorkflowByLevel(internalCommitteeWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error)
	GetInternalCommitteeWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error)
	GetNextLevel(internalCommitteeWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.InternalCommitteeWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	UpdateNextLevelParticipant(internalCommitteeID uint, currentUserID uint, dto *organizationDTOs.UpdateNextLevelParticipantDTO, tx ...utility.DatabaseTransaction) error

	DeleteInternalCommitteeWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteInternalCommitteeReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteInternalCommitteeWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
}

type InternalCommitteeWorkflowService struct {
	repo                 repositories.InternalCommitteeWorkflowRepositoryInterface
	activityLogSvc       logSvc.ActivityLogServiceInterface
	workflowSvc          settinsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              UserServiceInterface
	committeeSvc         InternalCommitteeServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	organizationUnitsSvc OrganizationUnitsServiceInterface
	settingsOrgUnitSvc   settinsSvc.OrganizationUnitServiceInterface
	roleSvc              settinsSvc.RoleServiceInterface
}

func NewInternalCommitteeWorkflowService(
	repo repositories.InternalCommitteeWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc UserServiceInterface,
	committeeSvc InternalCommitteeServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	organizationUnitsSvc OrganizationUnitsServiceInterface,
	settingsOrgUnitSvc settinsSvc.OrganizationUnitServiceInterface,
	roleSvc settinsSvc.RoleServiceInterface,
) InternalCommitteeWorkflowServiceInterface {
	return &InternalCommitteeWorkflowService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		workflowSvc:          workflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		userSvc:              userSvc,
		committeeSvc:         committeeSvc,
		workflowContextSvc:   workflowContextSvc,
		organizationUnitsSvc: organizationUnitsSvc,
		settingsOrgUnitSvc:   settingsOrgUnitSvc,
		roleSvc:              roleSvc,
	}
}

func (icws *InternalCommitteeWorkflowService) SetInternalCommitteeService(svc InternalCommitteeServiceInterface) {
	icws.committeeSvc = svc
}

func (icws *InternalCommitteeWorkflowService) GetInternalCommitteeWorkflowList(internalCommitteeID uint) (*organizationDTOs.InternalCommitteeWorkflowListDTO, error) {

	committee, err := icws.committeeSvc.GetLatestVersionByID(internalCommitteeID)
	if err != nil {
		return nil, err
	}

	// Get latest version ID
	var latestVersionID uint
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			break
		}
	}

	if latestVersionID == 0 {
		return nil, fmt.Errorf("no latest version found for committee")
	}

	committeeWorkflow, err := icws.repo.GetWorkflowList(committee.ID, latestVersionID)
	if err != nil {
		return nil, err
	}

	dto := &organizationDTOs.InternalCommitteeWorkflowListDTO{
		ID:                         committeeWorkflow.ID,
		InternalCommitteeID:        committeeWorkflow.ID,
		InternalCommitteeVersionID: 0,
		NextReviewLevel:            nil,
		SubmittedAt:                nil,
		SubmittedByUser:            nil,
		Workflows:                  []*organizationDTOs.WorkflowListDTO{},
	}

	// Get latest version to get submitted info
	if len(committeeWorkflow.InternalCommitteeVersions) > 0 {
		for _, version := range committeeWorkflow.InternalCommitteeVersions {
			if version.IsLatestVersion != nil && *version.IsLatestVersion {
				dto.InternalCommitteeVersionID = version.ID
				dto.NextReviewLevel = version.NextReviewLevel
				dto.SubmittedAt = version.SubmittedAt
				dto.SubmittedByUser = icws.mapUserToDTO(version.SubmittedByUser)
				break
			}
		}
	}

	// Calculate total levels to determine second last level
	var maxLevel int
	for _, workflow := range committeeWorkflow.InternalCommitteeWorkflows {
		if workflow.Level > maxLevel {
			maxLevel = workflow.Level
		}
	}

	for _, workflow := range committeeWorkflow.InternalCommitteeWorkflows {
		workflowDTO := &organizationDTOs.WorkflowListDTO{
			Level:             workflow.Level,
			ApprovalType:      workflow.ApprovalTypes,
			IsSecondLastLevel: workflow.Level == maxLevel-1,
			CurrentStatus:     icws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers:     []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:       []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRevertedAt time.Time

		// Determine if this level is currently under review or needs re-approval
		isCurrentLevel := dto.NextReviewLevel != nil && *dto.NextReviewLevel == workflow.Level

		// Check if this level needs re-approval using tagged switch
		// When Level N is reverted, NextReviewLevel becomes N-1
		// We need to clear: Level N-1 (current review level) and Level N (reverted level)
		needsReapproval := false
		if dto.NextReviewLevel != nil {
			switch workflow.Level {
			case *dto.NextReviewLevel:
				// Current review level needs reapproval
				needsReapproval = true
			case *dto.NextReviewLevel + 1:
				// Reverted level needs reapproval
				needsReapproval = true
			}
		}

		// Only show approved_users and rejected_user if workflow status matches
		// This allows us to keep review history while showing correct current state
		if workflow.WorkflowStatus != nil {
			for _, review := range workflow.InternalCommitteeWorkflowReviews {
				userDTO := icws.mapUserToDTO(review.CreatedByUser)

				// Show approved users based on workflow status and current review level
				if review.WorkflowStatus.Type == "approve" {
					// Show approved users if workflow is fully approved and doesn't need re-approval
					if workflow.WorkflowStatus.Type == "approve" && !needsReapproval {
						workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
						if lastApprovedAt.IsZero() || review.CreatedAt.After(lastApprovedAt) {
							lastApprovedAt = review.CreatedAt
						}
					} else if isCurrentLevel && (workflow.WorkflowStatus.Type == "pending" || workflow.WorkflowStatus.Type == "revert") {
						// Show partial approvals for current level when pending or reverted
						// Old approval reviews are deleted during revert, so only new approvals will show
						workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
						if lastApprovedAt.IsZero() || review.CreatedAt.After(lastApprovedAt) {
							lastApprovedAt = review.CreatedAt
						}
					}
				}

				// Show rejected user only if this level's workflow status is "revert"
				if review.WorkflowStatus.Type == "revert" && workflow.WorkflowStatus.Type == "revert" {
					workflowDTO.RejectedUser = userDTO
					if lastRevertedAt.IsZero() || review.CreatedAt.After(lastRevertedAt) {
						lastRevertedAt = review.CreatedAt
					}
				}
			}
		}

		for _, participant := range workflow.InternalCommitteeWorkflowParticipants {
			userDTO := icws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRevertedAt
		dto.Workflows = append(dto.Workflows, workflowDTO)
	}

	dto.History = icws.buildWorkflowHistory(committeeWorkflow)
	return dto, nil
}

func (icws *InternalCommitteeWorkflowService) GetPreviousVersionWorkflowHistory(internalCommitteeID uint) ([]*organizationDTOs.InternalCommitteeWorkflowHistoryDTO, error) {
	committee, err := icws.committeeSvc.GetLatestVersionByID(internalCommitteeID)
	if err != nil {
		return nil, err
	}

	// Get latest version ID
	var latestVersionID uint
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			break
		}
	}

	if latestVersionID == 0 {
		return nil, fmt.Errorf("no latest version found for committee")
	}

	// Get previous version workflows from repository
	previousWorkflows, err := icws.repo.GetWorkflowHistoryFromHistoryTable(internalCommitteeID, latestVersionID)
	if err != nil {
		return nil, err
	}

	// Build history from previous version workflows
	var history []*organizationDTOs.InternalCommitteeWorkflowHistoryDTO
	for _, workflow := range previousWorkflows {
		createdAt := workflow.CreatedAt
		history = append(history, &organizationDTOs.InternalCommitteeWorkflowHistoryDTO{
			ID:             workflow.ID,
			Level:          workflow.InternalCommitteeWorkflow.Level,
			WorkflowStatus: icws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			CreatedAt:      createdAt,
			CreatedByUser:  icws.mapUserToDTO(workflow.CreatedByUser),
		})
	}

	// Sort history by created_at in ascending order
	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.Before(history[j].CreatedAt)
	})

	return history, nil
}

func (icws *InternalCommitteeWorkflowService) CreateInternalCommitteeWorkflow(internalCommitteeID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	internalCommitee, err := icws.committeeSvc.GetFullDetailsByID(internalCommitteeID)
	if err != nil {
		return err
	}
	workflowContext, err := icws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := icws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := icws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := icws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	// Get latest version ID
	committee, err := icws.committeeSvc.GetLatestVersionByID(internalCommitteeID)
	if err != nil {
		return err
	}

	var latestVersionID uint
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			break
		}
	}

	if latestVersionID == 0 {
		return fmt.Errorf("no latest version found for committee")
	}

	for _, level := range workflowLevels.WorkflowLevels {
		committeeWorkflow := &models.InternalCommitteeWorkflow{
			InternalCommitteeID:        internalCommitteeID,
			InternalCommitteeVersionID: latestVersionID,
			Level:                      level.Level,
			WorkflowLevelTypeID:        level.WorkflowLevelType.ID,
			ApprovalTypes:              level.ApprovalType,
			WorkflowStatusID:           workflowStatus.ID,
			CreatedBYID:                currentUserID,
		}

		workflowID, err := icws.repo.CreateInternalCommitteeWorkflow(committeeWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := icws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = icws.createWorkflowParticipant(userId, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				err = icws.createWorkflowParticipant(userId, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.WorkflowCustomRole != nil {

				switch participant.WorkflowCustomRole.Type {
				case "department_head":
					if internalCommitee.RequestingEntityID == 0 {
						return fmt.Errorf("Requesting department not found for this committee")
					}
					deptHead, err := icws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(internalCommitee.RequestingEntityID)
					if err != nil {
						return err
					}
					if deptHead.HeadUser == nil {
						return utility.NewForbiddenError("Department head is not assigned for the requesting department")
					}
					userID := uint(deptHead.HeadUser.ID)
					if err = icws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}

				}
			} else if participant.Unit != nil {

				unitDetails, err := icws.settingsOrgUnitSvc.GetFullDetailsByID(uint(participant.Unit.ID))
				if err != nil {
					return err
				}

				deptHead, err := icws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}

				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(
						fmt.Sprintf(
							"No head has been assigned to the %s '%s'. Please assign a head and try again.",
							unitDetails.OrganizationLevel.TitleEn,
							unitDetails.TitleEn,
						),
					)
				}

				userID := uint(deptHead.HeadUser.ID)

				if err = icws.createWorkflowParticipant(
					userID,
					workflowID,
					&participant,
					currentUserID,
					tx...,
				); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (icws *InternalCommitteeWorkflowService) UpdateInternalCommitteeWorkflow(committeeWorkflow *models.InternalCommitteeWorkflow, tx ...utility.DatabaseTransaction) error {
	err := icws.repo.UpdateInternalCommitteeWorkflow(committeeWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (icws *InternalCommitteeWorkflowService) createWorkflowParticipant(userID uint, committeeWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	committeeWorkflowParticipant := &models.InternalCommitteeWorkflowParticipant{
		InternalCommitteeWorkflowID: committeeWorkflowID,
		UserID:                      userID,
		CreatedBYID:                 currentUserID,
	}

	if participant.Role != nil {
		committeeWorkflowParticipant.RoleID = &participant.Role.ID
	}
	if participant.WorkflowCustomRole != nil {
		committeeWorkflowParticipant.CustomRoleID = &participant.WorkflowCustomRole.ID
	}
	// if participant.Unit != nil {
	// 	unitID := uint(participant.Unit.ID)
	// 	committeeWorkflowParticipant.UnitID = &unitID
	// }
	err := icws.repo.CreateInternalCommitteeWorkflowParticipant(committeeWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (icws *InternalCommitteeWorkflowService) GetInternalCommitteeWorkflowByLevel(versionID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error) {
	committeeWorkflow, err := icws.repo.GetInternalCommitteeWorkflowByLevel(versionID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return committeeWorkflow, nil
}

func (icws *InternalCommitteeWorkflowService) GetInternalCommitteeWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error) {
	committeeWorkflow, err := icws.repo.GetInternalCommitteeWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return committeeWorkflow, nil
}

func (icws *InternalCommitteeWorkflowService) GetNextLevel(committeeID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeWorkflow, error) {
	nextWorkflow, err := icws.repo.GetNextLevel(committeeID, currentLevel, tx...)
	if err != nil {
		return nil, err // Handle or log error as needed
	}
	return nextWorkflow, nil
}

func (icws *InternalCommitteeWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.InternalCommitteeWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := icws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (icws *InternalCommitteeWorkflowService) DeleteInternalCommitteeWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(participantIds) == 0 {
		return nil
	}

	err := icws.repo.DeleteInternalCommitteeWorkflowParticipants(participantIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete internal committee workflow participants: %w", err)
	}

	return nil
}

func (icws *InternalCommitteeWorkflowService) DeleteInternalCommitteeReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := icws.repo.DeleteInternalCommitteeWorkflowReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete internal committee workflow reviews: %w", err)
	}

	return nil
}

func (icws *InternalCommitteeWorkflowService) DeleteInternalCommitteeWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(workflowIds) == 0 {
		return nil
	}

	err := icws.repo.DeleteInternalCommitteeWorkflows(workflowIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete internal committee workflows: %w", err)
	}

	return nil
}

func (icws *InternalCommitteeWorkflowService) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := icws.repo.DeleteInternalCommitteeWorkflowReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete workflow reviews: %w", err)
	}

	return nil
}

func (icws *InternalCommitteeWorkflowService) GetWorkflowListData(internalCommitteeID uint, latestVersionID uint, tx ...utility.DatabaseTransaction) (*models.InternalCommittee, error) {
	committeeWorkflow, err := icws.repo.GetWorkflowList(internalCommitteeID, latestVersionID, tx...)
	if err != nil {
		return nil, err
	}
	return committeeWorkflow, nil
}

func (icws *InternalCommitteeWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
	if user == nil {
		return &organizationDTOs.UserMinimalDTO{}
	}
	return &organizationDTOs.UserMinimalDTO{
		ID:               int(user.ID),
		FirstNameEN:      user.FirstNameEn,
		LastNameEN:       user.LastNameEn,
		FirstNameAR:      user.FirstNameAr,
		LastNameAR:       user.LastNameAr,
		ProfileImagePath: user.ProfileImagePath,
		Email:            user.Email,
		Initial:          user.Initial,
		Color:            user.Color,
	}
}

func (icws *InternalCommitteeWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
	if status == nil {
		return masterDTOs.WorkflowStatusDTO{}
	}
	return masterDTOs.WorkflowStatusDTO{
		ID:        status.ID,
		TitleEn:   status.TitleEn,
		TitleAr:   status.TitleAr,
		ColorCode: status.ColorCode,
		Status:    status.Status,
		Type:      status.Type,
	}
}

func (icws *InternalCommitteeWorkflowService) buildWorkflowHistory(committee *models.InternalCommittee) []*organizationDTOs.InternalCommitteeWorkflowHistoryDTO {
	var history []*organizationDTOs.InternalCommitteeWorkflowHistoryDTO

	// Get latest version ID
	var latestVersionID uint
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			break
		}
	}

	if latestVersionID == 0 {
		return history
	}

	// Fetch complete history from history table
	historyRecords, err := icws.repo.GetWorkflowHistoryFromHistoryTable(committee.ID, latestVersionID)
	if err != nil {
		return history
	}

	// Build history DTOs from history records
	for _, record := range historyRecords {
		history = append(history, &organizationDTOs.InternalCommitteeWorkflowHistoryDTO{
			ID:             record.ID,
			Level:          record.InternalCommitteeWorkflow.Level,
			WorkflowStatus: icws.mapWorkflowStatusToDTO(record.WorkflowStatus),
			CreatedAt:      record.CreatedAt,
			CreatedByUser:  icws.mapUserToDTO(record.CreatedByUser),
		})
	}

	return history
}

func (icws *InternalCommitteeWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := icws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := icws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}
	if workflow == nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	workflowLevels, err := icws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}
	if len(workflowLevels.WorkflowLevels) == 0 {
		return utility.NewForbiddenError("workflow_levels_not_found")
	}

	hasParticipants := false
	for _, level := range workflowLevels.WorkflowLevels {
		if len(level.Participants) > 0 {
			hasParticipants = true
			break
		}
	}

	if !hasParticipants {
		return utility.NewForbiddenError("participants_not_configured")
	}

	return nil
}

func (icws *InternalCommitteeWorkflowService) UpdateNextLevelParticipant(internalCommitteeID uint, currentUserID uint, dto *organizationDTOs.UpdateNextLevelParticipantDTO, tx ...utility.DatabaseTransaction) error {
	// Get latest committee version
	committee, err := icws.committeeSvc.GetLatestVersionByID(internalCommitteeID)
	if err != nil {
		return err
	}

	var latestVersionID uint
	for _, version := range committee.InternalCommitteeVersions {
		if version.IsLatestVersion != nil && *version.IsLatestVersion {
			latestVersionID = version.ID
			break
		}
	}

	if latestVersionID == 0 {
		return fmt.Errorf("no latest version found for committee")
	}

	// Validate current level exists
	_, err = icws.repo.GetInternalCommitteeWorkflowByLevel(latestVersionID, dto.CurrentLevel, tx...)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("current level %d not found", dto.CurrentLevel), "current_level_not_found", "")
	}

	// Calculate target level
	targetLevel := dto.CurrentLevel + 1

	// Get target workflow
	targetWorkflow, err := icws.repo.GetInternalCommitteeWorkflowByLevel(latestVersionID, targetLevel, tx...)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("target level %d not found", targetLevel), "target_level_not_found", "")
	}

	// Validate target level is not already approved
	workflowStatus, err := icws.workflowStatusSvc.GetByID(targetWorkflow.WorkflowStatusID)
	if err != nil {
		return err
	}
	if workflowStatus.Type == "approve" {
		return utility.NewValidationError(fmt.Errorf("target level %d is already approved", targetLevel), "target_level_already_approved", "")
	}

	// Validate role ID is allowed (CEO or Shared Services)
	role, err := icws.roleSvc.GetByID(dto.RoleID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("role not found"), "role_not_found", "")
	}

	// Check if role is CEO or Shared Services
	// Assuming role titles are stored in TitleEn field
	allowedRoles := map[string]bool{
		"ceo":             true,
		"shared-services": true,
	}

	if !allowedRoles[*role.Type] {
		return utility.NewValidationError(fmt.Errorf("role %s is not allowed for participant update", role.TitleEn), "role_not_allowed", "")
	}

	// Delete existing participants for target level
	var existingParticipantIds []*uint
	for _, participant := range targetWorkflow.InternalCommitteeWorkflowParticipants {
		participantId := participant.ID
		existingParticipantIds = append(existingParticipantIds, &participantId)
	}

	if len(existingParticipantIds) > 0 {
		if err := icws.repo.DeleteInternalCommitteeWorkflowParticipants(existingParticipantIds, tx...); err != nil {
			return err
		}
	}

	// Get users by role ID to create individual participants
	userIds, err := icws.userSvc.GetUserIdsByRole(dto.RoleID, tx...)
	if err != nil {
		return err
	}

	if len(userIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no users found for role %s", role.TitleEn), "no_users_for_role", "")
	}

	// Create participants for each user with the role
	for _, userId := range userIds {
		committeeWorkflowParticipant := &models.InternalCommitteeWorkflowParticipant{
			InternalCommitteeWorkflowID: targetWorkflow.ID,
			UserID:                      userId,
			RoleID:                      &dto.RoleID,
			CreatedBYID:                 currentUserID,
		}
		if err := icws.repo.CreateInternalCommitteeWorkflowParticipant(committeeWorkflowParticipant, tx...); err != nil {
			return err
		}
	}

	// Log activity asynchronously
	go func() {
		err := icws.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.InternalCommitteeWorkflow{}.TableName(),
			targetWorkflow.ID,
			targetWorkflow,
			currentUserID,
			"organization",
			"committee_management",
			"internal_committee_workflow",
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.InternalCommitteeWorkflow{}.TableName(),
				"recordID":  targetWorkflow.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
