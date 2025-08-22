package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interface
type OrganizationStructureRepositoryInterface interface {
	GetOrganizationStructure(sort string, sortBy string) ([]models.OrganizationUnit, error)
}

type OrganizationStructureRepository struct {
	db *gorm.DB
}

func NewOrganizationStructureRepository(db *gorm.DB) *OrganizationStructureRepository {
	return &OrganizationStructureRepository{
		db: db,
	}
}

func (repo *OrganizationStructureRepository) GetOrganizationStructure(sort string, sortBy string) ([]models.OrganizationUnit, error) {
	var units []models.OrganizationUnit

	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Build the base query
	query := repo.db.Table("organization_units").
		Select("organization_units.*").
		Joins("LEFT JOIN organization_levels ON organization_units.organization_level_id = organization_levels.id")

	// Execute the query
	query = query.Order(sortBy + " " + sort)
	err := query.Find(&units).Error

	return units, err
}
