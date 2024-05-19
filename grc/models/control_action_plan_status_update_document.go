package models

import (
	"time"
)

// ControlActionPlanStatusUpdateDocument represents the control_action_plan_status_update_documents table in the database.
type ControlActionPlanStatusUpdateDocument struct {
	ID                              uint       `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	ControlActionPlanStatusUpdateID uint       `json:"control_action_plan_status_update_id" gorm:"index not null" binding:"required"`
	URL                             string     `json:"url" gorm:"size:1000;not null" binding:"required"`
	Name                            string     `json:"name" gorm:"type:nvarchar(255)" binding:"required"`
	Size                            int        `json:"size" gorm:"not null" binding:"required"`
	Token                           string     `json:"token" gorm:"size:255;not null" binding:"required"`
	Type                            string     `json:"type" gorm:"size:255;not null" binding:"required"`
	Icon                            string     `json:"icon" gorm:"size:10;not null" binding:"required,oneof=pdf excel ppt image text other"`
	CreatedByID                     uint       `json:"created_by" gorm:"default:null; column:created_by" binding:"-"`
	UpdatedByID                     *uint      `json:"updated_by" gorm:"default:null; column:updated_by" binding:"-"`
	CreatedAt                       time.Time  `json:"created_at" gorm:"not null" binding:"required"`
	UpdatedAt                       *time.Time `json:"updated_at" gorm:"default:null" binding:"-"`

	// Relationships
	ControlActionPlanStatusUpdate *ControlActionPlanStatusUpdate `json:"control_action_plan_status_update" gorm:"foreignKey:ControlActionPlanStatusUpdateID"`
	Creator                       *User                          `json:"creator" gorm:"foreignKey:CreatedByID"`
	Updater                       *User                          `json:"updater" gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by GORM
func (ControlActionPlanStatusUpdateDocument) TableName() string {
	return "control_action_plan_status_update_documents"
}
