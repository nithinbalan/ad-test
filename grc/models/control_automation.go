package models

import "time"

type ControlAutomation struct {
	ID          uint      `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn     string    `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr     string    `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type        string    `json:"type" binding:"required,max=50" gorm:"type:nvarchar(50);not null"`
	Score       float64   `json:"score" gorm:"type:decimal(4,2);not null"`
	Status      string    `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;not null"`
	CreatedBYID uint      `json:"created_by_id" gorm:"column:created_by_id;not null"`

	//Relationships
	CreatedByUser *User `json:"created_by_user" gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ControlAutomation) TableName() string {
	return "control_automations"
}
