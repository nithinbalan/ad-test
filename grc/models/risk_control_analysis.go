package models

import (
	"time"
)

type RiskControlAnalysis struct {
	ID                            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type                          string    `json:"type" binding:"required" gorm:"not null;type:varchar(255);check:type in ('inherent', 'residual')"`
	RiskControlID                 uint      `json:"risk_control_id" binding:"required" gorm:"not null"`
	ControlTestID                 uint      `json:"control_test_id" binding:"required" gorm:"default:null"`
	ControlEfficiencyStatusID     *uint     `json:"control_efficiency_status_id" binding:"required" gorm:"not null"`
	ControlImplementationStatusID uint      `json:"control_implementation_status_id" binding:"required" gorm:"not null"`
	Score                         *uint     `json:"score,omitempty" gorm:"default:0" binding:"omitempty"`
	ControlEffectivenessScore     *float64  `json:"control_effectiveness_score" gorm:"type:decimal(3,2);default:null;check:control_effectiveness_score >= 0 AND control_effectiveness_score <= 5" binding:"omitempty"`
	ControlEffectivenessRatingID  *uint     `json:"control_effectiveness_rating_id,omitempty" gorm:"default:null" binding:"omitempty"`
	RiskReduction                 *float64  `json:"risk_reduction,omitempty" gorm:"type:decimal(5,2);default:null;check:risk_reduction >= 0 AND risk_reduction <= 100" binding:"omitempty"`
	CreatedByID                   uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt                     time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskControl                 *RiskControl                 `gorm:"foreignKey:RiskControlID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlEfficiencyStatus     *ControlEfficiencyStatus     `gorm:"foreignKey:ControlEfficiencyStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlImplementationStatus *ControlImplementationStatus `gorm:"foreignKey:ControlImplementationStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ControlEffectivenessRating  *ControlEffectivenessRating  `gorm:"foreignKey:ControlEffectivenessRatingID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser               *User                        `gorm:"foreignKey:CreatedByID"`
}

// TableName overrides the table name used by RiskControlAnalysis to `risk_control_analyses`
func (RiskControlAnalysis) TableName() string {
	return "risk_control_analyses"
}
