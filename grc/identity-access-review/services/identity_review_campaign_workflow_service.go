package services

import (
	"fmt"

	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	masterDTOs "jethur/grc/masters/dtos"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
)

type IdentityReviewCampaignWorkflowServiceInterface interface {
	GetWorkflowList(campaignID uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewCampaignWorkflowListDTO, error)
	ValidateWorkflowBeforeCreate(campaignID uint, tx ...utility.DatabaseTransaction) error
	CreateCampaignWorkflow(campaignID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error)
	GetWorkflowByLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error)
	GetNextLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, participants []*models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	UpdateWorkflow(workflow *models.IdentityReviewCampaignWorkflow, tx ...utility.DatabaseTransaction) error
	UpdateParticipant(participant *models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForCampaign(campaignID uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	ResetParticipantStatuses(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewCampaignWorkflowService struct {
	repo                repositories.IdentityReviewCampaignWorkflowRepositoryInterface
	workflowSvc         settingsServices.WorkflowServiceInterface
	workflowStatusSvc   masterServices.WorkflowStatusServiceInterface
	campaignSvc         IdentityReviewCampaignServiceInterface
	userSvc             organizationServices.UserServiceInterface
	organizationUnitSvc organizationServices.OrganizationUnitsServiceInterface
}

func NewIdentityReviewCampaignWorkflowService(
	repo repositories.IdentityReviewCampaignWorkflowRepositoryInterface,
	workflowSvc settingsServices.WorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	campaignSvc IdentityReviewCampaignServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc organizationServices.OrganizationUnitsServiceInterface,
) IdentityReviewCampaignWorkflowServiceInterface {
	return &IdentityReviewCampaignWorkflowService{
		repo:                repo,
		workflowSvc:         workflowSvc,
		workflowStatusSvc:   workflowStatusSvc,
		campaignSvc:         campaignSvc,
		userSvc:             userSvc,
		organizationUnitSvc: organizationUnitSvc,
	}
}

func (ircws *IdentityReviewCampaignWorkflowService) ValidateWorkflowBeforeCreate(campaignID uint, tx ...utility.DatabaseTransaction) error {
	campaign, err := ircws.campaignSvc.GetByID(campaignID, tx...)
	if err != nil {
		return err
	}

	if campaign.WorkflowID == 0 {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := ircws.workflowSvc.GetWorkflowLevels(
		campaign.WorkflowID,
		tx...,
	)
	if err != nil {
		return err
	}

	if len(workflowLevels.WorkflowLevels) == 0 {
		return utility.NewForbiddenError("workflow levels not found")
	}

	for _, level := range workflowLevels.WorkflowLevels {
		if len(level.Participants) > 0 {
			return nil
		}
	}

	return utility.NewForbiddenError("participants not configured")
}

func (ircws *IdentityReviewCampaignWorkflowService) CreateCampaignWorkflow(campaignID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	campaign, err := ircws.campaignSvc.GetByID(campaignID, tx...)
	if err != nil {
		return err
	}

	workflow, err := ircws.workflowSvc.GetWorkflowLevels(campaign.WorkflowID, tx...)
	if err != nil {
		return err
	}

	pending, err := ircws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflow.WorkflowLevels {
		campaignWorkflow := &models.IdentityReviewCampaignWorkflow{
			IdentityReviewCampaignID: campaignID,
			Level:                    level.Level,
			WorkflowLevelTypeID:      level.WorkflowLevelType.ID,
			ApprovalTypes:            level.ApprovalType,
			WorkflowStatusID:         pending.ID,
			CreatedByID:              currentUserID,
		}

		workflowID, err := ircws.repo.CreateWorkflow(campaignWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if err := ircws.createParticipants(campaign, workflowID, &participant, currentUserID, tx...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ircws *IdentityReviewCampaignWorkflowService) createParticipants(campaign *models.IdentityReviewCampaign, workflowID uint, participant *settingsDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	var userIDs []uint
	if participant.Role != nil {
		ids, err := ircws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
		if err != nil {
			return err
		}
		userIDs = ids
	} else if participant.User != nil {
		userIDs = []uint{uint(participant.User.ID)}
	} else if participant.WorkflowCustomRole != nil {
		userIDs = []uint{campaign.ReviewerID}
	} else if participant.Unit != nil {
		unit, err := ircws.organizationUnitSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
		if err != nil {
			return err
		}
		if unit.HeadUser == nil {
			return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
		}
		userIDs = []uint{uint(unit.HeadUser.ID)}
	}

	for _, userID := range userIDs {
		workflowParticipant := &models.IdentityReviewCampaignWorkflowParticipant{
			IdentityReviewCampaignWorkflowID: workflowID,
			UserID:                           userID,
			CreatedByID:                      currentUserID,
		}

		if participant.Role != nil {
			workflowParticipant.RoleID = &participant.Role.ID
		}

		if err := ircws.repo.CreateParticipant(workflowParticipant, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (ircws *IdentityReviewCampaignWorkflowService) GetWorkflowList(campaignID uint, tx ...utility.DatabaseTransaction) (*dtos.IdentityReviewCampaignWorkflowListDTO, error) {
	campaign, err := ircws.repo.GetWorkflowList(campaignID, tx...)
	if err != nil {
		return nil, err
	}

	overallWorkflowStatus, err := ircws.mapOverallWorkflowStatus(campaign)
	if err != nil {
		return nil, err
	}

	histories, err := ircws.repo.GetWorkflowHistories(campaignID, tx...)
	if err != nil {
		return nil, err
	}

	result := &dtos.IdentityReviewCampaignWorkflowListDTO{
		ID:              campaign.ID,
		CampaignID:      campaign.ID,
		SubmittedAt:     campaign.SubmittedAt,
		SubmittedByUser: ircws.mapUserToDTO(campaign.SubmittedByUser),
		NextReviewLevel: campaign.NextReviewLevel,
		CurrentStatus:   overallWorkflowStatus,
		Workflows:       []*dtos.IdentityReviewCampaignWorkflowLevelDTO{},
		History:         ircws.buildWorkflowHistory(histories),
	}

	maxLevel := 0
	for _, workflow := range campaign.CampaignWorkflows {
		if workflow.Level > maxLevel {
			maxLevel = workflow.Level
		}
	}

	for _, workflow := range campaign.CampaignWorkflows {
		workflowItem := &dtos.IdentityReviewCampaignWorkflowLevelDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: ircws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
		}

		for _, participant := range workflow.IdentityReviewCampaignWorkflowParticipants {
			workflowItem.ReviewUsers = append(workflowItem.ReviewUsers, ircws.mapUserToDTO(participant.User))
		}

		for _, review := range workflow.IdentityReviewCampaignWorkflowReviews {
			reviewUser := review.CreatedByUser
			if reviewUser == nil {
				for _, participant := range workflow.IdentityReviewCampaignWorkflowParticipants {
					if participant.UserID == review.CreatedByID {
						reviewUser = participant.User
						break
					}
				}
			}
			if review.WorkflowStatus != nil && review.WorkflowStatus.Type == "approve" {
				workflowItem.ApprovedUsers = append(workflowItem.ApprovedUsers, ircws.mapUserToDTO(reviewUser))
				t := review.CreatedAt
				if workflowItem.ApprovedAt == nil || t.After(*workflowItem.ApprovedAt) {
					workflowItem.ApprovedAt = &t
					workflowItem.ApprovedAtString = t.Format("02 Jan 2006 03:04 PM")
				}
				if workflow.Level == maxLevel && campaign.NextReviewLevel == nil {
					result.ApprovedAt = &t
					result.ApprovedByUser = ircws.mapUserToDTO(review.CreatedByUser)
				}
			}
			if review.WorkflowStatus != nil && review.WorkflowStatus.Type == "reject" {
				workflowItem.RejectedUser = ircws.mapUserToDTO(reviewUser)
				t := review.CreatedAt
				workflowItem.RejectedAt = &t
			}
		}

		result.Workflows = append(result.Workflows, workflowItem)
	}

	return result, nil
}

func (ircws *IdentityReviewCampaignWorkflowService) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error) {
	return ircws.repo.GetWorkflowByID(id, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) GetWorkflowByLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error) {
	return ircws.repo.GetWorkflowByLevel(campaignID, level, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) GetNextLevel(campaignID uint, level int, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignWorkflow, error) {
	return ircws.repo.GetNextLevel(campaignID, level, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) UpdateWorkflow(workflow *models.IdentityReviewCampaignWorkflow, tx ...utility.DatabaseTransaction) error {
	return ircws.repo.UpdateWorkflow(workflow, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) UpdateParticipant(participant *models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	return ircws.repo.UpdateParticipant(participant, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) DeleteAllWorkflowsForCampaign(campaignID uint, tx ...utility.DatabaseTransaction) error {
	return ircws.repo.DeleteAllWorkflowsForCampaign(campaignID, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return ircws.repo.DeleteReviewsByWorkflowID(workflowID, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) ResetParticipantStatuses(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return ircws.repo.ResetParticipantStatuses(workflowID, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	return ircws.repo.UpdateWorkflowStatus(workflowID, statusID, tx...)
}

func (ircws *IdentityReviewCampaignWorkflowService) CheckAllParticipantsApproved(workflowID uint, participants []*models.IdentityReviewCampaignWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := ircws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	if count < len(participants) {
		return false, nil
	}

	return true, nil
}

func (ircws *IdentityReviewCampaignWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
	if user == nil {
		return nil
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

func (ircws *IdentityReviewCampaignWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (ircws *IdentityReviewCampaignWorkflowService) mapOverallWorkflowStatus(campaign *models.IdentityReviewCampaign) (*masterDTOs.WorkflowStatusDTO, error) {
	if campaign.CampaignStatus == nil {
		return nil, nil
	}

	var statusType string
	switch campaign.CampaignStatus.Type {
	case "approved", "closed":
		statusType = "approve"
	case "rejected":
		statusType = "reject"
	case "submitted":
		if campaign.NextReviewLevel != nil && *campaign.NextReviewLevel == 1 {
			statusType = "submit"
		} else {
			statusType = "pending"
		}
	case "in_progress":
		statusType = "revert"
	default:
		return nil, nil
	}

	workflowStatus, err := ircws.workflowStatusSvc.GetByType(statusType)
	if err != nil {
		return nil, err
	}

	status := ircws.mapWorkflowStatusToDTO(workflowStatus)
	return &status, nil
}

func (ircws *IdentityReviewCampaignWorkflowService) buildWorkflowHistory(histories []*models.IdentityReviewCampaignWorkflowHistory) []*dtos.IdentityReviewCampaignWorkflowHistoryDTO {
	historyDTOs := make([]*dtos.IdentityReviewCampaignWorkflowHistoryDTO, 0, len(histories))

	for _, history := range histories {
		historyItem := &dtos.IdentityReviewCampaignWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: ircws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			CreatedByUser:  ircws.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
		}
		historyDTOs = append(historyDTOs, historyItem)
	}

	return historyDTOs
}
