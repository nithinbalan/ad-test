package models

import (
	"time"
)

type BcpProcessSystemRegister struct {
	ID               uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BcpProcessID     uint      `json:"bcp_process_id" gorm:"not null"`
	SystemRegisterID uint      `json:"system_register_id" gorm:"not null"`
	CreatedBYID      *uint     `json:"created_by_id"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	BcpProcess     *BcpProcess     `gorm:"foreignKey:BcpProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SystemRegister *SystemRegister `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (BcpProcessSystemRegister) TableName() string {
	return "bcp_process_system_registers"
}
