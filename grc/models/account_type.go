package models

type AccountType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,oneof=human admin service vendor shared" gorm:"type:varchar(20);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null"`
}

func (AccountType) TableName() string {
	return "account_types"
}
