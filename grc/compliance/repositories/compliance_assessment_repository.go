package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceAssessment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error)
	GetMinimalDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error)
	Create(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error

	IsAssessmentExists(organizationUnitID uint, compliancRegisterID uint, ids []uint) (bool, error)
	IsAssessmentExistExcludingID(id uint, organizationUnitID uint, compliancRegisterID uint, ids []uint) (bool, error)
	GetRequirementsForRegister(registerID uint) ([]*models.ComplianceRequirement, error)
	HasRequirementForRegister(registerID uint) (bool, error)
	CheckAllRequirementsHaveInherentScore(registerID uint) (bool, error)
	CheckAllRiskBasedRequirementsHaveInherentScore(registerID uint) (bool, error)

	CreateComplianceAssessmentRequirement(requirements *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteComplianceAssessmentRequirements(ids []*uint, tx ...utility.DatabaseTransaction) error
	CreateComplianceAssessmentRequirementDocumentTitle(requirementDocuments []*models.ComplianceAssessmentRequirementDocumentTitle, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentRequirementDocumentTitles(ids []*uint, tx ...utility.DatabaseTransaction) error

	CreateComplianceAssessmentRequirementResponsibleUnit(resUnits *models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error
	CreateComplianceAssessmentRequirementResponsibleUnits(resUnits []*models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentRequirementResponsibleUnits(ids []*uint, tx ...utility.DatabaseTransaction) error

	GetComplianceScoreSummaryByRequirement(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error)

	AssessmentRequirementCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error)
	GetComplianceStatusSummaryByRequirement(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error)

	UpdateComplianceRequirementResponsibleUintValue(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error
	UpdateComplianceRequirementStatusUpdate(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error
	UpdateComplianceStatus(complianceRegisterID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateComplianceAssessmentRequirementStatus(assessmentReqID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateComplianceAssessmentControlStatus(complianceControlID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateClauseComplianceStatus(clauseID uint, statusID uint, score float64, tx ...utility.DatabaseTransaction) error
	UpdateAverageResidualScoreAndRating(assessmentID uint, averageScore float64, ratingID *uint, tx ...utility.DatabaseTransaction) error
	UpdateRequirementAverageResidualScoreAndRating(requirementID uint, averageScore float64, ratingID *uint, tx ...utility.DatabaseTransaction) error
	UpdateRegisterAverageResidualScoreAndRating(registerID uint, averageScore float64, ratingID *uint, tx ...utility.DatabaseTransaction) error
	GetComplianceAssessmentRatings() ([]*models.ComplianceAssessmentRating, error)
	CreateComplianceAssessmentControl(control *models.ComplianceAssessmentControl, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentControls(ids []*uint, tx ...utility.DatabaseTransaction) error
	HasControlForRegister(registerID uint) (bool, error)

	GetRegisterIDsWithResponsibleUnitID(responsibleUnitID uint) ([]uint, error)
	CreateComplianceAssessmentCompliances(compliances []*models.ComplianceAssessmentCompliance, tx ...utility.DatabaseTransaction) error

	UpdateControlComplianceStatus(dto complianceDTOs.ControlSyncUpdateDTO, tx ...utility.DatabaseTransaction) error
	UpdateComplianceControlMeta(dto complianceDTOs.ControlSyncUpdateDTO, tx ...utility.DatabaseTransaction) error

	AssessmentControlCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error)
	GetComplianceStatusSummaryByControl(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error)
	GetComplianceScoreSummaryByControl(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error)
	GetFullDetailsByIDOptimized(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error)
	ComputeFlagsFast(assessmentID uint, tx ...utility.DatabaseTransaction) (isCompleted, isApproved bool, err error)

	// Action plan deletion methods
	DeleteActionPlansByAssessment(id uint, tx ...utility.DatabaseTransaction) error
	DeleteActionPlanCommentDocumentsByAssessment(id uint, tx ...utility.DatabaseTransaction) error
	DeleteActionPlanCommentsByAssessment(id uint, tx ...utility.DatabaseTransaction) error
	DeleteActionPlanStatusUpdateDocumentsByAssessment(id uint, tx ...utility.DatabaseTransaction) error
	DeleteActionPlanStatusUpdatesByAssessment(id uint, tx ...utility.DatabaseTransaction) error

	// Compliance assessment compliance methods
	GetComplianceAssessmentCompliaceByID(id uint) (*models.ComplianceAssessmentCompliance, error)
	UpdateComplianceAssessmentCompliance(compliance *models.ComplianceAssessmentCompliance, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssessmentCompliance(id uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRepository(db *gorm.DB) *ComplianceAssessmentRepository {
	return &ComplianceAssessmentRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentRepository) GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error) {

	var complianceAssessments []*complianceDTOs.ComplianceAssessmentListDTO
	var totalCount int64

	query := repo.db.Table("compliance_assessments").
		Select(`
			compliance_assessments.id,
			compliance_assessments.description,
			compliance_assessments.[type] AS type,
			compliance_assessments.created_at,
			compliance_assessments.updated_at,
			reg_agg.compliance_register_title_en AS compliance_register_title_en,
			reg_agg.compliance_register_title_ar AS compliance_register_title_ar,
			organization_units.id AS responsible_department_id,
			organization_units.title_en AS responsible_department_title_en,
			organization_units.title_ar AS responsible_department_title_ar,
			review_user.id AS review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,
			compliance_assessment_statuses.title_en AS compliance_assessment_status_title_en,
			compliance_assessment_statuses.title_ar AS compliance_assessment_status_title_ar,
			compliance_assessment_statuses.color_code AS compliance_assessment_status_color_code,
			compliance_assessment_statuses.type AS compliance_assessment_status_type,
			average_residual_rating.id AS average_residual_risk_rating_id,
			average_residual_rating.title_en AS average_residual_risk_rating_title_en,
			average_residual_rating.title_ar AS average_residual_risk_rating_title_ar,
			average_residual_rating.color_code AS average_residual_risk_rating_color_code,
			average_residual_rating.type AS average_residual_risk_rating_type,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN (
				SELECT
					cac.compliance_assessment_id,
					STRING_AGG(CAST(cr.title_en AS NVARCHAR(MAX)), ', ') AS compliance_register_title_en,
					STRING_AGG(CAST(cr.title_ar AS NVARCHAR(MAX)), ', ') AS compliance_register_title_ar
				FROM compliance_assessment_compliances cac
				JOIN compliance_register cr ON cr.id = cac.compliance_register_id
				GROUP BY cac.compliance_assessment_id
			) reg_agg ON reg_agg.compliance_assessment_id = compliance_assessments.id
			LEFT JOIN users AS review_user ON review_user.id = compliance_assessments.review_user_id
			LEFT JOIN compliance_assessment_statuses ON compliance_assessment_statuses.id = compliance_assessments.compliance_assessment_status_id
			LEFT JOIN compliance_assessment_ratings AS average_residual_rating ON average_residual_rating.id = compliance_assessments.average_residual_risk_rating_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_assessments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_assessments.updated_by_id
			LEFT JOIN organization_units ON organization_units.id = compliance_assessments.responsible_unit_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "draft":
			query = query.Where("compliance_assessment_statuses.type = ?", "draft")
		case "in_review":
			query = query.Where("compliance_assessment_statuses.type = ?", "in_review")
		case "approved":
			query = query.Where("compliance_assessment_statuses.type = ?", "approved")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(reg_agg.compliance_register_title_en LIKE ? OR
			reg_agg.compliance_register_title_ar LIKE ? OR
			review_user.first_name_en LIKE ? OR
			review_user.first_name_ar LIKE ? OR
			review_user.last_name_en LIKE ? OR
			review_user.last_name_ar LIKE ? OR
			review_user.email LIKE ? OR
			review_user.initial LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "compliance_assessments.created_at"
	}
	query = query.Order(sortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&complianceAssessments).Error
	if err != nil {
		return nil, err
	}

	if params.Export {
		return &complianceDTOs.ComplianceAssessmentListResponseDTO{
			Items: complianceAssessments,
		}, nil
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &complianceDTOs.ComplianceAssessmentListResponseDTO{
		Items: complianceAssessments,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceAssessmentRepository) GetByID(id uint) (*models.ComplianceAssessment, error) {
	var complianceAssessment *models.ComplianceAssessment

	result := repo.db.First(&complianceAssessment, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessment, nil

}

func (repo *ComplianceAssessmentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var m models.ComplianceAssessment

	if err := db.
		Preload("ComplianceAssessmentCompliances").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.Values").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.Values.ComplianceStatus").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ResponsibleUser").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.AverageResidualRiskRating").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities.IssuingAuthority").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities.Entity").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities.Document").
		Preload("AverageResidualRiskRating").
		Preload("AssessmentRequirements").
		Preload("AssessmentRequirements.ComplianceStatus").
		Preload("AssessmentRequirements.AverageResidualRiskRating").
		Preload("AssessmentRequirements.ResponsibleUnits").
		Preload("AssessmentRequirements.ResponsibleUnits.ResponsibleUnit").
		Preload("AssessmentRequirements.ResponsibleUnits.RequiredDocumentTitles").
		Preload("AssessmentRequirements.ResponsibleUnits.ComplianceAssessmentRequirementStatus").
		Preload("AssessmentRequirements.ResponsibleUnits.ComplianceAssessmentFrameworkValue").
		Preload("AssessmentRequirements.ResponsibleUnits.ResidualRiskAnalysis").
		Preload("AssessmentRequirements.ResponsibleUnits.ResidualRiskAnalysis.ComplianceAssessmentImpactType").
		Preload("AssessmentRequirements.ResponsibleUnits.ResidualRiskRating").
		Preload("AssessmentRequirements.ResponsibleUnits.ResidualControls").
		Preload("AssessmentRequirements.ResponsibleUnits.ResidualControls.ComplianceAssessmentImpactType").
		Preload("AssessmentRequirements.ResponsibleUnits.ResidualControls.Control").
		Preload("AssessmentRequirements.ComplianceRequirement").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceStatus").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceClause").
		Preload("AssessmentControls").
		Preload("AssessmentControls.ComplianceControl").
		Preload("AssessmentControls.ComplianceControl.Control").
		Preload("AssessmentControls.ComplianceControl.ComplianceClause"). // if nested
		Preload("AssessmentControls.ComplianceAssessmentFrameworkValue").
		Preload("AssessmentControls.ComplianceAssessmentRequirementStatus").
		Preload("AssessmentControls.CompletedByUser").
		Preload("ReviewUser").
		Preload("ComplianceAssessmentStatus").
		Preload("SubmittedByUser").
		Preload("ReviewedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ResponsibleOrganizationUnit").
		First(&m, id).Error; err != nil {
		return nil, err
	}

	return &m, nil
}

func (repo *ComplianceAssessmentRepository) GetMinimalDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceAssessment *models.ComplianceAssessment

	result := db.
		Preload("ComplianceAssessmentCompliances").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.Values").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.Values.ComplianceStatus").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ResponsibleUser").
		Preload("AssessmentRequirements").
		Preload("ReviewUser").
		Preload("ComplianceAssessmentStatus").
		Preload("SubmittedByUser").
		Preload("ReviewedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&complianceAssessment, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return complianceAssessment, nil
}
func (repo *ComplianceAssessmentRepository) Create(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceAssessment).Error; err != nil {
		return 0, err
	}

	return complianceAssessment.ID, nil
}

func (repo *ComplianceAssessmentRepository) Update(complianceAssessment *models.ComplianceAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceAssessment).Error; err != nil {
		fmt.Println("hoiiiiii", err)
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Compliance  Assessment
	if err := db.Delete(&models.ComplianceAssessment{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) IsAssessmentExists(organizationUnitID uint, complianceRegisterID uint, ids []uint) (bool, error) {
	db := repo.db

	var count int64
	query := db.Table("compliance_assessments").
		Where("compliance_register_id = ?", complianceRegisterID).
		Where("compliance_assessment_status_id IN ?", ids)

	if organizationUnitID > 0 {
		query = query.Where("organization_unit_id = ?", organizationUnitID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}
func (repo *ComplianceAssessmentRepository) HasRequirementForRegister(registerID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ?", registerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ComplianceAssessmentRepository) GetRequirementsForRegister(registerID uint) ([]*models.ComplianceRequirement, error) {
	var requirements []*models.ComplianceRequirement
	err := repo.db.
		Model(&models.ComplianceRequirement{}).
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Preload("ResponsibleUnits").
		Where("compliance_clauses.compliance_register_id = ?", registerID).
		Find(&requirements).Error
	if err != nil {
		return nil, err
	}

	return requirements, nil
}

func (repo *ComplianceAssessmentRepository) CheckAllRequirementsHaveInherentScore(registerID uint) (bool, error) {
	var totalCount int64
	var withScoreCount int64

	// Count total requirements for the register
	err := repo.db.
		Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ?", registerID).
		Count(&totalCount).Error

	if err != nil {
		return false, err
	}

	// Count requirements with inherent risk score
	err = repo.db.
		Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ? AND compliance_requirements.inherent_risk_score IS NOT NULL", registerID).
		Count(&withScoreCount).Error

	if err != nil {
		return false, err
	}

	// Return true only if all requirements have inherent risk scores and there's at least one requirement
	return totalCount > 0 && totalCount == withScoreCount, nil
}

func (repo *ComplianceAssessmentRepository) CheckAllRiskBasedRequirementsHaveInherentScore(registerID uint) (bool, error) {
	var totalCount int64
	var withScoreCount int64

	err := repo.db.
		Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Joins("JOIN compliance_requirement_responsible_units crru ON crru.compliance_requirement_id = compliance_requirements.id").
		Where("compliance_clauses.compliance_register_id = ? AND crru.is_risk_based_assessment = ?", registerID, true).
		Distinct("compliance_requirements.id").
		Count(&totalCount).Error
	if err != nil {
		return false, err
	}

	if totalCount == 0 {
		return true, nil
	}

	err = repo.db.
		Table("compliance_requirements").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Joins("JOIN compliance_requirement_responsible_units crru ON crru.compliance_requirement_id = compliance_requirements.id").
		Where("compliance_clauses.compliance_register_id = ? AND crru.is_risk_based_assessment = ? AND compliance_requirements.inherent_risk_score IS NOT NULL", registerID, true).
		Distinct("compliance_requirements.id").
		Count(&withScoreCount).Error
	if err != nil {
		return false, err
	}

	return totalCount == withScoreCount, nil
}
func (repo *ComplianceAssessmentRepository) IsAssessmentExistExcludingID(id uint, organizationUnitID uint, complianceRegisterID uint, ids []uint) (bool, error) {
	db := repo.db

	var count int64
	query := db.Table("compliance_assessments").
		Where("id <> ? AND compliance_register_id = ?", id, complianceRegisterID).
		Where("compliance_assessment_status_id IN ?", ids)

	if organizationUnitID > 0 {
		query = query.Where("organization_unit_id = ?", organizationUnitID)
	}

	err := query.Count(&count).Error

	return count > 0, err
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirement(requirements *models.ComplianceAssessmentRequirement, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&requirements).Error; err != nil {
		return 0, err
	}
	return requirements.ID, nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentRequirements(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentRequirement{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirementDocumentTitle(requirementDocuments []*models.ComplianceAssessmentRequirementDocumentTitle, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&requirementDocuments).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentRequirementDocumentTitles(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentRequirementDocumentTitle{}).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) AssessmentRequirementCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error) {
	var out complianceDTOs.AssessmentCompletionRollupDTO

	const q = `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN carru.is_completed = 1 THEN 1 ELSE 0 END) AS completed,
			SUM(CASE WHEN carru.is_completed = 0 THEN 1 ELSE 0 END) AS not_completed,
			CASE 
				WHEN COUNT(*) = 0 THEN 0
				ELSE CAST(SUM(CASE WHEN carru.is_completed = 1 THEN 1 ELSE 0 END) * 100.0 
				          / COUNT(*) AS DECIMAL(5,2))
			END AS completion_percentage
		FROM compliance_assessment_requirement_responsible_units carru
		JOIN compliance_assessment_requirements car 
			ON car.id = carru.compliance_assessment_requirement_id
		WHERE car.compliance_assessment_id = ?
	`

	if err := repo.db.Raw(q, assessmentID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (repo *ComplianceAssessmentRepository) GetComplianceStatusSummaryByRequirement(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error) {
	var out []complianceDTOs.ComplianceStatusCountDTO

	const q = `
		SELECT
			cs.id   AS compliance_status_id,
			cs.title_en,
			cs.title_ar,
			cs.type,
			cs.color_code,
			COALESCE(agg.cnt, 0) AS count
		FROM compliance_statuses cs
		LEFT JOIN (
			SELECT 
				cafv.compliance_status_id AS compliance_status_id,
				COUNT(*) AS cnt
			FROM compliance_assessment_requirement_responsible_units carru
			JOIN compliance_assessment_requirements car
				ON car.id = carru.compliance_assessment_requirement_id
			JOIN compliance_assessments ca
				ON car.compliance_assessment_id = ca.id
			JOIN compliance_assessment_statuses cas
				ON ca.compliance_assessment_status_id = cas.id
			LEFT JOIN compliance_assessment_framework_values cafv
				ON cafv.id = carru.compliance_assessment_framework_value_id
			WHERE car.compliance_assessment_id = ?
			GROUP BY cafv.compliance_status_id
		) agg
			ON agg.compliance_status_id = cs.id
		-- optional: only active statuses
		-- WHERE cs.status = 'active'
		ORDER BY cs.id
	`

	if err := repo.db.Raw(q, assessmentID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (repo *ComplianceAssessmentRepository) GetComplianceScoreSummaryByRequirement(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error) {
	var result complianceDTOs.ComplianceScoreSummaryDTO

	const query = `
		SELECT 
			COUNT(car.id) AS total_requirements, 
			COALESCE(AVG(CAST(car.score AS FLOAT)), 0) AS average_score
		FROM compliance_assessment_requirements car
		JOIN compliance_assessments ca 
			ON car.compliance_assessment_id = ca.id
		JOIN compliance_assessment_statuses cas 
			ON ca.compliance_assessment_status_id = cas.id
		JOIN compliance_statuses cs
			ON cs.id = car.compliance_status_id
		WHERE car.compliance_assessment_id = ? 
		  AND car.score IS NOT NULL
		  AND cas.type = 'approved'
		  AND cs.type <> 'not_applicable'
	`

	if err := repo.db.Raw(query, assessmentID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirementResponsibleUnits(resUnits []*models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&resUnits).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentRequirementResponsibleUnits(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentRequirementResponsibleUnit{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentRequirementResponsibleUnit(resUnits *models.ComplianceAssessmentRequirementResponsibleUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(resUnits).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) UpdateComplianceRequirementResponsibleUintValue(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ComplianceRequirementResponsibleUnit{}).Where("responsible_unit_id = ? AND compliance_requirement_id = ?", complianceRequirementValueUpdateDTO.ResponsibleUnitID, complianceRequirementValueUpdateDTO.ComplianceRequirementID).
		Updates(map[string]interface{}{
			"compliance_assessment_framework_value_id": complianceRequirementValueUpdateDTO.ComplianceAssessmentFrameworkValueID,
			"comments":             complianceRequirementValueUpdateDTO.Comments,
			"completed_by_user_id": complianceRequirementValueUpdateDTO.CompletedByUserID,
			"completed_at":         complianceRequirementValueUpdateDTO.CompletedAt,
		}).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceRequirementStatusUpdate(complianceRequirementValueUpdateDTO complianceDTOs.ComplianceRequirementValueUpdateDTO, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Model(&models.ComplianceRequirement{}).Where("id = ?", complianceRequirementValueUpdateDTO.ComplianceRequirementID).
		Update("compliance_status_id", complianceRequirementValueUpdateDTO.ComplianceStatusID).Update("score", complianceRequirementValueUpdateDTO.Score).Error; err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceStatus(complianceRegisterID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceRegister{}).
		Where("id = ?", complianceRegisterID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceAssessmentRequirementStatus(assessmentReqID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceAssessmentRequirement{}).
		Where("id = ?", assessmentReqID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRepository) UpdateComplianceAssessmentControlStatus(complianceControlID uint, complianceStatusID uint, score float64, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ComplianceAssessmentControl{}).
		Where("id = ?", complianceControlID).
		Update("compliance_status_id", complianceStatusID).Update("score", score).Error
}
func (repo *ComplianceAssessmentRepository) UpdateClauseComplianceStatus(clauseID uint, statusID uint, score float64, tx ...utility.DatabaseTransaction) error {
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

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentControl(control *models.ComplianceAssessmentControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(control).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentControls(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceAssessmentControl{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRepository) HasControlForRegister(registerID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table("compliance_controls").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_controls.compliance_clause_id").
		Where("compliance_clauses.compliance_register_id = ?", registerID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ComplianceAssessmentRepository) UpdateControlComplianceStatus(dto complianceDTOs.ControlSyncUpdateDTO, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Use map so NULLs (nil) are written
	return db.Model(&models.Control{}).
		Where("id = ?", dto.ControlID).
		Updates(map[string]interface{}{
			"compliance_status_id": dto.ComplianceStatusID,
		}).Error
}

func (repo *ComplianceAssessmentRepository) UpdateComplianceControlMeta(dto complianceDTOs.ControlSyncUpdateDTO, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Use map so NULLs (nil) are written
	return db.Model(&models.ComplianceControl{}).
		Where("id = ?", dto.ComplianceControlID).
		Updates(map[string]interface{}{
			"compliance_assessment_framework_value_id": dto.ComplianceAssessmentFrameworkValueID,
			"completed_by_user_id":                     dto.CompletedByUserID,
			"completed_at":                             dto.CompletedAt,
		}).Error
}

func (repo *ComplianceAssessmentRepository) AssessmentControlCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error) {
	var out complianceDTOs.AssessmentCompletionRollupDTO

	const q = `
		SELECT
			COUNT(*) AS total,
			SUM(CASE WHEN cac.is_completed = 1 THEN 1 ELSE 0 END) AS completed,
			SUM(CASE WHEN cac.is_completed = 0 THEN 1 ELSE 0 END) AS not_completed,
			CASE 
				WHEN COUNT(*) = 0 THEN 0
				ELSE CAST(SUM(CASE WHEN cac.is_completed = 1 THEN 1 ELSE 0 END) * 100.0 
				          / COUNT(*) AS DECIMAL(5,2))
			END AS completion_percentage
		FROM compliance_assessment_controls cac
		WHERE cac.compliance_assessment_id = ?
	`

	if err := repo.db.Raw(q, assessmentID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return &out, nil
}
func (repo *ComplianceAssessmentRepository) GetComplianceStatusSummaryByControl(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error) {
	var out []complianceDTOs.ComplianceStatusCountDTO

	const q = `
		SELECT
			cs.id   AS compliance_status_id,
			cs.title_en,
			cs.title_ar,
			cs.type,
			cs.color_code,
			COALESCE(agg.cnt, 0) AS count
		FROM compliance_statuses cs
		LEFT JOIN (
			SELECT 
				cafv.compliance_status_id AS compliance_status_id,
				COUNT(*) AS cnt
			FROM compliance_assessment_controls cac
			JOIN compliance_assessments ca
				ON ca.id = cac.compliance_assessment_id
			JOIN compliance_assessment_statuses cas
				ON cas.id = ca.compliance_assessment_status_id
			LEFT JOIN compliance_assessment_framework_values cafv
				ON cafv.id = cac.compliance_assessment_framework_value_id
			WHERE cac.compliance_assessment_id = ?
			  AND cas.type = 'approved'
				-- If you want to count only applicable controls, uncomment:
				-- AND cafv.is_applicable = 1
				-- If you want to ignore unset framework values, uncomment:
				-- AND cac.compliance_assessment_framework_value_id IS NOT NULL
			GROUP BY cafv.compliance_status_id
		) agg
			ON agg.compliance_status_id = cs.id
		ORDER BY cs.id
	`

	if err := repo.db.Raw(q, assessmentID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (repo *ComplianceAssessmentRepository) GetComplianceScoreSummaryByControl(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error) {
	db := repo.db

	var result complianceDTOs.ComplianceScoreSummaryDTO

	// average_score: average of clause-level scores (each clause weighted equally)
	// total_requirements: count of applicable controls for this assessment
	const q = `
        SELECT
            /* total applicable controls */
            (
                SELECT COUNT(cac.id)
                FROM compliance_assessment_controls cac
                JOIN compliance_assessments ca
                  ON ca.id = cac.compliance_assessment_id
                JOIN compliance_assessment_statuses cas
                  ON cas.id = ca.compliance_assessment_status_id
                LEFT JOIN compliance_assessment_framework_values cafv
                  ON cafv.id = cac.compliance_assessment_framework_value_id
                LEFT JOIN compliance_statuses cs
                  ON cs.id = cafv.compliance_status_id
                WHERE cac.compliance_assessment_id = ?
                  AND cas.type = 'approved'
                  AND cafv.is_applicable = 1             -- count only applicable controls
                  AND (cs.type IS NULL OR cs.type <> 'not_applicable')
            ) AS total_requirements,

            /* average across DISTINCT clauses using their persisted clause score */
            COALESCE(AVG(clause_score), 0) AS average_score
        FROM (
            SELECT
                c.id AS clause_id,
                COALESCE(c.score, 0) AS clause_score
            FROM compliance_clauses c
            JOIN compliance_controls cc
              ON cc.compliance_clause_id = c.id
            JOIN compliance_assessment_controls cac
              ON cac.compliance_control_id = cc.id
             AND cac.compliance_assessment_id = ?
            JOIN compliance_assessments ca
              ON ca.id = cac.compliance_assessment_id
            JOIN compliance_assessment_statuses cas
              ON cas.id = ca.compliance_assessment_status_id
            LEFT JOIN compliance_statuses ccs
              ON ccs.id = c.compliance_status_id
            WHERE cas.type = 'approved'
              AND (ccs.type IS NULL OR ccs.type <> 'not_applicable') -- drop NA clauses
            GROUP BY c.id, c.score
        ) clause_aggs;
    `

	if err := db.Raw(q, assessmentID, assessmentID).Scan(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

const chunkSize = 800 // < 2100 SQL Server parameter limit safe margin

func (repo *ComplianceAssessmentRepository) GetFullDetailsByIDOptimized(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	db = db.Session(&gorm.Session{
		QueryFields: true,
		PrepareStmt: true,
	})
	// 1) Load light header only
	var a *models.ComplianceAssessment
	if err := db.
		Preload("ComplianceAssessmentCompliances").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.Values").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ComplianceAssessmentFramework.Values.ComplianceStatus").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.ResponsibleUser").
		Preload("AverageResidualRiskRating").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities.IssuingAuthority").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities.Entity").
		Preload("ComplianceAssessmentCompliances.ComplianceRegister.IssuingAuthorities.Document").
		Preload("ReviewUser").
		Preload("ComplianceAssessmentStatus").
		Preload("SubmittedByUser").
		Preload("ReviewedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ResponsibleOrganizationUnit").
		First(&a, id).
		Error; err != nil {
		return nil, err
	}

	// 2) Collect requirement IDs first
	var reqIDs []uint
	if err := db.
		Model(&models.ComplianceAssessmentRequirement{}).
		Where("compliance_assessment_id = ?", id).
		Order("id ASC").
		Pluck("id", &reqIDs).Error; err != nil {
		return nil, err
	}
	if len(reqIDs) == 0 {
		a.AssessmentRequirements = []*models.ComplianceAssessmentRequirement{}
	}

	// 3) Hydrate requirements (lean) in chunks (no heavy preloads)
	var reqs []*models.ComplianceAssessmentRequirement
	{
		var tmp []models.ComplianceAssessmentRequirement
		for i := 0; i < len(reqIDs); i += chunkSize {
			end := min(i+chunkSize, len(reqIDs))
			var chunk []models.ComplianceAssessmentRequirement
			if err := db.
				Where("id IN ?", reqIDs[i:end]).
				Find(&chunk).Error; err != nil {
				return nil, err
			}
			tmp = append(tmp, chunk...)
		}
		reqs = make([]*models.ComplianceAssessmentRequirement, 0, len(tmp))
		for i := range tmp {
			reqs = append(reqs, &tmp[i])
		}
		a.AssessmentRequirements = reqs
	}

	// 4) Collect foreign keys to batch-load their relations
	var complianceReqIDs []uint
	for _, r := range reqs {
		complianceReqIDs = append(complianceReqIDs, r.ComplianceRequirementID)
	}

	// 4a) Batch load ComplianceRequirement (+Clause, +Status) for all reqs
	var allCR []models.ComplianceRequirement
	for i := 0; i < len(complianceReqIDs); i += chunkSize {
		end := min(i+chunkSize, len(complianceReqIDs))
		var chunk []models.ComplianceRequirement
		if err := db.
			Where("id IN ?", complianceReqIDs[i:end]).
			Find(&chunk).Error; err != nil {
			return nil, err
		}
		allCR = append(allCR, chunk...)
	}
	crByID := make(map[uint]*models.ComplianceRequirement, len(allCR))
	for i := range allCR {
		crByID[allCR[i].ID] = &allCR[i]
	}
	// Attach CR now (clause/status will be stitched next)
	for _, r := range reqs {
		r.ComplianceRequirement = crByID[r.ComplianceRequirementID]
	}

	// 4b) Load Clauses for those CRs
	var clauseIDs []uint
	for _, cr := range crByID {
		if cr != nil && cr.ComplianceClauseID != 0 {
			clauseIDs = append(clauseIDs, cr.ComplianceClauseID)
		}
	}
	var allClauses []models.ComplianceClause
	for i := 0; i < len(clauseIDs); i += chunkSize {
		end := min(i+chunkSize, len(clauseIDs))
		var chunk []models.ComplianceClause
		if err := db.Where("id IN ?", clauseIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		allClauses = append(allClauses, chunk...)
	}
	byClauseID := make(map[uint]*models.ComplianceClause, len(allClauses))
	for i := range allClauses {
		byClauseID[allClauses[i].ID] = &allClauses[i]
	}
	for _, cr := range crByID {
		if cr != nil && cr.ComplianceClauseID != 0 {
			cr.ComplianceClause = byClauseID[cr.ComplianceClauseID]
		}
	}

	// 4c) Load CR.ComplianceStatus (if used in UI)
	var crStatusIDs []uint
	for _, cr := range crByID {
		if cr != nil && cr.ComplianceStatusID != nil {
			crStatusIDs = append(crStatusIDs, *cr.ComplianceStatusID)
		}
	}
	if len(crStatusIDs) > 0 {
		var allSt []models.ComplianceStatus
		for i := 0; i < len(crStatusIDs); i += chunkSize {
			end := min(i+chunkSize, len(crStatusIDs))
			var chunk []models.ComplianceStatus
			if err := db.Where("id IN ?", crStatusIDs[i:end]).Find(&chunk).Error; err != nil {
				return nil, err
			}
			allSt = append(allSt, chunk...)
		}
		byStID := make(map[uint]*models.ComplianceStatus, len(allSt))
		for i := range allSt {
			byStID[allSt[i].ID] = &allSt[i]
		}
		for _, cr := range crByID {
			if cr != nil && cr.ComplianceStatusID != nil {
				cr.ComplianceStatus = byStID[*cr.ComplianceStatusID]
			}
		}
	}

	// 5) Load RUs for all requirement IDs in chunks (lean)
	var allRUs []models.ComplianceAssessmentRequirementResponsibleUnit
	for i := 0; i < len(reqIDs); i += chunkSize {
		end := min(i+chunkSize, len(reqIDs))
		var chunk []models.ComplianceAssessmentRequirementResponsibleUnit
		if err := db.
			Where("compliance_assessment_requirement_id IN ?", reqIDs[i:end]).
			Find(&chunk).Error; err != nil {
			return nil, err
		}
		allRUs = append(allRUs, chunk...)
	}
	// Group RUs by requirement
	ruByReq := make(map[uint][]*models.ComplianceAssessmentRequirementResponsibleUnit, len(reqIDs))
	for _, r := range reqIDs {
		ruByReq[r] = make([]*models.ComplianceAssessmentRequirementResponsibleUnit, 0, 2)
	}
	for i := range allRUs {
		rid := allRUs[i].ComplianceAssessmentRequirementID
		ruByReq[rid] = append(ruByReq[rid], &allRUs[i])
	}
	// Attach to requirements
	for _, r := range reqs {
		r.ResponsibleUnits = ruByReq[r.ID]
	}

	// 6) Batch-load RU relations: ResponsibleUnit, RequiredDocumentTitles, FrameworkValue, Status, CompletedByUser (if any)
	var respUnitIDs []uint
	var ruIDs []uint
	var fvIDs []uint
	var stIDs []uint
	var userIDs []uint

	for i := range allRUs {
		ruIDs = append(ruIDs, allRUs[i].ID)
		respUnitIDs = append(respUnitIDs, allRUs[i].ResponsibleUnitID)
		if allRUs[i].ComplianceAssessmentFrameworkValueID != nil {
			fvIDs = append(fvIDs, *allRUs[i].ComplianceAssessmentFrameworkValueID)
		}
		if allRUs[i].ComplianceAssessmentRequirementStatusID != nil {
			stIDs = append(stIDs, *allRUs[i].ComplianceAssessmentRequirementStatusID)
		}
		if allRUs[i].CompletedByUserID != nil {
			userIDs = append(userIDs, *allRUs[i].CompletedByUserID)
		}
	}
	// De-dup helpers
	respUnitIDs = uniq(respUnitIDs)
	fvIDs = uniq(fvIDs)
	stIDs = uniq(stIDs)
	userIDs = uniq(userIDs)

	// 6a) RequiredDocumentTitles (by RU)
	var allDocTitles []models.ComplianceAssessmentRequirementDocumentTitle
	for i := 0; i < len(ruIDs); i += chunkSize {
		end := min(i+chunkSize, len(ruIDs))
		var chunk []models.ComplianceAssessmentRequirementDocumentTitle
		if err := db.
			Where("compliance_assessment_requirement_responsible_unit_id IN ?", ruIDs[i:end]).
			Find(&chunk).Error; err != nil {
			return nil, err
		}
		allDocTitles = append(allDocTitles, chunk...)
	}
	docByRU := make(map[uint][]*models.ComplianceAssessmentRequirementDocumentTitle, len(ruIDs))
	for _, id := range ruIDs {
		docByRU[id] = make([]*models.ComplianceAssessmentRequirementDocumentTitle, 0, 2)
	}
	for i := range allDocTitles {
		id := allDocTitles[i].ComplianceAssessmentRequirementResponsibleUnitID
		docByRU[id] = append(docByRU[id], &allDocTitles[i])
	}

	// 6b) ResponsibleUnit (org unit)
	var allOrgUnits []models.OrganizationUnit
	for i := 0; i < len(respUnitIDs); i += chunkSize {
		end := min(i+chunkSize, len(respUnitIDs))
		var chunk []models.OrganizationUnit
		if err := db.Where("id IN ?", respUnitIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		allOrgUnits = append(allOrgUnits, chunk...)
	}
	byOrgUnit := make(map[uint]*models.OrganizationUnit, len(allOrgUnits))
	for i := range allOrgUnits {
		byOrgUnit[allOrgUnits[i].ID] = &allOrgUnits[i]
	}

	// 6c) Framework values
	var allFVs []models.ComplianceAssessmentFrameworkValue
	for i := 0; i < len(fvIDs); i += chunkSize {
		end := min(i+chunkSize, len(fvIDs))
		var chunk []models.ComplianceAssessmentFrameworkValue
		if err := db.Where("id IN ?", fvIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		allFVs = append(allFVs, chunk...)
	}
	byFV := make(map[uint]*models.ComplianceAssessmentFrameworkValue, len(allFVs))
	for i := range allFVs {
		byFV[allFVs[i].ID] = &allFVs[i]
	}

	// 6d) RU Status
	var allRUStatus []models.ComplianceAssessmentRequirementStatus
	for i := 0; i < len(stIDs); i += chunkSize {
		end := min(i+chunkSize, len(stIDs))
		var chunk []models.ComplianceAssessmentRequirementStatus
		if err := db.Where("id IN ?", stIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		allRUStatus = append(allRUStatus, chunk...)
	}
	byRUStatus := make(map[uint]*models.ComplianceAssessmentRequirementStatus, len(allRUStatus))
	for i := range allRUStatus {
		byRUStatus[allRUStatus[i].ID] = &allRUStatus[i]
	}

	// 6e) Users (CompletedBy)
	var allUsers []models.User
	for i := 0; i < len(userIDs); i += chunkSize {
		end := min(i+chunkSize, len(userIDs))
		var chunk []models.User
		if err := db.Where("id IN ?", userIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		allUsers = append(allUsers, chunk...)
	}
	byUser := make(map[uint]*models.User, len(allUsers))
	for i := range allUsers {
		byUser[allUsers[i].ID] = &allUsers[i]
	}

	// 6f) Stitch RU relations
	for i := range allRUs {
		ru := &allRUs[i]
		ru.ResponsibleUnit = byOrgUnit[ru.ResponsibleUnitID]
		if ru.ComplianceAssessmentFrameworkValueID != nil {
			ru.ComplianceAssessmentFrameworkValue = byFV[*ru.ComplianceAssessmentFrameworkValueID]
		}
		if ru.ComplianceAssessmentRequirementStatusID != nil {
			ru.ComplianceAssessmentRequirementStatus = byRUStatus[*ru.ComplianceAssessmentRequirementStatusID]
		}
		if ru.CompletedByUserID != nil {
			ru.CompletedByUser = byUser[*ru.CompletedByUserID]
		}
		ru.RequiredDocumentTitles = docByRU[ru.ID]
	}

	// === Controls (optimized)
	controls, err := repo.loadControlsOptimized(db, id)
	if err != nil {
		return nil, err
	}
	if len(controls) == 0 {
		a.AssessmentControls = []*models.ComplianceAssessmentControl{}
	}
	a.AssessmentControls = controls

	return a, nil
}

// Loads AssessmentControls and all their relations in chunks and stitches them.
func (repo *ComplianceAssessmentRepository) loadControlsOptimized(db *gorm.DB, assessmentID uint) ([]*models.ComplianceAssessmentControl, error) {
	// 1) Collect control row IDs (stable order)
	var ctrlIDs []uint
	if err := db.
		Model(&models.ComplianceAssessmentControl{}).
		Where("compliance_assessment_id = ?", assessmentID).
		Order("id ASC").
		Pluck("id", &ctrlIDs).Error; err != nil {
		return nil, err
	}
	if len(ctrlIDs) == 0 {
		return []*models.ComplianceAssessmentControl{}, nil
	}

	// 2) Hydrate control rows (lean) in chunks
	var rows []models.ComplianceAssessmentControl
	for i := 0; i < len(ctrlIDs); i += chunkSize {
		end := min(i+chunkSize, len(ctrlIDs))
		var chunk []models.ComplianceAssessmentControl
		if err := db.Where("id IN ?", ctrlIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		rows = append(rows, chunk...)
	}
	// pointer slice
	controls := make([]*models.ComplianceAssessmentControl, 0, len(rows))
	for i := range rows {
		controls = append(controls, &rows[i])
	}

	// 3) Collect FKs for batch loads
	var (
		complianceControlIDs []uint
		fvIDs                []uint
		statusIDs            []uint
		userIDs              []uint
	)
	for _, c := range controls {
		complianceControlIDs = append(complianceControlIDs, c.ComplianceControlID)
		if c.ComplianceAssessmentFrameworkValueID != nil {
			fvIDs = append(fvIDs, *c.ComplianceAssessmentFrameworkValueID)
		}
		if c.ComplianceAssessmentRequirementStatusID != nil {
			statusIDs = append(statusIDs, *c.ComplianceAssessmentRequirementStatusID)
		}
		if c.CompletedByUserID != nil {
			userIDs = append(userIDs, *c.CompletedByUserID)
		}
	}
	complianceControlIDs = uniq(complianceControlIDs)
	fvIDs = uniq(fvIDs)
	statusIDs = uniq(statusIDs)
	userIDs = uniq(userIDs)

	// 4) Load ComplianceControl (+ Control, + Clause)
	var ccRows []models.ComplianceControl
	for i := 0; i < len(complianceControlIDs); i += chunkSize {
		end := min(i+chunkSize, len(complianceControlIDs))
		var chunk []models.ComplianceControl
		if err := db.Where("id IN ?", complianceControlIDs[i:end]).Find(&chunk).Error; err != nil {
			return nil, err
		}
		ccRows = append(ccRows, chunk...)
	}
	ccByID := make(map[uint]*models.ComplianceControl, len(ccRows))
	for i := range ccRows {
		ccByID[ccRows[i].ID] = &ccRows[i]
	}

	// 4a) Load nested Control for those ComplianceControls
	var ctrlMasterIDs []uint
	for _, cc := range ccByID {
		if cc.ControlID != 0 {
			ctrlMasterIDs = append(ctrlMasterIDs, cc.ControlID)
		}
	}
	ctrlMasterIDs = uniq(ctrlMasterIDs)

	if len(ctrlMasterIDs) > 0 {
		var masters []models.Control
		for i := 0; i < len(ctrlMasterIDs); i += chunkSize {
			end := min(i+chunkSize, len(ctrlMasterIDs))
			var chunk []models.Control
			if err := db.Where("id IN ?", ctrlMasterIDs[i:end]).Find(&chunk).Error; err != nil {
				return nil, err
			}
			masters = append(masters, chunk...)
		}
		masterByID := make(map[uint]*models.Control, len(masters))
		for i := range masters {
			masterByID[masters[i].ID] = &masters[i]
		}
		for _, cc := range ccByID {
			if cc.ControlID != 0 {
				cc.Control = masterByID[cc.ControlID]
			}
		}
	}

	// 4b) Load Clause for those ComplianceControls
	// 4b) Load Clause for those ComplianceControls
	var clauseIDs []uint
	for _, cc := range ccByID {
		if cc.ComplianceClauseID != nil && *cc.ComplianceClauseID != 0 {
			clauseIDs = append(clauseIDs, *cc.ComplianceClauseID)
		}
	}
	clauseIDs = uniq(clauseIDs)

	if len(clauseIDs) > 0 {
		var clauseRows []models.ComplianceClause
		for i := 0; i < len(clauseIDs); i += chunkSize {
			end := min(i+chunkSize, len(clauseIDs))
			var chunk []models.ComplianceClause
			if err := db.Where("id IN ?", clauseIDs[i:end]).Find(&chunk).Error; err != nil {
				return nil, err
			}
			clauseRows = append(clauseRows, chunk...)
		}
		clauseByID := make(map[uint]*models.ComplianceClause, len(clauseRows))
		for i := range clauseRows {
			clauseByID[clauseRows[i].ID] = &clauseRows[i]
		}
		for _, cc := range ccByID {
			if cc.ComplianceClauseID != nil {
				cc.ComplianceClause = clauseByID[*cc.ComplianceClauseID]
			}
		}
	}

	// 5) Load RU-like relations on control row: FrameworkValue, Status, CompletedByUser
	// 5a) Framework values
	if len(fvIDs) > 0 {
		var fvRows []models.ComplianceAssessmentFrameworkValue
		for i := 0; i < len(fvIDs); i += chunkSize {
			end := min(i+chunkSize, len(fvIDs))
			var chunk []models.ComplianceAssessmentFrameworkValue
			if err := db.Where("id IN ?", fvIDs[i:end]).Find(&chunk).Error; err != nil {
				return nil, err
			}
			fvRows = append(fvRows, chunk...)
		}
		fvByID := make(map[uint]*models.ComplianceAssessmentFrameworkValue, len(fvRows))
		for i := range fvRows {
			fvByID[fvRows[i].ID] = &fvRows[i]
		}
		for _, c := range controls {
			if c.ComplianceAssessmentFrameworkValueID != nil {
				c.ComplianceAssessmentFrameworkValue = fvByID[*c.ComplianceAssessmentFrameworkValueID]
			}
		}
	}

	// 5b) Statuses
	if len(statusIDs) > 0 {
		var stRows []models.ComplianceAssessmentRequirementStatus
		for i := 0; i < len(statusIDs); i += chunkSize {
			end := min(i+chunkSize, len(statusIDs))
			var chunk []models.ComplianceAssessmentRequirementStatus
			if err := db.Where("id IN ?", statusIDs[i:end]).Find(&chunk).Error; err != nil {
				return nil, err
			}
			stRows = append(stRows, chunk...)
		}
		stByID := make(map[uint]*models.ComplianceAssessmentRequirementStatus, len(stRows))
		for i := range stRows {
			stByID[stRows[i].ID] = &stRows[i]
		}
		for _, c := range controls {
			if c.ComplianceAssessmentRequirementStatusID != nil {
				c.ComplianceAssessmentRequirementStatus = stByID[*c.ComplianceAssessmentRequirementStatusID]
			}
		}
	}

	// 5c) Users
	if len(userIDs) > 0 {
		var users []models.User
		for i := 0; i < len(userIDs); i += chunkSize {
			end := min(i+chunkSize, len(userIDs))
			var chunk []models.User
			if err := db.Where("id IN ?", userIDs[i:end]).Find(&chunk).Error; err != nil {
				return nil, err
			}
			users = append(users, chunk...)
		}
		userByID := make(map[uint]*models.User, len(users))
		for i := range users {
			userByID[users[i].ID] = &users[i]
		}
		for _, c := range controls {
			if c.CompletedByUserID != nil {
				c.CompletedByUser = userByID[*c.CompletedByUserID]
			}
		}
	}

	// 6) Stitch ComplianceControl onto each AssessmentControl
	for _, c := range controls {
		c.ComplianceControl = ccByID[c.ComplianceControlID]
	}

	return controls, nil
}

func (repo *ComplianceAssessmentRepository) ComputeFlagsFast(assessmentID uint, tx ...utility.DatabaseTransaction) (isCompleted bool, isApproved bool, err error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// helper: run SELECT CASE WHEN EXISTS (<inner SQL>) THEN 1 ELSE 0 END
	exists := func(innerSQL string, args ...any) (bool, error) {
		var flag int
		// No ORDER BY / LIMIT here; MSSQL will be happy.
		q := "SELECT CASE WHEN EXISTS (" + innerSQL + ") THEN 1 ELSE 0 END"
		if err := db.Raw(q, args...).Scan(&flag).Error; err != nil {
			return false, err
		}
		return flag == 1, nil
	}

	// Any RU not completed?
	ruIncomplete, err := exists(`
        SELECT 1
        FROM compliance_assessment_requirement_responsible_units AS ru
        JOIN compliance_assessment_requirements AS r
             ON r.id = ru.compliance_assessment_requirement_id
        WHERE r.compliance_assessment_id = ? AND ru.is_completed = 0
    `, assessmentID)
	if err != nil {
		return
	}

	// Any RU not approved? (NULL status is treated as not approved)
	ruNotApproved, err := exists(`
        SELECT 1
        FROM compliance_assessment_requirement_responsible_units AS ru
        JOIN compliance_assessment_requirements AS r
             ON r.id = ru.compliance_assessment_requirement_id
        LEFT JOIN compliance_assessment_requirement_statuses AS s
               ON s.id = ru.compliance_assessment_requirement_status_id
        WHERE r.compliance_assessment_id = ?
          AND (s.id IS NULL OR s.type <> 'approved')
    `, assessmentID)
	if err != nil {
		return
	}

	// Any Control not completed?
	ctrlIncomplete, err := exists(`
        SELECT 1
        FROM compliance_assessment_controls AS ac
        WHERE ac.compliance_assessment_id = ? AND ac.is_completed = 0
    `, assessmentID)
	if err != nil {
		return
	}

	// Any Control not approved? (NULL status counts as not approved)
	ctrlNotApproved, err := exists(`
        SELECT 1
        FROM compliance_assessment_controls AS ac
        LEFT JOIN compliance_assessment_requirement_statuses AS s
               ON s.id = ac.compliance_assessment_requirement_status_id
        WHERE ac.compliance_assessment_id = ?
          AND (s.id IS NULL OR s.type <> 'approved')
    `, assessmentID)
	if err != nil {
		return
	}

	isCompleted = !(ruIncomplete || ctrlIncomplete)
	isApproved = !(ruNotApproved || ctrlNotApproved)
	return
}

func uniq(xs []uint) []uint {
	m := make(map[uint]struct{}, len(xs))
	out := make([]uint, 0, len(xs))
	for _, v := range xs {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (repo *ComplianceAssessmentRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	/* ================= ACTION PLAN TREE ================= */

	queries := []string{
		`DELETE compliance_action_plan_status_update_documents
		 FROM compliance_action_plan_status_update_documents
		 JOIN compliance_action_plan_status_updates ON compliance_action_plan_status_updates.id = compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		 JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		 JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		 JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		 WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_comment_documents
		FROM compliance_action_plan_comment_documents
		JOIN compliance_action_plan_comments ON compliance_action_plan_comments.id = compliance_action_plan_comment_documents.compliance_action_plan_comment_id
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_assessment_requirement_required_documents
		FROM compliance_assessment_requirement_required_documents
		JOIN compliance_assessment_requirement_document_titles ON compliance_assessment_requirement_document_titles.id = compliance_assessment_requirement_required_documents.compliance_assessment_requirement_document_title_id
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_document_titles.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN compliance_action_plans ON compliance_action_plans.id = tasks.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_status_updates
		FROM compliance_action_plan_status_updates
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_comments
		FROM compliance_action_plan_comments
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_status_update_documents
		FROM compliance_action_plan_status_update_documents
		JOIN compliance_action_plan_status_updates ON compliance_action_plan_status_updates.id = compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_comment_documents
		FROM compliance_action_plan_comment_documents
		JOIN compliance_action_plan_comments ON compliance_action_plan_comments.id = compliance_action_plan_comment_documents.compliance_action_plan_comment_id
		JOIN compliance_action_plans ON compliance_action_plans.id =compliance_action_plan_comments.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE compliance_assessment_requirement_document_titles
		FROM compliance_assessment_requirement_document_titles
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id =compliance_assessment_requirement_document_titles.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id =compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_assessment_requirement_comments
		FROM compliance_assessment_requirement_comments
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_comments.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_assessment_requirement_action_plans
		FROM compliance_assessment_requirement_action_plans
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_action_plans
		FROM compliance_action_plans
		JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN compliance_action_plans
		ON compliance_action_plans.id = tasks.compliance_action_plan_id
		JOIN compliance_assessment_controls
		ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_status_updates
		FROM compliance_action_plan_status_updates
		JOIN compliance_action_plans
		ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_controls
		ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE compliance_action_plan_comments
		FROM compliance_action_plan_comments
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE compliance_assessment_requirement_responsible_units
		FROM compliance_assessment_requirement_responsible_units
		JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		WHERE compliance_assessment_requirements.compliance_assessment_id = ?`,

		`DELETE compliance_assessment_control_comments
		FROM compliance_assessment_control_comments
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_assessment_control_comments.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE compliance_action_plans
		FROM compliance_action_plans
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		WHERE compliance_assessment_controls.compliance_assessment_id = ?`,

		`DELETE FROM compliance_assessment_requirements
		 WHERE compliance_assessment_id = ?`,

		`DELETE FROM compliance_assessment_controls
		 WHERE compliance_assessment_id = ?`,
	}

	for _, q := range queries {
		if err := db.Exec(q, id).Error; err != nil {
			return err
		}
	}

	/* ================= ROOT ================= */
	return db.Delete(&models.ComplianceAssessment{}, id).Error
}

func (repo *ComplianceAssessmentRepository) UpdateAverageResidualScoreAndRating(assessmentID uint, averageScore float64, ratingID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{
		"average_residual_risk_score":     averageScore,
		"average_residual_risk_rating_id": ratingID,
	}

	fmt.Printf("[REPO] Updating assessment %d with score=%v, ratingID=%v\n", assessmentID, averageScore, ratingID)

	result := db.Model(&models.ComplianceAssessment{}).
		Where("id = ?", assessmentID).
		Updates(updates)

	fmt.Printf("[REPO] Assessment update - Rows affected: %d, Error: %v\n", result.RowsAffected, result.Error)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) UpdateRequirementAverageResidualScoreAndRating(requirementID uint, averageScore float64, ratingID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{
		"average_residual_risk_score":     averageScore,
		"average_residual_risk_rating_id": ratingID,
	}

	fmt.Printf("[REPO] Updating requirement %d with score=%v, ratingID=%v\n", requirementID, averageScore, ratingID)

	result := db.Model(&models.ComplianceAssessmentRequirement{}).
		Where("id = ?", requirementID).
		Updates(updates)

	fmt.Printf("[REPO] Requirement update - Rows affected: %d, Error: %v\n", result.RowsAffected, result.Error)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) UpdateRegisterAverageResidualScoreAndRating(registerID uint, averageScore float64, ratingID *uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	updates := map[string]interface{}{
		"average_residual_risk_score":     averageScore,
		"average_residual_risk_rating_id": ratingID,
	}

	fmt.Printf("[REPO] Updating register %d with score=%v, ratingID=%v\n", registerID, averageScore, ratingID)

	result := db.Model(&models.ComplianceRegister{}).
		Where("id = ?", registerID).
		Updates(updates)

	fmt.Printf("[REPO] Register update - Rows affected: %d, Error: %v\n", result.RowsAffected, result.Error)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) GetComplianceAssessmentRatings() ([]*models.ComplianceAssessmentRating, error) {
	var ratings []*models.ComplianceAssessmentRating

	result := repo.db.Where("status = ?", "active").
		Order("score_from ASC").
		Find(&ratings)

	if result.Error != nil {
		return nil, result.Error
	}

	return ratings, nil
}

func (repo *ComplianceAssessmentRepository) GetRegisterIDsWithResponsibleUnitID(responsibleUnitID uint) ([]uint, error) {
	var registerIDs []uint

	result := repo.db.Table("compliance_requirement_responsible_units").
		Select("DISTINCT compliance_clauses.compliance_register_id").
		Joins("JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_responsible_units.compliance_requirement_id").
		Joins("JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id").
		Where("compliance_requirement_responsible_units.responsible_unit_id = ?", responsibleUnitID).
		Pluck("compliance_clauses.compliance_register_id", &registerIDs)

	if result.Error != nil {
		return nil, result.Error
	}

	return registerIDs, nil
}

func (repo *ComplianceAssessmentRepository) CreateComplianceAssessmentCompliances(compliances []*models.ComplianceAssessmentCompliance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(&compliances)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceAssessmentRepository) DeleteActionPlanCommentDocumentsByAssessment(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		DELETE capcd
		FROM compliance_action_plan_comment_documents capcd
		JOIN compliance_action_plan_comments capc
			ON capc.id = capcd.compliance_action_plan_comment_id
		JOIN compliance_action_plans cap
			ON cap.id = capc.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ru
			ON ru.id = cap.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements r
			ON r.id = ru.compliance_assessment_requirement_id
		WHERE r.compliance_assessment_id = ?
	`, id).Error
}

func (repo *ComplianceAssessmentRepository) DeleteActionPlanCommentsByAssessment(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		DELETE capc
		FROM compliance_action_plan_comments capc
		JOIN compliance_action_plans cap
			ON cap.id = capc.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ru
			ON ru.id = cap.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements r
			ON r.id = ru.compliance_assessment_requirement_id
		WHERE r.compliance_assessment_id = ?
	`, id).Error
}

func (repo *ComplianceAssessmentRepository) DeleteActionPlanStatusUpdateDocumentsByAssessment(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		DELETE capsud
		FROM compliance_action_plan_status_update_documents capsud
		JOIN compliance_action_plan_status_updates capsu
			ON capsu.id = capsud.compliance_action_plan_status_update_id
		JOIN compliance_action_plans cap
			ON cap.id = capsu.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ru
			ON ru.id = cap.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements r
			ON r.id = ru.compliance_assessment_requirement_id
		WHERE r.compliance_assessment_id = ?
	`, id).Error
}

func (repo *ComplianceAssessmentRepository) DeleteActionPlanStatusUpdatesByAssessment(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		DELETE capsu
		FROM compliance_action_plan_status_updates capsu
		JOIN compliance_action_plans cap
			ON cap.id = capsu.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units ru
			ON ru.id = cap.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements r
			ON r.id = ru.compliance_assessment_requirement_id
		WHERE r.compliance_assessment_id = ?
	`, id).Error
}

func (repo *ComplianceAssessmentRepository) DeleteActionPlansByAssessment(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		DELETE cap
		FROM compliance_action_plans cap
		JOIN compliance_assessment_requirement_responsible_units ru
			ON ru.id = cap.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements r
			ON r.id = ru.compliance_assessment_requirement_id
		WHERE r.compliance_assessment_id = ?
	`, id).Error
}

func (repo *ComplianceAssessmentRepository) GetComplianceAssessmentCompliaceByID(id uint) (*models.ComplianceAssessmentCompliance, error) {
	var compliance models.ComplianceAssessmentCompliance
	err := repo.db.Where("id = ?", id).First(&compliance).Error
	if err != nil {
		return nil, err
	}
	return &compliance, nil
}

func (repo *ComplianceAssessmentRepository) UpdateComplianceAssessmentCompliance(compliance *models.ComplianceAssessmentCompliance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(compliance).Error
}

func (repo *ComplianceAssessmentRepository) DeleteComplianceAssessmentCompliance(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", id).Delete(&models.ComplianceAssessmentCompliance{}).Error
}
