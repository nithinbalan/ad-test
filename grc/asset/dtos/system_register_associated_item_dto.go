package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type SystemRegisterAssociatedItemListQueryParamsDTO struct {
	Sort        string   `form:"sort"`
	SortBy      string   `form:"sort_by"`
	Keyword     string   `form:"keyword"`
	Status      string   `form:"status"`
	Page        int      `form:"page"`
	PageSize    int      `form:"page_size"`
	QuickFilter []string `form:"quick_filter"`
}

type SystemRegisterSLAListResponseDTO struct {
	Items      []*SystemRegisterSLAListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type SystemRegisterSLAListDTO struct {
	ID                              uint   `json:"id"`
	SystemRegisterID                uint   `json:"system_register_id"`
	SLARegisterID                   uint   `json:"sla_register_id"`
	ReferenceNumber                 string `json:"reference_number"`
	LifecycleStatusID               uint   `json:"lifecycle_status_id"`
	LifecycleStatusTitleEn          string `json:"lifecycle_status_title_en"`
	LifecycleStatusTitleAr          string `json:"lifecycle_status_title_ar"`
	LifecycleStatusColorCode        string `json:"lifecycle_status_color_code"`
	TitleEn                         string `json:"title_en"`
	TitleAr                         string `json:"title_ar"`
	Description                     string `json:"description"`
	Client                          string `json:"client"`
	SLACategoryTitleEn              string `json:"sla_category_title_en"`
	SLACategoryTitleAr              string `json:"sla_category_title_ar"`
	SeverityLevelTitleEn            string `json:"severity_level_title_en"`
	SeverityLevelTitleAr            string `json:"severity_level_title_ar"`
	ResponsibleUnitTitleEn          string `json:"responsible_unit_title_en"`
	ResponsibleUnitTitleAr          string `json:"responsible_unit_title_ar"`
	ResponsibleUserID               uint   `json:"responsible_user_id"`
	ResponsibleUserFirstNameEn      string `json:"responsible_user_first_name_en"`
	ResponsibleUserFirstNameAr      string `json:"responsible_user_first_name_ar"`
	ResponsibleUserLastNameEn       string `json:"responsible_user_last_name_en"`
	ResponsibleUserLastNameAr       string `json:"responsible_user_last_name_ar"`
	ResponsibleUserInitial          string `json:"responsible_user_initial"`
	ResponsibleUserColor            string `json:"responsible_user_color"`
	ResponsibleUserEmail            string `json:"responsible_user_email"`
	ResponsibleUserProfileImagePath string `json:"responsible_user_profile_image_path"`
	ComplianceStatusTitleEn         string `json:"compliance_status_title_en"`
	ComplianceStatusTitleAr         string `json:"compliance_status_title_ar"`
	ComplianceStatusColorCode       string `json:"compliance_status_color_code"`
}

type SystemRegisterSLAFormDTO struct {
	SLARegisterIDs []uint `json:"sla_register_ids" binding:"required"`
}

type SystemRegisterProcessFormDTO struct {
	ProcessIDs []uint `json:"process_ids" binding:"required"`
}
