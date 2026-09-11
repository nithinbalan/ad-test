package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectWorkflowRepositoryInterface interface {
	GetWorkflowList(projectID uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	GetWorkflowHistories(projectID uint) ([]*models.ProjectWorkflowHistory, error)
	CreateWorkflow(workflow *models.ProjectWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateWorkflow(workflow *models.ProjectWorkflow, tx ...utility.DatabaseTransaction) error
	CreateWorkflowParticipant(participant *models.ProjectWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowParticipant(participant *models.ProjectWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowHistories(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error)
	GetWorkflowByLevel(projectID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error)
	GetNextLevel(projectID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type ProjectWorkflowRepository struct {
	db *gorm.DB
}

func NewProjectWorkflowRepository(db *gorm.DB) *ProjectWorkflowRepository {
	return &ProjectWorkflowRepository{db: db}
}

func (pwr *ProjectWorkflowRepository) GetWorkflowList(projectID uint, tx ...utility.DatabaseTransaction) (*models.Project, error) {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var project models.Project
	result := db.
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("RejectedByUser").
		Preload("WorkflowStatus").
		Preload("Workflows").
		Preload("Workflows.WorkflowStatus").
		Preload("Workflows.WorkflowLevelType").
		Preload("Workflows.ProjectWorkflowParticipants").
		Preload("Workflows.ProjectWorkflowParticipants.User").
		Preload("Workflows.ProjectWorkflowParticipants.Role").
		Preload("Workflows.ProjectWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("project_workflow_reviews.created_at ASC")
		}).
		Preload("Workflows.ProjectWorkflowReviews.WorkflowStatus").
		Preload("Workflows.ProjectWorkflowReviews.CreatedBy").
		Where("id = ?", projectID).
		Find(&project)

	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (pwr *ProjectWorkflowRepository) GetWorkflowHistories(projectID uint) ([]*models.ProjectWorkflowHistory, error) {
	var histories []*models.ProjectWorkflowHistory

	result := pwr.db.
		Preload("WorkflowStatus").
		Preload("CreatedBy").
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (pwr *ProjectWorkflowRepository) CreateWorkflow(workflow *models.ProjectWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflow).Error; err != nil {
		return 0, err
	}
	return workflow.ID, nil
}

func (pwr *ProjectWorkflowRepository) UpdateWorkflow(workflow *models.ProjectWorkflow, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(workflow).Error
}

func (pwr *ProjectWorkflowRepository) CreateWorkflowParticipant(participant *models.ProjectWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(participant).Error
}

func (pwr *ProjectWorkflowRepository) UpdateWorkflowParticipant(participant *models.ProjectWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(participant).Error
}

func (pwr *ProjectWorkflowRepository) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", participantIds).Delete(&models.ProjectWorkflowParticipant{}).Error
}

func (pwr *ProjectWorkflowRepository) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", reviewIds).Delete(&models.ProjectWorkflowReview{}).Error
}

func (pwr *ProjectWorkflowRepository) DeleteWorkflowHistories(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_workflow_id IN (?)", workflowIds).Delete(&models.ProjectWorkflowHistory{}).Error
}

func (pwr *ProjectWorkflowRepository) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_workflow_id = ?", workflowID).Delete(&models.ProjectWorkflowReview{}).Error
}

func (pwr *ProjectWorkflowRepository) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ProjectWorkflow{}).Where("id = ?", workflowID).Update("workflow_status_id", statusID).Error
}

func (pwr *ProjectWorkflowRepository) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", workflowIds).Delete(&models.ProjectWorkflow{}).Error
}

func (pwr *ProjectWorkflowRepository) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error) {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("ProjectWorkflowParticipants").
		Preload("ProjectWorkflowParticipants.User").
		Preload("ProjectWorkflowReviews").
		Where("id = ?", id).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (pwr *ProjectWorkflowRepository) GetWorkflowByLevel(projectID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error) {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("ProjectWorkflowParticipants").
		Preload("ProjectWorkflowParticipants.User").
		Preload("ProjectWorkflowParticipants.Role").
		Preload("ProjectWorkflowReviews").
		Where("project_id = ? AND level = ?", projectID, level).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (pwr *ProjectWorkflowRepository) GetNextLevel(projectID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectWorkflow, error) {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectWorkflow
	err := db.
		Where("project_id = ? AND level > ?", projectID, currentLevel).
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

func (pwr *ProjectWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := pwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM project_workflow_reviews
        WHERE project_workflow_id = ?
          AND workflow_status_id = (
              SELECT id
              FROM workflow_statuses
              WHERE type = 'resubmit'
          )
    `, workflowID).Scan(&resubmitID).Error

	if err != nil {
		return 0, err
	}

	if resubmitID == nil {
		resubmitID = new(uint)
	}

	var count int64
	query := db.Model(&models.ProjectWorkflowReview{}).
		Where("project_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
		Distinct("created_by_id")

	if *resubmitID > 0 {
		query = query.Where("id > ?", *resubmitID)
	}

	err = query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
