package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/settings/dtos"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ModuleServiceInterface interface {
	GetAll(params *dtos.ModuleListQueryParamsDTO) (*[]dtos.ModuleListDTO, error)
	GetModuleIDByIdentifier(identifier string) (uint, error)
	GetByID(id uint) (*models.Module, error)
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
