package models

import (
	"time"
)

type RiskTimeline struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID     uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	RiskTimelineTypeID uint       `json:"risk_timeline_type_id" binding:"required" gorm:"not null"`
	TimelineStatusID   uint       `json:"timeline_status_id" binding:"required" gorm:"not null"`
	PerformedByID      *uint      `json:"performed_by_id" gorm:"default:null"`
	PerformedAt        *time.Time `json:"performed_at" gorm:"type:datetime;null"`

	// Associations
	RiskRegister     RiskRegister     `json:"risk_register" gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskTimelineType RiskTimelineType `json:"risk_timeline_type" gorm:"foreignKey:RiskTimelineTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TimelineStatus   TimelineStatus   `json:"timeline_status" gorm:"foreignKey:TimelineStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PerformedBy      User             `json:"performed_by" gorm:"foreignKey:PerformedByID"`
}

// TableName sets the insert table name for this struct type
func (RiskTimeline) TableName() string {
	return "risk_timelines"
}
