package services

import (
	"fmt"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"
)

type AuditFindingAssociatedItemServiceInterface interface {
	GetAllAuditFindingDocuments(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateAuditFindingDocument(findingID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteAuditFindingDocument(findingID uint, documentID uint) error
	GetAllAuditFindingControls(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateAuditFindingControl(findingID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteAuditFindingControl(findingID uint, controlID uint) error

	GetAllAuditFindingProcesses(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateAuditFindingProcess(findingID uint, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteAuditFindingProcess(findingID uint, processID uint) error

	GetAllAuditFindingStrategicObjectives(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateAuditFindingStrategicObjective(findingID uint, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteAuditFindingStrategicObjective(findingID uint, strategicObjectiveID uint) error

	GetAllAuditFindingAssets(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateAuditFindingAsset(findingID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteAuditFindingAsset(findingID uint, assetID uint) error

	GetAllAuditFindingIncidents(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateAuditFindingIncident(findingID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteAuditFindingIncident(findingID uint, incidentID uint) error

	GetAllAuditFindingRisks(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateAuditFindingRisk(findingID uint, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error
	DeleteAuditFindingRisk(findingID uint, riskID uint) error
}

// Service and implementations
type AuditFindingAssociatedItemService struct {
	repo            repositories.AuditFindingAssociatedItemRepositoryInterface
	AuditFindingSvc AuditFindingServiceInterface
	validationSvc   commonServices.AssociatedItemValidationServiceInterface
}

func NewAuditFindingAssociatedItemService(repo repositories.AuditFindingAssociatedItemRepositoryInterface,
	AuditFindingSvc AuditFindingServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) AuditFindingAssociatedItemServiceInterface {
	return &AuditFindingAssociatedItemService{
		repo:            repo,
		AuditFindingSvc: AuditFindingSvc,
		validationSvc:   validationSvc,
	}
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingDocuments(assetID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingDocument(findingID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	if err := afais.validationSvc.ValidateDocumentMapping(findingID, dto.DocumentIDs, "audit_finding_associated_documents", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.AuditFindingAssociatedDocument{
			DocumentID:     documentID,
			AuditFindingID: findingID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingDocument(findingID uint, documentID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingDocument(findingID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingControls(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingControls(findingID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingControl(findingID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	if err := afais.validationSvc.ValidateControlMapping(findingID, dto.ControlIDs, "audit_finding_controls", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.AuditFindingControl{
			ControlID:      controlID,
			AuditFindingID: findingID,
			CreatedBYID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingControl(findingID uint, controlID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingControl(findingID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingProcesses(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingProcesses(findingID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingProcess(findingID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	if err := afais.validationSvc.ValidateProcessMapping(findingID, dto.ProcessIDs, "audit_finding_processes", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.AuditFindingProcess{
			ProcessID:      processID,
			AuditFindingID: findingID,
			CreatedByID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingProcess(findingID uint, processID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingProcess(findingID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingStrategicObjectives(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingStrategicObjectives(findingID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingStrategicObjective(findingID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	if err := afais.validationSvc.ValidateStrtegicObjectiveMapping(findingID, dto.StrategicObjectiveIDs, "audit_finding_strategic_objectives", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingStrategicObjective
	currentTime := time.Now()

	for _, StrategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.AuditFindingStrategicObjective{
			StrategicObjectiveID: StrategicObjectiveID,
			AuditFindingID:       findingID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingStrategicObjective(findingID uint, strategicObjectiveID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("AuditFinding does not exist for ID: %d", findingID), "not_exist", "audit_finding_id")
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingStrategicObjective(findingID, strategicObjectiveID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingAssets(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingAssets(findingID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingAsset(findingID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	if err := afais.validationSvc.ValidateAssetMapping(findingID, dto.AssetIDs, "audit_finding_assets", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingAsset
	currentTime := time.Now()

	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.AuditFindingAsset{
			AssetRegisterID: AssetID,
			AuditFindingID:  findingID,
			CreatedBYID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingAsset(findingID uint, assetID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingAsset(findingID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingIncidents(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingIncidents(findingID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingIncident(findingID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	if err := afais.validationSvc.ValidateIncidentMapping(findingID, dto.IncidentIDs, "audit_finding_incidents", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingIncident
	currentTime := time.Now()

	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.AuditFindingIncident{
			IncidentRegisterID: incidentID,
			AuditFindingID:     findingID,
			CreatedBYID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingIncident(findingID uint, incidentID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingIncident(findingID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) GetAllAuditFindingRisks(findingID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return afais.repo.GetAllAuditFindingRisks(findingID, params)
}

func (afais *AuditFindingAssociatedItemService) CreateAuditFindingRisk(findingID, userID uint, dto *commonDTOs.AssociatedRiskMapDTO) error {

	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	if err := afais.validationSvc.ValidateRiskMapping(findingID, dto.RiskIDs, "audit_finding_Risks", "audit_finding_id"); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()
	var mappings []*models.AuditFindingRisk
	currentTime := time.Now()

	for _, riskID := range dto.RiskIDs {
		mappings = append(mappings, &models.AuditFindingRisk{
			RiskRegisterID: riskID,
			AuditFindingID: findingID,
			CreatedBYID:    userID,
			CreatedAt:      currentTime,
		})
	}

	if err := afais.repo.CreateAuditFindingRisk(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (afais *AuditFindingAssociatedItemService) DeleteAuditFindingRisk(findingID uint, riskID uint) error {
	// Validate control existence
	if _, err := afais.AuditFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	tx := afais.repo.BeginTransaction()

	if err := afais.repo.DeleteAuditFindingRisk(findingID, riskID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
