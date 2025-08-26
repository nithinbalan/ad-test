package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetAssociatedItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAllControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
	GetAllRisksByRiskIDs(riskRegisterIDs []uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllVulnerabilities(assetID uint) ([]models.Vulnerability, error)
	GetAllClauses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)

	GetRisksByThreatID(threatID uint) ([]*assetDTOs.RiskDTO, error)

	GetVulnerabilityAssetExist(assetID, vulnerabilityID uint) (bool, error)
	CreateAssetVulnerabilities(vulnerabilities []*models.VulnerabilityAsset, tx ...utility.DatabaseTransaction) error
	DeleteAssetVulnerability(assetID, vulnerabilityID uint) error

	GetAssetComplianceClauseExist(assetID uint, clauseID uint) (bool, error)
	CreateAssetComplianceClauses(clauses []*models.AssetComplianceRequirement, tx ...utility.DatabaseTransaction) error
	GetAssetComplianceClauseByID(id uint) (*models.AssetComplianceRequirement, error)
	DeleteAssetComplianceClause(assetID uint, id uint) error

	GetVulnerabilityThreatExist(vulnerabilityId uint, threatID uint) (bool, error)
	CreateAssetVulnerabilityThreats(threats []*models.VulnerabilityThreat, tx ...utility.DatabaseTransaction) error
	DeleteVulnerabilityThreats(vulnerabilityID uint, threatId uint) error

	GetThreatRiskExist(threatId uint, riskID uint) (bool, error)
	CreateAssetThreatRisks(risks []*models.RiskThreat, tx ...utility.DatabaseTransaction) error
	DeleteAssetThreatRisks(riskID uint, threatId uint) error

	GetByRiskID(id uint) (*models.RiskRegister, error)
	GetAssetRiskAnalysis(assetID uint) ([]*assetDTOs.VulnerabilityDTO, []*assetDTOs.ThreatDTO, []*assetDTOs.RiskDTO, error)

	GetVulnerabilityIDsByAssetID(assetID uint) ([]uint, error)
	GetThreatIDsByVulnerabilityIDs(vulnIDs []uint) ([]uint, error)
	GetRiskRegisterIDsByThreatIDs(threatIDs []uint) ([]uint, error)

	CreateAssetClauses(clauses []*models.AssetComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteAssetClauses(assetID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error
	GetClausesByAssetId(assetID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
	CreateAssetRisk(risks []*models.RiskAsset, tx ...utility.DatabaseTransaction) error
	DeleteAssetRisk(assetID uint, riskID uint, tx ...utility.DatabaseTransaction) error

	GetAllAssetDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateAssetDocument(risks []*models.AssetAssociatedDocument, tx ...utility.DatabaseTransaction) error
	DeleteAssetDocument(assetID uint, documentID uint, tx ...utility.DatabaseTransaction) error

	CreateAssetControl(controls []*models.AssetControl, tx ...utility.DatabaseTransaction) error
	DeleteAssetControl(assetID uint, controlID uint, tx ...utility.DatabaseTransaction) error

	GetAllAssetProcess(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateAssetProcess(risks []*models.ProcessAsset, tx ...utility.DatabaseTransaction) error
	DeleteAssetProcess(assetID uint, processID uint, tx ...utility.DatabaseTransaction) error
	GetAllAssetStrategicObjectives(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateAssetStrategicObjective(objectives []*models.AssetStrategicObjective, tx ...utility.DatabaseTransaction) error
	DeleteAssetStrategicObjective(assetID uint, strategicObjectiveID uint, tx ...utility.DatabaseTransaction) error

	GetAllAssetKPI(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateAssetKPI(kpis []*models.KPIAsset, tx ...utility.DatabaseTransaction) error
	DeleteAssetKPI(assetID uint, kpiID uint, tx ...utility.DatabaseTransaction) error
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

func (repo *AssetAssociatedItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetAssociatedItemRepository) GetAllControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO
	var totalCount int64

	query := repo.db.Table("asset_controls").
		Select(`
		DISTINCT controls.id AS control_id,
		controls.*,
		CONVERT(VARCHAR, controls.last_test_date, 23) AS control_last_test_date,
		control_types.title_en AS control_type_title_en,
		control_types.title_ar AS control_type_title_ar,
		control_categories.title_en AS control_category_title_en,
		control_categories.title_ar AS control_category_title_ar,
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
		review_frequencies.title_en AS review_frequency_title_en,
		review_frequencies.title_ar AS review_frequency_title_ar,
		test_responsible_user.id AS test_responsible_user_id,
		test_responsible_user.first_name_en AS test_responsible_user_first_name_en,
		test_responsible_user.first_name_ar AS test_responsible_user_first_name_ar,
		test_responsible_user.last_name_en AS test_responsible_user_last_name_en,
		test_responsible_user.last_name_ar AS test_responsible_user_last_name_ar,
		test_responsible_user.initial AS test_responsible_user_initial,
		test_responsible_user.color AS test_responsible_user_color,
		test_responsible_user.email AS test_responsible_user_email,
		test_responsible_user.profile_image_path AS test_responsible_user_profile_image_path,
		test_frequencies.title_en AS test_frequency_title_en,
		test_frequencies.title_ar AS test_frequency_title_ar,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		control_implementation_statuses.title_en AS control_implementation_status_title_en,
		control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
		control_implementation_statuses.color_code AS control_implementation_status_color_code,
		control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
		control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
		control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		created_user.id AS created_user_id,
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
		LEFT JOIN controls ON controls.id = asset_controls.control_id
		LEFT JOIN control_types ON controls.control_type_id = control_types.id
		LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
		LEFT JOIN organization_units ON controls.responsible_unit_id = organization_units.id
		LEFT JOIN users AS owner ON controls.responsible_user_id = owner.id
		LEFT JOIN review_frequencies ON controls.review_frequency_id = review_frequencies.id
		LEFT JOIN users AS test_responsible_user ON controls.test_responsible_user_id = test_responsible_user.id
		LEFT JOIN review_frequencies AS test_frequencies ON controls.test_frequency_id = test_frequencies.id
		LEFT JOIN lifecycle_statuses ON controls.lifecycle_status_id = lifecycle_statuses.id
		LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
		LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
		LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = controls.risk_matrix_risk_rating_id
		LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = controls.risk_matrix_risk_rating_value_id
		LEFT JOIN users AS created_user ON created_user.id = controls.created_by_id
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
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedControlListResponseDTO{
		Items: controlList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
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

func (repo *AssetAssociatedItemRepository) GetAllRisksByRiskIDs(riskRegisterIDs []uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var totalCount int64
	var riskList []*commonDTOs.AssociatedRiskListDTO

	query := repo.db.Table("risk_register").
		Select(`
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
		Where("risk_register.id IN ?", riskRegisterIDs)

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

func (repo *AssetAssociatedItemRepository) GetAllClauses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {

	var clauseList []commonDTOs.AssociatedClauseRowDTO
	var totalCount int64
	query := repo.db.Table("asset_compliance_requirements").
		Select(`
			DISTINCT asset_compliance_requirements.id,
			compliance_clauses.id AS compliance_clause_id,
			compliance_clauses.clause_no,
			compliance_clauses.title_en AS clause_title_en,
			compliance_clauses.title_ar AS clause_title_ar,
			compliance_register.id AS compliance_register_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar
		`).
		Joins("LEFT JOIN compliance_clauses ON compliance_clauses.id = asset_compliance_requirements.compliance_clause_id").
		Joins("LEFT JOIN compliance_register ON compliance_register.id = compliance_clauses.compliance_register_id").
		Where("asset_compliance_requirements.asset_register_id = ?", assetID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(
			compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_clauses.clause_no LIKE ? OR
			compliance_clauses.title_en  LIKE ? OR
			compliance_clauses.title_ar  LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword,
		)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_register.title_en"
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&clauseList).Error; err != nil {
		return nil, err
	}

	clausesByComplianceRegister := make(map[uint]*commonDTOs.AssociatedClauseListDTO)
	for _, row := range clauseList {
		if _, exists := clausesByComplianceRegister[row.ComplianceRegisterID]; !exists {
			clausesByComplianceRegister[row.ComplianceRegisterID] = &commonDTOs.AssociatedClauseListDTO{
				ComplianceRegisterID:      row.ComplianceRegisterID,
				ComplianceRegisterTitleEn: row.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: row.ComplianceRegisterTitleAr,
				Clauses:                   []commonDTOs.AssociatedClauseDTO{},
			}
		}
		clausesByComplianceRegister[row.ComplianceRegisterID].Clauses = append(clausesByComplianceRegister[row.ComplianceRegisterID].Clauses, commonDTOs.AssociatedClauseDTO{
			ComplianceClauseID: row.ComplianceClauseID,
			ClauseNo:           row.ClauseNo,
			ClauseTitleEn:      row.ClauseTitleEn,
			ClauseTitleAr:      row.ClauseTitleAr,
		})
	}

	var groupedClauses []*commonDTOs.AssociatedClauseListDTO
	for _, v := range clausesByComplianceRegister {
		groupedClauses = append(groupedClauses, v)
	}

	if groupedClauses == nil {
		groupedClauses = []*commonDTOs.AssociatedClauseListDTO{}
	}

	totalCount = int64(len(groupedClauses))
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	start := (params.Page - 1) * params.PageSize
	if start > int(totalCount) {
		start = int(totalCount)
	}

	end := start + params.PageSize
	if end > int(totalCount) {
		end = int(totalCount)
	}

	pagedClauses := groupedClauses[start:end]
	response := &commonDTOs.AssociatedItemClauseListResponseDTO{
		Items: pagedClauses,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *AssetAssociatedItemRepository) CreateAssetVulnerabilities(vulnerabilities []*models.VulnerabilityAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(vulnerabilities).Error; err != nil {
		return err
	}
	return nil
}
func (repo *AssetAssociatedItemRepository) DeleteAssetVulnerability(assetID, vulnerabilityID uint) error {
	if err := repo.db.
		Where("asset_register_id = ? AND vulnerability_id = ?", assetID, vulnerabilityID).
		Delete(&models.VulnerabilityAsset{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetAssociatedItemRepository) CreateAssetVulnerabilityThreats(threats []*models.VulnerabilityThreat, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(threats).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetAssociatedItemRepository) DeleteVulnerabilityThreats(vulnerabilityID uint, threatId uint) error {
	if err := repo.db.
		Where("vulnerability_id = ? AND threat_id = ?", vulnerabilityID, threatId).
		Delete(&models.VulnerabilityThreat{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetAssociatedItemRepository) CreateAssetThreatRisks(risks []*models.RiskThreat, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(risks).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetAssociatedItemRepository) DeleteAssetThreatRisks(riskID uint, threatId uint) error {
	if err := repo.db.
		Where("risk_register_id = ? AND threat_id = ?", riskID, threatId).
		Delete(&models.RiskThreat{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetAssociatedItemRepository) GetByRiskID(id uint) (*models.RiskRegister, error) {
	var riskRegister *models.RiskRegister
	result := repo.db.First(&riskRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return riskRegister, nil
}

func (repo *AssetAssociatedItemRepository) GetAssetRiskAnalysis(assetID uint) ([]*assetDTOs.VulnerabilityDTO, []*assetDTOs.ThreatDTO, []*assetDTOs.RiskDTO, error) {
	var vulnerabilities []*assetDTOs.VulnerabilityDTO
	var threats []*assetDTOs.ThreatDTO
	var risks []*assetDTOs.RiskDTO
	// Step 1: Fetch Vulnerabilities linked to the asset
	err := repo.db.Table("vulnerability_assets").
		Select(`
            vulnerabilities.id,
            vulnerabilities.reference_number,
            vulnerabilities.title_en,
            vulnerabilities.title_ar,
            vulnerabilities.description,
            vulnerabilities.source,
            vulnerabilities.severity_level_id,
            severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
            severity_levels.color_code AS severity_level_color_code,
            vulnerabilities.vulnerability_status_id,
            vulnerability_statuses.title_en AS vulnerability_status_title_en,
            vulnerability_statuses.title_ar AS vulnerability_status_title_ar,
            vulnerability_statuses.color_code AS vulnerability_status_color_code,
            vulnerability_statuses.type AS vulnerability_status_type,
            vulnerability_types.title_en AS vulnerability_type_title_en,
            vulnerability_types.title_ar AS vulnerability_type_title_ar,
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
		Joins("LEFT JOIN vulnerabilities ON vulnerabilities.id = vulnerability_assets.vulnerability_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = vulnerabilities.severity_level_id").
		Joins("LEFT JOIN vulnerability_statuses ON vulnerability_statuses.id = vulnerabilities.vulnerability_status_id").
		Joins("LEFT JOIN vulnerability_types ON vulnerability_types.id = vulnerabilities.vulnerability_type_id").
		Joins("LEFT JOIN users ON users.id = vulnerabilities.responsible_user_id").
		Where("vulnerability_assets.asset_register_id = ?", assetID).
		Find(&vulnerabilities).Error

	if err != nil {
		return nil, nil, nil, err
	}

	err = repo.db.Table("vulnerability_threats").
		Select(`
            vulnerability_threats.vulnerability_id,
            vulnerability_threats.id,
            threats.title_en,
            threats.title_ar,
            threat_types.title_en AS threat_type_title_en,
            threat_types.title_ar AS threat_type_title_ar,
            threats.severity_level_id,
            severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
            severity_levels.color_code AS severity_level_color_code,
            threats.threat_status_id,
            threat_statuses.title_en AS threat_status_title_en,
            threat_statuses.title_ar AS threat_status_title_ar,
            threat_statuses.color_code AS threat_status_color_code,
            threat_statuses.type AS threat_status_type,
            threats.source,
            threats.description,
            threats.detection_date,
            threats.resolution_date,
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
		Joins("LEFT JOIN threats ON threats.id = vulnerability_threats.threat_id").
		Joins("LEFT JOIN threat_types ON threat_types.id = threats.threat_type_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = threats.severity_level_id").
		Joins("LEFT JOIN threat_statuses ON threat_statuses.id = threats.threat_status_id").
		Joins("LEFT JOIN users ON users.id = threats.responsible_user_id").
		Where("vulnerability_threats.vulnerability_id IN (?)", getVulnerabilityIDs(vulnerabilities)).
		Find(&threats).Error

	if err != nil {
		return nil, nil, nil, err
	}

	// Step 3: Fetch Risks related to these threats
	err = repo.db.Table("risk_threats").
		Select(`
        risk_threats.threat_id,
        risk_threats.id,
        risk_register.reference_number,
        risk_register.title,
        risk_register.description,
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

        inherent_likelihoods.title_en AS inherent_likelihood_title_en,
        inherent_likelihoods.title_ar AS inherent_likelihood_title_ar,
        inherent_impacts.title_en AS inherent_impact_title_en,
        inherent_impacts.title_ar AS inherent_impact_title_ar,
        inherent_risk_ratings.title_en AS inherent_risk_rating_title_en,
        inherent_risk_ratings.title_ar AS inherent_risk_rating_title_ar,
        inherent_risk_ratings.color_code AS inherent_risk_rating_color_code,

        residual_likelihoods.title_en AS residual_likelihood_title_en,
        residual_likelihoods.title_ar AS residual_likelihood_title_ar,
        residual_impacts.title_en AS residual_impact_title_en,
        residual_impacts.title_ar AS residual_impact_title_ar,
        residual_risk_ratings.title_en AS residual_risk_rating_title_en,
        residual_risk_ratings.title_ar AS residual_risk_rating_title_ar,
        residual_risk_ratings.color_code AS residual_risk_rating_color_code,

        review_frequencies.title_en AS review_frequency_title_en,
        review_frequencies.title_ar AS review_frequency_title_ar,
        CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
        risk_register.last_review_note
    `).
		Joins("LEFT JOIN risk_register ON risk_register.id = risk_threats.risk_register_id").
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN users ON users.id = risk_register.responsible_user_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS inherent_likelihoods ON inherent_likelihoods.id = risk_register.inherent_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS inherent_impacts ON inherent_impacts.id = risk_register.inherent_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_ratings ON inherent_risk_ratings.id = risk_register.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS residual_likelihoods ON residual_likelihoods.id = risk_register.residual_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS residual_impacts ON residual_impacts.id = risk_register.residual_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS residual_risk_ratings ON residual_risk_ratings.id = risk_register.residual_risk_rating_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id").
		Where("risk_threats.threat_id IN (?)", getThreatIDs(threats)).
		Find(&risks).Error

	if err != nil {
		return nil, nil, nil, err
	}

	return vulnerabilities, threats, risks, nil
}

// Helper function to extract vulnerability IDs
func getVulnerabilityIDs(vulnerabilities []*assetDTOs.VulnerabilityDTO) []uint {
	var ids []uint
	for _, v := range vulnerabilities {
		ids = append(ids, v.ID)
	}
	return ids
}

// Helper function to extract threat IDs
func getThreatIDs(threats []*assetDTOs.ThreatDTO) []uint {
	var ids []uint
	for _, t := range threats {
		ids = append(ids, t.ID)
	}
	return ids
}

func (repo *AssetAssociatedItemRepository) GetVulnerabilityAssetExist(assetID, vulnerabilityID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.VulnerabilityAsset{}).
		Where("asset_register_id = ? AND vulnerability_id = ?", assetID, vulnerabilityID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *AssetAssociatedItemRepository) GetAssetComplianceClauseExist(assetID uint, clauseID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.AssetComplianceRequirement{}).
		Where("asset_register_id = ? AND compliance_clause_id = ?", assetID, clauseID).
		Count(&count).Error
	return count > 0, err
}

func (repo *AssetAssociatedItemRepository) CreateAssetComplianceClauses(clauses []*models.AssetComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(clauses).Error; err != nil {
		return err
	}
	return nil
}

func (r *AssetAssociatedItemRepository) GetAssetComplianceClauseByID(id uint) (*models.AssetComplianceRequirement, error) {
	var record models.AssetComplianceRequirement
	err := r.db.First(&record, id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *AssetAssociatedItemRepository) DeleteAssetComplianceClause(assetID uint, id uint) error {
	return r.db.
		Where("id = ? AND asset_register_id = ?", id, assetID).
		Delete(&models.AssetComplianceRequirement{}).Error
}

func (repo *AssetAssociatedItemRepository) GetVulnerabilityThreatExist(vulnerabilityId uint, threatID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.VulnerabilityThreat{}).
		Where("vulnerability_id = ? AND threat_id = ?", vulnerabilityId, threatID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssetAssociatedItemRepository) GetThreatRiskExist(threatId uint, riskID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.RiskAsset{}).
		Where("risk_register_id = ? AND asset_register_id = ?", threatId, riskID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssetAssociatedItemRepository) GetVulnerabilityIDsByAssetID(assetID uint) ([]uint, error) {
	var ids []uint
	err := repo.db.
		Model(&models.VulnerabilityAsset{}).
		Where("asset_register_id = ?", assetID).
		Pluck("vulnerability_id", &ids).Error
	return ids, err
}

func (repo *AssetAssociatedItemRepository) GetThreatIDsByVulnerabilityIDs(vulnIDs []uint) ([]uint, error) {
	var ids []uint
	err := repo.db.
		Model(&models.VulnerabilityThreat{}).
		Where("vulnerability_id IN ?", vulnIDs).
		Pluck("threat_id", &ids).Error
	return ids, err
}

func (repo *AssetAssociatedItemRepository) GetRiskRegisterIDsByThreatIDs(threatIDs []uint) ([]uint, error) {
	var ids []uint
	err := repo.db.
		Model(&models.RiskThreat{}).
		Where("threat_id IN ?", threatIDs).
		Pluck("risk_register_id", &ids).Error
	return ids, err
}

func (repo *AssetAssociatedItemRepository) GetClausesByAssetId(assetID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clauseIDs []uint
	if err := db.
		Model(&models.AssetComplianceRequirement{}).
		Where("asset_register_id = ?", assetID).
		Pluck("compliance_clause_id", &clauseIDs).Error; err != nil {
		return nil, err
	}

	return clauseIDs, nil
}

func (repo *AssetAssociatedItemRepository) CreateAssetClauses(clauses []*models.AssetComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&clauses).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetAssociatedItemRepository) DeleteAssetClauses(assetID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.
		Where("asset_register_id = ? AND compliance_clause_id IN ?", assetID, clauseIDs).
		Delete(&models.AssetComplianceRequirement{}).Error; err != nil {
		return err
	}

	return nil
}
func (repo *AssetAssociatedItemRepository) CreateAssetRisk(risks []*models.RiskAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&risks).Error
}

func (repo *AssetAssociatedItemRepository) DeleteAssetRisk(assetID uint, riskId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND asset_register_id = ?", riskId, assetID).
		Delete(&models.RiskAsset{}).
		Error
}

func (repo *AssetAssociatedItemRepository) GetAllAssetDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("asset_associated_documents").
		Select(`
			asset_associated_documents.id AS id,
			asset_associated_documents.document_id,
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
			LEFT JOIN documents ON documents.id = asset_associated_documents.document_id
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
		Where("asset_associated_documents.asset_register_id = ?", assetID)

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

func (repo *AssetAssociatedItemRepository) CreateAssetDocument(documents []*models.AssetAssociatedDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (repo *AssetAssociatedItemRepository) DeleteAssetDocument(assetID uint, documentId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("document_id = ? AND asset_register_id = ?", documentId, assetID).
		Delete(&models.AssetAssociatedDocument{}).
		Error
}

func (repo *AssetAssociatedItemRepository) CreateAssetControl(controls []*models.AssetControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *AssetAssociatedItemRepository) DeleteAssetControl(assetID uint, controlId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("control_id = ? AND asset_register_id = ?", controlId, assetID).
		Delete(&models.AssetControl{}).
		Error
}

func (repo *AssetAssociatedItemRepository) GetAllAssetProcess(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	var processList []*commonDTOs.AssociatedProcessListDTO
	var totalCount int64

	query := repo.db.Table("process_assets").
		Select(`
		process_assets.id,
		process_assets.process_id,
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
			LEFT JOIN processes ON processes.id = process_assets.process_id
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
		`).Where("process_assets.asset_register_id = ?", assetID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			processes.reference_number LIKE ? OR
			process_groups.title_en LIKE ? OR
			process_groups.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ? OR
			lifecycle_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
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

func (repo *AssetAssociatedItemRepository) CreateAssetProcess(controls []*models.ProcessAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *AssetAssociatedItemRepository) DeleteAssetProcess(assetID uint, processId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("process_id = ? AND asset_register_id = ?", processId, assetID).
		Delete(&models.ProcessAsset{}).
		Error
}

func (repo *AssetAssociatedItemRepository) GetAllAssetStrategicObjectives(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	var objectiveList []*commonDTOs.AssociatedStrategicObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("asset_strategic_objectives").
		Select(`
			strategic_objectives.*,
			asset_strategic_objectives.strategic_objective_id,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,

			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,

			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

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
			LEFT JOIN strategic_objectives ON strategic_objectives.id = asset_strategic_objectives.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
			LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id
			LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id
        `).Where("asset_strategic_objectives.asset_register_id = ?", assetID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(strategic_objectives.title_en LIKE ? OR
			strategic_objectives.title_ar LIKE ? OR
			strategic_objectives.status LIKE ? )`,
			keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "strategic_objectives.title_en"
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&objectiveList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO{
		Items: objectiveList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *AssetAssociatedItemRepository) CreateAssetStrategicObjective(objectives []*models.AssetStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&objectives).Error
}

func (repo *AssetAssociatedItemRepository) DeleteAssetStrategicObjective(assetID uint, strategicObjectiveID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("strategic_objective_id = ? AND asset_register_id = ?", strategicObjectiveID, assetID).
		Delete(&models.AssetStrategicObjective{}).
		Error
}

func (repo *AssetAssociatedItemRepository) GetAllAssetKPI(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	var kpiRegister []*commonDTOs.AssociatedKPIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kpi_assets").
		Select(`
            kpi_register.*,
			kpi_assets.kpi_register_id,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			CONVERT(VARCHAR, kpi_register.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_register.end_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			calculation_types.title_en AS calculation_type_title_en,
			calculation_types.title_ar AS calculation_type_title_ar,
			calculation_types.type AS calculation_type_type,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
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
		    LEFT JOIN kpi_register ON kpi_register.id = kpi_assets.kpi_register_id
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
        `).Where("kpi_assets.asset_register_id = ?", assetID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(kpi_register.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "kpi_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Fetch the data
	if err := query.Find(&kpiRegister).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedKPIRegisterListResponseDTO{
		Items: kpiRegister,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AssetAssociatedItemRepository) CreateAssetKPI(kpis []*models.KPIAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&kpis).Error
}

func (repo *AssetAssociatedItemRepository) DeleteAssetKPI(assetID uint, kpiID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND asset_register_id = ?", kpiID, assetID).
		Delete(&models.KPIAsset{}).
		Error
}
