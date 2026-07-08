package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BiaProcessWorkflowRepositoryInterface interface {
	GetWorkflowList(biaProcessID uint) (*models.BiaProcess, error)
	CreateBiaProcessWorkflow(workflow *models.BiaProcessWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateBiaProcessWorkflow(workflow *models.BiaProcessWorkflow, tx ...utility.DatabaseTransaction) error
	CreateBiaProcessWorkflowParticipant(participant *models.BiaProcessWorkflowParticipant, tx ...utility.DatabaseTransaction) error

	GetBiaProcessWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error)
	GetBiaProcessWorkflowByLevel(biaProcessID uint, level int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error)
	GetNextLevel(biaProcessID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)

	DeleteBiaProcessWorkflowParticipants(participantIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteBiaProcessWorkflowReviews(reviewIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteBiaProcessWorkflows(workflowIDs []*uint, tx ...utility.DatabaseTransaction) error

	GetWorkflowHistories(biaProcessID uint) ([]*models.BiaProcessWorkflowHistory, error)
	CreateWorkflowHistory(history *models.BiaProcessWorkflowHistory, tx ...utility.DatabaseTransaction) error
}

type BiaProcessWorkflowRepository struct {
	db *gorm.DB
}

func NewBiaProcessWorkflowRepository(db *gorm.DB) BiaProcessWorkflowRepositoryInterface {
	return &BiaProcessWorkflowRepository{db: db}
}

func (repo *BiaProcessWorkflowRepository) GetWorkflowList(biaProcessID uint) (*models.BiaProcess, error) {
	var biaProcess models.BiaProcess
	result := repo.db.
		Preload("SubmittedBy").
		Preload("BiaProcessWorkflows").
		Preload("BiaProcessWorkflows.WorkflowStatus").
		Preload("BiaProcessWorkflows.BiaProcessWorkflowParticipants").
		Preload("BiaProcessWorkflows.BiaProcessWorkflowParticipants.User").
		Preload("BiaProcessWorkflows.BiaProcessWorkflowParticipants.Role").
		Preload("BiaProcessWorkflows.BiaProcessWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("bia_process_workflow_reviews.created_at ASC")
		}).
		Preload("BiaProcessWorkflows.BiaProcessWorkflowReviews.WorkflowStatus").
		Preload("BiaProcessWorkflows.BiaProcessWorkflowReviews.CreatedByUser").
		Where("id = ?", biaProcessID).
		Find(&biaProcess)

	if result.Error != nil {
		return nil, result.Error
	}

	return &biaProcess, nil
}

func (repo *BiaProcessWorkflowRepository) CreateBiaProcessWorkflow(workflow *models.BiaProcessWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflow).Error; err != nil {
		return 0, err
	}
	return workflow.ID, nil
}

func (repo *BiaProcessWorkflowRepository) UpdateBiaProcessWorkflow(workflow *models.BiaProcessWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(workflow).Error
}

func (repo *BiaProcessWorkflowRepository) CreateBiaProcessWorkflowParticipant(participant *models.BiaProcessWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(participant).Error
}

func (repo *BiaProcessWorkflowRepository) GetBiaProcessWorkflowByLevel(biaProcessID uint, level int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.BiaProcessWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("BiaProcessWorkflowParticipants").
		Preload("BiaProcessWorkflowParticipants.User").
		Preload("BiaProcessWorkflowReviews").
		Preload("BiaProcessWorkflowReviews.WorkflowStatus").
		Where("bia_process_id = ? AND level = ?", biaProcessID, level).
		Find(&workflow)

	if result.Error != nil {
		return nil, result.Error
	}
	return &workflow, nil
}

func (repo *BiaProcessWorkflowRepository) GetBiaProcessWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var workflow models.BiaProcessWorkflow
	result := db.First(&workflow, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &workflow, nil
}

func (repo *BiaProcessWorkflowRepository) GetNextLevel(biaProcessID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.BiaProcessWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.BiaProcessWorkflow{}).
		Where("bia_process_id = ? AND level > ?", biaProcessID, currentLevel).
		Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}

	var nextWorkflow models.BiaProcessWorkflow
	result := db.Where("bia_process_id = ? AND level > ?", biaProcessID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)
	if result.Error != nil {
		return nil, result.Error
	}
	return &nextWorkflow, nil
}

func (repo *BiaProcessWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
		SELECT MAX(id)
		FROM bia_process_workflow_reviews
		WHERE bia_process_workflow_id = ?
		  AND workflow_status_id = (
		      SELECT id FROM workflow_statuses WHERE type = 'resubmit'
		  )
	`, workflowID).Scan(&resubmitID).Error
	if err != nil {
		return 0, err
	}
	if resubmitID == nil {
		resubmitID = new(uint)
	}

	var count int64
	query := db.Model(&models.BiaProcessWorkflowReview{}).
		Where("bia_process_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *BiaProcessWorkflowRepository) DeleteBiaProcessWorkflowParticipants(participantIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", participantIDs).Delete(&models.BiaProcessWorkflowParticipant{}).Error
}

func (repo *BiaProcessWorkflowRepository) DeleteBiaProcessWorkflowReviews(reviewIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", reviewIDs).Delete(&models.BiaProcessWorkflowReview{}).Error
}

func (repo *BiaProcessWorkflowRepository) DeleteBiaProcessWorkflows(workflowIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", workflowIDs).Delete(&models.BiaProcessWorkflow{}).Error
}

func (repo *BiaProcessWorkflowRepository) GetWorkflowHistories(biaProcessID uint) ([]*models.BiaProcessWorkflowHistory, error) {
	var histories []*models.BiaProcessWorkflowHistory
	result := repo.db.
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("bia_process_id = ?", biaProcessID).
		Order("created_at DESC").
		Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil
}

func (repo *BiaProcessWorkflowRepository) CreateWorkflowHistory(history *models.BiaProcessWorkflowHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(history).Error
}
