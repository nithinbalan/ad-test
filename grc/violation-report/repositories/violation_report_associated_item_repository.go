package repositories

import (
	commonDTOs "jethur/grc/common/dtos"

	"jethur/utility"

	"jethur/grc/models"

	"gorm.io/gorm"
)

type ViolationReportAssociatedItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllControls(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	GetAllComplianceRegisters(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedComplianceListResponseDTO, error)
	GetAllRisks(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllDocuments(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	GetAllAssets(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateViolationReportRisk(risks []*models.ViolationReportRisk, tx ...utility.DatabaseTransaction) error
	DeleteViolationReportRisk(violationReportID uint, riskID uint, tx ...utility.DatabaseTransaction) error
	CreateViolationReportDocument(documents []*models.ViolationReportAssociatedDocument, tx ...utility.DatabaseTransaction) error
	DeleteViolationReportDocument(violationReportID uint, documentID uint, tx ...utility.DatabaseTransaction) error
	CreateViolationReportControl(controls []*models.ViolationReportControl, tx ...utility.DatabaseTransaction) error
	DeleteViolationReportControl(violationReportID uint, controlID uint, tx ...utility.DatabaseTransaction) error
	CreateViolationReportAsset(assets []*models.ViolationReportAsset, tx ...utility.DatabaseTransaction) error
	DeleteViolationReportAsset(violationReportID uint, assetID uint, tx ...utility.DatabaseTransaction) error
	CreateViolationReportComplianceRegister(complianceRegisters []*models.ViolationReportComplianceRegister, tx ...utility.DatabaseTransaction) error
	DeleteViolationReportComplianceRegister(violationReportID uint, complianceRegisterID uint, tx ...utility.DatabaseTransaction) error
}

type ViolationReportAssociatedItemRepository struct {
	db *gorm.DB
}

func NewViolationReportAssociatedItemRepository(db *gorm.DB) ViolationReportAssociatedItemRepositoryInterface {
	return &ViolationReportAssociatedItemRepository{
		db: db,
	}
}
func (repo *ViolationReportAssociatedItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}
func (repo *ViolationReportAssociatedItemRepository) GetAllControls(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO
	var totalCount int64

	query := repo.db.Table("violation_report_controls").
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
		LEFT JOIN controls ON controls.id = violation_report_controls.control_id
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
	`).Where("violation_report_controls.violation_report_id = ?", violationReportID)

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

func (repo *ViolationReportAssociatedItemRepository) GetAllComplianceRegisters(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedComplianceListResponseDTO, error) {
	var complianceList []*commonDTOs.AssociatedComplianceListDTO
	var totalCount int64

	query := repo.db.Table("violation_report_compliance_registers").
		Select(`
			violation_report_compliance_registers.id,
			violation_report_compliance_registers.compliance_register_id,
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
			compliance_register.reference_number,
			compliance_register.standard_code,
			compliance_register.title_en,
			compliance_register.title_ar,
			compliance_register.description_en,
			compliance_register.description_ar,
			compliance_register.version,
			CONVERT(VARCHAR, compliance_register.deadline, 23) AS deadline,
			CONVERT(VARCHAR, compliance_register.issue_date, 23) AS issue_date,
			CONVERT(VARCHAR, compliance_register.expiry_date, 23) AS expiry_date,
			CONVERT(VARCHAR, compliance_register.next_review_date, 23) AS next_review_date,
			CONVERT(VARCHAR, compliance_register.last_reviewed_date, 23) AS last_reviewed_date,
			CONVERT(VARCHAR, compliance_register.created_at, 23) AS created_at,
			CONVERT(VARCHAR, compliance_register.updated_at, 23) AS updated_at,
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
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			compliance_register.compliance_assessment_framework_id,
			compliance_assessment_frameworks.title_en AS compliance_assessment_framework_title_en,
			compliance_assessment_frameworks.title_ar AS compliance_assessment_framework_title_ar
        `).
		Joins(`
			LEFT JOIN compliance_register ON compliance_register.id = violation_report_compliance_registers.compliance_register_id
			LEFT JOIN compliance_categories ON compliance_register.compliance_category_id = compliance_categories.id
			LEFT JOIN compliance_statuses ON compliance_register.compliance_status_id = compliance_statuses.id
			LEFT JOIN organization_units ON compliance_register.responsible_unit_id = organization_units.id
			LEFT JOIN users AS responsible_user ON compliance_register.responsible_user_id = responsible_user.id
			LEFT JOIN review_frequencies ON compliance_register.review_frequency_id = review_frequencies.id
			LEFT JOIN users AS review_responsible_user ON compliance_register.review_responsible_user_id = review_responsible_user.id
			LEFT JOIN severity_levels ON compliance_register.severity_level_id = severity_levels.id
			LEFT JOIN users AS created_user ON compliance_register.created_by_id = created_user.id
			LEFT JOIN users AS updated_user ON compliance_register.updated_by_id = updated_user.id
			LEFT JOIN lifecycle_statuses ON compliance_register.lifecycle_status_id = lifecycle_statuses.id
			LEFT JOIN compliance_assessment_frameworks ON compliance_register.compliance_assessment_framework_id = compliance_assessment_frameworks.id
        `).
		Where("violation_report_compliance_registers.violation_report_id = ?", violationReportID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_register.description_en LIKE ? OR
			compliance_register.description_ar LIKE ? OR
			compliance_register.reference_number LIKE ? OR
			compliance_register.standard_code LIKE ? OR
			compliance_categories.title_en LIKE ? OR
			compliance_categories.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_register.title_en"
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&complianceList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedComplianceListResponseDTO{
		Items: complianceList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ViolationReportAssociatedItemRepository) GetAllRisks(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var incidentRisk []*commonDTOs.AssociatedRiskListDTO
	var totalCount int64

	query := repo.db.Table("violation_report_risks").
		Select(`
		    violation_report_risks.id,
			violation_report_risks.risk_register_id,
			violation_report_risks.violation_report_id,
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
			LEFT JOIN risk_register ON risk_register.id = violation_report_risks.risk_register_id
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
		Where("violation_report_risks.violation_report_id = ?", violationReportID)

	// Calculate total count for pagination

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.reference_number LIKE ? OR
			risk_register.title LIKE ? OR
			risk_register.description LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			risk_sources.title_en LIKE ? OR
			risk_sources.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_en LIKE ? OR
			users.last_name_ar LIKE ? OR
			risk_statuses.title_en LIKE ? OR
			risk_statuses.title_ar LIKE ? OR
			risk_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	if params.Value != 0 {
		query = query.Where("current_risk_rating_value.value = ?", params.Value)
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

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentRisk).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedRiskListResponseDTO{
		Items: incidentRisk,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ViolationReportAssociatedItemRepository) GetAllDocuments(incidentID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("violation_report_associated_documents").
		Select(`
			violation_report_associated_documents.id AS id,
			violation_report_associated_documents.document_id,
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
			LEFT JOIN documents ON documents.id = violation_report_associated_documents.document_id
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
		`).Where("violation_report_associated_documents.violation_report_id = ?", incidentID)

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

func (repo *ViolationReportAssociatedItemRepository) GetAllAssets(violationReportID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	var assetList []*commonDTOs.AssociatedAssetListDTO
	var totalCount int64

	query := repo.db.Table("violation_report_assets").
		Select(`
			violation_report_assets.id,
			violation_report_assets.asset_register_id,
			asset_register.reference_number,
			asset_register.title_en,
			asset_register.title_ar,
			asset_register.location,
			asset_register.description,
			CONVERT(VARCHAR, asset_register.acquisition_date, 23) AS acquisition_date,
			CONVERT(VARCHAR, asset_register.assignment_date, 23) AS assignment_date,
			CONVERT(VARCHAR, asset_register.warranty_expiry_date, 23) AS warranty_expiry_date,
			asset_register.manufacturer,
			asset_register.model,
			asset_register.serial_number,
			asset_categories.title_en AS asset_category_title_en,
			asset_categories.title_ar AS asset_category_title_ar,
			asset_sub_categories.title_en AS asset_sub_category_title_en,
			asset_sub_categories.title_ar AS asset_sub_category_title_ar,
			asset_statuses.title_en AS asset_status_title_en,
			asset_statuses.title_ar AS asset_status_title_ar,
			asset_statuses.color_code AS asset_status_color_code,
			asset_register.total_score,
			asset_criticality_matrix_rating_scales.title_en AS asset_criticality_matrix_rating_scale_title_en,
			asset_criticality_matrix_rating_scales.title_ar AS asset_criticality_matrix_rating_scale_title_ar,
			asset_criticality_matrix_rating_scales.color_code AS asset_criticality_matrix_rating_scale_color_code,
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
		Joins("LEFT JOIN asset_register ON asset_register.id = violation_report_assets.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN asset_criticality_matrix_rating_scales ON asset_criticality_matrix_rating_scales.id = asset_register.asset_criticality_matrix_rating_scale_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id").
		Where("violation_report_assets.violation_report_id = ?", violationReportID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			asset_register.reference_number LIKE ? OR
			asset_register.description LIKE ? OR
			asset_categories.title_en LIKE ? OR
			asset_categories.title_ar LIKE ? OR
			asset_statuses.title_en LIKE ? OR
			asset_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "asset_register.title_en"
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&assetList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemAssetListResponseDTO{
		Items: assetList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ViolationReportAssociatedItemRepository) CreateViolationReportRisk(risks []*models.ViolationReportRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&risks).Error
}

func (repo *ViolationReportAssociatedItemRepository) DeleteViolationReportRisk(violationReportID uint, riskId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("risk_register_id = ? AND violation_report_id = ?", riskId, violationReportID).
		Delete(&models.ViolationReportRisk{}).
		Error
}

func (repo *ViolationReportAssociatedItemRepository) CreateViolationReportDocument(documents []*models.ViolationReportAssociatedDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (repo *ViolationReportAssociatedItemRepository) DeleteViolationReportDocument(violationReportID uint, documentId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("document_id = ? AND violation_report_id = ?", documentId, violationReportID).
		Delete(&models.ViolationReportAssociatedDocument{}).
		Error
}
func (repo *ViolationReportAssociatedItemRepository) CreateViolationReportControl(controls []*models.ViolationReportControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *ViolationReportAssociatedItemRepository) DeleteViolationReportControl(violationReportID uint, controlId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("control_id = ? AND violation_report_id = ?", controlId, violationReportID).
		Delete(&models.ViolationReportControl{}).
		Error
}

func (repo *ViolationReportAssociatedItemRepository) CreateViolationReportAsset(assets []*models.ViolationReportAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&assets).Error
}

func (repo *ViolationReportAssociatedItemRepository) DeleteViolationReportAsset(violationReportID uint, assetID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("asset_register_id = ? AND violation_report_id = ?", assetID, violationReportID).
		Delete(&models.ViolationReportAsset{}).
		Error
}

func (repo *ViolationReportAssociatedItemRepository) CreateViolationReportComplianceRegister(complianceRegisters []*models.ViolationReportComplianceRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&complianceRegisters).Error
}

func (repo *ViolationReportAssociatedItemRepository) DeleteViolationReportComplianceRegister(violationReportID uint, complianceRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("id = ? AND violation_report_id = ?", complianceRegisterID, violationReportID).
		Delete(&models.ViolationReportComplianceRegister{}).
		Error
}
