package repositories

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftFindingConversationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditDraftFindingID uint, params *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error)
	GetByID(id uint) (*models.AuditDraftFindingConversation, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftFindingConversation, error)
	Create(auditDraftFindingConversation *models.AuditDraftFindingConversation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditDraftFindingConversation *models.AuditDraftFindingConversation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(docs []*models.AuditDraftFindingConversationDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditDraftFindingConversationDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

}

type AuditDraftFindingConversationRepository struct {
	db *gorm.DB
}

func NewAuditDraftFindingConversationRepository(db *gorm.DB) *AuditDraftFindingConversationRepository {
	return &AuditDraftFindingConversationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDraftFindingConversationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDraftFindingConversationRepository) GetAll(auditDraftFindingID uint, params *dtos.AuditDraftFindingConversationListQueryParamsDTO) ([]*dtos.AuditDraftFindingConversationListDTO, error) {

	var conversations []*dtos.AuditDraftFindingConversationListDTO

	// Initial query setup with joins
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
			LEFT JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
            LEFT JOIN users AS created_user ON created_user.id = audit_draft_finding_conversations.created_by_id
        `).
		Where("audit_draft_finding_conversations.audit_draft_finding_id = ?", auditDraftFindingID)

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
		params.SortBy = "description" // Default sort by title if no sortBy is specified
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

func (repo *AuditDraftFindingConversationRepository) GetByID(id uint) (*models.AuditDraftFindingConversation, error) {
	var conversation *models.AuditDraftFindingConversation
	result := repo.db.First(&conversation, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return conversation, nil
}

func (repo *AuditDraftFindingConversationRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftFindingConversation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditDraftFindingConversation *models.AuditDraftFindingConversation

	result := db.
		Preload("AuditDraftFinding").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditDraftFindingConversation, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftFindingConversation, nil
}


func (repo *AuditDraftFindingConversationRepository) Create(conversation *models.AuditDraftFindingConversation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(conversation).Error; err != nil {
		return 0, err
	}

	return conversation.ID, nil
}


func (repo *AuditDraftFindingConversationRepository) Update(conversation *models.AuditDraftFindingConversation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(conversation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftFindingConversationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Plan
	if err := db.Delete(&models.AuditDraftFindingConversation{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftFindingConversationRepository) CreateDocuments(docs []*models.AuditDraftFindingConversationDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *AuditDraftFindingConversationRepository) GetDocumentByID(id uint) (*models.AuditDraftFindingConversationDocument, error) {
	var document models.AuditDraftFindingConversationDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AuditDraftFindingConversationRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.AuditDraftFindingConversationDocument{}).Error; err != nil {
		return err
	}

	return nil
}
