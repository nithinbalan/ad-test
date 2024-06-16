package models

import (
	"time"
)

type Task struct {
	ID                     uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber        string     `json:"reference_number" binding:"required,max=45" gorm:"type:varchar(45);not null;unique"`
	ModuleID               uint       `json:"module_id" binding:"required"`
	ComplianceActionPlanID *uint      `json:"compliance_action_plan_id"`
	SLAActionPlanID        *uint      `json:"sla_action_plan_id"`
	IncidentActionPlanID   *uint      `json:"incident_action_plan_id"`
	ControlActionPlanID    *uint      `json:"control_action_plan_id"`
	Title                  string     `json:"title" binding:"required,max=255"`
	Description            *string    `json:"description" gorm:"type:nvarchar(max)"`
	TargetDate             *time.Time `json:"target_date"`
	URL                    string     `json:"url" binding:"required,max=1000"`
	SeverityLevelID        *uint      `json:"severity_level_id"`
	ResponsibleUserID      uint       `json:"responsible_user_id"`
	TaskStatusID           uint       `json:"task_status_id" binding:"required"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedBy              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBy              uint       `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Module               *Module               `gorm:"foreignkey:ModuleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceActionPlan *ComplianceActionPlan `gorm:"foreignkey:ComplianceActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SLAActionPlan        *SLAActionPlan        `gorm:"foreignkey:SLAActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentActionPlan   *IncidentActionPlan   `gorm:"foreignkey:IncidentActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ControlActionPlan    *ControlActionPlan    `gorm:"foreignkey:ControlActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SeverityLevel        *SeverityLevel        `gorm:"foreignkey:SeverityLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUser      *User                 `gorm:"foreignkey:ResponsibleUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	TaskStatus           *TaskStatus           `gorm:"foreignkey:TaskStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser        *User                 `gorm:"foreignKey:CreatedBy"`
	UpdatedByUser        *User                 `gorm:"foreignKey:UpdatedBy"`
}

// TableName overrides the table name used by Task to `tasks`
func (Task) TableName() string {
	return "tasks"
}
