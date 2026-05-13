package models

type BusinessContinuityDurationUnit struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	Minutes int    `json:"minutes" binding:"required" gorm:"type:int;not null"`
	Type    string `json:"type" binding:"required,oneof=minutes hours days weeks" gorm:"type:varchar(100);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`
}

func (BusinessContinuityDurationUnit) TableName() string {
	return "business_continuity_duration_units"
}
