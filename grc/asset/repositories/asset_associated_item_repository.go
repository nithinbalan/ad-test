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

	GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) ([]*assetDTOs.AssetControlListDTO, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
	GetAllRisks(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemRiskListResponseDTO, error)
	GetAllVulnerabilities(assetID uint) ([]models.Vulnerability, error)
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
            users.profile_image_path AS responsible_user_profile_image_path
        `).
		Joins("LEFT JOIN risk_register ON risk_register.id = risk_threats.risk_register_id").
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN users ON users.id = risk_register.responsible_user_id").
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
