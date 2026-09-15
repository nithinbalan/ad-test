package repositories

import (
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
	GetAll(params *documentDTOs.DocumentListQueryParamsDTO, currentUserID uint, isUserSuperAdmin bool, primaryOrgUnitField string) (*documentDTOs.DocumentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	Create(document *models.Document, tx ...utility.DatabaseTransaction) (uint, error)
	Update(document *models.Document, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error)
	CreateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error
	DeleteDocumentVersions(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateDocumentEditHistory(documentHistory *models.DocumentEditHistory, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteDocumentEditHistory(docHistoryEditIds []*uint, tx ...utility.DatabaseTransaction) error

	GetDocumentLatestEditHistory(documentVersionID uint) (*models.DocumentEditHistory, error)
	GetDocumentEditHistories(documentVersionID uint) ([]*models.DocumentEditHistory, error)

	CreateDocumentWorkflowReview(documentWorkflowReview *models.DocumentWorkflowReview, tx ...utility.DatabaseTransaction) error
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

func (repo *DocumentRepository) GetAll(params *documentDTOs.DocumentListQueryParamsDTO, currentUserID uint, isUserSuperAdmin bool, primaryOrgUnitField string) (*documentDTOs.DocumentListResponseDTO, error) {
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
	confidentiality_classifications.title_en AS confidentiality_classification_title_en,
	confidentiality_classifications.title_ar AS confidentiality_classification_title_ar,
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
	document_statuses.type AS document_status_type,
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
	published_user.profile_image_path AS published_user_profile_image_path,
	ISNULL(
		(
			SELECT STRING_AGG(access_type, ',')
			FROM (
				SELECT 'created user' as access_type WHERE documents.created_by_id = ?
				UNION ALL
				SELECT 'responsible' WHERE documents.responsible_user_id = ?
				UNION ALL
				SELECT 'review user' WHERE EXISTS (
					SELECT 1 FROM document_shared_users 
					WHERE document_shared_users.document_id = documents.id 
					AND document_shared_users.user_id = ? 
					AND document_shared_users.type = 'review-user'
				)
				UNION ALL
				SELECT 'shared user' WHERE EXISTS (
					SELECT 1 FROM document_shared_users 
					WHERE document_shared_users.document_id = documents.id 
					AND document_shared_users.user_id = ? 
					AND document_shared_users.type = 'shared-user'
				)
				UNION ALL
				SELECT 'shared unit user' WHERE EXISTS (
					SELECT 1 FROM document_shared_units dsu
					JOIN organization_units shared_unit ON shared_unit.id = dsu.organization_unit_id
					JOIN users u ON u.id = ?
					WHERE dsu.document_id = documents.id 
					AND shared_unit.id = (
						SELECT ou.id FROM organization_units ou 
						WHERE ou.id = u.` + primaryOrgUnitField + `
					)
				)
				UNION ALL
				SELECT 'public' WHERE documents.document_access_type_id = (SELECT id FROM document_access_types WHERE type = 'public')
			) AS access_types
		),
		''
	) AS user_access_types
`

	joinAcknowledge := false

	if params.EnableAcknowledgedStatus {
		joinAcknowledge = true
		selectClause += `, CASE WHEN document_acknowledge_users.id IS NULL THEN 0 ELSE 1 END AS acknowledged`
	}

	query := repo.db.Table("documents").Select(selectClause,
		currentUserID, // for creator check in array
		currentUserID, // for responsible check in array
		currentUserID, // for review-user check in array
		currentUserID, // for shared-user check in array
		currentUserID, // for unit access check in array
		currentUserID, // for creator check in single
		currentUserID, // for responsible check in single
		currentUserID, // for review-user check in single
		currentUserID, // for shared-user check in single
		currentUserID, // for unit access check in single
	).
		Joins(`
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
			LEFT JOIN confidentiality_classifications ON confidentiality_classifications.id = documents.confidentiality_classification_id
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
			LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
		`).Where("documents.is_internal = ?", false)

	if joinAcknowledge {
		query = query.Joins("LEFT JOIN document_acknowledge_users ON document_acknowledge_users.document_version_id = document_versions.id AND document_acknowledge_users.user_id = ?", currentUserID)
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

	// Apply date filters if provided
	if params.DateFrom != "" || params.DateTo != "" {
		isPublishedFilter := false
		for _, filter := range params.QuickFilter {
			if filter == "published" {
				isPublishedFilter = true
				break
			}
		}
		if !isPublishedFilter {
			for _, status := range params.DocumentStatus {
				if status == "published" {
					isPublishedFilter = true
					break
				}
			}
		}

		dateField := "documents.created_at"
		if isPublishedFilter {
			dateField = "documents.published_at"
		}

		if params.DateFrom != "" && params.DateTo != "" {
			query = query.Where(dateField+" BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		} else if params.DateFrom != "" {
			query = query.Where(dateField+" >= ?", params.DateFrom+" 00:00:00")
		} else if params.DateTo != "" {
			query = query.Where(dateField+" <= ?", params.DateTo+" 23:59:59")
		}
	}

	// Apply access control restrictions if not super admin
	if !isUserSuperAdmin {
		query = query.Where(`
			(documents.document_status_id = (SELECT id FROM document_statuses WHERE type = 'in-review') 
				AND (documents.created_by_id = ? OR documents.responsible_user_id = ? 
					OR EXISTS (
						SELECT 1 FROM document_shared_users 
						WHERE document_shared_users.document_id = documents.id 
						AND document_shared_users.user_id = ? 
						AND document_shared_users.type = 'review-user'
					)
				)
			)
			OR (documents.document_status_id = (SELECT id FROM document_statuses WHERE type = 'approved') 
				AND (documents.created_by_id = ? OR documents.responsible_user_id = ? 
					OR EXISTS (
						SELECT 1 FROM document_shared_users 
						WHERE document_shared_users.document_id = documents.id 
						AND document_shared_users.user_id = ? 
						AND document_shared_users.type = 'review-user'
					)
				)
			)
			OR (documents.document_status_id = (SELECT id FROM document_statuses WHERE type = 'published') 
				AND (
					(documents.document_access_type_id = (SELECT id FROM document_access_types WHERE type = 'public'))
					OR EXISTS (
						SELECT 1 FROM document_shared_users 
						WHERE document_shared_users.document_id = documents.id 
						AND document_shared_users.user_id = ? 
						AND document_shared_users.type IN ('shared-user', 'review-user')
					)
					OR EXISTS (
						SELECT 1 FROM document_shared_units dsu
						JOIN organization_units shared_unit ON shared_unit.id = dsu.organization_unit_id
						JOIN users u ON u.id = ?
						JOIN organization_units user_unit ON user_unit.id = u.organization_unit_2_id
						WHERE dsu.document_id = documents.id 
						AND shared_unit.id = (
							SELECT ou.id FROM organization_units ou 
							WHERE ou.id = u.`+primaryOrgUnitField+`
						)
					)
					OR documents.created_by_id = ?
					OR documents.responsible_user_id = ?
				)
			)
			OR (
				documents.document_status_id NOT IN (
					SELECT id FROM document_statuses WHERE type IN ('in-review', 'published', 'approved')
				)
				AND (documents.created_by_id = ? OR documents.responsible_user_id = ?)
			)
		`, currentUserID, currentUserID, currentUserID, // in-review access
			currentUserID, currentUserID, currentUserID, // approved access
			currentUserID, currentUserID, currentUserID, currentUserID, // published access
			currentUserID, currentUserID) // all other statuses access
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
				responsible_user.initial LIKE ? OR
				documents.created_at LIKE ? OR
				document_versions.version LIKE ? OR
				document_access_types.title_en LIKE ? OR
				document_access_types.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting once
	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		// Apply pagination without duplicate ordering
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
		Preload("ConfidentialityClassification").
		Preload("DocumentStatus").
		Preload("LatestDocumentVersion", "is_latest_version = ?", true).
		Preload("LatestDocumentVersion.CreatedByUser").
		Preload("LatestDocumentVersion.UpdatedByUser").
		Preload("LatestDocumentVersion.SubmittedByUser").
		Preload("LatestDocumentVersion.PublishedByUser").
		Preload("LatestDocumentVersion.DocumentWorkflows").
		Preload("LatestDocumentVersion.DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("LatestDocumentVersion.DocumentWorkflows.DocumentWorkflowReviews").
		Preload("DocumentVersions").
		Preload("ReviewFrequency").
		Preload("DocumentVersions.CreatedByUser").
		Preload("DocumentVersions.UpdatedByUser").
		Preload("DocumentVersions.PublishedByUser").
		Preload("DocumentVersions.DocumentEditHistory").
		Preload("DocumentVersions.DocumentWorkflows").
		Preload("DocumentVersions.DocumentWorkflows.DocumentWorkflowParticipants").
		Preload("DocumentVersions.DocumentWorkflows.DocumentWorkflowParticipants.User").
		Preload("DocumentVersions.DocumentWorkflows.DocumentWorkflowReviews").
		Preload("DocumentVersions.DocumentWorkflows.DocumentWorkflowReviews.DocumentWorkflowStatus").
		Preload("DocumentAuditLogs").
		Preload("DocumentSharedUsers").
		Preload("DocumentSharedUsers.User").
		Preload("DocumentSharedUsers.CreatedByUser").
		Preload("DocumentSharedUnits").
		Preload("DocumentSharedUnits.OrganizationUnit").
		Preload("DocumentRequirementGroups").
		Preload("DocumentRequirementGroups.DocumentRequirements").
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

func (repo *DocumentRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var registerIDs []uint

	err := db.Raw(`
		WITH RecursiveRegisters AS (
			SELECT id, parent_document_id, 0 AS lvl
			FROM documents
			WHERE id = ?

			UNION ALL

			SELECT cr.id, cr.parent_document_id, rr.lvl + 1
			FROM documents cr
			JOIN RecursiveRegisters rr
				ON cr.parent_document_id = rr.id
		)
		SELECT id
		FROM RecursiveRegisters
		ORDER BY lvl DESC, id DESC
		OPTION (MAXRECURSION 32767)
	`, id).Scan(&registerIDs).Error

	if err != nil {
		return err
	}

	queries := []string{
		`DELETE audit_finding_action_plan_status_update_documents
		FROM audit_finding_action_plan_status_update_documents
		JOIN audit_finding_action_plan_status_updates ON audit_finding_action_plan_status_updates.id = audit_finding_action_plan_status_update_documents.audit_finding_action_plan_status_update_id
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_action_plan_comment_documents
		FROM audit_finding_action_plan_comment_documents
		JOIN audit_finding_action_plan_comments ON audit_finding_action_plan_comments.id = audit_finding_action_plan_comment_documents.audit_finding_action_plan_comment_id
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_action_plan_comments
		FROM audit_finding_action_plan_comments
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_draft_finding_conversation_documents
		 FROM audit_draft_finding_conversation_documents 
		 JOIN audit_draft_finding_conversations ON audit_draft_finding_conversations.id = audit_draft_finding_conversation_documents.audit_draft_finding_conversation_id
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_action_plan_status_updates
		FROM audit_finding_action_plan_status_updates
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = tasks.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_action_plans
		FROM audit_finding_action_plans
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_assets
		FROM audit_finding_assets
		JOIN audit_findings ON audit_findings.id = audit_finding_assets.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE FROM audit_finding_associated_documents
		 WHERE audit_finding_associated_documents.document_id = ?`,

		`DELETE audit_finding_associated_documents
		FROM audit_finding_associated_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_associated_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_compliance_requirements
		FROM audit_finding_compliance_requirements
		JOIN audit_findings ON audit_findings.id = audit_finding_compliance_requirements.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_controls
		FROM audit_finding_controls
		JOIN audit_findings ON audit_findings.id = audit_finding_controls.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_documents
		FROM audit_finding_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_incidents
		FROM audit_finding_incidents
		JOIN audit_findings ON audit_findings.id = audit_finding_incidents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_kpis
		FROM audit_finding_kpis
		JOIN audit_findings ON audit_findings.id = audit_finding_kpis.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_processes
		FROM audit_finding_processes
		JOIN audit_findings ON audit_findings.id = audit_finding_processes.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_risks
		FROM audit_finding_risks
		JOIN audit_findings ON audit_findings.id = audit_finding_risks.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_root_cause_analyses
		FROM audit_finding_root_cause_analyses
		JOIN audit_findings ON audit_findings.id = audit_finding_root_cause_analyses.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_finding_strategic_objectives
		FROM audit_finding_strategic_objectives
		JOIN audit_findings ON audit_findings.id = audit_finding_strategic_objectives.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_draft_finding_conversations
		 FROM audit_draft_finding_conversations 
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_draft_finding_documents
		 FROM audit_draft_finding_documents 
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_documents.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_draft_finding_action_plans
		 FROM audit_draft_finding_action_plans 
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_action_plans.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_checklist_documents
		 FROM audit_checklist_documents 
		 JOIN audit_checklists ON audit_checklists.id = audit_checklist_documents.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_draft_findings
		 FROM audit_draft_findings 
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_findings
		FROM audit_findings
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		WHERE auditable_items.document_id = ?`,

		`DELETE audit_plan_auditable_item_risk_controls
		 FROM audit_plan_auditable_item_risk_controls 
		 JOIN audit_plan_auditable_item_risks ON audit_plan_auditable_item_risks.id = audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risk_id
		 JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE document_change_request_workflow_participants
		 FROM document_change_request_workflow_participants 
		 JOIN document_change_request_workflows ON document_change_request_workflows.id = document_change_request_workflow_participants.document_change_request_workflow_id
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_workflows.document_change_request_id
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_requests.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE document_change_request_workflow_reviews
		 FROM document_change_request_workflow_reviews 
		 JOIN document_change_request_workflows ON document_change_request_workflows.id = document_change_request_workflow_reviews.document_change_request_workflow_id
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_workflows.document_change_request_id
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_requests.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE policy_exception_risks
		 FROM policy_exception_risks 
		 JOIN policy_exceptions ON policy_exceptions.id = policy_exception_risks.policy_exception_id
		 JOIN document_version_clauses ON document_version_clauses.id = policy_exceptions.document_version_clause_id
		 JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE policy_exceptions
		 FROM policy_exceptions 
		 JOIN document_version_clauses ON document_version_clauses.id = policy_exceptions.document_version_clause_id
		 JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_change_request_workflows
		 FROM document_change_request_workflows 
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_workflows.document_change_request_id
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_requests.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE document_change_request_documents
		 FROM document_change_request_documents 
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_documents.document_change_request_id
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_requests.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE document_change_request_edit_histories
		 FROM document_change_request_edit_histories 
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_edit_histories.document_change_request_id
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_requests.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE document_version_clause_comments
		 FROM document_version_clause_comments 
		 JOIN document_version_clauses ON document_version_clauses.id = document_version_clause_comments.document_version_clause_id
		 JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_version_clause_comments
		 FROM document_version_clause_comments 
		 JOIN document_versions ON document_versions.id = document_version_clause_comments.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_version_clause_histories
		 FROM document_version_clause_histories 
		 JOIN document_version_clauses ON document_version_clauses.id = document_version_clause_histories.document_version_clause_id
		 JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE audit_plan_auditable_item_risks
		 FROM audit_plan_auditable_item_risks 
		 JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_checklists
		 FROM audit_checklists 
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE document_controls
		 FROM document_controls 
		 JOIN document_version_clauses ON document_version_clauses.id = document_controls.document_version_clause_id
		 JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_workflow_participants
		 FROM document_workflow_participants 
		 JOIN document_workflows ON document_workflows.id = document_workflow_participants.document_workflow_id
		 JOIN document_versions ON document_versions.id = document_workflows.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_workflow_reviews
		 FROM document_workflow_reviews 
		 JOIN document_workflows ON document_workflows.id = document_workflow_reviews.document_workflow_id
		 JOIN document_versions ON document_versions.id = document_workflows.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_workflows
		 FROM document_workflows 
		 JOIN document_versions ON document_versions.id = document_workflows.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_edit_histories
		 FROM document_edit_histories 
		 JOIN document_versions ON document_versions.id = document_edit_histories.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE audit_document_request_auditable_items
		 FROM audit_document_request_auditable_items 
		 JOIN auditable_items ON auditable_items.id = audit_document_request_auditable_items.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE document_audit_logs
		 FROM document_audit_logs 
		 JOIN document_versions ON document_versions.id = document_audit_logs.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE FROM document_acknowledge_users 
		 WHERE document_acknowledge_users.document_id = ?`,

		`DELETE document_acknowledge_users
		 FROM document_acknowledge_users 
		 JOIN document_versions ON document_versions.id = document_acknowledge_users.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE asset_policy_checklists
		 FROM asset_policy_checklists 
		 JOIN document_requirement_groups ON document_requirement_groups.id = asset_policy_checklists.document_requirement_group_id
		 WHERE document_requirement_groups.document_id = ?`,

		`DELETE asset_policy_comments
		 FROM asset_policy_comments 
		 JOIN asset_policies ON asset_policies.id = asset_policy_comments.asset_policy_id
		 WHERE asset_policies.document_id = ?`,

		`DELETE audit_plan_auditable_items
		 FROM audit_plan_auditable_items 
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_plan_checklists
		 FROM audit_plan_checklists 
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_plan_schedule_auditable_items
		 FROM audit_plan_schedule_auditable_items 
		 JOIN auditable_items ON auditable_items.id = audit_plan_schedule_auditable_items.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE audit_schedule_auditable_items
		 FROM audit_schedule_auditable_items 
		 JOIN auditable_items ON auditable_items.id = audit_schedule_auditable_items.auditable_item_id
		 WHERE auditable_items.document_id = ?`,

		`DELETE document_version_clauses
		 FROM document_version_clauses 
		 JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id
		 WHERE document_versions.document_id = ?`,

		`DELETE document_requirements
		 FROM document_requirements 
		 JOIN document_requirement_groups ON document_requirement_groups.id = document_requirements.document_requirement_group_id
		 WHERE document_requirement_groups.document_id = ?`,

		`DELETE document_review_support_documents
		 FROM document_review_support_documents 
		 JOIN document_reviews ON document_reviews.id = document_review_support_documents.document_review_id
		 WHERE document_reviews.document_id = ?`,

		`DELETE document_change_request_initiate_documents
		 FROM document_change_request_initiate_documents 
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_request_initiate_documents.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE document_change_requests
		 FROM document_change_requests 
		 JOIN document_change_request_initiates ON document_change_request_initiates.id = document_change_requests.document_change_request_initiate_id
		 WHERE document_change_request_initiates.document_id = ?`,

		`DELETE policy_report_acknowledgement_organization_units
		 FROM policy_report_acknowledgement_organization_units 
		 JOIN policy_reports ON policy_reports.id = policy_report_acknowledgement_organization_units.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_acknowledgement_users
		 FROM policy_report_acknowledgement_users 
		 JOIN policy_reports ON policy_reports.id = policy_report_acknowledgement_users.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_clauses
		 FROM policy_report_clauses 
		 JOIN policy_reports ON policy_reports.id = policy_report_clauses.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_compliances
		 FROM policy_report_compliances 
		 JOIN policy_reports ON policy_reports.id = policy_report_compliances.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_document_accesses
		 FROM policy_report_document_accesses 
		 JOIN policy_reports ON policy_reports.id = policy_report_document_accesses.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_document_approval_histories
		 FROM policy_report_document_approval_histories 
		 JOIN policy_reports ON policy_reports.id = policy_report_document_approval_histories.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_document_attachments
		 FROM policy_report_document_attachments 
		 JOIN policy_reports ON policy_reports.id = policy_report_document_attachments.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_document_versions
		 FROM policy_report_document_versions 
		 JOIN policy_reports ON policy_reports.id = policy_report_document_versions.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_key_metrics
		 FROM policy_report_key_metrics 
		 JOIN policy_reports ON policy_reports.id = policy_report_key_metrics.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_overviews
		 FROM policy_report_overviews 
		 JOIN policy_reports ON policy_reports.id = policy_report_overviews.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_risks
		 FROM policy_report_risks 
		 JOIN policy_reports ON policy_reports.id = policy_report_risks.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE policy_report_strategic_objectives
		 FROM policy_report_strategic_objectives 
		 JOIN policy_reports ON policy_reports.id = policy_report_strategic_objectives.policy_report_id
		 WHERE policy_reports.document_id = ?`,

		`DELETE document_controls
		 FROM document_controls 
		 WHERE document_controls.document_id = ?`,

		`DELETE document_change_request_workflow_participants
		 FROM document_change_request_workflow_participants 
		 JOIN document_change_request_workflows ON document_change_request_workflows.id = document_change_request_workflow_participants.document_change_request_workflow_id
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_workflows.document_change_request_id
		 WHERE document_change_requests.document_id = ?`,

		`DELETE document_change_request_workflow_reviews
		 FROM document_change_request_workflow_reviews 
		 JOIN document_change_request_workflows ON document_change_request_workflows.id = document_change_request_workflow_reviews.document_change_request_workflow_id
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_workflows.document_change_request_id
		 WHERE document_change_requests.document_id = ?`,

		`DELETE document_change_request_workflows
		 FROM document_change_request_workflows 
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_workflows.document_change_request_id
		 WHERE document_change_requests.document_id = ?`,

		`DELETE document_change_request_documents
		 FROM document_change_request_documents 
		 JOIN document_change_requests ON document_change_requests.id = document_change_request_documents.document_change_request_id
		 WHERE document_change_requests.document_id = ?`,

		`DELETE document_change_request_edit_histories
		FROM document_change_request_edit_histories 
		JOIN document_change_requests ON document_change_requests.id = document_change_request_edit_histories.document_change_request_id
		WHERE document_change_requests.document_id = ?`,

		`DELETE FROM document_change_requests WHERE document_change_requests.document_id = ?`,

		`DELETE FROM policy_reports WHERE policy_reports.document_id = ?`,
		`DELETE FROM process_associated_documents WHERE process_associated_documents.document_id = ?`,
		`DELETE FROM project_associated_documents WHERE project_associated_documents.document_id = ?`,
		`DELETE FROM risk_associated_documents WHERE risk_associated_documents.document_id = ?`,
		`DELETE FROM sla_associated_documents WHERE sla_associated_documents.document_id = ?`,
		`DELETE FROM strategic_objective_associated_documents WHERE strategic_objective_associated_documents.document_id = ?`,
		`DELETE FROM document_compliance_requirements WHERE document_compliance_requirements.document_id = ?`,
		`DELETE FROM document_change_request_initiates WHERE document_change_request_initiates.document_id = ?`,
		`DELETE FROM document_reviews WHERE document_reviews.document_id = ?`,
		`DELETE FROM document_shared_units WHERE document_shared_units.document_id = ?`,
		`DELETE FROM document_shared_users WHERE document_shared_users.document_id = ?`,
		`DELETE FROM document_versions WHERE document_versions.document_id = ?`,
		`DELETE FROM auditable_items WHERE auditable_items.document_id = ?`,
		`DELETE FROM bcp_associated_documents WHERE bcp_associated_documents.document_id = ?`,
		`DELETE FROM business_impact_assessment_recovery_objectives WHERE business_impact_assessment_recovery_objectives.document_id = ?`,
		`DELETE FROM compliance_requirement_associated_documents WHERE compliance_requirement_associated_documents.document_id = ?`,
		`DELETE FROM control_associated_documents WHERE control_associated_documents.document_id = ?`,
		`DELETE FROM asset_associated_documents WHERE asset_associated_documents.document_id = ?`,
		`DELETE FROM asset_policies WHERE asset_policies.document_id = ?`,
		`DELETE FROM document_associated_documents WHERE document_associated_documents.associated_document_id = ?`,
		`DELETE FROM document_associated_documents WHERE document_associated_documents.document_id = ?`,
		`DELETE FROM document_requirement_groups WHERE document_requirement_groups.document_id = ?`,
		`DELETE FROM incident_associated_documents WHERE incident_associated_documents.document_id = ?`,
		`DELETE FROM kpi_associated_documents WHERE kpi_associated_documents.document_id = ?`,
		`DELETE FROM kri_associated_documents WHERE kri_associated_documents.document_id = ?`,
		`DELETE FROM meeting_associated_documents WHERE meeting_associated_documents.document_id = ?`,
		`DELETE FROM third_party_associated_documents WHERE third_party_associated_documents.document_id = ?`,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, rid := range registerIDs {
			for i, q := range queries {
				if err := tx.Exec(q, rid).Error; err != nil {
					return fmt.Errorf("force delete failed for document %d at step %d: %w", rid, i, err)
				}
			}
			if err := tx.Delete(&models.Document{}, rid).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
