package models

import (
	"time"
)

type AuditPlanSchedule struct {
	ID                    uint       `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	AuditPlanID           uint       `json:"audit_plan_id" gorm:"not null"`
	AuditActivityTypeID   uint       `json:"audit_activity_type_id" gorm:"not null"`
	ResponsibleUserID     *uint      `json:"responsible_user_id" gorm:"default:null"`
	Title                 string     `json:"title" gorm:"type:nvarchar(255);not null"`
	Description           *string    `json:"description" gorm:"type:nvarchar(MAX);default:null"`
	StartDate             *time.Time `json:"start_date" gorm:"type:date;default:null"`
	TargetDate            *time.Time `json:"target_date" gorm:"type:date;default:null"`
	AuditScheduleStatusID uint       `json:"audit_schedule_status_id" gorm:"not null"`
	CreatedByID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	AuditPlan           *AuditPlan                        `gorm:"foreignKey:AuditPlanID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditActivityType   *AuditActivityType                `gorm:"foreignKey:AuditActivityTypeID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	ResponsibleUser     *User                             `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditScheduleStatus *AuditScheduleStatus              `gorm:"foreignKey:AuditScheduleStatusID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	AuditableItems      []*AuditPlanScheduleAuditableItem `gorm:"foreignKey:AuditPlanScheduleID"`
	CreatedByUser       *User                             `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	UpdatedByUser       *User                             `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by AuditSchedule to `audit_plan_schedules`
func (AuditPlanSchedule) TableName() string {
	return "audit_plan_schedules"
}
