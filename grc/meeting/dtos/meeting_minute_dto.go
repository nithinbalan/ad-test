package dtos

import (
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type MeetingMinuteListDTO struct {
	ID                           uint       `json:"id"`
	MeetingAgendaID              *uint      `json:"meeting_agenda_id"`
	MeetingID                    *uint      `json:"meeting_id"`
	MeetingTitle                 *string    `json:"meeting_title"`
	Minutes                      *string    `json:"minutes"`
	MeetingAgendaTitle           *string    `json:"meeting_agenda_title"`
	MeetingMinuteStatusID        uint       `json:"meeting_minute_status_id"`
	MeetingMinuteStatusTitleEn   string     `json:"meeting_minute_status_title_en"`
	MeetingMinuteStatusTitleAr   string     `json:"meeting_minute_status_title_ar"`
	MeetingMinuteStatusType      string     `json:"meeting_minute_status_type"`
	MeetingMinuteStatusColorCode string     `json:"meeting_minute_status_color_code"`
	CreatedAt                    time.Time  `json:"created_at"`
	UpdatedAt                    *time.Time `json:"updated_at"`
	CreatedUserID                uint       `json:"created_user_id"`
	CreatedUserFirstNameEn       string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr       string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn        string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr        string     `json:"created_user_last_name_ar"`
	CreatedUserInitial           string     `json:"created_user_initial"`
	CreatedUserColor             string     `json:"created_user_color"`
	CreatedUserEmail             string     `json:"created_user_email"`
	CreatedUserProfileImagePath  string     `json:"created_user_profile_image_path"`
	UpdatedUserID                *uint      `json:"updated_user_id"`
	UpdatedUserFirstNameEn       string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr       string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn        string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr        string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial           string     `json:"updated_user_initial"`
	UpdatedUserColor             string     `json:"updated_user_color"`
	UpdatedUserEmail             string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath  string     `json:"updated_user_profile_image_path"`
}

type MeetingMinuteListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}

type MeetingMinuteDetailDTO struct {
	ID                  uint                             `json:"id"`
	MeetingAgendaID     *uint                            `json:"meeting_agenda_id"`
	MeetingID           *uint                            `json:"meeting_id"`
	Minutes             *string                          `json:"minutes"`
	MeetingMinuteStatus MeetingMinuteStatusDTO           `json:"meeting_minute_status"`
	CreatedByUser       organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser       *organizationDTOs.UserMinimalDTO `json:"updated_by_id"`
	MeetingAgenda       *MeetingAgendaDetailDTO          `json:"agenda"`
	Meeting             *MeetingDetailDTO                `json:"meeting"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           *time.Time                       `json:"updated_at"`
}

type MeetingMinuteFormDTO struct {
	Minutes         string `json:"minutes"  binding:"required"`
	MeetingAgendaID *uint  `json:"meeting_agenda_id"`
}
