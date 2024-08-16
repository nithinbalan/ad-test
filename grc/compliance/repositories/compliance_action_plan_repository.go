package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ComplianceActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceActionPlan, error)
	Create(complianceActionPlan *models.ComplianceActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceActionPlan *models.ComplianceActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error

	CreateInitialComplianceActionPlanStatusUpdate(complianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteComplianceActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceActionPlanRepository struct {
	db *gorm.DB
}

func NewComplianceActionPlanRepository(db *gorm.DB) *ComplianceActionPlanRepository {
	return &ComplianceActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceActionPlanRepository) GetAll(complianceID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error) {

	var complianceActionPlans []*complianceDTOs.ComplianceActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("compliance_action_plans").
		Select(`
        compliance_action_plans.*,
		CONVERT(VARCHAR, compliance_action_plans.target_date, 23) AS target_date,
        compliance_requirements.title AS compliance_requirement_title,
        compliance_clauses.clause_no AS compliance_clause_no,
        compliance_clauses.title AS compliance_clause_title,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = compliance_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = compliance_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = compliance_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = compliance_action_plans.updated_by_id
        LEFT JOIN compliance_requirements ON compliance_requirements.id = compliance_action_plans.compliance_requirement_id
        LEFT JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
    `).Where(`compliance_action_plans.compliance_register_id`, complianceID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply target_date filters
	if params.TargetDateFrom != "" {
		query = query.Where("target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("target_date <= ?", params.TargetDateTo)
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
            (compliance_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            compliance_action_plans.target_date LIKE ? OR
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
			keyword, keyword, keyword, keyword,
			keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&complianceActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceActionPlanListResponseDTO{
		Items: complianceActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceActionPlanRepository) GetByID(id uint) (*models.ComplianceActionPlan, error) {
	var complianceActionPlan *models.ComplianceActionPlan
	result := repo.db.First(&complianceActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceActionPlan, nil
}

func (repo *ComplianceActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceActionPlan models.ComplianceActionPlan

	// Load the compliance action plan with all related data
	result := db.
		Preload("ComplianceRequirement.ComplianceClause").
		Preload("ComplianceRegister").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ActionPlanStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ActionPlanStatusUpdates").
		First(&complianceActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &complianceActionPlan, nil
}

func (repo *ComplianceActionPlanRepository) Create(complianceActionPlan *models.ComplianceActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceActionPlan).Error; err != nil {
		return 0, err
	}
	return complianceActionPlan.ID, nil
}

func (repo *ComplianceActionPlanRepository) Update(complianceActionPlan *models.ComplianceActionPlan, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(complianceActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanRepository) Delete(id uint) error {

	var complianceActionPlan *models.ComplianceActionPlan
	if err := repo.db.First(&complianceActionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceActionPlanRepository) CreateInitialComplianceActionPlanStatusUpdate(complianceActionPlanStatusUpdate *models.ComplianceActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return complianceActionPlanStatusUpdate.ID, nil
}

func (repo *ComplianceActionPlanRepository) DeleteComplianceActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}
