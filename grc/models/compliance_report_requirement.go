package models

import (
	"time"
)

type ComplianceReportRequirement struct {
	ID                        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceReportID        uint      `json:"compliance_report_id" gorm:"not null" binding:"required"`
	ComplianceReportClauseID  uint      `json:"compliance_report_clause_id" gorm:"not null" binding:"required"`
	TitleEn                   string    `json:"title_en" gorm:"type:nvarchar(500);not null" binding:"required"`
	TitleAr                   string    `json:"title_ar" gorm:"type:nvarchar(500);not null" binding:"required"`
	Score                     *int      `json:"score" gorm:"default:0"`
	AverageScore              *float64  `json:"average_score" gorm:"default:0" comment:"Calculated average score across all assessments"`
	ComplianceStatusID        *uint     `json:"compliance_status_id" gorm:"default:null"`
	AverageComplianceStatusID *uint     `json:"average_compliance_status_id" gorm:"default:null" comment:"Compliance status based on average score"`
	CreatedAt                 time.Time `json:"created_at" gorm:"not null"`
	CreatedBYID               uint      `json:"created_by_id" gorm:"not null" binding:"required"`
	UpdatedAt                 time.Time `json:"updated_at" gorm:"default:null"`
	UpdatedBYID               uint      `json:"updated_by_id" gorm:"default:null"`

	// Associations
	ComplianceReport            ComplianceReport              `json:"compliance_report" gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportClause      ComplianceReportClause        `json:"compliance_report_clause" gorm:"foreignKey:ComplianceReportClauseID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceStatus            ComplianceStatus              `json:"compliance_status" gorm:"foreignKey:ComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AverageComplianceStatus     ComplianceStatus              `json:"average_compliance_status" gorm:"foreignKey:AverageComplianceStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportActionPlans []*ComplianceReportActionPlan `json:"compliance_report_action_plans" gorm:"foreignKey:ComplianceReportRequirementID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser               User                          `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser               User                          `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (ComplianceReportRequirement) TableName() string {
	return "compliance_report_requirements"
}
