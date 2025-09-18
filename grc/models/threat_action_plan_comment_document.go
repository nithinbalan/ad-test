package models

import (
	"time"
)

// ThreatActionPlanCommentDocument represents a document attached to a Threat action plan comment.
type ThreatActionPlanCommentDocument struct {
	ID                            uint      `json:"id" gorm:"primary_key;auto_increment"`
	ThreatActionPlanCommentID uint      `json:"Threat_action_plan_comment_id" gorm:"not null"`
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
	ThreatActionPlanComment *ThreatActionPlanComment `gorm:"foreignKey:ThreatActionPlanCommentID;constraint:OnUpdate:RESTRICT,OnDelete:CASCADE;"`
	CreatedByUser           *User                    `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                    `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the default table name
func (ThreatActionPlanCommentDocument) TableName() string {
	return "Threat_action_plan_comment_documents"
}
