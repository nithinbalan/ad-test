package models

import (
	"time"
)

type IncidentActionPlanStatusUpdateDocument struct {
	ID                               uint      `json:"id" gorm:"primary_key;auto_increment"`
	IncidentActionPlanStatusUpdateID uint      `json:"incident_action_plan_status_update_id"`
	URL                              string    `json:"url"`
	Name                             string    `json:"name" binding:"required"`
	Size                             int       `json:"size"`
	Token                            string    `json:"token"`
	OriginalName                     string    `json:"original_name"`
	Type                             string    `json:"type"`
	CreatedBy                        uint      `json:"created_by" gorm:"not null"`
	UpdatedBy                        uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt                        time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                        time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentActionPlanStatusUpdate *IncidentActionPlanStatusUpdate `gorm:"foreignKey:IncidentActionPlanStatusUpdateID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedByUser                  *User                           `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser                  *User                           `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by IncidentActionPlanStatusUpdateDocument to `incident_action_plan_status_update_documents`
func (IncidentActionPlanStatusUpdateDocument) TableName() string {
	return "incident_action_plan_status_update_documents"
}
