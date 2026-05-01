package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskResponseWorkflowRepositoryInterface interface {
	GetWorkflowList(riskResponseID uint) (*models.RiskResponse, error)
	CreateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error
	CreateRiskResponseWorkflowParticipant(workflowParticipant *models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)

	GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type RiskResponseWorkflowRepository struct {
	db *gorm.DB
}

func NewRiskResponseWorkflowRepository(db *gorm.DB) *RiskResponseWorkflowRepository {
	return &RiskResponseWorkflowRepository{db: db}
}

func (repo *RiskResponseWorkflowRepository) GetWorkflowList(riskResponseID uint) (*models.RiskResponse, error) {
	var riskResponse models.RiskResponse
	result := repo.db.
		Preload("SubmittedBy").
		Preload("RiskResponseWorkflows").
		Preload("RiskResponseWorkflows.WorkflowStatus").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowParticipants").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowParticipants.User").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowParticipants.Role").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("risk_response_workflow_reviews.created_at ASC")
		}).
		Preload("RiskResponseWorkflows.RiskResponseWorkflowReviews.WorkflowStatus").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowReviews.CreatedByUser").
		Where("id = ?", riskResponseID).
		Find(&riskResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskResponse, nil
}

func (repo *RiskResponseWorkflowRepository) CreateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskResponseWorkflow).Error; err != nil {
		return 0, err
	}

	return riskResponseWorkflow.ID, nil
}

func (repo *RiskResponseWorkflowRepository) UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskResponseWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskResponseWorkflowRepository) CreateRiskResponseWorkflowParticipant(workflowParticipant *models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskResponseWorkflowRepository) GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var riskResponseWorkflow *models.RiskResponseWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("RiskResponseWorkflowParticipants").
		Preload("RiskResponseWorkflowReviews").
		Preload("RiskResponseWorkflowReviews.WorkflowStatus").
		Preload("RiskResponseWorkflowParticipants.User").
		Where("risk_response_id = ? AND level = ?", riskResponseID, currentApprovalLevel).
		Find(&riskResponseWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskResponseWorkflow, nil
}

func (repo *RiskResponseWorkflowRepository) GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponseWorkflow *models.RiskResponseWorkflow

	result := db.First(&riskResponseWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskResponseWorkflow, nil
}

func (repo *RiskResponseWorkflowRepository) GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.RiskResponseWorkflow{}).
		Where("risk_response_id = ? AND level > ?", riskResponseID, currentLevel).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	var nextWorkflow models.RiskResponseWorkflow
	result := db.Where("risk_response_id = ? AND level > ?", riskResponseID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *RiskResponseWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM risk_response_workflow_reviews
        WHERE risk_response_workflow_id = ?
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
	query := db.Model(&models.RiskResponseWorkflowReview{}).
		Where("risk_response_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *RiskResponseWorkflowRepository) DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.RiskResponseWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskResponseWorkflowRepository) DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.RiskResponseWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskResponseWorkflowRepository) DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.RiskResponseWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
