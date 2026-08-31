package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type IncidentReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *incidentDTOs.IncidentReportListQueryParamsDTO) (*incidentDTOs.IncidentReportListResponseDTO, error)
	GetFullReportDetailsByID(id uint) (*models.IncidentRegister, error)
	GetIncidentReportDetailsByID(id uint) (*models.IncidentReport, error)
	GetAllIncidentKPI(incidentID uint) ([]*commonDTOs.AssociatedKPIRegisterListDTO, error)
	GetAllIncidentStrategicObjectives(incidentID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error)
	GetAllIncidentActionPlan(incidentID uint) ([]*incidentDTOs.IncidentActionPlanListDTO, error)
	GetAllIncidentAssets(incidentID uint) ([]*commonDTOs.AssociatedAssetListDTO, error)
	GetAllIncidentAssociatedDocuments(incidentID uint) ([]*commonDTOs.AssociatedDocumentListDTO, error)
	GetAllIncidentControls(incidentID uint) ([]*commonDTOs.AssociatedControlListDTO, error)
	GetAllIncidentRisks(incidentID uint) ([]*commonDTOs.AssociatedRiskListDTO, error)
	GetAllIncidentProcesses(incidentID uint) ([]*commonDTOs.AssociatedProcessListDTO, error)
	GetAllIncidentComplianceRequirements(incidentID uint) ([]*commonDTOs.AssociatedClauseListDTO, error)
	GetIncidentsByIDs(ids []uint, tx utility.DatabaseTransaction) ([]incidentDTOs.IncidentReportIncidentDTO, error)
	GetReportIncidentsByDateRange(start time.Time, end time.Time, tx utility.DatabaseTransaction) ([]incidentDTOs.IncidentReportIncidentDTO, error)
	GetDocumentsByIncidentID(incidentID uint, tx utility.DatabaseTransaction) ([]incidentDTOs.IncidentReportDocumentDTO, error)
	GetEvidencesByIncidentID(incidentID uint, tx utility.DatabaseTransaction) ([]incidentDTOs.IncidentReportEvidenceDTO, error)
	Create(report *models.IncidentReport, tx utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.IncidentReportSection, tx utility.DatabaseTransaction) (uint, error)
	CreateReportIncident(incident *models.IncidentReportIncident, tx utility.DatabaseTransaction) error
	CreateIncidentDocuments(docs []*models.IncidentReportDocument, tx utility.DatabaseTransaction) error
	CreateIncidentEvidences(evidences []*models.IncidentReportEvidence, tx utility.DatabaseTransaction) error
	CreateIncidentRootCauseAnalyses(rcas []*models.IncidentReportRootCauseAnalysis, tx utility.DatabaseTransaction) error
	CreateIncidentReportActionPlans(plans []*models.IncidentReportActionPlan, tx utility.DatabaseTransaction) error

	GetIncidentReportByID(incidentReportID uint, tx ...utility.DatabaseTransaction) (*models.IncidentReport, error)
	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.IncidentReportSection, error)
	UpdateSection(section *models.IncidentReportSection, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	ReorderIncidentReportSections(incidentReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReportDocuments(incidentReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReportRootCauseAnalyses(incidentReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReportActionPlans(incidentReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReportIncidents(incidentReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReportSections(incidentReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteIncidentReport(incidentReportID uint, tx ...utility.DatabaseTransaction) error
}

type IncidentReportRepository struct {
	db *gorm.DB
}

func NewIncidentReportRepository(db *gorm.DB) *IncidentReportRepository {
	return &IncidentReportRepository{
		db: db,
	}
}

func (repo *IncidentReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentReportRepository) GetAll(params *incidentDTOs.IncidentReportListQueryParamsDTO) (*incidentDTOs.IncidentReportListResponseDTO, error) {
	var reports []*incidentDTOs.IncidentReportListDTO
	var totalCount int64

	query := repo.db.Table("incident_reports").
		Select(`
			incident_reports.id,
			incident_reports.title,
			CONVERT(VARCHAR, incident_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, incident_reports.end_date, 23) AS end_date,
			incident_reports.created_at,
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
		Joins("LEFT JOIN users AS created_user ON created_user.id = incident_reports.created_by_id")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(incident_reports.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword)
	}

	// Count before applying ORDER BY to avoid SQL Server OFFSET/FETCH issues
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Ensure ORDER BY is always present for SQL Server OFFSET/FETCH
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "incident_reports.id"
	}
	sort := params.Sort
	if sort == "" {
		sort = "DESC"
	}
	query = query.Order(sortBy + " " + sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &incidentDTOs.IncidentReportListResponseDTO{
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

func (repo *IncidentReportRepository) GetFullReportDetailsByID(id uint) (*models.IncidentRegister, error) {
	db := repo.db

	var incident *models.IncidentRegister

	result := db.
		Preload("IncidentCategory").
		Preload("IncidentSubCategory").
		Preload("IncidentType").
		Preload("IncidentStatus").
		Preload("ReportedByUser").
		Preload("SeverityLevel").
		Preload("ResponsibleUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("BusinessOwner").
		Preload("IncidentResponseTeam").
		Preload("IncidentResponseTeam.TeamLeadUser").
		Preload("IncidentResponseTeam.IncidentResponseTeamMembers").
		Preload("IncidentResponseTeam.IncidentResponseTeamMembers.User").
		Preload("ResponseTeamLead").
		Preload("ResponseTeamMembers").
		Preload("ResponseTeamMembers.User").
		Preload("IncidentScope").
		Preload("IncidentSensitivityClassification").
		Preload("ConfidentialitySeverityLevel").
		Preload("IntegritySeverityLevel").
		Preload("AvailabilitySeverityLevel").
		Preload("AcknowledgedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Investigation").
		Preload("Investigation.TeamLead").
		Preload("Investigation.Users").
		Preload("Investigation.Users.User").
		Preload("Investigation.Users.User.Designation").
		Preload("Investigation.Documents").
		Preload("Investigation.Findings").
		Preload("Documents").
		Preload("Timelines", func(db *gorm.DB) *gorm.DB {
			return db.Order("performed_at ASC")
		}).
		Preload("Timelines.PerformedBy").
		Preload("Timelines.TimelineStatus").
		Preload("Timelines.IncidentTimelineType").
		Preload("CommunicationLogs").
		Preload("IncidentClosure").
		Preload("IncidentPostReview").
		Preload("IncidentClosure.Lessons").
		Preload("IncidentClosure.SubmittedByUser").
		Preload("IncidentClosure.ApprovedByUser").
		Preload("IncidentClosure.WorkflowStatus").
		Preload("IncidentPostReview").
		Preload("IncidentPostReview.Users").
		Preload("IncidentPostReview.Users.User").
		Preload("IncidentPostReview.Recommendations").
		Preload("IncidentPostReview.Lessons").
		Preload("IncidentPostReview.IncidentPostReviewEffectiveness").
		Preload("IncidentRootCauseAnalyses").
		First(&incident, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil
}

func (repo *IncidentReportRepository) GetAllIncidentKPI(incidentID uint) ([]*commonDTOs.AssociatedKPIRegisterListDTO, error) {
	var kpiRegister []*commonDTOs.AssociatedKPIRegisterListDTO

	query := repo.db.Table("incident_kpis").
		Select(`
            kpi_register.*,
			incident_kpis.kpi_register_id,
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
		    LEFT JOIN kpi_register ON kpi_register.id = incident_kpis.kpi_register_id
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
        `).Where("incident_kpis.incident_register_id = ?", incidentID)

	err := query.Scan(&kpiRegister).Error
	if err != nil {
		return nil, err
	}

	return kpiRegister, nil

}

func (repo *IncidentReportRepository) GetAllIncidentStrategicObjectives(incidentID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error) {
	var objectiveList []*commonDTOs.AssociatedStrategicObjectiveListDTO

	query := repo.db.Table("incident_strategic_objectives").
		Select(`
			strategic_objectives.*,
			incident_strategic_objectives.strategic_objective_id,
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
			LEFT JOIN strategic_objectives ON strategic_objectives.id = incident_strategic_objectives.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
			LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id
			LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id
        `).Where("incident_strategic_objectives.incident_register_id = ?", incidentID)

	err := query.Scan(&objectiveList).Error
	if err != nil {
		return nil, err
	}

	return objectiveList, nil
}

func (repo *IncidentReportRepository) GetAllIncidentActionPlan(incidentID uint) ([]*incidentDTOs.IncidentActionPlanListDTO, error) {

	var incidentActionPlans []*incidentDTOs.IncidentActionPlanListDTO

	query := repo.db.Table("incident_action_plans").
		Select(`
        incident_action_plans.*,
		incident_register.reference_number AS incident_reference_number,
		incident_register.title AS incident_title,
		CONVERT(VARCHAR, incident_action_plans.target_date, 23) AS target_date,
        action_plan_statuses.title_en AS action_plan_status_title_en,
        action_plan_statuses.title_ar AS action_plan_status_title_ar,
		action_plan_statuses.color_code AS action_plan_status_color_code,
        severity_levels.title_en AS severity_level_title_en,
        severity_levels.title_ar AS severity_level_title_ar,
        severity_levels.color_code AS severity_level_color_code,
		responsible_user.id AS responsible_user_id,
        responsible_user.first_name_en AS responsible_user_first_name_en,
        responsible_user.first_name_ar AS responsible_user_first_name_ar,
        responsible_user.last_name_en AS responsible_user_last_name_en,
        responsible_user.last_name_ar AS responsible_user_last_name_ar,
		responsible_user.initial AS responsible_user_initial,
        responsible_user.color AS responsible_user_color,
        responsible_user.email AS responsible_user_email,
        responsible_user.profile_image_path AS responsible_user_profile_image_path, 
		review_user.id AS review_user_id,
    	review_user.first_name_en AS review_user_first_name_en,
		review_user.first_name_ar AS review_user_first_name_ar,
		review_user.last_name_en AS review_user_last_name_en,
		review_user.last_name_ar AS review_user_last_name_ar,
		review_user.initial AS review_user_initial,
        review_user.color AS review_user_color,
        review_user.email AS review_user_email,
        review_user.profile_image_path AS review_user_profile_image_path,
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
		LEFT JOIN incident_register ON incident_register.id = incident_action_plans.incident_register_id
        LEFT JOIN action_plan_statuses ON action_plan_statuses.id = incident_action_plans.action_plan_status_id
        LEFT JOIN severity_levels ON severity_levels.id = incident_action_plans.severity_level_id
        LEFT JOIN users AS responsible_user ON responsible_user.id = incident_action_plans.responsible_user_id
		LEFT JOIN users AS review_user ON review_user.id = incident_action_plans.review_user_id
        LEFT JOIN users AS created_user ON created_user.id = incident_action_plans.created_by_id
        LEFT JOIN users AS updated_user ON updated_user.id = incident_action_plans.updated_by_id
    `).Where("incident_action_plans.incident_register_id = ?", incidentID)
	err := query.Scan(&incidentActionPlans).Error
	if err != nil {
		return nil, err
	}

	return incidentActionPlans, nil
}

func (repo *IncidentReportRepository) GetAllIncidentAssets(incidentID uint) ([]*commonDTOs.AssociatedAssetListDTO, error) {
	var incidentAssets []*commonDTOs.AssociatedAssetListDTO

	query := repo.db.Table("incident_assets").
		Select(`
			incident_assets.id,
			incident_assets.asset_register_id,
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
		Joins("LEFT JOIN asset_register ON asset_register.id = incident_assets.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN asset_criticality_matrix_rating_scales ON asset_criticality_matrix_rating_scales.id = asset_register.asset_criticality_matrix_rating_scale_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id").
		Where("incident_assets.incident_register_id = ?", incidentID)

	err := query.Scan(&incidentAssets).Error
	if err != nil {
		return nil, err
	}

	return incidentAssets, nil
}

func (repo *IncidentReportRepository) GetAllIncidentAssociatedDocuments(incidentID uint) ([]*commonDTOs.AssociatedDocumentListDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO

	query := repo.db.Table("incident_associated_documents").
		Select(`
			incident_associated_documents.id AS id,
			incident_associated_documents.document_id,
			documents.reference_number AS reference_number,
			documents.document_type_id AS document_type_id,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			documents.last_checkout_at AS last_checkout_at,
			documents.last_checkin_at AS last_checkin_at,
			documents.is_locked AS is_locked,
			documents.description_en AS description_en,
			documents.description_ar AS description_ar,
			documents.responsible_user_id,
			documents.created_at,
			documents.updated_at,
			documents.published_at,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_access_types.title_en AS document_access_type_title_en,
			document_access_types.title_ar AS document_type_access_title_ar,
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
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name,
			document_versions.type AS document_type,
			document_versions.version AS document_version,
			document_statuses.title_en AS document_status_title_en,
			document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code,
			document_statuses.type AS document_status_type,
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
			published_user.first_name_en AS published_user_first_name_en,
			published_user.first_name_ar AS published_user_first_name_ar,
			published_user.last_name_en AS published_user_last_name_en,
			published_user.last_name_ar AS published_user_last_name_ar,
			published_user.initial AS published_user_initial,
			published_user.color AS published_user_color,
			published_user.email AS published_user_email,
			published_user.profile_image_path AS published_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN documents ON documents.id = incident_associated_documents.document_id
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
			LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
		`).Where("incident_associated_documents.incident_register_id = ?", incidentID)

	err := query.Scan(&documentList).Error
	if err != nil {
		return nil, err
	}

	return documentList, nil

}

func (repo *IncidentReportRepository) GetAllIncidentControls(incidentID uint) ([]*commonDTOs.AssociatedControlListDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO

	query := repo.db.Table("incident_controls").
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
		LEFT JOIN controls ON controls.id = incident_controls.control_id
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
	`).Where("incident_controls.incident_register_id = ?", incidentID)

	err := query.Scan(&controlList).Error
	if err != nil {
		return nil, err
	}

	return controlList, nil

}

func (repo *IncidentReportRepository) GetAllIncidentRisks(incidentID uint) ([]*commonDTOs.AssociatedRiskListDTO, error) {
	var incidentRisk []*commonDTOs.AssociatedRiskListDTO

	query := repo.db.Table("incident_risks").
		Select(`
		    incident_risks.id,
			incident_risks.risk_register_id,
			incident_risks.incident_register_id,
			 risk_register.reference_number,
            risk_register.title,
            risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
			risk_register.current_risk_rating_value_id,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
        	users.color AS responsible_user_color,
        	users.email AS responsible_user_email,
        	users.profile_image_path AS responsible_user_profile_image_path,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
            risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
            risk_register.last_review_note,
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
            current_likelihood.title_en AS current_likelihood_title_en,
            current_likelihood.title_ar AS current_likelihood_title_ar,
            current_impact.title_en AS current_impact_title_en,
            current_impact.title_ar AS current_impact_title_ar,
            current_risk_rating.title_en AS current_risk_rating_title_en,
            current_risk_rating.title_ar AS current_risk_rating_title_ar,
            current_risk_rating_value.value AS current_risk_rating_value,
            current_risk_rating.color_code AS current_risk_rating_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = incident_risks.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
            LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id`).
		Where("incident_risks.incident_register_id = ?", incidentID)

	err := query.Scan(&incidentRisk).Error
	if err != nil {
		return nil, err
	}

	return incidentRisk, nil
}

func (repo *IncidentReportRepository) GetAllIncidentProcesses(incidentID uint) ([]*commonDTOs.AssociatedProcessListDTO, error) {
	var processList []*commonDTOs.AssociatedProcessListDTO

	query := repo.db.Table("incident_processes").
		Select(`
		incident_processes.id,
		incident_processes.process_id,
		processes.reference_number,
		processes.title_en,
		processes.title_ar,
		processes.description,
		processes.marked_as_auditable_item,
		process_groups.title_en AS process_group_title_en,
		process_groups.title_ar AS process_group_title_ar,
		processes.rto_from_value as rto_from_value,
		processes.rto_to_value as rto_to_value,
		processes.rto_unit_id as rto_unit_id,
		rto_bc_unit.title_en AS rto_unit_title_en,
		rto_bc_unit.title_ar AS rto_unit_title_ar,
		processes.rto_condition as rto_condition,
		processes.mtpd_from_value as mtpd_from_value,
		processes.mtpd_to_value as mtpd_to_value,
		processes.mtpd_unit_id as mtpd_unit_id,
		mtpd_bc_unit.title_en AS mtpd_unit_title_en,
		mtpd_bc_unit.title_ar AS mtpd_unit_title_ar,
		processes.mtpd_condition as mtpd_condition,
		processes.rpo_from_value as rpo_from_value,
		processes.rpo_to_value as rpo_to_value,
		processes.rpo_unit_id as rpo_unit_id,
		rpo_bc_unit.title_en AS rpo_unit_title_en,
		rpo_bc_unit.title_ar AS rpo_unit_title_ar,
		processes.rpo_condition as rpo_condition,
		processes.rsp_from_value as rsp_from_value,
		processes.rsp_to_value as rsp_to_value,
		processes.rsp_unit_id as rsp_unit_id,
		rsp_bc_unit.title_en AS rsp_unit_title_en,
		rsp_bc_unit.title_ar AS rsp_unit_title_ar,
		processes.rsp_condition as rsp_condition,
		processes.impact_analysis_level_id AS impact_analysis_level_id,
		bcm_impact_analysis_levels.title_en AS impact_analysis_level_title_en,
		bcm_impact_analysis_levels.title_ar AS impact_analysis_level_title_ar,
		bcm_impact_analysis_levels.color_code AS impact_analysis_level_color_code,
		bcm_impact_analysis_levels.type AS impact_analysis_level_type,
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
			LEFT JOIN processes ON processes.id = incident_processes.process_id
		    LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id
			LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id
			LEFT JOIN business_continuity_duration_units AS rto_bc_unit ON rto_bc_unit.id = processes.rto_unit_id
			LEFT JOIN business_continuity_duration_units AS mtpd_bc_unit ON mtpd_bc_unit.id = processes.mtpd_unit_id
			LEFT JOIN business_continuity_duration_units AS rpo_bc_unit ON rpo_bc_unit.id = processes.rpo_unit_id
			LEFT JOIN business_continuity_duration_units AS rsp_bc_unit ON rsp_bc_unit.id = processes.rsp_unit_id
			LEFT JOIN bcm_impact_analysis_levels ON bcm_impact_analysis_levels.id = processes.impact_analysis_level_id
			LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id
			LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id
			LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id
			LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id
		`).Where("incident_processes.incident_register_id = ?", incidentID)

	err := query.Scan(&processList).Error
	if err != nil {
		return nil, err
	}

	return processList, nil
}

// func (repo *IncidentReportRepository) GetAllIncidentComplianceRequirements(incidentID uint) ([]*incidentDTOs.IncidentComplianceRequirementListDTO, error) {
// 	var complianceRequirementList []*incidentDTOs.IncidentComplianceRequirementListDTO

// 	query := repo.db.Table("incident_compliance_requirements").
// 		Select(`
// 			incident_compliance_requirements.compliance_requirement_id AS compliance_requirement_id,
//             compliance_clauses.title_en AS compliance_clause_title,
//             compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
//             compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
//             compliance_requirements.title_en AS compliance_requirement_title,
//             compliance_statuses.title_en AS compliance_status_title_en,
//             compliance_statuses.title_ar AS compliance_status_title_ar,
// 			compliance_statuses.color_code AS compliance_status_color_code
//         `).
// 		Joins(`
// 			LEFT JOIN compliance_requirements ON compliance_requirements.id = incident_compliance_requirements.compliance_requirement_id
//             LEFT JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id
//             LEFT JOIN compliance_statuses ON compliance_requirements.compliance_status_id = compliance_statuses.id
//             LEFT JOIN compliance_requirement_categories ON compliance_requirements.compliance_requirement_category_id = compliance_requirement_categories.id
//         `).
// 		Where("incident_compliance_requirements.incident_register_id = ?", incidentID)

// 	err := query.Scan(&complianceRequirementList).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return complianceRequirementList, nil

// }

func (repo *IncidentReportRepository) GetAllIncidentComplianceRequirements(incidentID uint) ([]*commonDTOs.AssociatedClauseListDTO, error) {
	var clauseList []commonDTOs.AssociatedClauseRowDTO

	query := repo.db.Table("incident_compliance_requirements").
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
			LEFT JOIN compliance_clauses ON incident_compliance_requirements.compliance_clause_id = compliance_clauses.id
			LEFT JOIN compliance_register ON incident_compliance_requirements.compliance_register_id = compliance_register.id
		`).
		Where("incident_compliance_requirements.incident_register_id = ?", incidentID)

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

func (repo *IncidentReportRepository) Create(report *models.IncidentReport, tx utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(report).Error; err != nil {
		return 0, err
	}

	return report.ID, nil
}

func (repo *IncidentReportRepository) GetIncidentsByIDs(
	ids []uint,
	tx utility.DatabaseTransaction,
) ([]incidentDTOs.IncidentReportIncidentDTO, error) {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	var incidents []incidentDTOs.IncidentReportIncidentDTO

	if err := db.Table("incident_register").
		Select(`
            incident_register.id,
            incident_register.reference_number,
            incident_register.title AS incident_title,
            incident_register.description,
            incident_register.quick_action,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_sub_categories.title_en AS incident_sub_category_title_en,
            incident_sub_categories.title_ar AS incident_sub_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
            severity_levels.title_en AS severity_title_en,
            severity_levels.title_ar AS severity_title_ar,
            severity_levels.color_code AS severity_color_code,
            incident_register.incident_at AS real_time,
            incident_register.reported_at AS detected_time,
            ISNULL(incident_register.location, '') AS incident_location,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
            incident_statuses.color_code AS incident_status_color_code,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            organization_units.color_code AS responsible_unit_color_code,
            business_owner.first_name_en + ' ' + business_owner.last_name_en AS business_owner_name_en,
            business_owner.first_name_ar + ' ' + business_owner.last_name_ar AS business_owner_name_ar,
            response_team_lead.first_name_en + ' ' + response_team_lead.last_name_en AS response_team_lead_name_en,
            response_team_lead.first_name_ar + ' ' + response_team_lead.last_name_ar AS response_team_lead_name_ar,
            incident_scopes.title_en AS incident_scope_title_en,
            incident_scopes.title_ar AS incident_scope_title_ar,
            incident_sensitivity_classifications.title_en AS incident_sensitivity_classification_title_en,
            incident_sensitivity_classifications.title_ar AS incident_sensitivity_classification_title_ar,
            confidentiality_severity.title_en AS confidentiality_severity_level_title_en,
            confidentiality_severity.title_ar AS confidentiality_severity_level_title_ar,
            confidentiality_severity.color_code AS confidentiality_severity_level_color_code,
            integrity_severity.title_en AS integrity_severity_level_title_en,
            integrity_severity.title_ar AS integrity_severity_level_title_ar,
            integrity_severity.color_code AS integrity_severity_level_color_code,
            availability_severity.title_en AS availability_severity_level_title_en,
            availability_severity.title_ar AS availability_severity_level_title_ar,
            availability_severity.color_code AS availability_severity_level_color_code,
            incident_register.is_acknowledge,
            incident_register.acknowledged_at,
            acknowledged_user.first_name_en + ' ' + acknowledged_user.last_name_en AS acknowledged_by_name_en,
            acknowledged_user.first_name_ar + ' ' + acknowledged_user.last_name_ar AS acknowledged_by_name_ar,
            incident_register.is_sla_time_achieved,
            incident_register.is_sla_time_exceeded,
            incident_investigations.description AS containment_strategy,
            incident_investigations.action_recommendation AS next_steps,
            incident_closures.completion_date AS resolved_date,
            incident_closures.resolution_summary AS recovery_measures,
            incident_closures.resolution_summary AS closure_resolution_summary,
            incident_closures.notes AS closure_notes,
            incident_closures.submitted_at AS closure_submitted_at,
            closure_submitted_user.first_name_en + ' ' + closure_submitted_user.last_name_en AS closure_submitted_by_name_en,
            closure_submitted_user.first_name_ar + ' ' + closure_submitted_user.last_name_ar AS closure_submitted_by_name_ar,
            incident_closures.approved_at AS closure_approved_at,
            closure_approved_user.first_name_en + ' ' + closure_approved_user.last_name_en AS closure_approved_by_name_en,
            closure_approved_user.first_name_ar + ' ' + closure_approved_user.last_name_ar AS closure_approved_by_name_ar,
            workflow_statuses.title_en AS closure_workflow_status_title_en,
            workflow_statuses.title_ar AS closure_workflow_status_title_ar,
            workflow_statuses.color_code AS closure_workflow_status_color_code,
            assigned_user.initial AS assigned_user_initial,
            assigned_user.first_name_en AS assigned_user_first_name_en,
            assigned_user.last_name_en AS assigned_user_last_name_en,
            assigned_user.first_name_ar AS assigned_user_first_name_ar,
            assigned_user.last_name_ar AS assigned_user_last_name_ar,
            ISNULL(responsible_user.first_name_en, '') + ' ' + ISNULL(responsible_user.last_name_en, '') AS responsible_user_name_en,
            ISNULL(responsible_user.first_name_ar, '') + ' ' + ISNULL(responsible_user.last_name_ar, '') AS responsible_user_name_ar,
            ISNULL(reported_user.first_name_en, '') + ' ' + ISNULL(reported_user.last_name_en, '') AS reported_user_name_en,
            ISNULL(reported_user.first_name_ar, '') + ' ' + ISNULL(reported_user.last_name_ar, '') AS reported_user_name_ar
        `).
		Joins("LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id").
		Joins("LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id").
		Joins("LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id").
		Joins("LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id").
		Joins("LEFT JOIN users AS business_owner ON business_owner.id = incident_register.business_owner_id").
		Joins("LEFT JOIN users AS response_team_lead ON response_team_lead.id = incident_register.response_team_lead_id").
		Joins("LEFT JOIN incident_scopes ON incident_scopes.id = incident_register.incident_scope_id").
		Joins("LEFT JOIN incident_sensitivity_classifications ON incident_sensitivity_classifications.id = incident_register.incident_sensitivity_classification_id").
		Joins("LEFT JOIN severity_levels AS confidentiality_severity ON confidentiality_severity.id = incident_register.confidentiality_severity_level_id").
		Joins("LEFT JOIN severity_levels AS integrity_severity ON integrity_severity.id = incident_register.integrity_severity_level_id").
		Joins("LEFT JOIN severity_levels AS availability_severity ON availability_severity.id = incident_register.availability_severity_level_id").
		Joins("LEFT JOIN users AS acknowledged_user ON acknowledged_user.id = incident_register.acknowledged_by_id").
		Joins("LEFT JOIN incident_investigations ON incident_investigations.incident_register_id = incident_register.id").
		Joins("LEFT JOIN users assigned_user ON assigned_user.id = incident_investigations.team_lead_id").
		Joins("LEFT JOIN incident_closures ON incident_closures.incident_register_id = incident_register.id").
		Joins("LEFT JOIN users AS closure_submitted_user ON closure_submitted_user.id = incident_closures.submitted_by_id").
		Joins("LEFT JOIN users AS closure_approved_user ON closure_approved_user.id = incident_closures.approved_by_id").
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = incident_closures.workflow_status_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id").
		Joins("LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by").
		Where("incident_register.id IN ?", ids).
		Order("incident_register.created_at DESC").
		Find(&incidents).Error; err != nil {

		return nil, err
	}

	return incidents, nil
}

func (repo *IncidentReportRepository) GetReportIncidentsByDateRange(
	start time.Time,
	end time.Time,
	tx utility.DatabaseTransaction,
) ([]incidentDTOs.IncidentReportIncidentDTO, error) {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	var incidents []incidentDTOs.IncidentReportIncidentDTO

	if err := db.Table("incident_register").
		Select(`
            incident_register.id,
            incident_register.reference_number,
            incident_register.title AS incident_title,
            incident_register.description,
            incident_register.quick_action,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_sub_categories.title_en AS incident_sub_category_title_en,
            incident_sub_categories.title_ar AS incident_sub_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
            severity_levels.title_en AS severity_title_en,
            severity_levels.title_ar AS severity_title_ar,
            severity_levels.color_code AS severity_color_code,
            incident_register.incident_at AS real_time,
            incident_register.reported_at AS detected_time,
            ISNULL(incident_register.location, '') AS incident_location,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
            incident_statuses.color_code AS incident_status_color_code,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            organization_units.color_code AS responsible_unit_color_code,
            business_owner.first_name_en + ' ' + business_owner.last_name_en AS business_owner_name_en,
            business_owner.first_name_ar + ' ' + business_owner.last_name_ar AS business_owner_name_ar,
            response_team_lead.first_name_en + ' ' + response_team_lead.last_name_en AS response_team_lead_name_en,
            response_team_lead.first_name_ar + ' ' + response_team_lead.last_name_ar AS response_team_lead_name_ar,
            incident_scopes.title_en AS incident_scope_title_en,
            incident_scopes.title_ar AS incident_scope_title_ar,
            incident_sensitivity_classifications.title_en AS incident_sensitivity_classification_title_en,
            incident_sensitivity_classifications.title_ar AS incident_sensitivity_classification_title_ar,
            confidentiality_severity.title_en AS confidentiality_severity_level_title_en,
            confidentiality_severity.title_ar AS confidentiality_severity_level_title_ar,
            confidentiality_severity.color_code AS confidentiality_severity_level_color_code,
            integrity_severity.title_en AS integrity_severity_level_title_en,
            integrity_severity.title_ar AS integrity_severity_level_title_ar,
            integrity_severity.color_code AS integrity_severity_level_color_code,
            availability_severity.title_en AS availability_severity_level_title_en,
            availability_severity.title_ar AS availability_severity_level_title_ar,
            availability_severity.color_code AS availability_severity_level_color_code,
            incident_register.is_acknowledge,
            incident_register.acknowledged_at,
            acknowledged_user.first_name_en + ' ' + acknowledged_user.last_name_en AS acknowledged_by_name_en,
            acknowledged_user.first_name_ar + ' ' + acknowledged_user.last_name_ar AS acknowledged_by_name_ar,
            incident_register.is_sla_time_achieved,
            incident_register.is_sla_time_exceeded,
            incident_investigations.description AS containment_strategy,
            incident_investigations.action_recommendation AS next_steps,
            incident_closures.completion_date AS resolved_date,
            incident_closures.resolution_summary AS recovery_measures,
            incident_closures.resolution_summary AS closure_resolution_summary,
            incident_closures.notes AS closure_notes,
            incident_closures.submitted_at AS closure_submitted_at,
            closure_submitted_user.first_name_en + ' ' + closure_submitted_user.last_name_en AS closure_submitted_by_name_en,
            closure_submitted_user.first_name_ar + ' ' + closure_submitted_user.last_name_ar AS closure_submitted_by_name_ar,
            incident_closures.approved_at AS closure_approved_at,
            closure_approved_user.first_name_en + ' ' + closure_approved_user.last_name_en AS closure_approved_by_name_en,
            closure_approved_user.first_name_ar + ' ' + closure_approved_user.last_name_ar AS closure_approved_by_name_ar,
            workflow_statuses.title_en AS closure_workflow_status_title_en,
            workflow_statuses.title_ar AS closure_workflow_status_title_ar,
            workflow_statuses.color_code AS closure_workflow_status_color_code,
            assigned_user.initial AS assigned_user_initial,
            assigned_user.first_name_en AS assigned_user_first_name_en,
            assigned_user.last_name_en AS assigned_user_last_name_en,
            assigned_user.first_name_ar AS assigned_user_first_name_ar,
            assigned_user.last_name_ar AS assigned_user_last_name_ar,
            ISNULL(responsible_user.first_name_en, '') + ' ' + ISNULL(responsible_user.last_name_en, '') AS responsible_user_name_en,
            ISNULL(responsible_user.first_name_ar, '') + ' ' + ISNULL(responsible_user.last_name_ar, '') AS responsible_user_name_ar,
            ISNULL(reported_user.first_name_en, '') + ' ' + ISNULL(reported_user.last_name_en, '') AS reported_user_name_en,
            ISNULL(reported_user.first_name_ar, '') + ' ' + ISNULL(reported_user.last_name_ar, '') AS reported_user_name_ar
        `).
		Joins("LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id").
		Joins("LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id").
		Joins("LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id").
		Joins("LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id").
		Joins("LEFT JOIN users AS business_owner ON business_owner.id = incident_register.business_owner_id").
		Joins("LEFT JOIN users AS response_team_lead ON response_team_lead.id = incident_register.response_team_lead_id").
		Joins("LEFT JOIN incident_scopes ON incident_scopes.id = incident_register.incident_scope_id").
		Joins("LEFT JOIN incident_sensitivity_classifications ON incident_sensitivity_classifications.id = incident_register.incident_sensitivity_classification_id").
		Joins("LEFT JOIN severity_levels AS confidentiality_severity ON confidentiality_severity.id = incident_register.confidentiality_severity_level_id").
		Joins("LEFT JOIN severity_levels AS integrity_severity ON integrity_severity.id = incident_register.integrity_severity_level_id").
		Joins("LEFT JOIN severity_levels AS availability_severity ON availability_severity.id = incident_register.availability_severity_level_id").
		Joins("LEFT JOIN users AS acknowledged_user ON acknowledged_user.id = incident_register.acknowledged_by_id").
		Joins("LEFT JOIN incident_investigations ON incident_investigations.incident_register_id = incident_register.id").
		Joins("LEFT JOIN users assigned_user ON assigned_user.id = incident_investigations.team_lead_id").
		Joins("LEFT JOIN incident_closures ON incident_closures.incident_register_id = incident_register.id").
		Joins("LEFT JOIN users AS closure_submitted_user ON closure_submitted_user.id = incident_closures.submitted_by_id").
		Joins("LEFT JOIN users AS closure_approved_user ON closure_approved_user.id = incident_closures.approved_by_id").
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = incident_closures.workflow_status_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id").
		Joins("LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by").
		Where("incident_register.reported_at BETWEEN ? AND ?", start, end).
		Order("incident_register.created_at DESC").
		Find(&incidents).Error; err != nil {
		return nil, err
	}

	return incidents, nil
}

func (repo *IncidentReportRepository) CreateSection(
	section *models.IncidentReportSection,
	tx utility.DatabaseTransaction,
) (uint, error) {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(section).Error; err != nil {
		return 0, err
	}

	return section.ID, nil
}

func (repo *IncidentReportRepository) CreateReportIncident(incident *models.IncidentReportIncident, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(incident).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentReportRepository) GetDocumentsByIncidentID(
	incidentID uint,
	tx utility.DatabaseTransaction,
) ([]incidentDTOs.IncidentReportDocumentDTO, error) {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	var documents []incidentDTOs.IncidentReportDocumentDTO

	query := db.Table("incident_documents").
		Select(`
            incident_documents.id,
            incident_documents.url,
            incident_documents.name,
            incident_documents.size,
            incident_documents.token,
            incident_documents.original_name,
            incident_documents.type
        `).
		Where("incident_documents.incident_register_id = ?", incidentID)

	if err := query.Find(&documents).Error; err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *IncidentReportRepository) GetEvidencesByIncidentID(
	incidentID uint,
	tx utility.DatabaseTransaction,
) ([]incidentDTOs.IncidentReportEvidenceDTO, error) {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	var evidences []incidentDTOs.IncidentReportEvidenceDTO

	query := db.Table("incident_evidences").
		Select(`
            incident_evidences.id,
            incident_evidences.url,
            incident_evidences.name,
            incident_evidences.size,
            incident_evidences.token,
            incident_evidences.original_name,
            incident_evidences.type
        `).
		Where("incident_evidences.incident_register_id = ?", incidentID)

	if err := query.Find(&evidences).Error; err != nil {
		return nil, err
	}

	return evidences, nil
}

func (repo *IncidentReportRepository) CreateIncidentDocuments(
	docs []*models.IncidentReportDocument,
	tx utility.DatabaseTransaction,
) error {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentReportRepository) CreateIncidentEvidences(
	evidences []*models.IncidentReportEvidence,
	tx utility.DatabaseTransaction,
) error {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(&evidences).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentReportRepository) CreateIncidentRootCauseAnalyses(
	rcas []*models.IncidentReportRootCauseAnalysis,
	tx utility.DatabaseTransaction,
) error {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(&rcas).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentReportRepository) CreateIncidentReportActionPlans(
	plans []*models.IncidentReportActionPlan,
	tx utility.DatabaseTransaction,
) error {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(&plans).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentReportRepository) GetIncidentReportDetailsByID(id uint) (*models.IncidentReport, error) {
	var report models.IncidentReport

	err := repo.db.
		Preload("CreatedByUser").
		Preload("Sections").
		Preload("Incidents").
		Preload("Incidents.Documents").
		Preload("Incidents.Evidences").
		Preload("Incidents.RootCauseAnalyses").
		Preload("Incidents.ActionPlans", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, incident_report_id, incident_report_incident_id, action_plan_title, action_plan_description, target_date, responsible_user_name_en, responsible_user_name_ar, severity_level_title_en, severity_level_title_ar, action_plan_status_title_en, action_plan_status_title_ar, action_plan_status_color_code")
		}).
		Preload("Incidents.IncidentRegister").
		Preload("Incidents.IncidentRegister.Investigation").
		Preload("Incidents.IncidentRegister.Investigation.TeamLead").
		Preload("Incidents.IncidentRegister.Investigation.Users").
		Preload("Incidents.IncidentRegister.Investigation.Users.User").
		Preload("Incidents.IncidentRegister.Investigation.Documents").
		Preload("Incidents.IncidentRegister.Investigation.Findings").
		Preload("Incidents.IncidentRegister.Investigation.CreatedByUser").
		Preload("Incidents.IncidentRegister.Investigation.UpdatedByUser").
		Preload("Incidents.IncidentRegister.Timelines").
		Preload("Incidents.IncidentRegister.Timelines.PerformedBy").
		Preload("Incidents.IncidentRegister.Timelines.TimelineStatus").
		Preload("Incidents.IncidentRegister.Timelines.IncidentTimelineType").
		Preload("Incidents.IncidentRegister.CommunicationLogs").
		Preload("Incidents.IncidentRegister.CommunicationLogs.CreatedByUser").
		Preload("Incidents.IncidentRegister.CommunicationLogs.UpdatedByUser").
		Preload("Incidents.IncidentRegister.IncidentClosure").
		Preload("Incidents.IncidentRegister.IncidentClosure.Lessons").
		Preload("Incidents.IncidentRegister.IncidentClosure.CreatedByUser").
		Preload("Incidents.IncidentRegister.IncidentClosure.UpdatedByUser").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.WorkflowLevelType").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.WorkflowStatus").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.CreatedBy").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.UpdatedBy").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowParticipants").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.User").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.DelegateeUser").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.Role").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowParticipants.WorkflowStatus").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowReviews").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowReviews.WorkflowStatus").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowReviews.CreatedByUser").
		Preload("Incidents.IncidentRegister.IncidentClosure.IncidentClosureWorkflows.IncidentClosureWorkflowReviews.UpdatedByUser").
		Preload("Incidents.IncidentRegister.IncidentPostReview").
		Preload("Incidents.IncidentRegister.IncidentPostReview.Users").
		Preload("Incidents.IncidentRegister.IncidentPostReview.Users.User").
		Preload("Incidents.IncidentRegister.IncidentPostReview.Recommendations").
		Preload("Incidents.IncidentRegister.IncidentPostReview.Lessons").
		Preload("Incidents.IncidentRegister.IncidentPostReview.IncidentPostReviewEffectiveness").
		Preload("Incidents.IncidentRegister.IncidentPostReview.CreatedByUser").
		Preload("Incidents.IncidentRegister.IncidentPostReview.UpdatedByUser").
		Where("id = ?", id).
		First(&report).Error

	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *IncidentReportRepository) GetIncidentReportByID(incidentReportID uint, tx ...utility.DatabaseTransaction) (*models.IncidentReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var report models.IncidentReport
	err := db.Where("id = ?", incidentReportID).First(&report).Error
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *IncidentReportRepository) DeleteIncidentReportDocuments(incidentReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("incident_report_incident_id IN (SELECT id FROM incident_report_incidents WHERE incident_report_id = ?)", incidentReportID).Delete(&models.IncidentReportDocument{}).Error
}

func (repo *IncidentReportRepository) DeleteIncidentReportRootCauseAnalyses(incidentReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("incident_report_id = ?", incidentReportID).Delete(&models.IncidentReportRootCauseAnalysis{}).Error
}

func (repo *IncidentReportRepository) DeleteIncidentReportActionPlans(incidentReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("incident_report_id = ?", incidentReportID).Delete(&models.IncidentReportActionPlan{}).Error
}

func (repo *IncidentReportRepository) DeleteIncidentReportIncidents(incidentReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("incident_report_id = ?", incidentReportID).Delete(&models.IncidentReportIncident{}).Error
}

func (repo *IncidentReportRepository) DeleteIncidentReportSections(incidentReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("incident_report_id = ?", incidentReportID).Delete(&models.IncidentReportSection{}).Error
}

func (repo *IncidentReportRepository) DeleteIncidentReport(incidentReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", incidentReportID).Delete(&models.IncidentReport{}).Error
}

func (repo *IncidentReportRepository) GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.IncidentReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section models.IncidentReportSection
	err := db.Where("id = ?", sectionID).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (repo *IncidentReportRepository) UpdateSection(section *models.IncidentReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (repo *IncidentReportRepository) DeleteIncidentReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.IncidentReportSection{}).Error
}

func (repo *IncidentReportRepository) ReorderIncidentReportSections(incidentReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Shift order numbers of remaining sections
	err := db.Exec(`
		UPDATE incident_report_sections
		SET order_num = order_num - 1
		WHERE incident_report_id = ? AND order_num > ?`,
		incidentReportID, deletedOrderNum).Error

	return err
}
