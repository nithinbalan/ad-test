package models

import "time"

type IncidentInvestigationUser struct {
	ID                      uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentInvestigationID uint       `json:"incident_investigation_id" binding:"required" gorm:"not null"`
	UserID                  uint       `json:"user_id" binding:"required" gorm:"not null"`
	CreatedBy               uint       `json:"created_by_id" binding:"required"`
	UpdatedBy               *uint      `json:"updated_by_id"`
	CreatedAt               time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	IncidentInvestigation IncidentInvestigation `gorm:"foreignkey:IncidentInvestigationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	User                  User                  `gorm:"foreignkey:UserID"`
	CreatedByUser         User                  `gorm:"foreignkey:CreatedBy"`
	UpdatedByUser         *User                 `gorm:"foreignkey:UpdatedBy"`
}

// TableName overrides the table name used by IncidentInvestigationUser to `incident_investigation_users`
func (IncidentInvestigationUser) TableName() string {
	return "incident_investigation_users"
}
