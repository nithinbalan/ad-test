package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, currentUserID uint, params *riskDTOs.RiskReviewCommentListQueryParamsDTO) (*riskDTOs.RiskReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.RiskReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskReviewComment, error)
	Create(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
}

type RiskReviewCommentRepository struct {
	db *gorm.DB
}

func NewRiskReviewCommentRepository(db *gorm.DB) RiskReviewCommentRepositoryInterface {
	return &RiskReviewCommentRepository{
		db: db,
	}
}

func (repo *RiskReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskReviewCommentRepository) GetAll(riskID uint, currentUserID uint, params *riskDTOs.RiskReviewCommentListQueryParamsDTO) (*riskDTOs.RiskReviewCommentListResponseDTO, error) {

	var comments []*riskDTOs.RiskReviewCommentListDTO
	var statusCounts []riskDTOs.StatusCount

	query := repo.db.
		Table("risk_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = risk_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = risk_review_comments.created_by_id").
		Where("risk_review_comments.risk_register_id = ? ", riskID)

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("risk_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("risk_review_comments.created_by_id = ?", currentUserID)
		}
	}

	// Select fields
	query = query.Select(`
        risk_review_comments.id,
        risk_review_comments.comment,
        risk_review_comments.risk_register_id,
        risk_review_comments.level,
        risk_review_comments.created_at,
        risk_review_comments.created_by_id,
        risk_review_comments.risk_review_comment_id AS risk_review_comment_id,
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
    `).Order("risk_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.
		Table("risk_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = risk_review_comments.comment_status_id").
		Where("risk_review_comments.risk_register_id = ? AND risk_review_comment_id IS NULL", riskID)

		// Apply same filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("risk_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("risk_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("risk_review_comments.created_by_id = ?", currentUserID)
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

	return &riskDTOs.RiskReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *RiskReviewCommentRepository) GetByID(id uint) (*models.RiskReviewComment, error) {
	var RiskReviewComment *models.RiskReviewComment
	result := repo.db.First(&RiskReviewComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return RiskReviewComment, nil
}

func (repo *RiskReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.RiskReviewComment
	result := db.
		Preload("RiskRegister").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *RiskReviewCommentRepository) Create(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *RiskReviewCommentRepository) Update(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.RiskReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}
