package models

type SwotPestelCategory struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	ColorCode string `json:"color_code" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Type      string `json:"type" binding:"required,max=25,oneof=new in_progress closed" gorm:"type:varchar(25);not null"`
	IsSwot    bool   `json:"is_swot" gorm:"type:tinyint(1);null;default:0"`
	IsPestel  bool   `json:"is_pestel" gorm:"type:tinyint(1);null;default:0"`
	Status    string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by SwotPestelCategory to `swot_pestel_categories`
func (SwotPestelCategory) TableName() string {
	return "swot_pestel_categories"
}
