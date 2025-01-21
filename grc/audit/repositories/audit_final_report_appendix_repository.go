package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFinalReportAppendixRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.AuditFinalReportAppendix, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFinalReportAppendix, error)
	Create(auditFinalReportAppendix *models.AuditFinalReportAppendix, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFinalReportAppendix *models.AuditFinalReportAppendix, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditFinalReportAppendixRepository struct {
	db *gorm.DB
}

func NewAuditFinalReportAppendixRepository(db *gorm.DB) *AuditFinalReportAppendixRepository {
	return &AuditFinalReportAppendixRepository{
		db: db,
	}
}

func (repo *AuditFinalReportAppendixRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFinalReportAppendixRepository) GetByID(id uint) (*models.AuditFinalReportAppendix, error) {
	var appendix models.AuditFinalReportAppendix
	result := repo.db.First(&appendix, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &appendix, nil
}

func (repo *AuditFinalReportAppendixRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFinalReportAppendix, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var appendix models.AuditFinalReportAppendix
	result := db.
		Preload("AuditFinalReport").
		Preload("AuditFinalReportSection").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&appendix, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &appendix, nil
}

func (repo *AuditFinalReportAppendixRepository) Create(appendix *models.AuditFinalReportAppendix, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(appendix).Error; err != nil {
		return 0, err
	}
	return appendix.ID, nil
}

func (repo *AuditFinalReportAppendixRepository) Update(appendix *models.AuditFinalReportAppendix, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(appendix).Error
}

func (repo *AuditFinalReportAppendixRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.AuditFinalReportAppendix{}, id).Error
}
