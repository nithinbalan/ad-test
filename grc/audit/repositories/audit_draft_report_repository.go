package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditID uint, params *auditDTOs.AuditDraftReportListQueryParamsDTO) (*auditDTOs.AuditDraftReportListResponseDTO, error)
	GetByID(id uint) (*models.AuditDraftReport, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftReport, error)
	Create(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditDraftReportRepository struct {
	db *gorm.DB
}

func NewAuditDraftReportRepository(db *gorm.DB) *AuditDraftReportRepository {
	return &AuditDraftReportRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDraftReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDraftReportRepository) GetAll(auditID uint, params *auditDTOs.AuditDraftReportListQueryParamsDTO) (*auditDTOs.AuditDraftReportListResponseDTO, error) {
	var auditDraftReports []*auditDTOs.AuditDraftReportListDTO
	var totalCount int64

	query := repo.db.Table("audit_draft_reports").
		Select(`
            audit_draft_reports.id,
            audit_draft_reports.audit_id,
            CONVERT(VARCHAR, audit_draft_reports.report_date, 23) AS report_date,
            audit_draft_reports.executive_summary,
            audit_draft_reports.conclusion,
            audit_draft_reports.created_by_id AS created_user_id,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,
            audit_draft_reports.updated_by_id AS updated_user_id,
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
            LEFT JOIN users AS created_user ON created_user.id = audit_draft_reports.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = audit_draft_reports.updated_by_id
        `).Where("audit_draft_reports.audit_id = ?", auditID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
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
            audit_draft_reports.executive_summary LIKE ? OR
            c.conclusion LIKE ?`,
			keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_draft_reports.report_date" // Default sort by report_date if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditDraftReports).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditDraftReportListResponseDTO{
		Items: auditDraftReports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditDraftReportRepository) GetByID(id uint) (*models.AuditDraftReport, error) {
	var auditDraftReport *models.AuditDraftReport

	result := repo.db.First(&auditDraftReport, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReport, nil
}

func (repo *AuditDraftReportRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditDraftReport *models.AuditDraftReport

	result := db.
		Preload("Audit").
		Preload("Audit.AuditDraftFindings").
		Preload("Audit.AuditDraftFindings.AuditDraftFindingStatus").
		Preload("Audit.AuditDraftFindings.Documents").
		Preload("Audit.AuditPlan").
		Preload("Audit.AuditPlan.AuditCategory").
		Preload("Audit.AuditPlan.AuditStatus").
		Preload("Audit.AuditPlan.TeamLead").
		Preload("Audit.AuditPlan.Members").
		Preload("Audit.AuditPlan.Members.User").
		Preload("Audit.CreatedByUser").
		Preload("Audit.UpdatedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditDraftReport, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReport, nil
}

func (repo *AuditDraftReportRepository) Create(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditDraftReport).Error; err != nil {
		return 0, err
	}

	return auditDraftReport.ID, nil
}

func (repo *AuditDraftReportRepository) Update(auditDraftReport *models.AuditDraftReport, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditDraftReport).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Schedule
	if err := db.Delete(&models.AuditDraftReport{}, id).Error; err != nil {
		return err
	}

	return nil
}
