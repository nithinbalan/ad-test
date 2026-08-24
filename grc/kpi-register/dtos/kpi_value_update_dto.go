package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	masterDTOs "jethur/grc/masters/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type KPIValueUpdateListDTO struct {
	ID                          uint     `json:"id"`
	KPIRegisterTitle            string   `json:"kpi_register_title"`
	KPIRegisterID               uint     `json:"kpi_register_id"`
	ReviewedDate                string   `json:"reviewed_date"`
	KPICategoryTitleEn          string   `json:"kpi_category_title_en"`
	KPICategoryTitleAr          string   `json:"kpi_category_title_ar"`
	Value                       *float64 `json:"value"`
	StartDate                   *string  `json:"start_date"`
	EndDate                     *string  `json:"end_date"`
	PeriodLabel                 *string  `json:"period_label"`
	Remarks                     string   `json:"remarks"`
	ReviewUserID                *uint    `json:"review_user_id"`
	ReviewUserFirstNameEn       *string  `json:"review_user_first_name_en"`
	ReviewUserLastNameEn        *string  `json:"review_user_last_name_en"`
	ReviewUserFirstNameAr       *string  `json:"review_user_first_name_ar"`
	ReviewUserLastNameAr        *string  `json:"review_user_last_name_ar"`
	ReviewUserInitial           *string  `json:"review_user_initial"`
	ReviewUserEmail             *string  `json:"review_user_email"`
	ReviewUserColor             *string  `json:"review_user_color"`
	ReviewUserProfileImagePath  *string  `json:"review_user_profile_image_path"`
	KPIReviewStatusID           uint     `json:"kpi_review_status_id"`
	KPIReviewStatusTitleEn      string   `json:"kpi_review_status_title_en"`
	KPIReviewStatusTitleAr      string   `json:"kpi_review_status_title_ar"`
	KPIReviewStatusColorCode    string   `json:"kpi_review_status_color_code"`
	KPIReviewStatusType         string   `json:"kpi_review_status_type"`
	IsLatest                    bool     `json:"is_latest"`
	CreatedAt                   string   `json:"created_at"`
	UpdatedAt                   string   `json:"updated_at"`
	CreatedUserID               uint     `json:"created_user_id"`
	CreatedUserFirstNameEn      string   `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string   `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string   `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       *string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          *string  `json:"created_user_initial"`
	CreatedUserColor            *string  `json:"created_user_color"`
	CreatedUserEmail            *string  `json:"created_user_email"`
	CreatedUserProfileImagePath *string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint    `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string  `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string  `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string  `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string  `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string  `json:"updated_user_initial"`
	UpdatedUserColor            *string  `json:"updated_user_color"`
	UpdatedUserEmail            *string  `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string  `json:"updated_user_profile_image_path"`
}

type KPIValueUpdateFormDTO struct {
	Value            float64                          `json:"value"`
	Remarks          *string                          `json:"remarks"`
	Documents        []*KPIValueUpdateDocumentFormDTO `json:"documents"`
	DataInputs       []*ValueUpdateInputDataFormDTO   `json:"data_inputs"`
	CorrectiveAction *string                          `json:"corrective_action"`
}

type ValueUpdateInputDataFormDTO struct {
	Variable string  `json:"variable"`
	Input    string  `json:"input"`
	Value    float64 `json:"value"`
}

type KPIValueUpdateListQueryParamsDTO struct {
	Sort          string
	SortBy        string
	Keyword       string
	Page          int
	PageSize      int
	Filters       map[string][]string
	QuickFilter   []string
	Export        bool
	IsPassiveUser bool
}

type KPIValueUpdateListResponseDTO struct {
	Items      []*KPIValueUpdateListDTO     `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type KPIValueUpdateDetailsDTO struct {
	ID                   uint                              `json:"id"`
	KPIRegisterID        uint                              `json:"kpi_register_id"`
	StartDateString      *string                           `json:"start_date"`   // ✅ new
	EndDateString        *string                           `json:"end_date"`     // ✅ new
	PeriodLabel          *string                           `json:"period_label"` // ✅ new
	ReviewedDateString   *string                           `json:"reviewed_date"`
	Value                *float64                          `json:"value"`
	Remarks              string                            `json:"remarks"`
	KPIRegister          KPIRegisterDetailDTO              `json:"kpi_register"`
	Unit                 *masterDTOs.UnitDTO               `json:"unit"`
	ReviewUser           *organizationDTOs.UserMinimalDTO  `json:"review_user"`
	KPIReviewStatus      *KPIReviewStatusDTO               `json:"kpi_review_status"`
	IsLatest             bool                              `json:"is_latest"`
	CreatedAt            time.Time                         `json:"created_at"`
	UpdatedAt            *time.Time                        `json:"updated_at"`
	ApprovedAt           *time.Time                        `json:"approved_at"`
	RevertedAt           *time.Time                        `json:"reverted_at"`
	Documents            []*KPIValueUpdateDocumentDTO      `json:"documents"`
	DataInputs           []*KPIValueUpdateDataInputDTO     `json:"data_inputs"`
	ValueUpdateHistories []*KPIValueUpdateHistoryDetailDTO `json:"value_update_histories"`
	ValueUpdateReviews   []*KPIValueUpdateReviewDetailDTO  `json:"value_update_reviews"`
	CreatedByUser        *organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser        *organizationDTOs.UserMinimalDTO  `json:"updated_by"`
}

type KPIValueUpdateDataInputDTO struct {
	ID       uint    `json:"id"`
	Variable string  `json:"variable"`
	Input    string  `json:"input"`
	Value    float64 `json:"value"`
}

type KPIValueUpdateParsedDates struct {
	ReviewedDate *time.Time
}

type KPIValueUpdateDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Token        string  `json:"token"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}
type KPIValueUpdateDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"` // Optional field for indicating if the file is new or updated
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type KPIValueUpdateHistoryInputDetailDTO struct {
	ID            uint                             `json:"id"`
	Variable      string                           `json:"variable"`
	Input         string                           `json:"input"`
	Value         float64                          `json:"value"`
	CreatedByUser *organizationDTOs.UserMinimalDTO `json:"created_by"`
}

type KPIValueUpdateHistoryDetailDTO struct {
	ID                          uint                                  `json:"id"`
	KPIValueUpdateID            uint                                  `json:"kpi_value_update_id"`
	Value                       *float64                              `json:"value"`
	Unit                        *masterDTOs.UnitDTO                   `json:"unit"`
	Remarks                     *string                               `json:"remarks"`
	CreatedAt                   time.Time                             `json:"created_at"`
	CreatedByUser               *organizationDTOs.UserMinimalDTO      `json:"created_by"`
	KPIValueUpdateHistoryInputs []KPIValueUpdateHistoryInputDetailDTO `json:"inputs"`
	CorrectiveAction            *string                               `json:"corrective_action"`
}

type KPIValueUpdateReviewDetailDTO struct {
	ID                uint                             `json:"id"`
	KPIValueUpdateID  uint                             `json:"kpi_value_update_id"`
	KPIReviewStatusID uint                             `json:"kpi_review_status_id"`
	KPIReviewStatus   *KPIReviewStatusDTO              `json:"kpi_review_status"`
	Comment           *string                          `json:"comment"`
	CreatedAt         time.Time                        `json:"created_at"`
	CreatedByUser     *organizationDTOs.UserMinimalDTO `json:"created_by"`
}

type KPIValueUpdateEvaluationFormDTO struct {
	DataInputs []KPIValueUpdateDataInputDTO `json:"data_inputs"`
}

type KPIValueUpdateEvaluationDetailDTO struct {
	DataInputs []KPIValueUpdateDataInputDTO `json:"data_inputs"`
	Value      float64                      `json:"value"`
}

// KPIValueUpdateReviewFormDTO is used for submit action (comment is optional)
type KPIValueUpdateReviewFormDTO struct {
	Comment *string `json:"comment"`
}

// KPIValueUpdateCommentFormDTO is used for approve and revert actions (comment is required)
type KPIValueUpdateCommentFormDTO struct {
	Comment string `json:"comment" binding:"required"`
}
