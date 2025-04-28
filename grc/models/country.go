package models

type Country struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string  `json:"title_en" gorm:"type:nvarchar(250);not null" binding:"required"`
	TitleAr   string  `json:"title_ar" gorm:"type:nvarchar(250);not null" binding:"required"`
	Code      string  `json:"code" gorm:"type:varchar(10);not null" binding:"required"`
	PhoneCode *string `json:"phone_code,omitempty" gorm:"type:varchar(10);default:null"`
	Status    string  `json:"status" gorm:"type:varchar(10);not null;default:'active'"`
}

// TableName overrides the default table name.
func (Country) TableName() string {
	return "countries"
}
