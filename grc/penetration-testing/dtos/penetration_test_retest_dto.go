package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
)

type PenetrationTestRetestListDTO struct {
	ID                          uint    `json:"id"`
	ReferenceNumber             string  `json:"reference_number"`
	PenetrationTestID           uint    `json:"penetration_test_id"`
	StartDate                   string  `json:"start_date"`
	EndDate                     *string `json:"end_date"`
	PenetrationTestExecutorID   uint    `json:"penetration_test_executor_id"`
	ExecutorTitleEn             string  `json:"executor_title_en"`
	ExecutorTitleAr             string  `json:"executor_title_ar"`
	PenetrationTestRetestStatusID uint  `json:"penetration_test_retest_status_id"`
	StatusTitleEn               string  `json:"status_title_en"`
	StatusTitleAr               string  `json:"status_title_ar"`
	StatusColorCode             string  `json:"status_color_code"`
	StatusType                  string  `json:"status_type"`
	Notes                       *string `json:"notes"`
	RemainingCritical           int     `json:"remaining_critical"`
	RemainingHigh               int     `json:"remaining_high"`
	RemainingMedium             int     `json:"remaining_medium"`
	RemainingLow                int     `json:"remaining_low"`
	RemainingTotal              int     `json:"remaining_total"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedByID                 uint    `json:"created_by_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
}

type PenetrationTestRetestFormDTO struct {
	PenetrationTestID         uint    `json:"penetration_test_id" binding:"required"`
	StartDate                 string  `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate                   *string `json:"end_date,omitempty" binding:"omitempty,datetime=2006-01-02"`
	PenetrationTestExecutorID uint    `json:"penetration_test_executor_id" binding:"required"`
	Notes                     *string `json:"notes,omitempty"`
	FindingIDs                []uint  `json:"finding_ids" binding:"omitempty"`
}

type PenetrationTestRetestFindingUpdateFormDTO struct {
	FindingIDs []uint `json:"finding_ids" binding:"required,min=1"`
}

type PenetrationTestRetestFindingSaveFormDTO struct {
	Comment   string                        `json:"comment" binding:"required"`
	Documents []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type PenetrationTestRetestFindingResultFormDTO struct {
	ResultID uint `json:"result_id" binding:"required"`
}

type PenetrationTestRetestResultDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
	Type    string `json:"type"`
}

type PenetrationTestRetestStatusDTO struct {
	ID        uint   `json:"id"`
	TitleEn   string `json:"title_en"`
	TitleAr   string `json:"title_ar"`
	ColorCode string `json:"color_code"`
	Type      string `json:"type"`
}

type PenetrationTestRetestFindingDTO struct {
	ID                       uint                                           `json:"id"`
	PenetrationTestFindingID uint                                           `json:"penetration_test_finding_id"`
	FindingReferenceNumber   string                                         `json:"finding_reference_number"`
	FindingTitle             string                                         `json:"finding_title"`
	Comment                  *string                                        `json:"comment"`
	CommentAt                *time.Time                                     `json:"comment_at"`
	CommentByUser            *organizationDTOs.UserMinimalDTO               `json:"comment_by_user"`
	Result                   *PenetrationTestRetestResultDTO                `json:"result"`
	ResultAt                 *time.Time                                     `json:"result_at"`
	ResultByUser             *organizationDTOs.UserMinimalDTO               `json:"result_by_user"`
	Documents                []*models.PenetrationTestRetestFindingDocument `json:"documents"`
}

type PenetrationTestRetestDetailDTO struct {
	ID                        uint                               `json:"id"`
	ReferenceNumber           string                             `json:"reference_number"`
	PenetrationTestID         uint                               `json:"penetration_test_id"`
	StartDate                 time.Time                          `json:"start_date"`
	EndDate                   *time.Time                         `json:"end_date"`
	PenetrationTestExecutorID uint                               `json:"penetration_test_executor_id"`
	PenetrationTestExecutor   PenetrationTestExecutorDTO         `json:"penetration_test_executor"`
	PenetrationTestRetestStatus *PenetrationTestRetestStatusDTO  `json:"penetration_test_retest_status"`
	Notes                     *string                            `json:"notes"`
	Findings                  []*PenetrationTestRetestFindingDTO `json:"findings"`
	RemainingCritical         int                                `json:"remaining_critical"`
	RemainingHigh             int                                `json:"remaining_high"`
	RemainingMedium           int                                `json:"remaining_medium"`
	RemainingLow              int                                `json:"remaining_low"`
	RemainingTotal            int                                `json:"remaining_total"`
	CreatedAt                 time.Time                          `json:"created_at"`
	UpdatedAt                 *time.Time                         `json:"updated_at"`
	CreatedByUser             organizationDTOs.UserMinimalDTO    `json:"created_by_user"`
	UpdatedByUser             *organizationDTOs.UserMinimalDTO   `json:"updated_by_user"`
}

type PenetrationTestRetestListQueryParamsDTO struct {
	Sort              string
	SortBy            string
	Keyword           string
	Status            string
	Page              int
	PageSize          int
	PenetrationTestID uint
	Filters           map[string][]string
	QuickFilter       []string
	DateFrom          string
	DateTo            string
}

type PenetrationTestRetestListResponseDTO struct {
	Items      []*PenetrationTestRetestListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}
