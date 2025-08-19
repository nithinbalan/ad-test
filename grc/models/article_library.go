package models

import (
	"time"
)

type ArticleLibrary struct {
	ID                uint      `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	ArticleCategoryID uint      `gorm:"not null" json:"article_category_id" binding:"required"`
	ImagePath         string    `gorm:"type:nvarchar(500);not null" json:"image_path"`
	TitleEn           string    `gorm:"type:nvarchar(max);not null" json:"title_en" binding:"required"`
	TitleAr           string    `gorm:"type:nvarchar(max);not null" json:"title_ar" binding:"required"`
	TableOfContents   string    `gorm:"type:nvarchar(max)" json:"table_of_contents"`
	Content           string    `gorm:"type:nvarchar(max);not null" json:"content"`
	CreatedBy         string    `gorm:"type:nvarchar(255);not null" json:"created_by"`
	CreatedAt         time.Time `gorm:"type:datetime;not null" json:"created_at"`

	// Relationship
	ArticleCategory *ArticleCategory `gorm:"foreignKey:ArticleCategoryID;constraint:OnUpdate:CASCADE,OnDelete:NO ACTION;" json:"article_category"`
}

func (ArticleLibrary) TableName() string {
	return "article_library"
}
