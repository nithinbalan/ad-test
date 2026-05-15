package services

import (
	"fmt"
	"time"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/utility"
)

type SystemRegisterAssociatedItemServiceInterface interface {
	// SLA methods
	GetAllSLAs(systemRegisterID uint, params *assetDTOs.SystemRegisterAssociatedItemListQueryParamsDTO) (*assetDTOs.SystemRegisterSLAListResponseDTO, error)
	CreateSLAs(systemRegisterID uint, currentUserID uint, formDTO *assetDTOs.SystemRegisterSLAFormDTO) error
	DeleteSLA(systemRegisterID uint, slaID uint) error

	// Process methods
	GetAllProcesses(systemRegisterID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateProcesses(systemRegisterID uint, currentUserID uint, formDTO *assetDTOs.SystemRegisterProcessFormDTO) error
	DeleteProcess(systemRegisterID uint, processID uint) error
}

type SystemRegisterAssociatedItemService struct {
	repo              repositories.SystemRegisterAssociatedItemRepositoryInterface
	systemRegisterSvc SystemRegisterServiceInterface
	validationSvc     commonServices.AssociatedItemValidationServiceInterface
}

func NewSystemRegisterAssociatedItemService(
	repo repositories.SystemRegisterAssociatedItemRepositoryInterface,
	systemRegisterSvc SystemRegisterServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) SystemRegisterAssociatedItemServiceInterface {
	return &SystemRegisterAssociatedItemService{
		repo:              repo,
		systemRegisterSvc: systemRegisterSvc,
		validationSvc:     validationSvc,
	}
}

// SLA Service Methods

func (s *SystemRegisterAssociatedItemService) GetAllSLAs(systemRegisterID uint, params *assetDTOs.SystemRegisterAssociatedItemListQueryParamsDTO) (*assetDTOs.SystemRegisterSLAListResponseDTO, error) {
	return s.repo.GetAllSLAs(systemRegisterID, params)
}

func (s *SystemRegisterAssociatedItemService) CreateSLAs(systemRegisterID uint, currentUserID uint, dto *assetDTOs.SystemRegisterSLAFormDTO) error {

	// Validate system register existence
	if _, err := s.systemRegisterSvc.GetByID(systemRegisterID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("system register does not exist for ID: %d", systemRegisterID),
			"not_exist",
			"system_register_id",
		)
	}

	// Validate SLA mapping (duplicate + already mapped check)
	if err := s.validationSvc.ValidateSLAMapping(
		systemRegisterID,
		dto.SLARegisterIDs,
		"system_register_slas",
		"system_register_id",
	); err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	var mappings []*models.SystemRegisterSLA
	currentTime := time.Now()

	for _, slaRegisterID := range dto.SLARegisterIDs {
		mappings = append(mappings, &models.SystemRegisterSLA{
			SystemRegisterID: systemRegisterID,
			SLARegisterID:    slaRegisterID,
			CreatedBYID:      currentUserID,
			CreatedAt:        currentTime,
		})
	}

	if err := s.repo.CreateSLAs(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SystemRegisterAssociatedItemService) DeleteSLA(systemRegisterID uint, slaID uint) error {

	// Validate system register existence
	if _, err := s.systemRegisterSvc.GetByID(systemRegisterID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("system register does not exist for ID: %d", systemRegisterID),
			"not_exist",
			"system_register_id",
		)
	}

	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteSystemRegisterSLA(systemRegisterID, slaID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Process Service Methods

func (s *SystemRegisterAssociatedItemService) GetAllProcesses(systemRegisterID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	return s.repo.GetAllProcesses(systemRegisterID, params)
}

func (s *SystemRegisterAssociatedItemService) CreateProcesses(
	systemRegisterID uint,
	currentUserID uint,
	dto *assetDTOs.SystemRegisterProcessFormDTO,
) error {

	// Validate system register existence
	if _, err := s.systemRegisterSvc.GetByID(systemRegisterID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("system register does not exist for ID: %d", systemRegisterID),
			"not_exist",
			"system_register_id",
		)
	}

	// Validate process mapping
	if err := s.validationSvc.ValidateProcessMapping(
		systemRegisterID,
		dto.ProcessIDs,
		"system_register_processes",
		"system_register_id",
	); err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	var mappings []*models.SystemRegisterProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.SystemRegisterProcess{
			SystemRegisterID: systemRegisterID,
			ProcessID:        processID,
			CreatedBYID:      currentUserID,
			CreatedAt:        currentTime,
		})
	}

	if err := s.repo.CreateProcesses(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *SystemRegisterAssociatedItemService) DeleteProcess(systemRegisterID uint, processID uint) error {

	// Validate system register existence
	if _, err := s.systemRegisterSvc.GetByID(systemRegisterID); err != nil {
		return utility.NewValidationError(
			fmt.Errorf("system register does not exist for ID: %d", systemRegisterID),
			"not_exist",
			"system_register_id",
		)
	}

	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteSystemRegisterProcess(systemRegisterID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
