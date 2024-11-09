package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ThreatActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(threatID uint, params *assetDTOs.ThreatActionPlanListQueryParamsDTO) (*assetDTOs.ThreatActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ThreatActionPlan, error)
	GetFullDetailsByID(id uint) (*models.ThreatActionPlan, error)
	Create(threatActionPlan *models.ThreatActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(threatActionPlan *models.ThreatActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateInitialThreatActionPlanStatusUpdate(threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteThreatActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type ThreatActionPlanRepository struct {
	db *gorm.DB
}

func NewThreatActionPlanRepository(db *gorm.DB) *ThreatActionPlanRepository {
	return &ThreatActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ThreatActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThreatActionPlanRepository) GetAll(threatID uint, params *assetDTOs.ThreatActionPlanListQueryParamsDTO) (*assetDTOs.ThreatActionPlanListResponseDTO, error) {
	var threatActionPlans []*assetDTOs.ThreatActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("threat_action_plans").
		Select(`
        threat_action_plans.*,
        threats.title_en AS threat_title_en,
        threats.title_ar AS threat_title_ar,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		CONVERT(VARCHAR, threat_action_plans.target_date, 23) AS target_date,
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
		LEFT JOIN threats ON threats.id = threat_action_plans.threat_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = threat_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = threat_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = threat_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = threat_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = threat_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = threat_action_plans.updated_by_id
    `)

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

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (threat_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            threat_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword)
	}
	// Conditionally apply threat ID filter
	if threatID != 0 {
		query = query.Where("threat_action_plans.threat_id = ?", threatID)
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

	if err := query.Find(&threatActionPlans).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &assetDTOs.ThreatActionPlanListResponseDTO{
			Items: threatActionPlans,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &assetDTOs.ThreatActionPlanListResponseDTO{
		Items: threatActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ThreatActionPlanRepository) GetByID(id uint) (*models.ThreatActionPlan, error) {
	var threatActionPlan *models.ThreatActionPlan
	result := repo.db.First(&threatActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return threatActionPlan, nil
}

func (repo *ThreatActionPlanRepository) GetFullDetailsByID(id uint) (*models.ThreatActionPlan, error) {
	var threatActionPlan *models.ThreatActionPlan

	result := repo.db.
		Preload("Threat").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("ActionPlanStatusUpdates").
		Preload("ReviewUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&threatActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return threatActionPlan, nil
}

func (repo *ThreatActionPlanRepository) Create(threatActionPlan *models.ThreatActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(threatActionPlan).Error; err != nil {
		return 0, err
	}
	return threatActionPlan.ID, nil
}

func (repo *ThreatActionPlanRepository) Update(threatActionPlan *models.ThreatActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(threatActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var threatActionPlan *models.ThreatActionPlan
	if err := db.First(&threatActionPlan, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.ThreatActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ThreatActionPlanRepository) CreateInitialThreatActionPlanStatusUpdate(threatActionPlanStatusUpdate *models.ThreatActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(threatActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return threatActionPlanStatusUpdate.ID, nil
}

func (repo *ThreatActionPlanRepository) DeleteThreatActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ThreatActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}
