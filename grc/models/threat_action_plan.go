package models

import (
	"time"
)

type ThreatActionPlan struct {
	ID                 uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ThreatID           uint       `json:"threat_id" binding:"required" gorm:"not null"`
	Title              string     `json:"title" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	Description        *string    `json:"description" gorm:"type:nvarchar(max);null"`
	SeverityLevelID    *uint      `json:"severity_level_id" binding:"omitempty" gorm:"not null"`
	ResponsibleUserID  uint       `json:"responsible_user_id" binding:"required" gorm:"not null"`
	Comments           *string    `json:"comments" gorm:"type:nvarchar(max);null"`
	TargetDate         *time.Time `json:"target_date" binding:"omitempty" gorm:"type:date;default:null"`
	ActionPlanStatusID uint       `json:"action_plan_status_id" binding:"required" gorm:"not null"`
	ReviewUserID       *uint      `json:"review_user_id" gorm:"default:null"`
	Percentage         *uint      `json:"percentage" binding:"omitempty,min=0,max=100" gorm:"not null;default:0"`
	CreatedAt          time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          *time.Time `json:"updated_at" gorm:"type:datetime"`
	CreatedByID        uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID        *uint      `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Threat                  *Threat                         `gorm:"foreignKey:ThreatID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	SeverityLevel           *SeverityLevel                  `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ResponsibleUser         *User                           `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ReviewUser              *User                           `gorm:"foreignKey:ReviewUserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatus        *ActionPlanStatus               `gorm:"foreignKey:ActionPlanStatusID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ActionPlanStatusUpdates []*ThreatActionPlanStatusUpdate `gorm:"foreignKey:ThreatActionPlanID"`
	CreatedByUser           *User                           `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser           *User                           `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by ThreatActionPlan to `threat_action_plans`
func (ThreatActionPlan) TableName() string {
	return "threat_action_plans"
}
