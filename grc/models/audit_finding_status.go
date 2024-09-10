package models

type AuditFindingStatus struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Type      string `json:"type" binding:"required,max=20,oneof=new in-progress resolved reopened closed deferred" gorm:"type:varchar(20);not null"`
	ColorCode string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	IsClosed  bool   `json:"is_closed" binding:"required" gorm:"type:tinyint(1);not null;default:0"`
}

// TableName overrides the table name used by AuditFindingStatus to `audit_finding_statuses`
func (AuditFindingStatus) TableName() string {
	return "audit_finding_statuses"
}
