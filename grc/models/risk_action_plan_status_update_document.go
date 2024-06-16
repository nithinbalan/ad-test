package models

import (
	"time"
)

type RiskActionPlanStatusUpdateDocument struct {
	ID                           uint      `json:"id" gorm:"primary_key;auto_increment"`
	RiskActionPlanStatusUpdateID uint      `json:"risk_action_plan_status_update_id"`
	URL                          string    `json:"url"`
	Name                         string    `json:"name" binding:"required"`
	Size                         int       `json:"size"`
	Token                        string    `json:"token"`
	Type                         string    `json:"type"`
	CreatedByID                  uint      `json:"created_by_id" gorm:"not null"`
	UpdatedBy                    uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                    time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                    time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskActionPlanStatusUpdate *RiskActionPlanStatusUpdate `gorm:"foreignKey:RiskActionPlanStatusUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser              *User                       `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser              *User                       `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by RiskActionPlanStatusUpdateDocument to `risk_action_plan_status_update_documents`
func (RiskActionPlanStatusUpdateDocument) TableName() string {
	return "risk_action_plan_status_update_documents"
}
