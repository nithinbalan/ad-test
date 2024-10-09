package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.AuditDraftReport, error)
	GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftReport, error)
	Create(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditDraftReportRepository struct {
	db *gorm.DB
}

func NewAuditDraftReportRepository(db *gorm.DB) *AuditDraftReportRepository {
	return &AuditDraftReportRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDraftReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDraftReportRepository) GetByID(id uint) (*models.AuditDraftReport, error) {
	var auditDraftReport *models.AuditDraftReport

	result := repo.db.First(&auditDraftReport, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReport, nil
}

func (repo *AuditDraftReportRepository) GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditDraftReport *models.AuditDraftReport

	result := db.Where("audit_id = ?", auditID).
		Preload("Audit").
		Preload("Audit.AuditDraftFindings").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingStatus").
		Preload("Audit.AuditDraftFindings.Documents").
		Preload("Audit.AuditDraftFindings.AuditableItem").
		Preload("Audit.AuditDraftFindings.AuditableItem.Process").
		Preload("Audit.AuditDraftFindings.AuditableItem.Control").
		Preload("Audit.AuditDraftFindings.AuditableItem.Document").
		Preload("Audit.AuditDraftFindings.AuditableItem.IncidentRegister").
		Preload("Audit.AuditDraftFindings.AuditableItem.ComplianceRequirement").
		Preload("Audit.AuditDraftFindings.AuditableItem.RiskRegister").
		Preload("Audit.AuditDraftFindings.AuditableItem.User").
		Preload("Audit.AuditDraftFindings.AuditFindingCategory").
		Preload("Audit.AuditDraftFindings.SeverityLevel").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingConversations").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingConversations.CreatedByUser").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingConversations.UpdatedByUser").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingActionPlans").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingActionPlans.ResponsibleUser").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingActionPlans.CreatedByUser").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingActionPlans.UpdatedByUser").
		Preload("Audit.AuditPlan").
		Preload("Audit.AuditPlan.AuditCategory").
		Preload("Audit.AuditPlan.AuditStatus").
		Preload("Audit.AuditPlan.TeamLead").
		Preload("Audit.AuditPlan.Members").
		Preload("Audit.AuditPlan.Members.User").
		Preload("Audit.CreatedByUser").
		Preload("Audit.UpdatedByUser").
		Preload("AuditDraftReportStatus").
		Preload("AuditReportWorkflows").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants.User").
		Preload("InternalReviewSubmittedByUser").
		Preload("AuditeeReviewSubmittedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditDraftReport)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReport, nil
}

func (repo *AuditDraftReportRepository) Create(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditDraftReport).Error; err != nil {
		return 0, err
	}

	return auditDraftReport.ID, nil
}

func (repo *AuditDraftReportRepository) Update(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditDraftReport).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Schedule
	if err := db.Delete(&models.AuditDraftReport{}, id).Error; err != nil {
		return err
	}

	return nil
}
