package models

type Unit struct {
	ID        uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn   string  `json:"title_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	TitleAr   string  `json:"title_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	CodeEn    string  `json:"code_en" binding:"required" gorm:"type:nvarchar(100);not null"`
	CodeAr    string  `json:"code_ar" binding:"required" gorm:"type:nvarchar(100);not null"`
	Type      string  `json:"type" binding:"required,oneof=length area volume mass time temperature speed pressure energy power frequency electrical light data other" gorm:"type:varchar(20);not null;default:'other';check:type IN ('length','area','volume','mass','time','temperature','speed','pressure','energy','power','frequency','electrical','light','data','other')"`
	ColorCode *string `json:"color_code" gorm:"type:varchar(10)"`
	Status    string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active', 'inactive')"`
}

func (Unit) TableName() string {
	return "units"
}
