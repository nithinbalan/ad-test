package models

import (
	"time"
)

type StakeholderActionPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StakeholderID      uint       `json:"bcp_stakeholder_id" binding:"required"`
	SeverityLevelID    *uint      `json:"severity_level_id"`
	Title              string     `json:"title"`
	Description        *string    `json:"description"`
	Requirement        *string    `json:"requirement"`
	ResponsibleUserID  *uint      `json:"responsible_user_id" gorm:"not null" binding:"required"`
	ResponsibleUnitID  *uint      `json:"responsible_unit_id" gorm:"not null" binding:"required"`
	TargetDate         *time.Time `json:"target_date"`
	ActionPlanStatusID *uint      `json:"action_plan_status_id"`
	AssignedOwnerID    *uint      `json:"assigned_owner_id" gorm:"default:null"`
	CreatedBYID        uint       `json:"created_by_id" binding:"required"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedBYID        *uint      `json:"updated_by_id"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`

	// Relationships
	Stakeholder      *Stakeholder                     `gorm:"foreignKey:StakeholderID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel    *SeverityLevel                   `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus *ActionPlanStatus                `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser  *User                            `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUnit  *OrganizationUnit                `gorm:"foreignKey:ResponsibleUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	AssignedOwner    *User                            `gorm:"foreignKey:AssignedOwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"` // New foreign key relationship
	CreatedByUser    *User                            `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser    *User                            `gorm:"foreignKey:UpdatedBYID"`
	Documents        []*StakeholderActionPlanDocument `gorm:"foreignKey:StakeholderActionPlanID"`
}

func (StakeholderActionPlan) TableName() string {
	return "stakeholder_action_plans"
}
