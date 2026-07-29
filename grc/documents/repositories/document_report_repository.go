package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type DocumentReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error)
	GetDocumentReportByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.DocumentReport, error)
	GetDocumentReportDetails(reportID uint, tx ...utility.DatabaseTransaction) (*models.DocumentReport, error)
	GetDocumentReportSections(reportID uint) ([]dtos.DocumentReportSectionResponseDTO, error)
	GetPublishedDocumentsForReporting(categoryID, typeID *uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportSourceDocumentDTO, error)
	GetDocumentClausesForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportClauseDTO, error)
	GetDocumentControlsForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportControlDTO, error)
	Create(report *models.DocumentReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.DocumentReportSection, tx ...utility.DatabaseTransaction) error
	CreateDocumentSnapshots(documents []*models.DocumentReportDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocumentReport(reportID uint, tx ...utility.DatabaseTransaction) error
}

type DocumentReportRepository struct {
	db *gorm.DB
}

func NewDocumentReportRepository(db *gorm.DB) DocumentReportRepositoryInterface {
	return &DocumentReportRepository{db: db}
}

func (repo *DocumentReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentReportRepository) GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error) {
	var reports []*dtos.DocumentReportListDTO
	var totalCount int64

	query := repo.db.Table("document_reports").
		Select(`
			document_reports.id, document_reports.title,
			document_reports.document_category_id, document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_reports.document_type_id, document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar, document_reports.created_at,
			created_user.id AS created_user_id, created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color, created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image,
			created_user.mobile_number AS created_user_mobile_number`).
		Joins("LEFT JOIN document_categories ON document_categories.id = document_reports.document_category_id").
		Joins("LEFT JOIN document_types ON document_types.id = document_reports.document_type_id").
		Joins("JOIN users AS created_user ON created_user.id = document_reports.created_by_id")
	if params.Keyword != "" {
		query = query.Where("document_reports.title LIKE ?", "%"+params.Keyword+"%")
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.DocumentReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *DocumentReportRepository) GetDocumentReportByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.DocumentReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var report models.DocumentReport
	if err := db.First(&report, reportID).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (repo *DocumentReportRepository) GetDocumentReportDetails(reportID uint, tx ...utility.DatabaseTransaction) (*models.DocumentReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var report models.DocumentReport
	err := db.
		Preload("CreatedByUser").
		Preload("Documents", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		First(&report, reportID).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (repo *DocumentReportRepository) GetDocumentReportSections(reportID uint) ([]dtos.DocumentReportSectionResponseDTO, error) {
	var sections []dtos.DocumentReportSectionResponseDTO
	err := repo.db.Table("document_report_sections").
		Select("id, type, title_en, title_ar, description_en, description_ar, order_num AS order_number").
		Where("document_report_id = ?", reportID).
		Order("order_num ASC, id ASC").
		Find(&sections).Error
	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (repo *DocumentReportRepository) GetPublishedDocumentsForReporting(categoryID, typeID *uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportSourceDocumentDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documents []dtos.DocumentReportSourceDocumentDTO
	err := repo.publishedDocumentsQuery(db, categoryID, typeID).Find(&documents).Error
	return documents, err
}

func (repo *DocumentReportRepository) publishedDocumentsQuery(db *gorm.DB, categoryID, typeID *uint) *gorm.DB {
	query := db.Table("documents").
		Select(`
			documents.id AS source_document_id, documents.reference_number, documents.title_en, documents.title_ar,
			documents.description_en, documents.description_ar,
			document_categories.title_en AS document_category_title_en, document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en, document_types.title_ar AS document_type_title_ar,
			document_statuses.title_en AS document_status_title_en, document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code, document_statuses.type AS document_status_type,
			document_access_types.title_en AS document_access_type_title_en,
			document_access_types.title_ar AS document_access_type_title_ar,
			document_access_types.color_code AS document_access_type_color_code,
			confidentiality_classifications.title_en AS confidentiality_classification_title_en,
			confidentiality_classifications.title_ar AS confidentiality_classification_title_ar,
			confidentiality_classifications.type AS confidentiality_classification_type,
			review_frequencies.title_en AS review_frequency_title_en, review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			organization_units.title_en AS responsible_unit_title_en, organization_units.title_ar AS responsible_unit_title_ar,
			CONCAT(responsible_user.first_name_en, ' ', responsible_user.last_name_en) AS responsible_user_name_en,
			CONCAT(responsible_user.first_name_ar, ' ', responsible_user.last_name_ar) AS responsible_user_name_ar,
			document_versions.version, COALESCE(document_versions.original_name, document_versions.name) AS file_name,
			document_versions.type AS file_type, document_versions.size AS file_size,
			documents.created_at AS document_created_at, documents.next_review_date,
			COALESCE(documents.published_at, document_versions.published_at) AS published_at`).
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id AND document_statuses.type = ?", "published").
		Joins("JOIN document_categories ON document_categories.id = documents.document_category_id").
		Joins("JOIN document_types ON document_types.id = documents.document_type_id").
		Joins("JOIN document_access_types ON document_access_types.id = documents.document_access_type_id").
		Joins("JOIN confidentiality_classifications ON confidentiality_classifications.id = documents.confidentiality_classification_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = documents.review_frequency_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id").
		Joins("LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = ?", true).
		Where("documents.is_folder = ?", false)
	if categoryID != nil {
		query = query.Where("documents.document_category_id = ?", *categoryID)
	}
	if typeID != nil {
		query = query.Where("documents.document_type_id = ?", *typeID)
	}
	return query.Order("documents.reference_number ASC")
}

func (repo *DocumentReportRepository) GetDocumentClausesForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportClauseDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	clauses := []dtos.DocumentReportClauseDTO{}
	err := db.Table("document_version_clauses").
		Select("document_version_clauses.id, document_version_clauses.clause_number, document_version_clauses.clause_title_en, document_version_clauses.clause_title_ar, document_version_clauses.description").
		Joins("JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id").
		Where("document_versions.document_id = ? AND document_versions.is_latest_version = ?", documentID, true).
		Order("document_version_clauses.id ASC").Find(&clauses).Error
	return clauses, err
}

func (repo *DocumentReportRepository) GetDocumentControlsForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportControlDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	controls := []dtos.DocumentReportControlDTO{}
	err := db.Table("document_controls").
		Select("controls.id, document_controls.document_version_clause_id AS clause_id, controls.reference_number, controls.code, controls.title_en, controls.title_ar, controls.description").
		Joins("JOIN controls ON controls.id = document_controls.control_id").
		Joins("JOIN document_version_clauses ON document_version_clauses.id = document_controls.document_version_clause_id").
		Joins("JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id").
		Where("document_controls.document_id = ? AND document_versions.document_id = ? AND document_versions.is_latest_version = ?", documentID, documentID, true).
		Order("document_version_clauses.id ASC, controls.reference_number ASC").
		Find(&controls).Error
	if err != nil {
		return nil, err
	}
	return controls, nil
}

func (repo *DocumentReportRepository) Create(report *models.DocumentReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(report).Error; err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (repo *DocumentReportRepository) CreateSection(section *models.DocumentReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(section).Error
}

func (repo *DocumentReportRepository) CreateDocumentSnapshots(documents []*models.DocumentReportDocument, tx ...utility.DatabaseTransaction) error {
	if len(documents) == 0 {
		return nil
	}

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (repo *DocumentReportRepository) DeleteDocumentReport(reportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("document_report_id = ?", reportID).Delete(&models.DocumentReportDocument{}).Error; err != nil {
		return err
	}
	if err := db.Where("document_report_id = ?", reportID).Delete(&models.DocumentReportSection{}).Error; err != nil {
		return err
	}
	return db.Where("id = ?", reportID).Delete(&models.DocumentReport{}).Error
}
