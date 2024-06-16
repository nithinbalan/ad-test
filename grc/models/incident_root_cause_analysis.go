package models

import (
	"time"
)

type IncidentRootCauseAnalysis struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	IncidentRegisterID uint       `json:"incident_register_id" binding:"required" gorm:"not null"`
	Why                uint       `json:"why" binding:"required" gorm:"not null"`
	Question           string     `json:"question" binding:"required" gorm:"not null"`
	Answer             string     `json:"answer,omitempty" gorm:"default:null"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy          *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser    *User             `gorm:"foreignKey:UpdatedBy"`
}

func (IncidentRootCauseAnalysis) TableName() string {
	return "incident_root_cause_analyses"
}
