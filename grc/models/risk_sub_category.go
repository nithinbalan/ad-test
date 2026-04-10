package models

type RiskSubCategory struct {
	ID             uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskCategoryID uint    `json:"risk_category_id" binding:"required" gorm:"not null"`
	TitleEn        string  `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr        string  `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn  *string `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr  *string `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	Status         string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:active"`

	RiskCategory *RiskCategory `json:"risk_category" gorm:"foreignKey:RiskCategoryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

func (RiskSubCategory) TableName() string {
	return "risk_sub_categories"
}
