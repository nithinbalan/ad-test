package models

type WorkflowCustomRole struct {
	ID               uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn          string `json:"title_en" gorm:"type:nvarchar(255);not null"`
	TitleAr          string `json:"title_ar" gorm:"type:nvarchar(255);not null"`
	Type             string `json:"type" gorm:"type:nvarchar(50);not null;check:type IN ('audit_group_head','audit_team_lead','audit_members')"`
	ModuleIdentifier string `json:"module_identifier" gorm:"type:nvarchar(100);not null"`
}

// TableName overrides the default table name used by GORM
func (WorkflowCustomRole) TableName() string {
	return "workflow_custom_roles"
}
