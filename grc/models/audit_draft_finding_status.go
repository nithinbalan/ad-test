package models

type AuditDraftFindingStatus struct {
	ID         uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn    string `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr    string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Type       string `json:"type" binding:"required,max=20,oneof=open accepted challenged dismissed" gorm:"type:varchar(20);not null"`
	ColorCode  string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	IsFinalized bool  `json:"is_finalized" binding:"required" gorm:"type:tinyint(1);not null;default:0"`
}

// TableName overrides the table name used by AuditDraftFindingStatus to `audit_draft_finding_statuses`
func (AuditDraftFindingStatus) TableName() string {
	return "audit_draft_finding_statuses"
}
