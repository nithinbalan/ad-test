package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
)

type AuditReportReviewListDTO struct {
	ID                      uint   `json:"id"`
	AuditDraftReportID      uint   `json:"audit_draft_report_id"`
	AuditReportWorkflowID   uint   `json:"audit_report_workflow_id"`
	WorkflowStatusID        uint   `json:"workflow_status_id"`
	WorkflowStatusTitleEn   string `json:"workflow_status_title_en"`
	WorkflowStatusTitleAr   string `json:"workflow_status_title_ar"`
	WorkflowStatusColorCode string `json:"workflow_status_color_code"`
	WorkflowStatusType      string `json:"workflow_status_type"`
	UserID                  uint   `json:"user_id"`
	UserFirstNameEn         string `json:"user_first_name_en"`
	UserFirstNameAr         string `json:"user_first_name_ar"`
	UserLastNameEn          string `json:"user_last_name_en"`
	UserLastNameAr          string `json:"user_last_name_ar"`
	UserInitial             string `json:"user_initial"`
	UserColor               string `json:"user_color"`
	UserEmail               string `json:"user_email"`
	UserProfileImage        string `json:"user_profile_image_path"`
	Comment                 string `json:"comment"`
	CreatedByID             uint   `json:"created_by_id"`
	CreatedUserFirstNameEn  string `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string `json:"created_user_last_name_ar"`
	CreatedUserInitial      string `json:"created_user_initial"`
	CreatedUserColor        string `json:"created_user_color"`
	CreatedUserEmail        string `json:"created_user_email"`
	CreatedUserProfileImage string `json:"created_user_profile_image_path"`
}

type AuditReportReviewListResponseDTO struct {
	Items      []*AuditReportReviewListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditReportReviewListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}

type AuditReportWorkflowSubmitFormDTO struct {
	UserIDs []*uint `json:"user_ids" binding:"required"`
	Comment *string `json:"comment"`
}

type AuditReportWorkflowApproveFormDTO struct {
	Comment *string `json:"comment"`
}