package dtos

type BcpTeamsListDTO struct {
	ID                              uint    `json:"id"`
	BusinessContinuityPlanID        uint    `json:"business_continuity_plan_id"`
	ResponsibleUserID               uint    `json:"responsible_user_id"`
	ResponsibleUnitID               uint    `json:"responsible_unit_id"`
	ResponsibleUnitTitleEn          string  `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string  `json:"responsible_unit_title_ar"`
	Responsibility                  string  `json:"responsibility"`
	BackupUserID                    *uint   `json:"backup_user_id"`
	CreatedByID                     uint    `json:"created_by_id"`
	ResponsibleUserFirstNameEn      string  `json:"responsible_user_first_name_en"`
	ResponsibleUserLastNameEn       string  `json:"responsible_user_last_name_en"`
	ResponsibleUserFirstNameAr      string  `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameAr       string  `json:"responsible_user_last_name_ar"`
	ResponsibleUserEmail            string  `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string  `json:"responsible_user_profile_image_path"`
	BackupUserFirstNameEn           string  `json:"backup_user_first_name_en"`
	BackupUserLastNameEn            string  `json:"backup_user_last_name_en"`
	BackupUserFirstNameAr           string  `json:"backup_user_first_name_ar"`
	BackupUserLastNameAr            string  `json:"backup_user_last_name_ar"`
	BackupUserEmail                 string  `json:"backup_user_email"`
	BackupUserProfileImagePath      string  `json:"backup_user_profile_image_path"`
	CreatedAt                       string  `json:"created_at"`
	UpdatedAt                       *string `json:"updated_at"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
}

type BcpTeamsFormDTO struct {
	ResponsibleUserID uint   `json:"responsible_user_id" binding:"required"`
	Responsibility    string `json:"responsibility" binding:"required,max=5000"`
	BackupUserID      *uint  `json:"backup_user_id"`
}

type BcpTeamsListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	Export      bool
}
