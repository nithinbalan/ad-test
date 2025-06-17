package models

import (
	"time"
)

type KPIPolicy struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID uint      `json:"kpi_register_id" binding:"required"`
	DocumentID    uint      `json:"document_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister   *KPIRegister `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Document      *Document    `gorm:"foreignKey:DocumentID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (KPIPolicy) TableName() string {
	return "kpi_policies"
}
