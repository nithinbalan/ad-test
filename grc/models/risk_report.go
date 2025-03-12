package models

import "time"

type RiskReport struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title       string     `json:"title" binding:"required,max=255"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     time.Time  `json:"end_date" binding:"required"`
	CreatedBYID uint       `json:"created_by_id" binding:"required"`
	CreatedAt   time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	CreatedByUser     *User                        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Sections          []RiskReportSection          `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Identifications   []RiskReportIdentification   `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ActionPlans       []RiskReportActionPlan       `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	KRIs              []RiskReportKRI              `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OrganizationUnits []RiskReportOrganizationUnit `gorm:"foreignKey:RiskReportID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName overrides the table name used by RiskReport to `risk_reports`
func (RiskReport) TableName() string {
	return "risk_reports"
}
