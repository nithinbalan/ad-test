package models

import (
	"time"
)

type AuditPlanStakeholderEngagement struct {
	ID                    uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ReferenceNumber       string     `json:"reference_number" gorm:"type:nvarchar(255);not null"`
	AuditPlanID           uint       `json:"audit_plan_id" gorm:"not null"`
	AuditEngagementTypeID uint       `json:"engagement_type_id" gorm:"not null"`
	Title                 string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Objectives            *string    `json:"objectives" gorm:"type:nvarchar(MAX);default:null"`
	ResponsibleUserID     *uint      `json:"responsible_user_id" gorm:"default:null"`
	AuditScheduleStatusID uint       `json:"audit_schedule_status_id" gorm:"not null"`
	CreatedByID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlan           *AuditPlan           `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditEngagementType *AuditEngagementType `gorm:"foreignKey:AuditEngagementTypeID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ResponsibleUser     *User                `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditScheduleStatus *AuditScheduleStatus `gorm:"foreignKey:AuditScheduleStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedByID;"`
	UpdatedByUser       *User                `gorm:"foreignKey:UpdatedByID;"`
}

// TableName overrides the table name used by AuditPlanStakeholderEngagement to `audit_plan_stakeholder_engagements`
func (AuditPlanStakeholderEngagement) TableName() string {
	return "audit_plan_stakeholder_engagements"
}
