package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type StakeholderActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderActionPlanListQueryParamsDTO, stakeholderAssessmentRequirementID uint) (*organizationalContextDTOs.StakeholderActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.StakeholderActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderActionPlan, error)
	Create(stakeholderActionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(stakeholderActionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error

	CreateInitialStakeholderActionPlanStatusUpdate(stakeholderActionPlanStatusUpdate *models.StakeholderActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteStakeholderActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteStakeholderActionPlanComments(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteStakeholderActionPlanCommentDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	DeleteStakeholderActionPlanStatusUpdateDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type StakeholderActionPlanRepository struct {
	db *gorm.DB
}

func NewStakeholderActionPlanRepository(db *gorm.DB) *StakeholderActionPlanRepository {
	return &StakeholderActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StakeholderActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderActionPlanRepository) GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderActionPlanListQueryParamsDTO, stakeholderAssessmentRequirementID uint) (*organizationalContextDTOs.StakeholderActionPlanListResponseDTO, error) {

	var stakeholderActionPlans []*organizationalContextDTOs.StakeholderActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("stakeholder_action_plans").
		Select(`
        stakeholder_action_plans.*,
		CONVERT(VARCHAR, stakeholder_action_plans.target_date, 23) AS target_date,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
		action_plan_statuses.type AS action_plan_status_type,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path,
		review_user.id AS review_user_id,
    	review_user.first_name_en AS review_user_first_name_en,
		review_user.first_name_ar AS review_user_first_name_ar,
		review_user.last_name_en AS review_user_last_name_en,
		review_user.last_name_ar AS review_user_last_name_ar,
		review_user.initial AS review_user_initial,
        review_user.color AS review_user_color,
        review_user.email AS review_user_email,
        review_user.profile_image_path AS review_user_profile_image_path,
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
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = stakeholder_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = stakeholder_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = stakeholder_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = stakeholder_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = stakeholder_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_action_plans.updated_by_id
        LEFT JOIN stakeholder_requirements ON stakeholder_requirements.id = stakeholder_action_plans.stakeholder_requirement_id
        LEFT JOIN stakeholder_clauses ON stakeholder_clauses.id = stakeholder_requirements.stakeholder_clause_id
    `)

	// Conditionally apply stakeholder ID filter
	if stakeholderID != 0 {
		query = query.Where("stakeholder_action_plans.stakeholder_register_id = ?", stakeholderID)
	}

	// Conditionally apply stakeholder ID filter
	if stakeholderAssessmentRequirementID != 0 {
		query = query.Where("stakeholder_action_plans.stakeholder_assessment_requirement_id = ?", stakeholderAssessmentRequirementID)
	}

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("stakeholder_action_plans.%s IN (?)", columnName), values)
		}
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
            (stakeholder_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            stakeholder_action_plans.target_date LIKE ? OR
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

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute the query
	err := query.Find(&stakeholderActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationalContextDTOs.StakeholderActionPlanListResponseDTO{
		Items: stakeholderActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *StakeholderActionPlanRepository) GetByID(id uint) (*models.StakeholderActionPlan, error) {
	var stakeholderActionPlan *models.StakeholderActionPlan
	result := repo.db.First(&stakeholderActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return stakeholderActionPlan, nil
}

func (repo *StakeholderActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var stakeholderActionPlan models.StakeholderActionPlan

	// Load the stakeholder action plan with all related data
	result := db.
		Preload("StakeholderRequirement.StakeholderClause").
		Preload("StakeholderAssessmentRequirement").
		Preload("StakeholderAssessmentRequirement.StakeholderRequirement").
		Preload("StakeholderRegister").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("ActionPlanStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ActionPlanStatusUpdates").
		Preload("ReviewUser").
		Preload("ReviewUser.OrganizationUnit1").
		First(&stakeholderActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &stakeholderActionPlan, nil
}

func (repo *StakeholderActionPlanRepository) Create(stakeholderActionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(stakeholderActionPlan).Error; err != nil {
		return 0, err
	}
	return stakeholderActionPlan.ID, nil
}

func (repo *StakeholderActionPlanRepository) Update(stakeholderActionPlan *models.StakeholderActionPlan, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(stakeholderActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StakeholderActionPlanRepository) Delete(id uint) error {

	var stakeholderActionPlan *models.StakeholderActionPlan
	if err := repo.db.First(&stakeholderActionPlan, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.StakeholderActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *StakeholderActionPlanRepository) CreateInitialStakeholderActionPlanStatusUpdate(stakeholderActionPlanStatusUpdate *models.StakeholderActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(stakeholderActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return stakeholderActionPlanStatusUpdate.ID, nil
}

func (repo *StakeholderActionPlanRepository) DeleteStakeholderActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.StakeholderActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanRepository) DeleteStakeholderActionPlanComments(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.StakeholderActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanRepository) DeleteStakeholderActionPlanCommentDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.StakeholderActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanRepository) DeleteStakeholderActionPlanStatusUpdateDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.StakeholderActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}
