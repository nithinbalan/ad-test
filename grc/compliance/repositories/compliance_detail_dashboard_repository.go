package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"time"

	"gorm.io/gorm"
)

type ComplianceDetailDashboardRepositoryInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error)
	GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error)
	GetTop5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetBottom5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetComplianceDomainWiseOverview(params *complianceDTOs.ComplianceDashboardQueryParamsDTO, complianceRegisterID uint) (*complianceDTOs.ComplianceDomainWiseOverviewDTO, error)
	GetPolicyByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyByOrganizationUnitDTO, error)
	GetControlEffectivenessSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlEffectivenessDTO, error)
	GetControlImplementationSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlImplementationSummaryDTO, error)
	GetControlSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlSummaryDTO, error)
	GetControlEfficiencySummary(complianceRegisterID uint) ([]complianceDTOs.ComplianceSummaryItemDTO, error)
	GetControlCoverageByDomain(complianceRegisterID uint) ([]complianceDTOs.ComplianceControlCoverageByDomainDTO, error)
}

type ComplianceDetailDashboardRepository struct {
	db *gorm.DB
}

func NewComplianceDetailDashboardRepository(db *gorm.DB) *ComplianceDetailDashboardRepository {
	return &ComplianceDetailDashboardRepository{
		db: db,
	}
}

func (repo *ComplianceDetailDashboardRepository) GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error) {
	result := &complianceDTOs.ComplianceDetailSummaryDTO{}

	// Overall Compliance Score
	var complianceRegister models.ComplianceRegister
	err := repo.db.Select("score").
		First(&complianceRegister, complianceRegisterID).Error
	if err != nil {
		return nil, err
	}

	// Total Clauses
	err = repo.db.Model(&models.ComplianceClause{}).
		Where("compliance_register_id = ?", complianceRegisterID).
		Count(&result.TotalClauses).Error
	if err != nil {
		return nil, err
	}

	// Total Requirements
	err = repo.db.Model(&models.ComplianceRequirement{}).
		Where("compliance_register_id = ?", complianceRegisterID).
		Count(&result.TotalRequirements).Error
	if err != nil {
		return nil, err
	}

	// Total Controls
	var efficientControls int64
	err = repo.db.Model(&models.ComplianceControl{}).
		Joins("JOIN controls ON controls.id = compliance_controls.control_id").
		Joins("JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id").
		Where("compliance_controls.compliance_register_id = ?", complianceRegisterID).
		Where("control_implementation_statuses.type = ?", "implemented").
		Select(`
			COUNT(DISTINCT controls.id) AS total_controls_implemented,
			COUNT(DISTINCT CASE WHEN control_efficiency_statuses.type = 'efficient' THEN controls.id END) AS efficient_controls
		`).
		Row().
		Scan(&result.TotalControlsImplemented, &efficientControls)
	if err != nil {
		return nil, err
	}

	if result.TotalControlsImplemented > 0 {
		result.OverallControlEfficiency = (float64(efficientControls) / float64(result.TotalControlsImplemented)) * 100
	} else {
		result.OverallControlEfficiency = 0
	}

	if complianceRegister.Score != nil {
		result.OverallComplianceScore = *complianceRegister.Score
	} else {
		result.OverallComplianceScore = 0
	}

	return result, nil
}

func (repo *ComplianceDetailDashboardRepository) GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error) {
	res := &complianceDTOs.ComplianceClauseSummaryDTO{}

	// Total Clauses
	if err := repo.db.Model(&models.ComplianceClause{}).
		Where("compliance_register_id = ?", complianceRegisterID).
		Count(&res.TotalClauses).Error; err != nil {
		return nil, err
	}

	// Fully Compliant
	if err := repo.db.Model(&models.ComplianceClause{}).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_clauses.compliance_status_id").
		Where("compliance_clauses.compliance_register_id = ?", complianceRegisterID).
		Where("compliance_statuses.type = ?", "compliant").
		Count(&res.FullyCompliantClauses).Error; err != nil {
		return nil, err
	}

	// Partially Compliant
	if err := repo.db.Model(&models.ComplianceClause{}).
		Joins("JOIN compliance_statuses ON compliance_statuses.id = compliance_clauses.compliance_status_id").
		Where("compliance_clauses.compliance_register_id = ?", complianceRegisterID).
		Where("compliance_statuses.type = ?", "partially_compliant").
		Count(&res.PartiallyCompliantClauses).Error; err != nil {
		return nil, err
	}

	// Not Assessed (score IS NULL)
	if err := repo.db.Model(&models.ComplianceClause{}).
		Where("compliance_register_id = ?", complianceRegisterID).
		Where("score IS NULL").
		Count(&res.NotAssessedClauses).Error; err != nil {
		return nil, err
	}

	// Overdue Review clauses
	currentTime := time.Now()
	if err := repo.db.Model(&models.ComplianceClause{}).
		Joins("JOIN compliance_register ON compliance_register.id = compliance_clauses.compliance_register_id").
		Where("compliance_clauses.compliance_register_id = ?", complianceRegisterID).
		Where("compliance_register.next_review_date IS NOT NULL").
		Where("compliance_register.next_review_date < ?", currentTime).
		Count(&res.OverdueReviewClauses).Error; err != nil {
		return nil, err
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

	if err := repo.db.Model(&models.RiskComplianceRequirement{}).
		Joins("JOIN risk_register ON risk_register.id = risk_compliance_requirements.risk_register_id").
		Joins("JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = risk_register.current_risk_rating_id").
		Where("risk_compliance_requirements.compliance_register_id = ?", complianceRegisterID).
		Where("risk_compliance_requirements.compliance_clause_id IS NOT NULL").
		Where("risk_matrix_risk_ratings.type = ?", highestRating).
		Distinct("risk_compliance_requirements.compliance_clause_id").
		Count(&res.HighRiskClauses).Error; err != nil {
		return nil, err
	}

	// Percents
	denom := float64(res.TotalClauses)
	if res.TotalClauses > 0 {
		res.FullyCompliantClausesPercent = float64(res.FullyCompliantClauses) / denom * 100
		res.PartiallyCompliantClausesPercent = float64(res.PartiallyCompliantClauses) / denom * 100
		res.NotAssessedClausesPercent = float64(res.NotAssessedClauses) / denom * 100
		res.OverdueReviewClausesPercent = float64(res.OverdueReviewClauses) / denom * 100
		res.HighRiskClausesPercent = float64(res.HighRiskClauses) / denom * 100
	}

	return res, nil
}

func (repo *ComplianceDetailDashboardRepository) GetTop5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error) {
	results := []complianceDTOs.ComplianceClausePerformanceDTO{}

	// first level clauses
	err := repo.db.Table("compliance_clauses AS child").
		Select(`
            parent.title_en AS compliance_clause_title_en,
			parent.title_ar as compliance_clause_title_ar,
            AVG(CAST(child.score AS FLOAT)) AS compliance_clause_score
        `).
		Joins("JOIN compliance_clauses AS parent ON parent.id = child.compliance_clause_id").
		Where("child.compliance_register_id = ?", complianceRegisterID).
		Where("child.compliance_clause_id IS NOT NULL").
		Where("child.score IS NOT NULL").
		Group("parent.title_en, parent.title_ar").
		Order("compliance_clause_score DESC").
		Limit(5).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *ComplianceDetailDashboardRepository) GetBottom5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error) {
	results := []complianceDTOs.ComplianceClausePerformanceDTO{}

	err := repo.db.Table("compliance_clauses AS child").
		Select(`
            parent.title_en AS compliance_clause_title_en,
			parent.title_ar as compliance_clause_title_ar,
            AVG(CAST(child.score AS FLOAT)) AS compliance_clause_score
        `).
		Joins("JOIN compliance_clauses AS parent ON parent.id = child.compliance_clause_id").
		Where("child.compliance_register_id = ?", complianceRegisterID).
		Where("child.compliance_clause_id IS NOT NULL").
		Where("child.score IS NOT NULL").
		Group("parent.title_en, parent.title_ar").
		Order("compliance_clause_score ASC").
		Limit(5).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (repo *ComplianceDetailDashboardRepository) GetComplianceDomainWiseOverview(params *complianceDTOs.ComplianceDashboardQueryParamsDTO, complianceRegisterID uint) (*complianceDTOs.ComplianceDomainWiseOverviewDTO, error) {
	results := []complianceDTOs.ComplianceDomainWiseOverviewItemDTO{}
	var totalCount int64

	query := repo.db.Table("compliance_clauses AS parent").
		Select(`
            parent.title_en AS domain_title_en,
            parent.title_ar AS domain_title_ar,
            SUM(CASE WHEN cs.type = 'compliant' THEN 1 ELSE 0 END) AS compliant,
            SUM(CASE WHEN cs.type = 'partially_compliant' THEN 1 ELSE 0 END) AS partial,
            SUM(CASE WHEN cs.type = 'not-implemented' THEN 1 ELSE 0 END) AS non_compliant,
            SUM(CASE WHEN cs.type = 'in_progress' THEN 1 ELSE 0 END) AS pending
        `).
		Joins("LEFT JOIN compliance_statuses AS cs ON cs.id = parent.compliance_status_id").
		Where("parent.compliance_register_id = ?", complianceRegisterID).
		Where("parent.compliance_clause_id IS NULL").
		Group("parent.title_en, parent.title_ar")

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&results).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	respnose := &complianceDTOs.ComplianceDomainWiseOverviewDTO{
		Items: results,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return respnose, nil
}

func (repo *ComplianceDetailDashboardRepository) GetPolicyByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyByOrganizationUnitDTO, error) {
	results := []complianceDTOs.CompliancePolicyByOrganizationUnitDTO{}

	err := repo.db.Table("document_compliance_requirements AS dcr").
		Select(`
            parent_unit.title_en AS organization_unit_title_en,
            parent_unit.title_ar AS organization_unit_title_ar,
            COUNT(DISTINCT dcr.document_id) AS total_policies
        `).
		Joins("JOIN documents doc ON doc.id = dcr.document_id").
		Joins("LEFT JOIN organization_units child_unit ON child_unit.id = doc.responsible_unit_id").
		Joins("LEFT JOIN organization_units parent_unit ON COALESCE(child_unit.organization_unit_id, child_unit.id) = parent_unit.id").
		Joins("JOIN organization_levels ol ON ol.id = parent_unit.organization_level_id AND ol.is_primary = ?", true).
		Where("dcr.compliance_register_id = ?", complianceRegisterID).
		Where("parent_unit.organization_unit_id IS NULL").
		Group("parent_unit.title_en, parent_unit.title_ar").
		Order("total_policies DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *ComplianceDetailDashboardRepository) GetControlEffectivenessSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlEffectivenessDTO, error) {
	result := &complianceDTOs.ComplianceControlEffectivenessDTO{}

	var totalControls int64
	var efficientControls int64

	query := repo.db.
		Table("compliance_controls AS cc").
		Select(`
        COUNT(DISTINCT cc.control_id) AS total_controls,
        COUNT(DISTINCT CASE WHEN es.type = 'efficient' THEN cc.control_id END) AS efficient_controls
    `).
		Joins("LEFT JOIN controls c ON c.id = cc.control_id").
		Joins("LEFT JOIN control_efficiency_statuses es ON es.id = c.control_efficiency_status_id").
		Where("cc.compliance_register_id = ?", complianceRegisterID).
		Row()

	if err := query.Scan(&totalControls, &efficientControls); err != nil {
		return nil, err
	}

	if totalControls > 0 {
		result.OverallEfficiency = float64(efficientControls) / float64(totalControls) * 100
	} else {
		result.OverallEfficiency = 0
	}

	return result, nil
}

func (repo *ComplianceDetailDashboardRepository) GetControlImplementationSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlImplementationSummaryDTO, error) {
	result := &complianceDTOs.ComplianceControlImplementationSummaryDTO{}

	err := repo.db.
		Model(&models.ComplianceControl{}).
		Where("compliance_register_id = ?", complianceRegisterID).
		Distinct("control_id").
		Count(&result.TotalControls).Error
	if err != nil {
		return nil, err
	}

	implementationStatuses := []complianceDTOs.ComplianceSummaryItemDTO{}

	err = repo.db.
		Model(&models.ComplianceControl{}).
		Select(`
			cis.id,
			cis.title_en,
			cis.title_ar,
			cis.color_code,
			cis.type AS status_type,
			COUNT(DISTINCT compliance_controls.control_id) AS count
		`).
		Joins("JOIN controls ON controls.id = compliance_controls.control_id").
		Joins("JOIN control_implementation_statuses cis ON cis.id = controls.control_implementation_status_id").
		Where("compliance_controls.compliance_register_id = ?", complianceRegisterID).
		Group(`
			cis.id, 
			cis.title_en, 
			cis.title_ar, 
			cis.color_code, 
			cis.type
		`).
		Order("cis.id").
		Scan(&implementationStatuses).Error
	if err != nil {
		return nil, err
	}

	for i := range implementationStatuses {
		if result.TotalControls > 0 {
			implementationStatuses[i].Percent = float64(implementationStatuses[i].Count) / float64(result.TotalControls) * 100
		} else {
			implementationStatuses[i].Percent = 0
		}
	}

	result.ImplementationStatuses = implementationStatuses
	return result, nil
}

func (repo *ComplianceDetailDashboardRepository) GetControlSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlSummaryDTO, error) {
	result := &complianceDTOs.ComplianceControlSummaryDTO{}

	err := repo.db.
		Table("compliance_controls AS cc").
		Select(`
			COUNT(DISTINCT cc.control_id) AS total_controls,
			COUNT(DISTINCT CASE 
				WHEN cc.completed_at IS NOT NULL 
				AND cc.completed_at < CAST(GETDATE() AS date) 
				THEN cc.control_id 
			END) AS overdue_implemented_controls,
			COUNT(DISTINCT CASE 
				WHEN cis.type = 'partially-implemented' 
				THEN cc.control_id 
			END) AS partially_implemented_controls,
			COUNT(DISTINCT CASE 
				WHEN cis.type = 'not-implemented' 
				THEN cc.control_id 
			END) AS not_implemented_controls
		`).
		Joins("LEFT JOIN controls c ON c.id = cc.control_id").
		Joins("LEFT JOIN control_implementation_statuses cis ON cis.id = c.control_implementation_status_id").
		Where("cc.compliance_register_id = ?", complianceRegisterID).
		Scan(result).Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repo *ComplianceDetailDashboardRepository) GetControlEfficiencySummary(complianceRegisterID uint) ([]complianceDTOs.ComplianceSummaryItemDTO, error) {
	summary := []complianceDTOs.ComplianceSummaryItemDTO{}

	var totalControls int64
	err := repo.db.
		Model(&models.ComplianceControl{}).
		Where("compliance_register_id = ?", complianceRegisterID).
		Distinct("control_id").
		Count(&totalControls).Error
	if err != nil {
		return nil, err
	}

	err = repo.db.Model(&models.ComplianceControl{}).
		Select(`
			ces.id,
			ces.title_en,
			ces.title_ar,
			ces.color_code,
			ces.type AS status_type,
			COUNT(DISTINCT compliance_controls.control_id) AS count
		`).
		Joins("JOIN controls ON controls.id = compliance_controls.control_id").
		Joins("JOIN control_efficiency_statuses ces ON ces.id = controls.control_efficiency_status_id").
		Where("compliance_controls.compliance_register_id = ?", complianceRegisterID).
		Group(`
			ces.id,
			ces.title_en,
			ces.title_ar,
			ces.color_code,
			ces.type
		`).
		Order("ces.id").
		Scan(&summary).Error
	if err != nil {
		return nil, err
	}

	for i := range summary {
		if totalControls > 0 {
			summary[i].Percent = float64(summary[i].Count) / float64(totalControls) * 100
		} else {
			summary[i].Percent = 0
		}
	}

	return summary, nil
}

func (repo *ComplianceDetailDashboardRepository) GetControlCoverageByDomain(complianceRegisterID uint) ([]complianceDTOs.ComplianceControlCoverageByDomainDTO, error) {
	results := []complianceDTOs.ComplianceControlCoverageByDomainDTO{}

	err := repo.db.Table("compliance_clauses AS parent").
		Select(`
            parent.title_en AS domain_title_en,
            parent.title_ar AS domain_title_ar,
            COUNT(DISTINCT cc.control_id) AS total_controls
        `).
		Joins(`
            LEFT JOIN compliance_controls AS cc
               ON cc.compliance_register_id = parent.compliance_register_id
              AND cc.compliance_clause_id = parent.id
        `).
		Where("parent.compliance_register_id = ?", complianceRegisterID).
		Where("parent.compliance_clause_id IS NULL").
		Group("parent.title_en, parent.title_ar").
		Order("total_controls DESC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}
