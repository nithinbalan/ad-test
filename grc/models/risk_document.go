package models

import (
	"time"
)

type RiskDocument struct {
	ID             uint       `json:"id" gorm:"primary_key;auto_increment"`
	RiskRegisterID uint       `json:"risk_register_id"`
	URL            string     `json:"url"`
	Name           string     `json:"name" binding:"required"`
	Size           int        `json:"size"`
	Token          string     `json:"token"`
	OriginalName   string     `json:"original_name"`
	Type           string     `json:"type"`
	CreatedByID    uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy      uint       `json:"updated_by_id" gorm:"default:null"`
	CreatedAt      time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	RiskRegister  *RiskRegister `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User         `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User         `gorm:"foreignKey:UpdatedBy"`
}

func (RiskDocument) TableName() string {
	return "risk_documents"
}
