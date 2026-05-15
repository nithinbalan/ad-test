package models

import "time"

type SystemRegisterSLA struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SystemRegisterID uint      `gorm:"not null" json:"system_register_id"`
	SLARegisterID    uint      `gorm:"not null" json:"sla_register_id"`
	CreatedBYID      uint      `gorm:"not null" json:"created_by_id"`
	CreatedAt        time.Time `gorm:"not null" json:"created_at"`

	SystemRegister *SystemRegister `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"system_register,omitempty"`
	SLARegister    *SLARegister    `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"sla_register,omitempty"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"created_by_user,omitempty"`
}

func (SystemRegisterSLA) TableName() string {
	return "system_register_slas"
}
