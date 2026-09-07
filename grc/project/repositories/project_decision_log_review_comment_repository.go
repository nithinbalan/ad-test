package repositories

import (
	"jethur/grc/models"
	projectDTOs "jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectDecisionLogReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(projectDecisionLogID uint, currentUserID uint, params *projectDTOs.ProjectDecisionLogReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectDecisionLogReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.ProjectDecisionLogReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogReviewComment, error)
	Create(comment *models.ProjectDecisionLogReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ProjectDecisionLogReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.ProjectDecisionLogReviewCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ProjectDecisionLogReviewCommentDocument, error)
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type ProjectDecisionLogReviewCommentRepository struct {
	db *gorm.DB
}

func NewProjectDecisionLogReviewCommentRepository(db *gorm.DB) ProjectDecisionLogReviewCommentRepositoryInterface {
	return &ProjectDecisionLogReviewCommentRepository{
		db: db,
	}
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pdlrcr.db)
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) GetAll(projectDecisionLogID uint, currentUserID uint, params *projectDTOs.ProjectDecisionLogReviewCommentListQueryParamsDTO) (*projectDTOs.ProjectDecisionLogReviewCommentListResponseDTO, error) {
	var comments []*projectDTOs.ProjectDecisionLogReviewCommentListDTO
	var statusCounts []struct {
		Status string
		Count  int
	}

	if params == nil {
		params = &projectDTOs.ProjectDecisionLogReviewCommentListQueryParamsDTO{}
	}

	query := pdlrcr.db.
		Table("project_decision_log_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = project_decision_log_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = project_decision_log_review_comments.created_by_id").
		Where("project_decision_log_review_comments.project_decision_log_id = ?", projectDecisionLogID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_decision_log_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("project_decision_log_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("project_decision_log_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
        project_decision_log_review_comments.id,
        project_decision_log_review_comments.comment,
        project_decision_log_review_comments.is_edited,
        project_decision_log_review_comments.project_decision_log_id,
        project_decision_log_review_comments.level,
        project_decision_log_review_comments.created_at,
        project_decision_log_review_comments.created_by_id AS created_user_id,
        project_decision_log_review_comments.parent_comment_id AS parent_comment_id,
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
    `).Order("project_decision_log_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*projectDTOs.ProjectDecisionLogReviewCommentDocumentDTO

		err := pdlrcr.db.
			Table("project_decision_log_review_comment_documents").
			Select("id, token, name, original_name, type, size").
			Where("project_decision_log_review_comment_id = ?", comment.ID).
			Scan(&documents).Error

		if err != nil {
			return nil, err
		}

		if documents == nil {
			documents = []*projectDTOs.ProjectDecisionLogReviewCommentDocumentDTO{}
		}

		comment.Documents = documents
	}

	commentMap := make(map[uint]*projectDTOs.ProjectDecisionLogReviewCommentListDTO)
	var rootComments []*projectDTOs.ProjectDecisionLogReviewCommentListDTO

	for _, comment := range comments {
		comment.Children = []*projectDTOs.ProjectDecisionLogReviewCommentListDTO{}
		commentMap[comment.ID] = comment
	}

	for _, comment := range comments {
		if comment.ParentCommentID == nil {
			rootComments = append(rootComments, comment)
		} else {
			if parent, exists := commentMap[*comment.ParentCommentID]; exists {
				parent.Children = append(parent.Children, comment)
			}
		}
	}

	countQuery := pdlrcr.db.
		Table("project_decision_log_review_comments").
		Select("comment_statuses.type AS status, COUNT(project_decision_log_review_comments.id) as count").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = project_decision_log_review_comments.comment_status_id").
		Where("project_decision_log_review_comments.project_decision_log_id = ? AND project_decision_log_review_comments.parent_comment_id IS NULL", projectDecisionLogID).
		Group("comment_statuses.type")

	if err := countQuery.Scan(&statusCounts).Error; err != nil {
		return nil, err
	}

	var openCount, resolvedCount int
	for _, sc := range statusCounts {
		switch sc.Status {
		case "open":
			openCount = sc.Count
		case "resolved":
			resolvedCount = sc.Count
		}
	}

	response := &projectDTOs.ProjectDecisionLogReviewCommentListResponseDTO{
		Comments:      rootComments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}

	return response, nil
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) GetByID(id uint) (*models.ProjectDecisionLogReviewComment, error) {
	var comment models.ProjectDecisionLogReviewComment
	result := pdlrcr.db.
		Preload("CommentStatus").
		Preload("Documents").
		Preload("Children").
		Where("id = ?", id).
		First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectDecisionLogReviewComment, error) {
	db := pdlrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var comment models.ProjectDecisionLogReviewComment
	result := db.
		Preload("CommentStatus").
		Preload("Documents").
		Preload("Children").
		Where("id = ?", id).
		First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) Create(comment *models.ProjectDecisionLogReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pdlrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) Update(comment *models.ProjectDecisionLogReviewComment, tx ...utility.DatabaseTransaction) error {
	db := pdlrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(comment).Error
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var childCommentIDs []uint
	if err := db.Model(&models.ProjectDecisionLogReviewComment{}).
		Where("parent_comment_id = ?", commentID).
		Pluck("id", &childCommentIDs).Error; err != nil {
		return err
	}

	allCommentIDs := append(childCommentIDs, commentID)

	if err := db.Where("project_decision_log_review_comment_id IN (?)", allCommentIDs).
		Delete(&models.ProjectDecisionLogReviewCommentDocument{}).Error; err != nil {
		return err
	}

	if len(childCommentIDs) > 0 {
		if err := db.Where("id IN (?)", childCommentIDs).
			Delete(&models.ProjectDecisionLogReviewComment{}).Error; err != nil {
			return err
		}
	}

	return db.Delete(&models.ProjectDecisionLogReviewComment{}, commentID).Error
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) CreateDocuments(documents []*models.ProjectDecisionLogReviewCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := pdlrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) GetDocumentByID(id uint) (*models.ProjectDecisionLogReviewCommentDocument, error) {
	var document models.ProjectDecisionLogReviewCommentDocument
	result := pdlrcr.db.Where("id = ?", id).First(&document)
	if result.Error != nil {
		return nil, result.Error
	}
	return &document, nil
}

func (pdlrcr *ProjectDecisionLogReviewCommentRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := pdlrcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", documentIDs).Delete(&models.ProjectDecisionLogReviewCommentDocument{}).Error
}
