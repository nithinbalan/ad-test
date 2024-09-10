package models

import (
	"time"
)

type AuditTimeline struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditID             uint       `json:"audit_id" binding:"required" gorm:"not null"`
	AuditTimelineTypeID uint       `json:"audit_timeline_type_id" binding:"required" gorm:"not null"`
	TimelineStatusID    uint       `json:"timeline_status_id" binding:"required" gorm:"not null"`
	PerformedByID       *uint      `json:"performed_by_id" gorm:"default:null"`
	PerformedAt         *time.Time `json:"performed_at" gorm:"type:datetime;null"`

	// Associations
	Audit             Audit             `json:"audit_register" gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AuditTimelineType AuditTimelineType `json:"audit_timeline_type" gorm:"foreignKey:AuditTimelineTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TimelineStatus    TimelineStatus    `json:"timeline_status" gorm:"foreignKey:TimelineStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PerformedBy       User              `json:"performed_by" gorm:"foreignKey:PerformedByID"`
}

// TableName sets the insert table name for this struct type
func (AuditTimeline) TableName() string {
	return "audit_timelines"
}
