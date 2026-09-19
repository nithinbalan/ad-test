package repositories

import (
	"errors"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetFullReportDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error)
	GetAllMilestones(projectID uint) ([]*dtos.ProjectMilestoneListDTO, error)
	GetAllControls(projectID uint) ([]*commonDTOs.AssociatedControlListDTO, error)
	GetAllRisks(projectID uint) ([]*commonDTOs.AssociatedRiskListDTO, error)
	GetAllProcesses(projectID uint) ([]*commonDTOs.AssociatedProcessListDTO, error)
	GetAllStrategicObjectives(projectID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error)
	GetAllKPIs(projectID uint) ([]*commonDTOs.AssociatedKPIRegisterListDTO, error)
	GetAllAssociatedDocuments(projectID uint) ([]*commonDTOs.AssociatedDocumentListDTO, error)
	GetAllAssets(projectID uint) ([]*commonDTOs.AssociatedAssetListDTO, error)
	GetAllIncidents(projectID uint) ([]*commonDTOs.AssociatedIncidentListDTO, error)
	GetAllComplianceRequirements(projectID uint) ([]*commonDTOs.AssociatedClauseListDTO, error)
	GetAllActionPlans(projectID uint) ([]*dtos.ProjectActionPlanListDTO, error)
	GetAllChangeRequests(projectID uint) ([]*dtos.ProjectChangeRequestListDTO, error)
	GetAllDecisionLogs(projectID uint) ([]*dtos.ProjectDecisionLogListDTO, error)
	GetAllSupportingDocuments(projectID uint) ([]*dtos.ProjectSupportingDocumentListDTO, error)
	GetProjectClosure(projectID uint) (*dtos.ProjectClosureReportDTO, error)

	GetAllReports(params *dtos.ProjectReportListQueryParamsDTO) (*dtos.ProjectReportListResponseDTO, error)
	GetProjectsByIDs(ids []uint, tx ...utility.DatabaseTransaction) ([]*models.Project, error)
	GetProjectsByDateRange(start time.Time, end time.Time, tx ...utility.DatabaseTransaction) ([]*models.Project, error)
	CreateReport(report *models.ProjectReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.ProjectReportSection, tx ...utility.DatabaseTransaction) (uint, error)
	CreateReportProject(project *models.ProjectReportProject, tx ...utility.DatabaseTransaction) error
	CreateProjectDocuments(docs []*models.ProjectReportDocument, tx ...utility.DatabaseTransaction) error
	GetReportByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.ProjectReport, error)
	GetReportDetailsByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.ProjectReport, error)
	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.ProjectReportSection, error)
	UpdateSection(section *models.ProjectReportSection, tx ...utility.DatabaseTransaction) error
	DeleteReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	ReorderReportSections(reportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error
	DeleteReportDocuments(reportID uint, tx ...utility.DatabaseTransaction) error
	DeleteReportProjects(reportID uint, tx ...utility.DatabaseTransaction) error
	DeleteReportSections(reportID uint, tx ...utility.DatabaseTransaction) error
	DeleteReport(reportID uint, tx ...utility.DatabaseTransaction) error
}

type ProjectReportRepository struct {
	db *gorm.DB
}

func NewProjectReportRepository(db *gorm.DB) *ProjectReportRepository {
	return &ProjectReportRepository{db: db}
}

func (r *ProjectReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(r.db)
}

func (r *ProjectReportRepository) GetFullReportDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Project, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var project models.Project
	result := db.
		Preload("ResponsibleUser").
		Preload("ProjectType").
		Preload("ProjectStatus").
		Preload("ProjectPhase").
		Preload("ProjectClassification").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("IncidentActionPlan").
		Preload("IncidentActionPlan.ActionPlanStatus").
		Preload("IncidentActionPlan.ResponsibleUser").
		Preload("IncidentActionPlan.IncidentRegister").
		First(&project, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (r *ProjectReportRepository) GetAllMilestones(projectID uint) ([]*dtos.ProjectMilestoneListDTO, error) {
	var milestones []*dtos.ProjectMilestoneListDTO
	err := r.db.Table("project_milestones").
		Select(`
			project_milestones.id,
			project_milestones.project_id,
			project_milestones.title_en,
			project_milestones.title_ar,
			project_milestones.description_en,
			project_milestones.description_ar,
			CONVERT(VARCHAR, project_milestones.planned_date, 23) AS planned_date,
			CONVERT(VARCHAR, project_milestones.actual_date, 23) AS actual_date,
			project_milestones.project_milestone_status_id,
			status.title_en AS project_milestone_status_title_en,
			status.title_ar AS project_milestone_status_title_ar,
			status.type AS project_milestone_status_type,
			status.color_code AS project_milestone_status_color_code,
			project_milestones.project_phase_id,
			phase.title_en AS project_phase_title_en,
			phase.title_ar AS project_phase_title_ar,
			phase.color_code AS project_phase_color_code,
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
		Joins("LEFT JOIN project_milestone_statuses AS status ON status.id = project_milestones.project_milestone_status_id").
		Joins("LEFT JOIN project_phases AS phase ON phase.id = project_milestones.project_phase_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_milestones.created_by_id").
		Where("project_milestones.project_id = ?", projectID).
		Order("project_milestones.id ASC").
		Scan(&milestones).Error

	return milestones, err
}

func (r *ProjectReportRepository) GetAllControls(projectID uint) ([]*commonDTOs.AssociatedControlListDTO, error) {
	var controls []*commonDTOs.AssociatedControlListDTO
	err := r.db.Table("project_controls").
		Select(`
			project_controls.id,
			project_controls.control_id,
			controls.reference_number,
			controls.code,
			controls.title_en,
			controls.title_ar,
			controls.description,
			controls.created_at,
			controls.updated_at,
			control_types.title_en AS control_type_title_en,
			control_types.title_ar AS control_type_title_ar,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.first_name_en AS owner_first_name_en,
			owner.last_name_en AS owner_last_name_en
		`).
		Joins("INNER JOIN controls ON controls.id = project_controls.control_id").
		Joins("LEFT JOIN control_types ON control_types.id = controls.control_type_id").
		Joins("LEFT JOIN control_categories ON control_categories.id = controls.control_category_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = controls.responsible_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = controls.responsible_user_id").
		Where("project_controls.project_id = ?", projectID).
		Order("project_controls.id ASC").
		Scan(&controls).Error

	return controls, err
}

func (r *ProjectReportRepository) GetAllRisks(projectID uint) ([]*commonDTOs.AssociatedRiskListDTO, error) {
	var risks []*commonDTOs.AssociatedRiskListDTO
	err := r.db.Table("project_risks").
		Select(`
			project_risks.id,
			project_risks.risk_register_id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			CONVERT(VARCHAR, risk_register.created_at, 120) AS created_at,
			CONVERT(VARCHAR, risk_register.updated_at, 120) AS updated_at,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
			residual_risk_rating.title_en AS residual_risk_rating_title_en
		`).
		Joins("INNER JOIN risk_register ON risk_register.id = project_risks.risk_register_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id").
		Where("project_risks.project_id = ?", projectID).
		Order("project_risks.id ASC").
		Scan(&risks).Error

	return risks, err
}

func (r *ProjectReportRepository) GetAllProcesses(projectID uint) ([]*commonDTOs.AssociatedProcessListDTO, error) {
	var processes []*commonDTOs.AssociatedProcessListDTO
	err := r.db.Table("project_processes").
		Select(`
			project_processes.id,
			project_processes.process_id,
			processes.reference_number,
			processes.title_en,
			processes.title_ar,
			processes.description,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar
		`).
		Joins("INNER JOIN processes ON processes.id = project_processes.process_id").
		Joins("LEFT JOIN process_groups ON process_groups.id = processes.process_group_id").
		Where("project_processes.project_id = ?", projectID).
		Order("project_processes.id ASC").
		Scan(&processes).Error

	return processes, err
}

func (r *ProjectReportRepository) GetAllStrategicObjectives(projectID uint) ([]*commonDTOs.AssociatedStrategicObjectiveListDTO, error) {
	var objectives []*commonDTOs.AssociatedStrategicObjectiveListDTO
	err := r.db.Table("project_strategic_objectives").
		Select(`
			project_strategic_objectives.id,
			project_strategic_objectives.strategic_objective_id,
			strategic_objectives.reference_number,
			strategic_objectives.title_en,
			strategic_objectives.title_ar,
			strategic_objectives.description_en,
			strategic_objectives.description_ar
		`).
		Joins("INNER JOIN strategic_objectives ON strategic_objectives.id = project_strategic_objectives.strategic_objective_id").
		Where("project_strategic_objectives.project_id = ?", projectID).
		Order("project_strategic_objectives.id ASC").
		Scan(&objectives).Error

	return objectives, err
}

func (r *ProjectReportRepository) GetAllKPIs(projectID uint) ([]*commonDTOs.AssociatedKPIRegisterListDTO, error) {
	var kpis []*commonDTOs.AssociatedKPIRegisterListDTO
	err := r.db.Table("kpi_projects").
		Select(`
			kpi_projects.id,
			kpi_projects.kpi_register_id,
			kpi_register.reference_number,
			kpi_register.title,
			kpi_register.description,
			kpi_register.target_value,
			CONVERT(VARCHAR, kpi_register.created_at, 120) AS created_at,
			CONVERT(VARCHAR, kpi_register.updated_at, 120) AS updated_at
		`).
		Joins("INNER JOIN kpi_register ON kpi_register.id = kpi_projects.kpi_register_id").
		Where("kpi_projects.project_id = ?", projectID).
		Order("kpi_projects.id ASC").
		Scan(&kpis).Error

	return kpis, err
}

func (r *ProjectReportRepository) GetAllAssociatedDocuments(projectID uint) ([]*commonDTOs.AssociatedDocumentListDTO, error) {
	var docs []*commonDTOs.AssociatedDocumentListDTO
	err := r.db.Table("project_associated_documents").
		Select(`
			project_associated_documents.id AS id,
			project_associated_documents.document_id,
			documents.reference_number AS reference_number,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			documents.description_en AS description_en,
			documents.description_ar AS description_ar,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en
		`).
		Joins("INNER JOIN documents ON documents.id = project_associated_documents.document_id").
		Joins("LEFT JOIN document_categories ON document_categories.id = documents.document_category_id").
		Joins("LEFT JOIN document_types ON document_types.id = documents.document_type_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id").
		Where("project_associated_documents.project_id = ?", projectID).
		Order("project_associated_documents.id ASC").
		Scan(&docs).Error

	return docs, err
}

func (r *ProjectReportRepository) GetAllAssets(projectID uint) ([]*commonDTOs.AssociatedAssetListDTO, error) {
	var assets []*commonDTOs.AssociatedAssetListDTO
	err := r.db.Table("project_assets").
		Select(`
			project_assets.id,
			project_assets.asset_register_id,
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
			responsible_user.profile_image_path AS responsible_user_profile_image_path
		`).
		Joins("INNER JOIN asset_register ON asset_register.id = project_assets.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN asset_criticality_matrix_rating_scales ON asset_criticality_matrix_rating_scales.id = asset_register.asset_criticality_matrix_rating_scale_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Where("project_assets.project_id = ?", projectID).
		Order("project_assets.id ASC").
		Scan(&assets).Error

	return assets, err
}

func (r *ProjectReportRepository) GetAllIncidents(projectID uint) ([]*commonDTOs.AssociatedIncidentListDTO, error) {
	var incidents []*commonDTOs.AssociatedIncidentListDTO
	err := r.db.Table("project_incidents").
		Select(`
			incident_register.*,
			project_incidents.id as id,
			project_incidents.incident_register_id,
			incident_categories.title_en AS incident_category_title_en,
			incident_categories.title_ar AS incident_category_title_ar,
			incident_types.title_en AS incident_type_title_en,
			incident_types.title_ar AS incident_type_title_ar,
			incident_statuses.title_en AS incident_status_title_en,
			incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
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
			responsible_user.profile_image_path AS responsible_user_profile_image_path
		`).
		Joins("INNER JOIN incident_register ON incident_register.id = project_incidents.incident_register_id").
		Joins("LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id").
		Joins("LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id").
		Joins("LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id").
		Joins("LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id").
		Where("project_incidents.project_id = ?", projectID).
		Order("project_incidents.id ASC").
		Scan(&incidents).Error

	return incidents, err
}

func (r *ProjectReportRepository) GetAllComplianceRequirements(projectID uint) ([]*commonDTOs.AssociatedClauseListDTO, error) {
	var clauseList []commonDTOs.AssociatedClauseRowDTO

	query := r.db.Table("project_compliance_requirements").
		Select(`
			compliance_clauses.title_en AS clause_title_en,
			compliance_clauses.title_ar AS clause_title_ar,
			compliance_clauses.clause_no,
			compliance_clauses.compliance_register_id,
			compliance_clauses.id AS compliance_clause_id,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar
		`).
		Joins("LEFT JOIN compliance_clauses ON project_compliance_requirements.compliance_clause_id = compliance_clauses.id").
		Joins("LEFT JOIN compliance_register ON project_compliance_requirements.compliance_register_id = compliance_register.id").
		Where("project_compliance_requirements.compliance_clause_id IS NOT NULL").
		Where("project_compliance_requirements.project_id = ?", projectID)

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

	return groupedClauses, nil
}

func (r *ProjectReportRepository) GetAllActionPlans(projectID uint) ([]*dtos.ProjectActionPlanListDTO, error) {
	var actionPlans []*dtos.ProjectActionPlanListDTO
	err := r.db.Table("project_action_plans").
		Select(`
			project_action_plans.*,
			CONVERT(VARCHAR, project_action_plans.target_date, 23) AS target_date,
			projects.title_en AS project_title_en,
			projects.title_ar AS project_title_ar,
			project_milestones.title_en AS project_milestone_title_en,
			project_milestones.title_ar AS project_milestone_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			action_plan_statuses.title_en AS action_plan_status_title_en,
			action_plan_statuses.title_ar AS action_plan_status_title_ar,
			action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar
		`).
		Joins("LEFT JOIN projects ON projects.id = project_action_plans.project_id").
		Joins("LEFT JOIN project_milestones ON project_milestones.id = project_action_plans.project_milestone_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = project_action_plans.severity_level_id").
		Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = project_action_plans.action_plan_status_id").
		Joins("LEFT JOIN users ON users.id = project_action_plans.responsible_user_id").
		Where("project_action_plans.project_id = ?", projectID).
		Order("project_action_plans.id ASC").
		Scan(&actionPlans).Error

	return actionPlans, err
}

func (r *ProjectReportRepository) GetAllChangeRequests(projectID uint) ([]*dtos.ProjectChangeRequestListDTO, error) {
	var changeRequests []*dtos.ProjectChangeRequestListDTO
	err := r.db.Table("project_change_requests").
		Select(`
			project_change_requests.id,
			project_change_requests.reference_number,
			project_change_requests.project_id,
			project_change_requests.title,
			project_change_requests.description,
			project_change_requests.business_justification,
			project_change_requests.anticipated_benefits,
			project_change_requests.risks_of_deferral,
			COALESCE(cr_impacts.impact, '') AS impact,
			project_change_requests.requested_by_user_id,
			CONVERT(VARCHAR, project_change_requests.requested_on, 23) AS requested_on,
			project_change_requests.workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			CONVERT(VARCHAR, project_change_requests.created_at, 120) AS created_at,
			users.id AS requested_by_user_id,
			users.first_name_en AS requested_by_user_first_name_en,
			users.first_name_ar AS requested_by_user_first_name_ar,
			users.last_name_en AS requested_by_user_last_name_en,
			users.last_name_ar AS requested_by_user_last_name_ar
		`).
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = project_change_requests.workflow_status_id").
		Joins(`LEFT JOIN (
			SELECT 
				pcia.project_change_request_id,
				STRING_AGG(CAST(pcia.title_en AS NVARCHAR(MAX)), ', ') AS impact
			FROM (
				SELECT DISTINCT 
					ia.project_change_request_id,
					piat.title_en
				FROM project_change_request_impact_assessments ia
				JOIN project_impact_assessment_types piat ON piat.id = ia.project_impact_assessment_type_id
			) pcia
			GROUP BY pcia.project_change_request_id
		) cr_impacts ON cr_impacts.project_change_request_id = project_change_requests.id`).
		Joins("LEFT JOIN users ON users.id = project_change_requests.requested_by_user_id").
		Where("project_change_requests.project_id = ?", projectID).
		Order("project_change_requests.id ASC").
		Scan(&changeRequests).Error

	return changeRequests, err
}

func (r *ProjectReportRepository) GetAllDecisionLogs(projectID uint) ([]*dtos.ProjectDecisionLogListDTO, error) {
	var decisionLogs []*dtos.ProjectDecisionLogListDTO
	err := r.db.Table("project_decision_logs").
		Select(`
			project_decision_logs.id,
			project_decision_logs.reference_number,
			project_decision_logs.project_id,
			project_decision_logs.decision,
			project_decision_logs.decision_notes,
			project_decision_logs.reason,
			CONVERT(VARCHAR, project_decision_logs.decision_date, 23) AS decision_date,
			project_decision_logs.requested_by_user_id,
			project_decision_logs.workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			CONVERT(VARCHAR, project_decision_logs.created_at, 120) AS created_at,
			users.id AS requested_by_user_id,
			users.first_name_en AS requested_by_user_first_name_en,
			users.first_name_ar AS requested_by_user_first_name_ar,
			users.last_name_en AS requested_by_user_last_name_en,
			users.last_name_ar AS requested_by_user_last_name_ar
		`).
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = project_decision_logs.workflow_status_id").
		Joins("LEFT JOIN users ON users.id = project_decision_logs.requested_by_user_id").
		Where("project_decision_logs.project_id = ?", projectID).
		Order("project_decision_logs.id ASC").
		Scan(&decisionLogs).Error

	return decisionLogs, err
}

func (r *ProjectReportRepository) GetAllSupportingDocuments(projectID uint) ([]*dtos.ProjectSupportingDocumentListDTO, error) {
	var supportingDocs []*dtos.ProjectSupportingDocumentListDTO
	err := r.db.Table("project_supporting_documents").
		Select(`
			project_supporting_documents.id,
			project_supporting_documents.project_id,
			project_supporting_documents.title_en,
			project_supporting_documents.title_ar,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			project_supporting_documents.created_at,
			project_supporting_documents.updated_at
		`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_supporting_documents.created_by_id").
		Where("project_supporting_documents.project_id = ?", projectID).
		Order("project_supporting_documents.id ASC").
		Scan(&supportingDocs).Error

	return supportingDocs, err
}

func (r *ProjectReportRepository) GetProjectClosure(projectID uint) (*dtos.ProjectClosureReportDTO, error) {
	var exists int64
	if err := r.db.Table("project_closures").Where("project_id = ?", projectID).Count(&exists).Error; err != nil || exists == 0 {
		return nil, nil
	}

	type closureScan struct {
		ID                    uint
		OverallProjectOutcome *string
		ClosureReason         *string
		ClosureComments       *string
		LessonLearned         *string
		ApproverID            uint
		IsApproved            bool
		ApproverFirstNameEn   *string
		ApproverLastNameEn    *string
		ApproverEmail         *string
		CreatedByID           uint
		CreatedByFirstNameEn  string
		CreatedByLastNameEn   string
		CreatedByEmail        string
		CreatedAt             string
	}

	var row closureScan
	err := r.db.Table("project_closures").
		Select(`
			project_closures.id,
			project_closures.overall_project_outcome,
			project_closures.closure_reason,
			project_closures.closure_comments,
			project_closures.lesson_learned,
			project_closures.approver_id,
			project_closures.is_approved,
			approver.first_name_en AS approver_first_name_en,
			approver.last_name_en AS approver_last_name_en,
			approver.email AS approver_email,
			project_closures.created_by_id,
			created_user.first_name_en AS created_by_first_name_en,
			created_user.last_name_en AS created_by_last_name_en,
			created_user.email AS created_by_email,
			CONVERT(VARCHAR, project_closures.created_at, 120) AS created_at
		`).
		Joins("LEFT JOIN users AS approver ON approver.id = project_closures.approver_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_closures.created_by_id").
		Where("project_closures.project_id = ?", projectID).
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	closureDTO := &dtos.ProjectClosureReportDTO{
		ID:                    row.ID,
		OverallProjectOutcome: row.OverallProjectOutcome,
		ClosureReason:         row.ClosureReason,
		ClosureComments:       row.ClosureComments,
		LessonLearned:         row.LessonLearned,
		ApproverID:            row.ApproverID,
		IsApproved:            row.IsApproved,
		CreatedAtFormatted:    row.CreatedAt,
		CreatedByUser: organizationDTOs.UserMinimalDTO{
			ID:          int(row.CreatedByID),
			FirstNameEN: row.CreatedByFirstNameEn,
			LastNameEN:  row.CreatedByLastNameEn,
			Email:       row.CreatedByEmail,
		},
	}

	if row.ApproverFirstNameEn != nil {
		closureDTO.ApproverUser = &organizationDTOs.UserMinimalDTO{
			ID:          int(row.ApproverID),
			FirstNameEN: *row.ApproverFirstNameEn,
			LastNameEN:  *row.ApproverLastNameEn,
			Email:       *row.ApproverEmail,
		}
	}

	return closureDTO, nil
}

func (r *ProjectReportRepository) GetAllReports(params *dtos.ProjectReportListQueryParamsDTO) (*dtos.ProjectReportListResponseDTO, error) {
	var reports []*dtos.ProjectReportListDTO
	var totalCount int64

	query := r.db.Table("project_reports").
		Select(`
			project_reports.id,
			project_reports.title,
			CONVERT(VARCHAR, project_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, project_reports.end_date, 23) AS end_date,
			project_reports.created_at,
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
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_reports.created_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(project_reports.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "project_reports.id"
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

	return &dtos.ProjectReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (r *ProjectReportRepository) GetProjectsByIDs(ids []uint, tx ...utility.DatabaseTransaction) ([]*models.Project, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var projects []*models.Project
	err := db.
		Preload("ProjectType").
		Preload("ProjectStatus").
		Preload("ProjectPhase").
		Preload("ProjectClassification").
		Preload("ResponsibleUser").
		Preload("IncidentActionPlan").
		Where("id IN ?", ids).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectReportRepository) GetProjectsByDateRange(start time.Time, end time.Time, tx ...utility.DatabaseTransaction) ([]*models.Project, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var projects []*models.Project
	err := db.
		Preload("ProjectType").
		Preload("ProjectStatus").
		Preload("ProjectPhase").
		Preload("ProjectClassification").
		Preload("ResponsibleUser").
		Preload("IncidentActionPlan").
		Where("start_date >= ? AND start_date <= ?", start, end).
		Find(&projects).Error
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (r *ProjectReportRepository) CreateReport(report *models.ProjectReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(report).Error; err != nil {
		return 0, err
	}
	return report.ID, nil
}

func (r *ProjectReportRepository) CreateSection(section *models.ProjectReportSection, tx ...utility.DatabaseTransaction) (uint, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(section).Error; err != nil {
		return 0, err
	}
	return section.ID, nil
}

func (r *ProjectReportRepository) CreateReportProject(project *models.ProjectReportProject, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(project).Error
}

func (r *ProjectReportRepository) CreateProjectDocuments(docs []*models.ProjectReportDocument, tx ...utility.DatabaseTransaction) error {
	if len(docs) == 0 {
		return nil
	}
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&docs).Error
}

func (r *ProjectReportRepository) GetReportByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.ProjectReport, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var report models.ProjectReport
	err := db.Where("id = ?", reportID).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ProjectReportRepository) GetReportDetailsByID(reportID uint, tx ...utility.DatabaseTransaction) (*models.ProjectReport, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var report models.ProjectReport
	err := db.
		Preload("CreatedByUser").
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_num ASC")
		}).
		Preload("Projects").
		Preload("Projects.Documents").
		Preload("Projects.Project").
		Preload("Projects.Project.ProjectStatus").
		Preload("Projects.Project.ProjectPhase").
		Preload("Projects.Project.IncidentActionPlan").
		Preload("Projects.Project.IncidentActionPlan.ActionPlanStatus").
		Preload("Projects.Project.IncidentActionPlan.ResponsibleUser").
		Preload("Projects.Project.IncidentActionPlan.IncidentRegister").
		Where("id = ?", reportID).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (r *ProjectReportRepository) GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.ProjectReportSection, error) {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section models.ProjectReportSection
	err := db.Where("id = ?", sectionID).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (r *ProjectReportRepository) UpdateSection(section *models.ProjectReportSection, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (r *ProjectReportRepository) DeleteReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.ProjectReportSection{}).Error
}

func (r *ProjectReportRepository) ReorderReportSections(reportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Exec(`
		UPDATE project_report_sections
		SET order_num = order_num - 1
		WHERE project_report_id = ? AND order_num > ?`,
		reportID, deletedOrderNum).Error
}

func (r *ProjectReportRepository) DeleteReportDocuments(reportID uint, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("project_report_project_id IN (SELECT id FROM project_report_projects WHERE project_report_id = ?)", reportID).
		Delete(&models.ProjectReportDocument{}).Error
}

func (r *ProjectReportRepository) DeleteReportProjects(reportID uint, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("project_report_id = ?", reportID).Delete(&models.ProjectReportProject{}).Error
}

func (r *ProjectReportRepository) DeleteReportSections(reportID uint, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("project_report_id = ?", reportID).Delete(&models.ProjectReportSection{}).Error
}

func (r *ProjectReportRepository) DeleteReport(reportID uint, tx ...utility.DatabaseTransaction) error {
	db := r.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", reportID).Delete(&models.ProjectReport{}).Error
}
