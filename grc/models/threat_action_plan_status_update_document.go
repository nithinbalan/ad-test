package models

import (
	"time"
)

type ThreatActionPlanStatusUpdateDocument struct {
	ID                             uint       `json:"id" gorm:"primary_key;auto_increment"`
	ThreatActionPlanStatusUpdateID uint       `json:"threat_action_plan_status_update_id"`
	URL                            string     `json:"url"`
	Name                           string     `json:"name" binding:"required"`
	Size                           int        `json:"size"`
	Token                          string     `json:"token"`
	OriginalName                   string     `json:"original_name"`
	Type                           string     `json:"type"`
	CreatedByID                    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID                    *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                      time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                      *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	ThreatActionPlanStatusUpdate *ThreatActionPlanStatusUpdate `gorm:"foreignKey:ThreatActionPlanStatusUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser                *User                         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                *User                         `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ThreatActionPlanStatusUpdateDocument to `threat_action_plan_status_update_documents`
func (ThreatActionPlanStatusUpdateDocument) TableName() string {
	return "threat_action_plan_status_update_documents"
}
