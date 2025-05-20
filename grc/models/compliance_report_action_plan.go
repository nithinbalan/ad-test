package models

import (
	"time"
)

type ComplianceReportActionPlan struct {
	ID                            uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	ComplianceReportID            uint       `json:"compliance_report_id" gorm:"not null" binding:"required"`
	ComplianceReportRequirementID uint       `json:"compliance_report_requirement_id" gorm:"not null" binding:"required"`
	Title                         string     `json:"title" gorm:"type:nvarchar(MAX);not null" binding:"required"`
	Description                   *string    `json:"description" gorm:"type:nvarchar(MAX);default:null"`
	ActionPlanStatusID            uint       `json:"action_plan_status_id" gorm:"not null" binding:"required"`
	TargetDate                    *time.Time `json:"target_date" gorm:"default:null;type:date"`
	SeverityLevelID               *uint      `json:"severity_level_id" gorm:"default:null"`
	CreatedAt                     time.Time  `json:"created_at" gorm:"not null"`
	CreatedBYID                   uint       `json:"created_by_id" gorm:"not null" binding:"required"`
	UpdatedAt                     *time.Time `json:"updated_at" gorm:"default:null"`
	UpdatedBYID                   *uint      `json:"updated_by_id" gorm:"default:null"`

	// Associations
	ComplianceReport            *ComplianceReport            `json:"compliance_report" gorm:"foreignKey:ComplianceReportID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ComplianceReportRequirement *ComplianceReportRequirement `json:"compliance_report_requirement" gorm:"foreignKey:ComplianceReportRequirementID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ActionPlanStatus            *ActionPlanStatus            `json:"action_plan_status" gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	SeverityLevel               *SeverityLevel               `json:"severity_level" gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser               *User                        `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser               *User                        `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

func (ComplianceReportActionPlan) TableName() string {
	return "compliance_report_action_plans"
}
