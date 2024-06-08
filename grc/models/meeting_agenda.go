package models

import (
	"time"
)

type MeetingAgenda struct {
	ID        uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID uint       `json:"meeting_id" gorm:"null"`
	TimeFrom  *string    `json:"time_from" gorm:"type:varchar(45);null"`
	TimeTo    *string    `json:"time_to" gorm:"type:varchar(45);null"`
	CreatedBy uint       `json:"created_by" gorm:"not null"`
	UpdatedBy *uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Meeting       *Meeting `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User    `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser *User    `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by MeetingAgenda to `meeting_agendas`
func (MeetingAgenda) TableName() string {
	return "meeting_agendas"
}
