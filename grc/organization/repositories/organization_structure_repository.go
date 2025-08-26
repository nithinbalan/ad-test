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
	GetComplianceStatusByUnit(unitID uint) (float64, error)
	GetOpenRiskCountByUnit(unitID uint) (int64, error)
	GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error)
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

func (repo *OrganizationStructureRepository) GetComplianceStatusByUnit(unitID uint) (float64, error) {
	var compliance float64
	query := `
		SELECT 
			CASE 
				WHEN COUNT(*) = 0 THEN 0
				ELSE CAST(SUM(CASE WHEN carru.is_completed = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(*) AS FLOAT)
			END
		FROM compliance_assessment_requirement_responsible_units carru
		WHERE carru.responsible_unit_id = ?
	`
	if err := repo.db.Raw(query, unitID).Scan(&compliance).Error; err != nil {
		return 0, err
	}
	return compliance, nil
}

func (repo *OrganizationStructureRepository) GetOpenRiskCountByUnit(unitID uint) (int64, error) {
	var count int64
	query := `
		SELECT COUNT(*)
		FROM risk_register r
		JOIN risk_statuses rs ON r.risk_status_id = rs.id
		WHERE r.responsible_unit_id = ? AND rs.type = 'open'
	`
	if err := repo.db.Raw(query, unitID).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *OrganizationStructureRepository) GetControlEfficiencyByUnit(unitID uint) (*organizationDTOs.ControlEfficiencyStatusDTO, error) {
	var avgScore float64
	if err := repo.db.Model(&models.Control{}).
		Select("COALESCE(AVG(efficiency_score),0)").
		Where("responsible_unit_id = ?", unitID).
		Scan(&avgScore).Error; err != nil {
		return nil, err
	}

	var status organizationDTOs.ControlEfficiencyStatusDTO
	err := repo.db.Table("control_efficiency_statuses").
		Select("id, title_en, title_ar, status, color_code, score_from, score_to").
		Where("? BETWEEN score_from AND score_to AND status = 'active'", avgScore).
		First(&status).Error
	if err != nil {
		return nil, err
	}

	return &status, nil
}
