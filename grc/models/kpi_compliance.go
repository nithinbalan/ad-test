package models

import (
	"time"
)

type KPICompliance struct {
	ID                   uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID        uint      `json:"kpi_register_id" binding:"required"`
	ComplianceRegisterID uint      `json:"compliance_id" binding:"required"`
	CreatedByID          uint      `json:"created_by_id" binding:"required"`
	CreatedAt            time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister        *KPIRegister        `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ComplianceRegister *ComplianceRegister `gorm:"foreignKey:ComplianceRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser      *User               `gorm:"foreignKey:CreatedByID"`
}

func (KPICompliance) TableName() string {
	return "kpi_compliances"
}
