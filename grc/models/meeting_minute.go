package models

import (
	"time"
)

type MeetingMinute struct {
	ID              uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingAgendaID *uint      `json:"meeting_agenda_id" gorm:"null"`
	Minutes         *string    `json:"minutes" gorm:"type:nvarchar(max);null"`
	CreatedBy       uint       `json:"created_by" gorm:"not null"`
	UpdatedBy       *uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt       *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	MeetingAgenda *MeetingAgenda `gorm:"foreignKey:MeetingAgendaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User          `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User          `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by MeetingMinute to `meeting_minutes`
func (MeetingMinute) TableName() string {
	return "meeting_minutes"
}
