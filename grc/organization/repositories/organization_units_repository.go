package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type OrganizationUnitsRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]models.OrganizationUnit, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, level uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	GetComplianceSummaryByUnit(unitID uint) (int64, int64, []organizationDTOs.ComplianceStatusCountDTO, error)
	GetOpenRiskCountByUnit(unitID uint) (int64, error)
	GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error)
	SetAllIsLatestFalse(orgUnitID uint, primaryField string, tx ...utility.DatabaseTransaction) error
	GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetParentOrganizationUnitWithHead(unitID uint, primaryField string) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error)
	GetUserRoleDistribution(orgUnitID uint, level uint) ([]organizationDTOs.OrganizationUnitUserRoleDistributionDTO, error)
	GetAllComplianceByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedComplianceListResponseDTO, error)
	GetAllProcessByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	GetOrganizationUnitOverview(orgUnitID uint, level uint) (*organizationDTOs.OrganizationUnitOverviewDTO, error)
	GetOrganizationUnitKeyMatrix(orgUnitID uint) (*organizationDTOs.OrganizationUnitKeyMatrixDTO, error)
	GetRiskRatingDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)
	GetPolicyStatusDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)
	GetComplianceStatusDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)
	GetProcessCriticalityDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error)
}

type OrganizationUnitsRepository struct {
	db *gorm.DB
}

func NewOrganizationUnitsRepository(db *gorm.DB) *OrganizationUnitsRepository {
	return &OrganizationUnitsRepository{
		db: db,
	}
}

func (repo *OrganizationUnitsRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitsHierarchies(sort string, sortBy string) ([]models.OrganizationUnit, error) {
	var units []struct {
		models.OrganizationUnit
		Level        uint   `gorm:"column:level"`
		LevelTitleEn string `gorm:"column:level_title_en"`
		LevelTitleAr string `gorm:"column:level_title_ar"`
		IsPrimary    bool   `gorm:"column:is_primary"`
	}

	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Build the query with explicit joins
	query := repo.db.Table("organization_units").
		Select(`
            organization_units.*,
			organization_levels.level,
            organization_levels.title_en as level_title_en,
            organization_levels.title_ar as level_title_ar,
            organization_levels.is_primary as is_primary
        `).
		Joins("LEFT JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id").
		Order(fmt.Sprintf("%s %s", sortBy, sort))

	err := query.Scan(&units).Error
	if err != nil {
		return nil, err
	}

	// Map the results to the expected return type
	result := make([]models.OrganizationUnit, len(units))
	for i, u := range units {
		result[i] = u.OrganizationUnit
		if u.OrganizationLevel == nil {
			result[i].OrganizationLevel = &models.OrganizationLevel{}
		}
		result[i].OrganizationLevel.TitleEn = u.LevelTitleEn
		result[i].OrganizationLevel.TitleAr = u.LevelTitleAr
		result[i].OrganizationLevel.IsPrimary = u.IsPrimary
		result[i].OrganizationLevel.Level = uint8(u.Level)
	}

	return result, nil
}

func (repo *OrganizationUnitsRepository) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, level uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	var users []organizationDTOs.OrganizationUnitUserPrimaryDetailsDTO
	var totalCount int64

	// Build the base query
	query := repo.db.Table("users").
		Select(`
		 users.id,
		 users.first_name_en,
		 users.last_name_en,
		 users.first_name_ar,
		 users.last_name_ar,
		 users.is_head,
		 users.initial,
		 users.email,
		 users.profile_image_path,
		 designations.title_en AS designation_en,
		 designations.title_ar AS designation_ar`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id")

	// Build the WHERE clause based on level
	if level > 0 {
		// Construct the field name using the level parameter (e.g., organization_unit_1_id)
		field := fmt.Sprintf("organization_unit_%d_id", level)
		query = query.Where(fmt.Sprintf("users.%s = ?", field), orgUnitID)
	}

	// Rest of the function remains the same...
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(int(offset)).Limit(int(params.PageSize))

	// Apply sorting
	sortOrder := "CASE WHEN users.is_head = 1 THEN 0 ELSE 1 END, " + params.SortBy + " " + params.Sort
	query = query.Order(sortOrder)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *OrganizationUnitsRepository) GetComplianceSummaryByUnit(unitID uint) (int64, int64, []organizationDTOs.ComplianceStatusCountDTO, error) {
	var compliant, total int64
	var statuses []organizationDTOs.ComplianceStatusCountDTO

	// Step 1: compliant + total
	const q1 = `
		SELECT 
			COALESCE(SUM(CASE WHEN cs.type = 'compliant' THEN 1 ELSE 0 END), 0) AS compliant,
			COALESCE(COUNT(*), 0) AS total
		FROM compliance_assessment_requirement_responsible_units carru
		JOIN compliance_assessment_framework_values cafv
			ON cafv.id = carru.compliance_assessment_framework_value_id
		JOIN compliance_statuses cs
			ON cs.id = cafv.compliance_status_id
		JOIN compliance_assessment_requirement_statuses carrs
			ON carrs.id = carru.compliance_assessment_requirement_status_id
		WHERE carru.responsible_unit_id = ?
		  AND carrs.type = 'approved'
		  AND cs.type NOT IN ('exempt','not_applicable')
	`
	if err := repo.db.Raw(q1, unitID).Row().Scan(&compliant, &total); err != nil {
		return 0, 0, nil, err
	}

	// Step 2: per-status (for thresholds)
	const q2 = `
		SELECT
			cs.id   AS compliance_status_id,
			cs.title_en,
			cs.title_ar,
			cs.type,
			cs.color_code,
			cs.score_from,
			cs.score_to,
			COALESCE(agg.cnt, 0) AS count
		FROM compliance_statuses cs
		LEFT JOIN (
			SELECT 
				cafv.compliance_status_id,
				COUNT(*) AS cnt
			FROM compliance_assessment_requirement_responsible_units carru
			JOIN compliance_assessment_framework_values cafv
				ON cafv.id = carru.compliance_assessment_framework_value_id
			JOIN compliance_assessment_requirement_statuses carrs
				ON carrs.id = carru.compliance_assessment_requirement_status_id
			WHERE carru.responsible_unit_id = ?
			  AND carrs.type = 'approved'
			GROUP BY cafv.compliance_status_id
		) agg
			ON agg.compliance_status_id = cs.id
		ORDER BY cs.id
	`
	if err := repo.db.Raw(q2, unitID).Scan(&statuses).Error; err != nil {
		return 0, 0, nil, err
	}

	return compliant, total, statuses, nil
}

func (repo *OrganizationUnitsRepository) GetOpenRiskCountByUnit(unitID uint) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM risk_register r
		JOIN risk_statuses rs ON r.risk_status_id = rs.id
		WHERE r.responsible_unit_id = ? AND rs.type = 'open'
	`
	if err := repo.db.Raw(query, unitID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *OrganizationUnitsRepository) GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error) {
	var avgScore float64
	if err := repo.db.Model(&models.Control{}).
		Select("COALESCE(AVG(efficiency_score),0)").
		Where("responsible_unit_id = ?", unitID).
		Scan(&avgScore).Error; err != nil {
		return nil, err
	}

	var status organizationDTOs.ControlEfficiencyStatusDTO
	err := repo.db.Table("control_efficiency_statuses").
		Select("id, title_en, title_ar, status, color_code, score_from, score_to").
		Where("? BETWEEN score_from AND score_to AND status = 'active'", avgScore).
		First(&status).Error
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func (repo *OrganizationUnitsRepository) SetAllIsLatestFalse(orgUnitID uint, primaryField string, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.User{}).
		Where(fmt.Sprintf("%s = ? AND is_head = ?", primaryField), orgUnitID, true).
		Update("is_head", false).Error
}

func (repo *OrganizationUnitsRepository) GetAllPoliciesByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("documents").
		Select(`
			documents.id AS id,
			documents.id AS document_id,
			documents.reference_number AS reference_number,
			documents.document_type_id AS document_type_id,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			documents.last_checkout_at AS last_checkout_at,
			documents.last_checkin_at AS last_checkin_at,
			documents.is_locked AS is_locked,
			documents.description_en AS description_en,
			documents.description_ar AS description_ar,
			documents.responsible_user_id,
			documents.created_at,
			documents.updated_at,
			documents.published_at,
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
            published_user.profile_image_path AS published_user_profile_image_path
		`).
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
        `).
		Where("documents.responsible_unit_id = ?", unitID).
		Where("document_types.type = 'policy'").
		Where("document_statuses.type = 'published'")

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			documents.reference_number LIKE ? OR
			document_categories.title_en LIKE ? OR
			document_categories.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			document_statuses.title_en LIKE ? OR
			document_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "documents.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&documentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedDocumentListResponseDTO{
		Items: documentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *OrganizationUnitsRepository) GetAllRiskByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var totalCount int64
	var riskList []*commonDTOs.AssociatedRiskListDTO

	query := repo.db.Table("risk_register").
		Select(`
		    risk_register.id as id,
			risk_register.id as risk_register_id,
			 risk_register.reference_number,
            risk_register.title,
            risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
			risk_register.current_risk_rating_value_id,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
        	users.color AS responsible_user_color,
        	users.email AS responsible_user_email,
        	users.profile_image_path AS responsible_user_profile_image_path,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
            risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
            risk_register.last_review_note,
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
            current_likelihood.title_en AS current_likelihood_title_en,
            current_likelihood.title_ar AS current_likelihood_title_ar,
            current_impact.title_en AS current_impact_title_en,
            current_impact.title_ar AS current_impact_title_ar,
            current_risk_rating.title_en AS current_risk_rating_title_en,
            current_risk_rating.title_ar AS current_risk_rating_title_ar,
            current_risk_rating_value.value AS current_risk_rating_value,
            current_risk_rating.color_code AS current_risk_rating_color_code,
            inherent_likelihood.title_en AS inherent_likelihood_title_en,
            inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
            inherent_impact.title_en AS inherent_impact_title_en,
            inherent_impact.title_ar AS inherent_impact_title_ar,
            inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
            inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
            inherent_risk_rating_value.value AS inherent_risk_rating_value,
            inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
            residual_likelihood.title_en AS residual_likelihood_title_en,
            residual_likelihood.title_ar AS residual_likelihood_title_ar,
            residual_impact.title_en AS residual_impact_title_en,
            residual_impact.title_ar AS residual_impact_title_ar,
            residual_risk_rating.title_en AS residual_risk_rating_title_en,
            residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
            residual_risk_rating_value.value AS residual_risk_rating_value,
            residual_risk_rating.color_code AS residual_risk_rating_color_code
		`).
		Joins(`
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
            LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id
            LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id
            LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id
		`).
		Where("risk_register.responsible_unit_id = ?", unitID).
		Where("risk_statuses.type = ?", "open")

	// Count total records

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_register.title LIKE ?", keyword)
	}

	// Filter by risk rating value
	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
	}

	// Sorting defaults
	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &commonDTOs.AssociatedRiskListResponseDTO{
		Items: riskList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

// GetParentOrganizationUnitWithHead retrieves the main organization unit and its head user
func (repo *OrganizationUnitsRepository) GetParentOrganizationUnitWithHead(unitID uint, primaryField string) (*organizationDTOs.ParentOrganizationUnitWithHeadDTO, error) {
	var result organizationDTOs.ParentOrganizationUnitWithHeadDTO

	// First, get the main organization unit (where organization_unit_id is null)
	parentQuery := `
        SELECT 
            COALESCE(parent_ou.id, child_ou.id) as id,
            COALESCE(parent_ou.title_en, child_ou.title_en) as title_en,
            COALESCE(parent_ou.title_ar, child_ou.title_ar) as title_ar
        FROM organization_units child_ou
        LEFT JOIN organization_units parent_ou ON parent_ou.id = child_ou.organization_unit_id
        WHERE child_ou.id = ?`

	var parentUnit organizationDTOs.OrganizationUnitDTO
	err := repo.db.Raw(parentQuery, unitID).Scan(&parentUnit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No parent unit found, return empty result
			return &result, nil
		}
		return nil, fmt.Errorf("failed to fetch parent organization unit: %w", err)
	}

	result.OrganizationUnit = parentUnit

	// If we found a parent unit, try to find its head user from primary organization unit
	if parentUnit.ID != 0 {
		headUserQuery := fmt.Sprintf(`
            SELECT
                u.id,
                u.first_name_en,
                u.last_name_en,
                u.first_name_ar,
                u.last_name_ar,
                u.profile_image_path,
                u.is_head,
				u.initial
            FROM users u
            WHERE u.%s = ? AND u.is_head = 1`, primaryField)

		var headUser organizationDTOs.OrganizationUnitUserDTO
		err = repo.db.Raw(headUserQuery, unitID).Scan(&headUser).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to fetch head user: %w", err)
		}

		// Only set head user if found
		if headUser.ID != 0 {
			result.HeadUser = &headUser
		}
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetUserRoleDistribution(orgUnitID uint, level uint) ([]organizationDTOs.OrganizationUnitUserRoleDistributionDTO, error) {
	results := []organizationDTOs.OrganizationUnitUserRoleDistributionDTO{}

	query := repo.db.Table("user_roles").
		Select(`
			roles.id,
			roles.title_en as role,
			COUNT(DISTINCT users.id) AS count`).
		Joins("JOIN users ON users.id = user_roles.user_id").
		Joins("JOIN roles ON roles.id = user_roles.role_id")

	if level > 0 {
		field := fmt.Sprintf("organization_unit_%d_id", level)
		query = query.Where(fmt.Sprintf("users.%s = ?", field), orgUnitID)
	}

	query = query.Group("roles.id, roles.title_en")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) GetAllComplianceByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedComplianceListResponseDTO, error) {
	var complianceRegisters []*commonDTOs.AssociatedComplianceListDTO
	var totalCount int64

	query := repo.db.Table("compliance_register").
		Select(`
            compliance_register.id,
            compliance_register.id AS compliance_register_id,
			compliance_register.reference_number,
			compliance_register.standard_code,
			compliance_register.title_en,
			compliance_register.title_ar,
			compliance_register.description_en,
			compliance_register.description_ar,
			compliance_register.version,
			compliance_register.total_score,
			compliance_register.score,
			CONVERT(VARCHAR, compliance_register.deadline, 23) AS deadline,
        	CONVERT(VARCHAR, compliance_register.issue_date, 23) AS issue_date,
        	CONVERT(VARCHAR, compliance_register.expiry_date, 23) AS expiry_date,
       		CONVERT(VARCHAR, compliance_register.next_review_date, 23) AS next_review_date,
			CONVERT(VARCHAR, compliance_register.last_reviewed_date, 23) AS last_reviewed_date,
			compliance_assessment_frameworks.id AS compliance_assessment_framework_id,
			compliance_assessment_frameworks.title_en AS compliance_assessment_framework_title_en,
			compliance_assessment_frameworks.title_ar AS compliance_assessment_framework_title_ar,
			compliance_categories.title_en AS compliance_category_title_en,
			compliance_categories.title_ar AS compliance_category_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
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
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_responsible_user.id AS review_responsible_user_id,
			review_responsible_user.first_name_en AS review_responsible_user_first_name_en,
            review_responsible_user.first_name_ar AS review_responsible_user_first_name_ar,
            review_responsible_user.last_name_en AS review_responsible_user_last_name_en,
            review_responsible_user.last_name_ar AS review_responsible_user_last_name_ar,
			review_responsible_user.initial AS review_responsible_user_initial,
        	review_responsible_user.color AS review_responsible_user_color,
        	review_responsible_user.email AS review_responsible_user_email,
        	review_responsible_user.profile_image_path AS review_responsible_user_profile_image_path,  
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
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
        	updated_user.profile_image_path AS updated_user_profile_image_path        
        `).
		Joins(`
		    LEFT JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id
		    LEFT JOIN organization_units ON organization_units.id = compliance_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_register.responsible_user_id
		    LEFT JOIN review_frequencies ON review_frequencies.id = compliance_register.review_frequency_id
			LEFT JOIN users AS review_responsible_user ON review_responsible_user.id = compliance_register.review_responsible_user_id
		    LEFT JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id
			LEFT JOIN compliance_assessment_frameworks ON compliance_assessment_frameworks.id = compliance_register.compliance_assessment_framework_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = compliance_register.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_register.updated_by_id
		`).
		Where("compliance_register.responsible_unit_id = ?", unitID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "compliant":
			query = query.Where("compliance_statuses.type = ?", "compliant")
		case "partially_compliant":
			query = query.Where("compliance_statuses.type = ?", "partially_compliant")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not-implemented")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("compliance_register.title_en LIKE ? OR compliance_register.title_ar LIKE ?", keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "compliance_register.id"
	}
	sortOrder := params.Sort
	if sortOrder == "" {
		sortOrder = "DESC"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	if params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&complianceRegisters).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	if totalPages == 0 {
		totalPages = 1
	}

	response := &commonDTOs.AssociatedComplianceListResponseDTO{
		Items: complianceRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *OrganizationUnitsRepository) GetAllProcessByOrganizationUnit(unitID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	var processList []*commonDTOs.AssociatedProcessListDTO
	var totalCount int64

	query := repo.db.Table("processes").
		Select(`
		processes.id,
		processes.id AS process_id,
		processes.reference_number,
		processes.title_en,
		processes.title_ar,
		processes.description,
		processes.marked_as_auditable_item,
		process_groups.title_en AS process_group_title_en,
		process_groups.title_ar AS process_group_title_ar,
		processes.rto as rto,
		processes.rto_unit_id as rto_unit_id,
		rto_time_unit.title_en AS rto_unit_title_en,
		rto_time_unit.title_ar AS rto_unit_title_ar,
		processes.mtpd as mtpd,
		processes.mtpd_unit_id as mtpd_unit_id,
		mtpd_time_unit.title_en AS mtpd_unit_title_en,
		mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
		processes.rpo as rpo,
		processes.rpo_unit_id as rpo_unit_id,
		rpo_time_unit.title_en AS rpo_unit_title_en,
		rpo_time_unit.title_ar AS rpo_unit_title_ar,
		processes.is_bia_performed,
		processes.last_assessed_at,
		last_assessed_user.id AS last_assessed_user_id,
		last_assessed_user.first_name_en AS last_assessed_user_first_name_en,
		last_assessed_user.first_name_ar AS last_assessed_user_first_name_ar,
		last_assessed_user.last_name_en AS last_assessed_user_last_name_en,
		last_assessed_user.last_name_ar AS last_assessed_user_last_name_ar,
		last_assessed_user.initial AS last_assessed_user_initial,
		last_assessed_user.color AS last_assessed_user_color,
		last_assessed_user.email AS last_assessed_user_email,
		last_assessed_user.profile_image_path AS last_assessed_user_profile_image_path,
		organization_units.title_en AS owner_unit_title_en,
		organization_units.title_ar AS owner_unit_title_ar,
		owner.id AS owner_id,
		owner.first_name_en AS owner_first_name_en,
		owner.first_name_ar AS owner_first_name_ar,
		owner.last_name_en AS owner_last_name_en,
		owner.last_name_ar AS owner_last_name_ar,
		owner.initial AS owner_initial,
		owner.color AS owner_color,
		owner.email AS owner_email,
		owner.profile_image_path AS owner_profile_image_path,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		process_criticality_levels.id AS process_criticality_level_id,
		process_criticality_levels.title_en AS process_criticality_level_title_en,
		process_criticality_levels.title_ar AS process_criticality_level_title_ar,
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
		processes.created_at,
		processes.updated_at,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		risk_levels.title_en AS risk_level_title_en,
        risk_levels.title_ar AS risk_level_title_ar,
		risk_levels.color_code AS risk_level_color_code,
		risk_levels.type AS risk_level_type
		`).
		Joins(`
			LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id
			LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = processes.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = processes.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = processes.rpo_unit_id
			LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id
		`).
		Where("processes.responsible_unit_id = ?", unitID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ?)`,
			keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "processes.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&processList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemProcessListResponseDTO{
		Items: processList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitOverview(orgUnitID uint, level uint) (*organizationDTOs.OrganizationUnitOverviewDTO, error) {

	var result organizationDTOs.OrganizationUnitOverviewDTO

	// Determine which organization_unit_X_id column to use dynamically
	field := fmt.Sprintf("organization_unit_%d_id", level)

	query := repo.db.Table("organization_units ou").
		Select(`
			ou.id AS organization_unit_id,
			ou.title_en AS organization_unit_title,
			CONCAT(head.first_name_en, ' ', head.last_name_en) AS organization_unit_head,
			head.email AS organization_unit_head_email,
			COUNT(DISTINCT users.id) AS total_employees
		`).
		Joins(fmt.Sprintf(` 
			LEFT JOIN users head 
				ON head.is_head = 1
				AND head.%s = ou.id
		`, field)).
		Joins(fmt.Sprintf(`
			LEFT JOIN users 
				ON users.%s = ou.id
		`, field)).
		Where("ou.status = ?", "active").
		Where("ou.id = ?", orgUnitID).
		Group("ou.id, ou.title_en, head.first_name_en, head.last_name_en, head.email")

	if err := query.Take(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *OrganizationUnitsRepository) GetOrganizationUnitKeyMatrix(unitID uint) (*organizationDTOs.OrganizationUnitKeyMatrixDTO, error) {
	var result *organizationDTOs.OrganizationUnitKeyMatrixDTO

	// risks
	if err := repo.db.Table("risk_register").
		Where("responsible_unit_id = ?", unitID).
		Count(&result.TotalRisks).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("documents").
		Joins("JOIN document_types ON document_types.id = documents.document_type_id").
		Where("documents.responsible_unit_id = ?", unitID).
		Where("document_types.type = ?", "policy").
		Count(&result.TotalPolicies).Error; err != nil {
		return nil, err
	}

	// processes
	if err := repo.db.Table("processes").
		Where("responsible_unit_id = ?", unitID).
		Count(&result.TotalProcesses).Error; err != nil {
		return nil, err
	}

	// compliances
	if err := repo.db.Table("compliance_register").
		Where("responsible_unit_id = ?", unitID).
		Count(&result.TotalCompliances).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *OrganizationUnitsRepository) GetRiskRatingDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	q := repo.db.Table("risk_register").
		Select(`
			ISNULL(current_risk_rating.id, 0) AS id,
			ISNULL(current_risk_rating.title_en, '') AS title_en,
			ISNULL(current_risk_rating.title_ar, '') AS title_ar,
			ISNULL(current_risk_rating.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN risk_statuses
				ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating
				ON current_risk_rating.id = risk_register.current_risk_rating_id
		`).
		Where("risk_register.responsible_unit_id = ?", unitID).
		Where("risk_statuses.type = ?", "open").
		Group(`
			current_risk_rating.id,
			current_risk_rating.title_en,
			current_risk_rating.title_ar,
			current_risk_rating.color_code
		`).
		Order("current_risk_rating.id ASC")

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil

}

func (repo *OrganizationUnitsRepository) GetPolicyStatusDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	q := repo.db.Table("documents").
		Select(`
			ISNULL(document_statuses.id, 0) AS id,
			ISNULL(document_statuses.title_en, '') AS title_en,
			ISNULL(document_statuses.title_ar, '') AS title_ar,
			ISNULL(document_statuses.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN document_types
				ON document_types.id = documents.document_type_id
			LEFT JOIN document_statuses
				ON document_statuses.id = documents.document_status_id
		`).
		Where("documents.responsible_unit_id = ?", unitID).
		Where("document_types.type = 'policy'").
		Group(`
			document_statuses.id,
			document_statuses.title_en,
			document_statuses.title_ar,
			document_statuses.color_code
		`).
		Order("document_statuses.title_en ASC")

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) GetProcessCriticalityDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	q := repo.db.Table("processes").
		Select(`
			ISNULL(process_criticality_levels.id, 0) AS id,
			ISNULL(process_criticality_levels.title_en, '') AS title_en,
			ISNULL(process_criticality_levels.title_ar, '') AS title_ar,
			ISNULL(process_criticality_levels.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN process_criticality_levels
				ON process_criticality_levels.id = processes.process_criticality_level_id
		`).
		Where("processes.responsible_unit_id = ?", unitID).
		Group(`
			process_criticality_levels.id,
			process_criticality_levels.title_en,
			process_criticality_levels.title_ar,
			process_criticality_levels.color_code
		`).
		Order("process_criticality_levels.title_en ASC")

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *OrganizationUnitsRepository) GetComplianceStatusDistributionByOrganizationUnit(unitID uint) ([]organizationDTOs.OrganizationUnitDistributionItemDTO, error) {
	results := []organizationDTOs.OrganizationUnitDistributionItemDTO{}

	q := repo.db.Table("compliance_register").
		Select(`
			ISNULL(compliance_statuses.id, 0) AS id,
			ISNULL(compliance_statuses.title_en, '') AS title_en,
			ISNULL(compliance_statuses.title_ar, '') AS title_ar,
			ISNULL(compliance_statuses.type, '') AS type,
			ISNULL(compliance_statuses.color_code, '') AS color_code,
			COUNT(*) AS count
		`).
		Joins(`
			LEFT JOIN compliance_statuses
				ON compliance_statuses.id = compliance_register.compliance_status_id
		`).
		Where("compliance_register.responsible_unit_id = ?", unitID).
		Group(`
			compliance_statuses.id,
			compliance_statuses.title_en,
			compliance_statuses.title_ar,
			compliance_statuses.type,
			compliance_statuses.color_code
		`).
		Order("compliance_statuses.title_en ASC")

	if err := q.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
