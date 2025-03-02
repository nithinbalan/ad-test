package models

import (
	"time"
)

type MeetingDocument struct {
	ID           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	MeetingID    uint       `json:"meeting_id" binding:"required" gorm:"not null"`
	URL          string     `json:"url" binding:"required" gorm:"type:varchar(1000);not null"`
	Name         string     `json:"name" binding:"required" gorm:"type:nvarchar(255);not null"`
	Size         int        `json:"size" binding:"required" gorm:"not null"`
	Token        string     `json:"token" binding:"required" gorm:"type:varchar(255);not null"`
	OriginalName *string    `json:"original_name" gorm:"type:nvarchar(255);null"`
	Type         string     `json:"type" binding:"required" gorm:"type:varchar(255);not null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"type:datetime;null"`

	// Relationships
	Meeting *Meeting `gorm:"foreignKey:MeetingID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by MeetingDocument to `meeting_documents`
func (MeetingDocument) TableName() string {
	return "meeting_documents"
}
