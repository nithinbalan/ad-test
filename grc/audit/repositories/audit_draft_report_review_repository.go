package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	CreateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	Create(auditReportReview *models.AuditReportReview, tx ...utility.DatabaseTransaction) (uint, error)
	CreateAuditReportWorkflowParticipants(participants []*models.AuditReportWorkflowParticipant, tx ...utility.DatabaseTransaction) error
}

type AuditDraftReportReviewRepository struct {
	db *gorm.DB
}

func NewAuditDraftReportReviewRepository(db *gorm.DB) *AuditDraftReportReviewRepository {
	return &AuditDraftReportReviewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDraftReportReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDraftReportReviewRepository) CreateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReportWorkflow).Error; err != nil {
		return 0, err
	}

	return auditReportWorkflow.ID, nil
}

func (repo *AuditDraftReportReviewRepository) Create(auditReportReview *models.AuditReportReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReportReview).Error; err != nil {
		return 0, err
	}

	return auditReportReview.ID, nil
}


func (repo *AuditDraftReportReviewRepository) CreateAuditReportWorkflowParticipants(participants []*models.AuditReportWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&participants).Error; err != nil {
		return err
	}

	return nil
}
