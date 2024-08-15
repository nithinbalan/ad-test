package models

import (
	"time"
)

type IncidentTimeline struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID     uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	IncidentTimelineTypeID uint       `json:"incident_timeline_type_id" binding:"required" gorm:"not null"`
	TimelineStatusID       uint       `json:"timeline_status_id" binding:"required" gorm:"not null"`
	PerformedByID          *uint      `json:"performed_by_id" gorm:"default:null"`
	PerformedAt            *time.Time `json:"performed_at" gorm:"type:datetime;null"`

	// Associations
	IncidentRegister     IncidentRegister     `json:"incident_register" gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentTimelineType IncidentTimelineType `json:"incident_timeline_type" gorm:"foreignKey:IncidentTimelineTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	TimelineStatus       TimelineStatus       `json:"timeline_status" gorm:"foreignKey:TimelineStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	PerformedBy          User                 `json:"performed_by" gorm:"foreignKey:PerformedByID"`
}

// TableName sets the insert table name for this struct type
func (IncidentTimeline) TableName() string {
	return "incident_timelines"
}
