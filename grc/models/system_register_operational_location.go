package models

import (
	"time"
)

type SystemRegisterOperationalLocation struct {
	ID                    uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SystemRegisterID      uint      `json:"system_register_id" binding:"required" gorm:"not null"`
	OperationalLocationID uint      `json:"operational_location_id" binding:"required" gorm:"not null"`
	CreatedBYID           uint      `json:"created_by_id" binding:"required" gorm:"not null"`
	CreatedAt             time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	SystemRegister      *SystemRegister `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	OperationalLocation *Location       `gorm:"foreignKey:OperationalLocationID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser       *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (SystemRegisterOperationalLocation) TableName() string {
	return "system_register_operational_locations"
}
