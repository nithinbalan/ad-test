package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceDetailDashboardRepositoryInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDashboardDetailSummaryDTO, error)
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
