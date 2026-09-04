package dtos

import (
	"time"
)

type ProjectChangeRequestReviewCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type ProjectChangeRequestReviewCommentListDTO struct {
	ID                          uint                                            `json:"id"`
	Comment                     string                                          `json:"comment"`
	IsEdited                    bool                                            `json:"is_edited"`
	CommentStatusID             uint                                            `json:"comment_status_id"`
	CommentStatusTitleEn        string                                          `json:"comment_status_title_en"`
	CommentStatusTitleAr        string                                          `json:"comment_status_title_ar"`
	CommentStatusType           string                                          `json:"comment_status_type"`
	CommentStatusColorCode      string                                          `json:"comment_status_color_code"`
	ProjectChangeRequestID      uint                                            `json:"project_change_request_id"`
	Level                       *int                                            `json:"level"`
	CreatedUserID               uint                                            `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                          `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                          `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                          `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                          `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                          `json:"created_user_initial"`
	CreatedUserColor            string                                          `json:"created_user_color"`
	CreatedUserEmail            string                                          `json:"created_user_email"`
	CreatedUserProfileImagePath *string                                         `json:"created_user_profile_image_path"`
	CreatedAt                   time.Time                                       `json:"created_at"`
	ParentCommentID             *uint                                           `json:"parent_comment_id"`
	Documents                   []*ProjectChangeRequestReviewCommentDocumentDTO `json:"documents"`
	Children                    []*ProjectChangeRequestReviewCommentListDTO     `json:"children"`
}

type ProjectChangeRequestReviewCommentListResponseDTO struct {
	Comments      []*ProjectChangeRequestReviewCommentListDTO `json:"comments"`
	OpenCount     int                                         `json:"open_count"`
	ResolvedCount int                                         `json:"resolved_count"`
}

type ProjectChangeRequestReviewCommentDocumentFormDTO struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	Type         string `json:"type"`
}

type ProjectChangeRequestReviewCommentFormDTO struct {
	Comment         string                                              `json:"comment"`
	ParentCommentID *uint                                               `json:"parent_comment_id"`
	Documents       []*ProjectChangeRequestReviewCommentDocumentFormDTO `json:"documents"`
}

type ProjectChangeRequestReviewCommentDetailDTO struct {
	ID                     uint                                            `json:"id"`
	Comment                string                                          `json:"comment"`
	IsEdited               bool                                            `json:"is_edited"`
	CommentStatusID        uint                                            `json:"comment_status_id"`
	ProjectChangeRequestID uint                                            `json:"project_change_request_id"`
	Level                  *int                                            `json:"level"`
	ParentCommentID        *uint                                           `json:"parent_comment_id"`
	CreatedBYID            uint                                            `json:"created_by_id"`
	UpdatedBYID            *uint                                           `json:"updated_by_id"`
	CreatedAt              time.Time                                       `json:"created_at"`
	UpdatedAt              *time.Time                                      `json:"updated_at"`
	Documents              []*ProjectChangeRequestReviewCommentDocumentDTO `json:"documents"`
}

type ProjectChangeRequestReviewCommentListQueryParamsDTO struct {
	Page        int
	PageSize    int
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}
