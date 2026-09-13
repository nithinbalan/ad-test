package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateChecklistExceptionWorkflowRepositoryInterface interface {
	GetWorkflowList(checklistID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklist, error)
	GetWorkflowHistories(checklistID uint, tx ...utility.DatabaseTransaction) ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error)
	CreateWorkflow(workflow *models.ProjectSecurityGateChecklistExceptionWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateWorkflow(workflow *models.ProjectSecurityGateChecklistExceptionWorkflow, tx ...utility.DatabaseTransaction) error
	CreateWorkflowParticipant(participant *models.ProjectSecurityGateChecklistExceptionWorkflowParticipant, tx ...utility.DatabaseTransaction) error

	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error)
	GetWorkflowByLevel(checklistID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error)
	GetNextLevel(checklistID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
	HasUserApproved(workflowID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error)

	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectSecurityGateChecklistExceptionWorkflowRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateChecklistExceptionWorkflowRepository(db *gorm.DB) ProjectSecurityGateChecklistExceptionWorkflowRepositoryInterface {
	return &ProjectSecurityGateChecklistExceptionWorkflowRepository{db: db}
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) GetWorkflowList(checklistID uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklist, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var checklist models.ProjectSecurityGateChecklist
	result := db.
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("ExceptionWorkflows", func(db *gorm.DB) *gorm.DB {
			return db.Order("project_security_gate_checklist_exception_workflows.level ASC")
		}).
		Preload("ExceptionWorkflows.WorkflowStatus").
		Preload("ExceptionWorkflows.WorkflowLevelType").
		Preload("ExceptionWorkflows.Participants").
		Preload("ExceptionWorkflows.Participants.User").
		Preload("ExceptionWorkflows.Participants.Role").
		Preload("ExceptionWorkflows.Reviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("project_security_gate_checklist_exception_workflow_reviews.created_at ASC")
		}).
		Preload("ExceptionWorkflows.Reviews.WorkflowStatus").
		Preload("ExceptionWorkflows.Reviews.CreatedByUser").
		Where("id = ?", checklistID).
		Find(&checklist)

	if result.Error != nil {
		return nil, result.Error
	}

	return &checklist, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) GetWorkflowHistories(checklistID uint, tx ...utility.DatabaseTransaction) ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var histories []*models.ProjectSecurityGateChecklistExceptionWorkflowHistory
	result := db.
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("checklist_id = ?", checklistID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) CreateWorkflow(workflow *models.ProjectSecurityGateChecklistExceptionWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(workflow).Error; err != nil {
		return 0, err
	}

	return workflow.ID, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) UpdateWorkflow(workflow *models.ProjectSecurityGateChecklistExceptionWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(workflow).Error
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) CreateWorkflowParticipant(participant *models.ProjectSecurityGateChecklistExceptionWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(participant).Error
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var workflow models.ProjectSecurityGateChecklistExceptionWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("Participants").
		Preload("Participants.User").
		Preload("Participants.Role").
		Where("id = ?", id).
		First(&workflow).Error

	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) GetWorkflowByLevel(checklistID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var workflow models.ProjectSecurityGateChecklistExceptionWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("Participants").
		Preload("Participants.User").
		Preload("Participants.Role").
		Preload("Reviews").
		Preload("Reviews.WorkflowStatus").
		Where("checklist_id = ? AND level = ?", checklistID, level).
		First(&workflow).Error

	if err != nil {
		return nil, err
	}

	return &workflow, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) GetNextLevel(checklistID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectSecurityGateChecklistExceptionWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var workflow models.ProjectSecurityGateChecklistExceptionWorkflow
	err := db.
		Where("checklist_id = ? AND level > ?", checklistID, currentLevel).
		Order("level ASC").
		First(&workflow).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &workflow, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.ProjectSecurityGateChecklistExceptionWorkflowReview{}).
		Where("exception_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
		Distinct("created_by_id").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) HasUserApproved(workflowID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.ProjectSecurityGateChecklistExceptionWorkflowReview{}).
		Where("exception_workflow_id = ? AND created_by_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", participantIds).Delete(&models.ProjectSecurityGateChecklistExceptionWorkflowParticipant{}).Error
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", reviewIds).Delete(&models.ProjectSecurityGateChecklistExceptionWorkflowReview{}).Error
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", workflowIds).Delete(&models.ProjectSecurityGateChecklistExceptionWorkflow{}).Error
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("exception_workflow_id = ?", workflowID).Delete(&models.ProjectSecurityGateChecklistExceptionWorkflowReview{}).Error
}

func (repo *ProjectSecurityGateChecklistExceptionWorkflowRepository) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ProjectSecurityGateChecklistExceptionWorkflow{}).Where("id = ?", workflowID).Update("workflow_status_id", statusID).Error
}
