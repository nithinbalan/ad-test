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

type ProjectChangeRequestWorkflowServiceInterface interface {
	SetProjectChangeRequestService(svc ProjectChangeRequestServiceInterface)
	GetWorkflowList(projectChangeRequestID uint) (*projectDTOs.ProjectChangeRequestWorkflowListDTO, error)
	CreateProjectChangeRequestWorkflow(projectChangeRequestID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateProjectChangeRequestWorkflow(workflow *models.ProjectChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error
	GetWorkflowByLevel(projectChangeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error)
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error)
	GetNextLevel(projectChangeRequestID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForProjectChangeRequest(projectChangeRequestID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectChangeRequestWorkflowService struct {
	repo                 repositories.ProjectChangeRequestWorkflowRepositoryInterface
	reviewRepo           repositories.ProjectChangeRequestReviewRepositoryInterface
	activityLogSvc       logSvc.ActivityLogServiceInterface
	workflowSvc          settingsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	projectCRService     ProjectChangeRequestServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
}

func NewProjectChangeRequestWorkflowService(
	repo repositories.ProjectChangeRequestWorkflowRepositoryInterface,
	reviewRepo repositories.ProjectChangeRequestReviewRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) ProjectChangeRequestWorkflowServiceInterface {
	return &ProjectChangeRequestWorkflowService{
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

func (pcrws *ProjectChangeRequestWorkflowService) SetProjectChangeRequestService(svc ProjectChangeRequestServiceInterface) {
	pcrws.projectCRService = svc
}

func (pcrws *ProjectChangeRequestWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (pcrws *ProjectChangeRequestWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (pcrws *ProjectChangeRequestWorkflowService) GetWorkflowList(projectChangeRequestID uint) (*projectDTOs.ProjectChangeRequestWorkflowListDTO, error) {
	pcr, err := pcrws.repo.GetWorkflowList(projectChangeRequestID)
	if err != nil {
		return nil, err
	}

	histories, err := pcrws.repo.GetWorkflowHistories(projectChangeRequestID)
	if err != nil {
		return nil, err
	}

	var currentStatus *masterDTOs.WorkflowStatusDTO
	if pcr.WorkflowStatus != nil {
		st := pcrws.mapWorkflowStatusToDTO(pcr.WorkflowStatus)
		currentStatus = &st
	}

	dto := &projectDTOs.ProjectChangeRequestWorkflowListDTO{
		ID:                     pcr.ID,
		ProjectChangeRequestID: pcr.ID,
		SubmittedAt:            pcr.SubmittedAt,
		SubmittedByUser:        pcrws.mapUserToDTO(pcr.SubmittedByUser),
		ApprovedAt:             pcr.ApprovedAt,
		ApprovedByUser:         pcrws.mapUserToDTO(pcr.ApprovedByUser),
		RejectedAt:             pcr.RejectedAt,
		RejectedByUser:         pcrws.mapUserToDTO(pcr.RejectedByUser),
		RevertedAt:             pcr.RevertedAt,
		RevertedByUser:         pcrws.mapUserToDTO(pcr.RevertedByUser),
		NextReviewLevel:        pcr.NextReviewLevel,
		CurrentStatus:          currentStatus,
		Workflows:              []*projectDTOs.ProjectChangeRequestWorkflowLevelDTO{},
		History:                []*projectDTOs.ProjectChangeRequestWorkflowHistoryDTO{},
	}

	for _, workflow := range pcr.Workflows {
		workflowDTO := &projectDTOs.ProjectChangeRequestWorkflowLevelDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: pcrws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
		}

		for _, participant := range workflow.ProjectChangeRequestWorkflowParticipants {
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, pcrws.mapUserToDTO(participant.User))
		}

		var lastApprovedAt *time.Time
		var lastRejectedAt *time.Time
		for _, review := range workflow.ProjectChangeRequestWorkflowReviews {
			userDTO := pcrws.mapUserToDTO(review.CreatedBy)
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
		historyItem := &projectDTOs.ProjectChangeRequestWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: pcrws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			Comment:        history.Comment,
			CreatedByUser:  pcrws.mapUserToDTO(history.CreatedBy),
			CreatedAt:      history.CreatedAt,
		}
		dto.History = append(dto.History, historyItem)
	}

	return dto, nil
}

func (pcrws *ProjectChangeRequestWorkflowService) ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error {
	workflowContext, err := pcrws.workflowContextSvc.GetByType("project_change_request")
	if err != nil {
		return err
	}

	workflow, err := pcrws.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil || workflow == nil {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := pcrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
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

func (pcrws *ProjectChangeRequestWorkflowService) createWorkflowParticipant(userID uint, workflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.ProjectChangeRequestWorkflowParticipant{
		ProjectChangeRequestWorkflowID: workflowID,
		UserID:                         userID,
		CreatedBYID:                    currentUserID,
	}

	if participant.Role != nil {
		roleID := uint(participant.Role.ID)
		workflowParticipant.RoleID = &roleID
	}

	return pcrws.repo.CreateWorkflowParticipant(workflowParticipant, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) CreateProjectChangeRequestWorkflow(projectChangeRequestID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := pcrws.workflowContextSvc.GetByType("project_change_request")
	if err != nil {
		return err
	}

	workflow, err := pcrws.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := pcrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := pcrws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		pcrWorkflow := &models.ProjectChangeRequestWorkflow{
			ProjectChangeRequestID: projectChangeRequestID,
			Level:                  level.Level,
			WorkflowLevelTypeID:    level.WorkflowLevelType.ID,
			ApprovalTypes:          level.ApprovalType,
			WorkflowStatusID:       workflowStatus.ID,
			CreatedBYID:            currentUserID,
		}

		workflowID, err := pcrws.repo.CreateWorkflow(pcrWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := pcrws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = pcrws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = pcrws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.Unit != nil {
				deptHead, err := pcrws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}
				userID := uint(deptHead.HeadUser.ID)
				err = pcrws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (pcrws *ProjectChangeRequestWorkflowService) UpdateProjectChangeRequestWorkflow(workflow *models.ProjectChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error {
	return pcrws.repo.UpdateWorkflow(workflow, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) GetWorkflowByLevel(projectChangeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error) {
	return pcrws.repo.GetWorkflowByLevel(projectChangeRequestID, currentApprovalLevel, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error) {
	return pcrws.repo.GetWorkflowByID(id, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) GetNextLevel(projectChangeRequestID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error) {
	return pcrws.repo.GetNextLevel(projectChangeRequestID, currentLevel, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := pcrws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}
	return count == len(workflowParticipants), nil
}

func (pcrws *ProjectChangeRequestWorkflowService) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pcrws.repo.DeleteWorkflowParticipants(participantIds, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pcrws.repo.DeleteWorkflowReviews(reviewIds, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return pcrws.repo.DeleteWorkflows(workflowIds, tx...)
}

func (pcrws *ProjectChangeRequestWorkflowService) DeleteAllWorkflowsForProjectChangeRequest(projectChangeRequestID uint, tx ...utility.DatabaseTransaction) error {
	pcr, err := pcrws.repo.GetWorkflowList(projectChangeRequestID)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range pcr.Workflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.ProjectChangeRequestWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.ProjectChangeRequestWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := pcrws.DeleteWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := pcrws.DeleteWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := pcrws.DeleteWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}
