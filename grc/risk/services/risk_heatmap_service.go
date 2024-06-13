package services

import (
	"fmt"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
)

type RiskHeatmapServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskHeatmap, error)
	GetByID(id uint) (*models.RiskHeatmap, error)
	GetByType(status string) (*models.RiskHeatmap, error)
}

// Service and implementations
type RiskHeatmapService struct {
	repo repositories.RiskHeatmapRepositoryInterface
}

func NewRiskHeatmapService(repo repositories.RiskHeatmapRepositoryInterface) RiskHeatmapServiceInterface {
	return &RiskHeatmapService{
		repo: repo,
	}
}

func (ps *RiskHeatmapService) validateProcessFormDTO(processFormDTO *organizationDTOs.ProcessFormDTO) error {
	if _, err := ps.processGroupSvc.GetByID(uint(processFormDTO.ProcessGroupID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("process group does not exist"), "not_exist", "process group")
	}

	if _, err := ps.userSvc.GetByID(uint(processFormDTO.OwnerID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("process owner does not exist"), "not_exist", "process owner")
	}

	controls := processFormDTO.ControlIds

	controlIDs := make(map[uint]bool)

	for _, control := range controls {
		if controlIDs[*control] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", control), "duplicate", "control")
		}

		if _, err := ps.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		controlIDs[*control] = true
	}

	return nil
}
