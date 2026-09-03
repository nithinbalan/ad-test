package models

type ProjectSecurityGateChecklistTemplate struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	StageID uint   `json:"stage_id" gorm:"not null"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null"`

	// Relationships
	Stage    *ProjectSecurityGateStage                     `json:"stage" gorm:"foreignKey:StageID"`
	Versions []ProjectSecurityGateChecklistTemplateVersion `json:"versions" gorm:"foreignKey:TemplateID"`
}

func (ProjectSecurityGateChecklistTemplate) TableName() string {
	return "project_security_gate_checklist_templates"
}
