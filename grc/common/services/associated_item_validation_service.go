package services

import (
	"fmt"
	"jethur/grc/common/repositories"
	"jethur/utility"
)

// Interface definition
type AssociatedItemValidationServiceInterface interface {
	ValidateRiskMapping(masterID uint, riskIDs []uint, masterName string, columnName string) error
}

// Service struct
type AssociatedItemValidationService struct {
	repo repositories.AssociatedItemValidationRepositoryInterface
}

// Constructor
func NewAssociatedItemValidationService(repo repositories.AssociatedItemValidationRepositoryInterface) AssociatedItemValidationServiceInterface {
	return &AssociatedItemValidationService{
		repo: repo,
	}
}

// Main validation logic
func (rms *AssociatedItemValidationService) ValidateRiskMapping(masterID uint, riskIDs []uint, masterName string, columnName string) error {
	if len(riskIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no risks provided"), "required", "risk_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range riskIDs {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate risk ID: %d", id), "duplicate", "risk_ids")
		}

		_, err := rms.repo.RiskExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("risk does not exist: %d", id), "not_exist", "risk_ids")
		}

		alreadyMapped, err := rms.repo.IsRiskMapped(masterName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("risk ID %d already mapped to this %s", id, masterName),
				"duplicate", "risk_ids")
		}

		seen[id] = true
	}

	return nil
}
