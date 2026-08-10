package repositories

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(identityReviewUserID uint, params *dtos.IdentityReviewActionPlanListQueryParamsDTO) (*dtos.IdentityReviewActionPlanListResponseDTO, error)
	GetByID(identityReviewUserID uint, id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewActionPlan, error)
	GetFullDetailsByID(identityReviewUserID uint, id uint) (*models.IdentityReviewActionPlan, error)
	Create(actionPlan *models.IdentityReviewActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(actionPlan *models.IdentityReviewActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateInitialStatusUpdate(update *models.IdentityReviewActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error
	DeleteIdentityReviewActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GetSummary(identityReviewUserID uint) (int, int, error)
}

type IdentityReviewActionPlanRepository struct {
	db *gorm.DB
}

func NewIdentityReviewActionPlanRepository(db *gorm.DB) *IdentityReviewActionPlanRepository {
	return &IdentityReviewActionPlanRepository{db: db}
}

func (repo *IdentityReviewActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IdentityReviewActionPlanRepository) GetAll(identityReviewUserID uint, params *dtos.IdentityReviewActionPlanListQueryParamsDTO) (*dtos.IdentityReviewActionPlanListResponseDTO, error) {
	var items []*dtos.IdentityReviewActionPlanListDTO
	var total int64

	query := repo.db.Table("identity_review_action_plans").
		Select(`
			identity_review_action_plans.id,
			identity_review_action_plans.identity_review_user_id,
			identity_review_users.username AS identity_review_username,
			identity_review_users.user_identifier,
			identity_review_action_plans.title,
			identity_review_action_plans.description,
			identity_review_action_plans.severity_level_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			identity_review_action_plans.responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			identity_review_action_plans.review_user_id,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,
			CONVERT(VARCHAR, identity_review_action_plans.target_date, 23) AS target_date,
			identity_review_action_plans.percentage,
			identity_review_action_plans.comments,
			identity_review_action_plans.action_plan_status_id,
			action_plan_statuses.title_en AS action_plan_status_title_en,
			action_plan_statuses.title_ar AS action_plan_status_title_ar,
			action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
			identity_review_action_plans.created_at,
			identity_review_action_plans.updated_at
		`).
		Joins("LEFT JOIN identity_review_users ON identity_review_users.id = identity_review_action_plans.identity_review_user_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = identity_review_action_plans.severity_level_id").
		Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = identity_review_action_plans.action_plan_status_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = identity_review_action_plans.responsible_user_id").
		Joins("LEFT JOIN users AS review_user ON review_user.id = identity_review_action_plans.review_user_id")
	if identityReviewUserID != 0 {
		query = query.Where("identity_review_action_plans.identity_review_user_id = ?", identityReviewUserID)
	}
	for field, values := range params.Filters {
		if len(values) == 0 {
			continue
		}
		switch field {
		case "severity_level_id":
			query = query.Where("identity_review_action_plans.severity_level_id IN (?)", values)
		case "action_plan_status_id":
			query = query.Where("identity_review_action_plans.action_plan_status_id IN (?)", values)
		case "responsible_user_id":
			query = query.Where("identity_review_action_plans.responsible_user_id IN (?)", values)
		case "review_user_id":
			query = query.Where("identity_review_action_plans.review_user_id IN (?)", values)
		case "identity_review_user_id":
			query = query.Where("identity_review_action_plans.identity_review_user_id IN (?)", values)
		}
	}
	if params.TargetDateFrom != "" {
		query = query.Where("identity_review_action_plans.target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("identity_review_action_plans.target_date <= ?", params.TargetDateTo)
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("identity_review_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	for _, filter := range params.QuickFilter {
		switch filter {
		case "overdue":
			query = query.Where("identity_review_action_plans.target_date < ? AND action_plan_statuses.type <> ?", time.Now(), "close")
		case "completed":
			query = query.Where("action_plan_statuses.type IN (?)", []string{"completed", "close"})
		case "in_progress":
			query = query.Where("action_plan_statuses.type = ?", "in_progress")
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(
			identity_review_action_plans.title LIKE ? OR
			identity_review_action_plans.description LIKE ? OR
			identity_review_users.username LIKE ? OR
			identity_review_users.user_identifier LIKE ? OR
			action_plan_statuses.title_en LIKE ? OR
			action_plan_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ?
		)`, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}
	query = query.Order(params.SortBy + " " + params.Sort)
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	if params.Export {
		return &dtos.IdentityReviewActionPlanListResponseDTO{Items: items}, nil
	}

	return &dtos.IdentityReviewActionPlanListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  int((total + int64(params.PageSize) - 1) / int64(params.PageSize)),
			TotalCount:  total,
		},
	}, nil
}

func (repo *IdentityReviewActionPlanRepository) GetSummary(identityReviewUserID uint) (int, int, error) {
	type summaryRow struct {
		OpenCount   int
		ClosedCount int
	}
	var row summaryRow
	query := repo.db.Table("identity_review_action_plans").
		Select(`
			COALESCE(SUM(CASE WHEN action_plan_statuses.is_closed = 0 THEN 1 ELSE 0 END), 0) AS open_count,
			COALESCE(SUM(CASE WHEN action_plan_statuses.is_closed = 1 THEN 1 ELSE 0 END), 0) AS closed_count
		`).
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = identity_review_action_plans.action_plan_status_id")
	if identityReviewUserID != 0 {
		query = query.Where("identity_review_action_plans.identity_review_user_id = ?", identityReviewUserID)
	}
	if err := query.Scan(&row).Error; err != nil {
		return 0, 0, err
	}
	return row.OpenCount, row.ClosedCount, nil
}

func (repo *IdentityReviewActionPlanRepository) GetByID(identityReviewUserID uint, id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var actionPlan models.IdentityReviewActionPlan
	if err := db.Where("identity_review_user_id = ? AND id = ?", identityReviewUserID, id).First(&actionPlan).Error; err != nil {
		return nil, err
	}
	return &actionPlan, nil
}

func (repo *IdentityReviewActionPlanRepository) GetFullDetailsByID(identityReviewUserID uint, id uint) (*models.IdentityReviewActionPlan, error) {
	var actionPlan models.IdentityReviewActionPlan
	err := repo.db.
		Preload("IdentityReviewUser").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ReviewUser").
		Preload("ActionPlanStatus").
		Preload("DelegateeUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ActionPlanStatusUpdates.ActionPlanStatus").
		Preload("ActionPlanStatusUpdates.CreatedByUser").
		Preload("ActionPlanStatusUpdates.UpdatedByUser").
		Preload("ActionPlanStatusUpdates.Documents").
		Where("identity_review_user_id = ? AND id = ?", identityReviewUserID, id).
		First(&actionPlan).Error
	if err != nil {
		return nil, err
	}
	return &actionPlan, nil
}

func (repo *IdentityReviewActionPlanRepository) Create(actionPlan *models.IdentityReviewActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(actionPlan).Error; err != nil {
		return 0, err
	}
	return actionPlan.ID, nil
}

func (repo *IdentityReviewActionPlanRepository) Update(actionPlan *models.IdentityReviewActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(actionPlan).Error
}

func (repo *IdentityReviewActionPlanRepository) CreateInitialStatusUpdate(update *models.IdentityReviewActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(update).Error
}

func (repo *IdentityReviewActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var actionPlan models.IdentityReviewActionPlan
	if err := db.First(&actionPlan, id).Error; err != nil {
		return err
	}

	return db.Delete(&models.IdentityReviewActionPlan{}, id).Error
}

func (repo *IdentityReviewActionPlanRepository) DeleteIdentityReviewActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	if len(ids) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", ids).Delete(&models.IdentityReviewActionPlanStatusUpdate{}).Error
}

func (repo *IdentityReviewActionPlanRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	queries := []string{
		`DELETE d FROM identity_review_action_plan_status_update_documents d JOIN identity_review_action_plan_status_updates u ON u.id = d.identity_review_action_plan_status_update_id WHERE u.identity_review_action_plan_id = ?`,
		`DELETE d FROM identity_review_action_plan_comment_documents d JOIN identity_review_action_plan_comments c ON c.id = d.identity_review_action_plan_comment_id WHERE c.identity_review_action_plan_id = ?`,
		`DELETE FROM tasks WHERE identity_review_action_plan_id = ?`,
		`DELETE FROM identity_review_action_plan_comments WHERE identity_review_action_plan_id = ?`,
		`DELETE FROM identity_review_action_plan_status_updates WHERE identity_review_action_plan_id = ?`,
	}
	for _, query := range queries {
		if err := db.Exec(query, id).Error; err != nil {
			return err
		}
	}
	return db.Delete(&models.IdentityReviewActionPlan{}, id).Error
}
