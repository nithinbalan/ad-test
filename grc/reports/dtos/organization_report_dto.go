package dtos

type UsersByStatusDTO struct {
	ID                  uint   `json:"id"`
	UserStatusTitleEn   string `json:"user_status_title_en"`
	UserStatusTitleAr   string `json:"user_status_title_ar"`
	UserStatusColorCode string `json:"user_status_color_code"`
	UserStatusType      string `json:"user_status_type"`
	Count               int    `json:"count"`
}

type UsersByRoleDTO struct {
	ID          uint   `json:"id"`
	RoleTitleEn string `json:"role_title_en"`
	RoleTitleAr string `json:"role_title_ar"`
	Count       int    `json:"count"`
}

type UsersByDesignationDTO struct {
	ID                 uint   `json:"id"`
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
