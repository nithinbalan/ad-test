package models

type ComplianceAssessmentImpactType struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string `json:"description_en" binding:"required" gorm:"type:nvarchar(max);not null"`
	DescriptionAr string `json:"description_ar" binding:"required" gorm:"type:nvarchar(max);not null"`
	Type          string `json:"type" binding:"required,max=10,oneof=regulatory reputational financial operational" gorm:"type:varchar(10);not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ComplianceAssessmentImpactType to `compliance_assessment_impact_types`
func (ComplianceAssessmentImpactType) TableName() string {
	return "compliance_assessment_impact_types"
}
