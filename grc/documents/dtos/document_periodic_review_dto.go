package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type DocumentPeriodicReviewListDTO struct {
	ID               uint      `json:"id"`
	DocumentID       uint      `json:"document_id"`
	DocumentTitle    string    `json:"document_title"`
	ReviewDate       time.Time `json:"review_date"`
	IsChangeRequired bool      `json:"is_change_required"`
	Description      *string   `json:"description"`

	CreatedByUserID               uint   `json:"created_by_user_id"`
	CreatedByUserFirstNameEn      string `json:"created_by_user_first_name_en"`
	CreatedByUserFirstNameAr      string `json:"created_by_user_first_name_ar"`
	CreatedByUserLastNameEn       string `json:"created_by_user_last_name_en"`
	CreatedByUserLastNameAr       string `json:"created_by_user_last_name_ar"`
	CreatedByUserEmail            string `json:"created_by_user_email"`
	CreatedByUserInitial          string `json:"created_by_user_initial"`
	CreatedByUserColor            string `json:"created_by_user_color"`
	CreatedByUserProfileImagePath string `json:"created_by_user_profile_image_path"`

	UpdatedByUserID               *uint   `json:"updated_by_user_id"`
	UpdatedByUserFirstNameEn      *string `json:"updated_by_user_first_name_en"`
	UpdatedByUserFirstNameAr      *string `json:"updated_by_user_first_name_ar"`
	UpdatedByUserLastNameEn       *string `json:"updated_by_user_last_name_en"`
	UpdatedByUserLastNameAr       *string `json:"updated_by_user_last_name_ar"`
	UpdatedByUserEmail            *string `json:"updated_by_user_email"`
	UpdatedByUserInitial          *string `json:"updated_by_user_initial"`
	UpdatedByUserColor            *string `json:"updated_by_user_color"`
	UpdatedByUserProfileImagePath *string `json:"updated_by_user_profile_image_path"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DocumentPeriodicReviewDetailsDTO struct {
	ID               uint      `json:"id"`
	DocumentID       uint      `json:"document_id"`
	ReviewDate       time.Time `json:"review_date"`
	IsChangeRequired bool      `json:"is_change_required"`
	Description      *string   `json:"description"`

	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`

	SupportDocuments      []*DocumentPeriodicReviewSupportDocumentDTO `json:"support_documents"`
	DocumentChangeRequest *DocumentChangeRequestDetailsDTO            `json:"document_change_request"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type DocumentPeriodicReviewFormDTO struct {
	ReviewDate       string                        `json:"review_date" binding:"required"`
	ISChangeRequired bool                          `json:"is_change_required"`
	Description      *string                       `json:"description" binding:"omitempty"`
	SupportDocuments []*commonDTOs.DocumentFormDTO `json:"support_documents" binding:"omitempty"`

	DocumentChangeRequestTypeID uint                       `json:"document_change_request_type_id"`
	ReasonForChange             *string                    `json:"reason_for_change"`
	Document                    commonDTOs.DocumentFormDTO `json:"document"`
}

type DocumentPeriodicReviewSupportDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
	Token        string  `json:"token"`
}

type DocumentPeriodicReviewListResponseDTO struct {
	Items      []*DocumentPeriodicReviewListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO     `json:"pagination"`
}

type DocumentPeriodicReviewQueryParamsDTO struct {
	Sort        string              `json:"sort"`
	SortBy      string              `json:"sort_by"`
	Keyword     string              `json:"keyword"`
	Status      string              `json:"status"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	Filters     map[string][]string `json:"filters"`
	QuickFilter []string            `json:"quick_filter"`
}

type DocumentPeriodicReviewParsedDates struct {
	ReviewDate *time.Time
}
