package services

import (
	"jethur/grc/documents/repositories"
	logSvc "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
)

type DocumentWorkflowServiceInterface interface {
	SetDocumentAccessService(svc DocumentAccessServiceInterface)
	CreateDocumentWorkflow(documentId uint, documentVersionId uint, documentTypeId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error
	createWorkflowParticipant(userID uint, documentWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetDocumentWorkflowByLevel(docVersionId uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	GetDocumentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	GetNextLevel(documentVersionID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.DocumentWorkflowParticipant) (bool, error)
}

type DocumentWorkflowService struct {
	repo                      repositories.DocumentWorkflowRepositoryInterface
	activityLogSvc            logSvc.ActivityLogServiceInterface
	workflowSvc               settinsSvc.WorkflowServiceInterface
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	documentAccessSvc         DocumentAccessServiceInterface
}

func NewDocumentWorkflowService(
	repo repositories.DocumentWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
) DocumentWorkflowServiceInterface {
	return &DocumentWorkflowService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		workflowSvc:               workflowSvc,
		documentWorkflowStatusSvc: documentWorkflowStatusSvc,
		userSvc:                   userSvc,
	}
}

func (dws *DocumentWorkflowService) SetDocumentAccessService(svc DocumentAccessServiceInterface) {
	dws.documentAccessSvc = svc
}

func (dws *DocumentWorkflowService) CreateDocumentWorkflow(documentId uint, documentVersionId uint, documentTypeId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflow, err := dws.workflowSvc.GetWorkflowByDocumentType(documentTypeId, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := dws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	documentWorkflowStatus, err := dws.documentWorkflowStatusSvc.GetByType("pending")

	if err != nil {
		return err
	}

	var sharedUsers []*models.DocumentSharedUser
	uniqueUsers := make(map[uint]bool) // Map to track added users

	for _, level := range workflowLevels.WorkflowLevels {
		documentWorkflow := &models.DocumentWorkflow{
			DocumentVersionID:        documentVersionId,
			Level:                    level.Level,
			WorkflowLevelTypeID:      level.WorkflowLevelType.ID,
			ApprovalTypes:            level.ApprovalType,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			CreatedBYID:              currentUserID,
		}

		docWorkflowId, err := dws.repo.CreateDocumentWorkflow(documentWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := dws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = dws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

					if _, exists := uniqueUsers[userId]; exists {
						continue
					}

					uniqueUsers[userId] = true

					// Append to sharedUsers slice
					sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
						DocumentID:  documentId,
						UserID:      userId,
						Type:        "review-user",
						CreatedBYID: currentUserID,
					})

				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				err = dws.createWorkflowParticipant(userId, docWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}

				if _, exists := uniqueUsers[userId]; exists {
					continue
				}

				uniqueUsers[userId] = true

				// Append to sharedUsers slice
				sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
					DocumentID:  documentId,
					UserID:      userId,
					Type:        "review-user",
					CreatedBYID: currentUserID,
				})
			}
		}
	}

	if err := dws.documentAccessSvc.CreateUsers(sharedUsers, tx...); err != nil {
		return err
	}

	return nil
}

func (dws *DocumentWorkflowService) UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error {
	err := dws.repo.UpdateDocumentWorkflow(documentWorkflow, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (dws *DocumentWorkflowService) createWorkflowParticipant(userID uint, documentWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	docWorkflowParticipant := &models.DocumentWorkflowParticipant{
		DocumentWorkflowID: documentWorkflowID,
		UserID:             userID,
		CreatedBYID:        currentUserID,
	}

	if participant.Role != nil {
		docWorkflowParticipant.RoleID = &participant.Role.ID
	}

	err := dws.repo.CreateDocumentWorkflowParticipant(docWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (dws *DocumentWorkflowService) GetDocumentWorkflowByLevel(docVersionId uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	documentWorkflow, err := dws.repo.GetDocumentWorkflowByLevel(docVersionId, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return documentWorkflow, nil
}

func (dws *DocumentWorkflowService) GetDocumentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	documentWorkflow, err := dws.repo.GetDocumentWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return documentWorkflow, nil
}

func (dws *DocumentWorkflowService) GetNextLevel(documentVersionID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	nextWorkflow, err := dws.repo.GetNextLevel(documentVersionID, currentLevel, tx...)
	if err != nil {
		return nil, err // Handle or log error as needed
	}
	return nextWorkflow, nil
}

func (dws *DocumentWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.DocumentWorkflowParticipant) (bool, error) {
	count, err := dws.repo.CountApprovedParticipants(workflowID)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}
