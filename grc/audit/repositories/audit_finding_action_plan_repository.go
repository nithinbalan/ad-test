package repositories

import (
	auditFindingDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type AuditFindingActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(findingID uint, params *auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditFindingActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.AuditFindingActionPlan, error)
	GetFullDetailsByID(id uint) (*models.AuditFindingActionPlan, error)
	Create(auditFindingActionPlan *models.AuditFindingActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFindingActionPlan *models.AuditFindingActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetCountByRiskID(findingID uint, tx ...utility.DatabaseTransaction) (int64, error)

	CreateInitialAuditFindingActionPlanStatusUpdate(auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteAuditFindingActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditFindingActionPlanRepository struct {
	db *gorm.DB
}

func NewAuditFindingActionPlanRepository(db *gorm.DB) *AuditFindingActionPlanRepository {
	return &AuditFindingActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFindingActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFindingActionPlanRepository) GetAll(findingID uint, params *auditFindingDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditFindingActionPlanListResponseDTO, error) {

	var auditFindingActionPlans []*auditFindingDTOs.AuditFindingActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("audit_finding_action_plans").
		Select(`
        audit_finding_action_plans.*,
		audit_findings.title AS audit_finding_title,
		audit_findings.reference_number AS audit_finding_reference_number,
		audit_findings.description AS audit_finding_description,
		CONVERT(VARCHAR, audit_finding_action_plans.target_date, 23) AS target_date,
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
		LEFT JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = audit_finding_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = audit_finding_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = audit_finding_action_plans.responsible_user_id
        LEFT JOIN users AS created_user ON created_user.id = audit_finding_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = audit_finding_action_plans.updated_by_id
    `).Where("audit_finding_action_plans.audit_finding_id = ?", findingID)

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
            (audit_finding_action_plans.title LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ? OR
            audit_finding_action_plans.target_date LIKE ? OR
            responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            responsible_user.last_name_en LIKE ? OR
            responsible_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword)
	}
	// Conditionally apply finding ID filter
	if findingID != 0 {
		query = query.Where("audit_finding_action_plans.audit_finding_id = ?", findingID)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditFindingActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditFindingDTOs.AuditFindingActionPlanListResponseDTO{
		Items: auditFindingActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditFindingActionPlanRepository) GetByID(id uint) (*models.AuditFindingActionPlan, error) {
	var auditFindingActionPlan *models.AuditFindingActionPlan
	result := repo.db.First(&auditFindingActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingActionPlan, nil
}

func (repo *AuditFindingActionPlanRepository) GetFullDetailsByID(id uint) (*models.AuditFindingActionPlan, error) {
	var auditFindingActionPlan *models.AuditFindingActionPlan

	result := repo.db.
		Preload("AuditFinding").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ActionPlanStatusUpdates").
		First(&auditFindingActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingActionPlan, nil
}

func (repo *AuditFindingActionPlanRepository) Create(auditFindingActionPlan *models.AuditFindingActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFindingActionPlan).Error; err != nil {
		return 0, err
	}
	return auditFindingActionPlan.ID, nil
}

func (repo *AuditFindingActionPlanRepository) Update(auditFindingActionPlan *models.AuditFindingActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFindingActionPlan).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditFindingActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditFindingActionPlan *models.AuditFindingActionPlan
	if err := db.First(&auditFindingActionPlan, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.AuditFindingActionPlan{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *AuditFindingActionPlanRepository) CreateInitialAuditFindingActionPlanStatusUpdate(auditFindingActionPlanStatusUpdate *models.AuditFindingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFindingActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}
	return auditFindingActionPlanStatusUpdate.ID, nil
}

func (repo *AuditFindingActionPlanRepository) DeleteAuditFindingActionPlanStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.AuditFindingActionPlanStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil

}

func (repo *AuditFindingActionPlanRepository) GetCountByRiskID(findingID uint, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	if err := db.Model(&models.AuditFindingActionPlan{}).
		Where("audit_finding_id = ?", findingID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
