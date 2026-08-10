package repositories

import (
	identityAccessReviewDTOs "jethur/grc/identity-access-review/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewCampaignReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(campaignID uint, currentUserID uint, params *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListQueryParamsDTO) (*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignReviewComment, error)
	Create(comment *models.IdentityReviewCampaignReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.IdentityReviewCampaignReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewCampaignReviewCommentRepository struct {
	db *gorm.DB
}

func NewIdentityReviewCampaignReviewCommentRepository(db *gorm.DB) IdentityReviewCampaignReviewCommentRepositoryInterface {
	return &IdentityReviewCampaignReviewCommentRepository{
		db: db,
	}
}

func (repo *IdentityReviewCampaignReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

type identityReviewCampaignReviewCommentStatusCount struct {
	Status string
	Count  int
}

func (repo *IdentityReviewCampaignReviewCommentRepository) GetAll(campaignID uint, currentUserID uint, params *identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListQueryParamsDTO) (*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListResponseDTO, error) {

	var comments []*identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListDTO
	var statusCounts []identityReviewCampaignReviewCommentStatusCount

	query := repo.db.
		Table("identity_review_campaign_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = identity_review_campaign_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = identity_review_campaign_review_comments.created_by_id").
		Where("identity_review_campaign_review_comments.identity_review_campaign_id = ? ", campaignID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("identity_review_campaign_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("identity_review_campaign_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("identity_review_campaign_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
        identity_review_campaign_review_comments.id,
        identity_review_campaign_review_comments.comment,
        identity_review_campaign_review_comments.identity_review_campaign_id,
        identity_review_campaign_review_comments.level,
        identity_review_campaign_review_comments.created_at,
        identity_review_campaign_review_comments.created_by_id,
        identity_review_campaign_review_comments.identity_review_campaign_review_comment_id AS identity_review_campaign_review_comment_id,
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
    `).Order("identity_review_campaign_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.
		Table("identity_review_campaign_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = identity_review_campaign_review_comments.comment_status_id").
		Where("identity_review_campaign_review_comments.identity_review_campaign_id = ? AND identity_review_campaign_review_comment_id IS NULL", campaignID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("identity_review_campaign_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("identity_review_campaign_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("identity_review_campaign_review_comments.created_by_id = ?", currentUserID)
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

	return &identityAccessReviewDTOs.IdentityReviewCampaignReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *IdentityReviewCampaignReviewCommentRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment *models.IdentityReviewCampaignReviewComment
	result := db.First(&comment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *IdentityReviewCampaignReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewCampaignReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.IdentityReviewCampaignReviewComment
	result := db.
		Preload("IdentityReviewCampaign").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *IdentityReviewCampaignReviewCommentRepository) Create(comment *models.IdentityReviewCampaignReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *IdentityReviewCampaignReviewCommentRepository) Update(comment *models.IdentityReviewCampaignReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.IdentityReviewCampaignReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IdentityReviewCampaignReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.IdentityReviewCampaignReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}
