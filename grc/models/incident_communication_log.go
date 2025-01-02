package models

import (
	"time"
)

type IncidentCommunicationLog struct {
	ID                  uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID  uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	CommunicatedAt      *time.Time `json:"communicated_at" binding:"required" gorm:"type:datetime;not null"`
	CommunicatedBy      string     `json:"communicated_by" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	CommunicatedTo      string     `json:"communicated_to" binding:"required,max=1000" gorm:"type:nvarchar(1000);not null"`
	CommunicationMethod string     `json:"communication_method" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Summary             *string    `json:"summary" gorm:"type:nvarchar(max);null"`
	CreatedAt           time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID         uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID         *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by IncidentCommunicationLog to `incident_communication_logs`
func (IncidentCommunicationLog) TableName() string {
	return "incident_communication_logs"
}
