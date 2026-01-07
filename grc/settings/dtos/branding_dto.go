package dtos

import "time"

type BrandingDetailDTO struct {
	ID                   uint            `json:"id"`
	LoginPageTitle       string          `json:"login_page_title"`
	LoginPageSubtitle    string          `json:"login_page_subtitle"`
	LoginPageBgImagePath *string         `json:"login_page_bg_image_path"`
	OrgLogoImagePath     *string         `json:"org_logo_image_path"`
	FaviconImagePath     *string         `json:"favicon_image_path"`
	CreatedByUser        *UserMinimalDTO `json:"created_by_user"`
	UpdatedByUser        *UserMinimalDTO `json:"updated_by_user"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            *time.Time      `json:"updated_at"`
}

type BrandingUpdateFormDTO struct {
	LoginPageTitle       string `json:"login_page_title" binding:"required"`
	LoginPageSubtitle    string `json:"login_page_subtitle" binding:"required"`
	LoginPageBgImagePath string `json:"login_page_bg_image_path" binding:"required"`
	OrgLogoImagePath     string `json:"org_logo_image_path" binding:"required"`
	FaviconImagePath     string `json:"favicon_image_path" binding:"required"`
}
