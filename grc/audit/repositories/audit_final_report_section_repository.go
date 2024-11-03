package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFinalReportSectionRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.AuditFinalReportSection, error)
	Update(auditFinalReportSection *models.AuditFinalReportSection, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetAllSubSectionsByParentID(parentID uint, tx ...utility.DatabaseTransaction) ([]models.AuditFinalReportSection, error)
}

type AuditFinalReportSectionRepository struct {
	db *gorm.DB
}

func NewAuditFinalReportSectionRepository(db *gorm.DB) *AuditFinalReportSectionRepository {
	return &AuditFinalReportSectionRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFinalReportSectionRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFinalReportSectionRepository) GetByID(id uint) (*models.AuditFinalReportSection, error) {
	var auditFinalReportSection *models.AuditFinalReportSection

	result := repo.db.First(&auditFinalReportSection, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditFinalReportSection, nil
}

func (repo *AuditFinalReportSectionRepository) Update(auditFinalReportSection *models.AuditFinalReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFinalReportSection).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFinalReportSectionRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AuditFinalReportSection{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFinalReportSectionRepository) GetAllSubSectionsByParentID(parentID uint, tx ...utility.DatabaseTransaction) ([]models.AuditFinalReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Define a slice to hold the sub-sections
	var subSections []models.AuditFinalReportSection

	// Query to find all sub-sections where the parent ID matches
	if err := db.Where("audit_final_report_section_id = ?", parentID).Order("[order] ASC").Find(&subSections).Error; err != nil {
		return nil, err
	}

	return subSections, nil
}
