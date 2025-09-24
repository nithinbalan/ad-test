package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskWorkflowRepositoryInterface interface {
	GetWorkflowList(riskID uint) (*models.RiskRegister, error)
	CreateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error
	CreateRiskWorkflowParticipant(docWorkflowParticipant *models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteRiskWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	GetRiskWorkflowByLevel(riskID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)

	GetNextLevel(riskID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type RiskWorkflowRepository struct {
	db *gorm.DB
}

func NewRiskWorkflowRepository(db *gorm.DB) *RiskWorkflowRepository {
	return &RiskWorkflowRepository{db: db}
}

func (repo *RiskWorkflowRepository) GetWorkflowList(docVersionId uint) (*models.RiskRegister, error) {
	var risk models.RiskRegister
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("RiskWorkflows").
		Preload("RiskWorkflows.WorkflowStatus").
		Preload("RiskWorkflows.RiskWorkflowParticipants").
		Preload("RiskWorkflows.RiskWorkflowParticipants.User").
		Preload("RiskWorkflows.RiskWorkflowParticipants.Role").
		Preload("RiskWorkflows.RiskWorkflowReviews", "workflow_status_id IN (?)", repo.db.Model(&models.WorkflowStatus{}).Select("id").Where("type IN (?)", []string{"approve", "reject"})).
		Preload("RiskWorkflows.RiskWorkflowReviews.CreatedByUser").
		Where("id = ?", docVersionId).
		Find(&risk)

	if result.Error != nil {
		return nil, result.Error
	}

	return &risk, nil
}

func (repo *RiskWorkflowRepository) CreateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskWorkflow).Error; err != nil {
		return 0, err
	}

	return riskWorkflow.ID, nil
}

func (repo *RiskWorkflowRepository) UpdateRiskWorkflow(riskWorkflow *models.RiskWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskWorkflowRepository) CreateRiskWorkflowParticipant(docWorkflowParticipant *models.RiskWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(docWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskWorkflowRepository) GetRiskWorkflowByLevel(riskID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
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

func (repo *RiskWorkflowRepository) GetRiskWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
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

func (repo *RiskWorkflowRepository) GetNextLevel(riskID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.RiskWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.RiskWorkflow
	result := db.Where("risk_register_id = ? AND level > ?", riskID, currentLevel).
		Order("level ASC").  // Ensure we are ascending by level
		First(&nextWorkflow) // Retrieve the first result only which will be the next level

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No further levels available
	} else if result.Error != nil {
		return nil, result.Error // Other error occurred
	}

	return &nextWorkflow, nil
}

// func (repo *RiskWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}
// 	var count int64
// 	err := db.Model(&models.RiskWorkflowReview{}).
// 		Where("risk_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
// 		Distinct("created_by_id").
// 		Count(&count).Error

// 	if err != nil {
// 		return 0, err
// 	}

//		return int(count), nil
//	}

func (repo *RiskWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint // Pointer to handle NULL values
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
		resubmitID = new(uint) // Initialize to zero if needed
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

func (repo *RiskWorkflowRepository) DeleteRiskWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.RiskWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskWorkflowRepository) DeleteRiskWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.RiskWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskWorkflowRepository) DeleteRiskWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.RiskWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
