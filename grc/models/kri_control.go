package models

import (
	"time"
)

type KRIControl struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	KRIRegisterID uint      `json:"kri_register_id" binding:"required"`
	ControlID     uint      `json:"control_id" binding:"required"`
	CreatedBYID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KRIRegister   *KRIRegister `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Control       *Control     `gorm:"foreignKey:ControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID"`
}

func (KRIControl) TableName() string {
	return "kri_controls"
}
