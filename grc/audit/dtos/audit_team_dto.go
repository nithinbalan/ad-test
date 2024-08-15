package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"

	"time"
)

type AuditTeamListDTO struct {
	ID                          uint                   `json:"id"`
	TitleEn                     string                 `json:"title_en"`
	TitleAr                     string                 `json:"title_ar"`
	TeamLeadFirstNameEn         string                 `json:"team_lead_first_name_en"`
	TeamLeadFirstNameAr         string                 `json:"team_lead_first_name_ar"`
	TeamLeadLastNameEn          string                 `json:"team_lead_last_name_en"`
	TeamLeadLastNameAr          string                 `json:"team_lead_last_name_ar"`
	TeamLeadEmail               string                 `json:"team_lead_email"`
	TeamLeadInitial             string                 `json:"team_lead_initial"`
	TeamLeadProfileImagePath    *string                `json:"team_lead_profile_image_path"`
	CreatedUserFirstNameEn      string                 `json:"created_user_first_name_en"`
	CreatedUserFirstNameAr      string                 `json:"created_user_first_name_ar"`
	CreatedUserLastNameEn       string                 `json:"created_user_last_name_en"`
	CreatedUserLastNameAr       string                 `json:"created_user_last_name_ar"`
	CreatedUserEmail            string                 `json:"created_user_email"`
	CreatedUserInitial          string                 `json:"created_user_initial"`
	CreatedUserProfileImagePath *string                `json:"created_user_profile_image_path"`
	UpdatedUserFirstNameEn      string                 `json:"updated_user_first_name_en"`
	UpdatedUserFirstNameAr      string                 `json:"updated_user_first_name_ar"`
	UpdatedUserLastNameEn       string                 `json:"updated_user_last_name_en"`
	UpdatedUserLastNameAr       string                 `json:"updated_user_last_name_ar"`
	UpdatedUserEmail            string                 `json:"updated_user_email"`
	UpdatedUserInitial          string                 `json:"updated_user_initial"`
	UpdatedUserProfileImagePath *string                `json:"updated_user_profile_image_path"`
	CreatedAt                   string                 `json:"created_at"`
	UpdatedAt                   string                 `json:"updated_at"`
	Users                       []*AuditTeamMembersDTO `json:"users" gorm:"-"`
}
type AuditTeamListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
}

type AuditTeamListResponseDTO struct {
	Items      []*AuditTeamListDTO          `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type AuditTeamFormDTO struct {
	TitleEn    string  `json:"title_en" binding:"required,max=500"`
	TitleAr    string  `json:"title_ar" binding:"required,max=500"`
	TeamLeadID uint    `json:"team_lead_id" binding:"required"`
	UserIDs    []*uint `json:"user_ids" binding:"required"`
}

type AuditTeamDetailDTO struct {
	ID            uint                             `json:"id"`
	TitleEn       string                           `json:"title_en"`
	TitleAr       string                           `json:"title_ar"`
	TeamLead      *organizationDTOs.UserMinimalDTO `json:"team_lead"`
	Users         []*AuditTeamMembersDTO           `json:"users"`
	CreatedAt     time.Time                        `json:"created_at"`
	UpdatedAt     *time.Time                       `json:"updated_at"`
	CreatedByUser *organizationDTOs.UserMinimalDTO `json:"created_by"`
	UpdatedByUser *organizationDTOs.UserMinimalDTO `json:"updated_by"`
}

type AuditTeamMembersDTO struct {
	ID          uint                            `json:"id"`
	AuditTeamID uint                            `json:"audit_team_id"`
	UserID      uint                            `json:"user_id"`
	User        organizationDTOs.UserMinimalDTO `json:"user"`
}
