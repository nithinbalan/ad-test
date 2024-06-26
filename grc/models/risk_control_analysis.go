package models

import (
	"time"
)

type RiskControlAnalysis struct {
	ID                            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type                          string    `json:"type" binding:"required" gorm:"not null;type:varchar(255);check:type in ('inherent', 'residual')"`
	RiskControlID                 uint      `json:"risk_control_id" binding:"required" gorm:"not null"`
	ControlEfficiencyStatusID     uint      `json:"control_efficiency_status_id" binding:"required" gorm:"not null"`
	ControlImplementationStatusID uint      `json:"control_implementation_status_id" binding:"required" gorm:"not null"`
	CreatedByID                   uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskControl                 *RiskControl                 `gorm:"foreignKey:RiskControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlEfficiencyStatus     *ControlEfficiencyStatus     `gorm:"foreignKey:ControlEfficiencyStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlImplementationStatus *ControlImplementationStatus `gorm:"foreignKey:ControlImplementationStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by RiskControlAnalysis to `risk_control_analyses`
func (RiskControlAnalysis) TableName() string {
	return "risk_control_analyses"
}
