package models

import (
	"time"
)

type IncidentKPI struct {
	ID                 uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID      uint      `json:"kpi_register_id" binding:"required"`
	IncidentRegisterID uint      `json:"incident_register_id" binding:"required"`
	CreatedByID        uint      `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister      *KPIRegister      `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	IncidentRegister *IncidentRegister `gorm:"foreignKey:IncidentRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser    *User             `gorm:"foreignKey:CreatedByID"`
}

func (IncidentKPI) TableName() string {
	return "incident_kpis"
}
