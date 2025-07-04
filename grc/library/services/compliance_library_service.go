package services

import (
	"fmt"
	"sort"

	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type ComplianceLibraryServiceInterface interface {
	GetNCADetails() (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	GetSAMADetails() (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)

	GetClauseTreeByVersion(versionID, complianceLibraryID uint) ([]libraryDTOs.ClauseWithSubClausesDTO, error)
	GetGroupedControlMappings() ([]libraryDTOs.ControlMappingDTO, error)
	GetGenericControlMappings() ([]libraryDTOs.ControlLibraryResponseDTO, error)
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

func (cls *ComplianceLibraryService) GetNCADetails() (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
	titleEn := "Essential Cybersecurity Controls (ECC – 2: 2024)"
	return cls.repo.GetComplianceLibraryDetailsByTitle(titleEn)
}

func (cls *ComplianceLibraryService) GetSAMADetails() (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
	titleEn := "Cyber Security Framework | Saudi Arabian Monetary Authority"
	return cls.repo.GetComplianceLibraryDetailsByTitle(titleEn)
}

func (cls *ComplianceLibraryService) GetClauseTreeByVersion(versionID, complianceLibraryID uint) ([]libraryDTOs.ClauseWithSubClausesDTO, error) {
	clauses, err := cls.repo.GetClauseTreeByVersion(versionID)
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
					controlTitle, err := cls.repo.GetControlLibraryTitle(ctrl.ControlLibraryID, complianceLibraryID)
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

func (cls *ComplianceLibraryService) GetGroupedControlMappings() ([]libraryDTOs.ControlMappingDTO, error) {
	rawData, err := cls.repo.FetchRawControlMappings()
	if err != nil {
		return nil, err
	}

	controlMap := make(map[uint]*libraryDTOs.ControlMappingDTO)

	for _, item := range rawData {
		// If this generic control is not added yet, add its base entry
		if _, exists := controlMap[item.GenericControlID]; !exists {
			controlMap[item.GenericControlID] = &libraryDTOs.ControlMappingDTO{
				ID:         fmt.Sprintf("CTRL-%03d", item.GenericControlID),
				TitleEn:    item.GenericControlTitleEn,
				TitleAr:    item.GenericControlTitleAr,
				Frameworks: []libraryDTOs.FrameworkControlsDTO{},
			}
		}

		// Find or add the framework inside that control
		mapping := controlMap[item.GenericControlID]
		found := false

		for i := range mapping.Frameworks {
			if mapping.Frameworks[i].Framework == item.FrameworkEn {
				mapping.Frameworks[i].Controls = append(mapping.Frameworks[i].Controls, libraryDTOs.ControlItemDTO{
					ControlCode:    item.ControlCode,
					ControlTitleEn: item.ControlTitleEn,
					ControlTitleAr: item.ControlTitleAr,
				})
				found = true
				break
			}
		}

		if !found {
			mapping.Frameworks = append(mapping.Frameworks, libraryDTOs.FrameworkControlsDTO{
				Framework: item.FrameworkEn,
				Controls: []libraryDTOs.ControlItemDTO{
					{
						ControlCode:    item.ControlCode,
						ControlTitleEn: item.ControlTitleEn,
						ControlTitleAr: item.ControlTitleAr,
					},
				},
			})
		}
	}

	// Convert map to slice
	var result []libraryDTOs.ControlMappingDTO
	for _, v := range controlMap {
		result = append(result, *v)
	}

	// Optional: sort by ID
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

func (cls *ComplianceLibraryService) GetGenericControlMappings() ([]libraryDTOs.ControlLibraryResponseDTO, error) {
	return cls.repo.GetGenericControlMappings()
}
