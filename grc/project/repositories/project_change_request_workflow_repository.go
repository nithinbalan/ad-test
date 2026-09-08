package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectChangeRequestWorkflowRepositoryInterface interface {
	GetWorkflowList(projectChangeRequestID uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequest, error)
	GetWorkflowHistories(projectChangeRequestID uint) ([]*models.ProjectChangeRequestWorkflowHistory, error)
	CreateWorkflow(workflow *models.ProjectChangeRequestWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateWorkflow(workflow *models.ProjectChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error
	CreateWorkflowParticipant(participant *models.ProjectChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowParticipant(participant *models.ProjectChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error
	GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error)
	GetWorkflowByLevel(projectChangeRequestID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error)
	GetNextLevel(projectChangeRequestID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type ProjectChangeRequestWorkflowRepository struct {
	db *gorm.DB
}

func NewProjectChangeRequestWorkflowRepository(db *gorm.DB) *ProjectChangeRequestWorkflowRepository {
	return &ProjectChangeRequestWorkflowRepository{db: db}
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) GetWorkflowList(projectChangeRequestID uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequest, error) {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var projectChangeRequest models.ProjectChangeRequest
	result := db.
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("RejectedByUser").
		Preload("WorkflowStatus").
		Preload("Workflows").
		Preload("Workflows.WorkflowStatus").
		Preload("Workflows.WorkflowLevelType").
		Preload("Workflows.ProjectChangeRequestWorkflowParticipants").
		Preload("Workflows.ProjectChangeRequestWorkflowParticipants.User").
		Preload("Workflows.ProjectChangeRequestWorkflowParticipants.Role").
		Preload("Workflows.ProjectChangeRequestWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("project_change_request_workflow_reviews.created_at ASC")
		}).
		Preload("Workflows.ProjectChangeRequestWorkflowReviews.WorkflowStatus").
		Preload("Workflows.ProjectChangeRequestWorkflowReviews.CreatedBy").
		Where("id = ?", projectChangeRequestID).
		Find(&projectChangeRequest)

	if result.Error != nil {
		return nil, result.Error
	}

	return &projectChangeRequest, nil
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) GetWorkflowHistories(projectChangeRequestID uint) ([]*models.ProjectChangeRequestWorkflowHistory, error) {
	var histories []*models.ProjectChangeRequestWorkflowHistory

	result := pcrwr.db.
		Preload("WorkflowStatus").
		Preload("CreatedBy").
		Where("project_change_request_id = ?", projectChangeRequestID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) CreateWorkflow(workflow *models.ProjectChangeRequestWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflow).Error; err != nil {
		return 0, err
	}
	return workflow.ID, nil
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) UpdateWorkflow(workflow *models.ProjectChangeRequestWorkflow, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(workflow).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) CreateWorkflowParticipant(participant *models.ProjectChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(participant).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) UpdateWorkflowParticipant(participant *models.ProjectChangeRequestWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(participant).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) DeleteWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", participantIds).Delete(&models.ProjectChangeRequestWorkflowParticipant{}).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) DeleteWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", reviewIds).Delete(&models.ProjectChangeRequestWorkflowReview{}).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) DeleteReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_change_request_workflow_id = ?", workflowID).Delete(&models.ProjectChangeRequestWorkflowReview{}).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) UpdateWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ProjectChangeRequestWorkflow{}).Where("id = ?", workflowID).Update("workflow_status_id", statusID).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) DeleteWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", workflowIds).Delete(&models.ProjectChangeRequestWorkflow{}).Error
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) GetWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error) {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectChangeRequestWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("ProjectChangeRequestWorkflowParticipants").
		Preload("ProjectChangeRequestWorkflowParticipants.User").
		Preload("ProjectChangeRequestWorkflowReviews").
		Where("id = ?", id).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) GetWorkflowByLevel(projectChangeRequestID uint, level int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error) {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectChangeRequestWorkflow
	err := db.
		Preload("WorkflowLevelType").
		Preload("ProjectChangeRequestWorkflowParticipants").
		Preload("ProjectChangeRequestWorkflowParticipants.User").
		Preload("ProjectChangeRequestWorkflowParticipants.Role").
		Preload("ProjectChangeRequestWorkflowReviews").
		Where("project_change_request_id = ? AND level = ?", projectChangeRequestID, level).
		First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

func (pcrwr *ProjectChangeRequestWorkflowRepository) GetNextLevel(projectChangeRequestID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestWorkflow, error) {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.ProjectChangeRequestWorkflow
	err := db.
		Where("project_change_request_id = ? AND level > ?", projectChangeRequestID, currentLevel).
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

func (pcrwr *ProjectChangeRequestWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := pcrwr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM project_change_request_workflow_reviews
        WHERE project_change_request_workflow_id = ?
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
	query := db.Model(&models.ProjectChangeRequestWorkflowReview{}).
		Where("project_change_request_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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
