package repositories

import (
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"
	"math"

	"gorm.io/gorm"
)

type GRCDashboardRepositoryInteface interface {
	GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error)
	GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error)
	GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error)
	GetGRCDocumentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentSummaryDTO, error)
	GetGRCOrganizationSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCOrganizationSummaryDTO, error)

	GetGRCComplianceOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceOverviewDTO, error)
	GetGRCRiskOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskOverviewDTO, error)

	GetTopComplianceRegisters(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopComplianceRegistersDTO, error)
	GetTopRiskRegisters(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopRiskRegistersDTO, error)
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
		summary.ComplianceRequirementFulfillmentRate = math.Round(float64(compliantRequirements) / float64(totalComplianceRequirements) * 100)
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
		summary.ComplianceActionPlanResolutionRate = math.Round(float64(completedActionPlans) / float64(totalActionPlans) * 100)
	}

	return summary, nil
}

func (repo *GRCDashboardRepository) GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error) {
	summary := &kpiDTOs.GRCRiskSummaryDTO{}

	var totalRisks, closedRisks, openHighRiskRatings int64

	totalRiskQuery := repo.db.Model(&models.RiskRegister{})
	if params.DateFrom != "" && params.DateTo != "" {
		totalRiskQuery = totalRiskQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err := totalRiskQuery.Count(&totalRisks).Error
	if err != nil {
		return nil, err
	}

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

	if totalRisks > 0 {
		summary.RiskClosureRate = math.Round(float64(closedRisks) / float64(totalRisks) * 100)
	}

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

	totalAuditQuery := repo.db.Model(&models.Audit{})
	if params.DateFrom != "" && params.DateTo != "" {
		totalAuditQuery = totalAuditQuery.Joins("JOIN audit_plans ON audit_plans.id = audits.audit_plan_id").
			Where("audits.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err := totalAuditQuery.Count(&totalAudits).Error
	if err != nil {
		return nil, err
	}

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

	var auditCompletionRate float64
	if totalAudits > 0 {
		auditCompletionRate = math.Round(float64(completedAudits) / float64(totalAudits) * 100)
	}

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

	totalFindingQuery := repo.db.Model(&models.AuditFinding{})
	if params.DateFrom != "" && params.DateTo != "" {
		totalFindingQuery = totalFindingQuery.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}
	err = totalFindingQuery.Count(&totalFindings).Error
	if err != nil {
		return nil, err
	}

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

	var auditFindingResolutionRate float64
	if totalFindings > 0 {
		auditFindingResolutionRate = math.Round(float64(closedFindings) / float64(totalFindings) * 100)
	}

	return &kpiDTOs.GRCAuditSummaryDTO{
		AuditCompletionRate:        auditCompletionRate,
		OpenMajorNCCount:           openMajorNCCount,
		AuditFindingResolutionRate: auditFindingResolutionRate,
	}, nil
}

func (repo *GRCDashboardRepository) GetGRCDocumentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentSummaryDTO, error) {
	var totalPolicies int64
	var totalOtherDocuments int64

	policyQuery := repo.db.Model(&models.Document{}).
		Joins("JOIN document_types ON document_types.id = documents.document_type_id").
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id").
		Where("document_statuses.type = ?", "published").
		Where("document_types.type IN ?", []string{"policy", "procedure"})

	if params.DateFrom != "" && params.DateTo != "" {
		policyQuery = policyQuery.Where("documents.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := policyQuery.Count(&totalPolicies).Error
	if err != nil {
		return nil, err
	}

	otherDocQuery := repo.db.Model(&models.Document{}).
		Joins("JOIN document_types ON document_types.id = documents.document_type_id").
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id").
		Where("document_statuses.type = ?", "published").
		Where("document_types.type = ?", "others")

	if params.DateFrom != "" && params.DateTo != "" {
		otherDocQuery = otherDocQuery.Where("documents.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err = otherDocQuery.Count(&totalOtherDocuments).Error
	if err != nil {
		return nil, err
	}

	return &kpiDTOs.GRCDocumentSummaryDTO{
		TotalPolicies:       totalPolicies,
		TotalOtherDocuments: totalOtherDocuments,
	}, nil
}

func (repo *GRCDashboardRepository) GetGRCOrganizationSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCOrganizationSummaryDTO, error) {
	var totalControls int64
	var efficientControls int64

	controlQuery := repo.db.Model(&models.Control{})

	if params.DateFrom != "" && params.DateTo != "" {
		controlQuery = controlQuery.Where("controls.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	err := controlQuery.Count(&totalControls).Error
	if err != nil {
		return nil, err
	}

	efficientControlQuery := controlQuery.
		Joins("JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id").
		Where("control_efficiency_statuses.type = ?", "efficient")

	err = efficientControlQuery.Count(&efficientControls).Error
	if err != nil {
		return nil, err
	}

	// Calculate effectiveness rate
	effectivenessRate := 0.0
	if totalControls > 0 {
		effectivenessRate = math.Round(float64(efficientControls) / float64(totalControls) * 100.0)
	}

	return &kpiDTOs.GRCOrganizationSummaryDTO{
		ControlEffectivenessRate: effectivenessRate,
	}, nil
}

func (repo *GRCDashboardRepository) GetGRCComplianceOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceOverviewDTO, error) {
	summary := &kpiDTOs.GRCComplianceOverviewDTO{}

	complianceQuery := repo.db.Model(&models.ComplianceRegister{}).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id")

	if params.DateFrom != "" && params.DateTo != "" {
		complianceQuery = complianceQuery.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	var totalComplaintsRegistered int64
	err := complianceQuery.Count(&totalComplaintsRegistered).Error
	if err != nil {
		return nil, err
	}

	var compliant, notCompliant, partiallyCompliant, others int64
	err = complianceQuery.Select("COUNT(*)").
		Where("compliance_statuses.type = ?", "compliant").
		Count(&compliant).Error
	if err != nil {
		return nil, err
	}

	err = complianceQuery.Select("COUNT(*)").
		Where("compliance_statuses.type = ?", "not_compliant").
		Count(&notCompliant).Error
	if err != nil {
		return nil, err
	}

	err = complianceQuery.Select("COUNT(*)").
		Where("compliance_statuses.type = ?", "partially_compliant").
		Count(&partiallyCompliant).Error
	if err != nil {
		return nil, err
	}

	err = complianceQuery.Select("COUNT(*)").
		Where("compliance_statuses.type NOT IN (?)", []string{"compliant", "not_compliant", "partially_compliant"}).
		Count(&others).Error
	if err != nil {
		return nil, err
	}

	summary.TotalComplaintsRegistered = totalComplaintsRegistered
	summary.Compliant = compliant
	summary.NotCompliant = notCompliant
	summary.PartiallyCompliant = partiallyCompliant
	summary.Others = others

	return summary, nil
}

func (repo *GRCDashboardRepository) GetGRCRiskOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskOverviewDTO, error) {
	summary := &kpiDTOs.GRCRiskOverviewDTO{}

	riskQuery := repo.db.Model(&models.RiskRegister{}).
		Joins("JOIN risk_matrix_risk_ratings ON risk_register.current_risk_rating_id = risk_matrix_risk_ratings.id")

	if params.DateFrom != "" && params.DateTo != "" {
		riskQuery = riskQuery.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	var totalRisksRegistered int64
	err := riskQuery.Count(&totalRisksRegistered).Error
	if err != nil {
		return nil, err
	}

	var riskRatings []models.RiskMatrixRiskRating
	err = repo.db.Model(&models.RiskMatrixRiskRating{}).Where("status = ?", "active").Find(&riskRatings).Error
	if err != nil {
		return nil, err
	}

	riskCountsByRating := make(map[string]int64)

	for _, rating := range riskRatings {
		var count int64
		err := repo.db.Model(&models.RiskRegister{}).
			Where("current_risk_rating_id = ?", rating.ID).
			Count(&count).Error
		if err != nil {
			return nil, err
		}
		riskCountsByRating[rating.Type] = count
	}

	summary.TotalRisksRegistered = totalRisksRegistered
	summary.RiskCountsByRating = riskCountsByRating

	return summary, nil
}

func (repo *GRCDashboardRepository) GetTopComplianceRegisters(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopComplianceRegistersDTO, error) {
	var complianceRegisters []models.ComplianceRegister
	query := repo.db.Model(&models.ComplianceRegister{}).
		Order("total_score DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	query = query.Limit(5)

	err := query.Find(&complianceRegisters).Error
	if err != nil {
		return nil, err
	}

	var topRegisters []*kpiDTOs.GRCTopComplianceRegistersDTO
	for _, register := range complianceRegisters {
		complianceRegisterDTO := &kpiDTOs.GRCTopComplianceRegistersDTO{
			ID:         register.ID,
			TitleEn:    register.TitleEn,
			TitleAr:    register.TitleAr,
			TotalScore: register.TotalScore,
		}
		topRegisters = append(topRegisters, complianceRegisterDTO)
	}

	return topRegisters, nil
}

func (repo *GRCDashboardRepository) GetTopRiskRegisters(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopRiskRegistersDTO, error) {
	topRiskRegisters := []*kpiDTOs.GRCTopRiskRegistersDTO{}

	query := repo.db.Table("risk_register").
		Select(`
			risk_register.id, risk_register.title, 
			risk_matrix_risk_rating_values.value AS risk_rating_value, 
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en, 
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar, 
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code
		`).
		Joins("JOIN risk_matrix_risk_rating_values ON risk_register.current_risk_rating_value_id = risk_matrix_risk_rating_values.id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_rating_values.risk_matrix_risk_rating_id = risk_matrix_risk_ratings.id").
		Order("risk_matrix_risk_rating_values.value DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("risk_register.created_at BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	query = query.Limit(5)

	if err := query.Find(&topRiskRegisters).Error; err != nil {
		return nil, err
	}

	return topRiskRegisters, nil
}
