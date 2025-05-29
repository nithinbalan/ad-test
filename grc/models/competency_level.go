package models

type CompetencyLevel struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn string `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	Status  string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null"`
}

func (CompetencyLevel) TableName() string {
	return "competency_levels"
}
