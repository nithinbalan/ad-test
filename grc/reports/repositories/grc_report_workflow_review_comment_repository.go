package repositories

import (
	"fmt"
	"jethur/grc/models"
	"jethur/grc/reports/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type GRCReportWorkflowReviewCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(grcReportID uint, params *dtos.GRCReportWorkflowReviewCommentListQueryParamsDTO) ([]*dtos.GRCReportWorkflowReviewCommentListDTO, error)
	GetByID(id uint) (*models.GRCReportWorkflowReviewComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.GRCReportWorkflowReviewComment, error)
	Create(comment *models.GRCReportWorkflowReviewComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.GRCReportWorkflowReviewComment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	UpdateSectionWorkflowStatus(risk *models.GRCReportWorkflow, tx ...utility.DatabaseTransaction) error
}

type GRCReportWorkflowReviewCommentRepository struct {
	db *gorm.DB
}

func NewGRCReportWorkflowReviewCommentRepository(db *gorm.DB) GRCReportWorkflowReviewCommentRepositoryInterface {
	return &GRCReportWorkflowReviewCommentRepository{
		db: db,
	}
}

func (repo *GRCReportWorkflowReviewCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *GRCReportWorkflowReviewCommentRepository) GetAll(grcReportID uint, params *dtos.GRCReportWorkflowReviewCommentListQueryParamsDTO) ([]*dtos.GRCReportWorkflowReviewCommentListDTO, error) {
	var comments []*dtos.GRCReportWorkflowReviewCommentListDTO
	query := repo.db.Table("grc_report_workflow_review_comments").
		Select(`
			grc_report_workflow_review_comments.id,
			grc_report_workflow_review_comments.comment,
			grc_report_workflow_review_comments.grc_report_workflow_id,
			grc_report_workflow_review_comments.created_at,
			grc_report_workflow_review_comments.created_by_id,
			users.first_name_en AS created_user_first_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_en AS created_user_last_name_en,
			users.last_name_ar AS created_user_last_name_ar,
			users.initial AS created_user_initial,
			users.color AS created_user_color,
			users.email AS created_user_email,
			users.profile_image_path AS created_user_profile_image_path,
			grc_report_comment_statuses.title_en AS grc_report_comment_status_title_en,
			grc_report_comment_statuses.title_ar AS grc_report_comment_status_title_ar,
			grc_report_comment_statuses.type AS grc_report_comment_status_type,
			grc_report_comment_statuses.color_code AS grc_report_comment_status_color_code
		`).
		Joins("LEFT JOIN users ON users.id = grc_report_workflow_review_comments.created_by_id").
		Joins("LEFT JOIN grc_report_comment_statuses ON grc_report_comment_statuses.id = grc_report_workflow_review_comments.grc_report_comment_status_id").
		Joins("LEFT JOIN grc_report_workflows ON grc_report_workflows.id = grc_report_workflow_review_comments.grc_report_workflow_id").
		Joins("LEFT JOIN grc_reports ON grc_reports.id = grc_report_workflows.grc_report_id").
		Where("grc_reports.id = ?", grcReportID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("grc_report_workflow_review_comments.%s IN (?)", columnName), values)
		}
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (grc_report_workflow_review_comments.comment LIKE ? OR
            grc_reports.title LIKE ?)`,
			keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "comment" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	if params.Section != "" {
		query = query.Where("grc_report_workflows.type >= ?", params.Section)
	}

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *GRCReportWorkflowReviewCommentRepository) GetByID(id uint) (*models.GRCReportWorkflowReviewComment, error) {
	var riskComment *models.GRCReportWorkflowReviewComment
	result := repo.db.First(&riskComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskComment, nil
}

func (repo *GRCReportWorkflowReviewCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.GRCReportWorkflowReviewComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.GRCReportWorkflowReviewComment
	result := db.
		Preload("GRCReportWorkflow").
		Preload("GRCReportCommentStatus").
		Preload("GRCReportWorkflowStatus").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *GRCReportWorkflowReviewCommentRepository) Create(comment *models.GRCReportWorkflowReviewComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *GRCReportWorkflowReviewCommentRepository) Update(comment *models.GRCReportWorkflowReviewComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.GRCReportWorkflowReviewComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GRCReportWorkflowReviewCommentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", id).Delete(&models.GRCReportWorkflowReviewComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *GRCReportWorkflowReviewCommentRepository) UpdateSectionWorkflowStatus(risk *models.GRCReportWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.GRCReportWorkflow{}).Where("id = ?", risk.ID).Updates(risk).Error; err != nil {
		return err
	}

	return nil
}
