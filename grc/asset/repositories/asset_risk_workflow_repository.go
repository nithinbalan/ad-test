package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskWorkflowRepositoryInterface interface {
	GetWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	CreateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error
	CreateRiskWorkflowParticipant(riskWorkflowParticipant *models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	CreateRiskReview(riskReview *models.RiskWorkflowReview, tx ...utility.DatabaseTransaction) error

	GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetRiskWorkflowByLevel(riskID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetNextLevel(riskID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type AssetRiskWorkflowRepository struct {
	db *gorm.DB
}

func NewAssetRiskWorkflowRepository(db *gorm.DB) AssetRiskWorkflowRepositoryInterface {
	return &AssetRiskWorkflowRepository{db: db}
}

func (repo *AssetRiskWorkflowRepository) GetWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var risk models.RiskRegister
	result := db.
		Preload("RiskStatus").
		Preload("SubmittedByUser").
		Preload("RiskWorkflows").
		Preload("RiskWorkflows.WorkflowStatus").
		Preload("RiskWorkflows.RiskWorkflowParticipants").
		Preload("RiskWorkflows.RiskWorkflowParticipants.User").
		Preload("RiskWorkflows.RiskWorkflowParticipants.Role").
		Preload("RiskWorkflows.RiskWorkflowReviews", "workflow_status_id IN (?)", repo.db.Model(&models.WorkflowStatus{}).Select("id").Where("type IN (?)", []string{"approve", "reject"})).
		Preload("RiskWorkflows.RiskWorkflowReviews.WorkflowStatus").
		Preload("RiskWorkflows.RiskWorkflowReviews.CreatedByUser").
		Where("id = ?", riskID).
		Find(&risk)

	if result.Error != nil {
		return nil, result.Error
	}

	return &risk, nil
}

func (repo *AssetRiskWorkflowRepository) CreateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskWorkflow).Error; err != nil {
		return 0, err
	}

	return riskWorkflow.ID, nil
}

func (repo *AssetRiskWorkflowRepository) UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskWorkflowRepository) CreateRiskWorkflowParticipant(riskWorkflowParticipant *models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskWorkflowRepository) CreateRiskReview(riskReview *models.RiskWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskReview).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskWorkflowRepository) GetRiskWorkflowByLevel(riskID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var riskWorkflow *models.RiskWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("RiskWorkflowParticipants").
		Preload("RiskWorkflowReviews").
		Preload("RiskWorkflowReviews.WorkflowStatus").
		Preload("RiskWorkflowParticipants.User").
		Where("risk_register_id = ? AND level = ?", riskID, currentApprovalLevel).
		Find(&riskWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskWorkflow, nil
}

func (repo *AssetRiskWorkflowRepository) GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskWorkflow *models.RiskWorkflow

	result := db.First(&riskWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskWorkflow, nil
}

func (repo *AssetRiskWorkflowRepository) GetNextLevel(riskID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.RiskWorkflow
	result := db.Where("risk_register_id = ? AND level > ?", riskID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *AssetRiskWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM risk_workflow_reviews
        WHERE risk_workflow_id = ?
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
	query := db.Model(&models.RiskWorkflowReview{}).
		Where("risk_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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
