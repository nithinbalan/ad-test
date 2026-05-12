package repositories

import (
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(riskRegisterUpdatePlanID uint, currentUserID uint, params *riskDTOs.RiskUpdatePlanReviewCommentListQueryParamsDTO) (*riskDTOs.RiskUpdatePlanReviewCommentListResponseDTO, error)
	GetByID(id uint) (*models.RiskUpdatePlanReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanReviewComment, error)
	Create(comment *models.RiskUpdatePlanReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.RiskUpdatePlanReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(commentID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.RiskUpdatePlanReviewCommentDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskUpdatePlanReviewCommentDocument, error)
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type RiskUpdatePlanReviewCommentRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanReviewCommentRepository(db *gorm.DB) RiskUpdatePlanReviewCommentRepositoryInterface {
	return &RiskUpdatePlanReviewCommentRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskUpdatePlanReviewCommentRepository) GetAll(riskRegisterUpdatePlanID uint, currentUserID uint, params *riskDTOs.RiskUpdatePlanReviewCommentListQueryParamsDTO) (*riskDTOs.RiskUpdatePlanReviewCommentListResponseDTO, error) {
	var comments []*riskDTOs.RiskUpdatePlanReviewCommentListDTO
	var statusCounts []riskDTOs.StatusCount

	if params == nil {
		params = &riskDTOs.RiskUpdatePlanReviewCommentListQueryParamsDTO{}
	}

	query := repo.db.
		Table("risk_update_plan_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = risk_update_plan_review_comments.comment_status_id").
		Joins("LEFT JOIN users ON users.id = risk_update_plan_review_comments.created_by_id").
		Where("risk_update_plan_review_comments.risk_register_update_plan_id = ? ", riskRegisterUpdatePlanID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("risk_update_plan_review_comments."+columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_update_plan_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			query = query.Where("comment_statuses.type = ?", "open")
		case "resolved":
			query = query.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			query = query.Where("risk_update_plan_review_comments.created_by_id = ?", currentUserID)
		}
	}

	query = query.Select(`
        risk_update_plan_review_comments.id,
        risk_update_plan_review_comments.comment,
        risk_update_plan_review_comments.is_edited,
        risk_update_plan_review_comments.risk_register_update_plan_id,
        risk_update_plan_review_comments.level,
        risk_update_plan_review_comments.created_at,
        risk_update_plan_review_comments.created_by_id AS created_user_id,
        risk_update_plan_review_comments.parent_comment_id AS parent_comment_id,
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
    `).Order("risk_update_plan_review_comments.created_at DESC")

	if err := query.Scan(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*riskDTOs.RiskUpdatePlanReviewCommentDocumentDTO

		err := repo.db.
			Table("risk_update_plan_review_comment_documents").
			Select("id, url, name, original_name, token, type, size").
			Where("risk_update_plan_review_comment_id = ?", comment.ID).
			Scan(&documents).Error

		if err != nil {
			return nil, err
		}

		if documents == nil {
			documents = []*riskDTOs.RiskUpdatePlanReviewCommentDocumentDTO{}
		}

		comment.Documents = documents
	}

	countQuery := repo.db.
		Table("risk_update_plan_review_comments").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = risk_update_plan_review_comments.comment_status_id").
		Where("risk_update_plan_review_comments.risk_register_update_plan_id = ? AND parent_comment_id IS NULL", riskRegisterUpdatePlanID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			countQuery = countQuery.Where("risk_update_plan_review_comments."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		countQuery = countQuery.Where("risk_update_plan_review_comments.comment LIKE ?", keyword)
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "open":
			countQuery = countQuery.Where("comment_statuses.type = ?", "open")
		case "resolved":
			countQuery = countQuery.Where("comment_statuses.type = ?", "resolved")
		case "my_comment":
			countQuery = countQuery.Where("risk_update_plan_review_comments.created_by_id = ?", currentUserID)
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

	return &riskDTOs.RiskUpdatePlanReviewCommentListResponseDTO{
		Comments:      comments,
		OpenCount:     openCount,
		ResolvedCount: resolvedCount,
	}, nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) GetByID(id uint) (*models.RiskUpdatePlanReviewComment, error) {
	var comment *models.RiskUpdatePlanReviewComment
	result := repo.db.First(&comment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.RiskUpdatePlanReviewComment
	result := db.
		Preload("RiskRegisterUpdatePlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) Create(comment *models.RiskUpdatePlanReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) Update(comment *models.RiskUpdatePlanReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskUpdatePlanReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) Delete(commentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", commentID).Delete(&models.RiskUpdatePlanReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) CreateDocuments(documents []*models.RiskUpdatePlanReviewCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&documents).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) GetDocumentByID(id uint) (*models.RiskUpdatePlanReviewCommentDocument, error) {
	var document models.RiskUpdatePlanReviewCommentDocument
	result := repo.db.First(&document, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &document, nil
}

func (repo *RiskUpdatePlanReviewCommentRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN ?", documentIDs).Delete(&models.RiskUpdatePlanReviewCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}
