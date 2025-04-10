package models

import "time"

type BcpRisk struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	BusinessContinuityPlanID uint       `json:"business_continuity_plan_id" binding:"required"`
	RiskRegisterID           uint       `json:"risk_register_id" binding:"required"`
	CreatedBYID              uint       `json:"created_by_id" binding:"required"`
	UpdatedBYID              *uint      `json:"updated_by_id"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                *time.Time `json:"updated_at"`

	// Relationships
	BusinessContinuityPlan *BusinessContinuityPlan `gorm:"foreignKey:BusinessContinuityPlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	RiskRegister           *RiskRegister           `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (BcpRisk) TableName() string {
	return "bcp_risks"
}
