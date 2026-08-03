package models

type ReviewType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,oneof=user privileged service external" gorm:"type:varchar(50);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null"`
}

func (ReviewType) TableName() string {
	return "review_types"
}
