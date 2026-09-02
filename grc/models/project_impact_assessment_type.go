package models

type ProjectImpactAssessmentType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,max=50,oneof=timeline budget resources technical_complexity" gorm:"type:varchar(50);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ProjectImpactAssessmentType to `project_impact_assessment_types`
func (ProjectImpactAssessmentType) TableName() string {
	return "project_impact_assessment_types"
}
