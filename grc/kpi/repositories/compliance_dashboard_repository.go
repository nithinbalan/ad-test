package repositories

import (
	"fmt"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceDashboardRepositoryInterface interface {
	GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error)
	GetLeadershipAttentionSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceLeadershipAttentionDTO, error)
	GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error)
	GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetLeastPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error)
	GetPolicyOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyOverviewDTO, error)
	GetPolicyStatusBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyStatusBreakdownDTO, error)
	GetPolicyByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePolicyByOrganizationUnitDTO, error)
	GetControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error)
	GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceEfficiencySummaryDTO, error)
	GetControlOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlOverviewDTO, error)
	GetOrganizationUnitWiseAssessmentCoverage(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceAssessmentControlOrganizationUnit, error)
	GetActionPlanOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanOverviewDTO, error)
	GetActionPlanStatusSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanStatusSummaryDTO, error)
	GetActionPlansByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO, error)
}

type ComplianceDashboardRepository struct {
	db *gorm.DB
}

func NewComplianceDashboardRepository(db *gorm.DB) *ComplianceDashboardRepository {
	return &ComplianceDashboardRepository{
		db: db,
	}
}

func (repo *ComplianceDashboardRepository) GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error) {
	overview := &kpiDTOs.ComplianceOverviewDTO{}

	registerQuery := repo.db.Model(&models.ComplianceRegister{}).
		Select(`
			COUNT(*) AS total_compliance_registers,
			COALESCE(AVG(score), 0) AS organization_compliance_score
		`)

	if params.DateFrom != "" && params.DateTo != "" {
		registerQuery = registerQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := registerQuery.Scan(&overview).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch compliance registers: %v", err)
	}

	requirementQuery := repo.db.Model(&models.ComplianceRequirement{})
	if params.DateFrom != "" && params.DateTo != "" {
		requirementQuery = requirementQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := requirementQuery.Count(&overview.TotalRequirements).Error; err != nil {
		return nil, err
	}

	assessmentQuery := repo.db.Model(&models.ComplianceAssessment{})
	if params.DateFrom != "" && params.DateTo != "" {
		assessmentQuery = assessmentQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := assessmentQuery.Count(&overview.TotalAssessments).Error; err != nil {
		return nil, err
	}

	return overview, nil
}

func (repo *ComplianceDashboardRepository) GetLeadershipAttentionSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceLeadershipAttentionDTO, error) {
	out := &kpiDTOs.ComplianceLeadershipAttentionDTO{}

	// Open Critical Requirements
	openQuery := repo.db.Table("compliance_requirements crq").
		Joins("LEFT JOIN severity_levels sl ON sl.id = crq.severity_level_id").
		Joins("LEFT JOIN compliance_statuses cs ON cs.id = crq.compliance_status_id").
		Where("sl.type = ?", "high").
		Where("(cs.id IS NULL OR cs.status <> ?)", "closed")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		openQuery = openQuery.Where("crq.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := openQuery.Count(&out.OpenCriticalRequirements).Error; err != nil {
		return nil, err
	}

	// Overdue Compliance Requirements (deadline passed and not closed)
	overdueQuery := repo.db.Table("compliance_requirements crq").
		Joins("LEFT JOIN compliance_statuses cs ON cs.id = crq.compliance_status_id").
		Where("crq.deadline IS NOT NULL").
		Where("crq.deadline < CAST(GETDATE() AS date)").
		Where("(cs.id IS NULL OR cs.status <> ?)", "closed")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		overdueQuery = overdueQuery.Where("crq.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := overdueQuery.Count(&out.OverdueComplianceRequirements).Error; err != nil {
		return nil, err
	}

	// Unassigned Compliance Items
	unassignedQuery := repo.db.Table("compliance_requirements crq").
		Where(`
		NOT EXISTS (
		SELECT 1 FROM compliance_requirement_responsible_units crru WHERE crru.compliance_requirement_id = crq.id
		)
		`)

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		unassignedQuery = unassignedQuery.Where("crq.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := unassignedQuery.Count(&out.UnassignedComplianceItems).Error; err != nil {
		return nil, err
	}

	// Past Due Actions (open action plans with target_date in the past)
	pastDueQuery := repo.db.Table("compliance_action_plans cap").
		Joins("JOIN action_plan_statuses aps ON aps.id = cap.action_plan_status_id").
		Where("aps.is_closed = ?", false).
		Where("cap.target_date IS NOT NULL AND cap.target_date < CAST(GETDATE() AS date)")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		pastDueQuery = pastDueQuery.Where("cap.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := pastDueQuery.Count(&out.PastDueActions).Error; err != nil {
		return nil, err
	}

	return out, nil
}

func (repo *ComplianceDashboardRepository) GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error) {
	complianceStatuses := []kpiDTOs.ComplianceSummaryItemDTO{}
	var total int64

	query := repo.db.Model(&models.ComplianceRegister{}).
		Select(`
            compliance_statuses.id,
            compliance_statuses.title_en,
            compliance_statuses.title_ar,
            compliance_statuses.color_code,
            compliance_statuses.type,
            COUNT(*) AS count
        `).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id").
		Group(`
            compliance_statuses.id, 
            compliance_statuses.title_en, 
            compliance_statuses.title_ar, 
            compliance_statuses.color_code, 
            compliance_statuses.type
        `).
		Order("compliance_statuses.id")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&complianceStatuses).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&models.ComplianceRegister{}).Count(&total).Error; err != nil {
		return nil, err
	}

	// percentage
	for i := range complianceStatuses {
		if total > 0 {
			complianceStatuses[i].Percent = (float64(complianceStatuses[i].Count) / float64(total)) * 100
		} else {
			complianceStatuses[i].Percent = 0
		}
	}

	result := &kpiDTOs.ComplianceBreakdownDTO{
		ComplianceStatuses: complianceStatuses,
		Total:              total,
	}

	return result, nil
}

func (repo *ComplianceDashboardRepository) GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error) {
	results := []kpiDTOs.CompliancePerformanceRadarDTO{}

	query := repo.db.Model(&models.ComplianceRegister{}).
		Select(`
            compliance_register.standard_code AS compliance_standard_code,
			compliance_register.title_en AS compliance_title_en,
			compliance_register.title_ar AS compliance_title_ar,
            AVG(COALESCE(compliance_register.score, 0)) AS score
        `).
		Where("compliance_register.score IS NOT NULL").
		Group("compliance_register.standard_code, compliance_register.title_en, compliance_register.title_ar").
		Having("AVG(compliance_register.score) > ?", 75).
		Order("score DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ComplianceDashboardRepository) GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePerformanceRadarDTO, error) {
	results := []kpiDTOs.CompliancePerformanceRadarDTO{}

	query := repo.db.Model(&models.ComplianceRegister{}).
		Select(`
            compliance_register.standard_code AS compliance_standard_code,
			compliance_register.title_en AS compliance_title_en,
			compliance_register.title_ar AS compliance_title_ar,
            AVG(COALESCE(compliance_register.score, 0)) AS score
        `).
		Where("compliance_register.score IS NOT NULL").
		Group("compliance_register.standard_code, compliance_register.title_en, compliance_register.title_ar").
		Having("AVG(compliance_register.score) <= ?", 75).
		Order("score ASC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ComplianceDashboardRepository) GetHighPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error) {
	results := []kpiDTOs.ComplianceClausePerformanceDTO{}

	query := repo.db.Table("compliance_clauses AS child").
		Select(`
			parent.clause_no AS compliance_clause_no,
            parent.title_en AS compliance_clause_title_en,
			parent.title_ar as compliance_clause_title_ar,
            AVG(CAST(child.score AS FLOAT)) AS compliance_clause_score
        `).
		Joins("JOIN compliance_clauses AS parent ON parent.id = child.compliance_clause_id").
		Where("child.compliance_clause_id IS NOT NULL").
		Where("child.score IS NOT NULL").
		Group("parent.clause_no, parent.title_en, parent.title_ar").
		Order("compliance_clause_score DESC").
		Limit(5)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_clauses.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *ComplianceDashboardRepository) GetLeastPerformanceComplianceClauses(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceClausePerformanceDTO, error) {
	results := []kpiDTOs.ComplianceClausePerformanceDTO{}

	query := repo.db.Table("compliance_clauses AS child").
		Select(`
			parent.clause_no AS compliance_clause_no,
            parent.title_en AS compliance_clause_title_en,
			parent.title_ar as compliance_clause_title_ar,
            AVG(CAST(child.score AS FLOAT)) AS compliance_clause_score
        `).
		Joins("JOIN compliance_clauses AS parent ON parent.id = child.compliance_clause_id").
		Where("child.compliance_clause_id IS NOT NULL").
		Where("child.score IS NOT NULL").
		Group("parent.clause_no, parent.title_en, parent.title_ar").
		Order("compliance_clause_score ASC").
		Limit(5)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_clauses.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *ComplianceDashboardRepository) GetPolicyOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyOverviewDTO, error) {
	out := &kpiDTOs.CompliancePolicyOverviewDTO{}

	// total policies
	totalQuery := repo.db.Table("document_compliance_requirements dcr").
		Distinct("dcr.document_id")

	if params.DateFrom != "" && params.DateTo != "" {
		totalQuery = totalQuery.Where("dcr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := totalQuery.Count(&out.TotalPolicies).Error; err != nil {
		return nil, err
	}

	// Acknowledged policies
	var acknowledged int64
	acknowledgedQuery := repo.db.Table("document_compliance_requirements dcr").
		Joins("JOIN document_acknowledge_users dcu ON dcr.document_id = dcu.document_id").
		Distinct("dcr.document_id")

	if params.DateFrom != "" && params.DateTo != "" {
		acknowledgedQuery = acknowledgedQuery.Where("dcr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := acknowledgedQuery.Count(&acknowledged).Error; err != nil {
		return nil, err
	}

	if out.TotalPolicies > 0 {
		out.AcknowledgementScore = float64(acknowledged) / float64(out.TotalPolicies) * 100
	} else {
		out.AcknowledgementScore = 0
	}

	// Determine the highest available risk rating dynamically
	var highestRating string
	err := repo.db.Table("risk_matrix_risk_ratings").
		Select("type").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &highestRating).Error
	if err != nil || highestRating == "" {
		return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
	}

	// High-risk policies
	highRiskQuery := repo.db.Table("document_compliance_requirements dcr").
		Joins("JOIN risk_associated_documents rad ON rad.document_id = dcr.document_id").
		Joins("JOIN risk_register rr ON rr.id = rad.risk_register_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings rmrr ON rmrr.id = rr.current_risk_rating_id").
		Where("rmrr.type = ?", highestRating).
		Distinct("dcr.document_id")

	if params.DateFrom != "" && params.DateTo != "" {
		highRiskQuery = highRiskQuery.Where("dcr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := highRiskQuery.Count(&out.HighRiskPolicies).Error; err != nil {
		return nil, err
	}

	// 4) Total policy exceptions raised for mapped policies
	policyExceptionsQuery := repo.db.Table("policy_exceptions pe").
		Where("pe.document_id IN (?)",
			repo.db.Table("document_compliance_requirements").Select("document_id"),
		)

	if params.DateFrom != "" && params.DateTo != "" {
		policyExceptionsQuery = policyExceptionsQuery.Where("pe.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := policyExceptionsQuery.Count(&out.PolicyExceptionsRaised).Error; err != nil {
		return nil, err
	}

	// Critical, unacknowledged policies
	criticalUnacknowledgedQuery := repo.db.Table("document_compliance_requirements dcr").
		Joins("JOIN risk_associated_documents rad ON rad.document_id = dcr.document_id").
		Joins("JOIN risk_register rr ON rr.id = rad.risk_register_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings rmrr ON rmrr.id = rr.current_risk_rating_id").
		Joins("LEFT JOIN document_acknowledge_users dcu ON dcu.document_id = dcr.document_id").
		Where("dcu.id IS NULL").
		Where("rmrr.type = ?", highestRating).
		Distinct("dcr.document_id")

	if params.DateFrom != "" && params.DateTo != "" {
		criticalUnacknowledgedQuery = criticalUnacknowledgedQuery.Where("dcr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := criticalUnacknowledgedQuery.Count(&out.CriticalPoliciesUnacknowledged).Error; err != nil {
		return nil, err
	}

	return out, nil
}

func (repo *ComplianceDashboardRepository) GetPolicyByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.CompliancePolicyByOrganizationUnitDTO, error) {
	results := []kpiDTOs.CompliancePolicyByOrganizationUnitDTO{}

	query := repo.db.Table("document_compliance_requirements AS dcr").
		Select(`
			parent_unit.title_en AS organization_unit_title_en,
			parent_unit.title_ar AS organization_unit_title_ar,
			parent_unit.color_code AS organization_unit_color_code,
			COUNT(DISTINCT dcr.document_id) AS total_policies
		`).
		Joins("JOIN documents doc ON doc.id = dcr.document_id").
		Joins("LEFT JOIN organization_units child_unit ON child_unit.id = doc.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(child_unit.organization_unit_id, child_unit.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = parent_unit.organization_level_id AND ol.is_primary = ?", true).
		Where("parent_unit.organization_unit_id IS NULL")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("dcr.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Group("parent_unit.title_en, parent_unit.title_ar, parent_unit.color_code").
		Order("total_policies DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ComplianceDashboardRepository) GetPolicyStatusBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePolicyStatusBreakdownDTO, error) {
	result := &kpiDTOs.CompliancePolicyStatusBreakdownDTO{}

	mappedDocs := repo.db.
		Table("document_compliance_requirements dcr").
		Select("DISTINCT dcr.document_id")

	query := repo.db.
		Table("(?) AS md", mappedDocs).
		Select(`
			COUNT(DISTINCT d.id) AS total_policies,
			SUM(CASE WHEN ds.status = 'active' THEN 1 ELSE 0 END) AS active_policies,
			SUM(CASE WHEN ds.status = 'inactive' THEN 1 ELSE 0 END) AS inactive_policies
		`).
		Joins("JOIN documents d ON d.id = md.document_id").
		Joins("LEFT JOIN document_statuses ds ON ds.id = d.document_status_id")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("d.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ComplianceDashboardRepository) GetControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error) {
	summary := &kpiDTOs.ComplianceControlSummaryDTO{}
	var efficientControlsCount int64

	countQuery := repo.db.Model(&models.ComplianceControl{}).
		Select(`
			COUNT(DISTINCT controls.id) AS total_controls,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'efficient' THEN controls.id END) AS efficient_controls
		`).
		Joins("JOIN controls ON controls.id = compliance_controls.control_id").
		Joins("LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id")

	if params.DateFrom != "" && params.DateTo != "" {
		countQuery = countQuery.Where("compliance_controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := countQuery.Row().Scan(&summary.TotalControls, &efficientControlsCount); err != nil {
		return nil, err
	}

	implementationStatuses := []kpiDTOs.ComplianceSummaryItemDTO{}
	implementationQuery := repo.db.Model(&models.ComplianceControl{}).
		Select(`
			control_implementation_statuses.id, 
			control_implementation_statuses.title_en, 
			control_implementation_statuses.title_ar, 
			control_implementation_statuses.color_code, 
			control_implementation_statuses.type, 
			COUNT(DISTINCT compliance_controls.control_id) AS count
		`).
		Joins("JOIN controls ON controls.id = compliance_controls.control_id").
		Joins("JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Group(`
			control_implementation_statuses.id, 
			control_implementation_statuses.title_en, 
			control_implementation_statuses.title_ar, 
			control_implementation_statuses.color_code, 
			control_implementation_statuses.type
		`).
		Order("control_implementation_statuses.id")
	if params.DateFrom != "" && params.DateTo != "" {
		implementationQuery = implementationQuery.Where("compliance_controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := implementationQuery.Scan(&implementationStatuses).Error; err != nil {
		return nil, err
	}
	for i := range implementationStatuses {
		if summary.TotalControls > 0 {
			implementationStatuses[i].Percent = float64(implementationStatuses[i].Count) / float64(summary.TotalControls) * 100
		}
	}
	summary.ImplementationStatuses = implementationStatuses

	efficiencyStatuses := []kpiDTOs.ComplianceSummaryItemDTO{}
	efficiencyQuery := repo.db.Model(&models.ComplianceControl{}).
		Select(`
			control_efficiency_statuses.id, 
			control_efficiency_statuses.title_en, 
			control_efficiency_statuses.title_ar, 
			control_efficiency_statuses.color_code, 
			control_efficiency_statuses.type, 
			COUNT(DISTINCT compliance_controls.control_id) AS count
		`).
		Joins("JOIN controls ON controls.id = compliance_controls.control_id").
		Joins("JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id").
		Group(`
			control_efficiency_statuses.id, 
			control_efficiency_statuses.title_en, 
			control_efficiency_statuses.title_ar, 
			control_efficiency_statuses.color_code, 
			control_efficiency_statuses.type
		`).
		Order("control_efficiency_statuses.id")

	if params.DateFrom != "" && params.DateTo != "" {
		efficiencyQuery = efficiencyQuery.Where("compliance_controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := efficiencyQuery.Scan(&efficiencyStatuses).Error; err != nil {
		return nil, err
	}
	for i := range efficiencyStatuses {
		if summary.TotalControls > 0 {
			efficiencyStatuses[i].Percent = float64(efficiencyStatuses[i].Count) / float64(summary.TotalControls) * 100
		}
	}
	summary.EfficiencyChart = efficiencyStatuses

	// overall efficiency
	if summary.TotalControls > 0 {
		summary.OverallEfficiency = (float64(efficientControlsCount) / float64(summary.TotalControls)) * 100
	} else {
		summary.OverallEfficiency = 0
	}

	return summary, nil
}

func (repo *ComplianceDashboardRepository) GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceEfficiencySummaryDTO, error) {

	var results []kpiDTOs.ComplianceEfficiencySummaryQueryResultDTO

	// Base query
	query := repo.db.Table("compliance_controls").
		Select(`
            compliance_register.standard_code AS compliance_standard_code,
            compliance_register.title_en AS compliance_title_en,
            compliance_register.title_ar AS compliance_title_ar,
            control_efficiency_statuses.id AS status_id,
            control_efficiency_statuses.title_en,
            control_efficiency_statuses.title_ar,
            control_efficiency_statuses.color_code,
            control_efficiency_statuses.type,
            COUNT(*) AS count
        `).
		Joins("JOIN controls ON compliance_controls.control_id = controls.id").
		Joins("JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Joins("JOIN compliance_register ON compliance_controls.compliance_register_id = compliance_register.id").
		Joins("JOIN control_implementation_statuses cis ON cis.id = controls.control_implementation_status_id").
		Where("cis.type = ?", "implemented")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"compliance_controls.created_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	query = query.Group(`
        compliance_register.standard_code,
        compliance_register.title_en,
        compliance_register.title_ar,
        control_efficiency_statuses.id,
        control_efficiency_statuses.title_en,
        control_efficiency_statuses.title_ar,
        control_efficiency_statuses.color_code,
        control_efficiency_statuses.type
    `).Order("compliance_register.standard_code")

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	// Map to group by compliance standard
	grouped := make(map[string]*kpiDTOs.ComplianceEfficiencySummaryDTO)

	for _, r := range results {
		item, exists := grouped[r.ComplianceStandardCode]
		if !exists {
			item = &kpiDTOs.ComplianceEfficiencySummaryDTO{
				ComplianceStandardCode: r.ComplianceStandardCode,
				ComplianceTitleEn:      r.ComplianceTitleEn,
				ComplianceTitleAr:      r.ComplianceTitleAr,
				EfficiencyStatuses:     make([]kpiDTOs.ComplianceSummaryItemDTO, 0),
			}
			grouped[r.ComplianceStandardCode] = item
		}

		item.EfficiencyStatuses = append(item.EfficiencyStatuses, kpiDTOs.ComplianceSummaryItemDTO{
			ID:        r.StatusID,
			TitleEn:   r.TitleEn,
			TitleAr:   r.TitleAr,
			ColorCode: r.ColorCode,
			Type:      r.Type,
			Count:     r.Count,
		})
	}

	// Prepare final result slice
	summaryList := make([]kpiDTOs.ComplianceEfficiencySummaryDTO, 0, len(grouped))

	for _, item := range grouped {
		var total int64
		for _, s := range item.EfficiencyStatuses {
			total += s.Count
		}

		if total > 0 {
			for i := range item.EfficiencyStatuses {
				item.EfficiencyStatuses[i].Percent = float64(item.EfficiencyStatuses[i].Count) / float64(total) * 100
			}
		}

		summaryList = append(summaryList, *item)
	}

	return summaryList, nil
}

func (repo *ComplianceDashboardRepository) GetControlOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlOverviewDTO, error) {
	result := &kpiDTOs.ComplianceControlOverviewDTO{}

	// date filter
	applyDateFilter := func(tx *gorm.DB, table string) *gorm.DB {
		if params.DateFrom != "" && params.DateTo != "" {
			tx = tx.Where(fmt.Sprintf("%s.reported_at BETWEEN ? AND ?", table), params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
		}
		return tx
	}

	// 1) Total assessments
	totalAssessmentsQuery := repo.db.Model(&models.ComplianceAssessmentControl{})
	totalAssessmentsQuery = applyDateFilter(totalAssessmentsQuery, "compliance_assessment_controls")

	if err := totalAssessmentsQuery.Count(&result.TotalAssessments).Error; err != nil {
		return nil, err
	}

	// 2) Assessment coverage
	var controlsAssessed, controlsMapped int64
	assessedQuery := repo.db.Model(&models.ComplianceAssessmentControl{}).Distinct("compliance_control_id")
	assessedQuery = applyDateFilter(assessedQuery, "compliance_assessment_controls")

	if err := assessedQuery.Count(&controlsAssessed).Error; err != nil {
		return nil, err
	}

	mappedQuery := repo.db.Model(&models.ComplianceControl{}).Distinct("id")
	mappedQuery = applyDateFilter(mappedQuery, "compliance_controls")

	if err := mappedQuery.Count(&controlsMapped).Error; err != nil {
		return nil, err
	}

	if controlsMapped > 0 {
		result.AssessmentCoverage = float64(controlsAssessed) / float64(controlsMapped) * 100
	}

	// 3) Completion percentage
	var completedAC int64
	completedQuery := repo.db.Model(&models.ComplianceAssessmentControl{}).Where("is_completed = ?", true)
	completedQuery = applyDateFilter(completedQuery, "compliance_assessment_controls")

	if err := completedQuery.Count(&completedAC).Error; err != nil {
		return nil, err
	}

	if result.TotalAssessments > 0 {
		result.Completion = float64(completedAC) / float64(result.TotalAssessments) * 100
	}

	// 4) High Risk Findings
	var highestRating string
	if err := repo.db.Table("risk_matrix_risk_ratings").
		Select("type").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &highestRating).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch highest risk rating: %v", err)
	}

	var highRiskAC int64
	highRiskQuery := repo.db.Table("compliance_assessment_controls AS cac").
		Joins("JOIN compliance_assessments AS ca ON ca.id = cac.compliance_assessment_id").
		Joins("JOIN risk_compliance_requirements AS rcr ON rcr.compliance_register_id = ca.compliance_register_id").
		Joins("JOIN risk_register AS rr ON rr.id = rcr.risk_register_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS rmrr ON rmrr.id = rr.current_risk_rating_id").
		Where("rmrr.type = ?", highestRating)

	highRiskQuery = applyDateFilter(highRiskQuery, "cac")

	if err := highRiskQuery.Count(&highRiskAC).Error; err != nil {
		return nil, err
	}
	result.HighRiskFindings = highRiskAC

	// 5) Registers Assessed
	registersAssessedXQuery := repo.db.Table("compliance_assessment_controls cac").
		Joins("JOIN compliance_assessments ca ON ca.id = cac.compliance_assessment_id").
		Distinct("ca.compliance_register_id")

	registersAssessedXQuery = applyDateFilter(registersAssessedXQuery, "cac")

	if err := registersAssessedXQuery.Count(&result.RegistersAssessedX).Error; err != nil {
		return nil, err
	}

	registersAssessedYQuery := repo.db.Model(&models.ComplianceRegister{})
	registersAssessedYQuery = applyDateFilter(registersAssessedYQuery, "compliance_registers")

	if err := registersAssessedYQuery.Count(&result.RegistersAssessedY).Error; err != nil {
		return nil, err
	}

	if result.RegistersAssessedY > 0 {
		result.RegistersAssessedPct = float64(result.RegistersAssessedX) / float64(result.RegistersAssessedY) * 100
	}

	// 6) Upcoming Assessments (created in last 14 days)
	upcomingQuery := repo.db.Model(&models.ComplianceAssessmentControl{}).
		Where("is_completed = ?", false).
		Where("created_at >= DATEADD(day, -14, CAST(GETDATE() AS date))")

	upcomingQuery = applyDateFilter(upcomingQuery, "compliance_assessment_controls")

	if err := upcomingQuery.Count(&result.UpcomingAssessments).Error; err != nil {
		return nil, err
	}

	// 7) Pending Remediation
	pendingQuery := repo.db.Table("compliance_action_plans cap").
		Joins("JOIN action_plan_statuses aps ON aps.id = cap.action_plan_status_id").
		Where("cap.compliance_assessment_control_id IS NOT NULL").
		Where("aps.is_closed = ?", false)

	pendingQuery = applyDateFilter(pendingQuery, "cap")

	if err := pendingQuery.Count(&result.PendingRemediation).Error; err != nil {
		return nil, err
	}

	// 8) Departments Not Assessed
	deptQuery := repo.db.Raw(`
		SELECT COUNT(*) FROM (
			SELECT cr.responsible_unit_id AS unit_id
			FROM compliance_register cr
			GROUP BY cr.responsible_unit_id
		) units
		WHERE NOT EXISTS (
			SELECT 1
			FROM compliance_assessment_controls cac
			JOIN compliance_assessments ca ON ca.id = cac.compliance_assessment_id
			JOIN compliance_register cr2 ON cr2.id = ca.compliance_register_id
			WHERE cr2.responsible_unit_id = units.unit_id
		)
	`)

	if err := deptQuery.Scan(&result.DepartmentsNotAssessed).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ComplianceDashboardRepository) GetOrganizationUnitWiseAssessmentCoverage(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceAssessmentControlOrganizationUnit, error) {
	results := []kpiDTOs.ComplianceAssessmentControlOrganizationUnit{}

	query := repo.db.Table("compliance_assessment_controls AS cac").
		Select(`
			parent_unit.title_en AS organization_unit_title_en,
			parent_unit.title_ar AS organization_unit_title_ar,
			parent_unit.color_code AS organization_unit_color_code,
			COUNT(DISTINCT cac.id) AS assessment_count
		`).
		Joins("JOIN compliance_assessments AS ca ON ca.id = cac.compliance_assessment_id").
		Joins("JOIN compliance_register AS cr ON cr.id = ca.compliance_register_id").
		Joins("LEFT JOIN organization_units child_unit ON child_unit.id = cr.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(child_unit.organization_unit_id, child_unit.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = parent_unit.organization_level_id AND ol.is_primary = ?", true)

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("cac.reported_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.
		Group("parent_unit.title_en, parent_unit.title_ar, parent_unit.color_code").
		Order("assessment_count DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *ComplianceDashboardRepository) GetActionPlanOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanOverviewDTO, error) {
	result := &kpiDTOs.ComplianceActionPlanOverviewDTO{}

	query := repo.db.Model(&models.ComplianceActionPlan{}).
		Joins("LEFT JOIN action_plan_statuses aps ON aps.id = compliance_action_plans.action_plan_status_id").
		Joins("LEFT JOIN severity_levels sl ON sl.id = compliance_action_plans.severity_level_id").
		Select(`
			COUNT(*) AS total_action_plans,
			SUM(CASE WHEN aps.type = 'completed' OR aps.is_closed = 1 THEN 1 ELSE 0 END) AS completed_plans,
			SUM(CASE WHEN aps.is_closed = 0 AND aps.type IN ('new', 'in_progress', 'on_hold', 'reopen') THEN 1 ELSE 0 END) AS pending_plans,
			SUM(
				CASE 
					WHEN aps.is_closed = 0 
						AND aps.type IN ('new', 'in_progress', 'on_hold', 'reopen') 
						AND compliance_action_plans.target_date IS NOT NULL 
						AND compliance_action_plans.target_date < CAST(GETDATE() AS date)
					THEN 1 
					ELSE 0 
				END
			) AS overdue_plans,
			SUM(CASE WHEN sl.type = 'high' THEN 1 ELSE 0 END) AS high_priority_plans
		`)

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(result).Error; err != nil {
		return nil, err
	}

	if result.TotalActionPlans > 0 {
		result.AvgCompletionPercent = float64(result.CompletedPlans) / float64(result.TotalActionPlans) * 100
	} else {
		result.AvgCompletionPercent = 0
	}

	return result, nil
}

func (repo *ComplianceDashboardRepository) GetActionPlanStatusSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceActionPlanStatusSummaryDTO, error) {
	summary := &kpiDTOs.ComplianceActionPlanStatusSummaryDTO{}
	actionPlanStatuses := []kpiDTOs.ComplianceSummaryItemDTO{}

	TotalActionPlansQuery := repo.db.Model(&models.ComplianceActionPlan{})
	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		TotalActionPlansQuery = TotalActionPlansQuery.Where("compliance_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if err := TotalActionPlansQuery.Count(&summary.TotalActionPlans).Error; err != nil {
		return nil, err
	}

	statusQuery := repo.db.Table("compliance_action_plans cap").
		Select(`
			aps.title_en AS title_en,
			aps.title_ar AS title_ar,
			aps.color_code AS color_code,
			aps.type AS type,
			COUNT(*) AS count
		`).
		Joins("JOIN action_plan_statuses aps ON aps.id = cap.action_plan_status_id").
		Group("aps.title_en, aps.title_ar, aps.color_code, aps.type").
		Order("aps.title_en ASC")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		statusQuery = statusQuery.Where("cap.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := statusQuery.Scan(&actionPlanStatuses).Error; err != nil {
		return nil, err
	}

	for i := range actionPlanStatuses {
		if summary.TotalActionPlans > 0 {
			actionPlanStatuses[i].Percent = float64(actionPlanStatuses[i].Count) / float64(summary.TotalActionPlans) * 100
		} else {
			actionPlanStatuses[i].Percent = 0
		}
	}
	summary.ActionPlanStatuses = actionPlanStatuses

	return summary, nil
}

func (repo *ComplianceDashboardRepository) GetActionPlansByOrganizationUnit(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) ([]kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO, error) {
	results := []kpiDTOs.ComplianceActionPlanByOrganizationUnitDTO{}

	query := repo.db.Table("compliance_action_plans AS cap").
		Select(`
			ou.title_en AS organization_unit_title_en,
			ou.title_ar AS organization_unit_title_ar,
			ou.color_code AS organization_unit_color_code,
			COUNT(DISTINCT cap.id) AS total_action_plans
		`).
		Joins(`
			JOIN users u ON u.id = cap.responsible_user_id
			OUTER APPLY (
				SELECT TOP 1 ou_sub.*
				FROM organization_units ou_sub
				JOIN organization_levels ol_sub ON ol_sub.id = ou_sub.organization_level_id
				WHERE ou_sub.id IN (
					u.organization_unit_1_id,
					u.organization_unit_2_id,
					u.organization_unit_3_id,
					u.organization_unit_4_id,
					u.organization_unit_5_id
				)
				AND ol_sub.is_primary = 1
			) AS ou
		`).
		Group("ou.title_en, ou.title_ar, ou.color_code").
		Order("total_action_plans DESC")

	if params != nil && params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("cap.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
