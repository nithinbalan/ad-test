package models

type KRISubCategory struct {
	ID            uint    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	KRICategoryID uint    `json:"kri_category_id" binding:"required" gorm:"not null"`
	TitleEn       string  `json:"title_en" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	TitleAr       string  `json:"title_ar" binding:"required,max=100" gorm:"type:nvarchar(100);not null"`
	DescriptionEn *string `json:"description_en" gorm:"type:nvarchar(MAX);default:null"`
	DescriptionAr *string `json:"description_ar" gorm:"type:nvarchar(MAX);default:null"`
	Status        string  `json:"status" binding:"required,max=10,oneof=active inactive" gorm:"type:varchar(10);not null;default:active"`

	KRICategory *KRICategory `json:"kri_category" gorm:"foreignKey:KRICategoryID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
}

func (KRISubCategory) TableName() string {
	return "kri_sub_categories"
}
