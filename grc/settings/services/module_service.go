package services

import (
	"jethur/grc/settings/repositories"
)

type ModuleServiceInterface interface {
	GetModuleIDByIdentifier(identifier string) (uint, error)
}

// Service and implementations
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
