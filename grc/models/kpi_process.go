package models

import (
	"time"
)

type KPIProcess struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID uint      `json:"kpi_register_id" binding:"required"`
	ProcessID     uint      `json:"process_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister   *KPIRegister `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process       *Process     `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (KPIProcess) TableName() string {
	return "kpi_processes"
}
