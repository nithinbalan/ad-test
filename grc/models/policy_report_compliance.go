package models

import "time"

type PolicyReportCompliance struct {
	ID                uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	PolicyReportID    uint      `json:"policy_report_id" gorm:"not null;index"`
	ComplianceTitleEn string    `json:"compliance_title_en" gorm:"type:nvarchar(255);not null"`
	ComplianceTitleAr string    `json:"compliance_title_ar" gorm:"type:nvarchar(255)"`
	Clauses           *string   `json:"clauses" gorm:"type:nvarchar(max)"`
	CreatedByID       uint      `json:"created_by_id" gorm:"not null;index"`
	CreatedAt         time.Time `json:"created_at" gorm:"not null"`

	PolicyReport  *PolicyReport `gorm:"foreignKey:PolicyReportID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (PolicyReportCompliance) TableName() string {
	return "policy_report_compliances"
}
