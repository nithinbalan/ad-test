package repositories

import (
	"jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentClosureReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(incidentClosureID uint, currentUserID uint, params *dtos.IncidentClosureReviewCommentListQueryParamsDTO) (*dtos.IncidentClosureReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.IncidentClosureReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosureReviewComment, error)
	Create(comment *models.IncidentClosureReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.IncidentClosureReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.IncidentClosureReviewCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IncidentClosureReviewCommentDocument, error)
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentClosureReviewCommentRepository struct {
	db *gorm.DB
}

func NewIncidentClosureReviewCommentRepository(db *gorm.DB) IncidentClosureReviewCommentRepositoryInterface {
	return &IncidentClosureReviewCommentRepository{
		db: db,
	}
}

func (repo *IncidentClosureReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentClosureReviewCommentRepository) GetAll(incidentClosureID uint, currentUserID uint, params *dtos.IncidentClosureReviewCommentListQueryParamsDTO) (*dtos.IncidentClosureReviewCommentListResponseDTO, error) {
	var comments []*dtos.IncidentClosureReviewCommentListDTO
	var statusCounts []dtos.StatusCount

	if params == nil {
		params = &dtos.IncidentClosureReviewCommentListQueryParamsDTO{}
	}

	query := repo.db.
		Table("incident_closure_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = incident_closure_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = incident_closure_review_comments.created_by_id").
		Where("incident_closure_review_comments.incident_closure_id = ? ", incidentClosureID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("incident_closure_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("incident_closure_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("incident_closure_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
        incident_closure_review_comments.id,
        incident_closure_review_comments.comment,
        incident_closure_review_comments.is_edited,
        incident_closure_review_comments.incident_closure_id,
        incident_closure_review_comments.level,
        incident_closure_review_comments.created_at,
        incident_closure_review_comments.created_by_id AS created_user_id,
        incident_closure_review_comments.incident_closure_review_comment_id AS parent_comment_id,
		comment_statuses.id AS comment_status_id,
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
    `).Order("incident_closure_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*dtos.IncidentClosureReviewCommentDocumentDTO

		err := repo.db.
			Table("incident_closure_review_comment_documents").
			Select("id, url, name, original_name, token, type, size").
			Where("incident_closure_review_comment_id = ?", comment.ID).
			Scan(&documents).Error

		if err != nil {
			return nil, err
		}

		if documents == nil {
			documents = []*dtos.IncidentClosureReviewCommentDocumentDTO{}
		}

		comment.Documents = documents
	}

	countQuery := repo.db.
		Table("incident_closure_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = incident_closure_review_comments.comment_status_id").
		Where("incident_closure_review_comments.incident_closure_id = ? AND incident_closure_review_comment_id IS NULL", incidentClosureID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("incident_closure_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("incident_closure_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("incident_closure_review_comments.created_by_id = ?", currentUserID)
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

	return &dtos.IncidentClosureReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *IncidentClosureReviewCommentRepository) GetByID(id uint) (*models.IncidentClosureReviewComment, error) {
	var comment *models.IncidentClosureReviewComment
	result := repo.db.First(&comment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *IncidentClosureReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentClosureReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.IncidentClosureReviewComment
	result := db.
		Preload("IncidentClosure").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *IncidentClosureReviewCommentRepository) Create(comment *models.IncidentClosureReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *IncidentClosureReviewCommentRepository) Update(comment *models.IncidentClosureReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.IncidentClosureReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.IncidentClosureReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureReviewCommentRepository) CreateDocuments(documents []*models.IncidentClosureReviewCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&documents).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentClosureReviewCommentRepository) GetDocumentByID(id uint) (*models.IncidentClosureReviewCommentDocument, error) {
	var document models.IncidentClosureReviewCommentDocument
	result := repo.db.First(&document, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &document, nil
}

func (repo *IncidentClosureReviewCommentRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN ?", documentIDs).Delete(&models.IncidentClosureReviewCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}
