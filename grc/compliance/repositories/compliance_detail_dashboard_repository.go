package repositories

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"time"

	"gorm.io/gorm"
)

type ComplianceDetailDashboardRepositoryInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardDetailSummaryDTO, error)
	GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClauseSummaryDTO, error)
	GetTop5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error)
	GetBottom5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error)
}

type ComplianceDetailDashboardRepository struct {
	db *gorm.DB
}

func NewComplianceDetailDashboardRepository(db *gorm.DB) *ComplianceDetailDashboardRepository {
	return &ComplianceDetailDashboardRepository{
		db: db,
	}
}

func (repo *ComplianceDetailDashboardRepository) GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardDetailSummaryDTO, error) {
	result := &complianceDTOs.ComplianceDashboardDetailSummaryDTO{}

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

func (repo *ComplianceDetailDashboardRepository) GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClauseSummaryDTO, error) {
	res := &complianceDTOs.ComplianceDashboardClauseSummaryDTO{}

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

func (repo *ComplianceDetailDashboardRepository) GetTop5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error) {
	results := []complianceDTOs.ComplianceDashboardClausePerformanceItemDTO{}

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
	return &complianceDTOs.ComplianceDashboardClausePerformanceDTO{
		Items: results,
	}, nil
}

func (repo *ComplianceDetailDashboardRepository) GetBottom5ClausePerformance(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardClausePerformanceDTO, error) {
	results := []complianceDTOs.ComplianceDashboardClausePerformanceItemDTO{}

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
	return &complianceDTOs.ComplianceDashboardClausePerformanceDTO{
		Items: results,
	}, nil
}
