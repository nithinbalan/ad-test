package repositories

import (
	organizationDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type ProcessAssociatedItemRepositoryInterface interface {
	GetAllControls(params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (map[string]map[string][]*organizationDTOs.ProcessControlListDTO, error)
}

type ProcessAssociatedItemRepository struct {
	db *gorm.DB
}

func NewProcessAssociatedItemRepository(db *gorm.DB) *ProcessAssociatedItemRepository {
	return &ProcessAssociatedItemRepository{
		db: db,
	}
}

func (repo *ProcessAssociatedItemRepository) GetAllControls(params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (map[string]map[string][]*organizationDTOs.ProcessControlListDTO, error) {
	var controlList []*organizationDTOs.ProcessControlListDTO

	query := repo.db.Table("controls").
		Select(`
			controls.id AS control_id,
            control_categories.title_en AS control_category_title_en,
            control_categories.title_ar AS control_category_title_ar,
            control_types.title_en AS control_type_title_en,
            control_types.title_ar AS control_type_title_ar,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            control_efficiency_statuses.title_en AS control_efficiency_status_en,
            control_efficiency_statuses.title_ar AS control_efficiency_status_ar,
            control_implementation_statuses.title_en AS control_implementation_status_en,
            control_implementation_statuses.title_ar AS control_implementation_status_ar
        `).
		Joins(`
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
        `)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			control_categories.title_en LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_efficiency_statuses.title_en LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword)
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

	// Grouping controls by implementation status and then by efficiency status
	groupedControls := make(map[string]map[string][]*organizationDTOs.ProcessControlListDTO)

	for _, control := range controlList {
		implementationStatus := control.ControlImplementationStatusEn
		efficiencyStatus := control.ControlEfficiencyStatusEn

		if _, exists := groupedControls[implementationStatus]; !exists {
			groupedControls[implementationStatus] = make(map[string][]*organizationDTOs.ProcessControlListDTO)
		}
		if _, exists := groupedControls[implementationStatus][efficiencyStatus]; !exists {
			groupedControls[implementationStatus][efficiencyStatus] = []*organizationDTOs.ProcessControlListDTO{}
		}
		groupedControls[implementationStatus][efficiencyStatus] = append(groupedControls[implementationStatus][efficiencyStatus], control)
	}

	return groupedControls, nil
}
