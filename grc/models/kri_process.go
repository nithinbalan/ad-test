package models

import (
	"time"
)

type KriProcess struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ProcessID     uint      `json:"process_id" binding:"required"`
	KRIRegisterID uint      `json:"kri_register_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Process       *Process     `gorm:"foreignKey:ProcessID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KRIRegister   *KRIRegister `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (KriProcess) TableName() string {
	return "kri_processes"
}
