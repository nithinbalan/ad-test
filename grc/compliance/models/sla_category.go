package models

type SlaCategory struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:nvarchar(50);not null"`
}

// TableName overrides the table name used by SlaCategory to `sla_categories`
func (SlaCategory) TableName() string {
	return "sla_categories"
}
