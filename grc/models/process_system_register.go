package models

import (
	"time"
)

type ProcessSystemRegister struct {
	ID               uint      `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ProcessID        uint      `json:"process_id" gorm:"not null"`
	SystemRegisterID uint      `json:"system_register_id" gorm:"not null"`
	CreatedBYID      uint      `json:"created_by_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process        *Process        `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegister *SystemRegister `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProcessSystemRegister) TableName() string {
	return "process_system_registers"
}
