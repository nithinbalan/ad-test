package models

import "time"

type RiskRegisterUpdatePlanTeamMember struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterUpdatePlanID uint       `json:"risk_register_update_plan_id" binding:"required"`
	UserID                   uint       `json:"user_id" binding:"required"`
	CreatedBYID              *uint      `json:"created_by_id" gorm:"default:null"`
	UpdatedBYID              *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt                time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	RiskRegisterUpdatePlan *RiskRegisterUpdatePlan `gorm:"foreignKey:RiskRegisterUpdatePlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	User                   *User                   `gorm:"foreignKey:UserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (RiskRegisterUpdatePlanTeamMember) TableName() string {
	return "risk_register_update_plan_team_members"
}
