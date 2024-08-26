package repositories

import (
	auditDTOs "jethur/grc/audit/dtos" // Adjust import path as needed
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditFindingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error)
	GetByID(id uint) (*models.AuditFinding, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFinding, error)
	Create(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditFindingRepository struct {
	db *gorm.DB
}

func NewAuditFindingRepository(db *gorm.DB) *AuditFindingRepository {
	return &AuditFindingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditFindingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditFindingRepository) GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error) {
	var auditFindings []*auditDTOs.AuditFindingListDTO
	var totalCount int64

	query := repo.db.Table("audit_findings").
		Select(`
			audit_findings.*,
			CONVERT(VARCHAR, audit_findings.start_date, 23) AS start_date,
			CONVERT(VARCHAR, audit_findings.end_date, 23) AS end_date,
			audit_plans.title AS audit_plan_title,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			team_lead.id AS team_lead_id,
			team_lead.first_name_en AS team_lead_first_name_en,
			team_lead.first_name_ar AS team_lead_first_name_ar,
			team_lead.last_name_en AS team_lead_last_name_en,
			team_lead.last_name_ar AS team_lead_last_name_ar,
			team_lead.email AS team_lead_email,
			team_lead.initial AS team_lead_initial,
			team_lead.color AS team_lead_color,
			team_lead.profile_image_path AS team_lead_profile_image_path,
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
			LEFT JOIN audit_plans ON audit_plans.id = audit_findings.audit_plan_id
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN users AS created_user ON created_user.id = audit_findings.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_findings.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("audit_findings.start_date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("audit_findings.end_date <= ?", params.EndDate)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_findings.start_date LIKE ? OR
			audit_findings.end_date LIKE ? OR
			audit_plans.title LIKE ? OR
			audit_categories.title_en LIKE ? OR
			audit_categories.title_ar LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_findings.start_date" // Default sort by start_date if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditFindings).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditFindingListResponseDTO{
		Items: auditFindings,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditFindingRepository) GetByID(id uint) (*models.AuditFinding, error) {
	var auditFinding *models.AuditFinding

	result := repo.db.First(&auditFinding, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditFinding, nil
}

func (repo *AuditFindingRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditFinding, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditFinding *models.AuditFinding

	result := db.
		Preload("AuditPlan").
		Preload("AuditPlan.AuditCategory").
		Preload("AuditPlan.AuditTeam").
		Preload("AuditPlan.TeamLead").
		Preload("AuditPlan.Members").
		Preload("AuditPlan.Members.User").
		Preload("AuditPlan.AuditStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditFinding, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFinding, nil
}

func (repo *AuditFindingRepository) Create(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditFinding).Error; err != nil {
		return 0, err
	}

	return auditFinding.ID, nil
}

func (repo *AuditFindingRepository) Update(auditFinding *models.AuditFinding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditFinding).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditFindingRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the AuditFinding
	if err := db.Delete(&models.AuditFinding{}, id).Error; err != nil {
		return err
	}

	return nil
}
