package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type ComplianceLibraryServiceInterface interface {
	GetComplianceLibraryDetailsByID(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	GetClausesByVersionID(versionID uint) ([]libraryDTOs.GroupedClauseDTO, error)
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

func (cls *ComplianceLibraryService) GetComplianceLibraryDetailsByID(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
	return cls.repo.GetComplianceLibraryDetailsByID(libraryID)
}

func (cls *ComplianceLibraryService) GetClausesByVersionID(versionID uint) ([]libraryDTOs.GroupedClauseDTO, error) {
	return cls.repo.GetClausesByVersionID(versionID)
}
