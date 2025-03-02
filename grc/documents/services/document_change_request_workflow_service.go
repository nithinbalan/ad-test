package services

import (
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logSvc "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingDTOs "jethur/grc/settings/dtos"
	settinsSvc "jethur/grc/settings/services"
	"jethur/utility"
	"time"
)

type DocumentChangeRequestWorkflowServiceInterface interface {
	GetDocumentCRWorkflowList(id uint) (*documentDTOs.DocumentCRWorkflowListDTO, error)
	CreateDocumentCRWorkflow(documentChangeRequestId uint, documentId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error
	CreateDocumentCRWorkflowParticipant(userID uint, documentCRWorkflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error

	GetDocumentCRWorkflowByLevel(changeRequestID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)
	GetDocumentCRWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)
	GetNextLevel(id uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error)

	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.DocumentChangeRequestWorkflowParticipant) (bool, error)

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
	documentChangeRequestSvc               DocumentChangeRequestServiceInterface
}

func NewDocumentChangeRequestWorkflowService(
	repo repositories.DocumentChangeRequestWorkflowRepositoryInterface,
	activityLogSvc logSvc.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	workflowSvc settinsSvc.WorkflowServiceInterface,
	documentChangeRequestWorkflowStatusSvc DocumentChangeRequestWorkflowStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	documentChangeRequestSvc DocumentChangeRequestServiceInterface,
) *DocumentChangeRequestWorkflowService {
	return &DocumentChangeRequestWorkflowService{
		repo:                                   repo,
		activityLogSvc:                         activityLogSvc,
		documentSvc:                            documentSvc,
		workflowSvc:                            workflowSvc,
		documentChangeRequestWorkflowStatusSvc: documentChangeRequestWorkflowStatusSvc,
		userSvc:                                userSvc,
		documentChangeRequestSvc:               documentChangeRequestSvc,
	}
}

func (dcrws *DocumentChangeRequestWorkflowService) CreateDocumentCRWorkflow(documentChangeRequestId uint, documentId uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {

	fmt.Println("12")
	document, err := dcrws.documentSvc.GetByID(documentId, tx...)
	if err != nil {
		return err
	}
	fmt.Println("13")
	workflow, err := dcrws.workflowSvc.GetWorkflowByDocumentType(document.DocumentTypeID, tx...)
	if err != nil {
		return err
	}
	fmt.Println("14")
	workflowLevels, err := dcrws.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}
	fmt.Println("15")
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
	fmt.Println("16")
	// if err := dws.documentAccessSvc.CreateUsers(sharedUsers, tx...); err != nil {
	// 	return err
	// }

	return nil
}

func (dcrws *DocumentChangeRequestWorkflowService) UpdateDocumentCRWorkflow(documentCRWorkflow *models.DocumentChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error {
	err := dcrws.repo.UpdateDocumentCRWorkflow(documentCRWorkflow, tx...)
	if err != nil {
		return err
	}
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

func (dcrws *DocumentChangeRequestWorkflowService) GetDocumentCRWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error) {
	documentCRWorkflow, err := dcrws.repo.GetDocumentCRWorkflowByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return documentCRWorkflow, nil
}

func (dcrws *DocumentChangeRequestWorkflowService) GetNextLevel(id uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentChangeRequestWorkflow, error) {
	nextWorkflow, err := dcrws.repo.GetNextLevel(id, currentLevel, tx...)
	if err != nil {
		return nil, err // Handle or log error as needed
	}
	return nextWorkflow, nil
}

func (dcrws *DocumentChangeRequestWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.DocumentChangeRequestWorkflowParticipant) (bool, error) {
	count, err := dcrws.repo.CountApprovedParticipants(workflowID)
	if err != nil {
		return false, err
	}

	// If the count of unique approvers is less than the number of participants, not all have approved.
	if count < len(workflowParticipants) {
		return false, nil
	}

	return true, nil
}

func (dcrws *DocumentChangeRequestWorkflowService) GetDocumentCRWorkflowList(id uint) (*documentDTOs.DocumentCRWorkflowListDTO, error) {
	// document, err := dcrws.documentChangeRequestSvc.GetFullDetailsByID(id)
	// if err != nil {
	// 	return nil, err
	// }

	documentChangeRequest, err := dcrws.repo.GetWorkflowList(id)
	if err != nil {
		return nil, err
	}

	dto := &documentDTOs.DocumentCRWorkflowListDTO{
		ID:                documentChangeRequest.ID,
		SubmittedAt:       documentChangeRequest.SubmittedAt,
		SubmittedByUser:   mapUserToDTO(documentChangeRequest.SubmittedByUser),
		DocumentWorkflows: []*documentDTOs.DocumentWorkflowListDTO{},
	}

	for _, workflow := range documentChangeRequest.DocumentChangeRequestWorkflows {
		workflowDTO := &documentDTOs.DocumentWorkflowListDTO{
			Level:         workflow.Level,
			ApprovalType:  workflow.ApprovalTypes,
			CurrentStatus: mapStatusToDTO((*models.DocumentWorkflowStatus)(workflow.DocumentChangeRequestWorkflowStatus)),
			ApprovedUsers: []*organizationDTOs.UserMinimalDTO{},
			ReviewUsers:   []*organizationDTOs.UserMinimalDTO{},
		}

		var lastApprovedAt time.Time
		for _, review := range workflow.DocumentChangeRequestWorkflowReviews {
			userDTO := mapUserToDTO(review.CreatedByUser)
			workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
			if review.CreatedAt.After(lastApprovedAt) {
				lastApprovedAt = review.CreatedAt
			}
		}

		for _, participant := range workflow.DocumentChangeRequestWorkflowParticipants {
			userDTO := mapUserToDTO(participant.User)
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, userDTO)
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		dto.DocumentWorkflows = append(dto.DocumentWorkflows, workflowDTO)
	}

	return dto, nil
}

// func mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
// 	if user == nil {
// 		return &organizationDTOs.UserMinimalDTO{}
// 	}
// 	return &organizationDTOs.UserMinimalDTO{
// 		ID:               int(user.ID),
// 		FirstNameEN:      user.FirstNameEn,
// 		LastNameEN:       user.LastNameEn,
// 		FirstNameAR:      user.FirstNameAr,
// 		LastNameAR:       user.LastNameAr,
// 		ProfileImagePath: user.ProfileImagePath,
// 		Email:            user.Email,
// 		Initial:          user.Initial,
// 		Color:            user.Color,
// 	}
// }

// func mapStatusToDTO(status *models.DocumentWorkflowStatus) documentDTOs.DocumentStatusDTO {
// 	if status == nil {
// 		return documentDTOs.DocumentStatusDTO{}
// 	}
// 	return documentDTOs.DocumentStatusDTO{
// 		ID:        int(status.ID),
// 		Type:      status.Type,
// 		TitleEn:   status.TitleEn,
// 		TitleAr:   status.TitleAr,
// 		ColorCode: status.ColorCode,
// 		Status:    status.Status,
// 	}
// }
