package models

type ControlType struct {
	ID      uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Type    string `json:"type" binding:"required,max=50" gorm:"type:nvarchar(50);not null" validate:"required,oneof=preventive detective corrective deterrent"`
	Score   int    `json:"score" gorm:"type:int;not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

func (ControlType) TableName() string {
	return "control_types"
}
