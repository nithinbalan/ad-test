package models

import "time"

type ProjectSecurityGateChecklistTemplateVersion struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID    uint       `json:"template_id" gorm:"not null"`
	Version       int        `json:"version" gorm:"not null"`
	Status        string     `json:"status" gorm:"type:varchar(10);not null"`
	CreatedByID   uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	PublishedByID *uint      `json:"published_by_id" gorm:"default:null"`
	PublishedAt   *time.Time `json:"published_at" gorm:"type:datetime;default:null"`

	// Relationships
	Template        *ProjectSecurityGateChecklistTemplate             `json:"template" gorm:"foreignKey:TemplateID"`
	CreatedByUser   *User                                             `json:"created_by_user" gorm:"foreignKey:CreatedByID"`
	PublishedByUser *User                                             `json:"published_by_user" gorm:"foreignKey:PublishedByID"`
	Requirements    []ProjectSecurityGateChecklistTemplateRequirement `json:"requirements" gorm:"foreignKey:TemplateVersionID"`
}

func (ProjectSecurityGateChecklistTemplateVersion) TableName() string {
	return "project_security_gate_checklist_template_versions"
}
