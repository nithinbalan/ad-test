package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type ComplianceLibraryServiceInterface interface {
	GetComplianceLibraryDetailsByID(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	// GetClausesByVersionID(versionID uint) ([]libraryDTOs.GroupedClauseDTO, error)
	GetClauseTreeByVersion(versionID, complianceLibraryID uint) ([]libraryDTOs.ClauseWithSubClausesDTO, error)
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

// func (cls *ComplianceLibraryService) GetClausesByVersionID(versionID uint) ([]libraryDTOs.GroupedClauseDTO, error) {
// 	clauses, err := cls.repo.GetAllClausesByVersionID(versionID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mainMap := make(map[uint]*libraryDTOs.GroupedClauseDTO)
// 	var result []*libraryDTOs.GroupedClauseDTO

// 	for _, clause := range clauses {
// 		if clause.ComplianceLibraryClauseID == nil {
// 			group := &libraryDTOs.GroupedClauseDTO{
// 				ID:            clause.ID,
// 				ClauseNumber:  clause.ClauseNumber,
// 				ClauseTitleEn: clause.ClauseTitleEn,
// 				ClauseTitleAr: *clause.ClauseTitleAr,
// 				CreatedAt:     clause.CreatedAt,
// 				UpdatedAt:     *clause.UpdatedAt,
// 				SubClauses:    []libraryDTOs.SubClauseDTO{},
// 			}
// 			mainMap[clause.ID] = group
// 			result = append(result, group)
// 		} else {
// 			parent := mainMap[*clause.ComplianceLibraryClauseID]
// 			if parent != nil {
// 				parent.SubClauses = append(parent.SubClauses, libraryDTOs.SubClauseDTO{
// 					ID:            clause.ID,
// 					ClauseNumber:  clause.ClauseNumber,
// 					ClauseTitleEn: clause.ClauseTitleEn,
// 					ClauseTitleAr: *clause.ClauseTitleAr,
// 				})
// 			}
// 		}
// 	}

// 	var final []libraryDTOs.GroupedClauseDTO
// 	for _, g := range result {
// 		final = append(final, *g)
// 	}
// 	return final, nil
// }

func (svc *ComplianceLibraryService) GetClauseTreeByVersion(versionID, complianceLibraryID uint) ([]libraryDTOs.ClauseWithSubClausesDTO, error) {
	clauses, err := svc.repo.GetClauseTreeByVersion(versionID)
	if err != nil {
		return nil, err
	}

	clauseMap := make(map[uint]*libraryDTOs.ClauseWithSubClausesDTO)
	var topLevel []*libraryDTOs.ClauseWithSubClausesDTO

	for _, clause := range clauses {
		if clause.ComplianceLibraryClauseID == nil {
			main := &libraryDTOs.ClauseWithSubClausesDTO{
				ID:            clause.ID,
				ClauseNumber:  clause.ClauseNumber,
				ClauseTitleEn: clause.ClauseTitleEn,
				ClauseTitleAr: *clause.ClauseTitleAr,
				CreatedAt:     clause.CreatedAt.Format("2006-01-02"),
				UpdatedAt:     clause.UpdatedAt.Format("2006-01-02"),
			}
			clauseMap[clause.ID] = main
			topLevel = append(topLevel, main)
		}
	}

	for _, clause := range clauses {
		if clause.ComplianceLibraryClauseID != nil {
			parent := clauseMap[*clause.ComplianceLibraryClauseID]
			if parent != nil {
				sub := libraryDTOs.SubClauseWithControlsDTO{
					ID:                  clause.ID,
					ClauseNumber:        clause.ClauseNumber,
					ClauseTitleEn:       clause.ClauseTitleEn,
					ClauseTitleAr:       *clause.ClauseTitleAr,
					ClauseDescriptionEn: *clause.ClauseDescriptionEn,
					ClauseDescriptionAr: *clause.ClauseDescriptionAr,
				}

				for _, ctrl := range clause.Controls {
					controlTitle, err := svc.repo.GetControlLibraryTitle(ctrl.ControlLibraryID, complianceLibraryID)
					if err == nil {
						sub.Controls = append(sub.Controls, libraryDTOs.ControlDTO{
							ID:             controlTitle.ControlLibraryID,
							ControlCode:    controlTitle.ControlCode,
							ControlTitleEn: controlTitle.ControlTitleEn,
							ControlTitleAr: controlTitle.ControlTitleAr,
						})
					}
				}
				parent.SubClauses = append(parent.SubClauses, sub)
			}
		}
	}

	var final []libraryDTOs.ClauseWithSubClausesDTO
	for _, cl := range topLevel {
		final = append(final, *cl)
	}
	return final, nil
}
