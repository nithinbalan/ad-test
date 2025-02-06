package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetByID(id uint) (*models.ComplianceAssessmentRequirement, error)
	GetFullDetailsByID(requirementId uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirement, error)
	Update(complianceAssessmentRequirement *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentRequirementRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRequirementRepository(db *gorm.DB) *ComplianceAssessmentRequirementRepository {
	return &ComplianceAssessmentRequirementRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceAssessmentRequirementRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentRequirementRepository) GetByID(id uint) (*models.ComplianceAssessmentRequirement, error) {
	var complianceAssessmentRequirement *models.ComplianceAssessmentRequirement

	result := repo.db.First(&complianceAssessmentRequirement, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessmentRequirement, nil

}

func (repo *ComplianceAssessmentRequirementRepository) GetFullDetailsByID(requirementId uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirement, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var complianceAssessmentRequirement *models.ComplianceAssessmentRequirement

	result := db.
		Preload("ComplianceAssessment").
		Preload("ComplianceAssessment.OrganizationUnit").
		Preload("ComplianceRequirement").
		Preload("ComplianceRequirement.ComplianceStatus").
		Preload("ComplianceRequirement.ComplianceClause").
		Preload("ComplianceRequirement.SeverityLevel").
		Preload("ComplianceStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&complianceAssessmentRequirement, requirementId)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessmentRequirement, nil
}

func (repo *ComplianceAssessmentRequirementRepository) Update(complianceAssessmentRequirement *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessmentRequirement).Error; err != nil {
		return err
	}

	return nil
}
