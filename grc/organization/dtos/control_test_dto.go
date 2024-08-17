package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ControlTestListDTO struct {
	ID                                   uint    `json:"id"`
	ControlID                            uint    `json:"control_id"`
	Title                                string  `json:"title"`
	ControlImplementationStatusTitleEn   string  `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string  `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string  `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       string  `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       string  `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     string  `json:"control_efficiency_status_color_code"`
	Description                          *string `json:"description"`
	Score                                *uint   `json:"score"`
	Comments                             *string `json:"comments"`
	TestedSource                         string  `json:"tested_source"`
	TestedByUserID                       uint    `json:"tested_by_user_id"`
	TestedByUserFirstNameEn              string  `json:"tested_by_user_first_name_en"`
	TestedByUserFirstNameAr              string  `json:"tested_by_user_first_name_ar"`
	TestedByUserLastNameEn               string  `json:"tested_by_user_last_name_en"`
	TestedByUserLastNameAr               string  `json:"tested_by_user_last_name_ar"`
	TestedByUserInitial                  string  `json:"tested_by_user_initial"`
	TestedByUserColor                    string  `json:"tested_by_user_color"`
	TestedByUserEmail                    string  `json:"tested_by_user_email"`
	TestedByUserProfileImagePath         string  `json:"tested_by_user_profile_image_path"`
	UpdatedUserID                        *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   *string `json:"updated_user_initial"`
	UpdatedUserColor                     *string `json:"updated_user_color"`
	UpdatedUserEmail                     *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath          *string `json:"updated_user_profile_image_path"`

	TestedAt  string     `json:"tested_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ControlTestListResponseDTO struct {
	Items      []*ControlTestListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ControlTestDetailsDTO struct {
	ID             uint       `json:"id"`
	ControlID      uint       `json:"control_id"`
	Title          string     `json:"title"`
	Description    *string    `json:"description"`
	Score          *uint      `json:"score"`
	Comments       *string    `json:"comments"`
	TestedSource   string     `json:"tested_source"`
	TestedAtString *string    `json:"tested_at"`
	UpdatedAt      *time.Time `json:"updated_at"`

	ControlImplementationStatus ControlImplementationStatusDTO `json:"control_implementation_status"`
	ControlEfficiencyStatus     ControlEfficiencyStatusDTO     `json:"control_efficiency_status"`
	TestedByUser                UserMinimalDTO                 `json:"tested_by_user"`
	UpdatedByUser               *UserMinimalDTO                `json:"updated_by_user"`
	Documents                   []*ControlTestDocumentDTO      `json:"documents"`
}

type ControlTestListQueryParamsDTO struct {
	Sort         string
	SortBy       string
	Keyword      string
	Status       string
	Page         int
	PageSize     int
	Filters      map[string][]string
	QuickFilter  []string
	TestedAtFrom string
	TestedAtTo   string
}

type ControlTestDocumentDTO struct {
	ID           *uint   `json:"id,omitempty"`
	URL          *string `json:"url"`
	Name         string  `json:"file_name"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type ControlTestFormDTO struct {
	ControlImplementationStatusID uint                          `json:"control_implementation_status_id" binding:"required"`
	ControlEfficiencyStatusID     uint                          `json:"control_efficiency_status_id" binding:"required"`
	Title                         string                        `json:"title" binding:"required"`
	Description                   *string                       `json:"description"  binding:"omitempty"`
	Score                         *uint                         `json:"score" binding:"omitempty"`
	Comments                      *string                       `json:"comments" binding:"omitempty"`
	TestedAt                      *string                       `json:"tested_at"` // assuming this will be a date in string format
	TestedBy                      uint                          `json:"tested_by" binding:"required"`
	Documents                     []*commonDTOs.DocumentFormDTO `json:"documents" binding:"omitempty"`
}

type ControlTestParsedDates struct {
	TestedAt *time.Time
}

type ControlTrendAnalysisDTO struct {
	TestedAtString string  `json:"tested_at"`
	Score          float64 `json:"score"`
	ColorCode      string  `json:"color_code"`
}
