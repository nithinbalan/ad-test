package repositories

import (
	"fmt"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceDashboardRepositoryInterface interface {
	GetComplianceOverview(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceOverviewDTO, error)
	GetComplianceBreakdown(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceBreakdownDTO, error)
	GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error)
	GetComplianceControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error)
	GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceEfficiencySummaryDTO, error)
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

func (repo *ComplianceDashboardRepository) GetHighPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error) {
	results := []kpiDTOs.CompliancePerformanceItemDTO{}

	query := repo.db.Model(&models.ComplianceRegister{}).
		Select(`
            compliance_register.standard_code as framework_name,
            AVG(COALESCE(compliance_register.score, 0)) as score
        `).
		Where("compliance_register.score IS NOT NULL").
		Group("compliance_register.standard_code").
		Having("AVG(compliance_register.score) > ?", 75).
		Order("score DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}
	summary := &kpiDTOs.CompliancePerformanceRadarDTO{
		Items: results,
	}

	return summary, nil
}

func (repo *ComplianceDashboardRepository) GetLeastPerformanceCompliances(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.CompliancePerformanceRadarDTO, error) {
	results := []kpiDTOs.CompliancePerformanceItemDTO{}

	query := repo.db.Model(&models.ComplianceRegister{}).
		Select(`
            compliance_register.standard_code as compliance_standard_code,
            AVG(COALESCE(compliance_register.score, 0)) as score
        `).
		Where("compliance_register.score IS NOT NULL").
		Group("compliance_register.standard_code").
		Having("AVG(compliance_register.score) <= ?", 75).
		Order("score ASC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Scan(&results).Error; err != nil {
		return nil, err
	}

	summary := &kpiDTOs.CompliancePerformanceRadarDTO{
		Items: results,
	}

	return summary, nil
}

func (repo *ComplianceDashboardRepository) GetComplianceControlSummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceControlSummaryDTO, error) {
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

func (repo *ComplianceDashboardRepository) GetComplianceEfficiencySummary(params *kpiDTOs.ComplianceDashboardQueryParamsDTO) (*kpiDTOs.ComplianceEfficiencySummaryDTO, error) {

	results := []kpiDTOs.ComplianceEfficiencySummaryQueryResultDTO{}

	query := repo.db.Table("compliance_controls").
		Select(`
            compliance_register.standard_code as compliance_standard_code,
            control_efficiency_statuses.id as status_id,
            control_efficiency_statuses.title_en,
            control_efficiency_statuses.title_ar,
            control_efficiency_statuses.color_code,
            control_efficiency_statuses.type,
            COUNT(*) as count
        `).
		Joins("JOIN controls ON compliance_controls.control_id = controls.id").
		Joins("JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id").
		Joins("JOIN compliance_register ON compliance_controls.compliance_register_id = compliance_register.id")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_controls.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	query = query.Group(`
		compliance_register.standard_code, 
		control_efficiency_statuses.id, 
		control_efficiency_statuses.title_en, 
		control_efficiency_statuses.title_ar, 
		control_efficiency_statuses.color_code, 
		control_efficiency_statuses.type
	`)

	if err := query.Order("compliance_register.standard_code").Scan(&results).Error; err != nil {
		return nil, err
	}

	// Transform flat results into nested JSON structure
	efficiencyMap := make(map[string][]kpiDTOs.ComplianceSummaryItemDTO)
	for _, item := range results {
		efficiencyMap[item.ComplianceStandardCode] = append(efficiencyMap[item.ComplianceStandardCode], kpiDTOs.ComplianceSummaryItemDTO{
			ID:        item.StatusID,
			TitleEn:   item.TitleEn,
			TitleAr:   item.TitleAr,
			ColorCode: item.ColorCode,
			Type:      item.Type,
			Count:     item.Count,
		})
	}

	summary := &kpiDTOs.ComplianceEfficiencySummaryDTO{
		Items: make([]kpiDTOs.ComplianceEfficiencySummaryItemDTO, 0),
	}

	for compliance, statuses := range efficiencyMap {
		var total int64
		for _, s := range statuses {
			total += s.Count
		}
		for i := range statuses {
			if total > 0 {
				statuses[i].Percent = float64(statuses[i].Count) / float64(total) * 100
			} else {
				statuses[i].Percent = 0
			}
		}
		summary.Items = append(summary.Items, kpiDTOs.ComplianceEfficiencySummaryItemDTO{
			ComplianceStandardCode: compliance,
			EfficiencyStatuses:     statuses,
		})
	}

	return summary, nil
}
