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
	complainceStatuses := []kpiDTOs.ComplianceStatusItemDTO{}
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

	if err := query.Scan(&complainceStatuses).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&models.ComplianceRegister{}).Count(&total).Error; err != nil {
		return nil, err
	}

	result := &kpiDTOs.ComplianceBreakdownDTO{
		CommplainceStatuses: complainceStatuses,
		Total:               total,
	}

	return result, nil
}
