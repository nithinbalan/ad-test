package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRegisterWorkflowRepositoryInterface interface {
	GetWorkflowList(complianceRegisterID uint) (*models.ComplianceRegister, error)
	GetWorkflowHistories(complianceRegisterID uint) ([]*models.ComplianceRegisterWorkflowHistory, error)
	CreateComplianceRegisterWorkflow(complianceRegisterWorkflow *models.ComplianceRegisterWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow *models.ComplianceRegisterWorkflow, tx ...utility.DatabaseTransaction) error
	CreateComplianceRegisterWorkflowParticipant(workflowParticipant *models.ComplianceRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRegisterWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetComplianceRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error)
	GetComplianceRegisterWorkflowByLevel(complianceRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error)

	GetNextLevel(complianceRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type ComplianceRegisterWorkflowRepository struct {
	db *gorm.DB
}

func NewComplianceRegisterWorkflowRepository(db *gorm.DB) *ComplianceRegisterWorkflowRepository {
	return &ComplianceRegisterWorkflowRepository{db: db}
}

func (repo *ComplianceRegisterWorkflowRepository) GetWorkflowList(complianceRegisterID uint) (*models.ComplianceRegister, error) {
	var complianceRegister models.ComplianceRegister
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("ComplianceRegisterWorkflows").
		Preload("ComplianceRegisterWorkflows.WorkflowStatus").
		Preload("ComplianceRegisterWorkflows.ComplianceRegisterWorkflowParticipants").
		Preload("ComplianceRegisterWorkflows.ComplianceRegisterWorkflowParticipants.User").
		Preload("ComplianceRegisterWorkflows.ComplianceRegisterWorkflowParticipants.Role").
		Preload("ComplianceRegisterWorkflows.ComplianceRegisterWorkflowReviews", func(db *gorm.DB) *gorm.DB {
			return db.Order("compliance_register_workflow_reviews.created_at ASC")
		}).
		Preload("ComplianceRegisterWorkflows.ComplianceRegisterWorkflowReviews.WorkflowStatus").
		Preload("ComplianceRegisterWorkflows.ComplianceRegisterWorkflowReviews.CreatedByUser").
		Where("id = ?", complianceRegisterID).
		Find(&complianceRegister)

	if result.Error != nil {
		return nil, result.Error
	}

	return &complianceRegister, nil
}

func (repo *ComplianceRegisterWorkflowRepository) GetWorkflowHistories(complianceRegisterID uint) ([]*models.ComplianceRegisterWorkflowHistory, error) {
	var histories []*models.ComplianceRegisterWorkflowHistory

	result := repo.db.
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Where("compliance_register_id = ?", complianceRegisterID).
		Order("created_at DESC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *ComplianceRegisterWorkflowRepository) CreateComplianceRegisterWorkflow(complianceRegisterWorkflow *models.ComplianceRegisterWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceRegisterWorkflow).Error; err != nil {
		return 0, err
	}

	return complianceRegisterWorkflow.ID, nil
}

func (repo *ComplianceRegisterWorkflowRepository) UpdateComplianceRegisterWorkflow(complianceRegisterWorkflow *models.ComplianceRegisterWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceRegisterWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterWorkflowRepository) CreateComplianceRegisterWorkflowParticipant(workflowParticipant *models.ComplianceRegisterWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(workflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterWorkflowRepository) GetComplianceRegisterWorkflowByLevel(complianceRegisterID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var complianceRegisterWorkflow *models.ComplianceRegisterWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("ComplianceRegisterWorkflowParticipants").
		Preload("ComplianceRegisterWorkflowReviews").
		Preload("ComplianceRegisterWorkflowReviews.WorkflowStatus").
		Preload("ComplianceRegisterWorkflowParticipants.User").
		Where("compliance_register_id = ? AND level = ?", complianceRegisterID, currentApprovalLevel).
		Find(&complianceRegisterWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRegisterWorkflow, nil
}

func (repo *ComplianceRegisterWorkflowRepository) GetComplianceRegisterWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceRegisterWorkflow *models.ComplianceRegisterWorkflow

	result := db.First(&complianceRegisterWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRegisterWorkflow, nil
}

func (repo *ComplianceRegisterWorkflowRepository) GetNextLevel(complianceRegisterID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.ComplianceRegisterWorkflow{}).
		Where("compliance_register_id = ? AND level > ?", complianceRegisterID, currentLevel).
		Count(&count).Error

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	var nextWorkflow models.ComplianceRegisterWorkflow
	result := db.Where("compliance_register_id = ? AND level > ?", complianceRegisterID, currentLevel).
		Order("level ASC").
		First(&nextWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return &nextWorkflow, nil
}

func (repo *ComplianceRegisterWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var resubmitID *uint
	err := db.Raw(`
        SELECT MAX(id)
        FROM compliance_register_workflow_reviews
        WHERE compliance_register_workflow_id = ?
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
	query := db.Model(&models.ComplianceRegisterWorkflowReview{}).
		Where("compliance_register_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
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

func (repo *ComplianceRegisterWorkflowRepository) DeleteComplianceRegisterWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.ComplianceRegisterWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterWorkflowRepository) DeleteComplianceRegisterWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.ComplianceRegisterWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterWorkflowRepository) DeleteComplianceRegisterWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.ComplianceRegisterWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
