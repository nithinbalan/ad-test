package models

type ProjectSecurityGateChecklistTemplateRequirement struct {
	ID                uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateVersionID uint   `json:"template_version_id" gorm:"not null"`
	TitleEn           string `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr           string `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`

	// Relationships
	TemplateVersion *ProjectSecurityGateChecklistTemplateVersion `json:"template_version" gorm:"foreignKey:TemplateVersionID"`
}

func (ProjectSecurityGateChecklistTemplateRequirement) TableName() string {
	return "project_security_gate_checklist_template_requirements"
}
