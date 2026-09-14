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

type ProjectWorkflowServiceInterface interface {
	SetProjectService(svc ProjectServiceInterface)
	GetWorkflowList(projectID uint) (*projectDTOs.ProjectWorkflowListDTO, error)
	CreateProjectWorkflow(projectID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateProjectWorkflow(workflow *models.ProjectWorkflow, tx ...utility.DatabaseTransaction) error
	GetWorkflowByLevel(projectID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error)
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error)
	GetNextLevel(projectID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowHistories(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForProject(projectID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectWorkflowService struct {
	repo                 repositories.ProjectWorkflowRepositoryInterface
	reviewRepo           repositories.ProjectReviewRepositoryInterface
	activityLogSvc       logSvc.ActivityLogServiceInterface
	workflowSvc          settingsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	projectService       ProjectServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
}

func NewProjectWorkflowService(
	repo repositories.ProjectWorkflowRepositoryInterface,
	reviewRepo repositories.ProjectReviewRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) ProjectWorkflowServiceInterface {
	return &ProjectWorkflowService{
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

func (pws *ProjectWorkflowService) SetProjectService(svc ProjectServiceInterface) {
	pws.projectService = svc
}

func (pws *ProjectWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (pws *ProjectWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (pws *ProjectWorkflowService) GetWorkflowList(projectID uint) (*projectDTOs.ProjectWorkflowListDTO, error) {
	project, err := pws.repo.GetWorkflowList(projectID)
	if err != nil {
		return nil, err
	}

	histories, err := pws.repo.GetWorkflowHistories(projectID)
	if err != nil {
		return nil, err
	}

	var workflowStatus *masterDTOs.WorkflowStatusDTO
	if project.WorkflowStatus != nil {
		st := pws.mapWorkflowStatusToDTO(project.WorkflowStatus)
		workflowStatus = &st
	}

	dto := &projectDTOs.ProjectWorkflowListDTO{
		ID:              project.ID,
		ProjectID:       project.ID,
		SubmittedAt:     project.SubmittedAt,
		SubmittedByUser: pws.mapUserToDTO(project.SubmittedByUser),
		ApprovedAt:      project.ApprovedAt,
		ApprovedByUser:  pws.mapUserToDTO(project.ApprovedByUser),
		RejectedAt:      project.RejectedAt,
		RejectedByUser:  pws.mapUserToDTO(project.RejectedByUser),
		NextReviewLevel: project.NextReviewLevel,
		WorkflowStatus:  workflowStatus,
		Workflows:       []*projectDTOs.ProjectWorkflowLevelDTO{},
		History:         []*projectDTOs.ProjectWorkflowHistoryDTO{},
	}

	for _, workflow := range project.Workflows {
		workflowDTO := &projectDTOs.ProjectWorkflowLevelDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			CurrentStatus: pws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
		}

		for _, participant := range workflow.ProjectWorkflowParticipants {
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, pws.mapUserToDTO(participant.User))
		}

		var lastApprovedAt *time.Time
		var lastRejectedAt *time.Time
		for _, review := range workflow.ProjectWorkflowReviews {
			userDTO := pws.mapUserToDTO(review.CreatedBy)
			if review.WorkflowStatus != nil {
				if review.WorkflowStatus.Type == "approve" {
					workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
					lastApprovedAt = &review.CreatedAt
				} else if review.WorkflowStatus.Type == "reject" {
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
		historyItem := &projectDTOs.ProjectWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: pws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			Comment:        history.Comment,
			CreatedByUser:  pws.mapUserToDTO(history.CreatedBy),
			CreatedAt:      history.CreatedAt,
		}
		dto.History = append(dto.History, historyItem)
	}

	return dto, nil
}

func (pws *ProjectWorkflowService) ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error {
	workflowContext, err := pws.workflowContextSvc.GetByType("project")
	if err != nil {
		return err
	}

	workflow, err := pws.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil || workflow == nil {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := pws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
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

func (pws *ProjectWorkflowService) createWorkflowParticipant(userID uint, workflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.ProjectWorkflowParticipant{
		ProjectWorkflowID: workflowID,
		UserID:            userID,
		CreatedBYID:       currentUserID,
	}

	if participant.Role != nil {
		roleID := uint(participant.Role.ID)
		workflowParticipant.RoleID = &roleID
	}

	return pws.repo.CreateWorkflowParticipant(workflowParticipant, tx...)
}

func (pws *ProjectWorkflowService) CreateProjectWorkflow(projectID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := pws.workflowContextSvc.GetByType("project")
	if err != nil {
		return err
	}

	workflow, err := pws.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := pws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := pws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		pWorkflow := &models.ProjectWorkflow{
			ProjectID:           projectID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		workflowID, err := pws.repo.CreateWorkflow(pWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := pws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = pws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = pws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.Unit != nil {
				deptHead, err := pws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}
				userID := uint(deptHead.HeadUser.ID)
				err = pws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (pws *ProjectWorkflowService) UpdateProjectWorkflow(workflow *models.ProjectWorkflow, tx ...utility.DatabaseTransaction) error {
	return pws.repo.UpdateWorkflow(workflow, tx...)
}

func (pws *ProjectWorkflowService) GetWorkflowByLevel(projectID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error) {
	return pws.repo.GetWorkflowByLevel(projectID, currentApprovalLevel, tx...)
}

func (pws *ProjectWorkflowService) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error) {
	return pws.repo.GetWorkflowByID(id, tx...)
}

func (pws *ProjectWorkflowService) GetNextLevel(projectID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error) {
	return pws.repo.GetNextLevel(projectID, currentLevel, tx...)
}

func (pws *ProjectWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := pws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}
	return count == len(workflowParticipants), nil
}

func (pws *ProjectWorkflowService) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pws.repo.DeleteWorkflowParticipants(participantIds, tx...)
}

func (pws *ProjectWorkflowService) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pws.repo.DeleteWorkflowReviews(reviewIds, tx...)
}

func (pws *ProjectWorkflowService) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return pws.repo.DeleteReviewsByWorkflowID(workflowID, tx...)
}

func (pws *ProjectWorkflowService) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	return pws.repo.UpdateWorkflowStatus(workflowID, statusID, tx...)
}

func (pws *ProjectWorkflowService) DeleteWorkflowHistories(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pws.repo.DeleteWorkflowHistories(workflowIds, tx...)
}

func (pws *ProjectWorkflowService) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pws.repo.DeleteWorkflows(workflowIds, tx...)
}

func (pws *ProjectWorkflowService) DeleteAllWorkflowsForProject(projectID uint, tx ...utility.DatabaseTransaction) error {
	project, err := pws.repo.GetWorkflowList(projectID, tx...)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range project.Workflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.ProjectWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.ProjectWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := pws.DeleteWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := pws.DeleteWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := pws.DeleteWorkflowHistories(workflowIds, tx...); err != nil {
			return err
		}
		if err := pws.DeleteWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}
