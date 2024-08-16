package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanBudgetRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanBudgetListQueryParamsDTO) (*auditDTOs.AuditPlanBudgetListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlanBudget, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanBudget, error)
	Create(auditPlanBudget *models.AuditPlanBudget, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditPlanBudget *models.AuditPlanBudget, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditPlanBudgetRepository struct {
	db *gorm.DB
}

func NewAuditPlanBudgetRepository(db *gorm.DB) *AuditPlanBudgetRepository {
	return &AuditPlanBudgetRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanBudgetRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanBudgetRepository) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanBudgetListQueryParamsDTO) (*auditDTOs.AuditPlanBudgetListResponseDTO, error) {
	var auditPlanBudgets []*auditDTOs.AuditPlanBudgetListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_budgets").
		Select(`
			audit_plan_budgets.id,
			audit_plan_budgets.audit_plan_id,
			audit_plan_budgets.estimated_cost,
			audit_plan_budgets.expense_category_id,
			expense_categories.title_en AS expense_category_title_en,
			expense_categories.title_ar AS expense_category_title_ar,
			expense_categories.description_en AS expense_category_description_en,
			expense_categories.description_ar AS expense_category_description_ar,
			audit_plan_budgets.title,
			audit_plan_budgets.comments,
			audit_plans.title AS audit_plan_title,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, audit_plan_budgets.created_at, 23) AS created_at,
			CONVERT(VARCHAR, audit_plan_budgets.updated_at, 23) AS updated_at
		`).
		Joins(`
			LEFT JOIN expense_categories ON expense_categories.id = audit_plan_budgets.expense_category_id
			LEFT JOIN audit_plans ON audit_plans.id = audit_plan_budgets.audit_plan_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_budgets.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_budgets.updated_by_id
		`).Where(`audit_plan_budgets.audit_plan_id`, auditPlanID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "active":
			query = query.Where("expense_categories.status = ?", "active")
		case "inactive":
			query = query.Where("expense_categories.status = ?", "inactive")
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
			audit_plan_budgets.title LIKE ? OR
			audit_plan_budgets.comments LIKE ? OR
			audit_plan_budgets.estimated_cost LIKE ? OR
			audit_plans.title LIKE ? OR
			expense_categories.title_en LIKE ? OR
			expense_categories.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditPlanBudgets).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanBudgetListResponseDTO{
		Items: auditPlanBudgets,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditPlanBudgetRepository) GetByID(id uint) (*models.AuditPlanBudget, error) {
	var auditPlanBudget *models.AuditPlanBudget

	result := repo.db.First(&auditPlanBudget, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanBudget, nil
}

func (repo *AuditPlanBudgetRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanBudget, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditPlanBudget *models.AuditPlanBudget

	result := db.
		Preload("ExpenseCategory").
		Preload("AuditPlan").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditPlanBudget, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanBudget, nil
}

func (repo *AuditPlanBudgetRepository) Create(auditPlanBudget *models.AuditPlanBudget, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanBudget).Error; err != nil {
		return 0, err
	}

	return auditPlanBudget.ID, nil
}

func (repo *AuditPlanBudgetRepository) Update(auditPlanBudget *models.AuditPlanBudget, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditPlanBudget).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanBudgetRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Plan Tool
	if err := db.Delete(&models.AuditPlanBudget{}, id).Error; err != nil {
		return err
	}

	return nil
}
