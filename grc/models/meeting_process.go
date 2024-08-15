package models

import (
	"time"
)

type MeetingProcess struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID   uint      `json:"meeting_id" binding:"required"`
	ProcessID   uint      `json:"process_id" binding:"required"`
	CreatedBYID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Meeting       *Meeting `gorm:"foreignkey:MeetingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process       *Process `gorm:"foreignkey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User    `gorm:"foreignkey:CreatedBYID"`
}

func (MeetingProcess) TableName() string {
	return "meeting_processes"
}
