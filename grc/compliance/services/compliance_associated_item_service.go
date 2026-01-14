package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"
)

type ComplianceAssociatedItemServiceInterface interface {
	GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error)
	GetAllRisks(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	GetAllIncidents(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemIncidentListResponseDTO, error)
	GetAllAssociatedCounts(complianceID uint) (*complianceDTOs.ComplianceAssociatedItemsCountResponseDTO, error)
	GetComplianceRequirementCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error)
	GetAllStrategicObjectives(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemStrategicObjectiveListResponseDTO, error)
	GetAllComplianceAssets(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	GetAllComplianceProcesses(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	GetAllComplianceKPI(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)

	CreateComplianceControls(complianceID, userID uint, dto *complianceDTOs.ComplianceRegisterControlMapDTO) error
	DeleteComplianceControl(complianceID, controlID uint) error
	CreateComplianceRisk(complianceID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteComplianceRisk(complianceID uint, riskID uint) error
	GetAllDocuments(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateComplianceDocument(complianceID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteComplianceDocument(complianceID uint, riskID uint) error
	GetComplianceControlCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error)
}

// Service and implementations
type ComplianceAssociatedItemService struct {
	repo                  complianceRepositories.ComplianceAssociatedItemRepositoryInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
	controlSvc            organizationServices.ControlServiceInterface
	validationSvc         commonServices.AssociatedItemValidationServiceInterface
}

func NewComplianceAssociatedItemService(
	repo complianceRepositories.ComplianceAssociatedItemRepositoryInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) ComplianceAssociatedItemServiceInterface {
	return &ComplianceAssociatedItemService{
		repo:                  repo,
		complianceRegisterSvc: complianceRegisterSvc,
		controlSvc:            controlSvc,
		validationSvc:         validationSvc,
	}
}

func (cais *ComplianceAssociatedItemService) GetAllControls(complianceID uint, params *complianceDTOs.ComplianceAssociatedItemListQueryParamsDTO) (*complianceDTOs.ComplianceAssociatedItemControlListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllControls(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) GetAllRisks(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
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

func (ccs *ComplianceAssociatedItemService) GetAllComplianceAssets(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	// Fetch associated controls through requirements from repository
	strategicObjectives, err := ccs.repo.GetAllComplianceAssets(complianceID, params)
	if err != nil {
		return nil, err
	}

	return strategicObjectives, nil
}

func (cais *ComplianceAssociatedItemService) GetAllComplianceProcesses(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	strategicObjectives, err := cais.repo.GetAllComplianceProcesses(complianceID, params)
	if err != nil {
		return nil, err
	}

	return strategicObjectives, nil
}

func (cais *ComplianceAssociatedItemService) GetAllComplianceKPI(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	kpis, err := cais.repo.GetAllComplianceKPI(complianceID, params)
	if err != nil {
		return nil, err
	}

	return kpis, nil
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

func (cais *ComplianceAssociatedItemService) CreateComplianceControls(complianceID, userID uint, dto *complianceDTOs.ComplianceRegisterControlMapDTO) error {
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

func (cais *ComplianceAssociatedItemService) DeleteComplianceControl(complianceID, controlID uint) error {
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

func (cais *ComplianceAssociatedItemService) CreateComplianceRisk(complianceID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance does not exist for ID: %d", complianceID), "not_exist", "compliance_register_id")
	}
	if err := cais.validationSvc.ValidateRiskMapping(complianceID, dto.RiskIDs, "risk_compliance_requirements", "compliance_register_id"); err != nil {
		return err
	}

	tx := cais.repo.BeginTransaction()
	var mappings []*models.RiskComplianceRequirement
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.RiskComplianceRequirement{
			RiskRegisterID:       riskID,
			ComplianceRegisterID: complianceID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := cais.repo.CreateComplianceRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (cais *ComplianceAssociatedItemService) DeleteComplianceRisk(complianceID uint, riskID uint) error {
	// Validate control existence
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance does not exist for ID: %d", complianceID), "not_exist", "compliance_register_id")
	}

	tx := cais.repo.BeginTransaction()

	if err := cais.repo.DeleteComplianceRisk(complianceID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (cais *ComplianceAssociatedItemService) GetAllDocuments(complianceID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return cais.repo.GetAllDocuments(complianceID, params)
}

func (cais *ComplianceAssociatedItemService) CreateComplianceDocument(complianceID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance does not exist for ID: %d", complianceID), "not_exist", "compliance_register_id")
	}
	if err := cais.validationSvc.ValidateDocumentMapping(complianceID, dto.DocumentIDs, "document_compliance_requirements", "compliance_register_id"); err != nil {
		return err
	}

	tx := cais.repo.BeginTransaction()
	var mappings []*models.DocumentComplianceRequirement
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.DocumentComplianceRequirement{
			DocumentID:           documentID,
			ComplianceRegisterID: complianceID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := cais.repo.CreateComplianceDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (cais *ComplianceAssociatedItemService) DeleteComplianceDocument(complianceID uint, documentID uint) error {
	// Validate control existence
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance does not exist for ID: %d", complianceID), "not_exist", "compliance_register_id")
	}

	tx := cais.repo.BeginTransaction()

	if err := cais.repo.DeleteComplianceDocument(complianceID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (cais *ComplianceAssociatedItemService) GetComplianceControlCountByStatus(complianceID uint) ([]complianceDTOs.ComplianceRequirementStatusCountResponseDTO, error) {
	if _, err := cais.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	response, err := cais.repo.GetComplianceControlCountByStatus(complianceID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
