package repositories

import (
	"fmt"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AuditFinalReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GenerateReferenceNumber() (string, error)
	GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*models.AuditFinalReport, error)
	GetByAuditID(auditID uint) (*models.AuditFinalReport, error)
	Create(auditFinalReport *models.AuditFinalReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateAuditFindingActionPlans(actionPlans []*models.AuditFindingActionPlan, tx ...utility.DatabaseTransaction) error
	CreateAuditFinalReportSections(sections []*models.AuditFinalReportSection, tx ...utility.DatabaseTransaction) error
	GetSectionByType(sectionType string, tx ...utility.DatabaseTransaction) (*models.AuditFinalReportSection, error)
	DeleteAuditFinalReportSections(sectionIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAuditFindingActionPlans(deleteActionPlanIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAuditFindings(deleteFindingIDs []*uint, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditFinalReportRepository struct {
	db *gorm.DB
}

func NewAuditFinalReportRepository(db *gorm.DB) *AuditFinalReportRepository {
	return &AuditFinalReportRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFinalReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFinalReportRepository) GetFullDetailsByID(auditID uint, tx ...utility.DatabaseTransaction) (*models.AuditFinalReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditFinalReport *models.AuditFinalReport

	result := db.Where("audit_id = ?", auditID).
		Preload("Audit").
		Preload("Audit.AuditFindings").
		Preload("Audit.AuditFindings.DepartmentResponseByUser").
		Preload("Audit.AuditFindings.Documents").
		Preload("Audit.AuditFindings.AuditFindingCategory").
		Preload("Audit.AuditFindings.SeverityLevel").
		Preload("Audit.AuditFindings.AuditFindingActionPlans").
		Preload("Audit.AuditFindings.AuditFindingActionPlans.ResponsibleUser").
		Preload("Audit.AuditFindings.AuditFindingActionPlans.ActionPlanStatus").
		Preload("Audit.AuditFindings.AuditableItem").
		Preload("Audit.AuditFindings.AuditableItem.ResponsibleUnit").
		Preload("Audit.AuditFindings.AuditableItem.Process").
		Preload("Audit.AuditFindings.AuditableItem.Control").
		Preload("Audit.AuditFindings.AuditableItem.Document").
		Preload("Audit.AuditFindings.AuditableItem.IncidentRegister").
		Preload("Audit.AuditFindings.AuditableItem.ComplianceRequirement").
		Preload("Audit.AuditFindings.AuditableItem.RiskRegister").
		Preload("Audit.AuditFindings.AuditableItem.User").
		Preload("Audit.AuditFindings.AuditableItem.ResponsibleUnit").
		Preload("Audit.AuditPlan").
		Preload("Audit.AuditPlan.AuditableItems").
		Preload("AuditFinalReportStatus").
		Preload("AuditFinalReportSections").
		Preload("AuditFinalReportSections.ParentSection").
		Preload("AuditFinalReportSections.CreatedByUser").
		Preload("AuditFinalReportSections.UpdatedByUser").
		Preload("AuditReportWorkflows").
		Preload("AuditReportWorkflows.AuditReportWorkflowParticipants").
		Preload("AuditReportWorkflows.AuditReportReviews").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditFinalReport)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFinalReport, nil
}

func (repo *AuditFinalReportRepository) Create(auditFinalReport *models.AuditFinalReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFinalReport).Error; err != nil {
		return 0, err
	}

	return auditFinalReport.ID, nil
}

func (repo *AuditFinalReportRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AuditPlan{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "AFR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *AuditFinalReportRepository) GetByAuditID(auditID uint) (*models.AuditFinalReport, error) {
	var auditFinalReport models.AuditFinalReport // Not a pointer, so it won't be nil on not found

	result := repo.db.Where("audit_id = ?", auditID).First(&auditFinalReport)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// No final report found for this audit
			return nil, nil
		}
		// Some other error occurred
		return nil, result.Error
	}

	return &auditFinalReport, nil
}

func (repo *AuditFinalReportRepository) CreateAuditFindingActionPlans(actionPlans []*models.AuditFindingActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&actionPlans).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFinalReportRepository) CreateAuditFinalReportSections(sections []*models.AuditFinalReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Insert the sections and automatically update their IDs
	if err := db.Create(&sections).Error; err != nil {
		return err
	}

	// The IDs of the inserted sections should now be available in each section object.
	return nil
}

func (repo *AuditFinalReportRepository) GetSectionByType(sectionType string, tx ...utility.DatabaseTransaction) (*models.AuditFinalReportSection, error) {
	var section models.AuditFinalReportSection
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("type = ? AND audit_final_report_section_id IS NULL", sectionType).First(&section).Error; err != nil {
		return nil, err
	}

	return &section, nil
}

func (repo *AuditFinalReportRepository) DeleteAuditFinalReportSections(sectionIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", sectionIDs).Delete(&models.AuditFinalReportSection{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFinalReportRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AuditFinalReport{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFinalReportRepository) DeleteAuditFindingActionPlans(deleteActionPlanIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", deleteActionPlanIDs).Delete(&models.AuditFindingActionPlan{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFinalReportRepository) DeleteAuditFindings(deleteFindingIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", deleteFindingIDs).Delete(&models.AuditFinding{}).Error; err != nil {
		return err
	}

	return nil
}
