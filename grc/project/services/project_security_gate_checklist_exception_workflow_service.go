package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
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

type ProjectSecurityGateChecklistExceptionWorkflowServiceInterface interface {
	GetWorkflowList(securityGateID uint, checklistID uint) (*projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowListDTO, error)
	GetDocumentPath(securityGateID uint, checklistID uint, historyID uint, token string) (string, error)
	CreateWorkflow(checklistID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflow(workflow *models.ProjectSecurityGateChecklistExceptionWorkflow, tx ...utility.DatabaseTransaction) error

	GetWorkflowByLevel(checklistID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error)
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error)
	GetNextLevel(checklistID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error)
	CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectSecurityGateChecklistExceptionWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error)
	HasUserApproved(workflowID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error)

	ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error

	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAllWorkflowsForChecklist(checklistID uint, tx ...utility.DatabaseTransaction) error

	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectSecurityGateChecklistExceptionWorkflowService struct {
	repo                 repositories.ProjectSecurityGateChecklistExceptionWorkflowRepositoryInterface
	checklistSvc         ProjectSecurityGateChecklistServiceInterface
	workflowSvc          settingsSvc.WorkflowServiceInterface
	workflowStatusSvc    masterSvc.WorkflowStatusServiceInterface
	workflowContextSvc   masterSvc.WorkflowContextServiceInterface
	userSvc              organizationSvc.UserServiceInterface
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface
}

func NewProjectSecurityGateChecklistExceptionWorkflowService(
	repo repositories.ProjectSecurityGateChecklistExceptionWorkflowRepositoryInterface,
	checklistSvc ProjectSecurityGateChecklistServiceInterface,
	workflowSvc settingsSvc.WorkflowServiceInterface,
	workflowStatusSvc masterSvc.WorkflowStatusServiceInterface,
	workflowContextSvc masterSvc.WorkflowContextServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	organizationUnitsSvc organizationSvc.OrganizationUnitsServiceInterface,
) ProjectSecurityGateChecklistExceptionWorkflowServiceInterface {
	return &ProjectSecurityGateChecklistExceptionWorkflowService{
		repo:                 repo,
		checklistSvc:         checklistSvc,
		workflowSvc:          workflowSvc,
		workflowStatusSvc:    workflowStatusSvc,
		workflowContextSvc:   workflowContextSvc,
		userSvc:              userSvc,
		organizationUnitsSvc: organizationUnitsSvc,
	}
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) GetWorkflowList(securityGateID uint, checklistID uint) (*projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowListDTO, error) {
	if _, err := psgcews.checklistSvc.GetByID(securityGateID, checklistID); err != nil {
		return nil, err
	}

	checklist, err := psgcews.repo.GetWorkflowList(checklistID)
	if err != nil {
		return nil, err
	}

	histories, err := psgcews.repo.GetWorkflowHistories(checklistID)
	if err != nil {
		return nil, err
	}

	dto := &projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowListDTO{
		ID:              checklist.ID,
		ChecklistID:     checklist.ID,
		NextReviewLevel: checklist.NextReviewLevel,
		SubmittedAt:     checklist.SubmittedAt,
		SubmittedByUser: psgcews.mapUserToDTO(checklist.SubmittedByUser),
		ApprovedAt:      checklist.ApprovedAt,
		ApprovedByUser:  psgcews.mapUserToDTO(checklist.ApprovedByUser),
		Workflows:       []*projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowLevelDTO{},
		History:         []*projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowHistoryItemDTO{},
	}

	if checklist.WorkflowStatus != nil {
		status := psgcews.mapWorkflowStatusToDTO(checklist.WorkflowStatus)
		dto.WorkflowStatus = &status
	}

	for _, workflow := range checklist.ExceptionWorkflows {
		workflowDTO := &projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowLevelDTO{
			Level:          workflow.Level,
			ApprovalType:   workflow.ApprovalTypes,
			WorkflowStatus: psgcews.mapWorkflowStatusToDTO(workflow.WorkflowStatus),
			ReviewUsers:    []*organizationDTOs.UserMinimalDTO{},
			ApprovedUsers:  []*organizationDTOs.UserMinimalDTO{},
		}

		for _, participant := range workflow.Participants {
			workflowDTO.ReviewUsers = append(workflowDTO.ReviewUsers, psgcews.mapUserToDTO(participant.User))
		}

		var lastApprovedAt *time.Time
		var lastRevertedAt *time.Time
		for _, review := range workflow.Reviews {
			if review.WorkflowStatus == nil {
				continue
			}

			userDTO := psgcews.mapUserToDTO(review.CreatedByUser)
			if review.WorkflowStatus.Type == "approve" {
				workflowDTO.ApprovedUsers = append(workflowDTO.ApprovedUsers, userDTO)
				if lastApprovedAt == nil || review.CreatedAt.After(*lastApprovedAt) {
					createdAt := review.CreatedAt
					lastApprovedAt = &createdAt
				}
			}
			if review.WorkflowStatus.Type == "revert" {
				workflowDTO.RevertedUser = userDTO
				if lastRevertedAt == nil || review.CreatedAt.After(*lastRevertedAt) {
					createdAt := review.CreatedAt
					lastRevertedAt = &createdAt
				}
			}
		}

		workflowDTO.ApprovedAt = lastApprovedAt
		workflowDTO.RevertedAt = lastRevertedAt
		dto.Workflows = append(dto.Workflows, workflowDTO)
	}

	for _, history := range histories {
		dto.History = append(dto.History, &projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowHistoryItemDTO{
			ID:             history.ID,
			Level:          history.Level,
			WorkflowStatus: psgcews.mapWorkflowStatusToDTO(history.WorkflowStatus),
			Comment:        history.Comment,
			Documents:      psgcews.mapDocumentsToDTO(history.Documents),
			CreatedByUser:  psgcews.mapUserToDTO(history.CreatedByUser),
			CreatedAt:      history.CreatedAt,
		})
	}

	return dto, nil
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) GetDocumentPath(securityGateID uint, checklistID uint, historyID uint, token string) (string, error) {
	if _, err := psgcews.checklistSvc.GetByID(securityGateID, checklistID); err != nil {
		return "", err
	}

	history, err := psgcews.repo.GetHistoryByID(historyID)
	if err != nil {
		return "", err
	}

	if history.ChecklistID != checklistID || len(history.Documents) == 0 {
		return "", &utility.ResourceNotFoundError{}
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", &utility.ResourceNotFoundError{}
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", &utility.ResourceNotFoundError{}
	}

	storagePath := filepath.Join(
		"project",
		"security-gates", strconv.FormatUint(uint64(securityGateID), 10),
		"checklists", strconv.FormatUint(uint64(checklistID), 10),
		"reviews", strconv.FormatUint(uint64(historyID), 10),
		"documents",
	)

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) ValidateWorkflowBeforeCreate(tx ...utility.DatabaseTransaction) error {
	workflowContext, err := psgcews.workflowContextSvc.GetByType("project_security_gate_exception")
	if err != nil {
		return err
	}

	workflow, err := psgcews.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil || workflow == nil {
		return utility.NewForbiddenError("workflow not configured")
	}

	workflowLevels, err := psgcews.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
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

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) CreateWorkflow(checklistID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowContext, err := psgcews.workflowContextSvc.GetByType("project_security_gate_exception")
	if err != nil {
		return err
	}

	workflow, err := psgcews.workflowSvc.GetWorkflowByContext(workflowContext.ID, tx...)
	if err != nil {
		return err
	}

	workflowLevels, err := psgcews.workflowSvc.GetWorkflowLevels(workflow.ID, tx...)
	if err != nil {
		return err
	}

	workflowStatus, err := psgcews.workflowStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	for _, level := range workflowLevels.WorkflowLevels {
		exceptionWorkflow := &models.ProjectSecurityGateChecklistExceptionWorkflow{
			ChecklistID:         checklistID,
			Level:               level.Level,
			WorkflowLevelTypeID: level.WorkflowLevelType.ID,
			ApprovalTypes:       level.ApprovalType,
			WorkflowStatusID:    workflowStatus.ID,
			CreatedBYID:         currentUserID,
			CreatedAt:           time.Now(),
		}

		workflowID, err := psgcews.repo.CreateWorkflow(exceptionWorkflow, tx...)
		if err != nil {
			return err
		}

		participantCount := 0
		for _, participant := range level.Participants {
			if participant.Role != nil {
				userIDs, err := psgcews.userSvc.GetUserIdsByRoleExcludingSystemUsers(participant.Role.ID, tx...)
				if err != nil {
					return err
				}

				for _, userID := range userIDs {
					if err := psgcews.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...); err != nil {
						return err
					}
					participantCount++
				}
			} else if participant.User != nil {
				userID := uint(participant.User.ID)
				if err := psgcews.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...); err != nil {
					return err
				}
				participantCount++
			} else if participant.Unit != nil {
				unitHead, err := psgcews.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(uint(participant.Unit.ID))
				if err != nil {
					return err
				}
				if unitHead.HeadUser == nil {
					return utility.NewForbiddenError(fmt.Sprintf("No head has been assigned to organization unit '%d'.", participant.Unit.ID))
				}

				userID := uint(unitHead.HeadUser.ID)
				if err := psgcews.createWorkflowParticipant(userID, workflowID, &participant, currentUserID, tx...); err != nil {
					return err
				}
				participantCount++
			}
		}

		if participantCount == 0 {
			return utility.NewForbiddenError(fmt.Sprintf("no approver found for workflow level %d", level.Level))
		}
	}

	return nil
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) createWorkflowParticipant(userID uint, workflowID uint, participant *settingDTOs.WorkflowLevelParticipantDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	workflowParticipant := &models.ProjectSecurityGateChecklistExceptionWorkflowParticipant{
		ExceptionWorkflowID: workflowID,
		UserID:              userID,
		CreatedBYID:         currentUserID,
		CreatedAt:           time.Now(),
	}

	if participant.Role != nil {
		workflowParticipant.RoleID = &participant.Role.ID
	}

	return psgcews.repo.CreateWorkflowParticipant(workflowParticipant, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) UpdateWorkflow(workflow *models.ProjectSecurityGateChecklistExceptionWorkflow, tx ...utility.DatabaseTransaction) error {
	return psgcews.repo.UpdateWorkflow(workflow, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) GetWorkflowByLevel(checklistID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error) {
	return psgcews.repo.GetWorkflowByLevel(checklistID, currentApprovalLevel, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error) {
	return psgcews.repo.GetWorkflowByID(id, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) GetNextLevel(checklistID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error) {
	return psgcews.repo.GetNextLevel(checklistID, currentLevel, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) CheckAllParticipantsApproved(workflowID uint, workflowParticipants []*models.ProjectSecurityGateChecklistExceptionWorkflowParticipant, tx ...utility.DatabaseTransaction) (bool, error) {
	count, err := psgcews.repo.CountApprovedParticipants(workflowID, tx...)
	if err != nil {
		return false, err
	}

	return count >= len(workflowParticipants), nil
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) HasUserApproved(workflowID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	return psgcews.repo.HasUserApproved(workflowID, userID, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	return psgcews.repo.DeleteWorkflowParticipants(participantIds, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	return psgcews.repo.DeleteWorkflowReviews(reviewIds, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	return psgcews.repo.DeleteWorkflows(workflowIds, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	return psgcews.repo.DeleteReviewsByWorkflowID(workflowID, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	return psgcews.repo.UpdateWorkflowStatus(workflowID, statusID, tx...)
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) DeleteAllWorkflowsForChecklist(checklistID uint, tx ...utility.DatabaseTransaction) error {
	checklist, err := psgcews.repo.GetWorkflowList(checklistID, tx...)
	if err != nil {
		return err
	}

	var workflowIds []*uint
	var participantIds []*uint
	var reviewIds []*uint

	for _, workflow := range checklist.ExceptionWorkflows {
		workflowIds = append(workflowIds, &workflow.ID)
		for _, participant := range workflow.Participants {
			participantIds = append(participantIds, &participant.ID)
		}
		for _, review := range workflow.Reviews {
			reviewIds = append(reviewIds, &review.ID)
		}
	}

	if len(reviewIds) > 0 {
		if err := psgcews.DeleteWorkflowReviews(reviewIds, tx...); err != nil {
			return err
		}
	}

	if len(participantIds) > 0 {
		if err := psgcews.DeleteWorkflowParticipants(participantIds, tx...); err != nil {
			return err
		}
	}

	if len(workflowIds) > 0 {
		if err := psgcews.DeleteWorkflows(workflowIds, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) mapUserToDTO(user *models.User) *organizationDTOs.UserMinimalDTO {
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

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) mapDocumentsToDTO(documents []*models.ProjectSecurityGateChecklistExceptionWorkflowReviewDocument) []*projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowReviewDocumentDTO {
	documentDTOs := make([]*projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowReviewDocumentDTO, 0, len(documents))

	for _, document := range documents {
		documentDTOs = append(documentDTOs, &projectDTOs.ProjectSecurityGateChecklistExceptionWorkflowReviewDocumentDTO{
			ID:           document.ID,
			URL:          document.URL,
			Name:         document.Name,
			Size:         document.Size,
			Token:        document.Token,
			OriginalName: document.OriginalName,
			Type:         document.Type,
		})
	}

	return documentDTOs
}

func (psgcews *ProjectSecurityGateChecklistExceptionWorkflowService) mapWorkflowStatusToDTO(status *models.WorkflowStatus) masterDTOs.WorkflowStatusDTO {
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
