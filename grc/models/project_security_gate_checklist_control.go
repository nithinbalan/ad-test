package models

import "time"

type ProjectSecurityGateChecklistControl struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ChecklistID uint      `json:"checklist_id" gorm:"not null"`
	ControlID   uint      `json:"control_id" gorm:"not null"`
	CreatedBYID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Checklist     *ProjectSecurityGateChecklist `json:"checklist" gorm:"foreignKey:ChecklistID"`
	Control       *Control                      `json:"control" gorm:"foreignKey:ControlID"`
	CreatedByUser *User                         `json:"created_by_user" gorm:"foreignKey:CreatedBYID"`
}

func (ProjectSecurityGateChecklistControl) TableName() string {
	return "project_security_gate_checklist_controls"
}
