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
	settingsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"sort"
	"time"
)

type KRIRegisterWorkflowServiceInterface interface {
	SetKRIRegisterService(svc KRIRegisterServiceInterface)
	GetKRIRegisterWorkflowList(kriRegisterID uint) (*riskDTOs.KRIRegisterWorkflowListDTO, error)
	CreateKRIRegisterWorkflow(kriRegisterID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateKRIRegisterWorkflow(kriRegisterWorkflow *models.KRIRegisterWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, kriRegisterWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetKRIRegisterWorkflowByLevel(kriRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error)
	GetKRIRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error)
	GetNextLevel(kriRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.KRIRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteKRIRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteKRIRegisterReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteKRIRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
}

type KRIRegisterWorkflowService struct {
	repo               repositories.KRIRegisterWorkflowRepositoryInterface
	activityLogSvc     logSvc.ActivityLogServiceInterface
	workflowSvc        settingsSvc.WorkflowServiceInterface
	workflowStatusSvc  masterSvc.WorkflowStatusServiceInterface
	userSvc            organizationSvc.UserServiceInterface
	kriRegisterSvc     KRIRegisterServiceInterface
	workflowContextSvc masterSvc.WorkflowContextServiceInterface
}

func NewKRIRegisterWorkflowService(
	repo repositories.KRIRegisterWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
) KRIRegisterWorkflowServiceInterface {
	return &KRIRegisterWorkflowService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		workflowSvc:        workflowSvc,
		workflowStatusSvc:  workflowStatusSvc,
		userSvc:            userSvc,
		workflowContextSvc: workflowContextSvc,
	}
}

func (arws *KRIRegisterWorkflowService) SetKRIRegisterService(svc KRIRegisterServiceInterface) {
	arws.kriRegisterSvc = svc
}

func (arws *KRIRegisterWorkflowService) GetKRIRegisterWorkflowList(kriRegisterID uint) (*riskDTOs.KRIRegisterWorkflowListDTO, error) {
	kriRegister, err := arws.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return nil, err
	}

	kriRegisterWorkflow, err := arws.repo.GetWorkflowList(kriRegister.ID)
	if err != nil {
		return nil, err
	}

	dto := &riskDTOs.KRIRegisterWorkflowListDTO{
		ID:                   kriRegisterWorkflow.ID,
		KRIRegisterID:        kriRegisterWorkflow.ID,
		NextReviewLevel:      kriRegisterWorkflow.NextReviewLevel,
		SubmittedAt:          kriRegisterWorkflow.SubmittedAt,
		SubmittedByUser:      arws.mapUserToDTO(kriRegisterWorkflow.SubmittedBy),
		KRIRegisterWorkflows: []*riskDTOs.WorkflowListDTO{},
		History:              []*riskDTOs.KRIRegisterWorkflowHistoryDTO{},
	}

	for _, workflow := range kriRegisterWorkflow.KRIRegisterWorkflows {
		workflowDTO := &riskDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: arws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.KRIRegisterWorkflowReviews {
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

		for _, participant := range workflow.KRIRegisterWorkflowParticipants {
			userDTO := arws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.KRIRegisterWorkflows = append(dto.KRIRegisterWorkflows, workflowDTO)
	}

	dto.History = arws.buildWorkflowHistory(kriRegisterWorkflow)
	return dto, nil
}

func (arws *KRIRegisterWorkflowService) CreateKRIRegisterWorkflow(kriRegisterID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
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
		kriRegisterWorkflow := &models.KRIRegisterWorkflow{
			KRIRegisterID:       kriRegisterID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		workflowID, err := arws.repo.CreateKRIRegisterWorkflow(kriRegisterWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := arws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = arws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = arws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (arws *KRIRegisterWorkflowService) UpdateKRIRegisterWorkflow(kriRegisterWorkflow *models.KRIRegisterWorkflow, tx ...utility.DatabaseTransaction) error {
	err := arws.repo.UpdateKRIRegisterWorkflow(kriRegisterWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *KRIRegisterWorkflowService) createWorkflowParticipant(userID uint, kriRegisterWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.KRIRegisterWorkflowParticipant{
		KRIRegisterWorkflowID: kriRegisterWorkflowID,
		UserID:                userID,
		CreatedBYID:           currentUserID,
	}

	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}

	err := arws.repo.CreateKRIRegisterWorkflowParticipant(workflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *KRIRegisterWorkflowService) GetKRIRegisterWorkflowByLevel(kriRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error) {
	kriRegisterWorkflow, err := arws.repo.GetKRIRegisterWorkflowByLevel(kriRegisterID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return kriRegisterWorkflow, nil
}

func (arws *KRIRegisterWorkflowService) GetKRIRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error) {
	kriRegisterWorkflow, err := arws.repo.GetKRIRegisterWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return kriRegisterWorkflow, nil
}

func (arws *KRIRegisterWorkflowService) GetNextLevel(kriRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.KRIRegisterWorkflow, error) {
	nextWorkflow, err := arws.repo.GetNextLevel(kriRegisterID, currentLevel, tx...)
	if err != nil {
		return nil, err
	}
	return nextWorkflow, nil
}

func (arws *KRIRegisterWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.KRIRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := arws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (arws *KRIRegisterWorkflowService) DeleteKRIRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(participantIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteKRIRegisterWorkflowParticipants(participantIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete kri workflow participants: %w", err)
	}

	return nil
}

func (arws *KRIRegisterWorkflowService) DeleteKRIRegisterReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteKRIRegisterWorkflowReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete kri workflow reviews: %w", err)
	}

	return nil
}

func (arws *KRIRegisterWorkflowService) DeleteKRIRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(workflowIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteKRIRegisterWorkflows(workflowIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete kri workflows: %w", err)
	}

	return nil
}

func (arws *KRIRegisterWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (arws *KRIRegisterWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (arws *KRIRegisterWorkflowService) buildWorkflowHistory(kriRegister *models.KRIRegister) []*riskDTOs.KRIRegisterWorkflowHistoryDTO {
	var history []*riskDTOs.KRIRegisterWorkflowHistoryDTO

	for _, workflow := range kriRegister.KRIRegisterWorkflows {
		for _, review := range workflow.KRIRegisterWorkflowReviews {
			history = append(history, &riskDTOs.KRIRegisterWorkflowHistoryDTO{
				ID:             review.ID,
				Level:          workflow.Level,
				WorkflowStatus: arws.mapWorkflowStatusToDTO(review.WorkflowStatus),
				CreatedAt:      review.CreatedAt,
				CreatedByUser:  arws.mapUserToDTO(review.CreatedByUser),
			})
		}
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].CreatedAt.Before(history[j].CreatedAt)
	})

	return history
}

func (arws *KRIRegisterWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
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
