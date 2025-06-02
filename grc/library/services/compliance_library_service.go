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
	clauses, err := cls.repo.GetAllClausesByVersionID(versionID)
	if err != nil {
		return nil, err
	}

	mainMap := make(map[uint]*libraryDTOs.GroupedClauseDTO)
	var result []*libraryDTOs.GroupedClauseDTO

	for _, clause := range clauses {
		if clause.ComplianceLibraryClauseID == nil {
			group := &libraryDTOs.GroupedClauseDTO{
				ID:            clause.ID,
				ClauseNumber:  clause.ClauseNumber,
				ClauseTitleEn: clause.ClauseTitleEn,
				ClauseTitleAr: *clause.ClauseTitleAr,
				CreatedAt:     clause.CreatedAt,
				UpdatedAt:     *clause.UpdatedAt,
				SubClauses:    []libraryDTOs.SubClauseDTO{},
			}
			mainMap[clause.ID] = group
			result = append(result, group)
		} else {
			parent := mainMap[*clause.ComplianceLibraryClauseID]
			if parent != nil {
				parent.SubClauses = append(parent.SubClauses, libraryDTOs.SubClauseDTO{
					ID:                  clause.ID,
					ClauseNumber:        clause.ClauseNumber,
					ClauseTitleEn:       clause.ClauseTitleEn,
					ClauseTitleAr:       *clause.ClauseTitleAr,
					ClauseDescriptionEn: *clause.ClauseDescriptionEn,
					ClauseDescriptionAr: *clause.ClauseDescriptionAr,
				})
			}
		}
	}

	var final []libraryDTOs.GroupedClauseDTO
	for _, g := range result {
		final = append(final, *g)
	}
	return final, nil
}
