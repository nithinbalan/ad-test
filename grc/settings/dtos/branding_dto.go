package dtos

type BrandingListDTO struct {
	ID                          uint    `json:"id"`
	LoginPageTitleEn            string  `json:"login_page_title_en"`
	LoginPageTitleAr            string  `json:"login_page_title_ar"`
	LoginPageSubtitleEn         string  `json:"login_page_subtitle_en"`
	LoginPageSubtitleAr         string  `json:"login_page_subtitle_ar"`
	LoginPageBgImagePath        *string `json:"login_page_bg_image_path"`
	OrgLogoImagePath            *string `json:"org_logo_image_path"`
	FaviconImagePath            *string `json:"favicon_image_path"`
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
type BrandingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Filters     map[string][]string
	QuickFilter []string
}
