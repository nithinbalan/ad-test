package repositories

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AuditInsightsDraftFindingConversationRepositoryInterface interface {
	GetAll(currentUserResponsibleUnitID uint, auditDraftFindingID uint, params *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error)
}

type AuditInsightsDraftFindingConversationRepository struct {
	db *gorm.DB
}

func NewAuditInsightsDraftFindingConversationRepository(db *gorm.DB) *AuditInsightsDraftFindingConversationRepository {
	return &AuditInsightsDraftFindingConversationRepository{
		db: db,
	}
}

func (repo *AuditInsightsDraftFindingConversationRepository) GetAll(currentUserResponsibleUnitID uint, auditDraftFindingID uint, params *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error) {
	var conversations []*dtos.AuditDraftFindingConversationListDTO

	// Query starting from auditable items filtered by responsible unit ID
	query := repo.db.Table("audit_draft_finding_conversations").
		Select(`
			audit_draft_finding_conversations.id,
			audit_draft_finding_conversations.audit_draft_finding_id,
			audit_draft_finding_conversations.comment,
			audit_draft_finding_conversations.is_accepted,
            audit_draft_finding_conversations.is_dismissed,
            audit_draft_finding_conversations.is_reopened,
			audit_draft_finding_conversations.created_at,
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
			INNER JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
			LEFT JOIN auditable_items ON auditable_items.id = audit_draft_findings.auditable_item_id
			LEFT JOIN users AS created_user ON created_user.id = audit_draft_finding_conversations.created_by_id
		`).
		Where("auditable_items.responsible_unit_id = ? AND audit_draft_finding_conversations.audit_draft_finding_id = ?", currentUserResponsibleUnitID, auditDraftFindingID).
		Group(`
			audit_draft_finding_conversations.id,
			audit_draft_finding_conversations.audit_draft_finding_id,
			audit_draft_finding_conversations.comment,
			audit_draft_finding_conversations.is_accepted,
            audit_draft_finding_conversations.is_dismissed,
            audit_draft_finding_conversations.is_reopened,
			audit_draft_finding_conversations.created_at,
			created_user.id,
			created_user.first_name_en,
			created_user.first_name_ar,
			created_user.last_name_en,
			created_user.last_name_ar,
			created_user.initial,
			created_user.color,
			created_user.email,
			created_user.profile_image_path
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Apply other filters as needed (e.g., date filters)
	if params.CreatedAt != "" {
		query = query.Where("created_at = ?", params.CreatedAt)
	}
	// Keyword filtering
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			audit_draft_finding_conversations.comment LIKE ? OR
			created_user.first_name_en LIKE ? OR
            created_user.first_name_ar LIKE ? OR
            created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ?`,
			keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "comment" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&conversations)

	// Fetching related documents for each conversations
	for _, conversation := range conversations {
		var documents []*dtos.AuditDraftFindingConversationDocumentDTO

		repo.db.Model(&models.AuditDraftFindingConversationDocument{}).
			Where("audit_draft_finding_conversation_id = ?", conversation.ID).
			Find(&documents)

		conversation.Documents = documents
	}

	return conversations, result.Error
}
