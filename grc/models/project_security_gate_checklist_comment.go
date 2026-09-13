package models

import "time"

type ProjectSecurityGateChecklistComment struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ChecklistID uint      `json:"checklist_id" gorm:"not null"`
	HistoryID   *uint     `json:"history_id" gorm:"default:null"`
	Comment     string    `json:"comment" binding:"required" gorm:"type:nvarchar(max);not null"`
	CreatedBYID uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Checklist     *ProjectSecurityGateChecklist                         `json:"checklist" gorm:"foreignKey:ChecklistID"`
	History       *ProjectSecurityGateChecklistExceptionWorkflowHistory `json:"history" gorm:"foreignKey:HistoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                                                 `json:"created_by_user" gorm:"foreignKey:CreatedBYID"`
	Documents     []ProjectSecurityGateChecklistCommentDocument         `json:"documents" gorm:"foreignKey:ProjectSecurityGateChecklistCommentID"`
}

func (ProjectSecurityGateChecklistComment) TableName() string {
	return "project_security_gate_checklist_comments"
}
