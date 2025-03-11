package models

import "time"

type RiskReportOrganizationUnit struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskReportID       uint      `json:"risk_report_id" gorm:"not null"`
	OrganizationUnitID uint      `json:"organization_unit_id" gorm:"not null"`
	CreatedBYID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskReport       *RiskReport       `gorm:"foreignKey:RiskReportID"`
	OrganizationUnit *OrganizationUnit `gorm:"foreignKey:OrganizationUnitID"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskReportOrganizationUnit) TableName() string {
	return "risk_report_organization_units"
}
