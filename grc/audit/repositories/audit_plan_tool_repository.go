package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditPlanToolRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditPlanID uint, params *auditDTOs.AuditPlanToolListQueryParamsDTO) (*auditDTOs.AuditPlanToolListResponseDTO, error)
	GetByID(id uint) (*models.AuditPlanTool, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanTool, error)
	Create(auditPlanTool *models.AuditPlanTool, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditPlanTool *models.AuditPlanTool, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditPlanToolRepository struct {
	db *gorm.DB
}

func NewAuditPlanToolRepository(db *gorm.DB) *AuditPlanToolRepository {
	return &AuditPlanToolRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditPlanToolRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditPlanToolRepository) GetAll(auditPlanID uint, params *auditDTOs.AuditPlanToolListQueryParamsDTO) (*auditDTOs.AuditPlanToolListResponseDTO, error) {
	var auditPlanTools []*auditDTOs.AuditPlanToolListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_tools").
		Select(`
			audit_plan_tools.id,
			audit_plan_tools.audit_plan_id,
			audit_plan_tools.tool_type_id,
			tool_types.title_en AS tool_type_title_en,
			tool_types.title_ar AS tool_type_title_ar,
			audit_plan_tools.title,
			audit_plan_tools.purpose,
			audit_plans.title AS audit_plan_title,
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
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, audit_plan_tools.created_at, 23) AS created_at,
			CONVERT(VARCHAR, audit_plan_tools.updated_at, 23) AS updated_at
		`).
		Joins(`
			LEFT JOIN tool_types ON tool_types.id = audit_plan_tools.tool_type_id
			LEFT JOIN audit_plans ON audit_plans.id = audit_plan_tools.audit_plan_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_tools.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_plan_tools.updated_by_id
		`).Where(`audit_plan_tools.audit_plan_id`, auditPlanID)

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
			query = query.Where("audit_plan_tools.status = ?", "active")
		case "inactive":
			query = query.Where("audit_plan_tools.status = ?", "inactive")
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
			audit_plan_tools.title LIKE ? OR
			audit_plan_tools.purpose LIKE ? OR
			audit_plans.title LIKE ? OR
			tool_types.title_en LIKE ? OR
			tool_types.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditPlanTools).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditPlanToolListResponseDTO{
		Items: auditPlanTools,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditPlanToolRepository) GetByID(id uint) (*models.AuditPlanTool, error) {
	var auditPlanTool *models.AuditPlanTool

	result := repo.db.First(&auditPlanTool, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanTool, nil
}

func (repo *AuditPlanToolRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditPlanTool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditPlanTool *models.AuditPlanTool

	result := db.
		Preload("ToolType").
		Preload("AuditPlan").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditPlanTool, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditPlanTool, nil
}

func (repo *AuditPlanToolRepository) Create(auditPlanTool *models.AuditPlanTool, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditPlanTool).Error; err != nil {
		return 0, err
	}

	return auditPlanTool.ID, nil
}

func (repo *AuditPlanToolRepository) Update(auditPlanTool *models.AuditPlanTool, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditPlanTool).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditPlanToolRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Plan Tool
	if err := db.Delete(&models.AuditPlanTool{}, id).Error; err != nil {
		return err
	}

	return nil
}
