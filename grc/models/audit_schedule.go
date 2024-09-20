package models

import (
	"time"
)

type AuditSchedule struct {
	ID                    uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID           uint       `json:"audit_plan_id" gorm:"not null"`
	AuditPlanScheduleID   *uint      `json:"audit_plan_schedule_id" gorm:"default:null"`
	AuditActivityTypeID   uint       `json:"audit_activity_type_id" gorm:"not null"`
	AuditID               uint       `json:"audit_id" gorm:"not null"`
	Title                 string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description           *string    `json:"description,omitempty" gorm:"type:nvarchar(max);default:null"`
	StartDate             *time.Time `json:"start_date,omitempty" gorm:"type:date;default:null"`
	TargetDate            *time.Time `json:"target_date,omitempty" gorm:"type:date;default:null"`
	ResponsibleUserID     *uint      `json:"responsible_user_id,omitempty" gorm:"default:null"`
	AuditScheduleStatusID uint       `json:"audit_schedule_status_id" gorm:"not null"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlan           *AuditPlan           `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditPlanSchedule   *AuditPlanSchedule   `gorm:"foreignKey:AuditPlanScheduleID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditActivityType   *AuditActivityType   `gorm:"foreignKey:AuditActivityTypeID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	Audit               *Audit               `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ResponsibleUser     *User                `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditScheduleStatus *AuditScheduleStatus `gorm:"foreignKey:AuditScheduleStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser       *User                `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`

	AuditableItems []*AuditScheduleAuditableItem `gorm:"foreignKey:AuditScheduleID"`
}

// TableName overrides the table name used by GORM to `audit_schedules`
func (AuditSchedule) TableName() string {
	return "audit_schedules"
}
