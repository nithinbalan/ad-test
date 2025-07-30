package services

import (
	"fmt"
	"jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"
)

type BcpAssociatedItemServiceInterface interface {
	GetAllBcpDocuments(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateBcpDocument(BcpID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteBcpDocument(BcpID uint, documentID uint) error
	GetAllBcpControl(BcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateBcpControl(BcpID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteBcpControl(BcpID uint, controlID uint) error
}

// Service and implementations
type BcpAssociatedItemService struct {
	repo          repositories.BcpAssociatedItemRepositoryInterface
	BcpSvc        BusinessContinuityPlanServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewBcpAssociatedItemService(repo repositories.BcpAssociatedItemRepositoryInterface,
	BcpSvc BusinessContinuityPlanServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) BcpAssociatedItemServiceInterface {
	return &BcpAssociatedItemService{
		repo:          repo,
		BcpSvc:        BcpSvc,
		validationSvc: validationSvc,
	}
}

func (pais *BcpAssociatedItemService) GetAllBcpDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := pais.BcpSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllBcpDocuments(assetID, params)
}

func (pais *BcpAssociatedItemService) CreateBcpDocument(BcpID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := pais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	if err := pais.validationSvc.ValidateDocumentMapping(BcpID, dto.DocumentIDs, "bcp_associated_documents", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.BCPAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.BCPAssociatedDocument{
			DocumentID:               documentID,
			BusinessContinuityPlanID: BcpID,
			CreatedByID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := pais.repo.CreateBcpDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *BcpAssociatedItemService) DeleteBcpDocument(BcpID uint, documentID uint) error {
	// Validate control existence
	if _, err := pais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteBcpDocument(BcpID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *BcpAssociatedItemService) GetAllBcpControl(bcpID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := pais.BcpSvc.GetByID(bcpID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllBcpControl(bcpID, params)
}

func (pais *BcpAssociatedItemService) CreateBcpControl(BcpID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := pais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	if err := pais.validationSvc.ValidateControlMapping(BcpID, dto.ControlIDs, "bcp_controls", "business_continuity_plan_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.BCPControl
	currentTime := time.Now()

	for _, ControlID := range dto.ControlIDs {
		mappings = append(mappings, &models.BCPControl{
			ControlID:                ControlID,
			BusinessContinuityPlanID: BcpID,
			CreatedBYID:              userID,
			CreatedAt:                currentTime,
		})
	}

	if err := pais.repo.CreateBcpControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *BcpAssociatedItemService) DeleteBcpControl(BcpID uint, controlID uint) error {
	// Validate control existence
	if _, err := pais.BcpSvc.GetByID(BcpID); err != nil {
		return utility.NewValidationError(fmt.Errorf("bcp does not exist for ID: %d", BcpID), "not_exist", "business_continuity_plan_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteBcpControl(BcpID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
