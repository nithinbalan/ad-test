package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentClosureWorkflowRepositoryInterface interface {
	GetWorkflowList(incidentClosureID uint) (*models.IncidentClosure, error)
	GetWorkflowHistories(incidentClosureID uint) ([]*models.IncidentClosureWorkflowHistory, error)
	CreateIncidentClosureWorkflow(incidentClosureWorkflow *models.IncidentClosureWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateIncidentClosureWorkflow(incidentClosureWorkflow *models.IncidentClosureWorkflow, tx ...utility.DatabaseTransaction) error
	CreateIncidentClosureWorkflowParticipant(workflowParticipant *models.IncidentClosureWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteIncidentClosureWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentClosureWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentClosureWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetIncidentClosureWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error)
	GetIncidentClosureWorkflowByLevel(incidentClosureID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error)

	GetNextLevel(incidentClosureID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type IncidentClosureWorkflowRepository struct {
	db *gorm.DB
}

func NewIncidentClosureWorkflowRepository(db *gorm.DB) *IncidentClosureWorkflowRepository {
	return &IncidentClosureWorkflowRepository{db: db}
}

func (repo *IncidentClosureWorkflowRepository) GetWorkflowList(incidentClosureID uint) (*models.IncidentClosure, error) {
	var incidentClosure models.IncidentClosure
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("IncidentClosureWorkflows").
		Preload("IncidentClosureWorkflows.WorkflowStatus").
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowParticipants").
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.User").
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.Role").
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.WorkflowStatus").
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("incident_closure_workflow_reviews.created_at ASC")
		}).
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowReviews.WorkflowStatus").
		Preload("IncidentClosureWorkflows.IncidentClosureWorkflowReviews.CreatedByUser").
		Where("id = ?", incidentClosureID).
		Find(&incidentClosure)

	if result.Error != nil {
		return nil, result.Error
	}

	return &incidentClosure, nil
}

func (repo *IncidentClosureWorkflowRepository) GetWorkflowHistories(incidentClosureID uint) ([]*models.IncidentClosureWorkflowHistory, error) {
	var histories []*models.IncidentClosureWorkflowHistory

	result := repo.db.
		Preload("IncidentClosure").
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("incident_closure_id = ?", incidentClosureID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *IncidentClosureWorkflowRepository) CreateIncidentClosureWorkflow(incidentClosureWorkflow *models.IncidentClosureWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentClosureWorkflow).Error; err != nil {
		return 0, err
	}

	return incidentClosureWorkflow.ID, nil
}

func (repo *IncidentClosureWorkflowRepository) UpdateIncidentClosureWorkflow(incidentClosureWorkflow *models.IncidentClosureWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentClosureWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureWorkflowRepository) CreateIncidentClosureWorkflowParticipant(workflowParticipant *models.IncidentClosureWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureWorkflowRepository) GetIncidentClosureWorkflowByLevel(incidentClosureID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incidentClosureWorkflow *models.IncidentClosureWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("IncidentClosureWorkflowParticipants").
		Preload("IncidentClosureWorkflowParticipants.User").
		Preload("IncidentClosureWorkflowReviews").
		Preload("IncidentClosureWorkflowReviews.WorkflowStatus").
		Where("incident_closure_id = ? AND level = ?", incidentClosureID, currentApprovalLevel).
		Find(&incidentClosureWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentClosureWorkflow, nil
}

func (repo *IncidentClosureWorkflowRepository) GetIncidentClosureWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var incidentClosureWorkflow *models.IncidentClosureWorkflow

	result := db.First(&incidentClosureWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return incidentClosureWorkflow, nil
}

func (repo *IncidentClosureWorkflowRepository) GetNextLevel(incidentClosureID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.IncidentClosureWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.IncidentClosureWorkflow{}).
		Where("incident_closure_id = ? AND level > ?", incidentClosureID, currentLevel).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	var nextWorkflow models.IncidentClosureWorkflow
	result := db.Where("incident_closure_id = ? AND level > ?", incidentClosureID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *IncidentClosureWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM incident_closure_workflow_reviews
        WHERE incident_closure_workflow_id = ?
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
	query := db.Model(&models.IncidentClosureWorkflowReview{}).
		Where("incident_closure_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *IncidentClosureWorkflowRepository) DeleteIncidentClosureWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.IncidentClosureWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureWorkflowRepository) DeleteIncidentClosureWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.IncidentClosureWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureWorkflowRepository) DeleteIncidentClosureWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.IncidentClosureWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
