package models

import (
	"time"
)

// ControlActionPlanCommentDocument represents a document attached to a control action plan comment.
type ControlActionPlanCommentDocument struct {
	ID                            uint      `json:"id" gorm:"primary_key;auto_increment"`
	ControlActionPlanCommentID uint      `json:"control_action_plan_comment_id" gorm:"not null"`
	URL                           string    `json:"url"`
	Name                          string    `json:"name" binding:"required"`
	Size                          int       `json:"size"`
	Token                         string    `json:"token"`
	OriginalName                  string    `json:"original_name"`
	Type                          string    `json:"type"`
	CreatedBYID                   uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBYID                   *uint     `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                     time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                     time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	ControlActionPlanComment *ControlActionPlanComment `gorm:"foreignKey:ControlActionPlanCommentID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser               *User                        `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the default table name
func (ControlActionPlanCommentDocument) TableName() string {
	return "control_action_plan_comment_documents"
}
