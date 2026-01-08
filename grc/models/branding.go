package models

import "time"

type Branding struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	LoginPageTitleEn     string     `json:"login_page_title_en" gorm:"type:nvarchar(max);not null"`
	LoginPageTitleAr     string     `json:"login_page_title_ar" gorm:"type:nvarchar(max);not null"`
	LoginPageSubtitleEn  string     `json:"login_page_subtitle_en" gorm:"type:nvarchar(max);not null"`
	LoginPageSubtitleAr  string     `json:"login_page_subtitle_ar" gorm:"type:nvarchar(max);not null"`
	LoginPageBgImagePath *string    `json:"login_page_bg_image_path" gorm:"type:varchar(500);default:null"`
	OrgLogoImagePath     *string    `json:"org_logo_image_path" gorm:"type:varchar(500);default:null"`
	FaviconImagePath     *string    `json:"favicon_image_path" gorm:"type:varchar(500);default:null"`
	LoaderGifPath        *string    `json:"loader_gif_path" gorm:"type:varchar(500);default:null"`
	CreatedBYID          uint       `json:"created_by_id" gorm:"not null"`
	UpdatedBYID          *uint      `json:"updated_by_id" gorm:"default:null"`
	CreatedAt            time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
	UpdatedAt            *time.Time `json:"updated_at" gorm:"null"`

	// Relationships
	CreatedByUser *User `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name
func (Branding) TableName() string {
	return "branding"
}
