package services

import (
	"fmt"
	"time"

	commonServices "jethur/grc/common/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settingsSvc "jethur/grc/settings/services"
	"jethur/utility"
)

type IncidentClosureWorkflowServiceInterface interface {
	SetIncidentClosureService(svc IncidentClosureServiceInterface)
	GetIncidentClosureWorkflowList(closureID uint) (*incidentDTOs.IncidentClosureWorkflowListDTO, error)
	CreateIncidentClosureWorkflow(incidentClosureID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateIncidentClosureWorkflow(incidentClosureWorkflow *models.IncidentClosureWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, incidentClosureWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetIncidentClosureWorkflowByLevel(incidentClosureID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error)
	GetIncidentClosureWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error)
	GetNextLevel(incidentClosureID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.IncidentClosureWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteIncidentClosureWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentClosureWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentClosureWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForIncidentClosure(incidentClosureID uint, tx ...utility.DatabaseTransaction) error
}

type IncidentClosureWorkflowService struct {
	repo                 incidentRepositories.IncidentClosureWorkflowRepositoryInterface
	activityLogSvc       logSvc.ActivityLogServiceInterface
	workflowSvc          settingsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	incidentClosureSvc   IncidentClosureServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
}

func NewIncidentClosureWorkflowService(
	repo incidentRepositories.IncidentClosureWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) IncidentClosureWorkflowServiceInterface {
	return &IncidentClosureWorkflowService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		workflowSvc:          workflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		userSvc:              userSvc,
		workflowContextSvc:   workflowContextSvc,
		notificationSvc:      notificationSvc,
		organizationUnitsSvc: organizationUnitsSvc,
	}
}

func (icws *IncidentClosureWorkflowService) SetIncidentClosureService(svc IncidentClosureServiceInterface) {
	icws.incidentClosureSvc = svc
}

func (icws *IncidentClosureWorkflowService) GetIncidentClosureWorkflowList(closureID uint) (*incidentDTOs.IncidentClosureWorkflowListDTO, error) {
	incidentClosure, err := icws.incidentClosureSvc.GetByClosureID(closureID)
	if err != nil {
		return nil, err
	}

	incidentClosureWorkflow, err := icws.repo.GetWorkflowList(incidentClosure.ID)
	if err != nil {
		return nil, err
	}

	histories, err := icws.repo.GetWorkflowHistories(incidentClosure.ID)
	if err != nil {
		return nil, err
	}

	dto := &incidentDTOs.IncidentClosureWorkflowListDTO{
		ID:                       incidentClosureWorkflow.ID,
		IncidentClosureID:        incidentClosureWorkflow.ID,
		NextReviewLevel:          incidentClosureWorkflow.NextReviewLevel,
		SubmittedAt:              incidentClosureWorkflow.SubmittedAt,
		SubmittedByUser:          icws.mapUserToDTO(incidentClosureWorkflow.SubmittedByUser),
		IncidentClosureWorkflows: []*incidentDTOs.WorkflowListDTO{},
		History:                  []*incidentDTOs.IncidentClosureWorkflowHistoryDTO{},
	}

	for _, workflow := range incidentClosureWorkflow.IncidentClosureWorkflows {
		workflowDTO := &incidentDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: icws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.IncidentClosureWorkflowReviews {
			userDTO := icws.mapUserToDTO(review.CreatedByUser)
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

		for _, participant := range workflow.IncidentClosureWorkflowParticipants {
			userDTO := icws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.IncidentClosureWorkflows = append(dto.IncidentClosureWorkflows, workflowDTO)
	}

	dto.History = icws.buildWorkflowHistory(histories)
	return dto, nil
}

func (icws *IncidentClosureWorkflowService) CreateIncidentClosureWorkflow(incidentClosureID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := icws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := icws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := icws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := icws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		incidentClosureWorkflow := &models.IncidentClosureWorkflow{
			IncidentClosureID:   incidentClosureID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		workflowID, err := icws.repo.CreateIncidentClosureWorkflow(incidentClosureWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := icws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = icws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = icws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.Unit != nil {
				deptHead, err := icws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}
				userID := uint(deptHead.HeadUser.ID)
				err = icws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (icws *IncidentClosureWorkflowService) UpdateIncidentClosureWorkflow(incidentClosureWorkflow *models.IncidentClosureWorkflow, tx ...utility.DatabaseTransaction) error {
	err := icws.repo.UpdateIncidentClosureWorkflow(incidentClosureWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (icws *IncidentClosureWorkflowService) createWorkflowParticipant(userID uint, incidentClosureWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	pendingStatus, err := icws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	workflowParticipant := &models.IncidentClosureWorkflowParticipant{
		IncidentClosureWorkflowID: incidentClosureWorkflowID,
		UserID:                    userID,
		WorkflowStatusID:          pendingStatus.ID,
		CreatedBYID:               currentUserID,
	}

	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}

	err = icws.repo.CreateIncidentClosureWorkflowParticipant(workflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (icws *IncidentClosureWorkflowService) GetIncidentClosureWorkflowByLevel(incidentClosureID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error) {
	incidentClosureWorkflow, err := icws.repo.GetIncidentClosureWorkflowByLevel(incidentClosureID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return incidentClosureWorkflow, nil
}

func (icws *IncidentClosureWorkflowService) GetIncidentClosureWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error) {
	incidentClosureWorkflow, err := icws.repo.GetIncidentClosureWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return incidentClosureWorkflow, nil
}

func (icws *IncidentClosureWorkflowService) GetNextLevel(incidentClosureID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error) {
	nextWorkflow, err := icws.repo.GetNextLevel(incidentClosureID, currentLevel, tx...)
	if err != nil {
		return nil, err
	}
	return nextWorkflow, nil
}

func (icws *IncidentClosureWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.IncidentClosureWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := icws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	return count == len(workflowParticipants), nil
}

func (icws *IncidentClosureWorkflowService) DeleteIncidentClosureWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return icws.repo.DeleteIncidentClosureWorkflowParticipants(participantIds, tx...)
}

func (icws *IncidentClosureWorkflowService) DeleteIncidentClosureWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return icws.repo.DeleteIncidentClosureWorkflowReviews(reviewIds, tx...)
}

func (icws *IncidentClosureWorkflowService) DeleteIncidentClosureWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return icws.repo.DeleteIncidentClosureWorkflows(workflowIds, tx...)
}

func (icws *IncidentClosureWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := icws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	_, err = icws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	return nil
}

func (icws *IncidentClosureWorkflowService) DeleteAllWorkflowsForIncidentClosure(incidentClosureID uint, tx ...utility.DatabaseTransaction) error {
	incidentClosure, err := icws.repo.GetWorkflowList(incidentClosureID)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range incidentClosure.IncidentClosureWorkflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.IncidentClosureWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.IncidentClosureWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := icws.DeleteIncidentClosureWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := icws.DeleteIncidentClosureWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := icws.DeleteIncidentClosureWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (icws *IncidentClosureWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (icws *IncidentClosureWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (icws *IncidentClosureWorkflowService) buildWorkflowHistory(histories []*models.IncidentClosureWorkflowHistory) []*incidentDTOs.IncidentClosureWorkflowHistoryDTO {
	var historyDTOs []*incidentDTOs.IncidentClosureWorkflowHistoryDTO

	for _, history := range histories {
		historyItem := &incidentDTOs.IncidentClosureWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: icws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			CreatedByUser:  icws.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
			Comment:        history.Comment,
		}
		historyDTOs = append(historyDTOs, historyItem)
	}

	return historyDTOs
}
