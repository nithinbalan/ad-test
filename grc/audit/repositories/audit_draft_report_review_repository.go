package repositories

import (
	"jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditDraftReportReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditDraftReportID uint, params *dtos.AuditReportReviewListQueryParamsDTO) ([]*dtos.AuditReportReviewListDTO, error)
	CreateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) error
	Create(auditReportReview *models.AuditReportReview, tx ...utility.DatabaseTransaction) (uint, error)
	CreateAuditReportWorkflowParticipants(participants []*models.AuditReportWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteAuditReportWorkflowParticipants(participants []*uint, tx ...utility.DatabaseTransaction) error
	GetAuditReportWorkflowByDraftReportID(auditDraftReportID uint, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error)
}

type AuditDraftReportReviewRepository struct {
	db *gorm.DB
}

func NewAuditDraftReportReviewRepository(db *gorm.DB) *AuditDraftReportReviewRepository {
	return &AuditDraftReportReviewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditDraftReportReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditDraftReportReviewRepository) GetAll(auditDraftReportID uint, params *dtos.AuditReportReviewListQueryParamsDTO) ([]*dtos.AuditReportReviewListDTO, error) {
	var reviews []*dtos.AuditReportReviewListDTO

	query := repo.db.Table("audit_report_reviews").
		Select(`
			audit_report_reviews.id,
			audit_report_reviews.audit_draft_report_id,
			audit_report_reviews.audit_report_workflow_id,
			audit_report_reviews.workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			audit_report_reviews.user_id,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.email AS review_user_email,
			review_user.initial AS review_user_initial,
        	review_user.color AS review_user_color,
			review_user.profile_image_path AS review_user_profile_image_path,
			audit_report_reviews.comment,
			audit_report_reviews.created_by_id,
			audit_report_reviews.created_at,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN users AS created_user ON created_user.id = audit_report_reviews.created_by_id
			LEFT JOIN users AS review_user ON review_user.id = audit_report_reviews.user_id
			LEFT JOIN workflow_statuses ON workflow_statuses.id = audit_report_reviews.workflow_status_id`).
		Where(`audit_report_reviews.audit_draft_report_id = ?`, auditDraftReportID)

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(audit_report_reviews.comment LIKE ? )`,
			"%"+params.Keyword+"%")
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&reviews)

	return reviews, result.Error
}

func (repo *AuditDraftReportReviewRepository) CreateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReportWorkflow).Error; err != nil {
		return 0, err
	}

	return auditReportWorkflow.ID, nil
}
func (repo *AuditDraftReportReviewRepository) UpdateAuditReportWorkflow(auditReportWorkflow *models.AuditReportWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditReportWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportReviewRepository) Create(auditReportReview *models.AuditReportReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditReportReview).Error; err != nil {
		return 0, err
	}

	return auditReportReview.ID, nil
}

func (repo *AuditDraftReportReviewRepository) CreateAuditReportWorkflowParticipants(participants []*models.AuditReportWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&participants).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportReviewRepository) DeleteAuditReportWorkflowParticipants(participants []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participants).Delete(&models.AuditReportWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditDraftReportReviewRepository) GetAuditReportWorkflowByDraftReportID(auditDraftReportID uint, tx ...utility.DatabaseTransaction) (*models.AuditReportWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditReportWorkflow models.AuditReportWorkflow
	// Query to find the audit report workflow by audit draft report ID
	result := db.
		Preload("AuditReportWorkflowParticipants").
		Preload("CreatedByUser").
		Where("audit_draft_report_id = ?", auditDraftReportID).
		First(&auditReportWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}
	return &auditReportWorkflow, nil
}
