package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditChecklistRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditId uint, params *auditDTOs.AuditChecklistListQueryParamsDTO) ([]*auditDTOs.AuditChecklistListDTO, error)
	GetByID(id uint) (*models.AuditChecklist, error)
	GetFullDetailsByID(id uint) (*models.AuditChecklist, error)
	GetDraftFindingsDetailsByChecklistID(checklistID uint) ([]*models.AuditDraftFinding, error)
	Create(auditChecklist *models.AuditChecklist, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditChecklist *models.AuditChecklist, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	// Document methods
	CreateDocuments(docs []*models.AuditChecklistDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditChecklistDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditChecklistRepository struct {
	db *gorm.DB
}

func NewAuditChecklistRepository(db *gorm.DB) *AuditChecklistRepository {
	return &AuditChecklistRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditChecklistRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditChecklistRepository) GetAll(auditId uint, params *auditDTOs.AuditChecklistListQueryParamsDTO) ([]*auditDTOs.AuditChecklistListDTO, error) {
	var auditChecklists []*auditDTOs.AuditChecklistListDTO

	query := repo.db.Table("audit_checklists").
		Select(`
			audit_checklists.id,
			audit_checklists.audit_id,
			audit_checklists.audit_plan_checklist_id,
			audit_checklists.question,
			audit_checklists.conformance_id,
			audit_checklists.auditable_item_id,
			auditable_items.title_en AS auditable_item_title_en,
			auditable_items.title_ar AS auditable_item_title_ar,
			auditable_items.type AS auditable_item_type,
			auditable_items.comment AS auditable_item_comment,
			auditable_items.process_id AS auditable_item_process_id,
			auditable_items.control_id AS auditable_item_control_id,
			auditable_items.document_id AS auditable_item_document_id,
			auditable_items.incident_register_id AS auditable_item_incident_register_id,
			auditable_items.compliance_requirement_id AS auditable_item_compliance_requirement_id,
			auditable_items.risk_register_id AS auditable_item_risk_register_id,
			auditable_items.user_id AS auditable_item_user_id,
			conformances.title_en AS conformance_title_en,
			conformances.title_ar AS conformance_title_ar,
			conformances.type AS conformance_type,
			conformances.color_code AS conformance_color_code,
			audit_checklists.comment,
			audit_checklists.answered_by_id,
			audit_checklists.answered_at,
			audit_checklists.created_by_id,
			audit_checklists.created_at,
			audit_checklists.updated_by_id,
			audit_checklists.updated_at,
			statuses.id AS status_id,
			statuses.title_en AS status_title_en,
			statuses.title_ar AS status_title_ar,
			statuses.type AS status_type,
			statuses.color_code AS status_color_code,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			controls.title_en AS control_title_en,
			controls.title_ar AS control_title_ar,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			incident_register.title AS incident_register_title,
			compliance_requirements.title AS compliance_requirement_title,
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
			answered_user.first_name_en AS answered_user_first_name_en,
			answered_user.first_name_ar AS answered_user_first_name_ar,
			answered_user.last_name_en AS answered_user_last_name_en,
			answered_user.last_name_ar AS answered_user_last_name_ar,
			answered_user.initial AS answered_user_initial,
			answered_user.color AS answered_user_color,
			answered_user.email AS answered_user_email,
			answered_user.profile_image_path AS answered_user_profile_image_path,
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
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_checklists.auditable_item_id
			LEFT JOIN processes ON processes.id = auditable_items.process_id
			LEFT JOIN controls ON controls.id = auditable_items.control_id
			LEFT JOIN documents ON documents.id = auditable_items.document_id
			LEFT JOIN incident_register ON incident_register.id = auditable_items.incident_register_id
			LEFT JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
			LEFT JOIN risk_register ON risk_register.id = auditable_items.risk_register_id
			LEFT JOIN users ON users.id = auditable_items.user_id
			LEFT JOIN statuses ON statuses.id = audit_checklists.status_id
			LEFT JOIN conformances ON conformances.id = audit_checklists.conformance_id
			LEFT JOIN users AS answered_user ON answered_user.id = audit_checklists.answered_by_id
			LEFT JOIN users AS created_user ON created_user.id = audit_checklists.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_checklists.updated_by_id
		`).Where("audit_checklists.audit_id = ?", auditId)

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "process":
			query = query.Where("auditable_items.process_id IS NOT NULL")
		case "control":
			query = query.Where("auditable_items.control_id IS NOT NULL")
		case "document":
			query = query.Where("auditable_items.document_id IS NOT NULL")
		case "risk":
			query = query.Where("auditable_items.risk_register_id IS NOT NULL")
		case "incident":
			query = query.Where("auditable_items.incident_register_id IS NOT NULL")
		case "compliance":
			query = query.Where("auditable_items.compliance_requirement_id IS NOT NULL")
		}
	}

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_checklists.question LIKE ? OR
			auditable_items.title_en LIKE ? OR
			auditable_items.title_ar LIKE ?)`,
			keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by created_at if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&auditChecklists).Error
	if err != nil {
		return nil, err
	}
	// Fetching related documents for each conversation
	for _, auditChecklist := range auditChecklists {
		var documents []*auditDTOs.AuditChecklistDocumentDTO

		repo.db.Model(&models.AuditChecklistDocument{}).
			Where("audit_checklist_id = ?", auditChecklist.ID).
			Find(&documents)

			auditChecklist.Documents = documents
	}
	return auditChecklists, nil
}

func (repo *AuditChecklistRepository) GetByID(id uint) (*models.AuditChecklist, error) {
	var auditChecklist *models.AuditChecklist

	result := repo.db.First(&auditChecklist, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditChecklist, nil
}

func (repo *AuditChecklistRepository) GetFullDetailsByID(id uint) (*models.AuditChecklist, error) {
	var checklist *models.AuditChecklist

	result := repo.db.
		Preload("Audit").
		Preload("AuditPlanChecklist").
		Preload("AuditPlanChecklist.AuditPlan").
		Preload("AuditableItem").
		Preload("Conformance").
		Preload("AuditDraftFindings").
		Preload("AuditDraftFindings.AuditableItem").
		Preload("AuditDraftFindings.AuditFindingCategory").
		Preload("AuditDraftFindings.SeverityLevel").
		Preload("AuditDraftFindings.AuditDraftFindingStatus").
		Preload("AuditDraftFindings.AuditChecklist").
		Preload("AuditDraftFindings.CreatedByUser").
		Preload("AuditDraftFindings.UpdatedByUser").
		Preload("AuditDraftFindings.Documents").
		Preload("Documents").
		Preload("AnsweredByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Status").
		First(&checklist, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return checklist, nil
}

func (repo *AuditChecklistRepository) GetDraftFindingsDetailsByChecklistID(checklistID uint) ([]*models.AuditDraftFinding, error) {
	var auditFindings []*models.AuditDraftFinding

	result := repo.db.
		Preload("AuditableItem").
		Preload("AuditableItem.ResponsibleUnit").
		Preload("AuditableItem.Process").
		Preload("AuditableItem.Control").
		Preload("AuditableItem.Document").
		Preload("AuditableItem.IncidentRegister").
		Preload("AuditableItem.ComplianceRequirement").
		Preload("AuditableItem.RiskRegister").
		Preload("AuditableItem.User").
		Preload("AuditFindingCategory").
		Preload("SeverityLevel").
		Preload("AuditDraftFindingStatus").
		Preload("AuditChecklist").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("audit_checklist_id = ?", checklistID).
		Find(&auditFindings) // Note: Find to get all matching records

	if result.Error != nil {
		return nil, result.Error
	}

	return auditFindings, nil
}

func (repo *AuditChecklistRepository) Create(auditChecklist *models.AuditChecklist, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditChecklist).Error; err != nil {
		return 0, err
	}

	return auditChecklist.ID, nil
}

func (repo *AuditChecklistRepository) Update(auditChecklist *models.AuditChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditChecklist).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditChecklistRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.AuditChecklist{}, id).Error; err != nil {
		return err
	}

	return nil
}


func (repo *AuditChecklistRepository) CreateDocuments(docs []*models.AuditChecklistDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *AuditChecklistRepository) GetDocumentByID(id uint) (*models.AuditChecklistDocument, error) {
	var document models.AuditChecklistDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AuditChecklistRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.AuditChecklistDocument{}).Error; err != nil {
		return err
	}

	return nil
}
