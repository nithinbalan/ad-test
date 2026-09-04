package models

import "time"

type ProjectSecurityGateChecklistResponseDocument struct {
	ID           uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ChecklistID  uint       `json:"checklist_id" gorm:"column:project_security_gate_checklist_id;not null"`
	URL          string     `json:"url"`
	Name         string     `json:"name" binding:"required"`
	Size         int        `json:"size"`
	Token        string     `json:"token"`
	OriginalName string     `json:"original_name"`
	Type         string     `json:"type"`
	CreatedBYID  uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID  *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt    time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt    *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	Checklist     *ProjectSecurityGateChecklist `json:"checklist" gorm:"foreignKey:ChecklistID"`
	CreatedByUser *User                         `json:"created_by_user" gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User                         `json:"updated_by_user" gorm:"foreignKey:UpdatedBYID"`
}

func (ProjectSecurityGateChecklistResponseDocument) TableName() string {
	return "project_security_gate_checklist_response_documents"
}
