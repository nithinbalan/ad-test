package repositories

import (
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type PolicyClauseCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllDocumentComments(documentID uint, versionID uint, currentUserID uint, params *policyDTOs.PolicyClauseCommentListQueryParamsDTO) (*policyDTOs.PolicyClauseCommentListResponseDTO, error)
	GetByID(id uint) (*models.DocumentVersionClauseComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClauseComment, error)
	Create(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
}

type PolicyClauseCommentRepository struct {
	db *gorm.DB
}

func NewPolicyClauseCommentRepository(db *gorm.DB) PolicyClauseCommentRepositoryInterface {
	return &PolicyClauseCommentRepository{
		db: db,
	}
}

func (repo *PolicyClauseCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PolicyClauseCommentRepository) GetAllDocumentComments(documentID uint, versionID uint, currentUserID uint, params *policyDTOs.PolicyClauseCommentListQueryParamsDTO) (*policyDTOs.PolicyClauseCommentListResponseDTO, error) {
	var comments []*policyDTOs.PolicyClauseCommentListDTO
	var statusCounts []policyDTOs.StatusCount

	query := repo.db.
		Table("document_version_clause_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = document_version_clause_comments.comment_status_id").
		Joins("LEFT JOIN document_version_clauses ON document_version_clauses.id = document_version_clause_comments.document_version_clause_id").
		Joins("LEFT JOIN document_versions ON document_versions.id = document_version_clause_comments.document_version_id").
		Joins("LEFT JOIN users ON users.id = document_version_clause_comments.created_by_id").
		Where("document_versions.id = ? ", versionID)

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("document_version_clause_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("document_version_clause_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("document_version_clause_comments.created_by_id = ?", currentUserID)
		}
	}

	// Select fields
	query = query.Select(`
        document_version_clause_comments.id,
        document_version_clause_comments.comment,
        document_version_clause_comments.document_version_clause_id,
        document_version_clause_comments.level,
        document_version_clause_comments.created_at,
        document_version_clause_comments.created_by_id,
        document_version_clauses.clause_number AS document_version_clause_number,
        document_version_clause_comments.document_version_clause_comment_id AS document_version_clause_comment_id,
		document_version_clause_comments.document_version_id,
        comment_statuses.title_en AS comment_status_title_en,
        comment_statuses.title_ar AS comment_status_title_ar,
        comment_statuses.type AS comment_status_type,
        comment_statuses.color_code AS comment_status_color_code,
        users.first_name_en AS created_user_first_name_en,
        users.first_name_ar AS created_user_first_name_ar,
        users.last_name_en AS created_user_last_name_en,
        users.last_name_ar AS created_user_last_name_ar,
        users.initial AS created_user_initial,
        users.color AS created_user_color,
        users.email AS created_user_email,
        users.profile_image_path AS created_user_profile_image_path
    `).Order("document_version_clause_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.
		Table("document_version_clause_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = document_version_clause_comments.comment_status_id").
		Joins("LEFT JOIN document_version_clauses ON document_version_clauses.id = document_version_clause_comments.document_version_clause_id").
		Joins("LEFT JOIN document_versions ON document_versions.id = document_version_clause_comments.document_version_id").
		Where("document_versions.id = ? AND document_version_clause_comment_id IS NULL", versionID)

		// Apply same filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("document_version_clause_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("document_version_clause_comments.comment LIKE ?", keyword)
	}

	// Group by status
	countQuery = countQuery.
		Select("comment_statuses.type AS status, COUNT(*) as count").
		Group("comment_statuses.type")

	if err := countQuery.Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	openCount, resolvedCount := 0, 0
	for _, sc := range statusCounts {
		if sc.Status == "open" {
			openCount = sc.Count
		}
		if sc.Status == "resolved" {
			resolvedCount = sc.Count
		}
	}

	return &policyDTOs.PolicyClauseCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *PolicyClauseCommentRepository) GetByID(id uint) (*models.DocumentVersionClauseComment, error) {
	var PolicyClauseComment *models.DocumentVersionClauseComment
	result := repo.db.First(&PolicyClauseComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return PolicyClauseComment, nil
}

func (repo *PolicyClauseCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClauseComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.DocumentVersionClauseComment
	result := db.
		Preload("DocumentVersionClause").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *PolicyClauseCommentRepository) Create(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *PolicyClauseCommentRepository) Update(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.DocumentVersionClauseComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PolicyClauseCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.DocumentVersionClauseComment{}).Error; err != nil {
		return err
	}

	return nil
}
