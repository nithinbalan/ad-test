package services

import (
	"fmt"
	commonServices "jethur/grc/common/services"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	settingDTOs "jethur/grc/settings/dtos"
	settingsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type ProjectDecisionLogWorkflowServiceInterface interface {
	SetProjectDecisionLogService(svc ProjectDecisionLogServiceInterface)
	GetWorkflowList(projectDecisionLogID uint) (*projectDTOs.ProjectDecisionLogWorkflowListDTO, error)
	CreateProjectDecisionLogWorkflow(projectDecisionLogID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateProjectDecisionLogWorkflow(workflow *models.ProjectDecisionLogWorkflow, tx ...utility.DatabaseTransaction) error
	GetWorkflowByLevel(projectDecisionLogID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error)
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error)
	GetNextLevel(projectDecisionLogID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectDecisionLogWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForProjectDecisionLog(projectDecisionLogID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectDecisionLogWorkflowService struct {
	repo                 repositories.ProjectDecisionLogWorkflowRepositoryInterface
	reviewRepo           repositories.ProjectDecisionLogReviewRepositoryInterface
	activityLogSvc       logSvc.ActivityLogServiceInterface
	workflowSvc          settingsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	projectDLService     ProjectDecisionLogServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
}

func NewProjectDecisionLogWorkflowService(
	repo repositories.ProjectDecisionLogWorkflowRepositoryInterface,
	reviewRepo repositories.ProjectDecisionLogReviewRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) ProjectDecisionLogWorkflowServiceInterface {
	return &ProjectDecisionLogWorkflowService{
		repo:                 repo,
		reviewRepo:           reviewRepo,
		activityLogSvc:       activityLogSvc,
		workflowSvc:          workflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		userSvc:              userSvc,
		workflowContextSvc:   workflowContextSvc,
		notificationSvc:      notificationSvc,
		organizationUnitsSvc: organizationUnitsSvc,
	}
}

func (pdlws *ProjectDecisionLogWorkflowService) SetProjectDecisionLogService(svc ProjectDecisionLogServiceInterface) {
	pdlws.projectDLService = svc
}

func (pdlws *ProjectDecisionLogWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (pdlws *ProjectDecisionLogWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
	if status == nil {
		return masterDTOs.WorkflowStatusDTO{}
	}
	return masterDTOs.WorkflowStatusDTO{
		ID:        status.ID,
		TitleEn:   status.TitleEn,
		TitleAr:   status.TitleAr,
		Type:      status.Type,
		ColorCode: status.ColorCode,
		Status:    status.Status,
	}
}

func (pdlws *ProjectDecisionLogWorkflowService) GetWorkflowList(projectDecisionLogID uint) (*projectDTOs.ProjectDecisionLogWorkflowListDTO, error) {
	pdl, err := pdlws.repo.GetWorkflowList(projectDecisionLogID)
	if err != nil {
		return nil, err
	}

	histories, err := pdlws.repo.GetWorkflowHistories(projectDecisionLogID)
	if err != nil {
		return nil, err
	}

	var currentStatus *masterDTOs.WorkflowStatusDTO
	if pdl.WorkflowStatus != nil {
		st := pdlws.mapWorkflowStatusToDTO(pdl.WorkflowStatus)
		currentStatus = &st
	}

	dto := &projectDTOs.ProjectDecisionLogWorkflowListDTO{
		ID:                   pdl.ID,
		ProjectDecisionLogID: pdl.ID,
		SubmittedAt:          pdl.SubmittedAt,
		SubmittedByUser:      pdlws.mapUserToDTO(pdl.SubmittedByUser),
		ApprovedAt:           pdl.ApprovedAt,
		ApprovedByUser:       pdlws.mapUserToDTO(pdl.ApprovedByUser),
		RejectedAt:           pdl.RejectedAt,
		RejectedByUser:       pdlws.mapUserToDTO(pdl.RejectedByUser),
		NextReviewLevel:      pdl.NextReviewLevel,
		CurrentStatus:        currentStatus,
		Workflows:            []*projectDTOs.ProjectDecisionLogWorkflowLevelDTO{},
		History:              []*projectDTOs.ProjectDecisionLogWorkflowHistoryDTO{},
	}

	for _, workflow := range pdl.Workflows {
		workflowDTO := &projectDTOs.ProjectDecisionLogWorkflowLevelDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: pdlws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
		}

		for _, participant := range workflow.ProjectDecisionLogWorkflowParticipants {
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, pdlws.mapUserToDTO(participant.User))
		}

		var lastApprovedAt *time.Time
		var lastRejectedAt *time.Time
		for _, review := range workflow.ProjectDecisionLogWorkflowReviews {
			userDTO := pdlws.mapUserToDTO(review.CreatedBy)
			if review.WorkflowStatus != nil {
				if review.WorkflowStatus.Type == "approve" {
					workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
					lastApprovedAt = &review.CreatedAt
				} else if review.WorkflowStatus.Type == "reject" || review.WorkflowStatus.Type == "revert" {
					lastRejectedAt = &review.CreatedAt
					workflowDTO.RejectedUser = userDTO
				}
			}
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		if lastApprovedAt != nil {
			workflowDTO.ApprovedAtString = lastApprovedAt.Format("2006-01-02T15:04:05.000Z")
		}
		workflowDTO.RejectedAt = lastRejectedAt
		dto.Workflows = append(dto.Workflows, workflowDTO)
	}

	for _, history := range histories {
		historyItem := &projectDTOs.ProjectDecisionLogWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: pdlws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			Comment:        history.Comment,
			CreatedByUser:  pdlws.mapUserToDTO(history.CreatedBy),
			CreatedAt:      history.CreatedAt,
		}
		dto.History = append(dto.History, historyItem)
	}

	return dto, nil
}

func (pdlws *ProjectDecisionLogWorkflowService) ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error {
	workflowContext, err := pdlws.workflowContextSvc.GetByType("project_decision_log")
	if err != nil {
		return err
	}

	workflow, err := pdlws.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil || workflow == nil {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := pdlws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
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

func (pdlws *ProjectDecisionLogWorkflowService) createWorkflowParticipant(userID uint, workflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.ProjectDecisionLogWorkflowParticipant{
		ProjectDecisionLogWorkflowID: workflowID,
		UserID:                       userID,
		CreatedBYID:                  currentUserID,
	}

	if participant.Role != nil {
		roleID := uint(participant.Role.ID)
		workflowParticipant.RoleID = &roleID
	}

	return pdlws.repo.CreateWorkflowParticipant(workflowParticipant, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) CreateProjectDecisionLogWorkflow(projectDecisionLogID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := pdlws.workflowContextSvc.GetByType("project_decision_log")
	if err != nil {
		return err
	}

	workflow, err := pdlws.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := pdlws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := pdlws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		pdlWorkflow := &models.ProjectDecisionLogWorkflow{
			ProjectDecisionLogID: projectDecisionLogID,
			Level:                level.Level,
			WorkflowLevelTypeID:  level.WorkflowLevelType.ID,
			ApprovalTypes:        level.ApprovalType,
			WorkflowStatusID:     workflowStatus.ID,
			CreatedBYID:          currentUserID,
		}

		workflowID, err := pdlws.repo.CreateWorkflow(pdlWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := pdlws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = pdlws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = pdlws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.Unit != nil {
				deptHead, err := pdlws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}
				userID := uint(deptHead.HeadUser.ID)
				err = pdlws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (pdlws *ProjectDecisionLogWorkflowService) UpdateProjectDecisionLogWorkflow(workflow *models.ProjectDecisionLogWorkflow, tx ...utility.DatabaseTransaction) error {
	return pdlws.repo.UpdateWorkflow(workflow, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) GetWorkflowByLevel(projectDecisionLogID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error) {
	return pdlws.repo.GetWorkflowByLevel(projectDecisionLogID, currentApprovalLevel, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error) {
	return pdlws.repo.GetWorkflowByID(id, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) GetNextLevel(projectDecisionLogID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error) {
	return pdlws.repo.GetNextLevel(projectDecisionLogID, currentLevel, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectDecisionLogWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := pdlws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}
	return count == len(workflowParticipants), nil
}

func (pdlws *ProjectDecisionLogWorkflowService) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pdlws.repo.DeleteWorkflowParticipants(participantIds, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pdlws.repo.DeleteWorkflowReviews(reviewIds, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pdlws.repo.DeleteWorkflows(workflowIds, tx...)
}

func (pdlws *ProjectDecisionLogWorkflowService) DeleteAllWorkflowsForProjectDecisionLog(projectDecisionLogID uint, tx ...utility.DatabaseTransaction) error {
	pdl, err := pdlws.repo.GetWorkflowList(projectDecisionLogID)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range pdl.Workflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.ProjectDecisionLogWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.ProjectDecisionLogWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := pdlws.DeleteWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := pdlws.DeleteWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := pdlws.DeleteWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}
