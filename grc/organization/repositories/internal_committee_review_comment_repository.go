package repositories

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type InternalCommitteeReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(committeeVersionID uint, currentUserID uint, params *dtos.InternalCommitteeReviewCommentListQueryParamsDTO) (*dtos.InternalCommitteeReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.InternalCommitteeReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeReviewComment, error)
	Create(comment *models.InternalCommitteeReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.InternalCommitteeReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
}

type InternalCommitteeReviewCommentRepository struct {
	db *gorm.DB
}

func NewInternalCommitteeReviewCommentRepository(db *gorm.DB) *InternalCommitteeReviewCommentRepository {
	return &InternalCommitteeReviewCommentRepository{db: db}
}

func (repo *InternalCommitteeReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *InternalCommitteeReviewCommentRepository) GetAll(committeeVersionID uint, currentUserID uint, params *dtos.InternalCommitteeReviewCommentListQueryParamsDTO) (*dtos.InternalCommitteeReviewCommentListResponseDTO, error) {
	var comments []*dtos.InternalCommitteeReviewCommentListDTO
	var statusCounts []dtos.StatusCount

	query := repo.db.
		Table("internal_committee_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = internal_committee_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = internal_committee_review_comments.created_by_id").
		Where("internal_committee_review_comments.internal_committee_version_id = ?", committeeVersionID)

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("internal_committee_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("internal_committee_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comments":
			query = query.Where("internal_committee_review_comments.created_by_id = ?", currentUserID)
		}
	}

	// Select fields
	query = query.Select(`
		internal_committee_review_comments.id,
		internal_committee_review_comments.comment,
		internal_committee_review_comments.internal_committee_version_id,
		internal_committee_review_comments.level,
		internal_committee_review_comments.created_at,
		internal_committee_review_comments.created_by_id,
		internal_committee_review_comments.internal_committee_review_comment_id AS parent_comment_id,
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
	`).Order("internal_committee_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	// Count query
	countQuery := repo.db.
		Table("internal_committee_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = internal_committee_review_comments.comment_status_id").
		Where("internal_committee_review_comments.internal_committee_version_id = ? AND internal_committee_review_comment_id IS NULL", committeeVersionID)

	// Apply same filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("internal_committee_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("internal_committee_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comments":
			countQuery = countQuery.Where("internal_committee_review_comments.created_by_id = ?", currentUserID)
		}
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

	return &dtos.InternalCommitteeReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *InternalCommitteeReviewCommentRepository) GetByID(id uint) (*models.InternalCommitteeReviewComment, error) {
	var comment models.InternalCommitteeReviewComment
	result := repo.db.First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (repo *InternalCommitteeReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeReviewComment, error) {
	var db *gorm.DB
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	} else {
		db = repo.db
	}

	var comment models.InternalCommitteeReviewComment
	result := db.
		Preload("CommentStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *InternalCommitteeReviewCommentRepository) Create(comment *models.InternalCommitteeReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(comment)
	if result.Error != nil {
		return 0, result.Error
	}

	return comment.ID, nil
}

func (repo *InternalCommitteeReviewCommentRepository) Update(comment *models.InternalCommitteeReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Save(comment)
	return result.Error
}

func (repo *InternalCommitteeReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Delete(&models.InternalCommitteeReviewComment{}, commentID)
	return result.Error
}
