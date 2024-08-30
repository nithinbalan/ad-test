package repositories

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDocumentRequestConversationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditDocumentRequestID uint, params *dtos.AuditDocumentRequestConversationListQueryParamsDTO) ([]*dtos.AuditDocumentRequestConversationListDTO, error)
	GetByID(id uint) (*models.AuditDocumentRequestConversation, error)
	Create(AuditDocumentRequestConversation *models.AuditDocumentRequestConversation, tx ...utility.DatabaseTransaction) (uint, error)

	// Document-related methods
	CreateDocuments(docs []*models.AuditDocumentRequestConversationDocument, tx ...utility.DatabaseTransaction) error
	AuditDocumentRequestStatusUpdate(auditDocumentRequest *models.AuditDocumentRequest, tx ...utility.DatabaseTransaction) error

}

type AuditDocumentRequestConversationRepository struct {
	db *gorm.DB
}

func NewAuditDocumentRequestConversationRepository(db *gorm.DB) *AuditDocumentRequestConversationRepository {
	return &AuditDocumentRequestConversationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDocumentRequestConversationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDocumentRequestConversationRepository) GetAll(auditDocumentRequestID uint, params *dtos.AuditDocumentRequestConversationListQueryParamsDTO) ([]*dtos.AuditDocumentRequestConversationListDTO, error) {

	var conversations []*dtos.AuditDocumentRequestConversationListDTO

	// Initial query setup with joins
	query := repo.db.Table("audit_document_request_conversations").
		Select(`
            audit_document_request_conversations.id,
            audit_document_request_conversations.audit_document_request_id,
            audit_document_request_conversations.comment,
            audit_document_request_conversations.commented_at,
			audit_document_requests.requested_by_id,
            commented_user.id AS commented_user_id,
            commented_user.first_name_en AS commented_user_first_name_en,
            commented_user.first_name_ar AS commented_user_first_name_ar,
            commented_user.last_name_en AS commented_user_last_name_en,
            commented_user.last_name_ar AS commented_user_last_name_ar,
            commented_user.initial AS commented_user_initial,
            commented_user.color AS commented_user_color,
            commented_user.email AS commented_user_email,
            commented_user.profile_image_path AS commented_user_profile_image_path
        `).
		Joins(`
			LEFT JOIN audit_document_requests ON audit_document_requests.id = audit_document_request_conversations.audit_document_request_id
            LEFT JOIN users AS commented_user ON commented_user.id = audit_document_request_conversations.commented_by_id
        `).
		Where("audit_document_request_conversations.audit_document_request_id = ?", auditDocumentRequestID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Apply other filters as needed (e.g., date filters)
	if params.CommentedAt != "" {
		query = query.Where("commented_at = ?", params.CommentedAt)
	}
	// Keyword filtering
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			audit_document_request_conversations.comment LIKE ? OR
			commented_user.first_name_en LIKE ? OR
            commented_user.first_name_ar LIKE ? OR
            commented_user.last_name_en LIKE ? OR
			commented_user.last_name_ar LIKE ?`,
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

	// Fetching related documents for each conversation
	for _, conversation := range conversations {
		var documents []*dtos.AuditDocumentRequestConversationDocumentDTO

		repo.db.Model(&models.AuditDocumentRequestConversationDocument{}).
			Where("audit_document_request_conversation_id = ?", conversation.ID).
			Find(&documents)

		conversation.Documents = documents
	}

	return conversations, result.Error
}

func (repo *AuditDocumentRequestConversationRepository) GetByID(id uint) (*models.AuditDocumentRequestConversation, error) {
	var conversation *models.AuditDocumentRequestConversation
	result := repo.db.First(&conversation, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return conversation, nil
}

func (repo *AuditDocumentRequestConversationRepository) Create(conversation *models.AuditDocumentRequestConversation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(conversation).Error; err != nil {
		return 0, err
	}

	return conversation.ID, nil
}

func (repo *AuditDocumentRequestConversationRepository) CreateDocuments(docs []*models.AuditDocumentRequestConversationDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *AuditDocumentRequestConversationRepository) AuditDocumentRequestStatusUpdate(auditDocumentRequest *models.AuditDocumentRequest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditDocumentRequest).Error; err != nil {
		return err
	}

	return nil
}
