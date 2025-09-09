package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"time"
)

type KRIAssociatedItemServiceInterface interface {
	GetAllKRIDocuments(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateKRIDocument(kriID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteKRIDocument(kriID uint, documentID uint) error
	GetAllKRIControls(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateKRIControl(kriID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteKRIControl(kriID uint, controlID uint) error
	GetAllKRIAssets(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateKRIAsset(kriID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteKRIAsset(kriID uint, assetID uint) error
	GetAllKRIIncidents(kriID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateKRIIncident(kriID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteKRIIncident(kriID uint, incidentID uint) error
}

// Service and implementations
type KRIAssociatedItemService struct {
	repo          repositories.KRIAssociatedItemRepositoryInterface
	KRISvc        KRIRegisterServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewKRIAssociatedItemService(repo repositories.KRIAssociatedItemRepositoryInterface,
	KRISvc KRIRegisterServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) KRIAssociatedItemServiceInterface {
	return &KRIAssociatedItemService{
		repo:          repo,
		KRISvc:        KRISvc,
		validationSvc: validationSvc,
	}
}

func (pais *KRIAssociatedItemService) GetAllKRIDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := pais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllKRIDocuments(assetID, params)
}

func (pais *KRIAssociatedItemService) CreateKRIDocument(kriID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	if err := pais.validationSvc.ValidateDocumentMapping(kriID, dto.DocumentIDs, "kri_associated_documents", "kri_register_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.KRIAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.KRIAssociatedDocument{
			DocumentID:    documentID,
			KRIRegisterID: kriID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := pais.repo.CreateKRIDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *KRIAssociatedItemService) DeleteKRIDocument(kriID uint, documentID uint) error {
	// Validate control existence
	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteKRIDocument(kriID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *KRIAssociatedItemService) GetAllKRIControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := pais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllKRIControls(assetID, params)
}

func (pais *KRIAssociatedItemService) CreateKRIControl(kriID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	if err := pais.validationSvc.ValidateControlMapping(kriID, dto.ControlIDs, "kri_controls", "kri_register_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.KRIControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.KRIControl{
			ControlID:     controlID,
			KRIRegisterID: kriID,
			CreatedBYID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := pais.repo.CreateKRIControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *KRIAssociatedItemService) DeleteKRIControl(kriID uint, controlID uint) error {
	// Validate control existence
	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", kriID), "not_exist", "kri_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteKRIControl(kriID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *KRIAssociatedItemService) GetAllKRIAssets(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := pais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllKRIAssets(assetID, params)
}

func (pais *KRIAssociatedItemService) CreateKRIAsset(kriID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := pais.validationSvc.ValidateAssetMapping(kriID, dto.AssetIDs, "asset_kris", "kri_register_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.AssetKri
	currentTime := time.Now()

	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.AssetKri{
			AssetRegisterID: AssetID,
			KriRegisterID:   kriID,
			CreatedBYID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := pais.repo.CreateKRIAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *KRIAssociatedItemService) DeleteKRIAsset(kriID uint, assetID uint) error {

	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteKRIAsset(kriID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *KRIAssociatedItemService) GetAllKRIIncidents(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := pais.KRISvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllKRIIncidents(assetID, params)
}

func (pais *KRIAssociatedItemService) CreateKRIIncident(kriID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	if err := pais.validationSvc.ValidateIncidentMapping(kriID, dto.IncidentIDs, "kri_incidents", "kri_register_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.KriIncident
	currentTime := time.Now()

	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.KriIncident{
			IncidentRegisterID: incidentID,
			KriRegisterID:      kriID,
			CreatedBYID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := pais.repo.CreateKRIIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *KRIAssociatedItemService) DeleteKRIIncident(kriID uint, incidentID uint) error {

	if _, err := pais.KRISvc.GetByID(kriID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteKRIIncident(kriID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
