package models

import (
	"time"
)

type ComplianceReportDepartment struct {
	ID                   uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceReportID   uint      `json:"compliance_report_id" gorm:"not null" binding:"required"`
	PrimaryUnitID        uint      `json:"primary_unit_id" gorm:"not null;column:organization_unit_id" binding:"required"`
	CompliancePercentage float64   `json:"compliance_percentage" gorm:"type:decimal(10,2);default:0.00"`
	CreatedAt            time.Time `json:"created_at" gorm:"not null"`
	CreatedBYID          uint      `json:"created_by_id" gorm:"not null" binding:"required"`

	// Associations
	ComplianceReport        *ComplianceReport         `json:"compliance_report" gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	OrganizationUnit        *OrganizationUnit         `json:"organization_unit" gorm:"foreignKey:PrimaryUnitID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportClauses []*ComplianceReportClause `json:"compliance_report_clauses" gorm:"foreignKey:ComplianceReportDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser           *User                     `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (ComplianceReportDepartment) TableName() string {
	return "compliance_report_departments"
}
