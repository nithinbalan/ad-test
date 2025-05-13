package models

import "time"

type GRCReportRiskSummary struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	GRCReportID         uint       `json:"grc_report_id" gorm:"not null"`
	ReviewUserID        *uint      `json:"review_user_id"`
	ApproveUserID       *uint      `json:"approve_user_id"`
	RiskIdentifiedCount int        `json:"risk_identified_count" gorm:"not null"`
	CriticalRisks       *string    `json:"critical_risks" gorm:"type:nvarchar(max)"`
	TreatmentStrategies *string    `json:"treatment_strategies" gorm:"type:nvarchar(max)"`
	RiskScoreSummary    *string    `json:"risk_score_summary" gorm:"type:nvarchar(max)"`
	KRIs                *string    `json:"kris" gorm:"column:kris;type:nvarchar(max)"`
	RiskOverviewChart   *string    `json:"risk_overview_chart" gorm:"type:nvarchar(max)"`
	RiskHeatmapChart    *string    `json:"risk_heatmap_chart" gorm:"type:nvarchar(max)"`
	TopRisks            *string    `json:"top_risks" gorm:"type:nvarchar(max)"`
	CreatedBYID         uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID         *uint      `json:"updated_by_id"`
	UpdatedAt           *time.Time `json:"updated_at"`

	GRCReport     *GRCReport `gorm:"foreignKey:GRCReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User      `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User      `gorm:"foreignKey:UpdatedBYID"`
	ReviewUser    *User      `gorm:"foreignKey:ReviewUserID"`
	ApproveUser   *User      `gorm:"foreignKey:ApproveUserID"`
}

func (GRCReportRiskSummary) TableName() string {
	return "grc_report_risk_summary"
}
