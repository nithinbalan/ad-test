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
            compliance_register.standard_code as framework_name,
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
