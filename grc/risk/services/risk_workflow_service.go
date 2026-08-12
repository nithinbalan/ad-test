package services

import (
	"fmt"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type RiskWorkflowServiceInterface interface {
	SetRiskService(svc RiskRegisterServiceInterface)
	GetRiskWorkflowList(riskID uint) (*riskDTOs.RiskWorkflowListDTO, error)
	CreateRiskWorkflow(riskWorkflowID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, riskWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetRiskWorkflowByLevel(riskWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetNextLevel(riskWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteRiskWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
}

type RiskWorkflowService struct {
	repo                 repositories.RiskWorkflowRepositoryInterface
	activityLogSvc       logSvc.ActivityLogServiceInterface
	workflowSvc          settinsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	riskRegisterSvc      RiskRegisterServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
}

func NewRiskWorkflowService(
	repo repositories.RiskWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) RiskWorkflowServiceInterface {
	return &RiskWorkflowService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		workflowSvc:          workflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		userSvc:              userSvc,
		workflowContextSvc:   workflowContextSvc,
		organizationUnitsSvc: organizationUnitsSvc,
	}
}

func (arws *RiskWorkflowService) SetRiskService(svc RiskRegisterServiceInterface) {
	arws.riskRegisterSvc = svc
}

func (arws *RiskWorkflowService) GetRiskWorkflowList(riskID uint) (*riskDTOs.RiskWorkflowListDTO, error) {

	risk, err := arws.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		return nil, err
	}

	riskWorkflow, err := arws.repo.GetWorkflowList(risk.ID)
	if err != nil {
		return nil, err
	}

	dto := &riskDTOs.RiskWorkflowListDTO{
		ID:              riskWorkflow.ID,
		RiskRegisterID:  riskWorkflow.ID,
		SubmittedAt:     riskWorkflow.SubmittedAt,
		SubmittedByUser: arws.mapUserToDTO(riskWorkflow.SubmittedByUser),
		RiskWorkflows:   []*riskDTOs.WorkflowListDTO{},
	}

	for _, workflow := range riskWorkflow.RiskWorkflows {
		workflowDTO := &riskDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: arws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.RiskWorkflowReviews {
			userDTO := arws.mapUserToDTO(review.CreatedByUser)
			if review.WorkflowStatus.Type == "approve" {
				workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
				if review.CreatedAt.After(lastApprovedAt) {
					lastApprovedAt = review.CreatedAt
				}
			}
			if review.WorkflowStatus.Type == "reject" {
				workflowDTO.RejectedUser = userDTO
				if review.CreatedAt.After(lastRejectedAt) {
					lastRejectedAt = review.CreatedAt
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

func (arws *RiskWorkflowService) CreateRiskWorkflow(riskWorkflowID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {

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
			RiskRegisterID:      riskWorkflowID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		docWorkflowId, err := arws.repo.CreateRiskWorkflow(riskWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := arws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = arws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				err = arws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.Unit != nil {
				deptHead, err := arws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}
				userId := uint(deptHead.HeadUser.ID)
				err = arws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (arws *RiskWorkflowService) UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error {
	err := arws.repo.UpdateRiskWorkflow(riskWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *RiskWorkflowService) createWorkflowParticipant(userID uint, riskWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	riskWorkflowParticipant := &models.RiskWorkflowParticipant{
		RiskWorkflowID: riskWorkflowID,
		UserID:         userID,
		CreatedBYID:    currentUserID,
	}

	if participant.Role != nil {
		riskWorkflowParticipant.RoleID = &participant.Role.ID
	}

	err := arws.repo.CreateRiskWorkflowParticipant(riskWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *RiskWorkflowService) GetRiskWorkflowByLevel(riskWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	riskWorkflow, err := arws.repo.GetRiskWorkflowByLevel(riskWorkflowID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return riskWorkflow, nil
}

func (arws *RiskWorkflowService) GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	riskWorkflow, err := arws.repo.GetRiskWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return riskWorkflow, nil
}

func (arws *RiskWorkflowService) GetNextLevel(riskWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	nextWorkflow, err := arws.repo.GetNextLevel(riskWorkflowID, currentLevel, tx...)
	if err != nil {
		return nil, err // Handle or log error as needed
	}
	return nextWorkflow, nil
}

func (arws *RiskWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := arws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (arws *RiskWorkflowService) DeleteRiskWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(participantIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteRiskWorkflowParticipants(participantIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete risk workflow participants: %w", err)
	}

	return nil
}

func (arws *RiskWorkflowService) DeleteRiskReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteRiskWorkflowReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete risk workflow reviews: %w", err)
	}

	return nil
}

func (arws *RiskWorkflowService) DeleteRiskWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(workflowIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteRiskWorkflows(workflowIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete risk workflows: %w", err)
	}

	return nil
}

func (arws *RiskWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (arws *RiskWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (arws *RiskWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := arws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := arws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}
	if workflow == nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	workflowLevels, err := arws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
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
