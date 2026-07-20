package services

import (
	"fmt"

	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type AssetRiskWorkflowServiceInterface interface {
	GetRiskWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*assetRiskDTOs.AssetRiskWorkflowListDTO, error)
	CreateRiskWorkflow(riskRegisterID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetRiskWorkflowByLevel(riskWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetNextLevel(riskWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error
	UpdateRiskWorkflowParticipant(riskWorkflowParticipant *models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	GetWorkflowListData(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error

	DeleteRiskWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForRisk(riskID uint, tx ...utility.DatabaseTransaction) error

	DeleteRiskWorkflowReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	ResetRiskWorkflowParticipantsStatusByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskWorkflowService struct {
	repo                 repositories.AssetRiskWorkflowRepositoryInterface
	workflowSvc          settinsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	riskSvc              AssetRiskRegisterServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
	settingsOrgUnitSvc   settinsSvc.OrganizationUnitServiceInterface
}

func NewAssetRiskWorkflowService(
	repo repositories.AssetRiskWorkflowRepositoryInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	riskSvc AssetRiskRegisterServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
	settingsOrgUnitSvc settinsSvc.OrganizationUnitServiceInterface,
) AssetRiskWorkflowServiceInterface {
	return &AssetRiskWorkflowService{
		repo:                 repo,
		workflowSvc:          workflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		userSvc:              userSvc,
		workflowContextSvc:   workflowContextSvc,
		riskSvc:              riskSvc,
		organizationUnitsSvc: organizationUnitsSvc,
		settingsOrgUnitSvc:   settingsOrgUnitSvc,
	}
}

func (arws *AssetRiskWorkflowService) GetRiskWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*assetRiskDTOs.AssetRiskWorkflowListDTO, error) {
	riskWorkflow, err := arws.repo.GetWorkflowList(riskID, tx...)
	if err != nil {
		return nil, err
	}

	overallStatus, err := arws.mapOverallWorkflowStatus(riskWorkflow)
	if err != nil {
		return nil, err
	}

	dto := &assetRiskDTOs.AssetRiskWorkflowListDTO{
		ID:              riskWorkflow.ID,
		RiskRegisterID:  riskWorkflow.ID,
		SubmittedAt:     riskWorkflow.SubmittedAt,
		SubmittedByUser: arws.mapUserToDTO(riskWorkflow.SubmittedByUser),
		NextReviewLevel: riskWorkflow.NextReviewLevel,
		CurrentStatus:   overallStatus,
		RiskWorkflows:   []*assetRiskDTOs.WorkflowListDTO{},
	}

	var maxLevel int
	for _, workflow := range riskWorkflow.RiskWorkflows {
		if workflow.Level > maxLevel {
			maxLevel = workflow.Level
		}
	}

	// Final ApprovedAt/ApprovedByUser only apply once every level has been approved.
	if riskWorkflow.RiskStatus != nil && riskWorkflow.RiskStatus.Type == "open" {
		for _, workflow := range riskWorkflow.RiskWorkflows {
			if workflow.Level != maxLevel {
				continue
			}
			for _, review := range workflow.RiskWorkflowReviews {
				if review.WorkflowStatus.Type != "approve" {
					continue
				}
				if dto.ApprovedAt == nil || review.CreatedAt.After(*dto.ApprovedAt) {
					t := review.CreatedAt
					dto.ApprovedAt = &t
					dto.ApprovedByUser = arws.mapUserToDTO(review.CreatedByUser)
				}
			}
		}
	}

	for _, workflow := range riskWorkflow.RiskWorkflows {
		workflowDTO := &assetRiskDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: arws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt *time.Time
		var lastRejectedAt *time.Time
		for _, review := range workflow.RiskWorkflowReviews {
			userDTO := arws.mapUserToDTO(review.CreatedByUser)
			if review.WorkflowStatus.Type == "approve" {
				workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
				if lastApprovedAt == nil || review.CreatedAt.After(*lastApprovedAt) {
					t := review.CreatedAt
					lastApprovedAt = &t
				}
			}
			if review.WorkflowStatus.Type == "reject" {
				workflowDTO.RejectedUser = userDTO
				if lastRejectedAt == nil || review.CreatedAt.After(*lastRejectedAt) {
					t := review.CreatedAt
					lastRejectedAt = &t
				}
			}
		}

		for _, participant := range workflow.RiskWorkflowParticipants {
			userDTO := arws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.RiskWorkflows = append(dto.RiskWorkflows, workflowDTO)
	}

	return dto, nil
}

func (arws *AssetRiskWorkflowService) CreateRiskWorkflow(riskRegisterID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := arws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := arws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := arws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := arws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		riskWorkflow := &models.RiskWorkflow{
			RiskRegisterID:      riskRegisterID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		riskWorkflowID, err := arws.repo.CreateRiskWorkflow(riskWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIds, err := arws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					if err := arws.createWorkflowParticipant(userId, riskWorkflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				if err := arws.createWorkflowParticipant(userId, riskWorkflowID, &participant, currentUserID, tx...); err != nil {
					return err
				}
			} else if participant.WorkflowCustomRole != nil {
				switch participant.WorkflowCustomRole.Type {
				case "department_head":
					riskRegister, err := arws.riskSvc.GetByID(riskRegisterID, tx...)
					if err != nil {
						return err
					}
					if riskRegister.ResponsibleUnitID == nil {
						return fmt.Errorf("responsible unit not found for this risk")
					}
					deptHead, err := arws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(*riskRegister.ResponsibleUnitID)
					if err != nil {
						return err
					}
					if deptHead.HeadUser == nil {
						return utility.NewForbiddenError("Department head is not assigned for the responsible unit")
					}
					userId := uint(deptHead.HeadUser.ID)
					if err := arws.createWorkflowParticipant(userId, riskWorkflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}
				case "risk_responsible_user":
					riskRegister, err := arws.riskSvc.GetByID(riskRegisterID, tx...)
					if err != nil {
						return err
					}
					if riskRegister.ResponsibleUserID == nil {
						return fmt.Errorf("responsible user not found for this risk")
					}
					userId := *riskRegister.ResponsibleUserID
					if err := arws.createWorkflowParticipant(userId, riskWorkflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}
				}
			} else if participant.Unit != nil {
				unitDetails, err := arws.settingsOrgUnitSvc.GetFullDetailsByID(uint(participant.Unit.ID))
				if err != nil {
					return err
				}

				deptHead, err := arws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
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

				userId := uint(deptHead.HeadUser.ID)
				if err := arws.createWorkflowParticipant(userId, riskWorkflowID, &participant, currentUserID, tx...); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (arws *AssetRiskWorkflowService) GetWorkflowListData(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	return arws.repo.GetWorkflowList(riskID, tx...)
}

func (arws *AssetRiskWorkflowService) UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error {
	return arws.repo.UpdateRiskWorkflow(riskWorkflow, tx...)
}

func (arws *AssetRiskWorkflowService) UpdateRiskWorkflowParticipant(riskWorkflowParticipant *models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	return arws.repo.UpdateRiskWorkflowParticipant(riskWorkflowParticipant, tx...)
}

func (arws *AssetRiskWorkflowService) createWorkflowParticipant(userID uint, riskWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	riskWorkflowParticipant := &models.RiskWorkflowParticipant{
		RiskWorkflowID: riskWorkflowID,
		UserID:         userID,
		CreatedBYID:    currentUserID,
	}

	if participant.Role != nil {
		riskWorkflowParticipant.RoleID = &participant.Role.ID
	}

	return arws.repo.CreateRiskWorkflowParticipant(riskWorkflowParticipant, tx...)
}

func (arws *AssetRiskWorkflowService) GetRiskWorkflowByLevel(riskWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	return arws.repo.GetRiskWorkflowByLevel(riskWorkflowID, currentApprovalLevel, tx...)
}

func (arws *AssetRiskWorkflowService) GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	return arws.repo.GetRiskWorkflowByID(id, tx...)
}

func (arws *AssetRiskWorkflowService) GetNextLevel(riskWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	return arws.repo.GetNextLevel(riskWorkflowID, currentLevel, tx...)
}

func (arws *AssetRiskWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := arws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (arws *AssetRiskWorkflowService) DeleteRiskWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return arws.repo.DeleteRiskWorkflowReviews(reviewIds, tx...)
}

func (arws *AssetRiskWorkflowService) DeleteRiskWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return arws.repo.DeleteRiskWorkflowParticipants(participantIds, tx...)
}

func (arws *AssetRiskWorkflowService) DeleteRiskWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return arws.repo.DeleteRiskWorkflows(workflowIds, tx...)
}

func (arws *AssetRiskWorkflowService) DeleteRiskWorkflowReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return arws.repo.DeleteRiskWorkflowReviewsByWorkflowID(workflowID, tx...)
}

func (arws *AssetRiskWorkflowService) ResetRiskWorkflowParticipantsStatusByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return arws.repo.ResetRiskWorkflowParticipantsStatusByWorkflowID(workflowID, tx...)
}

func (arws *AssetRiskWorkflowService) UpdateRiskWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	return arws.repo.UpdateRiskWorkflowStatus(workflowID, statusID, tx...)
}

func (arws *AssetRiskWorkflowService) DeleteAllWorkflowsForRisk(riskID uint, tx ...utility.DatabaseTransaction) error {
	riskRegister, err := arws.repo.GetWorkflowsForDeletion(riskID, tx...)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range riskRegister.RiskWorkflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.RiskWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.RiskWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := arws.DeleteRiskWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := arws.DeleteRiskWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := arws.DeleteRiskWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (arws *AssetRiskWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (arws *AssetRiskWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (arws *AssetRiskWorkflowService) mapOverallWorkflowStatus(riskRegister *models.RiskRegister) (*masterDTOs.WorkflowStatusDTO, error) {
	if riskRegister.RiskStatus == nil {
		return nil, nil
	}

	var statusType string
	switch riskRegister.RiskStatus.Type {
	case "rejected":
		statusType = "reject"
	case "open":
		statusType = "approve"
	case "submitted":
		if riskRegister.NextReviewLevel != nil && *riskRegister.NextReviewLevel == 1 {
			statusType = "submit"
		} else {
			statusType = "pending"
		}
	case "resubmitted":
		if riskRegister.NextReviewLevel != nil && *riskRegister.NextReviewLevel == 1 {
			statusType = "resubmit"
		} else {
			statusType = "pending"
		}
	default:
		return nil, nil
	}

	workflowStatus, err := arws.workflowStatusSvc.GetByType(statusType)
	if err != nil {
		return nil, err
	}

	status := arws.mapWorkflowStatusToDTO(workflowStatus)
	return &status, nil
}

func (arws *AssetRiskWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := arws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := arws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}
	if workflow == nil {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := arws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}
	if len(workflowLevels.WorkflowLevels) == 0 {
		return utility.NewForbiddenError("workflow levels not found")
	}

	hasParticipants := false
	for _, level := range workflowLevels.WorkflowLevels {
		if len(level.Participants) > 0 {
			hasParticipants = true
			break
		}
	}

	if !hasParticipants {
		return utility.NewForbiddenError("participants not configured")
	}

	return nil
}
