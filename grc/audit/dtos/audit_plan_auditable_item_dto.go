package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type AuditPlanAuditableItemListDTO struct {
	ID                          uint   `json:"id"`
	AuditPlanID                 uint   `json:"audit_plan_id"`
	AuditableItemID             uint   `json:"auditable_item_id"`
	CreatedByID                 uint   `json:"created_by_id"`
	CreatedAt                   string `json:"created_at"`
	ReferenceNumber             string `json:"reference_number"`
	AuditPlanTitle              string `json:"audit_plan_title"`
	AuditableItemTitleEn        string `json:"auditable_item_title_en"`
	AuditableItemTitleAr        string `json:"auditable_item_title_ar"`
	AuditCategoryTitleEn        string `json:"audit_category_title_en"`
	AuditCategoryTitleAr        string `json:"audit_category_title_ar"`
	AuditTeamTitleEn            string `json:"audit_team_title_en"`
	AuditTeamTitleAr            string `json:"audit_team_title_ar"`
	TeamLeadFirstNameEn         string `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr         string `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn          string `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr          string `json:"team_lead_last_name_ar"`
	TeamLeadEmail               string `json:"team_lead_email"`
	TeamLeadProfileImagePath    string `json:"team_lead_profile_image_path"`
	AuditStatusTitleEn          string `json:"audit_status_title_en"`
	AuditStatusTitleAr          string `json:"audit_status_title_ar"`
	AuditStatusType             string `json:"audit_status_type"`
	AuditStatusColorCode        string `json:"audit_status_color_code"`
	CreatedUserID               uint   `json:"created_user_id"`
	CreatedUserFirstNameEn      string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string `json:"created_user_last_name_ar"`
	CreatedUserInitial          string `json:"created_user_initial"`
	CreatedUserColor            string `json:"created_user_color"`
	CreatedUserEmail            string `json:"created_user_email"`
	CreatedUserProfileImagePath string `json:"created_user_profile_image_path"`
}

type AuditPlanAuditableItemFormDTO struct {
	AuditableItemIds []*uint `json:"auditable_item_ids" binding:"required"`
}

type AuditPlanAuditableItemListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditPlanAuditableItemListResponseDTO struct {
	Items      []*AuditPlanAuditableItemListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}
