package models

import (
	"time"
)

type BCPKPI struct {
	ID                       uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KPIRegisterID            uint      `json:"kpi_register_id" binding:"required"`
	BusinessContinuityPlanID uint      `json:"business_continuity_plan_id" binding:"required"`
	CreatedByID              uint      `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	KPIRegister            *KPIRegister            `gorm:"foreignKey:KPIRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID"`
}

func (BCPKPI) TableName() string {
	return "bcp_kpis"
}
