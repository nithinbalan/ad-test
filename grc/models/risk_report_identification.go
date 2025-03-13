package models

import "time"

type RiskReportIdentification struct {
	ID                     uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID           uint      `json:"risk_report_id" gorm:"not null"`
	ReferenceNumber        string    `json:"reference_number" gorm:"not null"`
	RiskCategoryTitleEn    string    `json:"risk_category_title_en" gorm:"not null"`
	RiskCategoryTitleAr    string    `json:"risk_category_title_ar" gorm:"not null"`
	Title                  string    `json:"title" gorm:"not null"`
	Description            *string   `json:"description"`
	RiskSourceTitleEn      *string   `json:"risk_source_title_en"`
	RiskSourceTitleAr      *string   `json:"risk_source_title_ar"`
	RiskReportedAt         time.Time `json:"risk_reported_at" gorm:"not null"`
	ResponsibleUnitTitleEn string    `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr string    `json:"responsible_unit_title_ar"`

	InherentRiskRatingTitleEn   *string `json:"inherent_risk_rating_title_en"`
	InherentRiskRatingTitleAr   *string `json:"inherent_risk_rating_title_ar"`
	InherentRiskRatingColorCode *string `json:"inherent_risk_rating_color_code"`

	InherentLikelihoodTitleEn *string `json:"inherent_likelihood_title_en"`
	InherentLikelihoodTitleAr *string `json:"inherent_likelihood_title_ar"`

	InherentImpactTitleEn *string `json:"inherent_impact_title_en"`
	InherentImpactTitleAr *string `json:"inherent_impact_title_ar"`

	ResidualRiskRatingTitleEn   *string `json:"residual_risk_rating_title_en"`
	ResidualRiskRatingTitleAr   *string `json:"residual_risk_rating_title_ar"`
	ResidualRiskRatingColorCode *string `json:"residual_risk_rating_color_code"`

	ResidualLikelihoodTitleEn *string `json:"residual_likelihood_title_en"`
	ResidualLikelihoodTitleAr *string `json:"residual_likelihood_title_ar"`

	ResidualImpactTitleEn *string `json:"residual_impact_title_en"`
	ResidualImpactTitleAr *string `json:"residual_impact_title_ar"`

	ExistingControls *string `json:"existing_controls"`
	ProposedControls *string `json:"proposed_controls"`

	RiskStatusTitleEn   string `json:"risk_status_title_en"`
	RiskStatusTitleAr   string `json:"risk_status_title_ar"`
	RiskStatusColorCode string `json:"risk_status_color_code"`

	RiskReport *RiskReport `gorm:"foreignKey:RiskReportID"`
}

func (RiskReportIdentification) TableName() string {
	return "risk_report_identifications"
}
