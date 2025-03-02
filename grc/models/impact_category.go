package models

type ImpactCategory struct {
	ID      uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string  `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string  `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type    *string `json:"type" gorm:"type:varchar(20);null"`
	Status  string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by ImpactCategory to `impact_categories`
func (ImpactCategory) TableName() string {
	return "impact_categories"
}
