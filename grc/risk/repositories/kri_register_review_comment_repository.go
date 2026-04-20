package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRIRegisterReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(kriRegisterID uint, currentUserID uint, params *riskDTOs.KRIRegisterReviewCommentListQueryParamsDTO) (*riskDTOs.KRIRegisterReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.KRIRegisterReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegisterReviewComment, error)
	Create(comment *models.KRIRegisterReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.KRIRegisterReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
}

type KRIRegisterReviewCommentRepository struct {
	db *gorm.DB
}

func NewKRIRegisterReviewCommentRepository(db *gorm.DB) KRIRegisterReviewCommentRepositoryInterface {
	return &KRIRegisterReviewCommentRepository{
		db: db,
	}
}

func (repo *KRIRegisterReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KRIRegisterReviewCommentRepository) GetAll(kriRegisterID uint, currentUserID uint, params *riskDTOs.KRIRegisterReviewCommentListQueryParamsDTO) (*riskDTOs.KRIRegisterReviewCommentListResponseDTO, error) {
	var comments []*riskDTOs.KRIRegisterReviewCommentListDTO
	var statusCounts []riskDTOs.StatusCount

	query := repo.db.
		Table("kri_register_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = kri_register_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = kri_register_review_comments.created_by_id").
		Where("kri_register_review_comments.kri_register_id = ? ", kriRegisterID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("kri_register_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("kri_register_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("kri_register_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
        kri_register_review_comments.id,
        kri_register_review_comments.comment,
        kri_register_review_comments.kri_register_id,
        kri_register_review_comments.level,
        kri_register_review_comments.created_at,
        kri_register_review_comments.created_by_id,
        kri_register_review_comments.parent_comment_id AS parent_comment_id,
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
    `).Order("kri_register_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	countQuery := repo.db.
		Table("kri_register_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = kri_register_review_comments.comment_status_id").
		Where("kri_register_review_comments.kri_register_id = ? AND parent_comment_id IS NULL", kriRegisterID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("kri_register_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("kri_register_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("kri_register_review_comments.created_by_id = ?", currentUserID)
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

	return &riskDTOs.KRIRegisterReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *KRIRegisterReviewCommentRepository) GetByID(id uint) (*models.KRIRegisterReviewComment, error) {
	var comment *models.KRIRegisterReviewComment
	result := repo.db.First(&comment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *KRIRegisterReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegisterReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.KRIRegisterReviewComment
	result := db.
		Preload("KRIRegister").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *KRIRegisterReviewCommentRepository) Create(comment *models.KRIRegisterReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *KRIRegisterReviewCommentRepository) Update(comment *models.KRIRegisterReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.KRIRegisterReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRIRegisterReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.KRIRegisterReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}
