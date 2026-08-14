package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
)

type AuditFindingReassessmentDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type AuditFindingReassessmentListDTO struct {
	ID                          uint                                   `json:"id"`
	AuditFindingID              uint                                   `json:"audit_finding_id"`
	Outcome                     *string                                `json:"outcome"`
	Comment                     *string                                `json:"comment"`
	CreatedAt                   string                                 `json:"created_at"`
	UpdatedAt                   *string                                `json:"updated_at"`
	CreatedUserID               uint                                   `json:"created_user_id"`
	CreatedUserFirstNameEn      string                                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                                 `json:"created_user_last_name_ar"`
	CreatedUserInitial          string                                 `json:"created_user_initial"`
	CreatedUserColor            string                                 `json:"created_user_color"`
	CreatedUserEmail            string                                 `json:"created_user_email"`
	CreatedUserProfileImagePath string                                 `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint                                  `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string                                `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string                                `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string                                `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string                                `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string                                `json:"updated_user_initial"`
	UpdatedUserColor            *string                                `json:"updated_user_color"`
	UpdatedUserEmail            *string                                `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string                                `json:"updated_user_profile_image_path"`
	IsLatest                    bool                                   `json:"is_latest" gorm:"-"`
	Documents                   []*AuditFindingReassessmentDocumentDTO `json:"documents" gorm:"-"`
}

type AuditFindingReassessmentDetailDTO struct {
	ID             uint                                   `json:"id"`
	AuditFindingID uint                                   `json:"audit_finding_id"`
	Outcome        *string                                `json:"outcome"`
	Comment        *string                                `json:"comment"`
	CreatedAt      time.Time                              `json:"created_at"`
	UpdatedAt      *time.Time                             `json:"updated_at"`
	CreatedByUser  organizationDTOs.UserMinimalDTO        `json:"created_by"`
	UpdatedByUser  *organizationDTOs.UserMinimalDTO       `json:"updated_by"`
	IsLatest       bool                                   `json:"is_latest"`
	Documents      []*AuditFindingReassessmentDocumentDTO `json:"documents"`
}

type AuditFindingReassessmentListResponseDTO struct {
	Items      []*AuditFindingReassessmentListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO       `json:"pagination"`
}

type AuditFindingReassessmentFormDTO struct {
	Outcome   string                        `json:"outcome" binding:"required,oneof=pass fail"`
	Comment   string                        `json:"comment" binding:"required"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type AuditFindingReassessmentListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
}
