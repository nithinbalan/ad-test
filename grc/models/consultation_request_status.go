package models

import "time"

type ConsultationRequestStatus struct {
	ID            uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn       string    `json:"title_en" binding:"required,max=100" gorm:"type:varchar(100);not null"`
	TitleAr       string    `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn string    `json:"description_en" binding:"required" gorm:"type:text"`
	DescriptionAr string    `json:"description_ar" binding:"required" gorm:"type:text"`
	Type          string    `json:"type" binding:"required,max=25" gorm:"type:varchar(25);not null"`
	Status        string    `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:datetime;not null"`
}

func (ConsultationRequestStatus) TableName() string {
	return "consultation_request_statuses"
}
