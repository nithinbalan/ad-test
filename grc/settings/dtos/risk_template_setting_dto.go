package dtos

import "time"

type RiskTemplateSettingListDTO struct {
	ID                          uint    `json:"id"`
	LogoImagePath               *string `json:"logo_image_path"`
	CoverPageBgImagePath        *string `json:"cover_page_bg_image_path"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string  `json:"updated_user_initial"`
	UpdatedUserColor            string  `json:"updated_user_color"`
	UpdatedUserEmail            string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath string  `json:"updated_user_profile_image_path"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
}

type RiskTemplateSettingDetailDTO struct {
	ID                   uint       `json:"id"`
	LogoImagePath        *string    `json:"logo_image_path"`
	CoverPageBgImagePath *string    `json:"cover_page_bg_image_path"`
	CreatedByID          uint       `json:"created_by_id"`
	CreatedByName        string     `json:"created_by_name"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedByID          *uint      `json:"updated_by_id"`
	UpdatedByName        *string    `json:"updated_by_name"`
	UpdatedAt            *time.Time `json:"updated_at"`
}

type RiskTemplateSettingFormDTO struct {
	Logo       *RiskTemplateSettingDocumentFormDTO `json:"logo,omitempty"`
	CoverImage *RiskTemplateSettingDocumentFormDTO `json:"cover_image,omitempty"`
}

type RiskTemplateSettingDocumentFormDTO struct {
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type RiskTemplateSettingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	Page        int
	PageSize    int
	QuickFilter []string
}
