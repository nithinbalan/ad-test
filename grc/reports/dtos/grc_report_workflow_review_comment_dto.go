package dtos

import "time"

type GRCReportWorkflowReviewCommentListDTO struct {
	ID                              uint    `json:"id"`
	Comments                        string  `json:"comment"`
	GRCReportWorkflowID             uint    `json:"grc_report_workflow_id"`
	CreatedAt                       string  `json:"created_at"`
	CreatedByID                     uint    `json:"created_by_id"`
	CreatedUserFirstNameEn          string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr          string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn           string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr           string  `json:"created_user_last_name_ar"`
	CreatedUserInitial              string  `json:"created_user_initial"`
	CreatedUserColor                string  `json:"created_user_color"`
	CreatedUserEmail                string  `json:"created_user_email"`
	CreatedUserProfileImage         string  `json:"created_user_profile_image"`
	GRCReportCommentStatusTitleEn   *string `json:"grc_report_comment_status_title_en"`
	GRCReportCommentStatusType      *string `json:"grc_report_comment_status_type"`
	GRCReportCommentStatusTitleAr   *string `json:"grc_report_comment_status_title_ar"`
	GRCReportCommentStatusColorCode *string `json:"grc_report_comment_status_color_code"`
}

type GRCReportWorkflowReviewCommentDetailDTO struct {
	ID                      uint                       `json:"id"`
	Comment                 string                     `json:"comment"`
	GRCReportWorkflowID     uint                       `json:"grc_report_workflow_id"`
	GRCReportCommentStatus  *GRCReportCommentStatusDTO `json:"grc_report_comment_status"`
	CreatedAt               time.Time                  `json:"created_at"`
	CreatedByID             uint                       `json:"created_by_id"`
	CreatedUserFirstNameEn  string                     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr  string                     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn   string                     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr   string                     `json:"created_user_last_name_ar"`
	CreatedUserInitial      string                     `json:"created_user_initial"`
	CreatedUserColor        string                     `json:"created_user_color"`
	CreatedUserEmail        string                     `json:"created_user_email"`
	CreatedUserProfileImage string                     `json:"created_user_profile_image"`
}

type GRCReportWorkflowReviewCommentFormDTO struct {
	Type    string `json:"type" binding:"required"`
	Comment string `json:"comment" binding:"required"`
}
