package models

type ViolationReportCategory struct {
	ID            uint    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	DepartmentID  uint    `json:"department_id" binding:"required" gorm:"type:int;not null"`
	TitleEn       string  `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	DescriptionEn *string `json:"description_en" gorm:"type:nvarchar(MAX)"`
	DescriptionAr *string `json:"description_ar" gorm:"type:nvarchar(MAX)"`
	Status        string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);default:active"`

	// Relationships
	Department *OrganizationUnit `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by ViolationReportCategory to `violation_report_categories`
func (ViolationReportCategory) TableName() string {
	return "violation_report_categories"
}
