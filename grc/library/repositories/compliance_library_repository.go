package repositories

import (
	"strings"

	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceLibraryRepositoryInterface interface {
	GetComplianceLibraryDetailsByTitle(titleEn string) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	GetClausesByVersion(versionID uint) ([]models.ComplianceLibraryClause, error)
	GetControlLibraryTitle(controlLibraryID, complianceLibraryID uint) (*models.ControlLibraryTitle, error)
	FetchRawControlMappings() ([]libraryDTOs.RawControlMapping, error)
	GetControlMappings(params *libraryDTOs.ControlMappingListQueryParamsDTO) (*libraryDTOs.ControlMappingListResponseDTO, error)
}

type ComplianceLibraryRepository struct {
	db *gorm.DB
}

func NewComplianceLibraryRepository(db *gorm.DB) *ComplianceLibraryRepository {
	return &ComplianceLibraryRepository{
		db: db,
	}
}

func (repo *ComplianceLibraryRepository) GetComplianceLibraryDetailsByTitle(titleEn string) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
	var result libraryDTOs.ComplianceLibraryWithVersionsDTO

	err := repo.db.
		Table("compliance_libraries").
		Select(`compliance_libraries.id AS id,
			compliance_libraries.icon,
			compliance_libraries.title_en,
			compliance_libraries.title_ar,
			compliance_libraries.description_en,
			compliance_libraries.description_ar,
			compliance_libraries.issued_by,
			compliance_libraries.classification,
			compliance_libraries.document_classification,
			compliance_libraries.disclaimer,
			compliance_library_versions.id AS version_id,
			compliance_library_versions.version,
			compliance_library_versions.year,
			compliance_library_versions.overview,
			compliance_library_versions.introduction,
			compliance_library_versions.domains`).
		Joins("LEFT JOIN compliance_library_versions ON compliance_libraries.id = compliance_library_versions.compliance_library_id").
		Where("compliance_libraries.title_en = ?", titleEn).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *ComplianceLibraryRepository) GetClausesByVersion(versionID uint) ([]models.ComplianceLibraryClause, error) {
	var clauses []models.ComplianceLibraryClause
	err := repo.db.Preload("Controls").Where("compliance_library_version_id = ?", versionID).Find(&clauses).Error
	return clauses, err
}

func (repo *ComplianceLibraryRepository) GetControlLibraryTitle(controlLibraryID, complianceLibraryID uint) (*models.ControlLibraryTitle, error) {
	var title models.ControlLibraryTitle
	err := repo.db.Where("control_library_id = ? AND compliance_library_id = ?", controlLibraryID, complianceLibraryID).First(&title).Error
	if err != nil {
		return nil, err
	}
	return &title, nil
}

func (repo *ComplianceLibraryRepository) FetchRawControlMappings() ([]libraryDTOs.RawControlMapping, error) {
	var results []libraryDTOs.RawControlMapping

	err := repo.db.Table("control_library_titles AS clt").
		Select(`
			cl.id AS generic_control_id,
			cl.title_en AS generic_control_title_en,
			cl.title_ar AS generic_control_title_ar,
			clt.control_code,
			clt.control_title_en,
			clt.control_title_ar,
			lib.title_en AS framework_en,
			lib.title_ar AS framework_ar
		`).
		Joins("JOIN control_library AS cl ON cl.id = clt.control_library_id").
		Joins("JOIN compliance_libraries AS lib ON lib.id = clt.compliance_library_id").
		Order("cl.id").
		Scan(&results).Error

	return results, err
}

func (repo *ComplianceLibraryRepository) GetControlMappings(params *libraryDTOs.ControlMappingListQueryParamsDTO) (*libraryDTOs.ControlMappingListResponseDTO, error) {
	var items []*libraryDTOs.ControlLibraryResponseDTO
	var totalCount int64

	// Base query from control_library table
	query := repo.db.
		Table("control_library").
		Select("id, title_en, title_ar")

	// Filtering
	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where("title_en LIKE ? OR title_ar LIKE ?", kw, kw)
	}
	if len(params.TitleEnList) > 0 {
		query = query.Where("title_en IN ?", params.TitleEnList)
	}

	// Count total
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := "asc"
	if params.Sort == "desc" {
		sortOrder = "desc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute base query
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}

	// Attach mapped controls per framework (SAMA/NCA) as arrays
	for i := range items {
		var titles []libraryDTOs.ControlTitleDTO
		subQuery := repo.db.
			Table("control_library_titles AS clt").
			Select(`
				clt.control_library_id,
				cpl.title_en AS compliance_library_title_en,
				cpl.title_ar AS compliance_library_title_ar,
				clt.control_code,
				clt.control_title_en,
				clt.control_title_ar
			`).
			Joins("LEFT JOIN compliance_libraries AS cpl ON cpl.id = clt.compliance_library_id").
			Where("clt.control_library_id = ?", items[i].ID)

		if len(params.ControlTitleEnList) > 0 {
			subQuery = subQuery.Where("clt.control_title_en IN ?", params.ControlTitleEnList)
		}
		if len(params.ComplianceLibraryTitlesEn) > 0 {
			subQuery = subQuery.Where("cpl.title_en IN ?", params.ComplianceLibraryTitlesEn)
		}

		if err := subQuery.Find(&titles).Error; err != nil {
			return nil, err
		}

		// Assign multiple mapped controls by framework
		for _, t := range titles {
			switch strings.ToLower(t.ComplianceLibraryTitleEn) {
			case "cyber security framework | saudi arabian monetary authority":
				items[i].SAMA = append(items[i].SAMA, t)
			case "essential cybersecurity controls (ecc – 2: 2024)":
				items[i].NCA = append(items[i].NCA, t)
			}
		}
	}

	// Final response
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	return &libraryDTOs.ControlMappingListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}
