package repositories

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceLibraryRepositoryInterface interface {
	GetComplianceLibraryDetailsByID(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	GetClausesByVersionID(versionID uint) ([]libraryDTOs.GroupedClauseDTO, error)
}

type ComplianceLibraryRepository struct {
	db *gorm.DB
}

func NewComplianceLibraryRepository(db *gorm.DB) *ComplianceLibraryRepository {
	return &ComplianceLibraryRepository{
		db: db,
	}
}

func (repo *ComplianceLibraryRepository) GetComplianceLibraryDetailsByID(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error) {
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
		Where("compliance_libraries.id = ?", libraryID).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *ComplianceLibraryRepository) GetClausesByVersionID(versionID uint) ([]libraryDTOs.GroupedClauseDTO, error) {
	var mainClauses []models.ComplianceLibraryClause
	err := repo.db.
		Where("compliance_library_version_id = ? AND compliance_library_clause_id IS NULL", versionID).
		Order(`
			TRY_CAST(LEFT(clause_number, CHARINDEX('.', clause_number + '.') - 1) AS INT)
		`).
		Find(&mainClauses).Error
	if err != nil {
		return nil, err
	}

	var result []libraryDTOs.GroupedClauseDTO
	for _, main := range mainClauses {
		group := libraryDTOs.GroupedClauseDTO{
			ID:            main.ID,
			ClauseNumber:  main.ClauseNumber,
			ClauseTitleEn: main.ClauseTitleEn,
			ClauseTitleAr: *main.ClauseTitleAr,
			CreatedAt:     main.CreatedAt,
			UpdatedAt:     *main.UpdatedAt,
		}

		var subClauses []models.ComplianceLibraryClause
		err = repo.db.
			Where("compliance_library_clause_id = ?", main.ID).
			Order(`
				TRY_CAST(SUBSTRING(clause_number, CHARINDEX('.', clause_number + '.') + 1, LEN(clause_number)) AS INT)
			`).
			Find(&subClauses).Error
		if err != nil {
			return nil, err
		}

		for _, sub := range subClauses {
			group.SubClauses = append(group.SubClauses, libraryDTOs.SubClauseDTO{
				ID:                  sub.ID,
				ClauseNumber:        sub.ClauseNumber,
				ClauseTitleEn:       sub.ClauseTitleEn,
				ClauseTitleAr:       *sub.ClauseTitleAr,
				ClauseDescriptionEn: *sub.ClauseDescriptionEn,
				ClauseDescriptionAr: *sub.ClauseDescriptionAr,
			})
		}

		result = append(result, group)
	}

	return result, nil
}
