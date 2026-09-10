package models

import "time"

type ProjectSecurityGateChecklistExceptionWorkflowReviewDocument struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	HistoryID    uint      `json:"history_id" gorm:"not null"`
	URL          string    `json:"url"`
	Name         string    `json:"name" binding:"required"`
	Size         int       `json:"size"`
	Token        string    `json:"token"`
	OriginalName string    `json:"original_name"`
	Type         string    `json:"type"`
	CreatedBYID  uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships.
	History       *ProjectSecurityGateChecklistExceptionWorkflowHistory `json:"history" gorm:"foreignKey:HistoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                                                 `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateChecklistExceptionWorkflowReviewDocument) TableName() string {
	return "project_security_gate_checklist_exception_workflow_review_documents"
}
