package dtos

import "time"

type BrandingDetailDTO struct {
	ID                   uint            `json:"id"`
	LoginPageTitleEn     string          `json:"login_page_title_en"`
	LoginPageTitleAr     string          `json:"login_page_title_ar"`
	LoginPageSubtitleEn  string          `json:"login_page_subtitle_en"`
	LoginPageSubtitleAr  string          `json:"login_page_subtitle_ar"`
	LoginPageBgImagePath *string         `json:"login_page_bg_image_path"`
	OrgLogoImagePath     *string         `json:"org_logo_image_path"`
	FaviconImagePath     *string         `json:"favicon_image_path"`
	LoaderGifPath        *string         `json:"loader_gif_path"`
	CreatedByUser        *UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser        *UserMinimalDTO `json:"updated_by_user"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            *time.Time      `json:"updated_at"`
}

type BrandingUpdateFormDTO struct {
	LoginPageTitleEn     string `json:"login_page_title_en" binding:"required"`
	LoginPageTitleAr     string `json:"login_page_title_ar" binding:"required"`
	LoginPageSubtitleEn  string `json:"login_page_subtitle_en" binding:"required"`
	LoginPageSubtitleAr  string `json:"login_page_subtitle_ar" binding:"required"`
	LoginPageBgImagePath string `json:"login_page_bg_image_path" binding:"required"`
	OrgLogoImagePath     string `json:"org_logo_image_path" binding:"required"`
	FaviconImagePath     string `json:"favicon_image_path" binding:"required"`
	LoaderGifPath        string `json:"loader_gif_path" binding:"required"`
}
