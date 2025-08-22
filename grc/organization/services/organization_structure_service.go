package services

import (
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
)

// Interface
type OrganizationStructureServiceInterface interface {
	GetOrganizationStructure(sort string, sortBy string) ([]dtos.OrganizationStructureDTO, error)
}

type OrganizationStructureService struct {
	repo repositories.OrganizationStructureRepositoryInterface
}

func NewOrganizationStructureService(repo repositories.OrganizationStructureRepositoryInterface) *OrganizationStructureService {
	return &OrganizationStructureService{
		repo: repo,
	}
}

func (s *OrganizationStructureService) GetOrganizationStructure(sort string, sortBy string) ([]dtos.OrganizationStructureDTO, error) {
	// Get all organization units
	units, err := s.repo.GetOrganizationStructure(sort, sortBy)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup of units by ID
	unitMap := make(map[uint]*dtos.OrganizationStructureDTO)
	var rootUnits []*dtos.OrganizationStructureDTO

	// First pass: create all nodes
	for _, unit := range units {
		structure := &dtos.OrganizationStructureDTO{
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
	result := make([]dtos.OrganizationStructureDTO, len(rootUnits))
	for i, structure := range rootUnits {
		result[i] = *structure
	}

	return result, nil
}
