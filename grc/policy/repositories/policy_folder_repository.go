package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type PolicyFolderRepositoryInterface interface {
	GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO, currentUserID uint, isUserSuperAdmin bool, primaryOrgUnitField string) (*policyDTOs.PolicyRegisterListResponseDTO, error)
	GetByID(id uint) (*models.Document, error)
	Create(folder *models.Document) (*models.Document, error)
	Update(id uint, folder *models.Document) (*models.Document, error)
	Delete(id uint) error
	ExistsByTitle(titleEn string, titleAr string, parentDocumentID *uint) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string, parentDocumentID *uint) (bool, error)
}

// Implementation
type PolicyFolderRepository struct {
	db *gorm.DB
}

func NewPolicyFolderRepository(db *gorm.DB) *PolicyFolderRepository {
	return &PolicyFolderRepository{
		db: db,
	}
}

func (repo *PolicyFolderRepository) GetAll(params *policyDTOs.PolicyRegisterListQueryParamsDTO, currentUserID uint, isUserSuperAdmin bool, primaryOrgUnitField string) (*policyDTOs.PolicyRegisterListResponseDTO, error) {
	var documents []*policyDTOs.PolicyRegisterListDTO
	var totalCount int64

	selectClause := `
	documents.*,
	document_categories.title_en AS document_category_title_en,
	document_categories.title_ar AS document_category_title_ar,
	document_types.title_en AS document_type_title_en,
	document_types.title_ar AS document_type_title_ar,
	document_access_types.title_en AS document_access_type_title_en,
	document_access_types.title_ar AS document_access_type_title_ar,
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
	) AS user_access_types,
	(SELECT COUNT(*) FROM documents d2 WHERE d2.parent_document_id = documents.id) AS document_count
`

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
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
			LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
		`).Where("documents.is_internal = ? AND (document_types.type IN (?) OR documents.is_folder = ?)", true, []string{"policy", "procedure"}, true)

	_, ok := params.Filters["parent_document_id"]
	if !ok {
		query.Where("documents.parent_document_id IS NULL")
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
			OR (documents.document_status_id = (SELECT id FROM document_statuses WHERE type = 'rejected') 
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
			OR (
				documents.is_folder = ?
			)
		`, currentUserID, currentUserID, currentUserID, // in-review access
			currentUserID, currentUserID, currentUserID, //reject case
			currentUserID, currentUserID, currentUserID, // approved access
			currentUserID, currentUserID, currentUserID, currentUserID, // published access
			currentUserID, currentUserID, true) // all other statuses access
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

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting once
	if params.SortBy == "" {
		query = query.Order("documents.id DESC")
	} else {
		query = query.Order("documents.is_folder DESC").Order(params.SortBy + " " + params.Sort)
	}

	if !params.Export {
		// Apply pagination without duplicate ordering
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&documents).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &policyDTOs.PolicyRegisterListResponseDTO{
			Items: documents,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &policyDTOs.PolicyRegisterListResponseDTO{
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

func (repo *PolicyFolderRepository) GetByID(id uint) (*models.Document, error) {
	var folder *models.Document
	result := repo.db.First(&folder, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return folder, nil
}

func (repo *PolicyFolderRepository) Create(folder *models.Document) (*models.Document, error) {
	result := repo.db.Create(&folder)

	return folder, result.Error
}

func (repo *PolicyFolderRepository) Update(id uint, updatedFolder *models.Document) (*models.Document, error) {
	var folder *models.Document

	if err := repo.db.First(&folder, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&folder).Updates(updatedFolder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (repo *PolicyFolderRepository) Delete(id uint) error {
	var folder *models.Document
	if err := repo.db.First(&folder, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Document{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PolicyFolderRepository) ExistsByTitle(titleEn, titleAr string, parentDocumentID *uint) (bool, error) {
	var count int64

	query := repo.db.Model(&models.Document{})

	// If parentDocumentID is provided, filter by it
	if parentDocumentID != nil {
		query = query.Where("parent_document_id = ?", *parentDocumentID)
	} else {
		query = query.Where("parent_document_id IS NULL")
	}

	// Check for English or Arabic title
	err := query.Where("title_en = ? OR title_ar = ?", titleEn, titleAr).Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *PolicyFolderRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string, parentDocumentID *uint) (bool, error) {
	var count int64

	query := repo.db.Model(&models.Document{}).Where("id <> ?", id)

	// Apply parent folder condition
	if parentDocumentID != nil {
		query = query.Where("parent_document_id = ?", *parentDocumentID)
	} else {
		query = query.Where("parent_document_id IS NULL")
	}

	// Check for duplicate titles (either English or Arabic)
	err := query.Where("title_en = ? OR title_ar = ?", titleEn, titleAr).Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
