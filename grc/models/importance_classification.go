package models

import "time"

type ImportanceClassification struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string    `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr       string    `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string    `json:"description_en" binding:"required,max=255" gorm:"type:varchar(255);not null"`
	DescriptionAr string    `json:"description_ar" binding:"required,max=255" gorm:"type:nvarchar(255);not null"`
	Status        string    `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`
}

func (ImportanceClassification) TableName() string {
	return "importance_classifications"
}
