package repositories

import (
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type ComplianceReportRepositoryInterface interface {
	GetComplianceByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error)
}

type ComplianceReportRepository struct {
	db *gorm.DB
}

func NewComplianceReportRepository(db *gorm.DB) *ComplianceReportRepository {
	return &ComplianceReportRepository{
		db: db,
	}
}

func (repo *ComplianceReportRepository) GetComplianceByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ComplianceByCategoryDTO, error) {
	results := []*reportDTOs.ComplianceByCategoryDTO{}
	query := repo.db.Table("compliance_register").
		Select("compliance_categories.id AS compliance_category_id, compliance_categories.title_en AS compliance_category_title_en, compliance_categories.title_ar AS compliance_category_title_ar, COUNT(compliance_register.id) AS count").
		Joins("JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id").
		Group("compliance_categories.id, compliance_categories.title_en, compliance_categories.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
