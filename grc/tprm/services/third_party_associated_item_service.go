package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/tprm/repositories"
	"jethur/utility"
	"time"
)

type ThirdPartyAssociatedItemServiceInterface interface {
	GetAllThirdPartyProcesses(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateThirdPartyProcess(ThirdPartyID uint, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteThirdPartyProcess(ThirdPartyID uint, processID uint) error

	GetAllThirdPartyStrategicObjectives(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateThirdPartyStrategicObjective(ThirdPartyID uint, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteThirdPartyStrategicObjective(ThirdPartyID uint, strategicObjectiveID uint) error

	GetAllThirdPartyAssets(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateThirdPartyAsset(ThirdPartyID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteThirdPartyAsset(ThirdPartyID uint, assetID uint) error

	GetAllThirdPartyIncidents(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateThirdPartyIncident(ThirdPartyID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteThirdPartyIncident(ThirdPartyID uint, incidentID uint) error
	GetAllThirdPartyDocuments(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateThirdPartyDocument(ThirdPartyID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteThirdPartyDocument(ThirdPartyID uint, incidentID uint) error
	GetAllThirdPartyControls(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateThirdPartyControl(ThirdPartyID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteThirdPartyControl(ThirdPartyID uint, incidentID uint) error
}

// Service and implementations
type ThirdPartyAssociatedItemService struct {
	repo          repositories.ThirdPartyAssociatedItemRepositoryInterface
	ThirdPartySvc ThirdPartyRegisterServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewThirdPartyAssociatedItemService(repo repositories.ThirdPartyAssociatedItemRepositoryInterface,
	ThirdPartySvc ThirdPartyRegisterServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) ThirdPartyAssociatedItemServiceInterface {
	return &ThirdPartyAssociatedItemService{
		repo:          repo,
		ThirdPartySvc: ThirdPartySvc,
		validationSvc: validationSvc,
	}
}

func (tpais *ThirdPartyAssociatedItemService) GetAllThirdPartyProcesses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := tpais.ThirdPartySvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return tpais.repo.GetAllThirdPartyProcesses(assetID, params)
}

func (tpais *ThirdPartyAssociatedItemService) CreateThirdPartyProcess(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("ThirdParty does not exist for ID: %d", ThirdPartyID), "not_exist", "asset_id")
	}

	if err := tpais.validationSvc.ValidateProcessMapping(ThirdPartyID, dto.ProcessIDs, "third_party_processes", "third_party_register_id"); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()
	var mappings []*models.ThirdPartyProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.ThirdPartyProcess{
			ProcessID:            processID,
			ThirdPartyRegisterID: ThirdPartyID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := tpais.repo.CreateThirdPartyProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) DeleteThirdPartyProcess(ThirdPartyID uint, processID uint) error {
	// Validate control existence
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("ThirdParty does not exist for ID: %d", ThirdPartyID), "not_exist", "asset_register_id")
	}

	tx := tpais.repo.BeginTransaction()

	if err := tpais.repo.DeleteThirdPartyProcess(ThirdPartyID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) GetAllThirdPartyStrategicObjectives(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return nil, err
	}
	return tpais.repo.GetAllThirdPartyStrategicObjectives(ThirdPartyID, params)
}

func (tpais *ThirdPartyAssociatedItemService) CreateThirdPartyStrategicObjective(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {

	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("ThirdParty does not exist for ID: %d", ThirdPartyID), "not_exist", "asset_id")
	}

	if err := tpais.validationSvc.ValidateStrtegicObjectiveMapping(ThirdPartyID, dto.StrategicObjectiveIDs, "third_party_strategic_objectives", "third_party_register_id"); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()
	var mappings []*models.ThirdPartyStrategicObjective
	currentTime := time.Now()

	for _, StrategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.ThirdPartyStrategicObjective{
			StrategicObjectiveID: StrategicObjectiveID,
			ThirdPartyRegisterID: ThirdPartyID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := tpais.repo.CreateThirdPartyStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) DeleteThirdPartyStrategicObjective(ThirdPartyID uint, strategicObjectiveID uint) error {
	// Validate control existence
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("ThirdParty does not exist for ID: %d", ThirdPartyID), "not_exist", "asset_register_id")
	}

	tx := tpais.repo.BeginTransaction()

	if err := tpais.repo.DeleteThirdPartyStrategicObjective(ThirdPartyID, strategicObjectiveID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) GetAllThirdPartyAssets(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return nil, err
	}
	return tpais.repo.GetAllThirdPartyAssets(ThirdPartyID, params)
}

func (tpais *ThirdPartyAssociatedItemService) CreateThirdPartyAsset(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	if err := tpais.validationSvc.ValidateAssetMapping(ThirdPartyID, dto.AssetIDs, "third_party_assets", "third_party_register_id"); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()
	var mappings []*models.ThirdPartyAsset
	currentTime := time.Now()

	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.ThirdPartyAsset{
			AssetRegisterID:      AssetID,
			ThirdPartyRegisterID: ThirdPartyID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := tpais.repo.CreateThirdPartyAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) DeleteThirdPartyAsset(ThirdPartyID uint, assetID uint) error {
	// Validate control existence
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()

	if err := tpais.repo.DeleteThirdPartyAsset(ThirdPartyID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) GetAllThirdPartyIncidents(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return nil, err
	}
	return tpais.repo.GetAllThirdPartyIncidents(ThirdPartyID, params)
}

func (tpais *ThirdPartyAssociatedItemService) CreateThirdPartyIncident(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	if err := tpais.validationSvc.ValidateIncidentMapping(ThirdPartyID, dto.IncidentIDs, "third_party_incidents", "third_party_register_id"); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()
	var mappings []*models.ThirdPartyIncident
	currentTime := time.Now()

	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.ThirdPartyIncident{
			IncidentRegisterID:   incidentID,
			ThirdPartyRegisterID: ThirdPartyID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := tpais.repo.CreateThirdPartyIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) DeleteThirdPartyIncident(ThirdPartyID uint, incidentID uint) error {
	// Validate control existence
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()

	if err := tpais.repo.DeleteThirdPartyIncident(ThirdPartyID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) GetAllThirdPartyDocuments(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return nil, err
	}
	return tpais.repo.GetAllThirdPartyDocuments(ThirdPartyID, params)
}

func (tpais *ThirdPartyAssociatedItemService) CreateThirdPartyDocument(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	if err := tpais.validationSvc.ValidateDocumentMapping(ThirdPartyID, dto.DocumentIDs, "third_party_associated_documents", "third_party_register_id"); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()
	var mappings []*models.ThirdPartyAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.ThirdPartyAssociatedDocument{
			DocumentID:           documentID,
			ThirdPartyRegisterID: ThirdPartyID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := tpais.repo.CreateThirdPartyDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) DeleteThirdPartyDocument(ThirdPartyID uint, documentID uint) error {
	// Validate control existence
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()

	if err := tpais.repo.DeleteThirdPartyDocument(ThirdPartyID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) GetAllThirdPartyControls(ThirdPartyID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return nil, err
	}
	return tpais.repo.GetAllThirdPartyControls(ThirdPartyID, params)
}

func (tpais *ThirdPartyAssociatedItemService) CreateThirdPartyControl(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	if err := tpais.validationSvc.ValidateControlMapping(ThirdPartyID, dto.ControlIDs, "third_party_controls", "third_party_register_id"); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()
	var mappings []*models.ThirdPartyControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.ThirdPartyControl{
			ControlID:            controlID,
			ThirdPartyRegisterID: ThirdPartyID,
			CreatedByID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := tpais.repo.CreateThirdPartyControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (tpais *ThirdPartyAssociatedItemService) DeleteThirdPartyControl(ThirdPartyID uint, controlID uint) error {
	// Validate control existence
	if _, err := tpais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return err
	}

	tx := tpais.repo.BeginTransaction()

	if err := tpais.repo.DeleteThirdPartyControl(ThirdPartyID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
