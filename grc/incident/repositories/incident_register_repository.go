package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type IncidentRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error)
	GetByID(id uint) (*models.IncidentRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentRegister, error)
	Create(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	//Document
	CreateDocuments(docs []*models.IncidentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IncidentDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateIncidentSLAs(slaRegisters []*models.IncidentSLA, tx ...utility.DatabaseTransaction) error
	DeleteIncidentSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentRegisterRepository struct {
	db *gorm.DB
}

func NewIncidentRegisterRepository(db *gorm.DB) *IncidentRegisterRepository {
	return &IncidentRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *IncidentRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentRegisterRepository) GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error) {
	var incidentRegisters []*incidentDTOs.IncidentRegisterListDTO
	var totalCount int64

	query := repo.db.Table("incident_register").
		Select(`
            incident_register.*,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			reported_user.id AS reported_by_id,
            reported_user.first_name_en AS reported_user_first_name_en,
            reported_user.first_name_ar AS reported_user_first_name_ar,
            reported_user.last_name_en AS reported_user_last_name_en,
            reported_user.last_name_ar AS reported_user_last_name_ar,
			reported_user.initial AS reported_user_initial,
            reported_user.color AS reported_user_color,
            reported_user.email AS reported_user_email,
            reported_user.profile_image_path AS reported_user_profile_image_path,
			CONVERT(VARCHAR, incident_investigations.initiated_at, 23) AS investigation_initiated_at,
			incident_investigations.description AS investigation_description,
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
            LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
            LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
            LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
            LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
		    LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
            LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id
			LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
			LEFT JOIN incident_investigations ON incident_investigations.incident_register_id = incident_register.id
			LEFT JOIN users AS created_user ON created_user.id = incident_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = incident_register.updated_by_id
        `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "incident_category_id" {
				query = query.Where("incident_register.incident_category_id IN (?)", values)
			} else if columnName == "reported_user_id" {
				query = query.Where("incident_register.reported_by IN (?)", values)
			} else {
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	// Apply other filters as needed (e.g., date filters)
	if params.ReportedAt != "" {
		query = query.Where("reported_at >= ?", params.ReportedAt)
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("incident_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(incident_register.title LIKE ? OR
			incident_register.description LIKE ? OR
			incident_register.reference_number LIKE ? OR
			incident_categories.title_en LIKE ? OR
			incident_categories.title_ar LIKE ? OR
			incident_types.title_en LIKE ? OR
			incident_types.title_ar LIKE ? OR
			incident_statuses.title_en LIKE ? OR
			incident_statuses.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			reported_user.first_name_en LIKE ? OR
			reported_user.first_name_ar LIKE ? OR
			reported_user.last_name_en LIKE ? OR
			reported_user.last_name_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&incidentRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &incidentDTOs.IncidentRegisterListResponseDTO{
		Items: incidentRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *IncidentRegisterRepository) GetByID(id uint) (*models.IncidentRegister, error) {
	var incident *models.IncidentRegister

	result := repo.db.First(&incident, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil

}

func (repo *IncidentRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var incident *models.IncidentRegister

	result := db.
		Preload("IncidentCategory").
		Preload("IncidentType").
		Preload("IncidentStatus").
		Preload("ReportedByUser").
		Preload("SeverityLevel").
		Preload("Documents").
		Preload("ComplianceRequirements").
		Preload("SLARegisters").
		Preload("AssociatedDocuments").
		Preload("ResponsibleUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&incident, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil
}

func (repo *IncidentRegisterRepository) Create(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(incidentRegister).Error; err != nil {
		return 0, err
	}

	return incidentRegister.ID, nil
}

func (repo *IncidentRegisterRepository) Update(incidentRegister *models.IncidentRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Incident register
	if err := db.Delete(&models.IncidentRegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.IncidentRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "IR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

// Document func
func (repo *IncidentRegisterRepository) GetDocumentByID(id uint) (*models.IncidentDocument, error) {
	var document models.IncidentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *IncidentRegisterRepository) CreateDocuments(docs []*models.IncidentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.IncidentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) CreateIncidentSLAs(slaRegisters []*models.IncidentSLA, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&slaRegisters).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentRegisterRepository) DeleteIncidentSLAs(slaRegisterIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete SLARegister by IDs
	if err := db.Where("id IN (?)", slaRegisterIds).Delete(&models.IncidentSLA{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentRegisterRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	queries := []string{
		`DELETE audit_finding_action_plan_status_update_documents
		FROM audit_finding_action_plan_status_update_documents
		JOIN audit_finding_action_plan_status_updates ON audit_finding_action_plan_status_updates.id = audit_finding_action_plan_status_update_documents.audit_finding_action_plan_status_update_id
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_action_plan_comment_documents
		FROM audit_finding_action_plan_comment_documents
		JOIN audit_finding_action_plan_comments ON audit_finding_action_plan_comments.id = audit_finding_action_plan_comment_documents.audit_finding_action_plan_comment_id
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = tasks.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_action_plan_status_updates
		FROM audit_finding_action_plan_status_updates
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_action_plan_comments
		FROM audit_finding_action_plan_comments
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_draft_finding_conversation_documents
		FROM audit_draft_finding_conversation_documents
		JOIN audit_draft_finding_conversations ON audit_draft_finding_conversations.id = audit_draft_finding_conversation_documents.audit_draft_finding_conversation_id
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_strategic_objectives
		FROM audit_finding_strategic_objectives
		JOIN audit_findings ON audit_findings.id = audit_finding_strategic_objectives.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_root_cause_analyses
		FROM audit_finding_root_cause_analyses
		JOIN audit_findings ON audit_findings.id = audit_finding_root_cause_analyses.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_risks
		FROM audit_finding_risks
		JOIN audit_findings ON audit_findings.id = audit_finding_risks.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_processes
		FROM audit_finding_processes
		JOIN audit_findings ON audit_findings.id = audit_finding_processes.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_kpis
		FROM audit_finding_kpis
		JOIN audit_findings ON audit_findings.id = audit_finding_kpis.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_incidents
		FROM audit_finding_incidents
		JOIN audit_findings ON audit_findings.id = audit_finding_incidents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_documents
		FROM audit_finding_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_controls
		FROM audit_finding_controls
		JOIN audit_findings ON audit_findings.id = audit_finding_controls.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_compliance_requirements
		FROM audit_finding_compliance_requirements
		JOIN audit_findings ON audit_findings.id = audit_finding_compliance_requirements.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_associated_documents
		FROM audit_finding_associated_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_associated_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_assets
		FROM audit_finding_assets
		JOIN audit_findings ON audit_findings.id = audit_finding_assets.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_finding_action_plans
		FROM audit_finding_action_plans
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_draft_finding_documents
		FROM audit_draft_finding_documents
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_documents.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_draft_finding_conversations
		FROM audit_draft_finding_conversations
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_draft_finding_action_plans
		FROM audit_draft_finding_action_plans
		JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_action_plans.audit_draft_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_plan_auditable_item_risk_controls
		FROM audit_plan_auditable_item_risk_controls
		JOIN audit_plan_auditable_item_risks ON audit_plan_auditable_item_risks.id = audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risk_id
		JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_findings
		FROM audit_findings
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_draft_findings
		FROM audit_draft_findings
		JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_checklist_documents
		FROM audit_checklist_documents
		JOIN audit_checklists ON audit_checklists.id = audit_checklist_documents.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE incident_action_plan_status_update_documents
		FROM incident_action_plan_status_update_documents
		JOIN incident_action_plan_status_updates ON incident_action_plan_status_updates.id = incident_action_plan_status_update_documents.incident_action_plan_status_update_id
		JOIN incident_action_plans ON incident_action_plans.id = incident_action_plan_status_updates.incident_action_plan_id
		WHERE incident_action_plans.incident_register_id = ?`,

		`DELETE incident_action_plan_comment_documents
		FROM incident_action_plan_comment_documents
		JOIN incident_action_plan_comments ON incident_action_plan_comments.id = incident_action_plan_comment_documents.incident_action_plan_comment_id
		JOIN incident_action_plans ON incident_action_plans.id = incident_action_plan_comments.incident_action_plan_id
		WHERE incident_action_plans.incident_register_id = ?`,

		`DELETE audit_checklists
		FROM audit_checklists
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_plan_auditable_item_risks
		FROM audit_plan_auditable_item_risks
		JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE incident_post_review_users
		FROM incident_post_review_users
		JOIN incident_post_reviews ON incident_post_reviews.id = incident_post_review_users.incident_post_review_id
		WHERE incident_post_reviews.incident_register_id = ?`,

		`DELETE incident_post_review_recommendations
		FROM incident_post_review_recommendations
		JOIN incident_post_reviews ON incident_post_reviews.id = incident_post_review_recommendations.incident_post_review_id
		WHERE incident_post_reviews.incident_register_id = ?`,

		`DELETE incident_post_review_lessons
		FROM incident_post_review_lessons
		JOIN incident_post_reviews ON incident_post_reviews.id = incident_post_review_lessons.incident_post_review_id
		WHERE incident_post_reviews.incident_register_id = ?`,

		`DELETE incident_investigation_users
		FROM incident_investigation_users
		JOIN incident_investigations ON incident_investigations.id = incident_investigation_users.incident_investigation_id
		WHERE incident_investigations.incident_register_id = ?`,

		`DELETE incident_investigation_findings
		FROM incident_investigation_findings
		JOIN incident_investigations ON incident_investigations.id = incident_investigation_findings.incident_investigation_id
		WHERE incident_investigations.incident_register_id = ?`,

		`DELETE incident_investigation_documents
		FROM incident_investigation_documents
		JOIN incident_investigations ON incident_investigations.id = incident_investigation_documents.incident_investigation_id
		WHERE incident_investigations.incident_register_id = ?`,

		`DELETE incident_closure_lessons
		FROM incident_closure_lessons
		JOIN incident_closures ON incident_closures.id = incident_closure_lessons.incident_closure_id
		WHERE incident_closures.incident_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN incident_action_plans ON incident_action_plans.id = tasks.incident_action_plan_id
		WHERE incident_action_plans.incident_register_id = ?`,

		`DELETE incident_action_plan_status_updates
		FROM incident_action_plan_status_updates
		JOIN incident_action_plans ON incident_action_plans.id = incident_action_plan_status_updates.incident_action_plan_id
		WHERE incident_action_plans.incident_register_id = ?`,

		`DELETE incident_action_plan_comments
		FROM incident_action_plan_comments
		JOIN incident_action_plans ON incident_action_plans.id = incident_action_plan_comments.incident_action_plan_id
		WHERE incident_action_plans.incident_register_id = ?`,

		`DELETE audit_schedule_auditable_items
		FROM audit_schedule_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_schedule_auditable_items.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_plan_schedule_auditable_items
		FROM audit_plan_schedule_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_plan_schedule_auditable_items.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_plan_checklists
		FROM audit_plan_checklists
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_plan_auditable_items
		FROM audit_plan_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE audit_document_request_auditable_items
		FROM audit_document_request_auditable_items
		JOIN auditable_items ON auditable_items.id = audit_document_request_auditable_items.auditable_item_id
		WHERE auditable_items.incident_register_id = ?`,

		`DELETE FROM third_party_incidents WHERE third_party_incidents.incident_register_id = ?`,
		`DELETE FROM project_incidents WHERE project_incidents.incident_register_id = ?`,
		`DELETE FROM meeting_incidents WHERE meeting_incidents.incident_register_id = ?`,
		`DELETE FROM kri_incidents WHERE kri_incidents.incident_register_id = ?`,
		`DELETE FROM incident_timelines WHERE incident_timelines.incident_register_id = ?`,
		`DELETE FROM incident_strategic_objectives WHERE incident_strategic_objectives.incident_register_id = ?`,
		`DELETE FROM incident_slas WHERE incident_slas.incident_register_id = ?`,
		`DELETE FROM incident_root_cause_analyses WHERE incident_root_cause_analyses.incident_register_id = ?`,
		`DELETE FROM incident_risks WHERE incident_risks.incident_register_id = ?`,
		`DELETE FROM incident_processes WHERE incident_processes.incident_register_id = ?`,
		`DELETE FROM incident_post_reviews WHERE incident_post_reviews.incident_register_id = ?`,
		`DELETE FROM incident_kpis WHERE incident_kpis.incident_register_id = ?`,
		`DELETE FROM incident_investigations WHERE incident_investigations.incident_register_id = ?`,
		`DELETE FROM incident_documents WHERE incident_documents.incident_register_id = ?`,
		`DELETE FROM incident_controls WHERE incident_controls.incident_register_id = ?`,
		`DELETE FROM incident_communication_logs WHERE incident_communication_logs.incident_register_id = ?`,
		`DELETE FROM incident_closures WHERE incident_closures.incident_register_id = ?`,
		`DELETE FROM incident_associated_documents WHERE incident_associated_documents.incident_register_id = ?`,
		`DELETE FROM incident_compliance_requirements WHERE incident_compliance_requirements.incident_register_id = ?`,
		`DELETE FROM incident_assets WHERE incident_assets.incident_register_id = ?`,
		`DELETE FROM incident_action_plans WHERE incident_action_plans.incident_register_id = ?`,
		`DELETE FROM bcp_incidents WHERE bcp_incidents.incident_register_id = ?`,
		`DELETE FROM auditable_items WHERE auditable_items.incident_register_id = ?`,
		`DELETE FROM audit_finding_incidents WHERE audit_finding_incidents.incident_register_id = ?`,
	}

	// Delete the Risk register

	return db.Transaction(func(tx *gorm.DB) error {
		for _, q := range queries {
			if err := tx.Exec(q, id).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&models.IncidentRegister{}, id).Error
	})

}
