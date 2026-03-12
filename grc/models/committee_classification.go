package models

type CommitteeClassification struct {
	ID            uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn       string  `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr       string  `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	DescriptionEn *string `json:"description_en" gorm:"type:nvarchar(max);default:null"`
	DescriptionAr *string `json:"description_ar" gorm:"type:nvarchar(max);default:null"`
	Status        string  `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active';check:status IN ('active','inactive')"`
}

func (CommitteeClassification) TableName() string {
	return "committee_classifications"
}
