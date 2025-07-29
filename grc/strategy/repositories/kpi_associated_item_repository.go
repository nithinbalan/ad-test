package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type KPIAssociatedItemRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllControls(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemControlListResponseDTO, error)
	CreateKPIControls(controls []*models.KPIControl, tx ...utility.DatabaseTransaction) error
	DeleteKPIControlMapping(kpiID, controlID uint, tx ...utility.DatabaseTransaction) error
	IsControlMappedToKPI(kpiID, controlID uint) (bool, error)

	GetAllProcesses(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemProcessListResponseDTO, error)
	CreateKPIProcesses(processes []*models.KPIProcess, tx ...utility.DatabaseTransaction) error
	DeleteKPIProcessMapping(kpiID, processID uint, tx ...utility.DatabaseTransaction) error
	IsProcessMappedToKPI(kpiID, processID uint) (bool, error)

	GetAllCompliances(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemComplianceListResponseDTO, error)
	CreateKPICompliances(compliances []*models.KPICompliance, tx ...utility.DatabaseTransaction) error
	DeleteKPIComplianceMapping(kpiID, complianceRegisterID uint, tx ...utility.DatabaseTransaction) error
	IsComplianceMappedToKPI(kpiID, complianceRegisterID uint) (bool, error)

	GetAllAssets(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemAssetListResponseDTO, error)
	CreateKPIAssets(assets []*models.KPIAsset, tx ...utility.DatabaseTransaction) error
	DeleteKPIAssetMapping(kpiID, assetRegisterID uint, tx ...utility.DatabaseTransaction) error
	IsAssetMappedToKPI(kpiID, assetRegisterID uint) (bool, error)

	GetAllRisks(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateKPIRisks(risks []*models.KPIRisk, tx ...utility.DatabaseTransaction) error
	DeleteKPIRiskMapping(kpiID, riskRegisterID uint, tx ...utility.DatabaseTransaction) error
	IsRiskMappedToKPI(kpiID, riskRegisterID uint) (bool, error)

	GetAllThirdParties(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemThirdPartyListResponseDTO, error)
	CreateKPIThirdParties(thirdParties []*models.KPIThirdParty, tx ...utility.DatabaseTransaction) error
	DeleteKPIThirdPartyMapping(kpiID, ThirdPartyRegisterID uint, tx ...utility.DatabaseTransaction) error
	IsThirdPartyMappedToKPI(kpiID, ThirdPartyRegisterID uint) (bool, error)

	GetAllPolicies(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateKPIPolicies(policies []*models.KPIAssociatedDocument, tx ...utility.DatabaseTransaction) error
	DeleteKPIPolicyMapping(kpiID, documentID uint, tx ...utility.DatabaseTransaction) error
	// IsPolicyMappedToKPI(kpiID, documentID uint) (bool, error)

	GetAllClauses(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error)
	CreateKPIClauses(clauses []*models.KPICompliance, tx ...utility.DatabaseTransaction) error
	DeleteKPIClauses(kpiID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error
	GetClausesByKPIID(kpiID uint, tx ...utility.DatabaseTransaction) ([]uint, error)
}

type KPIAssociatedItemRepository struct {
	db *gorm.DB
}

func NewKPIAssociatedItemRepository(db *gorm.DB) *KPIAssociatedItemRepository {
	return &KPIAssociatedItemRepository{
		db: db,
	}
}

func (repo *KPIAssociatedItemRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KPIAssociatedItemRepository) GetAllControls(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemControlListResponseDTO, error) {
	var controlList []*strategyDTOs.KPIControlListDTO
	var totalCount int64

	query := repo.db.Table("kpi_controls").
		Select(`
			DISTINCT controls.id AS control_id,
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
			LEFT JOIN controls ON controls.id = kpi_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
			LEFT JOIN users AS owner ON controls.responsible_user_id = owner.id
        `).Where("kpi_controls.kpi_register_id = ?", kpiID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &strategyDTOs.KPIAssociatedItemControlListResponseDTO{
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

func (repo *KPIAssociatedItemRepository) CreateKPIControls(controls []*models.KPIControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIControlMapping(kpiID, controlID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND control_id = ?", kpiID, controlID).
		Delete(&models.KPIControl{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsControlMappedToKPI(kpiID, controlID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPIControl{}).
		Where("kpi_register_id = ? AND control_id = ?", kpiID, controlID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *KPIAssociatedItemRepository) GetAllProcesses(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemProcessListResponseDTO, error) {
	var processList []*strategyDTOs.KPIProcessListDTO
	var totalCount int64

	query := repo.db.Table("kpi_processes").
		Select(`
			processes.id AS process_id,
			processes.reference_number,
			processes.title_en,
			processes.title_ar,
			processes.description,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS control_owner_initial,
            owner.color AS control_owner_color,
            owner.email AS control_owner_email,
            owner.profile_image_path AS control_owner_profile_image_path,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code
		`).
		Joins(`
			LEFT JOIN processes ON processes.id = kpi_processes.process_id
			LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
		`).Where("kpi_processes.kpi_register_id = ?", kpiID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&processList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &strategyDTOs.KPIAssociatedItemProcessListResponseDTO{
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

func (repo *KPIAssociatedItemRepository) CreateKPIProcesses(controls []*models.KPIProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&controls).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIProcessMapping(kpiID, processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND process_id = ?", kpiID, processID).
		Delete(&models.KPIProcess{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsProcessMappedToKPI(kpiID, processID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPIProcess{}).
		Where("kpi_register_id = ? AND process_id = ?", kpiID, processID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *KPIAssociatedItemRepository) GetAllCompliances(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemComplianceListResponseDTO, error) {
	var complianceList []*strategyDTOs.KPIComplianceListDTO
	var totalCount int64

	query := repo.db.Table("kpi_compliances").
		Select(`
			compliance_register.id AS compliance_register_id,
			compliance_register.reference_number AS compliance_register_reference_number,
			compliance_register.standard_code AS compliance_register_standard_code,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar,
			compliance_register.description_en AS compliance_register_description_en,
			compliance_register.description_ar AS compliance_register_description_ar,
			compliance_register.version AS compliance_register_version,
			CONVERT(VARCHAR, compliance_register.deadline, 23) AS deadline,
        	CONVERT(VARCHAR, compliance_register.issue_date, 23) AS issue_date,
        	CONVERT(VARCHAR, compliance_register.expiry_date, 23) AS expiry_date,
       		CONVERT(VARCHAR, compliance_register.next_review_date, 23) AS next_review_date,
			CONVERT(VARCHAR, compliance_register.last_reviewed_date, 23) AS last_reviewed_date,
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
			lifecycle_statuses.color_code AS lifecycle_status_color_code
        `).
		Joins(`
			LEFT JOIN compliance_register ON compliance_register.id = kpi_compliances.compliance_register_id
            LEFT JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id
		    LEFT JOIN organization_units ON organization_units.id = compliance_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_register.responsible_user_id
		    LEFT JOIN review_frequencies ON review_frequencies.id = compliance_register.review_frequency_id
			LEFT JOIN users AS review_responsible_user ON review_responsible_user.id = compliance_register.review_responsible_user_id
		    LEFT JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = compliance_register.lifecycle_status_id
        `).Where("kpi_compliances.kpi_register_id = ?", kpiID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_register.reference_number LIKE ? OR
			compliance_register.description_en LIKE ? OR
			compliance_register.description_ar LIKE ? OR
			compliance_register.standard_code LIKE ? OR
			compliance_categories.title_en LIKE ? OR
			compliance_categories.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_register.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&complianceList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &strategyDTOs.KPIAssociatedItemComplianceListResponseDTO{
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

func (repo *KPIAssociatedItemRepository) CreateKPICompliances(compliance []*models.KPICompliance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&compliance).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIComplianceMapping(kpiID, complianceID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND compliance_register_id = ?", kpiID, complianceID).
		Delete(&models.KPICompliance{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsComplianceMappedToKPI(kpiID, complianceID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPICompliance{}).
		Where("kpi_register_id = ? AND compliance_register_id = ?", kpiID, complianceID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *KPIAssociatedItemRepository) GetAllAssets(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemAssetListResponseDTO, error) {
	var assetList []*strategyDTOs.KPIAssetListDTO
	var totalCount int64

	query := repo.db.Table("kpi_assets").
		Select(`
		asset_register.id AS asset_register_id,
		asset_register.reference_number,
		asset_register.title_en,
		asset_register.title_ar,
		asset_register.location,
		asset_register.description,
		CONVERT(VARCHAR, asset_register.acquisition_date, 23) AS acquisition_date,
		asset_categories.title_en AS asset_category_title_en,
		asset_categories.title_ar AS asset_category_title_ar,
		asset_statuses.title_en AS asset_status_title_en,
		asset_statuses.title_ar AS asset_status_title_ar,
		asset_statuses.color_code AS asset_status_color_code,
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
		Joins("LEFT JOIN asset_register ON asset_register.id = kpi_assets.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id").
		Where("kpi_assets.kpi_register_id = ?", kpiID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&assetList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &strategyDTOs.KPIAssociatedItemAssetListResponseDTO{
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

func (repo *KPIAssociatedItemRepository) CreateKPIAssets(assets []*models.KPIAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&assets).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIAssetMapping(kpiID, assetID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND asset_register_id = ?", kpiID, assetID).
		Delete(&models.KPIAsset{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsAssetMappedToKPI(kpiID, assetID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPIAsset{}).
		Where("kpi_register_id = ? AND asset_register_id = ?", kpiID, assetID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *KPIAssociatedItemRepository) GetAllRisks(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var riskList []*commonDTOs.AssociatedRiskListDTO
	var totalCount int64

	query := repo.db.Table("kpi_risks").
		Select(`
			kpi_risks.id,
			kpi_risks.risk_register_id AS risk_register_id,
			kpi_risks.kpi_register_id,
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
			LEFT JOIN risk_register ON risk_register.id = kpi_risks.risk_register_id
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
		`).Where("kpi_risks.kpi_register_id = ?", kpiID)

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

	response := &commonDTOs.AssociatedRiskListResponseDTO{
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

func (repo *KPIAssociatedItemRepository) CreateKPIRisks(risks []*models.KPIRisk, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&risks).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIRiskMapping(kpiID, riskID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND risk_register_id = ?", kpiID, riskID).
		Delete(&models.KPIRisk{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsRiskMappedToKPI(kpiID, riskID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPIRisk{}).
		Where("kpi_register_id = ? AND risk_register_id = ?", kpiID, riskID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *KPIAssociatedItemRepository) GetAllThirdParties(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemThirdPartyListResponseDTO, error) {
	var thirdPartyList []*strategyDTOs.KPIThirdPartyListDTO
	var totalCount int64

	query := repo.db.Table("kpi_third_parties").
		Select(`
			third_party_register.id AS third_party_register_id,
			third_party_register.title,
			third_party_register.reference_number,
			third_party_register.address,
			third_party_register.website,
			third_party_register.organization_unit_id,
			third_party_register.responsible_user_id,
			third_party_register.phone_number,
			third_party_register.email,
			third_party_register.third_party_type_id,
			third_party_register.representative_name,
			third_party_register.representative_role,
			third_party_register.severity_level_id,
			CONVERT(VARCHAR, third_party_register.onboarding_date, 23) AS onboarding_date,
			CONVERT(VARCHAR, third_party_register.service_start_date, 23) AS service_start_date,
			CONVERT(VARCHAR, third_party_register.service_end_date, 23) AS service_end_date,
			countries.id AS country_id,
			countries.title_en AS country_title_en,
			countries.title_ar AS country_title_ar,
			third_party_types.title_en AS third_party_type_title_en,
			third_party_types.title_ar AS third_party_type_title_ar,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar
        `).
		Joins(`
			LEFT JOIN third_party_register ON third_party_register.id = kpi_third_parties.third_party_register_id
			LEFT JOIN third_party_types ON third_party_types.id = third_party_register.third_party_type_id            
            LEFT JOIN countries ON countries.id = third_party_register.country_id
            LEFT JOIN organization_units ON organization_units.id = third_party_register.organization_unit_id
			LEFT JOIN severity_levels ON severity_levels.id = third_party_register.severity_level_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = third_party_register.responsible_user_id
        `).Where("kpi_third_parties.kpi_register_id = ?", kpiID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			third_party_register.title LIKE ? OR
			third_party_register.reference_number LIKE ? OR
			third_party_register.representative_name LIKE ? OR
			third_party_register.representative_role LIKE ? OR
			third_party_register.email LIKE ? OR
			third_party_register.service_start_date LIKE ? OR
			third_party_register.service_end_date LIKE ? OR
	
			countries.title_en LIKE ? OR
			countries.title_ar LIKE ? OR
	
			third_party_types.title_en LIKE ? OR
			third_party_types.title_ar LIKE ? OR
	
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
	
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
	
			responsible_user.first_name_en LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.email LIKE ?
		`,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword,
			keyword, keyword, keyword,
			keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
		)
	}

	if params.SortBy == "" {
		params.SortBy = "third_party_register.title"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&thirdPartyList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &strategyDTOs.KPIAssociatedItemThirdPartyListResponseDTO{
		Items: thirdPartyList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KPIAssociatedItemRepository) CreateKPIThirdParties(thirdParties []*models.KPIThirdParty, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&thirdParties).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIThirdPartyMapping(kpiID, thirPartyRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND third_party_register_id = ?", kpiID, thirPartyRegisterID).
		Delete(&models.KPIThirdParty{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsThirdPartyMappedToKPI(kpiID, thirPartyRegisterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPIThirdParty{}).
		Where("kpi_register_id = ? AND third_party_register_id = ?", kpiID, thirPartyRegisterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *KPIAssociatedItemRepository) GetAllPolicies(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var policyList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("kpi_associated_documents").
		Select(`
			kpi_associated_documents.id AS id,
			kpi_associated_documents.document_id,
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
			LEFT JOIN documents ON documents.id = kpi_associated_documents.document_id
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
		Where("kpi_associated_documents.kpi_register_id = ?", kpiID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(documents.title_en LIKE ? OR
			documents.title_ar LIKE ?)`,
			keyword, keyword)
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
	err := query.Offset(offset).Limit(params.PageSize).Find(&policyList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedDocumentListResponseDTO{
		Items: policyList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KPIAssociatedItemRepository) CreateKPIPolicies(policies []*models.KPIAssociatedDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&policies).Error
}

func (repo *KPIAssociatedItemRepository) DeleteKPIPolicyMapping(kpiID, documentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.
		Where("kpi_register_id = ? AND document_id = ?", kpiID, documentID).
		Delete(&models.KPIAssociatedDocument{}).
		Error
}

// func (repo *KPIAssociatedItemRepository) IsPolicyMappedToKPI(kpiID, documentID uint) (bool, error) {
// 	var count int64
// 	err := repo.db.
// 		Model(&models.KPIPolicy{}).
// 		Where("kpi_register_id = ? AND document_id = ?", kpiID, documentID).
// 		Count(&count).Error

// 	if err != nil {
// 		return false, err
// 	}
// 	return count > 0, nil
// }

func (repo *KPIAssociatedItemRepository) GetAllClauses(kpiID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemClauseListResponseDTO, error) {

	var clauseList []commonDTOs.AssociatedClauseRowDTO
	var totalCount int64
	query := repo.db.Table("kpi_compliances").
		Select(`
			compliance_clauses.title_en AS clause_title_en,
			compliance_clauses.title_ar AS clause_title_ar,
			compliance_clauses.clause_no,
			compliance_clauses.compliance_register_id,
			compliance_clauses.id AS compliance_clause_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar
		`).
		Joins(`
			LEFT JOIN compliance_clauses ON kpi_compliances.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_register ON kpi_compliances.compliance_register_id = compliance_register.id
		`).
		Where("kpi_compliances.kpi_register_id = ?", kpiID)
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`(
		compliance_register.title_en LIKE ? OR
		compliance_register.title_ar LIKE ? OR
		compliance_clauses.clause_no LIKE ? OR
		compliance_clauses.title_en LIKE ? OR
		compliance_clauses.title_ar LIKE ?)`, keyword, keyword, keyword, keyword, keyword)
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

func (repo *KPIAssociatedItemRepository) GetClausesByKPIID(kpiID uint, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clauseIDs []uint
	if err := db.
		Model(&models.KPICompliance{}).
		Where("kpi_register_id = ?", kpiID).
		Pluck("compliance_clause_id", &clauseIDs).Error; err != nil {
		return nil, err
	}

	return clauseIDs, nil
}
func (repo *KPIAssociatedItemRepository) CreateKPIClauses(clauses []*models.KPICompliance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&clauses).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIAssociatedItemRepository) DeleteKPIClauses(kpiID uint, clauseIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.
		Where("kpi_register_id = ? AND compliance_clause_id IN ?", kpiID, clauseIDs).
		Delete(&models.KPICompliance{}).Error; err != nil {
		return err
	}

	return nil
}
