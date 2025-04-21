package models

import (
	"time"
)

type ComplianceReportClause struct {
	ID                           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceReportID           uint       `json:"compliance_report_id" gorm:"not null" binding:"required"`
	ComplianceReportDepartmentID uint       `json:"compliance_report_department_id" gorm:"not null" binding:"required"`
	ClauseNo                     string     `json:"clause_no" gorm:"type:nvarchar(45);not null" binding:"required"`
	ClauseTitleEn                string     `json:"clause_title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	ClauseTitleAr                *string    `json:"clause_title_ar" gorm:"type:nvarchar(500);default:null"`
	DescriptionEn                *string    `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr                *string    `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	CreatedAt                    time.Time  `json:"created_at" gorm:"not null"`
	CreatedBYID                  uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	UpdatedAt                    *time.Time `json:"updated_at" gorm:"default:null"`
	UpdatedBYID                  *uint      `json:"updated_by_id" gorm:"default:null"`

	// Associations
	ComplianceReport                         ComplianceReport               `json:"compliance_report" gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportDepartment               ComplianceReportDepartment     `json:"compliance_report_department" gorm:"foreignKey:ComplianceReportDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportRequirements             []*ComplianceReportRequirement `json:"compliance_report_requirements" gorm:"foreignKey:ComplianceReportClauseID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportNonCompliantRequirements []*ComplianceReportRequirement `json:"compliance_report_non_compliant_requirements" gorm:"foreignKey:ComplianceReportClauseID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser                            User                           `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser                            User                           `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (ComplianceReportClause) TableName() string {
	return "compliance_report_clauses"
}
