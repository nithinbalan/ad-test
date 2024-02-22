package models

type ComplianceRequirementCategory struct {
	ID      int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Status    string `json:"status" binding:"required,oneof=active inactive" gorm:"type:nvarchar(50);not null"`
}

// TableName overrides the table name used by ComplianceRequirementCategory to `compliance_requirement_categories`
func (ComplianceRequirementCategory) TableName() string {
	return "compliance_requirement_categories"
}
