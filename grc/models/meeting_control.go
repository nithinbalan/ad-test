package models

import (
	"time"
)

type MeetingControl struct {
	ID        uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID uint      `json:"meeting_id" binding:"required"`
	ControlID uint      `json:"control_id" binding:"required"`
	CreatedBy uint      `json:"created_by" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Meeting       *Meeting `gorm:"foreignkey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control `gorm:"foreignkey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User    `gorm:"foreignkey:CreatedBy"`
}

func (MeetingControl) TableName() string {
	return "meeting_controls"
}
