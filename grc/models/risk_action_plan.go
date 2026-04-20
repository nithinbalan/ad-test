package models

import (
	"time"
)

type RiskActionPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskRegisterID     uint       `json:"risk_register_id" binding:"required" gorm:"not null"`
	Title              string     `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description        *string    `json:"description" gorm:"type:nvarchar(max);null"`
	SeverityLevelID    *uint      `json:"severity_level_id" binding:"omitempty" gorm:"not null"`
	ResponsibleUnitID  *uint      `json:"responsible_unit_id" binding:"required" gorm:"not null"`
	ResponsibleUserID  uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	Comments           *string    `json:"comments" gorm:"type:nvarchar(max);null"`
	StartDate          *time.Time `json:"start_date" binding:"omitempty" gorm:"type:date,default:null"`
	TargetDate         *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date,default:null"`
	Cost               *int       `json:"cost" gorm:"type:int;default:null"`
	ActionPlanStatusID uint       `json:"action_plan_status_id" binding:"required" gorm:"not null"`
	ReviewUserID       *uint      `json:"review_user_id" gorm:"default:null"`
	Percentage         *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	Notes              *string    `json:"notes" gorm:"type:nvarchar(max);default:null"`
	DelegateeUserID    *uint      `json:"delegatee_user_id"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedBYID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID        *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	RiskRegister            *RiskRegister                 `gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel           *SeverityLevel                `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser         *User                         `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit         *OrganizationUnit             `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser              *User                         `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // New foreign key relationship
	ActionPlanStatus        *ActionPlanStatus             `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatusUpdates []*RiskActionPlanStatusUpdate `gorm:"foreignKey:RiskActionPlanID"`
	Controls                []*RiskActionPlanControl      `gorm:"foreignKey:RiskActionPlanID"`
	DelegateeUser           *User                         `gorm:"foreignKey:DelegateeUserID"`
	CreatedByUser           *User                         `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser           *User                         `gorm:"foreignKey:UpdatedBYID"`
}

// TableName overrides the table name used by RiskActionPlan to `risk_action_plans`
func (RiskActionPlan) TableName() string {
	return "risk_action_plans"
}
