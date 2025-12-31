package services

import (
	"fmt"
	logSvc "jethur/grc/logs/services"
	masterDTOs "jethur/grc/masters/dtos"
	masterSvc "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	tprmDTOs "jethur/grc/tprm/dtos"
	"jethur/grc/tprm/repositories"
	"jethur/utility"
	"time"
)

type ThirdPartyAssessmentWorkflowServiceInterface interface {
	SetThirdPartyAssessmentService(svc ThirdPartyAssessmentServiceInterface)
	GetTprmWorkflowList(assessmentID uint) (*tprmDTOs.ThirdPartyAssessmentWorkflowListDTO, error)
	CreateTprmAssessmentWorkflow(assessmentWorkflowID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateTprmAssessmentWorkflow(assessmentWorkflow *models.ThirdPartyAssessmentWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, tprmWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetTprmAssessmentWorkflowByLevel(assessmentWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error)
	GetTrpmAsessmentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error)
	GetNextLevel(tprmWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ThirdPartyAssessmentWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteTprmAssessmentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteTprmAssessmentReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteTprmAssessmentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error
}

type ThirdPartyAssessmentWorkflowService struct {
	repo              repositories.ThirdPartyAssessmentWorkflowRepositoryInterface
	activityLogSvc    logSvc.ActivityLogServiceInterface
	workflowSvc       settinsSvc.WorkflowServiceInterface
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface
	userSvc           organizationSvc.UserServiceInterface
	tprmAssessmentSvc ThirdPartyAssessmentServiceInterface
}

func NewThirdPartyAssessmentWorkflowService(
	repo repositories.ThirdPartyAssessmentWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
) ThirdPartyAssessmentWorkflowServiceInterface {
	return &ThirdPartyAssessmentWorkflowService{
		repo:              repo,
		activityLogSvc:    activityLogSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
	}
}

func (arws *ThirdPartyAssessmentWorkflowService) SetThirdPartyAssessmentService(svc ThirdPartyAssessmentServiceInterface) {
	arws.tprmAssessmentSvc = svc
}

func (arws *ThirdPartyAssessmentWorkflowService) GetTprmWorkflowList(assessmentID uint) (*tprmDTOs.ThirdPartyAssessmentWorkflowListDTO, error) {

	assessment, err := arws.tprmAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		return nil, err
	}

	assessmentWorkflow, err := arws.repo.GetWorkflowList(assessment.ID)
	if err != nil {
		return nil, err
	}

	dto := &tprmDTOs.ThirdPartyAssessmentWorkflowListDTO{
		ID:                            assessmentWorkflow.ID,
		ThirdPartyAssessmentID:        assessmentWorkflow.ID,
		SubmittedAt:                   assessmentWorkflow.SubmittedAt,
		SubmittedByUser:               arws.mapUserToDTO(assessment.SubmittedByUser),
		ThirdPartyAssessmentWorkflows: []*tprmDTOs.WorkflowListDTO{},
	}

	for _, workflow := range assessmentWorkflow.ThirdPartyAssessmentWorkflows {
		workflowDTO := &tprmDTOs.WorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: arws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		var lastRejectedAt time.Time
		for _, review := range workflow.ThirdPartyAssessmentWorkflowReviews {
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

		for _, participant := range workflow.ThirdPartyAssessmentWorkflowParticipants {
			userDTO := arws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RejectedAt = lastRejectedAt
		dto.ThirdPartyAssessmentWorkflows = append(dto.ThirdPartyAssessmentWorkflows, workflowDTO)
	}

	return dto, nil
}

func (arws *ThirdPartyAssessmentWorkflowService) CreateTprmAssessmentWorkflow(assessmentWorkflowID uint, contextType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {

	workflow, err := arws.workflowSvc.GetAllWorkflowsByType("tprm-assessment", "tprm")
	if err != nil {
		return err
	}

	if len(workflow) < 1 {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	workflowLevels, err := arws.workflowSvc.GetWorkflowLevels(workflow[0].ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := arws.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		assessmentWorkflow := &models.ThirdPartyAssessmentWorkflow{
			ThirdPartyAssessmentID: assessmentWorkflowID,
			Level:                  level.Level,
			WorkflowLevelTypeID:    level.WorkflowLevelType.ID,
			ApprovalTypes:          level.ApprovalType,
			WorkflowStatusID:       workflowStatus.ID,
			CreatedBYID:            currentUserID,
		}

		docWorkflowId, err := arws.repo.CreateTprmAssessmentWorkflow(assessmentWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := arws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = arws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				err = arws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (arws *ThirdPartyAssessmentWorkflowService) UpdateTprmAssessmentWorkflow(assessmentWorkflow *models.ThirdPartyAssessmentWorkflow, tx ...utility.DatabaseTransaction) error {
	err := arws.repo.UpdateTprmAssessmentWorkflow(assessmentWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *ThirdPartyAssessmentWorkflowService) createWorkflowParticipant(userID uint, tprmWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	riskWorkflowParticipant := &models.ThirdPartyAssessmentWorkflowParticipant{
		TprmAssessmentWorkflowID: tprmWorkflowID,
		UserID:                   userID,
		CreatedBYID:              currentUserID,
	}

	if participant.Role != nil {
		riskWorkflowParticipant.RoleID = &participant.Role.ID
	}

	err := arws.repo.CreateTprmAssessmentWorkflowParticipant(riskWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *ThirdPartyAssessmentWorkflowService) GetTprmAssessmentWorkflowByLevel(assessmentWorkflowID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error) {
	assessmentWorkflow, err := arws.repo.GetTprmAssessmentWorkflowByLevel(assessmentWorkflowID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return assessmentWorkflow, nil
}

func (arws *ThirdPartyAssessmentWorkflowService) GetTrpmAsessmentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error) {
	assessmentWorkflow, err := arws.repo.GetTprmAssessmentWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return assessmentWorkflow, nil
}

func (arws *ThirdPartyAssessmentWorkflowService) GetNextLevel(riskWorkflowID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessmentWorkflow, error) {
	nextWorkflow, err := arws.repo.GetNextLevel(riskWorkflowID, currentLevel, tx...)
	if err != nil {
		return nil, err // Handle or log error as needed
	}
	return nextWorkflow, nil
}

func (arws *ThirdPartyAssessmentWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ThirdPartyAssessmentWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := arws.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (arws *ThirdPartyAssessmentWorkflowService) DeleteTprmAssessmentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(participantIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteTprmAssessmentWorkflowParticipants(participantIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete third party assessment workflow participants: %w", err)
	}

	return nil
}

func (arws *ThirdPartyAssessmentWorkflowService) DeleteTprmAssessmentReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteTprmAssessmentWorkflowReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete third party assessment workflow reviews: %w", err)
	}

	return nil
}

func (arws *ThirdPartyAssessmentWorkflowService) DeleteTprmAssessmentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(workflowIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteTprmAssessmentWorkflows(workflowIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete third party assessment workflows: %w", err)
	}

	return nil
}

func (arws *ThirdPartyAssessmentWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (arws *ThirdPartyAssessmentWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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

func (arws *ThirdPartyAssessmentWorkflowService) ValidateWorkflowBeforeCreate(contextType string, tx ...utility.DatabaseTransaction) error {
	// workflowContext, err := arws.workflowContextSvc.GetByType(contextType)
	// if err != nil {
	// 	return err
	// }

	workflow, err := arws.workflowSvc.GetAllWorkflowsByType("tprm-assessment", "tprm")
	if err != nil {
		return err
	}

	if len(workflow) < 1 {
		return utility.NewForbiddenError("workflow_not_configured")
	}

	workflowLevels, err := arws.workflowSvc.GetWorkflowLevels(workflow[0].ID, tx...)
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
