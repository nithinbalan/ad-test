package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DocumentAcknowledgeUserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(documentID uint, primaryOrgUnitField string, params *documentDTOs.DocumentAcknowledgeListQueryParamsDTO) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error)
	GetAllForPublicDocument(documentID uint, primaryOrgUnitField string, params *documentDTOs.DocumentAcknowledgeListQueryParamsDTO) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error)
	Create(documentAcknowledgeUser *models.DocumentAcknowledgeUser, tx ...utility.DatabaseTransaction) (*models.DocumentAcknowledgeUser, error)
	Exists(documentID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error)
	Delete(documentID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	GetDocumentAcknowledgmentSummary(documentID uint, primaryOrgUnitField string) (*documentDTOs.DocumentAcknowledgmentSummaryDTO, error)
	GetPublicDocumentAcknowledgmentSummary(documentID uint) (*documentDTOs.DocumentAcknowledgmentSummaryDTO, error)
}

// Implementation
type DocumentAcknowledgeUserRepository struct {
	db *gorm.DB
}

func NewDocumentAcknowledgeUserRepository(db *gorm.DB) *DocumentAcknowledgeUserRepository {
	return &DocumentAcknowledgeUserRepository{
		db: db,
	}
}

func (repo *DocumentAcknowledgeUserRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentAcknowledgeUserRepository) GetAll(documentID uint, primaryOrgUnitField string, params *documentDTOs.DocumentAcknowledgeListQueryParamsDTO) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error) {
	var acknowledgementList []documentDTOs.DocumentAcknowledgementListDTO
	var totalCount int64

	// Prepare the keyword search filter
	keywordFilter := ""
	if params.Keyword != "" {
		keywordFilter = fmt.Sprintf(`
            AND (users.first_name_en LIKE '%%%s%%' 
            OR users.first_name_ar LIKE '%%%s%%' 
            OR users.last_name_en LIKE '%%%s%%' 
            OR users.last_name_ar LIKE '%%%s%%' 
            OR users.email LIKE '%%%s%%'
            OR organization_units.title_en LIKE '%%%s%%'
            OR organization_units.title_ar LIKE '%%%s%%')`,
			params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword)
	}

	// Construct the SQL query dynamically using the passed primary organization unit field, keyword filter, and sorting
	rawSQL := fmt.Sprintf(`
        SELECT 
            users.id, 
            users.first_name_en, 
            users.first_name_ar, 
            users.last_name_en, 
            users.last_name_ar, 
            users.email,
			users.initial,
            users.color,
            users.profile_image_path, 
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE 
                WHEN document_acknowledge_users.id IS NOT NULL THEN 1 
                ELSE 0 
            END AS acknowledged
        FROM document_shared_users
        JOIN users ON users.id = document_shared_users.user_id
        LEFT JOIN document_acknowledge_users ON document_acknowledge_users.user_id = users.id AND document_acknowledge_users.document_id = ?
        LEFT JOIN organization_units ON organization_units.id = users.%s
        WHERE document_shared_users.document_id = ? %s
        
        UNION
        
        SELECT 
            users.id, 
            users.first_name_en, 
            users.first_name_ar, 
            users.last_name_en, 
            users.last_name_ar, 
            users.email,
			users.initial,
            users.color,
            users.profile_image_path,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE 
                WHEN document_acknowledge_users.id IS NOT NULL THEN 1 
                ELSE 0 
            END AS acknowledged
        FROM document_shared_units
        JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id
        JOIN users ON users.%s = organization_units.id
        LEFT JOIN document_acknowledge_users ON document_acknowledge_users.user_id = users.id AND document_acknowledge_users.document_id = ?
        WHERE document_shared_units.document_id = ? %s
    `, primaryOrgUnitField, keywordFilter, primaryOrgUnitField, keywordFilter)

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS subquery", rawSQL)
	err := repo.db.Raw(countSQL, documentID, documentID, documentID, documentID).Scan(&totalCount).Error
	if err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	paginatedSQL := fmt.Sprintf("%s ORDER BY %s %s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", rawSQL, params.SortBy, params.Sort, offset, params.PageSize)

	err = repo.db.Raw(paginatedSQL, documentID, documentID, documentID, documentID).Scan(&acknowledgementList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &documentDTOs.DocumentAcknowledgementListResponseDTO{
		Items: acknowledgementList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *DocumentAcknowledgeUserRepository) GetAllForPublicDocument(documentID uint, primaryOrgUnitField string, params *documentDTOs.DocumentAcknowledgeListQueryParamsDTO) (*documentDTOs.DocumentAcknowledgementListResponseDTO, error) {
	var acknowledgementList []documentDTOs.DocumentAcknowledgementListDTO
	var totalCount int64

	// Prepare the keyword search filter
	keywordFilter := ""
	if params.Keyword != "" {
		keywordFilter = fmt.Sprintf(`
            AND (users.first_name_en LIKE '%%%s%%' 
            OR users.first_name_ar LIKE '%%%s%%' 
            OR users.last_name_en LIKE '%%%s%%' 
            OR users.last_name_ar LIKE '%%%s%%' 
            OR users.email LIKE '%%%s%%'
            OR organization_units.title_en LIKE '%%%s%%'
            OR organization_units.title_ar LIKE '%%%s%%')`,
			params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword)
	}

	// SQL for fetching all users and checking their acknowledgment status
	rawSQL := fmt.Sprintf(`
        SELECT 
            users.id, 
            users.first_name_en, 
            users.first_name_ar, 
            users.last_name_en, 
            users.last_name_ar, 
            users.email,
			users.initial,
            users.color,
            users.profile_image_path,  
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE 
                WHEN document_acknowledge_users.id IS NOT NULL THEN 1 
                ELSE 0 
            END AS acknowledged
        FROM users
        LEFT JOIN document_acknowledge_users ON document_acknowledge_users.user_id = users.id AND document_acknowledge_users.document_id = ?
        LEFT JOIN organization_units ON organization_units.id = users.%s
        WHERE 1=1 %s
    `, primaryOrgUnitField, keywordFilter)

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS subquery", rawSQL)
	err := repo.db.Raw(countSQL, documentID).Scan(&totalCount).Error
	if err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	paginatedSQL := fmt.Sprintf("%s ORDER BY %s %s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", rawSQL, params.SortBy, params.Sort, offset, params.PageSize)

	err = repo.db.Raw(paginatedSQL, documentID).Scan(&acknowledgementList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &documentDTOs.DocumentAcknowledgementListResponseDTO{
		Items: acknowledgementList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *DocumentAcknowledgeUserRepository) Create(documentAcknowledgeUser *models.DocumentAcknowledgeUser, tx ...utility.DatabaseTransaction) (*models.DocumentAcknowledgeUser, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.Create(&documentAcknowledgeUser)

	return documentAcknowledgeUser, result.Error
}

func (repo *DocumentAcknowledgeUserRepository) Exists(documentID uint, userID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.DocumentAcknowledgeUser{}).
		Where("document_id = ? AND user_id = ?", documentID, userID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *DocumentAcknowledgeUserRepository) Delete(documentID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	condition := map[string]interface{}{
		"document_id": documentID,
		"user_id":     currentUserID,
	}

	// Find the record based on the condition
	var documentAcknowledgeUser models.DocumentAcknowledgeUser
	if err := db.Where(condition).First(&documentAcknowledgeUser).Error; err != nil {
		return err
	}

	// Delete the record
	result := db.Delete(&documentAcknowledgeUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *DocumentAcknowledgeUserRepository) GetDocumentAcknowledgmentSummary(documentID uint, primaryOrgUnitField string) (*documentDTOs.DocumentAcknowledgmentSummaryDTO, error) {
	var totalUsers int64
	var acknowledgedUsers int64

	totalUsersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id 
			FROM document_shared_users 
			JOIN users ON users.id = document_shared_users.user_id
			JOIN user_statuses ON user_statuses.id = users.user_status_id
			WHERE document_shared_users.document_id = ? AND user_statuses.is_active = 1
			UNION
			SELECT users.id AS user_id 
			FROM document_shared_units
			JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id
			JOIN users ON users.%s = organization_units.id
			JOIN user_statuses ON user_statuses.id = users.user_status_id
			WHERE document_shared_units.document_id = ? AND user_statuses.is_active = 1
		) AS combined_users
	`, primaryOrgUnitField)

	acknowledgedUsersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT document_acknowledge_users.user_id 
			FROM document_acknowledge_users
			JOIN users ON users.id = document_acknowledge_users.user_id
			JOIN user_statuses ON user_statuses.id = users.user_status_id
			WHERE document_acknowledge_users.document_id = ? 
			AND user_statuses.is_active = 1
			AND EXISTS (
				SELECT 1 
				FROM document_shared_users 
				WHERE document_shared_users.user_id = document_acknowledge_users.user_id 
				AND document_shared_users.document_id = document_acknowledge_users.document_id
				UNION
				SELECT 1 
				FROM document_shared_units
				JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id
				JOIN users AS unit_users ON unit_users.%s = organization_units.id
				JOIN user_statuses AS unit_status ON unit_status.id = unit_users.user_status_id
				WHERE document_shared_units.document_id = document_acknowledge_users.document_id
				AND unit_users.id = document_acknowledge_users.user_id AND unit_status.is_active = 1
			)
		) AS acknowledged_users
	`, primaryOrgUnitField)

	err := repo.db.Raw(totalUsersQuery, documentID, documentID).Scan(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching total users: %v", err)
	}

	err = repo.db.Raw(acknowledgedUsersQuery, documentID).Scan(&acknowledgedUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching acknowledged users: %v", err)
	}

	notAcknowledgedUsers := totalUsers - acknowledgedUsers

	summary := &documentDTOs.DocumentAcknowledgmentSummaryDTO{
		TotalUsers:             totalUsers,
		AcknowledgedUsers:      acknowledgedUsers,
		NotAcknowledgedUsers:   notAcknowledgedUsers,
		AcknowledgedPercent:    float64(acknowledgedUsers) / float64(totalUsers) * 100,
		NotAcknowledgedPercent: float64(notAcknowledgedUsers) / float64(totalUsers) * 100,
	}

	return summary, nil
}

func (repo *DocumentAcknowledgeUserRepository) GetPublicDocumentAcknowledgmentSummary(documentID uint) (*documentDTOs.DocumentAcknowledgmentSummaryDTO, error) {
	var totalUsers int64
	var acknowledgedUsers int64

	// Query to count all active users
	totalUsersQuery := `
		SELECT COUNT(*)
		FROM users
		JOIN user_statuses ON user_statuses.id = users.user_status_id
		WHERE user_statuses.is_active = 1
	`

	// Query to count acknowledged users for the public document
	acknowledgedUsersQuery := `
		SELECT COUNT(*)
		FROM document_acknowledge_users
		JOIN users ON users.id = document_acknowledge_users.user_id
		JOIN user_statuses ON user_statuses.id = users.user_status_id
		WHERE document_acknowledge_users.document_id = ? AND user_statuses.is_active = 1
	`

	// Fetch total users
	err := repo.db.Raw(totalUsersQuery).Scan(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching total users for public document: %v", err)
	}

	// Fetch acknowledged users
	err = repo.db.Raw(acknowledgedUsersQuery, documentID).Scan(&acknowledgedUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching acknowledged users for public document: %v", err)
	}

	// Calculate not acknowledged users
	notAcknowledgedUsers := totalUsers - acknowledgedUsers

	// Prepare and return the summary
	summary := &documentDTOs.DocumentAcknowledgmentSummaryDTO{
		TotalUsers:             totalUsers,
		AcknowledgedUsers:      acknowledgedUsers,
		NotAcknowledgedUsers:   notAcknowledgedUsers,
		AcknowledgedPercent:    float64(acknowledgedUsers) / float64(totalUsers) * 100,
		NotAcknowledgedPercent: float64(notAcknowledgedUsers) / float64(totalUsers) * 100,
	}

	return summary, nil
}
