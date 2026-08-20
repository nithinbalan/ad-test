package dtos

import (
	"time"

	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	settingsDTOs "jethur/grc/settings/dtos"
)

type CSLogSocRoomAccessLogListDTO struct {
	ID                                 uint      `json:"id"`
	ReferenceNumber                    string    `json:"reference_number"`
	VisitorFullName                    string    `json:"visitor_full_name"`
	EntryTime                          time.Time `json:"entry_time"`
	ExitTime                           time.Time `json:"exit_time"`
	VisitorIDNumber                    string    `json:"visitor_id_number"`
	DepartmentID                       uint      `json:"department_id"`
	DepartmentTitleEn                  string    `json:"department_title_en"`
	DepartmentTitleAr                  string    `json:"department_title_ar"`
	ExternalPartyName                  string    `json:"external_party_name"`
	PurposeOfEntry                     string    `json:"purpose_of_entry"`
	AuthorizingOfficerID               uint      `json:"authorizing_officer_id"`
	AuthorizingOfficerFirstNameEn      string    `json:"authorizing_officer_first_name_en"`
	AuthorizingOfficerFirstNameAr      string    `json:"authorizing_officer_first_name_ar"`
	AuthorizingOfficerLastNameEn       string    `json:"authorizing_officer_last_name_en"`
	AuthorizingOfficerLastNameAr       string    `json:"authorizing_officer_last_name_ar"`
	AuthorizingOfficerInitial          string    `json:"authorizing_officer_initial"`
	AuthorizingOfficerColor            string    `json:"authorizing_officer_color"`
	AuthorizingOfficerProfileImagePath string    `json:"authorizing_officer_profile_image_path"`
	CreatedAt                          string    `json:"created_at"`
	UpdatedAt                          *string   `json:"updated_at"`
	CreatedUserID                      uint      `json:"created_user_id"`
	CreatedUserFirstNameEn             string    `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr             string    `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn              string    `json:"created_user_last_name_en"`
	CreatedUserLastNameAr              string    `json:"created_user_last_name_ar"`
	CreatedUserInitial                 string    `json:"created_user_initial"`
	CreatedUserColor                   string    `json:"created_user_color"`
	CreatedUserEmail                   string    `json:"created_user_email"`
	CreatedUserProfileImagePath        string    `json:"created_user_profile_image_path"`
	UpdatedUserID                      *uint     `json:"updated_user_id"`
	UpdatedUserFirstNameEn             *string   `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr             *string   `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn              *string   `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr              *string   `json:"updated_user_last_name_ar"`
	UpdatedUserInitial                 *string   `json:"updated_user_initial"`
	UpdatedUserColor                   *string   `json:"updated_user_color"`
	UpdatedUserEmail                   *string   `json:"updated_user_email"`
	UpdatedUserProfileImagePath        *string   `json:"updated_user_profile_image_path"`
}

type CSLogSocRoomAccessLogFormDTO struct {
	VisitorFullName      string                        `json:"visitor_full_name" binding:"required,max=255"`
	EntryTime            string                        `json:"entry_time" binding:"required"`
	ExitTime             string                        `json:"exit_time" binding:"required"`
	VisitorIDNumber      string                        `json:"visitor_id_number" binding:"required,max=100"`
	DepartmentID         uint                          `json:"department_id" binding:"required"`
	ExternalPartyName    string                        `json:"external_party_name" binding:"required,max=255"`
	PurposeOfEntry       string                        `json:"purpose_of_entry" binding:"required"`
	AuthorizingOfficerID uint                          `json:"authorizing_officer_id" binding:"required"`
	Notes                *string                       `json:"notes"`
	Documents            []*commonDTOs.DocumentFormDTO `json:"documents"`
}

type CSLogSocRoomAccessLogDocumentDTO struct {
	ID           *uint   `json:"id"`
	URL          *string `json:"url"`
	Name         string  `json:"name"`
	Token        string  `json:"token"`
	OriginalName string  `json:"original_name"`
	Type         string  `json:"type"`
	Size         int     `json:"size"`
}

type CSLogSocRoomAccessLogListQueryParamsDTO struct {
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

type CSLogSocRoomAccessLogListResponseDTO struct {
	Items      []*CSLogSocRoomAccessLogListDTO `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO    `json:"pagination"`
}

type CSLogSocRoomAccessLogDetailDTO struct {
	ID                 uint                                `json:"id"`
	ReferenceNumber    string                              `json:"reference_number"`
	VisitorFullName    string                              `json:"visitor_full_name"`
	EntryTime          time.Time                           `json:"entry_time"`
	ExitTime           time.Time                           `json:"exit_time"`
	VisitorIDNumber    string                              `json:"visitor_id_number"`
	ExternalPartyName  string                              `json:"external_party_name"`
	PurposeOfEntry     string                              `json:"purpose_of_entry"`
	Notes              *string                             `json:"notes"`
	Department         settingsDTOs.OrganizationUnitDTO    `json:"department"`
	AuthorizingOfficer organizationDTOs.UserMinimalDTO     `json:"authorizing_officer"`
	Documents          []*CSLogSocRoomAccessLogDocumentDTO `json:"documents"`
	CreatedAt          time.Time                           `json:"created_at"`
	UpdatedAt          *time.Time                          `json:"updated_at"`
	CreatedByUser      organizationDTOs.UserMinimalDTO     `json:"created_by"`
	UpdatedByUser      *organizationDTOs.UserMinimalDTO    `json:"updated_by"`
}

type CSLogSocRoomAccessLogParsedDates struct {
	EntryTime *time.Time
	ExitTime  *time.Time
}
