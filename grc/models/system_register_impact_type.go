package models

import (
	"time"
)

type SystemRegisterImpactType struct {
	ID               uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SystemRegisterID uint      `json:"system_register_id" gorm:"not null"`
	ImpactTypeID     uint      `json:"impact_type_id" gorm:"not null"`
	CreatedBYID      *uint     `json:"created_by_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	SystemRegister *SystemRegister  `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ImpactType     *ImpactCategory  `gorm:"foreignKey:ImpactTypeID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (SystemRegisterImpactType) TableName() string {
	return "system_register_impact_types"
}
