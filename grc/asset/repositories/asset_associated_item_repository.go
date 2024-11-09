package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"

	"gorm.io/gorm"
)

type AssetAssociatedItemRepositoryInterface interface {
	GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemControlListResponseDTO, error)
	GetAllIncidents(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemIncidentListResponseDTO, error)
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

func (repo *AssetAssociatedItemRepository) GetAllControls(assetID uint, params *assetDTOs.AssetAssociatedItemListQueryParamsDTO) (*assetDTOs.AssetAssociatedItemControlListResponseDTO, error) {
	var controlList []*assetDTOs.AssetControlListDTO
	var totalCount int64

	query := repo.db.Table("asset_controls").
		Select(`
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

	response := &assetDTOs.AssetAssociatedItemControlListResponseDTO{
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

	// Second Query: Fetch asset requirement incidents using the asset requirement IDs
	query := repo.db.Table("incident_assets").
		Select(`
            DISTINCT incident_assets.incident_register_id AS incident_register_id,
            incident_register.reference_number AS incident_reference_number,
            incident_register.title AS incident_register_title,
            incident_register.description AS incident_register_description,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_sub_categories.title_en AS incident_sub_category_title_en,
            incident_sub_categories.title_ar AS incident_sub_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
			severity_levels.title_en AS incident_severity_level_title_en,
			severity_levels.title_ar AS incident_severity_level_title_ar,
			severity_levels.color_code AS incident_severity_level_color_code,
			incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
            incident_statuses.color_code AS incident_status_color_code,
			reported_user.id AS incident_reported_by_id,
            reported_user.first_name_en AS incident_reported_user_first_name_en,
            reported_user.first_name_ar AS incident_reported_user_first_name_ar,
            reported_user.last_name_en AS incident_reported_user_last_name_en,
            reported_user.last_name_ar AS incident_reported_user_last_name_ar,
			reported_user.initial AS reported_user_initial,
        	reported_user.color AS reported_user_color,
        	reported_user.email AS reported_user_email,
        	reported_user.profile_image_path AS reported_user_profile_image_path
        `).
		Joins(`
            LEFT JOIN incident_register ON incident_register.id = incident_assets.incident_register_id
			LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
            LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
            LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
            LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
		    LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
            LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
        `).
		Where("incident_assets.asset_register_id = ?", assetID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("incident_register.title LIKE ?", keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "incident_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

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
