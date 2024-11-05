package repositories

import (
	"jethur/grc/models"
	integrationDTOs "jethur/grc/settings/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IntegrationRepositoryInterface interface {
	GetAll(params *integrationDTOs.IntegrationListQueryParamsDTO) ([]*integrationDTOs.IntegrationDTO, error)
	GetFullDetailsByID(id uint) (*models.Integration, error)
	GetCredentailByID(id uint) (*models.IntegrationCredential, error)
	GetByID(id uint) (*models.Integration, error)
	Create(integrationCredential *models.IntegrationCredential) (uint, error)
	Update(id uint, integrationCredential *models.IntegrationCredential) error
	Delete(id uint, ids []*uint, tx ...utility.DatabaseTransaction) error
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

	for _, integration := range integrations {
		var credentials []*integrationDTOs.IntegrationCredentialDTO

		repo.db.Model(&models.IntegrationCredential{}).Where("integration_id = ?", integration.ID).Find(&credentials)
		integration.Credentials = credentials
	}

	return integrations, result.Error

}

func (repo *IntegrationRepository) GetFullDetailsByID(id uint) (*models.Integration, error) {
	var integration *models.Integration // Declare the variable inside the function

	result := repo.db.
		Preload("Credentials").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&integration, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return integration, nil
}

func (repo *IntegrationRepository) GetByID(id uint) (*models.Integration, error) {
	var integration *models.Integration

	result := repo.db.First(&integration, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return integration, result.Error
}

func (repo *IntegrationRepository) GetCredentailByID(id uint) (*models.IntegrationCredential, error) {
	var integrationCredential *models.IntegrationCredential

	result := repo.db.First(&integrationCredential, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return integrationCredential, result.Error
}

func (repo *IntegrationRepository) Create(integrationCredential *models.IntegrationCredential) (uint, error) {
	db := repo.db
	result := repo.db.Create(&integrationCredential)

	err := db.Model(&models.Integration{}).
		Where("id = ?", integrationCredential.IntegrationID).
		UpdateColumn("is_enable", true).Error
	if err != nil {
		return 0, err
	}

	return integrationCredential.IntegrationID, result.Error
}

func (repo *IntegrationRepository) Update(id uint, updatedIntegrationCredential *models.IntegrationCredential) error {
	var integrationCredential *models.IntegrationCredential

	if err := repo.db.First(&integrationCredential, id).Error; err != nil {
		return err
	}

	if err := repo.db.Model(&integrationCredential).Updates(updatedIntegrationCredential).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IntegrationRepository) Delete(id uint, ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.Model(&models.Integration{}).
		Where("id = ?", id).
		UpdateColumn("is_enable", false).Error
	if err != nil {
		return err
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.IntegrationCredential{}).Error; err != nil {
		return err
	}

	return nil
}
