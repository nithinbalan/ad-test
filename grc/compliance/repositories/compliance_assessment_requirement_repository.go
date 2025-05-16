package repositories

import (
	"fmt"
	"strings"

	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementListDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentRequirement, error)
	GetFullDetailsByID(requirementId uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirement, error)
	Update(complianceAssessmentRequirement *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) error

	CreateDocuments(documents []*models.ComplianceAssessmentRequirementRequiredDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(documentId uint, tx ...utility.DatabaseTransaction) error

	CreateComment(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) error
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

func (repo *ComplianceAssessmentRequirementRepository) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceAssessmentRequirementListDTO, error) {
	var complianceRequirements []*complianceDTOs.ComplianceAssessmentRequirementListDTO

	query := repo.db.Table("compliance_assessment_requirements").
		Select(`
			compliance_assessment_requirements.*,
			compliance_clauses.id AS clause_id,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_clauses.title_en AS compliance_clause_title_en,
			compliance_clauses.title_ar AS compliance_clause_title_ar,
			compliance_requirements.title_en AS compliance_requirement_title_en,
			compliance_requirements.title_ar AS compliance_requirement_title_ar,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			compliance_assessment_requirement_statuses.title_en AS compliance_assessment_requirement_status_title_en,
			compliance_assessment_requirement_statuses.title_ar AS compliance_assessment_requirement_status_title_ar,
			compliance_assessment_requirement_statuses.color_code AS compliance_assessment_requirement_status_color_code,
			compliance_assessment_requirement_statuses.type AS compliance_assessment_requirement_status_type
		`).
		Joins(`
			LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
			LEFT JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id
			LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_assessment_requirements.compliance_status_id
			LEFT JOIN compliance_assessment_requirement_statuses ON compliance_assessment_requirement_statuses.id = compliance_assessment_requirements.compliance_assessment_requirement_status_id

		`).
		Where("compliance_assessment_requirements.compliance_assessment_id = ?", assessmentID)

	// Apply dynamic filters if provided
	for columnName, values := range params.Filters {
		if columnName == "compliance_status_id" {
			query = query.Where(fmt.Sprintf("compliance_assessment_requirements.%s IN (?)", columnName), values)
		} else if columnName == "compliance_clause_id" {
			query = query.Where(fmt.Sprintf("compliance_requirements.%s IN (?)", columnName), values)
		} else if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "completed":
			query = query.Where("compliance_assessment_requirements.is_completed = ?", true)
		case "not_completed":
			query = query.Where("compliance_assessment_requirements.is_completed = ?", false)
		}
	}

	// Apply keyword-based search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_clauses.clause_no LIKE ? OR
			compliance_clauses.title_en LIKE ? OR
			compliance_clauses.title_ar LIKE ? OR
			compliance_requirements.title_en LIKE ? OR
			compliance_requirements.title_ar LIKE ? OR
			compliance_requirement_categories.title_en LIKE ? OR
			compliance_requirement_categories.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? )
		`, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Ensure sort parameters are clean and valid
	sortBy := params.SortBy
	sort := strings.ToUpper(strings.TrimSpace(params.Sort))
	if sort != "ASC" && sort != "DESC" {
		sort = "ASC"
	}

	// Use TRIM on the database column in the ORDER BY clause
	query = query.Order(fmt.Sprintf("TRIM(%s) %s", sortBy, sort))

	// Execute the query
	err := query.Find(&complianceRequirements).Error
	if err != nil {
		return nil, err
	}

	return complianceRequirements, nil
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
		Preload("ComplianceAssessment.SubmittedByUser").
		Preload("ComplianceAssessment.OrganizationUnit").
		Preload("ComplianceRequirement").
		Preload("ComplianceRequirement.ComplianceStatus").
		Preload("ComplianceRequirement.ComplianceClause").
		Preload("ComplianceRequirement.SeverityLevel").
		Preload("RequiredDocumentTitles").
		Preload("RequiredDocumentTitles.RequiredDocumentTitle").
		Preload("RequiredDocumentTitles.Documents").
		Preload("ActionPlans").
		Preload("ActionPlans.ResponsibleUser").
		Preload("ActionPlans.ReviewUser").
		Preload("ActionPlans.ActionPlanStatus").
		Preload("RequirementComments").
		Preload("RequirementComments.CreatedByUser").
		Preload("ComplianceStatus").
		Preload("ComplianceAssessmentRequirementStatus").
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

func (repo *ComplianceAssessmentRequirementRepository) CreateDocuments(docs []*models.ComplianceAssessmentRequirementRequiredDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRequirementRepository) DeleteDocuments(docIds uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceAssessmentRequirementRequiredDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRequirementRepository) CreateComment(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Create(comment)
	return result.Error
}
