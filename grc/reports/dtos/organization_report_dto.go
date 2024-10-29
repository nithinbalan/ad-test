package dtos

type UsersByStatusDTO struct {
	UserStatusID        uint   `json:"user_status_id"`
	UserStatusTitleEn   string `json:"user_status_title_en"`
	UserStatusTitleAr   string `json:"user_status_title_ar"`
	UserStatusColorCode string `json:"user_status_color_code"`
	UserStatusType      string `json:"user_status_type"`
	Count               int    `json:"count"`
}

type UsersByRoleDTO struct {
	RoleID      uint   `json:"role_id"`
	RoleTitleEn string `json:"role_title_en"`
	RoleTitleAr string `json:"role_title_ar"`
	Count       int    `json:"count"`
}

type UsersByDesignationDTO struct {
	DesignationID      uint   `json:"designation_id"`
	DesignationTitleEn string `json:"designation_title_en"`
	DesignationTitleAr string `json:"designation_title_ar"`
	Count              int    `json:"count"`
}

type ProcessByUnitDTO struct {
	OwnerUnitID        uint   `json:"owner_unit_id"`
	OwnerUnitTitleEn   string `json:"owner_unit_title_en"`
	OwnerUnitTitleAr   string `json:"owner_unit_title_ar"`
	OwnerUnitColorCode string `json:"owner_unit_color_code"`
	OwnerUnitStatus    string `json:"owner_unit_status"`
	Count              int    `json:"count"`
}

type ProcessByRiskRatingDTO struct {
	RiskRatingID        uint   `json:"risk_rating_id"`
	RiskRatingTitleEn   string `json:"risk_rating_title_en"`
	RiskRatingTitleAr   string `json:"risk_rating_title_ar"`
	RiskRatingColorCode string `json:"risk_rating_color_code"`
	RiskRatingType      string `json:"risk_rating_type"`
	Count               int    `json:"count"`
}

type ProcessByResponsibleUserDTO struct {
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	Count                           int    `json:"count"`
}

type ProcessByProcessGroupDTO struct {
	ProcessGroupID      uint   `json:"process_group_id"`
	ProcessGroupTitleEn string `json:"process_group_title_en"`
	ProcessGroupTitleAr string `json:"process_group_title_ar"`
	Count               int    `json:"count"`
}
