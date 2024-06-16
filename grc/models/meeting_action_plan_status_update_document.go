package models

import (
	"time"
)

// MeetingActionPlanStatusUpdateDocument represents the meeting_action_plan_status_update_documents table in the database.
type MeetingActionPlanStatusUpdateDocument struct {
	ID                              uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	MeetingActionPlanStatusUpdateID uint       `json:"meeting_action_plan_status_update_id" gorm:"index not null" binding:"required"`
	URL                             string     `json:"url" gorm:"size:1000;not null" binding:"required"`
	Name                            string     `json:"name" gorm:"type:nvarchar(255)" binding:"required"`
	Size                            int        `json:"size" gorm:"not null" binding:"required"`
	Token                           string     `json:"token" gorm:"size:255;not null" binding:"required"`
	OriginalName                    string     `json:"original_name"`
	Type                            string     `json:"type" gorm:"size:255;not null" binding:"required"`
	CreatedByID                     uint       `json:"created_by_id,omitempty" gorm:"default:null; column:created_by_id" binding:"-"`
	UpdatedByID                     *uint      `json:"updated_by_id,omitempty" gorm:"default:null; column:updated_by_id" binding:"-"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"not null" binding:"required"`
	UpdatedAt                       *time.Time `json:"updated_at,omitempty" gorm:"default:null" binding:"-"`

	// Relationships
	MeetingActionPlanStatusUpdate *MeetingActionPlanStatusUpdate `json:"meeting_action_plan_status_update,omitempty" gorm:"foreignKey:MeetingActionPlanStatusUpdateID"`
	Creator                       *User                          `json:"creator,omitempty" gorm:"foreignKey:CreatedByID"`
	Updater                       *User                          `json:"updater,omitempty" gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by GORM
func (MeetingActionPlanStatusUpdateDocument) TableName() string {
	return "meeting_action_plan_status_update_documents"
}
