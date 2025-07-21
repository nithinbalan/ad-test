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

	GetAllPolicies(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemPolicyListResponseDTO, error)
	CreateKPIPolicies(policies []*models.KPIPolicy, tx ...utility.DatabaseTransaction) error
	DeleteKPIPolicyMapping(kpiID, documentID uint, tx ...utility.DatabaseTransaction) error
	IsPolicyMappedToKPI(kpiID, documentID uint) (bool, error)
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
			risk_statuses.color_code AS risk_status_color_code,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
            risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
			review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
        	updated_user.color AS updated_user_color,
        	updated_user.email AS updated_user_email,
        	updated_user.profile_image_path AS updated_user_profile_image_path,
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
            residual_risk_rating.color_code AS residual_risk_rating_color_code,
			risk_register.created_at,
            risk_register.updated_at
			
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = kpi_risks.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
			LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
			LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
			LEFT JOIN users ON users.id = risk_register.responsible_user_id
			LEFT JOIN risk_matrix_likelihoods ON risk_matrix_likelihoods.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts ON risk_matrix_impacts.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = risk_register.current_risk_rating_value_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
			LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
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

func (repo *KPIAssociatedItemRepository) GetAllPolicies(kpiID uint, params *strategyDTOs.KPIAssociatedItemListQueryParamsDTO) (*strategyDTOs.KPIAssociatedItemPolicyListResponseDTO, error) {
	var policyList []*strategyDTOs.KPIPolicyListDTO
	var totalCount int64

	query := repo.db.Table("kpi_policies").
		Select(`
			kpi_policies.id,
			kpi_policies.document_id,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name
		`).
		Joins(`
			LEFT JOIN documents ON documents.id = kpi_policies.document_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
        `).Where("kpi_policies.kpi_register_id = ?", kpiID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

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

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&policyList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &strategyDTOs.KPIAssociatedItemPolicyListResponseDTO{
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

func (repo *KPIAssociatedItemRepository) CreateKPIPolicies(policies []*models.KPIPolicy, tx ...utility.DatabaseTransaction) error {
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
		Delete(&models.KPIPolicy{}).
		Error
}

func (repo *KPIAssociatedItemRepository) IsPolicyMappedToKPI(kpiID, documentID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.KPIPolicy{}).
		Where("kpi_register_id = ? AND document_id = ?", kpiID, documentID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
