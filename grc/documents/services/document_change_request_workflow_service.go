package services

import (
	"fmt"
	"jethur/grc/documents/repositories"
	logSvc "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
)

type DocumentChangeRequestWorkflowServiceInterface interface {
	CreateDocumentCRWorkflow(documentChangeRequestId uint, documentId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	CreateDocumentCRWorkflowParticipant(userID uint, documentCRWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetDocumentCRWorkflowByLevel(changeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)

	DeleteDocumentCRWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentCRWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentCRWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
}

type DocumentChangeRequestWorkflowService struct {
	repo                                   repositories.DocumentChangeRequestWorkflowRepositoryInterface
	activityLogSvc                         logSvc.ActivityLogServiceInterface
	documentSvc                            DocumentServiceInterface
	workflowSvc                            settinsSvc.WorkflowServiceInterface
	documentChangeRequestWorkflowStatusSvc DocumentChangeRequestWorkflowStatusServiceInterface
	userSvc                                organizationSvc.UserServiceInterface
}

func NewDocumentChangeRequestWorkflowService(
	repo repositories.DocumentChangeRequestWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	documentChangeRequestWorkflowStatusSvc DocumentChangeRequestWorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
) *DocumentChangeRequestWorkflowService {
	return &DocumentChangeRequestWorkflowService{
		repo:                                   repo,
		activityLogSvc:                         activityLogSvc,
		documentSvc:                            documentSvc,
		workflowSvc:                            workflowSvc,
		documentChangeRequestWorkflowStatusSvc: documentChangeRequestWorkflowStatusSvc,
		userSvc:                                userSvc,
	}
}

func (dcrws *DocumentChangeRequestWorkflowService) CreateDocumentCRWorkflow(documentChangeRequestId uint, documentId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {

	document, err := dcrws.documentSvc.GetByID(documentId, tx...)
	if err != nil {
		return err
	}

	workflow, err := dcrws.workflowSvc.GetWorkflowByDocumentType(document.DocumentTypeID, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := dcrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	documentChangeRequestWorkflowStatus, err := dcrws.documentChangeRequestWorkflowStatusSvc.GetByType("pending")

	if err != nil {
		return err
	}

	// var sharedUsers []*models.DocumentSharedUser
	// uniqueUsers := make(map[uint]bool) // Map to track added users

	for _, level := range workflowLevels.WorkflowLevels {
		documentCRWorkflow := &models.DocumentChangeRequestWorkflow{
			DocumentChangeRequestID:               documentChangeRequestId,
			Level:                                 level.Level,
			WorkflowLevelTypeID:                   level.WorkflowLevelType.ID,
			ApprovalTypes:                         level.ApprovalType,
			DocumentChangeRequestWorkflowStatusID: documentChangeRequestWorkflowStatus.ID,
			CreatedBYID:                           currentUserID,
		}

		docCRWorkflowId, err := dcrws.repo.CreateDocumentCRWorkflow(documentCRWorkflow, tx...)
		if err != nil {
			return err
		}

		// Handle participants for each level
		for _, participant := range level.Participants {
			if participant.Role != nil {

				userIds, err := dcrws.userSvc.GetUserIdsByRole(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userId := range userIds {
					err = dcrws.CreateDocumentCRWorkflowParticipant(userId, docCRWorkflowId, &participant, currentUserID, tx...)
					if err != nil {
						return err
					}

					// if _, exists := uniqueUsers[userId]; exists {
					// 	continue
					// }

					// uniqueUsers[userId] = true

					// // Append to sharedUsers slice
					// sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
					// 	DocumentID:  documentId,
					// 	UserID:      userId,
					// 	Type:        "review-user",
					// 	CreatedBYID: currentUserID,
					// })

				}
			} else if participant.User != nil {
				userId := uint(participant.User.ID)
				err = dcrws.CreateDocumentCRWorkflowParticipant(userId, docCRWorkflowId, &participant, currentUserID, tx...)
				if err != nil {
					return err
				}

				// if _, exists := uniqueUsers[userId]; exists {
				// 	continue
				// }

				// uniqueUsers[userId] = true

				// // Append to sharedUsers slice
				// sharedUsers = append(sharedUsers, &models.DocumentSharedUser{
				// 	DocumentID:  documentId,
				// 	UserID:      userId,
				// 	Type:        "review-user",
				// 	CreatedBYID: currentUserID,
				// })
			}
		}
	}

	// if err := dws.documentAccessSvc.CreateUsers(sharedUsers, tx...); err != nil {
	// 	return err
	// }

	return nil
}

func (dcrws *DocumentChangeRequestWorkflowService) CreateDocumentCRWorkflowParticipant(userID uint, documentCRWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	docCRWorkflowParticipant := &models.DocumentChangeRequestWorkflowParticipant{
		DocumentChangeRequestWorkflowID: documentCRWorkflowID,
		UserID:                          userID,
		CreatedBYID:                     currentUserID,
	}

	if participant.Role != nil {
		docCRWorkflowParticipant.RoleID = &participant.Role.ID
	}

	err := dcrws.repo.CreateDocumentCRWorkflowParticipant(docCRWorkflowParticipant, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (dcrws *DocumentChangeRequestWorkflowService) DeleteDocumentCRWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(participantIds) == 0 {
		return nil
	}

	err := dcrws.repo.DeleteDocumentCRWorkflowParticipants(participantIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete document workflow participants: %w", err)
	}

	return nil
}

func (dcrws *DocumentChangeRequestWorkflowService) DeleteDocumentCRWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(reviewIds) == 0 {
		return nil
	}

	err := dcrws.repo.DeleteDocumentCRWorkflowReviews(reviewIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete document workflow reviews: %w", err)
	}

	return nil
}

func (dcrws *DocumentChangeRequestWorkflowService) DeleteDocumentCRWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	if len(workflowIds) == 0 {
		return nil
	}

	err := dcrws.repo.DeleteDocumentCRWorkflows(workflowIds, tx...)
	if err != nil {
		return fmt.Errorf("failed to delete document workflows: %w", err)
	}

	return nil
}

func (dcrws *DocumentChangeRequestWorkflowService) GetDocumentCRWorkflowByLevel(changeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error) {
	documentCRWorkflow, err := dcrws.repo.GetDocumentCRWorkflowByLevel(changeRequestID, currentApprovalLevel, tx...)
	if err != nil {
		return nil, err
	}
	return documentCRWorkflow, nil
}
