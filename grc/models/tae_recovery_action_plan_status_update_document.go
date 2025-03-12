package models

import (
	"time"
)

type TaeRecoveryActionPlanStatusUpdateDocument struct {
	ID                                  uint      `json:"id" gorm:"primary_key;auto_increment"`
	TaeRecoveryActionPlanStatusUpdateID uint      `json:"tae_recovery_action_plan_status_update_id"`
	URL                                 string    `json:"url"`
	Name                                string    `json:"name" binding:"required"`
	Size                                int       `json:"size"`
	Token                               string    `json:"token"`
	OriginalName                        string    `json:"original_name"`
	Type                                string    `json:"type"`
	CreatedByID                         uint      `json:"created_by_id" gorm:"not null"`
	UpdatedByID                         uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                           time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                           time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	TaeRecoveryActionPlanStatusUpdate *TaeRecoveryActionPlanStatusUpdate `gorm:"foreignKey:TaeRecoveryActionPlanStatusUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser                     *User                              `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser                     *User                              `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by TaeRecoveryActionPlanStatusUpdateDocument to `tae_recovery_action_plan_status_update_documents`
func (TaeRecoveryActionPlanStatusUpdateDocument) TableName() string {
	return "tae_recovery_action_plan_status_update_documents"
}
