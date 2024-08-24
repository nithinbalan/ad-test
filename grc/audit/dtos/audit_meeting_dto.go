package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type AuditMeetingListDTO struct {
	ID                          uint    `json:"id"`
	AuditsID                    uint    `json:"audits_id"`
	Title                       string  `json:"title"`
	Date                        string  `json:"date"`
	FromTime                    string  `json:"from_time"`
	ToTime                      string  `json:"to_time"`
	OtherParticipants           *string `json:"other_participants"`
	Mom                         *string `json:"mom"`
	CreatedAt                   string  `json:"created_at"`
	UpdatedAt                   *string `json:"updated_at"`
	CreatedUserID               uint    `json:"created_user_id"`
	CreatedUserFirstNameEn      string  `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string  `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string  `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string  `json:"created_user_last_name_ar"`
	CreatedUserInitial          string  `json:"created_user_initial"`
	CreatedUserColor            string  `json:"created_user_color"`
	CreatedUserEmail            string  `json:"created_user_email"`
	CreatedUserProfileImagePath string  `json:"created_user_profile_image_path"`
	UpdatedUserID               *uint   `json:"updated_user_id"`
	UpdatedUserFirstNameEn      *string `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      *string `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       *string `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       *string `json:"updated_user_last_name_ar"`
	UpdatedUserInitial          *string `json:"updated_user_initial"`
	UpdatedUserColor            *string `json:"updated_user_color"`
	UpdatedUserEmail            *string `json:"updated_user_email"`
	UpdatedUserProfileImagePath *string `json:"updated_user_profile_image_path"`
}

type AuditMeetingFormDTO struct {
	AuditID           uint    `json:"audit_id" binding:"required"`
	Title             string  `json:"title" binding:"required"`
	Date              string  `json:"date" binding:"required"`
	FromTime          string  `json:"from_time"`
	ToTime            string  `json:"to_time"`
	OtherParticipants *string `json:"other_participants"`
	Mom               *string `json:"mom"`
	ParticipantIDs    []uint  `json:"participant_ids"`
}

type AuditMeetingDetailDTO struct {
	ID                uint                               `json:"id"`
	AuditsID          uint                               `json:"audits_id"`
	Title             string                             `json:"title"`
	Date              string                             `json:"date"`
	FromTime          string                             `json:"from_time"`
	ToTime            string                             `json:"to_time"`
	OtherParticipants *string                            `json:"other_participants"`
	Mom               *string                            `json:"mom"`
	CreatedAt         time.Time                          `json:"created_at"`
	UpdatedAt         *time.Time                         `json:"updated_at"`
	CreatedByUser     organizationDTOs.UserMinimalDTO    `json:"created_by"`
	UpdatedByUser     *organizationDTOs.UserMinimalDTO   `json:"updated_by"`
	Participants      []*organizationDTOs.UserMinimalDTO `json:"participants"`
}

type AuditMeetingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	StartDate   string
	EndDate     string
}

type AuditMeetingListResponseDTO struct {
	Items      []*AuditMeetingListDTO       `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditMeetingParsedDates struct {
	Date     time.Time
	FromTime time.Time
	ToTime   time.Time
}
