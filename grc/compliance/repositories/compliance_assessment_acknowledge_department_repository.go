package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceAssessmentAcknowledgeDepartmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assessmentID uint, currentUserID *uint, primaryOrgUnitField string, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error)
	GetAllForPublicAssessment(assessmentID uint, currentUserID *uint, primaryOrgUnitField string, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error)
	GetAllForSharedAssessment(assessmentID uint, complianceRegisterID uint, currentUserID *uint, primaryOrgUnitField string, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error)
	Create(acknowledgeDepartment *models.ComplianceAssessmentAcknowledgeDepartment, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentAcknowledgeDepartment, error)
	DepartmentExists(assessmentID uint, departmentID uint, tx ...utility.DatabaseTransaction) (bool, error)
	Delete(assessmentID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	GetAssessmentAcknowledgmentSummary(assessmentID uint, primaryOrgUnitField string) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error)
	GetPublicAssessmentAcknowledgmentSummary(assessmentID uint) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error)
	GetSharedAssessmentAcknowledgmentSummary(assessmentID uint, complianceRegisterID uint) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error)
	GetOrgUnitAcknowledgmentSummary(assessmentID uint, primaryOrgUnitField string) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error)
	GetPublicAcknowledgmentSummaryByOrgUnit(assessmentID uint) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error)
}

// Implementation
type ComplianceAssessmentAcknowledgeDepartmentRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentAcknowledgeDepartmentRepository(db *gorm.DB) *ComplianceAssessmentAcknowledgeDepartmentRepository {
	return &ComplianceAssessmentAcknowledgeDepartmentRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetAll(assessmentID uint, currentUserID *uint, primaryOrgUnitField string, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error) {
	var acknowledgementList []complianceDTOs.ComplianceAssessmentAcknowledgementListDTO
	var totalCount int64

	// Prepare the keyword search filter
	keywordFilter := ""
	if params.Keyword != "" {
		keywordFilter = fmt.Sprintf(`
            AND (ack_user.first_name_en LIKE '%%%s%%' 
            OR ack_user.first_name_ar LIKE '%%%s%%' 
            OR ack_user.last_name_en LIKE '%%%s%%' 
            OR ack_user.last_name_ar LIKE '%%%s%%' 
            OR ack_user.email LIKE '%%%s%%' 
            OR organization_units.title_en LIKE '%%%s%%' 
            OR organization_units.title_ar LIKE '%%%s%%')`,
			params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword)
	}

	// Build the main query to show all departments with user details only if someone from that department acknowledged
	var query string
	if currentUserID != nil {
		// Apply user-based filtering
		query = fmt.Sprintf(`
        SELECT 
            COALESCE(ack_user.id, 0) AS id,
            COALESCE(ack_user.first_name_en, '') AS first_name_en,
            COALESCE(ack_user.first_name_ar, '') AS first_name_ar,
            COALESCE(ack_user.last_name_en, '') AS last_name_en,
            COALESCE(ack_user.last_name_ar, '') AS last_name_ar,
            COALESCE(ack_user.email, '') AS email,
            COALESCE(ack_user.initial, '') AS initial,
            COALESCE(ack_user.color, '') AS color,
            COALESCE(caad.comment, '') AS comment,
            COALESCE(ack_user.profile_image_path, '') AS profile_image_path,
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE WHEN caad.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        INNER JOIN users user_filter ON user_filter.id = %d AND (
            user_filter.%s = organization_units.id OR
            user_filter.organization_unit_2_id = organization_units.id OR
            user_filter.organization_unit_3_id = organization_units.id
        )
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        LEFT JOIN users ack_user ON caad.user_id = ack_user.id
        WHERE organization_units.status = 'active' %s
        ORDER BY organization_units.title_en, ack_user.first_name_en %s`,
			*currentUserID, primaryOrgUnitField, assessmentID, keywordFilter, params.Sort)
	} else {
		// No user filtering for super admin
		query = fmt.Sprintf(`
        SELECT 
            COALESCE(ack_user.id, 0) AS id,
            COALESCE(ack_user.first_name_en, '') AS first_name_en,
            COALESCE(ack_user.first_name_ar, '') AS first_name_ar,
            COALESCE(ack_user.last_name_en, '') AS last_name_en,
            COALESCE(ack_user.last_name_ar, '') AS last_name_ar,
            COALESCE(ack_user.email, '') AS email,
            COALESCE(ack_user.initial, '') AS initial,
            COALESCE(ack_user.color, '') AS color,
            COALESCE(caad.comment, '') AS comment,
            COALESCE(ack_user.profile_image_path, '') AS profile_image_path,
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE WHEN caad.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        LEFT JOIN users ack_user ON caad.user_id = ack_user.id
        WHERE organization_units.status = 'active' %s
        ORDER BY organization_units.title_en, ack_user.first_name_en %s`,
			assessmentID, keywordFilter, params.Sort)
	}

	// Execute count query for active departments only
	var countQuery string
	if currentUserID != nil {
		// Apply user-based filtering
		countQuery = fmt.Sprintf(`
        SELECT COUNT(DISTINCT organization_units.id)
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        INNER JOIN users user_filter ON user_filter.id = %d AND (
            user_filter.%s = organization_units.id OR
            user_filter.organization_unit_2_id = organization_units.id OR
            user_filter.organization_unit_3_id = organization_units.id
        )
        WHERE organization_units.status = 'active' %s`,
			*currentUserID, primaryOrgUnitField, keywordFilter)
	} else {
		// No user filtering for super admin
		countQuery = fmt.Sprintf(`
        SELECT COUNT(DISTINCT organization_units.id)
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        WHERE organization_units.status = 'active' %s`,
			keywordFilter)
	}

	if err := repo.db.Raw(countQuery).Scan(&totalCount).Error; err != nil {
		return nil, err
	}

	// Execute main query with pagination
	offset := (params.Page - 1) * params.PageSize
	paginatedQuery := fmt.Sprintf("%s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", query, offset, params.PageSize)

	if err := repo.db.Raw(paginatedQuery).Scan(&acknowledgementList).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO{
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

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetAllForPublicAssessment(assessmentID uint, currentUserID *uint, primaryOrgUnitField string, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error) {
	var acknowledgementList []complianceDTOs.ComplianceAssessmentAcknowledgementListDTO
	var totalCount int64

	// Prepare the keyword search filter
	keywordFilter := ""
	if params.Keyword != "" {
		keywordFilter = fmt.Sprintf(`
            AND (ack_user.first_name_en LIKE '%%%s%%' 
            OR ack_user.first_name_ar LIKE '%%%s%%' 
            OR ack_user.last_name_en LIKE '%%%s%%' 
            OR ack_user.last_name_ar LIKE '%%%s%%' 
            OR ack_user.email LIKE '%%%s%%' 
            OR organization_units.title_en LIKE '%%%s%%' 
            OR organization_units.title_ar LIKE '%%%s%%')`,
			params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword)
	}

	// Build main query for public assessment - show all departments with user details only if acknowledged
	var query string
	if currentUserID != nil {
		// Apply user-based filtering
		query = fmt.Sprintf(`
        SELECT 
            COALESCE(ack_user.id, 0) AS id,
            COALESCE(ack_user.first_name_en, '') AS first_name_en,
            COALESCE(ack_user.first_name_ar, '') AS first_name_ar,
            COALESCE(ack_user.last_name_en, '') AS last_name_en,
            COALESCE(ack_user.last_name_ar, '') AS last_name_ar,
            COALESCE(ack_user.email, '') AS email,
            COALESCE(ack_user.initial, '') AS initial,
            COALESCE(ack_user.color, '') AS color,
            COALESCE(caad.comment, '') AS comment,
            COALESCE(ack_user.profile_image_path, '') AS profile_image_path,
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE WHEN caad.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        INNER JOIN users user_filter ON user_filter.id = %d AND (
            user_filter.%s = organization_units.id OR
            user_filter.organization_unit_2_id = organization_units.id OR
            user_filter.organization_unit_3_id = organization_units.id
        )
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        LEFT JOIN users ack_user ON caad.user_id = ack_user.id
        WHERE organization_units.status = 'active' %s
        ORDER BY organization_units.title_en, ack_user.first_name_en %s`,
			*currentUserID, primaryOrgUnitField, assessmentID, keywordFilter, params.Sort)
	} else {
		// No user filtering for super admin
		query = fmt.Sprintf(`
        SELECT 
            COALESCE(ack_user.id, 0) AS id,
            COALESCE(ack_user.first_name_en, '') AS first_name_en,
            COALESCE(ack_user.first_name_ar, '') AS first_name_ar,
            COALESCE(ack_user.last_name_en, '') AS last_name_en,
            COALESCE(ack_user.last_name_ar, '') AS last_name_ar,
            COALESCE(ack_user.email, '') AS email,
            COALESCE(ack_user.initial, '') AS initial,
            COALESCE(ack_user.color, '') AS color,
            COALESCE(caad.comment, '') AS comment,
            COALESCE(ack_user.profile_image_path, '') AS profile_image_path,
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE WHEN caad.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        LEFT JOIN users ack_user ON caad.user_id = ack_user.id
        WHERE organization_units.status = 'active' %s
        ORDER BY organization_units.title_en, ack_user.first_name_en %s`,
			assessmentID, keywordFilter, params.Sort)
	}

	// Execute count query for active departments only
	var countQuery string
	if currentUserID != nil {
		// Apply user-based filtering
		countQuery = fmt.Sprintf(`
        SELECT COUNT(DISTINCT organization_units.id)
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        INNER JOIN users user_filter ON user_filter.id = %d AND (
            user_filter.%s = organization_units.id OR
            user_filter.organization_unit_2_id = organization_units.id OR
            user_filter.organization_unit_3_id = organization_units.id
        )
        WHERE organization_units.status = 'active' %s`,
			*currentUserID, primaryOrgUnitField, keywordFilter)
	} else {
		// No user filtering for super admin
		countQuery = fmt.Sprintf(`
        SELECT COUNT(DISTINCT organization_units.id)
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        WHERE organization_units.status = 'active' %s`,
			keywordFilter)
	}

	if err := repo.db.Raw(countQuery).Scan(&totalCount).Error; err != nil {
		return nil, err
	}

	// Execute main query with pagination
	offset := (params.Page - 1) * params.PageSize
	paginatedQuery := fmt.Sprintf("%s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", query, offset, params.PageSize)

	if err := repo.db.Raw(paginatedQuery).Scan(&acknowledgementList).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO{
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

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetAllForSharedAssessment(assessmentID uint, complianceRegisterID uint, currentUserID *uint, primaryOrgUnitField string, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error) {
	var acknowledgementList []complianceDTOs.ComplianceAssessmentAcknowledgementListDTO
	var totalCount int64

	// Prepare the keyword search filter
	keywordFilter := ""
	if params.Keyword != "" {
		keywordFilter = fmt.Sprintf(`
            AND (ack_user.first_name_en LIKE '%%%s%%' 
            OR ack_user.first_name_ar LIKE '%%%s%%' 
            OR ack_user.last_name_en LIKE '%%%s%%' 
            OR ack_user.last_name_ar LIKE '%%%s%%' 
            OR ack_user.email LIKE '%%%s%%' 
            OR organization_units.title_en LIKE '%%%s%%' 
            OR organization_units.title_ar LIKE '%%%s%%')`,
			params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword, params.Keyword)
	}

	// Build main query for shared assessment - show only shared departments with user details only if acknowledged
	var query string
	if currentUserID != nil {
		// Apply user-based filtering
		query = fmt.Sprintf(`
        SELECT 
            COALESCE(ack_user.id, 0) AS id,
            COALESCE(ack_user.first_name_en, '') AS first_name_en,
            COALESCE(ack_user.first_name_ar, '') AS first_name_ar,
            COALESCE(ack_user.last_name_en, '') AS last_name_en,
            COALESCE(ack_user.last_name_ar, '') AS last_name_ar,
            COALESCE(ack_user.email, '') AS email,
            COALESCE(ack_user.initial, '') AS initial,
            COALESCE(ack_user.color, '') AS color,
            COALESCE(caad.comment, '') AS comment,
            COALESCE(ack_user.profile_image_path, '') AS profile_image_path,
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE WHEN caad.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM organization_units
        INNER JOIN compliance_register_shared_units crsu ON crsu.department_id = organization_units.id AND crsu.compliance_register_id = %d
        INNER JOIN users user_filter ON user_filter.id = %d AND (
            user_filter.%s = organization_units.id OR
            user_filter.organization_unit_2_id = organization_units.id OR
            user_filter.organization_unit_3_id = organization_units.id
        )
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        LEFT JOIN users ack_user ON caad.user_id = ack_user.id
        WHERE organization_units.status = 'active' %s
        ORDER BY organization_units.title_en, ack_user.first_name_en %s`,
			complianceRegisterID, *currentUserID, primaryOrgUnitField, assessmentID, keywordFilter, params.Sort)
	} else {
		// No user filtering for super admin
		query = fmt.Sprintf(`
        SELECT 
            COALESCE(ack_user.id, 0) AS id,
            COALESCE(ack_user.first_name_en, '') AS first_name_en,
            COALESCE(ack_user.first_name_ar, '') AS first_name_ar,
            COALESCE(ack_user.last_name_en, '') AS last_name_en,
            COALESCE(ack_user.last_name_ar, '') AS last_name_ar,
            COALESCE(ack_user.email, '') AS email,
            COALESCE(ack_user.initial, '') AS initial,
            COALESCE(ack_user.color, '') AS color,
            COALESCE(caad.comment, '') AS comment,
            COALESCE(ack_user.profile_image_path, '') AS profile_image_path,
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            CASE WHEN caad.id IS NOT NULL THEN 1 ELSE 0 END AS acknowledged
        FROM organization_units
        INNER JOIN compliance_register_shared_units crsu ON crsu.department_id = organization_units.id AND crsu.compliance_register_id = %d
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        LEFT JOIN users ack_user ON caad.user_id = ack_user.id
        WHERE organization_units.status = 'active' %s
        ORDER BY organization_units.title_en, ack_user.first_name_en %s`,
			complianceRegisterID, assessmentID, keywordFilter, params.Sort)
	}

	// Execute count query for shared departments only
	var countQuery string
	if currentUserID != nil {
		// Apply user-based filtering
		countQuery = fmt.Sprintf(`
        SELECT COUNT(DISTINCT organization_units.id)
        FROM organization_units
        INNER JOIN compliance_register_shared_units crsu ON crsu.department_id = organization_units.id AND crsu.compliance_register_id = %d
        INNER JOIN users user_filter ON user_filter.id = %d AND (
            user_filter.%s = organization_units.id OR
            user_filter.organization_unit_2_id = organization_units.id OR
            user_filter.organization_unit_3_id = organization_units.id
        )
        WHERE organization_units.status = 'active' %s`,
			complianceRegisterID, *currentUserID, primaryOrgUnitField, keywordFilter)
	} else {
		// No user filtering for super admin
		countQuery = fmt.Sprintf(`
        SELECT COUNT(DISTINCT organization_units.id)
        FROM organization_units
        INNER JOIN compliance_register_shared_units crsu ON crsu.department_id = organization_units.id AND crsu.compliance_register_id = %d
        WHERE organization_units.status = 'active' %s`,
			complianceRegisterID, keywordFilter)
	}

	if err := repo.db.Raw(countQuery).Scan(&totalCount).Error; err != nil {
		return nil, err
	}

	// Execute main query with pagination
	offset := (params.Page - 1) * params.PageSize
	paginatedQuery := fmt.Sprintf("%s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", query, offset, params.PageSize)

	if err := repo.db.Raw(paginatedQuery).Scan(&acknowledgementList).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO{
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

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) Create(acknowledgeDepartment *models.ComplianceAssessmentAcknowledgeDepartment, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentAcknowledgeDepartment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(acknowledgeDepartment).Error; err != nil {
		return nil, err
	}

	return acknowledgeDepartment, nil
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) DepartmentExists(assessmentID uint, departmentID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.ComplianceAssessmentAcknowledgeDepartment{}).
		Where("compliance_assessment_id = ? AND department_id = ?", assessmentID, departmentID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) Delete(assessmentID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_assessment_id = ? AND user_id = ?", assessmentID, currentUserID).
		Delete(&models.ComplianceAssessmentAcknowledgeDepartment{}).Error
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetAssessmentAcknowledgmentSummary(assessmentID uint, primaryOrgUnitField string) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error) {
	var summary complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO

	query := fmt.Sprintf(`
        SELECT 
            COUNT(DISTINCT organization_units.id) AS total_departments,
            COUNT(DISTINCT caad.department_id) AS acknowledged_departments,
            (COUNT(DISTINCT organization_units.id) - COUNT(DISTINCT caad.department_id)) AS pending_departments
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        WHERE organization_units.status = 'active'`,
		assessmentID)

	if err := repo.db.Raw(query).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetPublicAssessmentAcknowledgmentSummary(assessmentID uint) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error) {
	var summary complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO

	query := fmt.Sprintf(`
        SELECT 
            COUNT(DISTINCT organization_units.id) AS total_departments,
            COUNT(DISTINCT caad.department_id) AS acknowledged_departments,
            (COUNT(DISTINCT organization_units.id) - COUNT(DISTINCT caad.department_id)) AS pending_departments
        FROM organization_units
        INNER JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id AND organization_levels.is_primary = 1
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        WHERE organization_units.status = 'active'`,
		assessmentID)

	if err := repo.db.Raw(query).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetSharedAssessmentAcknowledgmentSummary(assessmentID uint, complianceRegisterID uint) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error) {
	var summary complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO

	query := fmt.Sprintf(`
        SELECT 
            COUNT(DISTINCT organization_units.id) AS total_departments,
            COUNT(DISTINCT caad.department_id) AS acknowledged_departments,
            (COUNT(DISTINCT organization_units.id) - COUNT(DISTINCT caad.department_id)) AS pending_departments
        FROM organization_units
        INNER JOIN compliance_register_shared_units crsu ON crsu.department_id = organization_units.id AND crsu.compliance_register_id = %d
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.department_id = organization_units.id AND caad.compliance_assessment_id = %d
        WHERE organization_units.status = 'active'`,
		complianceRegisterID, assessmentID)

	if err := repo.db.Raw(query).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return &summary, nil
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetOrgUnitAcknowledgmentSummary(assessmentID uint, primaryOrgUnitField string) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error) {
	var summaries []complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO

	query := fmt.Sprintf(`
        SELECT 
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            COUNT(users.id) AS total_users,
            COUNT(caad.id) AS acknowledged_users,
            (COUNT(users.id) - COUNT(caad.id)) AS not_acknowledged_users,
            CASE 
                WHEN COUNT(users.id) = 0 THEN 0 
                ELSE CAST(COUNT(caad.id) AS FLOAT) / COUNT(users.id) * 100 
            END AS acknowledged_percent
        FROM organization_units
        INNER JOIN users ON users.%s = organization_units.id
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.user_id = users.id AND caad.compliance_assessment_id = %d
        GROUP BY organization_units.id, organization_units.title_en, organization_units.title_ar
        ORDER BY organization_units.title_en`,
		primaryOrgUnitField, assessmentID)

	if err := repo.db.Raw(query).Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}

func (repo *ComplianceAssessmentAcknowledgeDepartmentRepository) GetPublicAcknowledgmentSummaryByOrgUnit(assessmentID uint) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error) {
	var summaries []complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO

	query := fmt.Sprintf(`
        SELECT 
            organization_units.id AS organization_unit_id,
            organization_units.title_en AS organization_unit_title_en,
            organization_units.title_ar AS organization_unit_title_ar,
            COUNT(users.id) AS total_users,
            COUNT(caad.id) AS acknowledged_users,
            (COUNT(users.id) - COUNT(caad.id)) AS not_acknowledged_users,
            CASE 
                WHEN COUNT(users.id) = 0 THEN 0 
                ELSE CAST(COUNT(caad.id) AS FLOAT) / COUNT(users.id) * 100 
            END AS acknowledged_percent
        FROM organization_units
        INNER JOIN users ON users.organization_unit_1_id = organization_units.id
        LEFT JOIN compliance_assessment_acknowledge_departments caad ON caad.user_id = users.id AND caad.compliance_assessment_id = %d
        GROUP BY organization_units.id, organization_units.title_en, organization_units.title_ar
        ORDER BY organization_units.title_en`,
		assessmentID)

	if err := repo.db.Raw(query).Scan(&summaries).Error; err != nil {
		return nil, err
	}

	return summaries, nil
}
