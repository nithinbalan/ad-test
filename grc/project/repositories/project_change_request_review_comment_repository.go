package repositories

import (
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectChangeRequestReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(projectChangeRequestID uint, currentUserID uint, params *projectDTOs.ProjectChangeRequestReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectChangeRequestReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.ProjectChangeRequestReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestReviewComment, error)
	Create(comment *models.ProjectChangeRequestReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ProjectChangeRequestReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.ProjectChangeRequestReviewCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ProjectChangeRequestReviewCommentDocument, error)
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type ProjectChangeRequestReviewCommentRepository struct {
	db *gorm.DB
}

func NewProjectChangeRequestReviewCommentRepository(db *gorm.DB) ProjectChangeRequestReviewCommentRepositoryInterface {
	return &ProjectChangeRequestReviewCommentRepository{
		db: db,
	}
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pcrrcr.db)
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) GetAll(projectChangeRequestID uint, currentUserID uint, params *projectDTOs.ProjectChangeRequestReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectChangeRequestReviewCommentListResponseDTO, error) {
	var comments []*projectDTOs.ProjectChangeRequestReviewCommentListDTO
	var statusCounts []struct {
		Status string
		Count  int
	}

	if params == nil {
		params = &projectDTOs.ProjectChangeRequestReviewCommentListQueryParamsDTO{}
	}

	query := pcrrcr.db.
		Table("project_change_request_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = project_change_request_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = project_change_request_review_comments.created_by_id").
		Where("project_change_request_review_comments.project_change_request_id = ?", projectChangeRequestID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_change_request_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("project_change_request_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("project_change_request_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
        project_change_request_review_comments.id,
        project_change_request_review_comments.comment,
        project_change_request_review_comments.is_edited,
        project_change_request_review_comments.project_change_request_id,
        project_change_request_review_comments.level,
        project_change_request_review_comments.created_at,
        project_change_request_review_comments.created_by_id AS created_user_id,
        project_change_request_review_comments.parent_comment_id AS parent_comment_id,
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
    `).Order("project_change_request_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*projectDTOs.ProjectChangeRequestReviewCommentDocumentDTO

		err := pcrrcr.db.
			Table("project_change_request_review_comment_documents").
			Select("id, token, name, original_name, type, size").
			Where("project_change_request_review_comment_id = ?", comment.ID).
			Scan(&documents).Error

		if err != nil {
			return nil, err
		}

		if documents == nil {
			documents = []*projectDTOs.ProjectChangeRequestReviewCommentDocumentDTO{}
		}

		comment.Documents = documents
	}

	countQuery := pcrrcr.db.
		Table("project_change_request_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = project_change_request_review_comments.comment_status_id").
		Where("project_change_request_review_comments.project_change_request_id = ? AND parent_comment_id IS NULL", projectChangeRequestID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("project_change_request_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("project_change_request_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("project_change_request_review_comments.created_by_id = ?", currentUserID)
		}
	}

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

	return &projectDTOs.ProjectChangeRequestReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) GetByID(id uint) (*models.ProjectChangeRequestReviewComment, error) {
	var comment models.ProjectChangeRequestReviewComment
	result := pcrrcr.db.First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestReviewComment, error) {
	db := pcrrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.ProjectChangeRequestReviewComment
	result := db.
		Preload("ProjectChangeRequest").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) Create(comment *models.ProjectChangeRequestReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pcrrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) Update(comment *models.ProjectChangeRequestReviewComment, tx ...utility.DatabaseTransaction) error {
	db := pcrrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ProjectChangeRequestReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", commentID).Delete(&models.ProjectChangeRequestReviewComment{}).Error
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) CreateDocuments(documents []*models.ProjectChangeRequestReviewCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := pcrrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) GetDocumentByID(id uint) (*models.ProjectChangeRequestReviewCommentDocument, error) {
	var document models.ProjectChangeRequestReviewCommentDocument
	result := pcrrcr.db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &document, nil
}

func (pcrrcr *ProjectChangeRequestReviewCommentRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := pcrrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", documentIDs).Delete(&models.ProjectChangeRequestReviewCommentDocument{}).Error
}
