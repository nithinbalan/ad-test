package dtos

type CompetencyFormDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	TitleAr string `json:"title_ar" binding:"required,max=255" gorm:"type:nvarchar(100);not null"`
	Status  string `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null"`
}
