package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/risk/dtos"

	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(params *dtos.RiskReportListQueryParamsDTO) (*dtos.RiskReportListResponseDTO, error)
	// GetFullDetailsByID(riskReportID uint) (*models.RiskReport, error)

	GetRiskReportByID(riskReportID uint, tx ...utility.DatabaseTransaction) (*dtos.RiskReportListDTO, error)
	GetRiskReportSections(riskReportID uint) ([]dtos.RiskReportSectionResponseDTO, error)
	GetRiskIdentifications(riskReportID uint) ([]dtos.RiskIdentificationResponseDTO, error)
	GetRiskActionPlans(riskReportID uint) ([]dtos.RiskReportActionPlanResponseDTO, error)
	GetRiskKRIs(riskReportID uint) ([]dtos.RiskReportKriResponseDTO, error)
	GetRiskReportDocuments(riskReportID uint) ([]dtos.RiskReportDocumentDTO, error)

	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.RiskReportSection, error)

	Create(riskReport *models.RiskReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateRiskReportOrganizationUnit(orgUnit *models.RiskReportOrganizationUnit, tx ...utility.DatabaseTransaction) error

	GetRisksByDateRange(startDate time.Time, endDate time.Time, organizationUnitIDs []uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportIdentifiedRiskDTO, error)
	GetControlTitlesByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]string, error)
	GetActionPlansByRiskID(riskRegisterID uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportActionPlanDTO, error)
	GetActionPlanControlsByActionPlanID(actionPlanID uint, tx ...utility.DatabaseTransaction) (string, error)
	GetKRIsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportKRIDTO, error)
	GetDocumentsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportDocumentDTO, error)

	CreateSection(section *models.RiskReportSection, tx ...utility.DatabaseTransaction) error
	CreateRiskIdentification(identification *models.RiskReportIdentification, tx ...utility.DatabaseTransaction) error
	CreateRiskActionPlan(actionPlan *models.RiskReportActionPlan, tx ...utility.DatabaseTransaction) error
	CreateRiskKRI(riskKRI *models.RiskReportKRI, tx ...utility.DatabaseTransaction) error
	CreateRiskDocument(riskDocument *models.RiskReportDocument, tx ...utility.DatabaseTransaction) error

	UpdateSection(section *models.RiskReportSection, tx ...utility.DatabaseTransaction) error

	DeleteRiskReportOrganizationUnits(riskReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskKRIs(riskReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskActionPlans(riskReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskIdentifications(riskReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskReportSections(riskReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteRiskReport(riskReportID uint, tx ...utility.DatabaseTransaction) error

	DeleteRiskReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error

	ReorderRiskReportSections(riskReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error

	GetRiskReportDocumentByID(riskReportID uint, documentID uint, tx ...utility.DatabaseTransaction) (*models.RiskReportDocument, error)
	CreateRiskDocuments(documents []*models.RiskReportDocument, tx ...utility.DatabaseTransaction) error
	DeleteRiskReportDocument(riskReportID uint, documentID uint, tx ...utility.DatabaseTransaction) error
}

type RiskReportRepository struct {
	db *gorm.DB
}

func NewRiskReportRepository(db *gorm.DB) RiskReportRepositoryInterface {
	return &RiskReportRepository{db: db}
}

func (repo *RiskReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskReportRepository) GetAll(params *dtos.RiskReportListQueryParamsDTO) (*dtos.RiskReportListResponseDTO, error) {
	var reports []*dtos.RiskReportListDTO
	var totalCount int64

	query := repo.db.Table("risk_reports").
		Select(`
			risk_reports.id, 
			risk_reports.title, 
			CONVERT(VARCHAR, risk_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, risk_reports.end_date, 23) AS end_date,
			risk_reports.created_at,
			STRING_AGG(organization_units.title_en, ', ') AS organization_unit_titles_en,
			STRING_AGG(organization_units.title_ar, ', ') AS organization_unit_titles_ar
		`).
		Joins("LEFT JOIN risk_report_organization_units ON risk_report_organization_units.risk_report_id = risk_reports.id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_report_organization_units.organization_unit_id").
		Group("risk_reports.id, risk_reports.title, risk_reports.start_date, risk_reports.end_date, risk_reports.created_at")

	if params.Keyword != "" {
		query = query.Where(`
			(risk_reports.title LIKE ? OR 
			organization_units.title_en LIKE ? OR 
			organization_units.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
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

	response := &dtos.RiskReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

// func (repo *RiskReportRepository) GetFullDetailsByID(riskReportID uint) (*models.RiskReport, error) {
// 	var riskReport models.RiskReport

// 	err := repo.db.
// 		Preload("CreatedByUser").
// 		Preload("Sections").
// 		Preload("Identifications").
// 		Preload("Identifications.ActionPlans").
// 		Preload("Identifications.KRIs").
// 		Preload("OrganizationUnits.OrganizationUnit").
// 		Where("id = ?", riskReportID).
// 		First(&riskReport).Error

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &riskReport, nil
// }

func (repo *RiskReportRepository) GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.RiskReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskReportSection models.RiskReportSection
	err := db.
		Where("id = ?", sectionID).First(&riskReportSection).Error
	if err != nil {
		return nil, err
	}
	return &riskReportSection, nil
}

// Fetch main report details
// func (repo *RiskReportRepository) GetRiskReportByID(riskReportID uint, tx ...utility.DatabaseTransaction) (*models.RiskReport, error) {
// 	db := repo.db
// 	if len(tx) > 0 && tx[0] != nil {
// 		db = tx[0].DB()
// 	}
// 	var riskReport models.RiskReport
// 	err := db.
// 		Preload("CreatedByUser").
// 		Preload("OrganizationUnits.OrganizationUnit").
// 		Where("id = ?", riskReportID).First(&riskReport).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &riskReport, nil
// }

func (repo *RiskReportRepository) GetRiskReportByID(riskReportID uint, tx ...utility.DatabaseTransaction) (*dtos.RiskReportListDTO, error) {
	var report dtos.RiskReportListDTO

	query := repo.db.Table("risk_reports").
		Select(`
			risk_reports.id, 
			risk_reports.title, 
			CONVERT(VARCHAR, risk_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, risk_reports.end_date, 23) AS end_date,
			risk_reports.created_at,
			users.id AS created_by_id,
			users.first_name_en AS created_user_first_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_en AS created_user_last_name_en,
			users.last_name_ar AS created_user_last_name_ar,
			users.email AS created_user_email,
			COALESCE(STRING_AGG(organization_units.title_en, ', '), '') AS organization_unit_titles_en,
			COALESCE(STRING_AGG(organization_units.title_ar, ', '), '') AS organization_unit_titles_ar
		`).
		Joins("LEFT JOIN users ON users.id = risk_reports.created_by_id").
		Joins("LEFT JOIN risk_report_organization_units ON risk_report_organization_units.risk_report_id = risk_reports.id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_report_organization_units.organization_unit_id").
		Where("risk_reports.id = ?", riskReportID).
		Group("risk_reports.id, risk_reports.title, risk_reports.start_date, risk_reports.end_date, risk_reports.created_at, users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.email")

	if err := query.First(&report).Error; err != nil {

		return nil, err
	}

	return &report, nil
}

// Fetch sections for the report
func (repo *RiskReportRepository) GetRiskReportSections(riskReportID uint) ([]dtos.RiskReportSectionResponseDTO, error) {
	var sections []dtos.RiskReportSectionResponseDTO
	err := repo.db.Table("risk_report_sections").
		Select("id, type, title_en, title_ar, description_en, description_ar, order_num").
		Where("risk_report_id = ?", riskReportID).
		Order("order_num ASC").
		Find(&sections).Error

	if err != nil {
		return nil, err
	}
	return sections, nil
}

// Fetch identified risks for the report
func (repo *RiskReportRepository) GetRiskIdentifications(riskReportID uint) ([]dtos.RiskIdentificationResponseDTO, error) {
	var risks []dtos.RiskIdentificationResponseDTO
	err := repo.db.Table("risk_report_identifications").
		Select(`
			id, reference_number, title, description, risk_category_title_en, risk_category_title_ar,
			risk_source_title_en, risk_source_title_ar, risk_reported_at,
			responsible_unit_title_en, responsible_unit_title_ar,
			inherent_impact_title_en, inherent_impact_title_ar, inherent_likelihood_title_en, inherent_likelihood_title_ar,
			inherent_risk_rating_title_en, inherent_risk_rating_title_ar, inherent_risk_rating_color_code,
			residual_impact_title_en, residual_impact_title_ar, residual_likelihood_title_en, residual_likelihood_title_ar,
			residual_risk_rating_title_en, residual_risk_rating_title_ar, residual_risk_rating_color_code,
			risk_status_title_en, risk_status_title_ar, risk_status_color_code,
			existing_controls, proposed_controls
		`).
		Where("risk_report_id = ?", riskReportID).
		Order("risk_reported_at DESC").
		Find(&risks).Error

	if err != nil {
		return nil, err
	}
	return risks, nil
}

// Fetch risk action plans
func (repo *RiskReportRepository) GetRiskActionPlans(riskReportID uint) ([]dtos.RiskReportActionPlanResponseDTO, error) {
	var actionPlans []dtos.RiskReportActionPlanResponseDTO
	err := repo.db.Table("risk_report_action_plans").
		Select(`
			risk_report_action_plans.id, risk_report_action_plans.risk_report_identification_id,
			risk_report_identifications.reference_number AS risk_reference_number,
			risk_report_identifications.title AS risk_title,
			action_plan_title, action_plan_description,
			severity_level_title_en, severity_level_title_ar,
			responsible_user_name_en, responsible_user_name_ar,
			target_date,
			action_plan_status_title_en, action_plan_status_title_ar, action_plan_status_color_code,
			risk_report_action_plans.proposed_controls
		`).
		Joins("LEFT JOIN risk_report_identifications ON risk_report_identifications.id = risk_report_action_plans.risk_report_identification_id").
		Where("risk_report_action_plans.risk_report_id = ?", riskReportID).
		Order("target_date ASC").
		Find(&actionPlans).Error

	if err != nil {
		return nil, err
	}
	return actionPlans, nil
}

// Fetch risk KRIs
func (repo *RiskReportRepository) GetRiskKRIs(riskReportID uint) ([]dtos.RiskReportKriResponseDTO, error) {
	var kris []dtos.RiskReportKriResponseDTO
	err := repo.db.Table("risk_report_kris").
		Select(`
			risk_report_kris.id,
			risk_report_identifications.reference_number AS risk_reference_number,
			risk_report_identifications.title AS risk_title,
			risk_report_kris.title, risk_report_kris.description,
			risk_report_kris.low_threshold, risk_report_kris.medium_threshold,
			risk_report_kris.unit_title_en, risk_report_kris.unit_title_ar, risk_report_kris.current_value
		`).
		Joins("LEFT JOIN risk_report_identifications ON risk_report_identifications.id = risk_report_kris.risk_report_identification_id").
		Where("risk_report_kris.risk_report_id = ?", riskReportID).
		Order("risk_report_kris.title ASC").
		Find(&kris).Error

	if err != nil {
		return nil, err
	}
	return kris, nil
}

func (repo *RiskReportRepository) GetRiskReportDocuments(riskReportID uint) ([]dtos.RiskReportDocumentDTO, error) {
	var documents []dtos.RiskReportDocumentDTO
	err := repo.db.Table("risk_report_documents").
		Select("id, risk_report_id, risk_report_identification_id, url, name, size, original_name, type, token, created_at").
		Where("risk_report_id = ?", riskReportID).
		Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *RiskReportRepository) Create(riskReport *models.RiskReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(riskReport).Error; err != nil {
		return 0, err
	}
	return riskReport.ID, nil
}

func (repo *RiskReportRepository) CreateRiskIdentification(identification *models.RiskReportIdentification, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(identification).Error
}

func (repo *RiskReportRepository) GetRisksByDateRange(startDate time.Time, endDate time.Time, organizationUnitIDs []uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportIdentifiedRiskDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var risks []dtos.RiskReportIdentifiedRiskDTO

	query := db.Table("risk_register").
		Select(`
			risk_register.id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			risk_register.reported_at,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
			inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
			inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
			inherent_likelihood.title_en AS inherent_likelihood_title_en,
			inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
			inherent_impact.title_en AS inherent_impact_title_en,
			inherent_impact.title_ar AS inherent_impact_title_ar,
			residual_risk_rating.title_en AS residual_risk_rating_title_en,
			residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
			residual_risk_rating.color_code AS residual_risk_rating_color_code,
			residual_likelihood.title_en AS residual_likelihood_title_en,
			residual_likelihood.title_ar AS residual_likelihood_title_ar,
			residual_impact.title_en AS residual_impact_title_en,
			residual_impact.title_ar AS residual_impact_title_ar,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id").
		Where("risk_register.reported_at BETWEEN ? AND ?", startDate, endDate).
		// Where("risk_register.responsible_unit_id IN (?)", organizationUnitIDs).
		Order("risk_register.reported_at DESC")

	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	return risks, nil
}

func (repo *RiskReportRepository) GetControlTitlesByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var controlTitles []string
	err := db.Table("risk_controls").
		Select("controls.title_en").
		Joins("LEFT JOIN controls ON controls.id = risk_controls.control_id").
		Where("risk_controls.risk_register_id = ?", riskID).
		Pluck("controls.title_en", &controlTitles).Error

	if err != nil {
		return nil, err
	}
	return controlTitles, nil
}

func (repo *RiskReportRepository) CreateSection(section *models.RiskReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&section).Error
}

func (repo *RiskReportRepository) GetActionPlansByRiskID(riskRegisterID uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportActionPlanDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var actionPlans []dtos.RiskReportActionPlanDTO

	query := db.Table("risk_action_plans").
		Select(`
			risk_action_plans.id,
			risk_action_plans.risk_register_id,
			risk_action_plans.title AS action_plan_title,
			risk_action_plans.description AS action_plan_description,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			(CONCAT(users.first_name_en, ' ', users.last_name_en)) AS responsible_user_name_en,
			(CONCAT(users.first_name_ar, ' ', users.last_name_ar)) AS responsible_user_name_ar,
			risk_action_plans.target_date,
			action_plan_statuses.title_en AS action_plan_status_title_en,
			action_plan_statuses.title_ar AS action_plan_status_title_ar,
			action_plan_statuses.color_code AS action_plan_status_color_code
		`).
		Joins("LEFT JOIN severity_levels ON severity_levels.id = risk_action_plans.severity_level_id").
		Joins("LEFT JOIN users ON users.id = risk_action_plans.responsible_user_id").
		Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = risk_action_plans.action_plan_status_id").
		Where("risk_action_plans.risk_register_id = ?", riskRegisterID)

	if err := query.Find(&actionPlans).Error; err != nil {
		return nil, err
	}

	return actionPlans, nil
}

func (repo *RiskReportRepository) GetActionPlanControlsByActionPlanID(actionPlanID uint, tx ...utility.DatabaseTransaction) (string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var controlTitles []string
	err := db.Table("risk_action_plan_controls").
		Select("COALESCE(STRING_AGG(controls.title_en, ', '), '') AS control_titles").
		Joins("LEFT JOIN controls ON controls.id = risk_action_plan_controls.control_id").
		Where("risk_action_plan_controls.risk_action_plan_id = ?", actionPlanID).
		Group("risk_action_plan_controls.risk_action_plan_id").
		Find(&controlTitles).Error

	if err != nil {
		return "", err
	}

	if len(controlTitles) > 0 {
		return controlTitles[0], nil
	}

	return "", nil
}

func (repo *RiskReportRepository) CreateRiskActionPlan(actionPlan *models.RiskReportActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(actionPlan).Error
}

func (repo *RiskReportRepository) GetKRIsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportKRIDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var kris []dtos.RiskReportKRIDTO

	query := db.Table("risk_kris").
		Select(`
			kri_register.title,
			kri_register.description,
			kri_register.low_threshold,
			kri_register.medium_threshold,
			kri_register.current_value,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar
		`).
		Joins("LEFT JOIN kri_register ON kri_register.id = risk_kris.kri_register_id").
		Joins("LEFT JOIN units ON units.id = kri_register.unit_id").
		Where("risk_kris.risk_register_id = ?", riskID)

	if err := query.Find(&kris).Error; err != nil {
		return nil, err
	}

	return kris, nil
}

func (repo *RiskReportRepository) GetDocumentsByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]dtos.RiskReportDocumentDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documents []dtos.RiskReportDocumentDTO

	query := db.Table("risk_documents").
		Select(`
			risk_documents.url,
			risk_documents.name,
			risk_documents.size,
			risk_documents.token,
			risk_documents.original_name,
			risk_documents.type
		`).
		Where("risk_documents.risk_register_id = ?", riskID)

	if err := query.Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *RiskReportRepository) CreateRiskKRI(riskKRI *models.RiskReportKRI, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(riskKRI).Error
}

func (repo *RiskReportRepository) CreateRiskDocument(riskDocument *models.RiskReportDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(riskDocument).Error
}

func (repo *RiskReportRepository) CreateRiskReportOrganizationUnit(orgUnit *models.RiskReportOrganizationUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(orgUnit).Error
}

func (repo *RiskReportRepository) UpdateSection(section *models.RiskReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (repo *RiskReportRepository) DeleteRiskReportOrganizationUnits(riskReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_report_id = ?", riskReportID).Delete(&models.RiskReportOrganizationUnit{}).Error
}

func (repo *RiskReportRepository) DeleteRiskKRIs(riskReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_report_id = ?", riskReportID).Delete(&models.RiskReportKRI{}).Error
}

func (repo *RiskReportRepository) DeleteRiskActionPlans(riskReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_report_id = ?", riskReportID).Delete(&models.RiskReportActionPlan{}).Error
}

func (repo *RiskReportRepository) DeleteRiskIdentifications(riskReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_report_id = ?", riskReportID).Delete(&models.RiskReportIdentification{}).Error
}

func (repo *RiskReportRepository) DeleteRiskReportSections(riskReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("risk_report_id = ?", riskReportID).Delete(&models.RiskReportSection{}).Error
}

func (repo *RiskReportRepository) DeleteRiskReport(riskReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", riskReportID).Delete(&models.RiskReport{}).Error
}

func (repo *RiskReportRepository) DeleteRiskReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.RiskReportSection{}).Error
}

func (repo *RiskReportRepository) ReorderRiskReportSections(riskReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Shift order numbers of remaining sections
	err := db.Exec(`
		UPDATE risk_report_sections
		SET order_num = order_num - 1
		WHERE risk_report_id = ? AND order_num > ?`,
		riskReportID, deletedOrderNum).Error

	return err
}

func (repo *RiskReportRepository) CreateRiskDocuments(documents []*models.RiskReportDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if len(documents) == 0 {
		return nil // No documents to insert
	}

	if err := db.Create(&documents).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskReportRepository) DeleteRiskReportDocument(riskReportID uint, documentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ? AND risk_report_id = ?", documentID, riskReportID).Delete(&models.RiskReportDocument{}).Error
}

func (repo *RiskReportRepository) GetRiskReportDocumentByID(riskReportID uint, documentID uint, tx ...utility.DatabaseTransaction) (*models.RiskReportDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var riskReportDocument models.RiskReportDocument

	err := db.
		Where("id = ? AND risk_report_id = ?", documentID, riskReportID).
		First(&riskReportDocument).Error

	if err != nil {
		return nil, err
	}

	return &riskReportDocument, nil
}
