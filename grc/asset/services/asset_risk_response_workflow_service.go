package services

import (
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

type AssetRiskResponseWorkflowServiceInterface interface {
	GetRiskResponseWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*assetRiskDTOs.AssetRiskResponseWorkflowListDTO, error)
	CreateRiskResponseWorkflow(riskResponseID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error

	GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error

	DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForRiskResponse(riskID uint, tx ...utility.DatabaseTransaction) error

	DeleteRiskResponseWorkflowReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskResponseWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskResponseWorkflowService struct {
	repo               repositories.AssetRiskResponseWorkflowRepositoryInterface
	workflowSvc        settinsSvc.WorkflowServiceInterface
	workflowStatusSvc  masterSvc.WorkflowStatusServiceInterface
	userSvc            organizationSvc.UserServiceInterface
	workflowContextSvc masterSvc.WorkflowContextServiceInterface
}

func NewAssetRiskResponseWorkflowService(
	repo repositories.AssetRiskResponseWorkflowRepositoryInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
) AssetRiskResponseWorkflowServiceInterface {
	return &AssetRiskResponseWorkflowService{
		repo:               repo,
		workflowSvc:        workflowSvc,
		workflowStatusSvc:  workflowStatusSvc,
		userSvc:            userSvc,
		workflowContextSvc: workflowContextSvc,
	}
}

func (arrws *AssetRiskResponseWorkflowService) GetRiskResponseWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*assetRiskDTOs.AssetRiskResponseWorkflowListDTO, error) {
	riskResponse, err := arrws.repo.GetWorkflowList(riskID, tx...)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, riskID)
	}

	histories, err := arrws.repo.GetWorkflowHistories(riskID, tx...)
	if err != nil {
		return nil, err
	}

	dto := &assetRiskDTOs.AssetRiskResponseWorkflowListDTO{
		ID:                    riskResponse.ID,
		RiskResponseID:        riskResponse.ID,
		NextReviewLevel:       riskResponse.NextReviewLevel,
		SubmittedAt:           riskResponse.SubmittedAt,
		SubmittedByUser:       arrws.mapUserToDTO(riskResponse.SubmittedBy),
		ApprovedAt:            riskResponse.ApprovedAt,
		ApprovedByUser:        arrws.mapUserToDTO(riskResponse.ApprovedBy),
		RiskResponseWorkflows: []*assetRiskDTOs.WorkflowListDTO{},
		History:               []*assetRiskDTOs.AssetRiskResponseWorkflowHistoryDTO{},
	}

	if riskResponse.WorkflowStatus != nil {
		status := arrws.mapWorkflowStatusToDTO(riskResponse.WorkflowStatus)
		dto.CurrentStatus = &status
	}

	for _, workflow := range riskResponse.RiskResponseWorkflows {
		workflowDTO := &assetRiskDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: arrws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt *time.Time
		var lastRejectedAt *time.Time
		for _, review := range workflow.RiskResponseWorkflowReviews {
			userDTO := arrws.mapUserToDTO(review.CreatedByUser)
			if review.WorkflowStatus.Type == "approve" {
				workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
				if lastApprovedAt == nil || review.CreatedAt.After(*lastApprovedAt) {
					t := review.CreatedAt
					lastApprovedAt = &t
				}
			}
			if review.WorkflowStatus.Type == "revert" {
				workflowDTO.RejectedUser = userDTO
				if lastRejectedAt == nil || review.CreatedAt.After(*lastRejectedAt) {
					t := review.CreatedAt
					lastRejectedAt = &t
				}
			}
		}

		for _, participant := range workflow.RiskResponseWorkflowParticipants {
			userDTO := arrws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.RiskResponseWorkflows = append(dto.RiskResponseWorkflows, workflowDTO)
	}

	dto.History = arrws.buildWorkflowHistory(histories)

	return dto, nil
}

func (arrws *AssetRiskResponseWorkflowService) CreateRiskResponseWorkflow(riskResponseID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := arrws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := arrws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := arrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := arrws.workflowStatusSvc.GetByType("pending")
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

		riskResponseWorkflowID, err := arrws.repo.CreateRiskResponseWorkflow(riskResponseWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIds, err := arrws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					if err := arrws.createWorkflowParticipant(userId, riskResponseWorkflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				if err := arrws.createWorkflowParticipant(userId, riskResponseWorkflowID, &participant, currentUserID, tx...); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (arrws *AssetRiskResponseWorkflowService) createWorkflowParticipant(userID uint, riskResponseWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	riskResponseWorkflowParticipant := &models.RiskResponseWorkflowParticipant{
		RiskResponseWorkflowID: riskResponseWorkflowID,
		UserID:                 userID,
		CreatedBYID:            currentUserID,
	}

	if participant.Role != nil {
		riskResponseWorkflowParticipant.RoleID = &participant.Role.ID
	}

	return arrws.repo.CreateRiskResponseWorkflowParticipant(riskResponseWorkflowParticipant, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error {
	return arrws.repo.UpdateRiskResponseWorkflow(riskResponseWorkflow, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	return arrws.repo.GetRiskResponseWorkflowByLevel(riskResponseID, currentApprovalLevel, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	return arrws.repo.GetRiskResponseWorkflowByID(id, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	return arrws.repo.GetNextLevel(riskResponseID, currentLevel, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := arrws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	return count >= len(workflowParticipants), nil
}

func (arrws *AssetRiskResponseWorkflowService) DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return arrws.repo.DeleteRiskResponseWorkflowReviews(reviewIds, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return arrws.repo.DeleteRiskResponseWorkflowParticipants(participantIds, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return arrws.repo.DeleteRiskResponseWorkflows(workflowIds, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) DeleteRiskResponseWorkflowReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return arrws.repo.DeleteRiskResponseWorkflowReviewsByWorkflowID(workflowID, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) UpdateRiskResponseWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	return arrws.repo.UpdateRiskResponseWorkflowStatus(workflowID, statusID, tx...)
}

func (arrws *AssetRiskResponseWorkflowService) DeleteAllWorkflowsForRiskResponse(riskID uint, tx ...utility.DatabaseTransaction) error {
	riskResponse, err := arrws.repo.GetWorkflowsForDeletion(riskID, tx...)
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
		if err := arrws.DeleteRiskResponseWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := arrws.DeleteRiskResponseWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := arrws.DeleteRiskResponseWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (arrws *AssetRiskResponseWorkflowService) buildWorkflowHistory(histories []*models.RiskResponseWorkflowHistory) []*assetRiskDTOs.AssetRiskResponseWorkflowHistoryDTO {
	var historyDTOs []*assetRiskDTOs.AssetRiskResponseWorkflowHistoryDTO

	for _, history := range histories {
		historyDTOs = append(historyDTOs, &assetRiskDTOs.AssetRiskResponseWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: arrws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			CreatedByUser:  arrws.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
		})
	}

	return historyDTOs
}

func (arrws *AssetRiskResponseWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := arrws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := arrws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return utility.NewForbiddenError("workflow not configured")
	}
	if workflow == nil {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := arrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
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

func (arrws *AssetRiskResponseWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (arrws *AssetRiskResponseWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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
