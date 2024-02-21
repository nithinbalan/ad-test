package models

type ComplianceCategory struct {
	ID      int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:nvarchar(50);not null"`
}

// TableName overrides the table name used by ComplianceCategory to `compliance_categories`
func (ComplianceCategory) TableName() string {
	return "compliance_categories"
}
