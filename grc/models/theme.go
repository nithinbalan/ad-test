package models

import "time"

type Theme struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	TitleEn string `json:"title_en" gorm:"type:nvarchar(100);not null" binding:"required"`
	TitleAr string `json:"title_ar" gorm:"type:nvarchar(100);not null" binding:"required"`

	DescriptionEn *string `json:"description_en,omitempty" gorm:"type:nvarchar(max)"`
	DescriptionAr *string `json:"description_ar,omitempty" gorm:"type:nvarchar(max)"`

	PrimaryColor   string  `json:"primary_color" gorm:"type:varchar(20);not null"`
	SecondaryColor *string `json:"secondary_color" gorm:"type:varchar(20)"`
	LinkColor      *string `json:"link_color" gorm:"type:varchar(20)"`
	SidebarBgColor *string `json:"sidebar_bg_color" gorm:"type:varchar(20)"`

	IsSystem bool `json:"is_system" gorm:"not null;default:false"`

	CreatedBYID uint  `json:"created_by_id" gorm:"column:created_by_id;not null"`
	UpdatedBYID *uint `json:"updated_by_id" gorm:"column:updated_by_id;default:null"`

	CreatedAt time.Time  `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:datetime;default:null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID"`
}

func (Theme) TableName() string {
	return "themes"
}