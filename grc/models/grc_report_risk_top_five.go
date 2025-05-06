package models

import "time"

type GRCReportRiskTopFive struct {
	ID                          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID                 uint      `json:"grc_report_id" gorm:"not null"`
	RiskReferenceNumber         string    `json:"risk_reference_number" gorm:"type:nvarchar(50);not null"`
	RiskCategoryTitleEn         string    `json:"risk_category_title_en" gorm:"type:nvarchar(255);not null"`
	RiskCategoryTitleAr         string    `json:"risk_category_title_ar" gorm:"type:nvarchar(255);not null"`
	Title                       string    `json:"title" gorm:"type:nvarchar(255);not null"`
	Description                 *string   `json:"description" gorm:"type:nvarchar(max)"`
	ResponsibleUserNameEn       string    `json:"responsible_user_name_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUserNameAr       string    `json:"responsible_user_name_ar" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleEn      string    `json:"responsible_unit_title_en" gorm:"type:nvarchar(255);not null"`
	ResponsibleUnitTitleAr      string    `json:"responsible_unit_title_ar" gorm:"type:nvarchar(255);not null"`
	IsOpen                      bool      `json:"is_open" gorm:"not null;default:0"`
	InherentRiskRatingTitleEn   *string   `json:"inherent_risk_rating_title_en" gorm:"type:nvarchar(255)"`
	InherentRiskRatingTitleAr   *string   `json:"inherent_risk_rating_title_ar" gorm:"type:nvarchar(255)"`
	InherentRiskRatingColorCode *string   `json:"inherent_risk_rating_color_code" gorm:"type:nvarchar(10)"`
	InherentRiskRatingScore     *uint     `json:"inherent_risk_rating_score" gorm:"type:nvarchar(50)"`
	InherentLikelihoodTitleEn   *string   `json:"inherent_likelihood_title_en" gorm:"type:nvarchar(255)"`
	InherentLikelihoodTitleAr   *string   `json:"inherent_likelihood_title_ar" gorm:"type:nvarchar(255)"`
	InherentImpactTitleEn       *string   `json:"inherent_impact_title_en" gorm:"type:nvarchar(255)"`
	InherentImpactTitleAr       *string   `json:"inherent_impact_title_ar" gorm:"type:nvarchar(255)"`
	ResidualRiskRatingTitleEn   *string   `json:"residual_risk_rating_title_en" gorm:"type:nvarchar(255)"`
	ResidualRiskRatingTitleAr   *string   `json:"residual_risk_rating_title_ar" gorm:"type:nvarchar(255)"`
	ResidualRiskRatingColorCode *string   `json:"residual_risk_rating_color_code" gorm:"type:nvarchar(10)"`
	ResidualRiskRatingScore     *uint     `json:"residual_risk_rating_score" gorm:"type:nvarchar(50)"`
	ResidualLikelihoodTitleEn   *string   `json:"residual_likelihood_title_en" gorm:"type:nvarchar(255)"`
	ResidualLikelihoodTitleAr   *string   `json:"residual_likelihood_title_ar" gorm:"type:nvarchar(255)"`
	ResidualImpactTitleEn       *string   `json:"residual_impact_title_en" gorm:"type:nvarchar(255)"`
	ResidualImpactTitleAr       *string   `json:"residual_impact_title_ar" gorm:"type:nvarchar(255)"`
	ExistingControls            *string   `json:"existing_controls" gorm:"type:nvarchar(max)"`
	RiskResponseStrategyTitleEn *string   `json:"risk_response_strategy_title_en" gorm:"type:nvarchar(255)"`
	RiskResponseStrategyTitleAr *string   `json:"risk_response_strategy_title_ar" gorm:"type:nvarchar(255)"`
	ApprovalNote                *string   `json:"approval_note" gorm:"type:nvarchar(max)"`
	CreatedBYID                 uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	GRCReport     *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
}

func (GRCReportRiskTopFive) TableName() string {
	return "grc_report_risk_top_five"
}
