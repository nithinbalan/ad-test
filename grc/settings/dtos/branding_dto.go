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
type BrandingDocumentFormDTO struct {
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type BrandingFormDTO struct {
	LoginPageTitleEn    string                   `json:"login_page_title_en" binding:"required"`
	LoginPageTitleAr    string                   `json:"login_page_title_ar" binding:"required"`
	LoginPageSubtitleEn string                   `json:"login_page_subtitle_en" binding:"required"`
	LoginPageSubtitleAr string                   `json:"login_page_subtitle_ar" binding:"required"`
	LoginPageBgImage    *BrandingDocumentFormDTO `json:"login_page_bg_image,omitempty"`
	OrgLogoImage        *BrandingDocumentFormDTO `json:"org_logo_image,omitempty"`
	FaviconImage        *BrandingDocumentFormDTO `json:"favicon_image,omitempty"`
	LoaderGif           *BrandingDocumentFormDTO `json:"loader_gif,omitempty"`
}

type LoginBrandingDTO struct {
	LoginPageTitleEn     string  `json:"login_page_title_en"`
	LoginPageTitleAr     string  `json:"login_page_title_ar"`
	LoginPageSubtitleEn  string  `json:"login_page_subtitle_en"`
	LoginPageSubtitleAr  string  `json:"login_page_subtitle_ar"`
	LoginPageBgImagePath *string `json:"login_page_bg_image_path"`
	OrgLogoImagePath     *string `json:"org_logo_image_path"`
}
