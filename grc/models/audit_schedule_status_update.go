package models

import (
	"time"
)

type AuditScheduleStatusUpdate struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditScheduleID       uint       `json:"audit_schedule_id" gorm:"not null"`
	AuditScheduleStatusID uint       `json:"audit_schedule_status_id" gorm:"not null"`
	Comment               string     `json:"comment" gorm:"type:nvarchar(max);not null"`
	Percentage            *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	AuditSchedule       *AuditSchedule       `gorm:"foreignKey:AuditScheduleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	AuditScheduleStatus *AuditScheduleStatus `gorm:"foreignKey:AuditScheduleStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User                `gorm:"foreignKey:UpdatedBYID"`
}

func (AuditScheduleStatusUpdate) TableName() string {
	return "audit_schedule_status_updates"
}
