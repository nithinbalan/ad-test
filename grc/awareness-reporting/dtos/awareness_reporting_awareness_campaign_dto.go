package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
	"time"
)

type AwarenessCampaignListDTO struct {
	ID                               uint       `json:"id"`
	ReferenceNumber                  string     `json:"reference_number"`
	CampaignTitle                    string     `json:"campaign_title"`
	Objective                        string     `json:"objective"`
	TargetAudience                   string     `json:"target_audience"`
	StartDate                        time.Time  `json:"start_date"`
	EndDate                          *time.Time `json:"end_date"`
	LinkedSecurityControl            *string    `json:"linked_security_control"`
	AwarenessCampaignTypeID          uint       `json:"awareness_campaign_type_id"`
	AwarenessCampaignTypeTitleEn     string     `json:"awareness_campaign_type_title_en"`
	AwarenessCampaignTypeTitleAr     string     `json:"awareness_campaign_type_title_ar"`
	AwarenessPriorityID              *uint      `json:"awareness_priority_id"`
	AwarenessPriorityTitleEn         string     `json:"awareness_priority_title_en"`
	AwarenessPriorityTitleAr         string     `json:"awareness_priority_title_ar"`
	AwarenessCampaignStatusID        uint       `json:"awareness_campaign_status_id"`
	AwarenessCampaignStatusTitleEn   string     `json:"awareness_campaign_status_title_en"`
	AwarenessCampaignStatusTitleAr   string     `json:"awareness_campaign_status_title_ar"`
	AwarenessCampaignStatusColorCode string     `json:"awareness_campaign_status_color_code"`
	OrganizationUnitID               uint       `json:"organization_unit_id"`
	OrganizationUnitTitleEn          string     `json:"organization_unit_title_en"`
	OrganizationUnitTitleAr          string     `json:"organization_unit_title_ar"`
	OwnerID                          uint       `json:"owner_id"`
	OwnerFirstNameEn                 string     `json:"owner_first_name_en"`
	OwnerFirstNameAr                 string     `json:"owner_first_name_ar"`
	OwnerLastNameEn                  string     `json:"owner_last_name_en"`
	OwnerLastNameAr                  string     `json:"owner_last_name_ar"`
	OwnerInitial                     string     `json:"owner_initial"`
	OwnerColor                       string     `json:"owner_color"`
	OwnerEmail                       string     `json:"owner_email"`
	OwnerProfileImagePath            string     `json:"owner_profile_image_path"`
	CreatedAt                        string     `json:"created_at"`
	UpdatedAt                        string     `json:"updated_at"`
	CreatedUserID                    uint       `json:"created_user_id"`
	CreatedUserFirstNameEn           string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr           string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn            string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr            string     `json:"created_user_last_name_ar"`
	CreatedUserInitial               string     `json:"created_user_initial"`
	CreatedUserColor                 string     `json:"created_user_color"`
	CreatedUserEmail                 string     `json:"created_user_email"`
	CreatedUserProfileImagePath      string     `json:"created_user_profile_image_path"`
	UpdatedUserID                    *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn           string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr           string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn            string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr            string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial               string     `json:"updated_user_initial"`
	UpdatedUserColor                 string     `json:"updated_user_color"`
	UpdatedUserEmail                 string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath      string     `json:"updated_user_profile_image_path"`
}

type AwarenessCampaignListQueryParamsDTO struct {
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

type AwarenessCampaignListResponseDTO struct {
	Items      []*AwarenessCampaignListDTO  `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AwarenessCampaignFormDTO struct {
	AwarenessCampaignTypeID   uint                            `json:"awareness_campaign_type_id" binding:"required"`
	CampaignTitle             string                          `json:"campaign_title" binding:"required,max=255"`
	Objective                 string                          `json:"objective" binding:"required"`
	AwarenessPriorityID       *uint                           `json:"awareness_priority_id"`
	TargetAudience            string                          `json:"target_audience" binding:"required"`
	OwnerID                   uint                            `json:"owner_id" binding:"required"`
	OrganizationUnitID        uint                            `json:"organization_unit_id" binding:"required"`
	StartDate                 string                          `json:"start_date" binding:"required"`
	EndDate                   *string                         `json:"end_date"`
	AwarenessCampaignStatusID uint                            `json:"awareness_campaign_status_id" binding:"required"`
	LinkedSecurityControl     *string                         `json:"linked_security_control"`
	Detail                    *AwarenessCampaignDetailFormDTO `json:"detail"`
}

type AwarenessCampaignDetailFormDTO struct {
	AwarenessChannelIDs         []uint                              `json:"awareness_channel_ids"`
	AwarenessPackName           *string                             `json:"awareness_pack_name"`
	RecipientCount              *int                                `json:"recipient_count"`
	SentCount                   *int                                `json:"sent_count"`
	EmailOpenRate               *int                                `json:"email_open_rate"`
	ClickRate                   *int                                `json:"click_rate"`
	DataSubmissionRate          *int                                `json:"data_submission_rate"`
	ReportRate                  *int                                `json:"report_rate"`
	RiskClassificationID        *uint                               `json:"risk_classification_id"`
	CampaignResult              *string                             `json:"campaign_result"`
	ImprovementRecommendations  *string                             `json:"improvement_recommendations"`
	SubtypeRole                 *string                             `json:"subtype_role"`
	Presenter                   *string                             `json:"presenter"`
	Agenda                      *string                             `json:"agenda"`
	SessionStartTime            *string                             `json:"session_start_time"`
	SessionEndTime              *string                             `json:"session_end_time"`
	Venue                       *string                             `json:"venue"`
	AttendeeCount               *int                                `json:"attendee_count"`
	AwarenessEvaluationRatingID *uint                               `json:"awareness_evaluation_rating_id"`
	Documents                   []*AwarenessCampaignDocumentFormDTO `json:"documents"`
}

type AwarenessCampaignDocumentFormDTO struct {
	ID            *uint  `json:"id,omitempty"`
	URL           string `json:"url"`
	FileName      string `json:"file_name"`
	OriginalName  string `json:"original_name"`
	Type          string `json:"type"`
	Size          int    `json:"size"`
	IsNewFile     *bool  `json:"is_new_file,omitempty"`
	IsDeletedFile *bool  `json:"is_deleted_file,omitempty"`
}

type AwarenessCampaignDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type AwarenessEvaluationRatingDTO struct {
	ID      uint   `json:"id"`
	TitleEn string `json:"title_en"`
	TitleAr string `json:"title_ar"`
}

type AwarenessCampaignDetailDTO struct {
	ID                         uint                            `json:"id"`
	AwarenessCampaignID        uint                            `json:"awareness_campaign_id"`
	AwarenessChannels          []*AwarenessChannelDTO          `json:"awareness_channels"`
	AwarenessPackName          *string                         `json:"awareness_pack_name"`
	RecipientCount             *int                            `json:"recipient_count"`
	SentCount                  *int                            `json:"sent_count"`
	EmailOpenRate              *int                            `json:"email_open_rate"`
	ClickRate                  *int                            `json:"click_rate"`
	DataSubmissionRate         *int                            `json:"data_submission_rate"`
	ReportRate                 *int                            `json:"report_rate"`
	RiskClassification         *AwarenessPriorityDTO           `json:"risk_classification"`
	CampaignResult             *string                         `json:"campaign_result"`
	ImprovementRecommendations *string                         `json:"improvement_recommendations"`
	SubtypeRole                *string                         `json:"subtype_role"`
	Presenter                  *string                         `json:"presenter"`
	Agenda                     *string                         `json:"agenda"`
	SessionStartTime           *string                         `json:"session_start_time"`
	SessionEndTime             *string                         `json:"session_end_time"`
	Venue                      *string                         `json:"venue"`
	AttendeeCount              *int                            `json:"attendee_count"`
	AwarenessEvaluationRating  *AwarenessEvaluationRatingDTO   `json:"awareness_evaluation_rating"`
	Documents                  []*AwarenessCampaignDocumentDTO `json:"documents"`
	CreatedAt                  time.Time                       `json:"created_at"`
	UpdatedAt                  *time.Time                      `json:"updated_at"`
}

type AwarenessCampaignDetailResponseDTO struct {
	ID                      uint                             `json:"id"`
	ReferenceNumber         string                           `json:"reference_number"`
	CampaignTitle           string                           `json:"campaign_title"`
	Objective               string                           `json:"objective"`
	TargetAudience          string                           `json:"target_audience"`
	StartDateString         *string                          `json:"start_date"`
	EndDateString           *string                          `json:"end_date"`
	LinkedSecurityControl   *string                          `json:"linked_security_control"`
	AwarenessCampaignType   AwarenessCampaignTypeDTO         `json:"awareness_campaign_type"`
	AwarenessPriority       *AwarenessPriorityDTO            `json:"awareness_priority"`
	AwarenessCampaignStatus AwarenessCampaignStatusDTO       `json:"awareness_campaign_status"`
	Owner                   organizationDTOs.UserMinimalDTO  `json:"owner"`
	OrganizationUnit        settingsDTOs.OrganizationUnitDTO `json:"organization_unit"`
	Detail                  *AwarenessCampaignDetailDTO      `json:"detail"`
	CreatedAt               time.Time                        `json:"created_at"`
	UpdatedAt               *time.Time                       `json:"updated_at"`
	CreatedByUser           organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser           organizationDTOs.UserMinimalDTO  `json:"updated_by"`
}

type AwarenessCampaignParsedDates struct {
	StartDate *time.Time
	EndDate   *time.Time
}
