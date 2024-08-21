package models

import (
	"time"
)

type MeetingMinute struct {
	ID                    uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingAgendaID       *uint      `json:"meeting_agenda_id" gorm:"null"`
	MeetingMinuteStatusID uint       `json:"meeting_minute_status_id" gorm:"not null"`
	Minutes               *string    `json:"minutes" gorm:"type:nvarchar(max);null"`
	CreatedBYID           uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID           *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt             time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt             *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	MeetingAgenda       *MeetingAgenda       `gorm:"foreignKey:MeetingAgendaID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingMinuteStatus *MeetingMinuteStatus `gorm:"foreignKey:MeetingMinuteStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser       *User                `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser       *User                `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by MeetingMinute to `meeting_minutes`
func (MeetingMinute) TableName() string {
	return "meeting_minutes"
}
