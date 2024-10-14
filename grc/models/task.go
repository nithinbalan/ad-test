package models

import (
	"time"
)

type Task struct {
	ID                       uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ReferenceNumber          string     `json:"reference_number" binding:"required,max=45" gorm:"type:varchar(45);not null;unique"`
	ModuleID                 uint       `json:"module_id" binding:"required"`
	ComplianceActionPlanID   *uint      `json:"compliance_action_plan_id"`
	SLAActionPlanID          *uint      `json:"sla_action_plan_id"`
	IncidentActionPlanID     *uint      `json:"incident_action_plan_id"`
	ControlActionPlanID      *uint      `json:"control_action_plan_id"`
	RiskActionPlanID         *uint      `json:"risk_action_plan_id"`
	MeetingActionPlanID      *uint      `json:"meeting_action_plan_id"`
	AuditFindingActionPlanID *uint      `json:"audit_finding_action_plan_id"`
	DocumentID               *uint      `json:"document_id"`
	DocumentWorkflowID       *uint      `json:"document_workflow_id"`
	Title                    string     `json:"title" binding:"required,max=255"`
	Description              *string    `json:"description" gorm:"type:nvarchar(max)"`
	TargetDate               *time.Time `json:"target_date"`
	URLParam                 *string    `json:"url_param" gorm:"type:nvarchar(max)"`
	QueryParam               *string    `json:"query_param" gorm:"type:nvarchar(max)"`
	SeverityLevelID          *uint      `json:"severity_level_id"`
	ResponsibleUserID        uint       `json:"responsible_user_id"`
	TaskStatusID             uint       `json:"task_status_id" binding:"required"`
	CreatedAt                time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt                *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`
	CreatedByID              uint       `json:"created_by_id" gorm:"not null"`
	UpdatedByID              uint       `json:"updated_by_id" gorm:"default:null"`

	// Relationships
	Module                 *Module                 `gorm:"foreignKey:ModuleID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ComplianceActionPlan   *ComplianceActionPlan   `gorm:"foreignKey:ComplianceActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	SLAActionPlan          *SLAActionPlan          `gorm:"foreignKey:SLAActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	IncidentActionPlan     *IncidentActionPlan     `gorm:"foreignKey:IncidentActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ControlActionPlan      *ControlActionPlan      `gorm:"foreignKey:ControlActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskActionPlan         *RiskActionPlan         `gorm:"foreignKey:RiskActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`         // New relationship added
	AuditFindingActionPlan *AuditFindingActionPlan `gorm:"foreignKey:AuditFindingActionPlanID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"` // New relationship added
	SeverityLevel          *SeverityLevel          `gorm:"foreignKey:SeverityLevelID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	ResponsibleUser        *User                   `gorm:"foreignKey:ResponsibleUserID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	TaskStatus             *TaskStatus             `gorm:"foreignKey:TaskStatusID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser          *User                   `gorm:"foreignKey:CreatedByID"`
	UpdatedByUser          *User                   `gorm:"foreignKey:UpdatedByID"`
}

// TableName overrides the table name used by Task to `tasks`
func (Task) TableName() string {
	return "tasks"
}
