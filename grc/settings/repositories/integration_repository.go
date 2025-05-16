package repositories

import (
	"errors"
	"jethur/grc/models"
	integrationDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IntegrationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationListDTO, error)
	GetByID(id uint) (*models.Integration, error)
	GetFullDetailsByID(id uint) (*models.Integration, error)
	Update(integrationID uint, integration *models.Integration) (uint, error)
	CreateCredential(integrationCredential *models.IntegrationCredential) error
	UpdateCredential(integrationCredential *models.IntegrationCredential, tx ...utility.DatabaseTransaction) error
	DeleteCredential(id uint, tx ...utility.DatabaseTransaction) error

	GetByType(integrationType string) (*models.Integration, error)
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

// BeginTransaction starts a new transaction.
func (repo *IntegrationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IntegrationRepository) GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationListDTO, error) {
	var integrations []*integrationDTOs.IntegrationListDTO

	query := repo.db.Table("integrations").
		Select(`
            integrations.*, 
			integration_categories.title_en as integration_category_title_en,
			integration_categories.title_ar as integration_category_title_ar
			`).
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
	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("integrations.status = ?", params.Status)
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "is_enable":
			query = query.Where("integrations.is_enable = ?", true)
		}
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	result := query.Find(&integrations)

	return integrations, result.Error

}

func (repo *IntegrationRepository) GetByID(id uint) (*models.Integration, error) {
	var integration *models.Integration

	result := repo.db.First(&integration, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return integration, result.Error
}

func (repo *IntegrationRepository) GetFullDetailsByID(id uint) (*models.Integration, error) {
	var integration *models.Integration // Declare the variable inside the function

	result := repo.db.
		Preload("IntegrationCredential").
		Preload("CreatedBYUser").
		Preload("UpdatedBYUser").
		First(&integration, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return integration, nil
}

func (repo *IntegrationRepository) Update(integrationId uint, integration *models.Integration) (uint, error) {
	db := repo.db
	if err := db.Save(integration).Error; err != nil {
		return 0, err
	}

	return integrationId, nil
}

func (repo *IntegrationRepository) CreateCredential(integrationCredential *models.IntegrationCredential) error {

	result := repo.db.Create(&integrationCredential)
	return result.Error

}

func (repo *IntegrationRepository) UpdateCredential(integrationCredential *models.IntegrationCredential, tx ...utility.DatabaseTransaction) error {

	if err := repo.db.Model(&integrationCredential).Updates(integrationCredential).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IntegrationRepository) DeleteCredential(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.IntegrationCredential{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IntegrationRepository) GetByType(integrationType string) (*models.Integration, error) {
	var integration *models.Integration

	result := repo.db.Where("type = ?", integrationType).First(&integration)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Integration type with type "+integrationType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return integration, nil
}
