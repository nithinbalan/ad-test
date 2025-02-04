package models

import (
	"time"
)

type RiskSetting struct {
	ID           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReviewRoleID uint       `json:"review_role_id" binding:"required" gorm:"not null"`
	CreatedBYID  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID  *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	ReviewRole    *Role `gorm:"foreignKey:ReviewRoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the default table name to `risk_settings`
func (RiskSetting) TableName() string {
	return "risk_settings"
}
