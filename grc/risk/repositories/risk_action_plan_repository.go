package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.RiskActionPlan, error)
	GetFullDetailsByID(id uint) (*models.RiskActionPlan, error)
	Create(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateInitialRiskActionPlanStatusUpdate(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteRiskActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error

	CreateRiskActionPlanControl(riskActionPlanControl *models.RiskActionPlanControl, tx ...utility.DatabaseTransaction) error
	DeleteRiskActionPlanControl(riskActionPlanControlIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateRiskControls(controls *models.RiskControl, tx ...utility.DatabaseTransaction) error
	DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error
}

type RiskActionPlanRepository struct {
	db *gorm.DB
}

func NewRiskActionPlanRepository(db *gorm.DB) *RiskActionPlanRepository {
	return &RiskActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *RiskActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskActionPlanRepository) GetAll(riskID uint, params *riskDTOs.RiskActionPlanListQueryParamsDTO) (*riskDTOs.RiskActionPlanListResponseDTO, error) {

	var riskActionPlans []*riskDTOs.RiskActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("risk_action_plans").
		Select(`
        risk_action_plans.*,
		risk_register.reference_number AS risk_reference_number,
        risk_register.title AS risk_title,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		CONVERT(VARCHAR, risk_action_plans.target_date, 23) AS target_date,
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
		LEFT JOIN risk_register ON risk_register.id = risk_action_plans.risk_register_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = risk_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = risk_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = risk_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = risk_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = risk_action_plans.updated_by_id
    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("risk_action_plans.%s IN (?)", columnName), values)
		}
	}

	if params.TargetDateFrom != "" {
		query = query.Where("risk_action_plans.target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("risk_action_plans.target_date <= ?", params.TargetDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "critical":
			var highestRating string
			err := repo.db.Table("risk_matrix_risk_ratings").
				Select("type").
				Where("status = ?", "active").
				Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
				Limit(1).Pluck("type", &highestRating).Error
			if err != nil || highestRating == "" {
				return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
			}

			query = query.Joins(`
				JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
			`).Where("risk_matrix_risk_ratings.type = ?", highestRating)

		case "overdue":
			query = query.Where("risk_action_plans.target_date < ?", time.Now())
		case "pending":
			query = query.Where("action_plan_statuses.type IN (?)", []string{"in_progress", "new"})
		case "upcoming_deadlines":
			query = query.Where("risk_action_plans.target_date BETWEEN ? AND ?", time.Now(), time.Now().AddDate(0, 0, 7))
		case "closed":
			query = query.Where("action_plan_statuses.type = ?", "close")
		case "strategic_mitigations":
			query = query.Joins(`
				LEFT JOIN (
					SELECT DISTINCT risk_action_plans.id AS action_plan_id
					FROM risk_action_plans
					JOIN risk_register ON risk_register.id = risk_action_plans.risk_register_id
					JOIN risk_strategic_objectives ON risk_strategic_objectives.risk_register_id = risk_register.id
				) AS strategic_action_plans ON strategic_action_plans.action_plan_id = risk_action_plans.id
			`).Where("strategic_action_plans.action_plan_id IS NOT NULL")

		case "financial_impact":
			query = query.Joins(`
				JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			`).Where("risk_categories.title_en = ?", "Financial Risk")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (risk_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            risk_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword)
	}
	// Conditionally apply risk ID filter
	if riskID != 0 {
		query = query.Where("risk_action_plans.risk_register_id = ?", riskID)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&riskActionPlans).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &riskDTOs.RiskActionPlanListResponseDTO{
			Items: riskActionPlans,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskActionPlanListResponseDTO{
		Items: riskActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *RiskActionPlanRepository) GetByID(id uint) (*models.RiskActionPlan, error) {
	var riskActionPlan *models.RiskActionPlan
	result := repo.db.First(&riskActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlan, nil
}

func (repo *RiskActionPlanRepository) GetFullDetailsByID(id uint) (*models.RiskActionPlan, error) {
	var riskActionPlan *models.RiskActionPlan

	result := repo.db.
		Preload("RiskRegister").
		Preload("RiskRegister.Controls").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("ActionPlanStatusUpdates").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Controls.Control.ControlType").
		Preload("Controls.Control.Owner").
		Preload("Controls.Control.OwnerUnit").
		Preload("ReviewUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&riskActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlan, nil
}

func (repo *RiskActionPlanRepository) Create(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskActionPlan).Error; err != nil {
		return 0, err
	}
	return riskActionPlan.ID, nil
}

func (repo *RiskActionPlanRepository) Update(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var riskActionPlan *models.RiskActionPlan
	if err := db.First(&riskActionPlan, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.RiskActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskActionPlanRepository) CreateInitialRiskActionPlanStatusUpdate(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return riskActionPlanStatusUpdate.ID, nil
}

func (repo *RiskActionPlanRepository) DeleteRiskActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.RiskActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskActionPlanRepository) CreateRiskActionPlanControl(riskActionPlanControl *models.RiskActionPlanControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&riskActionPlanControl).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskActionPlanRepository) DeleteRiskActionPlanControl(riskActionPlanControlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", riskActionPlanControlIds).Delete(&models.RiskActionPlanControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskActionPlanRepository) CreateRiskControls(controls *models.RiskControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *RiskActionPlanRepository) DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", controlIds).Delete(&models.RiskControl{}).Error; err != nil {
		return err
	}

	return nil
}
