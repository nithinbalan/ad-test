package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanWorkflowRepositoryInterface interface {
	GetWorkflowList(riskRegisterUpdatePlanID uint) (*models.RiskRegisterUpdatePlan, error)
	GetWorkflowHistories(riskRegisterUpdatePlanID uint) ([]*models.RiskRegisterUpdatePlanWorkflowHistory, error)
	CreateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow, tx ...utility.DatabaseTransaction) error
	CreateRiskRegisterUpdatePlanWorkflowParticipant(workflowParticipant *models.RiskRegisterUpdatePlanWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteRiskRegisterUpdatePlanWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskRegisterUpdatePlanWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskRegisterUpdatePlanWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetRiskRegisterUpdatePlanWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error)
	GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error)

	GetNextLevel(riskRegisterUpdatePlanID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type RiskRegisterUpdatePlanWorkflowRepository struct {
	db *gorm.DB
}

func NewRiskRegisterUpdatePlanWorkflowRepository(db *gorm.DB) *RiskRegisterUpdatePlanWorkflowRepository {
	return &RiskRegisterUpdatePlanWorkflowRepository{db: db}
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) GetWorkflowList(riskRegisterUpdatePlanID uint) (*models.RiskRegisterUpdatePlan, error) {
	var riskRegisterUpdatePlan models.RiskRegisterUpdatePlan
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("RiskRegisterUpdatePlanWorkflows").
		Preload("RiskRegisterUpdatePlanWorkflows.WorkflowStatus").
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowParticipants").
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowParticipants.User").
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowParticipants.Role").
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowParticipants.WorkflowStatus").
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("risk_update_plan_workflow_reviews.created_at ASC")
		}).
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowReviews.WorkflowStatus").
		Preload("RiskRegisterUpdatePlanWorkflows.RiskRegisterUpdatePlanWorkflowReviews.CreatedByUser").
		Where("id = ?", riskRegisterUpdatePlanID).
		Find(&riskRegisterUpdatePlan)

	if result.Error != nil {
		return nil, result.Error
	}

	return &riskRegisterUpdatePlan, nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) GetWorkflowHistories(riskRegisterUpdatePlanID uint) ([]*models.RiskRegisterUpdatePlanWorkflowHistory, error) {
	var histories []*models.RiskRegisterUpdatePlanWorkflowHistory

	result := repo.db.
		Preload("RiskRegisterUpdatePlan").
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("risk_register_update_plan_id = ?", riskRegisterUpdatePlanID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) CreateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskRegisterUpdatePlanWorkflow).Error; err != nil {
		return 0, err
	}

	return riskRegisterUpdatePlanWorkflow.ID, nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) UpdateRiskRegisterUpdatePlanWorkflow(riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskRegisterUpdatePlanWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) CreateRiskRegisterUpdatePlanWorkflowParticipant(workflowParticipant *models.RiskRegisterUpdatePlanWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) GetRiskRegisterUpdatePlanWorkflowByLevel(riskRegisterUpdatePlanID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("RiskRegisterUpdatePlanWorkflowParticipants").
		Preload("RiskRegisterUpdatePlanWorkflowReviews").
		Preload("RiskRegisterUpdatePlanWorkflowReviews.WorkflowStatus").
		Preload("RiskRegisterUpdatePlanWorkflowParticipants.User").
		Where("risk_register_update_plan_id = ? AND level = ?", riskRegisterUpdatePlanID, currentApprovalLevel).
		Find(&riskRegisterUpdatePlanWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskRegisterUpdatePlanWorkflow, nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) GetRiskRegisterUpdatePlanWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskRegisterUpdatePlanWorkflow *models.RiskRegisterUpdatePlanWorkflow

	result := db.First(&riskRegisterUpdatePlanWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskRegisterUpdatePlanWorkflow, nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) GetNextLevel(riskRegisterUpdatePlanID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlanWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.RiskRegisterUpdatePlanWorkflow{}).
		Where("risk_register_update_plan_id = ? AND level > ?", riskRegisterUpdatePlanID, currentLevel).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	var nextWorkflow models.RiskRegisterUpdatePlanWorkflow
	result := db.Where("risk_register_update_plan_id = ? AND level > ?", riskRegisterUpdatePlanID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM risk_register_update_plan_workflow_reviews
        WHERE risk_register_update_plan_workflow_id = ?
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
	query := db.Model(&models.RiskRegisterUpdatePlanWorkflowReview{}).
		Where("risk_register_update_plan_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *RiskRegisterUpdatePlanWorkflowRepository) DeleteRiskRegisterUpdatePlanWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.RiskRegisterUpdatePlanWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) DeleteRiskRegisterUpdatePlanWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.RiskRegisterUpdatePlanWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanWorkflowRepository) DeleteRiskRegisterUpdatePlanWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.RiskRegisterUpdatePlanWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
