package models

type StrategicObjectiveCategory struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" binding:"required" gorm:"type:nvarchar(250);not null"`
	TitleAr string `json:"title_ar" binding:"required" gorm:"type:nvarchar(250);not null"`
	Status  string `json:"status" gorm:"type:varchar(10);not null;default:'active'"`
}

// TableName overrides the default table name
func (StrategicObjectiveCategory) TableName() string {
	return "strategic_objective_categories"
}
