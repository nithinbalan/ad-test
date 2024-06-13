package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"

	"gorm.io/gorm"
)

type ComplianceAssociatedItemRepositoryInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceControlListDTO, error)
}

type ComplianceAssociatedItemRepository struct {
	db *gorm.DB
}

func NewComplianceAssociatedItemRepository(db *gorm.DB) *ComplianceAssociatedItemRepository {
	return &ComplianceAssociatedItemRepository{
		db: db,
	}
}

func (repo *ComplianceAssociatedItemRepository) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) ([]*complianceDTOs.ComplianceControlListDTO, error) {
	var controlList []*complianceDTOs.ComplianceControlListDTO

	query := repo.db.Table("compliance_controls").
		Select(`
			controls.id AS control_id,
			controls.reference_number AS control_reference_number,
            controls.code AS control_code,
            controls.last_test_date AS control_last_test_date,
            control_categories.title_en AS control_category_title_en,
            control_categories.title_ar AS control_category_title_ar,
            control_types.title_en AS control_type_title_en,
            control_types.title_ar AS control_type_title_ar,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            control_efficiency_statuses.title_en AS control_efficiency_status_en,
            control_efficiency_statuses.title_ar AS control_efficiency_status_ar,
            control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
            control_efficiency_statuses.type AS control_efficiency_status_type,
            control_implementation_statuses.title_en AS control_implementation_status_en,
            control_implementation_statuses.title_ar AS control_implementation_status_ar
            control_implementation_statuses.type AS control_implementation_status_type
            control_implementation_statuses.color_code AS control_implementation_status_color_code,
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
			LEFT JOIN controls ON controls.id = compliance_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
			LEFT JOIN users AS owner ON controls.owner_id = owner.id
        `).Where("compliance_controls.compliance_register_id = ?", complianceID)

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
