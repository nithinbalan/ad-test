package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type ComplianceLibraryServiceInterface interface {
	GetComplianceLibraryDetails(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
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

func (cls *ComplianceLibraryService) GetComplianceLibraryDetails(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
	return cls.repo.GetComplianceLibraryDetails(libraryID)
}
