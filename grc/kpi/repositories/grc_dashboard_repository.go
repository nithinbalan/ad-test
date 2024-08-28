package repositories

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type GRCDashboardRepositoryInteface interface {
	GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error)
	GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error)
	GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error)
}

type GRCDashboardRepository struct {
	db *gorm.DB
}

func NewGRCDashboardRepository(db *gorm.DB) *GRCDashboardRepository {
	return &GRCDashboardRepository{
		db: db,
	}
}

func (repo *GRCDashboardRepository) GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error) {
	summary := &kpiDTOs.GRCComplianceSummaryDTO{}

	var totalComplianceRequirements, compliantRequirements int64
	var totalActionPlans, completedActionPlans int64

	// Calculate total number of compliance requirements
	complianceRequirementQuery := repo.db.Model(&models.ComplianceRequirement{})
	if params.DateFrom != "" && params.DateTo != "" {
		complianceRequirementQuery = complianceRequirementQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err := complianceRequirementQuery.Count(&totalComplianceRequirements).Error
	if err != nil {
		return nil, err
	}

	// Calculate number of compliant requirements
	compliantQuery := repo.db.Model(&models.ComplianceRequirement{}).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_requirements.compliance_status_id").
		Where("compliance_statuses.type = ?", "compliant")
	if params.DateFrom != "" && params.DateTo != "" {
		compliantQuery = compliantQuery.Where("compliance_requirements.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = compliantQuery.Count(&compliantRequirements).Error
	if err != nil {
		return nil, err
	}

	// Calculate Compliance Requirement Fulfillment Rate
	if totalComplianceRequirements > 0 {
		summary.ComplianceRequirementFulfillmentRate = float64(compliantRequirements) / float64(totalComplianceRequirements) * 100
	}

	// Calculate total number of compliance action plans
	actionPlanQuery := repo.db.Model(&models.ComplianceActionPlan{})
	if params.DateFrom != "" && params.DateTo != "" {
		actionPlanQuery = actionPlanQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = actionPlanQuery.Count(&totalActionPlans).Error
	if err != nil {
		return nil, err
	}

	// Calculate number of completed action plans
	completedQuery := repo.db.Model(&models.ComplianceActionPlan{}).
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = compliance_action_plans.action_plan_status_id").
		Where("action_plan_statuses.type = ?", "completed")
	if params.DateFrom != "" && params.DateTo != "" {
		completedQuery = completedQuery.Where("compliance_action_plans.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = completedQuery.Count(&completedActionPlans).Error
	if err != nil {
		return nil, err
	}

	// Calculate Compliance Action Plan Resolution Rate
	if totalActionPlans > 0 {
		summary.ComplianceActionPlanResolutionRate = float64(completedActionPlans) / float64(totalActionPlans) * 100
	}

	return summary, nil
}

func (repo *GRCDashboardRepository) GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error) {
	summary := &kpiDTOs.GRCRiskSummaryDTO{}

	var totalRisks, closedRisks, openHighRiskRatings int64

	// Calculate total risks
	totalRiskQuery := repo.db.Model(&models.RiskRegister{})
	if params.DateFrom != "" && params.DateTo != "" {
		totalRiskQuery = totalRiskQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err := totalRiskQuery.Count(&totalRisks).Error
	if err != nil {
		return nil, err
	}

	// Calculate closed risks
	closedRiskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_statuses.type = ?", "closed")
	if params.DateFrom != "" && params.DateTo != "" {
		closedRiskQuery = closedRiskQuery.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = closedRiskQuery.Count(&closedRisks).Error
	if err != nil {
		return nil, err
	}

	// Calculate Risk Closure Rate
	if totalRisks > 0 {
		summary.RiskClosureRate = float64(closedRisks) / float64(totalRisks) * 100
	}

	// Calculate count of open high-risk ratings
	highRiskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Joins("JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_matrix_risk_ratings.type = ?", "very_high").
		Where("risk_statuses.type = ?", "open")
	if params.DateFrom != "" && params.DateTo != "" {
		highRiskQuery = highRiskQuery.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = highRiskQuery.Count(&openHighRiskRatings).Error
	if err != nil {
		return nil, err
	}
	summary.CountOfOpenHighRiskRatings = openHighRiskRatings

	return summary, nil
}

func (repo *GRCDashboardRepository) GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error) {
	var totalAudits, completedAudits, openMajorNCCount, totalFindings, closedFindings int64

	// Calculate total number of audits
	totalAuditQuery := repo.db.Model(&models.Audit{})
	if params.DateFrom != "" && params.DateTo != "" {
		totalAuditQuery = totalAuditQuery.Joins("JOIN audit_plans ON audit_plans.id = audits.audit_plan_id").
			Where("audits.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err := totalAuditQuery.Count(&totalAudits).Error
	if err != nil {
		return nil, err
	}

	// Calculate number of completed audits
	completedAuditQuery := repo.db.Model(&models.Audit{}).
		Joins("JOIN audit_plans ON audit_plans.id = audits.audit_plan_id").
		Joins("JOIN audit_statuses ON audit_statuses.id = audit_plans.audit_status_id").
		Where("audit_statuses.type = ?", "completed")
	if params.DateFrom != "" && params.DateTo != "" {
		completedAuditQuery = completedAuditQuery.Where("audits.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = completedAuditQuery.Count(&completedAudits).Error
	if err != nil {
		return nil, err
	}

	// Calculate Audit Completion Rate
	var auditCompletionRate float64
	if totalAudits > 0 {
		auditCompletionRate = float64(completedAudits) / float64(totalAudits) * 100
	}

	// Count of Open Major Non-Conformities (NCs)
	openMajorNCQuery := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN audit_finding_categories ON audit_finding_categories.id = audit_findings.audit_finding_category_id").
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Where("audit_finding_categories.type = ?", "major-nc").
		Where("audit_finding_statuses.is_closed = ?", false)

	if params.DateFrom != "" && params.DateTo != "" {
		openMajorNCQuery = openMajorNCQuery.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err = openMajorNCQuery.Count(&openMajorNCCount).Error
	if err != nil {
		return nil, err
	}

	// Calculate total number of findings
	totalFindingQuery := repo.db.Model(&models.AuditFinding{})
	if params.DateFrom != "" && params.DateTo != "" {
		totalFindingQuery = totalFindingQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = totalFindingQuery.Count(&totalFindings).Error
	if err != nil {
		return nil, err
	}

	// Calculate number of closed findings
	closedFindingQuery := repo.db.Model(&models.AuditFinding{}).
		Joins("JOIN audit_finding_statuses ON audit_finding_statuses.id = audit_findings.audit_finding_status_id").
		Where("audit_finding_statuses.is_closed = ?", true)
	if params.DateFrom != "" && params.DateTo != "" {
		closedFindingQuery = closedFindingQuery.Where("audit_findings.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = closedFindingQuery.Count(&closedFindings).Error
	if err != nil {
		return nil, err
	}

	// Calculate Audit Finding Resolution Rate
	var auditFindingResolutionRate float64
	if totalFindings > 0 {
		auditFindingResolutionRate = float64(closedFindings) / float64(totalFindings) * 100
	}

	return &kpiDTOs.GRCAuditSummaryDTO{
		AuditCompletionRate:        auditCompletionRate,
		OpenMajorNCCount:           openMajorNCCount,
		AuditFindingResolutionRate: auditFindingResolutionRate,
	}, nil
}
