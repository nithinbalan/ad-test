package dtos

import (
	masterDTOs "jethur/grc/masters/dtos"
	"time"
)

type InternalCommitteeReviewCommentListDTO struct {
	ID                          uint                                     `json:"id"`
	Comment                     string                                   `json:"comment"`
	Level                       *int                                     `json:"level"`
	CommentStatusID             uint                                     `json:"comment_status_id"`
	CommentStatusTitleEn        string                                   `json:"comment_status_title_en"`
	CommentStatusTitleAr        string                                   `json:"comment_status_title_ar"`
	CommentStatusColorCode      string                                   `json:"comment_status_color_code"`
	InternalCommitteeVersionID  uint                                     `json:"internal_committee_version_id"`
	IsEdited                    bool                                     `json:"is_edited"`
	CreatedAt                   string                                   `json:"created_at"`
	CreatedUserID               uint                                     `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                   `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                   `json:"created_user_initial"`
	CreatedUserColor            string                                   `json:"created_user_color"`
	CreatedUserEmail            string                                   `json:"created_user_email"`
	CreatedUserProfileImagePath string                                   `json:"created_user_profile_image_path"`
	ParentCommentID             *uint                                    `json:"parent_comment_id"`
	Children                    []*InternalCommitteeReviewCommentListDTO `json:"children" gorm:"-"`
}

type InternalCommitteeReviewCommentListResponseDTO struct {
	Comments      []*InternalCommitteeReviewCommentListDTO `json:"comments"`
	OpenCount     int                                      `json:"open_count"`
	ResolvedCount int                                      `json:"resolved_count"`
}

type InternalCommitteeReviewCommentFormDTO struct {
	Comment         string `json:"comment" binding:"required"`
	ParentCommentID *uint  `json:"parent_comment_id"`
}

type InternalCommitteeReviewCommentDetailsDTO struct {
	ID                         uint                          `json:"id"`
	Comment                    string                        `json:"comment"`
	Level                      *int                          `json:"level"`
	CommentStatusID            uint                          `json:"comment_status_id"`
	CommentStatus              *masterDTOs.WorkflowStatusDTO `json:"comment_status"`
	InternalCommitteeVersionID uint                          `json:"internal_committee_version_id"`
	ParentCommentID            *uint                         `json:"parent_comment_id"`
	IsEdited                   bool                          `json:"is_edited"`
	CreatedBYID                uint                          `json:"created_by_id"`
	CreatedByUser              *UserMinimalDTO               `json:"created_by_user"`
	UpdatedBYID                *uint                         `json:"updated_by_id"`
	UpdatedByUser              *UserMinimalDTO               `json:"updated_by_user"`
	CreatedAt                  time.Time                     `json:"created_at"`
	UpdatedAt                  *time.Time                    `json:"updated_at"`
}

type InternalCommitteeReviewCommentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}

type StatusCount struct {
	Status string
	Count  int
}
