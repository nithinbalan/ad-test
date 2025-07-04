package repositories

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceLibraryRepositoryInterface interface {
	GetComplianceLibraryDetailsByTitle(titleEn string) (*libraryDTOs.ComplianceLibraryWithVersionsDTO, error)
	GetClausesByVersion(versionID uint) ([]models.ComplianceLibraryClause, error)
	GetControlLibraryTitle(controlLibraryID, complianceLibraryID uint) (*models.ControlLibraryTitle, error)
	FetchRawControlMappings() ([]libraryDTOs.RawControlMapping, error)
	GetGenericControlMappings() ([]libraryDTOs.ControlLibraryResponseDTO, error)
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

func (repo *ComplianceLibraryRepository) GetGenericControlMappings() ([]libraryDTOs.ControlLibraryResponseDTO, error) {
	rows, err := repo.db.
		Table("control_library AS cl").
		Select(`
			cl.id AS cl_id,
			cl.title_en,
			cl.title_ar,
			clt.control_code,
			clt.control_title_en,
			clt.control_title_ar,
			cpl.title_en AS compliance_library_title_en
		`).
		Joins("LEFT JOIN control_library_titles AS clt ON cl.id = clt.control_library_id").
		Joins("LEFT JOIN compliance_libraries AS cpl ON clt.compliance_library_id = cpl.id").
		Order("cl.id").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	controlMap := make(map[uint]*libraryDTOs.ControlLibraryResponseDTO)
	var controlOrder []uint

	for rows.Next() {
		var (
			clID                     uint
			titleEn, titleAr         string
			controlCode              *string
			controlTitleEn           *string
			controlTitleAr           *string
			complianceLibraryTitleEn *string
		)

		err := rows.Scan(
			&clID, &titleEn, &titleAr,
			&controlCode, &controlTitleEn, &controlTitleAr,
			&complianceLibraryTitleEn,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := controlMap[clID]; !exists {
			controlMap[clID] = &libraryDTOs.ControlLibraryResponseDTO{
				TitleEn: titleEn,
				TitleAr: titleAr,
				Titles:  []libraryDTOs.ControlTitleDTO{},
			}
			controlOrder = append(controlOrder, clID)
		}

		if controlCode != nil || controlTitleEn != nil || controlTitleAr != nil {
			controlMap[clID].Titles = append(controlMap[clID].Titles, libraryDTOs.ControlTitleDTO{
				ComplianceLibraryTitleEn: *complianceLibraryTitleEn,
				ControlCode:              *controlCode,
				ControlTitleEn:           *controlTitleEn,
				ControlTitleAr:           *controlTitleAr,
			})
		}
	}

	var result []libraryDTOs.ControlLibraryResponseDTO
	for _, clID := range controlOrder {
		result = append(result, *controlMap[clID])
	}

	return result, nil
}
