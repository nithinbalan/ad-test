package models

import "time"

type MeetingParticipant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	MeetingID uint      `json:"meeting_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"`

	Meeting *Meeting `gorm:"foreignkey:MeetingID"`
	User    *User    `gorm:"foreignkey:UserID"`
}

func (MeetingParticipant) TableName() string {
	return "meeting_participants"
}
