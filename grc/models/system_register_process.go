package models

import "time"

type SystemRegisterProcess struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SystemRegisterID uint      `gorm:"not null" json:"system_register_id"`
	ProcessID        uint      `gorm:"not null" json:"process_id"`
	CreatedBYID      uint      `gorm:"not null" json:"created_by_id"`
	CreatedAt        time.Time `gorm:"not null" json:"created_at"`

	SystemRegister *SystemRegister `gorm:"foreignKey:SystemRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"system_register,omitempty"`
	Process        *Process        `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"process,omitempty"`
	CreatedByUser  *User           `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION" json:"created_by_user,omitempty"`
}

func (SystemRegisterProcess) TableName() string {
	return "system_register_processes"
}
