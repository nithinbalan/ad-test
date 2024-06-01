package repositories

import (
	riskDTOs "jethur/grc/risk/dtos"

	"gorm.io/gorm"
)

type RiskAssociatedItemRepositoryInterface interface {
	GetAllControls(params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]*riskDTOs.RiskControlListDTO, error)
}

type RiskAssociatedItemRepository struct {
	db *gorm.DB
}

func NewRiskAssociatedItemRepository(db *gorm.DB) *RiskAssociatedItemRepository {
	return &RiskAssociatedItemRepository{
		db: db,
	}
}

func (repo *RiskAssociatedItemRepository) GetAllControls(params *riskDTOs.RiskAssociatedItemListQueryParamsDTO) ([]*riskDTOs.RiskControlListDTO, error) {
	var controlList []*riskDTOs.RiskControlListDTO

	query := repo.db.Table("controls").
		Select(`
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

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "open_high_risk":
			query = query.Where("controls.risk_status_title_en = ?", "open")
		}
	}

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

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "controls.title_en" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	if err := query.Find(&controlList).Error; err != nil {
		return nil, err
	}

	return controlList, nil

}
