package repositories

import (
	"fmt"
	"math"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectDecisionLogRepositoryInterface interface {
	GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error)
	GetFullDetailsByID(projectID uint, id uint) (*models.ProjectDecisionLog, error)
	GenerateReferenceNumber() (string, error)
	Create(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) (uint, error)
	Update(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowHistories(decisionLogID uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewComments(decisionLogID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.ProjectDecisionLogDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentsByDecisionLogID(decisionLogID uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ProjectDecisionLogDocument, error)
	GetSummary(projectID uint) (*dtos.ProjectDecisionLogSummaryDTO, error)
	BeginTransaction() utility.DatabaseTransaction
}

type ProjectDecisionLogRepository struct {
	db *gorm.DB
}

func NewProjectDecisionLogRepository(db *gorm.DB) *ProjectDecisionLogRepository {
	return &ProjectDecisionLogRepository{db: db}
}

func (pdlr *ProjectDecisionLogRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pdlr.db)
}

func (pdlr *ProjectDecisionLogRepository) GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error) {
	var decisionLogs []*dtos.ProjectDecisionLogListDTO
	var totalCount int64

	query := pdlr.db.Table("project_decision_logs").
		Select(`
			project_decision_logs.id,
			project_decision_logs.reference_number,
			project_decision_logs.project_id,
			project_decision_logs.decision,
			project_decision_logs.decision_notes,
			project_decision_logs.reason,
			CONVERT(VARCHAR, project_decision_logs.decision_date, 23) AS decision_date,
			project_decision_logs.requested_by_user_id,
			requested_user.first_name_en AS requested_by_user_first_name_en,
			requested_user.first_name_ar AS requested_by_user_first_name_ar,
			requested_user.last_name_en AS requested_by_user_last_name_en,
			requested_user.last_name_ar AS requested_by_user_last_name_ar,
			requested_user.initial AS requested_by_user_initial,
			requested_user.color AS requested_by_user_color,
			requested_user.email AS requested_by_user_email,
			requested_user.profile_image_path AS requested_by_user_profile_image_path,
			project_decision_logs.workflow_status_id AS workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			project_decision_logs.next_review_level AS next_review_level,
			CONVERT(VARCHAR, project_decision_logs.submitted_at, 23) AS submitted_at,
			submitted_user.id AS submitted_user_id,
			submitted_user.first_name_en AS submitted_user_first_name_en,
			submitted_user.first_name_ar AS submitted_user_first_name_ar,
			submitted_user.last_name_en AS submitted_user_last_name_en,
			submitted_user.last_name_ar AS submitted_user_last_name_ar,
			CONVERT(VARCHAR, project_decision_logs.approved_at, 23) AS approved_at,
			approved_user.id AS approved_user_id,
			approved_user.first_name_en AS approved_user_first_name_en,
			approved_user.first_name_ar AS approved_user_first_name_ar,
			approved_user.last_name_en AS approved_user_last_name_en,
			approved_user.last_name_ar AS approved_user_last_name_ar,
			CONVERT(VARCHAR, project_decision_logs.rejected_at, 23) AS rejected_at,
			rejected_user.id AS rejected_user_id,
			rejected_user.first_name_en AS rejected_user_first_name_en,
			rejected_user.first_name_ar AS rejected_user_first_name_ar,
			rejected_user.last_name_en AS rejected_user_last_name_en,
			rejected_user.last_name_ar AS rejected_user_last_name_ar,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, project_decision_logs.created_at, 120) AS created_at,
			CONVERT(VARCHAR, project_decision_logs.updated_at, 120) AS updated_at
		`).
		Joins(`LEFT JOIN projects ON project_decision_logs.project_id = projects.id`).
		Joins(`LEFT JOIN workflow_statuses ON workflow_statuses.id = project_decision_logs.workflow_status_id`).
		Joins(`LEFT JOIN users AS requested_user ON project_decision_logs.requested_by_user_id = requested_user.id`).
		Joins(`LEFT JOIN users AS submitted_user ON project_decision_logs.submitted_by_id = submitted_user.id`).
		Joins(`LEFT JOIN users AS approved_user ON project_decision_logs.approved_by_id = approved_user.id`).
		Joins(`LEFT JOIN users AS rejected_user ON project_decision_logs.rejected_by_id = rejected_user.id`).
		Joins(`LEFT JOIN users AS created_user ON project_decision_logs.created_by_id = created_user.id`).
		Joins(`LEFT JOIN users AS updated_user ON project_decision_logs.updated_by_id = updated_user.id`)

	if projectID != 0 {
		query = query.Where("project_decision_logs.project_id = ?", projectID)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_decision_logs."+columnName+" IN (?)", values)
		}
	}
	for _, filter := range params.QuickFilter {
		switch filter {
		case "pending":
			query = query.Where("workflow_statuses.type = ?", "pending")
		case "approved":
			query = query.Where("workflow_statuses.type = ?", "approve")
		case "rejected":
			query = query.Where("workflow_statuses.type = ?", "reject")
		case "reverted":
			query = query.Where("workflow_statuses.type = ?", "revert")
		case "draft":
			query = query.Where("project_decision_logs.workflow_status_id IS NULL")
		}
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_decision_logs.decision_date BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(project_decision_logs.reference_number LIKE ? OR
			project_decision_logs.decision LIKE ? OR
			project_decision_logs.decision_notes LIKE ? OR
			project_decision_logs.reason LIKE ?)
		`, keyword, keyword, keyword, keyword)
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&decisionLogs).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))

	response := &dtos.ProjectDecisionLogListResponseDTO{
		Data: decisionLogs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (pdlr *ProjectDecisionLogRepository) GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error) {
	var decisionLog models.ProjectDecisionLog
	result := pdlr.db.Where("project_id = ? AND id = ?", projectID, id).First(&decisionLog)
	if result.Error != nil {
		return nil, result.Error
	}
	return &decisionLog, nil
}

func (pdlr *ProjectDecisionLogRepository) GetFullDetailsByID(projectID uint, id uint) (*models.ProjectDecisionLog, error) {
	var decisionLog models.ProjectDecisionLog
	result := pdlr.db.
		Preload("Project").
		Preload("RequestedByUser").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("RejectedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("project_id = ? AND id = ?", projectID, id).
		First(&decisionLog)
	if result.Error != nil {
		return nil, result.Error
	}
	return &decisionLog, nil
}

func (pdlr *ProjectDecisionLogRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := pdlr.db.Model(&models.ProjectDecisionLog{}).
		Select("reference_number").
		Where("reference_number IS NOT NULL").
		Order("reference_number DESC").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("DL-%d-%04d", currentYear, lastReferenceID+1), nil
}

func (pdlr *ProjectDecisionLogRepository) Create(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Create(decisionLog)
	if result.Error != nil {
		return 0, result.Error
	}
	return decisionLog.ID, nil
}

func (pdlr *ProjectDecisionLogRepository) Update(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(decisionLog).Error
}

func (pdlr *ProjectDecisionLogRepository) DeleteWorkflowHistories(decisionLogID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_decision_log_id = ?", decisionLogID).Delete(&models.ProjectDecisionLogWorkflowHistory{}).Error
}

func (pdlr *ProjectDecisionLogRepository) DeleteReviewComments(decisionLogID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var commentIDs []uint
	if err := db.Model(&models.ProjectDecisionLogReviewComment{}).
		Where("project_decision_log_id = ?", decisionLogID).
		Pluck("id", &commentIDs).Error; err != nil {
		return err
	}
	if len(commentIDs) > 0 {
		if err := db.Where("project_decision_log_review_comment_id IN (?)", commentIDs).
			Delete(&models.ProjectDecisionLogReviewCommentDocument{}).Error; err != nil {
			return err
		}
		if err := db.Where("project_decision_log_id = ?", decisionLogID).
			Delete(&models.ProjectDecisionLogReviewComment{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (pdlr *ProjectDecisionLogRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ProjectDecisionLog{}, id).Error
}

func (pdlr *ProjectDecisionLogRepository) CreateDocuments(documents []*models.ProjectDecisionLogDocument, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (pdlr *ProjectDecisionLogRepository) DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", ids).Delete(&models.ProjectDecisionLogDocument{}).Error
}

func (pdlr *ProjectDecisionLogRepository) DeleteDocumentsByDecisionLogID(decisionLogID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_decision_log_id = ?", decisionLogID).Delete(&models.ProjectDecisionLogDocument{}).Error
}

func (pdlr *ProjectDecisionLogRepository) GetDocumentByID(id uint) (*models.ProjectDecisionLogDocument, error) {
	var document models.ProjectDecisionLogDocument
	result := pdlr.db.Where("id = ?", id).First(&document)
	if result.Error != nil {
		return nil, result.Error
	}
	return &document, nil
}

func (pdlr *ProjectDecisionLogRepository) GetSummary(projectID uint) (*dtos.ProjectDecisionLogSummaryDTO, error) {
	var summary dtos.ProjectDecisionLogSummaryDTO

	query := pdlr.db.Table("project_decision_logs").
		Select(`
			COUNT(*) AS total,
			COUNT(CASE WHEN workflow_statuses.type = 'submit' THEN 1 END) AS submitted,
			COUNT(CASE WHEN workflow_statuses.type = 'approve' THEN 1 END) AS approved,
			COUNT(CASE WHEN workflow_statuses.type = 'reject' THEN 1 END) AS rejected,
			COUNT(CASE WHEN workflow_statuses.type = 'pending' THEN 1 END) AS pending
		`).
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = project_decision_logs.workflow_status_id")

	if projectID != 0 {
		query = query.Where("project_decision_logs.project_id = ?", projectID)
	}

	if err := query.Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}
