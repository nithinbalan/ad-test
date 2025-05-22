package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"
)

type ComplianceAssociatedItemServiceInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error)
	GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(complianceID uint) (*complianceDTOs.ComplianceAssociatedItemsCountResponseDTO, error)
	GetComplianceRequirementCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error)
	GetAllStrategicObjectives(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateControlsToCompliance(complianceID, userID uint, dto *complianceDTOs.ComplianceRegisterControlMapDTO) error
	DeleteControlFromCompliance(complianceID, controlID uint) error
}

// Service and implementations
type ComplianceAssociatedItemService struct {
	repo                  complianceRepositories.ComplianceAssociatedItemRepositoryInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
	controlSvc            organizationServices.ControlServiceInterface
}

func NewComplianceAssociatedItemService(
	repo complianceRepositories.ComplianceAssociatedItemRepositoryInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
) ComplianceAssociatedItemServiceInterface {
	return &ComplianceAssociatedItemService{
		repo:                  repo,
		complianceRegisterSvc: complianceRegisterSvc,
		controlSvc:            controlSvc,
	}
}

func (cais *ComplianceAssociatedItemService) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllControls(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllRisks(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemRiskListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllRisks(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllIncidents(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllStrategicObjectives(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllStrategicObjectives(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllAssociatedCounts(complianceID uint) (*complianceDTOs.ComplianceAssociatedItemsCountResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	response, err := cais.repo.GetAllAssociatedCounts(complianceID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cais *ComplianceAssociatedItemService) GetComplianceRequirementCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	response, err := cais.repo.GetComplianceRequirementCountByStatus(complianceID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cais *ComplianceAssociatedItemService) CreateControlsToCompliance(complianceID, userID uint, dto *complianceDTOs.ComplianceRegisterControlMapDTO) error {
	if err := cais.validateControlIDs(complianceID, dto.ControlIDs); err != nil {
		return err
	}

	tx := cais.repo.BeginTransaction()
	var mappings []*models.ComplianceControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.ComplianceControl{
			ComplianceRegisterID: complianceID,
			ControlID:            controlID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := cais.repo.CreateComplianceControls(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (cais *ComplianceAssociatedItemService) DeleteControlFromCompliance(complianceID, controlID uint) error {
	// Validate control existence
	if _, err := cais.controlSvc.GetByID(controlID); err != nil {
		return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlID), "not_exist", "control_id")
	}

	tx := cais.repo.BeginTransaction()

	if err := cais.repo.DeleteComplianceControlMapping(complianceID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (cais *ComplianceAssociatedItemService) validateControlIDs(complianceID uint, controlIDs []uint) error {
	if len(controlIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no controls provided"), "required", "control_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range controlIDs {
		alreadyMapped, err := cais.repo.IsControlMappedToCompliance(complianceID, id)
		if err != nil {
			return fmt.Errorf("failed to check existing mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(fmt.Errorf("control ID %d already mapped to this compliance", id), "duplicate", "control_ids")
		}
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", id), "duplicate", "control_ids")
		}
		if _, err := cais.controlSvc.GetByID(id); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", id), "not_exist", "control_ids")
		}
		seen[id] = true
	}
	return nil
}
