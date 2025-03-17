package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

// AssetPolicyCommentListDTO represents a summarized comment for listing
type AssetPolicyCommentListDTO struct {
	ID                          uint       `json:"id"`
	Comments                    string     `json:"comments"`
	AssetPolicyID               uint       `json:"asset_policy_id"`
	CreatedAt                   time.Time  `json:"created_at"`
	CreatedByID                 uint       `json:"created_by_id"`
	UpdatedByID                 uint       `json:"updated_by_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn      *string    `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string    `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string    `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string    `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
	UpdatedAt                   *time.Time `json:"updated_at"`
}

// AssetPolicyCommentDetailDTO provides full details for a comment
type AssetPolicyCommentDetailDTO struct {
	ID            uint                             `json:"id"`
	Comments      string                           `json:"comments"`
	AssetPolicyID uint                             `json:"asset_policy_id"`
	AssetPolicy   AssetPolicyDetailDTO             `json:"asset_policy"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

// AssetPolicyCommentFormDTO is used for creating/updating a comment
type AssetPolicyCommentFormDTO struct {
	Comments string `json:"comments" binding:"required"`
}
