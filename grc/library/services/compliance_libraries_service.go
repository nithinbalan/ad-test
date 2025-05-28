package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type ComplianceLibraryServiceInterface interface {
	GetAll() ([]*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
}

// Service and implementations
type ComplianceLibraryService struct {
	repo repositories.ComplianceLibraryRepositoryInterface
}

func NewComplianceLibraryService(
	repo repositories.ComplianceLibraryRepositoryInterface,
) ComplianceLibraryServiceInterface {
	return &ComplianceLibraryService{
		repo: repo,
	}
}

func (cls *ComplianceLibraryService) GetAll() ([]*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
	return cls.repo.GetAll()
}
