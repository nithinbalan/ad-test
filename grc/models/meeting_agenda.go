package models

import (
	"time"
)

type MeetingAgenda struct {
	ID          uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID   uint       `json:"meeting_id" gorm:"null"`
	Title       string     `json:"title" gorm:"type:nvarchar(500);not null" binding:"required"`
	CreatedBYID uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Meeting        *Meeting         `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingMinutes []*MeetingMinute `gorm:"foreignKey:MeetingAgendaID"`
	CreatedByUser  *User            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser  *User            `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by MeetingAgenda to `meeting_agendas`
func (MeetingAgenda) TableName() string {
	return "meeting_agendas"
}
