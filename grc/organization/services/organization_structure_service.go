package services

import (
	"fmt"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingServices "jethur/grc/settings/services"
)

// Interface
type OrganizationStructureServiceInterface interface {
	GetOrganizationStructure(sort string, sortBy string) ([]organizationDTOs.OrganizationStructureDTO, error)
	GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error)
	GetOrganizationUnitDashboard(unitID uint) (*organizationDTOs.OrganizationUnitDashboardDTO, error)
}

type OrganizationStructureService struct {
	repo                 repositories.OrganizationStructureRepositoryInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
}

func NewOrganizationStructureService(repo repositories.OrganizationStructureRepositoryInterface, organizationLevelSvc settingServices.OrganizationLevelServiceInterface) *OrganizationStructureService {
	return &OrganizationStructureService{
		repo:                 repo,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (s *OrganizationStructureService) GetOrganizationStructure(sort string, sortBy string) ([]organizationDTOs.OrganizationStructureDTO, error) {
	// Get all organization units
	units, err := s.repo.GetOrganizationStructure(sort, sortBy)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of units by ID
	unitMap := make(map[uint]*organizationDTOs.OrganizationStructureDTO)
	var rootUnits []*organizationDTOs.OrganizationStructureDTO

	// First pass: create all nodes
	for _, unit := range units {
		structure := &organizationDTOs.OrganizationStructureDTO{
			ID:      unit.ID,
			TitleEn: unit.TitleEn,
			TitleAr: unit.TitleAr,
		}
		unitMap[unit.ID] = structure

		// If it's a root unit, add to root units
		if unit.OrganizationUnitID == nil {
			rootUnits = append(rootUnits, structure)
		}
	}

	// Second pass: build the tree
	for _, unit := range units {
		if unit.OrganizationUnitID != nil {
			if parent, exists := unitMap[*unit.OrganizationUnitID]; exists {
				parent.Children = append(parent.Children, *unitMap[unit.ID])
			}
		}
	}

	// Convert to slice of DTOs
	result := make([]organizationDTOs.OrganizationStructureDTO, len(rootUnits))
	for i, structure := range rootUnits {
		result[i] = *structure
	}

	return result, nil
}

func (s *OrganizationStructureService) GetUsersByOrganizationUnit(params *organizationDTOs.OrganizationUnitUserPrimaryDetailsQueryParamsDTO, orgUnitID uint) (*organizationDTOs.OrganizationUnitUserPrimaryDetailsResponseDTO, error) {
	var primaryOrgUnitField string

	// Get the primary organization level
	primaryLevel, err := s.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	// Determine the field name based on the primary level
	switch primaryLevel.Level {
	case 1:
		primaryOrgUnitField = "organization_unit_1_id"
	case 2:
		primaryOrgUnitField = "organization_unit_2_id"
	case 3:
		primaryOrgUnitField = "organization_unit_3_id"
	case 4:
		primaryOrgUnitField = "organization_unit_4_id"
	case 5:
		primaryOrgUnitField = "organization_unit_5_id"
	default:
		return nil, fmt.Errorf("unsupported primary level")
	}

	return s.repo.GetUsersByOrganizationUnit(params, orgUnitID, primaryOrgUnitField)
}

func (s *OrganizationStructureService) GetOrganizationUnitDashboard(unitID uint) (*organizationDTOs.OrganizationUnitDashboardDTO, error) {
	dto := &organizationDTOs.OrganizationUnitDashboardDTO{
		ResponsibleUnitID: unitID,
	}

	// Compliance %
	cs, err := s.repo.GetComplianceStatusByUnit(unitID)
	if err != nil {
		return nil, err
	}
	dto.ComplianceStatus = cs

	// Risks
	riskCount, err := s.repo.GetOpenRiskCountByUnit(unitID)
	if err != nil {
		return nil, err
	}
	dto.OpenRiskCount = riskCount

	// Control Efficiency (DTO)
	ce, err := s.repo.GetControlEfficiencyByUnit(unitID)
	if err != nil {
		return nil, err
	}
	dto.OverallControlEfficiency = ce

	return dto, nil
}
