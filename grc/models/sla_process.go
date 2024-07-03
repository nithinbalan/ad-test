package models

import (
	"time"
)

type SLAProcess struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	SLARegisterID uint      `json:"sla_register_id" binding:"required"`
	ProcessID     uint      `json:"process_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	SLARegister   *SLARegister `gorm:"foreignKey:SLARegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Process       *Process     `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (SLAProcess) TableName() string {
	return "sla_processes"
}
