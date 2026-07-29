package repositories

import (
	"sort"

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
	GetDocumentReportDetails(reportID uint) (*models.DocumentReport, error)
	GetDocumentReportSections(reportID uint) ([]dtos.DocumentReportSectionResponseDTO, error)
	GetPublishedDocumentsForReporting(categoryID, typeID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportSourceDocumentDTO, error)
	GetDocumentClausesForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportClauseDTO, error)
	GetDocumentControlsForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportControlDTO, error)
	GetDocumentAssociatedItemsForReporting(documentID uint, tx ...utility.DatabaseTransaction) (*dtos.DocumentReportAssociatedItemsDTO, error)
	Create(report *models.DocumentReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.DocumentReportSection, tx ...utility.DatabaseTransaction) error
	CreateDocumentSnapshots(documents []*models.DocumentReportDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocumentReport(reportID uint, tx ...utility.DatabaseTransaction) error
}

type DocumentReportRepository struct{ db *gorm.DB }

func NewDocumentReportRepository(db *gorm.DB) DocumentReportRepositoryInterface {
	return &DocumentReportRepository{db: db}
}

func (repo *DocumentReportRepository) getDB(tx []utility.DatabaseTransaction) *gorm.DB {
	if len(tx) > 0 && tx[0] != nil {
		return tx[0].DB()
	}
	return repo.db
}

func (repo *DocumentReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentReportRepository) GetAll(params *dtos.DocumentReportListQueryParamsDTO) (*dtos.DocumentReportListResponseDTO, error) {
	var items []*dtos.DocumentReportListDTO
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
			created_user.last_name_ar AS created_user_last_name_ar, created_user.initial AS created_user_initial`).
		Joins("LEFT JOIN document_categories ON document_categories.id = document_reports.document_category_id").
		Joins("LEFT JOIN document_types ON document_types.id = document_reports.document_type_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = document_reports.created_by_id")
	if params.Keyword != "" {
		query = query.Where("document_reports.title LIKE ?", "%"+params.Keyword+"%")
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort).
		Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	if err := query.Find(&items).Error; err != nil {
		return nil, err
	}
	return &dtos.DocumentReportListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *DocumentReportRepository) GetDocumentReportByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.DocumentReport, error) {
	var report models.DocumentReport
	if err := repo.getDB(tx).First(&report, reportID).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

func (repo *DocumentReportRepository) GetDocumentReportDetails(reportID uint) (*models.DocumentReport, error) {
	var report models.DocumentReport
	err := repo.db.Preload("CreatedByUser").
		Preload("Documents", func(db *gorm.DB) *gorm.DB { return db.Order("id ASC") }).
		First(&report, reportID).Error
	return &report, err
}

func (repo *DocumentReportRepository) GetDocumentReportSections(reportID uint) ([]dtos.DocumentReportSectionResponseDTO, error) {
	var sections []dtos.DocumentReportSectionResponseDTO
	err := repo.db.Table("document_report_sections").
		Select("id, type, title_en, title_ar, description_en, description_ar, order_num AS order_number").
		Where("document_report_id = ?", reportID).Order("order_num ASC, id ASC").Find(&sections).Error
	return sections, err
}

func (repo *DocumentReportRepository) GetPublishedDocumentsForReporting(categoryID, typeID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportSourceDocumentDTO, error) {
	var documents []dtos.DocumentReportSourceDocumentDTO
	err := repo.publishedDocumentsQuery(repo.getDB(tx), categoryID, typeID).Find(&documents).Error
	return documents, err
}

func (repo *DocumentReportRepository) publishedDocumentsQuery(db *gorm.DB, categoryID, typeID uint) *gorm.DB {
	return db.Table("documents").
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
		Where("documents.document_category_id = ? AND documents.document_type_id = ? AND documents.is_folder = ?", categoryID, typeID, false).
		Order("documents.reference_number ASC")
}

func (repo *DocumentReportRepository) GetDocumentClausesForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportClauseDTO, error) {
	clauses := []dtos.DocumentReportClauseDTO{}
	err := repo.getDB(tx).Table("document_version_clauses").
		Select("document_version_clauses.id, document_version_clauses.clause_number, document_version_clauses.clause_title_en, document_version_clauses.clause_title_ar, document_version_clauses.description").
		Joins("JOIN document_versions ON document_versions.id = document_version_clauses.document_version_id").
		Where("document_versions.document_id = ? AND document_versions.is_latest_version = ?", documentID, true).
		Order("document_version_clauses.id ASC").Find(&clauses).Error
	return clauses, err
}

func (repo *DocumentReportRepository) GetDocumentControlsForReporting(documentID uint, tx ...utility.DatabaseTransaction) ([]dtos.DocumentReportControlDTO, error) {
	db := repo.getDB(tx)
	clauseControls := []dtos.DocumentReportControlDTO{}
	err := db.Table("document_controls").
		Select("controls.id, document_controls.document_version_clause_id AS clause_id, controls.reference_number, controls.code, controls.title_en, controls.title_ar, controls.description").
		Joins("JOIN controls ON controls.id = document_controls.control_id").
		Joins("LEFT JOIN document_versions ON document_versions.id = document_controls.document_version_id").
		Where("document_controls.document_id = ? AND (document_controls.document_version_id IS NULL OR document_versions.is_latest_version = ?)", documentID, true).
		Order("controls.reference_number ASC").Find(&clauseControls).Error
	if err != nil {
		return nil, err
	}

	associatedControls := []dtos.DocumentReportControlDTO{}
	err = db.Table("control_associated_documents").
		Select("controls.id, controls.reference_number, controls.code, controls.title_en, controls.title_ar, controls.description").
		Joins("JOIN controls ON controls.id = control_associated_documents.control_id").
		Where("control_associated_documents.document_id = ?", documentID).
		Order("controls.reference_number ASC").Find(&associatedControls).Error
	if err != nil {
		return nil, err
	}

	return mergeDocumentReportControls(clauseControls, associatedControls), nil
}

func mergeDocumentReportControls(
	clauseControls []dtos.DocumentReportControlDTO,
	associatedControls []dtos.DocumentReportControlDTO,
) []dtos.DocumentReportControlDTO {
	controlsByID := make(map[uint]dtos.DocumentReportControlDTO, len(clauseControls)+len(associatedControls))
	for _, control := range associatedControls {
		controlsByID[control.ID] = control
	}
	for _, control := range clauseControls {
		// Prefer the clause-linked entry because it retains the clause relationship.
		controlsByID[control.ID] = control
	}
	controls := make([]dtos.DocumentReportControlDTO, 0, len(controlsByID))
	for _, control := range controlsByID {
		controls = append(controls, control)
	}
	sort.Slice(controls, func(i, j int) bool {
		return controls[i].ReferenceNumber < controls[j].ReferenceNumber
	})
	return controls
}

func (repo *DocumentReportRepository) GetDocumentAssociatedItemsForReporting(documentID uint, tx ...utility.DatabaseTransaction) (*dtos.DocumentReportAssociatedItemsDTO, error) {
	db := repo.getDB(tx)
	result := &dtos.DocumentReportAssociatedItemsDTO{
		Risks: []dtos.DocumentReportAssociatedItemDTO{}, Incidents: []dtos.DocumentReportAssociatedItemDTO{},
		Assets: []dtos.DocumentReportAssociatedItemDTO{}, StrategicObjectives: []dtos.DocumentReportAssociatedItemDTO{},
	}
	if err := db.Table("risk_associated_documents").
		Select("risk_register.id, COALESCE(risk_register.reference_number, '') AS reference_number, risk_register.title AS title_en, risk_register.title AS title_ar, risk_register.description").
		Joins("JOIN risk_register ON risk_register.id = risk_associated_documents.risk_register_id").
		Where("risk_associated_documents.document_id = ?", documentID).
		Order("risk_register.reference_number ASC, risk_register.id ASC").Find(&result.Risks).Error; err != nil {
		return nil, err
	}
	if err := db.Table("incident_associated_documents").
		Select("incident_register.id, incident_register.reference_number, incident_register.title AS title_en, incident_register.title AS title_ar, incident_register.description").
		Joins("JOIN incident_register ON incident_register.id = incident_associated_documents.incident_register_id").
		Where("incident_associated_documents.document_id = ?", documentID).
		Order("incident_register.reference_number ASC, incident_register.id ASC").Find(&result.Incidents).Error; err != nil {
		return nil, err
	}
	if err := db.Table("asset_associated_documents").
		Select("asset_register.id, asset_register.reference_number, asset_register.title_en, asset_register.title_ar, asset_register.description").
		Joins("JOIN asset_register ON asset_register.id = asset_associated_documents.asset_register_id").
		Where("asset_associated_documents.document_id = ?", documentID).
		Order("asset_register.reference_number ASC, asset_register.id ASC").Find(&result.Assets).Error; err != nil {
		return nil, err
	}
	if err := db.Table("strategic_objective_associated_documents").
		Select("strategic_objectives.id, strategic_objectives.reference_number, strategic_objectives.title_en, strategic_objectives.title_ar, strategic_objectives.description_en AS description").
		Joins("JOIN strategic_objectives ON strategic_objectives.id = strategic_objective_associated_documents.strategic_objective_id").
		Where("strategic_objective_associated_documents.document_id = ?", documentID).
		Order("strategic_objectives.reference_number ASC, strategic_objectives.id ASC").Find(&result.StrategicObjectives).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (repo *DocumentReportRepository) Create(report *models.DocumentReport, tx ...utility.DatabaseTransaction) (uint, error) {
	if err := repo.getDB(tx).Create(report).Error; err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (repo *DocumentReportRepository) CreateSection(section *models.DocumentReportSection, tx ...utility.DatabaseTransaction) error {
	return repo.getDB(tx).Create(section).Error
}

func (repo *DocumentReportRepository) CreateDocumentSnapshots(documents []*models.DocumentReportDocument, tx ...utility.DatabaseTransaction) error {
	if len(documents) == 0 {
		return nil
	}
	return repo.getDB(tx).Create(&documents).Error
}

func (repo *DocumentReportRepository) DeleteDocumentReport(reportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.getDB(tx)
	if err := db.Where("document_report_id = ?", reportID).Delete(&models.DocumentReportDocument{}).Error; err != nil {
		return err
	}
	if err := db.Where("document_report_id = ?", reportID).Delete(&models.DocumentReportSection{}).Error; err != nil {
		return err
	}
	return db.Where("id = ?", reportID).Delete(&models.DocumentReport{}).Error
}
