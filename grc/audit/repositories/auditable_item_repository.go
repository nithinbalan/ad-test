package repositories

import (
	"fmt"
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditableItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *auditDTOs.AuditableItemListQueryParamsDTO) (*auditDTOs.AuditableItemListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditableItem, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditableItem, error)
	Create(auditableItem *models.AuditableItem, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditableItem *models.AuditableItem, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditableItemRepository struct {
	db *gorm.DB
}

func NewAuditableItemRepository(db *gorm.DB) *AuditableItemRepository {
	return &AuditableItemRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditableItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditableItemRepository) GetAll(params *auditDTOs.AuditableItemListQueryParamsDTO) (*auditDTOs.AuditableItemListResponseDTO, error) {

	var auditableItems []*auditDTOs.AuditableItemListDTO
	var totalCount int64

	query := repo.db.Table("auditable_items").
		Select(`
			auditable_items.*,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			auditable_items.type AS auditable_item_type,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			incident_register.title AS incident_register_title,
			compliance_requirements.title_en AS compliance_requirement_title,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,   
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.email AS updated_user_email,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			audit_statuses.title_en AS audit_status_title_en,
			audit_statuses.title_ar AS audit_status_title_ar,
			audit_statuses.color_code AS audit_status_color_code,
			audit_statuses.type AS audit_status_type,
			auditor_risk_ratings.title_en AS auditor_risk_rating_title_en,
			auditor_risk_ratings.title_ar AS auditor_risk_rating_title_ar,
			auditor_risk_ratings.color_code AS auditor_risk_rating_color_code,
			auditor_risk_ratings.type AS auditor_risk_rating_type
		`).
		Joins(`
        LEFT JOIN processes ON processes.id = auditable_items.process_id
        LEFT JOIN controls ON controls.id = auditable_items.control_id
        LEFT JOIN documents ON documents.id = auditable_items.document_id
        LEFT JOIN incident_register ON incident_register.id = auditable_items.incident_register_id
        LEFT JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		LEFT JOIN organization_units ON organization_units.id = auditable_items.responsible_unit_id
        LEFT JOIN risk_register ON risk_register.id = auditable_items.risk_register_id
        LEFT JOIN users ON users.id = auditable_items.user_id
        LEFT JOIN users AS created_user ON created_user.id = auditable_items.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = auditable_items.updated_by_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = auditable_items.responsible_user_id
        LEFT JOIN audit_statuses ON audit_statuses.id = auditable_items.audit_status_id
        LEFT JOIN auditor_risk_ratings ON auditor_risk_ratings.id = auditable_items.auditor_risk_rating_id
    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s IN (?)", columnName), values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "draft":
			query = query.Where("audit_statuses.type = ?", "draft")
		case "completed":
			query = query.Where("audit_statuses.type = ?", "completed")
		case "in_review":
			query = query.Where("audit_statuses.type = ?", " in-review")
		case "in_progress":
			query = query.Where("audit_statuses.type = ?", "in-progress")
		case "published":
			query = query.Where("audit_statuses.type = ?", "published")
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
            (auditable_items.title_en LIKE ? OR
            auditable_items.title_ar LIKE ? OR
            auditor_risk_ratings.title_en LIKE ? OR
            auditor_risk_ratings.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
            responsible_user.first_name_ar LIKE ? OR
            users.first_name_en LIKE ? OR
            users.first_name_ar LIKE ? OR
            users.last_name_en LIKE ? OR
            users.last_name_ar LIKE ? OR
            users.email LIKE ? OR
            users.profile_image_path LIKE ? OR
            organization_units.title_en LIKE ? OR
            organization_units.title_ar LIKE ? OR
            audit_statuses.title_en LIKE ? OR
            audit_statuses.title_ar LIKE ? OR
            auditor_risk_ratings.title_en LIKE ? OR
            auditor_risk_ratings.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditableItems).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditableItemListResponseDTO{
		Items: auditableItems,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditableItemRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditableItem, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditableItem *models.AuditableItem
	result := db.First(&auditableItem, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditableItem, nil
}

func (repo *AuditableItemRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditableItem, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var auditableItem models.AuditableItem

	// Load the auditable item with all related data
	result := db.
		Preload("AuditStatus").
		Preload("AuditorRiskRating").
		Preload("Process").
		Preload("Control").
		Preload("Document").
		Preload("IncidentRegister").
		Preload("ComplianceRequirement").
		Preload("RiskRegister").
		Preload("User").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditableItem, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &auditableItem, nil
}

func (repo *AuditableItemRepository) Create(auditableItem *models.AuditableItem, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditableItem).Error; err != nil {
		return 0, err
	}
	return auditableItem.ID, nil
}

func (repo *AuditableItemRepository) Update(auditableItem *models.AuditableItem, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(auditableItem).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AuditableItemRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditableItem *models.AuditableItem
	if err := db.First(&auditableItem, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.AuditableItem{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
