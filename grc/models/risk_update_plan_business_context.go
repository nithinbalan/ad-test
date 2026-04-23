package models

import "time"

type RiskUpdatePlanBusinessContext struct {
	ID                       uint      `json:"id" gorm:"primary_key;auto_increment"`
	RiskRegisterUpdatePlanID uint      `json:"risk_register_update_plan_id" gorm:"not null"`
	BusinessContextID        *uint     `json:"business_context_id" gorm:"default:null"`
	ProcessID                *uint     `json:"process_id" gorm:"default:null"`
	ControlID                *uint     `json:"control_id" gorm:"default:null"`
	IsBusinessContextUpdated bool      `json:"is_business_context_updated" gorm:"default:false;not null"`
	CreatedBYID              *uint     `json:"created_by_id" gorm:"default:null"`
	CreatedAt                time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	RiskRegisterUpdatePlan *RiskRegisterUpdatePlan          `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	BusinessContext        *OrganizationUnitBusinessContext `gorm:"foreignKey:BusinessContextID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Process                *Process                         `gorm:"foreignKey:ProcessID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	Control                *Control                         `gorm:"foreignKey:ControlID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                            `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskUpdatePlanBusinessContext) TableName() string {
	return "risk_update_plan_business_contexts"
}
