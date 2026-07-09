package services

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityRepositories "jethur/grc/business-continuity/repositories"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settingsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type BiaProcessWorkflowServiceInterface interface {
	SetBiaProcessService(svc BiaProcessServiceInterface)
	GetBiaProcessWorkflowList(biaProcessID uint) (*businessContinuityDTOs.BiaProcessWorkflowListDTO, error)
	CreateBiaProcessWorkflow(biaProcessID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateBiaProcessWorkflow(workflow *models.BiaProcessWorkflow, tx ...utility.DatabaseTransaction) error

	GetBiaProcessWorkflowByLevel(biaProcessID uint, level int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error)
	GetBiaProcessWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error)
	GetNextLevel(biaProcessID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, participants []*models.BiaProcessWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteBiaProcessWorkflowParticipants(participantIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteBiaProcessWorkflowReviews(reviewIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteBiaProcessWorkflows(workflowIDs []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForBiaProcess(biaProcessID uint, tx ...utility.DatabaseTransaction) error
	CreateWorkflowHistory(biaProcessID uint, level int, workflowStatusID uint, comment string, currentUserID uint, tx ...utility.DatabaseTransaction) error
}

type BiaProcessWorkflowService struct {
	repo               businessContinuityRepositories.BiaProcessWorkflowRepositoryInterface
	reviewRepo         businessContinuityRepositories.BiaProcessReviewRepositoryInterface
	activityLogSvc     logSvc.ActivityLogServiceInterface
	workflowSvc        settingsSvc.WorkflowServiceInterface
	workflowStatusSvc  masterSvc.WorkflowStatusServiceInterface
	userSvc            organizationSvc.UserServiceInterface
	biaProcessSvc      BiaProcessServiceInterface
	workflowContextSvc masterSvc.WorkflowContextServiceInterface
}

func NewBiaProcessWorkflowService(
	repo businessContinuityRepositories.BiaProcessWorkflowRepositoryInterface,
	reviewRepo businessContinuityRepositories.BiaProcessReviewRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
) BiaProcessWorkflowServiceInterface {
	return &BiaProcessWorkflowService{
		repo:               repo,
		reviewRepo:         reviewRepo,
		activityLogSvc:     activityLogSvc,
		workflowSvc:        workflowSvc,
		workflowStatusSvc:  workflowStatusSvc,
		userSvc:            userSvc,
		workflowContextSvc: workflowContextSvc,
	}
}

func (svc *BiaProcessWorkflowService) SetBiaProcessService(s BiaProcessServiceInterface) {
	svc.biaProcessSvc = s
}

func (svc *BiaProcessWorkflowService) GetBiaProcessWorkflowList(biaProcessID uint) (*businessContinuityDTOs.BiaProcessWorkflowListDTO, error) {
	biaProcess, err := svc.biaProcessSvc.GetByID(biaProcessID)
	if err != nil {
		return nil, err
	}

	biaProcessWorkflow, err := svc.repo.GetWorkflowList(biaProcess.ID)
	if err != nil {
		return nil, err
	}

	histories, err := svc.repo.GetWorkflowHistories(biaProcess.ID)
	if err != nil {
		return nil, err
	}

	dto := &businessContinuityDTOs.BiaProcessWorkflowListDTO{
		ID:                  biaProcessWorkflow.ID,
		BiaProcessID:        biaProcessWorkflow.ID,
		NextReviewLevel:     biaProcessWorkflow.NextReviewLevel,
		SubmittedAt:         biaProcessWorkflow.SubmittedAt,
		SubmittedByUser:     svc.mapUserToDTO(biaProcessWorkflow.SubmittedBy),
		BiaProcessWorkflows: []*businessContinuityDTOs.BiaWorkflowListDTO{},
		History:             []*businessContinuityDTOs.BiaProcessWorkflowHistoryDTO{},
	}

	for _, workflow := range biaProcessWorkflow.BiaProcessWorkflows {
		workflowDTO := &businessContinuityDTOs.BiaWorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: svc.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.BiaProcessWorkflowReviews {
			userDTO := svc.mapUserToDTO(review.CreatedByUser)
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

		for _, participant := range workflow.BiaProcessWorkflowParticipants {
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, svc.mapUserToDTO(participant.User))
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.BiaProcessWorkflows = append(dto.BiaProcessWorkflows, workflowDTO)
	}

	dto.History = svc.buildWorkflowHistory(histories)
	return dto, nil
}

func (svc *BiaProcessWorkflowService) CreateBiaProcessWorkflow(biaProcessID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := svc.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := svc.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := svc.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := svc.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		biaProcessWorkflow := &models.BiaProcessWorkflow{
			BiaProcessID:        biaProcessID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		workflowID, err := svc.repo.CreateBiaProcessWorkflow(biaProcessWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := svc.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}
				for _, userID := range userIDs {
					if err := svc.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				if err := svc.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (svc *BiaProcessWorkflowService) createWorkflowParticipant(userID uint, biaProcessWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.BiaProcessWorkflowParticipant{
		BiaProcessWorkflowID: biaProcessWorkflowID,
		UserID:               userID,
		CreatedBYID:          currentUserID,
	}
	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}
	return svc.repo.CreateBiaProcessWorkflowParticipant(workflowParticipant, tx...)
}

func (svc *BiaProcessWorkflowService) UpdateBiaProcessWorkflow(workflow *models.BiaProcessWorkflow, tx ...utility.DatabaseTransaction) error {
	return svc.repo.UpdateBiaProcessWorkflow(workflow, tx...)
}

func (svc *BiaProcessWorkflowService) GetBiaProcessWorkflowByLevel(biaProcessID uint, level int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error) {
	return svc.repo.GetBiaProcessWorkflowByLevel(biaProcessID, level, tx...)
}

func (svc *BiaProcessWorkflowService) GetBiaProcessWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error) {
	return svc.repo.GetBiaProcessWorkflowByID(id, tx...)
}

func (svc *BiaProcessWorkflowService) GetNextLevel(biaProcessID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error) {
	return svc.repo.GetNextLevel(biaProcessID, currentLevel, tx...)
}

func (svc *BiaProcessWorkflowService) CheckAllParticipantsApproved(workflowID uint, participants []*models.BiaProcessWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := svc.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}
	return count == len(participants), nil
}

func (svc *BiaProcessWorkflowService) DeleteBiaProcessWorkflowParticipants(participantIDs []*uint, tx ...utility.DatabaseTransaction) error {
	return svc.repo.DeleteBiaProcessWorkflowParticipants(participantIDs, tx...)
}

func (svc *BiaProcessWorkflowService) DeleteBiaProcessWorkflowReviews(reviewIDs []*uint, tx ...utility.DatabaseTransaction) error {
	return svc.repo.DeleteBiaProcessWorkflowReviews(reviewIDs, tx...)
}

func (svc *BiaProcessWorkflowService) DeleteBiaProcessWorkflows(workflowIDs []*uint, tx ...utility.DatabaseTransaction) error {
	return svc.repo.DeleteBiaProcessWorkflows(workflowIDs, tx...)
}

func (svc *BiaProcessWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := svc.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}
	_, err = svc.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}
	return nil
}

func (svc *BiaProcessWorkflowService) DeleteAllWorkflowsForBiaProcess(biaProcessID uint, tx ...utility.DatabaseTransaction) error {
	biaProcess, err := svc.repo.GetWorkflowList(biaProcessID)
	if err != nil {
		return err
	}

	var workflowIDs []*uint
	var participantIDs []*uint
	var reviewIDs []*uint

	for _, workflow := range biaProcess.BiaProcessWorkflows {
		workflowIDs = append(workflowIDs, &workflow.ID)
		for _, participant := range workflow.BiaProcessWorkflowParticipants {
			participantIDs = append(participantIDs, &participant.ID)
		}
		for _, review := range workflow.BiaProcessWorkflowReviews {
			reviewIDs = append(reviewIDs, &review.ID)
		}
	}

	if len(reviewIDs) > 0 {
		if err := svc.DeleteBiaProcessWorkflowReviews(reviewIDs, tx...); err != nil {
			return err
		}
	}
	if len(participantIDs) > 0 {
		if err := svc.DeleteBiaProcessWorkflowParticipants(participantIDs, tx...); err != nil {
			return err
		}
	}
	if len(workflowIDs) > 0 {
		if err := svc.DeleteBiaProcessWorkflows(workflowIDs, tx...); err != nil {
			return err
		}
	}
	return nil
}

func (svc *BiaProcessWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (svc *BiaProcessWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (svc *BiaProcessWorkflowService) CreateWorkflowHistory(biaProcessID uint, level int, workflowStatusID uint, comment string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	now := time.Now()
	history := &models.BiaProcessWorkflowHistory{
		BiaProcessID:     biaProcessID,
		Level:            level,
		WorkflowStatusID: workflowStatusID,
		Comment:          comment,
		CreatedBYID:      currentUserID,
		CreatedAt:        now,
	}
	return svc.repo.CreateWorkflowHistory(history, tx...)
}

func (svc *BiaProcessWorkflowService) buildWorkflowHistory(histories []*models.BiaProcessWorkflowHistory) []*businessContinuityDTOs.BiaProcessWorkflowHistoryDTO {
	var historyDTOs []*businessContinuityDTOs.BiaProcessWorkflowHistoryDTO

	for _, history := range histories {
		historyItem := &businessContinuityDTOs.BiaProcessWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: svc.mapWorkflowStatusToDTO(history.WorkflowStatus),
			CreatedByUser:  svc.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
		}
		historyDTOs = append(historyDTOs, historyItem)
	}

	return historyDTOs
}
