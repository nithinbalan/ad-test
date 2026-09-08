package repositories

import (
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
	GetAllActionPlans(projectID uint) ([]*dtos.ProjectActionPlanListDTO, error)
	GetAllChangeRequests(projectID uint) ([]*dtos.ProjectChangeRequestListDTO, error)
	GetAllDecisionLogs(projectID uint) ([]*dtos.ProjectDecisionLogListDTO, error)
	GetAllSupportingDocuments(projectID uint) ([]*dtos.ProjectSupportingDocumentListDTO, error)
	GetProjectClosure(projectID uint) (*dtos.ProjectClosureReportDTO, error)
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
			compliance_controls.name,
			compliance_controls.description,
			compliance_controls.control_number,
			compliance_controls.created_at,
			compliance_controls.updated_at,
			compliance_control_types.title_en AS control_type_title_en,
			compliance_control_types.title_ar AS control_type_title_ar,
			compliance_control_classifications.title_en AS control_classification_title_en,
			compliance_control_classifications.title_ar AS control_classification_title_ar,
			compliance_control_priorities.title_en AS control_priority_title_en,
			compliance_control_priorities.title_ar AS control_priority_title_ar,
			compliance_control_phases.title_en AS control_phase_title_en,
			compliance_control_phases.title_ar AS control_phase_title_ar,
			compliance_control_maturity_levels.title_en AS control_maturity_level_title_en,
			compliance_control_maturity_levels.title_ar AS control_maturity_level_title_ar,
			compliance_control_enforcement_levels.title_en AS control_enforcement_level_title_en,
			compliance_control_enforcement_levels.title_ar AS control_enforcement_level_title_ar,
			organization_units.title_en AS control_owner_title_en,
			organization_units.title_ar AS control_owner_title_ar
		`).
		Joins("INNER JOIN compliance_controls ON compliance_controls.id = project_controls.control_id").
		Joins("LEFT JOIN compliance_control_types ON compliance_control_types.id = compliance_controls.control_type_id").
		Joins("LEFT JOIN compliance_control_classifications ON compliance_control_classifications.id = compliance_controls.control_classification_id").
		Joins("LEFT JOIN compliance_control_priorities ON compliance_control_priorities.id = compliance_controls.control_priority_id").
		Joins("LEFT JOIN compliance_control_phases ON compliance_control_phases.id = compliance_controls.control_phase_id").
		Joins("LEFT JOIN compliance_control_maturity_levels ON compliance_control_maturity_levels.id = compliance_controls.control_maturity_level_id").
		Joins("LEFT JOIN compliance_control_enforcement_levels ON compliance_control_enforcement_levels.id = compliance_controls.control_enforcement_level_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = compliance_controls.control_owner_id").
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
			project_risks.risk_id,
			risk_registers.reference_number,
			risk_registers.subject,
			risk_registers.description,
			risk_registers.inherent_risk_score,
			risk_registers.residual_risk_score,
			risk_registers.created_at,
			risk_registers.updated_at,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_types.title_en AS risk_type_title_en,
			risk_types.title_ar AS risk_type_title_ar,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar
		`).
		Joins("INNER JOIN risk_registers ON risk_registers.id = project_risks.risk_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_registers.risk_status_id").
		Joins("LEFT JOIN risk_types ON risk_types.id = risk_registers.risk_type_id").
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_registers.risk_category_id").
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
			kpi_projects.kpi_id,
			kpi_registers.reference_number,
			kpi_registers.title_en,
			kpi_registers.title_ar,
			kpi_registers.description_en,
			kpi_registers.description_ar,
			kpi_registers.current_value,
			kpi_registers.target_value,
			kpi_registers.created_at,
			kpi_registers.updated_at
		`).
		Joins("INNER JOIN kpi_registers ON kpi_registers.id = kpi_projects.kpi_id").
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
			project_change_requests.scope_impact,
			project_change_requests.schedule_impact,
			project_change_requests.cost_impact,
			project_change_requests.requested_by_user_id,
			project_change_requests.is_assigned,
			project_change_requests.workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			project_change_requests.created_at,
			users.id AS requested_by_user_id,
			users.first_name_en AS requested_by_user_first_name_en,
			users.first_name_ar AS requested_by_user_first_name_ar,
			users.last_name_en AS requested_by_user_last_name_en,
			users.last_name_ar AS requested_by_user_last_name_ar
		`).
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = project_change_requests.workflow_status_id").
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
			project_decision_logs.is_assigned,
			project_decision_logs.workflow_status_id,
			workflow_statuses.title_en AS workflow_status_title_en,
			workflow_statuses.title_ar AS workflow_status_title_ar,
			workflow_statuses.color_code AS workflow_status_color_code,
			workflow_statuses.type AS workflow_status_type,
			project_decision_logs.created_at,
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
		return nil, nil
	}

	closureDTO := &dtos.ProjectClosureReportDTO{
		ID:                    row.ID,
		OverallProjectOutcome: row.OverallProjectOutcome,
		ClosureReason:         row.ClosureReason,
		ClosureComments:       row.ClosureComments,
		LessonLearned:         row.LessonLearned,
		ApproverID:            row.ApproverID,
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
