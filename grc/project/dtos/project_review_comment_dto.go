package dtos

import (
	"time"
)

type ProjectReviewCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int64  `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type ProjectReviewCommentListDTO struct {
	ID                          uint                               `json:"id"`
	Comment                     string                             `json:"comment"`
	IsEdited                    bool                               `json:"is_edited"`
	CommentStatusID             uint                               `json:"comment_status_id"`
	CommentStatusTitleEn        string                             `json:"comment_status_title_en"`
	CommentStatusTitleAr        string                             `json:"comment_status_title_ar"`
	CommentStatusType           string                             `json:"comment_status_type"`
	CommentStatusColorCode      string                             `json:"comment_status_color_code"`
	ProjectID                   uint                               `json:"project_id"`
	Level                       *int                               `json:"level"`
	CreatedUserID               uint                               `json:"created_user_id"`
	CreatedUserFirstNameEn      string                             `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                             `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                             `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                             `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                             `json:"created_user_initial"`
	CreatedUserColor            string                             `json:"created_user_color"`
	CreatedUserEmail            string                             `json:"created_user_email"`
	CreatedUserProfileImagePath *string                            `json:"created_user_profile_image_path"`
	CreatedAt                   time.Time                          `json:"created_at"`
	ParentCommentID             *uint                              `json:"parent_comment_id"`
	Documents                   []*ProjectReviewCommentDocumentDTO `json:"documents" gorm:"-"`
	Children                    []*ProjectReviewCommentListDTO     `json:"children" gorm:"-"`
}

type ProjectReviewCommentListResponseDTO struct {
	Comments      []*ProjectReviewCommentListDTO `json:"comments"`
	OpenCount     int                            `json:"open_count"`
	ResolvedCount int                            `json:"resolved_count"`
}

type ProjectReviewCommentDocumentFormDTO struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	Type         string `json:"type"`
}

type ProjectReviewCommentFormDTO struct {
	Comment         string                                 `json:"comment"`
	ParentCommentID *uint                                  `json:"parent_comment_id"`
	Documents       []*ProjectReviewCommentDocumentFormDTO `json:"documents"`
}

type ProjectReviewCommentDetailDTO struct {
	ID              uint                               `json:"id"`
	Comment         string                             `json:"comment"`
	IsEdited        bool                               `json:"is_edited"`
	CommentStatusID uint                               `json:"comment_status_id"`
	ProjectID       uint                               `json:"project_id"`
	Level           *int                               `json:"level"`
	ParentCommentID *uint                              `json:"parent_comment_id"`
	CreatedBYID     uint                               `json:"created_by_id"`
	UpdatedBYID     *uint                              `json:"updated_by_id"`
	CreatedAt       time.Time                          `json:"created_at"`
	UpdatedAt       *time.Time                         `json:"updated_at"`
	Documents       []*ProjectReviewCommentDocumentDTO `json:"documents"`
}

type ProjectReviewCommentListQueryParamsDTO struct {
	Page        int
	PageSize    int
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}
