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
