package dtos

import (
	"time"
)

type ThemeListDTO struct {
	ID             int        `json:"id"`
	TitleEn        string     `json:"title_en"`
	TitleAr        string     `json:"title_ar"`
	PrimaryColor   string     `json:"primary_color"`
	SecondaryColor *string    `json:"secondary_color"`
	LinkColor      *string    `json:"link_color"`
	SidebarBgColor *string    `json:"sidebar_bg_color"`
	IsSystem       bool       `json:"is_system"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`

	CreatedUserID           uint   `json:"created_user_id"`
	CreatedUserFirstNameEn  string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string `json:"created_user_last_name_ar"`
	CreatedUserInitial      string `json:"created_user_initial"`
	CreatedUserColor        string `json:"created_user_color"`
	CreatedUserEmail        string `json:"created_user_email"`
	CreatedUserProfileImage string `json:"created_user_profile_image_path"`

	UpdatedUserID           *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn  *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr  *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn   *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr   *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial      *string `json:"updated_user_initial"`
	UpdatedUserColor        *string `json:"updated_user_color"`
	UpdatedUserEmail        *string `json:"updated_user_email"`
	UpdatedUserProfileImage *string `json:"updated_user_profile_image_path"`
}

type ThemeFormDTO struct {
	TitleEn        string  `json:"title_en" binding:"required"`
	TitleAr        string  `json:"title_ar" binding:"required"`
	DescriptionEn  *string `json:"description_en"`
	DescriptionAr  *string `json:"description_ar"`
	PrimaryColor   string  `json:"primary_color" binding:"required"`
	SecondaryColor string  `json:"secondary_color" binding:"required"`
	LinkColor      string  `json:"link_color" binding:"required"`
	SidebarBgColor string  `json:"sidebar_bg_color" binding:"required"`
	IsSystem       bool    `json:"is_system"`
}

type ThemeDetailDTO struct {
	ID             uint       `json:"id"`
	TitleEn        string     `json:"title_en"`
	TitleAr        string     `json:"title_ar"`
	DescriptionEn  *string    `json:"description_en"`
	DescriptionAr  *string    `json:"description_ar"`
	PrimaryColor   string     `json:"primary_color"`
	SecondaryColor *string    `json:"secondary_color"`
	LinkColor      *string    `json:"link_color"`
	SidebarBgColor *string    `json:"sidebar_bg_color"`
	IsSystem       bool       `json:"is_system"`
	CreatedByID    uint       `json:"created_by_id"`
	UpdatedByID    *uint      `json:"updated_by_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`

	CreatedByUser UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *UserMinimalDTO `json:"updated_by_user"`
}
