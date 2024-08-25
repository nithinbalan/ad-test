package models

import (
	"time"
)

type AuditMeeting struct {
	ID                uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AuditID           uint       `json:"audit_id" binding:"required" gorm:"not null"`
	Title             string     `json:"title" binding:"required" gorm:"type:nvarchar(255);not null"`
	Date              time.Time  `json:"date" binding:"required" gorm:"type:date;not null"`
	FromTime          *string    `json:"from_time" gorm:"type:varchar(45);null"`
	ToTime            *string    `json:"to_time" gorm:"type:varchar(45);null"`
	OtherParticipants *string    `json:"other_participants,omitempty" gorm:"type:nvarchar(max);default:null"`
	Mom               *string    `json:"mom,omitempty" gorm:"type:nvarchar(max);default:null"`
	CreatedBYID       uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBYID       *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt         time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	Audit         *Audit `gorm:"foreignKey:AuditID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION"`
	CreatedByUser *User  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User  `gorm:"foreignKey:UpdatedBYID"`

	Participants []*AuditMeetingParticipant `gorm:"foreignKey:AuditMeetingID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// TableName overrides the default table name used by GORM
func (AuditMeeting) TableName() string {
	return "audit_meetings"
}
