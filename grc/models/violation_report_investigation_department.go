package models

import (
	"time"
)

type ViolationReportInvestigationDepartment struct {
	ID                             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ViolationReportInvestigationID uint      `json:"violation_report_investigation_id" binding:"required"`
	DepartmentID                   uint      `json:"department_id" binding:"required"`
	CreatedBYID                    uint      `json:"created_by_id" binding:"required"`
	CreatedAt                      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	ViolationReportInvestigation *ViolationReportInvestigation `gorm:"foreignKey:ViolationReportInvestigationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Department                   *OrganizationUnit             `gorm:"foreignKey:DepartmentID"`
	CreatedByUser                *User                         `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (ViolationReportInvestigationDepartment) TableName() string {
	return "violation_report_investigation_departments"
}
