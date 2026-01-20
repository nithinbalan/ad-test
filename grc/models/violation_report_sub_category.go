package models

type ViolationReportSubCategory struct {
	ID                        uint    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	ViolationReportCategoryID uint    `json:"violation_report_category_id" binding:"required" gorm:"type:int;not null"`
	TitleEn                   string  `json:"title_en" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	TitleAr                   string  `json:"title_ar" binding:"required,max=500" gorm:"type:nvarchar(500);not null"`
	DescriptionEn             *string `json:"description_en" gorm:"type:nvarchar(MAX)"`
	DescriptionAr             *string `json:"description_ar" gorm:"type:nvarchar(MAX)"`
	Status                    string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);default:active"`

	// Relationships
	ViolationReportCategory *ViolationReportCategory `gorm:"foreignKey:ViolationReportCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by ViolationReportSubCategory to `violation_report_sub_categories`
func (ViolationReportSubCategory) TableName() string {
	return "violation_report_sub_categories"
}
