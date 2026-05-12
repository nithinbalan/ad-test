package services

import (
	commonServices "jethur/grc/common/services"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settingsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type RiskRegisterUpdatePlanWorkflowServiceInterface interface {
	SetRiskRegisterUpdatePlanService(svc RiskRegisterUpdatePlanServiceInterface)
	GetRiskRegisterUpdatePlanWorkflowList(riskRegisterUpdatePlanID uint) (*riskDTOs.RiskRegisterUpdatePlanWorkflowListDTO, error)
	CreateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, riskRegisterUpdatePlanWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error)
	GetRiskRegisterUpdatePlanWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error)
	GetNextLevel(riskRegisterUpdatePlanID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskRegisterUpdatePlanWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteRiskRegisterUpdatePlanWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskRegisterUpdatePlanWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskRegisterUpdatePlanWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForRiskRegisterUpdatePlan(riskRegisterUpdatePlanID uint, tx ...utility.DatabaseTransaction) error
}

type RiskRegisterUpdatePlanWorkflowService struct {
	repo                      repositories.RiskRegisterUpdatePlanWorkflowRepositoryInterface
	activityLogSvc            logSvc.ActivityLogServiceInterface
	workflowSvc               settingsSvc.WorkflowServiceInterface
	workflowStatusSvc         masterSvc.WorkflowStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	riskRegisterUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface
	workflowContextSvc        masterSvc.WorkflowContextServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
}

func NewRiskRegisterUpdatePlanWorkflowService(
	repo repositories.RiskRegisterUpdatePlanWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) RiskRegisterUpdatePlanWorkflowServiceInterface {
	return &RiskRegisterUpdatePlanWorkflowService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		workflowSvc:        workflowSvc,
		workflowStatusSvc:  workflowStatusSvc,
		userSvc:            userSvc,
		workflowContextSvc: workflowContextSvc,
		notificationSvc:    notificationSvc,
	}
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) SetRiskRegisterUpdatePlanService(svc RiskRegisterUpdatePlanServiceInterface) {
	rrupws.riskRegisterUpdatePlanSvc = svc
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) GetRiskRegisterUpdatePlanWorkflowList(riskRegisterUpdatePlanID uint) (*riskDTOs.RiskRegisterUpdatePlanWorkflowListDTO, error) {
	riskRegisterUpdatePlan, err := rrupws.riskRegisterUpdatePlanSvc.GetByID(riskRegisterUpdatePlanID)
	if err != nil {
		return nil, err
	}

	riskRegisterUpdatePlanWorkflow, err := rrupws.repo.GetWorkflowList(riskRegisterUpdatePlan.ID)
	if err != nil {
		return nil, err
	}

	// Get workflow histories
	histories, err := rrupws.repo.GetWorkflowHistories(riskRegisterUpdatePlan.ID)
	if err != nil {
		return nil, err
	}

	dto := &riskDTOs.RiskRegisterUpdatePlanWorkflowListDTO{
		ID:                              riskRegisterUpdatePlanWorkflow.ID,
		RiskRegisterUpdatePlanID:        riskRegisterUpdatePlanWorkflow.ID,
		NextReviewLevel:                 riskRegisterUpdatePlanWorkflow.NextReviewLevel,
		SubmittedAt:                     riskRegisterUpdatePlanWorkflow.SubmittedAt,
		SubmittedByUser:                 rrupws.mapUserToDTO(riskRegisterUpdatePlanWorkflow.SubmittedByUser),
		RiskRegisterUpdatePlanWorkflows: []*riskDTOs.WorkflowListDTO{},
		History:                         []*riskDTOs.RiskRegisterUpdatePlanWorkflowHistoryDTO{},
	}

	for _, workflow := range riskRegisterUpdatePlanWorkflow.RiskRegisterUpdatePlanWorkflows {
		workflowDTO := &riskDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: rrupws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.RiskRegisterUpdatePlanWorkflowReviews {
			userDTO := rrupws.mapUserToDTO(review.CreatedByUser)
			if review.WorkflowStatus.Type == "approve" {
				workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
				if review.CreatedAt.After(lastApprovedAt) {
					lastApprovedAt = review.CreatedAt
				}
			}
			if review.WorkflowStatus.Type == "revert" {
				workflowDTO.RejectedUser = userDTO
				if review.CreatedAt.After(lastRejectedAt) {
					lastRejectedAt = review.CreatedAt
				}
			}
		}

		for _, participant := range workflow.RiskRegisterUpdatePlanWorkflowParticipants {
			userDTO := rrupws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.RiskRegisterUpdatePlanWorkflows = append(dto.RiskRegisterUpdatePlanWorkflows, workflowDTO)
	}

	dto.History = rrupws.buildWorkflowHistory(histories)
	return dto, nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) CreateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := rrupws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := rrupws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := rrupws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := rrupws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		riskRegisterUpdatePlanWorkflow := &models.RiskRegisterUpdatePlanWorkflow{
			RiskRegisterUpdatePlanID: riskRegisterUpdatePlanID,
			Level:                    level.Level,
			WorkflowLevelTypeID:      level.WorkflowLevelType.ID,
			ApprovalTypes:            level.ApprovalType,
			WorkflowStatusID:         workflowStatus.ID,
			CreatedBYID:              currentUserID,
		}

		workflowID, err := rrupws.repo.CreateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := rrupws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = rrupws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = rrupws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow, tx ...utility.DatabaseTransaction) error {
	err := rrupws.repo.UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) createWorkflowParticipant(userID uint, riskRegisterUpdatePlanWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	// Get the "pending" workflow status
	pendingStatus, err := rrupws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	workflowParticipant := &models.RiskRegisterUpdatePlanWorkflowParticipant{
		RiskRegisterUpdatePlanWorkflowID: riskRegisterUpdatePlanWorkflowID,
		UserID:                           userID,
		WorkflowStatusID:                 pendingStatus.ID,
		CreatedBYID:                      currentUserID,
	}

	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}

	err = rrupws.repo.CreateRiskRegisterUpdatePlanWorkflowParticipant(workflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error) {
	riskRegisterUpdatePlanWorkflow, err := rrupws.repo.GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return riskRegisterUpdatePlanWorkflow, nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) GetRiskRegisterUpdatePlanWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error) {
	riskRegisterUpdatePlanWorkflow, err := rrupws.repo.GetRiskRegisterUpdatePlanWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return riskRegisterUpdatePlanWorkflow, nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) GetNextLevel(riskRegisterUpdatePlanID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error) {
	nextWorkflow, err := rrupws.repo.GetNextLevel(riskRegisterUpdatePlanID, currentLevel, tx...)
	if err != nil {
		return nil, err
	}
	return nextWorkflow, nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskRegisterUpdatePlanWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := rrupws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	return count == len(workflowParticipants), nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) DeleteRiskRegisterUpdatePlanWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return rrupws.repo.DeleteRiskRegisterUpdatePlanWorkflowParticipants(participantIds, tx...)
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) DeleteRiskRegisterUpdatePlanWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return rrupws.repo.DeleteRiskRegisterUpdatePlanWorkflowReviews(reviewIds, tx...)
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) DeleteRiskRegisterUpdatePlanWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return rrupws.repo.DeleteRiskRegisterUpdatePlanWorkflows(workflowIds, tx...)
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := rrupws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	_, err = rrupws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	return nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) DeleteAllWorkflowsForRiskRegisterUpdatePlan(riskRegisterUpdatePlanID uint, tx ...utility.DatabaseTransaction) error {
	riskRegisterUpdatePlan, err := rrupws.repo.GetWorkflowList(riskRegisterUpdatePlanID)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range riskRegisterUpdatePlan.RiskRegisterUpdatePlanWorkflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.RiskRegisterUpdatePlanWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.RiskRegisterUpdatePlanWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := rrupws.DeleteRiskRegisterUpdatePlanWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := rrupws.DeleteRiskRegisterUpdatePlanWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := rrupws.DeleteRiskRegisterUpdatePlanWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (rrupws *RiskRegisterUpdatePlanWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
	if status == nil {
		return masterDTOs.WorkflowStatusDTO{}
	}
	return masterDTOs.WorkflowStatusDTO{
		ID:        status.ID,
		TitleEn:   status.TitleEn,
		TitleAr:   status.TitleAr,
		Type:      status.Type,
		ColorCode: status.ColorCode,
	}
}

func (rrupws *RiskRegisterUpdatePlanWorkflowService) buildWorkflowHistory(histories []*models.RiskRegisterUpdatePlanWorkflowHistory) []*riskDTOs.RiskRegisterUpdatePlanWorkflowHistoryDTO {
	var historyDTOs []*riskDTOs.RiskRegisterUpdatePlanWorkflowHistoryDTO

	for _, history := range histories {
		historyItem := &riskDTOs.RiskRegisterUpdatePlanWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: rrupws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			CreatedByUser:  rrupws.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
			Comment:        history.Comment,
		}
		historyDTOs = append(historyDTOs, historyItem)
	}

	return historyDTOs
}
