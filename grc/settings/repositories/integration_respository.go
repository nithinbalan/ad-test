package repositories

import (
	integrationDTOs "jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

// Interfaces
type IntegrationRepositoryInterface interface {
	GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationDTO, error)
}

// implementation
type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{
		db: db,
	}
}

func (repo *IntegrationRepository) GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationDTO, error) {
	var integrations []*integrationDTOs.IntegrationDTO

	query := repo.db.Table("integrations").
		Select(`
            integrations.*, integration_categories.title_en as integration_category_title_en, integration_categories.title_ar as integration_category_title_ar`).
		Joins(`
		    LEFT JOIN integration_categories ON integration_categories.id = integrations.integration_category_id
		`)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.Keyword != "" {
		query = query.Where(`
				(integrations.title_en LIKE ? OR
				integrations.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&integrations)

	return integrations, result.Error

}
