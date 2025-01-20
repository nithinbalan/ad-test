package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditListQueryParamsDTO) (*auditDTOs.AuditListResponseDTO, error)
	GetByID(id uint) (*models.Audit, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Audit, error)
	Create(audit *models.Audit, tx ...utility.DatabaseTransaction) (uint, error)
	Update(audit *models.Audit, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateAuditSchedule(auditSchedule *models.AuditSchedule, tx ...utility.DatabaseTransaction) (*models.AuditSchedule, error)
	CreateAuditScheduleStatusUpdate(statusUpdate *models.AuditScheduleStatusUpdate, tx ...utility.DatabaseTransaction) error
	CreateAuditScheduleAuditableItem(item *models.AuditScheduleAuditableItem, tx ...utility.DatabaseTransaction) error
	CreateAuditChecklists(auditChecklists []*models.AuditChecklist, tx ...utility.DatabaseTransaction) error

	GetAllAuditableItems(auditID uint, params *auditDTOs.AuditAuditableItemListQueryParamsDTO) (*auditDTOs.AuditAuditableItemListResponseDTO, error)
	GetAuditByAuditPlanID(auditPlanID uint) (*models.Audit, error)
}

type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditRepository) GetAll(params *auditDTOs.AuditListQueryParamsDTO) (*auditDTOs.AuditListResponseDTO, error) {
	var audits []*auditDTOs.AuditListDTO
	var totalCount int64

	query := repo.db.Table("audits").
		Select(`
			audits.*,
			CONVERT(VARCHAR, audits.start_date, 23) AS start_date,
			CONVERT(VARCHAR, audits.end_date, 23) AS end_date,
			audit_plans.title AS audit_plan_title,
			audit_plans.description AS audit_plan_description,
			audit_plans.purpose AS audit_plan_purpose,
			audit_plans.objective AS audit_plan_objective,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.type AS audit_status_type,
			audit_statuses.color_code AS audit_status_color_code,
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
			LEFT JOIN audit_statuses ON audit_statuses.id = audits.audit_status_id
			LEFT JOIN audit_plans ON audit_plans.id = audits.audit_plan_id
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN users AS created_user ON created_user.id = audits.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audits.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			//audits and audits plans using same fields as audit_status_id, so attaching table name in to the field
			if columnName == "audit_status_id" {
				// Prepend the table name for "audit_status_id"
				query = query.Where("audits.audit_status_id IN (?)", values)
			} else {
				// Use the column name as-is for other filters
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.StartDate != "" {
		query = query.Where("audits.start_date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("audits.end_date <= ?", params.EndDate)
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "draft":
			query = query.Where("audit_statuses.type = ?", "draft")
		case "published":
			query = query.Where("audit_statuses.type = ?", "published")
		case "completed":
			query = query.Where("audit_statuses.type = ?", "completed")
		case "in_review":
			query = query.Where("audit_statuses.type = ?", "in-review")
		case "in_progress":
			query = query.Where("audit_statuses.type = ?", "in-progress")
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
			(audits.start_date LIKE ? OR
			audits.end_date LIKE ? OR
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

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&audits).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}
	response := &auditDTOs.AuditListResponseDTO{
		Items: audits,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditRepository) GetByID(id uint) (*models.Audit, error) {
	var audit *models.Audit

	result := repo.db.First(&audit, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return audit, nil
}

func (repo *AuditRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Audit, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var audit *models.Audit

	result := db.
		Preload("AuditPlan").
		Preload("AuditPlan.AuditCategory").
		Preload("AuditPlan.AuditTeam").
		Preload("AuditPlan.TeamLead").
		Preload("AuditPlan.Members").
		Preload("AuditPlan.Members.User").
		Preload("AuditPlan.AuditStatus").
		Preload("AuditPlan.AuditableItems.AuditableItem.ResponsibleUnit").
		Preload("AuditPlan.AuditableItems.AuditableItem.Process").
		Preload("AuditPlan.AuditableItems.AuditableItem.Control").
		Preload("AuditPlan.AuditableItems.AuditableItem.RiskRegister").
		Preload("AuditPlan.AuditableItems.AuditableItem.IncidentRegister").
		Preload("AuditPlan.AuditableItems.AuditableItem.ComplianceRequirement").
		Preload("AuditPlan.AuditableItems.AuditableItem.User").
		Preload("AuditPlan.AuditableItems.AuditableItem.Document").
		Preload("AuditPlan.AuditableItems.AuditableItem.AuditorRiskRating").
		Preload("AuditPlan.AuditableItems.AuditableItem.AuditStatus").
		Preload("AuditChecklists").
		Preload("AuditDraftFindings").
		Preload("AuditStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&audit, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return audit, nil
}

func (repo *AuditRepository) Create(audit *models.Audit, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(audit).Error; err != nil {
		return 0, err
	}

	return audit.ID, nil
}

func (repo *AuditRepository) Update(audit *models.Audit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(audit).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit
	if err := db.Delete(&models.Audit{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditRepository) CreateAuditSchedule(auditSchedule *models.AuditSchedule, tx ...utility.DatabaseTransaction) (*models.AuditSchedule, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(auditSchedule).Error; err != nil {
		return nil, err
	}

	return auditSchedule, nil
}

func (repo *AuditRepository) CreateAuditScheduleStatusUpdate(statusUpdate *models.AuditScheduleStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 {
		db = tx[0].DB()
	}
	if err := db.Create(statusUpdate).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditRepository) CreateAuditScheduleAuditableItem(item *models.AuditScheduleAuditableItem, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(item).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditRepository) CreateAuditChecklists(auditChecklists []*models.AuditChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&auditChecklists).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditRepository) GetAllAuditableItems(auditID uint, params *auditDTOs.AuditAuditableItemListQueryParamsDTO) (*auditDTOs.AuditAuditableItemListResponseDTO, error) {
	var auditAuditableItems []*auditDTOs.AuditAuditableItemListDTO
	var totalCount int64

	query := repo.db.Table("audit_plan_auditable_items").
		Select(`
			audit_plan_auditable_items.id,
			audit_plan_auditable_items.audit_plan_id,
			audit_plan_auditable_items.auditable_item_id,
			audit_plan_auditable_items.created_by_id,
			audit_plan_auditable_items.created_at,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			incident_register.title AS incident_register_title,
			compliance_requirements.title AS compliance_requirement_title,
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			risk_register.title AS risk_register_title,
			users.id AS user_id,
			users.first_name_en AS user_first_name_en,
			users.first_name_ar AS user_first_name_ar,
			users.last_name_en AS user_last_name_en,
			users.last_name_ar AS user_last_name_ar,
			users.email AS user_email,
			users.initial AS user_initial,
			users.color AS user_color,
			users.profile_image_path AS user_profile_image_path,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			audit_plans.reference_number,
			audit_plans.title AS audit_plan_title,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			auditable_items.type AS auditable_item_type,
			auditable_items.comment AS auditable_item_comment,
			audit_categories.title_en AS audit_category_title_en,
			audit_categories.title_ar AS audit_category_title_ar,
			audit_teams.title_en AS audit_team_title_en,
			audit_teams.title_ar AS audit_team_title_ar,
			team_lead.id AS team_lead_id,
			team_lead.first_name_en AS team_lead_first_name_en,
			team_lead.first_name_ar AS team_lead_first_name_ar,
			team_lead.last_name_en AS team_lead_last_name_en,
			team_lead.last_name_ar AS team_lead_last_name_ar,
			team_lead.email AS team_lead_email,
			team_lead.profile_image_path AS team_lead_profile_image_path,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.type AS audit_status_type,
			audit_statuses.color_code AS audit_status_color_code,
			auditor_risk_ratings.title_en AS auditor_risk_rating_title_en,
			auditor_risk_ratings.title_ar AS auditor_risk_rating_title_ar,
			auditor_risk_ratings.color_code AS auditor_risk_rating_color_code,
			auditor_risk_ratings.type AS auditor_risk_rating_type,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN audit_plans ON audit_plans.id = audit_plan_auditable_items.audit_plan_id
			LEFT JOIN audits ON audit_plans.id = audits.audit_plan_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
			LEFT JOIN processes ON processes.id = auditable_items.process_id
			LEFT JOIN controls ON controls.id = auditable_items.control_id
			LEFT JOIN documents ON documents.id = auditable_items.document_id
			LEFT JOIN incident_register ON incident_register.id = auditable_items.incident_register_id
			LEFT JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
			LEFT JOIN organization_units ON organization_units.id = auditable_items.responsible_unit_id
			LEFT JOIN risk_register ON risk_register.id = auditable_items.risk_register_id
			LEFT JOIN users ON users.id = auditable_items.user_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = auditable_items.responsible_user_id
			LEFT JOIN audit_categories ON audit_categories.id = audit_plans.audit_category_id
			LEFT JOIN audit_teams ON audit_teams.id = audit_plans.audit_team_id
			LEFT JOIN users AS team_lead ON team_lead.id = audit_plans.team_lead_id
			LEFT JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id
        	LEFT JOIN auditor_risk_ratings ON auditor_risk_ratings.id = auditable_items.auditor_risk_rating_id
			LEFT JOIN users AS created_user ON created_user.id = audit_plan_auditable_items.created_by_id
		`).Where("audits.id = ?", auditID)

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
			(audit_plans.title LIKE ? OR
			audit_plans.description LIKE ? OR
			audit_plans.reference_number LIKE ? OR
			audit_categories.title_en LIKE ? OR
			audit_categories.title_ar LIKE ? OR
			auditable_items.title_en LIKE ? OR
			auditable_items.title_ar LIKE ? OR
			audit_teams.title_en LIKE ? OR
			audit_teams.title_ar LIKE ? OR
			audit_statuses.title_en LIKE ? OR
			audit_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			team_lead.first_name_en LIKE ? OR
			team_lead.first_name_ar LIKE ? OR
			team_lead.last_name_en LIKE ? OR
			team_lead.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "audit_plans.title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditAuditableItems).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditAuditableItemListResponseDTO{
		Items: auditAuditableItems,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditRepository) GetAuditByAuditPlanID(auditPlanID uint) (*models.Audit, error) {
	var audit models.Audit

	// Query the database to find if an audit exists with the given AuditPlanID
	result := repo.db.Where("audit_plan_id = ?", auditPlanID).First(&audit)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return nil if no audit is found (not an error, meaning it's safe to create a new audit)
			return nil, nil
		}
		// Return an error if something went wrong during the query
		return nil, result.Error
	}

	// Return the found audit if it exists
	return &audit, nil
}
