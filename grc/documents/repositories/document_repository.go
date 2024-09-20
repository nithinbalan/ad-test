package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Interfaces
type DocumentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *documentDTOs.DocumentListQueryParamsDTO, currentUserID uint) (*documentDTOs.DocumentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	Create(document *models.Document, tx ...utility.DatabaseTransaction) (uint, error)
	Update(document *models.Document, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	DeleteSharedUsers(sharedUserIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error
	DeleteSharedUnits(sharedUnitIds []*uint, tx ...utility.DatabaseTransaction) error

	GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error)
	CreateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error
	DeleteDocumentVersions(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateDocumentEditHistory(documentHistory *models.DocumentEditHistory, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteDocumentEditHistory(docHistoryEditIds []*uint, tx ...utility.DatabaseTransaction) error

	GetDocumentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	CreateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error
	CreateDocumentWorkflowParticipant(docWorkflowParticipant *models.DocumentWorkflowParticipant, tx ...utility.DatabaseTransaction) error
	DeleteDocumentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error

	GetComments(docVersionId uint, currentApprovalLevel int) (*models.DocumentVersion, error)
	GetDocumentWorkflowByLevel(docVersionId uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	GetNextLevel(documentVersionID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error)
	GetWorkflowList(docVersionId uint) (*models.DocumentVersion, error)
	DeleteDocumentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error
	CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error)
	DeleteDocumentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error

	GetSharedUsers(documentId uint, page uint, pageSize uint) ([]*models.User, commonDTOs.PaginationInfoDTO, error)
	GetReviewUsers(documentId uint) ([]*models.User, error)
	GetSharedUnits(documentId uint) ([]*models.OrganizationUnit, error)

	GetDocumentLatestEditHistory(documentVersionID uint) (*models.DocumentEditHistory, error)
	GetDocumentEditHistories(documentVersionID uint) ([]*models.DocumentEditHistory, error)
}

// Implementation
type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (repo *DocumentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentRepository) GetAll(params *documentDTOs.DocumentListQueryParamsDTO, currentUserID uint) (*documentDTOs.DocumentListResponseDTO, error) {
	var documents []*documentDTOs.DocumentListDTO
	var totalCount int64
	selectClause := `
	documents.*,
	document_categories.title_en AS document_category_title_en,
	document_categories.title_ar AS document_category_title_ar,
	document_types.title_en AS document_type_title_en,
	document_types.title_ar AS document_type_title_ar,
	document_access_types.title_en AS document_access_type_title_en,
	document_access_types.title_ar AS document_type_access_title_ar,
	organization_units.title_en AS responsible_unit_title_en,
	organization_units.title_ar AS responsible_unit_title_ar,
	responsible_user.id AS responsible_user_id,
	responsible_user.first_name_en AS responsible_user_first_name_en,
	responsible_user.first_name_ar AS responsible_user_first_name_ar,
	responsible_user.last_name_en AS responsible_user_last_name_en,
	responsible_user.last_name_ar AS responsible_user_last_name_ar,
	responsible_user.initial AS responsible_user_initial,
	responsible_user.color AS responsible_user_color,
	responsible_user.email AS responsible_user_email,
	responsible_user.profile_image_path AS responsible_user_profile_image_path,
	document_versions.url AS document_url,
	document_versions.name AS document_name,
	document_versions.size AS document_size,
	document_versions.token AS document_token,
	document_versions.original_name AS document_original_name,
	document_versions.type AS document_type,
	document_versions.version AS document_version,
	document_statuses.title_en AS document_status_title_en,
	document_statuses.title_ar AS document_status_title_ar,
	document_statuses.color_code AS document_status_color_code,
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
	published_user.first_name_en AS published_user_first_name_en,
	published_user.first_name_ar AS published_user_first_name_ar,
	published_user.last_name_en AS published_user_last_name_en,
	published_user.last_name_ar AS published_user_last_name_ar,
	published_user.initial AS published_user_initial,
	published_user.color AS published_user_color,
	published_user.email AS published_user_email,
	published_user.profile_image_path AS published_user_profile_image_path
	
`

	joinAcknowledge := false

	if params.EnableAcknowledgedStatus {
		joinAcknowledge = true
		selectClause += `, CASE WHEN document_acknowledge_users.id IS NULL THEN 0 ELSE 1 END AS acknowledged`
	}

	query := repo.db.Table("documents").Select(selectClause).
		Joins(`
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
			LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
		`)

	if joinAcknowledge {
		query = query.Joins("LEFT JOIN document_acknowledge_users ON document_acknowledge_users.document_id = documents.id AND document_acknowledge_users.user_id = ?", currentUserID)
	}

	// Filter by document types if provided
	if len(params.DocumentTypes) > 0 {
		query = query.Where("document_types.type IN (?)", params.DocumentTypes)
	}
	// Filter by document status if provided
	if len(params.DocumentStatus) > 0 {
		query = query.Where("document_statuses.type IN (?)", params.DocumentStatus)
	}
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "policy":
			query = query.Where("document_types.type = ?", "policy")
		case "procedure":
			query = query.Where("document_types.type = ?", "procedure")
		case "in_review":
			query = query.Where("document_statuses.type = ?", "in-review")
		case "draft":
			query = query.Where("document_statuses.type = ?", "draft")
		case "published":
			query = query.Where("document_statuses.type = ?", "published")
		case "acknowledged":
			if joinAcknowledge {
				query = query.Where("document_acknowledge_users.id IS NOT NULL")
			}
		case "unacknowledged":
			if joinAcknowledge {
				query = query.Where("document_acknowledge_users.id IS NULL")
			}
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
				(documents.title_en LIKE ? OR
				documents.title_ar LIKE ? OR
				documents.reference_number LIKE ? OR
				documents.description_en LIKE ? OR
				documents.description_ar LIKE ? OR
				document_categories.title_en LIKE ? OR
				document_categories.title_ar LIKE ? OR
				document_statuses.title_en LIKE ? OR
				document_statuses.title_ar LIKE ? OR
				document_types.title_en LIKE ? OR
				document_types.title_ar LIKE ? OR
				organization_units.title_en LIKE ? OR
				organization_units.title_ar LIKE ? OR
				responsible_user.first_name_en LIKE ? OR
				responsible_user.first_name_ar LIKE ? OR
				responsible_user.last_name_en LIKE ? OR
				responsible_user.last_name_ar LIKE ? OR
				responsible_user.email LIKE ? OR
				responsible_user.initial LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order

	query = query.Order(params.SortBy + " " + params.Sort)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&documents).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &documentDTOs.DocumentListResponseDTO{
			Items: documents,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &documentDTOs.DocumentListResponseDTO{
		Items: documents,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil

}

func (repo *DocumentRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var document *models.Document

	result := db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return document, nil
}

func (repo *DocumentRepository) GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document *models.Document
	result := db.
		Preload("DocumentType").
		Preload("DocumentCategory").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("DocumentAccessType").
		Preload("DocumentStatus").
		Preload("LatestDocumentVersion", "is_latest_version = ?", true).
		Preload("LatestDocumentVersion.CreatedByUser").
		Preload("LatestDocumentVersion.UpdatedByUser").
		Preload("LatestDocumentVersion.SubmittedByUser").
		Preload("DocumentVersions").
		Preload("DocumentVersions.CreatedByUser").
		Preload("DocumentVersions.UpdatedByUser").
		Preload("DocumentVersions.DocumentEditHistory").
		Preload("DocumentVersions.DocumentWorkflows").
		Preload("DocumentVersions.DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("DocumentVersions.DocumentWorkflows.DocumentWorkflowReviews").
		Preload("DocumentAuditLogs").
		Preload("DocumentSharedUsers").
		Preload("DocumentSharedUsers.User").
		Preload("DocumentSharedUsers.CreatedByUser").
		Preload("DocumentSharedUnits").
		Preload("DocumentSharedUnits.OrganizationUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&document, id)

	if result.Error != nil {
		return nil, result.Error
	}

	// If currentUserID is provided, check the acknowledgment status
	if currentUserID != nil {
		var count int64
		db.Model(&models.DocumentAcknowledgeUser{}).
			Where("document_id = ? AND user_id = ?", id, *currentUserID).
			Count(&count)
		document.Acknowledged = count > 0
	}

	return document, nil
}

func (repo *DocumentRepository) Create(document *models.Document, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *DocumentRepository) Update(document *models.Document, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(document).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete document
	if err := db.Delete(&models.Document{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documentVersion *models.DocumentVersion

	result := db.First(&documentVersion, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return documentVersion, nil
}

func (repo *DocumentRepository) CreateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentVersion).Error; err != nil {
		return 0, err
	}
	return documentVersion.ID, nil
}

func (repo *DocumentRepository) UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentVersion).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) DeleteDocumentVersions(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.DocumentVersion{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&sharedUsers).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentRepository) DeleteSharedUsers(sharedUserIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", sharedUserIds).Delete(&models.DocumentSharedUser{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&sharedUnits).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentRepository) DeleteSharedUnits(sharedUnitIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", sharedUnitIds).Delete(&models.DocumentSharedUnit{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.Document{}).
		Select("reference_number").
		Where("reference_number IS NOT NULL"). // Ensure to filter out NULL values
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "DC", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *DocumentRepository) CreateDocumentEditHistory(documentHistory *models.DocumentEditHistory, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentHistory).Error; err != nil {
		return 0, err
	}
	return documentHistory.ID, nil
}

func (repo *DocumentRepository) DeleteDocumentEditHistory(docEditHistoryIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docEditHistoryIds).Delete(&models.DocumentEditHistory{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) GetDocumentWorkflowByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documentWorkflow *models.DocumentWorkflow

	result := db.First(&documentWorkflow, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return documentWorkflow, nil
}

func (repo *DocumentRepository) CreateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentWorkflow).Error; err != nil {
		return 0, err
	}

	return documentWorkflow.ID, nil
}

func (repo *DocumentRepository) UpdateDocumentWorkflow(documentWorkflow *models.DocumentWorkflow, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentWorkflow).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) CreateDocumentWorkflowParticipant(docWorkflowParticipant *models.DocumentWorkflowParticipant, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(docWorkflowParticipant).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) GetComments(docVersionId uint, currentApprovalLevel int) (*models.DocumentVersion, error) {
	var documentVersion models.DocumentVersion
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("DocumentWorkflows", "level <= ?", currentApprovalLevel).
		Preload("DocumentWorkflows.DocumentWorkflowStatus").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants.User").
		Preload("DocumentWorkflows.DocumentWorkflowReviews").
		Preload("DocumentWorkflows.DocumentWorkflowReviews.DocumentWorkflowStatus").
		Preload("DocumentWorkflows.DocumentWorkflowReviews.CreatedByUser").
		Preload("DocumentWorkflows.DocumentWorkflowReviews.UpdatedByUser").
		Where("id = ?", docVersionId).
		Find(&documentVersion)

	if result.Error != nil {
		return nil, result.Error
	}

	return &documentVersion, nil
}

func (repo *DocumentRepository) GetDocumentWorkflowByLevel(docVersionId uint, currentApprovalLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var documentWorkflow *models.DocumentWorkflow
	result := db.
		Preload("WorkflowLevelType").
		Preload("DocumentWorkflowParticipants").
		Preload("DocumentWorkflowReviews").
		Preload("DocumentWorkflowParticipants.User").
		Where("document_version_id = ? AND level = ?", docVersionId, currentApprovalLevel).
		Find(&documentWorkflow)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentWorkflow, nil
}

func (repo *DocumentRepository) GetNextLevel(documentVersionID uint, currentLevel int, tx ...utility.DatabaseTransaction) (*models.DocumentWorkflow, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nextWorkflow models.DocumentWorkflow
	result := db.Where("document_version_id = ? AND level > ?", documentVersionID, currentLevel).
		Order("level ASC").  // Ensure we are ascending by level
		First(&nextWorkflow) // Retrieve the first result only which will be the next level

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // No further levels available
	} else if result.Error != nil {
		return nil, result.Error // Other error occurred
	}

	return &nextWorkflow, nil
}

func (repo *DocumentRepository) CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentWorkflowReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentRepository) CountApprovedParticipants(workflowID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var count int64
	err := db.Model(&models.DocumentWorkflowReview{}).
		Where("document_workflow_id = ? AND document_workflow_status_id = (SELECT id FROM document_workflow_statuses WHERE type = 'approve')", workflowID).
		Distinct("created_by_id").
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (repo *DocumentRepository) GetWorkflowList(docVersionId uint) (*models.DocumentVersion, error) {
	var documentVersion models.DocumentVersion
	result := repo.db.
		Preload("SubmittedByUser").
		Preload("DocumentWorkflows").
		Preload("DocumentWorkflows.DocumentWorkflowStatus").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants.User").
		Preload("DocumentWorkflows.DocumentWorkflowParticipants.Role").
		Preload("DocumentWorkflows.DocumentWorkflowReviews", "document_workflow_status_id IN (?)", repo.db.Model(&models.DocumentWorkflowStatus{}).Select("id").Where("type = ?", "approve")).
		Preload("DocumentWorkflows.DocumentWorkflowReviews.CreatedByUser").
		Where("id = ?", docVersionId).
		Find(&documentVersion)

	if result.Error != nil {
		return nil, result.Error
	}

	return &documentVersion, nil
}

func (repo *DocumentRepository) GetSharedUsers(documentId uint, page uint, pageSize uint) ([]*models.User, commonDTOs.PaginationInfoDTO, error) {
	var users []*models.User
	var totalCount int64

	// First, count the total eligible records
	countQuery := repo.db.Model(&models.DocumentSharedUser{}).
		Where("document_id = ? AND type = 'shared-user'", documentId)
	countQuery.Count(&totalCount)

	// Then fetch the paginated results
	offset := (page - 1) * pageSize
	err := countQuery.
		Joins("JOIN users ON users.id = document_shared_users.user_id").
		Select("users.*").
		Offset(int(offset)).
		Limit(int(pageSize)).
		Find(&users).Error

	if err != nil {
		return nil, commonDTOs.PaginationInfoDTO{}, err
	}

	pagination := commonDTOs.PaginationInfoDTO{
		CurrentPage: int(page),
		PageSize:    int(pageSize),
		TotalCount:  totalCount,
		TotalPages:  int((totalCount + int64(pageSize) - 1) / int64(pageSize)),
	}

	return users, pagination, nil
}

func (repo *DocumentRepository) GetReviewUsers(documentId uint) ([]*models.User, error) {
	var users []*models.User
	err := repo.db.
		Model(&models.DocumentSharedUser{}).
		Where("document_id = ? AND type = 'review-user'", documentId).
		Joins("JOIN users ON users.id = document_shared_users.user_id").
		Select("users.*").
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *DocumentRepository) GetSharedUnits(documentId uint) ([]*models.OrganizationUnit, error) {
	var units []*models.OrganizationUnit
	err := repo.db.
		Model(&models.DocumentSharedUnit{}).
		Where("document_id = ?", documentId).
		Joins("JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id").
		Select("organization_units.*").
		Find(&units).Error
	if err != nil {
		return nil, err
	}
	return units, nil
}

func (repo *DocumentRepository) GetDocumentLatestEditHistory(documentVersionID uint) (*models.DocumentEditHistory, error) {
	var history models.DocumentEditHistory
	err := repo.db.
		Preload("CreatedByUser").
		Where("document_version_id = ?", documentVersionID).
		Order("created_at DESC").
		First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (repo *DocumentRepository) GetDocumentEditHistories(documentVersionID uint) ([]*models.DocumentEditHistory, error) {
	var histories []*models.DocumentEditHistory
	err := repo.db.
		Where("document_version_id = ?", documentVersionID).
		Order("created_at DESC").
		Find(&histories).Error
	if err != nil {
		return nil, err
	}
	return histories, nil
}

func (repo *DocumentRepository) DeleteDocumentWorkflows(workflowIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", workflowIds).Delete(&models.DocumentWorkflow{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) DeleteDocumentWorkflowParticipants(participantIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", participantIds).Delete(&models.DocumentWorkflowParticipant{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) DeleteDocumentWorkflowReviews(reviewIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", reviewIds).Delete(&models.DocumentWorkflowReview{}).Error; err != nil {
		return err
	}

	return nil
}
