package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type AssetRiskResponseReviewCommentListDTO struct {
	ID                          uint                                         `json:"id"`
	Comment                     string                                       `json:"comment"`
	Level                       *int                                         `json:"level"`
	CommentStatusID             uint                                         `json:"comment_status_id"`
	CommentStatusTitleEn        string                                       `json:"comment_status_title_en"`
	CommentStatusTitleAr        string                                       `json:"comment_status_title_ar"`
	CommentStatusColorCode      string                                       `json:"comment_status_color_code"`
	CommentStatusType           string                                       `json:"comment_status_type"`
	RiskResponseID              uint                                         `json:"risk_response_id"`
	IsEdited                    bool                                         `json:"is_edited"`
	CreatedAt                   string                                       `json:"created_at"`
	CreatedUserID               uint                                         `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                       `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                       `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                       `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                       `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                       `json:"created_user_initial"`
	CreatedUserColor            string                                       `json:"created_user_color"`
	CreatedUserEmail            string                                       `json:"created_user_email"`
	CreatedUserProfileImagePath string                                       `json:"created_user_profile_image_path"`
	ParentCommentID             *uint                                        `json:"parent_comment_id"`
	Children                    []*AssetRiskResponseReviewCommentListDTO     `json:"children" gorm:"-"`
	Documents                   []*AssetRiskResponseReviewCommentDocumentDTO `json:"documents" gorm:"-"`
}

type AssetRiskResponseReviewCommentListResponseDTO struct {
	Comments      []*AssetRiskResponseReviewCommentListDTO `json:"comments"`
	OpenCount     int                                      `json:"open_count"`
	ResolvedCount int                                      `json:"resolved_count"`
}

type AssetRiskResponseReviewCommentFormDTO struct {
	Comment         string                        `json:"comment" binding:"required,max=1000"`
	ParentCommentID *uint                         `json:"parent_comment_id"`
	Documents       []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AssetRiskResponseReviewCommentDetailDTO struct {
	ID              uint                                         `json:"id"`
	Comment         string                                       `json:"comment"`
	Level           *int                                         `json:"level"`
	CommentStatusID uint                                         `json:"comment_status_id"`
	RiskResponseID  uint                                         `json:"risk_response_id"`
	ParentCommentID *uint                                        `json:"parent_comment_id"`
	IsEdited        bool                                         `json:"is_edited"`
	CreatedBYID     uint                                         `json:"created_by_id"`
	CreatedByUser   *organizationDTOs.UserMinimalDTO             `json:"created_by_user"`
	UpdatedBYID     *uint                                        `json:"updated_by_id"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO             `json:"updated_by_user"`
	Documents       []*AssetRiskResponseReviewCommentDocumentDTO `json:"documents"`
	CreatedAt       time.Time                                    `json:"created_at"`
	UpdatedAt       *time.Time                                   `json:"updated_at"`
}

type AssetRiskResponseReviewCommentListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
}

type AssetRiskResponseReviewCommentDocumentDTO struct {
	ID           uint   `json:"id"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	Size         int    `json:"size"`
	Token        string `json:"token"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}
