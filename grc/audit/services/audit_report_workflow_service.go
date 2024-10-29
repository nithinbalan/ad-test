package services

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logSvc "jethur/grc/logs/services"
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

type AuditReportWorkflowServiceInterface interface {
	SetAuditFinalReportService(svc AuditFinalReportServiceInterface)
	GetAuditReportWorkflowList(auditID uint) (*auditDTOs.AuditFinalReportWorkflowListDTO, error)
	CreateAuditReportWorkflow(auditFinalReportID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, auditReportWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetAuditReportWorkflowByLevel(auditFinalReportID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)
	GetAuditReportWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)
	GetNextLevel(auditFinalReportID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.AuditReportWorkflowParticipant) (bool, error)

	DeleteAuditReportWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAuditReportReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAuditReportWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditReportWorkflowService struct {
	repo                repositories.AuditReportWorkflowRepositoryInterface
	activityLogSvc      logSvc.ActivityLogServiceInterface
	workflowSvc         settinsSvc.WorkflowServiceInterface
	workflowStatusSvc   masterSvc.WorkflowStatusServiceInterface
	userSvc             organizationSvc.UserServiceInterface
	auditFinalReportSvc AuditFinalReportServiceInterface
	auditSvc                  AuditServiceInterface
}

func NewAuditReportWorkflowService(
	repo repositories.AuditReportWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	auditSvc AuditServiceInterface,
) AuditReportWorkflowServiceInterface {
	return &AuditReportWorkflowService{
		repo:              repo,
		activityLogSvc:    activityLogSvc,
		workflowSvc:       workflowSvc,
		workflowStatusSvc: workflowStatusSvc,
		userSvc:           userSvc,
		auditSvc:           auditSvc,
	}
}

func (arws *AuditReportWorkflowService) SetAuditFinalReportService(svc AuditFinalReportServiceInterface) {
	arws.auditFinalReportSvc = svc
}

func (arws *AuditReportWorkflowService) GetAuditReportWorkflowList(auditID uint) (*auditDTOs.AuditFinalReportWorkflowListDTO, error) {
	_, err := arws.auditSvc.GetByID(auditID)
	if err != nil {
		return nil, err
	}
	finalReport, err := arws.auditFinalReportSvc.GetByAuditID(auditID)
	if err != nil {
		return nil, err
	}

	auditFinalReport, err := arws.repo.GetWorkflowList(finalReport.ID)
	if err != nil {
		return nil, err
	}

	dto := &auditDTOs.AuditFinalReportWorkflowListDTO{
		ID:                   auditFinalReport.ID,
		AuditID:              auditFinalReport.AuditID,
		AuditFinalReportID:   auditFinalReport.ID,
		SubmittedAt:          auditFinalReport.SubmittedAt,
		SubmittedByUser:      arws.mapUserToDTO(auditFinalReport.SubmittedByUser),
		AuditReportWorkflows: []*auditDTOs.AuditReportWorkflowListDTO{},
	}

	for _, workflow := range auditFinalReport.AuditReportWorkflows {
		workflowDTO := &auditDTOs.AuditReportWorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalType,
			CurrentStatus: arws.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		for _, review := range workflow.AuditReportReviews {
			userDTO := arws.mapUserToDTO(review.CreatedByUser)
			workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
			if review.CreatedAt.After(lastApprovedAt) {
				lastApprovedAt = review.CreatedAt
			}
		}

		for _, participant := range workflow.AuditReportWorkflowParticipants {
			userDTO := arws.mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		dto.AuditReportWorkflows = append(dto.AuditReportWorkflows, workflowDTO)
	}

	return dto, nil
}

func (arws *AuditReportWorkflowService) CreateAuditReportWorkflow(auditFinalReportID uint, workflowType string, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflow, err := arws.workflowSvc.GetAllWorkflowsByType(workflowType, "audit")
	if err != nil {
		return err
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
		auditReportWorkflow := &models.AuditReportWorkflow{
			AuditFinalReportID:  &auditFinalReportID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalType:        level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
		}

		docWorkflowId, err := arws.repo.CreateAuditReportWorkflow(auditReportWorkflow, tx...)
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

func (arws *AuditReportWorkflowService) UpdateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) error {
	err := arws.repo.UpdateAuditReportWorkflow(auditReportWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *AuditReportWorkflowService) createWorkflowParticipant(userID uint, auditReportWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	auditReportWorkflowParticipant := &models.AuditReportWorkflowParticipant{
		AuditReportWorkflowID: auditReportWorkflowID,
		UserID:                userID,
		CreatedBYID:           currentUserID,
	}

	if participant.Role != nil {
		auditReportWorkflowParticipant.RoleID = &participant.Role.ID
	}

	err := arws.repo.CreateAuditReportWorkflowParticipant(auditReportWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (arws *AuditReportWorkflowService) GetAuditReportWorkflowByLevel(auditFinalReportID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	auditReportWorkflow, err := arws.repo.GetAuditReportWorkflowByLevel(auditFinalReportID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return auditReportWorkflow, nil
}

func (arws *AuditReportWorkflowService) GetAuditReportWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	auditReportWorkflow, err := arws.repo.GetAuditReportWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return auditReportWorkflow, nil
}

func (arws *AuditReportWorkflowService) GetNextLevel(auditFinalReportID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	nextWorkflow, err := arws.repo.GetNextLevel(auditFinalReportID, currentLevel, tx...)
	if err != nil {
		return nil, err // Handle or log error as needed
	}
	return nextWorkflow, nil
}

func (arws *AuditReportWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.AuditReportWorkflowParticipant) (bool, error) {
	count, err := arws.repo.CountApprovedParticipants(workflowID)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (arws *AuditReportWorkflowService) DeleteAuditReportWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(participantIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteAuditReportWorkflowParticipants(participantIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete audit report workflow participants: %w", err)
	}

	return nil
}

func (arws *AuditReportWorkflowService) DeleteAuditReportReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteAuditReportReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete audit report workflow reviews: %w", err)
	}

	return nil
}

func (arws *AuditReportWorkflowService) DeleteAuditReportWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(workflowIds) == 0 {
		return nil
	}

	err := arws.repo.DeleteAuditReportWorkflows(workflowIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete audit report workflows: %w", err)
	}

	return nil
}

func (arws *AuditReportWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (arws *AuditReportWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
	if status == nil {
		return masterDTOs.WorkflowStatusDTO{}
	}
	return masterDTOs.WorkflowStatusDTO{
		ID:        status.ID,
		TitleEn:   status.TitleEn,
		TitleAr:   status.TitleAr,
		ColorCode: status.ColorCode,
		Status:    status.Status,
	}
}
