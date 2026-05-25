package repositories

import (
	"fmt"
	"sort"
	"strings"

	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAllClauses(assessmentID uint, currentUserID uint) ([]*models.ComplianceClause, error)
	GetClauseRequirementCounts(assessmentID uint) (map[uint]complianceDTOs.ClauseRequirementCountDTO, error)
	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementListResponseDTO, error)
	GetResponsibleUnitRowsByAssessmentRequirementIDs(assessmentRequirementIDs []uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO, out *[]complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitListRowDTO) error
	GetByID(id uint) (*models.ComplianceAssessmentRequirement, error)
	GetFullDetailsByID(requirementId uint, responsibleUnitID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirement, error)
	Update(complianceAssessmentRequirement *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) error

	CreateDocuments(documents []*models.ComplianceAssessmentRequirementRequiredDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(documentId uint, tx ...utility.DatabaseTransaction) error

	CreateComment(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) error
	GetComplianceAssessmentRequirementResponsibleUnitByID(id uint) (*models.ComplianceAssessmentRequirementResponsibleUnit, error)
	UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit *models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error
	ComplianceAssessmentUpdateStatus(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error
	UpdateComplianceStatus(complianceRegisterID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateComplianceAssessmentRequirementStatus(assessmentReqID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateClauseComplianceStatus(clauseID uint, statusID uint, score float64, tx ...utility.DatabaseTransaction) error
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

func (repo *ComplianceAssessmentRequirementRepository) GetAllClauses(assessmentID uint, currentUserID uint) ([]*models.ComplianceClause, error) {
	var clauses []*models.ComplianceClause

	// Return all clauses with hierarchical structure for users with full access
	err := repo.db.
		Preload("ChildClauses", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ChildClauses") // Recursive preload
		}).
		Joins("JOIN compliance_register ON compliance_register.id = compliance_clauses.compliance_register_id").
		Joins("JOIN compliance_assessments ON compliance_assessments.compliance_register_id = compliance_register.id").
		Where("compliance_assessments.id = ?", assessmentID).
		Find(&clauses).Error

	if err != nil {
		return nil, err
	}

	// Sort clauses correctly using numeric sorting
	sort.SliceStable(clauses, func(i, j int) bool {
		return compareClauseNumbers(clauses[i].ClauseNo, clauses[j].ClauseNo)
	})

	// Recursively sort child clauses
	for _, clause := range clauses {
		sortChildClauses(clause)
	}

	return clauses, nil
}

func (repo *ComplianceAssessmentRequirementRepository) GetClauseRequirementCounts(assessmentID uint) (map[uint]complianceDTOs.ClauseRequirementCountDTO, error) {
	var results []complianceDTOs.ClauseRequirementCountDTO

	err := repo.db.Table("compliance_clauses cc").
		Select(`
			cc.id as id,
			COUNT(carru.id) as total_requirements,
			SUM(CASE WHEN carru.is_completed = 1 THEN 1 ELSE 0 END) as completed_requirements
		`).
		Joins("JOIN compliance_requirements cr ON cr.compliance_clause_id = cc.id").
		Joins("JOIN compliance_assessment_requirements car ON car.compliance_requirement_id = cr.id").
		Joins("JOIN compliance_assessments ca ON ca.id = car.compliance_assessment_id").
		Joins("LEFT JOIN compliance_assessment_requirement_responsible_units carru ON carru.compliance_assessment_requirement_id = car.id").
		Where("ca.id = ?", assessmentID).
		Group("cc.id").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	countMap := make(map[uint]complianceDTOs.ClauseRequirementCountDTO)
	for _, r := range results {
		countMap[r.ID] = r
	}
	return countMap, nil
}

func (repo *ComplianceAssessmentRequirementRepository) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementListResponseDTO, error) {
	var complianceRequirements []*complianceDTOs.ComplianceAssessmentRequirementListDTO
	var total int64

	query := repo.db.Table("compliance_assessment_requirements").
		Select(`
			compliance_assessment_requirements.*,
			compliance_clauses.id AS clause_id,
			compliance_clauses.description_en AS clause_description_en,
			compliance_clauses.description_ar AS clause_description_ar,
			compliance_clauses.clause_no AS compliance_clause_no,
			compliance_clauses.title_en AS compliance_clause_title_en,
			compliance_clauses.title_ar AS compliance_clause_title_ar,
			compliance_requirements.id AS compliance_requirement_id,
			compliance_requirements.title_en AS compliance_requirement_title_en,
			compliance_requirements.title_ar AS compliance_requirement_title_ar,
			compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
			compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
			severity_levels.id AS severity_level_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.type AS severity_level_type,
			compliance_assessment_requirement_statuses.title_en AS compliance_assessment_requirement_status_title_en,
			compliance_assessment_requirement_statuses.title_ar AS compliance_assessment_requirement_status_title_ar,
			compliance_assessment_requirement_statuses.color_code AS compliance_assessment_requirement_status_color_code,
			compliance_assessment_requirement_statuses.type AS compliance_assessment_requirement_status_type,
			average_residual_rating.id AS average_residual_risk_rating_id,
			average_residual_rating.title_en AS average_residual_risk_rating_title_en,
			average_residual_rating.title_ar AS average_residual_risk_rating_title_ar,
			average_residual_rating.color_code AS average_residual_risk_rating_color_code,
			average_residual_rating.type AS average_residual_risk_rating_type
		`).
		Joins(`
			LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
			LEFT JOIN compliance_requirement_categories ON compliance_requirement_categories.id = compliance_requirements.compliance_requirement_category_id
			LEFT JOIN severity_levels ON severity_levels.id = compliance_requirements.severity_level_id
			LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
			LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_assessment_requirements.compliance_status_id
			LEFT JOIN compliance_assessment_requirement_statuses ON compliance_assessment_requirement_statuses.id = compliance_assessment_requirements.compliance_assessment_requirement_status_id
			LEFT JOIN compliance_assessment_ratings AS average_residual_rating ON average_residual_rating.id = compliance_assessment_requirements.average_residual_risk_rating_id

		`).
		Where("compliance_assessment_requirements.compliance_assessment_id = ?", assessmentID)

	// Apply clause filter if provided
	if params.ClauseID != 0 {
		query = query.Where("compliance_requirements.compliance_clause_id = ?", params.ClauseID)
	}
	if params.IsRegulationNotificationAssessment != nil {
		query = query.Where(
			`EXISTS (
				SELECT 1
				FROM compliance_assessment_requirement_responsible_units carru
				JOIN compliance_requirement_responsible_units crru
					ON crru.compliance_requirement_id = compliance_assessment_requirements.compliance_requirement_id
					AND crru.responsible_unit_id = carru.responsible_unit_id
				WHERE carru.compliance_assessment_requirement_id = compliance_assessment_requirements.id
				  AND crru.is_regulation_notification_assessment = ?
			)`,
			*params.IsRegulationNotificationAssessment,
		)
	}
	if params.IsRiskBasedAssessment != nil {
		query = query.Where(
			`EXISTS (
				SELECT 1
				FROM compliance_assessment_requirement_responsible_units carru
				JOIN compliance_requirement_responsible_units crru
					ON crru.compliance_requirement_id = compliance_assessment_requirements.compliance_requirement_id
					AND crru.responsible_unit_id = carru.responsible_unit_id
				WHERE carru.compliance_assessment_requirement_id = compliance_assessment_requirements.id
				  AND crru.is_risk_based_assessment = ?
			)`,
			*params.IsRiskBasedAssessment,
		)
	}
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
			query = query.Where("compliance_assessment_requirements.is_all_completed = ?", true)
		case "not_completed":
			query = query.Where("compliance_assessment_requirements.is_all_completed = ?", false)
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

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Use TRIM on the database column in the ORDER BY clause
	query = query.Order(fmt.Sprintf("TRIM(%s) %s", sortBy, sort))

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Execute the query
	err := query.Find(&complianceRequirements).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((total + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssessmentRequirementListResponseDTO{
		Items: complianceRequirements,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  total,
		},
	}
	return response, nil
}

func (repo *ComplianceAssessmentRequirementRepository) GetResponsibleUnitRowsByAssessmentRequirementIDs(assessmentRequirementIDs []uint, params *complianceDTOs.ComplianceAssessmentRequirementListQueryParamsDTO, out *[]complianceDTOs.ComplianceAssessmentRequirementResponsibleUnitListRowDTO) error {
	if len(assessmentRequirementIDs) == 0 {
		return nil
	}

	query := repo.db.Table("compliance_assessment_requirement_responsible_units carru").
		Select(`
			carru.id,
			carru.compliance_assessment_requirement_id as assessment_requirement_id,
			carru.responsible_unit_id,
			crru.is_regulation_notification_assessment,
			crru.is_risk_based_assessment,
			ou.title_en,
			ou.title_ar,
			ou.status,
			ou.color_code,
			ou.organization_unit_id as parent_id
		`).
		Joins("JOIN compliance_assessment_requirements car ON car.id = carru.compliance_assessment_requirement_id").
		Joins("JOIN compliance_requirement_responsible_units crru ON crru.compliance_requirement_id = car.compliance_requirement_id AND crru.responsible_unit_id = carru.responsible_unit_id").
		Joins("JOIN organization_units ou ON ou.id = carru.responsible_unit_id").
		Where("carru.compliance_assessment_requirement_id IN ?", assessmentRequirementIDs)

	if params != nil && params.IsRegulationNotificationAssessment != nil {
		query = query.Where("crru.is_regulation_notification_assessment = ?", *params.IsRegulationNotificationAssessment)
	}
	if params != nil && params.IsRiskBasedAssessment != nil {
		query = query.Where("crru.is_risk_based_assessment = ?", *params.IsRiskBasedAssessment)
	}

	return query.Scan(out).Error
}

func (repo *ComplianceAssessmentRequirementRepository) GetByID(id uint) (*models.ComplianceAssessmentRequirement, error) {
	var complianceAssessmentRequirement *models.ComplianceAssessmentRequirement

	result := repo.db.First(&complianceAssessmentRequirement, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessmentRequirement, nil
}

func (repo *ComplianceAssessmentRequirementRepository) GetFullDetailsByID(requirementId uint, responsibleUnitID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirement, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var car models.ComplianceAssessmentRequirement

	query := db.
		Preload("ComplianceAssessment").
		Preload("ComplianceAssessment.SubmittedByUser").
		Preload("ComplianceRequirement").
		Preload("ComplianceRequirement.ComplianceStatus").
		Preload("ComplianceRequirement.ComplianceClause").
		Preload("ComplianceRequirement.SeverityLevel").
		Preload("ComplianceRequirement.InherentRiskAnalysis").
		Preload("ComplianceRequirement.InherentRiskAnalysis.ComplianceAssessmentImpactType").
		Preload("ComplianceRequirement.InherentRiskRating").
		Preload("ComplianceStatus").
		Preload("ComplianceAssessmentRequirementStatus").
		Preload("AverageResidualRiskRating").
		Preload("CreatedByUser").
		Preload("UpdatedByUser")

	if responsibleUnitID != 0 {
		query = query.Preload("ResponsibleUnits", "responsible_unit_id = ?", responsibleUnitID)
	} else {
		query = query.Preload("ResponsibleUnits")
	}

	query = query.
		Preload("ResponsibleUnits.CompletedByUser").
		Preload("ResponsibleUnits.ResponsibleUnit").
		Preload("ResponsibleUnits.ComplianceAssessmentFrameworkValue").
		Preload("ResponsibleUnits.ComplianceAssessmentFrameworkValue.ComplianceAssessmentFramework").
		Preload("ResponsibleUnits.ComplianceAssessmentFrameworkValue.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ResponsibleUnits.ComplianceAssessmentFrameworkValue.ComplianceStatus").
		Preload("ResponsibleUnits.ComplianceAssessmentRequirementStatus").
		Preload("ResponsibleUnits.RequiredDocumentTitles").
		Preload("ResponsibleUnits.RequiredDocumentTitles.RequiredDocumentTitle").
		Preload("ResponsibleUnits.RequiredDocumentTitles.Documents").
		Preload("ResponsibleUnits.ActionPlans").
		Preload("ResponsibleUnits.ActionPlans.ResponsibleUser").
		Preload("ResponsibleUnits.ActionPlans.ReviewUser").
		Preload("ResponsibleUnits.ActionPlans.ActionPlanStatus").
		Preload("ResponsibleUnits.RequirementComments").
		Preload("ResponsibleUnits.RequirementComments.ComplianceAssessmentRequirementStatus").
		Preload("ResponsibleUnits.RequirementComments.CreatedByUser").
		Preload("ResponsibleUnits.ResidualRiskAnalysis").
		Preload("ResponsibleUnits.ResidualRiskAnalysis.ComplianceAssessmentImpactType").
		Preload("ResponsibleUnits.ResidualControls").
		Preload("ResponsibleUnits.ResidualControls.ComplianceAssessmentImpactType").
		Preload("ResponsibleUnits.ResidualControls.Control").
		Preload("ResponsibleUnits.ResidualRiskRating")

	if err := query.First(&car, requirementId).Error; err != nil {
		return nil, err
	}
	return &car, nil
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
func (repo *ComplianceAssessmentRequirementRepository) GetComplianceAssessmentRequirementResponsibleUnitByID(id uint) (*models.ComplianceAssessmentRequirementResponsibleUnit, error) {
	var complianceAssessmentRequirementResponsibleUnit *models.ComplianceAssessmentRequirementResponsibleUnit

	result := repo.db.Preload("ResponsibleUnit").Preload("ComplianceAssessmentRequirement").First(&complianceAssessmentRequirementResponsibleUnit, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessmentRequirementResponsibleUnit, nil
}

func (repo *ComplianceAssessmentRequirementRepository) UpdateComplianceAssessmentRequirementResponsibleUnit(complianceAssessmentRequirementResponsibleUnit *models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(complianceAssessmentRequirementResponsibleUnit).Error
}

func (repo *ComplianceAssessmentRequirementRepository) ComplianceAssessmentUpdateStatus(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessment).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRequirementRepository) UpdateComplianceStatus(complianceRegisterID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceRegister{}).
		Where("id = ?", complianceRegisterID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRequirementRepository) UpdateComplianceAssessmentRequirementStatus(assessmentReqID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceAssessmentRequirement{}).
		Where("id = ?", assessmentReqID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRequirementRepository) UpdateClauseComplianceStatus(clauseID uint, statusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	res := db.Model(&models.ComplianceClause{}).
		Where("id = ?", clauseID).
		Updates(map[string]interface{}{
			"compliance_status_id": statusID,
			"score":                score,
		})

	if res.Error != nil {
		return res.Error
	}

	return nil
}
