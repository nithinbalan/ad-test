package models

type AwarenessMaturityLevel struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,oneof=excellent very_good good acceptable" gorm:"type:varchar(50);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null"`
}

func (AwarenessMaturityLevel) TableName() string {
	return "awareness_maturity_levels"
}