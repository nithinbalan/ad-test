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

	// Methods for Risk Register PDF Export (following incident pattern)
	GetAllRiskAssets(riskID uint) ([]*commonDTOs.AssociatedAssetListDTO, error)
	GetAllRiskControls(riskID uint) ([]*commonDTOs.AssociatedControlListDTO, error)
	GetAllRiskDocuments(riskID uint) ([]*commonDTOs.AssociatedDocumentListDTO, error)
	GetAllRiskProcesses(riskID uint) ([]*commonDTOs.AssociatedProcessListDTO, error)
	GetAllRiskKPIs(riskID uint) ([]*commonDTOs.AssociatedKPIRegisterListDTO, error)
	GetAllRiskStrategicObjectives(riskID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error)
	GetAllRiskComplianceRequirements(riskID uint) ([]*commonDTOs.AssociatedClauseListDTO, error)
	GetAllRiskIncidents(riskID uint) ([]*commonDTOs.AssociatedIncidentListDTO, error)
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
			STRING_AGG(organization_units.title_ar, ', ') AS organization_unit_titles_ar,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN risk_report_organization_units ON risk_report_organization_units.risk_report_id = risk_reports.id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_report_organization_units.organization_unit_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = risk_reports.created_by_id").
		Group(`
			risk_reports.id, risk_reports.title, risk_reports.start_date, risk_reports.end_date, risk_reports.created_at,
			created_user.id, created_user.first_name_en, created_user.first_name_ar, 
			created_user.last_name_en, created_user.last_name_ar, created_user.email,
			created_user.initial, created_user.color, created_user.profile_image_path
		`)

	if params.Keyword != "" {
		query = query.Where(`
			(risk_reports.title LIKE ? OR 
			organization_units.title_en LIKE ? OR 
			organization_units.title_ar LIKE ? OR 
			created_user.first_name_en LIKE ? OR 
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
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

// Methods for Risk Register PDF Export (following incident pattern)

func (repo *RiskReportRepository) GetAllRiskAssets(riskID uint) ([]*commonDTOs.AssociatedAssetListDTO, error) {
	var riskAssets []*commonDTOs.AssociatedAssetListDTO

	query := repo.db.Table("risk_assets").
		Select(`
			risk_assets.id,
			risk_assets.asset_register_id,
			asset_register.reference_number,
			asset_register.title_en,
			asset_register.title_ar,
			asset_register.location,
			asset_register.description,
			CONVERT(VARCHAR, asset_register.acquisition_date, 23) AS acquisition_date,
			CONVERT(VARCHAR, asset_register.assignment_date, 23) AS assignment_date,
			CONVERT(VARCHAR, asset_register.warranty_expiry_date, 23) AS warranty_expiry_date,
			asset_register.manufacturer,
			asset_register.model,
			asset_register.serial_number,
			asset_categories.title_en AS asset_category_title_en,
			asset_categories.title_ar AS asset_category_title_ar,
			asset_sub_categories.title_en AS asset_sub_category_title_en,
			asset_sub_categories.title_ar AS asset_sub_category_title_ar,
			asset_statuses.title_en AS asset_status_title_en,
			asset_statuses.title_ar AS asset_status_title_ar,
			asset_statuses.color_code AS asset_status_color_code,
			asset_register.total_score,
			asset_criticality_matrix_rating_scales.title_en AS asset_criticality_matrix_rating_scale_title_en,
			asset_criticality_matrix_rating_scales.title_ar AS asset_criticality_matrix_rating_scale_title_ar,
			asset_criticality_matrix_rating_scales.color_code AS asset_criticality_matrix_rating_scale_color_code,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN asset_register ON asset_register.id = risk_assets.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN asset_criticality_matrix_rating_scales ON asset_criticality_matrix_rating_scales.id = asset_register.asset_criticality_matrix_rating_scale_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id").
		Where("risk_assets.risk_register_id = ?", riskID)

	err := query.Scan(&riskAssets).Error
	if err != nil {
		return nil, err
	}

	return riskAssets, nil
}

func (repo *RiskReportRepository) GetAllRiskControls(riskID uint) ([]*commonDTOs.AssociatedControlListDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO

	query := repo.db.Table("risk_controls").
		Select(`
		DISTINCT controls.id AS control_id,
		controls.*,
		CONVERT(VARCHAR, controls.last_test_date, 23) AS control_last_test_date,
		control_types.title_en AS control_type_title_en,
		control_types.title_ar AS control_type_title_ar,
		control_categories.title_en AS control_category_title_en,
		control_categories.title_ar AS control_category_title_ar,
		organization_units.title_en AS owner_unit_title_en,
		organization_units.title_ar AS owner_unit_title_ar,
		owner.id AS owner_id,
		owner.first_name_en AS owner_first_name_en,
		owner.first_name_ar AS owner_first_name_ar,
		owner.last_name_en AS owner_last_name_en,
		owner.last_name_ar AS owner_last_name_ar,
		owner.initial AS owner_initial,
		owner.color AS owner_color,
		owner.email AS owner_email,
		owner.profile_image_path AS owner_profile_image_path,
		review_frequencies.title_en AS review_frequency_title_en,
		review_frequencies.title_ar AS review_frequency_title_ar,
		test_responsible_user.id AS test_responsible_user_id,
		test_responsible_user.first_name_en AS test_responsible_user_first_name_en,
		test_responsible_user.first_name_ar AS test_responsible_user_first_name_ar,
		test_responsible_user.last_name_en AS test_responsible_user_last_name_en,
		test_responsible_user.last_name_ar AS test_responsible_user_last_name_ar,
		test_responsible_user.initial AS test_responsible_user_initial,
		test_responsible_user.color AS test_responsible_user_color,
		test_responsible_user.email AS test_responsible_user_email,
		test_responsible_user.profile_image_path AS test_responsible_user_profile_image_path,
		test_frequencies.title_en AS test_frequency_title_en,
		test_frequencies.title_ar AS test_frequency_title_ar,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		control_implementation_statuses.title_en AS control_implementation_status_title_en,
		control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
		control_implementation_statuses.color_code AS control_implementation_status_color_code,
		control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
		control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
		control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		created_user.id AS created_user_id,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		created_user.initial AS created_user_initial,
		created_user.color AS created_user_color,
		created_user.email AS created_user_email,
		created_user.profile_image_path AS created_user_profile_image_path
	`).
		Joins(`
		LEFT JOIN controls ON controls.id = risk_controls.control_id
		LEFT JOIN control_types ON controls.control_type_id = control_types.id
		LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
		LEFT JOIN organization_units ON controls.responsible_unit_id = organization_units.id
		LEFT JOIN users AS owner ON controls.responsible_user_id = owner.id
		LEFT JOIN review_frequencies ON controls.review_frequency_id = review_frequencies.id
		LEFT JOIN users AS test_responsible_user ON controls.test_responsible_user_id = test_responsible_user.id
		LEFT JOIN review_frequencies AS test_frequencies ON controls.test_frequency_id = test_frequencies.id
		LEFT JOIN lifecycle_statuses ON controls.lifecycle_status_id = lifecycle_statuses.id
		LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
		LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
		LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = controls.risk_matrix_risk_rating_id
		LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = controls.risk_matrix_risk_rating_value_id
		LEFT JOIN users AS created_user ON created_user.id = controls.created_by_id
	`).Where("risk_controls.risk_register_id = ?", riskID)

	err := query.Scan(&controlList).Error
	if err != nil {
		return nil, err
	}

	return controlList, nil
}

func (repo *RiskReportRepository) GetAllRiskDocuments(riskID uint) ([]*commonDTOs.AssociatedDocumentListDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO

	query := repo.db.Table("risk_documents").
		Select(`
			risk_documents.id AS id,
			risk_documents.id AS document_id,
			'' AS reference_number,
			0 AS document_type_id,
			risk_documents.original_name AS title_en,
			risk_documents.original_name AS title_ar,
			NULL AS last_checkout_at,
			NULL AS last_checkin_at,
			0 AS is_locked,
			NULL AS description_en,
			NULL AS description_ar,
			0 AS responsible_user_id,
			risk_documents.created_at,
			risk_documents.updated_at,
			NULL AS published_at,
			'' AS document_category_title_en,
			'' AS document_category_title_ar,
			'' AS document_type_title_en,
			'' AS document_type_title_ar,
			'' AS document_access_type_title_en,
			'' AS document_type_access_title_ar,
			'' AS responsible_unit_title_en,
			'' AS responsible_unit_title_ar,
			0 AS responsible_user_id,
			'' AS responsible_user_first_name_en,
			'' AS responsible_user_first_name_ar,
			'' AS responsible_user_last_name_en,
			'' AS responsible_user_last_name_ar,
			'' AS responsible_user_initial,
			'' AS responsible_user_color,
			'' AS responsible_user_email,
			'' AS responsible_user_profile_image_path,
			risk_documents.url AS document_url,
			risk_documents.name AS document_name,
			risk_documents.size AS document_size,
			risk_documents.token AS document_token,
			risk_documents.original_name AS document_original_name,
			risk_documents.type AS document_type,
			1.0 AS document_version,
			'' AS document_status_title_en,
			'' AS document_status_title_ar,
			'' AS document_status_color_code,
			'' AS document_status_type,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			'' AS published_user_first_name_en,
			'' AS published_user_first_name_ar,
			'' AS published_user_last_name_en,
			'' AS published_user_last_name_ar,
			'' AS published_user_initial,
			'' AS published_user_color,
			'' AS published_user_email,
			'' AS published_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN users AS created_user ON created_user.id = risk_documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = risk_documents.updated_by_id
		`).Where("risk_documents.risk_register_id = ?", riskID)

	err := query.Scan(&documentList).Error
	if err != nil {
		return nil, err
	}

	return documentList, nil
}

func (repo *RiskReportRepository) GetAllRiskProcesses(riskID uint) ([]*commonDTOs.AssociatedProcessListDTO, error) {
	var processList []*commonDTOs.AssociatedProcessListDTO

	query := repo.db.Table("risk_processes").
		Select(`
		risk_processes.id,
		risk_processes.process_id,
		processes.reference_number,
		processes.title_en,
		processes.title_ar,
		processes.description,
		processes.marked_as_auditable_item,
		process_groups.title_en AS process_group_title_en,
		process_groups.title_ar AS process_group_title_ar,
		processes.rto as rto,
		processes.rto_unit_id as rto_unit_id,
		rto_time_unit.title_en AS rto_unit_title_en,
		rto_time_unit.title_ar AS rto_unit_title_ar,
		processes.mtpd as mtpd,
		processes.mtpd_unit_id as mtpd_unit_id,
		mtpd_time_unit.title_en AS mtpd_unit_title_en,
		mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
		processes.rpo as rpo,
		processes.rpo_unit_id as rpo_unit_id,
		rpo_time_unit.title_en AS rpo_unit_title_en,
		rpo_time_unit.title_ar AS rpo_unit_title_ar,
		processes.is_bia_performed,
		processes.last_assessed_at,
		last_assessed_user.id AS last_assessed_user_id,
		last_assessed_user.first_name_en AS last_assessed_user_first_name_en,
		last_assessed_user.first_name_ar AS last_assessed_user_first_name_ar,
		last_assessed_user.last_name_en AS last_assessed_user_last_name_en,
		last_assessed_user.last_name_ar AS last_assessed_user_last_name_ar,
		last_assessed_user.initial AS last_assessed_user_initial,
		last_assessed_user.color AS last_assessed_user_color,
		last_assessed_user.email AS last_assessed_user_email,
		last_assessed_user.profile_image_path AS last_assessed_user_profile_image_path,
		organization_units.title_en AS owner_unit_title_en,
		organization_units.title_ar AS owner_unit_title_ar,
		owner.id AS owner_id,
		owner.first_name_en AS owner_first_name_en,
		owner.first_name_ar AS owner_first_name_ar,
		owner.last_name_en AS owner_last_name_en,
		owner.last_name_ar AS owner_last_name_ar,
		owner.initial AS owner_initial,
		owner.color AS owner_color,
		owner.email AS owner_email,
		owner.profile_image_path AS owner_profile_image_path,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
		process_criticality_levels.id AS process_criticality_level_id,
		process_criticality_levels.title_en AS process_criticality_level_title_en,
		process_criticality_levels.title_ar AS process_criticality_level_title_ar,
		created_user.id AS created_user_id,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		created_user.initial AS created_user_initial,
		created_user.color AS created_user_color,
		created_user.email AS created_user_email,
		created_user.profile_image_path AS created_user_profile_image_path,
		updated_user.id AS updated_user_id,
		updated_user.first_name_en AS updated_user_first_name_en,
		updated_user.first_name_ar AS updated_user_first_name_ar,
		updated_user.last_name_en AS updated_user_last_name_en,
		updated_user.last_name_ar AS updated_user_last_name_ar,
		updated_user.initial AS updated_user_initial,
		updated_user.color AS updated_user_color,
		updated_user.email AS updated_user_email,
		updated_user.profile_image_path AS updated_user_profile_image_path,
		processes.created_at,
		processes.updated_at,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		risk_levels.title_en AS risk_level_title_en,
        risk_levels.title_ar AS risk_level_title_ar,
		risk_levels.color_code AS risk_level_color_code,
		risk_levels.type AS risk_level_type
		`).
		Joins(`
			LEFT JOIN processes ON processes.id = risk_processes.process_id
		    LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id
			LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = processes.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = processes.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = processes.rpo_unit_id
			LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id
		`).Where("risk_processes.risk_register_id = ?", riskID)

	err := query.Scan(&processList).Error
	if err != nil {
		return nil, err
	}

	return processList, nil
}

func (repo *RiskReportRepository) GetAllRiskKPIs(riskID uint) ([]*commonDTOs.AssociatedKPIRegisterListDTO, error) {
	var kpiRegister []*commonDTOs.AssociatedKPIRegisterListDTO

	query := repo.db.Table("risk_kris").
		Select(`
            kpi_register.*,
			risk_kris.kri_register_id,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			CONVERT(VARCHAR, kpi_register.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_register.end_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			calculation_types.title_en AS calculation_type_title_en,
			calculation_types.title_ar AS calculation_type_title_ar,
			calculation_types.type AS calculation_type_type,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path  
        `).
		Joins(`
		    LEFT JOIN kpi_register ON kpi_register.id = risk_kris.kri_register_id
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
        `).Where("risk_kris.risk_register_id = ?", riskID)

	err := query.Scan(&kpiRegister).Error
	if err != nil {
		return nil, err
	}

	return kpiRegister, nil
}

func (repo *RiskReportRepository) GetAllRiskStrategicObjectives(riskID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error) {
	var objectiveList []*commonDTOs.AssociatedStrategicObjectiveListDTO

	query := repo.db.Table("risk_strategic_objectives").
		Select(`
			strategic_objectives.*,
			risk_strategic_objectives.strategic_objective_id,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
        `).
		Joins(`
			LEFT JOIN strategic_objectives ON strategic_objectives.id = risk_strategic_objectives.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
			LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id
			LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id
        `).Where("risk_strategic_objectives.risk_register_id = ?", riskID)

	err := query.Scan(&objectiveList).Error
	if err != nil {
		return nil, err
	}

	return objectiveList, nil
}

func (repo *RiskReportRepository) GetAllRiskComplianceRequirements(riskID uint) ([]*commonDTOs.AssociatedClauseListDTO, error) {
	var clauseList []commonDTOs.AssociatedClauseRowDTO

	query := repo.db.Table("risk_compliance_requirements").
		Select(`
			compliance_clauses.title_en AS clause_title_en,
			compliance_clauses.title_ar AS clause_title_ar,
			compliance_clauses.clause_no,
			compliance_clauses.compliance_register_id,
			compliance_clauses.id AS compliance_clause_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar
		`).
		Joins(`
			LEFT JOIN compliance_clauses ON risk_compliance_requirements.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_register ON risk_compliance_requirements.compliance_register_id = compliance_register.id
		`).
		Where("risk_compliance_requirements.risk_register_id = ?", riskID)

	if err := query.Find(&clauseList).Error; err != nil {
		return nil, err
	}

	clausesByComplianceRegister := make(map[uint]*commonDTOs.AssociatedClauseListDTO)
	for _, row := range clauseList {
		if _, exists := clausesByComplianceRegister[row.ComplianceRegisterID]; !exists {
			clausesByComplianceRegister[row.ComplianceRegisterID] = &commonDTOs.AssociatedClauseListDTO{
				ComplianceRegisterID:      row.ComplianceRegisterID,
				ComplianceRegisterTitleEn: row.ComplianceRegisterTitleEn,
				ComplianceRegisterTitleAr: row.ComplianceRegisterTitleAr,
				Clauses:                   []commonDTOs.AssociatedClauseDTO{},
			}
		}
		clausesByComplianceRegister[row.ComplianceRegisterID].Clauses = append(
			clausesByComplianceRegister[row.ComplianceRegisterID].Clauses,
			commonDTOs.AssociatedClauseDTO{
				ComplianceClauseID: row.ComplianceClauseID,
				ClauseNo:           row.ClauseNo,
				ClauseTitleEn:      row.ClauseTitleEn,
				ClauseTitleAr:      row.ClauseTitleAr,
			},
		)
	}

	var groupedClauses []*commonDTOs.AssociatedClauseListDTO
	for _, v := range clausesByComplianceRegister {
		groupedClauses = append(groupedClauses, v)
	}

	if groupedClauses == nil {
		groupedClauses = []*commonDTOs.AssociatedClauseListDTO{}
	}

	return groupedClauses, nil
}

func (repo *RiskReportRepository) GetAllRiskIncidents(riskID uint) ([]*commonDTOs.AssociatedIncidentListDTO, error) {
	var incidentList []*commonDTOs.AssociatedIncidentListDTO

	query := repo.db.Table("incident_risks").
		Select(`
			incident_risks.id,
			incident_risks.incident_register_id,
			incident_register.reference_number,
			incident_register.title,
			incident_register.description,
			incident_register.incident_at,
			incident_register.location,
			incident_register.reported_at,
			incident_register.quick_action,
			incident_register.incident_category_id,
			incident_register.incident_type_id,
			incident_register.severity_level_id,
			incident_register.incident_status_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			incident_categories.title_en AS incident_category_title_en,
			incident_categories.title_ar AS incident_category_title_ar,
			incident_types.title_en AS incident_type_title_en,
			incident_types.title_ar AS incident_type_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			incident_statuses.title_en AS incident_status_title_en,
			incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code,
			severity_levels.color_code AS severity_level_color_code,
			reported_user.id AS reported_user_id,
			reported_user.first_name_en AS reported_user_first_name_en,
			reported_user.first_name_ar AS reported_user_first_name_ar,
			reported_user.last_name_en AS reported_user_last_name_en,
			reported_user.last_name_ar AS reported_user_last_name_ar,
			reported_user.initial AS reported_user_initial,
			reported_user.color AS reported_user_color,
			reported_user.email AS reported_user_email,
			reported_user.profile_image_path AS reported_user_profile_image_path,
			CONVERT(VARCHAR, investigation.initiated_at, 23) AS investigation_initiated_at,
			investigation.description AS investigation_description,
			CONVERT(VARCHAR, incident_register.created_at, 23) AS created_at,
			CONVERT(VARCHAR, incident_register.updated_at, 23) AS updated_at,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN incident_register ON incident_register.id = incident_risks.incident_register_id
			LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
			LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
			LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
			LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
			LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id
			LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
			LEFT JOIN users AS created_user ON created_user.id = incident_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = incident_register.updated_by_id
			LEFT JOIN incident_investigations AS investigation ON investigation.incident_register_id = incident_register.id
		`).Where("incident_risks.risk_register_id = ?", riskID)

	err := query.Scan(&incidentList).Error
	if err != nil {
		return nil, err
	}

	return incidentList, nil
}
