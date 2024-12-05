package models

type AuditFindingCategory struct {
	ID            uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string  `gorm:"type:nvarchar(45);not null" json:"title_en"`
	TitleAr       string  `gorm:"type:nvarchar(45);not null" json:"title_ar"`
	DescriptionEn *string `gorm:"type:nvarchar(max)" json:"description_en"`
	DescriptionAr *string `gorm:"type:nvarchar(max)" json:"description_ar"`
	ColorCode     string  `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type          string  `json:"type" binding:"required,max=20,oneof=major-nc minor-nc ofi observation" gorm:"type:varchar(20);not null"`
}

func (AuditFindingCategory) TableName() string {
	return "audit_finding_categories"
}
