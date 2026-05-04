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
	"sort"
	"time"
)

type RiskResponseWorkflowServiceInterface interface {
	SetRiskResponseService(svc RiskResponseServiceInterface)
	GetRiskResponseWorkflowList(riskResponseID uint) (*riskDTOs.RiskResponseWorkflowListDTO, error)
	CreateRiskResponseWorkflow(riskResponseID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, riskResponseWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForRiskResponse(riskResponseID uint, tx ...utility.DatabaseTransaction) error
}

type RiskResponseWorkflowService struct {
	repo               repositories.RiskResponseWorkflowRepositoryInterface
	reviewRepo         repositories.RiskResponseReviewRepositoryInterface
	activityLogSvc     logSvc.ActivityLogServiceInterface
	workflowSvc        settingsSvc.WorkflowServiceInterface
	workflowStatusSvc  masterSvc.WorkflowStatusServiceInterface
	userSvc            organizationSvc.UserServiceInterface
	riskResponseSvc    RiskResponseServiceInterface
	workflowContextSvc masterSvc.WorkflowContextServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
}

func NewRiskResponseWorkflowService(
	repo repositories.RiskResponseWorkflowRepositoryInterface,
	reviewRepo repositories.RiskResponseReviewRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) RiskResponseWorkflowServiceInterface {
	return &RiskResponseWorkflowService{
		repo:               repo,
		reviewRepo:         reviewRepo,
		activityLogSvc:     activityLogSvc,
		workflowSvc:        workflowSvc,
		workflowStatusSvc:  workflowStatusSvc,
		userSvc:            userSvc,
		workflowContextSvc: workflowContextSvc,
		notificationSvc:    notificationSvc,
	}
}

func (rrws *RiskResponseWorkflowService) SetRiskResponseService(svc RiskResponseServiceInterface) {
	rrws.riskResponseSvc = svc
}

func (rrws *RiskResponseWorkflowService) GetRiskResponseWorkflowList(riskResponseID uint) (*riskDTOs.RiskResponseWorkflowListDTO, error) {
	riskResponse, err := rrws.riskResponseSvc.GetByID(riskResponseID)
	if err != nil {
		return nil, err
	}

	riskResponseWorkflow, err := rrws.repo.GetWorkflowList(riskResponse.ID)
	if err != nil {
		return nil, err
	}

	dto := &riskDTOs.RiskResponseWorkflowListDTO{
		ID:                    riskResponseWorkflow.ID,
		RiskResponseID:        riskResponseWorkflow.ID,
		NextReviewLevel:       riskResponseWorkflow.NextReviewLevel,
		SubmittedAt:           riskResponseWorkflow.SubmittedAt,
		SubmittedByUser:       rrws.mapUserToDTO(riskResponseWorkflow.SubmittedBy),
		RiskResponseWorkflows: []*riskDTOs.WorkflowListDTO{},
		History:               []*riskDTOs.RiskResponseWorkflowHistoryDTO{},
	}

	for _, workflow := range riskResponseWorkflow.RiskResponseWorkflows {
		workflowDTO := &riskDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: rrws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.RiskResponseWorkflowReviews {
			userDTO := rrws.mapUserToDTO(review.CreatedByUser)
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

		for _, participant := range workflow.RiskResponseWorkflowParticipants {
			userDTO := rrws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.RiskResponseWorkflows = append(dto.RiskResponseWorkflows, workflowDTO)
	}

	dto.History = rrws.buildWorkflowHistory(riskResponseWorkflow)
	return dto, nil
}

func (rrws *RiskResponseWorkflowService) CreateRiskResponseWorkflow(riskResponseID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := rrws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := rrws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := rrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := rrws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		riskResponseWorkflow := &models.RiskResponseWorkflow{
			RiskResponseID:      riskResponseID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		workflowID, err := rrws.repo.CreateRiskResponseWorkflow(riskResponseWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := rrws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = rrws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = rrws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (rrws *RiskResponseWorkflowService) UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error {
	err := rrws.repo.UpdateRiskResponseWorkflow(riskResponseWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (rrws *RiskResponseWorkflowService) createWorkflowParticipant(userID uint, riskResponseWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.RiskResponseWorkflowParticipant{
		RiskResponseWorkflowID: riskResponseWorkflowID,
		UserID:                 userID,
		CreatedBYID:            currentUserID,
	}

	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}

	err := rrws.repo.CreateRiskResponseWorkflowParticipant(workflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (rrws *RiskResponseWorkflowService) GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	riskResponseWorkflow, err := rrws.repo.GetRiskResponseWorkflowByLevel(riskResponseID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return riskResponseWorkflow, nil
}

func (rrws *RiskResponseWorkflowService) GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	riskResponseWorkflow, err := rrws.repo.GetRiskResponseWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return riskResponseWorkflow, nil
}

func (rrws *RiskResponseWorkflowService) GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	nextWorkflow, err := rrws.repo.GetNextLevel(riskResponseID, currentLevel, tx...)
	if err != nil {
		return nil, err
	}
	return nextWorkflow, nil
}

func (rrws *RiskResponseWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := rrws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	return count == len(workflowParticipants), nil
}

func (rrws *RiskResponseWorkflowService) DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return rrws.repo.DeleteRiskResponseWorkflowParticipants(participantIds, tx...)
}

func (rrws *RiskResponseWorkflowService) DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return rrws.repo.DeleteRiskResponseWorkflowReviews(reviewIds, tx...)
}

func (rrws *RiskResponseWorkflowService) DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return rrws.repo.DeleteRiskResponseWorkflows(workflowIds, tx...)
}

func (rrws *RiskResponseWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := rrws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	_, err = rrws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	return nil
}

func (rrws *RiskResponseWorkflowService) DeleteAllWorkflowsForRiskResponse(riskResponseID uint, tx ...utility.DatabaseTransaction) error {
	riskResponse, err := rrws.repo.GetWorkflowList(riskResponseID)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range riskResponse.RiskResponseWorkflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.RiskResponseWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.RiskResponseWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := rrws.DeleteRiskResponseWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := rrws.DeleteRiskResponseWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := rrws.DeleteRiskResponseWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (rrws *RiskResponseWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (rrws *RiskResponseWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (rrws *RiskResponseWorkflowService) buildWorkflowHistory(riskResponse *models.RiskResponse) []*riskDTOs.RiskResponseWorkflowHistoryDTO {
	var history []*riskDTOs.RiskResponseWorkflowHistoryDTO

	for _, workflow := range riskResponse.RiskResponseWorkflows {
		for _, review := range workflow.RiskResponseWorkflowReviews {
			historyItem := &riskDTOs.RiskResponseWorkflowHistoryDTO{
				ID:             review.ID,
				Level:          workflow.Level,
				WorkflowStatus: rrws.mapWorkflowStatusToDTO(review.WorkflowStatus),
				CreatedByUser:  rrws.mapUserToDTO(review.CreatedByUser),
				CreatedAt:      review.CreatedAt,
			}
			history = append(history, historyItem)
		}
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.After(history[j].CreatedAt)
	})

	return history
}
