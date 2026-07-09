package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BiaProcessReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(biaProcessID uint, currentUserID uint, params *businessContinuityDTOs.BiaProcessReviewCommentListQueryParamsDTO) (*businessContinuityDTOs.BiaProcessReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.BiaProcessReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcessReviewComment, error)
	Create(comment *models.BiaProcessReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.BiaProcessReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.BiaProcessReviewCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.BiaProcessReviewCommentDocument, error)
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type BiaProcessReviewCommentRepository struct {
	db *gorm.DB
}

func NewBiaProcessReviewCommentRepository(db *gorm.DB) BiaProcessReviewCommentRepositoryInterface {
	return &BiaProcessReviewCommentRepository{db: db}
}

func (repo *BiaProcessReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BiaProcessReviewCommentRepository) GetAll(biaProcessID uint, currentUserID uint, params *businessContinuityDTOs.BiaProcessReviewCommentListQueryParamsDTO) (*businessContinuityDTOs.BiaProcessReviewCommentListResponseDTO, error) {
	var comments []*businessContinuityDTOs.BiaProcessReviewCommentListDTO
	var statusCounts []businessContinuityDTOs.StatusCount

	if params == nil {
		params = &businessContinuityDTOs.BiaProcessReviewCommentListQueryParamsDTO{}
	}

	query := repo.db.
		Table("bia_process_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = bia_process_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = bia_process_review_comments.created_by_id").
		Where("bia_process_review_comments.bia_process_id = ?", biaProcessID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("bia_process_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("bia_process_review_comments.comment LIKE ?", keyword)
	}
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("bia_process_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
		bia_process_review_comments.id,
		bia_process_review_comments.comment,
		bia_process_review_comments.is_edited,
		bia_process_review_comments.bia_process_id,
		bia_process_review_comments.level,
		bia_process_review_comments.created_at,
		bia_process_review_comments.created_by_id AS created_user_id,
		bia_process_review_comments.parent_comment_id AS parent_comment_id,
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
	`).Order("bia_process_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	if len(comments) > 0 {
		var commentIDs []uint
		for _, comment := range comments {
			commentIDs = append(commentIDs, comment.ID)
		}

		type docWithCommentID struct {
			BiaProcessReviewCommentID uint `gorm:"column:bia_process_review_comment_id"`
			businessContinuityDTOs.BiaProcessReviewCommentDocumentDTO
		}
		var allDocuments []docWithCommentID
		err := repo.db.
			Table("bia_process_review_comment_documents").
			Select("id, bia_process_review_comment_id, url, name, original_name, token, type, size").
			Where("bia_process_review_comment_id IN ?", commentIDs).
			Scan(&allDocuments).Error
		if err != nil {
			return nil, err
		}

		docMap := make(map[uint][]*businessContinuityDTOs.BiaProcessReviewCommentDocumentDTO)
		for _, doc := range allDocuments {
			d := doc.BiaProcessReviewCommentDocumentDTO
			docMap[doc.BiaProcessReviewCommentID] = append(docMap[doc.BiaProcessReviewCommentID], &d)
		}

		for _, comment := range comments {
			if docs, ok := docMap[comment.ID]; ok {
				comment.Documents = docs
			} else {
				comment.Documents = []*businessContinuityDTOs.BiaProcessReviewCommentDocumentDTO{}
			}
		}
	}

	countQuery := repo.db.
		Table("bia_process_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = bia_process_review_comments.comment_status_id").
		Where("bia_process_review_comments.bia_process_id = ? AND parent_comment_id IS NULL", biaProcessID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("bia_process_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("bia_process_review_comments.comment LIKE ?", keyword)
	}
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("bia_process_review_comments.created_by_id = ?", currentUserID)
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

	return &businessContinuityDTOs.BiaProcessReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *BiaProcessReviewCommentRepository) GetByID(id uint) (*models.BiaProcessReviewComment, error) {
	var comment models.BiaProcessReviewComment
	result := repo.db.First(&comment, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (repo *BiaProcessReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BiaProcessReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var comment models.BiaProcessReviewComment
	result := db.
		Preload("BiaProcess").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &comment, nil
}

func (repo *BiaProcessReviewCommentRepository) Create(comment *models.BiaProcessReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (repo *BiaProcessReviewCommentRepository) Update(comment *models.BiaProcessReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.BiaProcessReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error
}

func (repo *BiaProcessReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", commentID).Delete(&models.BiaProcessReviewComment{}).Error
}

func (repo *BiaProcessReviewCommentRepository) CreateDocuments(documents []*models.BiaProcessReviewCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (repo *BiaProcessReviewCommentRepository) GetDocumentByID(id uint) (*models.BiaProcessReviewCommentDocument, error) {
	var document models.BiaProcessReviewCommentDocument
	result := repo.db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &document, nil
}

func (repo *BiaProcessReviewCommentRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var ids []uint
	for _, id := range documentIDs {
		if id != nil {
			ids = append(ids, *id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	return db.Where("id IN ?", ids).Delete(&models.BiaProcessReviewCommentDocument{}).Error
}
