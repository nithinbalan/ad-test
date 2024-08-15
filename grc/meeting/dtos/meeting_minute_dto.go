package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type MeetingMinuteListDTO struct {
	ID                          uint       `json:"id"`
	MeetingAgendaID             *uint      `json:"meeting_agenda_id"`
	Minutes                     *string    `json:"minutes"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   *time.Time `json:"updated_at"`
	CreatedUserFirstNameEn      string     `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string     `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string     `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string     `json:"created_user_last_name_ar"`
	CreatedUserInitial          string     `json:"created_user_initial"`
	CreatedUserColor            string     `json:"created_user_color"`
	CreatedUserEmail            string     `json:"created_user_email"`
	CreatedUserProfileImagePath string     `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn      string     `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string     `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string     `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string     `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          string     `json:"updated_user_initial"`
	UpdatedUserColor            string     `json:"updated_user_color"`
	UpdatedUserEmail            string     `json:"updated_user_email"`
	UpdatedUserProfileImagePath string     `json:"updated_user_profile_image_path"`
}

type MeetingMinuteListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type MeetingMinuteListResponseDTO struct {
	Items      []*MeetingMinuteListDTO      `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingMinuteDetailDTO struct {
	ID              uint                             `json:"id"`
	MeetingAgendaID *uint                            `json:"meeting_agenda_id"`
	Minutes         *string                          `json:"minutes"`
	CreatedByUser   organizationDTOs.UserMinimalDTO  `json:"created_by_id"`
	UpdatedByUser   *organizationDTOs.UserMinimalDTO `json:"updated_by_id"`
	MeetingAgenda   *MeetingAgendaDetailDTO          `json:"agenda"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       *time.Time                       `json:"updated_at"`
}

type MeetingMinuteFormDTO struct {
	Minutes         string `json:"minutes"`
	TimeTo          string `json:"time_to"`
	MeetingAgendaID *uint  `json:"meeting_agenda_id"`
}
