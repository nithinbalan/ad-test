package models

import "time"

type RiskUpdatePlanErmDeptEvaluation struct {
	ID                       uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	RiskRegisterUpdatePlanID uint      `json:"risk_register_update_plan_id" binding:"required"`
	RiskRegisterID           uint      `json:"risk_register_id" binding:"required"`
	CreatedBYID              uint      `json:"created_by_id" binding:"required"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegisterUpdatePlan *RiskRegisterUpdatePlan `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskRegister           *RiskRegister           `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by GORM
func (RiskUpdatePlanErmDeptEvaluation) TableName() string {
	return "risk_update_plan_erm_dept_evaluations"
}
