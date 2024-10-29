package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditReportWorkflowRepositoryInterface interface {
	GetWorkflowList(auditFinalReportID uint) (*models.AuditFinalReport, error)
	CreateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) error
	CreateAuditReportWorkflowParticipant(auditWorkflowParticipant *models.AuditReportWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteAuditReportWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAuditReportWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAuditReportWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	GetAuditReportWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)
	GetAuditReportWorkflowByLevel(auditFinalReportID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)

	GetNextLevel(auditFinalReportID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
}

type AuditReportWorkflowRepository struct {
	db *gorm.DB
}

func NewAuditReportWorkflowRepository(db *gorm.DB) *AuditReportWorkflowRepository {
	return &AuditReportWorkflowRepository{db: db}
}

func (repo *AuditReportWorkflowRepository) GetWorkflowList(auditFinalReportID uint) (*models.AuditFinalReport, error) {
	var auditFinalReport models.AuditFinalReport
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("AuditReportWorkflows").
		Preload("AuditReportWorkflows.WorkflowStatus").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants.User").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants.Role").
		Preload("AuditReportWorkflows.AuditReportReviews", "workflow_status_id IN (?)", repo.db.Model(&models.WorkflowStatus{}).Select("id").Where("type = ?", "approve")).
		Preload("AuditReportWorkflows.AuditReportReviews.CreatedByUser").
		Where("id = ?", auditFinalReportID).
		Find(&auditFinalReport)

	if result.Error != nil {
		return nil, result.Error
	}

	return &auditFinalReport, nil
}

func (repo *AuditReportWorkflowRepository) CreateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReportWorkflow).Error; err != nil {
		return 0, err
	}

	return auditReportWorkflow.ID, nil
}

func (repo *AuditReportWorkflowRepository) UpdateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditReportWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditReportWorkflowRepository) CreateAuditReportWorkflowParticipant(auditWorkflowParticipant *models.AuditReportWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditReportWorkflowRepository) GetAuditReportWorkflowByLevel(auditFinalReportID uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditReportWorkflow *models.AuditReportWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("AuditReportWorkflowParticipants").
		Preload("AuditReportWorkflowReviews").
		Preload("AuditReportWorkflowParticipants.User").
		Where("document_version_id = ? AND level = ?", auditFinalReportID, currentApprovalLevel).
		Find(&auditReportWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditReportWorkflow, nil
}

func (repo *AuditReportWorkflowRepository) GetAuditReportWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var auditReportWorkflow *models.AuditReportWorkflow

	result := db.First(&auditReportWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditReportWorkflow, nil
}

func (repo *AuditReportWorkflowRepository) GetNextLevel(auditFinalReportID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.AuditReportWorkflow
	result := db.Where("audit_final_report_id = ? AND level > ?", auditFinalReportID, currentLevel).
		Order("level ASC").  // Ensure we are ascending by level
		First(&nextWorkflow) // Retrieve the first result only which will be the next level

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No further levels available
	} else if result.Error != nil {
		return nil, result.Error // Other error occurred
	}

	return &nextWorkflow, nil
}

func (repo *AuditReportWorkflowRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	err := db.Model(&models.AuditReportReview{}).
		Where("audit_report_workflow_id = ? AND workflow_status_id = (SELECT id FROM workflow_statuses WHERE type = 'approve')", workflowID).
		Distinct("created_by_id").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *AuditReportWorkflowRepository) DeleteAuditReportWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.AuditReportWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditReportWorkflowRepository) DeleteAuditReportWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.AuditReportReview{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditReportWorkflowRepository) DeleteAuditReportWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.AuditReportWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}
