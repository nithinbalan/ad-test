package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ProjectChangeRequestRepositoryInterface interface {
	GetAll(projectID uint, params *dtos.ProjectChangeRequestListQueryParamsDTO) (*dtos.ProjectChangeRequestListResponseDTO, error)
	GetByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequest, error)
	GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequest, error)
	GenerateReferenceNumber() (string, error)
	Create(changeRequest *models.ProjectChangeRequest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(changeRequest *models.ProjectChangeRequest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeleteImpactAssessments(changeRequestID uint, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentsByChangeRequestID(changeRequestID uint, tx ...utility.DatabaseTransaction) error
	DeleteWorkflowHistories(changeRequestID uint, tx ...utility.DatabaseTransaction) error
	DeleteReviewComments(changeRequestID uint, tx ...utility.DatabaseTransaction) error
	CreateImpactAssessments(impactAssessments []*models.ProjectImpactAssessment, tx ...utility.DatabaseTransaction) error
	GetImpactAssessments(changeRequestID uint, tx ...utility.DatabaseTransaction) ([]*models.ProjectImpactAssessment, error)
	CreateDocuments(documents []*models.ProjectChangeRequestDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestDocument, error)
	GetSummary(projectID uint) (*dtos.ProjectChangeRequestSummaryDTO, error)
	BeginTransaction() utility.DatabaseTransaction
}

type ProjectChangeRequestRepository struct {
	db *gorm.DB
}

func NewProjectChangeRequestRepository(db *gorm.DB) *ProjectChangeRequestRepository {
	return &ProjectChangeRequestRepository{db: db}
}

func (pcrr *ProjectChangeRequestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pcrr.db)
}

func (pcrr *ProjectChangeRequestRepository) GetAll(projectID uint, params *dtos.ProjectChangeRequestListQueryParamsDTO) (*dtos.ProjectChangeRequestListResponseDTO, error) {
	var changeRequests []*dtos.ProjectChangeRequestListDTO
	var totalCount int64

	query := pcrr.db.Table("project_change_requests").
		Select(`
			project_change_requests.id,
			project_change_requests.reference_number,
			project_change_requests.project_id,
			project_change_requests.title,
			project_change_requests.description,
			project_change_requests.business_justification,
			project_change_requests.anticipated_benefits,
			project_change_requests.risks_of_deferral,
			COALESCE(cr_impacts.impact, '') AS impact,
			project_change_requests.requested_by_user_id,
			requested_user.first_name_en AS requested_by_user_first_name_en,
			requested_user.first_name_ar AS requested_by_user_first_name_ar,
			requested_user.last_name_en AS requested_by_user_last_name_en,
			requested_user.last_name_ar AS requested_by_user_last_name_ar,
			requested_user.initial AS requested_by_user_initial,
			requested_user.color AS requested_by_user_color,
			requested_user.email AS requested_by_user_email,
			requested_user.profile_image_path AS requested_by_user_image_path,
			CONVERT(VARCHAR, project_change_requests.requested_on, 23) AS requested_on,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			project_change_requests.workflow_status_id AS workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			project_change_requests.next_review_level AS next_review_level,
			CONVERT(VARCHAR, project_change_requests.submitted_at, 23) AS submitted_at,
			submitted_user.id AS submitted_user_id,
			submitted_user.first_name_en AS submitted_user_first_name_en,
			submitted_user.first_name_ar AS submitted_user_first_name_ar,
			submitted_user.last_name_en AS submitted_user_last_name_en,
			submitted_user.last_name_ar AS submitted_user_last_name_ar,
			CONVERT(VARCHAR, project_change_requests.approved_at, 23) AS approved_at,
			approved_user.id AS approved_user_id,
			approved_user.first_name_en AS approved_user_first_name_en,
			approved_user.first_name_ar AS approved_user_first_name_ar,
			approved_user.last_name_en AS approved_user_last_name_en,
			approved_user.last_name_ar AS approved_user_last_name_ar,
			CONVERT(VARCHAR, project_change_requests.rejected_at, 23) AS rejected_at,
			rejected_user.id AS rejected_user_id,
			rejected_user.first_name_en AS rejected_user_first_name_en,
			rejected_user.first_name_ar AS rejected_user_first_name_ar,
			rejected_user.last_name_en AS rejected_user_last_name_en,
			rejected_user.last_name_ar AS rejected_user_last_name_ar,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`LEFT JOIN projects ON project_change_requests.project_id = projects.id`).
		Joins(`LEFT JOIN workflow_statuses ON workflow_statuses.id = project_change_requests.workflow_status_id`).
		Joins(`LEFT JOIN (
			SELECT 
				pcia.project_change_request_id,
				STRING_AGG(CAST(pcia.title_en AS NVARCHAR(MAX)), ', ') AS impact
			FROM (
				SELECT DISTINCT 
					ia.project_change_request_id,
					piat.title_en
				FROM project_change_request_impact_assessments ia
				JOIN project_impact_assessment_types piat ON piat.id = ia.project_impact_assessment_type_id
			) pcia
			GROUP BY pcia.project_change_request_id
		) cr_impacts ON cr_impacts.project_change_request_id = project_change_requests.id`).
		Joins(`LEFT JOIN users AS requested_user ON project_change_requests.requested_by_user_id = requested_user.id`).
		Joins(`LEFT JOIN users AS submitted_user ON project_change_requests.submitted_by_id = submitted_user.id`).
		Joins(`LEFT JOIN users AS approved_user ON project_change_requests.approved_by_id = approved_user.id`).
		Joins(`LEFT JOIN users AS rejected_user ON project_change_requests.rejected_by_id = rejected_user.id`).
		Joins(`LEFT JOIN users AS created_user ON project_change_requests.created_by_id = created_user.id`).
		Joins(`LEFT JOIN users AS updated_user ON project_change_requests.updated_by_id = updated_user.id`)

	// Filter by project_id if provided
	if projectID != 0 {
		query = query.Where("project_change_requests.project_id = ?", projectID)
	}

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_change_requests."+columnName+" IN (?)", values)
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
			query = query.Where("project_change_requests.workflow_status_id IS NULL")
		}
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_change_requests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply search keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(project_change_requests.title LIKE ? OR
			project_change_requests.reference_number LIKE ? OR
			project_change_requests.description LIKE ? OR
			project_change_requests.business_justification LIKE ? OR
			cr_impacts.impact LIKE ?)
		`, keyword, keyword, keyword, keyword, keyword)
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&changeRequests).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.ProjectChangeRequestListResponseDTO{
		Items: changeRequests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (pcrr *ProjectChangeRequestRepository) GetByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequest, error) {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var changeRequest models.ProjectChangeRequest
	result := db.Where("project_id = ? AND id = ?", projectID, id).First(&changeRequest)
	if result.Error != nil {
		return nil, result.Error
	}
	return &changeRequest, nil
}

func (pcrr *ProjectChangeRequestRepository) GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequest, error) {
	var changeRequest models.ProjectChangeRequest
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.
		Preload("Project").
		Preload("RequestedByUser").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("RejectedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ImpactAssessments").
		Preload("ImpactAssessments.ProjectImpactAssessmentType").
		Preload("Documents").
		Where("project_id = ? AND id = ?", projectID, id).
		First(&changeRequest)
	if result.Error != nil {
		return nil, result.Error
	}
	return &changeRequest, nil
}

func (pcrr *ProjectChangeRequestRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := pcrr.db.Model(&models.ProjectChangeRequest{}).
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

	return fmt.Sprintf("PCR-%d-%04d", currentYear, lastReferenceID+1), nil
}

func (pcrr *ProjectChangeRequestRepository) Create(changeRequest *models.ProjectChangeRequest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(changeRequest).Error; err != nil {
		return 0, err
	}
	return changeRequest.ID, nil
}

func (pcrr *ProjectChangeRequestRepository) Update(changeRequest *models.ProjectChangeRequest, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ProjectChangeRequest{}).
		Where("id = ?", changeRequest.ID).
		Select("title", "description", "business_justification", "anticipated_benefits", "risks_of_deferral", "requested_by_user_id", "requested_on", "updated_by_id", "updated_at").
		Updates(changeRequest).Error
}

func (pcrr *ProjectChangeRequestRepository) CreateImpactAssessments(impactAssessments []*models.ProjectImpactAssessment, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&impactAssessments).Error
}

func (pcrr *ProjectChangeRequestRepository) DeleteImpactAssessments(changeRequestID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_change_request_id = ?", changeRequestID).Delete(&models.ProjectImpactAssessment{}).Error
}

func (pcrr *ProjectChangeRequestRepository) CreateDocuments(documents []*models.ProjectChangeRequestDocument, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (pcrr *ProjectChangeRequestRepository) GetDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.ProjectChangeRequestDocument, error) {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document models.ProjectChangeRequestDocument
	result := db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &document, nil
}

func (pcrr *ProjectChangeRequestRepository) DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", ids).Delete(&models.ProjectChangeRequestDocument{}).Error
}

func (pcrr *ProjectChangeRequestRepository) DeleteDocumentsByChangeRequestID(changeRequestID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_change_request_id = ?", changeRequestID).Delete(&models.ProjectChangeRequestDocument{}).Error
}

func (pcrr *ProjectChangeRequestRepository) DeleteWorkflowHistories(changeRequestID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_change_request_id = ?", changeRequestID).Delete(&models.ProjectChangeRequestWorkflowHistory{}).Error
}

func (pcrr *ProjectChangeRequestRepository) DeleteReviewComments(changeRequestID uint, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var commentIDs []uint
	if err := db.Model(&models.ProjectChangeRequestReviewComment{}).
		Where("project_change_request_id = ?", changeRequestID).
		Pluck("id", &commentIDs).Error; err != nil {
		return err
	}
	if len(commentIDs) > 0 {
		if err := db.Where("project_change_request_review_comment_id IN (?)", commentIDs).
			Delete(&models.ProjectChangeRequestReviewCommentDocument{}).Error; err != nil {
			return err
		}
		if err := db.Where("project_change_request_id = ?", changeRequestID).
			Delete(&models.ProjectChangeRequestReviewComment{}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (pcrr *ProjectChangeRequestRepository) GetImpactAssessments(changeRequestID uint, tx ...utility.DatabaseTransaction) ([]*models.ProjectImpactAssessment, error) {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var impactAssessments []*models.ProjectImpactAssessment
	err := db.
		Preload("ProjectImpactAssessmentType").
		Where("project_change_request_id = ?", changeRequestID).
		Find(&impactAssessments).Error
	if err != nil {
		return nil, err
	}
	return impactAssessments, nil
}

func (pcrr *ProjectChangeRequestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := pcrr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ProjectChangeRequest{}, id).Error
}

func (pcrr *ProjectChangeRequestRepository) GetSummary(projectID uint) (*dtos.ProjectChangeRequestSummaryDTO, error) {
	var summary dtos.ProjectChangeRequestSummaryDTO

	query := pcrr.db.Table("project_change_requests").
		Select(`
			COUNT(*) AS total,
			COUNT(CASE WHEN workflow_statuses.type = 'submit' THEN 1 END) AS submitted,
			COUNT(CASE WHEN workflow_statuses.type = 'approve' THEN 1 END) AS approved,
			COUNT(CASE WHEN workflow_statuses.type = 'reject' THEN 1 END) AS rejected,
			COUNT(CASE WHEN workflow_statuses.type = 'pending' THEN 1 END) AS pending
		`).
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = project_change_requests.workflow_status_id")

	if projectID != 0 {
		query = query.Where("project_change_requests.project_id = ?", projectID)
	}

	if err := query.Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}
