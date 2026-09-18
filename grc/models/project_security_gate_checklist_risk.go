package models

import "time"

type ProjectSecurityGateChecklistRisk struct {
	ID             uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	ChecklistID    uint      `json:"checklist_id" gorm:"not null"`
	RiskRegisterID uint      `json:"risk_register_id" gorm:"not null"`
	CreatedBYID    uint      `json:"created_by_id" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime;not null"`

	// Relationships
	Checklist     *ProjectSecurityGateChecklist `json:"checklist" gorm:"foreignKey:ChecklistID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	RiskRegister  *RiskRegister                 `json:"risk_register" gorm:"foreignKey:RiskRegisterID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User                         `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectSecurityGateChecklistRisk) TableName() string {
	return "project_security_gate_checklist_risks"
}
