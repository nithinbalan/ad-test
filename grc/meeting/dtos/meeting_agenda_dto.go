package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type MeetingAgendaListDTO struct {
	ID                     uint       `json:"id"`
	MeetingID              uint       `json:"meeting_id"`
	TimeFrom               *string    `json:"time_from"`
	TimeTo                 *string    `json:"time_to"`
	MeetingTitle           string     `json:"meeting_title"`
	MeetingLocation        *string    `json:"meeting_location"`
	MeetingDate            time.Time  `json:"meeting_date"`
	MeetingStartTime       time.Time  `json:"meeting_start_time"`
	MeetingEndTime         time.Time  `json:"meeting_end_time"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at"`
	CreatedUserFirstNameEn string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn  string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr  string     `json:"created_user_last_name_ar"`
	UpdatedUserFirstNameEn string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn  string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr  string     `json:"updated_user_last_name_ar"`
}

type MeetingAgendaListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type MeetingAgendaListResponseDTO struct {
	Items      []*MeetingAgendaListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingAgendaDetailDTO struct {
	ID            uint                             `json:"id"`
	MeetingID     uint                             `json:"meeting_id"`
	TimeFrom      *string                          `json:"time_from"`
	TimeTo        *string                          `json:"time_to"`
	CreatedByUser organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
}

type MeetingAgendaFormDTO struct {
	TimeFrom string `json:"time_from"`
	TimeTo   string `json:"time_to"`
}
