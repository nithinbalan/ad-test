package repositories

import (
	"jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type BcpRiskAnalysisRepositoryInterface interface {
	GetAllRisksByBcpId(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error)
	AddRiskToBcp(bcpRisks []models.BcpRisk) error
	DeleteRiskFromBcp(bcpId uint, bcpRiskForm dtos.BcpRiskForm) error
	GetAllProcessRisks(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error)
}

type BcpRiskAnalysisRepository struct {
	db *gorm.DB
}

func NewBcpRiskAnalysisRepository(db *gorm.DB) BcpRiskAnalysisRepositoryInterface {
	return &BcpRiskAnalysisRepository{
		db: db,
	}
}

func (repo *BcpRiskAnalysisRepository) GetAllRisksByBcpId(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error) {
	var riskList []*dtos.BcpRiskListDTO
	var totalCount int64

	query := repo.db.Table("bcp_risks").
		Select(`
			bcp_risks.id,
			bcp_risks.risk_register_id,
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
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = bcp_risks.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
			LEFT JOIN risk_matrix_likelihoods ON risk_matrix_likelihoods.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = risk_register.current_risk_rating_value_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).Where("bcp_risks.business_continuity_plan_id = ?", bcpId)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	if params.Value != 0 {
		query = query.Where("risk_matrix_risk_rating_values.value = ?", params.Value)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.BcpRiskListResponseDTO{
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

func (bpr *BcpRiskAnalysisRepository) AddRiskToBcp(bcpRisks []models.BcpRisk) error {
	return bpr.db.Create(bcpRisks).Error
}

func (bpr *BcpRiskAnalysisRepository) DeleteRiskFromBcp(bcpId uint, bcpRiskForm dtos.BcpRiskForm) error {
	return bpr.db.Where("business_continuity_plan_id = ? AND risk_register_id IN (?)", bcpId, bcpRiskForm.RiskRegisterIDs).Delete(&models.BcpRisk{}).Error
}

func (repo *BcpRiskAnalysisRepository) GetAllProcessRisks(bcpId uint, params *dtos.BcpRiskListQueryParamsDTO) (*dtos.BcpRiskListResponseDTO, error) {
	var riskList []*dtos.BcpRiskListDTO
	var totalCount int64

	query := repo.db.Table("bcp_processes").
		Select(`		
			bcp_processes.process_id,
			processes.title_en as process_title_en,
			processes.title_ar as process_title_ar,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			risk_register.id as risk_register_id,
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
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			INNER JOIN risk_processes ON risk_processes.process_id = bcp_processes.process_id
			INNER JOIN risk_register ON risk_register.id = risk_processes.risk_register_id
			LEFT JOIN processes ON processes.id = bcp_processes.process_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
			LEFT JOIN risk_matrix_likelihoods ON risk_matrix_likelihoods.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = risk_register.current_risk_rating_value_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).
		Where("bcp_processes.business_continuity_plan_id = ?", bcpId)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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
			risk_statuses.color_code LIKE ? OR
			processes.title_en LIKE ? OR
			processes.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.Value != 0 {
		query = query.Where("risk_matrix_risk_rating_values.value = ?", params.Value)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.BcpRiskListResponseDTO{
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
