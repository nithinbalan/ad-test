package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ModuleServiceInterface interface {
	GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error)
	GetModuleIDByIdentifier(identifier string) (uint, error)
	GetByID(id uint) (*models.Module, error)
	GetFullDetailsByID(id uint) (*dtos.ModuleDetailDTO, error)
	EnableModule(id uint) (*models.Module, error)
	DisableModule(id uint) (*models.Module, error)
}

type ModuleService struct {
	repo repositories.ModuleRepositoryInterface
}

func NewModuleService(repo repositories.ModuleRepositoryInterface) ModuleServiceInterface {
	return &ModuleService{
		repo: repo,
	}
}

func (ms *ModuleService) GetModuleIDByIdentifier(identifier string) (uint, error) {
	return ms.repo.GetModuleIDByIdentifier(identifier)
}

func (ms *ModuleService) GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error) {
	return ms.repo.GetAll(params)
}

func (ms *ModuleService) GetByID(id uint) (*models.Module, error) {
	module, err := ms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return module, nil
}

func (ms ModuleService) GetFullDetailsByID(moduleId uint) (*dtos.ModuleDetailDTO, error) {
	// Initialize your DTO
	var moduleDTO dtos.ModuleDetailDTO

	integration, err := ms.repo.GetFullDetailsByID(moduleId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, moduleId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(integration, &moduleDTO)
	if err != nil {
		return nil, err
	}

	return &moduleDTO, nil
}

func (ms *ModuleService) EnableModule(id uint) (*models.Module, error) {
	module, err := ms.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	module.IsEnable = true
	_, err = ms.repo.Update(module.ID, module)
	if err != nil {
		return nil, err
	}
	return module, nil
}

func (ms *ModuleService) DisableModule(id uint) (*models.Module, error) {
	module, err := ms.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	module.IsEnable = false
	_, err = ms.repo.Update(module.ID, module)
	if err != nil {
		return nil, err
	}
	return module, nil
}
