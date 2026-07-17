package repositories

import (
	assetRiskDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskID uint, currentUserID uint, params *assetRiskDTOs.AssetRiskReviewCommentListQueryParamsDTO) (*assetRiskDTOs.AssetRiskReviewCommentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskReviewComment, error)
	Create(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
}

type AssetRiskReviewCommentRepository struct {
	db *gorm.DB
}

func NewAssetRiskReviewCommentRepository(db *gorm.DB) AssetRiskReviewCommentRepositoryInterface {
	return &AssetRiskReviewCommentRepository{
		db: db,
	}
}

func (repo *AssetRiskReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

type assetRiskReviewCommentStatusCount struct {
	Status string
	Count  int
}

func (repo *AssetRiskReviewCommentRepository) GetAll(riskID uint, currentUserID uint, params *assetRiskDTOs.AssetRiskReviewCommentListQueryParamsDTO) (*assetRiskDTOs.AssetRiskReviewCommentListResponseDTO, error) {

	var comments []*assetRiskDTOs.AssetRiskReviewCommentListDTO
	var statusCounts []assetRiskReviewCommentStatusCount

	query := repo.db.
		Table("risk_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = risk_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = risk_review_comments.created_by_id").
		Where("risk_review_comments.risk_register_id = ? ", riskID)

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

	return &assetRiskDTOs.AssetRiskReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *AssetRiskReviewCommentRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment *models.RiskReviewComment
	result := db.First(&comment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *AssetRiskReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskReviewComment, error) {
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

func (repo *AssetRiskReviewCommentRepository) Create(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *AssetRiskReviewCommentRepository) Update(comment *models.RiskReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.RiskReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}
