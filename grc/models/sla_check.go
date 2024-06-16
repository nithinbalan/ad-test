package models

import (
	"time"
)

type SLACheck struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLARegisterID uint       `json:"sla_register_id" binding:"required" gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	CreatedByID   uint       `json:"created_by_id" binding:"required" gorm:"not null"`
	UpdatedBy     *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`

	// Relationships
	SLARegister   *SLARegister `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by SLACheck to `sla_checks`
func (SLACheck) TableName() string {
	return "sla_checks"
}
