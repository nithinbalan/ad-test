package models

import "time"

type MeetingParticipant struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	MeetingID   uint      `json:"meeting_id"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedByID uint      `json:"created_by_id"`

	Meeting *Meeting `gorm:"foreignKey:MeetingID"`
	User    *User    `gorm:"foreignKey:UserID"`
}

func (MeetingParticipant) TableName() string {
	return "meeting_participants"
}
