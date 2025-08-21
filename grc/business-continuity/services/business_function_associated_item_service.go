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

type BusinessFunctionAssociatedItemServiceInterface interface {
	GetAllBusinessFunctionProcesses(BusinessFunctionID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateBusinessFunctionProcess(BusinessFunctionID uint, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteBusinessFunctionProcess(BusinessFunctionID uint, processID uint) error
}

// Service and implementations
type BusinessFunctionAssociatedItemService struct {
	repo                repositories.BusinessFunctionAssociatedItemRepositoryInterface
	BusinessFunctionSvc BusinessFunctionServiceInterface
	validationSvc       commonServices.AssociatedItemValidationServiceInterface
}

func NewBusinessFunctionAssociatedItemService(repo repositories.BusinessFunctionAssociatedItemRepositoryInterface,
	BusinessFunctionSvc BusinessFunctionServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) BusinessFunctionAssociatedItemServiceInterface {
	return &BusinessFunctionAssociatedItemService{
		repo:                repo,
		BusinessFunctionSvc: BusinessFunctionSvc,
		validationSvc:       validationSvc,
	}
}

func (pais *BusinessFunctionAssociatedItemService) GetAllBusinessFunctionProcesses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := pais.BusinessFunctionSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllBusinessFunctionProcesses(assetID, params)
}

func (pais *BusinessFunctionAssociatedItemService) CreateBusinessFunctionProcess(BusinessFunctionID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := pais.BusinessFunctionSvc.GetByID(BusinessFunctionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("BusinessFunction does not exist for ID: %d", BusinessFunctionID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateProcessMapping(BusinessFunctionID, dto.ProcessIDs, "business_function_processes", "business_function_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.BusinessFunctionProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.BusinessFunctionProcess{
			ProcessID:          processID,
			BusinessFunctionID: BusinessFunctionID,
			CreatedBYID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := pais.repo.CreateBusinessFunctionProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *BusinessFunctionAssociatedItemService) DeleteBusinessFunctionProcess(BusinessFunctionID uint, processID uint) error {
	// Validate control existence
	if _, err := pais.BusinessFunctionSvc.GetByID(BusinessFunctionID); err != nil {
		return utility.NewValidationError(fmt.Errorf("BusinessFunction does not exist for ID: %d", BusinessFunctionID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteBusinessFunctionProcess(BusinessFunctionID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
