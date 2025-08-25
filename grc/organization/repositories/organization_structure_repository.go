package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

// Interface
type OrganizationStructureRepositoryInterface interface {
	GetOrganizationStructure(sort string, sortBy string) ([]models.OrganizationUnit, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, primaryOrgUnitField string) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
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

func (repo *OrganizationStructureRepository) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint, primaryOrgUnitField string) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	var users []organizationDTOs.OrganizationUnitUserPrimaryDetailsDTO
	var totalCount int64

	// Base query with joins
	joinClause := fmt.Sprintf("LEFT JOIN organization_units AS primary_unit ON primary_unit.id = users.%s", primaryOrgUnitField)
	query := repo.db.Table("users").
		Select(`
            users.id,
            users.first_name_en,
            users.last_name_en,
            users.first_name_ar,
            users.last_name_ar,
            designations.title_en AS designation_en,
            designations.title_ar AS designation_ar`).
		Joins("LEFT JOIN designations ON designations.id = users.designation_id").
		Joins(joinClause).
		Where("users."+primaryOrgUnitField+" = ?", orgUnitID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(int(offset)).Limit(int(params.PageSize))

	// Apply sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	response := &organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO{
		Items: users,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}
