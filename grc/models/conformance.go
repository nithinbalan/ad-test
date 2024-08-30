package models

type Conformance struct {
	ID        uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn   string `json:"title_en" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	TitleAr   string `json:"title_ar" binding:"required,max=45" gorm:"type:nvarchar(45);not null"`
	Type      string `json:"type" binding:"required,max=20,oneof=yes no not-applicable" gorm:"type:varchar(20);not null"`
	ColorCode string `json:"color_code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
}

// TableName overrides the table name used by Conformance to `conformances`
func (Conformance) TableName() string {
	return "conformances"
}
