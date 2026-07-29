package repositories

import (
	"fmt"
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, params *assetDTOs.AssetRiskActionPlanListQueryParamsDTO) (*assetDTOs.AssetRiskActionPlanListResponseDTO, error)
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

type AssetRiskActionPlanRepository struct {
	db *gorm.DB
}

func NewAssetRiskActionPlanRepository(db *gorm.DB) *AssetRiskActionPlanRepository {
	return &AssetRiskActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AssetRiskActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskActionPlanRepository) GetAll(riskID uint, params *assetDTOs.AssetRiskActionPlanListQueryParamsDTO) (*assetDTOs.AssetRiskActionPlanListResponseDTO, error) {

	var riskActionPlans []*assetDTOs.AssetRiskActionPlanListDTO
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
		responsible_unit.id AS responsible_unit_id,
		responsible_unit.title_en AS responsible_unit_title_en,
		responsible_unit.title_ar AS responsible_unit_title_ar,
		CONVERT(VARCHAR, risk_action_plans.start_date, 23) AS start_date,
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
		LEFT JOIN organization_units AS responsible_unit ON responsible_unit.id = risk_action_plans.responsible_unit_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = risk_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = risk_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = risk_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = risk_action_plans.updated_by_id
    `).
		Where("risk_register.is_cyber_risk = ?", true)

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
			query = query.Where("severity_levels.type = ?", "high")
		case "overdue":
			query = query.Where("risk_action_plans.target_date < ?", time.Now())
		case "pending":
			query = query.Where("action_plan_statuses.type IN (?)", []string{"in_progress", "new"})
		case "upcoming_deadlines":
			query = query.Where("risk_action_plans.target_date BETWEEN ? AND ?", time.Now(), time.Now().AddDate(0, 0, 7))
		case "closed":
			query = query.Where("action_plan_statuses.type = ?", "close")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (risk_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			risk_register.title LIKE ? OR
            risk_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword)
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
		return &assetDTOs.AssetRiskActionPlanListResponseDTO{
			Items: riskActionPlans,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &assetDTOs.AssetRiskActionPlanListResponseDTO{
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

func (repo *AssetRiskActionPlanRepository) GetByID(id uint) (*models.RiskActionPlan, error) {
	var riskActionPlan *models.RiskActionPlan
	result := repo.db.First(&riskActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlan, nil
}

func (repo *AssetRiskActionPlanRepository) GetFullDetailsByID(id uint) (*models.RiskActionPlan, error) {
	var riskActionPlan *models.RiskActionPlan

	result := repo.db.
		Preload("RiskRegister").
		Preload("RiskRegister.Controls").
		Preload("SeverityLevel").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ActionPlanStatus").
		Preload("ActionPlanStatusUpdates").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Controls.Control.ControlType").
		Preload("Controls.Control.ResponsibleUser").
		Preload("Controls.Control.ResponsibleUnit").
		Preload("ReviewUser").
		Preload("DelegateeUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&riskActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlan, nil
}

func (repo *AssetRiskActionPlanRepository) Create(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskActionPlan).Error; err != nil {
		return 0, err
	}
	return riskActionPlan.ID, nil
}

func (repo *AssetRiskActionPlanRepository) Update(riskActionPlan *models.RiskActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
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

func (repo *AssetRiskActionPlanRepository) CreateInitialRiskActionPlanStatusUpdate(riskActionPlanStatusUpdate *models.RiskActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return riskActionPlanStatusUpdate.ID, nil
}

func (repo *AssetRiskActionPlanRepository) DeleteRiskActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.RiskActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskActionPlanRepository) CreateRiskActionPlanControl(riskActionPlanControl *models.RiskActionPlanControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&riskActionPlanControl).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskActionPlanRepository) DeleteRiskActionPlanControl(riskActionPlanControlIds []*uint, tx ...utility.DatabaseTransaction) error {
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

func (repo *AssetRiskActionPlanRepository) CreateRiskControls(controls *models.RiskControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskActionPlanRepository) DeleteRiskControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", controlIds).Delete(&models.RiskControl{}).Error; err != nil {
		return err
	}

	return nil
}
