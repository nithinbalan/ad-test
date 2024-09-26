package repositories

import (
	"jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.AuditDraftReportCommentListQueryParamsDTO) (*dtos.AuditDraftReportCommentResponseDTO, error)
	GetByID(id uint) (*models.AuditDraftReportComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftReportComment, error)
	Create(auditDraftReportComment *models.AuditDraftReportComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditDraftReportComment *models.AuditDraftReportComment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.AuditDraftReportCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.AuditDraftReportCommentDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type AuditDraftReportCommentRepository struct {
	db *gorm.DB
}

func NewAuditDraftReportCommentRepository(db *gorm.DB) *AuditDraftReportCommentRepository {
	return &AuditDraftReportCommentRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDraftReportCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDraftReportCommentRepository) GetAll(params *dtos.AuditDraftReportCommentListQueryParamsDTO) (*dtos.AuditDraftReportCommentResponseDTO, error) {

	var comments []*dtos.AuditDraftReportCommentListDTO
	var totalCount int64

	query := repo.db.Table("audit_draft_report_comments").
		Select(`
			audit_draft_report_comments.*,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN users AS created_user ON created_user.id = audit_draft_report_comments.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_draft_report_comments.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(audit_draft_report_comments.comment LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.last_name_en LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination and sorting
	// if !params.Export {
	// 	offset := (params.Page - 1) * params.PageSize
	// 	query = query.Offset(offset).Limit(params.PageSize)
	// }
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)
	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := 1
	// if !params.Export {
	// 	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	// }
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Prepare the response
	response := &dtos.AuditDraftReportCommentResponseDTO{
		Items: comments,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditDraftReportCommentRepository) GetByID(id uint) (*models.AuditDraftReportComment, error) {
	var auditDraftReportComment *models.AuditDraftReportComment
	result := repo.db.First(&auditDraftReportComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReportComment, nil
}

func (repo *AuditDraftReportCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditDraftReportComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditDraftReportComment *models.AuditDraftReportComment

	result := db.
		Preload("AuditDraftFinding").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditDraftReportComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReportComment, nil
}

func (repo *AuditDraftReportCommentRepository) Create(auditDraftReportComment *models.AuditDraftReportComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditDraftReportComment).Error; err != nil {
		return 0, err
	}

	return auditDraftReportComment.ID, nil
}

func (repo *AuditDraftReportCommentRepository) Update(auditDraftReportComment *models.AuditDraftReportComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditDraftReportComment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportCommentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit Plan
	if err := db.Delete(&models.AuditDraftReportComment{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportCommentRepository) CreateDocuments(docs []*models.AuditDraftReportCommentDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *AuditDraftReportCommentRepository) GetDocumentByID(id uint) (*models.AuditDraftReportCommentDocument, error) {
	var document models.AuditDraftReportCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *AuditDraftReportCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.AuditDraftReportCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}
