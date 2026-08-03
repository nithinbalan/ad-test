package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskResponseWorkflowRepositoryInterface interface {
	GetWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	GetWorkflowHistories(riskID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskResponseWorkflowHistory, error)
	CreateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error
	CreateRiskResponseWorkflowParticipant(workflowParticipant *models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	CreateRiskResponseReview(riskResponseReview *models.RiskResponseWorkflowReview, tx ...utility.DatabaseTransaction) error

	GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)

	GetWorkflowsForDeletion(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error)
	DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	DeleteRiskResponseWorkflowReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error
	UpdateRiskResponseWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskResponseWorkflowRepository struct {
	db *gorm.DB
}

func NewAssetRiskResponseWorkflowRepository(db *gorm.DB) AssetRiskResponseWorkflowRepositoryInterface {
	return &AssetRiskResponseWorkflowRepository{db: db}
}

func (repo *AssetRiskResponseWorkflowRepository) GetWorkflowList(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	result := db.
		Preload("WorkflowStatus").
		Preload("SubmittedBy").
		Preload("ApprovedBy").
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
		Where("risk_register_id = ?", riskID).
		Find(&riskResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskResponse, nil
}

func (repo *AssetRiskResponseWorkflowRepository) GetWorkflowHistories(riskID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskResponseWorkflowHistory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var histories []*models.RiskResponseWorkflowHistory
	result := db.
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("risk_response_id IN (?)", db.Model(&models.RiskResponse{}).Select("id").Where("risk_register_id = ?", riskID)).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *AssetRiskResponseWorkflowRepository) CreateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskResponseWorkflow).Error; err != nil {
		return 0, err
	}

	return riskResponseWorkflow.ID, nil
}

func (repo *AssetRiskResponseWorkflowRepository) UpdateRiskResponseWorkflow(riskResponseWorkflow *models.RiskResponseWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskResponseWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) CreateRiskResponseWorkflowParticipant(workflowParticipant *models.RiskResponseWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) CreateRiskResponseReview(riskResponseReview *models.RiskResponseWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskResponseReview).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) GetRiskResponseWorkflowByLevel(riskResponseID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var riskResponseWorkflow *models.RiskResponseWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("RiskResponseWorkflowParticipants").
		Preload("RiskResponseWorkflowParticipants.User").
		Preload("RiskResponseWorkflowReviews").
		Preload("RiskResponseWorkflowReviews.WorkflowStatus").
		Where("risk_response_id = ? AND level = ?", riskResponseID, currentApprovalLevel).
		Find(&riskResponseWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskResponseWorkflow, nil
}

func (repo *AssetRiskResponseWorkflowRepository) GetRiskResponseWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
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

func (repo *AssetRiskResponseWorkflowRepository) GetNextLevel(riskResponseID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskResponseWorkflow, error) {
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

func (repo *AssetRiskResponseWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
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

func (repo *AssetRiskResponseWorkflowRepository) GetWorkflowsForDeletion(riskID uint, tx ...utility.DatabaseTransaction) (*models.RiskResponse, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskResponse models.RiskResponse
	result := db.
		Preload("RiskResponseWorkflows").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowParticipants").
		Preload("RiskResponseWorkflows.RiskResponseWorkflowReviews").
		Where("risk_register_id = ?", riskID).
		Find(&riskResponse)

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskResponse, nil
}

func (repo *AssetRiskResponseWorkflowRepository) DeleteRiskResponseWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.RiskResponseWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) DeleteRiskResponseWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.RiskResponseWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) DeleteRiskResponseWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.RiskResponseWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) DeleteRiskResponseWorkflowReviewsByWorkflowID(workflowID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_response_workflow_id = ?", workflowID).Delete(&models.RiskResponseWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskResponseWorkflowRepository) UpdateRiskResponseWorkflowStatus(workflowID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskResponseWorkflow{}).Where("id = ?", workflowID).Update("workflow_status_id", statusID).Error; err != nil {
		return err
	}

	return nil
}
