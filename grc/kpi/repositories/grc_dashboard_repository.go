package repositories

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type GRCDashboardRepositoryInteface interface {
	GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error)
	GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error)
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
