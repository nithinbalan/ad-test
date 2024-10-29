package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditReportReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetComments(auditFinalReportID uint, currentApprovalLevel int) (*models.AuditFinalReport, error)
	CreateAuditReportReview(auditReview *models.AuditReportReview, tx ...utility.DatabaseTransaction) error
}

type AuditReportReviewRepository struct {
	db *gorm.DB
}

func NewAuditReportReviewRepository(db *gorm.DB) *AuditReportReviewRepository {
	return &AuditReportReviewRepository{db: db}
}

// BeginTransaction starts a new transaction.
func (repo *AuditReportReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}


func (repo *AuditReportReviewRepository) GetComments(auditFinalReportID uint, currentApprovalLevel int) (*models.AuditFinalReport, error) {
	var auditFinalReport models.AuditFinalReport
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("AuditReportWorkflows", "level <= ?", currentApprovalLevel).
		Preload("AuditReportWorkflows.WorkflowStatus").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants.User").
		Preload("AuditReportWorkflows.AuditReportReviews").
		Preload("AuditReportWorkflows.AuditReportReviews.WorkflowStatus").
		Preload("AuditReportWorkflows.AuditReportReviews.CreatedByUser").
		Where("id = ?", auditFinalReportID).
		Find(&auditFinalReport)

	if result.Error != nil {
		return nil, result.Error
	}

	return &auditFinalReport, nil
}

func (repo *AuditReportReviewRepository) CreateAuditReportReview(auditReview *models.AuditReportReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReview).Error; err != nil {
		return err
	}
	return nil
}
