package models

import "time"

type PolicyReportRisk struct {
	ID                    uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID        uint      `json:"policy_report_id" gorm:"not null;index"`
	RiskTitle             string    `json:"risk_title" gorm:"type:nvarchar(255);not null"`
	Description           *string   `json:"description" gorm:"type:nvarchar(max)"`
	RiskCategoryTitleEn   string    `json:"risk_category_title_en" gorm:"type:nvarchar(255)"`
	RiskCategoryTitleAr   string    `json:"risk_category_title_ar" gorm:"type:nvarchar(255)"`
	RiskStatusTitleEn     string    `json:"risk_status_title_en" gorm:"type:nvarchar(100)"`
	RiskStatusTitleAr     string    `json:"risk_status_title_ar" gorm:"type:nvarchar(100)"`
	RiskStatusColorCode   string    `json:"risk_status_color_code" gorm:"type:varchar(9)"`
	RiskLikelihoodTitleEn *string   `json:"risk_likelihood_title_en" gorm:"type:nvarchar(100)"`
	RiskLikelihoodTitleAr *string   `json:"risk_likelihood_title_ar" gorm:"type:nvarchar(100)"`
	RiskImpactTitleEn     *string   `json:"risk_impact_title_en" gorm:"type:nvarchar(100)"`
	RiskImpactTitleAr     *string   `json:"risk_impact_title_ar" gorm:"type:nvarchar(100)"`
	RiskRatingTitleEn     *string   `json:"risk_rating_title_en" gorm:"type:nvarchar(100)"`
	RiskRatingTitleAr     *string   `json:"risk_rating_title_ar" gorm:"type:nvarchar(100)"`
	RiskRatingColorCode   *string   `json:"risk_rating_color_code" gorm:"type:varchar(9)"`
	CreatedByID           uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt             time.Time `json:"created_at" gorm:"not null"`

	// Relationships
	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportRisk) TableName() string {
	return "policy_report_risks"
}
