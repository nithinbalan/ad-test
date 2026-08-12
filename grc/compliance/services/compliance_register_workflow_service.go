package services

import (
	"fmt"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
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

type ComplianceRegisterWorkflowServiceInterface interface {
	SetComplianceRegisterService(svc ComplianceRegisterServiceInterface)
	GetComplianceRegisterWorkflowList(complianceRegisterID uint) (*complianceDTOs.ComplianceRegisterWorkflowListDTO, error)
	CreateComplianceRegisterWorkflow(complianceRegisterID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow *models.ComplianceRegisterWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, complianceRegisterWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetComplianceRegisterWorkflowByLevel(complianceRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error)
	GetComplianceRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error)
	GetNextLevel(complianceRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ComplianceRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteComplianceRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRegisterWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForComplianceRegister(complianceRegisterID uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceRegisterWorkflowService struct {
	repo                  repositories.ComplianceRegisterWorkflowRepositoryInterface
	reviewRepo            repositories.ComplianceRegisterReviewRepositoryInterface
	activityLogSvc        logSvc.ActivityLogServiceInterface
	workflowSvc           settingsSvc.WorkflowServiceInterface
	workflowStatusSvc     masterSvc.WorkflowStatusServiceInterface
	userSvc               organizationSvc.UserServiceInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
	workflowContextSvc    masterSvc.WorkflowContextServiceInterface
	notificationSvc       commonServices.NotificationServiceInterface
	organizationUnitsSvc  organizationSvc.OrganizationUnitsServiceInterface
}

func NewComplianceRegisterWorkflowService(
	repo repositories.ComplianceRegisterWorkflowRepositoryInterface,
	reviewRepo repositories.ComplianceRegisterReviewRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) ComplianceRegisterWorkflowServiceInterface {
	return &ComplianceRegisterWorkflowService{
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

func (crws *ComplianceRegisterWorkflowService) SetComplianceRegisterService(svc ComplianceRegisterServiceInterface) {
	crws.complianceRegisterSvc = svc
}

func (crws *ComplianceRegisterWorkflowService) GetComplianceRegisterWorkflowList(complianceRegisterID uint) (*complianceDTOs.ComplianceRegisterWorkflowListDTO, error) {
	complianceRegister, err := crws.complianceRegisterSvc.GetByID(complianceRegisterID)
	if err != nil {
		return nil, err
	}

	complianceRegisterWorkflow, err := crws.repo.GetWorkflowList(complianceRegister.ID)
	if err != nil {
		return nil, err
	}

	// Get workflow histories
	histories, err := crws.repo.GetWorkflowHistories(complianceRegister.ID)
	if err != nil {
		return nil, err
	}

	dto := &complianceDTOs.ComplianceRegisterWorkflowListDTO{
		ID:                          complianceRegisterWorkflow.ID,
		ComplianceRegisterID:        complianceRegisterWorkflow.ID,
		NextReviewLevel:             complianceRegisterWorkflow.NextReviewLevel,
		SubmittedAt:                 complianceRegisterWorkflow.SubmittedAt,
		SubmittedByUser:             crws.mapUserToDTO(complianceRegisterWorkflow.SubmittedByUser),
		ComplianceRegisterWorkflows: []*complianceDTOs.WorkflowListDTO{},
		History:                     []*complianceDTOs.ComplianceRegisterWorkflowHistoryDTO{},
	}

	for _, workflow := range complianceRegisterWorkflow.ComplianceRegisterWorkflows {
		workflowDTO := &complianceDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: crws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.ComplianceRegisterWorkflowReviews {
			userDTO := crws.mapUserToDTO(review.CreatedByUser)
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

		for _, participant := range workflow.ComplianceRegisterWorkflowParticipants {
			userDTO := crws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.ComplianceRegisterWorkflows = append(dto.ComplianceRegisterWorkflows, workflowDTO)
	}

	dto.History = crws.buildWorkflowHistory(histories)
	return dto, nil
}

func (crws *ComplianceRegisterWorkflowService) CreateComplianceRegisterWorkflow(complianceRegisterID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := crws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return err
	}

	workflow, err := crws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return err
	}

	workflowLevels, err := crws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := crws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		complianceRegisterWorkflow := &models.ComplianceRegisterWorkflow{
			ComplianceRegisterID: complianceRegisterID,
			Level:                level.Level,
			WorkflowLevelTypeID:  level.WorkflowLevelType.ID,
			ApprovalTypes:        level.ApprovalType,
			WorkflowStatusID:     workflowStatus.ID,
			CreatedBYID:          currentUserID,
		}

		workflowID, err := crws.repo.CreateComplianceRegisterWorkflow(complianceRegisterWorkflow, tx...)
		if err != nil {
			return err
		}

		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := crws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					err = crws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				err = crws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			} else if participant.Unit != nil {
				deptHead, err := crws.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if deptHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}
				userID := uint(deptHead.HeadUser.ID)
				err = crws.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (crws *ComplianceRegisterWorkflowService) UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow *models.ComplianceRegisterWorkflow, tx ...utility.DatabaseTransaction) error {
	err := crws.repo.UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (crws *ComplianceRegisterWorkflowService) createWorkflowParticipant(userID uint, complianceRegisterWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.ComplianceRegisterWorkflowParticipant{
		ComplianceRegisterWorkflowID: complianceRegisterWorkflowID,
		UserID:                       userID,
		CreatedBYID:                  currentUserID,
	}

	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}

	err := crws.repo.CreateComplianceRegisterWorkflowParticipant(workflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (crws *ComplianceRegisterWorkflowService) GetComplianceRegisterWorkflowByLevel(complianceRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error) {
	complianceRegisterWorkflow, err := crws.repo.GetComplianceRegisterWorkflowByLevel(complianceRegisterID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return complianceRegisterWorkflow, nil
}

func (crws *ComplianceRegisterWorkflowService) GetComplianceRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error) {
	complianceRegisterWorkflow, err := crws.repo.GetComplianceRegisterWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return complianceRegisterWorkflow, nil
}

func (crws *ComplianceRegisterWorkflowService) GetNextLevel(complianceRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error) {
	nextWorkflow, err := crws.repo.GetNextLevel(complianceRegisterID, currentLevel, tx...)
	if err != nil {
		return nil, err
	}
	return nextWorkflow, nil
}

func (crws *ComplianceRegisterWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ComplianceRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := crws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	return count == len(workflowParticipants), nil
}

func (crws *ComplianceRegisterWorkflowService) DeleteComplianceRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return crws.repo.DeleteComplianceRegisterWorkflowParticipants(participantIds, tx...)
}

func (crws *ComplianceRegisterWorkflowService) DeleteComplianceRegisterWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return crws.repo.DeleteComplianceRegisterWorkflowReviews(reviewIds, tx...)
}

func (crws *ComplianceRegisterWorkflowService) DeleteComplianceRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return crws.repo.DeleteComplianceRegisterWorkflows(workflowIds, tx...)
}

func (crws *ComplianceRegisterWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := crws.workflowContextSvc.GetByType(contextType)
	if err != nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	_, err = crws.workflowSvc.GetWorkflowByContext(workflowContext.ID)
	if err != nil {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	return nil
}

func (crws *ComplianceRegisterWorkflowService) DeleteAllWorkflowsForComplianceRegister(complianceRegisterID uint, tx ...utility.DatabaseTransaction) error {
	complianceRegister, err := crws.repo.GetWorkflowList(complianceRegisterID)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range complianceRegister.ComplianceRegisterWorkflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.ComplianceRegisterWorkflowParticipants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.ComplianceRegisterWorkflowReviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := crws.DeleteComplianceRegisterWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := crws.DeleteComplianceRegisterWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := crws.DeleteComplianceRegisterWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (crws *ComplianceRegisterWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (crws *ComplianceRegisterWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (crws *ComplianceRegisterWorkflowService) buildWorkflowHistory(histories []*models.ComplianceRegisterWorkflowHistory) []*complianceDTOs.ComplianceRegisterWorkflowHistoryDTO {
	var historyDTOs []*complianceDTOs.ComplianceRegisterWorkflowHistoryDTO

	for _, history := range histories {
		historyItem := &complianceDTOs.ComplianceRegisterWorkflowHistoryDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: crws.mapWorkflowStatusToDTO(history.WorkflowStatus),
			CreatedByUser:  crws.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
		}
		historyDTOs = append(historyDTOs, historyItem)
	}

	return historyDTOs
}
