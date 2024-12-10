package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssetAssociatedItemRepositoryInterface interface {
	GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]*assetDTOs.AssetControlListDTO, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemRiskListResponseDTO, error)
	GetAllVulnerabilities(assetID uint) ([]models.Vulnerability, error)
	GetRisksByThreatID(threatID uint) ([]*assetDTOs.RiskDTO, error)
}

type AssetAssociatedItemRepository struct {
	db *gorm.DB
}

// Constructor function for the repository
func NewAssetAssociatedItemRepository(db *gorm.DB) *AssetAssociatedItemRepository {
	return &AssetAssociatedItemRepository{
		db: db,
	}
}

func (repo *AssetAssociatedItemRepository) GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]*assetDTOs.AssetControlListDTO, error) {
	var controlList []*assetDTOs.AssetControlListDTO

	query := repo.db.Table("asset_controls").
		Select(`
			asset_controls.id AS id,
			asset_controls.asset_register_id AS asset_register_id,
			controls.id AS control_id,
			controls.reference_number AS control_reference_number,
            controls.code AS control_code,
			CONVERT(VARCHAR, controls.last_test_date, 23) AS control_last_test_date,
            control_categories.title_en AS control_category_title_en,
            control_categories.title_ar AS control_category_title_ar,
            control_types.title_en AS control_type_title_en,
            control_types.title_ar AS control_type_title_ar,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
            control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
            control_efficiency_statuses.type AS control_efficiency_status_type,
            control_implementation_statuses.title_en AS control_implementation_status_title_en,
            control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.type AS control_implementation_status_type,
            control_implementation_statuses.color_code AS control_implementation_status_color_code,
			owner.id AS control_owner_id,
			owner.first_name_en AS control_owner_first_name_en,
			owner.first_name_ar AS control_owner_first_name_ar,
			owner.last_name_en AS control_owner_last_name_en,
			owner.last_name_ar AS control_owner_last_name_ar,
			owner.initial AS control_owner_initial,
            owner.color AS control_owner_color,
            owner.email AS control_owner_email,
            owner.profile_image_path AS control_owner_profile_image_path
        `).
		Joins(`
			LEFT JOIN controls ON controls.id = asset_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
			LEFT JOIN users AS owner ON controls.owner_id = owner.id
        `).Where("asset_controls.asset_register_id = ?", assetID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			controls.code LIKE ? OR
			control_categories.title_en LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_efficiency_statuses.title_en LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			owner.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "controls.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&controlList).Error; err != nil {
		return nil, err
	}

	return controlList, nil
}

func (repo *AssetAssociatedItemRepository) GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error) {
	var totalCount int64
	var incidentList []*assetDTOs.AssetIncidentListDTO

	// Query to fetch asset incidents with fields matching AssetIncidentListDTO
	query := repo.db.Table("incident_assets").
		Select(`
			DISTINCT incident_register.id AS id,
			incident_assets.incident_register_id AS incident_register_id,
			incident_assets.asset_register_id AS asset_register_id,
			incident_register.reference_number AS reference_number,
			incident_register.title AS title,
			incident_register.description AS description,
			incident_register.incident_at AS incident_at,
			incident_register.location AS location,
			incident_register.reported_at AS reported_at,
			incident_register.quick_action AS quick_action,
			incident_categories.title_en AS incident_category_title_en,
			incident_categories.title_ar AS incident_category_title_ar,
			incident_sub_categories.title_en AS incident_sub_category_title_en,
			incident_sub_categories.title_ar AS incident_sub_category_title_ar,
			incident_types.title_en AS incident_type_title_en,
			incident_types.title_ar AS incident_type_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			incident_statuses.title_en AS incident_status_title_en,
			incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code,
			reported_user.id AS incident_reported_user_id,
			reported_user.first_name_en AS incident_reported_user_first_name_en,
			reported_user.first_name_ar AS incident_reported_user_first_name_ar,
			reported_user.last_name_en AS incident_reported_user_last_name_en,
			reported_user.last_name_ar AS incident_reported_user_last_name_ar,
			reported_user.initial AS incident_reported_user_initial,
			reported_user.color AS incident_reported_user_color,
			reported_user.email AS incident_reported_user_email,
			reported_user.profile_image_path AS incident_reported_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN incident_register ON incident_register.id = incident_assets.incident_register_id
			LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
			LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
			LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
			LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
			LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
			LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
		`).
		Where("incident_assets.asset_register_id = ?", assetID)

	// Count total records for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if provided
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("incident_register.title LIKE ?", keyword)
	}

	// Set default sorting if not provided
	if params.SortBy == "" {
		params.SortBy = "incident_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	// Apply sorting and pagination
	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize

	// Fetch paginated results
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentList).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Prepare response
	response := &assetDTOs.AssetAssociatedItemIncidentListResponseDTO{
		Items: incidentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AssetAssociatedItemRepository) GetAllRisks(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemRiskListResponseDTO, error) {
	var totalCount int64
	var riskList []*assetDTOs.AssetRiskListDTO

	// Query: Fetch asset-associated risks using the asset IDs
	query := repo.db.Table("risk_assets").
		Select(`
            DISTINCT risk_assets.risk_register_id AS risk_register_id,
            risk_register.reference_number AS risk_reference_number,
            risk_register.title AS risk_register_title,
            risk_register.description AS risk_register_description,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
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
            risk_statuses.color_code AS risk_status_color_code,
            review_frequencies.title_en AS risk_review_frequency_title_en,
            review_frequencies.title_ar AS risk_review_frequency_title_ar,
            risk_matrix_likelihoods.title_en AS risk_matrix_likelihood_title_en,
            risk_matrix_likelihoods.title_ar AS risk_matrix_likelihood_title_ar,
            risk_matrix_likelihoods.score AS risk_matrix_likelihood_score,
            risk_matrix_impacts.title_en AS risk_matrix_impact_title_en,
            risk_matrix_impacts.title_ar AS risk_matrix_impact_title_ar,
            risk_matrix_impacts.score AS risk_matrix_impact_score,
            risk_matrix_risk_ratings.title_en AS risk_matrix_risk_rating_title_en,
            risk_matrix_risk_ratings.title_ar AS risk_matrix_risk_rating_title_ar,
            risk_matrix_risk_ratings.type AS risk_matrix_risk_rating_type,
            risk_matrix_risk_ratings.color_code AS risk_matrix_risk_rating_color_code,
            risk_matrix_risk_rating_values.value AS risk_matrix_risk_rating_value,
            risk_register.next_review_date AS risk_next_review_date
        `).
		Joins(`
            LEFT JOIN risk_register ON risk_register.id = risk_assets.risk_register_id
            LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
            LEFT JOIN risk_matrix_likelihoods ON risk_matrix_likelihoods.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = risk_register.current_risk_rating_value_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
        `).
		Where("risk_assets.asset_register_id = ?", assetID)

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("risk_register.title LIKE ?", keyword)
	}

	// Filter by risk rating value
	if params.Value != 0 {
		query = query.Where("risk_matrix_risk_rating_values.value = ?", params.Value)
	}

	// Apply sorting
	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize

	// Fetch the results
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	// Build the response
	response := &assetDTOs.AssetAssociatedItemRiskListResponseDTO{
		Items: riskList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AssetAssociatedItemRepository) GetAllVulnerabilities(assetID uint) ([]models.Vulnerability, error) {
	var vulnerabilities []models.Vulnerability

	// Fetch vulnerabilities with their related threats
	err := repo.db.
		Preload("Threats.ThreatType").
		Preload("Threats.SeverityLevel").
		Preload("Threats.ThreatStatus").
		Preload("Threats.ResponsibleUser").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("VulnerabilityType").
		Preload("VulnerabilityStatus").
		Where("id IN (SELECT vulnerability_id FROM vulnerability_assets WHERE asset_register_id = ?)", assetID).
		Find(&vulnerabilities).Error

	return vulnerabilities, err
}

func (repo *AssetAssociatedItemRepository) GetRisksByThreatID(threatID uint) ([]*assetDTOs.RiskDTO, error) {
	var risks []*assetDTOs.RiskDTO

	query := repo.db.Table("risk_threats").
		Select(`
			DISTINCT risk_threats.risk_register_id AS id,
			risk_register.reference_number AS reference_number,
			risk_register.title AS title,
			risk_register.description AS risk_description,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
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
			users.profile_image_path AS responsible_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = risk_threats.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
		`).
		Where("risk_threats.threat_id = ?", threatID)

	// Execute query
	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	return risks, nil
}
