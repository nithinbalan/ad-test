package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AwarenessPostCampaignSurveyListDTO struct {
	ID                          uint       `json:"id"`
	ReferenceNumber             string     `json:"reference_number"`
	LinkedCampaignID            *uint      `json:"linked_campaign_id"`
	LinkedCampaignTitle         string     `json:"linked_campaign_title"`
	SendDate                    *time.Time `json:"send_date"`
	TargetAudience              *string    `json:"target_audience"`
	Recipients                  *int       `json:"recipients"`
	ResponseRate                *int       `json:"response_rate"`
	QuestionCount               *int       `json:"question_count"`
	AverageRating               *int       `json:"average_rating"`
	MaturityLevelID             *uint      `json:"maturity_level_id"`
	MaturityLevelTitleEn        string     `json:"maturity_level_title_en"`
	MaturityLevelTitleAr        string     `json:"maturity_level_title_ar"`
	Recommendations             *string    `json:"recommendations"`
	ExecutingParty              *string    `json:"executing_party"`
	ResultsLink                 *string    `json:"results_link"`
	CreatedAt                   string     `json:"created_at"`
	UpdatedAt                   string     `json:"updated_at"`
	CreatedUserID               uint       `json:"created_user_id"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn      string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
}

type AwarenessPostCampaignSurveyListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
	Export   bool
	DateFrom string
	DateTo   string
}

type AwarenessPostCampaignSurveyListResponseDTO struct {
	Items      []*AwarenessPostCampaignSurveyListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO          `json:"pagination"`
}

type AwarenessPostCampaignSurveyFormDTO struct {
	LinkedCampaignID            *uint      `json:"linked_campaign_id"`
	SendDate                    *time.Time `json:"send_date"`
	TargetAudience              *string    `json:"target_audience"`
	Recipients                  *int       `json:"recipients"`
	ResponseRate                *int       `json:"response_rate"`
	QuestionCount               *int       `json:"question_count"`
	AverageRating               *int       `json:"average_rating"`
	MaturityLevelID             *uint      `json:"maturity_level_id"`
	Recommendations             *string    `json:"recommendations"`
	ExecutingParty              *string    `json:"executing_party"`
	ResultsLink                 *string    `json:"results_link"`
	AwarenessAssessmentThemeIDs []uint     `json:"awareness_assessment_theme_ids"`
}

type AwarenessPostCampaignSurveyDetailResponseDTO struct {
	ID              uint                                `json:"id"`
	ReferenceNumber string                              `json:"reference_number"`
	LinkedCampaign  *AwarenessCampaignDetailResponseDTO `json:"linked_campaign"`
	SendDate        *time.Time                          `json:"send_date"`
	TargetAudience  *string                             `json:"target_audience"`
	Recipients      *int                                `json:"recipients"`
	ResponseRate    *int                                `json:"response_rate"`
	QuestionCount   *int                                `json:"question_count"`
	AverageRating   *int                                `json:"average_rating"`
	MaturityLevel   *AwarenessMaturityLevelDTO          `json:"maturity_level"`
	Recommendations *string                             `json:"recommendations"`
	ExecutingParty  *string                             `json:"executing_party"`
	ResultsLink     *string                             `json:"results_link"`
	SurveyThemes    []*AwarenessAssessmentThemeDTO      `json:"survey_themes"`
	CreatedAt       time.Time                           `json:"created_at"`
	UpdatedAt       *time.Time                          `json:"updated_at"`
	CreatedByUser   organizationDTOs.UserMinimalDTO     `json:"created_by"`
	UpdatedByUser   organizationDTOs.UserMinimalDTO     `json:"updated_by"`
}
