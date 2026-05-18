package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	"time"
)

type ControlTestListDTO struct {
	ID                                   uint                      `json:"id"`
	ControlID                            uint                      `json:"control_id"`
	Title                                string                    `json:"title"`
	ControlImplementationStatusTitleEn   string                    `json:"control_implementation_status_title_en"`
	ControlImplementationStatusTitleAr   string                    `json:"control_implementation_status_title_ar"`
	ControlImplementationStatusColorCode string                    `json:"control_implementation_status_color_code"`
	ControlEfficiencyStatusTitleEn       *string                   `json:"control_efficiency_status_title_en"`
	ControlEfficiencyStatusTitleAr       *string                   `json:"control_efficiency_status_title_ar"`
	ControlEfficiencyStatusColorCode     *string                   `json:"control_efficiency_status_color_code"`
	Description                          *string                   `json:"description"`
	Score                                *uint                     `json:"score"`
	ControlEffectivenessScore            *float64                  `json:"control_effectiveness_score"`
	ControlEffectivenessRatingID         *uint                     `json:"control_effectiveness_rating_id"`
	ControlEffectivenessRatingTitleEn    *string                   `json:"control_effectiveness_rating_title_en"`
	ControlEffectivenessRatingTitleAr    *string                   `json:"control_effectiveness_rating_title_ar"`
	ControlEffectivenessRatingMinScore   *float64                  `json:"control_effectiveness_rating_min_score"`
	ControlEffectivenessRatingMaxScore   *float64                  `json:"control_effectiveness_rating_max_score"`
	ControlEffectivenessRatingType       *string                   `json:"control_effectiveness_rating_type"`
	ControlEffectivenessRatingColorCode  *string                   `json:"control_effectiveness_rating_color_code"`
	RiskReduction                        *float64                  `json:"risk_reduction"`
	Comments                             *string                   `json:"comments"`
	ControlAutomationID                  *uint                     `json:"control_automation_id"`
	ControlAutomationTitleEn             *string                   `json:"control_automation_title_en"`
	ControlAutomationTitleAr             *string                   `json:"control_automation_title_ar"`
	ControlAutomationType                *string                   `json:"control_automation_type"`
	ControlAutomationScore               *float64                  `json:"control_automation_score"`
	ControlAutomationStatus              *string                   `json:"control_automation_status"`
	ControlConsistencyID                 *uint                     `json:"control_consistency_id"`
	ControlConsistencyTitleEn            *string                   `json:"control_consistency_title_en"`
	ControlConsistencyTitleAr            *string                   `json:"control_consistency_title_ar"`
	ControlConsistencyType               *string                   `json:"control_consistency_type"`
	ControlConsistencyScore              *float64                  `json:"control_consistency_score"`
	ControlConsistencyStatus             *string                   `json:"control_consistency_status"`
	ControlConsistencyDescriptionEn      *string                   `json:"control_consistency_description_en"`
	TestedSource                         string                    `json:"tested_source"`
	TestedByUserID                       uint                      `json:"tested_by_user_id"`
	TestedByUserFirstNameEn              string                    `json:"tested_by_user_first_name_en"`
	TestedByUserFirstNameAr              string                    `json:"tested_by_user_first_name_ar"`
	TestedByUserLastNameEn               string                    `json:"tested_by_user_last_name_en"`
	TestedByUserLastNameAr               string                    `json:"tested_by_user_last_name_ar"`
	TestedByUserInitial                  string                    `json:"tested_by_user_initial"`
	TestedByUserColor                    string                    `json:"tested_by_user_color"`
	TestedByUserEmail                    string                    `json:"tested_by_user_email"`
	TestedByUserProfileImagePath         string                    `json:"tested_by_user_profile_image_path"`
	UpdatedUserID                        *uint                     `json:"updated_user_id"`
	UpdatedUserFirstNameEn               *string                   `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr               *string                   `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn                *string                   `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr                *string                   `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                   *string                   `json:"updated_user_initial"`
	UpdatedUserColor                     *string                   `json:"updated_user_color"`
	UpdatedUserEmail                     *string                   `json:"updated_user_email"`
	UpdatedUserProfileImagePath          *string                   `json:"updated_user_profile_image_path"`
	Documents                            []*ControlTestDocumentDTO `json:"documents" gorm:"-"`

	TestedAt  string     `json:"tested_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ControlTestListResponseDTO struct {
	Items      []*ControlTestListDTO        `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type ControlTestDetailsDTO struct {
	ID                         uint                           `json:"id"`
	ControlID                  uint                           `json:"control_id"`
	Title                      string                         `json:"title"`
	Description                *string                        `json:"description"`
	Score                      *uint                          `json:"score"`
	ControlEffectivenessScore  *float64                       `json:"control_effectiveness_score"`
	ControlEffectivenessRating *ControlEffectivenessRatingDTO `json:"control_effectiveness_rating"`
	RiskReduction              *float64                       `json:"risk_reduction"`
	Comments                   *string                        `json:"comments"`
	ControlAutomation          *ControlAutomationDetailDTO    `json:"control_automation"`
	ControlConsistency         *ControlConsistencyDetailDTO   `json:"control_consistency"`
	TestedSource               string                         `json:"tested_source"`
	TestedAtString             *string                        `json:"tested_at"`
	UpdatedAt                  *time.Time                     `json:"updated_at"`

	ControlImplementationStatus ControlImplementationStatusDTO `json:"control_implementation_status"`
	ControlEfficiencyStatus     *ControlEfficiencyStatusDTO    `json:"control_efficiency_status"`
	TestedByUser                UserMinimalDTO                 `json:"tested_by_user"`
	UpdatedByUser               *UserMinimalDTO                `json:"updated_by_user"`
	Documents                   []ControlTestDocumentDTO       `json:"documents" binding:"required"`
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
	ID            uint       `json:"id"`
	ControlTestID uint       `json:"control_test_id"`
	URL           string     `json:"url"`
	Name          string     `json:"name"`
	Size          int        `json:"size"`
	Token         string     `json:"token"`
	OriginalName  string     `json:"original_name"`
	Type          string     `json:"type"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}

type ControlTestFormDTO struct {
	ControlImplementationStatusID uint                          `json:"control_implementation_status_id" binding:"required"`
	ControlEfficiencyStatusID     *uint                         `json:"control_efficiency_status_id,omitempty" binding:"omitempty"`
	Title                         string                        `json:"title" binding:"required"`
	Description                   *string                       `json:"description"  binding:"omitempty"`
	Score                         *uint                         `json:"score" binding:"omitempty"`
	Comments                      *string                       `json:"comments" binding:"omitempty"`
	TestedAt                      *string                       `json:"tested_at" binding:"required"` // assuming this will be a date in string format
	TestedBy                      uint                          `json:"tested_by" binding:"required"`
	ControlAutomationID           uint                          `json:"control_automation_id,omitempty" binding:"required"`
	ControlConsistencyID          uint                          `json:"control_consistency_id,omitempty" binding:"required"`
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
