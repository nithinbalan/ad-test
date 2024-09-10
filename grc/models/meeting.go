package models

import (
	"time"
)

type Meeting struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber    string     `json:"reference_number" binding:"omitempty,max=45" gorm:"size:45"`
	MeetingID          *uint      `json:"meeting_id" gorm:"null"`
	MeetingTypeID      uint       `json:"meeting_type_id" binding:"required" gorm:"not null"`
	MeetingFrequencyID *uint      `json:"meeting_frequency_id" gorm:"null"`
	Title              string     `json:"title" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	Description        *string    `json:"description" gorm:"type:nvarchar(max)"`
	Location           *string    `json:"location" gorm:"type:varchar(100);null"`
	Date               time.Time  `json:"date" binding:"required" gorm:"type:date;not null"`
	StartTime          *string    `json:"start_time" gorm:"type:varchar(45);null"`
	EndTime            *string    `json:"end_time" gorm:"type:varchar(45);null"`
	OrganizerID        uint       `json:"organizer_id" binding:"required" gorm:"not null"`
	IsPublic           bool       `json:"is_public" gorm:"not null;default:false"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID        *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	MeetingStatusID    uint       `json:"meeting_status_id" binding:"required" gorm:"not null"`

	// Relationships
	Agendas                []*MeetingAgenda                `gorm:"foreignKey:MeetingID"`
	Controls               []*MeetingControl               `gorm:"foreignKey:MeetingID"`
	ComplianceRequirements []*MeetingComplianceRequirement `gorm:"foreignKey:MeetingID"`
	Risks                  []*MeetingRisk                  `gorm:"foreignKey:MeetingID"`
	Processes              []*MeetingProcess               `gorm:"foreignKey:MeetingID"`
	Incidents              []*MeetingIncident              `gorm:"foreignKey:MeetingID"`
	MeetingType            *MeetingType                    `gorm:"foreignKey:MeetingTypeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Meeting                *Meeting                        `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingFrequency       *MeetingFrequency               `gorm:"foreignKey:MeetingFrequencyID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Organizer              *User                           `gorm:"foreignKey:OrganizerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	MeetingStatus          *MeetingStatus                  `gorm:"foreignKey:MeetingStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                           `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                           `gorm:"foreignKey:UpdatedByID"`
	Participants           []*MeetingParticipant           `gorm:"foreignKey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

// TableName overrides the table name used by Meeting to `meetings`
func (Meeting) TableName() string {
	return "meetings"
}
