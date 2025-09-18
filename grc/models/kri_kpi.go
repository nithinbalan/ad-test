package models

import (
	"time"
)

type KriKpi struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID uint      `json:"kpi_register_id" binding:"required"`
	KRIRegisterID uint      `json:"kri_register_id" binding:"required"`
	CreatedByID   uint      `json:"created_by_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister   *KPIRegister `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	KRIRegister   *KRIRegister `gorm:"foreignKey:KRIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID"`
}

func (KriKpi) TableName() string {
	return "kri_kpis"
}
