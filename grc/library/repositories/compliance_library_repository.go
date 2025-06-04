package repositories

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceLibraryRepositoryInterface interface {
	GetComplianceLibraryDetailsByID(libraryID uint) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	GetClauseTreeByVersion(versionID uint) ([]models.ComplianceLibraryClause, error)
	GetControlLibraryTitle(controlLibraryID, complianceLibraryID uint) (*models.ControlLibraryTitle, error)
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

func (repo *ComplianceLibraryRepository) GetClauseTreeByVersion(versionID uint) ([]models.ComplianceLibraryClause, error) {
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
