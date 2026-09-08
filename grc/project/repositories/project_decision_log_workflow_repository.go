package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectDecisionLogWorkflowRepositoryInterface interface {
	GetWorkflowList(projectDecisionLogID uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLog, error)
	GetWorkflowHistories(projectDecisionLogID uint) ([]*models.ProjectDecisionLogWorkflowHistory, error)
	CreateWorkflow(workflow *models.ProjectDecisionLogWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateWorkflow(workflow *models.ProjectDecisionLogWorkflow, tx ...utility.DatabaseTransaction) error
	CreateWorkflowParticipant(participant *models.ProjectDecisionLogWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowParticipant(participant *models.ProjectDecisionLogWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error)
	GetWorkflowByLevel(projectDecisionLogID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error)
	GetNextLevel(projectDecisionLogID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type ProjectDecisionLogWorkflowRepository struct {
	db *gorm.DB
}

func NewProjectDecisionLogWorkflowRepository(db *gorm.DB) *ProjectDecisionLogWorkflowRepository {
	return &ProjectDecisionLogWorkflowRepository{db: db}
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) GetWorkflowList(projectDecisionLogID uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLog, error) {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var projectDecisionLog models.ProjectDecisionLog
	result := db.
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("RejectedByUser").
		Preload("WorkflowStatus").
		Preload("Workflows").
		Preload("Workflows.WorkflowStatus").
		Preload("Workflows.WorkflowLevelType").
		Preload("Workflows.ProjectDecisionLogWorkflowParticipants").
		Preload("Workflows.ProjectDecisionLogWorkflowParticipants.User").
		Preload("Workflows.ProjectDecisionLogWorkflowParticipants.Role").
		Preload("Workflows.ProjectDecisionLogWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("project_decision_log_workflow_reviews.created_at ASC")
		}).
		Preload("Workflows.ProjectDecisionLogWorkflowReviews.WorkflowStatus").
		Preload("Workflows.ProjectDecisionLogWorkflowReviews.CreatedBy").
		Where("id = ?", projectDecisionLogID).
		Find(&projectDecisionLog)

	if result.Error != nil {
		return nil, result.Error
	}

	return &projectDecisionLog, nil
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) GetWorkflowHistories(projectDecisionLogID uint) ([]*models.ProjectDecisionLogWorkflowHistory, error) {
	var histories []*models.ProjectDecisionLogWorkflowHistory

	result := pdlwr.db.
		Preload("WorkflowStatus").
		Preload("CreatedBy").
		Where("project_decision_log_id = ?", projectDecisionLogID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) CreateWorkflow(workflow *models.ProjectDecisionLogWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflow).Error; err != nil {
		return 0, err
	}
	return workflow.ID, nil
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) UpdateWorkflow(workflow *models.ProjectDecisionLogWorkflow, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(workflow).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) CreateWorkflowParticipant(participant *models.ProjectDecisionLogWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(participant).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) UpdateWorkflowParticipant(participant *models.ProjectDecisionLogWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(participant).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", participantIds).Delete(&models.ProjectDecisionLogWorkflowParticipant{}).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", reviewIds).Delete(&models.ProjectDecisionLogWorkflowReview{}).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_decision_log_workflow_id = ?", workflowID).Delete(&models.ProjectDecisionLogWorkflowReview{}).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ProjectDecisionLogWorkflow{}).Where("id = ?", workflowID).Update("workflow_status_id", statusID).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", workflowIds).Delete(&models.ProjectDecisionLogWorkflow{}).Error
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error) {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectDecisionLogWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("ProjectDecisionLogWorkflowParticipants").
		Preload("ProjectDecisionLogWorkflowParticipants.User").
		Preload("ProjectDecisionLogWorkflowParticipants.Role").
		Where("id = ?", id).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) GetWorkflowByLevel(projectDecisionLogID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error) {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectDecisionLogWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("ProjectDecisionLogWorkflowParticipants").
		Preload("ProjectDecisionLogWorkflowParticipants.User").
		Preload("ProjectDecisionLogWorkflowParticipants.Role").
		Where("project_decision_log_id = ? AND level = ?", projectDecisionLogID, level).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (pdlwr *ProjectDecisionLogWorkflowRepository) GetNextLevel(projectDecisionLogID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogWorkflow, error) {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectDecisionLogWorkflow
	err := db.
		Where("project_decision_log_id = ? AND level > ?", projectDecisionLogID, currentLevel).
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

func (pdlwr *ProjectDecisionLogWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := pdlwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM project_decision_log_workflow_reviews
        WHERE project_decision_log_workflow_id = ?
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
	query := db.Model(&models.ProjectDecisionLogWorkflowReview{}).
		Where("project_decision_log_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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
