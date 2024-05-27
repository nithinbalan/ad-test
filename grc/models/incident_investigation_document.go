package models

import (
	"time"
)

type IncidentInvestigationDocument struct {
	ID                      uint      `json:"id" gorm:"primary_key;auto_increment"`
	IncidentInvestigationID uint      `json:"incident_investigation_id"`
	URL                     string    `json:"url"`
	Name                    string    `json:"name" binding:"required"`
	Size                    int       `json:"size"`
	Token                   string    `json:"token"`
	OriginalName            string    `json:"original_name"`
	Type                    string    `json:"type"`
	CreatedBy               uint      `json:"created_by" gorm:"not null"`
	UpdatedBy               uint      `json:"updated_by" gorm:"default:null"`
	CreatedAt               time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt               time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	IncidentInvestigation *IncidentInvestigation `gorm:"foreignkey:IncidentInvestigationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBy"`
}

func (IncidentInvestigationDocument) TableName() string {
	return "incident_investigation_documents"
}
