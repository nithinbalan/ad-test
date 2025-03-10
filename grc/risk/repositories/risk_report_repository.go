package repositories

import (
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"time"

	"gorm.io/gorm"
)

type RiskReportRepositoryInterface interface {
	Create(riskReport *models.RiskReport, tx *gorm.DB) (uint, error)
	CreateRiskIdentification(identification *models.RiskReportIdentification, tx *gorm.DB) error
	GetRisksByDateRange(startDate, endDate time.Time, tx *gorm.DB) ([]dtos.RiskRegisterDTO, error)
	CreateSection(section *models.RiskReportSection, tx *gorm.DB) error
	BeginTransaction() *gorm.DB
}

type RiskReportRepository struct {
	db *gorm.DB
}

func NewRiskReportRepository(db *gorm.DB) RiskReportRepositoryInterface {
	return &RiskReportRepository{db: db}
}

// BeginTransaction starts a new transaction
func (r *RiskReportRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *RiskReportRepository) Create(riskReport *models.RiskReport, tx *gorm.DB) (uint, error) {
	if err := tx.Create(riskReport).Error; err != nil {
		return 0, err
	}
	return riskReport.ID, nil
}

func (repo *RiskReportRepository) CreateRiskIdentification(identification *models.RiskReportIdentification, tx *gorm.DB) error {
	return tx.Create(identification).Error
}

func (repo *RiskReportRepository) GetRisksByDateRange(startDate, endDate time.Time, tx *gorm.DB) ([]dtos.RiskRegisterDTO, error) {
	var risks []dtos.RiskRegisterDTO

	query := tx.Table("risk_register").
		Select(`
			risk_register.id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			risk_register.reported_at,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			risk_sources.title_en AS risk_source_title_en,
			risk_sources.title_ar AS risk_source_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
			inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
			inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
			inherent_likelihood.title_en AS inherent_likelihood_title_en,
			inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
			inherent_impact.title_en AS inherent_impact_title_en,
			inherent_impact.title_ar AS inherent_impact_title_ar,
			residual_risk_rating.title_en AS residual_risk_rating_title_en,
			residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
			residual_risk_rating.color_code AS residual_risk_rating_color_code,
			residual_likelihood.title_en AS residual_likelihood_title_en,
			residual_likelihood.title_ar AS residual_likelihood_title_ar,
			residual_impact.title_en AS residual_impact_title_en,
			residual_impact.title_ar AS residual_impact_title_ar,
			
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id").
		Where("risk_register.reported_at BETWEEN ? AND ?", startDate, endDate).
		Order("risk_register.reported_at DESC")

	if err := query.Find(&risks).Error; err != nil {
		return nil, err
	}

	return risks, nil
}

func (repo *RiskReportRepository) CreateSection(section *models.RiskReportSection, tx *gorm.DB) error {
	if tx != nil {
		return tx.Create(&section).Error
	}
	return repo.db.Create(&section).Error
}
