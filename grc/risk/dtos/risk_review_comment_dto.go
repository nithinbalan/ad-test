package dtos

import (
	"time"

	organizationDTOs "jethur/grc/organization/dtos"
)

type RiskReviewCommentListDTO struct {
	ID                          uint                        `json:"id"`
	Comment                     string                      `json:"comment" binding:"required"`
	CommentStatusTitleEn        string                      `json:"comment_status_title_en"`
	CommentStatusTitleAr        string                      `json:"comment_status_title_ar"`
	CommentStatusType           string                      `json:"comment_status_type"`
	CommentStatusColorCode      string                      `json:"comment_status_color_code"`
	RiskRegisterID              uint                        `json:"risk_register_id"`
	Level                       int                         `json:"level"`
	CreatedAt                   string                      `json:"created_at"`
	CreatedBYID                 uint                        `json:"created_by_id"`
	CreatedUserFirstNameEn      string                      `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                      `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                      `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                      `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                      `json:"created_user_initial"`
	CreatedUserColor            string                      `json:"created_user_color"`
	CreatedUserEmail            string                      `json:"created_user_email"`
	CreatedUserProfileImagePath string                      `json:"created_user_profile_image_path"`
	RiskReviewCommentID         *uint                       `json:"risk_review_comment_id"`
	Children                    []*RiskReviewCommentListDTO `json:"children" gorm:"-"`
}

type RiskReviewCommentListResponseDTO struct {
	Comments      []*RiskReviewCommentListDTO `json:"comments"`
	OpenCount     int                         `json:"open_count"`
	ResolvedCount int                         `json:"resolved_count"`
}

type StatusCount struct {
	Status string
	Count  int
}

type RiskReviewCommentFormDTO struct {
	Comment         string `json:"comment" binding:"required,max=1000"`
	ParentCommentID *uint  `json:"parent_comment_id"`
}

type RiskReviewCommentDetailDTO struct {
	ID             uint                            `json:"id"`
	Comment        string                          `json:"comment"`
	RiskRegisterID uint                            `json:"risk_register_id"`
	CreatedAt      time.Time                       `json:"created_at"`
	CreatedByUser  organizationDTOs.UserMinimalDTO `json:"created_by"`
}

type RiskReviewCommentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}
