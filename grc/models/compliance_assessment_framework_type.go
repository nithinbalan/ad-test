package models

type ComplianceAssessmentFrameworkType struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Type          string `json:"type" binding:"required,max=25,oneof=binary maturity point percentage" gorm:"type:varchar(25);not null"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string `json:"description_en" binding:"required,max=255" gorm:"type:text;not null"`
	DescriptionAr string `json:"description_ar" binding:"required,max=255" gorm:"type:text;not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ComplianceAssessmentFrameworkType to `compliance_assessment_framework_types`
func (ComplianceAssessmentFrameworkType) TableName() string {
	return "compliance_assessment_framework_types"
}
