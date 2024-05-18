package services

import (
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationRepositories "jethur/grc/organization/repositories"
)

type ProcessServiceInterface interface {
	// GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	// GetByID(id uint) (*models.Process, error)
	// GetFullDetailsByID(id uint) (*organizationDTOs.ProcessDetailDTO, error)
	// Create(currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) (uint, error)
	// Update(id uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) error
	// Delete(ids []uint, currentUserID uint) ([]uint, error)

	// Archive(ids []uint, currentUserID uint) ([]uint, error)
	// Restore(ids []uint, currentUserID uint) ([]uint, error)
}

type ProcessService struct {
	repo               organizationRepositories.ControlRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	processGroupSvc    ProcessGroupServiceInterface
	userSvc            UserServiceInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
}

func NewProcessService(
	repo organizationRepositories.ControlRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processGroupSvc ProcessGroupServiceInterface,
	userSvc UserServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface) ProcessServiceInterface {
	return &ProcessService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		processGroupSvc:    processGroupSvc,
		userSvc:            userSvc,
		lifecycleStatusSvc: lifecycleStatusSvc,
	}
}
