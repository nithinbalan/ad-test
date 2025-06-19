package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type StrategicObjectivePeriodicReviewListDTO struct {
	ID                        uint      `json:"id"`
	StrategicObjectiveID      uint      `json:"strategic_objective_id"`
	StrategicObjectiveTitleEn string    `json:"strategic_objective_title_en"`
	StrategicObjectiveTitleAr string    `json:"strategic_objective_title_ar"`
	ReviewDate                time.Time `json:"review_date"`
	Description               *string   `json:"description"`

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

	SupportDocuments []*StrategicObjectivePeriodicReviewSupportDocumentDTO `json:"documents" gorm:"-"`
}

type StrategicObjectivePeriodicReviewDetailsDTO struct {
	ID                   uint    `json:"id"`
	StrategicObjectiveID uint    `json:"strategic_objective_id"`
	ReviewDateString     string  `json:"review_date"`
	Description          *string `json:"description"`

	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by_user"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by_user"`

	SupportDocuments []*StrategicObjectivePeriodicReviewSupportDocumentDTO `json:"documents"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type StrategicObjectivePeriodicReviewFormDTO struct {
	ReviewDate       string                        `json:"review_date" binding:"required"`
	Description      *string                       `json:"description" binding:"omitempty"`
	SupportDocuments []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`

	Document commonDTOs.DocumentFormDTO `json:"document"`
}

type StrategicObjectivePeriodicReviewSupportDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
	Token        string  `json:"token"`
}

type StrategicObjectivePeriodicReviewListResponseDTO struct {
	Items      []*StrategicObjectivePeriodicReviewListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO               `json:"pagination"`
}

type StrategicObjectivePeriodicReviewQueryParamsDTO struct {
	Sort        string              `json:"sort"`
	SortBy      string              `json:"sort_by"`
	Keyword     string              `json:"keyword"`
	Status      string              `json:"status"`
	Page        int                 `json:"page"`
	PageSize    int                 `json:"page_size"`
	Filters     map[string][]string `json:"filters"`
	QuickFilter []string            `json:"quick_filter"`
}

type StrategicObjectivePeriodicReviewParsedDates struct {
	ReviewDate *time.Time
}
