package models

import (
	"time"
)

type Meeting struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID          *uint      `json:"meeting_id" gorm:"null"`
	MeetingTypeID      uint       `json:"meeting_type_id" binding:"required" gorm:"not null"`
	MeetingFrequencyID *uint      `json:"meeting_frequency_id" gorm:"null"`
	Title              string     `json:"title" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	Location           *string    `json:"location" gorm:"type:varchar(100);null"`
	Date               time.Time  `json:"date" binding:"required" gorm:"type:date;not null"`
	StartTime          time.Time  `json:"start_time" binding:"required" gorm:"type:time;not null"`
	EndTime            time.Time  `json:"end_time" binding:"required" gorm:"type:time;not null"`
	OrganizerID        uint       `json:"organizer_id" binding:"required" gorm:"not null"`
	CreatedBy          uint       `json:"created_by" gorm:"not null"`
	UpdatedBy          *uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	MeetingStatusID    uint       `json:"meeting_status_id" binding:"required" gorm:"not null"`

	// Relationships
	MeetingType      *MeetingType          `gorm:"foreignKey:MeetingTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Meeting          *Meeting              `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingFrequency *MeetingFrequency     `gorm:"foreignKey:MeetingFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Organizer        *User                 `gorm:"foreignKey:OrganizerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingStatus    *MeetingStatus        `gorm:"foreignKey:MeetingStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User                 `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser    *User                 `gorm:"foreignKey:UpdatedBy"`
	Participants     []*MeetingParticipant `gorm:"foreignkey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by Meeting to `meetings`
func (Meeting) TableName() string {
	return "meetings"
}
