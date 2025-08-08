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

func (pais *ThirdPartyAssociatedItemService) GetAllThirdPartyProcesses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := pais.ThirdPartySvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllThirdPartyProcesses(assetID, params)
}

func (pais *ThirdPartyAssociatedItemService) CreateThirdPartyProcess(ThirdPartyID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := pais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("ThirdParty does not exist for ID: %d", ThirdPartyID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateProcessMapping(ThirdPartyID, dto.ProcessIDs, "third_party_processes", "third_party_register_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
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

	if err := pais.repo.CreateThirdPartyProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ThirdPartyAssociatedItemService) DeleteThirdPartyProcess(ThirdPartyID uint, processID uint) error {
	// Validate control existence
	if _, err := pais.ThirdPartySvc.GetByID(ThirdPartyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("ThirdParty does not exist for ID: %d", ThirdPartyID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteThirdPartyProcess(ThirdPartyID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
