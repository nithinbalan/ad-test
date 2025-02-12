package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceAssessmentRequirementID uint, params *complianceDTOs.ComplianceAssessmentRequirementActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentRequirementActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirementActionPlan, error)
	Create(complianceAssessmentRequirementActionPlan *models.ComplianceAssessmentRequirementActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceAssessmentRequirementActionPlan *models.ComplianceAssessmentRequirementActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
}

type ComplianceAssessmentRequirementActionPlanRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRequirementActionPlanRepository(db *gorm.DB) *ComplianceAssessmentRequirementActionPlanRepository {
	return &ComplianceAssessmentRequirementActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceAssessmentRequirementActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentRequirementActionPlanRepository) GetAll(complianceAssessmentRequirementID uint, params *complianceDTOs.ComplianceAssessmentRequirementActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentRequirementActionPlanListResponseDTO, error) {

	var complianceAssessmentRequirementActionPlans []*complianceDTOs.ComplianceAssessmentRequirementActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("compliance_assessment_requirement_action_plans").
		Select(`
        compliance_assessment_requirement_action_plans.*,
		CONVERT(VARCHAR, compliance_assessment_requirement_action_plans.target_date, 23) AS target_date,
        compliance_requirements.title_en AS compliance_requirement_title_en,
        compliance_requirements.title_ar AS compliance_requirement_title_ar,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
		action_plan_statuses.type AS action_plan_status_type,
		responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = compliance_assessment_requirement_action_plans.action_plan_status_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_assessment_requirement_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = compliance_assessment_requirement_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = compliance_assessment_requirement_action_plans.updated_by_id
        LEFT JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_action_plans.compliance_assessment_requirement_id
		LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
    `).Where("compliance_assessment_requirement_action_plans.compliance_assessment_requirement_id = ?", complianceAssessmentRequirementID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("compliance_assessment_requirement_action_plans.%s IN (?)", columnName), values)
		}
	}

	// Apply target_date filters
	if params.TargetDateFrom != "" {
		query = query.Where("target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("target_date <= ?", params.TargetDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_assessment_requirement_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "overdue":
			query = query.Where("target_date < ?", time.Now())
		case "completed":
			query = query.Where("action_plan_statuses.type = ?", "completed")
		case "in_progress":
			query = query.Where("action_plan_statuses.type = ?", "in_progress")
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (compliance_assessment_requirement_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            compliance_assessment_requirement_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	err := query.Find(&complianceAssessmentRequirementActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssessmentRequirementActionPlanListResponseDTO{
		Items: complianceAssessmentRequirementActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceAssessmentRequirementActionPlanRepository) GetByID(id uint) (*models.ComplianceAssessmentRequirementActionPlan, error) {
	var complianceAssessmentRequirementActionPlan *models.ComplianceAssessmentRequirementActionPlan
	result := repo.db.First(&complianceAssessmentRequirementActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessmentRequirementActionPlan, nil
}

func (repo *ComplianceAssessmentRequirementActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirementActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceAssessmentRequirementActionPlan models.ComplianceAssessmentRequirementActionPlan

	// Load the compliance action plan with all related data
	result := db.
		Preload("ComplianceAssessmentRequirement").
		Preload("ComplianceAssessmentRequirement.ComplianceRequirement").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("ActionPlanStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&complianceAssessmentRequirementActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &complianceAssessmentRequirementActionPlan, nil
}

func (repo *ComplianceAssessmentRequirementActionPlanRepository) Create(complianceAssessmentRequirementActionPlan *models.ComplianceAssessmentRequirementActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceAssessmentRequirementActionPlan).Error; err != nil {
		return 0, err
	}
	return complianceAssessmentRequirementActionPlan.ID, nil
}

func (repo *ComplianceAssessmentRequirementActionPlanRepository) Update(complianceAssessmentRequirementActionPlan *models.ComplianceAssessmentRequirementActionPlan, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(complianceAssessmentRequirementActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceAssessmentRequirementActionPlanRepository) Delete(id uint) error {

	var complianceAssessmentRequirementActionPlan *models.ComplianceAssessmentRequirementActionPlan
	if err := repo.db.First(&complianceAssessmentRequirementActionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceAssessmentRequirementActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
