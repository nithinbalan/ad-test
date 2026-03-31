package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ViolationReportActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportActionPlanListQueryParamsDTO) (*violationReportDTOs.ViolationReportActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.ViolationReportActionPlan, error)
	GetFullDetailsByID(id uint) (*models.ViolationReportActionPlan, error)
	//GetUsedActionPlanMasterIDsByViolationReportID(violationReportID uint) ([]uint, error)
	Create(violationReportActionPlan *models.ViolationReportActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(violationReportActionPlan *models.ViolationReportActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	//GetCountByRiskID(violationReportID uint, tx ...utility.DatabaseTransaction) (int64, error)
	//HasActionPlanMasterInViolationReport(violationReportID uint, actionPlanMasterID uint) (bool, error)
	GetAllByViolationReportID(violationReportID uint, tx ...utility.DatabaseTransaction) ([]*models.ViolationReportActionPlan, error)

	CreateInitialViolationReportActionPlanStatusUpdate(statusUpdate *models.ViolationReportActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteViolationReportActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type ViolationReportActionPlanRepository struct {
	db *gorm.DB
}

func NewViolationReportActionPlanRepository(db *gorm.DB) *ViolationReportActionPlanRepository {
	return &ViolationReportActionPlanRepository{db: db}
}

func (repo *ViolationReportActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ViolationReportActionPlanRepository) GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportActionPlanListQueryParamsDTO) (*violationReportDTOs.ViolationReportActionPlanListResponseDTO, error) {
	var items []*violationReportDTOs.ViolationReportActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("violation_report_action_plans").
		Select(`
		violation_report_action_plans.*,
		violation_reports.reference_number AS violation_report_reference_number,
		violation_reports.title_en AS violation_report_title_en,
		violation_reports.title_ar AS violation_report_title_ar,
		CONVERT(VARCHAR, violation_report_action_plans.target_date, 23) AS target_date,
		action_plan_statuses.title_en AS action_plan_status_title_en,
		action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
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
		violation_report_action_plans.created_by_id,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		created_user.initial AS created_user_initial,
		created_user.color AS created_user_color,
		created_user.email AS created_user_email,
		created_user.profile_image_path AS created_user_profile_image_path,

		violation_report_action_plans.updated_by_id,
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
		LEFT JOIN violation_reports ON violation_reports.id = violation_report_action_plans.violation_report_id
		LEFT JOIN action_plan_statuses ON action_plan_statuses.id = violation_report_action_plans.action_plan_status_id
		LEFT JOIN severity_levels ON severity_levels.id = violation_report_action_plans.severity_level_id
		LEFT JOIN users AS responsible_user ON responsible_user.id = violation_report_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = violation_report_action_plans.review_user_id
		LEFT JOIN users AS created_user ON created_user.id = violation_report_action_plans.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = violation_report_action_plans.updated_by_id
		`)

	fieldMappings := map[string]string{
		"severity_level_id":     "violation_report_action_plans.severity_level_id",
		"action_plan_status_id": "violation_report_action_plans.action_plan_status_id",
		"responsible_user_id":   "violation_report_action_plans.responsible_user_id",
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if mapped, ok := fieldMappings[columnName]; ok {
				query = query.Where(mapped+" IN (?)", values)
			} else {
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	if params.TargetDateFrom != "" {
		query = query.Where("target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("target_date <= ?", params.TargetDateTo)
	}

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

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
		violation_report_action_plans.title LIKE ? OR
		action_plan_statuses.title_en LIKE ? OR
		action_plan_statuses.title_ar LIKE ?`, kw, kw, kw)
	}

	if violationReportID != 0 {
		query = query.Where("violation_report_action_plans.violation_report_id = ?", violationReportID)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &violationReportDTOs.ViolationReportActionPlanListResponseDTO{Items: items}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &violationReportDTOs.ViolationReportActionPlanListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ViolationReportActionPlanRepository) GetByID(id uint) (*models.ViolationReportActionPlan, error) {
	var plan *models.ViolationReportActionPlan
	if err := repo.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

func (repo *ViolationReportActionPlanRepository) GetFullDetailsByID(id uint) (*models.ViolationReportActionPlan, error) {
	var plan *models.ViolationReportActionPlan
	if err := repo.db.
		Preload("ViolationReport").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("ActionPlanStatus").
		Preload("ActionPlanStatusUpdates").
		Preload("ReviewUser").
		First(&plan, id).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

func (repo *ViolationReportActionPlanRepository) Create(plan *models.ViolationReportActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(plan).Error; err != nil {
		return 0, err
	}
	return plan.ID, nil
}

func (repo *ViolationReportActionPlanRepository) Update(plan *models.ViolationReportActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(plan).Error
}

func (repo *ViolationReportActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ViolationReportActionPlan{}, id).Error
}

func (repo *ViolationReportActionPlanRepository) HasActionPlanMasterInViolationReport(violationReportID uint, actionPlanMasterID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportActionPlan{}).
		Where("violation_report_id = ? AND action_plan_master_id = ?", violationReportID, actionPlanMasterID).
		Count(&count).Error
	return count > 0, err
}

func (repo *ViolationReportActionPlanRepository) GetAllByViolationReportID(violationReportID uint, tx ...utility.DatabaseTransaction) ([]*models.ViolationReportActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var plans []*models.ViolationReportActionPlan
	err := db.Where("violation_report_id = ?", violationReportID).Order("created_at ASC").Find(&plans).Error
	return plans, err
}

func (repo *ViolationReportActionPlanRepository) CreateInitialViolationReportActionPlanStatusUpdate(update *models.ViolationReportActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(update).Error; err != nil {
		return 0, err
	}
	return update.ID, nil
}

func (repo *ViolationReportActionPlanRepository) DeleteViolationReportActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", ids).Delete(&models.ViolationReportActionPlanStatusUpdate{}).Error
}

func (repo *ViolationReportActionPlanRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	/* ================= ACTION PLAN TREE ================= */

	queries := []string{

		// Action plan comment documents
		`DELETE violation_report_action_plan_status_update_documents
		 FROM violation_report_action_plan_status_update_documents 
		 JOIN violation_report_action_plan_status_updates ON violation_report_action_plan_status_updates.id = violation_report_action_plan_status_update_documents.violation_report_action_plan_status_update_id
		 WHERE violation_report_action_plan_status_updates.violation_report_action_plan_id = ?`,

		`DELETE violation_report_action_plan_comment_documents
		 FROM violation_report_action_plan_comment_documents 
		 JOIN violation_report_action_plan_comments ON violation_report_action_plan_comments.id = violation_report_action_plan_comment_documents.violation_report_action_plan_comment_id
		 WHERE violation_report_action_plan_comments.violation_report_action_plan_id = ?`,
		`DELETE FROM tasks WHERE tasks.violation_report_action_plan_id = ?`,
		`DELETE FROM violation_report_action_plan_comments WHERE violation_report_action_plan_comments.violation_report_action_plan_id = ?`,
		`DELETE FROM violation_report_action_plan_status_updates WHERE violation_report_action_plan_status_updates.violation_report_action_plan_id = ?`,
	}

	for _, q := range queries {
		if err := db.Exec(q, id).Error; err != nil {
			return err
		}
	}

	/* ================= ROOT ================= */
	return db.Delete(&models.ViolationReportActionPlan{}, id).Error
}
