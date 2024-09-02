package models

type ComplianceStatus struct {
	ID            uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr       string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string `json:"description_en" binding:"required" gorm:"type:text"`
	DescriptionAr string `json:"description_ar" binding:"required" gorm:"type:text"`
	ColorCode     string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type          string `json:"type" binding:"required,max=25,oneof=not-implemented in_progress partially_compliant compliant exempt not_applicable" gorm:"type:varchar(25);not null"`
	Status        string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	BaseScore     int    `json:"base_score" gorm:"not null;default:0"`

}

// TableName overrides the table name used by ComplianceStatus to `compliance_statuses`
func (ComplianceStatus) TableName() string {
	return "compliance_statuses"
}
