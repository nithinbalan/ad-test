package models

type ArticleCategory struct {
	ID      uint   `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	TitleEn string `gorm:"type:nvarchar(255);not null" json:"title_en" binding:"required"`
	TitleAr string `gorm:"type:nvarchar(255);not null" json:"title_ar" binding:"required"`
}

func (ArticleCategory) TableName() string {
	return "article_categories"
}
